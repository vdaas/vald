# Operations

This page introduces best practices for operating a Vald cluster.

Table of Contents
---

- [Deployment](#deployment)
    - [Kubernetes cluster](#kubernetes-cluster)
    - [On multi-tenant cluster](#on-multi-tenant-cluster)
- [Monitoring](#monitoring)
    - [Logging](#logging)
    - [Observability features of Vald](#observability-features-of-vald)
    - [Enabling observability feature](#enabling-observability-feature)
    - [Monitoring Vald cluster using Prometheus and Grafana](#monitoring-vald-cluster-using-prometheus-and-grafana)
- [Upgrading](#upgrading)
    - [In case of manual deploy](#in-case-of-manual-deploy)
    - [In case of using Helm](#in-case-of-using-helm)
    - [In case of using Vald-Helm-Operator](#in-case-of-using-vald-helm-operator)
- [References](#references)


## Deployment

### Kubernetes cluster

Since Vald agents hold vector data on their memory, unexpected disruption or eviction of agents may cause loss of indices.
Also, disruption or deletion of worker nodes that have Vald agents may cause loss of indices.
To minimize the loss of indices, it is better to increase number of nodes and pods.

However, to maximize the efficiency of search operations, it is better to have a certain amount of vectors in each NGT vector space.

It is recommended to have more than 3 worker nodes with larger than 16 GB RAM.
It is better to deploy 2 or 3 Vald agent pods to each worker node.

For example:

- 10 worker nodes with 24 GB RAM and 3 Vald agents on each worker node (total: 240 GB RAM, 30 Vald agents)
- 20 worker nodes with 16 GB RAM and 2 Vald agents on each worker node (total: 320 GB RAM, 40 Vald agents)

### On multi-tenant cluster

If you're going to deploy Vald on multi-tenant cluster, please take care about the followings.

- It is recommended to define PriorityClasses for agents not to be evicted.
    - For more info, please visit the page [Pod Priority and Preemption][pod-priority-preemption].
    - If you are using [the Vald chart][vald-helm-chart], PriorityClasses are defined by default.
- It is recommended to define unique namespaces for each Vald and the other apps.
- Then, please define ResourceQuotas for the namespace for the other apps to limit the memory usage of them.
    - For more info, please visit tha page [Resource Quotas][resource-quota].

## Monitoring

### Logging

The logging level of Vald components can be configured by using `defaults.logging.level` (or `[component].logging.level`) field in Helm Chart values.
The level must be a one of "debug", "info", "warn", "error" and "fatal".
The levels are defined in [the Coding Style document][coding-style-logging].

### Observability features of Vald

The observability features are useful for monitoring Vald components.
Vald has various types of exporters, such as Prometheus, Jaeger or Stackdriver.
Using this feature, you can observe and visualize the internal stats or the events like the number of NGT index, when to create index, or the number of RPCs.

### Enabling observability feature

By setting `defaults.observability.enabled` (or `[component].observability.enabled`) in the Helm Chart values to the value `true`, the observability features become enabled.
If observability features are enabled, the metrics will be collected periodically.
The duration can be set on `observability.collector.duration`.

If you'd like to use tracing feature, you should enable it by setting `observability.trace.enabled` to the value `true`. The sampling rate can be configured with `observability.trace.sampling_rate`

### Monitoring Vald cluster using Prometheus and Grafana

In this section, an example of monitoring Vald cluster using [Prometheus][prometheus-io] and [Grafana][grafana] will be shown.

#### Vald configuration

To use Prometheus exporter, you should enable it by setting both `observability.prometheus.enabled` and `server_config.metrics.prometheus.enabled` to the value `true`.
The exporter port and endpoint are specified in each `server_config.metrics.prometheus.port` and `observability.prometheus.endpoint`.

Now it's ready to scrape Vald metrics.
Please deploy Prometheus and Grafana to your cluster.

#### Deploy Prometheus

Prometheus can be installed using one of the followings.

- [Prometheus Operator][prometheus-operator]
- [Prometheus deployments in Vald repository][vald-prometheus]

If you use Prometheus Operator, it is required to set configurations properly along with [Prometheus Configuration][prometheus-configuration] page.
It is recommended to use the endpoints role of the service discovery.

#### Deploy Grafana

Grafana can be installed using one of the followings.

- [Grafana Operator][grafana-operator]
- [Grafana deployments in Vald repository][vald-grafana]

It is required to set your Prometheus to a data source.

#### Example dashboard

Now you can construct your own Grafana dashboard to monitor Vald metrics.
This is an example of a custom dashboard. It is based on [our standard dashboard settings][vald-grafana-dashboards].

[![Grafana dashboard example](../../assets/docs/guides/operations/grafana-example.png)](../../assets/docs/guides/operations/grafana-example.png)


## Upgrading

Our versioning strategy is based on [Semantic Versioning][semver].
Upgrading to new minor or major version may require you to change your configurations.
Please read the [CHANGELOG][CHANGELOG] before upgrading.

### In case of manual deploy

In case of manual deploy, generally, it is required to update your configmaps first.
After that, please update the image tags of Vald components in your deployments.

### In case of using Helm

In case of using Helm and Vald's chart, please update `defaults.image.tag` field and install it.

### In case of using Vald-Helm-Operator

In case of using Vald-Helm-Operator, please update the operator first.
If you're using `valdhelmoperatorrelease` (or `vhor`) resource, please update the `spec.image.tag` field of it.
On the other hand, please update the operator's deployment manually.

After that, please update `image.tag` field in your valdrelease (or `vr`) resource.
The operator will automatically detect the changes and update the deployed Vald cluster.




### References

[vald-helm-chart]: https://github.com/vdaas/vald/tree/master/charts/vald
[vald-helm-operator-chart]: https://github.com/vdaas/vald/tree/master/charts/vald-helm-operator
[CHANGELOG]: https://vald.vdaas.org/docs/release/changelog/

[pod-priority-preemption]: https://kubernetes.io/docs/concepts/configuration/pod-priority-preemption/
[resource-quota]: https://kubernetes.io/docs/concepts/policy/resource-quotas/

[coding-style-logging]: ../contributing/coding-style.md#logging

[prometheus-io]: https://prometheus.io/
[grafana]: https://grafana.com
[vald-prometheus]: https://github.com/vdaas/vald/tree/master/k8s/metrics/prometheus
[prometheus-operator]: https://github.com/coreos/prometheus-operator
[prometheus-configuration]: https://prometheus.io/docs/prometheus/latest/configuration/configuration/
[vald-grafana]: https://github.com/vdaas/vald/tree/master/k8s/metrics/grafana
[vald-grafana-dashboards]: https://github.com/vdaas/vald/tree/master/k8s/metrics/grafana/dashboards
[grafana-operator]: https://operatorhub.io/operator/grafana-operator
[semver]: https://semver.org/
