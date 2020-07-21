package harbor

import (
	"github.com/goharbor/harbor-cluster-operator/controllers/k8s"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	certv1 "github.com/jetstack/cert-manager/pkg/apis/certmanager/v1alpha2"
)

// getRegistryIssuerNamespacedName returns CertificateIssuerRef name and namespace
func (harbor *HarborReconciler) getRegistryIssuerNamespacedName() types.NamespacedName {
	return types.NamespacedName{
		Namespace: harbor.HarborCluster.Namespace,
		Name:      harbor.HarborCluster.Spec.CertificateIssuerRef.Name,
	}
}

// CheckIssuer check issuer has exist, if not will create issuer
func (harbor *HarborReconciler) CheckIssuer() error {
	var issuer certv1.Issuer
	namespacedName := harbor.getRegistryIssuerNamespacedName()

	err := harbor.Get(namespacedName, &issuer)
	if err != nil {
		if errors.IsNotFound(err) {
			return provisionIssuer(harbor.Client, namespacedName.Name, namespacedName.Namespace)
		}
	}
	return err
}

// provisionIssuer will create issuer for registry
func provisionIssuer(client k8s.Client, name, namespace string) error {
	issuer := newIssuer(name, namespace)
	return client.Create(issuer)
}

// newIssuer returns new Issuer object
func newIssuer(name, namespace string) *certv1.Issuer {
	return &certv1.Issuer{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: certv1.IssuerSpec{
			IssuerConfig: certv1.IssuerConfig{
				SelfSigned: &certv1.SelfSignedIssuer{},
			},
		},
	}
}
