# Setup Suomitek-appboard testing environment

This guide explains how to setup your environment to test Suomitek-appboard integration with other services.

## Background

Suomitek-appboard can be integrated with other services to extend its capabilities. Find more information about these integrations in the links below:

- [Using Private App Repositories with Suomitek-appboard](../user/private-app-repository.md).
- [Kubernetes Service Catalog Suomitek-appboard Integration](../user/service-catalog.md).

This guide aims to provide the instructions to easily setup the environment to test these integrations.

## Prerequisites

- [Kubernetes cluster (v1.12+)](https://kubernetes.io/docs/setup/pick-right-solution/).
- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/).
- [Helm client](https://helm.sh/docs/intro/install/).

## Environment setup

We are providing scripts to automatically setup both Suomitek-appboard and the services to integrate on a K8s cluster. Find them under the [script](../../script) directory.

Currently supported integrations:

- Suomitek-appboard integration with Harbor.

### Suomitek-appboard integration with Harbor

You can setup environment to test Suomitek-appboard integration with Harbor using the scripts below:

- [setup-suomitek-appboard](../../script/setup-suomitek-appboard.sh).
- [setup-harbor](../../script/setup-harbor.sh).

These scripts will create the necessary namespaces, install the charts, wait for them to be available, and perform any extra action that might be needed. Find detailed information about how to use these scripts running the commands below:

```bash
./setup-suomitek-appboard.sh --help
./setup-harbor.sh --help
```

You can also use the [setup-suomitek-appboard-harbor](../../script/setup-suomitek-appboard-harbor.sh) script which is a wrapper that uses both the scripts mentioned above with some default values:

- Install Harbor under the `harbor` namespace.
- Install Suomitek-appboard under the `suomitek-appboard` namespace.
- Adds Harbor as an extra initial repository to Suomitek-appboard, based on its service hostname.

#### Cleaning up the environment

You can use the scripts [delete-suomitek-appboard](../../script/delete-suomitek-appboard.sh) and [delete-harbor](../../script/delete-harbor.sh) to uninstall Suomitek-appboard and Harbor releases from the cluster, respectively. These scripts will also remove the associated namespaces and resources.

> Note: you can use the [delete-suomitek-appboard-harbor](../../script/delete-suomitek-appboard-harbor.sh) script to clean up the environment if you used the [setup-suomitek-appboard-harbor](../../script/setup-suomitek-appboard-harbor.sh) script to setup the environment.
