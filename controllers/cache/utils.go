package cache

import (
	"bytes"
	"context"
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
func (d *defaultCache) GetRedisName() string {
	return generateName(ReidsType, d.Harbor.Name)
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

// getRedisPassword is get redis password
func (d *defaultCache) GetRedisPassword() (string, error) {
	var redisPassWord string
	redisPassMap, err := d.GetRedisSecret()
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

// GetRedisSecret get the Redis Password Secret
func (d *defaultCache) GetRedisSecret() (map[string][]byte, error) {
	secret := &corev1.Secret{}
	name := d.Request.Name
	namespace := d.Request.Namespace

	err := d.Client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, secret)
	if err != nil {
		return nil, err
	}
	opts := &client.ListOptions{}
	set := labels.SelectorFromSet(secret.Labels)
	opts.LabelSelector = set

	sc := &corev1.SecretList{}
	err = d.Client.List(context.TODO(), sc, opts)
	if err != nil {
		return nil, err
	}
	var redisPw map[string][]byte
	for _, rp := range sc.Items {
		redisPw = rp.Data
	}

	return redisPw, nil
}

func (d *defaultCache) GetDeploymentPods() (*appsv1.Deployment, *corev1.PodList, error) {
	deploy := &appsv1.Deployment{}
	name := fmt.Sprintf("%s-%s", "rfs", d.Request.Name)
	namespace := d.Request.Namespace
	fmt.Println(name, namespace)
	err := d.Client.Get(d.CTX, types.NamespacedName{Name: name, Namespace: namespace}, deploy)
	if err != nil {
		return nil, nil, err
	}

	opts := &client.ListOptions{}
	set := labels1.SelectorFromSet(deploy.Spec.Selector.MatchLabels)
	opts.LabelSelector = set

	pod := &corev1.PodList{}
	err = d.Client.List(d.CTX, pod, opts)
	if err != nil {
		d.Log.Error(err, "fail to get pod.", "namespace", namespace, "name", name)
		return nil, nil, err
	}
	return deploy, pod, nil
}

func (d *defaultCache) GetStatefulSetPods() (*appsv1.StatefulSet, *corev1.PodList, error) {
	sts := &appsv1.StatefulSet{}
	name := fmt.Sprintf("%s-%s", "rfr", d.Request.Name)
	namespace := d.Request.Namespace
	fmt.Println(name, namespace)
	err := d.Client.Get(d.CTX, types.NamespacedName{Name: name, Namespace: namespace}, sts)
	if err != nil {
		return nil, nil, err
	}

	opts := &client.ListOptions{}
	set := labels1.SelectorFromSet(sts.Spec.Selector.MatchLabels)
	opts.LabelSelector = set

	pod := &corev1.PodList{}
	err = d.Client.List(d.CTX, pod, opts)
	if err != nil {
		d.Log.Error(err, "fail to get pod.", "namespace", namespace, "name", name)
		return nil, nil, err
	}
	return sts, pod, nil
}

func (d *defaultCache) GetServiceUrl(pods []corev1.Pod) string {
	var url string
	_, err := rest.InClusterConfig()
	if err != nil {
		randomPod := pods[rand.Intn(len(pods))]
		url = randomPod.Status.PodIP
	} else {
		url = fmt.Sprintf("%s-%s.svc", "rfs", d.Request.Name)
	}

	return url
}
