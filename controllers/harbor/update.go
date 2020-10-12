package harbor

import (
	goharborv1 "github.com/goharbor/harbor-cluster-operator/apis/goharbor.io/v1alpha1"
	"github.com/goharbor/harbor-cluster-operator/lcm"
)

func (harbor *HarborReconciler) Update(spec *goharborv1.HarborCluster) (*lcm.CRStatus, error) {
	desiredHarborCR := harbor.newHarborCR()
	err := harbor.Client.Update(desiredHarborCR)
	if err != nil {
		return harborClusterCRUnknownStatus(UpdateHarborCRError, err.Error()), err
	}
	return harborClusterCRStatus(desiredHarborCR), nil
}
