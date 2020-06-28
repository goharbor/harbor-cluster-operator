package storage

import (
	"fmt"
	"github.com/goharbor/harbor-cluster-operator/controllers/common"
	"github.com/goharbor/harbor-cluster-operator/lcm"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/intstr"
	"reflect"
	"strconv"

	goharborv1 "github.com/goharbor/harbor-cluster-operator/api/v1"
	minio "github.com/minio/minio-operator/pkg/apis/operator.min.io/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (m *MinIOReconciler) ProvisionInClusterSecretAsS3(minioInstamnce *minio.MinIOInstance) (*lcm.CRStatus, error) {
	inClusterSecret, err := m.generateInClusterSecret(minioInstamnce)
	if err != nil {
		return minioNotReadyStatus(GetMinIOSecretError, err.Error()), err
	}
	err = m.KubeClient.Create(inClusterSecret)

	p := &lcm.Property{
		Name:  s3Storage + ExternalStorageSecretSuffix,
		Value: inClusterSecret.Name,
	}
	properties := &lcm.Properties{p}
	return minioReadyStatus(properties), err
}

func (m *MinIOReconciler) generateInClusterSecret(minioInstamnce *minio.MinIOInstance) (*corev1.Secret, error) {
	labels := m.getLabels()
	labels[LabelOfStorageType] = inClusterStorage
	accessKey, secretKey, err := m.getCredsFromSecret()
	if err != nil {
		return nil, err
	}

	inClusterSecret := &corev1.Secret{
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
			"accesskey":      accessKey,
			"secretkey":      secretKey,
			"region":         []byte(DefaultRegion),
			"bucket":         []byte(DefaultBucket),
			"regionendpoint": []byte(m.getServiceName() + "." + m.HarborCluster.Namespace),
			"encrypt":        []byte("false"),
			"secure":         []byte("false"),
			"v4auth":         []byte("false"),
		},
	}

	return inClusterSecret, nil
}

func (m *MinIOReconciler) ProvisionExternalStorage() (*lcm.CRStatus, error) {
	exSecret, err := m.generateExternalSecret()
	if err != nil {
		return minioNotReadyStatus(err.Error(), err.Error()), err
	}

	err = m.KubeClient.Create(exSecret)
	if err != nil {
		minioNotReadyStatus(CreateExternalSecretError, err.Error())
	}

	p := &lcm.Property{
		Name:  m.HarborCluster.Spec.Storage.Kind + ExternalStorageSecretSuffix,
		Value: m.getExternalSecretName(),
	}
	properties := &lcm.Properties{p}

	return minioReadyStatus(properties), nil
}

func (m *MinIOReconciler) generateExternalSecret() (*corev1.Secret, error) {
	var exSecret *corev1.Secret
	labels := m.getLabels()

	switch m.HarborCluster.Spec.Storage.Kind {
	case azureStorage:
		labels[LabelOfStorageType] = azureStorage
		exSecret = m.generateAzureSecret(labels)
	case gcsStorage:
		labels[LabelOfStorageType] = gcsStorage
		exSecret = m.generateGcsSecret(labels)
	case s3Storage:
		labels[LabelOfStorageType] = s3Storage
		exSecret = m.generateS3Secret(labels)
	case swiftStorage:
		labels[LabelOfStorageType] = swiftStorage
		exSecret = m.generateSwiftSecret(labels)
	case ossStorage:
		labels[LabelOfStorageType] = ossStorage
		exSecret = m.generateOssSecret(labels)
	default:
		return exSecret, fmt.Errorf(NotSupportType)
	}

	return exSecret, nil
}

func (m *MinIOReconciler) generateS3Secret(labels map[string]string) *corev1.Secret {
	return &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Secret",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        m.getExternalSecretName(),
			Namespace:   m.HarborCluster.Namespace,
			Labels:      labels,
			Annotations: m.generateAnnotations(),
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(m.HarborCluster, goharborv1.HarborClusterGVK),
			},
		},
		Type: corev1.SecretTypeOpaque,
		Data: map[string][]byte{
			"region":         []byte(m.HarborCluster.Spec.Storage.S3.Region),
			"bucket":         []byte(m.HarborCluster.Spec.Storage.S3.Bucket),
			"accesskey":      []byte(m.HarborCluster.Spec.Storage.S3.AccessKey),
			"secretkey":      []byte(m.HarborCluster.Spec.Storage.S3.SecretKey),
			"regionendpoint": []byte(m.HarborCluster.Spec.Storage.S3.RegionEndpoint),
			"encrypt":        []byte(strconv.FormatBool(m.HarborCluster.Spec.Storage.S3.Encrypt)),
			"keyid":          []byte(m.HarborCluster.Spec.Storage.S3.KeyId),
			"secure":         []byte(strconv.FormatBool(m.HarborCluster.Spec.Storage.S3.Secure)),
			"chunksize":      []byte(m.HarborCluster.Spec.Storage.S3.ChunkSize),
			"rootdirectory":  []byte(m.HarborCluster.Spec.Storage.S3.RootDirectory),
			"storageclass":   []byte(m.HarborCluster.Spec.Storage.S3.StorageClass),
			"v4auth":         []byte(strconv.FormatBool(m.HarborCluster.Spec.Storage.S3.V4Auth))},
	}
}

func (m *MinIOReconciler) generateAzureSecret(labels map[string]string) *corev1.Secret {
	return &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Secret",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        m.getExternalSecretName(),
			Namespace:   m.HarborCluster.Namespace,
			Labels:      labels,
			Annotations: m.generateAnnotations(),
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(m.HarborCluster, goharborv1.HarborClusterGVK),
			},
		},
		Type: corev1.SecretTypeOpaque,
		Data: map[string][]byte{
			"realm":       []byte(m.HarborCluster.Spec.Storage.Azure.Realm),
			"accountname": []byte(m.HarborCluster.Spec.Storage.Azure.AccountName),
			"accountkey":  []byte(m.HarborCluster.Spec.Storage.Azure.AccountKey),
			"container":   []byte(m.HarborCluster.Spec.Storage.Azure.Container),
		},
	}
}

func (m *MinIOReconciler) generateGcsSecret(labels map[string]string) *corev1.Secret {
	return &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Secret",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        m.getExternalSecretName(),
			Namespace:   m.HarborCluster.Namespace,
			Labels:      labels,
			Annotations: m.generateAnnotations(),
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(m.HarborCluster, goharborv1.HarborClusterGVK),
			},
		},
		Type: corev1.SecretTypeOpaque,
		Data: map[string][]byte{
			"bucket":        []byte(m.HarborCluster.Spec.Storage.Gcs.Bucket),
			"encodedkey":    []byte(m.HarborCluster.Spec.Storage.Gcs.EncodedKey),
			"rootdirectory": []byte(m.HarborCluster.Spec.Storage.Gcs.RootDirectory),
			"chunksize":     []byte(m.HarborCluster.Spec.Storage.Gcs.ChunkSize)},
	}
}

func (m *MinIOReconciler) generateSwiftSecret(labels map[string]string) *corev1.Secret {
	return &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Secret",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        m.getExternalSecretName(),
			Namespace:   m.HarborCluster.Namespace,
			Labels:      labels,
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
			"tenantid":            []byte(m.HarborCluster.Spec.Storage.Swift.TenantId),
			"domain":              []byte(m.HarborCluster.Spec.Storage.Swift.Domain),
			"domainid":            []byte(m.HarborCluster.Spec.Storage.Swift.DomainId),
			"trustid":             []byte(m.HarborCluster.Spec.Storage.Swift.TrustId),
			"insecureskipverify":  []byte(strconv.FormatBool(m.HarborCluster.Spec.Storage.Swift.InsecureSkipVerify)),
			"prefix":              []byte(m.HarborCluster.Spec.Storage.Swift.Prefix),
			"secretkey":           []byte(m.HarborCluster.Spec.Storage.Swift.SecretKey),
			"authversion":         []byte(string(m.HarborCluster.Spec.Storage.Swift.AuthVersion)),
			"endpointtype":        []byte(m.HarborCluster.Spec.Storage.Swift.EndpointType),
			"tempurlcontainerkey": []byte(strconv.FormatBool(m.HarborCluster.Spec.Storage.Swift.TempurlContainerkey)),
			"tempurlmethods":      []byte(m.HarborCluster.Spec.Storage.Swift.TempurlMethods),
			"chunksize":           []byte(m.HarborCluster.Spec.Storage.Swift.ChunkSize)},
	}
}

func (m *MinIOReconciler) generateOssSecret(labels map[string]string) *corev1.Secret {
	return &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Secret",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        m.getExternalSecretName(),
			Namespace:   m.HarborCluster.Namespace,
			Labels:      labels,
			Annotations: m.generateAnnotations(),
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(m.HarborCluster, goharborv1.HarborClusterGVK),
			},
		},
		Type: corev1.SecretTypeOpaque,
		Data: map[string][]byte{
			"accesskeyid":     []byte(m.HarborCluster.Spec.Storage.Oss.AccessKeyId),
			"accesskeysecret": []byte(m.HarborCluster.Spec.Storage.Oss.AccessKeySecret),
			"region":          []byte(m.HarborCluster.Spec.Storage.Oss.Region),
			"bucket":          []byte(m.HarborCluster.Spec.Storage.Oss.Bucket),
			"endpoint":        []byte(m.HarborCluster.Spec.Storage.Oss.Region),
			"internal":        []byte(m.HarborCluster.Spec.Storage.Oss.Internal),
			"encrypt":         []byte(m.HarborCluster.Spec.Storage.Oss.Encrypt),
			"secure":          []byte(m.HarborCluster.Spec.Storage.Oss.Secure),
			"chunksize":       []byte(m.HarborCluster.Spec.Storage.Oss.ChunkSize),
			"rootdirectory":   []byte(m.HarborCluster.Spec.Storage.Oss.RootDirectory),
		},
	}
}

func (m *MinIOReconciler) Provision() (*lcm.CRStatus, error) {
	credsSecret := m.generateCredsSecret()
	err := m.KubeClient.Create(credsSecret)
	if err != nil {
		return minioNotReadyStatus(CreateMinIOSecretError, err.Error()), err
	}
	service := m.generateService()
	err = m.KubeClient.Create(service)
	if err != nil {
		return minioNotReadyStatus(CreateMinIOServiceError, err.Error()), err
	}

	err = m.KubeClient.Create(m.DesiredMinIOCR)
	if err != nil {
		return minioNotReadyStatus(CreateMinIOError, err.Error()), err
	}

	var minioCR minio.MinIOInstance
	err = m.KubeClient.Get(m.getMinIONamespacedName(), &minioCR)
	if err != nil {
		return minioNotReadyStatus(CreateMinIOError, err.Error()), err
	}

	credsSecret.OwnerReferences = []metav1.OwnerReference{
		*metav1.NewControllerRef(&minioCR, HarborClusterMinIOGVK),
	}
	err = m.KubeClient.Update(credsSecret)
	if err != nil {
		return minioNotReadyStatus(CreateMinIOError, err.Error()), err
	}

	service.OwnerReferences = []metav1.OwnerReference{
		*metav1.NewControllerRef(&minioCR, HarborClusterMinIOGVK),
	}
	err = m.KubeClient.Update(service)
	if err != nil {
		return minioNotReadyStatus(CreateMinIOServiceError, err.Error()), err
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
				{
					Name:    m.HarborCluster.Name + "-" + DefaultZone,
					Servers: m.HarborCluster.Spec.Storage.InCluster.Spec.Replicas,
				},
			},
			VolumesPerServer:    1,
			Mountpath:           minio.MinIOVolumeMountPath,
			VolumeClaimTemplate: m.getVolumeClaimTemplate(),
			CredsSecret: &corev1.LocalObjectReference{
				Name: m.getMinIOSecretNamespacedName().Name,
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
				{
					Port:       9000,
					TargetPort: intstr.FromInt(9000),
					Protocol:   corev1.ProtocolTCP,
				},
			},
		},
	}
}

func (m *MinIOReconciler) generateCredsSecret() *corev1.Secret {
	credsAccesskey := common.RandomString(8, "a")
	credsSecretkey := common.RandomString(8, "a")

	return &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Secret",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        m.getMinIOSecretNamespacedName().Name,
			Namespace:   m.HarborCluster.Namespace,
			Labels:      m.getLabels(),
			Annotations: m.generateAnnotations(),
		},
		Type: corev1.SecretTypeOpaque,
		Data: map[string][]byte{
			"accesskey": []byte(credsAccesskey),
			"secretkey": []byte(credsSecretkey)},
	}
}

func (m *MinIOReconciler) getCredsFromSecret() ([]byte, []byte, error) {
	var minIOSecret corev1.Secret
	namespaced := m.getMinIOSecretNamespacedName()
	err := m.KubeClient.Get(namespaced, &minIOSecret)
	return minIOSecret.Data["accesskey"], minIOSecret.Data["secretkey"], err
}
