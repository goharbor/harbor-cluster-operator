package database

import (
	"context"
	"fmt"
	"github.com/go-logr/logr"
	goharborv1 "github.com/goharbor/harbor-cluster-operator/api/v1"
	"github.com/goharbor/harbor-cluster-operator/controllers/k8s"
	"github.com/goharbor/harbor-cluster-operator/lcm"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

	crdClient := postgres.DClient.WithResource(databaseFailoversGVR).WithNamespace(postgres.HarborCluster.Namespace)
	name := fmt.Sprintf("%s-%s", postgres.HarborCluster.Namespace, postgres.HarborCluster.Name)

	actualCR, err := crdClient.Get(name, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		return postgres.Provision()
	} else if err != nil {
		return databaseNotReadyStatus(GetDatabaseCrError, err.Error()), err
	}
	expectCR, err := postgres.generatePostgresCR()
	if err != nil {
		return databaseNotReadyStatus(GenerateDatabaseCrError, err.Error()), err
	}

	postgres.ActualCR = actualCR
	postgres.ExpectCR = expectCR

	crStatus, err := postgres.Readiness()
	if err != nil {
		return databaseNotReadyStatus(CheckDatabaseHealthError, err.Error()), err
	}

	crStatus, err = postgres.Update()
	if err != nil {
		return crStatus, err
	}

	return crStatus, nil
}

func (postgres *PostgreSQLReconciler) Provision() (*lcm.CRStatus, error) {
	return postgres.Deploy()
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
