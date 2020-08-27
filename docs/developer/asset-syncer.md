# Suomitek-appboard asset-syncer Developer Guide

The `asset-syncer` component is a tool that scans a Helm chart repository and populates chart metadata in the database. This metadata is then served by the `assetsvc` component.

## Prerequisites

- [Git](https://git-scm.com/)
- [Make](https://www.gnu.org/software/make/)
- [Go programming language](https://golang.org/dl/)
- [Docker CE](https://www.docker.com/community-edition)
- [Kubernetes cluster (v1.8+)](https://kubernetes.io/docs/setup/pick-right-solution/). [Minikube](https://github.com/kubernetes/minikbue) is recommended.
- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
- [Telepresence](https://telepresence.io)

## Environment

```bash
export GOPATH=~/gopath
export PATH=$GOPATH/bin:$PATH
export KUBEAPPS_DIR=$GOPATH/src/github.com/suomitek/suomitek-appboard
```

## Download the Suomitek-appboard source code

```bash
git clone https://github.com/suomitek/suomitek-appboard $KUBEAPPS_DIR
```

The `asset-syncer` sources are located under the `cmd/asset-syncer/` directory.

### Install Suomitek-appboard in your cluster

Suomitek-appboard is a Kubernetes-native application. To develop and test Suomitek-appboard components we need a Kubernetes cluster with Suomitek-appboard already installed. Follow the [Suomitek-appboard installation guide](../../chart/suomitek-appboard/README.md) to install Suomitek-appboard in your cluster.

### Building the `asset-syncer` image

```bash
cd $KUBEAPPS_DIR
make suomitek-appboard/asset-syncer
```

This builds the `asset-syncer` Docker image.

### Running in development

When using MongoDB:

```bash
export DB_PASSWORD=$(kubectl get secret --namespace suomitek-appboard suomitek-appboard-mongodb -o go-template='{{index .data "mongodb-root-password" | base64decode}}')
telepresence --namespace suomitek-appboard --docker-run -e DB_PASSWORD=$DB_PASSWORD --rm -ti suomitek-appboard/asset-syncer /asset-syncer sync --database-user=root --database-url=suomitek-appboard-mongodb --database-type=mongodb --database-name=charts stable https://kubernetes-charts.storage.googleapis.com
```

When using PostgreSQL:

```bash
export DB_PASSWORD=$(kubectl get secret --namespace suomitek-appboard suomitek-appboard-db -o go-template='{{index .data "postgresql-password" | base64decode}}')
telepresence --namespace suomitek-appboard --docker-run -e DB_PASSWORD=$DB_PASSWORD --rm -ti suomitek-appboard/asset-syncer /asset-syncer sync --database-user=postgres --database-url=suomitek-appboard-postgresql:5432 --database-type=postgresql --database-name=assets stable https://kubernetes-charts.storage.googleapis.com
```

Note that the asset-syncer should be rebuilt for new changes to take effect.

### Running tests

You can run the asset-syncer tests along with the tests for the Suomitek-appboard project:

```bash
go test -v ./...
```
