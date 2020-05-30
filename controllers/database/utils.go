package database

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

// GetDatabaseConn is get database connection
func (postgre *PostgreSQLReconciler) GetDatabaseConn() (*Connect, error) {
	var (
		host     string
		port     string
		username string
		password string
		database string
	)
	secret, err := postgre.GetSecret()
	if err != nil {
		return nil, err
	}
	for k, v := range secret {
		switch k {
		case "host":
			host = string(v)
		case "port":
			port = string(v)
		case "username":
			username = string(v)
		case "password":
			password = string(v)
		case "database":
			database = string(v)
		}
	}

	conn := &Connect{
		Host:     host,
		Port:     port,
		Password: password,
		Username: username,
		Database: database,
	}

	return conn, nil
}

// GetSecret returns the database connection Secret
func (postgre *PostgreSQLReconciler) GetSecret() (map[string][]byte, error) {
	secret := &corev1.Secret{}
	err := postgre.Client.Get(types.NamespacedName{Name: postgre.Name, Namespace: postgre.Namespace}, secret)
	if err != nil {
		return nil, err
	}
	data := secret.Data
	return data, nil
}
