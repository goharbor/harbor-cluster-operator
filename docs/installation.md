# Installation

Follow the guide shown below to deploy the Harbor operators and the relevant dependant services.

## Prerequisite

### Cert-Manager

[Cert-Manager] is used to manage the related certificates of Harbor. Use the following command to install:

```shell script
# Kubernetes 1.15+
$ kubectl apply --validate=false -f https://github.com/jetstack/cert-manager/releases/download/v0.16.1/cert-manager.yaml
```

### Ingress Controller

Deployed Harbor services are exposed with ingress way. An ingress controller should be installed in the target K8s cluster.
Nginx ingress controller is regular option. Based on your environment, follow the guide 
shown [here](https://kubernetes.github.io/ingress-nginx/deploy/) to install it.

## Install Operators

A kustomization template is provided to install all the related operators required by deploying all-in-one harbor cluster.
Use the command shown below to start the installation:

```shell script
kubectl -f manifests/all-in-one.yaml
```

or 

```shell script
kustomize build manifests/ | kubectl apply -f -
```

## Uninstall Operators

Use K8s delete command and the deployment manifest to uninstall all resources of operators.

```shell script
kubectl delete -f manifests/all-in-one.yaml
```

## Other References

- Follow guide shown [here](./installation_local.md) to deploy harbor operators on the local cluster (kind) and deploy 
sample Harbor with in-cluster dependant services
- Follow guide shown [here](./installation_external.md) to deploy harbor operators and deploy sample Harbor with external services
- Follow [sample deployment guide](./sample_deploy_guide.md) to deploy a sample Harbor