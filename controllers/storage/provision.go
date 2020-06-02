package storage

import (
	"github.com/goharbor/harbor-cluster-operator/lcm"
	"k8s.io/apimachinery/pkg/api/resource"
	"reflect"
	"strconv"

	goharborv1 "github.com/goharbor/harbor-cluster-operator/api/v1"
	minio "github.com/minio/minio-operator/pkg/apis/operator.min.io/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	DefaultExternalSecretSuffix = "harbor-cluster-storage"
	S3Secret                    = "s3Secret"
	AzureSecret                 = "azureSecret"
	GcsSecret                   = "gcsSecret"
	SwiftSecret                 = "swiftSecret"
	OssSecret                   = "ossSecret"
	DefaultCredsSecret          = "minio-creds-secret"
	DefaultMcsSecret            = "minio-mcs-secret"
	CredsAccesskey              = "bWluaW8="
	CredsSecretkey              = "bWluaW8xMjM="
	DefaultZone                 = "zone-harbor"
	DefaultMinIO                = "minio"
	DefaultRegion               = "us-east-1"
	DefaultBucket               = "harbor"
)

func (m *MinIOReconciler) ProvisionInClusterSecretAsOss(minioInstamnce *minio.MinIOInstance) (*lcm.CRStatus, error) {
	inClusterSecret := m.generateInClusterSecret(minioInstamnce)
	err := m.KubeClient.Create(inClusterSecret)

	p := &lcm.Property{
		Name:  OssSecret,
		Value: inClusterSecret.Name,
	}
	properties := &lcm.Properties{p}
	return minioReadyStatus(properties), err
}

func (m *MinIOReconciler) generateInClusterSecret(minioInstamnce *minio.MinIOInstance) *corev1.Secret {
	return &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Secret",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        m.getServiceName(),
			Namespace:   m.HarborCluster.Namespace,
			Labels:      m.getLabels(),
			Annotations: m.generateAnnotations(),
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(minioInstamnce, HarborClusterMinIOGVK),
			},
		},
		Type: corev1.SecretTypeOpaque,
		Data: map[string][]byte{
			// TODO how about using random passwords ??
			"accesskeyid":     []byte("minio"),
			"accesskeysecret": []byte("minio123"),
			"region":          []byte(DefaultRegion),
			"bucket":          []byte(DefaultBucket),
			"endpoint":        []byte(m.getServiceName() + "." + m.HarborCluster.Namespace),
		},
	}
}

func (m *MinIOReconciler) ProvisionExternalStorage() (*lcm.CRStatus, error) {
	switch m.HarborCluster.Spec.Storage.Kind {
	case azureStorage:
		properties, err := m.ProvisionAzure()
		if err != nil {
			minioNotReadyStatus(ErrorReason1, err.Error())
		}
		return minioReadyStatus(properties), nil
	case gcsStorage:
		properties, err := m.ProvisionGcs()
		if err != nil {
			minioNotReadyStatus(ErrorReason1, err.Error())
		}
		return minioReadyStatus(properties), nil
	case s3Storage:
		properties, err := m.ProvisionS3()
		if err != nil {
			minioNotReadyStatus(ErrorReason1, err.Error())
		}
		return minioReadyStatus(properties), nil
	case swiftStorage:
		properties, err := m.ProvisionSwift()
		if err != nil {
			minioNotReadyStatus(ErrorReason1, err.Error())
		}
		return minioReadyStatus(properties), nil
	case ossStorage:
		properties, err := m.ProvisionOss()
		if err != nil {
			minioNotReadyStatus(ErrorReason1, err.Error())
		}
		return minioReadyStatus(properties), nil
	default:
		return minioNotReadyStatus(ErrorReason3, ErrorReason3), nil
	}
}

func (m *MinIOReconciler) ProvisionS3() (*lcm.Properties, error) {
	s3Secret := m.generateS3Secret()
	err := m.KubeClient.Create(s3Secret)
	p := &lcm.Property{
		Name:  S3Secret,
		Value: s3Secret.Name,
	}
	properties := &lcm.Properties{p}
	return properties, err
}

func (m *MinIOReconciler) generateS3Secret() *corev1.Secret {
	return &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Secret",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        m.HarborCluster.Name + "-" + DefaultExternalSecretSuffix,
			Namespace:   m.HarborCluster.Namespace,
			Labels:      m.getLabels(),
			Annotations: m.generateAnnotations(),
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(m.HarborCluster, goharborv1.HarborClusterGVK),
			},
		},
		Type: corev1.SecretTypeOpaque,
		Data: map[string][]byte{
			"region":         []byte(m.HarborCluster.Spec.Storage.S3.Region),
			"bucket":         []byte(m.HarborCluster.Spec.Storage.S3.Bucket),
			"accesskey":      []byte(m.HarborCluster.Spec.Storage.S3.Accesskey),
			"secretkey":      []byte(m.HarborCluster.Spec.Storage.S3.Secretkey),
			"regionendpoint": []byte(m.HarborCluster.Spec.Storage.S3.Regionendpoint),
			"encrypt":        []byte(strconv.FormatBool(m.HarborCluster.Spec.Storage.S3.Encrypt)),
			"keyid":          []byte(m.HarborCluster.Spec.Storage.S3.Keyid),
			"secure":         []byte(strconv.FormatBool(m.HarborCluster.Spec.Storage.S3.Secure)),
			"chunksize":      []byte(m.HarborCluster.Spec.Storage.S3.Chunksize),
			"rootdirectory":  []byte(m.HarborCluster.Spec.Storage.S3.Rootdirectory),
			"storageclass":   []byte(m.HarborCluster.Spec.Storage.S3.Storageclass),
			"v4auth":         []byte(strconv.FormatBool(m.HarborCluster.Spec.Storage.S3.V4auth))},
	}
}

func (m *MinIOReconciler) ProvisionAzure() (*lcm.Properties, error) {
	azureSecret := m.generateAzureSecret()
	err := m.KubeClient.Create(azureSecret)
	p := &lcm.Property{
		Name:  AzureSecret,
		Value: azureSecret.Name,
	}
	properties := &lcm.Properties{p}
	return properties, err
}

func (m *MinIOReconciler) generateAzureSecret() *corev1.Secret {
	return &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Secret",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        m.HarborCluster.Name + "-" + DefaultExternalSecretSuffix,
			Namespace:   m.HarborCluster.Namespace,
			Labels:      m.getLabels(),
			Annotations: m.generateAnnotations(),
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(m.HarborCluster, goharborv1.HarborClusterGVK),
			},
		},
		Type: corev1.SecretTypeOpaque,
		Data: map[string][]byte{
			"authurl":   []byte(m.HarborCluster.Spec.Storage.Swift.Authurl),
			"username":  []byte(m.HarborCluster.Spec.Storage.Swift.Username),
			"password":  []byte(m.HarborCluster.Spec.Storage.Swift.Password),
			"container": []byte(m.HarborCluster.Spec.Storage.Azure.Container),
		},
	}
}

func (m *MinIOReconciler) ProvisionGcs() (*lcm.Properties, error) {
	gcsSecret := m.generateGcsSecret()
	err := m.KubeClient.Create(gcsSecret)
	p := &lcm.Property{
		Name:  GcsSecret,
		Value: gcsSecret.Name,
	}
	properties := &lcm.Properties{p}
	return properties, err
}

func (m *MinIOReconciler) generateGcsSecret() *corev1.Secret {
	return &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Secret",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        m.HarborCluster.Name + "-" + DefaultExternalSecretSuffix,
			Namespace:   m.HarborCluster.Namespace,
			Labels:      m.getLabels(),
			Annotations: m.generateAnnotations(),
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(m.HarborCluster, goharborv1.HarborClusterGVK),
			},
		},
		Type: corev1.SecretTypeOpaque,
		Data: map[string][]byte{
			"bucket":        []byte(m.HarborCluster.Spec.Storage.Gcs.Bucket),
			"encodedkey":    []byte(m.HarborCluster.Spec.Storage.Gcs.Encodedkey),
			"rootdirectory": []byte(m.HarborCluster.Spec.Storage.Gcs.Rootdirectory),
			"chunksize":     []byte(m.HarborCluster.Spec.Storage.Gcs.Chunksize)},
	}
}

func (m *MinIOReconciler) ProvisionSwift() (*lcm.Properties, error) {
	swiftSecret := m.generateSwiftSecret()
	err := m.KubeClient.Create(swiftSecret)
	p := &lcm.Property{
		Name:  SwiftSecret,
		Value: swiftSecret.Name,
	}
	properties := &lcm.Properties{p}
	return properties, err
}

func (m *MinIOReconciler) generateSwiftSecret() *corev1.Secret {
	return &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Secret",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        m.HarborCluster.Name + "-" + DefaultExternalSecretSuffix,
			Namespace:   m.HarborCluster.Namespace,
			Labels:      m.getLabels(),
			Annotations: m.generateAnnotations(),
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(m.HarborCluster, goharborv1.HarborClusterGVK),
			},
		},
		Type: corev1.SecretTypeOpaque,
		Data: map[string][]byte{
			"authurl":             []byte(m.HarborCluster.Spec.Storage.Swift.Authurl),
			"username":            []byte(m.HarborCluster.Spec.Storage.Swift.Username),
			"password":            []byte(m.HarborCluster.Spec.Storage.Swift.Password),
			"container":           []byte(m.HarborCluster.Spec.Storage.Swift.Container),
			"region":              []byte(m.HarborCluster.Spec.Storage.Swift.Region),
			"tenant":              []byte(m.HarborCluster.Spec.Storage.Swift.Tenant),
			"tenantid":            []byte(m.HarborCluster.Spec.Storage.Swift.Tenantid),
			"domain":              []byte(m.HarborCluster.Spec.Storage.Swift.Domain),
			"domainid":            []byte(m.HarborCluster.Spec.Storage.Swift.Domainid),
			"trustid":             []byte(m.HarborCluster.Spec.Storage.Swift.Trustid),
			"insecureskipverify":  []byte(strconv.FormatBool(m.HarborCluster.Spec.Storage.Swift.Insecureskipverify)),
			"prefix":              []byte(m.HarborCluster.Spec.Storage.Swift.Prefix),
			"secretkey":           []byte(m.HarborCluster.Spec.Storage.Swift.Secretkey),
			"authversion":         []byte(string(m.HarborCluster.Spec.Storage.Swift.AuthVersion)),
			"endpointtype":        []byte(m.HarborCluster.Spec.Storage.Swift.EndpointType),
			"tempurlcontainerkey": []byte(strconv.FormatBool(m.HarborCluster.Spec.Storage.Swift.TempurlContainerkey)),
			"tempurlmethods":      []byte(m.HarborCluster.Spec.Storage.Swift.TempurlMethods),
			"chunksize":           []byte(m.HarborCluster.Spec.Storage.Swift.Chunksize)},
	}
}

func (m *MinIOReconciler) ProvisionOss() (*lcm.Properties, error) {
	ossSecret := m.generateOssSecret()
	err := m.KubeClient.Create(ossSecret)
	p := &lcm.Property{
		Name:  OssSecret,
		Value: ossSecret.Name,
	}
	properties := &lcm.Properties{p}
	return properties, err
}

func (m *MinIOReconciler) generateOssSecret() *corev1.Secret {
	return &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Secret",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        m.HarborCluster.Name + "-" + DefaultExternalSecretSuffix,
			Namespace:   m.HarborCluster.Namespace,
			Labels:      m.getLabels(),
			Annotations: m.generateAnnotations(),
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(m.HarborCluster, goharborv1.HarborClusterGVK),
			},
		},
		Type: corev1.SecretTypeOpaque,
		Data: map[string][]byte{
			"accesskeyid":     []byte(m.HarborCluster.Spec.Storage.Oss.Accesskeyid),
			"accesskeysecret": []byte(m.HarborCluster.Spec.Storage.Oss.Accesskeysecret),
			"region":          []byte(m.HarborCluster.Spec.Storage.Oss.Region),
			"bucket":          []byte(m.HarborCluster.Spec.Storage.Oss.Bucket),
			"endpoint":        []byte(m.HarborCluster.Spec.Storage.Oss.Region),
			"internal":        []byte(m.HarborCluster.Spec.Storage.Oss.Internal),
			"encrypt":         []byte(m.HarborCluster.Spec.Storage.Oss.Encrypt),
			"secure":          []byte(m.HarborCluster.Spec.Storage.Oss.Secure),
			"chunksize":       []byte(m.HarborCluster.Spec.Storage.Oss.Secure),
			"rootdirectory":   []byte(m.HarborCluster.Spec.Storage.Oss.RootDirectory),
		},
	}
}

func (m *MinIOReconciler) Provision() (*lcm.CRStatus, error) {
	// TODO remove mcs secret ref https://github.com/minio/minio-operator/blob/master/examples/minioinstance.yaml
	//mcsSecret := m.generateMcsSecret()
	//err := m.KubeClient.Create(mcsSecret)
	//if err != nil {
	//	return minioNotReadyStatus(ErrorReason2, err.Error()), err
	//}
	credsSecret := m.generateCredsSecret()
	err := m.KubeClient.Create(credsSecret)
	if err != nil {
		return minioNotReadyStatus(ErrorReason2, err.Error()), err
	}
	service := m.generateService()
	err = m.KubeClient.Create(service)
	if err != nil {
		return minioNotReadyStatus(ErrorReason4, err.Error()), err
	}

	minioCreate := m.generateMinIOCR()
	err = m.KubeClient.Create(minioCreate)
	if err != nil {
		return minioNotReadyStatus(ErrorReason5, err.Error()), err
	}
	var minioCR minio.MinIOInstance
	err = m.KubeClient.Get(m.getMinIONamespacedName(), &minioCR)
	if err != nil {
		return minioNotReadyStatus(ErrorReason5, err.Error()), err
	}

	credsSecret.OwnerReferences = []metav1.OwnerReference{
		*metav1.NewControllerRef(minioCreate, HarborClusterMinIOGVK),
	}
	err = m.KubeClient.Update(credsSecret)
	if err != nil {
		return minioNotReadyStatus(ErrorReason5, err.Error()), err
	}

	service.OwnerReferences = []metav1.OwnerReference{
		*metav1.NewControllerRef(minioCreate, HarborClusterMinIOGVK),
	}
	err = m.KubeClient.Update(service)
	if err != nil {
		return minioNotReadyStatus(ErrorReason4, err.Error()), err
	}

	return minioUnknownStatus(), nil
}

func (m *MinIOReconciler) generateMinIOCR() *minio.MinIOInstance {
	return &minio.MinIOInstance{
		TypeMeta: metav1.TypeMeta{
			Kind:       minio.MinIOCRDResourceKind,
			APIVersion: minio.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        m.getServiceName(),
			Namespace:   m.HarborCluster.Namespace,
			Labels:      m.getLabels(),
			Annotations: m.generateAnnotations(),
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(m.HarborCluster, goharborv1.HarborClusterGVK),
			},
		},
		Spec: minio.MinIOInstanceSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: m.getLabels(),
			},
			Metadata: &metav1.ObjectMeta{
				Labels:      m.getLabels(),
				Annotations: m.generateAnnotations(),
			},
			ServiceName: m.getServiceName(),
			Image:       "minio/minio:" + m.HarborCluster.Spec.Storage.InCluster.Spec.Version,
			Zones: []minio.Zone{
				minio.Zone{
					Name:    m.HarborCluster.Name + "-" + DefaultZone,
					Servers: m.HarborCluster.Spec.Storage.InCluster.Spec.Replicas,
				},
			},
			VolumesPerServer:    1,
			Mountpath:           minio.MinIOVolumeMountPath,
			VolumeClaimTemplate: m.getVolumeClaimTemplate(),
			CredsSecret: &corev1.LocalObjectReference{
				Name: m.HarborCluster.Name + "-" + DefaultCredsSecret,
			},
			PodManagementPolicy: "Parallel",
			RequestAutoCert:     false,
			CertConfig: &minio.CertificateConfig{
				CommonName:       "",
				OrganizationName: []string{},
				DNSNames:         []string{},
			},
			Env: []corev1.EnvVar{
				corev1.EnvVar{
					Name:  "MINIO_BROWSER",
					Value: "on",
				},
			},
			Resources: *m.getResourceRequirements(), //m.HarborCluster.Spec.Storage.InCluster.Spec.Resources,
			Liveness: &minio.Liveness{
				InitialDelaySeconds: 120,
				PeriodSeconds:       60,
			},
			Readiness: &minio.Readiness{
				InitialDelaySeconds: 120,
				PeriodSeconds:       60,
			},
		},
	}
}

func (m *MinIOReconciler) getServiceName() string {
	return m.HarborCluster.Name + "-" + DefaultMinIO
}

func (m *MinIOReconciler) getResourceRequirements() *corev1.ResourceRequirements {
	isEmpty := reflect.DeepEqual(m.HarborCluster.Spec.Storage.InCluster.Spec.Resources, corev1.ResourceRequirements{})
	if !isEmpty {
		return &m.HarborCluster.Spec.Storage.InCluster.Spec.Resources
	}
	limits := map[corev1.ResourceName]resource.Quantity{
		corev1.ResourceCPU:    resource.MustParse("250m"),
		corev1.ResourceMemory: resource.MustParse("512Mi"),
	}
	requests := map[corev1.ResourceName]resource.Quantity{
		corev1.ResourceCPU:    resource.MustParse("250m"),
		corev1.ResourceMemory: resource.MustParse("512Mi"),
	}
	return &corev1.ResourceRequirements{
		Limits:   limits,
		Requests: requests,
	}
}

func (m *MinIOReconciler) getVolumeClaimTemplate() *corev1.PersistentVolumeClaim {
	isEmpty := reflect.DeepEqual(m.HarborCluster.Spec.Storage.InCluster.Spec.VolumeClaimTemplate, corev1.PersistentVolumeClaim{})
	if !isEmpty {
		return &m.HarborCluster.Spec.Storage.InCluster.Spec.VolumeClaimTemplate
	}
	defaultStorageClass := "default"
	return &corev1.PersistentVolumeClaim{
		Spec: corev1.PersistentVolumeClaimSpec{
			StorageClassName: &defaultStorageClass,
			AccessModes:      []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
			Resources: corev1.ResourceRequirements{
				Requests: map[corev1.ResourceName]resource.Quantity{
					corev1.ResourceStorage: resource.MustParse("10Gi"),
				},
			},
		},
	}
}

func (m *MinIOReconciler) getLabels() map[string]string {
	return map[string]string{"type": "harbor-cluster-minio", "app": "minio"}
}

func (m *MinIOReconciler) generateAnnotations() map[string]string {
	// TODO
	return nil
}

func (m *MinIOReconciler) generateService() *corev1.Service {
	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        m.getServiceName(),
			Namespace:   m.HarborCluster.Namespace,
			Labels:      m.getLabels(),
			Annotations: m.generateAnnotations(),
		},
		Spec: corev1.ServiceSpec{
			Type:     corev1.ServiceTypeClusterIP,
			Selector: m.getLabels(),
			Ports: []corev1.ServicePort{
				corev1.ServicePort{
					Port:       9000,
					TargetPort: intstr.FromInt(9000),
					Protocol:   corev1.ProtocolTCP,
				},
			},
		},
	}
}

func (m *MinIOReconciler) generateMcsSecret() *corev1.Secret {
	return &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Secret",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        m.HarborCluster.Name + "-" + DefaultMcsSecret,
			Namespace:   m.HarborCluster.Namespace,
			Labels:      m.getLabels(),
			Annotations: m.generateAnnotations(),
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(m.HarborCluster, goharborv1.HarborClusterGVK),
			},
		},
		Type: corev1.SecretTypeOpaque,
		Data: map[string][]byte{
			"mcshmacjwt":         []byte("WU9VUkpXVFNJR05JTkdTRUNSRVQ="),
			"mcspbkdfpassphrase": []byte("U0VDUkVU"),
			"mcspbkdfsalt":       []byte("U0VDUkVU"),
			"mcssecretkey":       []byte("WU9VUk1DU1NFQ1JFVA")},
	}
}

func (m *MinIOReconciler) generateCredsSecret() *corev1.Secret {
	return &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Secret",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        m.HarborCluster.Name + "-" + DefaultCredsSecret,
			Namespace:   m.HarborCluster.Namespace,
			Labels:      m.getLabels(),
			Annotations: m.generateAnnotations(),
		},
		Type: corev1.SecretTypeOpaque,
		Data: map[string][]byte{
			"accesskey": []byte(CredsAccesskey),
			"secretkey": []byte(CredsSecretkey)},
	}
}
