# Multi-cluster support for Suomitek-appboard

## Objective

Enable Suomitek-appboard users to be able to **install apps into other configured Kubernetes clusters** in addition to the cluster on which Suomitek-appboard is installed.

This design aims to utilize the fact that to install apps on additional clusters, all that is required are user credentials with the required RBAC and network access from Suomitek-appboard' pods to the Kubernetes API of the additional cluster. With Helm 3, we no longer need to run any infrastructure (ie. tiller or suomitek-appboard services) on the additional cluster(s) to be able to deploy applications, as we are just sending YAML to the api server using the users credentials.

## User Stories

* As an operator of multiple Kubernetes clusters, I want to configure Suomitek-appboard **on one cluster** so that **additional clusters can be targeted** when users install applications, as long as the user's credentials authorize them to do so.
* As a user of Suomitek-appboard, I want to login to my teams' Suomitek-appboard instance and install ApplicationX into my namespace on the teams sandbox cluster to test before installing on the staging cluster.

## Explicit Constraints
The following constraints are for wider discussion. While we cannot achieve a solution which fits all use-cases, we want a solution which suits many use cases, is within Suomitek-appboard current scope (ie. doesn't require a cluster-admin agent running in additional clusters) and within Kubernetes design principles:

* At least initially, we **will not support K8s service account tokens as a user authentication option** for users in multiple clusters. Suomitek-appboard initially supported only service tokens for user authentication until support for OIDC/SSO with an appropriately configured Kubernetes API server was later added. At least initially, the multi-cluster support will **only be available when using SSO such that a single user can be authenticated with multiple clusters using the same credential** (eg. an OIDC id_token). We may later decide to support service account tokens for users, but it is not recommended for use (even with a single cluster).
* **We should not need to run parts of Suomitek-appboard infrastructure in each additional cluster**: Ideally we will be able to achieve our objective without the requirement that human operators deploy and maintain extra Suomitek-appboard services in additional clusters, which would move away significantly from `helm install`ing Suomitek-appboard. Suomitek-appboard allows easy configuration of Kubernetes Apps, resulting in YAML which can be applied to the K8s API server with the users credentials.
* **Network access from the Suomitek-appboard cluster to the additional clusters’ API server**. For Suomitek-appboard to support an additional cluster, the cluster operator would need to ensure that the additional cluster's API server is reachable from the Suomitek-appboard’ pods on the initial cluster. This is normally the case for hosted clusters which have public endpoints requiring authorization or private endpoints within a common private network or even on multiple private networks which can be bridged.

## Design overview

The overview displayed below shows the simpest scenario of the multi-cluster support (ie. without private app repository support), which is discussed further in the [design doc](https://docs.google.com/document/d/1-6cKxOsW6K5u3lK7Om2zQeVYVPxzHT6dVwej5wy3_9A/edit).

![Suomitek-appboard Multi-cluster Overview](img/Suomitek-appboard-Multi-cluster-simple.png)

Similarly, the proposed extension including private repositories on additional clusters, though due to the current transition (in the Helm community) from chart repositories to OCI repositories, we may delay the private repository support until we implement OCI repository support:

![Suomitek-appboard Multi-cluster support with private repositories](img/Suomitek-appboard-Multi-cluster-private-repo.png)


## Details and discussion

More details, design considerations and discussion is in the separate [design doc](https://docs.google.com/document/d/1-6cKxOsW6K5u3lK7Om2zQeVYVPxzHT6dVwej5wy3_9A/edit).