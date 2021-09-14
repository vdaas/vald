Vald
===

This is a Helm chart to install Vald components.

Current chart version is `v1.2.2`

Table of Contents
---

- [Install](#install)
- [Configuration](#configuration)
    - [Overview](#overview)
    - [Parameters](#parameters)

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
    - configurations of gateways
        - `filter`
        - `lb`
- `agent`
    - configurations of vald-agent
- `discoverer`
    - configurations of vald-discoverer
- `manager`
    - configurations of manager components
        - `index`

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
| agent.ngt.auto_create_index_pool_size | int | `10000` | batch process pool size of automatic create index operation |
| agent.ngt.auto_index_check_duration | string | `"30m"` | check duration of automatic indexing |
| agent.ngt.auto_index_duration_limit | string | `"24h"` | limit duration of automatic indexing |
| agent.ngt.auto_index_length | int | `100` | number of cache to trigger automatic indexing |
| agent.ngt.auto_save_index_duration | string | `"35m"` | duration of automatic save index |
| agent.ngt.bulk_insert_chunk_size | int | `10` | bulk insert chunk size |
| agent.ngt.creation_edge_size | int | `20` | creation edge size |
| agent.ngt.default_epsilon | float | `0.1` | default epsilon used for search |
| agent.ngt.default_pool_size | int | `10000` | default create index batch pool size |
| agent.ngt.default_radius | float | `-1` | default radius used for search |
| agent.ngt.dimension | int | `4096` | vector dimension |
| agent.ngt.distance_type | string | `"l2"` | distance type. it should be `l1`, `l2`, `angle`, `hamming`, `cosine`, `normalizedangle`, `normalizedcosine` or `jaccard`. for further details about NGT libraries supported distance is https://github.com/yahoojapan/NGT/wiki/Command-Quick-Reference and vald agent's supported NGT distance type is https://pkg.go.dev/github.com/vdaas/vald/internal/core/algorithm/ngt#pkg-constants |
| agent.ngt.enable_in_memory_mode | bool | `true` | in-memory mode enabled |
| agent.ngt.enable_proactive_gc | bool | `true` | enable proactive GC call for reducing heap memory allocation |
| agent.ngt.index_path | string | `""` | path to index data |
| agent.ngt.initial_delay_max_duration | string | `"3m"` | maximum duration for initial delay |
| agent.ngt.kvsdb.concurrency | int | `20` | kvsdb processing concurrency |
| agent.ngt.load_index_timeout_factor | string | `"1ms"` | a factor of load index timeout. timeout duration will be calculated by (index count to be loaded) * (factor). |
| agent.ngt.max_load_index_timeout | string | `"10m"` | maximum duration of load index timeout |
| agent.ngt.min_load_index_timeout | string | `"3m"` | minimum duration of load index timeout |
| agent.ngt.object_type | string | `"float"` | object type. it should be `float` or `uint8`. for further details: https://github.com/yahoojapan/NGT/wiki/Command-Quick-Reference |
| agent.ngt.search_edge_size | int | `10` | search edge size |
| agent.ngt.vqueue.delete_buffer_pool_size | int | `5000` | delete slice pool buffer size |
| agent.ngt.vqueue.insert_buffer_pool_size | int | `10000` | insert slice pool buffer size |
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
| agent.podSecurityContext | object | `{"fsGroup":65532,"fsGroupChangePolicy":"OnRootMismatch","runAsGroup":65532,"runAsNonRoot":true,"runAsUser":65532}` | security context for pod |
| agent.progressDeadlineSeconds | int | `600` | progress deadline seconds |
| agent.resources | object | `{"requests":{"cpu":"300m","memory":"4Gi"}}` | compute resources. recommended setting of memory requests = cluster memory * 0.4 / number of agent pods |
| agent.revisionHistoryLimit | int | `2` | number of old history to retain to allow rollback |
| agent.rollingUpdate.maxSurge | string | `"25%"` | max surge of rolling update |
| agent.rollingUpdate.maxUnavailable | string | `"25%"` | max unavailable of rolling update |
| agent.rollingUpdate.partition | int | `0` | StatefulSet partition |
| agent.securityContext | object | `{"allowPrivilegeEscalation":false,"capabilities":{"drop":["ALL"]},"privileged":false,"readOnlyRootFilesystem":false,"runAsGroup":65532,"runAsNonRoot":true,"runAsUser":65532}` | security context for container |
| agent.server_config | object | `{"healths":{"liveness":{"enabled":false},"readiness":{}},"metrics":{"pprof":{},"prometheus":{}},"servers":{"grpc":{},"rest":{}}}` | server config (overrides defaults.server_config) |
| agent.service.annotations | object | `{}` | service annotations |
| agent.service.labels | object | `{}` | service labels |
| agent.serviceType | string | `"ClusterIP"` | service type: ClusterIP, LoadBalancer or NodePort |
| agent.sidecar.config.auto_backup_duration | string | `"24h"` | auto backup duration |
| agent.sidecar.config.auto_backup_enabled | bool | `true` | auto backup triggered by timer is enabled |
| agent.sidecar.config.blob_storage.bucket | string | `""` | bucket name |
| agent.sidecar.config.blob_storage.cloud_storage.client.credentials_file_path | string | `""` | credentials file path |
| agent.sidecar.config.blob_storage.cloud_storage.client.credentials_json | string | `""` | credentials json |
| agent.sidecar.config.blob_storage.cloud_storage.url | string | `""` | cloud storage url |
| agent.sidecar.config.blob_storage.cloud_storage.write_buffer_size | int | `0` | bytes of the chunks for upload |
| agent.sidecar.config.blob_storage.cloud_storage.write_cache_control | string | `""` | Cache-Control of HTTP Header |
| agent.sidecar.config.blob_storage.cloud_storage.write_content_disposition | string | `""` | Content-Disposition of HTTP Header |
| agent.sidecar.config.blob_storage.cloud_storage.write_content_encoding | string | `""` | the encoding of the blob's content |
| agent.sidecar.config.blob_storage.cloud_storage.write_content_language | string | `""` | the language of blob's content |
| agent.sidecar.config.blob_storage.cloud_storage.write_content_type | string | `""` | MIME type of the blob |
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
| agent.sidecar.config.client.net.dialer.dual_stack_enabled | bool | `false` | HTTP client TCP dialer dual stack enabled |
| agent.sidecar.config.client.net.dialer.keep_alive | string | `"5m"` | HTTP client TCP dialer keep alive |
| agent.sidecar.config.client.net.dialer.timeout | string | `"5s"` | HTTP client TCP dialer connect timeout |
| agent.sidecar.config.client.net.dns.cache_enabled | bool | `true` | HTTP client TCP DNS cache enabled |
| agent.sidecar.config.client.net.dns.cache_expiration | string | `"24h"` |  |
| agent.sidecar.config.client.net.dns.refresh_duration | string | `"1h"` | HTTP client TCP DNS cache expiration |
| agent.sidecar.config.client.net.socket_option.ip_recover_destination_addr | bool | `false` | server listen socket option for ip_recover_destination_addr functionality |
| agent.sidecar.config.client.net.socket_option.ip_transparent | bool | `false` | server listen socket option for ip_transparent functionality |
| agent.sidecar.config.client.net.socket_option.reuse_addr | bool | `true` | server listen socket option for reuse_addr functionality |
| agent.sidecar.config.client.net.socket_option.reuse_port | bool | `true` | server listen socket option for reuse_port functionality |
| agent.sidecar.config.client.net.socket_option.tcp_cork | bool | `false` | server listen socket option for tcp_cork functionality |
| agent.sidecar.config.client.net.socket_option.tcp_defer_accept | bool | `true` | server listen socket option for tcp_defer_accept functionality |
| agent.sidecar.config.client.net.socket_option.tcp_fast_open | bool | `true` | server listen socket option for tcp_fast_open functionality |
| agent.sidecar.config.client.net.socket_option.tcp_no_delay | bool | `true` | server listen socket option for tcp_no_delay functionality |
| agent.sidecar.config.client.net.socket_option.tcp_quick_ack | bool | `true` | server listen socket option for tcp_quick_ack functionality |
| agent.sidecar.config.client.net.tls.ca | string | `"/path/to/ca"` | TLS ca path |
| agent.sidecar.config.client.net.tls.cert | string | `"/path/to/cert"` | TLS cert path |
| agent.sidecar.config.client.net.tls.enabled | bool | `false` | TLS enabled |
| agent.sidecar.config.client.net.tls.insecure_skip_verify | bool | `false` | enable/disable skip SSL certificate verification |
| agent.sidecar.config.client.net.tls.key | string | `"/path/to/key"` | TLS key path |
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
| defaults.grpc.client.dial_option.backoff_base_delay | string | `"1s"` | gRPC client dial option base backoff delay |
| defaults.grpc.client.dial_option.backoff_jitter | float | `0.2` | gRPC client dial option base backoff delay |
| defaults.grpc.client.dial_option.backoff_max_delay | string | `"120s"` | gRPC client dial option max backoff delay |
| defaults.grpc.client.dial_option.backoff_multiplier | float | `1.6` | gRPC client dial option base backoff delay |
| defaults.grpc.client.dial_option.enable_backoff | bool | `false` | gRPC client dial option backoff enabled |
| defaults.grpc.client.dial_option.initial_connection_window_size | int | `0` | gRPC client dial option initial connection window size |
| defaults.grpc.client.dial_option.initial_window_size | int | `0` | gRPC client dial option initial window size |
| defaults.grpc.client.dial_option.insecure | bool | `true` | gRPC client dial option insecure enabled |
| defaults.grpc.client.dial_option.keep_alive.permit_without_stream | bool | `false` | gRPC client keep alive permit without stream |
| defaults.grpc.client.dial_option.keep_alive.time | string | `""` | gRPC client keep alive time |
| defaults.grpc.client.dial_option.keep_alive.timeout | string | `""` | gRPC client keep alive timeout |
| defaults.grpc.client.dial_option.max_msg_size | int | `0` | gRPC client dial option max message size |
| defaults.grpc.client.dial_option.min_connection_timeout | string | `"20s"` | gRPC client dial option minimum connection timeout |
| defaults.grpc.client.dial_option.net.dialer.dual_stack_enabled | bool | `true` | gRPC client TCP dialer dual stack enabled |
| defaults.grpc.client.dial_option.net.dialer.keep_alive | string | `""` | gRPC client TCP dialer keep alive |
| defaults.grpc.client.dial_option.net.dialer.timeout | string | `""` | gRPC client TCP dialer timeout |
| defaults.grpc.client.dial_option.net.dns.cache_enabled | bool | `true` | gRPC client TCP DNS cache enabled |
| defaults.grpc.client.dial_option.net.dns.cache_expiration | string | `"1h"` | gRPC client TCP DNS cache expiration |
| defaults.grpc.client.dial_option.net.dns.refresh_duration | string | `"30m"` | gRPC client TCP DNS cache refresh duration |
| defaults.grpc.client.dial_option.net.socket_option.ip_recover_destination_addr | bool | `false` | server listen socket option for ip_recover_destination_addr functionality |
| defaults.grpc.client.dial_option.net.socket_option.ip_transparent | bool | `false` | server listen socket option for ip_transparent functionality |
| defaults.grpc.client.dial_option.net.socket_option.reuse_addr | bool | `true` | server listen socket option for reuse_addr functionality |
| defaults.grpc.client.dial_option.net.socket_option.reuse_port | bool | `true` | server listen socket option for reuse_port functionality |
| defaults.grpc.client.dial_option.net.socket_option.tcp_cork | bool | `false` | server listen socket option for tcp_cork functionality |
| defaults.grpc.client.dial_option.net.socket_option.tcp_defer_accept | bool | `true` | server listen socket option for tcp_defer_accept functionality |
| defaults.grpc.client.dial_option.net.socket_option.tcp_fast_open | bool | `true` | server listen socket option for tcp_fast_open functionality |
| defaults.grpc.client.dial_option.net.socket_option.tcp_no_delay | bool | `true` | server listen socket option for tcp_no_delay functionality |
| defaults.grpc.client.dial_option.net.socket_option.tcp_quick_ack | bool | `true` | server listen socket option for tcp_quick_ack functionality |
| defaults.grpc.client.dial_option.net.tls.ca | string | `"/path/to/ca"` | TLS ca path |
| defaults.grpc.client.dial_option.net.tls.cert | string | `"/path/to/cert"` | TLS cert path |
| defaults.grpc.client.dial_option.net.tls.enabled | bool | `false` | TLS enabled |
| defaults.grpc.client.dial_option.net.tls.insecure_skip_verify | bool | `false` | enable/disable skip SSL certificate verification |
| defaults.grpc.client.dial_option.net.tls.key | string | `"/path/to/key"` | TLS key path |
| defaults.grpc.client.dial_option.read_buffer_size | int | `0` | gRPC client dial option read buffer size |
| defaults.grpc.client.dial_option.timeout | string | `""` | gRPC client dial option timeout |
| defaults.grpc.client.dial_option.write_buffer_size | int | `0` | gRPC client dial option write buffer size |
| defaults.grpc.client.health_check_duration | string | `"1s"` | gRPC client health check duration |
| defaults.grpc.client.tls.ca | string | `"/path/to/ca"` | TLS ca path |
| defaults.grpc.client.tls.cert | string | `"/path/to/cert"` | TLS cert path |
| defaults.grpc.client.tls.enabled | bool | `false` | TLS enabled |
| defaults.grpc.client.tls.insecure_skip_verify | bool | `false` | enable/disable skip SSL certificate verification |
| defaults.grpc.client.tls.key | string | `"/path/to/key"` | TLS key path |
| defaults.image.tag | string | `"v1.2.2"` | docker image tag |
| defaults.ingress.usev1beta1 | bool | `false` | use networking.k8s.io/v1beta1 instead of v1 for ingresses. This option will be removed once k8s 1.22 is released. |
| defaults.logging.format | string | `"raw"` | logging format. logging format must be `raw` or `json` |
| defaults.logging.level | string | `"debug"` | logging level. logging level must be `debug`, `info`, `warn`, `error` or `fatal`. |
| defaults.logging.logger | string | `"glg"` | logger name. currently logger must be `glg` or `zap`. |
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
| defaults.server_config.healths.liveness.server.network | string | `"tcp"` | mysql network |
| defaults.server_config.healths.liveness.server.probe_wait_time | string | `"3s"` | liveness server probe wait time |
| defaults.server_config.healths.liveness.server.socket_option.ip_recover_destination_addr | bool | `false` | server listen socket option for ip_recover_destination_addr functionality |
| defaults.server_config.healths.liveness.server.socket_option.ip_transparent | bool | `false` | server listen socket option for ip_transparent functionality |
| defaults.server_config.healths.liveness.server.socket_option.reuse_addr | bool | `true` | server listen socket option for reuse_addr functionality |
| defaults.server_config.healths.liveness.server.socket_option.reuse_port | bool | `true` | server listen socket option for reuse_port functionality |
| defaults.server_config.healths.liveness.server.socket_option.tcp_cork | bool | `false` | server listen socket option for tcp_cork functionality |
| defaults.server_config.healths.liveness.server.socket_option.tcp_defer_accept | bool | `true` | server listen socket option for tcp_defer_accept functionality |
| defaults.server_config.healths.liveness.server.socket_option.tcp_fast_open | bool | `true` | server listen socket option for tcp_fast_open functionality |
| defaults.server_config.healths.liveness.server.socket_option.tcp_no_delay | bool | `true` | server listen socket option for tcp_no_delay functionality |
| defaults.server_config.healths.liveness.server.socket_option.tcp_quick_ack | bool | `true` | server listen socket option for tcp_quick_ack functionality |
| defaults.server_config.healths.liveness.server.socket_path | string | `""` | mysql socket_path |
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
| defaults.server_config.healths.readiness.server.network | string | `"tcp"` | mysql network |
| defaults.server_config.healths.readiness.server.probe_wait_time | string | `"3s"` | readiness server probe wait time |
| defaults.server_config.healths.readiness.server.socket_option.ip_recover_destination_addr | bool | `false` | server listen socket option for ip_recover_destination_addr functionality |
| defaults.server_config.healths.readiness.server.socket_option.ip_transparent | bool | `false` | server listen socket option for ip_transparent functionality |
| defaults.server_config.healths.readiness.server.socket_option.reuse_addr | bool | `true` | server listen socket option for reuse_addr functionality |
| defaults.server_config.healths.readiness.server.socket_option.reuse_port | bool | `true` | server listen socket option for reuse_port functionality |
| defaults.server_config.healths.readiness.server.socket_option.tcp_cork | bool | `false` | server listen socket option for tcp_cork functionality |
| defaults.server_config.healths.readiness.server.socket_option.tcp_defer_accept | bool | `true` | server listen socket option for tcp_defer_accept functionality |
| defaults.server_config.healths.readiness.server.socket_option.tcp_fast_open | bool | `true` |  |
| defaults.server_config.healths.readiness.server.socket_option.tcp_no_delay | bool | `true` | server listen socket option for tcp_no_delay functionality |
| defaults.server_config.healths.readiness.server.socket_option.tcp_quick_ack | bool | `true` | server listen socket option for tcp_quick_ack functionality |
| defaults.server_config.healths.readiness.server.socket_path | string | `""` | mysql socket_path |
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
| defaults.server_config.metrics.pprof.server.network | string | `"tcp"` | mysql network |
| defaults.server_config.metrics.pprof.server.probe_wait_time | string | `"3s"` | pprof server probe wait time |
| defaults.server_config.metrics.pprof.server.socket_option.ip_recover_destination_addr | bool | `false` | server listen socket option for ip_recover_destination_addr functionality |
| defaults.server_config.metrics.pprof.server.socket_option.ip_transparent | bool | `false` | server listen socket option for ip_transparent functionality |
| defaults.server_config.metrics.pprof.server.socket_option.reuse_addr | bool | `true` | server listen socket option for reuse_addr functionality |
| defaults.server_config.metrics.pprof.server.socket_option.reuse_port | bool | `true` | server listen socket option for reuse_port functionality |
| defaults.server_config.metrics.pprof.server.socket_option.tcp_cork | bool | `false` | server listen socket option for tcp_cork functionality |
| defaults.server_config.metrics.pprof.server.socket_option.tcp_defer_accept | bool | `true` | server listen socket option for tcp_defer_accept functionality |
| defaults.server_config.metrics.pprof.server.socket_option.tcp_fast_open | bool | `true` | server listen socket option for tcp_fast_open functionality |
| defaults.server_config.metrics.pprof.server.socket_option.tcp_no_delay | bool | `true` | server listen socket option for tcp_no_delay functionality |
| defaults.server_config.metrics.pprof.server.socket_option.tcp_quick_ack | bool | `true` | server listen socket option for tcp_quick_ack functionality |
| defaults.server_config.metrics.pprof.server.socket_path | string | `""` | mysql socket_path |
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
| defaults.server_config.metrics.prometheus.server.network | string | `"tcp"` | mysql network |
| defaults.server_config.metrics.prometheus.server.probe_wait_time | string | `"3s"` | prometheus server probe wait time |
| defaults.server_config.metrics.prometheus.server.socket_option.ip_recover_destination_addr | bool | `false` | server listen socket option for ip_recover_destination_addr functionality |
| defaults.server_config.metrics.prometheus.server.socket_option.ip_transparent | bool | `false` | server listen socket option for ip_transparent functionality |
| defaults.server_config.metrics.prometheus.server.socket_option.reuse_addr | bool | `true` | server listen socket option for reuse_addr functionality |
| defaults.server_config.metrics.prometheus.server.socket_option.reuse_port | bool | `true` | server listen socket option for reuse_port functionality |
| defaults.server_config.metrics.prometheus.server.socket_option.tcp_cork | bool | `false` | server listen socket option for tcp_cork functionality |
| defaults.server_config.metrics.prometheus.server.socket_option.tcp_defer_accept | bool | `true` | server listen socket option for tcp_defer_accept functionality |
| defaults.server_config.metrics.prometheus.server.socket_option.tcp_fast_open | bool | `true` | server listen socket option for tcp_fast_open functionality |
| defaults.server_config.metrics.prometheus.server.socket_option.tcp_no_delay | bool | `true` | server listen socket option for tcp_no_delay functionality |
| defaults.server_config.metrics.prometheus.server.socket_option.tcp_quick_ack | bool | `true` | server listen socket option for tcp_quick_ack functionality |
| defaults.server_config.metrics.prometheus.server.socket_path | string | `""` | mysql socket_path |
| defaults.server_config.metrics.prometheus.servicePort | int | `6061` | prometheus server service port |
| defaults.server_config.servers.grpc.enabled | bool | `true` | gRPC server enabled |
| defaults.server_config.servers.grpc.host | string | `"0.0.0.0"` | gRPC server host |
| defaults.server_config.servers.grpc.port | int | `8081` | gRPC server port |
| defaults.server_config.servers.grpc.server.grpc.bidirectional_stream_concurrency | int | `20` | gRPC server bidirectional stream concurrency |
| defaults.server_config.servers.grpc.server.grpc.connection_timeout | string | `""` | gRPC server connection timeout |
| defaults.server_config.servers.grpc.server.grpc.enable_reflection | bool | `true` | gRPC server reflection option |
| defaults.server_config.servers.grpc.server.grpc.header_table_size | int | `0` | gRPC server header table size |
| defaults.server_config.servers.grpc.server.grpc.initial_conn_window_size | int | `0` | gRPC server initial connection window size |
| defaults.server_config.servers.grpc.server.grpc.initial_window_size | int | `0` | gRPC server initial window size |
| defaults.server_config.servers.grpc.server.grpc.interceptors | list | `["RecoverInterceptor"]` | gRPC server interceptors |
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
| defaults.server_config.servers.grpc.server.network | string | `"tcp"` | mysql network |
| defaults.server_config.servers.grpc.server.probe_wait_time | string | `"3s"` | gRPC server probe wait time |
| defaults.server_config.servers.grpc.server.restart | bool | `true` | gRPC server restart |
| defaults.server_config.servers.grpc.server.socket_option.ip_recover_destination_addr | bool | `false` | server listen socket option for ip_recover_destination_addr functionality |
| defaults.server_config.servers.grpc.server.socket_option.ip_transparent | bool | `false` | server listen socket option for ip_transparent functionality |
| defaults.server_config.servers.grpc.server.socket_option.reuse_addr | bool | `true` | server listen socket option for reuse_addr functionality |
| defaults.server_config.servers.grpc.server.socket_option.reuse_port | bool | `true` | server listen socket option for reuse_port functionality |
| defaults.server_config.servers.grpc.server.socket_option.tcp_cork | bool | `false` | server listen socket option for tcp_cork functionality |
| defaults.server_config.servers.grpc.server.socket_option.tcp_defer_accept | bool | `true` | server listen socket option for tcp_defer_accept functionality |
| defaults.server_config.servers.grpc.server.socket_option.tcp_fast_open | bool | `true` | server listen socket option for tcp_fast_open functionality |
| defaults.server_config.servers.grpc.server.socket_option.tcp_no_delay | bool | `true` | server listen socket option for tcp_no_delay functionality |
| defaults.server_config.servers.grpc.server.socket_option.tcp_quick_ack | bool | `true` | server listen socket option for tcp_quick_ack functionality |
| defaults.server_config.servers.grpc.server.socket_path | string | `""` | mysql socket_path |
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
| defaults.server_config.servers.rest.server.network | string | `"tcp"` | mysql network |
| defaults.server_config.servers.rest.server.probe_wait_time | string | `"3s"` | REST server probe wait time |
| defaults.server_config.servers.rest.server.socket_option.ip_recover_destination_addr | bool | `false` | server listen socket option for ip_recover_destination_addr functionality |
| defaults.server_config.servers.rest.server.socket_option.ip_transparent | bool | `false` | server listen socket option for ip_transparent functionality |
| defaults.server_config.servers.rest.server.socket_option.reuse_addr | bool | `true` | server listen socket option for reuse_addr functionality |
| defaults.server_config.servers.rest.server.socket_option.reuse_port | bool | `true` | server listen socket option for reuse_port functionality |
| defaults.server_config.servers.rest.server.socket_option.tcp_cork | bool | `false` | server listen socket option for tcp_cork functionality |
| defaults.server_config.servers.rest.server.socket_option.tcp_defer_accept | bool | `true` | server listen socket option for tcp_defer_accept functionality |
| defaults.server_config.servers.rest.server.socket_option.tcp_fast_open | bool | `true` | server listen socket option for tcp_fast_open functionality |
| defaults.server_config.servers.rest.server.socket_option.tcp_no_delay | bool | `true` | server listen socket option for tcp_no_delay functionality |
| defaults.server_config.servers.rest.server.socket_option.tcp_quick_ack | bool | `true` | server listen socket option for tcp_quick_ack functionality |
| defaults.server_config.servers.rest.server.socket_path | string | `""` | mysql socket_path |
| defaults.server_config.servers.rest.servicePort | int | `8080` | REST server service port |
| defaults.server_config.tls.ca | string | `"/path/to/ca"` | TLS ca path |
| defaults.server_config.tls.cert | string | `"/path/to/cert"` | TLS cert path |
| defaults.server_config.tls.enabled | bool | `false` | TLS enabled |
| defaults.server_config.tls.insecure_skip_verify | bool | `false` | enable/disable skip SSL certificate verification |
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
| discoverer.discoverer.net.dialer.dual_stack_enabled | bool | `false` | TCP dialer dual stack enabled |
| discoverer.discoverer.net.dialer.keep_alive | string | `"10m"` | TCP dialer keep alive |
| discoverer.discoverer.net.dialer.timeout | string | `"30s"` | TCP dialer timeout |
| discoverer.discoverer.net.dns.cache_enabled | bool | `true` | TCP DNS cache enabled |
| discoverer.discoverer.net.dns.cache_expiration | string | `"24h"` | TCP DNS cache expiration |
| discoverer.discoverer.net.dns.refresh_duration | string | `"5m"` | TCP DNS cache refresh duration |
| discoverer.discoverer.net.socket_option.ip_recover_destination_addr | bool | `false` | server listen socket option for ip_recover_destination_addr functionality |
| discoverer.discoverer.net.socket_option.ip_transparent | bool | `false` | server listen socket option for ip_transparent functionality |
| discoverer.discoverer.net.socket_option.reuse_addr | bool | `true` | server listen socket option for reuse_addr functionality |
| discoverer.discoverer.net.socket_option.reuse_port | bool | `true` | server listen socket option for reuse_port functionality |
| discoverer.discoverer.net.socket_option.tcp_cork | bool | `false` | server listen socket option for tcp_cork functionality |
| discoverer.discoverer.net.socket_option.tcp_defer_accept | bool | `true` | server listen socket option for tcp_defer_accept functionality |
| discoverer.discoverer.net.socket_option.tcp_fast_open | bool | `true` | server listen socket option for tcp_fast_open functionality |
| discoverer.discoverer.net.socket_option.tcp_no_delay | bool | `true` | server listen socket option for tcp_no_delay functionality |
| discoverer.discoverer.net.socket_option.tcp_quick_ack | bool | `true` | server listen socket option for tcp_quick_ack functionality |
| discoverer.discoverer.net.tls.ca | string | `"/path/to/ca"` | TLS ca path |
| discoverer.discoverer.net.tls.cert | string | `"/path/to/cert"` | TLS cert path |
| discoverer.discoverer.net.tls.enabled | bool | `false` | TLS enabled |
| discoverer.discoverer.net.tls.insecure_skip_verify | bool | `false` | enable/disable skip SSL certificate verification |
| discoverer.discoverer.net.tls.key | string | `"/path/to/key"` | TLS key path |
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
| discoverer.podSecurityContext | object | `{"fsGroup":65532,"fsGroupChangePolicy":"OnRootMismatch","runAsGroup":65532,"runAsNonRoot":true,"runAsUser":65532}` | security context for pod |
| discoverer.progressDeadlineSeconds | int | `600` | progress deadline seconds |
| discoverer.resources | object | `{"limits":{"cpu":"600m","memory":"200Mi"},"requests":{"cpu":"200m","memory":"65Mi"}}` | compute resources |
| discoverer.revisionHistoryLimit | int | `2` | number of old history to retain to allow rollback |
| discoverer.rollingUpdate.maxSurge | string | `"25%"` | max surge of rolling update |
| discoverer.rollingUpdate.maxUnavailable | string | `"25%"` | max unavailable of rolling update |
| discoverer.securityContext | object | `{"allowPrivilegeEscalation":false,"capabilities":{"drop":["ALL"]},"privileged":false,"readOnlyRootFilesystem":true,"runAsGroup":65532,"runAsNonRoot":true,"runAsUser":65532}` | security context for container |
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
| gateway.filter.affinity.nodeAffinity.preferredDuringSchedulingIgnoredDuringExecution | list | `[]` | node affinity preferred scheduling terms |
| gateway.filter.affinity.nodeAffinity.requiredDuringSchedulingIgnoredDuringExecution.nodeSelectorTerms | list | `[]` | node affinity required node selectors |
| gateway.filter.affinity.podAffinity.preferredDuringSchedulingIgnoredDuringExecution | list | `[]` | pod affinity preferred scheduling terms |
| gateway.filter.affinity.podAffinity.requiredDuringSchedulingIgnoredDuringExecution | list | `[]` | pod affinity required scheduling terms |
| gateway.filter.affinity.podAntiAffinity.preferredDuringSchedulingIgnoredDuringExecution | list | `[{"podAffinityTerm":{"labelSelector":{"matchExpressions":[{"key":"app","operator":"In","values":["vald-filter-gateway"]}]},"topologyKey":"kubernetes.io/hostname"},"weight":100}]` | pod anti-affinity preferred scheduling terms |
| gateway.filter.affinity.podAntiAffinity.requiredDuringSchedulingIgnoredDuringExecution | list | `[]` | pod anti-affinity required scheduling terms |
| gateway.filter.annotations | object | `{}` | deployment annotations |
| gateway.filter.enabled | bool | `false` | gateway enabled |
| gateway.filter.env | list | `[]` | environment variables |
| gateway.filter.externalTrafficPolicy | string | `""` | external traffic policy (can be specified when service type is LoadBalancer or NodePort) : Cluster or Local |
| gateway.filter.gateway_config.egress_filter | object | `{"client":{},"distance_filters":[],"object_filters":[]}` | gRPC client config for egress filter |
| gateway.filter.gateway_config.egress_filter.client | object | `{}` | gRPC client config for egress filter (overrides defaults.grpc.client) |
| gateway.filter.gateway_config.egress_filter.distance_filters | list | `[]` | distance egress vector filter targets |
| gateway.filter.gateway_config.egress_filter.object_filters | list | `[]` | object egress vector filter targets |
| gateway.filter.gateway_config.gateway_client | object | `{}` | gRPC client for next gateway (overrides defaults.grpc.client) |
| gateway.filter.gateway_config.ingress_filter | object | `{"client":{},"insert_filters":[],"search_filters":[],"update_filters":[],"upsert_filters":[],"vectorizer":""}` | gRPC client config for ingress filter |
| gateway.filter.gateway_config.ingress_filter.client | object | `{}` | gRPC client for ingress filter (overrides defaults.grpc.client) |
| gateway.filter.gateway_config.ingress_filter.insert_filters | list | `[]` | insert ingress vector filter targets |
| gateway.filter.gateway_config.ingress_filter.search_filters | list | `[]` | search ingress vector filter targets |
| gateway.filter.gateway_config.ingress_filter.update_filters | list | `[]` | update ingress vector filter targets |
| gateway.filter.gateway_config.ingress_filter.upsert_filters | list | `[]` | upsert ingress vector filter targets |
| gateway.filter.gateway_config.ingress_filter.vectorizer | string | `""` | object ingress vectorize filter targets |
| gateway.filter.hpa.enabled | bool | `true` | HPA enabled |
| gateway.filter.hpa.targetCPUUtilizationPercentage | int | `80` | HPA CPU utilization percentage |
| gateway.filter.image.pullPolicy | string | `"Always"` | image pull policy |
| gateway.filter.image.repository | string | `"vdaas/vald-filter-gateway"` | image repository |
| gateway.filter.image.tag | string | `""` | image tag (overrides defaults.image.tag) |
| gateway.filter.ingress.annotations | object | `{"nginx.ingress.kubernetes.io/grpc-backend":"true"}` | annotations for ingress |
| gateway.filter.ingress.enabled | bool | `false` | gateway ingress enabled |
| gateway.filter.ingress.host | string | `"filter.gateway.vald.vdaas.org"` | ingress hostname |
| gateway.filter.ingress.pathType | string | `"ImplementationSpecific"` | gateway ingress pathType |
| gateway.filter.ingress.servicePort | string | `"grpc"` | service port to be exposed by ingress |
| gateway.filter.initContainers | list | `[{"image":"busybox","name":"wait-for-gateway-lb","sleepDuration":2,"target":"gateway-lb","type":"wait-for"}]` | init containers |
| gateway.filter.kind | string | `"Deployment"` | deployment kind: Deployment or DaemonSet |
| gateway.filter.logging | object | `{}` | logging config (overrides defaults.logging) |
| gateway.filter.maxReplicas | int | `9` | maximum number of replicas. if HPA is disabled, this value will be ignored. |
| gateway.filter.maxUnavailable | string | `"50%"` | maximum number of unavailable replicas |
| gateway.filter.minReplicas | int | `3` | minimum number of replicas. if HPA is disabled, the replicas will be set to this value |
| gateway.filter.name | string | `"vald-filter-gateway"` | name of filter gateway deployment |
| gateway.filter.nodeName | string | `""` | node name |
| gateway.filter.nodeSelector | object | `{}` | node selector |
| gateway.filter.observability | object | `{"jaeger":{"service_name":"vald-filter-gateway"},"stackdriver":{"profiler":{"service":"vald-filter-gateway"}}}` | observability config (overrides defaults.observability) |
| gateway.filter.podAnnotations | object | `{}` | pod annotations |
| gateway.filter.podPriority.enabled | bool | `true` | gateway pod PriorityClass enabled |
| gateway.filter.podPriority.value | int | `1000000` | gateway pod PriorityClass value |
| gateway.filter.podSecurityContext | object | `{"fsGroup":65532,"fsGroupChangePolicy":"OnRootMismatch","runAsGroup":65532,"runAsNonRoot":true,"runAsUser":65532}` | security context for pod |
| gateway.filter.progressDeadlineSeconds | int | `600` | progress deadline seconds |
| gateway.filter.resources | object | `{"limits":{"cpu":"2000m","memory":"700Mi"},"requests":{"cpu":"200m","memory":"150Mi"}}` | compute resources |
| gateway.filter.revisionHistoryLimit | int | `2` | number of old history to retain to allow rollback |
| gateway.filter.rollingUpdate.maxSurge | string | `"25%"` | max surge of rolling update |
| gateway.filter.rollingUpdate.maxUnavailable | string | `"25%"` | max unavailable of rolling update |
| gateway.filter.securityContext | object | `{"allowPrivilegeEscalation":false,"capabilities":{"drop":["ALL"]},"privileged":false,"readOnlyRootFilesystem":true,"runAsGroup":65532,"runAsNonRoot":true,"runAsUser":65532}` | security context for container |
| gateway.filter.server_config | object | `{"healths":{"liveness":{},"readiness":{}},"metrics":{"pprof":{},"prometheus":{}},"servers":{"grpc":{},"rest":{}}}` | server config (overrides defaults.server_config) |
| gateway.filter.service.annotations | object | `{}` | service annotations |
| gateway.filter.service.labels | object | `{}` | service labels |
| gateway.filter.serviceType | string | `"ClusterIP"` | service type: ClusterIP, LoadBalancer or NodePort |
| gateway.filter.terminationGracePeriodSeconds | int | `30` | duration in seconds pod needs to terminate gracefully |
| gateway.filter.time_zone | string | `""` | Time zone |
| gateway.filter.tolerations | list | `[]` | tolerations |
| gateway.filter.topologySpreadConstraints | list | `[]` | topology spread constraints of gateway pods |
| gateway.filter.version | string | `"v0.0.0"` | version of gateway config |
| gateway.filter.volumeMounts | list | `[]` | volume mounts |
| gateway.filter.volumes | list | `[]` | volumes |
| gateway.lb.affinity.nodeAffinity.preferredDuringSchedulingIgnoredDuringExecution | list | `[]` | node affinity preferred scheduling terms |
| gateway.lb.affinity.nodeAffinity.requiredDuringSchedulingIgnoredDuringExecution.nodeSelectorTerms | list | `[]` | node affinity required node selectors |
| gateway.lb.affinity.podAffinity.preferredDuringSchedulingIgnoredDuringExecution | list | `[]` | pod affinity preferred scheduling terms |
| gateway.lb.affinity.podAffinity.requiredDuringSchedulingIgnoredDuringExecution | list | `[]` | pod affinity required scheduling terms |
| gateway.lb.affinity.podAntiAffinity.preferredDuringSchedulingIgnoredDuringExecution | list | `[{"podAffinityTerm":{"labelSelector":{"matchExpressions":[{"key":"app","operator":"In","values":["vald-lb-gateway"]}]},"topologyKey":"kubernetes.io/hostname"},"weight":100}]` | pod anti-affinity preferred scheduling terms |
| gateway.lb.affinity.podAntiAffinity.requiredDuringSchedulingIgnoredDuringExecution | list | `[]` | pod anti-affinity required scheduling terms |
| gateway.lb.annotations | object | `{}` | deployment annotations |
| gateway.lb.enabled | bool | `true` | gateway enabled |
| gateway.lb.env | list | `[{"name":"MY_POD_NAMESPACE","valueFrom":{"fieldRef":{"fieldPath":"metadata.namespace"}}}]` | environment variables |
| gateway.lb.externalTrafficPolicy | string | `""` | external traffic policy (can be specified when service type is LoadBalancer or NodePort) : Cluster or Local |
| gateway.lb.gateway_config.agent_namespace | string | `"_MY_POD_NAMESPACE_"` | agent namespace |
| gateway.lb.gateway_config.discoverer.agent_client_options | object | `{}` | gRPC client options for agents (overrides defaults.grpc.client) |
| gateway.lb.gateway_config.discoverer.client | object | `{}` | gRPC client for discoverer (overrides defaults.grpc.client) |
| gateway.lb.gateway_config.discoverer.duration | string | `"200ms"` |  |
| gateway.lb.gateway_config.index_replica | int | `5` | number of index replica |
| gateway.lb.gateway_config.node_name | string | `""` | node name |
| gateway.lb.hpa.enabled | bool | `true` | HPA enabled |
| gateway.lb.hpa.targetCPUUtilizationPercentage | int | `80` | HPA CPU utilization percentage |
| gateway.lb.image.pullPolicy | string | `"Always"` | image pull policy |
| gateway.lb.image.repository | string | `"vdaas/vald-lb-gateway"` | image repository |
| gateway.lb.image.tag | string | `""` | image tag (overrides defaults.image.tag) |
| gateway.lb.ingress.annotations | object | `{"nginx.ingress.kubernetes.io/grpc-backend":"true"}` | annotations for ingress |
| gateway.lb.ingress.enabled | bool | `false` | gateway ingress enabled |
| gateway.lb.ingress.host | string | `"lb.gateway.vald.vdaas.org"` | ingress hostname |
| gateway.lb.ingress.pathType | string | `"ImplementationSpecific"` | gateway ingress pathType |
| gateway.lb.ingress.servicePort | string | `"grpc"` | service port to be exposed by ingress |
| gateway.lb.initContainers | list | `[{"image":"busybox","name":"wait-for-discoverer","sleepDuration":2,"target":"discoverer","type":"wait-for"},{"image":"busybox","name":"wait-for-agent","sleepDuration":2,"target":"agent","type":"wait-for"}]` | init containers |
| gateway.lb.kind | string | `"Deployment"` | deployment kind: Deployment or DaemonSet |
| gateway.lb.logging | object | `{}` | logging config (overrides defaults.logging) |
| gateway.lb.maxReplicas | int | `9` | maximum number of replicas. if HPA is disabled, this value will be ignored. |
| gateway.lb.maxUnavailable | string | `"50%"` | maximum number of unavailable replicas |
| gateway.lb.minReplicas | int | `3` | minimum number of replicas. if HPA is disabled, the replicas will be set to this value |
| gateway.lb.name | string | `"vald-lb-gateway"` | name of gateway deployment |
| gateway.lb.nodeName | string | `""` | node name |
| gateway.lb.nodeSelector | object | `{}` | node selector |
| gateway.lb.observability | object | `{"jaeger":{"service_name":"vald-lb-gateway"},"stackdriver":{"profiler":{"service":"vald-lb-gateway"}}}` | observability config (overrides defaults.observability) |
| gateway.lb.podAnnotations | object | `{}` | pod annotations |
| gateway.lb.podPriority.enabled | bool | `true` | gateway pod PriorityClass enabled |
| gateway.lb.podPriority.value | int | `1000000` | gateway pod PriorityClass value |
| gateway.lb.podSecurityContext | object | `{"fsGroup":65532,"fsGroupChangePolicy":"OnRootMismatch","runAsGroup":65532,"runAsNonRoot":true,"runAsUser":65532}` | security context for pod |
| gateway.lb.progressDeadlineSeconds | int | `600` | progress deadline seconds |
| gateway.lb.resources | object | `{"limits":{"cpu":"2000m","memory":"700Mi"},"requests":{"cpu":"200m","memory":"150Mi"}}` | compute resources |
| gateway.lb.revisionHistoryLimit | int | `2` | number of old history to retain to allow rollback |
| gateway.lb.rollingUpdate.maxSurge | string | `"25%"` | max surge of rolling update |
| gateway.lb.rollingUpdate.maxUnavailable | string | `"25%"` | max unavailable of rolling update |
| gateway.lb.securityContext | object | `{"allowPrivilegeEscalation":false,"capabilities":{"drop":["ALL"]},"privileged":false,"readOnlyRootFilesystem":true,"runAsGroup":65532,"runAsNonRoot":true,"runAsUser":65532}` | security context for container |
| gateway.lb.server_config | object | `{"healths":{"liveness":{},"readiness":{}},"metrics":{"pprof":{},"prometheus":{}},"servers":{"grpc":{},"rest":{}}}` | server config (overrides defaults.server_config) |
| gateway.lb.service.annotations | object | `{}` | service annotations |
| gateway.lb.service.labels | object | `{}` | service labels |
| gateway.lb.serviceType | string | `"ClusterIP"` | service type: ClusterIP, LoadBalancer or NodePort |
| gateway.lb.terminationGracePeriodSeconds | int | `30` | duration in seconds pod needs to terminate gracefully |
| gateway.lb.time_zone | string | `""` | Time zone |
| gateway.lb.tolerations | list | `[]` | tolerations |
| gateway.lb.topologySpreadConstraints | list | `[]` | topology spread constraints of gateway pods |
| gateway.lb.version | string | `"v0.0.0"` | version of gateway config |
| gateway.lb.volumeMounts | list | `[]` | volume mounts |
| gateway.lb.volumes | list | `[]` | volumes |
| initializer.cassandra.configmap.backup.enabled | bool | `true` | backup table enabled |
| initializer.cassandra.configmap.backup.name | string | `"backup_vector"` | name of backup table |
| initializer.cassandra.configmap.enabled | bool | `false` | cassandra schema configmap will be created |
| initializer.cassandra.configmap.filename | string | `"init.cql"` | cassandra schema filename |
| initializer.cassandra.configmap.keyspace | string | `"vald"` | cassandra keyspace |
| initializer.cassandra.configmap.meta | object | `{"enabled":true,"name":{"kv":"kv","vk":"vk"}}` | cassandra settings for metadata store |
| initializer.cassandra.configmap.meta.enabled | bool | `true` | meta table enabled |
| initializer.cassandra.configmap.meta.name.kv | string | `"kv"` | name of KV table |
| initializer.cassandra.configmap.meta.name.vk | string | `"vk"` | name of VK table |
| initializer.cassandra.configmap.name | string | `"cassandra-initdb"` | cassandra schema configmap name |
| initializer.cassandra.configmap.replication_class | string | `"SimpleStrategy"` | cassandra replication class |
| initializer.cassandra.configmap.replication_factor | int | `3` | cassandra replication factor |
| initializer.cassandra.configmap.user | string | `"root"` | cassandra user |
| initializer.cassandra.enabled | bool | `false` | cassandra initializer job enabled |
| initializer.cassandra.env | list | `[{"name":"CASSANDRA_HOST","value":"cassandra.default.svc.cluster.local"}]` | environment variables |
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
| manager.index.affinity.nodeAffinity.preferredDuringSchedulingIgnoredDuringExecution | list | `[]` | node affinity preferred scheduling terms |
| manager.index.affinity.nodeAffinity.requiredDuringSchedulingIgnoredDuringExecution.nodeSelectorTerms | list | `[]` | node affinity required node selectors |
| manager.index.affinity.podAffinity.preferredDuringSchedulingIgnoredDuringExecution | list | `[]` | pod affinity preferred scheduling terms |
| manager.index.affinity.podAffinity.requiredDuringSchedulingIgnoredDuringExecution | list | `[]` | pod affinity required scheduling terms |
| manager.index.affinity.podAntiAffinity.preferredDuringSchedulingIgnoredDuringExecution | list | `[]` | pod anti-affinity preferred scheduling terms |
| manager.index.affinity.podAntiAffinity.requiredDuringSchedulingIgnoredDuringExecution | list | `[]` | pod anti-affinity required scheduling terms |
| manager.index.annotations | object | `{}` | deployment annotations |
| manager.index.enabled | bool | `true` | index manager enabled |
| manager.index.env | list | `[{"name":"MY_POD_NAMESPACE","valueFrom":{"fieldRef":{"fieldPath":"metadata.namespace"}}}]` | environment variables |
| manager.index.externalTrafficPolicy | string | `""` | external traffic policy (can be specified when service type is LoadBalancer or NodePort) : Cluster or Local |
| manager.index.image.pullPolicy | string | `"Always"` | image pull policy |
| manager.index.image.repository | string | `"vdaas/vald-manager-index"` | image repository |
| manager.index.image.tag | string | `""` | image tag (overrides defaults.image.tag) |
| manager.index.indexer.agent_namespace | string | `"_MY_POD_NAMESPACE_"` | namespace of agent pods to manage |
| manager.index.indexer.auto_index_check_duration | string | `"1m"` | check duration of automatic indexing |
| manager.index.indexer.auto_index_duration_limit | string | `"30m"` | limit duration of automatic indexing |
| manager.index.indexer.auto_index_length | int | `100` | number of cache to trigger automatic indexing |
| manager.index.indexer.auto_save_index_duration_limit | string | `"3h"` | limit duration of automatic index saving |
| manager.index.indexer.auto_save_index_wait_duration | string | `"10m"` | duration of automatic index saving wait duration for next saving |
| manager.index.indexer.concurrency | int | `1` | concurrency |
| manager.index.indexer.creation_pool_size | int | `10000` | number of pool size of create index processing |
| manager.index.indexer.discoverer.agent_client_options | object | `{"dial_option":{"net":{"dialer":{"keep_alive":"15m"}}}}` | gRPC client options for agents (overrides defaults.grpc.client) |
| manager.index.indexer.discoverer.client | object | `{}` | gRPC client for discoverer (overrides defaults.grpc.client) |
| manager.index.indexer.discoverer.duration | string | `"500ms"` | refresh duration to discover |
| manager.index.indexer.node_name | string | `""` | node name |
| manager.index.initContainers | list | `[{"image":"busybox","name":"wait-for-agent","sleepDuration":2,"target":"agent","type":"wait-for"},{"image":"busybox","name":"wait-for-discoverer","sleepDuration":2,"target":"discoverer","type":"wait-for"}]` | init containers |
| manager.index.kind | string | `"Deployment"` | deployment kind: Deployment or DaemonSet |
| manager.index.logging | object | `{}` | logging config (overrides defaults.logging) |
| manager.index.maxUnavailable | string | `"50%"` | maximum number of unavailable replicas |
| manager.index.name | string | `"vald-manager-index"` | name of index manager deployment |
| manager.index.nodeName | string | `""` | node name |
| manager.index.nodeSelector | object | `{}` | node selector |
| manager.index.observability | object | `{"jaeger":{"service_name":"vald-manager-index"},"stackdriver":{"profiler":{"service":"vald-manager-index"}}}` | observability config (overrides defaults.observability) |
| manager.index.podAnnotations | object | `{}` | pod annotations |
| manager.index.podPriority.enabled | bool | `true` | index manager pod PriorityClass enabled |
| manager.index.podPriority.value | int | `1000000` | index manager pod PriorityClass value |
| manager.index.podSecurityContext | object | `{"fsGroup":65532,"fsGroupChangePolicy":"OnRootMismatch","runAsGroup":65532,"runAsNonRoot":true,"runAsUser":65532}` | security context for pod |
| manager.index.progressDeadlineSeconds | int | `600` | progress deadline seconds |
| manager.index.replicas | int | `1` | number of replicas |
| manager.index.resources | object | `{"limits":{"cpu":"1000m","memory":"500Mi"},"requests":{"cpu":"200m","memory":"80Mi"}}` | compute resources |
| manager.index.revisionHistoryLimit | int | `2` | number of old history to retain to allow rollback |
| manager.index.rollingUpdate.maxSurge | string | `"25%"` | max surge of rolling update |
| manager.index.rollingUpdate.maxUnavailable | string | `"25%"` | max unavailable of rolling update |
| manager.index.securityContext | object | `{"allowPrivilegeEscalation":false,"capabilities":{"drop":["ALL"]},"privileged":false,"readOnlyRootFilesystem":true,"runAsGroup":65532,"runAsNonRoot":true,"runAsUser":65532}` | security context for container |
| manager.index.server_config | object | `{"healths":{"liveness":{},"readiness":{}},"metrics":{"pprof":{},"prometheus":{}},"servers":{"grpc":{},"rest":{}}}` | server config (overrides defaults.server_config) |
| manager.index.service.annotations | object | `{}` | service annotations |
| manager.index.service.labels | object | `{}` | service labels |
| manager.index.serviceType | string | `"ClusterIP"` | service type: ClusterIP, LoadBalancer or NodePort |
| manager.index.terminationGracePeriodSeconds | int | `30` | duration in seconds pod needs to terminate gracefully |
| manager.index.time_zone | string | `""` | Time zone |
| manager.index.tolerations | list | `[]` | tolerations |
| manager.index.topologySpreadConstraints | list | `[]` | topology spread constraints of index manager pods |
| manager.index.version | string | `"v0.0.0"` | version of index manager config |
| manager.index.volumeMounts | list | `[]` | volume mounts |
| manager.index.volumes | list | `[]` | volumes |
