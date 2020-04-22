package storage

import (
	goharborv1 "github.com/goharbor/harbor-cluster-operator/api/v1"
	"github.com/goharbor/harbor-cluster-operator/lcm"
)

type MinIOReconciler struct {
	HarborCluster *goharborv1.HarborCluster
}

// Reconciler implements the reconcile logic of minIO service
func (minio *MinIOReconciler) Reconcile() (*lcm.CRStatus, error) {
	// TODO
	return nil, nil
}

func (minio *MinIOReconciler) provision(spec *goharborv1.HarborCluster) (*lcm.CRStatus, error) {
	panic("implement me")
}

func (minio *MinIOReconciler) delete() (*lcm.CRStatus, error) {
	panic("implement me")
}

func (minio *MinIOReconciler) scaleUp(newReplicas uint64) (*lcm.CRStatus, error) {
	panic("implement me")
}

func (minio *MinIOReconciler) scaleDown(newReplicas uint64) (*lcm.CRStatus, error) {
	panic("implement me")
}

func (minio *MinIOReconciler) update(spec *goharborv1.HarborCluster) (*lcm.CRStatus, error) {
	panic("implement me")
}
