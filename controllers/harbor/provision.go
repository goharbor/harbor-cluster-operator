package harbor

import (
	"fmt"
	goharborv1 "github.com/goharbor/harbor-cluster-operator/api/v1"
	"github.com/goharbor/harbor-cluster-operator/lcm"
	"github.com/goharbor/harbor-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Provision will create a new Harbor CR.
func (harbor *HarborReconciler) Provision() (*lcm.CRStatus, error) {
	harborCR := harbor.newHarborCR()
	err := harbor.Create(harborCR)
	if err != nil {
		return harborCRNotReadyStatus("", ""), err
	}
	return harborCRUnknownStatus(), err
}

// newHarborCR will create a new Harbor CR controlled by harbor-operator
func (harbor *HarborReconciler) newHarborCR() *v1alpha1.Harbor {
	namespacedName := harbor.getHarborCRNamespacedName()

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
				Registry:    harbor.newRegistryComponentIfNecessary(),
				JobService:  harbor.newJobServiceComponentIfNecessary(),
				ChartMuseum: harbor.newChartMuseumComponentIfNecessary(),
				Clair:       harbor.newClairComponentIfNecessary(),
				Notary:      harbor.newNotaryComponentIfNecessary(),
			},
			AdminPasswordSecret:  harbor.HarborCluster.Spec.AdminPasswordSecret,
			Priority:             harbor.HarborCluster.Spec.Priority,
			CertificateIssuerRef: harbor.HarborCluster.Spec.CertificateIssuerRef,
		},
	}
}

func (harbor *HarborReconciler) newCoreComponentIfNecessary() *v1alpha1.CoreComponent {
	return &v1alpha1.CoreComponent{
		HarborDeployment: v1alpha1.HarborDeployment{
			Replicas:         IntToInt32Ptr(harbor.HarborCluster.Spec.Replicas),
			Image:            harbor.getCoreComponentImage(),
			NodeSelector:     nil,
			ImagePullSecrets: harbor.getImagePullSecrets(),
		},
		DatabaseSecret: harbor.getCoreDatabaseSecret(),
	}
}

func (harbor *HarborReconciler) newPortalComponentIfNecessary() *v1alpha1.PortalComponent {
	return &v1alpha1.PortalComponent{
		HarborDeployment: v1alpha1.HarborDeployment{
			Replicas:         IntToInt32Ptr(harbor.HarborCluster.Spec.Replicas),
			Image:            harbor.getCoreComponentImage(),
			NodeSelector:     nil,
			ImagePullSecrets: harbor.getImagePullSecrets(),
		},
	}
}

func (harbor *HarborReconciler) newRegistryComponentIfNecessary() *v1alpha1.RegistryComponent {
	return &v1alpha1.RegistryComponent{
		HarborDeployment: v1alpha1.HarborDeployment{
			Replicas:         IntToInt32Ptr(harbor.HarborCluster.Spec.Replicas),
			Image:            harbor.getRegistryComponentImage(),
			NodeSelector:     nil,
			ImagePullSecrets: harbor.getImagePullSecrets(),
		},
		Controller: v1alpha1.RegistryControllerComponent{
			Image: harbor.getRegistryControllerComponentImage(),
		},
		StorageSecret: harbor.getRegistryStorageSecret(),
		CacheSecret:   harbor.getRegistryCacheSecret(),
	}
}

func (harbor *HarborReconciler) newJobServiceComponentIfNecessary() *v1alpha1.JobServiceComponent {
	if harbor.HarborCluster.Spec.JobService != nil {
		return &v1alpha1.JobServiceComponent{
			HarborDeployment: v1alpha1.HarborDeployment{
				Replicas:         IntToInt32Ptr(harbor.HarborCluster.Spec.Replicas),
				Image:            harbor.getJobServiceComponentImage(),
				NodeSelector:     nil,
				ImagePullSecrets: harbor.getImagePullSecrets(),
			},
			RedisSecret: harbor.getJobServiceCacheSecret(),
			WorkerCount: harbor.HarborCluster.Spec.JobService.WorkerCount,
		}
	}
	return nil
}

func (harbor *HarborReconciler) newChartMuseumComponentIfNecessary() *v1alpha1.ChartMuseumComponent {
	if harbor.HarborCluster.Spec.ChartMuseum != nil {
		return &v1alpha1.ChartMuseumComponent{
			HarborDeployment: v1alpha1.HarborDeployment{
				Replicas:         IntToInt32Ptr(harbor.HarborCluster.Spec.Replicas),
				Image:            harbor.getChartMuseumComponentImage(),
				NodeSelector:     nil,
				ImagePullSecrets: harbor.getImagePullSecrets(),
			},
			StorageSecret: harbor.getStorageSecret(),
			CacheSecret:   harbor.getChartMuseumCacheSecret(),
		}
	}
	return nil
}

func (harbor *HarborReconciler) newClairComponentIfNecessary() *v1alpha1.ClairComponent {
	if harbor.HarborCluster.Spec.Clair != nil {
		return &v1alpha1.ClairComponent{
			HarborDeployment: v1alpha1.HarborDeployment{
				Replicas:         IntToInt32Ptr(harbor.HarborCluster.Spec.Replicas),
				Image:            harbor.getClairComponentImage(),
				NodeSelector:     nil,
				ImagePullSecrets: harbor.getImagePullSecrets(),
			},
			DatabaseSecret:       harbor.getClairDatabaseSecret(),
			VulnerabilitySources: harbor.HarborCluster.Spec.Clair.VulnerabilitySources,
			Adapter: v1alpha1.ClairAdapterComponent{
				Image:       harbor.getClairAdapterComponentImage(),
				RedisSecret: harbor.getClairCacheSecret(),
			},
		}
	}
	return nil
}

// newNotaryComponentIfNecessary will return a NotaryComponent in harbor CRD.
// TODO our HarborCluster CRD not define Notary spec.
func (harbor *HarborReconciler) newNotaryComponentIfNecessary() *v1alpha1.NotaryComponent {
	return nil
}

// getRegistryCacheSecret will get a name of k8s secret which stores cache info
func (harbor *HarborReconciler) getRegistryCacheSecret() string {
	p := harbor.getStorageProperty("registrySecret")
	if p == nil {
		return ""
	}
	return p.ToString()
}

// getJobServiceCacheSecret will get a name of k8s secret which stores cache info
func (harbor *HarborReconciler) getJobServiceCacheSecret() string {
	p := harbor.getStorageProperty("jobServiceSecret")
	if p == nil {
		return ""
	}
	return p.ToString()
}

// getClairCacheSecret will get a name of k8s secret which stores cache info
func (harbor *HarborReconciler) getClairCacheSecret() string {
	p := harbor.getStorageProperty("clairSecret")
	if p == nil {
		return ""
	}
	return p.ToString()
}

// getChartMuseumCacheSecret will get a name of k8s secret which stores cache info
func (harbor *HarborReconciler) getChartMuseumCacheSecret() string {
	p := harbor.getStorageProperty("chartMuseumSecret")
	if p == nil {
		return ""
	}
	return p.ToString()
}

// getStorageSecret will get a name of k8s secret which stores storage info
func (harbor *HarborReconciler) getStorageSecret() string {
	switch harbor.HarborCluster.Spec.Storage.Kind {
	case "azure":
		return "azureSecret"
	case "gcs":
		return "gcsSecret"
	case "swift":
		return "swiftSecret"
	case "s3":
		return "s3Secret"
	case "oss":
		return "ossSecret"
	}
	return ""
}

// getStorageSecret will get a name of k8s secret which stores storage info
func (harbor *HarborReconciler) getRegistryStorageSecret() string {
	p := harbor.getStorageProperty("registrySecret")
	if p == nil {
		return ""
	}
	return p.ToString()
}

// getCoreDatabaseSecret will get a name of k8s secret which stores core-database info
func (harbor *HarborReconciler) getCoreDatabaseSecret() string {
	p := harbor.getDatabaseProperty("coreSecret")
	if p == nil {
		return ""
	}
	return ""
}

// getClairDatabaseSecret will get a name of k8s secret which stores clair-database info
func (harbor *HarborReconciler) getClairDatabaseSecret() string {
	p := harbor.getDatabaseProperty("clairSecret")
	if p == nil {
		return ""
	}
	return ""
}

// getNotaryServerDatabaseSecret will get a name of k8s secret which stores notary-server-database info
func (harbor *HarborReconciler) getNotaryServerDatabaseSecret() string {
	p := harbor.getDatabaseProperty("notaryServerSecret")
	if p == nil {
		return ""
	}
	return p.ToString()
}

// getClairDatabaseSecret will get a name of k8s secret which stores clair-database info
func (harbor *HarborReconciler) getNotarySignerDatabaseSecret() string {
	p := harbor.getDatabaseProperty("notarySignerSecret")
	if p == nil {
		return ""
	}
	return p.ToString()
}

func (harbor *HarborReconciler) getStorageProperty(name string) *lcm.Property {
	return harbor.getProperty(goharborv1.ComponentStorage, name)
}

func (harbor *HarborReconciler) getDatabaseProperty(name string) *lcm.Property {
	return harbor.getProperty(goharborv1.ComponentDatabase, name)
}

func (harbor *HarborReconciler) getCacheProperty(name string) *lcm.Property {
	return harbor.getProperty(goharborv1.ComponentCache, name)
}

func (harbor *HarborReconciler) getProperty(component goharborv1.Component, name string) *lcm.Property {
	if harbor.ComponentToCRStatus == nil {
		return nil
	}
	crStatus := harbor.ComponentToCRStatus[component]
	if crStatus == nil || len(crStatus.Properties) == 0 {
		return nil
	}
	return crStatus.Properties.Get(name)
}

func (harbor *HarborReconciler) getCoreComponentImage() *string {
	return harbor.getComponentImage("goharbor/harbor-core")
}

func (harbor *HarborReconciler) getPortalComponentImage() *string {
	return harbor.getComponentImage("goharbor/harbor-portal")
}

func (harbor *HarborReconciler) getRegistryComponentImage() *string {
	return harbor.getComponentImage("goharbor/registry-photon")
}

func (harbor *HarborReconciler) getRegistryControllerComponentImage() *string {
	return harbor.getComponentImage("goharbor/harbor-registry")
}

func (harbor *HarborReconciler) getJobServiceComponentImage() *string {
	return harbor.getComponentImage("goharbor/harbor-jobservice")
}

func (harbor *HarborReconciler) getChartMuseumComponentImage() *string {
	return harbor.getComponentImage("goharbor/chartmuseum-photon")
}

func (harbor *HarborReconciler) getClairComponentImage() *string {
	return harbor.getComponentImage("holyhope/clair-adapter-with-config")
}

func (harbor *HarborReconciler) getClairAdapterComponentImage() *string {
	return harbor.getComponentImage("goharbor/clair-photon")
}

func (harbor *HarborReconciler) getComponentImage(imageRepo string) *string {
	var image string
	if harbor.HarborCluster.Spec.ImageSource == nil && harbor.HarborCluster.Spec.ImageSource.Registry == "" {
		image = fmt.Sprintf("%s:%s", imageRepo, harbor.HarborCluster.Spec.Version)
	} else {
		image = fmt.Sprintf("%s/%s:%s", harbor.HarborCluster.Spec.ImageSource.Registry, imageRepo, harbor.HarborCluster.Spec.Version)
	}
	return &image
}

func (harbor *HarborReconciler) getImagePullSecrets() []corev1.LocalObjectReference {
	if harbor.HarborCluster.Spec.ImageSource != nil {
		return []corev1.LocalObjectReference{{
			harbor.HarborCluster.Spec.ImageSource.ImagePullSecret,
		}}
	}
	return nil
}

func IntToInt32Ptr(value int) *int32 {
	int32Val := int32(value)
	return &int32Val
}
