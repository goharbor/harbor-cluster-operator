package harbor

import (
	"context"
	"fmt"
	goharborv1 "github.com/goharbor/harbor-cluster-operator/api/v1"
	"github.com/goharbor/harbor-cluster-operator/controllers/image"
	"github.com/goharbor/harbor-cluster-operator/controllers/k8s"
	"github.com/goharbor/harbor-cluster-operator/lcm"
	"github.com/goharbor/harbor-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
)

type HarborReconciler struct {
	k8s.Client
	Ctx                 context.Context
	HarborCluster       *goharborv1.HarborCluster
	CurrentHarborCR     *v1alpha1.Harbor
	ImageGetter         image.ImageGetter
	ComponentToCRStatus map[goharborv1.Component]*lcm.CRStatus
}

// Reconciler implements the reconcile logic of services
func (harbor *HarborReconciler) Reconcile() (*lcm.CRStatus, error) {
	var harborCR v1alpha1.Harbor
	err := harbor.Get(harbor.getHarborCRNamespacedName(), &harborCR)
	if err != nil {
		if errors.IsNotFound(err) {
			return harbor.Provision()
		} else {
			return harborClusterCRUnknownStatus(), err
		}
	}
	harbor.CurrentHarborCR = &harborCR
	isScalingEvent := harbor.isScalingEvent(harbor.HarborCluster, &harborCR)
	if isScalingEvent {
		return harbor.Scale()
	}

	isUpdatingEvent := harbor.isUpdatingEvent(harbor.HarborCluster, &harborCR)
	if isUpdatingEvent {
		return harbor.Update(harbor.HarborCluster)
	}

	err = harbor.Get(harbor.getHarborCRNamespacedName(), &harborCR)
	if err != nil {
		return harborClusterCRUnknownStatus(), err
	}
	return harborClusterCRStatus(&harborCR), nil
}

func (harbor *HarborReconciler) Delete() (*lcm.CRStatus, error) {
	panic("implement me")
}

func (harbor *HarborReconciler) ScaleUp(newReplicas uint64) (*lcm.CRStatus, error) {
	panic("implement me")
}

func (harbor *HarborReconciler) ScaleDown(newReplicas uint64) (*lcm.CRStatus, error) {
	panic("implement me")
}

func (harbor *HarborReconciler) Update(spec *goharborv1.HarborCluster) (*lcm.CRStatus, error) {
	panic("implement me")
}

func harborClusterCRNotReadyStatus(reason, message string) *lcm.CRStatus {
	return lcm.New(goharborv1.ServiceReady).WithStatus(corev1.ConditionFalse).WithReason(reason).WithMessage(message)
}

func harborClusterCRUnknownStatus() *lcm.CRStatus {
	return lcm.New(goharborv1.ServiceReady).WithStatus(corev1.ConditionUnknown)
}

func harborClusterCRReadyStatus() *lcm.CRStatus {
	return lcm.New(goharborv1.ServiceReady).WithStatus(corev1.ConditionTrue)
}

func harborClusterCRStatus(harbor *v1alpha1.Harbor) *lcm.CRStatus {
	for _, condition := range harbor.Status.Conditions {
		if condition.Type == v1alpha1.ReadyConditionType {
			return lcm.New(goharborv1.ServiceReady).WithStatus(condition.Status).WithMessage(condition.Message).WithReason(condition.Reason)
		}
	}
	return harborClusterCRUnknownStatus()
}

func (harbor *HarborReconciler) getHarborCRNamespacedName() types.NamespacedName {
	return types.NamespacedName{
		Namespace: harbor.HarborCluster.Namespace,
		Name:      fmt.Sprintf("%s-harbor", harbor.HarborCluster.Name),
	}
}
