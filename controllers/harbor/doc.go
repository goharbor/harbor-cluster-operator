package harbor

//CRStatus properties:
//
// cache CRStatus properties MAY include:
// - url
// - namespace
//
// database CRStatus properties MAY include:
// - name: host
//   value: "..."
// - name: port
//   value: "..."
// - name: database
//   value: "..."
// - name: username
//   value: "..."
// - name: password
//   value: "..."
//
// or include a secret
// - name: secret
//   value: "k8s secret name in same namespace"

//- storageSecret
//-
//
//
//apiVersion: goharbor.io/v1alpha2
//kind: Harbor
//metadata:
//...
//spec:
//...
//components:
//...
//registry:
//...
//storageSecret: registry-backend
//cacheSecret: registry-redis
//jobService:
//...
//redisSecret: jobservice-redis
//chartMuseum:
//...
//cacheSecret: chartmuseum-redis
//core:
//databaseSecret: core-database
//...
//clair:
//databaseSecret: clair-database
//adapter:
//...
//redisSecret: clair-redis
//...
//notary:
//...
//server:
//databaseSecret: notary-server-database
//...
//signer:
//databaseSecret: notary-signer-database
//...
