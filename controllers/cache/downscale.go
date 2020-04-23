package cache

import (
	"fmt"
	redisCli "github.com/spotahome/redis-operator/api/redisfailover/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	harborCluster "src/github.com/goharbor/harbor-cluster-operator/api/v1"

	//appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

func (d *defaultCache) DownScale(status *harborCluster.CRStatus) (*harborCluster.CRStatus, error) {
	name := d.Request.Name
	nameSpace := d.Request.Namespace
	crdClient := d.DClient.Resource(virtualServiceGVR)
	if d.ExpectCR == nil {
		return status, nil
	}

	var actualCR redisCli.RedisFailover
	var expectCR redisCli.RedisFailover

	if err := runtime.DefaultUnstructuredConverter.
		FromUnstructured(d.ActualCR.UnstructuredContent(), &actualCR); err != nil {
		return status, err
	}

	if err := runtime.DefaultUnstructuredConverter.
		FromUnstructured(d.ExpectCR.UnstructuredContent(), &expectCR); err != nil {
		return status, err
	}

	expectReplica := expectCR.Spec.Redis.Replicas
	actualReplica := actualCR.Spec.Redis.Replicas

	if expectReplica < actualReplica {
		msg := fmt.Sprintf(UpdateMessageRedisCluster, name)
		d.Recorder.Event(d.Harbor, corev1.EventTypeNormal, RedisDownScaling, msg)

		d.Log.Info(
			"Scaling replicas down",
			"from", actualReplica,
			"to", expectReplica,
		)

		msg = fmt.Sprintf(MessageRedisDownScaling, actualReplica, expectReplica)
		d.Recorder.Event(d.Harbor, corev1.EventTypeNormal, RedisDownScaling, msg)

		expectCR.ObjectMeta.SetResourceVersion(actualCR.ObjectMeta.GetResourceVersion())

		data, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&expectCR)
		if err != nil {
			return status, err
		}
		d.ExpectCR = &unstructured.Unstructured{Object: data}

		_, err = crdClient.Namespace(nameSpace).Update(d.ExpectCR, metav1.UpdateOptions{})
		if err != nil {
			return status, err
		}
		status.Phase = harborCluster.UpgradingPhase
	}
	return status, nil
}
