vald-helm-operator
===

This is a Helm chart to install vald-helm-operator.

Current chart version is `v0.0.37`

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
apiVersion: vald.vdaas.org/v1alpha1
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
apiVersion: vald.vdaas.org/v1alpha1
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
| image.pullPolicy | string | `"Always"` | image pull policy |
| image.repository | string | `"vdaas/vald-helm-operator"` | image repository |
| image.tag | string | `"v0.0.37"` | image tag |
| logging.format | string | `"console"` | logging format of operator (console or json) |
| logging.level | string | `"info"` | logging level of operator (debug, info, or error) |
| logging.stacktraceLevel | string | `"error"` | minimum log level triggers stacktrace generation |
| logging.timeEncoding | string | `"iso8601"` | logging time format of operator (epoch, millis, nano, or iso8601) |
| maxWorkers | int | `1` | number of workers inside one operator pod |
| name | string | `"vald-helm-operator"` | name of the deployment |
| nodeSelector | object | `{}` | node labels for pod assignment |
| rbac.create | bool | `true` | required roles and rolebindings will be created |
| rbac.name | string | `"vald-helm-operator"` | name of roles and rolebindings |
| reconcilePeriod | string | `"1m"` | reconcile duration of operator |
| replicas | int | `2` | number of replicas |
| resources | object | `{}` | k8s resources of pod |
| serviceAccount.create | bool | `true` | service account will be created |
| serviceAccount.name | string | `"vald-helm-operator"` | name of service account |
| tolerations | list | `[]` | tolerations |
