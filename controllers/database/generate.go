package database

import (
	"fmt"

	"github.com/goharbor/harbor-cluster-operator/controllers/database/api"
	pg "github.com/zalando/postgres-operator/pkg/apis/acid.zalan.do/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

var (
	databaseFailoversGVR = pg.SchemeGroupVersion.WithResource(pg.PostgresCRDResourcePlural)
)

// generatePostgreCR returns PostgreSqls CRs
func (postgres *PostgreSQLReconciler) generatePostgresCR() (*unstructured.Unstructured, error) {
	resource := postgres.GetPostgreResource()
	replica := postgres.GetPostgreReplica()
	storageSize := postgres.GetPostgreStorageSize()
	version := postgres.GetPostgreVersion()
	name := fmt.Sprintf("%s-%s", postgres.HarborCluster.Namespace, postgres.HarborCluster.Name)

	conf := &api.Postgresql{
		TypeMeta: metav1.TypeMeta{
			Kind:       "postgresql",
			APIVersion: "acid.zalan.do/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: postgres.HarborCluster.Namespace,
			Labels:    postgres.Labels,
		},
		Spec: api.PostgresSpec{
			Volume: api.Volume{
				Size: storageSize,
			},
			TeamID:            postgres.HarborCluster.Namespace,
			NumberOfInstances: replica,
			Users: map[string]api.UserFlags{
				"zalando": {
					"superuser",
					"createdb",
				},
				"foo_user": {},
			},
			Patroni: api.Patroni{
				InitDB: map[string]string{
					"encoding":       "UTF8",
					"locale":         "en_US.UTF-8",
					"data-checksums": "true",
				},
				PgHba: []string{
					"hostssl all all 0.0.0.0/0 md5",
					"host    all all 0.0.0.0/0 md5",
				},
			},
			Databases: map[string]string{
				"foo": "zalando",
			},
			PostgresqlParam: api.PostgresqlParam{
				PgVersion: version,
			},
			Resources: resource,
		},
	}

	mapResult, err := runtime.DefaultUnstructuredConverter.ToUnstructured(conf)
	if err != nil {
		return nil, err
	}

	data := unstructured.Unstructured{Object: mapResult}

	return &data, nil
}

//generateHarborDatabaseSecret returns database connection secret
func (postgres *PostgreSQLReconciler) generateHarborDatabaseSecret(conn *Connect, secretName string) *corev1.Secret {

	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: postgres.HarborCluster.Namespace,
			Labels:    postgres.Labels,
		},
		StringData: map[string]string{
			"host":     conn.Host,
			"port":     conn.Port,
			"database": conn.Database,
			"username": conn.Username,
			"password": conn.Password,
		},
	}
}
