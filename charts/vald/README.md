Vald
===

This is a Helm chart to install Vald components.

Current chart version is `v0.0.25`

Install
---

Add Vald Helm repository

    $ helm repo add vald https://vald.vdaas.org/charts

Run the following command to install the chart,

    $ helm install --generate-name vald/vald


Configuration
---

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| agent.hpa.enabled | bool | `false` |  |
| agent.hpa.targetCPUUtilizationPercentage | int | `80` |  |
| agent.image.pullPolicy | string | `"Always"` |  |
| agent.image.repository | string | `"vdaas/vald-agent-ngt"` |  |
| agent.kind | string | `"StatefulSet"` |  |
| agent.maxReplicas | int | `300` |  |
| agent.maxUnavailable | int | `1` |  |
| agent.minReplicas | int | `20` |  |
| agent.name | string | `"vald-agent-ngt"` |  |
| agent.ngt.auto_index_check_duration | string | `"30m"` |  |
| agent.ngt.auto_index_length | int | `100` |  |
| agent.ngt.auto_index_limit | string | `"24h"` |  |
| agent.ngt.bulk_insert_chunk_size | int | `10` |  |
| agent.ngt.creation_edge_size | int | `20` |  |
| agent.ngt.dimension | int | `4096` |  |
| agent.ngt.distance_type | string | `"l2"` |  |
| agent.ngt.enable_in_memory_mode | bool | `true` |  |
| agent.ngt.object_type | string | `"float"` |  |
| agent.ngt.search_edge_size | int | `10` |  |
| agent.observability.jaeger.service_name | string | `"vald-agent-ngt"` |  |
| agent.podManagementPolicy | string | `"OrderedReady"` |  |
| agent.podPriority.enabled | bool | `true` |  |
| agent.podPriority.value | int | `1000000000` |  |
| agent.progressDeadlineSeconds | int | `600` |  |
| agent.resources.requests.cpu | string | `"300m"` |  |
| agent.resources.requests.memory | string | `"4Gi"` |  |
| agent.revisionHistoryLimit | int | `2` |  |
| agent.rollingUpdate.maxSurge | string | `"25%"` |  |
| agent.rollingUpdate.maxUnavailable | string | `"25%"` |  |
| agent.rollingUpdate.partition | int | `0` |  |
| agent.server_config.full_shutdown_duration | string | `"600s"` |  |
| agent.server_config.healths.liveness.enabled | bool | `false` |  |
| agent.server_config.healths.readiness.enabled | bool | `false` |  |
| agent.server_config.metrics.pprof.enabled | bool | `false` |  |
| agent.server_config.metrics.prometheus.enabled | bool | `false` |  |
| agent.server_config.prefix | string | `"agent"` |  |
| agent.server_config.servers.grpc.enabled | bool | `false` |  |
| agent.server_config.servers.rest.enabled | bool | `false` |  |
| agent.server_config.tls.enabled | bool | `false` |  |
| agent.serviceType | string | `"ClusterIP"` |  |
| agent.terminationGracePeriodSeconds | int | `30` |  |
| agent.version | string | `"v0.0.0"` |  |
| backupManager.cassandra.config.connect_timeout | string | `"600ms"` |  |
| backupManager.cassandra.config.consistency | string | `"quorum"` |  |
| backupManager.cassandra.config.cql_version | string | `"3.0.0"` |  |
| backupManager.cassandra.config.default_idempotence | bool | `false` |  |
| backupManager.cassandra.config.default_timestamp | bool | `true` |  |
| backupManager.cassandra.config.disable_initial_host_lookup | bool | `false` |  |
| backupManager.cassandra.config.disable_node_status_events | bool | `false` |  |
| backupManager.cassandra.config.disable_skip_metadata | bool | `false` |  |
| backupManager.cassandra.config.disable_topology_events | bool | `false` |  |
| backupManager.cassandra.config.enable_host_verification | bool | `false` |  |
| backupManager.cassandra.config.hosts[0] | string | `"cassandra-0.cassandra.default.svc.cluster.local"` |  |
| backupManager.cassandra.config.hosts[1] | string | `"cassandra-1.cassandra.default.svc.cluster.local"` |  |
| backupManager.cassandra.config.hosts[2] | string | `"cassandra-2.cassandra.default.svc.cluster.local"` |  |
| backupManager.cassandra.config.ignore_peer_addr | bool | `false` |  |
| backupManager.cassandra.config.keyspace | string | `"vald"` |  |
| backupManager.cassandra.config.max_prepared_stmts | int | `1000` |  |
| backupManager.cassandra.config.max_routing_key_info | int | `1000` |  |
| backupManager.cassandra.config.max_wait_schema_agreement | string | `"1m"` |  |
| backupManager.cassandra.config.meta_table | string | `"meta_vector"` |  |
| backupManager.cassandra.config.num_conns | int | `2` |  |
| backupManager.cassandra.config.page_size | int | `5000` |  |
| backupManager.cassandra.config.password | string | `"_CASSANDRA_PASSWORD_"` |  |
| backupManager.cassandra.config.pool_config.data_center | string | `""` |  |
| backupManager.cassandra.config.pool_config.dc_aware_routing | bool | `false` |  |
| backupManager.cassandra.config.pool_config.non_local_replicas_fallback | bool | `false` |  |
| backupManager.cassandra.config.pool_config.shuffle_replicas | bool | `false` |  |
| backupManager.cassandra.config.port | int | `9042` |  |
| backupManager.cassandra.config.proto_version | int | `0` |  |
| backupManager.cassandra.config.reconnect_interval | string | `"1m"` |  |
| backupManager.cassandra.config.reconnection_policy.initial_interval | string | `"1m"` |  |
| backupManager.cassandra.config.reconnection_policy.max_retries | int | `3` |  |
| backupManager.cassandra.config.retry_policy.max_duration | string | `"30s"` |  |
| backupManager.cassandra.config.retry_policy.min_duration | string | `"1s"` |  |
| backupManager.cassandra.config.retry_policy.num_retries | int | `3` |  |
| backupManager.cassandra.config.socket_keepalive | string | `"0s"` |  |
| backupManager.cassandra.config.tcp.dialer.dual_stack_enabled | bool | `false` |  |
| backupManager.cassandra.config.tcp.dialer.keep_alive | string | `"10m"` |  |
| backupManager.cassandra.config.tcp.dialer.timeout | string | `"30s"` |  |
| backupManager.cassandra.config.tcp.dns.cache_enabled | bool | `true` |  |
| backupManager.cassandra.config.tcp.dns.cache_expiration | string | `"24h"` |  |
| backupManager.cassandra.config.tcp.dns.refresh_duration | string | `"5m"` |  |
| backupManager.cassandra.config.timeout | string | `"600ms"` |  |
| backupManager.cassandra.config.tls.ca | string | `"/path/to/ca"` |  |
| backupManager.cassandra.config.tls.cert | string | `"/path/to/cert"` |  |
| backupManager.cassandra.config.tls.enabled | bool | `false` |  |
| backupManager.cassandra.config.tls.key | string | `"/path/to/key"` |  |
| backupManager.cassandra.config.username | string | `"root"` |  |
| backupManager.cassandra.config.write_coalesce_wait_time | string | `"200ms"` |  |
| backupManager.cassandra.enabled | bool | `false` |  |
| backupManager.env[0].name | string | `"MYSQL_PASSWORD"` |  |
| backupManager.env[0].valueFrom.secretKeyRef.key | string | `"password"` |  |
| backupManager.env[0].valueFrom.secretKeyRef.name | string | `"mysql-secret"` |  |
| backupManager.hpa.enabled | bool | `true` |  |
| backupManager.hpa.targetCPUUtilizationPercentage | int | `80` |  |
| backupManager.image.pullPolicy | string | `"Always"` |  |
| backupManager.image.repository | string | `"vdaas/vald-manager-backup-mysql"` |  |
| backupManager.initContainers[0].env[0].name | string | `"MYSQL_PASSWORD"` |  |
| backupManager.initContainers[0].env[0].valueFrom.secretKeyRef.key | string | `"password"` |  |
| backupManager.initContainers[0].env[0].valueFrom.secretKeyRef.name | string | `"mysql-secret"` |  |
| backupManager.initContainers[0].image | string | `"mysql:latest"` |  |
| backupManager.initContainers[0].mysql.hosts[0] | string | `"mysql.default.svc.cluster.local"` |  |
| backupManager.initContainers[0].mysql.options[0] | string | `"-uroot"` |  |
| backupManager.initContainers[0].mysql.options[1] | string | `"-p${MYSQL_PASSWORD}"` |  |
| backupManager.initContainers[0].name | string | `"wait-for-mysql"` |  |
| backupManager.initContainers[0].sleepDuration | int | `2` |  |
| backupManager.initContainers[0].type | string | `"wait-for-mysql"` |  |
| backupManager.kind | string | `"Deployment"` |  |
| backupManager.maxReplicas | int | `15` |  |
| backupManager.maxUnavailable | string | `"50%"` |  |
| backupManager.minReplicas | int | `3` |  |
| backupManager.mysql.config.conn_max_life_time | string | `"30s"` |  |
| backupManager.mysql.config.db | string | `"mysql"` |  |
| backupManager.mysql.config.host | string | `"mysql.default.svc.cluster.local"` |  |
| backupManager.mysql.config.max_idle_conns | int | `100` |  |
| backupManager.mysql.config.max_open_conns | int | `100` |  |
| backupManager.mysql.config.name | string | `"vald"` |  |
| backupManager.mysql.config.pass | string | `"_MYSQL_PASSWORD_"` |  |
| backupManager.mysql.config.port | int | `3306` |  |
| backupManager.mysql.config.tcp.dialer.dual_stack_enabled | bool | `false` |  |
| backupManager.mysql.config.tcp.dialer.keep_alive | string | `"5m"` |  |
| backupManager.mysql.config.tcp.dialer.timeout | string | `"5s"` |  |
| backupManager.mysql.config.tcp.dns.cache_enabled | bool | `true` |  |
| backupManager.mysql.config.tcp.dns.cache_expiration | string | `"24h"` |  |
| backupManager.mysql.config.tcp.dns.refresh_duration | string | `"1h"` |  |
| backupManager.mysql.config.tcp.tls.ca | string | `"/path/to/ca"` |  |
| backupManager.mysql.config.tcp.tls.cert | string | `"/path/to/cert"` |  |
| backupManager.mysql.config.tcp.tls.enabled | bool | `false` |  |
| backupManager.mysql.config.tcp.tls.key | string | `"/path/to/key"` |  |
| backupManager.mysql.config.tls.ca | string | `"/path/to/ca"` |  |
| backupManager.mysql.config.tls.cert | string | `"/path/to/cert"` |  |
| backupManager.mysql.config.tls.enabled | bool | `false` |  |
| backupManager.mysql.config.tls.key | string | `"/path/to/key"` |  |
| backupManager.mysql.config.user | string | `"root"` |  |
| backupManager.mysql.enabled | bool | `true` |  |
| backupManager.name | string | `"vald-manager-backup"` |  |
| backupManager.observability.jaeger.service_name | string | `"vald-manager-backup"` |  |
| backupManager.progressDeadlineSeconds | int | `600` |  |
| backupManager.resources.limits.cpu | string | `"500m"` |  |
| backupManager.resources.limits.memory | string | `"150Mi"` |  |
| backupManager.resources.requests.cpu | string | `"100m"` |  |
| backupManager.resources.requests.memory | string | `"50Mi"` |  |
| backupManager.revisionHistoryLimit | int | `2` |  |
| backupManager.rollingUpdate.maxSurge | string | `"25%"` |  |
| backupManager.rollingUpdate.maxUnavailable | string | `"25%"` |  |
| backupManager.server_config.full_shutdown_duration | string | `"600s"` |  |
| backupManager.server_config.healths.liveness.enabled | bool | `false` |  |
| backupManager.server_config.healths.readiness.enabled | bool | `false` |  |
| backupManager.server_config.metrics.pprof.enabled | bool | `false` |  |
| backupManager.server_config.metrics.prometheus.enabled | bool | `false` |  |
| backupManager.server_config.prefix | string | `"backup-manager"` |  |
| backupManager.server_config.servers.grpc.enabled | bool | `false` |  |
| backupManager.server_config.servers.rest.enabled | bool | `false` |  |
| backupManager.server_config.tls.enabled | bool | `false` |  |
| backupManager.serviceType | string | `"ClusterIP"` |  |
| backupManager.terminationGracePeriodSeconds | int | `30` |  |
| backupManager.version | string | `"v0.0.0"` |  |
| compressor.backup.client | object | `{}` |  |
| compressor.compress.buffer | int | `100` |  |
| compressor.compress.compress_algorithm | string | `"zstd"` |  |
| compressor.compress.compression_level | int | `10` |  |
| compressor.compress.concurrent_limit | int | `10` |  |
| compressor.hpa.enabled | bool | `true` |  |
| compressor.hpa.targetCPUUtilizationPercentage | int | `80` |  |
| compressor.image.pullPolicy | string | `"Always"` |  |
| compressor.image.repository | string | `"vdaas/vald-manager-compressor"` |  |
| compressor.initContainers[0].image | string | `"busybox"` |  |
| compressor.initContainers[0].name | string | `"wait-for-manager-backup"` |  |
| compressor.initContainers[0].sleepDuration | int | `2` |  |
| compressor.initContainers[0].target | string | `"manager-backup"` |  |
| compressor.initContainers[0].type | string | `"wait-for"` |  |
| compressor.kind | string | `"Deployment"` |  |
| compressor.maxReplicas | int | `15` |  |
| compressor.maxUnavailable | string | `"50%"` |  |
| compressor.minReplicas | int | `3` |  |
| compressor.name | string | `"vald-manager-compressor"` |  |
| compressor.observability.jaeger.service_name | string | `"vald-manager-compressor"` |  |
| compressor.progressDeadlineSeconds | int | `600` |  |
| compressor.resources.limits.cpu | string | `"800m"` |  |
| compressor.resources.limits.memory | string | `"500Mi"` |  |
| compressor.resources.requests.cpu | string | `"300m"` |  |
| compressor.resources.requests.memory | string | `"50Mi"` |  |
| compressor.revisionHistoryLimit | int | `2` |  |
| compressor.rollingUpdate.maxSurge | string | `"25%"` |  |
| compressor.rollingUpdate.maxUnavailable | string | `"25%"` |  |
| compressor.server_config.full_shutdown_duration | string | `"600s"` |  |
| compressor.server_config.healths.liveness.enabled | bool | `false` |  |
| compressor.server_config.healths.readiness.enabled | bool | `false` |  |
| compressor.server_config.metrics.pprof.enabled | bool | `false` |  |
| compressor.server_config.metrics.prometheus.enabled | bool | `false` |  |
| compressor.server_config.prefix | string | `"manager-compressor"` |  |
| compressor.server_config.servers.grpc.enabled | bool | `false` |  |
| compressor.server_config.servers.rest.enabled | bool | `false` |  |
| compressor.server_config.tls.enabled | bool | `false` |  |
| compressor.serviceType | string | `"ClusterIP"` |  |
| compressor.terminationGracePeriodSeconds | int | `30` |  |
| compressor.version | string | `"v0.0.0"` |  |
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
| defaults.grpc.client.connection_pool | int | `3` | number of gRPC client connection pool |
| defaults.grpc.client.dial_option.enable_backoff | bool | `false` | gRPC client dial option backoff enabled |
| defaults.grpc.client.dial_option.initial_connection_window_size | int | `0` | gRPC client dial option initial connection window size |
| defaults.grpc.client.dial_option.initial_window_size | int | `0` | gRPC client dial option initial window size |
| defaults.grpc.client.dial_option.insecure | bool | `true` | gRPC client dial option insecure enabled |
| defaults.grpc.client.dial_option.keep_alive.permit_without_stream | bool | `false` | gRPC client keep alive permit without stream |
| defaults.grpc.client.dial_option.keep_alive.time | string | `""` | gRPC client keep alive time |
| defaults.grpc.client.dial_option.keep_alive.timeout | string | `""` | gRPC client keep alive timeout |
| defaults.grpc.client.dial_option.max_backoff_delay | string | `""` | gRPC client dial option max backoff delay |
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
| defaults.image.tag | string | `"v0.0.25"` | image tag |
| defaults.logging.format | string | `"raw"` | logging format |
| defaults.logging.level | string | `"debug"` | logging level |
| defaults.logging.logger | string | `"glg"` | logger name |
| defaults.observability.collector.duration | string | `"5s"` | observability collector collect duration |
| defaults.observability.collector.metrics.enable_cgo | bool | `true` | observability collector cgo metrics enabled |
| defaults.observability.collector.metrics.enable_goroutine | bool | `true` | observability collector goroutine metrics enabled |
| defaults.observability.collector.metrics.enable_memory | bool | `true` | observability collector memory metrics enabled |
| defaults.observability.collector.metrics.enable_version_info | bool | `true` | observability collector version info enabled |
| defaults.observability.enabled | bool | `false` | observability enabled |
| defaults.observability.jaeger.agent_endpoint | string | `"jaeger-agent.default.svc.cluster.local:6831"` | Jaeger agent endpoint |
| defaults.observability.jaeger.buffer_max_count | int | `10` | Jaeger buffer max count |
| defaults.observability.jaeger.collector_endpoint | string | `""` | Jaeger collector endpoint |
| defaults.observability.jaeger.enabled | bool | `false` | Jaeger exporter enabled |
| defaults.observability.jaeger.password | string | `""` | Jaeger password |
| defaults.observability.jaeger.service_name | string | `"vald"` | Jaeger service name |
| defaults.observability.jaeger.username | string | `""` | Jaeger username |
| defaults.observability.prometheus.enabled | bool | `false` | Prometheus exporter enabled |
| defaults.observability.trace.enabled | bool | `false` | trace enabled |
| defaults.observability.trace.sampling_rate | float | `1` | trace sampling rate |
| defaults.server_config.full_shutdown_duration | string | `"600s"` | server full shutdown duration |
| defaults.server_config.healths.liveness.enabled | bool | `true` | liveness server enabled |
| defaults.server_config.healths.liveness.host | string | `"0.0.0.0"` | liveness server host |
| defaults.server_config.healths.liveness.livenessProbe.failureThreshold | int | `2` |  |
| defaults.server_config.healths.liveness.livenessProbe.httpGet.path | string | `"/liveness"` | liveness probe path |
| defaults.server_config.healths.liveness.livenessProbe.httpGet.port | string | `"liveness"` | liveness probe port |
| defaults.server_config.healths.liveness.livenessProbe.httpGet.scheme | string | `"HTTP"` | liveness probe scheme |
| defaults.server_config.healths.liveness.livenessProbe.initialDelaySeconds | int | `5` |  |
| defaults.server_config.healths.liveness.livenessProbe.periodSeconds | int | `3` |  |
| defaults.server_config.healths.liveness.livenessProbe.successThreshold | int | `1` |  |
| defaults.server_config.healths.liveness.livenessProbe.timeoutSeconds | int | `2` |  |
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
| defaults.server_config.healths.readiness.server.http.shutdown_duration | string | `"5s"` | readiness server shutdown duration |
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
| discoverer.clusterRole.enabled | bool | `true` |  |
| discoverer.clusterRole.name | string | `"discoverer"` |  |
| discoverer.clusterRoleBinding.enabled | bool | `true` |  |
| discoverer.clusterRoleBinding.name | string | `"discoverer"` |  |
| discoverer.discoverer.cache_sync_duration | string | `"3s"` |  |
| discoverer.discoverer.name | string | `""` |  |
| discoverer.discoverer.namespace | string | `"_MY_POD_NAMESPACE_"` |  |
| discoverer.env[0].name | string | `"MY_POD_NAMESPACE"` |  |
| discoverer.env[0].valueFrom.fieldRef.fieldPath | string | `"metadata.namespace"` |  |
| discoverer.image.pullPolicy | string | `"Always"` |  |
| discoverer.image.repository | string | `"vdaas/vald-discoverer-k8s"` |  |
| discoverer.kind | string | `"Deployment"` |  |
| discoverer.maxReplicas | int | `2` |  |
| discoverer.maxUnavailable | string | `"50%"` |  |
| discoverer.minReplicas | int | `1` |  |
| discoverer.name | string | `"vald-discoverer"` |  |
| discoverer.observability.jaeger.service_name | string | `"vald-discoverer"` |  |
| discoverer.progressDeadlineSeconds | int | `600` |  |
| discoverer.resources.limits.cpu | string | `"600m"` |  |
| discoverer.resources.limits.memory | string | `"200Mi"` |  |
| discoverer.resources.requests.cpu | string | `"200m"` |  |
| discoverer.resources.requests.memory | string | `"65Mi"` |  |
| discoverer.revisionHistoryLimit | int | `2` |  |
| discoverer.rollingUpdate.maxSurge | string | `"25%"` |  |
| discoverer.rollingUpdate.maxUnavailable | string | `"25%"` |  |
| discoverer.server_config.full_shutdown_duration | string | `"600s"` |  |
| discoverer.server_config.healths.liveness.enabled | bool | `false` |  |
| discoverer.server_config.healths.readiness.enabled | bool | `false` |  |
| discoverer.server_config.metrics.pprof.enabled | bool | `false` |  |
| discoverer.server_config.metrics.prometheus.enabled | bool | `false` |  |
| discoverer.server_config.prefix | string | `"discoverer"` |  |
| discoverer.server_config.servers.grpc.enabled | bool | `false` |  |
| discoverer.server_config.servers.rest.enabled | bool | `false` |  |
| discoverer.server_config.tls.enabled | bool | `false` |  |
| discoverer.serviceAccount.enabled | bool | `true` |  |
| discoverer.serviceAccount.name | string | `"vald"` |  |
| discoverer.serviceType | string | `"ClusterIP"` |  |
| discoverer.terminationGracePeriodSeconds | int | `30` |  |
| discoverer.version | string | `"v0.0.0"` |  |
| gateway.annotations | list | `nil` | deployment annotations |
| gateway.env | list | `[{"name":"MY_POD_NAMESPACE","valueFrom":{"fieldRef":{"fieldPath":"metadata.namespace"}}}]` | environment variables |
| gateway.externalTrafficPolicy | string | `nil` | external traffic policy (can be specified when service type is LoadBalancer or NodePort) : Cluster or Local |
| gateway.filter.egress | list | `[""]` | egress filters |
| gateway.filter.ingress | list | `[""]` | ingress filters |
| gateway.gateway_config.agent_namespace | string | `"_MY_POD_NAMESPACE_"` | agent namespace |
| gateway.gateway_config.backup.client | object | `{}` | gRPC client for backup (overrides defaults.grpc.client) |
| gateway.gateway_config.discoverer.agent_client | object | `{}` | gRPC client for agents (overrides defaults.grpc.client) |
| gateway.gateway_config.discoverer.discover_client | object | `{}` | gRPC client for discoverer (overrides defaults.grpc.client) |
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
| gateway.image.tag | string | `nil` | image tag (overrides defaults.image.tag) |
| gateway.ingress.annotations | object | `{"nginx.ingress.kubernetes.io/grpc-backend":"true"}` | annotations for ingress |
| gateway.ingress.host | string | `"vald.gateway.vdaas.org"` | ingress hostname |
| gateway.ingress.servicePort | string | `"grpc"` | service port to be exposed by ingress |
| gateway.initContainers | list | `[{"image":"busybox","name":"wait-for-manager-compressor","sleepDuration":2,"target":"compressor","type":"wait-for"},{"image":"busybox","name":"wait-for-meta","sleepDuration":2,"target":"meta","type":"wait-for"},{"image":"busybox","name":"wait-for-discoverer","sleepDuration":2,"target":"discoverer","type":"wait-for"},{"image":"busybox","name":"wait-for-agent","sleepDuration":2,"target":"agent","type":"wait-for"}]` | init containers |
| gateway.kind | string | `"Deployment"` | deployment kind: Deployment or DaemonSet |
| gateway.maxReplicas | int | `9` | maximum number of replicas |
| gateway.maxUnavailable | string | `"50%"` | maximum number of unavailable replicas |
| gateway.minReplicas | int | `3` | minimum number of replicas |
| gateway.name | string | `"vald-gateway"` | name of vald-gateway |
| gateway.nodeName | string | `nil` | node name |
| gateway.nodeSelector | object | `nil` | node selector |
| gateway.observability | object | `{"jaeger":{"service_name":"vald-gateway"}}` | observability config (overrides defaults.observability) |
| gateway.podAnnotations | list | `nil` | pod annotations |
| gateway.progressDeadlineSeconds | int | `600` | progress deadline seconds |
| gateway.resources | object | `{"limits":{"cpu":"2000m","memory":"700Mi"},"requests":{"cpu":"200m","memory":"150Mi"}}` | compute resources |
| gateway.revisionHistoryLimit | int | `2` | number of old history to retain to allow rollback |
| gateway.rollingUpdate.maxSurge | string | `"25%"` | max surge of rolling update |
| gateway.rollingUpdate.maxUnavailable | string | `"25%"` | max unavailable of rolling update |
| gateway.server_config | object | `{"full_shutdown_duration":"600s","healths":{"liveness":{"enabled":false},"readiness":{"enabled":false}},"metrics":{"pprof":{"enabled":false},"prometheus":{"enabled":false}},"prefix":"gateway","servers":{"grpc":{"enabled":false},"rest":{"enabled":false}},"tls":{"enabled":false}}` | server config (overrides defaults.server_config) |
| gateway.service.annotations | list | `nil` | service annotations |
| gateway.service.labels | list | `nil` | service labels |
| gateway.serviceType | string | `"ClusterIP"` | service type: ClusterIP, LoadBalancer or NodePort |
| gateway.terminationGracePeriodSeconds | int | `30` | duration in seconds pod needs to terminate gracefully |
| gateway.version | string | `"v0.0.0"` | version of gateway config |
| gateway.volumeMounts | list | `nil` | volume mounts |
| gateway.volumes | list | `nil` | volumes |
| indexManager.env[0].name | string | `"MY_POD_NAMESPACE"` |  |
| indexManager.env[0].valueFrom.fieldRef.fieldPath | string | `"metadata.namespace"` |  |
| indexManager.image.pullPolicy | string | `"Always"` |  |
| indexManager.image.repository | string | `"vdaas/vald-manager-index"` |  |
| indexManager.indexer.agent_namespace | string | `"_MY_POD_NAMESPACE_"` |  |
| indexManager.indexer.auto_index_check_duration | string | `"1m"` |  |
| indexManager.indexer.auto_index_duration_limit | string | `"30m"` |  |
| indexManager.indexer.auto_index_length | int | `100` |  |
| indexManager.indexer.concurrency | int | `1` |  |
| indexManager.indexer.discoverer.agent_client.dial_option.tcp.dialer.keep_alive | string | `"15m"` |  |
| indexManager.indexer.discoverer.discover_client | object | `{}` |  |
| indexManager.indexer.discoverer.duration | string | `"500ms"` |  |
| indexManager.indexer.node_name | string | `""` |  |
| indexManager.initContainers[0].image | string | `"busybox"` |  |
| indexManager.initContainers[0].name | string | `"wait-for-agent"` |  |
| indexManager.initContainers[0].sleepDuration | int | `2` |  |
| indexManager.initContainers[0].target | string | `"agent"` |  |
| indexManager.initContainers[0].type | string | `"wait-for"` |  |
| indexManager.initContainers[1].image | string | `"busybox"` |  |
| indexManager.initContainers[1].name | string | `"wait-for-discoverer"` |  |
| indexManager.initContainers[1].sleepDuration | int | `2` |  |
| indexManager.initContainers[1].target | string | `"discoverer"` |  |
| indexManager.initContainers[1].type | string | `"wait-for"` |  |
| indexManager.kind | string | `"Deployment"` |  |
| indexManager.maxUnavailable | string | `"50%"` |  |
| indexManager.name | string | `"vald-manager-index"` |  |
| indexManager.observability.jaeger.service_name | string | `"vald-manager-index"` |  |
| indexManager.progressDeadlineSeconds | int | `600` |  |
| indexManager.replicas | int | `1` |  |
| indexManager.resources.limits.cpu | int | `1` |  |
| indexManager.resources.limits.memory | string | `"500Mi"` |  |
| indexManager.resources.requests.cpu | string | `"200m"` |  |
| indexManager.resources.requests.memory | string | `"80Mi"` |  |
| indexManager.revisionHistoryLimit | int | `2` |  |
| indexManager.rollingUpdate.maxSurge | string | `"25%"` |  |
| indexManager.rollingUpdate.maxUnavailable | string | `"25%"` |  |
| indexManager.server_config.full_shutdown_duration | string | `"600s"` |  |
| indexManager.server_config.healths.liveness.enabled | bool | `false` |  |
| indexManager.server_config.healths.readiness.enabled | bool | `false` |  |
| indexManager.server_config.metrics.pprof.enabled | bool | `false` |  |
| indexManager.server_config.metrics.prometheus.enabled | bool | `false` |  |
| indexManager.server_config.prefix | string | `"index-manager"` |  |
| indexManager.server_config.servers.grpc.enabled | bool | `false` |  |
| indexManager.server_config.servers.rest.enabled | bool | `false` |  |
| indexManager.server_config.tls.enabled | bool | `false` |  |
| indexManager.serviceType | string | `"ClusterIP"` |  |
| indexManager.terminationGracePeriodSeconds | int | `30` |  |
| indexManager.version | string | `"v0.0.0"` |  |
| initializer.cassandra.configmap.backup.enabled | bool | `true` |  |
| initializer.cassandra.configmap.backup.name | string | `"meta_vector"` |  |
| initializer.cassandra.configmap.enabled | bool | `false` |  |
| initializer.cassandra.configmap.filename | string | `"init.cql"` |  |
| initializer.cassandra.configmap.keyspace | string | `"vald"` |  |
| initializer.cassandra.configmap.meta.enabled | bool | `true` |  |
| initializer.cassandra.configmap.meta.name.kv | string | `"kv"` |  |
| initializer.cassandra.configmap.meta.name.vk | string | `"vk"` |  |
| initializer.cassandra.configmap.name | string | `"cassandra-initdb"` |  |
| initializer.cassandra.configmap.replication_class | string | `"SimpleStrategy"` |  |
| initializer.cassandra.configmap.replication_factor | int | `3` |  |
| initializer.cassandra.enabled | bool | `false` |  |
| initializer.cassandra.env[0].name | string | `"CASSANDRA_HOST"` |  |
| initializer.cassandra.env[0].value | string | `"cassandra.default.svc.cluster.local"` |  |
| initializer.cassandra.env[1].name | string | `"CASSANDRA_USER"` |  |
| initializer.cassandra.env[1].value | string | `"root"` |  |
| initializer.cassandra.env[2].name | string | `"CASSANDRA_PASSWORD"` |  |
| initializer.cassandra.env[2].valueFrom.secretKeyRef.key | string | `"password"` |  |
| initializer.cassandra.env[2].valueFrom.secretKeyRef.name | string | `"cassandra-secret"` |  |
| initializer.cassandra.image.pullPolicy | string | `"Always"` |  |
| initializer.cassandra.image.repository | string | `"cassandra"` |  |
| initializer.cassandra.image.tag | string | `"latest"` |  |
| initializer.cassandra.name | string | `"cassandra-init"` |  |
| initializer.cassandra.restartPolicy | string | `"Never"` |  |
| initializer.cassandra.secret.data.password | string | `"cGFzc3dvcmQ="` |  |
| initializer.cassandra.secret.enabled | bool | `false` |  |
| initializer.cassandra.secret.name | string | `"cassandra-secret"` |  |
| initializer.mysql.configmap.enabled | bool | `false` |  |
| initializer.mysql.configmap.filename | string | `"ddl.sql"` |  |
| initializer.mysql.configmap.name | string | `"mysql-config"` |  |
| initializer.mysql.configmap.schema | string | `"vald"` |  |
| initializer.mysql.enabled | bool | `false` |  |
| initializer.mysql.env[0].name | string | `"MYSQL_HOST"` |  |
| initializer.mysql.env[0].value | string | `"mysql.default.svc.cluster.local"` |  |
| initializer.mysql.env[1].name | string | `"MYSQL_USER"` |  |
| initializer.mysql.env[1].value | string | `"root"` |  |
| initializer.mysql.env[2].name | string | `"MYSQL_PASSWORD"` |  |
| initializer.mysql.env[2].valueFrom.secretKeyRef.key | string | `"password"` |  |
| initializer.mysql.env[2].valueFrom.secretKeyRef.name | string | `"mysql-secret"` |  |
| initializer.mysql.image.pullPolicy | string | `"Always"` |  |
| initializer.mysql.image.repository | string | `"mysql"` |  |
| initializer.mysql.image.tag | string | `"latest"` |  |
| initializer.mysql.name | string | `"mysql-init"` |  |
| initializer.mysql.restartPolicy | string | `"Never"` |  |
| initializer.mysql.secret.data.password | string | `"cGFzc3dvcmQ="` |  |
| initializer.mysql.secret.enabled | bool | `false` |  |
| initializer.mysql.secret.name | string | `"mysql-secret"` |  |
| initializer.redis.enabled | bool | `false` |  |
| initializer.redis.env[0].name | string | `"REDIS_HOST"` |  |
| initializer.redis.env[0].value | string | `"redis.default.svc.cluster.local"` |  |
| initializer.redis.env[1].name | string | `"REDIS_PASSWORD"` |  |
| initializer.redis.env[1].valueFrom.secretKeyRef.key | string | `"password"` |  |
| initializer.redis.env[1].valueFrom.secretKeyRef.name | string | `"redis-secret"` |  |
| initializer.redis.image.pullPolicy | string | `"Always"` |  |
| initializer.redis.image.repository | string | `"redis"` |  |
| initializer.redis.image.tag | string | `"latest"` |  |
| initializer.redis.name | string | `"redis-init"` |  |
| initializer.redis.restartPolicy | string | `"Never"` |  |
| initializer.redis.secret.data.password | string | `"cGFzc3dvcmQ="` |  |
| initializer.redis.secret.enabled | bool | `false` |  |
| initializer.redis.secret.name | string | `"redis-secret"` |  |
| meta.cassandra.config.connect_timeout | string | `"600ms"` |  |
| meta.cassandra.config.consistency | string | `"quorum"` |  |
| meta.cassandra.config.cql_version | string | `"3.0.0"` |  |
| meta.cassandra.config.default_idempotence | bool | `false` |  |
| meta.cassandra.config.default_timestamp | bool | `true` |  |
| meta.cassandra.config.disable_initial_host_lookup | bool | `false` |  |
| meta.cassandra.config.disable_node_status_events | bool | `false` |  |
| meta.cassandra.config.disable_skip_metadata | bool | `false` |  |
| meta.cassandra.config.disable_topology_events | bool | `false` |  |
| meta.cassandra.config.enable_host_verification | bool | `false` |  |
| meta.cassandra.config.hosts[0] | string | `"cassandra-0.cassandra.default.svc.cluster.local"` |  |
| meta.cassandra.config.hosts[1] | string | `"cassandra-1.cassandra.default.svc.cluster.local"` |  |
| meta.cassandra.config.hosts[2] | string | `"cassandra-2.cassandra.default.svc.cluster.local"` |  |
| meta.cassandra.config.ignore_peer_addr | bool | `false` |  |
| meta.cassandra.config.keyspace | string | `"vald"` |  |
| meta.cassandra.config.kv_table | string | `"kv"` |  |
| meta.cassandra.config.max_prepared_stmts | int | `1000` |  |
| meta.cassandra.config.max_routing_key_info | int | `1000` |  |
| meta.cassandra.config.max_wait_schema_agreement | string | `"1m"` |  |
| meta.cassandra.config.num_conns | int | `2` |  |
| meta.cassandra.config.page_size | int | `5000` |  |
| meta.cassandra.config.password | string | `"_CASSANDRA_PASSWORD_"` |  |
| meta.cassandra.config.pool_config.data_center | string | `""` |  |
| meta.cassandra.config.pool_config.dc_aware_routing | bool | `false` |  |
| meta.cassandra.config.pool_config.non_local_replicas_fallback | bool | `false` |  |
| meta.cassandra.config.pool_config.shuffle_replicas | bool | `false` |  |
| meta.cassandra.config.port | int | `9042` |  |
| meta.cassandra.config.proto_version | int | `0` |  |
| meta.cassandra.config.reconnect_interval | string | `"1m"` |  |
| meta.cassandra.config.reconnection_policy.initial_interval | string | `"1m"` |  |
| meta.cassandra.config.reconnection_policy.max_retries | int | `3` |  |
| meta.cassandra.config.retry_policy.max_duration | string | `"30s"` |  |
| meta.cassandra.config.retry_policy.min_duration | string | `"1s"` |  |
| meta.cassandra.config.retry_policy.num_retries | int | `3` |  |
| meta.cassandra.config.socket_keepalive | string | `"0s"` |  |
| meta.cassandra.config.tcp.dialer.dual_stack_enabled | bool | `false` |  |
| meta.cassandra.config.tcp.dialer.keep_alive | string | `"10m"` |  |
| meta.cassandra.config.tcp.dialer.timeout | string | `"30s"` |  |
| meta.cassandra.config.tcp.dns.cache_enabled | bool | `true` |  |
| meta.cassandra.config.tcp.dns.cache_expiration | string | `"24h"` |  |
| meta.cassandra.config.tcp.dns.refresh_duration | string | `"5m"` |  |
| meta.cassandra.config.timeout | string | `"600ms"` |  |
| meta.cassandra.config.tls.ca | string | `"/path/to/ca"` |  |
| meta.cassandra.config.tls.cert | string | `"/path/to/cert"` |  |
| meta.cassandra.config.tls.enabled | bool | `false` |  |
| meta.cassandra.config.tls.key | string | `"/path/to/key"` |  |
| meta.cassandra.config.username | string | `"root"` |  |
| meta.cassandra.config.vk_table | string | `"vk"` |  |
| meta.cassandra.config.write_coalesce_wait_time | string | `"200ms"` |  |
| meta.cassandra.enabled | bool | `false` |  |
| meta.env[0].name | string | `"REDIS_PASSWORD"` |  |
| meta.env[0].valueFrom.secretKeyRef.key | string | `"password"` |  |
| meta.env[0].valueFrom.secretKeyRef.name | string | `"redis-secret"` |  |
| meta.hpa.enabled | bool | `true` |  |
| meta.hpa.targetCPUUtilizationPercentage | int | `80` |  |
| meta.image.pullPolicy | string | `"Always"` |  |
| meta.image.repository | string | `"vdaas/vald-meta-redis"` |  |
| meta.initContainers[0].env[0].name | string | `"REDIS_PASSWORD"` |  |
| meta.initContainers[0].env[0].valueFrom.secretKeyRef.key | string | `"password"` |  |
| meta.initContainers[0].env[0].valueFrom.secretKeyRef.name | string | `"redis-secret"` |  |
| meta.initContainers[0].image | string | `"redis:latest"` |  |
| meta.initContainers[0].name | string | `"wait-for-redis"` |  |
| meta.initContainers[0].redis.hosts[0] | string | `"redis.default.svc.cluster.local"` |  |
| meta.initContainers[0].redis.options[0] | string | `"-a ${REDIS_PASSWORD}"` |  |
| meta.initContainers[0].sleepDuration | int | `2` |  |
| meta.initContainers[0].type | string | `"wait-for-redis"` |  |
| meta.kind | string | `"Deployment"` |  |
| meta.maxReplicas | int | `10` |  |
| meta.maxUnavailable | string | `"50%"` |  |
| meta.minReplicas | int | `2` |  |
| meta.name | string | `"vald-meta"` |  |
| meta.observability.jaeger.service_name | string | `"vald-meta"` |  |
| meta.progressDeadlineSeconds | int | `600` |  |
| meta.redis.config.addrs[0] | string | `"redis.default.svc.cluster.local:6379"` |  |
| meta.redis.config.db | int | `0` |  |
| meta.redis.config.dial_timeout | string | `"5s"` |  |
| meta.redis.config.idle_check_frequency | string | `"1m"` |  |
| meta.redis.config.idle_timeout | string | `"5m"` |  |
| meta.redis.config.key_pref | string | `""` |  |
| meta.redis.config.kv_prefix | string | `""` |  |
| meta.redis.config.max_conn_age | string | `"0s"` |  |
| meta.redis.config.max_redirects | int | `3` |  |
| meta.redis.config.max_retries | int | `0` |  |
| meta.redis.config.max_retry_backoff | string | `"512ms"` |  |
| meta.redis.config.min_idle_conns | int | `0` |  |
| meta.redis.config.min_retry_backoff | string | `"8ms"` |  |
| meta.redis.config.password | string | `"_REDIS_PASSWORD_"` |  |
| meta.redis.config.pool_size | int | `10` |  |
| meta.redis.config.pool_timeout | string | `"4s"` |  |
| meta.redis.config.prefix_delimiter | string | `""` |  |
| meta.redis.config.read_only | bool | `false` |  |
| meta.redis.config.read_timeout | string | `"3s"` |  |
| meta.redis.config.route_by_latency | bool | `false` |  |
| meta.redis.config.route_randomly | bool | `true` |  |
| meta.redis.config.tcp.dialer.dual_stack_enabled | bool | `false` |  |
| meta.redis.config.tcp.dialer.keep_alive | string | `"5m"` |  |
| meta.redis.config.tcp.dialer.timeout | string | `"5s"` |  |
| meta.redis.config.tcp.dns.cache_enabled | bool | `true` |  |
| meta.redis.config.tcp.dns.cache_expiration | string | `"24h"` |  |
| meta.redis.config.tcp.dns.refresh_duration | string | `"1h"` |  |
| meta.redis.config.tcp.tls.enabled | bool | `false` |  |
| meta.redis.config.tls.enabled | bool | `false` |  |
| meta.redis.config.vk_prefix | string | `""` |  |
| meta.redis.config.write_timeout | string | `"3s"` |  |
| meta.redis.enabled | bool | `true` |  |
| meta.resources.limits.cpu | string | `"300m"` |  |
| meta.resources.limits.memory | string | `"100Mi"` |  |
| meta.resources.requests.cpu | string | `"100m"` |  |
| meta.resources.requests.memory | string | `"40Mi"` |  |
| meta.revisionHistoryLimit | int | `2` |  |
| meta.rollingUpdate.maxSurge | string | `"25%"` |  |
| meta.rollingUpdate.maxUnavailable | string | `"25%"` |  |
| meta.server_config.full_shutdown_duration | string | `"600s"` |  |
| meta.server_config.healths.liveness.enabled | bool | `false` |  |
| meta.server_config.healths.readiness.enabled | bool | `false` |  |
| meta.server_config.metrics.pprof.enabled | bool | `false` |  |
| meta.server_config.metrics.prometheus.enabled | bool | `false` |  |
| meta.server_config.prefix | string | `"meta"` |  |
| meta.server_config.servers.grpc.enabled | bool | `false` |  |
| meta.server_config.servers.rest.enabled | bool | `false` |  |
| meta.server_config.tls.enabled | bool | `false` |  |
| meta.serviceType | string | `"ClusterIP"` |  |
| meta.terminationGracePeriodSeconds | int | `30` |  |
| meta.version | string | `"v0.0.0"` |  |
