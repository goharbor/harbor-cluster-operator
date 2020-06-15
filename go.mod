module github.com/goharbor/harbor-cluster-operator

go 1.13

require (
	github.com/go-logr/logr v0.1.0
	github.com/go-redis/redis v6.15.8+incompatible
	github.com/goharbor/harbor-operator v0.5.0
	github.com/jackc/pgx/v4 v4.6.0
	github.com/jetstack/cert-manager v0.14.2
	github.com/minio/minio-operator v0.0.0-20200515184103-acad2ca18729
	github.com/onsi/ginkgo v1.12.0
	github.com/onsi/gomega v1.9.0
	github.com/spotahome/redis-operator v1.0.0
	github.com/zalando/postgres-operator v1.5.0
	k8s.io/api v0.18.2
	k8s.io/apimachinery v0.18.2
	k8s.io/client-go v11.0.0+incompatible
	sigs.k8s.io/controller-runtime v0.6.0
)

replace k8s.io/client-go v11.0.0+incompatible => k8s.io/client-go v0.18.2
