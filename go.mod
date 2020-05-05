module github.com/goharbor/harbor-cluster-operator

go 1.13

require (
	github.com/go-logr/logr v0.1.0
	github.com/jetstack/cert-manager v0.14.2
	github.com/minio/minio-operator v1.0.9
	github.com/onsi/ginkgo v1.12.0
	github.com/onsi/gomega v1.9.0
	k8s.io/api v0.18.2
	k8s.io/apimachinery v0.18.2
	k8s.io/client-go v11.0.0+incompatible
	sigs.k8s.io/controller-runtime v0.6.0
)
