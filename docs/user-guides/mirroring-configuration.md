# Mirror Configuration

This page describes how to enable mirroring features on the Vald cluster.

Before you use the mirroring functions, please look at [the Vald Mirror Gateway document](../overview/component/mirror-gateway.md) for what you can do.

## Requirement

- Vald version: v1.8
- The number of Vald clusters: 2~

## Configuration

This chapter shows how to configure values.yaml to enable Vald Mirror Gateway and how to interconnect Vald Mirror Gateways.

The setting points are the followings:

1. Enable the Vald Mirror Gateway using Helm values configuration
2. Interconnect Vald Mirror Gateways using the Custom Resource configuration

### Helm Values Configuration

The Helm values configuration is required for each Vald cluster to be deployed.

It is easy to enable the mirroring feature.

```yaml
---
gateway:
  mirror:
    enabled: true
```

If you want to make more detailed settings, please set the following parameters.

```yaml
gateway:
  mirror:
  ...
    gateway_config:
    ...
      # gRPC client configuration (overrides defaults.grpc.client)
      client: {}
      # The duration to register other Mirror Gateways.
      register_duration: "1s"
      # The target namespace to discover ValdMirrorTarget (CR) resource.
      # The default value is its own namespace.
      namespace: "vald"
      # The group name of the Mirror Gateways (optional).
      # It is used to discover ValdMirrorTarget resources (CR) with the same group name.
      # The default value is empty.
      group: "group1"
      # The duration to discover other mirror gateways in the same group.
      discovery_duration: 1s
      # The colocation name of the data center (optional).
      colocation: "dc1"
```

The cluster role configuration is required when you deploy Vald clusters with Vald Mirror Gateway on multiple Namespaces in the Kubernetes cluster.

Please refer to [Cluster Role Configuration](./cluster-role-binding.md) about cluster role settings for Mirror Gateway.

### Custom Resource Configuration

The Mirror Gateway is not connected to other mirror gateways when deployed.

The Vald Mirror Gateway connects to another Mirror Gateway component specified in the `ValdMirrorTarget` resource (Custom Resource).

Based on this resource, if the connection succeeds, the Mirror Gateway will interconnect with another.

```yaml
apiVersion: vald.vdaas.org/v1
kind: ValdMirrorTarget
metadata:
  name: mirror-target-01
  namespace: vald-03
  labels:
    # The group name of the Mirror Gateways.
    group: mirror-group-01
spec:
  # Colocation name. (optional)
  colocation: dc1
  # The connection target to another mirror gateway.
  target:
    # The hostname.
    host: vald-mirror-gateway.vald-01.svc.cluster.local
    # The port number
    port: 8081
```
