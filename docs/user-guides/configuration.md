# Configurations

This page introduces best practices for setting up values for the Vald Helm Chart.

Before reading, please read the overview of Vald Helm Chart in [its README][vald-helm-chart].

<div class="notice">
This page shows the notable fields in Vald Helm Chart.<BR>
It is highly recommended to verify before deployment.
</div>

## General

### Specify image tag

It is highly recommended to specify the Vald version.
You can specify the image version by setting `image.tag` field in each component (`[component].image.tag`) or `defaults` section.

```yaml
defaults:
  image:
    tag: v1.5.6
```

or you can use the older image only for a target component, e.g., the agent,

```yaml
agent:
  image:
    tag: v1.5.5
```

### Specify appropriate logging level and format

The default logging levels and formats are configured in `defaults.logging.level` and `defaults.logging.format`.
You can also specify them in each component section (`[component].logging`).

```yaml
defaults:
  logging:
    level: info
    format: raw
```

You can specify log level `debug` and JSON format for lb-gateway by the followings:

```yaml
gateway:
  lb:
    logging:
      level: debug
      format: json
```

The logging level is defined in [the Coding Style Guide](../contributing/coding-style.md#logging).

### Servers

Each Vald component has several types of servers.
They can be configured by specifying the values in `defaults.server_config`.

Examples:

```yaml
defaults:
  server_config:
    servers:
      grpc:
        enabled: true
        host: 0.0.0.0
        port: 8081
        servicePort: 8081
        server:
          mode: GRPC
          ...
```

In addition, they can be overwritten by setting each `[component].server_config`, e.g., `gateway.lb.server_config` is following.

```yaml
gateway:
  lb:
    server_config:
      servers:
        rest:
          enabled: true
          host: 0.0.0.0
          port: 8080
          servicePort: 8080
          server:
            mode: REST
            ...
```

#### gRPC server

gRPC server should be enabled because all Vald components use gRPC to communicate with others.
The API specs are placed in [Vald APIs](../api).

#### REST server

REST server is optional.
The swagger specs are placed in [Vald APIs Swagger][vald-swagger-specs].

#### Health check servers

There are two built-in health check servers: liveness and readiness.
They are used as servers for [Kubernetes liveness and readiness probe][kubernetes-liveness-readiness].
The liveness health server is disabled by default due to the liveness probe may accidentally kill the Vald Agent component.

```yaml
agent:
  server_config:
    healths:
      liveness:
        enabled: false
```

### Metrics servers

The metrics server enables easier debugging and monitoring of Vald components.
There are two types of metrics servers: pprof and Prometheus.

pprof server is implemented using Go's `net/http/pprof` package.
You can use [google's pprof][google-pprof] to analyze the exported profile result.

Prometheus server is a [Prometheus][prometheus-io] exporter.
It is required to set the `observability` section on each Vald component to enable the monitoring using Prometheus.
Please refer to the next section.

### Observability

The observability features are useful for monitoring Vald components.
These settings can be enabled by setting the `defaults.observability.enabled` field to the value `true` or by overriding it in each component (`[component].observability.enabled`).
And also, enable each feature by setting the value `true` on its `enabled` field.

If observability features are enabled, the metrics will be collected periodically.
The duration can be set on `observability.collector.duration`.
Please refer to [the Vald operation guide](../user-guides/configuration.md) for more detail.

## Component basic configuration

### Agents

#### NGT

Vald Agent NGT uses [yahoojapan/NGT][yj-ngt] as a core library for searching vectors.
The behaviors of NGT can be configured by setting `agent.ngt` field object.

The important parameters are the followings:

- `agent.ngt.dimension`
- `agent.ngt.distance_type`
- `agent.ngt.object_type`

Users should configure these parameters first to fit their use case.
For further details, please read [the NGT wiki][yj-ngt-wiki].

Vald Agent NGT has a feature to start indexing automatically.
The behavior of this feature can be configured with these parameters:

- `agent.ngt.auto_index_duration_limit`
- `agent.ngt.auto_index_check_duration`
- `agent.ngt.auto_index_length`

<div class="notice">
While the Vald Agent NGT is in the process of creating indexes, it will ignore all search requests to the target pods.
</div>

<div class="warning">
When deploying Vald Index Manager, the above parameters should be set much longer than the Vald Index Manager settings (Please refer to the Vald Index Manager section).<BR>
E.g., set agent.ngt.auto_index_duration_limit to "720h" and agent.ngt.auto_index_check_duration to "24h".<BR>
This is because the Vald Index Manager accurately grasps the index information of each Vald Agent NGT and controls the execution timing of indexing.<BR><BR>
When the setting parameter of Vald Agent NGT is shorter than the setting value of Vald Index Manager, Vald Agent NGT may start indexing by itself without the execution command from Vald Index Manager.
If this happens, the Index Manager may not function properly.
</div>

#### Resource requests and limits, Pod priorities

Because the Vald Agent pod places indexes on memory, termination of agent pods causes loss of indexes.
It is important to set the resource requests and limits appropriately to avoid terminating the Vald Agent pods.

Requesting 40% of cluster memory for agent pods is highly recommended.
Also, it is highly recommended not to set the resource limits for the Vald Agent pods.

Pod priorities are also useful for saving agent pods from eviction.
By default, very high priority is set to agent pods in the Chart.

[The capacity planning page](../user-guides/capacity-planning.md) helps to estimate the resources.

#### Pod scheduling

It is recommended to schedule agent pods on different nodes as much as possible.

<div class="warning">
The affinity setting for Vald Agent is significant for the Vald cluster.<BR>
Please DO NOT remove the default settings.
</div>

To achieve this, the following [podAntiAffinity][kubernetes-affinity-antiaffinity] is set by default.

```yaml
agent:
  affinity:
    podAntiAffinity:
      preferredDuringSchedulingIgnoredDuringExecution:
        - weight: 100
          podAffinityTerm:
            topologyKey: kubernetes.io/hostname
            labelSelector:
              matchExpressions:
                - key: app
                  operator: In
                  values:
                    - vald-agent-ngt
```

It can also be achieved by using [pod topology spread constraints][kubernetes-topology-spread-constraints].

```yaml
agent:
  topologySpreadConstraints:
    - topologyKey: node
      maxSkew: 1
      whenUnsatisfiable: ScheduleAnyway
      labelSelector:
        matchLabels:
          app: vald-agent-ngt
  affinity:
    podAntiAffinity:
      preferredDuringSchedulingIgnoredDuringExecution: [] # to disable default settings
```

### Gateway LB

#### Ingress

Ingress for gateways can be configured by `gateway.{filter,lb}.ingress` field object.
It is important to set your host to `gateway.{filter,lb}.ingress.host` field.
`gateway.{filter,lb}.ingress.servicePort` should be `grpc` or `rest`.

```yaml
gateway:
  lb:
    ingress:
      enabled: true
      host: vald.vdaas.org # Set correct hostname here
      servicePort: grpc
```

#### Index replica

`gateway.lb.gateway_config.index_replica` represents how many Vald Agent pods a vector will be inserted into.
The maximum value of the index replica should be 30% of the Vald Agent pods deployed.

```yaml
gateway:
  lb:
    gateway_config:
      index_replica: 3 // By setting the index replica to 3, the number of Vald Agent pods deployed should be more than 9 (3 / 0.3).
```

#### Resource requests and limits

The gateway's resource requests and limits depend on the request traffic and available resources.
If the request traffic varies largely, enabling HPA for the gateway and adjusting the resource requests is recommended.

#### Discoverer request duration

`gateway.lb.gateway_config.discoverer.duration` represents the frequency of sending requests to the discoverer.
If the discoverer's CPU utilization is too high, make this value longer or reduce the number of LB gateway pods.

```yaml
gateway:
  lb:
    gateway_config:
      discoverer:
        duration: 2s
```

### Discoverer

#### Cluster Role

Vald Discoverer gets the Node and Pod metrics from [kube-apiserver](https://kubernetes.io/docs/reference/command-line-tools-reference/kube-apiserver/) as described in [Vald Discoverer](../overview/component/discoverer.md).
Vald's Helm deployment supports RBAC as default, and the default configuration is the following.

```yaml
discoverer:
  clusterRole:
    # if true, the clusterRole configuration will be created.
    enabled: true
    name: discoverer
  clusterRoleBinding:
    # if true, the clusterRoleBinding configuration will be created.
    enabled: true
    name: discoverer
  serviceAccount:
    # if true, the serviceAccount configuration will be created.
    enabled: true
    name: vald
```

When `RBAC` is unavailable in your environment, or you would like to put some restrictions, please modify it and grant the permissions to the user executing the discoverer.
Each configuration file is the following:

- [Cluster role](https://github.com/vdaas/vald/blob/main/k8s/discoverer/clusterrole.yaml)
- [Cluster role binding](https://github.com/vdaas/vald/blob/main/k8s/discoverer/clusterrolebinding.yaml)
- [Service account](https://github.com/vdaas/vald/blob/main/k8s/discoverer/serviceaccount.yaml)

#### Resource requests and limits

The number of discoverer pods and resource limits can be estimated by the configurations of your LB gateways and index managers because its APIs are called by them.
Discoverer CPU loads almost depend on API request traffic.

```bash
# The API traffic formula
(the number of LB gateways x its request frequency) + (the number of index managers x its request frequency).
```

### Index Manager

#### Execution index command to Vald Agent

Vald Index Manager controls the indexing timing for all Vald Agent pods in the Vald cluster.
These parameters are related to the control process.

```yaml
manager:
  index:
    indexer:
      # namespace of agent pods to manage
      agent_namespace: vald # namespace of agent pods to manage
      # check duration of automatic indexing
      auto_index_check_duration: "1m"
      # limit duration of automatic indexing
      auto_index_duration_limit: "30m"
      # number of caches to trigger automatic indexing
      auto_index_length: 100
      # limit duration of automatic index saving
      auto_save_index_duration_limit: "3h"
      # duration of automatic index saving wait duration for next saving
      auto_save_index_wait_duration: "10m"
      # the number of Agent Pods indexing at the same time
      concurrency: 1
      # number of pool size of creating index processing
      creation_pool_size: 10000
```

#### Discoverer request duration

Same as LB gateway, `manager.index.indexer.discoverer.duration` represents the frequency of sending requests to the discoverer.

## References

For further details, there are references to the Helm values in the Vald GitHub repository.

- [README of Vald Helm Chart][vald-helm-chart]
- [README of Vald-Helm-Operator Chart][vald-helm-operator-chart]

<!-- TODO: add related document(pullugable options) -->

[vald-helm-chart]: https://github.com/vdaas/vald/tree/main/charts/vald
[vald-helm-operator-chart]: https://github.com/vdaas/vald/tree/main/charts/vald-helm-operator
[vald-apis-docs]: https://github.com/vdaas/vald/tree/main/apis/docs
[vald-swagger-specs]: https://github.com/vdaas/vald/tree/main/apis/swagger
[google-pprof]: https://github.com/google/pprof
[prometheus-io]: https://prometheus.io/
[kubernetes-liveness-readiness]: https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/
[kubernetes-affinity-antiaffinity]: https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/#affinity-and-anti-affinity
[kubernetes-topology-spread-constraints]: https://kubernetes.io/docs/concepts/workloads/pods/pod-topology-spread-constraints/
[yj-ngt]: https://github.com/yahoojapan/NGT
[yj-ngt-wiki]: https://github.com/yahoojapan/NGT/wiki
