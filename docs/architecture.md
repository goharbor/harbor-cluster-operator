# Harbor cluster operator design

## Overview

The overall design of the operator is shown below:

![harbor-cluster-operator](./assets/harbor-cluster-operator.png)

- A top CR(Customize Resource) `harbor-cluster` is defined to represent the full harbor stack. It is governed by the 
corresponding controller.
- The harbor service components are created and managed through the [harbor operator](https://github.com/goharbor/harbor-operator) and the CR `harbor`.
- The following related CRs of the dependent services will be created and managed by their operator controllers and bound 
to the `Harbor` CR via [secret](https://kubernetes.io/docs/concepts/configuration/secret/) or [configMaps](https://kubernetes.io/docs/tasks/configure-pod-container/configure-pod-configmap/) ways. Sometimes, the configurations of dependencies may be directly populated 
into the pods of CR `Harbor` to reduce the difficulties.
  - Database service (PostgresSQL)
  - Cache service (Redis)
  - Storage service (minIO)
- The dependent CRs are owned by the `harbor-cluster` CR to enable the cascaded deletion and GC when the CR `harbor-cluster` 
is deleted.

## CRD

Below is the spec defined for the CR `harbor-cluster`:

```yaml
# harbor version to be deployed
# this version determines the image tags of harbor service components
version: 1.10

# external URL for access Harbor registry
publicURL: https://harbor.registry.com

# password for the root admin 
adminPasswordSecret: adminSecret

# secret reference for the TLS certs
tlsSecret: tlsSecret

# certificate issuers
certificateIssuerRef: cert_issuer

paused: false
priority: 30

# pod instance number
replicas: 3

# cache service(Redis) configurations
redis:
  connections: conSecret1
  mode: AOT/RDB
  nodes: 3
  version: 5.2
  # More options
  #……

# database service (PostgresSQL) configuration
database:
  connection: conSecret2
  replicas: 3
  version: 9.6
  connectionOptions:
    ssl_mode: disable
    max_idle_conns: 2
    max_open_conns: 0
  # more configuration options
  # ……

# storage service configurations
# might be external cloud storage services or
# default in-cluster storage (minIO)
storage:
  kind: local/remote
  options:
    host: s3.com

# source registry of images
imageSource:
  registry: harbor.com
  imagePullSecrets:
   - pSecret
# log configurations
log:
  level: debug

# set proxy
proxy:
  http_proxy: 10.10.20.2
  https_proxy: 10.10.20.2
  no_proxy: 10.123.111.10
    components:
    - core
    - jobservice
    - clair

# extra configuration options for jobservices
jobService:
  workerCount: 10
  replicas: 5

# extra configuration options for clair scanner
clair:
  updateInterval: 10
  vulnerabilitySources:
    - ubuntu
    - alphine

# extra configuration options for trivy scanner
trivy:
  github_token: 123

# extra configuration options for chartmeseum
chartMuseum:
  absoluteURL: true
```

## Reconcile flow

[TBD]

## Misc

[TBD]