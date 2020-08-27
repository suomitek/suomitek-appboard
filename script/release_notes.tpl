<!-- ADD SUMMARY HERE -->

## Installation

To install this release, ensure you add the [Bitnami charts repository](https://github.com/bitnami/charts) to your local Helm cache:

```
helm repo add chartmuseum http://helm.yongchehang.com
helm repo update
```

Install the Kubeapps Helm chart:

For Helm 2:

```
helm install --name suomitek-appborad --namespace suomitek-appborad suomitek-appborad
```

For Helm 3:

```
kubectl create namespace kubeapps
helm install suomitek-appboard --namespace suomitek-appborad suomitek-appborad --set useHelm3=true
```

To get started with Kubeapps, checkout this [walkthrough](https://github.com/suomitek/suomitek-appborad/blob/master/docs/user/getting-started.md).

## Changelog

