package database

import (
	"errors"
	"fmt"

	goharborv1 "github.com/goharbor/harbor-cluster-operator/apis/goharbor.io/v1alpha1"
	"github.com/goharbor/harbor-cluster-operator/controllers/k8s"
	"github.com/goharbor/harbor-cluster-operator/lcm"
	"github.com/jackc/pgx/v4"
	corev1 "k8s.io/api/core/v1"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	labels1 "k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const (
	HarborCore         = "core"
	HarborClair        = "clair"
	HarborNotaryServer = "notaryServer"
	HarborNotarySigner = "notarySigner"

	CoreDatabase         = "core"
	ClairDatabase        = "clair"
	NotaryServerDatabase = "notaryserver"
	NotarySignerDatabase = "notarysigner"

	CoreSecretName         = "core"
	ClairSecretName        = "clair"
	NotaryServerSecretName = "notary-server"
	NotarySignerSecretName = "notary-signer"
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
		conn, client, err = postgres.GetInClusterDatabaseInfo()
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
	components := map[string]string{
		HarborCore: CoreSecretName,
	}
	if postgres.HarborCluster.Spec.Clair != nil {
		components[HarborClair] = ClairSecretName
	}
	if postgres.HarborCluster.Spec.Notary != nil {
		components[HarborNotaryServer] = NotaryServerSecretName
		components[HarborNotarySigner] = NotarySignerSecretName
	}

	for key, component := range components {
		secretName := getComponentSecretName(component)
		propertyName := getPropertyName(key)
		if err := postgres.DeployComponentSecret(conn, component, secretName, key); err != nil {
			return nil, err
		}
		properties.Add(propertyName, secretName)
	}

	crStatus := lcm.New(goharborv1.DatabaseReady).
		WithStatus(corev1.ConditionTrue).
		WithReason("database already ready").
		WithMessage("harbor component database secrets are already create.").
		WithProperties(*properties)
	return crStatus, nil
}

func getPropertyName(key string) string {
	return fmt.Sprintf("%sSecret", key)
}

func getComponentSecretName(component string) string {
	return fmt.Sprintf("%s-database", component)
}

// DeployComponentSecret deploy harbor component database secret
func (postgres *PostgreSQLReconciler) DeployComponentSecret(conn *Connect, component, secretName, propertyName string) error {
	secret := &corev1.Secret{}
	sc := postgres.generateHarborDatabaseSecret(conn, secretName, propertyName)

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

// GetInClusterDatabaseInfo returns inCluster database connection client
func (postgres *PostgreSQLReconciler) GetInClusterDatabaseInfo() (*Connect, *pgx.Conn, error) {
	var (
		connect *Connect
		client  *pgx.Conn
		err     error
	)

	pw, err := postgres.GetInClusterDatabasePassword()
	if err != nil {
		return connect, client, err
	}

	if connect, err = postgres.GetInClusterDatabaseConn(postgres.GetDatabaseName(), pw); err != nil {
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

// GetInClusterDatabaseConn returns inCluster database connection info
func (postgres *PostgreSQLReconciler) GetInClusterDatabaseConn(name, pw string) (*Connect, error) {
	host, err := postgres.GetInClusterHost(name)
	if err != nil {
		return nil, err
	}
	conn := &Connect{
		Host:     host,
		Port:     InClusterDatabasePort,
		Password: pw,
		Username: InClusterDatabaseUserName,
		Database: InClusterDatabaseName,
	}
	return conn, nil
}

func GenInClusterPasswordSecretName(teamID, name string) string {
	return fmt.Sprintf("postgres.%s-%s.credentials", teamID, name)
}

// GetInClusterHost returns the Database master pod ip or service name
func (postgres *PostgreSQLReconciler) GetInClusterHost(name string) (string, error) {
	var (
		url string
		err error
	)
	_, err = rest.InClusterConfig()
	if err != nil {
		url, err = postgres.GetMasterPodsIP()
		if err != nil {
			return url, err
		}
	} else {
		url = fmt.Sprintf("%s.%s.svc", name, postgres.HarborCluster.Namespace)
	}

	return url, nil
}

func (postgres *PostgreSQLReconciler) GetDatabaseName() string {
	return fmt.Sprintf("%s-%s", postgres.HarborCluster.Namespace, postgres.HarborCluster.Name)
}

// GetInClusterDatabasePassword is get inCluster postgresql password
func (postgres *PostgreSQLReconciler) GetInClusterDatabasePassword() (string, error) {
	var pw string

	secretName := GenInClusterPasswordSecretName(postgres.HarborCluster.Namespace, postgres.HarborCluster.Name)
	secret, err := postgres.GetSecret(secretName)
	if err != nil {
		return pw, err
	}

	for k, v := range secret {
		if k == InClusterDatabasePasswordKey {
			pw = string(v)
			return pw, nil
		}
	}
	return pw, nil
}

// GetStatefulSetPods returns the postgresql master pod
func (postgres *PostgreSQLReconciler) GetStatefulSetPods() (*corev1.PodList, error) {
	name := postgres.GetDatabaseName()
	label := map[string]string{
		"application":  "spilo",
		"cluster-name": name,
		"spilo-role":   "master",
	}

	opts := &client.ListOptions{}
	set := labels1.SelectorFromSet(label)
	opts.LabelSelector = set
	pod := &corev1.PodList{}

	if err := postgres.Client.List(opts, pod); err != nil {
		postgres.Log.Error(err, "fail to get pod.",
			"namespace", postgres.HarborCluster.Namespace, "name", name)
		return nil, err
	}
	return pod, nil
}

// GetMasterPodsIP returns postgresql master node ip
func (postgres *PostgreSQLReconciler) GetMasterPodsIP() (string, error) {
	var masterIP string
	podList, err := postgres.GetStatefulSetPods()
	if err != nil {
		return masterIP, err
	}
	if len(podList.Items) > 1 {
		return masterIP, errors.New("the number of master node copies cannot exceed 1")
	}
	for _, p := range podList.Items {
		if p.DeletionTimestamp != nil {
			continue
		}
		masterIP = p.Status.PodIP
	}
	return masterIP, nil
}
