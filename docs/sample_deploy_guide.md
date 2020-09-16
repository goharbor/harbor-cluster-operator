# Guide for deploying sample

## Prerequisite

Make sure the operators have been deployed to the target K8s cluster by following the [installation guide](./installation.md).

## Create deployment manifest

Refer the `HarborCluster` [spec](./cr_HarborCluster_spec.md) to create a deployment manifest. e.g.(including namespace and secrets):

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: harbor
---
# A secret of harbor admin password.
apiVersion: v1
kind: Secret
metadata:
  name: admin-secret
  namespace: harbor
data:
  password: SGFyYm9yMTIzNDU=
type: Opaque
---
apiVersion: cert-manager.io/v1alpha2
kind: Issuer
metadata:
  name: selfsigned-issuer
  namespace: sample
spec:
  selfSigned: {}
---
apiVersion: cert-manager.io/v1alpha2
kind: Certificate
metadata:
  name: public-certificate
  namespace: sample
spec:
  secretName: public-certificate
  dnsNames:
  - sample.goharbor.io
  - notary.goharbor.io
  issuerRef:
    name: selfsigned-issuer
    kind: Issuer
---
apiVersion: goharbor.io/v1
kind: HarborCluster
metadata:
  name: sz-harbor-cluster
  namespace: harbor
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
  publicURL: "https://harbor.goharbor.io"
  replicas: 3
  notary:
    publicUrl: "https://notary-harbor.goharbor.io"
  disableRedirect: true
  jobService:
    workerCount: 10
    replicas: 3
  chartMuseum:
    absoluteURL: true
  clair:
    updateInterval: 10
    vulnerabilitySources:
    - ubuntu
    - alphine
  storage:
    kind: "inCluster"
    options:
      provider: minIO
      spec:
        replicas: 3
        volumesPerServer: 2
        version: RELEASE.2020-08-13T02-39-50Z
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
  version: 1.10.4
```

## Deploy

Apply the manifest to the K8s cluster:

```shell script
kubectl apply -f your_manifest.yaml
```

Check the status of the `HarborCluster` CR:

```shell script
kubectl get harborcluster -n your_namespace -o wide

# Output
# NAME                VERSION   PUBLIC URL                   SERVICE READY   CACHE READY   DATABASE READY   STORAGE READY
# sz-harbor-cluster   1.10.0    https://harbor.goharbor.io   Unknown         True          True             True
```

Check the related K8s deployments/pods/statefulsets owned by the `HarborCluster` CR:

```shell script
kubectl get all -n your_namespace
```

Output sample:

```shell script
NAME                                                         READY   STATUS    RESTARTS   AGE
pod/harbor-sz-harbor-cluster-0                               1/1     Running   0          4d18h
pod/harbor-sz-harbor-cluster-1                               1/1     Running   0          4d18h
pod/rfr-sz-harbor-cluster-0                                  1/1     Running   0          4d18h
pod/rfs-sz-harbor-cluster-bd4bdcdcf-cqz5j                    1/1     Running   0          4d18h
pod/sz-harbor-cluster-harbor-chartmuseum-68d7fc88bb-42w7l    1/1     Running   0          4d18h
pod/sz-harbor-cluster-harbor-chartmuseum-68d7fc88bb-tzzm6    1/1     Running   0          4d18h
pod/sz-harbor-cluster-harbor-chartmuseum-68d7fc88bb-xbr8s    1/1     Running   0          4d18h
pod/sz-harbor-cluster-harbor-clair-5f65c876f8-5zfnq          2/2     Running   0          4d18h
pod/sz-harbor-cluster-harbor-clair-5f65c876f8-crcnf          2/2     Running   0          4d18h
pod/sz-harbor-cluster-harbor-clair-5f65c876f8-zzjq6          2/2     Running   0          4d18h
pod/sz-harbor-cluster-harbor-core-86cb58b8dd-dfmd2           1/1     Running   0          4d18h
pod/sz-harbor-cluster-harbor-core-86cb58b8dd-hrd4s           1/1     Running   0          4d18h
pod/sz-harbor-cluster-harbor-core-86cb58b8dd-tc6xk           1/1     Running   0          4d18h
pod/sz-harbor-cluster-harbor-jobservice-7f45584546-nks8s     1/1     Running   1          4d18h
pod/sz-harbor-cluster-harbor-jobservice-7f45584546-nml6z     1/1     Running   1          4d18h
pod/sz-harbor-cluster-harbor-jobservice-7f45584546-x4ltr     1/1     Running   1          4d18h
pod/sz-harbor-cluster-harbor-notary-server-f89d84df7-bhvpl   1/1     Running   0          4d18h
pod/sz-harbor-cluster-harbor-notary-server-f89d84df7-kkmpd   1/1     Running   0          4d18h
pod/sz-harbor-cluster-harbor-notary-server-f89d84df7-vvvp6   1/1     Running   0          4d18h
pod/sz-harbor-cluster-harbor-notary-signer-b544646b6-277gb   1/1     Running   0          4d18h
pod/sz-harbor-cluster-harbor-notary-signer-b544646b6-jqrdp   1/1     Running   0          4d18h
pod/sz-harbor-cluster-harbor-notary-signer-b544646b6-ptfzp   1/1     Running   0          4d18h
pod/sz-harbor-cluster-harbor-portal-844c4f9f55-4xtf8         1/1     Running   0          4d18h
pod/sz-harbor-cluster-harbor-portal-844c4f9f55-g5ff6         1/1     Running   0          4d18h
pod/sz-harbor-cluster-harbor-portal-844c4f9f55-qvw9j         1/1     Running   0          4d18h
pod/sz-harbor-cluster-harbor-registry-665f8c845-2vcdl        2/2     Running   0          4d18h
pod/sz-harbor-cluster-harbor-registry-665f8c845-b5d88        2/2     Running   0          4d18h
pod/sz-harbor-cluster-harbor-registry-665f8c845-hsnpn        2/2     Running   0          4d18h
pod/sz-harbor-cluster-minio-zone-harbor-0                    1/1     Running   0          4d18h
pod/sz-harbor-cluster-minio-zone-harbor-1                    1/1     Running   0          4d18h
pod/sz-harbor-cluster-minio-zone-harbor-2                    1/1     Running   0          4d18h

NAME                                             TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)                    AGE
service/cluster-sz-harbor-cluster                ClusterIP   10.96.210.42     <none>        6379/TCP                   4d18h
service/harbor-sz-harbor-cluster                 ClusterIP   10.109.45.47     <none>        5432/TCP                   4d18h
service/harbor-sz-harbor-cluster-config          ClusterIP   None             <none>        <none>                     4d18h
service/harbor-sz-harbor-cluster-repl            ClusterIP   10.102.13.237    <none>        5432/TCP                   4d18h
service/rfs-sz-harbor-cluster                    ClusterIP   10.109.190.52    <none>        26379/TCP                  4d18h
service/sz-harbor-cluster-harbor-chartmuseum     ClusterIP   10.103.127.162   <none>        80/TCP                     4d18h
service/sz-harbor-cluster-harbor-clair           ClusterIP   10.96.68.203     <none>        80/TCP,6061/TCP,8080/TCP   4d18h
service/sz-harbor-cluster-harbor-core            ClusterIP   10.98.217.196    <none>        80/TCP                     4d18h
service/sz-harbor-cluster-harbor-jobservice      ClusterIP   10.110.158.46    <none>        80/TCP                     4d18h
service/sz-harbor-cluster-harbor-notary-server   ClusterIP   10.98.207.239    <none>        80/TCP                     4d18h
service/sz-harbor-cluster-harbor-notary-signer   ClusterIP   10.99.195.39     <none>        80/TCP                     4d18h
service/sz-harbor-cluster-harbor-portal          ClusterIP   10.101.19.211    <none>        80/TCP                     4d18h
service/sz-harbor-cluster-harbor-registry        ClusterIP   10.106.146.109   <none>        80/TCP,5001/TCP,8080/TCP   4d18h
service/sz-harbor-cluster-minio                  ClusterIP   10.101.162.130   <none>        9000/TCP                   4d18h
service/sz-harbor-cluster-minio-hl               ClusterIP   None             <none>        9000/TCP                   4d18h

NAME                                                     READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/rfs-sz-harbor-cluster                    1/1     1            1           4d18h
deployment.apps/sz-harbor-cluster-harbor-chartmuseum     3/3     3            3           4d18h
deployment.apps/sz-harbor-cluster-harbor-clair           3/3     3            3           4d18h
deployment.apps/sz-harbor-cluster-harbor-core            3/3     3            3           4d18h
deployment.apps/sz-harbor-cluster-harbor-jobservice      3/3     3            3           4d18h
deployment.apps/sz-harbor-cluster-harbor-notary-server   3/3     3            3           4d18h
deployment.apps/sz-harbor-cluster-harbor-notary-signer   3/3     3            3           4d18h
deployment.apps/sz-harbor-cluster-harbor-portal          3/3     3            3           4d18h
deployment.apps/sz-harbor-cluster-harbor-registry        3/3     3            3           4d18h

NAME                                                               DESIRED   CURRENT   READY   AGE
replicaset.apps/rfs-sz-harbor-cluster-bd4bdcdcf                    1         1         1       4d18h
replicaset.apps/sz-harbor-cluster-harbor-chartmuseum-68d7fc88bb    3         3         3       4d18h
replicaset.apps/sz-harbor-cluster-harbor-clair-5f65c876f8          3         3         3       4d18h
replicaset.apps/sz-harbor-cluster-harbor-core-86cb58b8dd           3         3         3       4d18h
replicaset.apps/sz-harbor-cluster-harbor-jobservice-7f45584546     3         3         3       4d18h
replicaset.apps/sz-harbor-cluster-harbor-notary-server-f89d84df7   3         3         3       4d18h
replicaset.apps/sz-harbor-cluster-harbor-notary-signer-b544646b6   3         3         3       4d18h
replicaset.apps/sz-harbor-cluster-harbor-portal-844c4f9f55         3         3         3       4d18h
replicaset.apps/sz-harbor-cluster-harbor-registry-665f8c845        3         3         3       4d18h

NAME                                                   READY   AGE
statefulset.apps/harbor-sz-harbor-cluster              2/2     4d18h
statefulset.apps/rfr-sz-harbor-cluster                 1/1     4d18h
statefulset.apps/sz-harbor-cluster-minio-zone-harbor   3/3     4d18h

NAME                                                      AGE
redisfailover.databases.spotahome.com/sz-harbor-cluster   4d18h
```
Check other relevant K8s resources:

```shell script
# Ingress
kubectl get ingress -n your_namespace

# Certificate

kubectl get cert -n your_namespace

# Secrets

kubectl get secret -n your_namespace

# ConfigMap

kubectl get cm -n your_namespace
```

## Delete 

```shell script
kubectl delete -f your_manifest.yaml
```
