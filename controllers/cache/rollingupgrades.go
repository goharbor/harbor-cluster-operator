package cache

import (
	"errors"
	"fmt"
	"github.com/goharbor/harbor-cluster-operator/lcm"
	redisCli "github.com/spotahome/redis-operator/api/redisfailover/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"

	//appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

// RollingUpgrades reconcile will rolling upgrades Redis sentinel cluster if resource upscale.
// It does:
// - check resource
// - update RedisFailovers CR resource
func (redis *RedisReconciler) RollingUpgrades(crStatus *lcm.CRStatus) (*lcm.CRStatus, error) {

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
	expectResource := expectCR.Spec.Redis.Resources.String()
	actualResource := actualCR.Spec.Redis.Resources.String()

	_, redisPodList, err := redis.GetStatefulSetPods()
	if err != nil {
		redis.Log.Error(err, "Fail to get deployment pods.")
		return crStatus, err
	}

	if len(redisPodList.Items) < int(actualCR.Spec.Redis.Replicas) {
		redis.Log.Info(
			"Some pods still need to be created/deleted.",
			"namespace", redis.Namespace, "name", redis.Name,
			"expected_pods_num", expectReplica, "actual_pods_num", len(redisPodList.Items),
		)
		return crStatus, errors.New("some pods still create/delete, need to requeue")
	}

	if isUpgradeResource(&expectCR, &actualCR) {
		msg := fmt.Sprintf(UpdateMessageRedisCluster, redis.Name)
		redis.Recorder.Event(redis.HarborCluster, corev1.EventTypeNormal, RedisUpScaling, msg)

		redis.Log.Info(
			"RollingUpgrades Redis resource",
			"namespace", redis.Namespace, "name", redis.Name,
			"expected_resource", expectResource, "actual_resource", actualResource,
		)

		msg = fmt.Sprintf(MessageRedisRollingUpgrades, actualResource, expectResource)
		redis.Recorder.Event(redis.HarborCluster, corev1.EventTypeNormal, RedisRollingUpgrades, msg)

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

