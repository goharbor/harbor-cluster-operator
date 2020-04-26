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

// TODO
type CRStatus struct {
}
