package cache

import (
	"fmt"
	"github.com/goharbor/harbor-cluster-operator/lcm"
	redisCli "github.com/spotahome/redis-operator/api/redisfailover/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"

	//appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

// UpScale reconcile will upscale Redis sentinel cluster if replicas upscale.
// It does:
// - check replicas
// - update RedisFailovers CR replicas
func (redis *RedisReconciler) UpScale(crStatus *lcm.CRStatus) (*lcm.CRStatus, error) {

	crdClient := redis.DClient.WithResource(virtualServiceGVR).WithNamespace(redis.Namespace)
	if redis.ExpectCR == nil {
		return crStatus, nil
	}

	var actualCR redisCli.RedisFailover
	var expectCR redisCli.RedisFailover

	if redis.ExpectCR == nil {
		return crStatus, nil
	}

	if err := runtime.DefaultUnstructuredConverter.
		FromUnstructured(redis.ActualCR.UnstructuredContent(), &actualCR); err != nil {
		return crStatus, err
	}

	if err := runtime.DefaultUnstructuredConverter.
		FromUnstructured(redis.ExpectCR.UnstructuredContent(), &expectCR); err != nil {
		return crStatus, err
	}

	expectReplica := expectCR.Spec.Redis.Replicas
	actualReplica := actualCR.Spec.Redis.Replicas

	if expectReplica > actualReplica {
		msg := fmt.Sprintf(UpdateMessageRedisCluster, redis.Name)
		redis.Recorder.Event(redis.HarborCluster, corev1.EventTypeNormal, RedisUpScaling, msg)

		redis.Log.Info(
			"Scaling replicas up",
			"from", actualReplica,
			"to", expectReplica,
		)

		msg = fmt.Sprintf(MessageRedisUpScaling, actualReplica, expectReplica)
		redis.Recorder.Event(redis.HarborCluster, corev1.EventTypeNormal, RedisUpScaling, msg)

		expectCR.ObjectMeta.SetResourceVersion(actualCR.ObjectMeta.GetResourceVersion())

		data, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&expectCR)
		if err != nil {
			return crStatus, err
		}
		redis.ExpectCR = &unstructured.Unstructured{Object: data}

		_, err = crdClient.Update(redis.ExpectCR, metav1.UpdateOptions{})
		if err != nil {
			return crStatus, err
		}
		crStatus.Phase = lcm.UpgradingPhase
	}
	return crStatus, nil
}

