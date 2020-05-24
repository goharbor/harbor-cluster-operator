package cache

import (
	"errors"
	"fmt"
	rediscli "github.com/go-redis/redis"
	goharborv1 "github.com/goharbor/harbor-cluster-operator/api/v1"
	"github.com/goharbor/harbor-cluster-operator/lcm"
	corev1 "k8s.io/api/core/v1"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"strings"
)

const (
	HarborChartMuseum = "chartmuseum"
	HarborClair       = "clair"
	HarborJobService  = "jobservice"
	HarborRegistry    = "registry"
)

var (
	components = []string{
		HarborChartMuseum,
		HarborClair,
		HarborJobService,
		HarborRegistry,
	}
)

// Readiness reconcile will check Redis sentinel cluster if that has available.
// It does:
// - create redis connection pool
// - ping redis server
// - return redis properties if redis has available
func (redis *RedisReconciler) Readiness() error {
	var (
		client *rediscli.Client
		err    error
	)

	switch redis.HarborCluster.Spec.Redis.Kind {
	case "external":
		client, err = redis.GetExternalRedisInfo()
	case "inCluster":
		client, err = redis.GetInClusterRedisInfo()
	}

	if err != nil {
		return err
	}

	defer client.Close()

	if err := client.Ping().Err(); err != nil {
		redis.Log.Error(err, "Fail to check Redis.", "namespace", redis.Namespace, "name", redis.Name)
		return err
	}

	redis.Log.Info("Redis already ready.", "namespace", redis.Namespace, "name", redis.Name)

	for _, component := range components {
		url := redis.RedisConnect.GenRedisConnURL()
		if err := redis.DeployComponentSecret(component, url, ""); err != nil {
			return err
		}
	}
	return nil
}

// DeployComponentSecret deploy harbor component redis secret
func (redis *RedisReconciler) DeployComponentSecret(component, url, namespace string) error {
	secret := &corev1.Secret{}
	secretName := fmt.Sprintf("%s-redis", component)
	sc := redis.generateHarborCacheSecret(component, secretName, url, namespace)

	if err := controllerutil.SetControllerReference(redis.HarborCluster, sc, redis.Scheme); err != nil {
		return err
	}
	err := redis.Client.Get(types.NamespacedName{Name: secretName, Namespace: redis.Namespace}, secret)
	if err != nil && kerr.IsNotFound(err) {
		redis.Log.Info("Creating Harbor Component Secret",
			"namespace", redis.Namespace,
			"name", secretName,
			"component", component)
		err = redis.Client.Create(sc)
		if err != nil {
			return err
		}

		redis.Properties.New(component, secretName)
	} else if err != nil {
		return err
	}

	return nil
}

func (redis *RedisReconciler) GetExternalRedisInfo() (*rediscli.Client, error) {
	var (
		connect  *RedisConnect
		endpoint []string
		port     string
		client   *rediscli.Client
		err      error
		pw       string
	)
	spec := redis.HarborCluster.Spec.Redis.Spec
	switch spec.Schema {
	case "sentinel":
		if len(spec.Hosts) < 1 || spec.GroupName == "" {
			return nil, errors.New(".redis.spec.hosts or .redis.spec.groupName is invalid")
		}

		endpoint, port = GetExternalRedisHost(spec)

		if spec.SecretName != "" {
			pw, err = GetExternalRedisPassword(spec, redis.Namespace)
		}

		connect = &RedisConnect{
			Endpoint:  strings.Join(endpoint[:], ","),
			Port:      port,
			Password:  pw,
			GroupName: spec.GroupName,
		}

		redis.RedisConnect = connect
		client = connect.NewRedisPool()
	case "redis":
		if len(spec.Hosts) != 1 {
			return nil, errors.New(".redis.spec.hosts is invalid")
		}
		endpoint, port = GetExternalRedisHost(spec)

		if spec.SecretName != "" {
			pw, err = GetExternalRedisPassword(spec, redis.Namespace)
		}

		connect = &RedisConnect{
			Endpoint:  fmt.Sprintf("%s:%s", endpoint, port),
			Port:      port,
			Password:  pw,
			GroupName: spec.GroupName,
		}
		redis.RedisConnect = connect
		client = connect.NewRedisClient()
	}

	if err != nil {
		return nil, err
	}

	return client, nil
}

// GetExternalRedisHost returns external redis host list and port
func GetExternalRedisHost(spec *goharborv1.RedisSpec) ([]string, string) {
	var (
		endpoint []string
		port     string
	)
	for _, host := range spec.Hosts {
		sp := host.Host + ":" + host.Port
		endpoint = append(endpoint, sp)
		port = host.Port
	}
	return endpoint, port
}

// GetExternalRedisPassword returns external redis password
func GetExternalRedisPassword(spec *goharborv1.RedisSpec, namespace string) (string, error) {
	external := &RedisReconciler{
		Name:      spec.SecretName,
		Namespace: namespace,
	}

	pw, err := external.GetRedisPassword()
	if err != nil {
		return "", err
	}

	return pw, err
}

// GetInClusterRedisInfo returns inCluster redis sentinel pool client
func (redis *RedisReconciler) GetInClusterRedisInfo() (*rediscli.Client, error) {
	password, err := redis.GetRedisPassword()
	if err != nil {
		return nil, err
	}

	_, sentinelPodList, err := redis.GetDeploymentPods()
	if err != nil {
		redis.Log.Error(err, "Fail to get deployment pods.")
		return nil, err
	}

	_, redisPodList, err := redis.GetStatefulSetPods()
	if err != nil {
		redis.Log.Error(err, "Fail to get deployment pods.")
		return nil, err
	}

	if len(sentinelPodList.Items) == 0 || len(redisPodList.Items) == 0 {
		redis.Log.Info("pod list is empty，pls wait.")
		return nil, nil
	}

	sentinelPodArray := sentinelPodList.Items

	_, currentSentinelPods := redis.GetPodsStatus(sentinelPodArray)

	if len(currentSentinelPods) == 0 {
		return nil, errors.New("need to requeue")
	}

	endpoint := redis.GetSentinelServiceUrl(currentSentinelPods)

	connect := &RedisConnect{
		Endpoint:  endpoint,
		Port:      RedisSentinelConnPort,
		Password:  password,
		GroupName: RedisSentinelConnGroup,
	}

	redis.RedisConnect = connect

	client := connect.NewRedisPool()

	return client, nil
}

func cacheNotReadyStatus(reason, message string) *lcm.CRStatus {
	return &lcm.CRStatus{
		Condition: goharborv1.HarborClusterCondition{
			Type:               goharborv1.CacheReady,
			Status:             corev1.ConditionFalse,
			LastTransitionTime: metav1.Now(),
			Reason:             reason,
			Message:            message,
		},
		Properties: nil,
	}
}

func cacheReadyStatus(properties *lcm.Properties) *lcm.CRStatus {
	return &lcm.CRStatus{
		Condition: goharborv1.HarborClusterCondition{
			Type:               goharborv1.CacheReady,
			Status:             corev1.ConditionFalse,
			LastTransitionTime: metav1.Now(),
			Reason:             "reason",
			Message:            "cache already ready",
		},
		Properties: *properties,
	}
}
