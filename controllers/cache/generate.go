package cache

import (
	redisCli "github.com/spotahome/redis-operator/api/redisfailover/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	virtualServiceGVR = schema.GroupVersionResource{
		Group:    "databases.spotahome.com",
		Version:  "v1",
		Resource: "redisfailovers",
	}
)

// generateRedisCR returns RedisFailovers CRs
func (redis *RedisReconciler) generateRedisCR() (*unstructured.Unstructured, error) {
	resource := redis.GetRedisResource()
	redisRep := redis.GetHRedisServerReplica()
	sentinelRep := redis.GetHRedisSentinelReplica()
	storageSize := redis.HarborCluster.Spec.Redis.Spec.Server.Storage

	conf := &redisCli.RedisFailover{
		TypeMeta: v1.TypeMeta{
			Kind:       "RedisFailover",
			APIVersion: "databases.spotahome.com/v1",
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      redis.Name,
			Namespace: redis.Namespace,
		},
		Spec: redisCli.RedisFailoverSpec{
			Redis: redisCli.RedisSettings{
				Replicas: redisRep,
				Resources: corev1.ResourceRequirements{
					Requests: resource,
					Limits:   resource,
				},
			},
			Sentinel: redisCli.SentinelSettings{
				Replicas: sentinelRep,
				Resources: corev1.ResourceRequirements{
					Requests: resource,
					Limits:   resource,
				},
			},
			Auth: redisCli.AuthSettings{SecretPath: redis.Name},
		},
	}

	if redis.HarborCluster.Spec.Redis.Spec.Server.Storage != "" {
		conf.Spec.Redis.Storage.PersistentVolumeClaim = generateRedisStorage(storageSize, redis.Name, redis.Labels)
	}

	mapResult, err := runtime.DefaultUnstructuredConverter.ToUnstructured(conf)
	if err != nil {
		return nil, err
	}
	data := unstructured.Unstructured{Object: mapResult}

	return &data, nil
}

//generateRedisSecret returns redis password secret
func (redis *RedisReconciler) generateRedisSecret(labels map[string]string) *corev1.Secret {
	labels = MergeLabels(labels, generateLabels(RoleName, redis.Name))

	passStr := RandomString(8, "a")

	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      redis.Name,
			Namespace: redis.Namespace,
			Labels:    labels,
		},
		StringData: map[string]string{
			"password": passStr,
		},
	}
}

func generateRedisStorage(size, name string, labels map[string]string) *corev1.PersistentVolumeClaim {
	storage, _ := resource.ParseQuantity(size)
	return &corev1.PersistentVolumeClaim{
		ObjectMeta: v1.ObjectMeta{
			Name:   name,
			Labels: labels,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{
				corev1.ReadWriteOnce,
			},
			Selector: nil,
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					"storage": storage,
				},
			},
		},
	}
}
