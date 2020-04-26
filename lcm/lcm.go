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
	Phase      Phase       `json:"phase"`
	Properties []*Property `json:"properties"`
}

//Phase is the current state of component
type Phase string

const (
	//PendingPhase represents pending state of component.
	//The reconcile has created CR, but the CR has not been scheduled yet.
	PendingPhase Phase = "Pending"
	//DeployingPhase represents deploying state of component.
	//The reconcile is deploying CR, but the service is not available.
	DeployingPhase Phase = "Deploying"
	//ReadyPhase represents ready state of component.
	//The CR has been deployed, the number of nodes is as expected, and the service is available.
	ReadyPhase Phase = "Ready"
	//UpgradingPhase represents upgrade state of component.
	//The CR is processing upscale、downscale、rolling update.
	UpgradingPhase Phase = "Upgrading"
	//DestroyingPhase represents delete state of component.
	DestroyingPhase Phase = "Destroying"
)

const (
	//ProperConn represents the connection info of the component.
	ProperConn = "Connection"
	//ProperPort represents the connection port of the component.
	ProperPort = "Port"
	//ProperUser represents the connection user of the component.
	ProperUser = "Username"
	//ProperPass represents the connection password of the component.
	ProperPass = "Password"
	//ProperNodes represent the available nodes of the component.
	ProperNodes = "AvailableNodes"
)

//Property is the current property of component.
type Property struct {
	//Property name, e.p: Connection,Port.
	Name string
	//Property type, e.p: int, string, secret.
	Type string
	//Property value, e.p: "rfs-harborcluster-sample.svc"
	Value interface{}
}

//NewProperty create an Property
func NewProperty(name string, typ string, value interface{}) *Property {
	return &Property{
		Name:  name,
		Type:  typ,
		Value: value,
	}
}
