package storage

import (
	"encoding/json"
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
	k8serror "k8s.io/apimachinery/pkg/api/errors"
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

	data := map[string]string{
		"accesskey":      string(accessKey),
		"secretkey":      string(secretKey),
		"region":         DefaultRegion,
		"bucket":         DefaultBucket,
		"regionendpoint": m.getServiceName() + "." + m.HarborCluster.Namespace,
		"encrypt":        "false",
		"secure":         "false",
		"v4auth":         "false",
	}
	dataJson, _ := json.Marshal(&data)
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
			s3Storage: dataJson,
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
	data := map[string]string{
		"region":         m.HarborCluster.Spec.Storage.S3.Region,
		"bucket":         m.HarborCluster.Spec.Storage.S3.Bucket,
		"accesskey":      m.HarborCluster.Spec.Storage.S3.AccessKey,
		"secretkey":      m.HarborCluster.Spec.Storage.S3.SecretKey,
		"regionendpoint": m.HarborCluster.Spec.Storage.S3.RegionEndpoint,
		"encrypt":        strconv.FormatBool(m.HarborCluster.Spec.Storage.S3.Encrypt),
		"keyid":          m.HarborCluster.Spec.Storage.S3.KeyId,
		"secure":         strconv.FormatBool(m.HarborCluster.Spec.Storage.S3.Secure),
		"chunksize":      m.HarborCluster.Spec.Storage.S3.ChunkSize,
		"rootdirectory":  m.HarborCluster.Spec.Storage.S3.RootDirectory,
		"storageclass":   m.HarborCluster.Spec.Storage.S3.StorageClass,
		"v4auth":         strconv.FormatBool(m.HarborCluster.Spec.Storage.S3.V4Auth),
	}
	dataJson, _ := json.Marshal(&data)
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
			s3Storage: dataJson,
		},
	}
}

func (m *MinIOReconciler) generateAzureSecret(labels map[string]string) *corev1.Secret {
	data := map[string]string{
		"realm":       m.HarborCluster.Spec.Storage.Azure.Realm,
		"accountname": m.HarborCluster.Spec.Storage.Azure.AccountName,
		"accountkey":  m.HarborCluster.Spec.Storage.Azure.AccountKey,
		"container":   m.HarborCluster.Spec.Storage.Azure.Container,
	}
	dataJson, _ := json.Marshal(&data)
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
			azureStorage: dataJson,
		},
	}
}

func (m *MinIOReconciler) generateGcsSecret(labels map[string]string) *corev1.Secret {
	data := map[string]string{
		"bucket":        m.HarborCluster.Spec.Storage.Gcs.Bucket,
		"encodedkey":    m.HarborCluster.Spec.Storage.Gcs.EncodedKey,
		"rootdirectory": m.HarborCluster.Spec.Storage.Gcs.RootDirectory,
		"chunksize":     m.HarborCluster.Spec.Storage.Gcs.ChunkSize,
	}
	dataJson, _ := json.Marshal(&data)
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
			gcsStorage: dataJson,
		},
	}
}

func (m *MinIOReconciler) generateSwiftSecret(labels map[string]string) *corev1.Secret {
	data := map[string]string{
		"authurl":             m.HarborCluster.Spec.Storage.Swift.Authurl,
		"username":            m.HarborCluster.Spec.Storage.Swift.Username,
		"password":            m.HarborCluster.Spec.Storage.Swift.Password,
		"container":           m.HarborCluster.Spec.Storage.Swift.Container,
		"region":              m.HarborCluster.Spec.Storage.Swift.Region,
		"tenant":              m.HarborCluster.Spec.Storage.Swift.Tenant,
		"tenantid":            m.HarborCluster.Spec.Storage.Swift.TenantId,
		"domain":              m.HarborCluster.Spec.Storage.Swift.Domain,
		"domainid":            m.HarborCluster.Spec.Storage.Swift.DomainId,
		"trustid":             m.HarborCluster.Spec.Storage.Swift.TrustId,
		"insecureskipverify":  strconv.FormatBool(m.HarborCluster.Spec.Storage.Swift.InsecureSkipVerify),
		"prefix":              m.HarborCluster.Spec.Storage.Swift.Prefix,
		"secretkey":           m.HarborCluster.Spec.Storage.Swift.SecretKey,
		"authversion":         string(m.HarborCluster.Spec.Storage.Swift.AuthVersion),
		"endpointtype":        m.HarborCluster.Spec.Storage.Swift.EndpointType,
		"tempurlcontainerkey": strconv.FormatBool(m.HarborCluster.Spec.Storage.Swift.TempurlContainerkey),
		"tempurlmethods":      m.HarborCluster.Spec.Storage.Swift.TempurlMethods,
		"chunksize":           m.HarborCluster.Spec.Storage.Swift.ChunkSize,
	}
	dataJson, _ := json.Marshal(&data)
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
			swiftStorage: dataJson,
		},
	}
}

func (m *MinIOReconciler) generateOssSecret(labels map[string]string) *corev1.Secret {
	data := map[string]string{
		"accesskeyid":     m.HarborCluster.Spec.Storage.Oss.AccessKeyId,
		"accesskeysecret": m.HarborCluster.Spec.Storage.Oss.AccessKeySecret,
		"region":          m.HarborCluster.Spec.Storage.Oss.Region,
		"bucket":          m.HarborCluster.Spec.Storage.Oss.Bucket,
		"endpoint":        m.HarborCluster.Spec.Storage.Oss.Region,
		"internal":        m.HarborCluster.Spec.Storage.Oss.Internal,
		"encrypt":         m.HarborCluster.Spec.Storage.Oss.Encrypt,
		"secure":          m.HarborCluster.Spec.Storage.Oss.Secure,
		"chunksize":       m.HarborCluster.Spec.Storage.Oss.ChunkSize,
		"rootdirectory":   m.HarborCluster.Spec.Storage.Oss.RootDirectory,
	}
	dataJson, _ := json.Marshal(&data)
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
			ossStorage: dataJson,
		},
	}
}

func (m *MinIOReconciler) Provision() (*lcm.CRStatus, error) {
	credsSecret := m.generateCredsSecret()
	err := m.KubeClient.Create(credsSecret)
	if err != nil && !k8serror.IsAlreadyExists(err) {
		return minioNotReadyStatus(CreateMinIOSecretError, err.Error()), err
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

	service := m.generateService()
	err = m.KubeClient.Create(service)
	if err != nil && !k8serror.IsAlreadyExists(err) {
		return minioNotReadyStatus(CreateMinIOServiceError, err.Error()), err
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
