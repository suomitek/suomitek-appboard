# Suomitek-appboard

[![CircleCI](https://circleci.com/gh/suomitek/suomitek-appboard/tree/master.svg?style=svg)](https://circleci.com/gh/suomitek/suomitek-appboard/tree/master)

[Suomitek-appboard](https://suomitek.com) is a web-based UI for deploying and managing applications in Kubernetes clusters. Suomitek-appboard allows you to:

- Browse and deploy [Helm](https://github.com/helm/helm) charts from chart repositories
- Inspect, upgrade and delete Helm-based applications installed in the cluster
- Add custom and private chart repositories (supports [ChartMuseum](https://github.com/helm/chartmuseum) and [JFrog Artifactory](https://www.jfrog.com/confluence/display/RTF/Helm+Chart+Repositories))
- Browse and provision external services from the [Service Catalog](https://github.com/kubernetes-incubator/service-catalog) and available Service Brokers
- Connect Helm-based applications to external services with Service Catalog Bindings
- Secure authentication to Suomitek-appboard using an [OAuth2/OIDC provider](https://github.com/suomitek/suomitek-appboard/blob/master/docs/user/using-an-OIDC-provider.md)
- Secure authorization based on Kubernetes [Role-Based Access Control](https://github.com/suomitek/suomitek-appboard/blob/master/docs/user/access-control.md)

## TL;DR

For Helm 2:

```bash
helm repo add chartmuseum http://helm.yongchehang.com
helm install --name suomitek-appboard --namespace suomitek-appboard chartmuseum/suomitek-appboard
```

For Helm 3:

```bash
helm repo add chartmuseum http://helm.yongchehang.com
kubectl create namespace suomitek-appboard
helm install suomitek-appboard --namespace suomitek-appboard chartmuseum/suomitek-appboard --set useHelm3=true
```

## Introduction

This chart bootstraps a [Suomitek-appboard](https://suomitek.com) deployment on a [Kubernetes](http://kubernetes.io) cluster using the [Helm](https://helm.sh) package manager.

It also packages the [Bitnami MongoDB chart](https://github.com/bitnami/charts/tree/master/bitnami/mongodb) or the [Bitnami PostgreSQL chart](https://github.com/bitnami/charts/tree/master/bitnami/postgresql) which is required for bootstrapping a deployment for the database requirements of the Suomitek-appboard application.

## Prerequisites

- Kubernetes 1.8+ (tested with Azure Kubernetes Service, Google Kubernetes Engine, minikube and Docker for Desktop Kubernetes)
- Helm 2.14.0+
- Administrative access to the cluster to create Custom Resource Definitions (CRDs)

## Installing the Chart

To install the chart with the release name `suomitek-appboard`:

For Helm 2:

```bash
helm repo add chartmuseum http://helm.yongchehang.com
helm install --name suomitek-appboard --namespace suomitek-appboard chartmuseum/suomitek-appboard
```

> **IMPORTANT** This assumes an insecure Helm 2 installation, which is not recommended in production. See [the documentation to learn how to secure Helm 2 and Suomitek-appboard in production](https://github.com/suomitek/suomitek-appboard/blob/master/docs/user/securing-suomitek-appboard.md).

For Helm 3:

```bash
helm repo add chartmuseum http://helm.yongchehang.com
kubectl create namespace suomitek-appboard
helm install suomitek-appboard --namespace suomitek-appboard chartmuseum/suomitek-appboard --set useHelm3=true
```

The command deploys Suomitek-appboard on the Kubernetes cluster in the `suomitek-appboard` namespace. The [Parameters](#parameters) section lists the parameters that can be configured during installation.

> **Caveat**: Only one Suomitek-appboard installation is supported per namespace

Once you have installed Suomitek-appboard follow the [Getting Started Guide](https://github.com/suomitek/suomitek-appboard/blob/master/docs/user/getting-started.md) for additional information on how to access and use Suomitek-appboard.

## Parameters

For a full list of configuration parameters of the Suomitek-appboard chart, see the [values.yaml](values.yaml) file.

Specify each parameter using the `--set key=value[,key=value]` argument to `helm install`. For example,

```bash
helm install suomitek-appboard --namespace suomitek-appboard \
  --set assetsvc.service.port=9090 \
    chartmuseum/suomitek-appboard
```

The above command sets the port for the assetsvc Service to 9090.

Alternatively, a YAML file that specifies the values for parameters can be provided while installing the chart. For example,

```bash
helm install suomitek-appboard --namespace suomitek-appboard -f custom-values.yaml chartmuseum/suomitek-appboard
```

## Configuration and installation details

### Configuring Initial Repositories

By default, Suomitek-appboard will track the [community Helm charts](https://github.com/helm/charts) and the [Kubernetes Service Catalog charts](https://github.com/kubernetes-incubator/service-catalog). To change these defaults, override with your desired parameters the `apprepository.initialRepos` object present in the [values.yaml](values.yaml) file.

### Configuring the database to use

Suomitek-appboard supports two database types: MongoDB or PostgreSQL. By default MongoDB is installed. If you want to enable PostgreSQL instead set the following values when installing the application: `mongodb.enabled=false` and `postresql.enabled=true`.

> **Note**: Changing the database type when upgrading is not supported.

### Enabling Operators

Since v1.9.0, Suomitek-appboard supports to deploy and manage Operators within its dashboard. To enable this feature, set the flag `featureFlags.operators=true`. More information about how to enable and use this feature can be found in [this guide](https://github.com/suomitek/suomitek-appboard/blob/master/docs/user/operators.md).

### [Only for Helm 2] Configuring connection to a custom namespace Tiller instance

By default, Suomitek-appboard connects to the Tiller Service in the `kube-system` namespace, the default install location for Helm.

If your instance of Tiller is running in a different namespace or you want to have different instances of Suomitek-appboard connected to different Tiller instances, you can achieve it by setting the `tillerProxy.host` parameter. For example, you can set `tillerProxy.host=tiller-deploy.my-custom-namespace:44134`

### [Only for Helm 2] Configuring connection to a secure Tiller instance

In production, we strongly recommend setting up a [secure installation of Tiller](https://docs.helm.sh/using_helm/#using-ssl-between-helm-and-tiller), the Helm server side component.

Learn more about how to secure your Suomitek-appboard installation [here](https://github.com/suomitek/suomitek-appboard/blob/master/docs/user/securing-suomitek-appboard.md).

### Exposing Externally

> **Note**: The Suomitek-appboard frontend sets up a proxy to the Kubernetes API service which means that when exposing the Suomitek-appboard service to a network external to the Kubernetes cluster (perhaps on an internal or public network), the Kubernetes API will also be exposed for authenticated requests from that network. If you explicitly [use an OAuth2/OIDC provider with Suomitek-appboard](https://github.com/suomitek/suomitek-appboard/blob/master/docs/user/using-an-OIDC-provider.md) (recommended), then only the configured users trusted by your Identity Provider will be able to reach the Kubernetes API. See [#1111](https://github.com/suomitek/suomitek-appboard/issues/1111) for more details.

#### LoadBalancer Service

The simplest way to expose the Suomitek-appboard Dashboard is to assign a LoadBalancer type to the Suomitek-appboard frontend Service. For example, you can use the following parameter: `frontend.service.type=LoadBalancer`

Wait for your cluster to assign a LoadBalancer IP or Hostname to the `suomitek-appboard` Service and access it on that address:

```bash
kubectl get services --namespace suomitek-appboard --watch
```

#### Ingress

This chart provides support for ingress resources. If you have an ingress controller installed on your cluster, such as [nginx-ingress](https://hub.suomitek.com/charts/stable/nginx-ingress) or [traefik](https://hub.suomitek.com/charts/stable/traefik) you can utilize the ingress controller to expose Suomitek-appboard.

To enable ingress integration, please set `ingress.enabled` to `true`

##### Hosts

Most likely you will only want to have one hostname that maps to this Suomitek-appboard installation (use the `ingress.hostname` parameter to set the hostname), however, it is possible to have more than one host. To facilitate this, the `ingress.extraHosts` object is an array.

##### Annotations

For annotations, please see [this document](https://github.com/kubernetes/ingress-nginx/blob/master/docs/user-guide/nginx-configuration/annotations.md). Not all annotations are supported by all ingress controllers, but this document does a good job of indicating which annotation is supported by many popular ingress controllers. Annotations can be set using `ingress.annotations`.

##### TLS

To enable TLS, please set `ingress.tls` to `true`. When enabling this parameter, the TLS certificates will be retrieved from a TLS secret with name *INGRESS_HOSTNAME-tls* (where *INGRESS_HOSTNAME* is a placeholder to be replaced with the hostname you set using the `ingress.hostname` parameter).

You can use the `ingress.extraTls` to provide the TLS configuration for the extra hosts you set using the `ingress.extraHosts` array. Please see [this example](https://kubernetes.github.io/ingress-nginx/examples/tls-termination/) for more information.

You can provide your own certificates using the `ingress.secrets` object. If your cluster has a [cert-manager](https://github.com/jetstack/cert-manager) add-on to automate the management and issuance of TLS certificates, set `ingress.certManager` boolean to true to enable the corresponding annotations for cert-manager. For a full list of configuration parameters related to configuring TLS can see the [values.yaml](values.yaml) file.

## Upgrading Suomitek-appboard

You can upgrade Suomitek-appboard from the Suomitek-appboard web interface. Select the namespace in which Suomitek-appboard is installed (`suomitek-appboard` if you followed the instructions in this guide) and click on the "Upgrade" button. Select the new version and confirm.

You can also use the Helm CLI to upgrade Suomitek-appboard, first ensure you have updated your local chart repository cache:

```bash
helm repo update
```

Now upgrade Suomitek-appboard:

```bash
export RELEASE_NAME=suomitek-appboard
helm upgrade $RELEASE_NAME chartmuseum/suomitek-appboard
```

If you find issues upgrading Suomitek-appboard, check the [troubleshooting](#error-while-upgrading-the-chart) section.

## Uninstalling the Chart

To uninstall/delete the `suomitek-appboard` deployment:

```bash
# For Helm 2
helm delete --purge suomitek-appboard

# For Helm 3
helm uninstall suomitek-appboard

# Optional: Only if there are no more instances of Suomitek-appboard
kubectl delete crd apprepositories.suomitek.com
```

The first command removes most of the Kubernetes components associated with the chart and deletes the release. After that, if there are no more instances of Suomitek-appboard in the cluster you can manually delete the `apprepositories.suomitek.com` CRD used by Suomitek-appboard that is shared for the entire cluster.

> **NOTE**: If you delete the CRD for `apprepositories.suomitek.com` it will delete the repositories for **all** the installed instances of `suomitek-appboard`. This will break existing installations of `suomitek-appboard` if they exist.

If you have dedicated a namespace only for Suomitek-appboard you can completely clean remaining completed/failed jobs or any stale resources by deleting the namespace

```bash
kubectl delete namespace suomitek-appboard
```

## Troubleshooting

### Nginx Ipv6 error

When starting the application with the `--set enableIPv6=true` option, the Nginx server present in the services `suomitek-appboard` and `suomitek-appboard-internal-dashboard` may fail with the following:

```
nginx: [emerg] socket() [::]:8080 failed (97: Address family not supported by protocol)
```

This usually means that your cluster is not compatible with IPv6. To disable it, install suomitek-appboard with the flag: `--set enableIPv6=false`.

### Forbidden error while installing the Chart

If during installation you run into an error similar to:

```
Error: release suomitek-appboard failed: clusterroles.rbac.authorization.k8s.io "suomitek-appboard-apprepository-controller" is forbidden: attempt to grant extra privileges: [{[get] [batch] [cronjobs] [] []...
```

Or:

```
Error: namespaces "suomitek-appboard" is forbidden: User "system:serviceaccount:kube-system:default" cannot get namespaces in the namespace "suomitek-appboard"
```

This usually is an indication that Tiller was not installed with enough permissions to create the resources required by Suomitek-appboard. In order to install Suomitek-appboard, tiller will need to be able to install Custom Resource Definitions cluster-wide, as well as manage app repositories in your suomitek-appboard namespace. The easiest way to enable this in a development environment is install Tiller with elevated permissions (e.g. as a cluster-admin). For example:

```bash
kubectl -n kube-system create sa tiller
kubectl create clusterrolebinding tiller --clusterrole cluster-admin --serviceaccount=kube-system:tiller
helm init --service-account tiller
```

but for a production environment you can assign the specific permissions so that tiller can [manage CRDs on the cluster](https://github.com/suomitek/suomitek-appboard/blob/master/docs/user/manifests/openshift-tiller-with-crd-rbac.yaml) as well as [create app repositories in your Suomitek-appboard namespace](https://github.com/suomitek/suomitek-appboard/blob/master/docs/user/manifests/openshift-tiller-with-apprepository-rbac.yaml) (examples are from our in development support for OpenShift).

It is also possible, though less common, that your cluster does not have Role Based Access Control (RBAC) enabled. To check if your cluster has RBAC you can execute:

```bash
kubectl api-versions
```

If the above command does not include entries for `rbac.authorization.k8s.io` you should perform the chart installation by setting `rbac.create=false`:

```bash
helm install --name suomitek-appboard --namespace suomitek-appboard chartmuseum/suomitek-appboard --set rbac.create=false
```

### Error while upgrading the Chart

It is possible that when upgrading Suomitek-appboard an error appears. That can be caused by a breaking change in the new chart or because the current chart installation is in an inconsistent state. If you find issues upgrading Suomitek-appboard you can follow these steps:

> Note: This steps assume that you have installed Suomitek-appboard in the namespace `suomitek-appboard` using the name `suomitek-appboard`. If that is not the case replace the command with your namespace and/or name.

1.  (Optional) Backup your personal repositories (if you have any):

```bash
kubectl get apprepository --namespace suomitek-appboard -o yaml <repo name> > <repo name>.yaml
```

2.  Delete Suomitek-appboard:

```bash
helm del --purge suomitek-appboard
```

3.  (Optional) Delete the App Repositories CRD:

> **Warning**: Don't execute this step if you have more than one Suomitek-appboard installation in your cluster.

```bash
kubectl delete crd apprepositories.suomitek.com
```

4.  (Optional) Clean the Suomitek-appboard namespace:

> **Warning**: Don't execute this step if you have workloads other than Suomitek-appboard in the `suomitek-appboard` namespace.

```bash
kubectl delete namespace suomitek-appboard
```

5.  Install the latest version of Suomitek-appboard (using any custom modifications you need):

```bash
helm repo update
helm install --name suomitek-appboard --namespace suomitek-appboard chartmuseum/suomitek-appboard
```

6.  (Optional) Restore any repositories you backed up in the first step:

```bash
kubectl apply -f <repo name>.yaml
```

After that you should be able to access the new version of Suomitek-appboard. If the above doesn't work for you or you run into any other issues please open an [issue](https://github.com/suomitek/suomitek-appboard/issues/new).
