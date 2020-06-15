package database

import (
	"fmt"
	"k8s.io/apimachinery/pkg/api/errors"
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
func (postgres *PostgreSQLReconciler) Deploy() error {

	if postgres.HarborCluster.Spec.Database.Kind == "external" {
		return nil
	}

	var (
		actualCR *unstructured.Unstructured
		expectCR *unstructured.Unstructured
	)

	name := fmt.Sprintf("%s-%s", postgres.HarborCluster.Namespace, postgres.HarborCluster.Name)

	crdClient := postgres.DClient.WithResource(databaseFailoversGVR).WithNamespace(postgres.HarborCluster.Namespace)

	expectCR, err := postgres.generatePostgresCR()
	if err != nil {
		return err
	}

	if err := controllerutil.SetControllerReference(postgres.HarborCluster, expectCR, postgres.Scheme); err != nil {
		return err
	}

	actualCR, err = crdClient.Get(name, metav1.GetOptions{})
	if err != nil && errors.IsNotFound(err) {

		postgres.Log.Info("Creating Database.", "namespace", postgres.HarborCluster.Namespace, "name", name)
		_, err = crdClient.Create(expectCR, metav1.CreateOptions{})
		if err != nil {
			return err
		}

		postgres.Log.Info("Database create complete.", "namespace", postgres.HarborCluster.Namespace, "name", name)
		return nil
	}

	if err != nil {
		return err
	}

	postgres.Log.Info("Database has existed.", "namespace", postgres.HarborCluster.Namespace, "name", name)

	postgres.ExpectCR = expectCR
	postgres.ActualCR = actualCR
	return nil
}
