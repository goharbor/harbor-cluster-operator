package database

import (
	"fmt"
	"github.com/goharbor/harbor-cluster-operator/controllers/k8s"
	"github.com/goharbor/harbor-cluster-operator/lcm"
	"github.com/google/go-cmp/cmp"
	pg "github.com/zalando/postgres-operator/pkg/apis/acid.zalan.do/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

// RollingUpgrades reconcile will rolling upgrades PostgreSQL if resource upscale.
// It does:
// - check resource
// - update RedisFailovers CR resource
func (postgres *PostgreSQLReconciler) RollingUpgrades() (*lcm.CRStatus, error) {

	name := fmt.Sprintf("%s-%s", postgres.HarborCluster.Namespace, postgres.HarborCluster.Name)

	crdClient := postgres.DClient.WithResource(databaseFailoversGVR).WithNamespace(postgres.HarborCluster.Namespace)
	if postgres.ExpectCR == nil {
		return databaseUnknownStatus(), nil
	}

	var actualCR pg.Postgresql
	var expectCR pg.Postgresql

	if err := runtime.DefaultUnstructuredConverter.
		FromUnstructured(postgres.ActualCR.UnstructuredContent(), &actualCR); err != nil {
		return databaseNotReadyStatus(DefaultUnstructuredConverterError, err.Error()), err
	}

	if err := runtime.DefaultUnstructuredConverter.
		FromUnstructured(postgres.ExpectCR.UnstructuredContent(), &expectCR); err != nil {
		return databaseNotReadyStatus(DefaultUnstructuredConverterError, err.Error()), err
	}

	if IsEqual(expectCR, actualCR) {
		msg := fmt.Sprintf(MessageDatabaseUpdate, name)
		postgres.Recorder.Event(postgres.HarborCluster, corev1.EventTypeNormal, RollingUpgradesDatabase, msg)

		postgres.Log.Info(
			"Update Redis resource",
			"namespace", postgres.HarborCluster.Namespace, "name", name,
		)

		if err := Update(crdClient, actualCR, expectCR); err != nil {
			return databaseNotReadyStatus(UpdateDatabaseCrError, err.Error()), err
		}
	}
	return databaseUnknownStatus(), nil
}

// isEqual check whether cache cr is equal expect.
func IsEqual(actualCR, expectCR pg.Postgresql) bool {
	return cmp.Equal(expectCR.DeepCopy().Spec, actualCR.DeepCopy().Spec)
}

func Update(crdClient k8s.DClient, actualCR, expectCR pg.Postgresql) error {
	expectCR.ObjectMeta.SetResourceVersion(actualCR.ObjectMeta.GetResourceVersion())

	data, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&expectCR)
	if err != nil {
		return err
	}

	_, err = crdClient.Update(&unstructured.Unstructured{Object: data}, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	return nil
}
