package harbor

import (
	goharborv1 "github.com/goharbor/harbor-cluster-operator/api/v1"
	"github.com/goharbor/harbor-operator/api/v1alpha1"
)

func (harbor *HarborReconciler) isUpdatingEvent(desired *goharborv1.HarborCluster, current *v1alpha1.Harbor) bool {

}
