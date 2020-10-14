package controllers

import (
	"context"
	"github.com/go-logr/logr"
	goharborv1 "github.com/goharbor/harbor-cluster-operator/apis/goharbor.io/v1alpha1"
	"github.com/goharbor/harbor-cluster-operator/controllers/cache"
	"github.com/goharbor/harbor-cluster-operator/controllers/database"
	"github.com/goharbor/harbor-cluster-operator/controllers/harbor"
	"github.com/goharbor/harbor-cluster-operator/controllers/image"
	"github.com/goharbor/harbor-cluster-operator/controllers/k8s"
	"github.com/goharbor/harbor-cluster-operator/controllers/storage"
	"github.com/goharbor/harbor-cluster-operator/lcm"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
)

type Reconciler interface {
	// Reconcile the dependent service.
	Reconcile() (*lcm.CRStatus, error)
}

type ServiceGetter interface {
	// For Redis
	Cache(ctx context.Context, harborCluster *goharborv1.HarborCluster, options *GetOptions) Reconciler

	// For database
	Database(ctx context.Context, harborCluster *goharborv1.HarborCluster, options *GetOptions) Reconciler

	// For storage
	Storage(ctx context.Context, harborCluster *goharborv1.HarborCluster, options *GetOptions) Reconciler

	// For harbor itself
	Harbor(ctx context.Context, harborCluster *goharborv1.HarborCluster, componentToCRStatus map[goharborv1.Component]*lcm.CRStatus, options *GetOptions) Reconciler
}

type GetOptions struct {
	Client      k8s.Client
	Recorder    record.EventRecorder
	Log         logr.Logger
	DClient     k8s.DClient
	Scheme      *runtime.Scheme
	ImageGetter image.Getter
}

type ServiceGetterImpl struct {
}

func (impl *ServiceGetterImpl) Cache(ctx context.Context, harborCluster *goharborv1.HarborCluster, options *GetOptions) Reconciler {
	return &cache.RedisReconciler{
		HarborCluster: harborCluster,
		Client:        options.Client,
		Recorder:      options.Recorder,
		Log:           options.Log,
		DClient:       options.DClient,
		Scheme:        options.Scheme,
		CXT:           ctx,
	}
}

func (impl *ServiceGetterImpl) Database(ctx context.Context, harborCluster *goharborv1.HarborCluster, options *GetOptions) Reconciler {
	return &database.PostgreSQLReconciler{
		HarborCluster: harborCluster,
		Client:        options.Client,
		Recorder:      options.Recorder,
		Log:           options.Log,
		DClient:       options.DClient,
		Scheme:        options.Scheme,
		Ctx:           ctx,
	}
}

func (impl *ServiceGetterImpl) Storage(ctx context.Context, harborCluster *goharborv1.HarborCluster, options *GetOptions) Reconciler {
	return &storage.MinIOReconciler{
		HarborCluster: harborCluster,
		KubeClient:    options.Client,
		Ctx:           ctx,
		Log:           options.Log,
		Recorder:      options.Recorder,
	}
}

func (impl *ServiceGetterImpl) Harbor(ctx context.Context, harborCluster *goharborv1.HarborCluster, componentToCRStatus map[goharborv1.Component]*lcm.CRStatus, options *GetOptions) Reconciler {
	return &harbor.HarborReconciler{
		HarborCluster:       harborCluster,
		Client:              options.Client,
		ImageGetter:         options.ImageGetter,
		Ctx:                 ctx,
		ComponentToCRStatus: componentToCRStatus,
	}
}
