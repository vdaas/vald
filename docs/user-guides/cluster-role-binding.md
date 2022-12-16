# Cluster Role Configuration

The [cluster role](https://kubernetes.io/docs/reference/access-authn-authz/rbac/) feature provided by Kubernetes contains rules that represent a set of permission to grant access to a specific target depending on the binding rule.

This page describes why we need a cluster role for the Vald cluster and how to configure it.

## What are cluster role and cluster role binding for the Vald cluster?

Vald applies the distributed index system across the Kubernetes cluster depending on the resource usage of the Kubernetes Node, it requires configuration to grant permission to a specific role to retrieve cluster information on Kubernetes.

By default, the cluster role configurations are deployed automatically when using Helm.

The following manifest will be deployed by default.

- [clusterrole.yaml](https://github.com/vdaas/vald/blob/main/k8s/discoverer/clusterrole.yaml)
- [clusterrolebinding.yaml](https://github.com/vdaas/vald/blob/main/k8s/discoverer/clusterrolebinding.yaml)

These configurations allow the service account `discoverer`, which is for the Vald Discoverer components, to access different resources in the Kubernetes cluster.

This service account allows [Vald Discoverer](../overview/component/discoverer.md) to retrieve Node and Pod resource usage from [kube-apiserver](https://kubernetes.io/docs/reference/command-line-tools-reference/kube-apiserver/), and share it to other components in the Vald cluster.

For example, [Vald LB Gateway](../overview/component/lb-gateway.md) will control which Vald Agent to insert based on the Node and Pod resource usage retrieved by Vald Discoverer.

If you are interested, please refer to the [insert data flow](../overview/data-flow.md#insert) for more detail.

## Configuration for Vald Discoverer

As described in the above section, Vald Discoverer requires configuration on cluster role and cluster role binding to retrieve Node and Pod information from the Kubernetes Cluster.

In this section, we will describe how to configure it and how to customize these configurations.

### Cluster role configuration for Vald Discoverer

By looking at the [cluster role configuration](https://github.com/vdaas/vald/blob/main/k8s/discoverer/clusterrole.yaml), the access right of the following resources are granted to the cluster role `discoverer`.

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: discoverer
---
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

All of these rules are required to retrieve Node and Pod resource usage from [kube-apiserver](https://kubernetes.io/docs/reference/command-line-tools-reference/kube-apiserver/) and also used to discover new Vald Agent Pods or Nodes created on the cluster.

### Cluster role binding configuration for Vald Discoverer

The cluster role binding configuration binds the cluster role `discoverer` described in the previous section to the service account `vald` according to the [configuration file](https://github.com/vdaas/vald/blob/main/k8s/discoverer/clusterrolebinding.yaml).

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: discoverer
  ...
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: discoverer
subjects:
  - kind: ServiceAccount
    name: vald
    namespace: default
```

When the role binds to the service account, the access right of the role will be granted to the service account.
In this case, all the access rights of the role `discoverer` will be granted to the service account `vald`.

The service account `vald` is required for [Vald Discoverer](https://github.com/vdaas/vald/blob/main/k8s/discoverer/deployment.yaml#L155) to retrieve the required information to operate the Vald cluster.

For more information about Vald Discoverer, please refer [here](../overview/component/discoverer.md).

## Customize cluster role and cluster role binding configuration on Helm chart for Vald Discoverer

To customize the cluster role configuration on the Helm chart for Vald Discoverer, you may need to change the `discoverer.clusterRole` configuration on the Helm chart file. The cluster role configurations are enabled by default.

```yaml
discoverer:
---
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

<div class="warning">
If you disable these configurations, the Vald Discoverer will not work, and the Vald cluster will not be functional.
</div>

If you want to modify or disable these configurations, you need to grant the [cluster role configuration](https://github.com/vdaas/vald/blob/main/k8s/discoverer/clusterrole.yaml) and bind it to the Vald Discoverer to retrieve required information to operate the Vald cluster.

## Customize cluster role configuration on Cloud Providers

Please refer to the official guidelines to configure cluster role configuration for your cloud provider, and configure the service account name for Vald Discoverer.

For example:

- [Amazon EKS](https://docs.aws.amazon.com/eks/latest/userguide/add-user-role.html)
- [GKE](https://cloud.google.com/kubernetes-engine/docs/how-to/role-based-access-control)

For other cloud providers, you may need to find the related document on their official website, or you can enable the cluster role and the cluster role binding configurations on the Helm chart.

## Related Documents

- [Vald Discoverer](../overview/component/discoverer.md)
- [Data Flow](../overview/data-flow.md)
- [Using RBAC Authorization (External Site)](https://kubernetes.io/docs/reference/access-authn-authz/rbac/)
