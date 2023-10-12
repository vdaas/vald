# Vald Mirror Gateway

Vald Mirror Gateway is an optional component of the Vald, allowing the vector data to be synchronized across multiple Vald clusters.

This component makes it possible to enhance availability during a cluster failure.

<img src="../../../assets/docs/overview/component/mirror-gateway/mirror-gateway.png">

## Responsibility

Vald Mirror Gateway is responsible for the followings:

- Forward user requests ([Insert](https://vald.vdaas.org/docs/api/insert/) / [Upsert](https://vald.vdaas.org/docs/api/upsert/) / [Update](https://vald.vdaas.org/docs/api/update/) / [Remove](https://vald.vdaas.org/docs/api/remove/)) to the other Vald Mirror Gateways in the same group.
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

For more information about `ValdMirrorTarget` configuration, please refer to [Custom Resource Configuration](https://vald.vdaas.org/docs/user-guides/mirroring-configuration/).

### Request forwarding

<img src="../../../assets/docs/overview/component/mirror-gateway/request-forwarding.png">

The Vald Mirror Gateway forwards the incoming user request ([Insert](https://vald.vdaas.org/docs/api/insert/) / [Upsert](https://vald.vdaas.org/docs/api/upsert/) / [Update](https://vald.vdaas.org/docs/api/update/) / [Remove](https://vald.vdaas.org/docs/api/remove/)) to other Vald Mirror Gateways.
Then, while forwarding the user request, the Vald Mirror Gateway bypasses the incoming user request to Vald LB Gateway in its own cluster.

On the other hand, if the incoming user request is an [Object API](https://vald.vdaas.org/docs/api/object/) or [Search API](https://vald.vdaas.org/docs/api/search/), it is bypassed to only a Vald LB Gateway in its own cluster without forwarding it to other Vald Mirror Gateways.

### Continuous processing on failure

The request may fail at the forwarding destination or the bypass destination.

If some of the requests fails, the processing continues based on their status code.

The following is an overview of the process for each request.

- Insert Request

  - If the target host returns a status of `ALREADY_EXISTS`, the Update request is sent to this host.
  - If all target hosts returns a status of `ALREADY_EXISTS`, the Mirror Gateway returns `ALREADY_EXISTS` to the user without continuous processing.
  - If all target hosts returns a status of `OK` or `ALREADY_EXISTS`, the Mirror Gateway returns `OK` to the user without continuous processing.

- Update Request

  - If the target host returns a status of `NOT_FOUND`, the Insert request is sent to this host.
  - If all target hosts returns a status of `ALREADY_EXISTS`, the Mirror Gateway returns `ALREADY_EXISTS` to the user without continuous processing.
  - If all target hosts returns a status of `OK` or `ALREADY_EXISTS`, the Mirror Gateway returns `OK` to the user without continuous processing.

- Upsert Request

  - If all target hosts returns a status of `OK` or `ALREADY_EXISTS`, the Mirror Gateway returns `OK` to the user without continuous processing.
  - If all target hosts returns a status of `ALREADY_EXISTS`, the Mirror Gateway returns `ALREADY_EXISTS` to the user without continuous processing.

- Remove Request

  - If all target hosts returns a status of `OK` or `NOT_FOUND`, the Mirror Gateway returns `OK` to the user without continuous processing.
  - If all target hosts returns a status of `NOT_FOUND`, the Mirror Gateway returns `NOT_FOUND` to the user without continuous processing.

- RemoveByTimestamp Request

  - Same as `Remove Request`.

For more information, please refer to [Mirror Gateway Troubleshooting](https://vald.vdaas.org/docs/troubleshooting/mirror-gateway/).
