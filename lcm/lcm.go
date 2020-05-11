package lcm

import v1 "github.com/goharbor/harbor-cluster-operator/api/v1"

// This package container interface of harbor cluster service lifecycle manage.

type Controller interface {
	// Provision the new dependent service by the related section of cluster spec.
	Provision(spec *v1.HarborCluster) (*CRStatus, error)

	// Delete the service
	Delete() (*CRStatus, error)

	// Scale up
	ScaleUp(newReplicas uint64) (*CRStatus, error)

	// Scale down
	ScaleDown(newReplicas uint64) (*CRStatus, error)

	// Update the service
	Update(spec *v1.HarborCluster) (*CRStatus, error)

	// More...
}

type CRStatus struct {
	Condition  v1.HarborClusterCondition `json:"condition"`
	Properties Properties                `json:"properties"`
}

//Phase is the current state of component
type Phase string

const (
	//PendingPhase represents pending state of component.
	//The reconcile has created CR, but the CR has not been scheduled yet.
	PendingPhase Phase = "Pending"
	//CreatingPhase represents creating state of component.
	//The reconcile is creating CR, but the service is not available.
	CreatingPhase Phase = "Creating"
	//ReadyPhase represents ready state of component.
	//The CR has been deployed, the number of nodes is as expected, and the service is available.
	ReadyPhase Phase = "Ready"
	//UpgradingPhase represents upgrade state of component.
	//The CR is processing upscale、downscale、rolling update.
	UpgradingPhase Phase = "Upgrading"
	//DestroyingPhase represents delete state of component.
	DestroyingPhase Phase = "Destroying"
)
