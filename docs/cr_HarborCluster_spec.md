## Harbor-Cluster CRD

Below is the spec defined for the CR `harbor-cluster`:

```yaml
# harbor version to be deployed
# this version determines the image tags of harbor service components
# required
version: 1.10

# external URL for access Harbor registry
# required
publicURL: https://harbor.registry.com

# password for the root admin 
# required
adminPasswordSecret: adminSecret

# secret reference for the TLS certs
tlsSecret: tlsSecret

# certificate issuers
certificateIssuerRef: 
  name: cert_issuer

paused: false

priority: 30

# pod instance number
# required
replicas: 3

# source registry of images
imageSource:
  registry: harbor.com
  imagePullSecrets:
   - pSecret

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

notary:
  publicURL: "http://.."  
  
# cache service(Redis) configurations
# might be external redis services or inCluster redis services
# required
redis:
  # set the kind of which redis service to be used, inCluster or external.
  # setting up a harbor-cluster with external redis service should provide client params to communicate. The difference between inCluster redis and external redis is that the inCluster redis installed automatically. the params of external kind are in the following comments.
  # kind: external
  #   // the secret must contains "address:port","usernane" and "password".
  #   // required
  #   secretName: secret
  #.  // Maximum number of socket connections.
  #   // Default is 10 connections per every CPU as reported by runtime.NumCPU.
  #   // optional
  #   poolSize: 10
  #   // TLS Config to use. When set TLS will be negotiated.
  #   // set the secret which type of Opaque, and contains "tls.key","tls.crt","ca.crt".
  #   // optional
  #   tlsConfig: secretName
  kind: inCluster
  server:
    replicas: 3
    # optional
    resources:
      requests:
        memory: 2048Mi
        cpu: 2000m
      limit:
        memory: 2048Mi
        cpu: 2000m
    # optional
    storageClassName: default
    storage: 5Gi
  sentinel:
    replicas: 3

# database service (PostgresSQL) configuration
# required
database:
  # set the kind of which redis service to be used, inCluster or external.
  # a sample of external kind.
  # kind: external
  #.  // the secret must contains "address:port","usernane" and "password".
  #   // required
  #   secretName: secret
  #   // TLS Config to use. When set TLS will be negotiated.
  #   // set the secret which type of Opaque, and contains "tls.key","tls.crt","ca.crt".
  #   // optional
  #   sslConfig: secretName
  #   connect_timeout: 10
  kind: inCluster
    storage: 1Gi
    replicas: 2
    version: "12"
    # optional
    storageClassName: default
    # optional
    resources:
      limits:
        cpu: 500m
        memory: 500Mi
      requests:
        cpu: 100m
        memory: 250Mi

# storage service configurations
# might be external cloud storage services or inCluster storage (minIO)
# required
storage:
  # set the kind of which storage service to be used. Set the kind as "azure",
  # "gcs", "s3", "oss", "swift" or "inCluster" and fill the information
  # in the options section. inCluster indicates the local storage service of harbor-cluster. We use minIO as a default built-in object storage service. All of kind and option parameters are in the following comments.
  # azure:
  #   accountname: accountname
  #   accountkey: base64encodedaccountkey
  #   container: containername
  #   realm: core.windows.net
  # gcs:
  #   bucket: bucketname
  #   # The base64 encoded json file which contains the key
  #   encodedkey: base64-encoded-json-key-file
  #   rootdirectory: /gcs/object/name/prefix
  #   chunksize: "5242880"
  # s3:
  #   region: us-west-1
  #   bucket: bucketname
  #   accesskey: awsaccesskey
  #   secretkey: awssecretkey
  #   regionendpoint: http://myobjects.local
  #   encrypt: false
  #   keyid: mykeyid
  #   secure: true
  #   v4auth: true
  #   chunksize: "5242880"
  #   rootdirectory: /s3/object/name/prefix
  #   storageclass: STANDARD
  # swift:
  #   authurl: https://storage.myprovider.com/v3/auth
  #   username: username
  #   password: password
  #   container: containername
  #   region: fr
  #   tenant: tenantname
  #   tenantid: tenantid
  #   domain: domainname
  #   domainid: domainid
  #   trustid: trustid
  #   insecureskipverify: false
  #   chunksize: 5M
  #   prefix:
  #   secretkey: secretkey
  #   accesskey: accesskey
  #   authversion: 3
  #   endpointtype: public
  #   tempurlcontainerkey: false
  #   tempurlmethods:
  # oss:
  #   accesskeyid: accesskeyid
  #   accesskeysecret: accesskeysecret
  #   region: regionname
  #   bucket: bucketname
  #   endpoint: endpoint
  #   internal: false
  #   encrypt: false
  #   secure: true
  #   chunksize: 10M
  #   rootdirectory: rootdirectory
  # Here is a sample of how to use inCluster kind to provide storage service.
  kind: inCluster
  options:
    provider: minIO
    spec:
      # Supply number of replicas.
      # For standalone mode, supply 1. For distributed mode, supply 4 or more (should be even).
      # Note that the operator does not support upgrading from standalone to distributed mode.
      replicas: 4
      version: RELEASE.2020-01-03T19-12-21Z
      # VolumeClaimTemplate allows a user to specify how volumes inside a MinIOInstance
      volumeClaimTemplate:
        spec:
          # optional
          storageClassName: default
          accessModes:
            - ReadWriteOnce
          resources:
            requests:
              storage: 10Gi
      # optional
      resources:
        requests:
          memory: 512Mi
          cpu: 250m
        limits:
          memory: 512Mi
          cpu: 250m
```

