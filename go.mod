module github.com/goharbor/harbor-cluster-operator

go 1.13

require (
	github.com/go-logr/logr v0.2.1-0.20200730175230-ee2de8da5be6
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/goharbor/harbor-operator v0.5.2-0.20200817115335-b421dca2f798
	github.com/google/go-cmp v0.5.1
	github.com/jackc/pgx/v4 v4.8.1
	github.com/jetstack/cert-manager v0.16.1
	github.com/minio/minio-go/v6 v6.0.57
	github.com/minio/operator v0.0.0-20200814200655-60bf757aac60
	github.com/onsi/ginkgo v1.12.1
	github.com/onsi/gomega v1.10.1
	github.com/spotahome/redis-operator v1.0.0
	github.com/zalando/postgres-operator v1.5.0
	k8s.io/api v0.19.0-rc.3
	k8s.io/apimachinery v0.19.0-rc.3
	k8s.io/client-go v11.0.0+incompatible
	sigs.k8s.io/controller-runtime v0.6.1-0.20200804124940-17eebbff0d48
)

replace k8s.io/client-go v11.0.0+incompatible => k8s.io/client-go v0.0.0-20200813012017-e7a1d9ada0d5
