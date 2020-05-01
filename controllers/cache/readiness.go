package cache

import (
	"errors"
	rediscli "github.com/go-redis/redis"
	"github.com/goharbor/harbor-cluster-operator/lcm"
	corev1 "k8s.io/api/core/v1"
	"strings"
	"time"
)

const (
	RedisDownScaling = "RedisDownScaling"
	RedisUpScaling   = "RedisUpScaling"

	MessageRedisCluster = "Redis  %s already created."

	UpdateMessageRedisCluster = "Redis  %s already update."

	MessageRedisDownScaling = "Redis downscale from %d to %d"
	MessageRedisUpScaling   = "Redis upscale from %d to %d"

	RedisSentinelConnPort = "26379"
)

type RedisConnect struct {
	SentinelEndpoint string
	SentinelPort     string
	Password         string
	GroupName        string
}

// NewRedisConnection returns redis connection
func NewRedisConnection(endpoint, port, password, groupName string) *RedisConnect {
	return &RedisConnect{
		SentinelEndpoint: endpoint,
		SentinelPort:     port,
		Password:         password,
		GroupName:        groupName,
	}
}

// Readiness reconcile will check Redis sentinel cluster if that has available.
// It does:
// - create redis connection pool
// - ping redis server
// - return redis properties if redis has available
func (redis *RedisReconciler) Readiness(crStatus *lcm.CRStatus) (*lcm.CRStatus, error) {
	password, err := redis.GetRedisPassword()
	if err != nil {
		return crStatus, err
	}

	_, sentinelPodList, err := redis.GetDeploymentPods()
	if err != nil {
		redis.Log.Error(err, "Fail to get deployment pods.")
		return crStatus, err
	}

	_, redisPodList, err := redis.GetStatefulSetPods()
	if err != nil {
		redis.Log.Error(err, "Fail to get deployment pods.")
		return crStatus, err
	}

	if len(sentinelPodList.Items) == 0 || len(redisPodList.Items) == 0 {
		redis.Log.Info("pod list is empty，pls wait.")
		return crStatus, nil
	}

	sentinelPodArray := sentinelPodList.Items
	redisPodArray := redisPodList.Items

	_, currentSentinelPods := redis.GetPodsStatus(sentinelPodArray)
	_, currentRedisPods := redis.GetPodsStatus(redisPodArray)

	if len(currentSentinelPods) == 0 {
		return crStatus, errors.New("Need to Requeue")
	}
	endpoint := redis.GetServiceUrl(currentSentinelPods)
	connect := NewRedisConnection(endpoint, "26379", password, "mymaster")
	client := connect.NewRedisPool()
	defer client.Close()

	if err := client.Ping().Err(); err != nil {
		redis.Log.Error(err, "Fail to check Redis.", "namespace", redis.Namespace, "name", redis.Name)
		return crStatus, err
	}

	crStatus.Phase = lcm.ReadyPhase
	properties := lcm.Properties{}
	conn := properties.New(lcm.ProperConn, endpoint)
	port := properties.New(lcm.ProperPort, RedisSentinelConnPort)
	nodes := properties.New(lcm.ProperNodes, len(currentRedisPods))

	properties = append(properties, conn, port, nodes)

	crStatus.Properties = properties

	return crStatus, nil
}

// NewRedisPool returns redis client
func (c *RedisConnect) NewRedisPool() *rediscli.Client {

	return BuildRedisPool(c.SentinelEndpoint, c.SentinelPort, c.Password, c.GroupName, 0)
}

// BuildRedisPool returns redis connection pool client
func BuildRedisPool(redisSentinelIP, redisSentinelPort, redisSentinelPassword, redisGroupName string, redisIndex int) *rediscli.Client {

	var sentinelsInfo []string
	sentinels := strings.Split(redisSentinelIP, ",")
	if len(sentinels) > 0 {
		for _, s := range sentinels {
			sp := s + ":" + redisSentinelPort
			sentinelsInfo = append(sentinelsInfo, sp)
		}
	}

	options := &rediscli.FailoverOptions{
		MasterName:         redisGroupName,
		SentinelAddrs:      sentinelsInfo,
		Password:           redisSentinelPassword,
		DB:                 redisIndex,
		PoolSize:           100,
		DialTimeout:        10 * time.Second,
		ReadTimeout:        30 * time.Second,
		WriteTimeout:       30 * time.Second,
		PoolTimeout:        30 * time.Second,
		IdleTimeout:        time.Millisecond,
		IdleCheckFrequency: time.Millisecond,
	}

	client := rediscli.NewFailoverClient(options)

	return client

}

// GetPodsStatus returns deleting  and current pod list
func (redis *RedisReconciler) GetPodsStatus(podArray []corev1.Pod) ([]corev1.Pod, []corev1.Pod) {
	deletingPods := make([]corev1.Pod, 0)
	currentPods := make([]corev1.Pod, 0, len(podArray))
	currentPodsByPhase := make(map[corev1.PodPhase][]corev1.Pod)

	for _, p := range podArray {
		if p.DeletionTimestamp != nil {
			deletingPods = append(deletingPods, p)
			continue
		}
		currentPods = append(currentPods, p)
		podsInPhase, ok := currentPodsByPhase[p.Status.Phase]
		if !ok {
			podsInPhase = []corev1.Pod{p}
		} else {
			podsInPhase = append(podsInPhase, p)
		}
		currentPodsByPhase[p.Status.Phase] = podsInPhase
	}
	return deletingPods, currentPods
}

