# harbor-cluster-operator

**NOTES:** The `master` branch may be in an *unstable or even broken state* during development. Currently, please use branch [release-0.5.0](https://github.com/goharbor/harbor-cluster-operator/tree/release-0.5.0) instead of the `master` branch in order to get a stable deployment.

[Harbor](https://github.com/goharbor/harbor) is a CNCF hosted open source trusted cloud native registry project that stores, signs, and scans content.

The project [harbor-operator](https://github.com/goharbor/harbor-operator) is created to cover both Day1 and Day2 operations of an enterprise-grade Harbor deployment.
The harbor-operator extends the usual K8s resources with Harbor-related custom ones. The Kubernetes API can then be used in a declarative way to manage Harbor and 
ensure its high-availability operation, thanks to the [Kubernetes control loop](https://kubernetes.io/docs/concepts/#kubernetes-control-plane).

This operator aims to build a Kubernetes operator on top of the harbor-operator to deploy and manage a full Harbor service stack including both the harbor service components 
and its relevant dependent services such as database(PostgresSQL), Cache(Redis) and default in-cluster storage(minIO) services in a scalable and high-available way. It provides 
a solid unified solution to cover the lifecycle management of Harbor service. The customize resource `HarborCluster` is defined to describe the full harbor stack that includes 
the dependencies. The CR `HarborCluster` owns the underlying CRs like `Harbor` managed by harbor-operator, `Postgresql` managed by PostgresSQL operator, `RedisFailover` managed by Redis 
operator and `MinIOInstance` managed by minIO operator. The reconcile process of `HarborCluster` wil make sure the actual state of the `HarborCluster` CR matches the designed state set 
in the spec by the users. The reconcile process also takes care of the service creation and ready order to reflect the real service dependent topology to avoid starting failures issues.

Project codebase is scaffolded by [kubebuilder](https://kubebuilder.io/) V2(.2).

## Features

With this operator, you're able to deploy and manage a full Harbor stack:

- Provision a full Harbor stack including the relevant dependent services like database(PostgresSQL), cache(Redis) and 
in-cluster storage(minIO) services in a scalable and high-available way.
- Inherit deployment customization capabilities from the underlying harbor-operator, the following components could be optional:
  - ChartMuseum
  - Notary
  - Clair
  - Trivy
- Update the spec of the deployed Harbor stack to do adjustments like replicas (scalability) and service properties.
- Upgrade the deployed Harbor stack to a newer version.
- Delete the Harbor stack and all the related resources owned by the stack.

## Design

Diagram below shows the overall design of this operator,

![harbor-cluster-operator](./docs/assets/harbor-cluster-operator.png)

For more design details, check the [architecture](./docs/architecture.md) document.

## Installation

You can follow the [installation guide](docs/installation.md) to deploy this operator to your K8s clusters.

Additionally, follow [sample deployment guide](./docs/sample_deploy_guide.md) to have a try of deploying the sample to your K8s clusters.

## Versioning & Dependencies

| Component \ Versions |  0.5.0 | 1.0.0 | 1.1.0 |
|----------------------|--------|-------|-------|
| **Harbor**           | 1.10.x | [TBD] | [TBD] |
|                      |        |               |
| harbor-operator      | 0.5.2  | [TBD] | [TBD] |
| PostgresSQL operator | [TBD]  | [TBD] | [TBD] |
| Redis operator       | [TBD]  | [TBD] | [TBD] |
| minIO operator       | 3.0.13 | [TBD] | [TBD] |

## Compatibilities

| Kubernetes / Versions |  0.5.0  |  1.0.0  | 1.1.0 |
|-----------------------|---------|---------|------|
|     1.17              |    +    | [TBD] | [TBD] |
|     1.18              |    +    | [TBD] | [TBD] |
|     1.19              |    +    | [TBD] | [TBD] |

**Notes:** `+`= verified `-`= not verified


## Development

Interested in contributions? Follow the [CONTRIBUTING](./docs/CONTRIBUTING.md) guide to start on this project. Your contributions will be highly appreciated and creditable.

## Community

* Slack channel #harbor-operator-dev at [CNCF Workspace](https://slack.cncf.io)
* Send mail to Harbor dev mail group:  harbor-dev@lists.cncf.io

## Documents

See documents [here](./docs).

## Additional Documents

* [Kubernetes Operator](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/)
* [Kubebuilder](https://book.kubebuilder.io/)
* [Custom Resource Definition](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/)
* [Underlying harbor-operator](https://github.com/goharbor/harbor-operator)
* [Underlying PostgresSQL operator](https://github.com/zalando/postgres-operator)
* [Underlying Redis operator](https://github.com/spotahome/redis-operator)
* [Underlying Storage operator](https://github.com/minio/minio-operator)


## License

[Apache-2.0](https://github.com/goharbor/harbor-cluster-operator/blob/master/LICENSE)
