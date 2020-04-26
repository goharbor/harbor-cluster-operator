package controllers

import (
	goharborv1 "github.com/goharbor/harbor-cluster-operator/api/v1"
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
	Harbor(harborCluster *goharborv1.HarborCluster) Reconciler
}
