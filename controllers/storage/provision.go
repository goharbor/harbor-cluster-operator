package storage

import (
	"encoding/json"
	"fmt"
	"reflect"

	goharborv1 "github.com/goharbor/harbor-cluster-operator/api/v1"
	"github.com/goharbor/harbor-cluster-operator/controllers/common"
	minio "github.com/goharbor/harbor-cluster-operator/controllers/storage/minio/api/v1"
	"github.com/goharbor/harbor-cluster-operator/lcm"
	corev1 "k8s.io/api/core/v1"
	k8serror "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (m *MinIOReconciler) ProvisionInClusterSecretAsS3(minioInstamnce *minio.Tenant) (*lcm.CRStatus, error) {
	inClusterSecret, chartMuseumSecret, err := m.generateInClusterSecret(minioInstamnce)
	if err != nil {
		return minioNotReadyStatus(GetMinIOSecretError, err.Error()), err
	}
	err = m.KubeClient.Create(inClusterSecret)
	if err != nil && !k8serror.IsAlreadyExists(err) {
		return minioNotReadyStatus(GetMinIOSecretError, err.Error()), err
	}

	err = m.KubeClient.Create(chartMuseumSecret)
	if err != nil && !k8serror.IsAlreadyExists(err) {
		return minioNotReadyStatus(CreateChartMuseumStorageSecretError, err.Error()), err
	}

	properties := &lcm.Properties{}
	properties.Add(lcm.InClusterSecretForStorage, inClusterSecret.Name)
	if m.HarborCluster.Spec.ChartMuseum != nil {
		properties.Add(lcm.ChartMuseumSecretForStorage, m.getChartMuseumSecretName())
	}

	return minioReadyStatus(properties), nil
}

func (m *MinIOReconciler) generateInClusterSecret(minioInstance *minio.Tenant) (inClusterSecret *corev1.Secret, chartMuseumSecret *corev1.Secret, err error) {
	labels := m.getLabels()
	labels[LabelOfStorageType] = inClusterStorage
	accessKey, secretKey, err := m.getCredsFromSecret()
	if err != nil {
		return nil, nil, err
	}

	endpoint := fmt.Sprintf("http://%s.%s.svc:%s", m.getServiceName(), m.HarborCluster.Namespace, "9000")

	data := map[string]string{
		"accesskey":      string(accessKey),
		"secretkey":      string(secretKey),
		"region":         DefaultRegion,
		"bucket":         DefaultBucket,
		"regionendpoint": endpoint,
		"encrypt":        "false",
		"secure":         "false",
		"v4auth":         "false",
	}
	dataJson, _ := json.Marshal(&data)
	inClusterSecret = &corev1.Secret{
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
				*metav1.NewControllerRef(minioInstance, HarborClusterMinIOGVK),
			},
		},
		Type: corev1.SecretTypeOpaque,
		Data: map[string][]byte{
			s3Storage: dataJson,
		},
	}

	chartMuseumSecret = &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Secret",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        m.getChartMuseumSecretName(),
			Namespace:   m.HarborCluster.Namespace,
			Labels:      m.getLabels(),
			Annotations: m.generateAnnotations(),
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(minioInstance, goharborv1.HarborClusterGVK),
			},
		},
		Type: corev1.SecretTypeOpaque,
		Data: map[string][]byte{
			"kind":                  []byte("amazon"),
			"AWS_ACCESS_KEY_ID":     accessKey,
			"AWS_SECRET_ACCESS_KEY": secretKey,
			// use same bucket.
			"AMAZON_BUCKET":   []byte(DefaultBucket),
			"AMAZON_PREFIX":   []byte(fmt.Sprintf("%s-subfloder", DefaultBucket)),
			"AMAZON_REGION":   []byte(DefaultRegion),
			"AMAZON_ENDPOINT": []byte(endpoint),
		},
	}

	return inClusterSecret, chartMuseumSecret, nil
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

	chartMuseumSecret, err := m.generateSecretForChartMuseum()
	if err != nil {
		return minioNotReadyStatus(GenerateChartMuseumStorageSecretError, err.Error()), err
	}

	err = m.KubeClient.Create(chartMuseumSecret)
	if err != nil {
		minioNotReadyStatus(CreateChartMuseumStorageSecretError, err.Error())
	}

	properties := &lcm.Properties{}
	properties.Add(m.HarborCluster.Spec.Storage.Kind+ExternalStorageSecretSuffix, m.getExternalSecretName())
	if m.HarborCluster.Spec.ChartMuseum != nil {
		properties.Add(lcm.ChartMuseumSecretForStorage, chartMuseumSecret.Name)
	}

	return minioReadyStatus(properties), nil
}

func (m *MinIOReconciler) generateExternalSecret() (exSecret *corev1.Secret, err error) {
	labels := m.getLabels()

	switch m.HarborCluster.Spec.Storage.Kind {
	case azureStorage:
		labels[LabelOfStorageType] = azureStorage
		exSecret, err = m.generateAzureSecret(labels)
	case gcsStorage:
		labels[LabelOfStorageType] = gcsStorage
		exSecret, err = m.generateGcsSecret(labels)
	case s3Storage:
		labels[LabelOfStorageType] = s3Storage
		exSecret, err = m.generateS3Secret(labels)
	case swiftStorage:
		labels[LabelOfStorageType] = swiftStorage
		exSecret, err = m.generateSwiftSecret(labels)
	case ossStorage:
		labels[LabelOfStorageType] = ossStorage
		exSecret, err = m.generateOssSecret(labels)
	default:
		return exSecret, fmt.Errorf(NotSupportType)
	}

	return exSecret, err
}

func (m *MinIOReconciler) generateS3Secret(labels map[string]string) (*corev1.Secret, error) {
	dataJson, err := json.Marshal(m.HarborCluster.Spec.Storage.S3)
	if err != nil {
		return nil, err
	}
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
	}, nil
}

func (m *MinIOReconciler) generateSecretForChartMuseum() (secret *corev1.Secret, err error) {
	if m.HarborCluster.Spec.ChartMuseum == nil {
		return secret, nil
	}
	labels := m.getLabels()
	switch m.HarborCluster.Spec.Storage.Kind {
	case s3Storage:
		labels[LabelOfStorageType] = s3Storage
		secret = m.generateS3SecretForChartMuseum(labels)
	default:
		return secret, fmt.Errorf(NotSupportType)
	}
	return secret, nil
}

func (m *MinIOReconciler) generateS3SecretForChartMuseum(labels map[string]string) *corev1.Secret {
	return &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Secret",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        m.getChartMuseumSecretName(),
			Namespace:   m.HarborCluster.Namespace,
			Labels:      labels,
			Annotations: m.generateAnnotations(),
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(m.HarborCluster, goharborv1.HarborClusterGVK),
			},
		},
		Type: corev1.SecretTypeOpaque,
		Data: map[string][]byte{
			"kind":                  []byte("amazon"),
			"AWS_ACCESS_KEY_ID":     []byte(m.HarborCluster.Spec.Storage.S3.AccessKey),
			"AWS_SECRET_ACCESS_KEY": []byte(m.HarborCluster.Spec.Storage.S3.SecretKey),
			"AMAZON_BUCKET":         []byte(m.HarborCluster.Spec.Storage.S3.Bucket),
			"AMAZON_PREFIX":         []byte(fmt.Sprintf("%s-subfloder", m.HarborCluster.Spec.Storage.S3.Bucket)),
			"AMAZON_REGION":         []byte(m.HarborCluster.Spec.Storage.S3.Region),
			"AMAZON_ENDPOINT":       []byte(m.HarborCluster.Spec.Storage.S3.RegionEndpoint),
		},
	}
}

func (m *MinIOReconciler) generateAzureSecret(labels map[string]string) (*corev1.Secret, error) {
	dataJson, err := json.Marshal(m.HarborCluster.Spec.Storage.Azure)
	if err != nil {
		return nil, err
	}
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
	}, nil
}

func (m *MinIOReconciler) generateGcsSecret(labels map[string]string) (*corev1.Secret, error) {
	dataJson, err := json.Marshal(m.HarborCluster.Spec.Storage.Gcs)
	if err != nil {
		return nil, err
	}
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
	}, nil
}

func (m *MinIOReconciler) generateSwiftSecret(labels map[string]string) (*corev1.Secret, error) {
	dataJson, err := json.Marshal(m.HarborCluster.Spec.Storage.Swift)
	if err != nil {
		return nil, err
	}
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
	}, nil
}

func (m *MinIOReconciler) generateOssSecret(labels map[string]string) (*corev1.Secret, error) {
	dataJson, err := json.Marshal(m.HarborCluster.Spec.Storage.Oss)
	if err != nil {
		return nil, err
	}
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
	}, nil
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

	var minioCR minio.Tenant
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

func (m *MinIOReconciler) generateMinIOCR() *minio.Tenant {
	return &minio.Tenant{
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
		Spec: minio.TenantSpec{
			Metadata: &metav1.ObjectMeta{
				Labels:      m.getLabels(),
				Annotations: m.generateAnnotations(),
			},
			ServiceName: m.getServiceName(),
			Image:       "minio/minio:" + m.HarborCluster.Spec.Storage.InCluster.Spec.Version,
			Zones: []minio.Zone{
				{
					Name:                DefaultZone,
					Servers:             m.HarborCluster.Spec.Storage.InCluster.Spec.Replicas,
					VolumesPerServer:    m.HarborCluster.Spec.Storage.InCluster.Spec.VolumesPerServer,
					VolumeClaimTemplate: m.getVolumeClaimTemplate(),
					Resources:           *m.getResourceRequirements(),
				},
			},
			Mountpath: minio.MinIOVolumeMountPath,
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
			Liveness: &minio.Liveness{
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
