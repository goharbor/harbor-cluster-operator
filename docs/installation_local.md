# Deploy Harbor with cluster operator on kind cluster

**NOTES:** 
>Currently, we do not have a proper package to deploy the cluster operator as well its dependant operators
with a simple and unified way. So, the guideline shown here is a temporary solution. We're working on to provide a formal
deployment solution latter.

## Prerequisites

* A VM with linux OS (MEM: 4G+, DISK: 50GB+)
* [Docker](https://docs.docker.com/engine/install/) installed (Version: v19.03.12+)
* [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) installed (Version: v1.18+)
* [kustomize](https://sigs.k8s.io/kustomize/docs/INSTALL.md) installed (Version: v3.1.0+)
* [kind](https://kind.sigs.k8s.io/docs/user/quick-start/) installed (Version: v0.8.1+)

## Create kind cluster

Use the following configuration to create a kind cluster with multiple worker nodes:

```yaml
# a cluster with 3 control-plane nodes and 3 workers
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
  kubeadmConfigPatches:
  - |
    kind: InitConfiguration
    nodeRegistration:
      kubeletExtraArgs:
        node-labels: "ingress-ready=true"
  extraPortMappings:
  - containerPort: 80
    hostPort: 80
    protocol: TCP
  - containerPort: 443
    hostPort: 443
    protocol: TCP
- role: worker
- role: worker
- role: worker
```

Execute command:

```shell script
kind create cluster --name myk8s --config kind.yaml
```

Check kind cluster:

```shell script
kubectl cluster-info
```

>Kubernetes master is running at https://127.0.0.1:43415
 KubeDNS is running at https://127.0.0.1:43415/api/v1/namespaces/kube-system/services/kube-dns:dns/proxy
 
>To further debug and diagnose cluster problems, use 'kubectl cluster-info dump'.

## Deploy dependant components

### Ingress controller

Deploy nginx ingress controller with the command shown below:

```shell script
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/master/deploy/static/provider/kind/deploy.yaml
```

Check if it is ready:

```shell script
kubectl wait --namespace ingress-nginx \
  --for=condition=ready pod \
  --selector=app.kubernetes.io/component=controller \
  --timeout=90s
```

As optional steps, you can try to deploy sample apps and access the ingress routes.

Deploy sample apps:
```shell script
kubectl apply -f https://kind.sigs.k8s.io/examples/ingress/usage.yaml
```
Verify that the ingress works:
```shell script
# should output "foo"
curl localhost/foo
# should output "bar"
curl localhost/bar
```

### Cert-manager

Follow the guide shown [here](https://cert-manager.io/docs/installation/kubernetes/) to deploy the cert-manager into the kind cluster.

## Deploy operator

### Deploy dependant operators

#### PostgreSQL

Follow the installation guide shown [here](https://github.com/zalando/postgres-operator/blob/master/docs/quickstart.md#configuration-options) to install the PostgreSQL operator.
It can be used with kubectl 1.14 or newer as easy as:
```shell script
kubectl apply -k github.com/zalando/postgres-operator/manifests
```

#### Redis

Follow the deployment guide shown [here](https://github.com/spotahome/redis-operator#operator-deployment-on-kubernetes) to deploy Redis operator to the kind cluster.

A simple way is:
```shell script
kubectl create -f https://raw.githubusercontent.com/spotahome/redis-operator/master/example/operator/all-redis-operator-resources.yaml
```

**NOTES:**
> If encounter RBAC permission issue, as a simple way, upgrade the privileges of the default service account under the
>deploying namespace. Try:

```shell script
# <NAMESPACE> is PLACEHOLDER

kubectl create clusterrolebinding --clusterrole=cluster-admin  --user=system:serviceaccount:<NAMESPACE>:default --clusterrole=cluster-admin --user=system:serviceaccount rds-admin-binding
```

#### Storage(minIO)

**NOTES:**
> The cluster operator is rely on a lower version of minIO operator and we cannot follow the regualr guide to install minIO
>operator. To make compatible with current cluster operator, here we'll install it from the source code.
>
>The related fixing work has been started.

Clone the repo:

```shell script
git clone https://github.com/minio/operator
```

Deploy the operator:

```shell script
kustomize build | kubectl apply -f -
```

### Deploy Harbor core operator

Deploy core operator from source code.

Clone the repo:

```shell script
https://github.com/goharbor/harbor-operator.git
```

Build the controller image:

```shell script
# cd harbor-operator
# IMG ?= goharbor/harbor-operator:dev
make docker-build
```

Load the image into kind cluster nodes:

```shell script
# my k8s is cluster name
kind load --name myk8s docker-image goharbor/harbor-operator:dev
```

Deploy the operator:

```shell script
make deploy
```

### Deploy Harbor cluster operator

Deploy cluster operator from source code.

Clone the repo:

```shell script
git clone https://github.com/goharbor/harbor-cluster-operator.git
```

Build the controller image:

```shell script
# cd harbor-operator
export IMG=goharbor/harbor-cluster-operator:dev

make docker-build
```

Load the image into kind cluster nodes:

```shell script
# my k8s is cluster name
kind load --name myk8s docker-image goharbor/harbor-cluster-operator:dev
```

Deploy the operator:

```shell script
make deploy
```

## Deploy Harbor
Create a `sample` namespace:

```shell script
kubectl create ns sample
```

Create an admin secret with the manifest like:

```shell script
cat <<EOF | kubectl apply -f -
# A secret of harbor admin password.
# Password is encoded with base64.
apiVersion: v1
kind: Secret
metadata:
  name: admin-secret
  namespace: sample
data:
  password: SGFyYm9yMTIzNDU=
type: Opaque
EOF
```

Create a self-signed issuer:

```shell script
cat <<EOF | kubectl apply -f -
apiVersion: cert-manager.io/v1alpha2
kind: Issuer
metadata:
  name: selfsigned-issuer
  namespace: sample
spec:
  selfSigned: {}
EOF
```

Here is a sample manifest for deploying a Harbor with all in-cluster services(use `sample.goharbor.io` as public URL 
and `notary.goharbor.io` as Notary public URL):

```shell script
cat <<EOF | kubectl apply -f -
apiVersion: goharbor.io/v1
kind: HarborCluster
metadata:
  name: sz-harbor-cluster
  namespace: sample
spec:
  redis:
    kind: "inCluster"
    spec:
      server:
        replicas: 1
        resources:
          requests:
            cpu: "1"
            memory: "2Gi"
        storage: "10Gi"
      sentinel:
        replicas: 1
      schema: "redis"
  adminPasswordSecret: "admin-secret"
  certificateIssuerRef:
    name: selfsigned-issuer
  tlsSecret: public-certificate
  database:
    kind: "inCluster"
    spec:
      replicas: 2
      resources:
        requests:
          cpu: "1"
          memory: "2Gi"
        limits:
          cpu: "1"
          memory: "2Gi"
  publicURL: "https://sample.goharbor.io"
  replicas: 2
  notary:
    publicUrl: "https://notary.goharbor.io"
  disableRedirect: true
  jobService:
    workerCount: 10
    replicas: 2
  storage:
    kind: "inCluster"
    options:
      provider: minIO
      spec:
        replicas: 4
        version: RELEASE.2020-01-03T19-12-21Z
        volumeClaimTemplate:
          spec:
            storageClassName: standard
            accessModes:
              - ReadWriteOnce
            resources:
              requests:
                storage: 10Gi
        resources:
          requests:
            memory: 1Gi
            cpu: 500m
          limits:
            memory: 1Gi
            cpu: 1000m
  version: 1.10.0
EOF
```

After a while, the harbor cluster (`HarborCluster`) should be ready:

```shell script
kubectl get HarborCluster -n sample -o wide
```

you can get the output:
```shell script
NAME                VERSION   PUBLIC URL                   SERVICE READY   CACHE READY   DATABASE READY   STORAGE READY
sz-harbor-cluster   1.10.0    https://sample.goharbor.io   Unknown         True          True             True
```

## Post actions

As an easy and quick way, add host mappings into the `/ect/hosts` of the host that's used to access the new deployed Harbor.

```shell script
<KIND_HOST_IP> sample.goharbor.io 
<KIND_HOST_IP> notary.goharbor.io 
```

Try the API server first:

```shell script
curl -k https://sample.goharbor.io/api/systeminfo
```
There will be some JSON data output like:
```json
{"with_notary":true,"with_admiral":false,"admiral_endpoint":"NA","auth_mode":"db_auth","registry_url":"sample.goharbor.io","external_url":"https://sample.goharbor.io","project_creation_restriction":"everyone","self_registration":false,"has_ca_root":false,"harbor_version":"v1.10.0-6b84b62f","registry_storage_provider_name":"memory","read_only":false,"with_chartmuseum":false,"notification_enable":true}
```

Try to push images:

```
docker login sample.goharbor.io -u admin -p <PASSWORD>

docker tag nginx:latest sample.goharbor.io/library/nginx:latest

docker push sample.goharbor.io/library/nginx:latest
```

Open browser and navigate to `https://sample.goharbor.io` to open web UI of Harbor.
