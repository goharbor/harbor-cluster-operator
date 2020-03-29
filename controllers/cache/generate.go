package cache

import (
	"encoding/json"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	redisCli "src/github.com/goharbor/harbor-cluster-operator/controllers/cache/client/api/v1"
	"src/github.com/goharbor/harbor-cluster-operator/controllers/utils"
)

var (
	virtualServiceGVR = schema.GroupVersionResource{
		Group:    "databases.spotahome.com",
		Version:  "v1",
		Resource: "redisfailovers",
	}

	groupVersionKind = schema.GroupVersionKind{
		Group:   "databases.spotahome.com",
		Version: "v1",
		Kind:    "RedisFailover",
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
		},
	}

	var mapResult map[string]interface{}
	confBytes, _ := json.Marshal(conf)
	if err := json.Unmarshal(confBytes, &mapResult); err != nil {
		return nil, err
	}
	data := unstructured.Unstructured{Object: mapResult}

	return &data, nil
}
