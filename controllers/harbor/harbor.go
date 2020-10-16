package harbor

import (
	"context"
	"fmt"
	goharborv1 "github.com/goharbor/harbor-cluster-operator/apis/goharbor.io/v1alpha2"
	"github.com/goharbor/harbor-cluster-operator/controllers/image"
	"github.com/goharbor/harbor-cluster-operator/controllers/k8s"
	"github.com/goharbor/harbor-cluster-operator/lcm"
	"github.com/goharbor/harbor-operator/api/v1alpha1"
	"github.com/google/go-cmp/cmp"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
)

const (
	ScalingEvent  = "Scaling"
	UpdatingEvent = "Updating"
)

type HarborReconciler struct {
	k8s.Client
	Ctx                 context.Context
	HarborCluster       *goharborv1.HarborCluster
	CurrentHarborCR     *v1alpha1.Harbor
	DesiredHarborCR     *v1alpha1.Harbor
	ImageGetter         image.Getter
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
			return harborClusterCRNotReadyStatus(GetHarborCRError, err.Error()), err
		}
	}

	harbor.CurrentHarborCR = &harborCR
	harbor.DesiredHarborCR = harbor.newHarborCR()

	event := harbor.checkReconcileEvent(harbor.HarborCluster, &harborCR)
	switch event {
	case ScalingEvent:
		return harbor.Scale()
	case UpdatingEvent:
		return harbor.Update(harbor.HarborCluster)
	}

	err = harbor.Get(harbor.getHarborCRNamespacedName(), &harborCR)
	if err != nil {
		return harborClusterCRUnknownStatus(GetHarborCRError, err.Error()), err
	}
	return harborClusterCRStatus(&harborCR), nil
}

// unsetReplicas will set replicas to nil for all components in v1alpha2.Harbor.
// This is a helper method to check whether harbor cr is equal expect replicas.
func unsetReplicas(harbor *v1alpha1.Harbor) {
	if harbor.Spec.Components.Core != nil {
		harbor.Spec.Components.Core.Replicas = nil
	}

	if harbor.Spec.Components.Portal != nil {
		harbor.Spec.Components.Portal.Replicas = nil
	}

	if harbor.Spec.Components.Registry != nil {
		harbor.Spec.Components.Registry.Replicas = nil
	}

	if harbor.Spec.Components.Clair != nil {
		harbor.Spec.Components.Clair.Replicas = nil
	}

	if harbor.Spec.Components.ChartMuseum != nil {
		harbor.Spec.Components.ChartMuseum.Replicas = nil
	}

	if harbor.Spec.Components.Notary != nil {
		harbor.Spec.Components.Notary.Server.Replicas = nil
		harbor.Spec.Components.Notary.Signer.Replicas = nil
	}

	if harbor.Spec.Components.JobService != nil {
		harbor.Spec.Components.JobService.Replicas = nil
	}
}

// isEqualExpectReplicas check whether harbor cr is equal expect replicas.
func isEqualExpectReplicas(desiredHarborCR *v1alpha1.Harbor, currentHarborCR *v1alpha1.Harbor) bool {
	desiredHarborCRCopied := desiredHarborCR.DeepCopy()
	currentHarborCRCopied := currentHarborCR.DeepCopy()

	unsetReplicas(desiredHarborCRCopied)
	unsetReplicas(currentHarborCRCopied)

	return cmp.Equal(desiredHarborCRCopied.Spec, currentHarborCRCopied.Spec)
}

func (harbor *HarborReconciler) checkReconcileEvent(desired *goharborv1.HarborCluster, current *v1alpha1.Harbor) string {
	isEqualExpectReplicas := isEqualExpectReplicas(harbor.DesiredHarborCR, current)
	if !isEqualExpectReplicas {
		return UpdatingEvent
	}
	if harbor.isScalingEvent(desired, current) {
		return ScalingEvent
	}
	return ""
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

func harborClusterCRNotReadyStatus(reason, message string) *lcm.CRStatus {
	return lcm.New(goharborv1.ServiceReady).WithStatus(corev1.ConditionFalse).WithReason(reason).WithMessage(message)
}

func harborClusterCRUnknownStatus(reason, message string) *lcm.CRStatus {
	return lcm.New(goharborv1.ServiceReady).WithStatus(corev1.ConditionUnknown).WithReason(reason).WithMessage(message)
}

// harborClusterCRStatus will assembly the harbor cluster status according the v1alpha2.Harbor status
func harborClusterCRStatus(harbor *v1alpha1.Harbor) *lcm.CRStatus {
	for _, condition := range harbor.Status.Conditions {
		if condition.Type == v1alpha1.ReadyConditionType {
			return lcm.New(goharborv1.ServiceReady).WithStatus(condition.Status).WithMessage(condition.Message).WithReason(condition.Reason)
		}
	}
	return harborClusterCRUnknownStatus(EmptyHarborCRStatusError, "The ready condition of harbor.goharbor.io is empty. Please wait for minutes.")
}

func (harbor *HarborReconciler) getHarborCRNamespacedName() types.NamespacedName {
	return types.NamespacedName{
		Namespace: harbor.HarborCluster.Namespace,
		Name:      fmt.Sprintf("%s-harbor", harbor.HarborCluster.Name),
	}
}
