# Cluster Role Configuration

Cluster role contains rules that representing a set of permission to grant access to specfic target depending on the binding rule.

This page describe why do we need cluster role for Vald cluster.

## What is cluster role and cluster role binding for Vald cluster?

In Vald, the index is distributed across the cluster depending on the resource usage of the cluster, it requires settings to grant permission to specific role to retrieve cluster information on Kuberenetes.

In Vald, the settings are deployed automatically when using helm to deploy.

- [clusterrole.yaml](https://github.com/vdaas/vald/blob/main/k8s/discoverer/clusterrole.yaml)
- [clusterrolebinding.yaml](https://github.com/vdaas/vald/blob/main/k8s/discoverer/clusterrolebinding.yaml)

These configurations allows service account `discoverer` to access different resources in Kuberenetes cluster.

### Cluster role settings

By looking at the [cluster role settings](https://github.com/vdaas/vald/blob/main/k8s/discoverer/clusterrole.yaml), the access right of the following resources are granted to service account `discoverer`.

```yaml
rules:
  - apiGroups:
      - apps
    resources:
      - replicasets
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - ""
    resources:
      - nodes
      - nodes/proxy
      - services
      - endpoints
      - pods
    verbs:
      - get
      - list
      - watch
  - nonResourceURLs:
      - /metrics
    verbs:
      - get
  - apiGroups:
      - "metrics.k8s.io"
    resources:
      - nodes
      - pods
    verbs:
      - get
      - list
```

All of these rules are required to retrieve Nodes and Pods resource usage from [kube-apiserver](https://kubernetes.io/docs/reference/command-line-tools-reference/kube-apiserver/), and also used to discover new Vald Agents or Nodes created on the cluster.

### Cluster role binding settings

The above cluster role settings will be bind to the service account `vald` according to the [configuration file](https://github.com/vdaas/vald/blob/main/k8s/discoverer/clusterrolebinding.yaml).

```yaml
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: discoverer
subjects:
  - kind: ServiceAccount
    name: vald
    namespace: default
```

This service account will be [used in Vald Discoverer](https://github.com/vdaas/vald/blob/main/k8s/discoverer/deployment.yaml#L155) to helps automate operations for the Vald cluster. For more information about Vald Discoverer, please refer [here](../overview/component/discoverer.md).

## Customize cluster role configuration

To customize cluster role configuration, you may need to change the `discoverer.clusterRole` configuration on helm chart. The cluster role configurations are enabled by default.

```yaml
  clusterRole:
    # discoverer.clusterRole.enabled -- creates clusterRole resource
    enabled: true
    # discoverer.clusterRole.name -- name of clusterRole
    name: discoverer
  clusterRoleBinding:
    # discoverer.clusterRoleBinding.enabled -- creates clusterRoleBinding resource
    enabled: true
    # discoverer.clusterRoleBinding.name -- name of clusterRoleBinding
    name: discoverer
  serviceAccount:
    # discoverer.serviceAccount.enabled -- creates service account
    enabled: true
    # discoverer.serviceAccount.name -- name of service account
    name: vald
```

Please note that the [Vald Discoverer](../overview/component/discoverer.md) will not work after disabling these settings, and Vald cluster will not work.

If you really want to modify or disable these settings, you need to grant the [cluster role settings](https://github.com/vdaas/vald/blob/main/k8s/discoverer/clusterrole.yaml) and bind it to the Vald Discoverer to retrieve required information to operate the Vald cluster.

## Related Document

- [Vald Discoverer](../overview/component/discoverer.md)
