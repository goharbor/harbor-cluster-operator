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
func (postgres *PostgreSQLReconciler) Readiness() (*lcm.CRStatus, error) {
	var (
		conn   *Connect
		client *pgx.Conn
		err    error
	)

	switch postgres.HarborCluster.Spec.Database.Kind {
	case "external":
		conn, client, err = postgres.GetExternalDatabaseInfo()
	case "inCluster":

	default:
		return nil, errors.New("fail to check database kind")
	}

	if err != nil {
		return nil, err
	}

	defer client.Close(postgres.Ctx)

	if err := client.Ping(postgres.Ctx); err != nil {
		postgres.Log.Error(err, "Fail to check Database.", "namespace", postgres.HarborCluster.Namespace, "name", postgres.HarborCluster.Name)
		return nil, err
	}
	postgres.Log.Info("Database already ready.", "namespace", postgres.HarborCluster.Namespace, "name", postgres.HarborCluster.Name)

	properties := &lcm.Properties{}
	for _, component := range components {
		secretName := fmt.Sprintf("%s-database", component)
		propertyName := fmt.Sprintf("%sSecret", component)
		if err := postgres.DeployComponentSecret(conn, component, secretName); err != nil {
			return nil, err
		}
		properties = properties.Add(propertyName, secretName)
	}

	crStatus := lcm.New(goharborv1.DatabaseReady).
		WithStatus(corev1.ConditionTrue).
		WithReason("database already ready").
		WithMessage("harbor component database secrets are already create.").
		WithProperties(*properties)
	return crStatus, nil
}

// DeployComponentSecret deploy harbor component database secret
func (postgres *PostgreSQLReconciler) DeployComponentSecret(conn *Connect, component, secretName string) error {
	secret := &corev1.Secret{}
	sc := postgres.generateHarborDatabaseSecret(conn, secretName)

	if err := controllerutil.SetControllerReference(postgres.HarborCluster, sc, postgres.Scheme); err != nil {
		return err
	}
	err := postgres.Client.Get(types.NamespacedName{Name: secretName, Namespace: postgres.HarborCluster.Namespace}, secret)
	if err != nil {
		if kerr.IsNotFound(err) {
			postgres.Log.Info("Creating Harbor Component Secret",
				"namespace", postgres.HarborCluster.Namespace,
				"name", secretName,
				"component", component)
			err = postgres.Client.Create(sc)
			if err != nil {
				return err
			}

			return nil
		}
		return err
	}
	return nil
}

// GetExternalDatabaseInfo returns external database connection client
func (postgres *PostgreSQLReconciler) GetExternalDatabaseInfo() (*Connect, *pgx.Conn, error) {
	var (
		connect *Connect
		client  *pgx.Conn
		err     error
	)
	spec := postgres.HarborCluster.Spec.Database.Spec
	if spec.SecretName == "" {
		return connect, client, errors.New(".database.spec.secretName is invalid")
	}

	if connect, err = postgres.GetExternalDatabaseConn(spec.SecretName, postgres.Client); err != nil {
		return connect, client, err
	}

	url := connect.GenDatabaseUrl()

	client, err = pgx.Connect(postgres.Ctx, url)
	if err != nil {
		postgres.Log.Error(err, "Unable to connect to database")
		return connect, client, err
	}

	return connect, client, nil
}

// GetExternalDatabaseConn returns external database connection info
func (postgres *PostgreSQLReconciler) GetExternalDatabaseConn(secretName string, client k8s.Client) (*Connect, error) {

	conn, err := postgres.GetDatabaseConn(secretName)
	if err != nil {
		return nil, err
	}

	return conn, err
}
