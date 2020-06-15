package database

import (
	"errors"
	"fmt"
	goharborv1 "github.com/goharbor/harbor-cluster-operator/api/v1"
	"github.com/goharbor/harbor-cluster-operator/controllers/k8s"
	"github.com/goharbor/harbor-cluster-operator/lcm"
	"github.com/jackc/pgx/v4"
	corev1 "k8s.io/api/core/v1"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const (
	HarborCore         = "core"
	HarborClair        = "clair"
	HarborNotaryServer = "notary-server"
	HarborNotarySigner = "notary-signer"
)

var (
	components = []string{
		HarborCore,
		HarborClair,
		HarborNotaryServer,
		HarborNotarySigner,
	}
)

// Readiness reconcile will check postgre sql cluster if that has available.
// It does:
// - create postgre connection pool
// - ping postgre server
// - return postgre properties if postgre has available
func (postgre *PostgreSQLReconciler) Readiness() error {
	var (
		client *pgx.Conn
		err    error
	)

	switch postgre.HarborCluster.Spec.Database.Kind {
	case "external":
		client, err = postgre.GetExternalDatabaseInfo()
	}

	if err != nil {
		return err
	}

	defer client.Close(postgre.Ctx)

	if err := client.Ping(postgre.Ctx); err != nil {
		postgre.Log.Error(err, "Fail to check Database.", "namespace", postgre.Namespace, "name", postgre.Name)
		return err
	}
	postgre.Log.Info("Database already ready.", "namespace", postgre.Namespace, "name", postgre.Name)

	for _, component := range components {
		if err := postgre.DeployComponentSecret(component); err != nil {
			return err
		}
	}

	postgre.CRStatus = lcm.New(goharborv1.DatabaseReady).
		WithStatus(corev1.ConditionTrue).
		WithReason("database already ready").
		WithMessage("harbor component database secrets are already create.").
		WithProperties(*postgre.Properties)
	return nil
}

// DeployComponentSecret deploy harbor component database secret
func (postgre *PostgreSQLReconciler) DeployComponentSecret(component string) error {
	secret := &corev1.Secret{}
	secretName := fmt.Sprintf("%s-database", component)
	propertyName := fmt.Sprintf("%sSecret", component)
	sc := postgre.generateHarborDatabaseSecret(secretName)

	if err := controllerutil.SetControllerReference(postgre.HarborCluster, sc, postgre.Scheme); err != nil {
		return err
	}
	err := postgre.Client.Get(types.NamespacedName{Name: secretName, Namespace: postgre.Namespace}, secret)
	if err != nil {
		if kerr.IsNotFound(err) {
			postgre.Log.Info("Creating Harbor Component Secret",
				"namespace", postgre.Namespace,
				"name", secretName,
				"component", component)
			err = postgre.Client.Create(sc)
			if err != nil {
				return err
			}
			postgre.Properties = postgre.Properties.New(propertyName, secretName)
			return nil
		}
		return err
	}
	postgre.Properties = postgre.Properties.New(propertyName, secretName)
	return nil
}

// GetExternalDatabaseInfo returns external database connection client
func (postgre *PostgreSQLReconciler) GetExternalDatabaseInfo() (*pgx.Conn, error) {
	var (
		connect *Connect
		client  *pgx.Conn
		err     error
	)
	spec := postgre.HarborCluster.Spec.Database.Spec
	if spec.SecretName == "" {
		return nil, errors.New(".database.spec.secretName is invalid")
	}

	if connect, err = GetExternalDatabaseConn(spec, postgre.Namespace, postgre.Client); err != nil {
		return nil, err
	}

	postgre.Connect = connect

	url := connect.GenDatabaseUrl()

	client, err = pgx.Connect(postgre.Ctx, url)
	if err != nil {
		postgre.Log.Error(err, "Unable to connect to database")
		return nil, err
	}

	return client, nil
}

// GetExternalDatabaseConn returns external database connection info
func GetExternalDatabaseConn(spec *goharborv1.PostgresSQL, namespace string, client k8s.Client) (*Connect, error) {
	external := &PostgreSQLReconciler{
		Name:      spec.SecretName,
		Namespace: namespace,
		Client:    client,
	}

	conn, err := external.GetDatabaseConn()
	if err != nil {
		return nil, err
	}

	return conn, err
}
