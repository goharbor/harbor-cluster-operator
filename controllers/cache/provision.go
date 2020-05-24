package cache

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// Deploy reconcile will deploy Redis sentinel cluster if that does not exist.
// It does:
// - check redis does exist
// - create any new RedisFailovers CRs
// - create redis password secret
// It does not:
// - perform any RedisFailovers downscale (left for downscale phase)
// - perform any RedisFailovers upscale (left for upscale phase)
// - perform any pod upgrade (left for rolling upgrade phase)
func (redis *RedisReconciler) Deploy() error {

	if redis.HarborCluster.Spec.Redis.Kind == "external" {
		return nil
	}

	var actualCR *unstructured.Unstructured
	var expectCR *unstructured.Unstructured

	crdClient := redis.DClient.WithResource(virtualServiceGVR).WithNamespace(redis.Namespace)

	expectCR, err := redis.generateRedisCR()
	if err != nil {
		return err
	}

	if err := controllerutil.SetControllerReference(redis.HarborCluster, expectCR, redis.Scheme); err != nil {
		return err
	}

	actualCR, err = crdClient.Get(redis.Name, metav1.GetOptions{})
	if err != nil && errors.IsNotFound(err) {

		if err := redis.DeploySecret(); err != nil {
			return err
		}

		redis.Log.Info("Creating Redis.", "namespace", redis.Namespace, "name", redis.Name)
		_, err = crdClient.Create(expectCR, metav1.CreateOptions{})
		if err != nil {
			return err
		}

		redis.Log.Info("Redis create complete.", "namespace", redis.Namespace, "name", redis.Name)
	} else if err != nil {
		return err
	} else {
		redis.ExpectCR = expectCR
		redis.ActualCR = actualCR
	}

	return nil
}

// DeploySecret deploy the Redis Password Secret
func (redis *RedisReconciler) DeploySecret() error {
	secret := &corev1.Secret{}
	sc := redis.generateRedisSecret()

	if err := controllerutil.SetControllerReference(redis.HarborCluster, sc, redis.Scheme); err != nil {
		return err
	}
	err := redis.Client.Get(types.NamespacedName{Name: redis.Name, Namespace: redis.Namespace}, secret)
	if err != nil && errors.IsNotFound(err) {
		redis.Log.Info("Creating Redis Password Secret", "namespace", redis.Namespace, "name", redis.Name)
		err = redis.Client.Create(sc)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	return nil
}
