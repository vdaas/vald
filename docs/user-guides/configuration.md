# Configurations

This page introduces best practices for setting up values for Vald Helm Chart.
Before reading, please read the overview of Vald Helm Chart in [its README][vald-helm-chart].

<div class="notice">
This page shows the notable fields in Vald Helm Chart.<BR>
It is highly recommended before deployment.
</div>

## General

### Specify image tag

It is highly recommended to specify Vald version.
You can specify image version by set `image.tag` field in each component (`[component].image.tag`) or `defaults` section.

```yaml
defaults:
  image:
    tag: v1.5.6
```

or you can use the older image only for agent,

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

you can specify log level `debug` and JSON format for lb-gateway by the followings:

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
They can be overwritten by specifying `[component].server_config`.

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

gRPC server should be enabled, because all Vald components use gRPC to communicate with others.
The API specs are placed in [apis/docs][vald-apis-docs].

#### REST server

REST server is optional.
The swagger specs are placed in [apis/swagger][vald-swagger-specs].

#### Health check servers

There are two types of built-in health check servers, liveness and readiness.
They are used as servers for [Kubernetes liveness and readiness probe][kubernetes-liveness-readiness].
By default, liveness servers are disabled for agent, because the liveness probes may accidentally kill it.

```yaml
agent:
  server_config:
    healths:
      liveness:
        enabled: false
```

### Metrics servers

Metrics servers are useful for debugging and monitoring Vald components.
There are two types of metrics servers, pprof and Prometheus.

pprof server is implemented using Go's `net/http/pprof` package.
You can use [google's pprof][google-pprof] to analyze the profiling data exported from it.

Prometheus server is a [Prometheus][prometheus-io] exporter.
It is required to set the `observability` section on each Vald component to enable the monitoring using Prometheus.
Please refer to the next section.

### Observability

The observability features are useful for monitoring Vald components.
They can be enabled by setting the value `true` on the `defaults.observability.enabled` field or override it in each component (`[component].observability.enabled`).
And also, enable each feature by setting the value `true` on its `enabled` field.

If observability features are enabled, the metrics will be collected periodically.
The duration can be set on `observability.collector.duration`.
Please refer to [Vald operation guide](../user-guides/configuration.md) for more detail.

## Component basic configuration

### Agents

#### NGT

Agent-NGT uses [yahoojapan/NGT][yj-ngt] as a core library for searching vector.
The behaviors of NGT can be configured by setting `agent.ngt` field object.

The important parameters are the followings:

- `agent.ngt.dimension`
- `agent.ngt.distance_type`
- `agent.ngt.object_type`

Users should configure these parameters first to fit to their use case.

For further details, please read [NGT wiki][yj-ngt-wiki].

Agent-NGT has a feature to start indexing automatically.
The behavior of this feature can be configured with these parameters:

- `agent.ngt.auto_index_duration_limit`
- `agent.ngt.auto_index_check_duration`
- `agent.ngt.auto_index_length`

#### Resource requests and limits, Pod priorities

Because agent places indices on memory, termination of agent pods causes loss of indices.
It is important to set resource requests and limits appropriately not to terminate agent pods.

It is highly recommended to request a totally 40% of cluster memory for agent pods.
And also it is highly recommended not to set resource limits to agent pods.

Pod priorities are also useful for saving agent pods from eviction.
By default, very high priority is set to agent pods in the Chart.

[The capacity planning](../user-guides/capacity-planning.md) helps to estimate the resources.

#### Pod scheduling

It is recommended to schedule agent pods on different nodes as much as possible.

<div class="warning">
The affinity setting for Vald Agent is the significant for the Vald cluster.<BR>
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

It can be also achieved by using [pod topology spread constraints][kubernetes-topology-spread-constraints].

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

`gateway.lb.gateway_config.index_replica` represents how many Vald Agent pods that a vector will be inserted into.
We recommend set it as one third of the number of Vald Agent pods.

```yaml
gateway:
  lb:
    gateway_config:
      index_replica: 3 // The number of Vald Agent pods should be larger than 9.
```

#### Resource requests and limits

The gateway's resource requests and limits depend on the request traffic and available resources.
If the request traffic varies largely, it is recommended to enable HPA for gateway and adjust the resource requests.

#### Discoverer request duration

`gateway.lb.gateway_config.discoverer.duration` represents a frequency to send requests to discoverer.
If discoverer's CPU utilization is too high, make this value longer or reduce the number of LB gateway pods.

```yaml
gateway:
  lb:
    gateway_config:
      discoverer:
        duration: 2s
```

### Discoverer

#### Cluster Role

Vald Discoverer gets the Node and Por metrics from [kube-apiserver](https://kubernetes.io/ja/docs/reference/command-line-tools-reference/kube-apiserver/) as described in [Vald Discoverer](../overview/component/discoverer.md).
Vald's Helm deployment supports RBAC as default, and the default configuration is the following.

```yaml
discoverer:
  clusterRole:
    enabled: true
    name: discoverer
  clusterRoleBinding:
    enabled: true
    name: discoverer
  serviceAccount:
    enabled: true
    name: vald
```

When `RBAC` is unavailable in your environment, or you would like to put some restrictions, please modify it and grant the permissions to the user executing the discoverer.
Each configuration file is the following:

- [clusterRole](https://github.com/vdaas/vald/blob/master/k8s/discoverer/clusterrole.yaml)
- [clusterRoleBinding](https://github.com/vdaas/vald/blob/master/k8s/discoverer/clusterrolebinding.yaml)
- [serviceAccount](https://github.com/vdaas/vald/blob/master/k8s/discoverer/serviceaccount.yaml)


#### Resource requests and limits

The number of discoverer pods and resource limits can be estimated by the configurations of your LB gateways and index managers because its APIs are called by them.
Discoverer CPU loads almost depend on API request traffic = (the number of LB gateways x its request frequency) + (the number of index managers x its request frequency).

### Index Manager

#### Discoverer request duration

Same as LB gateway, `manager.index.indexer.discoverer.duration` represents a frequency to send requests to discoverer.

## References

For further details, there are references of Helm values in GitHub Vald repository.

- [README of Vald Helm Chart][vald-helm-chart]
- [README of Vald-Helm-Operator Chart][vald-helm-operator-chart]

<!-- TODO: add related document(pullugable options) -->


[vald-helm-chart]: https://github.com/vdaas/vald/tree/master/charts/vald
[vald-helm-operator-chart]: https://github.com/vdaas/vald/tree/master/charts/vald-helm-operator

[vald-apis-docs]: https://github.com/vdaas/vald/tree/master/apis/docs
[vald-swagger-specs]: https://github.com/vdaas/vald/tree/master/apis/swagger
[google-pprof]: https://github.com/google/pprof
[prometheus-io]: https://prometheus.io/
[kubernetes-liveness-readiness]: https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/
[kubernetes-affinity-antiaffinity]: https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/#affinity-and-anti-affinity
[kubernetes-topology-spread-constraints]: https://kubernetes.io/docs/concepts/workloads/pods/pod-topology-spread-constraints/
[yj-ngt]: https://github.com/yahoojapan/NGT
[yj-ngt-wiki]: https://github.com/yahoojapan/NGT/wiki
