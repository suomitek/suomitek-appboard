# Access Control in Suomitek-appboard

Suomitek-appboard requires users to login with a Kubernetes API token in order to make
requests to the Kubernetes API server as the user. This ensures that a certain
user of Suomitek-appboard is only permitted to view and manage applications that they
have access to (for example, within a specific namespace). If a user does not
have access to a particular resource, Suomitek-appboard will display an error describing
the required roles to access the resource.

If your cluster supports [Token
Authentication](https://kubernetes.io/docs/admin/authentication/) you may login
with the same tokens. Alternatively, you can create Service Accounts for
Suomitek-appboard users. The examples below use a Service Account, as it is the most
common scenario.

## Service Accounts

To create a Service Account for a user "example" in the "default" namespace, run
the following:

```bash
kubectl create -n default serviceaccount example
```

To get the API token for this Service Account, run the following:

```bash
kubectl get -n default secret $(kubectl get -n default serviceaccount example -o jsonpath='{.secrets[].name}') -o go-template='{{.data.token | base64decode}}' && echo
```

## Assigning Suomitek-appboard User Roles

Suomitek-appboard will install a set of preset Roles and ClusterRoles in your cluster
that you can bind to user or Service Accounts. Each Role and ClusterRole
pertains to a certain operation within Suomitek-appboard. This documentation describes
the roles that should be applied to a user in order to perform operations within
Suomitek-appboard.

### Applications

#### Read access to Applications within a namespace

In order to list and view Applications in a namespace, first we will create a `ClusterRole` with read-access to **all** the possible resources. In case you want
to limit this access, create a custom cluster role or use one of the [default ones](https://kubernetes.io/docs/reference/access-authn-authz/rbac/#user-facing-roles). Then we will bind that cluster role to our service account.

```bash
kubectl apply -f https://raw.githubusercontent.com/suomitek/suomitek-appboard/master/docs/user/manifests/suomitek-appboard-applications-read.yaml
kubectl create -n default rolebinding example-view \
  --clusterrole=suomitek-appboard-applications-read \
  --serviceaccount default:example
```

#### Write access to Applications within a namespace

In order to create, update and delete Applications in a namespace, apply the
`edit` ClusterRole in the desired namespace. The `edit` ClusterRole should be
available in most Kubernetes distributions, you can find more information about
that role
[here](https://kubernetes.io/docs/reference/access-authn-authz/rbac/#user-facing-roles).

```bash
kubectl create -n default rolebinding example-edit \
  --clusterrole=edit \
  --serviceaccount default:example
```

### Service Catalog, Service Instances and Bindings

#### Read access to Service Instances and Bindings within a namespaces

Service Brokers, Classes and Plans in the Service Catalog are cluster-scoped
resources, but Service Instances and Bindings can be restricted to a namespace.
We'll need to define two roles (`suomitek-appboard-service-catalog-browse` and
`suomitek-appboard-service-catalog-read`) to separate the roles required to view Service
Instances and Bindings so that they can be applied to desired namespaces.

In order to list and view Service Instances in a namespace, we'll create the
`suomitek-appboard-service-catalog-browse` ClusterRole in all namespaces and the
`suomitek-appboard-service-catalog-read` in the desired namespace.

```bash
kubectl apply -f https://raw.githubusercontent.com/suomitek/suomitek-appboard/master/docs/user/manifests/suomitek-appboard-service-catalog-browse.yaml
kubectl create clusterrolebinding example-suomitek-appboard-service-catalog-browse --clusterrole=suomitek-appboard-service-catalog-browse --serviceaccount default:example

kubectl apply -f https://raw.githubusercontent.com/suomitek/suomitek-appboard/master/docs/user/manifests/suomitek-appboard-service-catalog-read.yaml
kubectl create -n default rolebinding example-suomitek-appboard-service-catalog-read --clusterrole=suomitek-appboard-service-catalog-read --serviceaccount default:example
```

#### Write access to Service Instances and Bindings within a namespace

In order to create and delete Service Instances and Bindings in a namespace,
create and bind the `suomitek-appboard-service-catalog-write` ClusterRole in the desired namespace.

```bash
kubectl apply -f https://raw.githubusercontent.com/suomitek/suomitek-appboard/master/docs/user/manifests/suomitek-appboard-service-catalog-write.yaml
kubectl create -n default rolebinding example-suomitek-appboard-service-catalog-write --clusterrole=suomitek-appboard-service-catalog-write --serviceaccount default:example
```

#### Admin access to configure Service Brokers

In order to resync Service Brokers from the Service Brokers Configuration page,
create and apply the `suomitek-appboard-service-catalog-admin` ClusterRole in all namespaces.

```bash
kubectl apply -f https://raw.githubusercontent.com/suomitek/suomitek-appboard/master/docs/user/manifests/suomitek-appboard-service-catalog-admin.yaml
kubectl create clusterrolebinding example-suomitek-appboard-service-catalog-admin --clusterrole=suomitek-appboard-service-catalog-admin --serviceaccount default:example
```

### App Repositories

#### Read access to App Repositories

In order to list the configured App Repositories in Suomitek-appboard, [bind users/groups Subjects](https://kubernetes.io/docs/reference/access-authn-authz/rbac/#command-line-utilities) to the Suomitek-appboard `apprepositories-read` clusterrole in the namespace Suomitek-appboard was installed into by the helm chart.

```bash
export KUBEAPPS_NAMESPACE=suomitek-appboard
kubectl create -n $KUBEAPPS_NAMESPACE rolebinding example-suomitek-appboard-repositories-read \
  --clusterrole=suomitek-appboard:$KUBEAPPS_NAMESPACE:apprepositories-read \
  --serviceaccount default:example
```

#### Write access to App Repositories

Likewise to the read access bind users/group Subjects to the
Suomitek-appboard `apprepositories-write` ClusterRole in the namespace Suomitek-appboard is installed in
for users to create and refresh App Repositories in Suomitek-appboard

```bash
export KUBEAPPS_NAMESPACE=suomitek-appboard
kubectl create -n $KUBEAPPS_NAMESPACE rolebinding example-suomitek-appboard-repositories-write \
  --clusterrole=suomitek-appboard:$KUBEAPPS_NAMESPACE:apprepositories-write \
  --serviceaccount default:example
```

### Assigning roles across multiple namespaces

To give permissions in multiple namespaces, simply create the same RoleBindings
in each namespace you want to configure access for. For example, to give the
"example" user permissions to manage Applications in the "example" namespace:

```bash
kubectl create -n example rolebinding example-suomitek-appboard-applications-write --clusterrole=suomitek-appboard-applications-read --serviceaccount default:example
kubectl create -n example rolebinding example-suomitek-appboard-applications-write --clusterrole=suomitek-appboard-applications-write --serviceaccount default:example
```

Note that there's no need to recreate the RoleBinding in the namespace Suomitek-appboard
is installed in that is also needed, since that has already been created.

If you want to give access for every namespace, simply create a
ClusterRoleBinding instead of a RoleBinding. For example, to give the "example" user permissions to manage Applications in _any_ namespace:

```bash
kubectl create clusterrolebinding example-suomitek-appboard-applications-write --clusterrole=suomitek-appboard-applications-read --serviceaccount default:example
kubectl create clusterrolebinding example-suomitek-appboard-applications-write --clusterrole=suomitek-appboard-applications-write --serviceaccount default:example
```

## Using a cluster-admin user (not recommended)

A simpler way to configure access for Suomitek-appboard would be to give the user
cluster-admin access (effectively disabling RBAC). This is not recommended, but
useful for quick demonstrations or evaluations.

```bash
kubectl create serviceaccount suomitek-appboard-operator
kubectl create clusterrolebinding suomitek-appboard-operator --clusterrole=cluster-admin --serviceaccount=default:suomitek-appboard-operator
```
