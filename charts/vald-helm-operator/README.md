vald-helm-operator
===

This is a Helm chart to install vald-helm-operator.

Current chart version is `v0.0.66`

Table of Contents
---

- [Install](#install)
- [Custom Resources](#custom-resources)
    - [ValdRelease](#valdrelease)
    - [ValdHelmOperatorRelease](#valdhelmoperatorrelease)
- [Configuration](#configuration)

Install
---

Add Vald Helm repository

    $ helm repo add vald https://vald.vdaas.org/charts

Run the following command to install the chart,

    $ helm install vald-helm-operator-release vald/vald-helm-operator

Custom Resources
---

### ValdRelease

This is a custom resource that represents values of the Vald Helm chart.

Example:

```yaml
apiVersion: vald.vdaas.org/v1
kind: ValdRelease
metadata:
  name: vald-cluster
# the values of Helm chart for Vald can be placed under the `spec` field.
spec: {}
```

### ValdHelmOperatorRelease

This is a custom resource that represents values of the vald-helm-operator Helm chart.

Example:

```yaml
apiVersion: vald.vdaas.org/v1
kind: ValdHelmOperatorRelease
metadata:
  name: vald-helm-operator-release
# the values of Helm chart for vald-helm-operator can be placed under the `spec` field.
spec: {}
```

Configuration
---

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| affinity | object | `{}` | affinity |
| annotations | object | `{}` | deployment annotations |
| enableLeaderElection | bool | `true` | enable leader election for controller manager. |
| enableMetrics | bool | `true` | enable metrics endpoint |
| image.pullPolicy | string | `"Always"` | image pull policy |
| image.repository | string | `"vdaas/vald-helm-operator"` | image repository |
| image.tag | string | `"v0.0.66"` | image tag |
| leaderElectionID | string | `"vald-helm-operator"` | name of the configmap that is used for holding the leader lock. |
| logging.format | string | `"console"` | logging format of operator (console or json) |
| logging.level | string | `"info"` | logging level of operator (debug, info, or error) |
| logging.stacktraceLevel | string | `"error"` | minimum log level triggers stacktrace generation |
| maxConcurrentReconciles | int | `1` | max number of concurrent reconciles |
| metricsPort | int | `6061` | port of metrics endpoint |
| name | string | `"vald-helm-operator"` | name of the deployment |
| namespaced | bool | `true` | if it is true, operator will behave as a namespace-scoped operator, if it is false, it will behave as a cluster-scoped operator. |
| nodeSelector | object | `{}` | node labels for pod assignment |
| podAnnotations | object | `{}` | pod annotations |
| rbac.create | bool | `true` | required roles and rolebindings will be created |
| rbac.name | string | `"vald-helm-operator"` | name of roles and rolebindings |
| reconcilePeriod | string | `"1m"` | reconcile duration of operator |
| replicas | int | `2` | number of replicas |
| resources | object | `{}` | kubernetes resources of pod |
| service.annotations | object | `{}` | service annotations |
| service.enabled | bool | `true` | service enabled |
| service.externalTrafficPolicy | string | `""` | external traffic policy (can be specified when service type is LoadBalancer or NodePort) : Cluster or Local |
| service.labels | object | `{}` | service labels |
| service.type | string | `"ClusterIP"` | service type: ClusterIP, LoadBalancer or NodePort |
| serviceAccount.create | bool | `true` | service account will be created |
| serviceAccount.name | string | `"vald-helm-operator"` | name of service account |
| tolerations | list | `[]` | tolerations |
| watchNamespaces | string | `""` | comma separated names of namespaces to watch, if it is empty, the namespace that the operator exists in is used. |
