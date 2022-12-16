# Vald Discoverer

Vald Discoverer runs individually and helps automate operations for the Vald cluster.

It is also the only component that communicates with [kube-apiserver](https://kubernetes.io/docs/reference/command-line-tools-reference/kube-apiserver/).

## Responsibility

Vald Discoverer is responsible for retrieving each Node and Pod resource usage from [kube-apiserver](https://kubernetes.io/docs/reference/command-line-tools-reference/kube-apiserver/) and sharing it with other components in the Vald cluster.

## Feature

### Getting Node and Pod metrics

Vald Discoverer requires [kube-apiserver](https://kubernetes.io/docs/reference/command-line-tools-reference/kube-apiserver/) to get Node and Pod metrics in the Vald cluster.

It synchronizes the metrics when it changes by the `Reconcile function`, one of the internal functions.

The kind of metrics are here:

- Node

  | Key name     | Description                     |
  | ------------ | ------------------------------- |
  | Name         | node name                       |
  | InternalAddr | node internal address           |
  | ExternalAddr | node external address           |
  | CPU          | CPU limit and CPU request       |
  | Memory       | Memory limit and Memory request |
  | Pods         | Pod information in the node     |

- Pod

  | Key name  | Description                     |
  | --------- | ------------------------------- |
  | AppName   | AppName of pod                  |
  | Name      | pod name                        |
  | Namespace | namespace where pod runs        |
  | IP        | IP address                      |
  | CPU       | CPU limit and CPU request       |
  | Memory    | Memory limit and Memory request |

When syncing success, Vald Discoverer chooses the necessary metrics from the result and stores them into their four kinds of `Map` on their local memory:

| Map name        | Description                          |
| --------------- | ------------------------------------ |
| podsByNode      | pod metrics list group by node       |
| podsByNamespace | pod metrics list group by namespace  |
| podsByName      | pod metrics list group by pod name   |
| nodebyName      | node metrics list group by node name |

<!-- TODO:image -->

### Sharing Node and Pod metrics

Vald Discoverer provides the internal client for another component.

The component, which requires the metrics information, uses when it needs.

Vald LB Gateway and Vald Manager Index in the Vald cluster require it to achieve its responsibility.

For example, Vald LB Gateway creates the new `discoverer client` with config parameters and binds it to the new `gateway lb client` when Vald LB Gateway starts the container.

When the `initContainer` successes, Vald LB Gateway can get the metrics asynchronously according to its set parameters.

<!-- TODO:image -->

### Cluster role configurations

Please refer [here](../../user-guides/cluster-role-binding.md) for more information about the cluster role configuration.
