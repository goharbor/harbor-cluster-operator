package harbor

import (
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
			Image:            harbor.ImageGetter.CoreImage(),
			NodeSelector:     nil,
			ImagePullSecrets: harbor.getImagePullSecrets(),
		},
	}
}

func (harbor *HarborReconciler) newPortalComponentIfNecessary() *v1alpha1.PortalComponent {
	return &v1alpha1.PortalComponent{
		HarborDeployment: v1alpha1.HarborDeployment{
			Replicas:         IntToInt32Ptr(harbor.HarborCluster.Spec.Replicas),
			Image:            harbor.ImageGetter.PortalImage(),
			NodeSelector:     nil,
			ImagePullSecrets: harbor.getImagePullSecrets(),
		},
	}
}

func (harbor *HarborReconciler) newRegistryComponentIfNecessary() *v1alpha1.RegistryComponent {
	return &v1alpha1.RegistryComponent{
		HarborDeployment: v1alpha1.HarborDeployment{
			Replicas:         IntToInt32Ptr(harbor.HarborCluster.Spec.Replicas),
			Image:            harbor.ImageGetter.RegistryImage(),
			NodeSelector:     nil,
			ImagePullSecrets: harbor.getImagePullSecrets(),
		},
		Controller: v1alpha1.RegistryControllerComponent{
			Image: harbor.ImageGetter.RegistryControllerImage(),
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
				Image:            harbor.ImageGetter.JobServiceImage(),
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
				Image:            harbor.ImageGetter.ChartMuseumImage(),
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
				Image:            harbor.ImageGetter.ClairImage(),
				NodeSelector:     nil,
				ImagePullSecrets: harbor.getImagePullSecrets(),
			},
			DatabaseSecret:       harbor.getClairDatabaseSecret(),
			VulnerabilitySources: harbor.HarborCluster.Spec.Clair.VulnerabilitySources,
			Adapter: v1alpha1.ClairAdapterComponent{
				Image:       harbor.ImageGetter.ClairAdapterImage(),
				RedisSecret: harbor.getClairCacheSecret(),
			},
		}
	}
	return nil
}

// newNotaryComponentIfNecessary will return a NotaryComponent in harbor CRD.
func (harbor *HarborReconciler) newNotaryComponentIfNecessary() *v1alpha1.NotaryComponent {
	if harbor.HarborCluster.Spec.Notary != nil {
		return &v1alpha1.NotaryComponent{
			PublicURL: harbor.HarborCluster.Spec.Notary.PublicURL,
			DBMigrator: v1alpha1.NotaryDBMigrator{
				Image: harbor.ImageGetter.NotaryDBMigratorImage(),
			},
			Signer: v1alpha1.NotarySignerComponent{
				HarborDeployment: v1alpha1.HarborDeployment{
					Replicas:         IntToInt32Ptr(harbor.HarborCluster.Spec.Replicas),
					Image:            harbor.ImageGetter.NotarySingerImage(),
					NodeSelector:     nil,
					ImagePullSecrets: harbor.getImagePullSecrets(),
				},
				DatabaseSecret: harbor.getNotarySignerDatabaseSecret(),
			},
			Server: v1alpha1.NotaryServerComponent{
				HarborDeployment: v1alpha1.HarborDeployment{
					Replicas:         IntToInt32Ptr(harbor.HarborCluster.Spec.Replicas),
					Image:            harbor.ImageGetter.NotaryServerImage(),
					NodeSelector:     nil,
					ImagePullSecrets: harbor.getImagePullSecrets(),
				},
				DatabaseSecret: harbor.getNotaryServerDatabaseSecret(),
			},
		}
	}
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
