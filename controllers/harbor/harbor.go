package harbor

import (
	goharborv1 "github.com/goharbor/harbor-cluster-operator/api/v1"
	"github.com/goharbor/harbor-cluster-operator/lcm"
)

type HarborReconciler struct {
	HarborCluster       *goharborv1.HarborCluster
	ComponentToCRStatus map[goharborv1.Component]*lcm.CRStatus
}

// Reconciler implements the reconcile logic of minIO service
func (harbor *HarborReconciler) Reconcile() (*lcm.CRStatus, error) {
	// TODO
	return nil, nil
}

func (harbor *HarborReconciler) Provision() (*lcm.CRStatus, error) {
	panic("implement me")
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
