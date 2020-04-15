package cache

import (
	"errors"
	rediscli "github.com/go-redis/redis"
	corev1 "k8s.io/api/core/v1"
	"strings"
	"time"
)

type RedisConnect struct {
	SentinelEndpoint string
	SentinelPort     string
	Password         string
	GroupName        string
}

func NewRedisConnection(endpoint, port, password, groupName string) *RedisConnect {
	return &RedisConnect{
		SentinelEndpoint: endpoint,
		SentinelPort:     port,
		Password:         password,
		GroupName:        groupName,
	}
}

func (d *defaultCache) Readiness() error {
	password, err := d.GetRedisPassword()
	if err != nil {
		return err
	}

	_, podList, err := d.GetDeploymentPods()
	if err != nil {
		d.Log.Error(err, "Fail to get deployment pods.")
		return err
	}

	if len(podList.Items) == 0 {
		d.Log.Info("pod list is empty，pls wait.")
		return nil
	}

	podArray := podList.Items

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

	if len(currentPods) == 0 {
		return errors.New("Need to Requeue")
	}
	endpoint := d.GetServiceUrl(currentPods)
	connect := NewRedisConnection(endpoint, "26379", password, "mymaster")
	client := connect.NewRedisPool()
	defer client.Close()

	if err := client.Ping().Err(); err != nil {
		d.Log.Error(err, "Fail to check Redis.", "namespace", d.Request.Namespace, "name", d.Request.Name)
		return err
	}
	d.Log.Info("Redis is ready.")
	return nil
}

func (c *RedisConnect) NewRedisPool() *rediscli.Client {

	return BuildRedisPool(c.SentinelEndpoint, c.SentinelPort, c.Password, c.GroupName, 0)
}

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
