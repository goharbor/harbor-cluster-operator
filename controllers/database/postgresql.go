package database

import (
	goharborv1 "github.com/goharbor/harbor-cluster-operator/api/v1"
	"github.com/goharbor/harbor-cluster-operator/lcm"
)

type PostgreSQLReconciler struct {
	HarborCluster *goharborv1.HarborCluster
}

// Reconciler implements the reconcile logic of postgreSQL service
func (postgresql *PostgreSQLReconciler) Reconcile() (*lcm.CRStatus, error) {
	// TODO
	return nil, nil
}

func (postgresql *PostgreSQLReconciler) provision(spec *goharborv1.HarborCluster) (*lcm.CRStatus, error) {
	panic("implement me")
}

func (postgresql *PostgreSQLReconciler) delete() (*lcm.CRStatus, error) {
	panic("implement me")
}

func (postgresql *PostgreSQLReconciler) scaleUp(newReplicas uint64) (*lcm.CRStatus, error) {
	panic("implement me")
}

func (postgresql *PostgreSQLReconciler) scaleDown(newReplicas uint64) (*lcm.CRStatus, error) {
	panic("implement me")
}

func (postgresql *PostgreSQLReconciler) update(spec *goharborv1.HarborCluster) (*lcm.CRStatus, error) {
	panic("implement me")
}
