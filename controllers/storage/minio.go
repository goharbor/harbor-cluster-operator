package storage

import (
	"context"
	"github.com/go-logr/logr"
	goharborv1 "github.com/goharbor/harbor-cluster-operator/api/v1"
	"github.com/goharbor/harbor-cluster-operator/controllers/k8s"
	"github.com/goharbor/harbor-cluster-operator/lcm"
	minio "github.com/minio/minio-operator/pkg/apis/operator.min.io/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	k8serror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
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
	HarborCluster  *goharborv1.HarborCluster
	KubeClient     k8s.Client
	Ctx            context.Context
	Log            logr.Logger
	Scheme         *runtime.Scheme
	Recorder       record.EventRecorder
	CurrentMinIOCR *minio.MinIOInstance
}

var (
	HarborClusterMinIOGVK = schema.GroupVersionKind{
		Group:   minio.SchemeGroupVersion.Group,
		Version: minio.SchemeGroupVersion.Version,
		Kind:    minio.MinIOCRDResourceKind,
	}
)

// Reconciler implements the reconcile logic of minIO service
func (m *MinIOReconciler) Reconcile() (*lcm.CRStatus, error) {
	var minioCR minio.MinIOInstance

	if m.HarborCluster.Spec.Storage.Kind != inClusterStorage {
		return m.ProvisionExternalStorage()
	}

	err := m.KubeClient.Get(m.getMinIONamespacedName(), &minioCR)
	if k8serror.IsNotFound(err) {
		return m.Provision()
	} else if err != nil {
		return minioNotReadyStatus(GetMinIOError, err.Error()), err
	}

	m.CurrentMinIOCR = &minioCR

	isScale := m.checkMinIOScale()
	if isScale {
		return m.Scale()
	}

	isUpdate := m.checkMinIOUpdate()
	if isUpdate {
		return m.Update()
	}

	isReady, err := m.checkMinIOReady()
	if err != nil {
		return minioNotReadyStatus(GetMinIOError, err.Error()), err
	}

	if isReady {
		err := createDefaultBucket()
		if err != nil {
			return minioNotReadyStatus(CreateDefaultBucketeError, err.Error()), err
		}
		return m.ProvisionInClusterSecretAsOss(&minioCR)
	}

	return nil, nil
}

func createDefaultBucket() error {
	panic("implement me")
}

func (m *MinIOReconciler) checkMinIOUpdate() bool {
	panic("implement me")
}

func (m *MinIOReconciler) checkMinIOScale() bool {
	panic("implement me")
}

func (m *MinIOReconciler) checkMinIOReady() (bool, error) {
	var minioStatefulSet appsv1.StatefulSet
	err := m.KubeClient.Get(m.getMinIONamespacedName(), &minioStatefulSet)

	if minioStatefulSet.Status.ReadyReplicas == m.HarborCluster.Spec.Storage.InCluster.Spec.Replicas {
		return true, err
	}

	return false, err
}

func (m *MinIOReconciler) getMinIONamespacedName() types.NamespacedName {
	return types.NamespacedName{
		Namespace: m.HarborCluster.Namespace,
		Name:      m.getServiceName(),
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

// TODO Deprecated
func (m *MinIOReconciler) ScaleUp(newReplicas uint64) (*lcm.CRStatus, error) {
	panic("implement me")
}

// TODO Deprecated
func (m *MinIOReconciler) ScaleDown(newReplicas uint64) (*lcm.CRStatus, error) {
	panic("implement me")
}
