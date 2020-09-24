# Deploy operator to your K8s clusters

## Environment

OS Debian 9, 8G mem 4 CPU

## Prerequisites

Kubernetes API running. Or create a local k8s using KIND.

### Minimal

- tools: git, go, kubectl, kustomize, helmv3, docker
- cert-manager
- harbor-operator

### Optional
- kind
- redis
- minio
- postgresql

## Install Tools

If you have installed these, ingore this step.

```shell script
apt-get install -y git

curl -LO https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl
mv ./kubectl /usr/local/bin/kubectl

curl -sSL https://github.com/kubernetes-sigs/kustomize/releases/download/kustomize%2Fv3.8.1/kustomize_v3.8.1_linux_amd64.tar.gz | tar -zx -C /usr/local/bin/ 

wget https://dl.google.com/go/go1.14.1.linux-amd64.tar.gz && tar -xzvf go1.14.1.linux-amd64.tar.gz -C /usr/local/ && ln -sfv /usr/local/go/bin/go /usr/local/bin/

curl https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3 | bash
```

install docker reference: https://docs.docker.com/engine/install/debian/

## Install KIND(Optional)

Before using kind, be sure docker has been installed.

```shell script
curl -Lo ./kind https://github.com/kubernetes-sigs/kind/releases/download/v0.7.0/kind-$(uname)-amd64
mv ./kind /usr/local/bin/kind
# create a local k8s in docker
kind create cluster
```

Deploy ingress-nginx-controller

```shell script
helm install nginx stable/nginx-ingress \
   --set-string 'controller.config.proxy-body-size'=0 \
   --set 'controller.service.type'=NodePort \
   --set 'controller.tolerations[0].key'=node-role.kubernetes.io/master \
   --set 'controller.tolerations[0].operator'=Equal \
   --set 'controller.tolerations[0].effect'=NoSchedule
```


## Install Cert-Manager
```shell script
helm repo add jetstack https://charts.jetstack.io
helm repo add bitnami https://charts.bitnami.com/bitnami

helm repo update

kubectl create namespace cert-manager

helm install cert-manager jetstack/cert-manager --namespace cert-manager --version v0.15.1 --set installCRDs=true
```

## Install Redis(Optional)

Deploy an external redis cluster in k8s for test. 
If you want harbor cluster operator to auto deploy a redis cluster, ignore it.

```shell script
helm install all-in-one-redis bitnami/redis --set usePassword=false
```

The helm chart will create a service `all-in-one-redis-master` (port is 6379), We can use it to visit redis service.


## Install MinIO(Optional)

Deploy an external minio cluster in k8s for test. 
If you want harbor cluster operator to auto deploy a redis cluster, or you already have an external storage service. ignore it.

```shell script
# deploy a minio cluster using helm chart. create a default bucket harbor（in default namespace）
helm install all-in-one-minio stable/minio --set defaultBucket.enabled=true,defaultBucket.bucketName=harbor
```

The helm chart will create a secret contains assesskey and accesssecret.

## Install PostgreSQL(Optional)

Deploy an external postgresql instance in k8s for test. 
If you want harbor cluster operator to auto deploy a postgresql instance, or you already have an external postgresql database. ignore it.

```shell script
helm install all-in-one-database bitnami/postgresql
```

The helm chart will create a secret which contains the password, We should use the password to generate a secret which contains datebase connection info.
```shell script
export PG_PASSWORD="$(kubectl get secret "all-in-one-database-postgresql" -o jsonpath='{.data.postgresql-password}' | base64 --decode)"

// creaet a new secret which contains database connnect info(host, port, database, username and password).
kubectl create secret generic "db-secret" \
   --from-literal host="all-in-one-database-postgresql.default" \
   --from-literal port='5432' \
   --from-literal database='postgres' \
   --from-literal username='postgres' \
   --from-literal password="$PG_PASSWORD"
```

## Install Harbor-Operator

Ref: https://github.com/goharbor/harbor-operator/blob/master/docs/installation/installation.md

## Install Harbor-Cluster-Operator

- clone source code

- build image, or use office image

```shell script
registry=
make docker-build IMG=${registry}/goharbor/harbor-cluster-operator:dev
make docker-push  IMG=${registry}/goharbor/harbor-cluster-operator:dev
```

- deploy
```shell script
kubectl create namespace harbor-cluster-operator-system
make deploy
```

## Deploy Sample

```yaml
---
# A secret of harbor admin password.
apiVersion: v1
kind: Secret
metadata:
  name: admin-secret
  namespace: default
data:
  password: MTIzNDU2
type: Opaque

# Using all external service.
---
apiVersion: goharbor.io/v1
kind: HarborCluster
metadata:
  name: harborcluster-sample
  namespace: default
spec:
  adminPasswordSecret: admin-secret
  certificateIssuerRef:
    name: selfsigned-issuer
  database:
    kind: external
    spec:
      resources: {}
      # we create it when install postgresql
      secretName: db-secret
  jobService:
    replicas: 1
    workerCount: 1
  publicURL: https://test.harbor.local.com
  redis:
    kind: external
    spec:
      hosts:
      - host: all-in-one-redis-master.default
        port: "6379"
      poolSize: 10
      schema: redis
  replicas: 1
  storage:
    kind: s3
    s3:
      accesskey: AKIAIOSFODNN7EXAMPLE
      bucket: harbor
      chunksize: "5242880"
      region: us-west-1
      regionendpoint: http://all-in-one-minio.default.svc.cluster.local:9000
      secretkey: wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
      storageclass: STANDARD
  version: 1.10.0
```