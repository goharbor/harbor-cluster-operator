package storage

import "github.com/goharbor/harbor-cluster-operator/lcm"

func (m *MinIOReconciler) Scale() (*lcm.CRStatus, error) {
	minioCR := m.CurrentMinIOCR
	minioCR.Spec.Zones[0].Servers = m.HarborCluster.Spec.Storage.InCluster.Spec.Replicas

	err := m.KubeClient.Update(minioCR)

	return minioUnknownStatus(), err
}
