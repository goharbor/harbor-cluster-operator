package cache

import (
	"github.com/go-logr/logr"
	goharborv1 "github.com/goharbor/harbor-cluster-operator/api/v1"
	"github.com/goharbor/harbor-cluster-operator/controllers/k8s"
	"github.com/goharbor/harbor-cluster-operator/lcm"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
)

// RedisReconciler implement the Reconciler interface and lcm.Controller interface.
type RedisReconciler struct {
	HarborCluster *goharborv1.HarborCluster
	Client        k8s.Client
	Recorder      record.EventRecorder
	Log           logr.Logger
	DClient       k8s.DClient
	Scheme        *runtime.Scheme
	ExpectCR      *unstructured.Unstructured
	ActualCR      *unstructured.Unstructured
	Labels        map[string]string
	Name          string
	Namespace     string
}

// Reconciler implements the reconcile logic of redis service
func (redis *RedisReconciler) Reconcile() (*lcm.CRStatus, error) {
	redis.Labels = redis.NewLabels()
	redis.Name = redis.GetHarborClusterName()
	redis.Namespace = redis.GetHarborClusterNamespace()

	crStatus := &lcm.CRStatus{
		Phase:      lcm.PendingPhase,
		Properties: nil,
	}

	crStatus, err := redis.Provision(crStatus)
	if err != nil {
		return crStatus, err
	}

	return nil, nil
}

func (redis *RedisReconciler) Provision(crStatus *lcm.CRStatus) (*lcm.CRStatus, error) {
	crStatus, err := redis.Deploy(crStatus)
	if err != nil {
		return crStatus, err
	}

	crStatus, err = redis.Readiness(crStatus)
	if err != nil {
		return crStatus, err
	}
	return crStatus, nil
}

func (redis *RedisReconciler) Delete() (*lcm.CRStatus, error) {
	panic("implement me")
}

func (redis *RedisReconciler) ScaleUp(newReplicas uint64) (*lcm.CRStatus, error) {
	panic("implement me")
}

func (redis *RedisReconciler) ScaleDown(newReplicas uint64) (*lcm.CRStatus, error) {
	panic("implement me")
}

func (redis *RedisReconciler) Update(spec *goharborv1.HarborCluster) (*lcm.CRStatus, error) {
	panic("implement me")
}

