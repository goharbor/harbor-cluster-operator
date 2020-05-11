package harbor

import (
	goharborv1 "github.com/goharbor/harbor-cluster-operator/api/v1"
	"github.com/goharbor/harbor-cluster-operator/lcm"
	"github.com/goharbor/harbor-operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (harbor *HarborReconciler) Provision(harborCluster *goharborv1.HarborCluster) (*lcm.CRStatus, error) {
	harborCR := harbor.newHarborCR()
	err := harbor.Client.Create(harbor.Ctx, harborCR)
	if err != nil {
		return harborCRNotReadyStatus("", ""), err
	}
	return harborCRUnknownStatus(), err
}

// newHarborCR will create a new Harbor CR controlled by harbor-operator
func (harbor *HarborReconciler) newHarborCR() *v1alpha1.Harbor {
	namespacedName := harbor.getHarborCRNamespacedName()
	databaseSecret := harbor.getDatabaseSecret()
	cacehSecret := harbor.getCacheSecret()
	storageSecret := harbor.getStorageSecret()

	return &v1alpha1.Harbor{
		ObjectMeta: metav1.ObjectMeta{
			Name:      namespacedName.Name,
			Namespace: namespacedName.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(harbor.HarborCluster, goharborv1.HarborClusterGVK),
			},
		},
		Spec: v1alpha1.HarborSpec{
			HarborVersion: harbor.HarborCluster.Spec.Version,
			PublicURL:     harbor.HarborCluster.Spec.PublicURL,
			TLSSecretName: harbor.HarborCluster.Spec.TLSSecret,
			Components: v1alpha1.HarborComponents{
				Core:        harbor.newCoreComponentIfNecessary(),
				Portal:      harbor.newPortalComponentIfNecessary(),
				Registry:    harbor.newRegistryComponentIfNecessary(storageSecret, cacehSecret),
				JobService:  harbor.newJobServiceComponentIfNecessary(cacehSecret),
				ChartMuseum: harbor.newChartMuseumComponentIfNecessary(storageSecret, cacehSecret),
				Clair:       harbor.newClairComponentIfNecessary(databaseSecret),
				Notary:      harbor.newNotaryComponentIfNecessary(),
			},
			AdminPasswordSecret:  harbor.HarborCluster.Spec.AdminPasswordSecret,
			Priority:             harbor.HarborCluster.Spec.Priority,
			CertificateIssuerRef: harbor.HarborCluster.Spec.CertificateIssuerRef,
		},
	}
}

func (harbor *HarborReconciler) newCoreComponentIfNecessary() *v1alpha1.CoreComponent {
	return nil
}

func (harbor *HarborReconciler) newPortalComponentIfNecessary() *v1alpha1.PortalComponent {
	return nil
}

func (harbor *HarborReconciler) newRegistryComponentIfNecessary(storageSecret, cacheSecret string) *v1alpha1.RegistryComponent {
	return nil
}

func (harbor *HarborReconciler) newJobServiceComponentIfNecessary(cacheSecret string) *v1alpha1.JobServiceComponent {
	return nil
}

func (harbor *HarborReconciler) newChartMuseumComponentIfNecessary(storageSecret, cacheSecret string) *v1alpha1.ChartMuseumComponent {
	return nil
}

func (harbor *HarborReconciler) newClairComponentIfNecessary(databaseSecret string) *v1alpha1.ClairComponent {
	return nil
}

func (harbor *HarborReconciler) newNotaryComponentIfNecessary() *v1alpha1.NotaryComponent {
	return nil
}

// getCacheSecret will get a name of k8s secret which stores cache info
func (harbor *HarborReconciler) getCacheSecret() string {
	// TODO
	return ""
}

// getStorageSecret will get a name of k8s secret which stores storage info
func (harbor *HarborReconciler) getStorageSecret() string {
	// TODO
	return ""
}

// getStorageSecret will get a name of k8s secret which stores storage info
func (harbor *HarborReconciler) getDatabaseSecret() string {
	// TODO
	return ""
}
