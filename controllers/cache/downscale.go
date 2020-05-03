package cache

import (
	"errors"
	"fmt"
	rediscli "github.com/go-redis/redis"
	"github.com/goharbor/harbor-cluster-operator/lcm"
	redisCli "github.com/spotahome/redis-operator/api/redisfailover/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"strings"

	"k8s.io/apimachinery/pkg/runtime"
	"reflect"
)

// DownScale reconcile will downscale Redis sentinel cluster if replicas downscale.
// It does:
// - check resource have been updated, if updated priority to perform rolling upgrade
// - get a list of nodes that need to be leave
// - get redis server and sentinel pod list
// - get redis server password
// - check whether the leave node is the master node, if is a master node, needs to manual failover
func (redis *RedisReconciler) DownScale(crStatus *lcm.CRStatus) (*lcm.CRStatus, error) {

	crdClient := redis.DClient.WithResource(virtualServiceGVR).WithNamespace(redis.Namespace)
	if redis.ExpectCR == nil {
		return crStatus, nil
	}

	var actualCR redisCli.RedisFailover
	var expectCR redisCli.RedisFailover

	if err := runtime.DefaultUnstructuredConverter.
		FromUnstructured(redis.ActualCR.UnstructuredContent(), &actualCR); err != nil {
		return crStatus, err
	}

	if err := runtime.DefaultUnstructuredConverter.
		FromUnstructured(redis.ExpectCR.UnstructuredContent(), &expectCR); err != nil {
		return crStatus, err
	}

	if isUpgradeResource(&expectCR, &actualCR) {
		redis.Log.Info("Resources have been updated, priority to perform rolling upgrade",
			"namespace", redis.Namespace, "name", redis.Name)
	}

	expectReplica := expectCR.Spec.Redis.Replicas
	actualReplica := actualCR.Spec.Redis.Replicas

	downscales := calculateDownscales(expectCR, actualCR)
	leavingNodes := leavingNodeNames(downscales)

	if len(leavingNodes) != 0 {
		redis.Log.Info("Migrating data away from nodes", "nodes", leavingNodes)
	}

	if expectReplica < actualReplica {

		_, redisPodList, err := redis.GetStatefulSetPods()
		if err != nil {
			redis.Log.Error(err, "Fail to get deployment pods.")
			return crStatus, err
		}

		_, sentinelPodList, err := redis.GetDeploymentPods()
		if err != nil {
			redis.Log.Error(err, "Fail to get deployment pods.")
			return crStatus, err
		}

		sentinelPodArray := sentinelPodList.Items
		redisPodArray := redisPodList.Items

		_, currentSentinelPods := redis.GetPodsStatus(sentinelPodArray)
		_, currentRedisPods := redis.GetPodsStatus(redisPodArray)

		if len(currentSentinelPods) == 0 {
			return crStatus, errors.New("Need to Requeue")
		}

		password, err := redis.GetRedisPassword()
		if err != nil {
			return crStatus, err
		}

		var master []string
		masterMap := map[string]string{}
		for _, pod := range currentRedisPods {
			ok, err := IsMaster(pod.Status.PodIP, password)
			if err != nil {
				return crStatus, err
			}

			if ok {
				master = append(master, pod.Name)
				masterMap[pod.Name] = pod.Status.PodIP
			}
		}

		if len(master) > 1 {
			crStatus.Phase = lcm.FailedPhase
			return crStatus, errors.New("master node more than 1, need requeue")
		}

		if len(Intersect(leavingNodes, master)) > 0 {
			redis.Log.Info("Redis need to manual failover",
				"namespace", redis.Namespace, "name", redis.Name,
				"master", master)
			endpoint := redis.GetSentinelServiceUrl(currentSentinelPods)
			if err := ManualFailoverSentinel(endpoint); err != nil {
				redis.Log.Error(err, "Failed to redis manual failover ",
					"namespace", redis.Namespace, "name", redis.Name,
					"master", master)

				return crStatus, err
			}
			redis.Log.Info("Success to redis manual failover, need to requeue.",
				"namespace", redis.Namespace, "name", redis.Name,
				"master", master)
			return crStatus, errors.New("need to requeue")

		}

		msg := fmt.Sprintf(UpdateMessageRedisCluster, redis.Name)
		redis.Recorder.Event(redis.HarborCluster, corev1.EventTypeNormal, RedisDownScaling, msg)

		redis.Log.Info(
			"Scaling replicas down",
			"from", actualReplica,
			"to", expectReplica,
		)

		msg = fmt.Sprintf(MessageRedisDownScaling, actualReplica, expectReplica)
		redis.Recorder.Event(redis.HarborCluster, corev1.EventTypeNormal, RedisDownScaling, msg)

		expectCR.ObjectMeta.SetResourceVersion(actualCR.ObjectMeta.GetResourceVersion())

		data, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&expectCR)
		if err != nil {
			return crStatus, err
		}
		redis.ExpectCR = &unstructured.Unstructured{Object: data}

		_, err = crdClient.Update(redis.ExpectCR, metav1.UpdateOptions{})
		if err != nil {
			return crStatus, err
		}
		crStatus.Phase = lcm.UpgradingPhase
	}
	return crStatus, nil
}

// calculateDownscales returns expected replicas and actual replicas.
func calculateDownscales(expected redisCli.RedisFailover, actual redisCli.RedisFailover) []redisDownscale {
	var downscales []redisDownscale
	actualReplicas := GetReplicas(actual)
	shouldExist := GetByName(actual.Name, expected.Name)
	expectedReplicas := int32(0)
	if shouldExist {
		expectedReplicas = GetReplicas(expected)
	}
	if expectedReplicas == 0 ||
		expectedReplicas < actualReplicas {
		downscales = append(downscales, redisDownscale{
			redisFailover:   expected,
			initialReplicas: actualReplicas,
			targetReplicas:  expectedReplicas,
		})
	}
	return downscales
}

// GetReplicas returns redis replicas.
func GetReplicas(redis redisCli.RedisFailover) int32 {
	if redis.Spec.Redis.Replicas != 0 {
		return redis.Spec.Redis.Replicas
	}
	return 0
}

func GetByName(actualName string, expectName string) bool {
	if expectName == actualName {
		return true
	}
	return false
}

type redisDownscale struct {
	redisFailover   redisCli.RedisFailover
	initialReplicas int32
	targetReplicas  int32
}

// leavingNodeNames returns the names of all nodes that should leave the cluster.
func leavingNodeNames(downscales []redisDownscale) []string {
	var leavingNodes []string
	for _, d := range downscales {
		leavingNodes = append(leavingNodes, d.leavingNodeNames()...)
	}
	return leavingNodes
}

// leavingNodeNames returns names of the nodes that are supposed to leave the Redis cluster.
func (d redisDownscale) leavingNodeNames() []string {
	if d.targetReplicas >= d.initialReplicas {
		return nil
	}
	leavingNodes := make([]string, 0, d.initialReplicas-d.targetReplicas)
	for i := d.initialReplicas - 1; i >= d.targetReplicas; i-- {
		leavingNodes = append(leavingNodes, PodName(d.redisFailover.Name, i))
	}
	return leavingNodes
}

// PodName returns the name of the pod with the given ordinal for this StatefulSet.
func PodName(name string, ordinal int32) string {
	return fmt.Sprintf("%s-%s-%d", "rfr", name, ordinal)
}

// isUpgradeResource returns whether the resource has been updated.
func isUpgradeResource(expectRedis *redisCli.RedisFailover, actualRedis *redisCli.RedisFailover) bool {
	expectEsResource := expectRedis.Spec.Redis.Resources
	actualEsResource := actualRedis.Spec.Redis.Resources

	if reflect.DeepEqual(expectEsResource, actualEsResource) {
		return false
	}

	return true
}

// IsMaster returns whether the IP is the Redis master node.
func IsMaster(ip string, password string) (bool, error) {
	options := &rediscli.Options{
		Addr:     fmt.Sprintf("%s:%s", ip, "6379"),
		Password: password,
		DB:       0,
	}
	rClient := rediscli.NewClient(options)

	defer rClient.Close()
	info, err := rClient.Info("replication").Result()
	if err != nil {
		return false, err
	}
	return strings.Contains(info, redisRoleMaster), nil
}

// Intersect returns the same value in both arrays.
func Intersect(nums1 []string, nums2 []string) []string {
	m := make(map[string]int)
	for _, v := range nums1 {
		m[v]++
	}

	var intersect []string
	for _, v := range nums2 {
		_, ok := m[v]
		if ok {
			intersect = append(intersect, v)
		}
	}
	return intersect
}

// ManualFailoverSentinel sends a sentinel manual failover the given sentinel
func ManualFailoverSentinel(ip string) error {
	options := &rediscli.Options{
		Addr: fmt.Sprintf("%s:%s", ip, RedisSentinelPort),
		DB:   0,
	}
	rClient := rediscli.NewClient(options)
	defer rClient.Close()
	cmd := rediscli.NewStringCmd("SENTINEL", "failover", "mymaster")
	rClient.Process(cmd)
	_, err := cmd.Result()
	if err != nil {
		return err
	}
	return nil
}

