package cache

import (
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	redisCli "src/github.com/goharbor/harbor-cluster-operator/controllers/cache/client/api/v1"
)

func (d *defaultCache) Deploy() error {
	var actualCR redisCli.RedisFailover
	var expectCR redisCli.RedisFailover
	name := d.Request.Name
	nameSpace := d.Request.Namespace
	crdClient := d.DClient.Resource(virtualServiceGVR)

	uData, err := d.generateRedisCR()
	if err != nil {
		return err
	}
	if err := controllerutil.SetControllerReference(d.Harbor, uData, d.Scheme); err != nil {
		return err
	}

	resp, err := crdClient.Namespace(nameSpace).Get(name, metav1.GetOptions{})
	if err != nil && errors.IsNotFound(err) {

		d.Log.Info("Creating Redis.", "namespace", nameSpace, "name", name)
		_, err = crdClient.Namespace(nameSpace).Create(uData, metav1.CreateOptions{})
		if err != nil {
			return err
		}

		d.Log.Info("Redis create complete.", "namespace", nameSpace, "name", name)
	} else if err != nil {
		return err
	} else {
		if err := runtime.DefaultUnstructuredConverter.
			FromUnstructured(resp.UnstructuredContent(), &actualCR); err != nil {
			return err
		}

		if err := runtime.DefaultUnstructuredConverter.
			FromUnstructured(uData.UnstructuredContent(), &expectCR); err != nil {
			return err
		}

		d.ExpectCR = &expectCR
		d.ActualCR = &actualCR
	}

	return nil
}
