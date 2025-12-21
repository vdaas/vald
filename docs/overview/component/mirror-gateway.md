# Vald Mirror Gateway

Vald Mirror Gateway is an optional component of the Vald, allowing the vector data to be synchronized across multiple Vald clusters.

This component makes it possible to enhance availability during a cluster failure.

<img src="../../../assets/docs/overview/component/mirror-gateway/mirror-gateway.png">

## Responsibility

Vald Mirror Gateway is responsible for the followings:

- Forward user requests ([Insert](../../api/insert.md) / [Upsert](../../api/upsert.md) / [Update](../../api/update.md) / [Remove](../../api/remove.md)) to the other Vald Mirror Gateways in the same group.
- Manages the state of indexes stored in all clusters to ensure they are consistent.

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

For more information about `ValdMirrorTarget` configuration, please refer to [Custom Resource Configuration](../../user-guides/mirroring-configuration.md).

### Request forwarding

<img src="../../../assets/docs/overview/component/mirror-gateway/request-forwarding.png">

The Vald Mirror Gateway forwards the incoming user request ([Insert](../../api/insert.md) / [Upsert](../../api/upsert.md) / [Update](../../api/update.md) / [Remove](../../api/remove.md)) to other Vald Mirror Gateways.
Then, while forwarding the user request, the Vald Mirror Gateway bypasses the incoming user request to Vald LB Gateway in its own cluster.

On the other hand, if the incoming user request is an [Object API](../../api/object.md) or [Search API](../../api/search.md), it is bypassed to only a Vald LB Gateway in its own cluster without forwarding it to other Vald Mirror Gateways.

### Continuous processing on failure

The request may fail at the forwarding destination or the bypass destination.

If some of the requests fails, the processing continues based on their status code.

Here's an overview of how the Mirror Gateway handles failures for each type of request.

For more information about status code, please refer to [Mirror Gateway Troubleshooting](../../troubleshooting/mirror-gateway.md).

- Insert Request

  - If the target host returns a status code of `ALREADY_EXISTS`, the Update request is sent to this host.
  - If the target host returns a status code other than `OK`, `ALREADY_EXISTS`, the Mirror Gateway returns that status code without continuous processing.
  - If all target hosts return a status code `ALREADY_EXISTS`, the Mirror Gateway returns `ALREADY_EXISTS`.
  - If all target hosts return a status code `OK` or `ALREADY_EXISTS`, the Mirror Gateway returns `OK`.

- Update Request

  - If the target host returns a status code `NOT_FOUND`, the Insert request is sent to this host.
  - If the target host returns a status code other than `OK`, `ALREADY_EXISTS`, the Mirror Gateway returns that status code without continuous processing.
  - If all target hosts return a status code `ALREADY_EXISTS`, the Mirror Gateway returns `ALREADY_EXISTS`.
  - If all target hosts return a status code `OK` or `ALREADY_EXISTS`, the Mirror Gateway returns `OK`.

- Upsert Request

  - If all target hosts return a status code `ALREADY_EXISTS`, the Mirror Gateway returns `ALREADY_EXISTS`.
  - If the target host returns a status code other than `OK` or `ALREADY_EXISTS`, the Mirror Gateway returns that status code without continuous processing.
  - If all target hosts return a status code `OK` or `ALREADY_EXISTS`, the Mirror Gateway returns `OK`.

- Remove/RemoveByTimestamp Request

  - If all target hosts return a status code `NOT_FOUND`, the Mirror Gateway returns `NOT_FOUND`.
  - If the target host returns a status code other than `OK` or `NOT_FOUND`, the Mirror Gateway returns that status code without continuous processing.
  - If all target hosts return a status code `OK` or `NOT_FOUND`, the Mirror Gateway returns `OK`.
