package storage

import (
	goharborv1 "github.com/goharbor/harbor-cluster-operator/api/v1"
	"github.com/goharbor/harbor-cluster-operator/lcm"
	minio "github.com/minio/minio-operator/pkg/apis/miniooperator.min.io/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	virtualServiceGVR = schema.GroupVersionResource{
		Group:    "databases.spotahome.com",
		Version:  "v1",
		Resource: "redisfailovers",
	}
)

func (m *MinIOReconciler) Provision() (*lcm.CRStatus, error) {

	panic("implement me")
}

func (m *MinIOReconciler) generateMinIOCR() *minio.MinIOInstance {
	return &minio.MinIOInstance{
		ObjectMeta: metav1.ObjectMeta{
			Name: m.HarborCluster.Name,
			Namespace: m.HarborCluster.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(m.HarborCluster, goharborv1.HarborClusterGVK),
			},
		},
		Spec: minio.MinIOInstanceSpec{
			
		},
	}
}