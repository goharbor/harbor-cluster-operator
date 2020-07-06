package database

import (
	"fmt"
	"github.com/goharbor/harbor-cluster-operator/lcm"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// Deploy reconcile will deploy database cluster if that does not exist.
// It does:
// - check postgres.does exist
// - create any new postgresqls.acid.zalan.do CRs
// - create postgres connection secret
// It does not:
// - perform any postgresqls downscale (left for downscale phase)
// - perform any postgresqls upscale (left for upscale phase)
// - perform any pod upgrade (left for rolling upgrade phase)
func (postgres *PostgreSQLReconciler) Deploy() (*lcm.CRStatus, error) {

	if postgres.HarborCluster.Spec.Database.Kind == "external" {
		return databaseUnknownStatus(), nil
	}

	var expectCR *unstructured.Unstructured

	name := fmt.Sprintf("%s-%s", postgres.HarborCluster.Namespace, postgres.HarborCluster.Name)

	crdClient := postgres.DClient.WithResource(databaseFailoversGVR).WithNamespace(postgres.HarborCluster.Namespace)

	expectCR, err := postgres.generatePostgresCR()
	if err != nil {
		return databaseNotReadyStatus(GenerateDatabaseCrError, err.Error()), err
	}

	if err := controllerutil.SetControllerReference(postgres.HarborCluster, expectCR, postgres.Scheme); err != nil {
		return databaseNotReadyStatus(SetOwnerReferenceError, err.Error()), err
	}

	postgres.Log.Info("Creating Database.", "namespace", postgres.HarborCluster.Namespace, "name", name)
	_, err = crdClient.Create(expectCR, metav1.CreateOptions{})
	if err != nil {
		return databaseNotReadyStatus(CreateDatabaseCrError, err.Error()), err
	}

	postgres.Log.Info("Database create complete.", "namespace", postgres.HarborCluster.Namespace, "name", name)
	return databaseUnknownStatus(), nil
}
