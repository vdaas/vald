# Configurations

This page introduces best practices to set up values for Vald Helm Chart.

## Vald Helm Chart Overview

Vald Helm Chart's `values.yaml` is composed of the following sections:

- `defaults`
    - default configurations of common parts
    - be overridden by the fields in each section.
- `gateway`
    - configurations of vald-gateway
- `agent`
    - configurations of vald-agent
- `discoverer`
    - configurations of vald-discoverer
- `compressor`
    - configurations of vald-manager-compressor
- `backupManager`
    - configurations of vald-manager-backup
- `indexManager`
    - configurations of vald-manager-index
- `meta`
    - configurations of vald-meta
- `initializer`
    - configurations of MySQL, Cassandra and Redis initializer jobs

In each section, users can configure the deployments and behaviors of each component.

The detailed descriptions of each value can be found in [README of Vald Helm Chart][vald-helm-chart].


## Notable values in Vald Helm Chart

### Basics

#### Specify image tag

It is highly recommended to specify Vald version.
You can specify image version by set `image.tag` field in each component (`[component].image.tag`) or `defaults` section.

```yaml
defaults:
  image:
    tag: v0.0.33
```

or you can use the older image only for agent,

```yaml
agent:
  image:
    tag: v0.0.31
```

#### Specify appropriate logging level and format

The default logging levels and formats are configured in `defaults.logging.level` and `defaults.logging.format`.
You can also specify logging levels and formats in each component section (`[component].logging`).

```yaml
defaults:
  logging:
    level: info
    format: raw
```

you can specify log level `debug` and json format for gateway by the followings:

```yaml
gateway:
  logging:
    level: debug
    format: json
```

#### Servers

Each Vald component has several types of servers.
They can be configured by specifying the values in `defaults.server_config` and can be overwritten by specifying `[component].server_config`.

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

##### gRPC server

gRPC server should be enabled, because all Vald components are using gRPC to communicate with other Vald components.

The API specs are placed in [apis/docs][vald-apis-docs].


##### REST server

REST server is optional.

The swagger specs are placed in [apis/swagger][vald-swagger-specs].


##### Health check servers

There are two types of health check servers are built-in, liveness and readiness.
They are used as servers for [K8s liveness and readiness probe][k8s-liveness-readiness].

By default, liveness servers are disabled for agent and compressor, because the liveness probes may accidentally kill these components.

```yaml
agent:
  server_config:
    healths:
      liveness:
        enabled: false
```

##### Metrics servers

Metrics servers are useful for debugging and monitor Vald components.
There are two types of metrics servers, pprof and Prometheus.

pprof server is a server that implemented using golang's net/http/pprof package.
You can use [google's pprof][google-pprof] to analyze the profiling data exported from it.

Prometheus server is a [Prometheus][prometheus-io] exporter.
It is required to set the `observability` section on each Vald component to enable the monitoring using Prometheus.

#### Observability

The observability features are useful for monitoring Vald components.
They can be enabled by setting the value `true` on the `defaults.observability.enabled` field or override it in each component (`[component].observability.enabled`). And also, enable each feature by setting the value `true` on its `enabled` field.

If observability features are enabled, the metrics will be collected periodically. the duration can be set on `observability.collector.duration`.

Please refer to [Vald operation guide][vald-operation-guide] for more detail.


### Agents

#### NGT

Agent-NGT uses [yahoojapan/NGT][yj-ngt] as a core library for searching vector.
The behaviors of NGT can be configured by setting `agent.ngt` field object.

The important parameters are the followings:

- `agent.ngt.dimension`
- `agent.ngt.distance_type`
- `agent.ngt.object_type`

Users should configure these parameters first for their use case.

For further details, please read [NGT wiki][yj-ngt-wiki].

Agent-NGT has a feature to start indexing automatically.
The behavior of this feature can be configured with these parameters:

- `agent.ngt.auto_index_duration_limit`
- `agent.ngt.auto_index_check_duration`
- `agent.ngt.auto_index_length`


#### Resource requests and limits, Pod priorities

Because agent places indices on memory, termination of agent pods mean loss of indices.
It is important to set resource requests and limits appropriately not to terminate agent pods.

It is highly recommended to request a totally 40% of cluster memory for agent pods.
And also it is highly recommended not to set resource limits to agent pods.

Pod priorities are also useful for saving agent pods from eviction.
By default, very high priority is set to agent pods in the Chart.


#### Pod scheduling

It is recommended to schedule agent pods on different nodes as much as possible.
To achieve this, the following [podAntiAffinity][k8s-affinity-antiaffinity] is set by default.

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

It can be also achieved by using [pod topology spread constraints][k8s-topology-spread-constraints].

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

### Gateway

#### Ingress

Ingress for gateway can be configured by `gateway.ingress` field object.
It is important to set your host to `gateway.ingress.host` field.
`gateway.ingress.servicePort` should be `grpc` or `rest`.

If you're using Vald-Helm-operator, you can check the ingress host by using kubectl command.

```bash
$ kubectl get valdrelease
NAME           INGRESS                  AGE
vald-cluster   gateway.vald.vdaas.org   9d
```

#### Index replica

`gateway.gateway_config.index_replica` means how many number of agent pods that a vector will be inserted.

#### Discoverer request duration

`gateway.gateway_config.discoverer.duration` means a frequency to ask agent pod's IPs to the discoverer.
If discoverer's CPU utilization is too high, try to make this value longer or reduce the number of gateway pods.


#### Meta cache

Gateway has a cache functionality for metadata.
It can be enabled by `gateway.gateway_config.meta.enable_cache` and the behaviors controlled by `gateway.gateway_config.meta.cache_expiration` and `gateway.gateway_config.meta.expired_cache_check_duration`.

#### Resource requests and limits

Gateway's resource requests and limits depend on the request traffic and available resources.
If the request traffic varies largely, it is recommended to enable HPA for gateway and adjust the resource requests.

#### Init containers

Gateway should wait for discoverer, agent, meta, and compressor to be ready because it depends on these components.
For this purpose, "wait-for" type initContainers are provided in the Chart.

```yaml
  initContainers:
    - type: wait-for
      name: wait-for-manager-compressor
      target: compressor
      image: busybox
      sleepDuration: 2
      ...
```

"wait-for" type initContainers check readiness port of the target component is ok or not every "sleepDuration" seconds.
Once it became ready, the initContainer returns zero and become "Completed" status.

The definitions can be found in `_helpers.tpl` in Chart's templates directory.


### Discoverer

#### Resource requests and limits

The number of discoverer pods and resource limits are determined by the configurations of your gateways and index managers because APIs of discoverers are called by gateways and index managers.
Discoverer CPU loads depend on API request traffic = (the number of gateways x gateway's request duration) + (the number of index managers x index manager's request duration).


### Index Manager

#### Init containers

Index managers depend on discoverer and agents.
It is recommended to use initContainers to wait for these components to be ready.


#### Discoverer request duration

Same as gateway, `indexManager.indexer.discoverer.duration` means a frequency to ask agent pod IPs to discoverer.


### Replication Manager

TBW


### Meta, Backup Manager

#### Init containers

Meta and backup manager depends on their backend databases such as Cassandra, MySQL, Redis, etc...
The Chart provides useful initContainers for waiting for these databases.
They can be used as follows:

```yaml
  initContainers:
    - type: wait-for-mysql
      name: wait-for-mysql
      image: mysql:latest
      mysql:
        hosts:
          - mysql.default.svc.cluster.local
        options:
          - "-uroot"
          - "-p${MYSQL_PASSWORD}"
      sleepDuration: 2
      env:
      - name: MYSQL_PASSWORD
        valueFrom:
          secretKeyRef:
            name: mysql-secret
            key: password
```

```yaml
  initContainers:
    - type: wait-for-redis
      name: wait-for-redis
      image: redis:latest
      redis:
        hosts:
          - redis.default.svc.cluster.local
        options:
          - "-a ${REDIS_PASSWORD}"
      sleepDuration: 2
      env:
        - name: REDIS_PASSWORD
          valueFrom:
            secretKeyRef:
              name: redis-secret
              key: password
```

```yaml
  initContainers:
    - type: wait-for-cassandra
      name: wait-for-cassandra
      image: cassandra:latest
      cassandra:
        hosts:
          - cassandra-0.cassandra.default.svc.cluster.local
          - cassandra-1.cassandra.default.svc.cluster.local
          - cassandra-2.cassandra.default.svc.cluster.local
        options:
          - "-uroot"
          - "-p${CASSANDRA_PASSWORD}"
      sleepDuration: 2
      env:
      - name: CASSANDRA_PASSWORD
        valueFrom:
          secretKeyRef:
            name: cassandra-secret
            key: password
```

The definitions can be found in `_helpers.tpl` in Chart's templates directory.


## Advanced


### Ingress/Egress Filters

TBW


### References

For further details, there are references of Helm values in GitHub Vald repository.

- [README of Vald Helm Chart][vald-helm-chart]
- [README of Vald-Helm-Operator Chart][vald-helm-operator-chart]


[vald-helm-chart]: https://github.com/vdaas/vald/tree/v0.0.33/charts/vald
[vald-helm-operator-chart]: https://github.com/vdaas/vald/tree/v0.0.33/charts/vald-helm-operator

[vald-operation-guide]: ./operations.md

[vald-apis-docs]: https://github.com/vdaas/vald/tree/v0.0.33/apis/docs
[vald-swagger-specs]: https://github.com/vdaas/vald/tree/v0.0.33/apis/swagger
[google-pprof]: https://github.com/google/pprof
[prometheus-io]: https://prometheus.io/
[k8s-liveness-readiness]: https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/
[k8s-affinity-antiaffinity]: https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/#affinity-and-anti-affinity
[k8s-topology-spread-constraints]: https://kubernetes.io/docs/concepts/workloads/pods/pod-topology-spread-constraints/
[yj-ngt]: https://github.com/yahoojapan/NGT
[yj-ngt-wiki]: https://github.com/yahoojapan/NGT/wiki
