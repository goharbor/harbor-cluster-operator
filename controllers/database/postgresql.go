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

func (postgresql *PostgreSQLReconciler) Provision() (*lcm.CRStatus, error) {
	panic("implement me")
}

func (postgresql *PostgreSQLReconciler) Delete() (*lcm.CRStatus, error) {
	panic("implement me")
}

func (postgresql *PostgreSQLReconciler) Scale() (*lcm.CRStatus, error) {
	panic("implement me")
}

func (postgresql *PostgreSQLReconciler) ScaleUp(newReplicas uint64) (*lcm.CRStatus, error) {
	panic("implement me")
}

func (postgresql *PostgreSQLReconciler) ScaleDown(newReplicas uint64) (*lcm.CRStatus, error) {
	panic("implement me")
}

func (postgresql *PostgreSQLReconciler) Update(spec *goharborv1.HarborCluster) (*lcm.CRStatus, error) {
	panic("implement me")
}
