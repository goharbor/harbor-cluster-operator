package storage

import (
	"github.com/goharbor/harbor-cluster-operator/lcm"
)

func (m *MinIOReconciler) Update() (*lcm.CRStatus, error) {
	m.CurrentMinIOCR.Spec = m.DesiredMinIOCR.Spec
	err := m.KubeClient.Update(m.CurrentMinIOCR)
	if err != nil {
		return minioNotReadyStatus(UpdateMinIOError, err.Error()), err
	}

	return minioUnknownStatus(), nil
}

func (m *MinIOReconciler) ExternalUpdate() (*lcm.CRStatus, error) {
	currentSecret := m.CurrentExternalSecret
	currentSecret.Labels = m.DesiredExternalSecret.Labels
	currentSecret.Data = m.DesiredExternalSecret.Data

	err := m.KubeClient.Update(currentSecret)
	if err != nil {
		return minioNotReadyStatus(UpdateExternalSecretError, err.Error()), err
	}

	p := &lcm.Property{
		Name:  m.HarborCluster.Spec.Storage.Kind + ExternalStorageSecretSuffix,
		Value: m.getExternalSecretName(),
	}
	properties := &lcm.Properties{p}

	return minioReadyStatus(properties), nil
}
