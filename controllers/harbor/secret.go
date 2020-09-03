package harbor

import (
	"fmt"
	"github.com/goharbor/harbor-cluster-operator/controllers/common"
	"github.com/goharbor/harbor-cluster-operator/controllers/k8s"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

const (
	_adminPasswordSecretNameTemplate = "%s-admin-password-secret"
	_adminPasswordKey                = "password"
)

// CheckAdminPasswordSecret will check whether .spec.adminPassword is empty in HarborCluster.
// If empty, auto generate a secret.
func (harbor *HarborReconciler) CheckAdminPasswordSecret() error {
	// if not empty, check successfully
	if harbor.HarborCluster.Spec.AdminPasswordSecret != "" {
		return nil
	}

	// auto generate a secret for admin password.
	adminPasswordSecretName := fmt.Sprintf(_adminPasswordSecretNameTemplate, harbor.HarborCluster.Name)
	adminPasswordSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      adminPasswordSecretName,
			Namespace: harbor.HarborCluster.Namespace,
			Labels: map[string]string{
				k8s.HarborClusterNameLabel: harbor.HarborCluster.Name,
			},
		},
		StringData: map[string]string{
			_adminPasswordKey: common.RandomString(16, common.UpperStringRandomType),
		},
	}

	err := harbor.Client.Create(adminPasswordSecret)

	if errors.IsAlreadyExists(err) {
		err = harbor.Client.Get(types.NamespacedName{
			Namespace: harbor.HarborCluster.Namespace,
			Name:      adminPasswordSecretName,
		}, adminPasswordSecret)
	}

	if err != nil {
		return err
	}
	harbor.HarborCluster.Spec.AdminPasswordSecret = adminPasswordSecretName
	return nil
}
