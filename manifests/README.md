# Deployment manifests

Manifest yaml files here are `kustomization` templates used to deploy the harbor cluster operator as well its dependant operators.

## Notes

* Deployments of `Cert Manager` and `Ingress Controller` are not covered yet.
* Per the limitation of `kustomization`, the deployment manifest of Harbor(core) operator cannot be referred via 
remote `kustomization` template file. Here keeping a copy of deployment manifest generated through the `kustomization` 
templates at the `Harbor-operator` repository.
* Other relying on operators are linked via remote manifests.
* `all-in-one.yaml` is the all-in-one operator deployment manifest yaml file built by the `kustomization` templates.

## Usage

Directly use all-in-one yaml:

```shell script
kubectl apply -f ./all-in-one.yaml
```

Use `kustomize`:

```shell script
kustomize build . | kubectl apply -f -
```

## Attentions

If the CRD of harbor operator is updated, the referred harbor core manifest here should be updated too.

```shell script
cd harbor-operator/

make manifests

kustomize build config/default -o all-core-operator-resources.yaml

mv ./all-core-operator-resources.yaml ../harbor-cluster-operator/manifests/core-oeprator
```

If the CRD of harbor cluster operator is changed, update the related `kustomization` files first.

```shell script
cd harbor-cluster-operator

make manifests

cd manifests
# generate all-in-one yaml

kustomize build . -o all-in-one.yaml
```



