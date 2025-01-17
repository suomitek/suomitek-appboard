# Suomitek-appboard Dashboard Developer Guide

The dashboard is the main UI component of the Suomitek-appboard project. Written in Javascript, the dashboard uses the React Javascript library for the frontend.

## Prerequisites

- [Git](https://git-scm.com/)
- [Node 8.x](https://nodejs.org/)
- [Yarn](https://yarnpkg.com)
- [Kubernetes cluster (v1.8+)](https://kubernetes.io/docs/setup/pick-right-solution/)
- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
- [Docker CE](https://www.docker.com/community-edition)
- [Telepresence](https://telepresence.io)

*Telepresence is not a hard requirement, but is recommended for a better developer experience*

## Environment

```bash
export GOPATH=~/gopath
export PATH=$GOPATH/bin:$PATH
export KUBEAPPS_DIR=$GOPATH/src/github.com/suomitek/suomitek-appboard
```
## Download the suomitek-appboard source code

```bash
git clone --recurse-submodules https://github.com/suomitek/suomitek-appboard $KUBEAPPS_DIR
```

The dashboard application source is located under the `dashboard/` directory of the repository.

```bash
cd $KUBEAPPS_DIR/dashboard
```

### Install Suomitek-appboard in your cluster

Suomitek-appboard is a Kubernetes-native application. To develop and test Suomitek-appboard components we need a Kubernetes cluster with Suomitek-appboard already installed. Follow the [Suomitek-appboard installation guide](../../chart/suomitek-appboard/README.md) to install Suomitek-appboard in your cluster.

### Running the dashboard in development

[Telepresence](https://www.telepresence.io/) is a local development tool for Kubernetes microservices. As the dashboard is a service running in the Kubernetes cluster we use telepresence to proxy requests to the dashboard running in your cluster to your local development host.

First install the dashboard dependency packages:

```bash
yarn install
```

Next, create a `telepresence` shell to swap the `suomitek-appboard-internal-dashboard` deployment in the `suomitek-appboard` namespace, forwarding local port `3000` to port `8080` of the `suomitek-appboard-internal-dashboard` pod.

```bash
telepresence --namespace suomitek-appboard --method inject-tcp --swap-deployment suomitek-appboard-internal-dashboard --expose 3000:8080 --run-shell
```

> **NOTE**: If you encounter issues getting this setup working correctly, please try switching the telepresence proxying method in the above command to `vpn-tcp`. Refer to [the telepresence docs](https://www.telepresence.io/reference/methods) to learn more about the available proxying methods and their limitations. If this doesn't work you can use the [Telepresence alternative](#telepresence-alternative).

Finally, launch the dashboard within the telepresence shell:

```bash
export TELEPRESENCE_CONTAINER_NAMESPACE=suomitek-appboard
yarn run start
```

> **NOTE**: The commands above assume you install Suomitek-appboard in the `suomitek-appboard` namespace. Please update the environment variable `TELEPRESENCE_CONTAINER_NAMESPACE` if you are using a different namespace.

#### Telepresence alternative

As an alternative to using [Telepresence](https://www.telepresence.io/) you can use the default [Create React App API proxy](https://create-react-app.dev/docs/proxying-api-requests-in-development/) functionality.

First add the desired host:port to the package.json:

```patch
-  }
+  },
+  "proxy": "http://127.0.0.1:8080"
```

> **NOTE**: Add the [proxy](../../dashboard/package.json#L176) `key:value` to the end of the `package.json`. For convenience, you can change the `host:port` values to meet your needs.

To use this a run Suomitek-appboard per the [getting-started documentation](../../docs/user/getting-started.md#step-3-start-the-suomitek-appboard-dashboard). This will start Suomitek-appboard running on port `8080`.

Next you can launch the dashboard.

```bash
yarn run start
```

You can now access the local development server simply by accessing the dashboard as you usually would (e.g. doing a port-forward or accesing the Ingress URL).

#### Troubleshooting

In some cases, the 'Create React App' scripts keep listening on the 3000 port, even when you disconnect telepresence. If you see that `localhost:3000` is still serving the dashboard, even with your telepresence down, check if there is a 'Create React App' script process running (`ps aux | grep react`) and kill it.

### Running tests

Execute the following command within the dashboard directory to start the test runner which will watch for changes and automatically re-run the tests when changes are detected.

```bash
yarn run test
```

> **NOTE**: macOS users may need to install watchman (https://facebook.github.io/watchman/).

