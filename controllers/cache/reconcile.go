package cache

import (
	"context"
	"fmt"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	controller "sigs.k8s.io/controller-runtime/pkg/reconcile"
	harborCluster "src/github.com/goharbor/harbor-cluster-operator/api/v1"
	"src/github.com/goharbor/harbor-cluster-operator/controllers/reconciler"
	"time"
)

var (
	defaultRequeue = controller.Result{Requeue: true, RequeueAfter: 10 * time.Second}
)

type Cache interface {
	Reconcile() (*harborCluster.CRStatus, *reconciler.Results)
}

// NewDefaultCache returns the default cache implementation.
func NewDefaultCache(cache ReconcileCache) Cache {
	return &defaultCache{ReconcileCache: cache}
}

type ReconcileCache struct {
	Client   client.Client
	Recorder record.EventRecorder
	Log      logr.Logger
	CTX      context.Context
	Request  controller.Request
	DClient  dynamic.Interface
	Harbor   *harborCluster.HarborCluster
	Scheme   *runtime.Scheme
	ExpectCR *unstructured.Unstructured
	ActualCR *unstructured.Unstructured
	Labels   map[string]string
}

type defaultCache struct {
	ReconcileCache
}

func (d *defaultCache) Reconcile() (*harborCluster.CRStatus, *reconciler.Results) {
	var err error
	results := &reconciler.Results{}
	status := &harborCluster.CRStatus{
		Phase:           harborCluster.PendingPhase,
		ExternalService: "",
		AvailableNodes:  0,
	}
	fmt.Println("Reconcile is Running....")

	d.Labels = d.mergeLabels()

	status, err = d.Deploy(status)
	if err != nil {
		return status, results.WithError(err)
	}

	status, err = d.Readiness(status)
	if err != nil {
		return status, results.WithResult(defaultRequeue)
	}

	status, err = d.UpScale(status)
	if err != nil {
		return status, results.WithError(err)
	}

	status, err = d.DownScale(status)
	if err != nil {
		return status, results.WithError(err)
	}

	return status, results
}
