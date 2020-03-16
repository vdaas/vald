vald-helm-operator
===

This is a Helm chart to install vald-helm-operator.


Install
---

Add Vald Helm repository

    $ helm repo add vald https://vald.vdaas.org/charts

Run the following command to install the chart,

    $ helm install --generate-name vald/vald-helm-operator


Configuration
---

| Parameter | Description | Default |
|-----------|-------------|---------|
| `name` | name of the deployment | `vald-helm-operator` |
| `replicas` | number of replicas | `1` |
| `image.repository` | image repository | `vdaas/vald-helm-operator` |
| `image.tag` | image tag | version |
| `image.pullPolicy` | image pull policy | `Always` |
| `crd.create` | ValdRelease crd will be created | `true` |
| `vald.create` | ValdRelease resource will be created | `false` |
| `vald.name` | name of ValdRelease resource | `vald-cluster` |
| `vald.spec` | spec field of ValdRelease resource = the values of Helm chart for Vald | `{}` |
| `rbac.create` | required roles and rolebindings will be created | `true` |
| `rbac.name` | name of roles and rolebindings | `vald-helm-operator` |
| `serviceAccount.create` | service account will be created | `false` |
| `serviceAccount.name` | name of service account | `vald-helm-operator` |
| `resources` | k8s resources of pod | `{}` |
| `nodeSelector` | node labels for pod assignment | `{}` |
