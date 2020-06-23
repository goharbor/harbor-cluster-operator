package harbor

import (
	goharborv1 "github.com/goharbor/harbor-cluster-operator/api/v1"
	"github.com/goharbor/harbor-cluster-operator/lcm"
	"github.com/goharbor/harbor-operator/api/v1alpha1"
)

// ScaleUp will update replicas of all components, expect job service.
func (harbor *HarborReconciler) Scale() (*lcm.CRStatus, error) {
	current := harbor.CurrentHarborCR
	desiredReplicas := int32(harbor.HarborCluster.Spec.Replicas)
	if current.Spec.Components.Core != nil {
		current.Spec.Components.Core.Replicas = &desiredReplicas
	}

	if current.Spec.Components.Portal != nil {
		current.Spec.Components.Portal.Replicas = &desiredReplicas
	}

	if current.Spec.Components.Registry != nil {
		current.Spec.Components.Registry.Replicas = &desiredReplicas
	}

	if current.Spec.Components.Clair != nil {
		current.Spec.Components.Clair.Replicas = &desiredReplicas
	}

	if current.Spec.Components.ChartMuseum != nil {
		current.Spec.Components.ChartMuseum.Replicas = &desiredReplicas
	}

	if current.Spec.Components.Notary != nil {
		current.Spec.Components.Notary.Server.Replicas = &desiredReplicas
		current.Spec.Components.Notary.Signer.Replicas = &desiredReplicas
	}

	if harbor.HarborCluster.Spec.JobService != nil {
		desiredJobReplicas := int32(harbor.HarborCluster.Spec.JobService.Replicas)
		if current.Spec.Components.JobService != nil {
			current.Spec.Components.JobService.Replicas = &desiredJobReplicas
		}
	}

	err := harbor.Client.Update(current)
	if err != nil {
		return harborClusterCRUnknownStatus(), err
	}
	// TODO declare the detail CRStatus
	return harborClusterCRReadyStatus(), nil
}

// isScalingEvent will compare the actual replicas of any components with the desired replicas.
// return true if the actual replicas is not equal to desired replicas of any component.
// return false if the actual replicas of all components are equal to desired replicas.
func (harbor *HarborReconciler) isScalingEvent(desired *goharborv1.HarborCluster, current *v1alpha1.Harbor) bool {
	desiredReplicas := int32(desired.Spec.Replicas)

	if current.Spec.Components.Core != nil && current.Spec.Components.Core.Replicas != nil {
		coreReplicas := current.Spec.Components.Core.Replicas
		if desiredReplicas != *coreReplicas {
			return true
		}
	}

	if current.Spec.Components.Portal != nil && current.Spec.Components.Portal.Replicas != nil {
		portalReplicas := current.Spec.Components.Portal.Replicas
		if desiredReplicas != *portalReplicas {
			return true
		}
	}

	if current.Spec.Components.Registry != nil && current.Spec.Components.Registry.Replicas != nil {
		registryReplicas := current.Spec.Components.Registry.Replicas
		if desiredReplicas != *registryReplicas {
			return true
		}
	}

	if current.Spec.Components.Clair != nil && current.Spec.Components.Clair.Replicas != nil {
		clairReplicas := current.Spec.Components.Registry.Replicas
		if desiredReplicas != *clairReplicas {
			return true
		}
	}

	if current.Spec.Components.ChartMuseum != nil && current.Spec.Components.ChartMuseum.Replicas != nil {
		chartMuseumReplicas := current.Spec.Components.ChartMuseum.Replicas
		if desiredReplicas != *chartMuseumReplicas {
			return true
		}
	}

	if current.Spec.Components.Notary != nil {
		if current.Spec.Components.Notary.Server.Replicas != nil {
			notaryServerReplicas := current.Spec.Components.Notary.Server.Replicas
			if desiredReplicas != *notaryServerReplicas {
				return true
			}
		}
		if current.Spec.Components.Notary.Signer.Replicas != nil {
			notarySignerReplicas := current.Spec.Components.Notary.Signer.Replicas
			if desiredReplicas != *notarySignerReplicas {
				return true
			}
		}
	}

	desiredJobServiceReplicas := int32(desired.Spec.JobService.Replicas)
	if current.Spec.Components.JobService != nil && current.Spec.Components.JobService.Replicas != nil {
		jobServiceReplicas := current.Spec.Components.JobService.Replicas
		if desiredJobServiceReplicas != *jobServiceReplicas {
			return true
		}
	}

	return false
}
