package database

import (
	"context"
	"encoding/json"
	"fmt"
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
	CXT           context.Context
	Client        k8s.Client
	Recorder      record.EventRecorder
	Log           logr.Logger
	DClient       k8s.DClient
	Scheme        *runtime.Scheme
	ExpectCR      *unstructured.Unstructured
	ActualCR      *unstructured.Unstructured
	Labels        map[string]string
	Name          string
	Namespace     string
	CRStatus      *lcm.CRStatus
	Connect       *Connect
	Properties    *lcm.Properties
}

// Reconciler implements the reconcile logic of postgreSQL service
func (postgre *PostgreSQLReconciler) Reconcile() (*lcm.CRStatus, error) {

	postgre.Client.WithContext(postgre.CXT)
	postgre.DClient.WithContext(postgre.CXT)

	crStatus, err := postgre.Provision()
	if err != nil {
		return crStatus, err
	}

	return crStatus, nil
}

func (postgre *PostgreSQLReconciler) Provision() (*lcm.CRStatus, error) {
	if err := postgre.Deploy(); err != nil {
		return postgre.CRStatus, err
	}

	if err := postgre.Readiness(); err != nil {
		return postgre.CRStatus, err
	}
	return postgre.CRStatus, nil
}

func (postgre *PostgreSQLReconciler) Delete() (*lcm.CRStatus, error) {
	panic("implement me")
}

func (postgre *PostgreSQLReconciler) ScaleUp(newReplicas uint64) (*lcm.CRStatus, error) {
	panic("implement me")
}

func (postgre *PostgreSQLReconciler) ScaleDown(newReplicas uint64) (*lcm.CRStatus, error) {
	panic("implement me")
}

func (postgre *PostgreSQLReconciler) Update(spec *goharborv1.HarborCluster) (*lcm.CRStatus, error) {
	panic("implement me")
}
