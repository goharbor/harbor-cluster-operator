package storage

import (
	"context"
	"github.com/go-logr/logr"
	goharborv1 "github.com/goharbor/harbor-cluster-operator/api/v1"
	"github.com/goharbor/harbor-cluster-operator/lcm"
	minioClient "github.com/minio/minio-operator/pkg/client/clientset/versioned"
	k8serror "k8s.io/apimachinery/pkg/api/errors"
	apisv1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type MinIOReconciler struct {
	HarborCluster *goharborv1.HarborCluster

	// kubeClientSet is a standard kubernetes clientset
	kubeClientSet kubernetes.Interface
	// minioClientSet is a clientset for our own API group
	minioClientSet minioClient.Interface
	//// minioInstancesLister lists MinIOInstance from a shared informer's
	//// store.
	//minioInstancesLister miniolisters.MinIOInstanceLister
	//// minioInstancesSynced returns true if the StatefulSet shared informer
	//// has synced at least once.
	//minioInstancesSynced cache.InformerSynced

	ctx context.Context

	Log logr.Logger

	Namespace string

}

// Reconciler implements the reconcile logic of minIO service
func (minio *MinIOReconciler) Reconcile() (*lcm.CRStatus, error) {
	_,err :=minio.minioClientSet.MiniooperatorV1beta1().MinIOInstances(minio.Namespace).Get(minio.ctx,minio.HarborCluster.Name,apisv1.GetOptions{})
	if k8serror.IsNotFound(err) {
		return minio.Provision(minio.HarborCluster)
	} else if err != nil {
		return nil, err
	}

	return nil, nil
}

func (minio *MinIOReconciler) Delete() (*lcm.CRStatus, error) {
	panic("implement me")
}

func (minio *MinIOReconciler) ScaleUp(newReplicas uint64) (*lcm.CRStatus, error) {
	panic("implement me")
}

func (minio *MinIOReconciler) ScaleDown(newReplicas uint64) (*lcm.CRStatus, error) {
	panic("implement me")
}

func (minio *MinIOReconciler) Update(spec *goharborv1.HarborCluster) (*lcm.CRStatus, error) {
	panic("implement me")
}

func (minio *MinIOReconciler) generateMinIO(spec *goharborv1.HarborCluster) (*lcm.CRStatus, error) {
	panic("implement me")
}
