package harbor

import (
	goharborv1 "github.com/goharbor/harbor-cluster-operator/apis/goharbor.io/v1alpha2"
	"github.com/goharbor/harbor-cluster-operator/controllers/image"
	"github.com/goharbor/harbor-cluster-operator/controllers/k8s"
	"github.com/goharbor/harbor-cluster-operator/lcm"
	"github.com/goharbor/harbor-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Provision will create a new Harbor CR.
func (harbor *HarborReconciler) Provision() (*lcm.CRStatus, error) {

	err := harbor.CheckIssuer()
	if err != nil {
		return harborClusterCRNotReadyStatus(CreateRegistryCertError, err.Error()), err
	}

	err = harbor.CheckAdminPasswordSecret()
	if err != nil {
		return harborClusterCRNotReadyStatus(AutoGenerateAdminPasswordError, err.Error()), err
	}

	harborCR := harbor.newHarborCR()
	err = harbor.Create(harborCR)
	if err != nil {
		return harborClusterCRNotReadyStatus(CreateHarborCRError, err.Error()), err
	}
	return harborClusterCRStatus(harborCR), err
}

// newHarborCR will create a new Harbor CR controlled by harbor-operator
func (harbor *HarborReconciler) newHarborCR() *v1alpha1.Harbor {
	namespacedName := harbor.getHarborCRNamespacedName()

	return &v1alpha1.Harbor{
		ObjectMeta: metav1.ObjectMeta{
			Name:      namespacedName.Name,
			Namespace: namespacedName.Namespace,
			Labels: map[string]string{
				k8s.HarborClusterNameLabel: harbor.HarborCluster.Name,
			},
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(harbor.HarborCluster, goharborv1.HarborClusterGVK),
			},
		},
		Spec: v1alpha1.HarborSpec{
			HarborVersion: harbor.HarborCluster.Spec.Version,
			PublicURL:     harbor.HarborCluster.Spec.PublicURL,
			TLSSecretName: harbor.HarborCluster.Spec.TLSSecret,
			Components: v1alpha1.HarborComponents{
				Core:        harbor.newCoreComponent(),
				Portal:      harbor.newPortalComponent(),
				Registry:    harbor.newRegistryComponent(),
				JobService:  harbor.newJobServiceComponent(),
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

func (harbor *HarborReconciler) newCoreComponent() *v1alpha1.CoreComponent {
	return &v1alpha1.CoreComponent{
		HarborDeployment: v1alpha1.HarborDeployment{
			Replicas:         IntToInt32Ptr(harbor.HarborCluster.Spec.Replicas),
			Image:            image.String(harbor.ImageGetter.CoreImage()),
			NodeSelector:     nil,
			ImagePullSecrets: harbor.getImagePullSecrets(),
		},
		DatabaseSecret: harbor.getDatabaseSecret(lcm.CoreSecretForDatabase),
		CacheSecret:    harbor.getCacheSecret(lcm.CoreURLSecretForCache),
	}
}

func (harbor *HarborReconciler) newPortalComponent() *v1alpha1.PortalComponent {
	return &v1alpha1.PortalComponent{
		HarborDeployment: v1alpha1.HarborDeployment{
			Replicas:         IntToInt32Ptr(harbor.HarborCluster.Spec.Replicas),
			Image:            image.String(harbor.ImageGetter.PortalImage()),
			NodeSelector:     nil,
			ImagePullSecrets: harbor.getImagePullSecrets(),
		},
	}
}

func (harbor *HarborReconciler) newRegistryComponent() *v1alpha1.RegistryComponent {
	return &v1alpha1.RegistryComponent{
		HarborDeployment: v1alpha1.HarborDeployment{
			Replicas:         IntToInt32Ptr(harbor.HarborCluster.Spec.Replicas),
			Image:            image.String(harbor.ImageGetter.RegistryImage()),
			NodeSelector:     nil,
			ImagePullSecrets: harbor.getImagePullSecrets(),
		},
		Controller: v1alpha1.RegistryControllerComponent{
			Image: image.String(harbor.ImageGetter.RegistryControllerImage()),
		},
		StorageSecret:   harbor.getStorageSecret(),
		CacheSecret:     harbor.getCacheSecret(lcm.RegisterSecretForCache),
		DisableRedirect: harbor.HarborCluster.Spec.DisableRedirect,
	}
}

func (harbor *HarborReconciler) newJobServiceComponent() *v1alpha1.JobServiceComponent {
	if harbor.HarborCluster.Spec.JobService != nil {
		return &v1alpha1.JobServiceComponent{
			HarborDeployment: v1alpha1.HarborDeployment{
				Replicas:         IntToInt32Ptr(harbor.HarborCluster.Spec.JobService.Replicas),
				Image:            image.String(harbor.ImageGetter.JobServiceImage()),
				NodeSelector:     nil,
				ImagePullSecrets: harbor.getImagePullSecrets(),
			},
			RedisSecret: harbor.getCacheSecret(lcm.JobServiceSecretForCache),
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
				Image:            image.String(harbor.ImageGetter.ChartMuseumImage()),
				NodeSelector:     nil,
				ImagePullSecrets: harbor.getImagePullSecrets(),
			},
			StorageSecret: harbor.getStorageSecretForChartMuseum(),
			CacheSecret:   harbor.getCacheSecret(lcm.ChartMuseumSecretForCache),
		}
	}
	return nil
}

func (harbor *HarborReconciler) newClairComponentIfNecessary() *v1alpha1.ClairComponent {
	if harbor.HarborCluster.Spec.Clair != nil {
		return &v1alpha1.ClairComponent{
			HarborDeployment: v1alpha1.HarborDeployment{
				Replicas:         IntToInt32Ptr(harbor.HarborCluster.Spec.Replicas),
				Image:            image.String(harbor.ImageGetter.ClairImage()),
				NodeSelector:     nil,
				ImagePullSecrets: harbor.getImagePullSecrets(),
			},
			DatabaseSecret:       harbor.getDatabaseSecret(lcm.ClairSecretForDatabase),
			VulnerabilitySources: harbor.HarborCluster.Spec.Clair.VulnerabilitySources,
			Adapter: v1alpha1.ClairAdapterComponent{
				Image:       image.String(harbor.ImageGetter.ClairAdapterImage()),
				RedisSecret: harbor.getCacheSecret(lcm.ClairSecretForCache),
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
				Image: image.String(harbor.ImageGetter.NotaryDBMigratorImage()),
			},
			Signer: v1alpha1.NotarySignerComponent{
				HarborDeployment: v1alpha1.HarborDeployment{
					Replicas:         IntToInt32Ptr(harbor.HarborCluster.Spec.Replicas),
					Image:            image.String(harbor.ImageGetter.NotarySingerImage()),
					NodeSelector:     nil,
					ImagePullSecrets: harbor.getImagePullSecrets(),
				},
				DatabaseSecret: harbor.getDatabaseSecret(lcm.NotarySignerSecretForDatabase),
			},
			Server: v1alpha1.NotaryServerComponent{
				HarborDeployment: v1alpha1.HarborDeployment{
					Replicas:         IntToInt32Ptr(harbor.HarborCluster.Spec.Replicas),
					Image:            image.String(harbor.ImageGetter.NotaryServerImage()),
					NodeSelector:     nil,
					ImagePullSecrets: harbor.getImagePullSecrets(),
				},
				DatabaseSecret: harbor.getDatabaseSecret(lcm.NotaryServerSecretForDatabase),
			},
		}
	}
	return nil
}

// getCacheSecret will get a name of k8s secret which stores cache info
func (harbor *HarborReconciler) getCacheSecret(name string) string {
	p := harbor.getProperty(goharborv1.ComponentCache, name)
	if p != nil {
		return p.ToString()
	}
	return ""
}

// getDatabaseSecret will get a name of k8s secret which stores database info
func (harbor *HarborReconciler) getDatabaseSecret(name string) string {
	p := harbor.getProperty(goharborv1.ComponentDatabase, name)
	if p != nil {
		return p.ToString()
	}
	return ""
}

// getStorageSecret will get a name of k8s secret which stores storage info
func (harbor *HarborReconciler) getStorageSecret() string {
	var name string
	switch harbor.HarborCluster.Spec.Storage.Kind {
	case "azure":
		name = lcm.AzureSecretForStorage
	case "gcs":
		name = lcm.GcsSecretForStorage
	case "swift":
		name = lcm.SwiftSecretForStorage
	case "s3":
		name = lcm.S3SecretForStorage
	case "oss":
		name = lcm.OssSecretForStorage
	default:
		// default in cluster storage has been provided.
		name = lcm.InClusterSecretForStorage
	}
	p := harbor.getProperty(goharborv1.ComponentStorage, name)
	if p != nil {
		return p.ToString()
	}
	return ""
}

// getStorageSecretForChartMuseum will get the secret name of chart museum storage config.
func (harbor *HarborReconciler) getStorageSecretForChartMuseum() string {
	p := harbor.getProperty(goharborv1.ComponentStorage, lcm.ChartMuseumSecretForStorage)
	if p != nil {
		return p.ToString()
	}
	return ""
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
	if harbor.HarborCluster.Spec.ImageSource != nil && harbor.HarborCluster.Spec.ImageSource.ImagePullSecret != "" {
		return []corev1.LocalObjectReference{
			{Name: harbor.HarborCluster.Spec.ImageSource.ImagePullSecret},
		}
	}
	return nil
}

func IntToInt32Ptr(value int) *int32 {
	int32Val := int32(value)
	return &int32Val
}
