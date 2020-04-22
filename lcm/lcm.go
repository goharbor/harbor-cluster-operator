package lcm

import v1 "github.com/goharbor/harbor-cluster-operator/api/v1"

// This package container interface of harbor cluster service lifecycle manage.

type Controller interface {
	// Provision the new dependent service by the related section of cluster spec.
	provision(spec *v1.HarborCluster) (*CRStatus, error)

	// Delete the service
	delete() (*CRStatus, error)

	// Scale up
	scaleUp(newReplicas uint64) (*CRStatus, error)

	// Scale down
	scaleDown(newReplicas uint64) (*CRStatus, error)

	// Update the service
	update(spec *v1.HarborCluster) (*CRStatus, error)

	// More...
}

// TODO
type CRStatus struct {
}
