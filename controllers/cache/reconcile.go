package cache

import (
	"context"
	"fmt"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	controller "sigs.k8s.io/controller-runtime/pkg/reconcile"
	harborCluster "src/github.com/goharbor/harbor-cluster-operator/api/v1"
	redisCli "src/github.com/goharbor/harbor-cluster-operator/controllers/cache/client/api/v1"
	"src/github.com/goharbor/harbor-cluster-operator/controllers/reconciler"
	"time"
)

var (
	defaultRequeue = controller.Result{Requeue: true, RequeueAfter: 10 * time.Second}
)

type Cache interface {
	Reconcile() *reconciler.Results
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
	ExpectCR *redisCli.RedisFailover
	ActualCR *redisCli.RedisFailover
}

type defaultCache struct {
	ReconcileCache
}

func (d *defaultCache) Reconcile() *reconciler.Results {
	results := &reconciler.Results{}
	fmt.Println("Reconcile is Running....")

	//deploy redis
	if err := d.Deploy(); err != nil {
		return results.WithError(err)
	}
	//

	return results.WithResult(defaultRequeue)
}
