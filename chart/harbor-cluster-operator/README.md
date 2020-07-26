# harbor-cluster-operator

## Install

```bash
helm install harbor-cluster-operator https://github.com/goharbor/harbor-cluster-operator/releases/download/v0.5.0/harbor-cluster-operator-chart.tgz
```

you will see follow:

```
NAME: harbor-cluster-operator
LAST DEPLOYED: Mon Jan  6 14:47:48 2020
NAMESPACE: default
STATUS: deployed
REVISION: 1
TEST SUITE: None
```

## Uninstall

```bash
$ helm delete harbor-cluster-operator
release "harbor-cluster-operator" uninstalled
```

## Configuration

The following table lists the configurable parameters of the harbor-cluster-operator chart and their default values.

| Parameter                                 | Description                                                        | Default                             |
|-------------------------------------------|--------------------------------------------------------------------|-------------------------------------|
| `manager.resources.limits.cpu`            | CPU resource limit of manager container                            | `100m`                              |
| `manager.resources.limits.memory`         | Memory resource limit of manager container                         | `256Mi`                             |
| `manager.resources.requests.cpu`          | CPU resource request of manager container                          | `100m`                              |
| `manager.resources.requests.memory`       | Memory resource request of manager container                       | `256Mi`                             |
| `manager.metrics.addr`                    | Addr of metrics served                                             | `localhost`                         |
| `manager.metrics.port`                    | Port of metrics served                                             | `8080`                              |
| `manager.proxy.image`                     | Custom resources kube-rbac-proxy image                             | `gcr.io/kubebuilder/kube-rbac-proxy:v0.4.1`       |
| `manager.manager.image`                   | Custom resources manager image                                     | `controller:latest`       |