package database

import (
	goharborv1 "github.com/goharbor/harbor-cluster-operator/api/v1"
	"github.com/goharbor/harbor-cluster-operator/controllers/database/api"
	"github.com/goharbor/harbor-cluster-operator/lcm"

	//pg "github.com/zalando/postgres-operator/pkg/apis/acid.zalan.do/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

const (
	DefaultDatabaseReplica = 3
	DefaultDatabaseMemory  = "1Gi"
	DefaultDatabaseVersion = "12"
)

func (postgres *PostgreSQLReconciler) GetDatabases() map[string]string {
	databases := map[string]string{
		CoreDatabase: "zalando",
	}

	if postgres.HarborCluster.Spec.Clair != nil {
		databases[ClairDatabase] = "zalando"
	}

	if postgres.HarborCluster.Spec.Notary != nil {
		databases[NotaryServerDatabase] = "zalando"
		databases[NotarySignerDatabase] = "zalando"
	}

	return databases
}

// GetDatabaseConn is getting database connection
func (postgres *PostgreSQLReconciler) GetDatabaseConn(secretName string) (*Connect, error) {
	secret, err := postgres.GetSecret(secretName)
	if err != nil {
		return nil, err
	}

	conn := &Connect{
		Host:     string(secret["host"]),
		Port:     string(secret["port"]),
		Password: string(secret["password"]),
		Username: string(secret["username"]),
		Database: string(secret["database"]),
	}

	return conn, nil
}

// GetSecret returns the database connection Secret
func (postgres *PostgreSQLReconciler) GetSecret(secretName string) (map[string][]byte, error) {
	secret := &corev1.Secret{}
	err := postgres.Client.Get(types.NamespacedName{Name: secretName, Namespace: postgres.HarborCluster.Namespace}, secret)
	if err != nil {
		return nil, err
	}
	data := secret.Data
	return data, nil
}

// GetPostgreResource returns postgres resource
func (postgres *PostgreSQLReconciler) GetPostgreResource() api.Resources {
	resources := api.Resources{}

	if postgres.HarborCluster.Spec.Database.Spec == nil {
		resources.ResourceRequests = api.ResourceDescription{
			CPU:    "1",
			Memory: "1Gi",
		}
		resources.ResourceRequests = api.ResourceDescription{
			CPU:    "2",
			Memory: "2Gi",
		}
		return resources
	}

	cpu := postgres.HarborCluster.Spec.Database.Spec.Resources.Requests.Cpu()
	mem := postgres.HarborCluster.Spec.Database.Spec.Resources.Requests.Memory()

	request := api.ResourceDescription{}
	if cpu != nil {
		request.CPU = cpu.String()
	}
	if mem != nil {
		request.Memory = mem.String()
	}
	resources.ResourceRequests = request
	resources.ResourceLimits = request

	return resources
}

// GetRedisServerReplica returns postgres replicas
func (postgres *PostgreSQLReconciler) GetPostgreReplica() int32 {
	if postgres.HarborCluster.Spec.Database.Spec == nil {
		return DefaultDatabaseReplica
	}

	if postgres.HarborCluster.Spec.Database.Spec.Replicas == 0 {
		return DefaultDatabaseReplica
	}
	return int32(postgres.HarborCluster.Spec.Database.Spec.Replicas)
}

// GetPostgreStorageSize returns Postgre storage size
func (postgres *PostgreSQLReconciler) GetPostgreStorageSize() string {
	if postgres.HarborCluster.Spec.Database.Spec == nil {
		return DefaultDatabaseMemory
	}

	if postgres.HarborCluster.Spec.Database.Spec.Storage == "" {
		return DefaultDatabaseMemory
	}
	return postgres.HarborCluster.Spec.Database.Spec.Storage
}

func (postgres *PostgreSQLReconciler) GetPostgreVersion() string {
	if postgres.HarborCluster.Spec.Database.Spec == nil {
		return DefaultDatabaseVersion
	}

	if postgres.HarborCluster.Spec.Database.Spec.Version == "" {
		return DefaultDatabaseVersion
	}

	return postgres.HarborCluster.Spec.Database.Spec.Version
}

func databaseNotReadyStatus(reason, message string) *lcm.CRStatus {
	return lcm.New(goharborv1.DatabaseReady).
		WithStatus(corev1.ConditionFalse).
		WithReason(reason).
		WithMessage(message)
}

func databaseUnknownStatus() *lcm.CRStatus {
	return lcm.New(goharborv1.DatabaseReady).
		WithStatus(corev1.ConditionUnknown)
}
