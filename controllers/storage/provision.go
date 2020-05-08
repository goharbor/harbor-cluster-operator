package storage

import (
	goharborv1 "github.com/goharbor/harbor-cluster-operator/api/v1"
	"github.com/goharbor/harbor-cluster-operator/lcm"
	minio "github.com/minio/minio-operator/pkg/apis/miniooperator.min.io/v1beta1"
)

func (minio *MinIOReconciler) Provision(spec *goharborv1.HarborCluster) (*lcm.CRStatus, error) {

	panic("implement me")
}

func (minio *MinIOReconciler) generateMinIOCR(spec *goharborv1.HarborCluster) *minio.MinIOInstance {

	panic("implement me")
}