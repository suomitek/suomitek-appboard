# The Suomitek-appboard Overview

This document describes the Suomitek-appboard architecture at a high level.

## Components

### Suomitek-appboard dashboard

At the heart of Suomitek-appboard is an in-cluster Kubernetes dashboard that provides you a simple browse and click experience for installing and managing Kubernetes applications packaged as Helm charts.

Additionally, the dashboard integrates with the [Kubernetes service catalog](https://github.com/kubernetes-incubator/service-catalog) and enables you to browse and provision cloud services via the [Open Service Broker API](https://github.com/openservicebrokerapi/servicebroker).

The dashboard is written in the JavaScript programming language and is developed using the React JavaScript library.

### Tiller proxy (deprecated)

In order to secure the access to Tiller and allow the dashboard to contact the Helm Tiller server we deploy a proxy that handles the communication with Tiller. The goal of this proxy is to validate that the user doing the request has sufficent permissions to create or delete all the resources that are part of the specific chart being deployed or deleted.

This proxy is written in Go. Check more details about the implementation in this [document](/cmd/tiller-proxy/README.md).

### Kubeops

Kubeops is the successor of Tiller proxy. It's the service in charge of communicating both with the Helm (v3) API and other k8s resources like AppRepositories or Secrets.
Check more details about the implementation in [this document](/docs/developer/kubeops.md).

### Apprepository CRD and Controller

Chart repositories in Suomitek-appboard are managed with a `CustomResourceDefinition` called `apprepositories.suomitek.com`. Each repository added to Suomitek-appboard is an object of type `AppRepository` and the `apprepository-controller` will watch for changes on those type of objects to update the list of available charts to deploy.

### `asset-syncer`

The `asset-syncer` component is tool that scans a Helm chart repository and populates chart metadata in a database. This metadata is then served by the `assetsvc` component. Check more details about the implementation in this [document](/docs/developer/asset-syncer.md).

### `assetsvc`

The `assetsvc` component is a micro-service that creates an API endpoint for accessing the metadata for charts and other resources that's populated in a database. Check more details about the implementation in this [document](/docs/developer/asset-syncer.md).
