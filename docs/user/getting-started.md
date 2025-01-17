# Get Started with Suomitek-appboard

This guide will walk you through the process of deploying Suomitek-appboard for your cluster and installing an example application.

## Prerequisites

Suomitek-appboard assumes a working Kubernetes cluster (v1.8+), [`Helm`](https://helm.sh/) (2.14.0+) installed in your cluster and [`kubectl`](https://kubernetes.io/docs/tasks/tools/install-kubectl/) installed and configured to talk to your Kubernetes cluster. Suomitek-appboard has been tested with Azure Kubernetes Service (AKS), Google Kubernetes Engine (GKE), `minikube` and Docker for Desktop Kubernetes. Suomitek-appboard works on RBAC-enabled clusters and this configuration is encouraged for a more secure install.

> On GKE, you must either be an "Owner" or have the "Container Engine Admin" role in order to install Suomitek-appboard.

## Step 1: Install Suomitek-appboard

Use the Helm chart to install the latest version of Suomitek-appboard:

For Helm 2:

```bash
helm repo add chartmuseum http://helm.yongchehang.com
helm install --name suomitek-appboard --namespace suomitek-appboard chartmuseum/suomitek-appboard
```

If you are using Helm 3, you need to set an extra flag to enable it:

```bash
helm repo add chartmuseum http://helm.yongchehang.com
kubectl create namespace suomitek-appboard
helm install suomitek-appboard --namespace suomitek-appboard chartmuseum/suomitek-appboard --set useHelm3=true
```

For detailed information on installing, configuring and upgrading Suomitek-appboard, checkout the [chart README](../../chart/suomitek-appboard/README.md).

The above commands will deploy Suomitek-appboard into the `suomitek-appboard` namespace in your cluster. It may take a few minutes to execute. Once it has been deployed and the Suomitek-appboard pods are running, continue to step 2.

## Step 2: Create a Kubernetes API token

For trying out Suomitek-appboard, access to the Dashboard requires a Kubernetes API token to authenticate with the Kubernetes API server as shown below, but for any real installation of Suomitek-appboard you should instead [configure an OAuth2/OIDC provider](using-an-OIDC-provider.md).

```bash
kubectl create serviceaccount suomitek-appboard-operator
kubectl create clusterrolebinding suomitek-appboard-operator --clusterrole=cluster-admin --serviceaccount=default:suomitek-appboard-operator
```

> **NOTE** It's not recommended to create `cluster-admin` users for Suomitek-appboard production usage. Please refer to the [Access Control](/docs/user/access-control.md) documentation to configure fine-grained access control for users.

To retrieve the token,

### On Linux/macOS:

```bash
kubectl get secret $(kubectl get serviceaccount suomitek-appboard-operator -o jsonpath='{range .secrets[*]}{.name}{"\n"}{end}' | grep suomitek-appboard-operator-token) -o jsonpath='{.data.token}' -o go-template='{{.data.token | base64decode}}' && echo
```

### On Windows:

Create a file called `GetDashToken.cmd` with the following lines in it:

```bat
@ECHO OFF
REM Get the Service Account
kubectl get serviceaccount suomitek-appboard-operator -o jsonpath={.secrets[].name} > s.txt
SET /p ks=<s.txt
DEL s.txt

REM Get the Base64 encoded token
kubectl get secret %ks% -o jsonpath={.data.token} > b64.txt

REM Decode The Token
DEL token.txt
certutil -decode b64.txt token.txt
```

Open a command prompt and run the `GetDashToken.cmd` Your token can be found in the `token.txt` file.

## Step 3: Start the Suomitek-appboard Dashboard

Once Suomitek-appboard is installed, securely access the Suomitek-appboard Dashboard from your system by running:

```bash
kubectl port-forward -n suomitek-appboard svc/suomitek-appboard 8080:80
```

This will start an HTTP proxy for secure access to the Suomitek-appboard Dashboard. Visit http://127.0.0.1:8080/ in your preferred web browser to open the Dashboard. Here's what you should see:

![Dashboard login page](../img/dashboard-login.png)

Paste the token generated in the previous step to authenticate and access the Suomitek-appboard dashboard for Kubernetes.

![Dashboard main page](../img/dashboard-home.png)

***Note:*** If you are setting up Suomitek-appboard for other people to access, you will want to use a different service type or setup Ingress rather than using the above `kubectl port-forward`. For detailed information on installing, configuring and upgrading Suomitek-appboard, checkout the [chart README](../../chart/suomitek-appboard/README.md).

## Step 4: Deploy WordPress

Once you have the Suomitek-appboard Dashboard up and running, you can start deploying applications into your cluster.

- Use the "Deploy App" or click on the "Catalog" page in the Dashboard to select an application from the list of charts in any of the configured Helm chart repositories. This example assumes you want to deploy WordPress.

  ![WordPress chart](../img/wordpress-search.png)

- Click the "Deploy" button.

  ![WordPress chart](../img/wordpress-chart.png)

- You will be prompted for the release name and values for the application. The form is populated by the values (YAML), which you can see in the adjacent tab.

  ![WordPress installation](../img/wordpress-installation.png)

- Click the "Submit" button. The application will be deployed. You will be able to track the new Helm deployment directly from the browser. The status will be shown at the top and you can also look at the individual resources lower in the page. It will also show the number of ready pods. If you run your cursor over the status, you can see the workloads and number of ready and total pods within them.

  ![WordPress deployment](../img/wordpress-deployment.png)

To access your new WordPress site, you can run the commands in the "Notes" section to get the URLs or simply click a URL (HTTP and HTTPS) shown.

***Note:*** Depending on your cloud provider of choice, it may take some time for an access URL to be available for the application and the Service will stay in a "Pending" state until a URL is assigned. If using Minikube, you will need to run `minikube tunnel` in your terminal in order for an IP address to be assigned to your application.

![WordPress deployment notes](../img/wordpress-url.png)

To get the credentials for logging into your WordPress account, refer to the "Notes" section. You can also get the WordPress password by scrolling down to "Secrets" and clicking the eye next to `wordpress-password`.

![WordPress deployment notes](../img/wordpress-credentials.png)

## [Optional] Step 5: Uninstall/Delete WordPress
If you want to uninstall/delete your WordPress application, you can do so by clicking the "Delete" button. You can choose to click the checkbox for "Purge Release" (default action with the Helm 3 CLI). If you do not click it, the Helm chart history will remain (default action with Helm 2). This is fine, so long as you don't attempt to install another chart with the same name in the same namespace.

![WordPress uninstall](../img/wordpress-uninstall.png)

## Next Steps

Learn more about Suomitek-appboard with the links below:

- [Detailed installation instructions](../../chart/suomitek-appboard/README.md)
- [Deploying Operators](./operators.md)
- [Suomitek-appboard Dashboard documentation](dashboard.md)
- [Suomitek-appboard components](../architecture/overview.md)
- [Roadmap](https://github.com/suomitek/suomitek-appboard/wiki/Roadmap)
