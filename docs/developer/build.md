# The Kubeapps Build Guide

This guide explains how to build Kubeapps.

## Prerequisites

- [Git](https://git-scm.com/)
- [Make](https://www.gnu.org/software/make/)
- [Go programming language](https://golang.org/)
- [kubecfg](https://github.com/ksonnet/kubecfg)
- [Docker CE](https://www.docker.com/community-edition)

## Environment setup

```bash
export GOPATH=~/gopath
export PATH=$GOPATH/bin:$PATH
export KUBEAPPS_DIR=$GOPATH/src/github.com/suomitek/suomitek-appboard
```

## Download suomitek-appboard source code

```bash
git clone --recurse-submodules https://github.com/suomitek/suomitek-appboard $KUBEAPPS_DIR
cd $KUBEAPPS_DIR
```

## Build suomitek-appboard

Kubeapps consists of a number of in-cluster components. To build all these components in one go:

```bash
make IMAGE_TAG=myver all
```

Or if you wish to build specific component(s):

```bash
# to build the suomitek-appboard binary
make IMAGE_TAG=myver suomitek-appboard

# to build the suomitek-appboard/dashboard docker image
make IMAGE_TAG=myver suomitek-appboard/dashboard

# to build the suomitek-appboard/apprepository-controller docker image
make IMAGE_TAG=myver suomitek-appboard/apprepository-controller

# to build the suomitek-appboard/tiller-proxy docker image
make IMAGE_TAG=myver suomitek-appboard/tiller-proxy
```

## Running tests

To test all the components:

```bash
make test
```

Or if you wish to test specific component(s):

```bash
# to test the suomitek-appboard binary
make test-suomitek-appboard

# to test suomitek-appboard/dashboard
make test-dashboard

# to test the cmd/apprepository-controller package
make test-apprepository-controller

# to test the cmd/tiller-proxy package
make test-tiller-proxy
```
