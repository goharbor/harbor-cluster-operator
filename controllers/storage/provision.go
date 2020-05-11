package storage

import (
	goharborv1 "github.com/goharbor/harbor-cluster-operator/api/v1"
	"github.com/goharbor/harbor-cluster-operator/lcm"
	minio "github.com/minio/minio-operator/pkg/apis/miniooperator.min.io/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	Kind               = "MinIOInstance"
	ApiVersion         = "miniooperator.min.io/v1beta1"
	DefaultZone        = "zone-harbor"
	DefaultCredsSecret = "minio-creds-secret"
	DefaultMcsSecret   = "minio-mcs-secret"
	CredsAccesskey     = "bWluaW8="
	CredsSecretkey     = "bWluaW8xMjM="
)

func (m *MinIOReconciler) Provision() (*lcm.CRStatus, error) {
	mcsSecret := m.generateMcsSecret()
	credsSecret := m.generateCredsSecret()

	panic("implement me")
}

func (m *MinIOReconciler) generateMinIOCR() *minio.MinIOInstance {
	return &minio.MinIOInstance{
		TypeMeta: metav1.TypeMeta{
			Kind:       Kind,
			APIVersion: ApiVersion,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.HarborCluster.Name,
			Namespace: m.HarborCluster.Namespace,
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
			Image: "minio/minio:" + m.HarborCluster.Spec.Stroage.InCluster.Spec.Version,
			Zones: []minio.Zone{
				minio.Zone{
					Name:    m.HarborCluster.Name+"-"+DefaultZone,
					Servers: m.HarborCluster.Spec.Stroage.InCluster.Spec.Replicas,
				},
			},
			VolumesPerServer:    1,
			Mountpath:           "/export",
			VolumeClaimTemplate: m.getVolumeClaimTemplate(),
			CredsSecret: &corev1.LocalObjectReference{
				Name: m.HarborCluster.Name + "-" + DefaultCredsSecret,
			},
			PodManagementPolicy: "Parallel",
			RequestAutoCert: false,
			CertConfig: &minio.CertificateConfig{
				CommonName: "",
				OrganizationName: []string{},
				DNSNames: []string{},
			},
			Env: []corev1.EnvVar{
				corev1.EnvVar{
					Name: "MINIO_BROWSER",
					Value: "on",
				},
			},
			Resources: *m.getResourceRequirements(),//m.HarborCluster.Spec.Stroage.InCluster.Spec.Resources,
			Liveness: &corev1.Probe{
				Handler: corev1.Handler{
					HTTPGet: &corev1.HTTPGetAction{
						Path: "/minio/health/live",
						Port: intstr.IntOrString{
							IntVal: 9000,
						},

					},
				},
				InitialDelaySeconds: 120,
				PeriodSeconds: 60,
			},
			Readiness: &corev1.Probe{
				Handler: corev1.Handler{
					HTTPGet: &corev1.HTTPGetAction{
						Path: "/minio/health/ready",
						Port: intstr.IntOrString{
							IntVal: 9000,
						},

					},
				},
				InitialDelaySeconds: 120,
				PeriodSeconds: 60,
			},
		},
	}
}

func (m *MinIOReconciler) getResourceRequirements() *corev1.ResourceRequirements {
	// TODO
	return nil
}

func (m *MinIOReconciler) getVolumeClaimTemplate() *corev1.PersistentVolumeClaim {
	// TODO
	return nil
}

// TODO
func (m *MinIOReconciler) getLabels() map[string]string {
	// TODO
	return nil
}

func (m *MinIOReconciler) generateAnnotations() map[string]string {
	// TODO
	return nil
}

func (m *MinIOReconciler) generateService() *minio.MinIOInstance {
	// TODO
	return nil
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
		},
		Type: corev1.SecretTypeOpaque,
		// TODO
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
