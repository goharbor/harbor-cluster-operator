package storage

import (
	"context"
	"fmt"
	"github.com/go-logr/logr"
	goharborv1 "github.com/goharbor/harbor-cluster-operator/api/v1"
	"github.com/goharbor/harbor-cluster-operator/controllers/k8s"
	"github.com/goharbor/harbor-cluster-operator/lcm"
	minio "github.com/minio/minio-operator/pkg/apis/operator.min.io/v1"
	corev1 "k8s.io/api/core/v1"
	k8serror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
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
	KubeClient    k8s.Client
	Ctx           context.Context
	Log           logr.Logger
	Scheme        *runtime.Scheme
	Recorder      record.EventRecorder
}

// Reconciler implements the reconcile logic of minIO service
func (m *MinIOReconciler) Reconcile() (*lcm.CRStatus, error) {
	var minioCR minio.MinIOInstance

	if m.HarborCluster.Spec.Storage.Kind != inClusterStorage {
		return m.ProvisionExternalStorage()
	}

	err := m.KubeClient.Get(m.getminIONamespacedName(), &minioCR)
	if k8serror.IsNotFound(err) {
		// TODO need test
		return m.Provision()
	} else if err != nil {
		return minioNotReadyStatus(ErrorReason0, err.Error()), err
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

func minioReadyStatus(properties *lcm.Properties) *lcm.CRStatus {
	return &lcm.CRStatus{
		Condition: goharborv1.HarborClusterCondition{
			Type:               goharborv1.StorageReady,
			Status:             corev1.ConditionTrue,
			LastTransitionTime: metav1.Now(),
			Reason:             "",
			Message:            "",
		},
		Properties: *properties,
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
