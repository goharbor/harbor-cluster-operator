package cache

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	//redisCli "src/github.com/goharbor/harbor-cluster-operator/controllers/cache/client/api/v1"
	redisCli "github.com/spotahome/redis-operator/api/redisfailover/v1"
	"src/github.com/goharbor/harbor-cluster-operator/controllers/utils"
)

var (
	virtualServiceGVR = schema.GroupVersionResource{
		Group:    "databases.spotahome.com",
		Version:  "v1",
		Resource: "redisfailovers",
	}
)

func (d *defaultCache) generateRedisCR() (*unstructured.Unstructured, error) {
	name := d.Request.Name
	nameSpace := d.Request.Namespace

	resource := utils.GetRequests("1", "1Gi")
	conf := &redisCli.RedisFailover{
		TypeMeta: v1.TypeMeta{
			Kind:       "RedisFailover",
			APIVersion: "databases.spotahome.com/v1",
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      name,
			Namespace: nameSpace,
		},
		Spec: redisCli.RedisFailoverSpec{
			Redis: redisCli.RedisSettings{
				Replicas: 3,
				Resources: corev1.ResourceRequirements{
					Requests: resource,
					Limits:   resource,
				},
			},
			Sentinel: redisCli.SentinelSettings{
				Replicas: 3,
				Resources: corev1.ResourceRequirements{
					Requests: resource,
					Limits:   resource,
				},
			},
			Auth: redisCli.AuthSettings{SecretPath: name},
		},
	}

	mapResult, err := runtime.DefaultUnstructuredConverter.ToUnstructured(conf)
	if err != nil {
		return nil, err
	}
	data := unstructured.Unstructured{Object: mapResult}

	return &data, nil
}

func (d *defaultCache) generateRedisSecret(labels map[string]string) *corev1.Secret {
	name := d.Request.Name
	namespace := d.Request.Namespace

	labels = MergeLabels(labels, generateLabels(RoleName, name))

	passStr := RandomString(8, "a")

	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    labels,
		},
		StringData: map[string]string{
			"password": passStr,
		},
	}
}
