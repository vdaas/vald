# Vald Mirror Gateway

Vald Mirror Gateway is an optional component of the Vald, allowing the vector data to be synchronized across multiple Vald clusters.

This component makes it possible to enhance availability during a cluster failure.

<img src="../../../assets/docs/overview/component/mirror-gateway/mirror-gateway.png">

## Responsibility

Vald Mirror Gateway is responsible for the followings:

- Forward user requests ([Insert](https://vald.vdaas.org/docs/api/insert/) / [Upsert](https://vald.vdaas.org/docs/api/upsert/) / [Update](https://vald.vdaas.org/docs/api/update/) / [Remove](https://vald.vdaas.org/docs/api/remove/)) to the other Vald Mirror Gateways in the same group.
- Manage so that each index state will be the same.

## Features

This chapter shows the main features to fulfill Vald Mirror Gateway’s role:

- Full mesh connection
- Request forwarding
- Automatic rollback on failure

### Full mesh connection

<img src="../../../assets/docs/overview/component/mirror-gateway/full-mesh-connection.png">

The Vald Mirror Gateway is designed to interconnect with Vald Mirror Gateways in other Vald clusters.

Vald Mirror Gateway uses a Custom Resource called the `ValdMirrorTarget` to manage the connection destination information between Vald Mirror Gateways.

The `ValdMirrorTarget` is a Custom Resource related to the connection destination to other Vald Mirror Gateway.

When two Vald clusters contain Vald Mirror Gateways, Vald Mirror Gateways can send the request to each other by applying `ValdMirrorTarget`.

For more information about `ValdMirrorTarget` configuration, please refer to [Custom Resource Configuration](https://vald.vdaas.org/docs/user-guides/mirroring-configuration/).

### Request forwarding

<img src="../../../assets/docs/overview/component/mirror-gateway/request-forwarding.png">

The Vald Mirror Gateway forwards the incoming user request ([Insert](https://vald.vdaas.org/docs/api/insert/) / [Upsert](https://vald.vdaas.org/docs/api/upsert/) / [Update](https://vald.vdaas.org/docs/api/update/) / [Remove](https://vald.vdaas.org/docs/api/remove/)) to other Vald Mirror Gateways.
Then, while forwarding the user request, the Vald Mirror Gateway bypasses the incoming user request to Vald LB Gateway in its own cluster.

On the other hand, if the incoming user request is an [Object API](https://vald.vdaas.org/docs/api/object/) or [Search API](https://vald.vdaas.org/docs/api/search/), it is bypassed to only a Vald LB Gateway in its own cluster without forwarding it to other Vald Mirror Gateways.

### Automatic rollback on failure

The request may fail at the forwarding destination or the bypass destination.

If some requests fail, the vector data will not be consistent across Vald clusters.

To keep index state consistency, the Vald Mirror Gateway will send the rollback request for the failed request. After the rollback request succeeds, the index state will be the same as before requesting.

The following is the list of rollback types.

- Insert Request
  - Rollback Condition/Trigger
    - Status code other than `ALREADY_EXISTS` exists
  - Rollback request to the successful request
    - REMOVE request
- Remove Request
  - Rollback Condition/Trigger
    - Status code other than `NOT_FOUND` exists
  - Rollback request to the successful request
    - UPSERT Request with old vector
- Update Request
  - Rollback Condition/Trigger
    - Status code other than `ALREADY_EXISTS` exists
  - Rollback request to the successful request
    - REMOVE Request if there is no old vector data
    - UPDATE Request if there is old vector data
- Upsert Request
  - Rollback Condition/Trigger
    - Status code other than `ALREADY_EXISTS` exists
  - Rollback request to the successful request
    - REMOVE Request if there is no old vector data
    - UPDATE Request if there is old vector data
