package storage

import (
	"context"
	"fmt"
	"github.com/go-logr/logr"
	goharborv1 "github.com/goharbor/harbor-cluster-operator/api/v1"
	"github.com/goharbor/harbor-cluster-operator/lcm"
	minio "github.com/minio/minio-operator/pkg/apis/operator.min.io/v1"
	corev1 "k8s.io/api/core/v1"
	k8serror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	inClusterStorage = "inCluster"
	azureStorage     = "azure"
	gcsStorage       = "gcs"
	s3Storage        = "s3"
	swiftStorage     = "swift"
	ossStorage       = "oss"
)

type MinIOReconciler struct {
	HarborCluster *goharborv1.HarborCluster

	KubeClient client.Client

	Ctx context.Context

	Log logr.Logger
}

// Reconciler implements the reconcile logic of minIO service
func (m *MinIOReconciler) Reconcile() (*lcm.CRStatus, error) {
	// TODO
	var minioCR minio.MinIOInstance

	switch m.HarborCluster.Spec.Storage.Kind {
	case inClusterStorage:
		// TODO external storage service
		err := m.KubeClient.Get(m.Ctx, m.getminIONamespacedName(), &minioCR)
		if k8serror.IsNotFound(err) {
			return m.Provision()
		} else if err != nil {
			return minioNotReadyStatus(ErrorReason0, err.Error()), err
		}
	case azureStorage:
		minioStatus, err := m.ProvisionAzure()
		return minioStatus, err
	case gcsStorage:
		minioStatus, err := m.ProvisionGcs()
		return minioStatus, err
	case s3Storage:
		err := m.ProvisionS3()
		if err != nil {
			minioNotReadyStatus(ErrorReason1, err.Error())
		}
		return minioReadyStatus(), err
	case swiftStorage:
		minioStatus, err := m.ProvisionSwift()
		return minioStatus, err
	case ossStorage:
		minioStatus, err := m.ProvisionOss()
		return minioStatus, err
	default:
		// TODO
		return nil, nil
	}
	return nil, nil
}

func (m *MinIOReconciler) getminIONamespacedName() types.NamespacedName {
	return types.NamespacedName{
		Namespace: m.HarborCluster.Namespace,
		Name:      fmt.Sprintf("%s-minio", m.HarborCluster.Name),
	}
}

func minioNotReadyStatus(reason, message string) *lcm.CRStatus {
	return &lcm.CRStatus{
		Condition: goharborv1.HarborClusterCondition{
			Type:               goharborv1.StorageReady,
			Status:             corev1.ConditionFalse,
			LastTransitionTime: metav1.Now(),
			Reason:             reason,
			Message:            message,
		},
		Properties: nil,
	}
}

func minioUnknownStatus() *lcm.CRStatus {
	return &lcm.CRStatus{
		Condition: goharborv1.HarborClusterCondition{
			Type:               goharborv1.StorageReady,
			Status:             corev1.ConditionUnknown,
			LastTransitionTime: metav1.Now(),
			Reason:             "",
			Message:            "",
		},
		Properties: nil,
	}
}

func minioReadyStatus() *lcm.CRStatus {
	return &lcm.CRStatus{
		Condition: goharborv1.HarborClusterCondition{
			Type:               goharborv1.StorageReady,
			Status:             corev1.ConditionTrue,
			LastTransitionTime: metav1.Now(),
			Reason:             "",
			Message:            "",
		},
		Properties: nil,
	}
}

func (m *MinIOReconciler) Delete() (*lcm.CRStatus, error) {
	panic("implement me")
}

func (m *MinIOReconciler) ScaleUp(newReplicas uint64) (*lcm.CRStatus, error) {
	panic("implement me")
}

func (m *MinIOReconciler) ScaleDown(newReplicas uint64) (*lcm.CRStatus, error) {
	panic("implement me")
}

func (m *MinIOReconciler) Update(spec *goharborv1.HarborCluster) (*lcm.CRStatus, error) {
	panic("implement me")
}

func (minio *MinIOReconciler) generateMinIO(spec *goharborv1.HarborCluster) (*lcm.CRStatus, error) {
	panic("implement me")
}
