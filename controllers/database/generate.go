package database

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//generateHarborDatabaseSecret returns database connection secret
func (postgre *PostgreSQLReconciler) generateHarborDatabaseSecret(secretName string) *corev1.Secret {

	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: postgre.Namespace,
			Labels:    postgre.Labels,
		},
		StringData: map[string]string{
			"host":     postgre.Connect.Host,
			"port":     postgre.Connect.Port,
			"database": postgre.Connect.Database,
			"username": postgre.Connect.Username,
			"password": postgre.Connect.Password,
		},
	}
}
