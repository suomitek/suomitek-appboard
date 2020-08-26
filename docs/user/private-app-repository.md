# Using a Private Repository with Kubeapps

It is possible to use a private Helm repository to store your own Helm charts and deploy them using Kubeapps. In this guide we will show how you can do that with some of the solutions available right now:

- [ChartMuseum](#chartmuseum)
- [Harbor](#harbor)
- [Artifactory](#artifactory) (Pro)

But first, a note about Kubeapps AppRepository resources:

## Per Namespace App Repositories

Previously, once an App Repository was created in Kubeapps, the charts indexed by that repository were then available cluster-wide to all Kubeapps users. This was changed in Kubeapps 1.10 to allow creating App Repositories that are available only in specific namespaces, which is more inline with the Kubernetes RBAC model where an account can have roles in specific namespaces. This change also enables Kubeapps to support deploying charts with images from private docker registries (more below).

A Kubeapps AppRepository can be created by anyone with the required RBAC for that namespace. If you have cluster-wide RBAC for creating AppRepositories, you can still create an App Repository whose charts will be available to users in all namespaces by selecting "All Namespaces" when creating the repository.

To give a specific user `USERNAME` the ability to create App Repositories in a specific namespace named `custom-namespace`, grant them both read and write RBAC for AppRepositories in that namespace:

```bash
kubectl -n custom-namespace create rolebinding username-apprepositories-read --user $USERNAME --clusterrole kubeapps:$KUBEAPPS_NAMESPACE:apprepositories-read
kubectl -n custom-namespace create rolebinding username-apprepositories-write --user $USERNAME --clusterrole kubeapps:$KUBEAPPS_NAMESPACE:apprepositories-write
```

or to allow other users the ability to deploy charts from App Repositories in a specific namespace, grant the read access only.

## Associating docker image pull secrets to an AppRepository - Helm 3 only

When creating an AppRepository in Kubeapps, you can now additionally choose (or create) an [imagePullSecret](https://kubernetes.io/docs/concepts/containers/images/#specifying-imagepullsecrets-on-a-pod) to be associated with the AppRepository:

<img src="../img/app-repo-pull-secret.png" alt="AppRepository with imagePullSecret" width="600px">

When Kubeapps deploys any chart from this AppRepository, if a referenced docker image within the chart is from a docker registry server matching one of the secrets associated with the AppRepository, then Kubeapps with Helm 3 will automatically append the corresponding imagePullSecret so that image can be pulled from the private registry. Note that the user deploying the chart will need to be able to read secrets in that namespace, which is usually the case when deploying to a namespace.

> **Note**: The automatic inclusion of associated image pull secrets is a feature specific to Helm 3. We do not intend to add support for this for Kubeapps configured with Helm 2.

There will be further work to enable private AppRepositories to be available in multiple namespaces. Details about the design can be read on the [design document](https://docs.google.com/document/d/1YEeKC6nPLoq4oaxs9v8_UsmxrRfWxB6KCyqrh2-Q8x0/edit?ts=5e2adf87).

## ChartMuseum

[ChartMuseum](https://chartmuseum.com) is an open-source Helm Chart Repository written in Go (Golang), with support for cloud storage backends, including Google Cloud Storage, Amazon S3, Microsoft Azure Blob Storage, Alibaba Cloud OSS Storage and OpenStack Object Storage.

To use ChartMuseum with Kubeapps, first deploy its Helm chart from the `stable` repository:

<img src="../img/chartmuseum-chart.png" alt="ChartMuseum Chart" width="300px">

In the deployment form we should change at least two things:

- `env.open.DISABLE_API`: We should set this value to `false` so we can use the ChartMuseum API to push new charts.
- `persistence.enabled`: We will set this value to `true` to enable persistence for the charts we store. Note that this will create a [Kubernetes Persistent Volume Claim](https://kubernetes.io/docs/concepts/storage/persistent-volumes/#lifecycle-of-a-volume-and-claim) so depending on your Kubernetes provider you may need to manually allocate the required Persistent Volume to satisfy the claim. Some Kubernetes providers will automatically create PVs for you so setting this value to `true` will be enough.

<img src="../img/chartmuseum-deploy-form.png" alt="ChartMuseum Deploy Form" width="600px">

### ChartMuseum: Upload a Chart

Once ChartMuseum is deployed you will be able to upload a chart. In one terminal open a port-forward tunnel to the application:

```console
$ export POD_NAME=$(kubectl get pods --namespace default -l "app=chartmuseum" -l "release=my-chartrepo" -o jsonpath="{.items[0].metadata.name}")
$ kubectl port-forward $POD_NAME 8080:8080 --namespace default
Forwarding from 127.0.0.1:8080 -> 8080
Forwarding from [::1]:8080 -> 8080
```

And in a different terminal you can push your chart:

```console
$ helm package /path/to/my/chart
Successfully packaged chart and saved it to: /path/to/my/chart/my-chart-1.0.0.tgz
curl --data-binary "@my-chart-1.0.0.tgz" http://localhost:8080/api/charts
{"saved":true}
```

### ChartMuseum: Configure the repository in Kubeapps

To add your private repository to Kubeapps, select the Kubernetes namespace to which you want to add the repository (or "All Namespaces" if you want it available to users in all namespaces), go to `Configuration > App Repositories` and click on "Add App Repository". You will need to add your repository using the Kubernetes DNS name for the ChartMuseum service. This will be `<release_name>-chartmuseum.<namespace>:8080`:

<img src="../img/chartmuseum-repository.png" alt="ChartMuseum App Repository" width="300px">

Once you create the repository you can click on the link for the specific repository and you will be able to deploy your own applications using Kubeapps.

### ChartMuseum: Authentication/Authorization

It is possible to configure ChartMuseum to use authentication with two different mechanism:

- Using HTTP [basic authentication](https://chartmuseum.com/docs/#basic-auth) (user/password). To use this feature, it's needed to:
  - Specify the parameters `secret.AUTH_USER` and `secret.AUTH_PASS` when deploying the ChartMuseum.
  - Select `Basic Auth` when adding the repository to Kubeapps specifying that user and password.
- Using a [JWT token](https://github.com/chartmuseum/auth-server-example). Once you obtain a valid token you can select `Bearer Token` in the form and add the token in the dedicated field.

## Harbor

[Harbor](https://github.com/goharbor/harbor) is an open source trusted cloud native registry project that stores, signs, and scans content, e.g. Docker images. Harbor is hosted by the [Cloud Native Computing Foundation](https://cncf.io/). Since version 1.6.0, Harbor is a composite cloud native registry which supports both container image management and Helm chart management. Harbor integrates [ChartMuseum](https://chartmuseum.com) to provide the Helm chart repository functionality. The access to Helm Charts in a Harbor Chart Repository can be controlled via Role-Based Access Control.

To use Harbor with Kubeapps, first deploy the [Bitnami Harbor Helm Chart](https://github.com/bitnami/charts/tree/master/bitnami/harbor) from the `bitnami` repository (alternatively you can deploy Harbor using [Harbor offline installer](https://github.com/goharbor/harbor/blob/master/docs/installation_guide.md#downloading-the-installer)):

<img src="../img/harbor-chart.png" alt="Harbor Chart" width="300px">

In the deployment form we should change the parameter below:

- `service.tls.enabled`: We should set this value to `false` so we don't need to configure the TLS settings. Alternatively, you can provide valid TSL certificates (check [Bitnami Harbor Helm Chart documentation](https://github.com/bitnami/charts/tree/master/bitnami/harbor#parameters) for more information).

<img src="../img/harbor-deploy-form.png" alt="Harbor Deploy Form" width="600px">

Deploy the chart and wait for it te be ready.

<img src="../img/harbor-ready.png" alt="Harbor Chart Ready" width="600px">

### Harbor: Upload a Chart

First create a Helm chart package:

```console
$ helm package /path/to/my/chart
Successfully packaged chart and saved it to: /path/to/my/chart/my-chart-1.0.0.tgz
```

Second login Harbor admin portal following the instructions in the chart notes:

```console
1. Get the Harbor URL:

  echo "Harbor URL: https://127.0.0.1:8080/"
  kubectl port-forward --namespace default svc/my-harbor 8080:80 &

2. Login with the following credentials to see your Harbor application

  echo Username: "admin"
  echo Password: $(kubectl get secret --namespace default my-harbor-core-envvars -o jsonpath="{.data.HARBOR_ADMIN_PASSWORD}" | base64 --decode)
```

Create a new Project named 'my-helm-repo' with public access. Each project will serve as a Helm chart repository.

<img src="../img/harbor-new-project.png" width="300px">

Click the project name to view the project details page, then click 'Helm Charts' tab to list all helm charts.

<img src="../img/harbor-list-charts.png" width="600px">

Click 'UPLOAD' button to upload the Helm chart you previously created. You can also use helm command to upload chart too.

<img src="../img/harbor-upload-chart.png" width="500px">

Please refer to ['Manage Helm Charts in Harbor'](https://github.com/goharbor/harbor/blob/master/docs/user_guide.md#manage-helm-charts) for more details.

### Harbor: Configure the repository in Kubeapps

To add Harbor as the private chart repository in Kubeapps, select the Kubernetes namespace to which you want to add the repository (or "All Namespaces" if you want it available to users in all namespaces), go to `Configuration > App Repositories` and click on "Add App Repository" and use the Harbor helm repository URL `http://harbor.default.svc.cluster.local/chartrepo/my-helm-repo`

<img src="../img/harbor-add-repo.png" width="600px">

Once you create the repository you can click on the link for the specific repository and you will be able to deploy your own applications using Kubeapps.

### Harbor: Authentication/Authorization

It is possible to configure Harbor to use HTTP basic authentication:

  - When creating a new project for serving as the helm chart repository in Harbor, set the `Access Level` of the project to non public. This enforces authentication to access the charts in the chart repository via Helm CLI or other clients.
  - When `Adding App Repository` in Kubeapps, select `Basic Auth` for `Authorization` and specifiy the username and password for Harbor.

## Artifactory

JFrog Artifactory is a Repository Manager supporting all major packaging formats, build tools and CI servers.

> **Note**: In order to use the Helm repository feature, it's necessary to use an Artifactory Pro account.

To install Artifactory with Kubeapps first add the JFrog repository to Kubeapps. Go to `Configuration > App Repositories` and add their repository:

<img src="../img/jfrog-repository.png" alt="JFrog repository" width="300px">

Then click on the JFrog repository and deploy Artifactory. For detailed installation instructions, check its [README](https://github.com/jfrog/charts/tree/master/stable/artifactory). If you don't have any further requirement, the default values will work.

When deployed, in the setup wizard, select "Helm" to initialize a repository:

<img src="../img/jfrog-wizard.png" alt="JFrog repository" width="600px">

By default, Artifactory creates a chart repository called `helm`. That is the one you can use to store your applications.

### Artifactory: Upload a chart

First, you will need to obtain the user and password of the Helm repository. To do so, click on the `helm` repository and in the `Set Me Up` menu enter your password. After that you will be able to see the repository user and password.

Once you have done that, you will be able to upload a chart:

```bash
curl -u{USER}:{PASSWORD} -T /path/to/chart.tgz "http://{REPO_URL}/artifactory/helm/"
```

### Artifactory: Configure the repository in Kubeapps

To be able able to access private charts with Kubeapps first you need to generate a token. You can do that with the Artifactory API:

```bash
curl -u{USER}:{PASSWORD} -XPOST "http://{REPO_URL}/artifactory/api/security/token?expires_in=0" -d "username=kubeapps" -d "scope=member-of-groups:readers"
{
  "scope" : "member-of-groups:readers api:*",
  "access_token" : "TOKEN CONTENT",
  "token_type" : "Bearer"
}
```

The above command creates a token with read-only permissions. Now you can select the namespace to which you want to add the repository (or "All Namespaces" if you want it available to users in all namespaces), go to the `Configuration > App Repositories` menu and add your personal repository:

<img src="../img/jfrog-custom-repo.png" alt="JFrog custom repository" width="400px">

After submitting the repository, you will be able to click on the new repository and see the chart you uploaded in the previous step.

## Modifying the synchronization job

Kubeapps runs a periodic job (CronJob) to populate and synchronize the charts existing in each repository. Since Kubeapps v1.4.0, it's possible to modify the spec of this job. This is useful if you need to run the Pod in a certain Kubernetes node, or set some environment variables. To do so you can edit (or create) an AppRepository and specify the `syncJobPodTemplate` field. For example:

```yaml
apiVersion: kubeapps.com/v1alpha1
kind: AppRepository
metadata:
  name: my-repo
  namespace: kubeapps
spec:
  syncJobPodTemplate:
    metadata:
      labels:
        my-repo: "isPrivate"
    spec:
      containers:
        - env:
            - name: FOO
              value: BAR
  url: https://my.charts.com/
```

The above will generate a Pod with the label `my-repo: isPrivate` and the environment variable `FOO=BAR`.
