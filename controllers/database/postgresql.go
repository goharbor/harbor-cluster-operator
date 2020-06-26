package database

import (
	"context"
	"github.com/go-logr/logr"
	goharborv1 "github.com/goharbor/harbor-cluster-operator/api/v1"
	"github.com/goharbor/harbor-cluster-operator/controllers/k8s"
	"github.com/goharbor/harbor-cluster-operator/lcm"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
)

type PostgreSQLReconciler struct {
	HarborCluster *goharborv1.HarborCluster
	Ctx           context.Context
	Client        k8s.Client
	Recorder      record.EventRecorder
	Log           logr.Logger
	DClient       k8s.DClient
	Scheme        *runtime.Scheme
	ExpectCR      *unstructured.Unstructured
	ActualCR      *unstructured.Unstructured
	Labels        map[string]string
}

// Reconciler implements the reconcile logic of postgreSQL service
func (postgres *PostgreSQLReconciler) Reconcile() (*lcm.CRStatus, error) {

	postgres.Client.WithContext(postgres.Ctx)
	postgres.DClient.WithContext(postgres.Ctx)

	crStatus, err := postgres.Provision()
	if err != nil {
		return crStatus, err
	}

	return crStatus, nil
}

func (postgres *PostgreSQLReconciler) Provision() (*lcm.CRStatus, error) {
	if err := postgres.Deploy(); err != nil {
		return databaseNotReadyStatus(CreateDatabaseCrError, err.Error()), err
	}

	crStatus, err := postgres.Readiness()
	if err != nil {
		return databaseNotReadyStatus(CheckDatabaseHealthError, err.Error()), err
	}

	return crStatus, nil
}

func (postgres *PostgreSQLReconciler) Delete() (*lcm.CRStatus, error) {
	panic("implement me")
}

func (postgres *PostgreSQLReconciler) Scale() (*lcm.CRStatus, error) {
	panic("implement me")
}

func (postgres *PostgreSQLReconciler) ScaleUp(newReplicas uint64) (*lcm.CRStatus, error) {
	panic("implement me")
}

func (postgres *PostgreSQLReconciler) ScaleDown(newReplicas uint64) (*lcm.CRStatus, error) {
	panic("implement me")
}

func (postgres *PostgreSQLReconciler) Update(spec *goharborv1.HarborCluster) (*lcm.CRStatus, error) {
	panic("implement me")
}
