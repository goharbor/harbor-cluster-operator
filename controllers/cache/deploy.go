package cache

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
	"log"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (d *defaultCache) Deploy() error {
	var actualCR *unstructured.Unstructured
	var expectCR *unstructured.Unstructured
	name := d.Request.Name
	nameSpace := d.Request.Namespace
	crdClient := d.DClient.Resource(virtualServiceGVR)

	expectCR, err := d.generateRedisCR()
	if err != nil {
		return err
	}
	if err := controllerutil.SetControllerReference(d.Harbor, expectCR, d.Scheme); err != nil {
		return err
	}

	actualCR, err = crdClient.Namespace(nameSpace).Get(name, metav1.GetOptions{})
	if err != nil && errors.IsNotFound(err) {

		if err := d.DeploySecret(); err != nil {
			return err
		}

		d.Log.Info("Creating Redis.", "namespace", nameSpace, "name", name)
		_, err = crdClient.Namespace(nameSpace).Create(expectCR, metav1.CreateOptions{})
		if err != nil {
			return err
		}

		d.Log.Info("Redis create complete.", "namespace", nameSpace, "name", name)
	} else if err != nil {
		return err
	} else {
		d.ExpectCR = expectCR
		d.ActualCR = actualCR
	}

	return nil
}

// DeploySecret deploy the Redis Password Secret
func (d *defaultCache) DeploySecret() error {
	secret := &corev1.Secret{}
	name := d.Request.Name
	namespace := d.Request.Namespace
	sc := d.generateRedisSecret(d.Labels)
	if err := controllerutil.SetControllerReference(d.Harbor, sc, d.Scheme); err != nil {
		return err
	}
	err := d.Client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, secret)
	if err != nil && errors.IsNotFound(err) {
		log.Printf("Creating Redis Password Secret %s/%s\n", namespace, namespace)
		err = d.Client.Create(context.TODO(), sc)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	return nil
}
