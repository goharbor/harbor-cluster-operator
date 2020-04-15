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
	Reconcile() *reconciler.Results
	Deploy() error
	Readiness() error
	DownScale() error
	RollingUpgrades()
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

func (d *defaultCache) Reconcile() *reconciler.Results {
	results := &reconciler.Results{}
	fmt.Println("Reconcile is Running....")

	d.Labels = d.mergeLabels()

	//deploy redis
	if err := d.Deploy(); err != nil {
		return results.WithError(err)
	}

	if err := d.Readiness(); err != nil {
		return results.WithResult(defaultRequeue)
	}

	//if err := d.Observer(); err != nil {
	//	return results.WithError(err)
	//}
	//
	//if err := d.Finalizers(); err != nil {
	//	return results.WithError(err)
	//}

	if err := d.UpScale(); err != nil {
		return results.WithError(err)
	}

	if err := d.DownScale(); err != nil {
		return results.WithError(err)
	}

	//if err := d.RollingUpgrades(); err != nil {
	//	return results.WithError(err)
	//}

	return results
}
