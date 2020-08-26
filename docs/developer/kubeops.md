# Kubeapps Kubeops Developer Guide

The `kubeops` component is a micro-service that creates an API endpoint for accessing the Helm API and Kubernetes resources.

## Prerequisites

- [Git](https://git-scm.com/)
- [Make](https://www.gnu.org/software/make/)
- [Go programming language](https://golang.org/dl/)
- [Docker CE](https://www.docker.com/community-edition)
- [Kubernetes cluster (v1.8+)](https://kubernetes.io/docs/setup/pick-right-solution/). [Minikube](https://github.com/kubernetes/minikbue) is recommended.
- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)

## Environment

```bash
export GOPATH=~/gopath
export PATH=$GOPATH/bin:$PATH
export KUBEAPPS_DIR=$GOPATH/src/github.com/kubeapps/kubeapps
```

## Download the Kubeapps source code

```bash
git clone --recurse-submodules https://github.com/kubeapps/kubeapps $KUBEAPPS_DIR
```

The `kubeops` sources are located under `cmd/kubeops/` and use packages from the `pkg` directory.

### Install Kubeapps in your cluster

Kubeapps is a Kubernetes-native application. To develop and test Kubeapps components we need a Kubernetes cluster with Kubeapps already installed. Follow the [Kubeapps installation guide](../../chart/kubeapps/README.md) to install Kubeapps in your cluster.

### Building the `kubeops` binary

```bash
cd $KUBEAPPS_DIR/cmd/kubeops
go build
```

This builds the `kubeops` binary in the working directory.

### Running in development

If you are using Minikube it is important to start the cluster enabling RBAC (on by default in Minikube 0.26+) in order to check the authorization features:

```bash
minikube start
eval $(minikube docker-env)
```

Note: By default, Kubeapps will try to fetch the latest version of the image so in order to make this workflow work in Minikube you will need to update the imagePullPolicy first:

```bash
kubectl patch deployment kubeapps-internal-kubeops -n kubeapps --type=json -p='[{"op": "replace", "path": "/spec/template/spec/containers/0/imagePullPolicy", "value": "IfNotPresent"}]'
```

The easiest way to create the `kubeops` image is to execute the Makefile task to do so:

```bash
IMAGE_TAG=dev make kubeapps/kubeops
```

This will generate an image `kubeapps/kubeops:dev` that you can use in the current deployment:

```bash
kubectl set image -n kubeapps deployment kubeapps-internal-kubeops kubeops=kubeapps/kubeops:dev
```

For further redeploys you can change the version to deploy a different tag or rebuild the same image and restart the pod executing:

```bash
kubectl delete pod -n kubeapps -l app=kubeapps-internal-kubeops
```

Note: If you are using a cloud provider to develop the service you will need to retag the image and push it to a public registry.

### Running tests

You can run the kubeops tests along with the tests of all the projects:

```bash
make test
```
