package controllers

import (
	goharborv1 "github.com/goharbor/harbor-cluster-operator/api/v1"
	"github.com/goharbor/harbor-cluster-operator/controllers/cache"
	"github.com/goharbor/harbor-cluster-operator/controllers/database"
	"github.com/goharbor/harbor-cluster-operator/controllers/harbor"
	"github.com/goharbor/harbor-cluster-operator/controllers/storage"
	"github.com/goharbor/harbor-cluster-operator/lcm"
)

type Reconciler interface {
	// Reconcile the dependent service.
	Reconcile() (*lcm.CRStatus, error)
}

type ServiceGetter interface {
	// For Redis
	Cache(harborCluster *goharborv1.HarborCluster) Reconciler

	// For database
	Database(harborCluster *goharborv1.HarborCluster) Reconciler

	// For storage
	Storage(harborCluster *goharborv1.HarborCluster) Reconciler

	// For harbor itself
	Harbor(harborCluster *goharborv1.HarborCluster, componentToCRStatus map[goharborv1.Component]*lcm.CRStatus) Reconciler
}

type ServiceGetterImpl struct {
	HarborCluster *goharborv1.HarborCluster
	Client        k8s.Client
	Recorder      record.EventRecorder
	Log           logr.Logger
	DClient       k8s.DClient
	Scheme        *runtime.Scheme
}

func (impl *ServiceGetterImpl) Cache(harborCluster *goharborv1.HarborCluster) Reconciler {
	return &cache.RedisReconciler{
		HarborCluster: harborCluster,
		Client:        impl.Client,
		Recorder:      impl.Recorder,
		Log:           impl.Log,
		DClient:       impl.DClient,
		Scheme:        impl.Scheme,
	}
}

func (impl *ServiceGetterImpl) Database(harborCluster *goharborv1.HarborCluster) Reconciler {
	return &database.PostgreSQLReconciler{
		HarborCluster: harborCluster,
	}
}

func (impl *ServiceGetterImpl) Storage(harborCluster *goharborv1.HarborCluster) Reconciler {
	return &storage.MinIOReconciler{
		HarborCluster: harborCluster,
	}
}

func (impl *ServiceGetterImpl) Harbor(harborCluster *goharborv1.HarborCluster, componentToCRStatus map[goharborv1.Component]*lcm.CRStatus) Reconciler {
	return &harbor.HarborReconciler{
		HarborCluster:       harborCluster,
		ComponentToCRStatus: componentToCRStatus,
	}
}
