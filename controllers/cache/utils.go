package cache

import (
	"bytes"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	labels1 "k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	"math/rand"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
	"time"
)

const (
	ReidsType    = "rfr"
	SentinelType = "rfs"

	RoleName          = "harbor-cluster"
	RedisSentinelPort = "26379"
)

// GetRedisName returns the name for redis resources
func (redis *RedisReconciler) GetRedisName() string {
	return generateName(ReidsType, redis.GetHarborClusterName())
}

func generateName(typeName, metaName string) string {
	return fmt.Sprintf("%s-%s", typeName, metaName)
}

func RandomString(randLength int, randType string) (result string) {
	var num string = "0123456789"
	var lower string = "abcdefghijklmnopqrstuvwxyz"
	var upper string = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	b := bytes.Buffer{}
	if strings.Contains(randType, "0") {
		b.WriteString(num)
	}
	if strings.Contains(randType, "a") {
		b.WriteString(lower)
	}
	if strings.Contains(randType, "A") {
		b.WriteString(upper)
	}
	var str = b.String()
	var strLen = len(str)
	if strLen == 0 {
		result = ""
		return
	}

	rand.Seed(time.Now().UnixNano())
	b = bytes.Buffer{}
	for i := 0; i < randLength; i++ {
		b.WriteByte(str[rand.Intn(strLen)])
	}
	result = b.String()
	return
}

// GetRedisPassword is get redis password
func (redis *RedisReconciler) GetRedisPassword() (string, error) {
	var redisPassWord string
	redisPassMap, err := redis.GetRedisSecret()
	if err != nil {
		return "", err
	}
	for k, v := range redisPassMap {
		if k == "password" {
			redisPassWord = string(v)
			return redisPassWord, nil
		}
	}
	return redisPassWord, nil
}

// GetRedisSecret returns the Redis Password Secret
func (redis *RedisReconciler) GetRedisSecret() (map[string][]byte, error) {
	secret := &corev1.Secret{}

	err := redis.Client.Get(types.NamespacedName{Name: redis.Name, Namespace: redis.Namespace}, secret)
	if err != nil {
		return nil, err
	}
	opts := &client.ListOptions{}
	set := labels.SelectorFromSet(secret.Labels)
	opts.LabelSelector = set

	sc := &corev1.SecretList{}
	err = redis.Client.List(opts, sc)
	if err != nil {
		return nil, err
	}
	var redisPw map[string][]byte
	for _, rp := range sc.Items {
		redisPw = rp.Data
	}

	return redisPw, nil
}

// GetDeploymentPods returns the Redis Sentinel pod list
func (redis *RedisReconciler) GetDeploymentPods() (*appsv1.Deployment, *corev1.PodList, error) {
	deploy := &appsv1.Deployment{}
	name := fmt.Sprintf("%s-%s", "rfs", redis.Name)

	err := redis.Client.Get(types.NamespacedName{Name: name, Namespace: redis.Namespace}, deploy)
	if err != nil {
		return nil, nil, err
	}

	opts := &client.ListOptions{}
	set := labels1.SelectorFromSet(deploy.Spec.Selector.MatchLabels)
	opts.LabelSelector = set

	pod := &corev1.PodList{}
	err = redis.Client.List(opts, pod)
	if err != nil {
		redis.Log.Error(err, "fail to get pod.", "namespace", redis.Namespace, "name", name)
		return nil, nil, err
	}
	return deploy, pod, nil
}

// GetStatefulSetPods returns the Redis Server pod list
func (redis *RedisReconciler) GetStatefulSetPods() (*appsv1.StatefulSet, *corev1.PodList, error) {
	sts := &appsv1.StatefulSet{}
	name := fmt.Sprintf("%s-%s", "rfr", redis.Name)

	err := redis.Client.Get(types.NamespacedName{Name: name, Namespace: redis.Namespace}, sts)
	if err != nil {
		return nil, nil, err
	}

	opts := &client.ListOptions{}
	set := labels1.SelectorFromSet(sts.Spec.Selector.MatchLabels)
	opts.LabelSelector = set

	pod := &corev1.PodList{}
	err = redis.Client.List(opts, pod)
	if err != nil {
		redis.Log.Error(err, "fail to get pod.", "namespace", redis.Namespace, "name", name)
		return nil, nil, err
	}
	return sts, pod, nil
}

// GetServiceUrl returns the Redis Sentinel pod ip or service name
func (redis *RedisReconciler) GetServiceUrl(pods []corev1.Pod) string {
	var url string
	_, err := rest.InClusterConfig()
	if err != nil {
		randomPod := pods[rand.Intn(len(pods))]
		url = randomPod.Status.PodIP
	} else {
		url = fmt.Sprintf("%s-%s.svc", "rfs", redis.GetHarborClusterName())
	}

	return url
}

// GetHarborClusterName returns harbor cluster name
func (redis *RedisReconciler) GetHarborClusterName() string {
	return redis.HarborCluster.Name
}

// GetHarborClusterNamespace returns harbor cluster namespace
func (redis *RedisReconciler) GetHarborClusterNamespace() string {
	return redis.HarborCluster.Namespace
}

// GetRedisResource returns redis resource
func (redis *RedisReconciler) GetRedisResource() corev1.ResourceList {
	resources := corev1.ResourceList{}
	cpu := redis.HarborCluster.Spec.Redis.Spec.Server.Resources.Cpu()
	mem := redis.HarborCluster.Spec.Redis.Spec.Server.Resources.Memory()
	if cpu != nil {
		resources[corev1.ResourceCPU] = *cpu
	}
	if mem != nil {
		resources[corev1.ResourceMemory] = *mem
	}
	return resources
}

// GetHRedisServerReplica returns redis server replicas
func (redis *RedisReconciler) GetHRedisServerReplica() int32 {
	return int32(redis.HarborCluster.Spec.Redis.Spec.Server.Replicas)
}

// GetHRedisSentinelReplica returns redis sentinel replicas
func (redis *RedisReconciler) GetHRedisSentinelReplica() int32 {
	return int32(redis.HarborCluster.Spec.Redis.Spec.Sentinel.Replicas)
}

