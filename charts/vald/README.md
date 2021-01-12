Vald
===

This is a Helm chart to install Vald components.

Current chart version is `v0.0.66`

Table of Contents
---

- [Install](#install)
- [Configuration](#configuration)
    - [Overview](#overview)
    - [Parameters](#parameters)
- [Miscellaneous](#miscellaneous)
    - [Standalone Vald agent NGT deployment](#standalone-vald-agent-ngt-deployment)

Install
---

Add Vald Helm repository

    $ helm repo add vald https://vald.vdaas.org/charts

Run the following command to install the chart,

    $ helm install vald-cluster vald/vald

Configuration
---

### Overview

`values.yaml` is composed of the following sections:

- `defaults`
    - default configurations of common parts
    - be overridden by the fields in each components' configurations
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

### Parameters

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| agent.affinity.nodeAffinity.preferredDuringSchedulingIgnoredDuringExecution | list | `[]` | node affinity preferred scheduling terms |
| agent.affinity.nodeAffinity.requiredDuringSchedulingIgnoredDuringExecution.nodeSelectorTerms | list | `[]` | node affinity required node selectors |
| agent.affinity.podAffinity.preferredDuringSchedulingIgnoredDuringExecution | list | `[]` | pod affinity preferred scheduling terms |
| agent.affinity.podAffinity.requiredDuringSchedulingIgnoredDuringExecution | list | `[]` | pod affinity required scheduling terms |
| agent.affinity.podAntiAffinity.preferredDuringSchedulingIgnoredDuringExecution | list | `[{"podAffinityTerm":{"labelSelector":{"matchExpressions":[{"key":"app","operator":"In","values":["vald-agent-ngt"]}]},"topologyKey":"kubernetes.io/hostname"},"weight":100}]` | pod anti-affinity preferred scheduling terms |
| agent.affinity.podAntiAffinity.requiredDuringSchedulingIgnoredDuringExecution | list | `[]` | pod anti-affinity required scheduling terms |
| agent.annotations | object | `{}` | deployment annotations |
| agent.enabled | bool | `true` | agent enabled |
| agent.env | list | `[]` | environment variables |
| agent.externalTrafficPolicy | string | `""` | external traffic policy (can be specified when service type is LoadBalancer or NodePort) : Cluster or Local |
| agent.hpa.enabled | bool | `false` | HPA enabled |
| agent.hpa.targetCPUUtilizationPercentage | int | `80` | HPA CPU utilization percentage |
| agent.image.pullPolicy | string | `"Always"` | image pull policy |
| agent.image.repository | string | `"vdaas/vald-agent-ngt"` | image repository |
| agent.image.tag | string | `""` | image tag (overrides defaults.image.tag) |
| agent.initContainers | list | `[]` | init containers |
| agent.kind | string | `"StatefulSet"` | deployment kind: Deployment, DaemonSet or StatefulSet |
| agent.logging | object | `{}` | logging config (overrides defaults.logging) |
| agent.maxReplicas | int | `300` | maximum number of replicas. if HPA is disabled, this value will be ignored. |
| agent.maxUnavailable | string | `"1"` | maximum number of unavailable replicas |
| agent.minReplicas | int | `20` | minimum number of replicas. if HPA is disabled, the replicas will be set to this value |
| agent.name | string | `"vald-agent-ngt"` | name of agent deployment |
| agent.ngt.auto_index_check_duration | string | `"30m"` | check duration of automatic indexing |
| agent.ngt.auto_index_duration_limit | string | `"24h"` | limit duration of automatic indexing |
| agent.ngt.auto_index_length | int | `100` | number of cache to trigger automatic indexing |
| agent.ngt.auto_save_index_duration | string | `"35m"` | duration of automatic save index |
| agent.ngt.bulk_insert_chunk_size | int | `10` | bulk insert chunk size |
| agent.ngt.creation_edge_size | int | `20` | creation edge size |
| agent.ngt.default_epsilon | float | `0.01` | default epsilon used for search |
| agent.ngt.default_pool_size | int | `10000` | default create index batch pool size |
| agent.ngt.default_radius | float | `-1` | default radius used for search |
| agent.ngt.dimension | int | `4096` | vector dimension |
| agent.ngt.distance_type | string | `"l2"` | distance type. it should be `l1`, `l2`, `angle`, `hamming`, `cosine`, `normalizedangle`, `normalizedcosine` or `jaccard`. for further details about NGT libraries supported distance is https://github.com/yahoojapan/NGT/wiki/Command-Quick-Reference and vald agent's supported NGT distance type is https://pkg.go.dev/github.com/vdaas/vald/internal/core/algorithm/ngt#pkg-constants |
| agent.ngt.enable_in_memory_mode | bool | `true` | in-memory mode enabled |
| agent.ngt.enable_proactive_gc | bool | `true` | enable proactive GC call for reducing heap memory allocation |
| agent.ngt.index_path | string | `""` | path to index data |
| agent.ngt.initial_delay_max_duration | string | `"3m"` | maximum duration for initial delay |
| agent.ngt.load_index_timeout_factor | string | `"1ms"` | a factor of load index timeout. timeout duration will be calculated by (index count to be loaded) * (factor). |
| agent.ngt.max_load_index_timeout | string | `"10m"` | maximum duration of load index timeout |
| agent.ngt.min_load_index_timeout | string | `"3m"` | minimum duration of load index timeout |
| agent.ngt.object_type | string | `"float"` | object type. it should be `float` or `uint8`. for further details: https://github.com/yahoojapan/NGT/wiki/Command-Quick-Reference |
| agent.ngt.search_edge_size | int | `10` | search edge size |
| agent.nodeName | string | `""` | node name |
| agent.nodeSelector | object | `{}` | node selector |
| agent.observability | object | `{"jaeger":{"service_name":"vald-agent-ngt"},"stackdriver":{"profiler":{"service":"vald-agent-ngt"}}}` | observability config (overrides defaults.observability) |
| agent.persistentVolume.accessMode | string | `"ReadWriteOnce"` | agent pod storage accessMode |
| agent.persistentVolume.enabled | bool | `false` | enables PVC. It is required to enable if agent pod's file store functionality is enabled with non in-memory mode |
| agent.persistentVolume.size | string | `"100Gi"` | size of agent pod volume |
| agent.persistentVolume.storageClass | string | `"vald-sc"` | storageClass name for agent pod volume |
| agent.podAnnotations | object | `{}` | pod annotations |
| agent.podManagementPolicy | string | `"OrderedReady"` | pod management policy: OrderedReady or Parallel |
| agent.podPriority.enabled | bool | `true` | agent pod PriorityClass enabled |
| agent.podPriority.value | int | `1000000000` | agent pod PriorityClass value |
| agent.progressDeadlineSeconds | int | `600` | progress deadline seconds |
| agent.resources | object | `{"requests":{"cpu":"300m","memory":"4Gi"}}` | compute resources. recommended setting of memory requests = cluster memory * 0.4 / number of agent pods |
| agent.revisionHistoryLimit | int | `2` | number of old history to retain to allow rollback |
| agent.rollingUpdate.maxSurge | string | `"25%"` | max surge of rolling update |
| agent.rollingUpdate.maxUnavailable | string | `"25%"` | max unavailable of rolling update |
| agent.rollingUpdate.partition | int | `0` | StatefulSet partition |
| agent.server_config | object | `{"healths":{"liveness":{"enabled":false},"readiness":{}},"metrics":{"pprof":{},"prometheus":{}},"servers":{"grpc":{},"rest":{}}}` | server config (overrides defaults.server_config) |
| agent.service.annotations | object | `{}` | service annotations |
| agent.service.labels | object | `{}` | service labels |
| agent.serviceType | string | `"ClusterIP"` | service type: ClusterIP, LoadBalancer or NodePort |
| agent.sidecar.config.auto_backup_duration | string | `"24h"` | auto backup duration |
| agent.sidecar.config.auto_backup_enabled | bool | `true` | auto backup triggered by timer is enabled |
| agent.sidecar.config.blob_storage.bucket | string | `""` | bucket name |
| agent.sidecar.config.blob_storage.s3.access_key | string | `"_AWS_ACCESS_KEY_"` | s3 access key |
| agent.sidecar.config.blob_storage.s3.enable_100_continue | bool | `true` | enable AWS SDK adding the 'Expect: 100-Continue' header to PUT requests over 2MB of content. |
| agent.sidecar.config.blob_storage.s3.enable_content_md5_validation | bool | `true` | enable the S3 client to add MD5 checksum to upload API calls. |
| agent.sidecar.config.blob_storage.s3.enable_endpoint_discovery | bool | `false` | enable endpoint discovery |
| agent.sidecar.config.blob_storage.s3.enable_endpoint_host_prefix | bool | `true` | enable prefixing request endpoint hosts with modeled information |
| agent.sidecar.config.blob_storage.s3.enable_param_validation | bool | `true` | enables semantic parameter validation |
| agent.sidecar.config.blob_storage.s3.enable_ssl | bool | `true` | enable ssl for s3 session |
| agent.sidecar.config.blob_storage.s3.endpoint | string | `""` | s3 endpoint |
| agent.sidecar.config.blob_storage.s3.force_path_style | bool | `false` | use path-style addressing |
| agent.sidecar.config.blob_storage.s3.max_chunk_size | string | `"64mb"` | s3 download max chunk size |
| agent.sidecar.config.blob_storage.s3.max_part_size | string | `"64mb"` | s3 multipart upload max part size |
| agent.sidecar.config.blob_storage.s3.max_retries | int | `3` | maximum number of retries of s3 client |
| agent.sidecar.config.blob_storage.s3.region | string | `""` | s3 region |
| agent.sidecar.config.blob_storage.s3.secret_access_key | string | `"_AWS_SECRET_ACCESS_KEY_"` | s3 secret access key |
| agent.sidecar.config.blob_storage.s3.token | string | `""` | s3 token |
| agent.sidecar.config.blob_storage.s3.use_accelerate | bool | `false` | enable s3 accelerate feature |
| agent.sidecar.config.blob_storage.s3.use_arn_region | bool | `false` | s3 service client to use the region specified in the ARN |
| agent.sidecar.config.blob_storage.s3.use_dual_stack | bool | `false` | use dual stack |
| agent.sidecar.config.blob_storage.storage_type | string | `"s3"` | storage type |
| agent.sidecar.config.client.tcp.dialer.dual_stack_enabled | bool | `false` | HTTP client TCP dialer dual stack enabled |
| agent.sidecar.config.client.tcp.dialer.keep_alive | string | `"5m"` | HTTP client TCP dialer keep alive |
| agent.sidecar.config.client.tcp.dialer.timeout | string | `"5s"` | HTTP client TCP dialer connect timeout |
| agent.sidecar.config.client.tcp.dns.cache_enabled | bool | `true` | HTTP client TCP DNS cache enabled |
| agent.sidecar.config.client.tcp.dns.cache_expiration | string | `"24h"` |  |
| agent.sidecar.config.client.tcp.dns.refresh_duration | string | `"1h"` | HTTP client TCP DNS cache expiration |
| agent.sidecar.config.client.tcp.tls.ca | string | `"/path/to/ca"` | HTTP client TCP TLS ca path |
| agent.sidecar.config.client.tcp.tls.cert | string | `"/path/to/cert"` | HTTP client TCP TLS cert path |
| agent.sidecar.config.client.tcp.tls.enabled | bool | `false` | HTTP client TCP TLS enabled |
| agent.sidecar.config.client.tcp.tls.key | string | `"/path/to/key"` | HTTP client TCP TLS key path |
| agent.sidecar.config.client.transport.backoff.backoff_factor | float | `1.1` | backoff backoff factor |
| agent.sidecar.config.client.transport.backoff.backoff_time_limit | string | `"5s"` | backoff time limit |
| agent.sidecar.config.client.transport.backoff.enable_error_log | bool | `true` | backoff error log enabled |
| agent.sidecar.config.client.transport.backoff.initial_duration | string | `"5ms"` | backoff initial duration |
| agent.sidecar.config.client.transport.backoff.jitter_limit | string | `"100ms"` | backoff jitter limit |
| agent.sidecar.config.client.transport.backoff.maximum_duration | string | `"5s"` | backoff maximum duration |
| agent.sidecar.config.client.transport.backoff.retry_count | int | `100` | backoff retry count |
| agent.sidecar.config.client.transport.round_tripper.expect_continue_timeout | string | `"5s"` | expect continue timeout |
| agent.sidecar.config.client.transport.round_tripper.force_attempt_http_2 | bool | `true` | force attempt HTTP2 |
| agent.sidecar.config.client.transport.round_tripper.idle_conn_timeout | string | `"90s"` | timeout for idle connections |
| agent.sidecar.config.client.transport.round_tripper.max_conns_per_host | int | `10` | maximum count of connections per host |
| agent.sidecar.config.client.transport.round_tripper.max_idle_conns | int | `100` | maximum count of idle connections |
| agent.sidecar.config.client.transport.round_tripper.max_idle_conns_per_host | int | `10` | maximum count of idle connections per host |
| agent.sidecar.config.client.transport.round_tripper.max_response_header_size | int | `0` | maximum response header size |
| agent.sidecar.config.client.transport.round_tripper.read_buffer_size | int | `0` | read buffer size |
| agent.sidecar.config.client.transport.round_tripper.response_header_timeout | string | `"5s"` | timeout for response header |
| agent.sidecar.config.client.transport.round_tripper.tls_handshake_timeout | string | `"5s"` | TLS handshake timeout |
| agent.sidecar.config.client.transport.round_tripper.write_buffer_size | int | `0` | write buffer size |
| agent.sidecar.config.compress.compress_algorithm | string | `"gzip"` | compression algorithm. must be `gob`, `gzip`, `lz4` or `zstd` |
| agent.sidecar.config.compress.compression_level | int | `-1` | compression level. value range relies on which algorithm is used. `gob`: level will be ignored. `gzip`: -1 (default compression), 0 (no compression), or 1 (best speed) to 9 (best compression). `lz4`: >= 0, higher is better compression. `zstd`: 1 (fastest) to 22 (best), however implementation relies on klauspost/compress. |
| agent.sidecar.config.filename | string | `"_MY_POD_NAME_"` | backup filename |
| agent.sidecar.config.filename_suffix | string | `".tar.gz"` | suffix for backup filename |
| agent.sidecar.config.post_stop_timeout | string | `"2m"` | timeout for observing file changes during post stop |
| agent.sidecar.config.restore_backoff.backoff_factor | float | `1.2` | restore backoff factor |
| agent.sidecar.config.restore_backoff.backoff_time_limit | string | `"30m"` | restore backoff time limit |
| agent.sidecar.config.restore_backoff.enable_error_log | bool | `true` | restore backoff log enabled |
| agent.sidecar.config.restore_backoff.initial_duration | string | `"1s"` | restore backoff initial duration |
| agent.sidecar.config.restore_backoff.jitter_limit | string | `"10s"` | restore backoff jitter limit |
| agent.sidecar.config.restore_backoff.maximum_duration | string | `"1m"` | restore backoff maximum duration |
| agent.sidecar.config.restore_backoff.retry_count | int | `100` | restore backoff retry count |
| agent.sidecar.config.restore_backoff_enabled | bool | `false` | restore backoff enabled |
| agent.sidecar.config.watch_enabled | bool | `true` | auto backup triggered by file changes is enabled |
| agent.sidecar.enabled | bool | `false` | sidecar enabled |
| agent.sidecar.env | list | `[{"name":"MY_POD_NAME","valueFrom":{"fieldRef":{"fieldPath":"metadata.name"}}},{"name":"AWS_ACCESS_KEY","valueFrom":{"secretKeyRef":{"key":"access-key","name":"aws-secret"}}},{"name":"AWS_SECRET_ACCESS_KEY","valueFrom":{"secretKeyRef":{"key":"secret-access-key","name":"aws-secret"}}}]` | environment variables |
| agent.sidecar.image.pullPolicy | string | `"Always"` | image pull policy |
| agent.sidecar.image.repository | string | `"vdaas/vald-agent-sidecar"` | image repository |
| agent.sidecar.image.tag | string | `""` | image tag (overrides defaults.image.tag) |
| agent.sidecar.initContainerEnabled | bool | `false` | sidecar on initContainer mode enabled. |
| agent.sidecar.logging | object | `{}` | logging config (overrides defaults.logging) |
| agent.sidecar.name | string | `"vald-agent-sidecar"` | name of agent sidecar |
| agent.sidecar.observability | object | `{"jaeger":{"service_name":"vald-agent-sidecar"},"stackdriver":{"profiler":{"service":"vald-agent-sidecar"}}}` | observability config (overrides defaults.observability) |
| agent.sidecar.resources | object | `{"requests":{"cpu":"100m","memory":"100Mi"}}` | compute resources. |
| agent.sidecar.server_config | object | `{"healths":{"liveness":{"enabled":false,"port":13000,"servicePort":13000},"readiness":{"enabled":false,"port":13001,"servicePort":13001}},"metrics":{"pprof":{"port":16060,"servicePort":16060},"prometheus":{"port":16061,"servicePort":16061}},"servers":{"grpc":{"enabled":false,"port":18081,"servicePort":18081},"rest":{"enabled":false,"port":18080,"servicePort":18080}}}` | server config (overrides defaults.server_config) |
| agent.sidecar.service.annotations | object | `{}` | agent sidecar service annotations |
| agent.sidecar.service.enabled | bool | `false` | agent sidecar service enabled |
| agent.sidecar.service.externalTrafficPolicy | string | `""` | external traffic policy (can be specified when service type is LoadBalancer or NodePort) : Cluster or Local |
| agent.sidecar.service.labels | object | `{}` | agent sidecar service labels |
| agent.sidecar.service.type | string | `"ClusterIP"` | service type: ClusterIP, LoadBalancer or NodePort |
| agent.sidecar.time_zone | string | `""` | Time zone |
| agent.sidecar.version | string | `"v0.0.0"` | version of agent sidecar config |
| agent.terminationGracePeriodSeconds | int | `120` | duration in seconds pod needs to terminate gracefully |
| agent.time_zone | string | `""` | Time zone |
| agent.tolerations | list | `[]` | tolerations |
| agent.topologySpreadConstraints | list | `[]` | topology spread constraints for agent pods |
| agent.version | string | `"v0.0.0"` | version of agent config |
| agent.volumeMounts | list | `[]` | volume mounts |
| agent.volumes | list | `[]` | volumes |
| backupManager.affinity.nodeAffinity.preferredDuringSchedulingIgnoredDuringExecution | list | `[]` | node affinity preferred scheduling terms |
| backupManager.affinity.nodeAffinity.requiredDuringSchedulingIgnoredDuringExecution.nodeSelectorTerms | list | `[]` | node affinity required node selectors |
| backupManager.affinity.podAffinity.preferredDuringSchedulingIgnoredDuringExecution | list | `[]` | pod affinity preferred scheduling terms |
| backupManager.affinity.podAffinity.requiredDuringSchedulingIgnoredDuringExecution | list | `[]` | pod affinity required scheduling terms |
| backupManager.affinity.podAntiAffinity.preferredDuringSchedulingIgnoredDuringExecution | list | `[]` | pod anti-affinity preferred scheduling terms |
| backupManager.affinity.podAntiAffinity.requiredDuringSchedulingIgnoredDuringExecution | list | `[]` | pod anti-affinity required scheduling terms |
| backupManager.annotations | object | `{}` | deployment annotations |
| backupManager.cassandra.config.connect_timeout | string | `"3s"` | connect timeout |
| backupManager.cassandra.config.consistency | string | `"quorum"` | consistency type |
| backupManager.cassandra.config.cql_version | string | `"3.0.0"` | cassandra CQL version |
| backupManager.cassandra.config.default_idempotence | bool | `false` | default idempotence enabled |
| backupManager.cassandra.config.default_timestamp | bool | `true` | default timestamp enabled |
| backupManager.cassandra.config.disable_initial_host_lookup | bool | `false` | initial host lookup disabled |
| backupManager.cassandra.config.disable_node_status_events | bool | `false` | node status events disabled |
| backupManager.cassandra.config.disable_skip_metadata | bool | `false` | skip metadata disabled |
| backupManager.cassandra.config.disable_topology_events | bool | `false` | topology events disabled |
| backupManager.cassandra.config.enable_host_verification | bool | `false` | host verification enabled |
| backupManager.cassandra.config.host_filter.data_center | string | `""` | name of data center of filtering target |
| backupManager.cassandra.config.host_filter.enabled | bool | `false` | enables host filtering |
| backupManager.cassandra.config.host_filter.white_list | list | `[]` | list of white_list which allows each connection |
| backupManager.cassandra.config.hosts | list | `["cassandra-0.cassandra.default.svc.cluster.local","cassandra-1.cassandra.default.svc.cluster.local","cassandra-2.cassandra.default.svc.cluster.local"]` | cassandra hosts |
| backupManager.cassandra.config.ignore_peer_addr | bool | `false` | ignore peer addresses |
| backupManager.cassandra.config.keyspace | string | `"vald"` | cassandra keyspace |
| backupManager.cassandra.config.max_prepared_stmts | int | `1000` | maximum number of prepared statements |
| backupManager.cassandra.config.max_routing_key_info | int | `1000` | maximum number of routing key info |
| backupManager.cassandra.config.max_wait_schema_agreement | string | `"1m"` | maximum duration to wait for schema agreement |
| backupManager.cassandra.config.vector_backup_table | string | `"backup_vector"` | table name of backup |
| backupManager.cassandra.config.num_conns | int | `2` | number of connections per hosts |
| backupManager.cassandra.config.page_size | int | `5000` | page size |
| backupManager.cassandra.config.password | string | `"_CASSANDRA_PASSWORD_"` | cassandra password |
| backupManager.cassandra.config.pool_config.data_center | string | `""` | name of data center |
| backupManager.cassandra.config.pool_config.dc_aware_routing | bool | `false` | data center aware routine enabled |
| backupManager.cassandra.config.pool_config.non_local_replicas_fallback | bool | `false` | non-local replica fallback enabled |
| backupManager.cassandra.config.pool_config.shuffle_replicas | bool | `false` | shuffle replica enabled |
| backupManager.cassandra.config.pool_config.token_aware_host_policy | bool | `false` | token aware host policy enabled |
| backupManager.cassandra.config.port | int | `9042` | cassandra port |
| backupManager.cassandra.config.proto_version | int | `0` | cassandra proto version |
| backupManager.cassandra.config.reconnect_interval | string | `"100ms"` | interval of reconnection |
| backupManager.cassandra.config.reconnection_policy.initial_interval | string | `"100ms"` | initial interval to reconnect |
| backupManager.cassandra.config.reconnection_policy.max_retries | int | `3` | maximum number of retries to reconnect |
| backupManager.cassandra.config.retry_policy.max_duration | string | `"1s"` | maximum duration to retry |
| backupManager.cassandra.config.retry_policy.min_duration | string | `"10ms"` | minimum duration to retry |
| backupManager.cassandra.config.retry_policy.num_retries | int | `3` | number of retries |
| backupManager.cassandra.config.serial_consistency | string | `"localserial"` | read consistency type |
| backupManager.cassandra.config.socket_keepalive | string | `"0s"` | socket keep alive time |
| backupManager.cassandra.config.tcp.dialer.dual_stack_enabled | bool | `false` | TCP dialer dual stack enabled |
| backupManager.cassandra.config.tcp.dialer.keep_alive | string | `"10m"` | TCP dialer keep alive |
| backupManager.cassandra.config.tcp.dialer.timeout | string | `"30s"` | TCP dialer timeout |
| backupManager.cassandra.config.tcp.dns.cache_enabled | bool | `true` | TCP DNS cache enabled |
| backupManager.cassandra.config.tcp.dns.cache_expiration | string | `"24h"` | TCP DNS cache expiration |
| backupManager.cassandra.config.tcp.dns.refresh_duration | string | `"5m"` | TCP DNS cache refresh duration |
| backupManager.cassandra.config.timeout | string | `"600ms"` | timeout |
| backupManager.cassandra.config.tls.ca | string | `"/path/to/ca"` | path to TLS ca |
| backupManager.cassandra.config.tls.cert | string | `"/path/to/cert"` | path to TLS cert |
| backupManager.cassandra.config.tls.enabled | bool | `false` | TLS enabled |
| backupManager.cassandra.config.tls.key | string | `"/path/to/key"` | path to TLS key |
| backupManager.cassandra.config.username | string | `"root"` | cassandra username |
| backupManager.cassandra.config.write_coalesce_wait_time | string | `"200µs"` | write coalesce wait time |
| backupManager.cassandra.enabled | bool | `false` | cassandra config enabled |
| backupManager.enabled | bool | `true` | backup manager enabled |
| backupManager.env | list | `[{"name":"MYSQL_PASSWORD","valueFrom":{"secretKeyRef":{"key":"password","name":"mysql-secret"}}}]` | environment variables |
| backupManager.externalTrafficPolicy | string | `""` | external traffic policy (can be specified when service type is LoadBalancer or NodePort) : Cluster or Local |
| backupManager.hpa.enabled | bool | `true` | HPA enabled |
| backupManager.hpa.targetCPUUtilizationPercentage | int | `80` | HPA CPU utilization percentage |
| backupManager.image.pullPolicy | string | `"Always"` | image pull policy |
| backupManager.image.repository | string | `"vdaas/vald-manager-backup-mysql"` | image repository |
| backupManager.image.tag | string | `""` | image tag (overrides defaults.image.tag) |
| backupManager.initContainers | list | `[{"env":[{"name":"MYSQL_PASSWORD","valueFrom":{"secretKeyRef":{"key":"password","name":"mysql-secret"}}}],"image":"mysql:latest","mysql":{"hosts":["mysql.default.svc.cluster.local"],"options":["-uroot","-p${MYSQL_PASSWORD}"]},"name":"wait-for-mysql","sleepDuration":2,"type":"wait-for-mysql"}]` | init containers |
| backupManager.kind | string | `"Deployment"` | deployment kind: Deployment or DaemonSet |
| backupManager.logging | object | `{}` | logging config (overrides defaults.logging) |
| backupManager.maxReplicas | int | `15` | maximum number of replicas. if HPA is disabled, this value will be ignored. |
| backupManager.maxUnavailable | string | `"50%"` | maximum number of unavailable replicas |
| backupManager.minReplicas | int | `3` | minimum number of replicas. if HPA is disabled, the replicas will be set to this value |
| backupManager.mysql.config.conn_max_life_time | string | `"30s"` | connection maximum life time |
| backupManager.mysql.config.db | string | `"mysql"` | mysql db: mysql, postgres or sqlite3 |
| backupManager.mysql.config.host | string | `"mysql.default.svc.cluster.local"` | mysql hostname |
| backupManager.mysql.config.max_idle_conns | int | `100` | maximum number of idle connections |
| backupManager.mysql.config.max_open_conns | int | `100` | maximum number of open connections |
| backupManager.mysql.config.name | string | `"vald"` | mysql db name |
| backupManager.mysql.config.pass | string | `"_MYSQL_PASSWORD_"` | mysql password |
| backupManager.mysql.config.port | int | `3306` | mysql port |
| backupManager.mysql.config.tcp.dialer.dual_stack_enabled | bool | `false` | TCP dialer dual stack enabled |
| backupManager.mysql.config.tcp.dialer.keep_alive | string | `"5m"` | TCP dialer keep alive |
| backupManager.mysql.config.tcp.dialer.timeout | string | `"5s"` | TCP dialer timeout |
| backupManager.mysql.config.tcp.dns.cache_enabled | bool | `true` | TCP DNS cache enabled |
| backupManager.mysql.config.tcp.dns.cache_expiration | string | `"24h"` | TCP DNS cache expiration |
| backupManager.mysql.config.tcp.dns.refresh_duration | string | `"1h"` | TCP DNS cache refresh duration |
| backupManager.mysql.config.tcp.tls.ca | string | `"/path/to/ca"` | path to TCP TLS ca |
| backupManager.mysql.config.tcp.tls.cert | string | `"/path/to/cert"` | path to TCP TLS cert |
| backupManager.mysql.config.tcp.tls.enabled | bool | `false` | TCP TLS enabled |
| backupManager.mysql.config.tcp.tls.key | string | `"/path/to/key"` | path to TCP TLS key |
| backupManager.mysql.config.tls.ca | string | `"/path/to/ca"` | path to TLS ca |
| backupManager.mysql.config.tls.cert | string | `"/path/to/cert"` | path to TLS cert |
| backupManager.mysql.config.tls.enabled | bool | `false` | TLS enabled |
| backupManager.mysql.config.tls.key | string | `"/path/to/key"` | path to TLS key |
| backupManager.mysql.config.user | string | `"root"` | mysql username |
| backupManager.mysql.enabled | bool | `true` | mysql config enabled |
| backupManager.name | string | `"vald-manager-backup"` | name of backup manager deployment |
| backupManager.nodeName | string | `""` | node name |
| backupManager.nodeSelector | object | `{}` | node selector |
| backupManager.observability | object | `{"jaeger":{"service_name":"vald-manager-backup"},"stackdriver":{"profiler":{"service":"vald-manager-backup"}}}` | observability config (overrides defaults.observability) |
| backupManager.podAnnotations | object | `{}` | pod annotations |
| backupManager.podPriority.enabled | bool | `true` | backup manager pod PriorityClass enabled |
| backupManager.podPriority.value | int | `1000000` | backup manager pod PriorityClass value |
| backupManager.progressDeadlineSeconds | int | `600` | progress deadline seconds |
| backupManager.resources | object | `{"limits":{"cpu":"500m","memory":"150Mi"},"requests":{"cpu":"100m","memory":"50Mi"}}` | compute resources |
| backupManager.revisionHistoryLimit | int | `2` | number of old history to retain to allow rollback |
| backupManager.rollingUpdate.maxSurge | string | `"25%"` | max surge of rolling update |
| backupManager.rollingUpdate.maxUnavailable | string | `"25%"` | max unavailable of rolling update |
| backupManager.server_config | object | `{"healths":{"liveness":{},"readiness":{}},"metrics":{"pprof":{},"prometheus":{}},"servers":{"grpc":{},"rest":{}}}` | server config (overrides defaults.server_config) |
| backupManager.service.annotations | object | `{}` | service annotations |
| backupManager.service.labels | object | `{}` | service labels |
| backupManager.serviceType | string | `"ClusterIP"` | service type: ClusterIP, LoadBalancer or NodePort |
| backupManager.terminationGracePeriodSeconds | int | `30` | duration in seconds pod needs to terminate gracefully |
| backupManager.time_zone | string | `""` | Time zone |
| backupManager.tolerations | list | `[]` | tolerations |
| backupManager.topologySpreadConstraints | list | `[]` | topology spread constraints of backup manager pods |
| backupManager.version | string | `"v0.0.0"` | version of backup manager config |
| backupManager.volumeMounts | list | `[]` | volume mounts |
| backupManager.volumes | list | `[]` | volumes |
| compressor.affinity.nodeAffinity.preferredDuringSchedulingIgnoredDuringExecution | list | `[]` | node affinity preferred scheduling terms |
| compressor.affinity.nodeAffinity.requiredDuringSchedulingIgnoredDuringExecution.nodeSelectorTerms | list | `[]` | node affinity required node selectors |
| compressor.affinity.podAffinity.preferredDuringSchedulingIgnoredDuringExecution | list | `[]` | pod affinity preferred scheduling terms |
| compressor.affinity.podAffinity.requiredDuringSchedulingIgnoredDuringExecution | list | `[]` | pod affinity required scheduling terms |
| compressor.affinity.podAntiAffinity.preferredDuringSchedulingIgnoredDuringExecution | list | `[]` | pod anti-affinity preferred scheduling terms |
| compressor.affinity.podAntiAffinity.requiredDuringSchedulingIgnoredDuringExecution | list | `[]` | pod anti-affinity required scheduling terms |
| compressor.annotations | object | `{}` | deployment annotations |
| compressor.backup.client | object | `{}` | grpc client for backup (overrides defaults.grpc.client) |
| compressor.compress.compress_algorithm | string | `"zstd"` | compression algorithm. must be `gob`, `gzip`, `lz4` or `zstd` |
| compressor.compress.compression_level | int | `3` | compression level. value range relies on which algorithm is used. `gob`: level will be ignored. `gzip`: -1 (default compression), 0 (no compression), or 1 (best speed) to 9 (best compression). `lz4`: >= 0, higher is better compression. `zstd`: 1 (fastest) to 22 (best), however implementation relies on klauspost/compress. |
| compressor.compress.concurrent_limit | int | `10` | concurrency limit for compress/decompress processes |
| compressor.compress.queue_check_duration | string | `"200ms"` |  |
| compressor.enabled | bool | `true` | compressor enabled |
| compressor.env | list | `[{"name":"MY_POD_IP","valueFrom":{"fieldRef":{"fieldPath":"status.podIP"}}}]` | environment variables |
| compressor.externalTrafficPolicy | string | `""` | external traffic policy (can be specified when service type is LoadBalancer or NodePort) : Cluster or Local |
| compressor.hpa.enabled | bool | `true` | HPA enabled |
| compressor.hpa.targetCPUUtilizationPercentage | int | `80` | HPA CPU utilization percentage |
| compressor.image.pullPolicy | string | `"Always"` | image pull policy |
| compressor.image.repository | string | `"vdaas/vald-manager-compressor"` | image repository |
| compressor.image.tag | string | `""` | image tag (overrides defaults.image.tag) |
| compressor.initContainers | list | `[{"image":"busybox","name":"wait-for-manager-backup","sleepDuration":2,"target":"manager-backup","type":"wait-for"}]` | init containers |
| compressor.kind | string | `"Deployment"` | deployment kind: Deployment or DaemonSet |
| compressor.logging | object | `{}` | logging config (overrides defaults.logging) |
| compressor.maxReplicas | int | `15` | maximum number of replicas. if HPA is disabled, this value will be ignored. |
| compressor.maxUnavailable | string | `"1"` | maximum number of unavailable replicas |
| compressor.minReplicas | int | `3` | minimum number of replicas. if HPA is disabled, the replicas will be set to this value |
| compressor.name | string | `"vald-manager-compressor"` | name of compressor deployment |
| compressor.nodeName | string | `""` | node name |
| compressor.nodeSelector | object | `{}` | node selector |
| compressor.observability | object | `{"jaeger":{"service_name":"vald-manager-compressor"},"stackdriver":{"profiler":{"service":"vald-manager-compressor"}}}` | observability config (overrides defaults.observability) |
| compressor.podAnnotations | object | `{}` | pod annotations |
| compressor.podPriority.enabled | bool | `true` | compressor pod PriorityClass enabled |
| compressor.podPriority.value | int | `100000000` | compressor pod PriorityClass value |
| compressor.progressDeadlineSeconds | int | `600` | progress deadline seconds |
| compressor.registerer.compressor.client | object | `{}` | gRPC client for compressor (overrides defaults.grpc.client) |
| compressor.registerer.concurrent_limit | int | `10` | concurrency limit for registering vector processes |
| compressor.registerer.queue_check_duration | string | `"200ms"` |  |
| compressor.resources | object | `{"limits":{"cpu":"800m","memory":"500Mi"},"requests":{"cpu":"300m","memory":"50Mi"}}` | compute resources |
| compressor.revisionHistoryLimit | int | `2` | number of old history to retain to allow rollback |
| compressor.rollingUpdate.maxSurge | string | `"25%"` | max surge of rolling update |
| compressor.rollingUpdate.maxUnavailable | string | `"25%"` | max unavailable of rolling update |
| compressor.server_config | object | `{"healths":{"liveness":{"server":{"http":{"shutdown_duration":"2m"}}},"readiness":{}},"metrics":{"pprof":{},"prometheus":{}},"servers":{"grpc":{},"rest":{}}}` | server config (overrides defaults.server_config) |
| compressor.service.annotations | object | `{}` | service annotations |
| compressor.service.labels | object | `{}` | service labels |
| compressor.serviceType | string | `"ClusterIP"` | service type: ClusterIP, LoadBalancer or NodePort |
| compressor.terminationGracePeriodSeconds | int | `120` | duration in seconds pod needs to terminate gracefully |
| compressor.time_zone | string | `""` | Time zone |
| compressor.tolerations | list | `[]` | tolerations |
| compressor.topologySpreadConstraints | list | `[]` | topology spread constraints of compressor pods |
| compressor.version | string | `"v0.0.0"` | version of compressor config |
| compressor.volumeMounts | list | `[]` | volume mounts |
| compressor.volumes | list | `[]` | volumes |
| defaults.grpc.client.addrs | list | `[]` | gRPC client addresses |
| defaults.grpc.client.backoff.backoff_factor | float | `1.1` | gRPC client backoff factor |
| defaults.grpc.client.backoff.backoff_time_limit | string | `"5s"` | gRPC client backoff time limit |
| defaults.grpc.client.backoff.enable_error_log | bool | `true` | gRPC client backoff log enabled |
| defaults.grpc.client.backoff.initial_duration | string | `"5ms"` | gRPC client backoff initial duration |
| defaults.grpc.client.backoff.jitter_limit | string | `"100ms"` | gRPC client backoff jitter limit |
| defaults.grpc.client.backoff.maximum_duration | string | `"5s"` | gRPC client backoff maximum duration |
| defaults.grpc.client.backoff.retry_count | int | `100` | gRPC client backoff retry count |
| defaults.grpc.client.call_option.max_recv_msg_size | int | `0` | gRPC client call option max receive message size |
| defaults.grpc.client.call_option.max_retry_rpc_buffer_size | int | `0` | gRPC client call option max retry rpc buffer size |
| defaults.grpc.client.call_option.max_send_msg_size | int | `0` | gRPC client call option max send message size |
| defaults.grpc.client.call_option.wait_for_ready | bool | `true` | gRPC client call option wait for ready |
| defaults.grpc.client.connection_pool.enable_dns_resolver | bool | `true` | enables gRPC client connection pool dns resolver, when enabled vald uses ip handshake exclude dns discovery which improves network performance |
| defaults.grpc.client.connection_pool.enable_rebalance | bool | `true` | enables gRPC client connection pool rebalance |
| defaults.grpc.client.connection_pool.old_conn_close_duration | string | `"3s"` | makes delay before gRPC client connection closing during connection pool rebalance |
| defaults.grpc.client.connection_pool.rebalance_duration | string | `"30m"` | gRPC client connection pool rebalance duration |
| defaults.grpc.client.connection_pool.size | int | `3` | gRPC client connection pool size |
| defaults.grpc.client.dial_option.enable_backoff | bool | `false` | gRPC client dial option backoff enabled |
| defaults.grpc.client.dial_option.initial_connection_window_size | int | `0` | gRPC client dial option initial connection window size |
| defaults.grpc.client.dial_option.initial_window_size | int | `0` | gRPC client dial option initial window size |
| defaults.grpc.client.dial_option.insecure | bool | `true` | gRPC client dial option insecure enabled |
| defaults.grpc.client.dial_option.keep_alive.permit_without_stream | bool | `false` | gRPC client keep alive permit without stream |
| defaults.grpc.client.dial_option.keep_alive.time | string | `""` | gRPC client keep alive time |
| defaults.grpc.client.dial_option.keep_alive.timeout | string | `""` | gRPC client keep alive timeout |
| defaults.grpc.client.dial_option.backoff_max_delay | string | `""` | gRPC client dial option max backoff delay |
| defaults.grpc.client.dial_option.max_msg_size | int | `0` | gRPC client dial option max message size |
| defaults.grpc.client.dial_option.read_buffer_size | int | `0` | gRPC client dial option read buffer size |
| defaults.grpc.client.dial_option.tcp.dialer.dual_stack_enabled | bool | `true` | gRPC client TCP dialer dual stack enabled |
| defaults.grpc.client.dial_option.tcp.dialer.keep_alive | string | `""` | gRPC client TCP dialer keep alive |
| defaults.grpc.client.dial_option.tcp.dialer.timeout | string | `""` | gRPC client TCP dialer timeout |
| defaults.grpc.client.dial_option.tcp.dns.cache_enabled | bool | `true` | gRPC client TCP DNS cache enabled |
| defaults.grpc.client.dial_option.tcp.dns.cache_expiration | string | `"1h"` | gRPC client TCP DNS cache expiration |
| defaults.grpc.client.dial_option.tcp.dns.refresh_duration | string | `"30m"` | gRPC client TCP DNS cache refresh duration |
| defaults.grpc.client.dial_option.tcp.tls.ca | string | `"/path/to/ca"` | gRPC client TCP TLS ca path |
| defaults.grpc.client.dial_option.tcp.tls.cert | string | `"/path/to/cert"` | gRPC client TCP TLS cert path |
| defaults.grpc.client.dial_option.tcp.tls.enabled | bool | `false` | gRPC client TCP TLS enabled |
| defaults.grpc.client.dial_option.tcp.tls.key | string | `"/path/to/key"` | gRPC client TCP TLS key path |
| defaults.grpc.client.dial_option.timeout | string | `""` | gRPC client dial option timeout |
| defaults.grpc.client.dial_option.write_buffer_size | int | `0` | gRPC client dial option write buffer size |
| defaults.grpc.client.health_check_duration | string | `"1s"` | gRPC client health check duration |
| defaults.grpc.client.tls.ca | string | `"/path/to/ca"` | gRPC client TLS ca path |
| defaults.grpc.client.tls.cert | string | `"/path/to/cert"` | gRPC client TLS cert path |
| defaults.grpc.client.tls.enabled | bool | `false` | gRPC client TLS enabled |
| defaults.grpc.client.tls.key | string | `"/path/to/key"` | gRPC client TLS key path |
| defaults.image.tag | string | `"v0.0.66"` | docker image tag |
| defaults.logging.format | string | `"raw"` | logging format. logging format must be `raw` or `json` |
| defaults.logging.level | string | `"debug"` | logging level. logging level must be `debug`, `info`, `warn`, `error` or `fatal`. |
| defaults.logging.logger | string | `"glg"` | logger name. currently logger must be `glg`. |
| defaults.observability.collector.duration | string | `"5s"` | metrics collect duration. if it is set as 5s, enabled metrics are collected every 5 seconds. |
| defaults.observability.collector.metrics.enable_cgo | bool | `true` | CGO metrics enabled |
| defaults.observability.collector.metrics.enable_goroutine | bool | `true` | goroutine metrics enabled |
| defaults.observability.collector.metrics.enable_memory | bool | `true` | memory metrics enabled |
| defaults.observability.collector.metrics.enable_version_info | bool | `true` | version info metrics enabled |
| defaults.observability.collector.metrics.version_info_labels | list | `["vald_version","server_name","git_commit","build_time","go_version","go_os","go_arch","ngt_version"]` | enabled label names of version info |
| defaults.observability.enabled | bool | `false` | observability features enabled |
| defaults.observability.jaeger.agent_endpoint | string | `"jaeger-agent.default.svc.cluster.local:6831"` | Jaeger agent endpoint |
| defaults.observability.jaeger.buffer_max_count | int | `10` | Jaeger buffer max count |
| defaults.observability.jaeger.collector_endpoint | string | `""` | Jaeger collector endpoint |
| defaults.observability.jaeger.enabled | bool | `false` | Jaeger exporter enabled |
| defaults.observability.jaeger.password | string | `""` | Jaeger password |
| defaults.observability.jaeger.service_name | string | `"vald"` | Jaeger service name |
| defaults.observability.jaeger.username | string | `""` | Jaeger username |
| defaults.observability.prometheus.enabled | bool | `false` | Prometheus exporter enabled |
| defaults.observability.prometheus.endpoint | string | `"/metrics"` | Prometheus exporter endpoint |
| defaults.observability.prometheus.namespace | string | `"vald"` | prefix of exported metrics name |
| defaults.observability.stackdriver.client.api_key | string | `""` | API key to be used as the basis for authentication. |
| defaults.observability.stackdriver.client.audiences | list | `[]` | to be used as the audience field ("aud") for the JWT token authentication. |
| defaults.observability.stackdriver.client.authentication_enabled | bool | `true` | enables authentication. |
| defaults.observability.stackdriver.client.credentials_file | string | `""` | service account or refresh token JSON credentials file. |
| defaults.observability.stackdriver.client.credentials_json | string | `""` | service account or refresh token JSON credentials. |
| defaults.observability.stackdriver.client.endpoint | string | `""` | overrides the default endpoint to be used for a service. |
| defaults.observability.stackdriver.client.quota_project | string | `""` | the project used for quota and billing purposes. |
| defaults.observability.stackdriver.client.request_reason | string | `""` | a reason for making the request, which is intended to be recorded in audit logging. |
| defaults.observability.stackdriver.client.scopes | list | `[]` | overrides the default OAuth2 scopes to be used for a service. |
| defaults.observability.stackdriver.client.telemetry_enabled | bool | `true` | enables default telemetry settings on gRPC and HTTP clients. |
| defaults.observability.stackdriver.client.user_agent | string | `""` | sets the User-Agent. |
| defaults.observability.stackdriver.exporter.bundle_count_threshold | int | `0` | how many view data events or trace spans can be buffered. |
| defaults.observability.stackdriver.exporter.bundle_delay_threshold | string | `"0"` | the max amount of time the exporter can wait before uploading data. |
| defaults.observability.stackdriver.exporter.location | string | `""` | identifier of the GCP or AWS cloud region/zone the data is stored. |
| defaults.observability.stackdriver.exporter.metric_prefix | string | `"vald.vdaas.org"` | the prefix of a stackdriver metric names. |
| defaults.observability.stackdriver.exporter.monitoring_enabled | bool | `false` | stackdriver monitoring enabled |
| defaults.observability.stackdriver.exporter.number_of_workers | int | `1` | number of workers |
| defaults.observability.stackdriver.exporter.reporting_interval | string | `"1m"` | interval between reporting metrics |
| defaults.observability.stackdriver.exporter.skip_cmd | bool | `false` | skip all the CreateMetricDescriptor calls |
| defaults.observability.stackdriver.exporter.timeout | string | `"5s"` | timeout for all API calls |
| defaults.observability.stackdriver.exporter.trace_spans_buffer_max_bytes | int | `0` | maximum size of spans that will be buffered. |
| defaults.observability.stackdriver.exporter.tracing_enabled | bool | `false` | stackdriver tracing enabled |
| defaults.observability.stackdriver.profiler.alloc_force_gc | bool | `false` | forces GC before the collection of each heap profile. |
| defaults.observability.stackdriver.profiler.alloc_profiling | bool | `true` | enables allocation profiling. |
| defaults.observability.stackdriver.profiler.api_addr | string | `""` | HTTP endpoint to use to connect to the profiler agent API. |
| defaults.observability.stackdriver.profiler.cpu_profiling | bool | `true` | enables CPU profiling. |
| defaults.observability.stackdriver.profiler.debug_logging | bool | `false` | enables detailed logging from profiler. |
| defaults.observability.stackdriver.profiler.enabled | bool | `false` | stackdriver profiler enabled. |
| defaults.observability.stackdriver.profiler.goroutine_profiling | bool | `true` | enables goroutine profiling. |
| defaults.observability.stackdriver.profiler.heap_profiling | bool | `true` | enables heap profiling. |
| defaults.observability.stackdriver.profiler.instance | string | `""` | the name of Compute Engine instance. This is normally determined from the Compute Engine metadata server and doesn't need to be initialized. |
| defaults.observability.stackdriver.profiler.mutex_profiling | bool | `true` | enables mutex profiling. |
| defaults.observability.stackdriver.profiler.service | string | `""` | the name of the service. |
| defaults.observability.stackdriver.profiler.service_version | string | `""` | the version of the service. |
| defaults.observability.stackdriver.profiler.zone | string | `""` | the zone of Compute Engine instance. This is normally determined from the Compute Engine metadata server and doesn't need to be initialized. |
| defaults.observability.stackdriver.project_id | string | `""` | project id for uploading the stats data |
| defaults.observability.trace.enabled | bool | `false` | trace enabled |
| defaults.observability.trace.sampling_rate | float | `1` | trace sampling rate. must be between 0.0 to 1.0. |
| defaults.server_config.full_shutdown_duration | string | `"600s"` | server full shutdown duration |
| defaults.server_config.healths.liveness.enabled | bool | `true` | liveness server enabled |
| defaults.server_config.healths.liveness.host | string | `"0.0.0.0"` | liveness server host |
| defaults.server_config.healths.liveness.livenessProbe.failureThreshold | int | `2` | liveness probe failure threshold |
| defaults.server_config.healths.liveness.livenessProbe.httpGet.path | string | `"/liveness"` | liveness probe path |
| defaults.server_config.healths.liveness.livenessProbe.httpGet.port | string | `"liveness"` | liveness probe port |
| defaults.server_config.healths.liveness.livenessProbe.httpGet.scheme | string | `"HTTP"` | liveness probe scheme |
| defaults.server_config.healths.liveness.livenessProbe.initialDelaySeconds | int | `5` | liveness probe initial delay seconds |
| defaults.server_config.healths.liveness.livenessProbe.periodSeconds | int | `3` | liveness probe period seconds |
| defaults.server_config.healths.liveness.livenessProbe.successThreshold | int | `1` | liveness probe success threshold |
| defaults.server_config.healths.liveness.livenessProbe.timeoutSeconds | int | `2` | liveness probe timeout seconds |
| defaults.server_config.healths.liveness.port | int | `3000` | liveness server port |
| defaults.server_config.healths.liveness.server.http.handler_timeout | string | `""` | liveness server handler timeout |
| defaults.server_config.healths.liveness.server.http.idle_timeout | string | `""` | liveness server idle timeout |
| defaults.server_config.healths.liveness.server.http.read_header_timeout | string | `""` | liveness server read header timeout |
| defaults.server_config.healths.liveness.server.http.read_timeout | string | `""` | liveness server read timeout |
| defaults.server_config.healths.liveness.server.http.shutdown_duration | string | `"5s"` | liveness server shutdown duration |
| defaults.server_config.healths.liveness.server.http.write_timeout | string | `""` | liveness server write timeout |
| defaults.server_config.healths.liveness.server.mode | string | `""` | liveness server mode |
| defaults.server_config.healths.liveness.server.probe_wait_time | string | `"3s"` | liveness server probe wait time |
| defaults.server_config.healths.liveness.servicePort | int | `3000` | liveness server service port |
| defaults.server_config.healths.readiness.enabled | bool | `true` | readiness server enabled |
| defaults.server_config.healths.readiness.host | string | `"0.0.0.0"` | readiness server host |
| defaults.server_config.healths.readiness.port | int | `3001` | readiness server port |
| defaults.server_config.healths.readiness.readinessProbe.failureThreshold | int | `2` | readiness probe failure threshold |
| defaults.server_config.healths.readiness.readinessProbe.httpGet.path | string | `"/readiness"` | readiness probe path |
| defaults.server_config.healths.readiness.readinessProbe.httpGet.port | string | `"readiness"` | readiness probe port |
| defaults.server_config.healths.readiness.readinessProbe.httpGet.scheme | string | `"HTTP"` | readiness probe scheme |
| defaults.server_config.healths.readiness.readinessProbe.initialDelaySeconds | int | `10` | readiness probe initial delay seconds |
| defaults.server_config.healths.readiness.readinessProbe.periodSeconds | int | `3` | readiness probe period seconds |
| defaults.server_config.healths.readiness.readinessProbe.successThreshold | int | `1` | readiness probe success threshold |
| defaults.server_config.healths.readiness.readinessProbe.timeoutSeconds | int | `2` | readiness probe timeout seconds |
| defaults.server_config.healths.readiness.server.http.handler_timeout | string | `""` | readiness server handler timeout |
| defaults.server_config.healths.readiness.server.http.idle_timeout | string | `""` | readiness server idle timeout |
| defaults.server_config.healths.readiness.server.http.read_header_timeout | string | `""` | readiness server read header timeout |
| defaults.server_config.healths.readiness.server.http.read_timeout | string | `""` | readiness server read timeout |
| defaults.server_config.healths.readiness.server.http.shutdown_duration | string | `"0s"` | readiness server shutdown duration |
| defaults.server_config.healths.readiness.server.http.write_timeout | string | `""` | readiness server write timeout |
| defaults.server_config.healths.readiness.server.mode | string | `""` | readiness server mode |
| defaults.server_config.healths.readiness.server.probe_wait_time | string | `"3s"` | readiness server probe wait time |
| defaults.server_config.healths.readiness.servicePort | int | `3001` | readiness server service port |
| defaults.server_config.metrics.pprof.enabled | bool | `false` | pprof server enabled |
| defaults.server_config.metrics.pprof.host | string | `"0.0.0.0"` | pprof server host |
| defaults.server_config.metrics.pprof.port | int | `6060` | pprof server port |
| defaults.server_config.metrics.pprof.server.http.handler_timeout | string | `"5s"` | pprof server handler timeout |
| defaults.server_config.metrics.pprof.server.http.idle_timeout | string | `"2s"` | pprof server idle timeout |
| defaults.server_config.metrics.pprof.server.http.read_header_timeout | string | `"1s"` | pprof server read header timeout |
| defaults.server_config.metrics.pprof.server.http.read_timeout | string | `"1s"` | pprof server read timeout |
| defaults.server_config.metrics.pprof.server.http.shutdown_duration | string | `"5s"` | pprof server shutdown duration |
| defaults.server_config.metrics.pprof.server.http.write_timeout | string | `"1s"` | pprof server write timeout |
| defaults.server_config.metrics.pprof.server.mode | string | `"REST"` | pprof server mode |
| defaults.server_config.metrics.pprof.server.probe_wait_time | string | `"3s"` | pprof server probe wait time |
| defaults.server_config.metrics.pprof.servicePort | int | `6060` | pprof server service port |
| defaults.server_config.metrics.prometheus.enabled | bool | `false` | prometheus server enabled |
| defaults.server_config.metrics.prometheus.host | string | `"0.0.0.0"` | prometheus server host |
| defaults.server_config.metrics.prometheus.port | int | `6061` | prometheus server port |
| defaults.server_config.metrics.prometheus.server.http.handler_timeout | string | `"5s"` | prometheus server handler timeout |
| defaults.server_config.metrics.prometheus.server.http.idle_timeout | string | `"2s"` | prometheus server idle timeout |
| defaults.server_config.metrics.prometheus.server.http.read_header_timeout | string | `"1s"` | prometheus server read header timeout |
| defaults.server_config.metrics.prometheus.server.http.read_timeout | string | `"1s"` | prometheus server read timeout |
| defaults.server_config.metrics.prometheus.server.http.shutdown_duration | string | `"5s"` | prometheus server shutdown duration |
| defaults.server_config.metrics.prometheus.server.http.write_timeout | string | `"1s"` | prometheus server write timeout |
| defaults.server_config.metrics.prometheus.server.mode | string | `"REST"` | prometheus server mode |
| defaults.server_config.metrics.prometheus.server.probe_wait_time | string | `"3s"` | prometheus server probe wait time |
| defaults.server_config.metrics.prometheus.servicePort | int | `6061` | prometheus server service port |
| defaults.server_config.servers.grpc.enabled | bool | `true` | gRPC server enabled |
| defaults.server_config.servers.grpc.host | string | `"0.0.0.0"` | gRPC server host |
| defaults.server_config.servers.grpc.port | int | `8081` | gRPC server port |
| defaults.server_config.servers.grpc.server.grpc.bidirectional_stream_concurrency | int | `20` | gRPC server bidirectional stream concurrency |
| defaults.server_config.servers.grpc.server.grpc.connection_timeout | string | `""` | gRPC server connection timeout |
| defaults.server_config.servers.grpc.server.grpc.header_table_size | int | `0` | gRPC server header table size |
| defaults.server_config.servers.grpc.server.grpc.initial_conn_window_size | int | `0` | gRPC server initial connection window size |
| defaults.server_config.servers.grpc.server.grpc.initial_window_size | int | `0` | gRPC server initial window size |
| defaults.server_config.servers.grpc.server.grpc.interceptors | list | `[]` | gRPC server interceptors |
| defaults.server_config.servers.grpc.server.grpc.keepalive.max_conn_age | string | `""` | gRPC server keep alive max connection age |
| defaults.server_config.servers.grpc.server.grpc.keepalive.max_conn_age_grace | string | `""` | gRPC server keep alive max connection age grace |
| defaults.server_config.servers.grpc.server.grpc.keepalive.max_conn_idle | string | `""` | gRPC server keep alive max connection idle |
| defaults.server_config.servers.grpc.server.grpc.keepalive.time | string | `""` | gRPC server keep alive time |
| defaults.server_config.servers.grpc.server.grpc.keepalive.timeout | string | `""` | gRPC server keep alive timeout |
| defaults.server_config.servers.grpc.server.grpc.max_header_list_size | int | `0` | gRPC server max header list size |
| defaults.server_config.servers.grpc.server.grpc.max_receive_message_size | int | `0` | gRPC server max receive message size |
| defaults.server_config.servers.grpc.server.grpc.max_send_message_size | int | `0` | gRPC server max send message size |
| defaults.server_config.servers.grpc.server.grpc.read_buffer_size | int | `0` | gRPC server read buffer size |
| defaults.server_config.servers.grpc.server.grpc.write_buffer_size | int | `0` | gRPC server write buffer size |
| defaults.server_config.servers.grpc.server.mode | string | `"GRPC"` | gRPC server server mode |
| defaults.server_config.servers.grpc.server.probe_wait_time | string | `"3s"` | gRPC server probe wait time |
| defaults.server_config.servers.grpc.server.restart | bool | `true` | gRPC server restart |
| defaults.server_config.servers.grpc.servicePort | int | `8081` | gRPC server service port |
| defaults.server_config.servers.rest.enabled | bool | `false` | REST server enabled |
| defaults.server_config.servers.rest.host | string | `"0.0.0.0"` | REST server host |
| defaults.server_config.servers.rest.port | int | `8080` | REST server port |
| defaults.server_config.servers.rest.server.http.handler_timeout | string | `"5s"` | REST server handler timeout |
| defaults.server_config.servers.rest.server.http.idle_timeout | string | `"2s"` | REST server idle timeout |
| defaults.server_config.servers.rest.server.http.read_header_timeout | string | `"1s"` | REST server read header timeout |
| defaults.server_config.servers.rest.server.http.read_timeout | string | `"1s"` | REST server read timeout |
| defaults.server_config.servers.rest.server.http.shutdown_duration | string | `"5s"` | REST server shutdown duration |
| defaults.server_config.servers.rest.server.http.write_timeout | string | `"1s"` | REST server write timeout |
| defaults.server_config.servers.rest.server.mode | string | `"REST"` | REST server server mode |
| defaults.server_config.servers.rest.server.probe_wait_time | string | `"3s"` | REST server probe wait time |
| defaults.server_config.servers.rest.servicePort | int | `8080` | REST server service port |
| defaults.server_config.tls.ca | string | `"/path/to/ca"` | TLS ca path |
| defaults.server_config.tls.cert | string | `"/path/to/cert"` | TLS cert path |
| defaults.server_config.tls.enabled | bool | `false` | TLS enabled |
| defaults.server_config.tls.key | string | `"/path/to/key"` | TLS key path |
| defaults.time_zone | string | `"UTC"` | Time zone |
| discoverer.affinity.nodeAffinity.preferredDuringSchedulingIgnoredDuringExecution | list | `[]` | node affinity preferred scheduling terms |
| discoverer.affinity.nodeAffinity.requiredDuringSchedulingIgnoredDuringExecution.nodeSelectorTerms | list | `[]` | node affinity required node selectors |
| discoverer.affinity.podAffinity.preferredDuringSchedulingIgnoredDuringExecution | list | `[]` | pod affinity preferred scheduling terms |
| discoverer.affinity.podAffinity.requiredDuringSchedulingIgnoredDuringExecution | list | `[]` | pod affinity required scheduling terms |
| discoverer.affinity.podAntiAffinity.preferredDuringSchedulingIgnoredDuringExecution | list | `[{"podAffinityTerm":{"labelSelector":{"matchExpressions":[{"key":"app","operator":"In","values":["vald-discoverer"]}]},"topologyKey":"kubernetes.io/hostname"},"weight":100}]` | pod anti-affinity preferred scheduling terms |
| discoverer.affinity.podAntiAffinity.requiredDuringSchedulingIgnoredDuringExecution | list | `[]` | pod anti-affinity required scheduling terms |
| discoverer.annotations | object | `{}` | deployment annotations |
| discoverer.clusterRole.enabled | bool | `true` | creates clusterRole resource |
| discoverer.clusterRole.name | string | `"discoverer"` | name of clusterRole |
| discoverer.clusterRoleBinding.enabled | bool | `true` | creates clusterRoleBinding resource |
| discoverer.clusterRoleBinding.name | string | `"discoverer"` | name of clusterRoleBinding |
| discoverer.discoverer.discovery_duration | string | `"3s"` | duration to discovery |
| discoverer.discoverer.name | string | `""` | name to discovery |
| discoverer.discoverer.namespace | string | `"_MY_POD_NAMESPACE_"` | namespace to discovery |
| discoverer.enabled | bool | `true` | discoverer enabled |
| discoverer.env | list | `[{"name":"MY_POD_NAMESPACE","valueFrom":{"fieldRef":{"fieldPath":"metadata.namespace"}}}]` | environment variables |
| discoverer.externalTrafficPolicy | string | `""` | external traffic policy (can be specified when service type is LoadBalancer or NodePort) : Cluster or Local |
| discoverer.hpa.enabled | bool | `false` | HPA enabled |
| discoverer.hpa.targetCPUUtilizationPercentage | int | `80` | HPA CPU utilization percentage |
| discoverer.image.pullPolicy | string | `"Always"` | image pull policy |
| discoverer.image.repository | string | `"vdaas/vald-discoverer-k8s"` | image repository |
| discoverer.image.tag | string | `""` | image tag (overrides defaults.image.tag) |
| discoverer.initContainers | list | `[]` | init containers |
| discoverer.kind | string | `"Deployment"` | deployment kind: Deployment or DaemonSet |
| discoverer.logging | object | `{}` | logging config (overrides defaults.logging) |
| discoverer.maxReplicas | int | `2` | maximum number of replicas. if HPA is disabled, this value will be ignored. |
| discoverer.maxUnavailable | string | `"50%"` | maximum number of unavailable replicas |
| discoverer.minReplicas | int | `1` | minimum number of replicas. if HPA is disabled, the replicas will be set to this value |
| discoverer.name | string | `"vald-discoverer"` | name of discoverer deployment |
| discoverer.nodeName | string | `""` | node name |
| discoverer.nodeSelector | object | `{}` | node selector |
| discoverer.observability | object | `{"jaeger":{"service_name":"vald-discoverer"},"stackdriver":{"profiler":{"service":"vald-discoverer"}}}` | observability config (overrides defaults.observability) |
| discoverer.podAnnotations | object | `{}` | pod annotations |
| discoverer.podPriority.enabled | bool | `true` | discoverer pod PriorityClass enabled |
| discoverer.podPriority.value | int | `1000000` | discoverer pod PriorityClass value |
| discoverer.progressDeadlineSeconds | int | `600` | progress deadline seconds |
| discoverer.resources | object | `{"limits":{"cpu":"600m","memory":"200Mi"},"requests":{"cpu":"200m","memory":"65Mi"}}` | compute resources |
| discoverer.revisionHistoryLimit | int | `2` | number of old history to retain to allow rollback |
| discoverer.rollingUpdate.maxSurge | string | `"25%"` | max surge of rolling update |
| discoverer.rollingUpdate.maxUnavailable | string | `"25%"` | max unavailable of rolling update |
| discoverer.server_config | object | `{"healths":{"liveness":{},"readiness":{}},"metrics":{"pprof":{},"prometheus":{}},"servers":{"grpc":{},"rest":{}}}` | server config (overrides defaults.server_config) |
| discoverer.service.annotations | object | `{}` | service annotations |
| discoverer.service.labels | object | `{}` | service labels |
| discoverer.serviceAccount.enabled | bool | `true` | creates service account |
| discoverer.serviceAccount.name | string | `"vald"` | name of service account |
| discoverer.serviceType | string | `"ClusterIP"` | service type: ClusterIP, LoadBalancer or NodePort |
| discoverer.terminationGracePeriodSeconds | int | `30` | duration in seconds pod needs to terminate gracefully |
| discoverer.time_zone | string | `""` | Time zone |
| discoverer.tolerations | list | `[]` | tolerations |
| discoverer.topologySpreadConstraints | list | `[]` | topology spread constraints of discoverer pods |
| discoverer.version | string | `"v0.0.0"` | version of discoverer config |
| discoverer.volumeMounts | list | `[]` | volume mounts |
| discoverer.volumes | list | `[]` | volumes |
| gateway.affinity.nodeAffinity.preferredDuringSchedulingIgnoredDuringExecution | list | `[]` | node affinity preferred scheduling terms |
| gateway.affinity.nodeAffinity.requiredDuringSchedulingIgnoredDuringExecution.nodeSelectorTerms | list | `[]` | node affinity required node selectors |
| gateway.affinity.podAffinity.preferredDuringSchedulingIgnoredDuringExecution | list | `[]` | pod affinity preferred scheduling terms |
| gateway.affinity.podAffinity.requiredDuringSchedulingIgnoredDuringExecution | list | `[]` | pod affinity required scheduling terms |
| gateway.affinity.podAntiAffinity.preferredDuringSchedulingIgnoredDuringExecution | list | `[{"podAffinityTerm":{"labelSelector":{"matchExpressions":[{"key":"app","operator":"In","values":["vald-gateway"]}]},"topologyKey":"kubernetes.io/hostname"},"weight":100}]` | pod anti-affinity preferred scheduling terms |
| gateway.affinity.podAntiAffinity.requiredDuringSchedulingIgnoredDuringExecution | list | `[]` | pod anti-affinity required scheduling terms |
| gateway.annotations | object | `{}` | deployment annotations |
| gateway.enabled | bool | `true` | gateway enabled |
| gateway.env | list | `[{"name":"MY_POD_NAMESPACE","valueFrom":{"fieldRef":{"fieldPath":"metadata.namespace"}}}]` | environment variables |
| gateway.externalTrafficPolicy | string | `""` | external traffic policy (can be specified when service type is LoadBalancer or NodePort) : Cluster or Local |
| gateway.filter.egress | list | `[""]` | egress filters |
| gateway.filter.ingress | list | `[""]` | ingress filters |
| gateway.gateway_config.agent_namespace | string | `"_MY_POD_NAMESPACE_"` | agent namespace |
| gateway.gateway_config.backup.client | object | `{}` | gRPC client for backup (overrides defaults.grpc.client) |
| gateway.gateway_config.discoverer.agent_client_options | object | `{}` | gRPC client for agents (overrides defaults.grpc.client) |
| gateway.gateway_config.discoverer.client | object | `{}` | gRPC client for discoverer (overrides defaults.grpc.client) |
| gateway.gateway_config.discoverer.duration | string | `"200ms"` | discoverer duration |
| gateway.gateway_config.index_replica | int | `5` | number of index replica |
| gateway.gateway_config.meta.cache_expiration | string | `"30m"` | meta cache expire duration |
| gateway.gateway_config.meta.client | object | `{}` | gRPC client for meta (overrides defaults.grpc.client) |
| gateway.gateway_config.meta.enable_cache | bool | `true` | meta cache enabled |
| gateway.gateway_config.meta.expired_cache_check_duration | string | `"3m"` | meta cache expired check duration |
| gateway.gateway_config.node_name | string | `""` | node name |
| gateway.hpa.enabled | bool | `true` | HPA enabled |
| gateway.hpa.targetCPUUtilizationPercentage | int | `80` | HPA CPU utilization percentage |
| gateway.image.pullPolicy | string | `"Always"` | image pull policy |
| gateway.image.repository | string | `"vdaas/vald-gateway"` | image repository |
| gateway.image.tag | string | `""` | image tag (overrides defaults.image.tag) |
| gateway.ingress.annotations | object | `{"nginx.ingress.kubernetes.io/grpc-backend":"true"}` | annotations for ingress |
| gateway.ingress.enabled | bool | `true` | gateway ingress enabled |
| gateway.ingress.host | string | `"vald.gateway.vdaas.org"` | ingress hostname |
| gateway.ingress.servicePort | string | `"grpc"` | service port to be exposed by ingress |
| gateway.initContainers | list | `[{"image":"busybox","name":"wait-for-manager-compressor","sleepDuration":2,"target":"compressor","type":"wait-for"},{"image":"busybox","name":"wait-for-meta","sleepDuration":2,"target":"meta","type":"wait-for"},{"image":"busybox","name":"wait-for-discoverer","sleepDuration":2,"target":"discoverer","type":"wait-for"},{"image":"busybox","name":"wait-for-agent","sleepDuration":2,"target":"agent","type":"wait-for"}]` | init containers |
| gateway.kind | string | `"Deployment"` | deployment kind: Deployment or DaemonSet |
| gateway.logging | object | `{}` | logging config (overrides defaults.logging) |
| gateway.maxReplicas | int | `9` | maximum number of replicas. if HPA is disabled, this value will be ignored. |
| gateway.maxUnavailable | string | `"50%"` | maximum number of unavailable replicas |
| gateway.minReplicas | int | `3` | minimum number of replicas. if HPA is disabled, the replicas will be set to this value |
| gateway.name | string | `"vald-gateway"` | name of gateway deployment |
| gateway.nodeName | string | `""` | node name |
| gateway.nodeSelector | object | `{}` | node selector |
| gateway.observability | object | `{"jaeger":{"service_name":"vald-gateway"},"stackdriver":{"profiler":{"service":"vald-gateway"}}}` | observability config (overrides defaults.observability) |
| gateway.podAnnotations | object | `{}` | pod annotations |
| gateway.podPriority.enabled | bool | `true` | gateway pod PriorityClass enabled |
| gateway.podPriority.value | int | `1000000` | gateway pod PriorityClass value |
| gateway.progressDeadlineSeconds | int | `600` | progress deadline seconds |
| gateway.resources | object | `{"limits":{"cpu":"2000m","memory":"700Mi"},"requests":{"cpu":"200m","memory":"150Mi"}}` | compute resources |
| gateway.revisionHistoryLimit | int | `2` | number of old history to retain to allow rollback |
| gateway.rollingUpdate.maxSurge | string | `"25%"` | max surge of rolling update |
| gateway.rollingUpdate.maxUnavailable | string | `"25%"` | max unavailable of rolling update |
| gateway.server_config | object | `{"healths":{"liveness":{},"readiness":{}},"metrics":{"pprof":{},"prometheus":{}},"servers":{"grpc":{},"rest":{}}}` | server config (overrides defaults.server_config) |
| gateway.service.annotations | object | `{}` | service annotations |
| gateway.service.labels | object | `{}` | service labels |
| gateway.serviceType | string | `"ClusterIP"` | service type: ClusterIP, LoadBalancer or NodePort |
| gateway.terminationGracePeriodSeconds | int | `30` | duration in seconds pod needs to terminate gracefully |
| gateway.time_zone | string | `""` | Time zone |
| gateway.tolerations | list | `[]` | tolerations |
| gateway.topologySpreadConstraints | list | `[]` | topology spread constraints of gateway pods |
| gateway.version | string | `"v0.0.0"` | version of gateway config |
| gateway.volumeMounts | list | `[]` | volume mounts |
| gateway.volumes | list | `[]` | volumes |
| indexManager.affinity.nodeAffinity.preferredDuringSchedulingIgnoredDuringExecution | list | `[]` | node affinity preferred scheduling terms |
| indexManager.affinity.nodeAffinity.requiredDuringSchedulingIgnoredDuringExecution.nodeSelectorTerms | list | `[]` | node affinity required node selectors |
| indexManager.affinity.podAffinity.preferredDuringSchedulingIgnoredDuringExecution | list | `[]` | pod affinity preferred scheduling terms |
| indexManager.affinity.podAffinity.requiredDuringSchedulingIgnoredDuringExecution | list | `[]` | pod affinity required scheduling terms |
| indexManager.affinity.podAntiAffinity.preferredDuringSchedulingIgnoredDuringExecution | list | `[]` | pod anti-affinity preferred scheduling terms |
| indexManager.affinity.podAntiAffinity.requiredDuringSchedulingIgnoredDuringExecution | list | `[]` | pod anti-affinity required scheduling terms |
| indexManager.annotations | object | `{}` | deployment annotations |
| indexManager.enabled | bool | `true` | index manager enabled |
| indexManager.env | list | `[{"name":"MY_POD_NAMESPACE","valueFrom":{"fieldRef":{"fieldPath":"metadata.namespace"}}}]` | environment variables |
| indexManager.externalTrafficPolicy | string | `""` | external traffic policy (can be specified when service type is LoadBalancer or NodePort) : Cluster or Local |
| indexManager.image.pullPolicy | string | `"Always"` | image pull policy |
| indexManager.image.repository | string | `"vdaas/vald-manager-index"` | image repository |
| indexManager.image.tag | string | `""` | image tag (overrides defaults.image.tag) |
| indexManager.indexer.agent_namespace | string | `"_MY_POD_NAMESPACE_"` | namespace of agent pods to manage |
| indexManager.indexer.auto_index_check_duration | string | `"1m"` | check duration of automatic indexing |
| indexManager.indexer.auto_index_duration_limit | string | `"30m"` | limit duration of automatic indexing |
| indexManager.indexer.auto_index_length | int | `100` | number of cache to trigger automatic indexing |
| indexManager.indexer.concurrency | int | `1` | concurrency |
| indexManager.indexer.creation_pool_size | int | `10000` | number of pool size of create index processing |
| indexManager.indexer.discoverer.agent_client_options | object | `{"dial_option":{"tcp":{"dialer":{"keep_alive":"15m"}}}}` | gRPC client for agents (overrides defaults.grpc.client) |
| indexManager.indexer.discoverer.client | object | `{}` | gRPC client for discoverer (overrides defaults.grpc.client) |
| indexManager.indexer.discoverer.duration | string | `"500ms"` | refresh duration to discover |
| indexManager.indexer.node_name | string | `""` | node name |
| indexManager.initContainers | list | `[{"image":"busybox","name":"wait-for-agent","sleepDuration":2,"target":"agent","type":"wait-for"},{"image":"busybox","name":"wait-for-discoverer","sleepDuration":2,"target":"discoverer","type":"wait-for"}]` | init containers |
| indexManager.kind | string | `"Deployment"` | deployment kind: Deployment or DaemonSet |
| indexManager.logging | object | `{}` | logging config (overrides defaults.logging) |
| indexManager.maxUnavailable | string | `"50%"` | maximum number of unavailable replicas |
| indexManager.name | string | `"vald-manager-index"` | name of index manager deployment |
| indexManager.nodeName | string | `""` | node name |
| indexManager.nodeSelector | object | `{}` | node selector |
| indexManager.observability | object | `{"jaeger":{"service_name":"vald-manager-index"},"stackdriver":{"profiler":{"service":"vald-manager-index"}}}` | observability config (overrides defaults.observability) |
| indexManager.podAnnotations | object | `{}` | pod annotations |
| indexManager.podPriority.enabled | bool | `true` | index manager pod PriorityClass enabled |
| indexManager.podPriority.value | int | `1000000` | index manager pod PriorityClass value |
| indexManager.progressDeadlineSeconds | int | `600` | progress deadline seconds |
| indexManager.replicas | int | `1` | number of replicas |
| indexManager.resources | object | `{"limits":{"cpu":"1000m","memory":"500Mi"},"requests":{"cpu":"200m","memory":"80Mi"}}` | compute resources |
| indexManager.revisionHistoryLimit | int | `2` | number of old history to retain to allow rollback |
| indexManager.rollingUpdate.maxSurge | string | `"25%"` | max surge of rolling update |
| indexManager.rollingUpdate.maxUnavailable | string | `"25%"` | max unavailable of rolling update |
| indexManager.server_config | object | `{"healths":{"liveness":{},"readiness":{}},"metrics":{"pprof":{},"prometheus":{}},"servers":{"grpc":{},"rest":{}}}` | server config (overrides defaults.server_config) |
| indexManager.service.annotations | object | `{}` | service annotations |
| indexManager.service.labels | object | `{}` | service labels |
| indexManager.serviceType | string | `"ClusterIP"` | service type: ClusterIP, LoadBalancer or NodePort |
| indexManager.terminationGracePeriodSeconds | int | `30` | duration in seconds pod needs to terminate gracefully |
| indexManager.time_zone | string | `""` | Time zone |
| indexManager.tolerations | list | `[]` | tolerations |
| indexManager.topologySpreadConstraints | list | `[]` | topology spread constraints of index manager pods |
| indexManager.version | string | `"v0.0.0"` | version of index manager config |
| indexManager.volumeMounts | list | `[]` | volume mounts |
| indexManager.volumes | list | `[]` | volumes |
| initializer.cassandra.configmap.backup.enabled | bool | `true` | backup table enabled |
| initializer.cassandra.configmap.backup.name | string | `"backup_vector"` | name of backup table |
| initializer.cassandra.configmap.enabled | bool | `false` | cassandra schema configmap will be created |
| initializer.cassandra.configmap.filename | string | `"init.cql"` | cassandra schema filename |
| initializer.cassandra.configmap.keyspace | string | `"vald"` | cassandra keyspace |
| initializer.cassandra.configmap.meta.enabled | bool | `true` | meta table enabled |
| initializer.cassandra.configmap.meta.name.kv | string | `"kv"` | name of KV table |
| initializer.cassandra.configmap.meta.name.vk | string | `"vk"` | name of VK table |
| initializer.cassandra.configmap.name | string | `"cassandra-initdb"` | cassandra schema configmap name |
| initializer.cassandra.configmap.replication_class | string | `"SimpleStrategy"` | cassandra replication class |
| initializer.cassandra.configmap.replication_factor | int | `3` | cassandra replication factor |
| initializer.cassandra.enabled | bool | `false` | cassandra initializer job enabled |
| initializer.cassandra.env | list | `[{"name":"CASSANDRA_HOST","value":"cassandra.default.svc.cluster.local"},{"name":"CASSANDRA_USER","value":"root"},{"name":"CASSANDRA_PASSWORD","valueFrom":{"secretKeyRef":{"key":"password","name":"cassandra-secret"}}}]` | environment variables |
| initializer.cassandra.image.pullPolicy | string | `"Always"` | image pull policy |
| initializer.cassandra.image.repository | string | `"cassandra"` | image repository |
| initializer.cassandra.image.tag | string | `"latest"` | image tag |
| initializer.cassandra.name | string | `"cassandra-init"` | cassandra initializer job name |
| initializer.cassandra.restartPolicy | string | `"Never"` | restart policy |
| initializer.cassandra.secret.data | object | `{"password":"cGFzc3dvcmQ="}` | cassandra secret data |
| initializer.cassandra.secret.enabled | bool | `false` | cassandra secret will be created |
| initializer.cassandra.secret.name | string | `"cassandra-secret"` | cassandra secret name |
| initializer.mysql.configmap.enabled | bool | `false` | mysql schema configmap will be created |
| initializer.mysql.configmap.filename | string | `"ddl.sql"` | mysql schema filename |
| initializer.mysql.configmap.name | string | `"mysql-config"` | mysql schema configmap name |
| initializer.mysql.configmap.schema | string | `"vald"` | mysql schema name |
| initializer.mysql.enabled | bool | `false` | mysql initializer job enabled |
| initializer.mysql.env | list | `[{"name":"MYSQL_HOST","value":"mysql.default.svc.cluster.local"},{"name":"MYSQL_USER","value":"root"},{"name":"MYSQL_PASSWORD","valueFrom":{"secretKeyRef":{"key":"password","name":"mysql-secret"}}}]` | environment variables |
| initializer.mysql.image.pullPolicy | string | `"Always"` | image pull policy |
| initializer.mysql.image.repository | string | `"mysql"` | image repository |
| initializer.mysql.image.tag | string | `"latest"` | image tag |
| initializer.mysql.name | string | `"mysql-init"` | mysql initializer job name |
| initializer.mysql.restartPolicy | string | `"Never"` | restart policy |
| initializer.mysql.secret.data | object | `{"password":"cGFzc3dvcmQ="}` | mysql secret data |
| initializer.mysql.secret.enabled | bool | `false` | mysql secret will be created |
| initializer.mysql.secret.name | string | `"mysql-secret"` | mysql secret name |
| initializer.redis.enabled | bool | `false` | redis initializer job enabled |
| initializer.redis.env | list | `[{"name":"REDIS_HOST","value":"redis.default.svc.cluster.local"},{"name":"REDIS_PASSWORD","valueFrom":{"secretKeyRef":{"key":"password","name":"redis-secret"}}}]` | environment variables |
| initializer.redis.image.pullPolicy | string | `"Always"` | image pull policy |
| initializer.redis.image.repository | string | `"redis"` | image repository |
| initializer.redis.image.tag | string | `"latest"` | image tag |
| initializer.redis.name | string | `"redis-init"` | redis initializer job name |
| initializer.redis.restartPolicy | string | `"Never"` | restart policy |
| initializer.redis.secret.data | object | `{"password":"cGFzc3dvcmQ="}` | redis secret data |
| initializer.redis.secret.enabled | bool | `false` | redis secret will be created |
| initializer.redis.secret.name | string | `"redis-secret"` | redis secret name |
| meta.affinity.nodeAffinity.preferredDuringSchedulingIgnoredDuringExecution | list | `[]` | node affinity preferred scheduling terms |
| meta.affinity.nodeAffinity.requiredDuringSchedulingIgnoredDuringExecution.nodeSelectorTerms | list | `[]` | node affinity required node selectors |
| meta.affinity.podAffinity.preferredDuringSchedulingIgnoredDuringExecution | list | `[]` | pod affinity preferred scheduling terms |
| meta.affinity.podAffinity.requiredDuringSchedulingIgnoredDuringExecution | list | `[]` | pod affinity required scheduling terms |
| meta.affinity.podAntiAffinity.preferredDuringSchedulingIgnoredDuringExecution | list | `[]` | pod anti-affinity preferred scheduling terms |
| meta.affinity.podAntiAffinity.requiredDuringSchedulingIgnoredDuringExecution | list | `[]` | pod anti-affinity required scheduling terms |
| meta.annotations | object | `{}` | deployment annotations |
| meta.cassandra.config.connect_timeout | string | `"3s"` | connect timeout |
| meta.cassandra.config.consistency | string | `"quorum"` | consistency type |
| meta.cassandra.config.cql_version | string | `"3.0.0"` | cassandra CQL version |
| meta.cassandra.config.default_idempotence | bool | `false` | default idempotence enabled |
| meta.cassandra.config.default_timestamp | bool | `true` | default timestamp enabled |
| meta.cassandra.config.disable_initial_host_lookup | bool | `false` | initial host lookup disabled |
| meta.cassandra.config.disable_node_status_events | bool | `false` | node status events disabled |
| meta.cassandra.config.disable_skip_metadata | bool | `false` | skip metadata disabled |
| meta.cassandra.config.disable_topology_events | bool | `false` | topology events disabled |
| meta.cassandra.config.enable_host_verification | bool | `false` | host verification enabled |
| meta.cassandra.config.host_filter.data_center | string | `""` | name of data center of filtering target |
| meta.cassandra.config.host_filter.enabled | bool | `false` | enables host filtering |
| meta.cassandra.config.host_filter.white_list | list | `[]` | list of white_list which allows each connection |
| meta.cassandra.config.hosts | list | `["cassandra-0.cassandra.default.svc.cluster.local","cassandra-1.cassandra.default.svc.cluster.local","cassandra-2.cassandra.default.svc.cluster.local"]` | cassandra hosts |
| meta.cassandra.config.ignore_peer_addr | bool | `false` | ignore peer addresses |
| meta.cassandra.config.keyspace | string | `"vald"` | cassandra keyspace |
| meta.cassandra.config.max_prepared_stmts | int | `1000` | maximum number of prepared statements |
| meta.cassandra.config.max_routing_key_info | int | `1000` | maximum number of routing key info |
| meta.cassandra.config.max_wait_schema_agreement | string | `"1m"` | maximum duration to wait for schema agreement |
| meta.cassandra.config.vector_backup_table | string | `"backup_vector"` | table name of backup |
| meta.cassandra.config.num_conns | int | `2` | number of connections per hosts |
| meta.cassandra.config.page_size | int | `5000` | page size |
| meta.cassandra.config.password | string | `"_CASSANDRA_PASSWORD_"` | cassandra password |
| meta.cassandra.config.pool_config.data_center | string | `""` | name of data center |
| meta.cassandra.config.pool_config.dc_aware_routing | bool | `false` | data center aware routine enabled |
| meta.cassandra.config.pool_config.non_local_replicas_fallback | bool | `false` | non-local replica fallback enabled |
| meta.cassandra.config.pool_config.shuffle_replicas | bool | `false` | shuffle replica enabled |
| meta.cassandra.config.pool_config.token_aware_host_policy | bool | `false` | token aware host policy enabled |
| meta.cassandra.config.port | int | `9042` | cassandra port |
| meta.cassandra.config.proto_version | int | `0` | cassandra proto version |
| meta.cassandra.config.reconnect_interval | string | `"100ms"` | interval of reconnection |
| meta.cassandra.config.reconnection_policy.initial_interval | string | `"100ms"` | initial interval to reconnect |
| meta.cassandra.config.reconnection_policy.max_retries | int | `3` | maximum number of retries to reconnect |
| meta.cassandra.config.retry_policy.max_duration | string | `"1s"` | maximum duration to retry |
| meta.cassandra.config.retry_policy.min_duration | string | `"10ms"` | minimum duration to retry |
| meta.cassandra.config.retry_policy.num_retries | int | `3` | number of retries |
| meta.cassandra.config.serial_consistency | string | `"localserial"` | read consistency type |
| meta.cassandra.config.socket_keepalive | string | `"0s"` | socket keep alive time |
| meta.cassandra.config.tcp.dialer.dual_stack_enabled | bool | `false` | TCP dialer dual stack enabled |
| meta.cassandra.config.tcp.dialer.keep_alive | string | `"10m"` | TCP dialer keep alive |
| meta.cassandra.config.tcp.dialer.timeout | string | `"30s"` | TCP dialer timeout |
| meta.cassandra.config.tcp.dns.cache_enabled | bool | `true` | TCP DNS cache enabled |
| meta.cassandra.config.tcp.dns.cache_expiration | string | `"24h"` | TCP DNS cache expiration |
| meta.cassandra.config.tcp.dns.refresh_duration | string | `"5m"` | TCP DNS cache refresh duration |
| meta.cassandra.config.timeout | string | `"600ms"` | timeout |
| meta.cassandra.config.tls.ca | string | `"/path/to/ca"` | path to TLS ca |
| meta.cassandra.config.tls.cert | string | `"/path/to/cert"` | path to TLS cert |
| meta.cassandra.config.tls.enabled | bool | `false` | TLS enabled |
| meta.cassandra.config.tls.key | string | `"/path/to/key"` | path to TLS key |
| meta.cassandra.config.username | string | `"root"` | cassandra username |
| meta.cassandra.config.write_coalesce_wait_time | string | `"200µs"` | write coalesce wait time |
| meta.cassandra.enabled | bool | `false` | cassandra config enabled |
| meta.enabled | bool | `true` | meta enabled |
| meta.env | list | `[{"name":"REDIS_PASSWORD","valueFrom":{"secretKeyRef":{"key":"password","name":"redis-secret"}}}]` | environment variables |
| meta.externalTrafficPolicy | string | `""` | external traffic policy (can be specified when service type is LoadBalancer or NodePort) : Cluster or Local |
| meta.hpa.enabled | bool | `true` | HPA enabled |
| meta.hpa.targetCPUUtilizationPercentage | int | `80` | HPA CPU utilization percentage |
| meta.image.pullPolicy | string | `"Always"` | image pull policy |
| meta.image.repository | string | `"vdaas/vald-meta-redis"` | image repository |
| meta.image.tag | string | `""` | image tag (overrides defaults.image.tag) |
| meta.initContainers | list | `[{"env":[{"name":"REDIS_PASSWORD","valueFrom":{"secretKeyRef":{"key":"password","name":"redis-secret"}}}],"image":"redis:latest","name":"wait-for-redis","redis":{"hosts":["redis.default.svc.cluster.local"],"options":["-a ${REDIS_PASSWORD}"]},"sleepDuration":2,"type":"wait-for-redis"}]` | init containers |
| meta.kind | string | `"Deployment"` | deployment kind: Deployment or DaemonSet |
| meta.logging | object | `{}` | logging config (overrides defaults.logging) |
| meta.maxReplicas | int | `10` | maximum number of replicas. if HPA is disabled, this value will be ignored. |
| meta.maxUnavailable | string | `"50%"` | maximum number of unavailable replicas |
| meta.minReplicas | int | `2` | minimum number of replicas. if HPA is disabled, the replicas will be set to this value. |
| meta.name | string | `"vald-meta"` | name of meta deployment |
| meta.nodeName | string | `""` | node name |
| meta.nodeSelector | object | `{}` | node selector |
| meta.observability | object | `{"jaeger":{"service_name":"vald-meta"},"stackdriver":{"profiler":{"service":"vald-meta"}}}` | observability config (overrides defaults.observability) |
| meta.podAnnotations | object | `{}` | pod annotations |
| meta.podPriority.enabled | bool | `true` | meta pod PriorityClass enabled |
| meta.podPriority.value | int | `1000000` | meta pod PriorityClass value |
| meta.progressDeadlineSeconds | int | `600` | progress deadline seconds |
| meta.redis.config.addrs | list | `["redis.default.svc.cluster.local:6379"]` | redis hosts and ports |
| meta.redis.config.db | int | `0` | database to be selected |
| meta.redis.config.dial_timeout | string | `"5s"` | dial timeout |
| meta.redis.config.idle_check_frequency | string | `"1m"` | idle check frequency |
| meta.redis.config.idle_timeout | string | `"5m"` | idle timeout |
| meta.redis.config.key_pref | string | `""` | key prefix |
| meta.redis.config.kv_prefix | string | `""` | KV prefix |
| meta.redis.config.max_conn_age | string | `"0s"` | max connection age |
| meta.redis.config.max_redirects | int | `3` | max redirects |
| meta.redis.config.max_retries | int | `0` | max retries |
| meta.redis.config.max_retry_backoff | string | `"512ms"` | max retry backoff |
| meta.redis.config.min_idle_conns | int | `0` | min idle connections |
| meta.redis.config.min_retry_backoff | string | `"8ms"` | min retry backoff |
| meta.redis.config.password | string | `"_REDIS_PASSWORD_"` | redis password |
| meta.redis.config.pool_size | int | `10` | pool size |
| meta.redis.config.pool_timeout | string | `"4s"` | pool timeout |
| meta.redis.config.prefix_delimiter | string | `""` | prefix delimiter |
| meta.redis.config.read_only | bool | `false` | read only enabled |
| meta.redis.config.read_timeout | string | `"3s"` | read timeout |
| meta.redis.config.route_by_latency | bool | `false` | latency based routing enabled |
| meta.redis.config.route_randomly | bool | `true` | random routing enabled |
| meta.redis.config.tcp.dialer.dual_stack_enabled | bool | `false` | TCP dialer dual stack enabled |
| meta.redis.config.tcp.dialer.keep_alive | string | `"5m"` | TCP dialer keep alive |
| meta.redis.config.tcp.dialer.timeout | string | `"5s"` | TCP dialer timeout |
| meta.redis.config.tcp.dns.cache_enabled | bool | `true` | TCP DNS cache enabled |
| meta.redis.config.tcp.dns.cache_expiration | string | `"24h"` | TCP DNS cache expiration |
| meta.redis.config.tcp.dns.refresh_duration | string | `"1h"` | TCP DNS cache refresh duration |
| meta.redis.config.tcp.tls.ca | string | `"/path/to/ca"` | path to TCP TLS ca |
| meta.redis.config.tcp.tls.cert | string | `"/path/to/cert"` | path to TCP TLS cert |
| meta.redis.config.tcp.tls.enabled | bool | `false` | TCP TLS enabled |
| meta.redis.config.tcp.tls.key | string | `"/path/to/key"` | path to TCP TLS key |
| meta.redis.config.tls.ca | string | `"/path/to/ca"` | path to TLS ca |
| meta.redis.config.tls.cert | string | `"/path/to/cert"` | path to TLS cert |
| meta.redis.config.tls.enabled | bool | `false` | TLS enabled |
| meta.redis.config.tls.key | string | `"/path/to/key"` | path to TLS key |
| meta.redis.config.vk_prefix | string | `""` | VK prefix |
| meta.redis.config.write_timeout | string | `"3s"` | write timeout |
| meta.redis.enabled | bool | `true` | redis config enabled |
| meta.resources | object | `{"limits":{"cpu":"300m","memory":"100Mi"},"requests":{"cpu":"100m","memory":"40Mi"}}` | compute resources |
| meta.revisionHistoryLimit | int | `2` | number of old history to retain to allow rollback |
| meta.rollingUpdate.maxSurge | string | `"25%"` | max surge of rolling update |
| meta.rollingUpdate.maxUnavailable | string | `"25%"` | max unavailable of rolling update |
| meta.server_config | object | `{"healths":{"liveness":{},"readiness":{}},"metrics":{"pprof":{},"prometheus":{}},"servers":{"grpc":{},"rest":{}}}` | server config (overrides defaults.server_config) |
| meta.service.annotations | object | `{}` | service annotations |
| meta.service.labels | object | `{}` | service labels |
| meta.serviceType | string | `"ClusterIP"` | service type: ClusterIP, LoadBalancer or NodePort |
| meta.terminationGracePeriodSeconds | int | `30` | duration in seconds pod needs to terminate gracefully |
| meta.time_zone | string | `""` | Time zone |
| meta.tolerations | list | `[]` | tolerations |
| meta.topologySpreadConstraints | list | `[]` | topology spread constraints of meta pods |
| meta.version | string | `"v0.0.0"` | version of meta config |
| meta.volumeMounts | list | `[]` | volume mounts |
| meta.volumes | list | `[]` | volumes |

Miscellaneous
---

### Standalone Vald agent NGT deployment

Each component can be disabled by setting the value `false` to the `[component].enabled` field.
This is useful for deploying only Vald agent NGT pods.

There is an example yaml [values-agent-ngt-standalone.yaml][agent-ngt-standalone-yaml] to deploy standalone agent NGT.
Please run the following command to install the chart with this values yaml,

    $ helm repo add vald https://vald.vdaas.org/charts
    $ helm install --values values-agent-ngt-standalone.yaml vald-agent-ngt vald/vald

If you'd like to access the agents from out of the Kubernetes cluster, it is recommended to create an [Ingress][kubernetes-ingress].

[agent-ngt-standalone-yaml]: ./values-agent-ngt-standalone.yaml
[kubernetes-ingress]: https://kubernetes.io/docs/concepts/services-networking/ingress/
