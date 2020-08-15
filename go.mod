module github.com/goharbor/harbor-cluster-operator

go 1.13

require (
	github.com/go-logr/logr v0.2.1-0.20200730175230-ee2de8da5be6
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/goharbor/harbor-operator v0.5.1
	github.com/google/go-cmp v0.5.1
	github.com/jackc/pgx/v4 v4.8.1
	github.com/jessevdk/go-flags v1.4.0 // indirect
	github.com/jetstack/cert-manager v0.16.1
	github.com/minio/minio-go/v6 v6.0.57
	github.com/minio/operator v0.0.0-20200814200655-60bf757aac60
	github.com/spotahome/redis-operator v1.0.0
	github.com/zalando/postgres-operator v1.5.0
	k8s.io/api v0.19.0-rc.3
	k8s.io/apimachinery v0.19.0-rc.3
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/kube-openapi v0.0.0-20200805222855-6aeccd4b50c6 // indirect
	k8s.io/utils v0.0.0-20200729134348-d5654de09c73 // indirect
	sigs.k8s.io/controller-runtime v0.6.1-0.20200804124940-17eebbff0d48
)

replace k8s.io/client-go v11.0.0+incompatible => k8s.io/client-go v0.0.0-20200813012017-e7a1d9ada0d5
