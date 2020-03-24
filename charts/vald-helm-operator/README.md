vald-helm-operator
===

This is a Helm chart to install vald-helm-operator.

Current chart version is `v0.0.25`

Install
---

Add Vald Helm repository

    $ helm repo add vald https://vald.vdaas.org/charts

Run the following command to install the chart,

    $ helm install --generate-name vald/vald-helm-operator


Configuration
---

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| crd.create | bool | `true` | ValdRelease crd will be created |
| image.pullPolicy | string | `"Always"` | image pull policy |
| image.repository | string | `"vdaas/vald-helm-operator"` | image repository |
| image.tag | string | `"v0.0.25"` | image tag |
| name | string | `"vald-helm-operator"` | name of the deployment |
| nodeSelector | object | `{}` | node labels for pod assignment |
| rbac.create | bool | `true` | required roles and rolebindings will be created |
| rbac.name | string | `"vald-helm-operator"` | name of roles and rolebindings |
| replicas | int | `1` | number of replicas |
| resources | object | `{}` | k8s resources of pod |
| serviceAccount.create | bool | `true` | service account will be created |
| serviceAccount.name | string | `"vald-helm-operator"` | name of service account |
| vald.create | bool | `false` | ValdRelease resource will be created |
| vald.name | string | `"vald-cluster"` | name of ValdRelease resource |
| vald.spec | object | `{}` | spec field of ValdRelease resource = the values of Helm chart for Vald |
