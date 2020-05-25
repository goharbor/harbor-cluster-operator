package harbor

import (
	"context"
	"fmt"
	goharborv1 "github.com/goharbor/harbor-cluster-operator/api/v1"
	"github.com/goharbor/harbor-cluster-operator/controllers/k8s"
	"github.com/goharbor/harbor-cluster-operator/lcm"
	"github.com/goharbor/harbor-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

type HarborReconciler struct {
	k8s.Client
	Ctx                 context.Context
	HarborCluster       *goharborv1.HarborCluster
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
			// TODO given clear reason and message.
			return harborCRNotReadyStatus("", ""), err
		}
	}
	return nil, nil
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

func harborCRNotReadyStatus(reason, message string) *lcm.CRStatus {
	return &lcm.CRStatus{
		Condition: goharborv1.HarborClusterCondition{
			Type:               goharborv1.ServiceReady,
			Status:             corev1.ConditionFalse,
			LastTransitionTime: metav1.Now(),
			Reason:             reason,
			Message:            message,
		},
		Properties: nil,
	}
}

func harborCRUnknownStatus() *lcm.CRStatus {
	return &lcm.CRStatus{
		Condition: goharborv1.HarborClusterCondition{
			Type:               goharborv1.ServiceReady,
			Status:             corev1.ConditionUnknown,
			LastTransitionTime: metav1.Now(),
			Reason:             "",
			Message:            "",
		},
		Properties: nil,
	}
}

func (harbor *HarborReconciler) getHarborCRNamespacedName() types.NamespacedName {
	return types.NamespacedName{
		Namespace: harbor.HarborCluster.Namespace,
		Name:      fmt.Sprintf("%s-harbor", harbor.HarborCluster.Name),
	}
}
