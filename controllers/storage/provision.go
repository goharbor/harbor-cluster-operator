package storage

import (
	goharborv1 "github.com/goharbor/harbor-cluster-operator/api/v1"
	"github.com/goharbor/harbor-cluster-operator/lcm"
	minio "github.com/minio/minio-operator/pkg/apis/miniooperator.min.io/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	Kind       = "MinIOInstance"
	ApiVersion = "miniooperator.min.io/v1beta1"
	DefaultZone = "zone-harbor"
)

func (m *MinIOReconciler) Provision() (*lcm.CRStatus, error) {

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
					Name: DefaultZone,
					Servers: m.HarborCluster.Spec.Stroage.InCluster.Spec.Replicas,
				},
			},
			VolumesPerServer: 1,
			Mountpath: "/export",
			VolumeClaimTemplate: m.getVolumeClaimTemplate(),


		},
	}
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

func (m *MinIOReconciler) generateHeadlessService() *minio.MinIOInstance {
	// TODO
	return nil
}

func (m *MinIOReconciler) generateMcsSecret() *minio.MinIOInstance {
	// TODO
	return nil
}

func (m *MinIOReconciler) generateCredsSecret() *minio.MinIOInstance {
	// TODO
	return nil
}

func (m *MinIOReconciler) generateService() *minio.MinIOInstance {
	// TODO
	return nil
}
