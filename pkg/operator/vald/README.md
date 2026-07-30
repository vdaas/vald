# vald-operator

The vald-operator reconciles `ValdOperatorRelease` custom resources into `ValdRelease`
resources. It follows the standard Vald component layout
(`cmd/operator/vald` + `pkg/operator/vald/{config,service,usecase}`)
and is configured through a single YAML file (see
`cmd/operator/vald/sample.yaml` and `k8s/operator/vald/configmap.yaml`).

## Package layout

- `pkg/operator/vald/config`: YAML schema (`operator:` section) and the runtime
  `Config` derived from it.
- `pkg/operator/vald/service`: the controller itself — reconcile pipeline
  (`reconciler.go`: fetch → capability → phase loop (build → sync → check) →
  status update), the ValdRelease builder (`builder*.go`), node-pool rules
  (`rules.go`, `capability.go`) and the manager wiring (`operator.go`).
- `pkg/operator/vald/usecase`: runner assembly (servers, observability,
  service).

Kubernetes-domain code independent of the operator business logic lives under
`internal/k8s`:

- `internal/k8s/vald/operator/api/{v1,valdrelease,metadata}`: the ValdOperatorRelease
  and ValdRelease CRD models (same placement convention as
  `internal/k8s/vald/{benchmark,mirror}/api/v1`).
- `internal/k8s/resource`: generic object helpers, `Syncer` (owner-scoped
  CreateOrUpdate + prune) and `metav1.Condition` upsert helpers.
- `internal/k8s/kustomize`: unstructured overlay merge via an in-memory
  kustomize run.

## Responsibility split

- **Default VRS template** (`vrs.default_vrs_path`, distributed as a ConfigMap):
  controls the _contents_ of the generated ValdRelease (component images,
  resources, and any field the operator does not compute itself).
- **Operator configuration** (`operator:` section, described below): controls
  the values the operator logic computes or overrides, and the behavior of the
  Kubernetes controller manager itself.

Values the builder injects (agent rolling-update, gateway HPA/ingress, indexer
duration limits, logging) override the template, which is why they are exposed
here instead of in the template.

## Configuration reference (`operator:` section)

### Top level

| Key         | Type   | Default | Description                                            |
| ----------- | ------ | ------- | ------------------------------------------------------ |
| `name`      | string | `""`    | Controller name (used as the manager controller name). |
| `namespace` | string | `""`    | Namespace the operator is deployed in.                 |

### `controller` — manager / reconcile loop

| Key                              | Type     | Default         | Description                                                                                         |
| -------------------------------- | -------- | --------------- | --------------------------------------------------------------------------------------------------- |
| `leader_election.enabled`        | bool     | `false`         | Enables leader election.                                                                            |
| `leader_election.id`             | string   | `""`            | Name of the leader election lease resource (required when enabled).                                 |
| `leader_election.namespace`      | string   | `""`            | Namespace of the lease resource.                                                                    |
| `leader_election.lease_duration` | duration | `""` (15s)      | How long non-leader candidates wait before acquiring an expired lease.                              |
| `leader_election.renew_deadline` | duration | `""` (10s)      | How long the leader retries refreshing its lease before giving up.                                  |
| `leader_election.retry_period`   | duration | `""` (2s)       | Wait between leader election actions.                                                               |
| `metrics_address`                | string   | `""` (`:8080`)  | Bind address of the manager metrics server. `"0"` disables it.                                      |
| `max_concurrent_reconciles`      | int      | `0` (1)         | Reconcile worker count of the ValdOperatorRelease controller.                                       |
| `sync_period`                    | duration | `""` (~10h)     | Informer cache resync interval.                                                                     |
| `cache_namespaces`               | []string | `[]` (cluster)  | Restricts the manager cache (watches/lists) to these namespaces.                                    |
| `requeue.success`                | duration | `""` (disabled) | Periodic requeue after a successful reconcile.                                                      |
| `requeue.on_error`               | duration | `""` (backoff)  | Fixed retry interval on reconcile failure; replaces the exponential backoff.                        |
| `requeue.not_found`              | duration | `""` (disabled) | Retry interval when the reconciled object is not found. Note: deleted CRs keep requeuing while set. |

### `vrs` — ValdRelease rendering

| Key                                      | Type   | Default              | Description                                                       |
| ---------------------------------------- | ------ | -------------------- | ----------------------------------------------------------------- |
| `default_vrs_path`                       | string | `""`                 | Path of the default ValdRelease template yaml.                    |
| `log_level`                              | string | `""`                 | Log level passed through to the underlying vald deployment.       |
| `log_format`                             | string | `raw`                | Log format passed through to the underlying vald deployment.      |
| `logger`                                 | string | `glg`                | Logger passed through to the underlying vald deployment.          |
| `managed_generation_label`               | string | `managed-generation` | Label key recording the owner generation on managed resources.    |
| `agent.max_surge`                        | string | `"1"`                | Agent StatefulSet rolling-update max surge (intstr string).       |
| `agent.max_unavailable`                  | string | `"1"`                | Agent StatefulSet rolling-update max unavailable (intstr string). |
| `agent.enable_in_memory_mode`            | bool   | `true`               | Enables the NGT in-memory mode.                                   |
| `gateway.hpa_target_cpu_utilization`     | int    | `80`                 | Gateway lb HPA target CPU utilization percentage.                 |
| `gateway.ingress_service_port`           | string | `grpc`               | Service port name referenced by the gateway ingress.              |
| `gateway.ingress_path_type`              | string | `Prefix`             | Ingress path type: `Prefix`, `Exact` or `ImplementationSpecific`. |
| `indexer.auto_index_duration_limit`      | string | `1h`                 | Index manager auto index duration limit.                          |
| `indexer.auto_save_index_duration_limit` | string | `-1h`                | Index manager auto save index duration limit.                     |

### `node_pool` — node-pool matching

| Key                   | Type   | Default | Description                                                         |
| --------------------- | ------ | ------- | ------------------------------------------------------------------- |
| `require_match`       | bool   | `false` | Generates VRS only when matching nodepools exist in the cluster.    |
| `label_prefix`        | string | `""`    | Optional prefix for nodepool-related label keys (`vald.vdaas.org`). |
| `agent_pods_per_node` | int    | `0`     | Number of agent pods per node.                                      |

### `persistent_volume` — Agent.PersistentVolume fallback and sizing

| Key                     | Type    | Default | Description                                                     |
| ----------------------- | ------- | ------- | --------------------------------------------------------------- |
| `default_storage_class` | string  | `""`    | Fallback StorageClass when the CR omits it.                     |
| `default_access_mode`   | string  | `""`    | Fallback AccessMode when the CR omits it.                       |
| `buffer_ratio`          | float64 | `0`     | PV sizing: `max(memoryRequest * buffer_ratio, min_size_bytes)`. |
| `min_size_bytes`        | int64   | `0`     | Minimum PV size in bytes.                                       |

### `networking` — gateway / discoverer networking

| Key                                    | Type              | Default         | Description                                                             |
| -------------------------------------- | ----------------- | --------------- | ----------------------------------------------------------------------- |
| `enable_ingress`                       | bool              | `false`         | Enables the ingress for the gateway.                                    |
| `gateway_ingress_annotations`          | map[string]string | `{}`            | Annotations applied to the gateway ingress.                             |
| `gateway_service_type`                 | string            | `""` (NodePort) | Service type of the gateway: `NodePort`, `ClusterIP` or `LoadBalancer`. |
| `discoverer_daemonset_max_surge`       | string            | `""`            | Discoverer DaemonSet rolling-update max surge (intstr percent).         |
| `discoverer_daemonset_max_unavailable` | string            | `""`            | Discoverer DaemonSet rolling-update max unavailable.                    |

The gateway fields resolve with the priority **CR spec > operator config > hardcoded default**:
the service type falls back from the CR `serviceType` to `gateway_service_type` and then `NodePort`;
the ingress renders only when both `enable_ingress` and the CR spec enable it; and
`gateway_ingress_annotations` seed the ingress annotations, with CR overlay entries winning per key.

Duration values use Go `time.ParseDuration` syntax (`100ms`, `1s`, `10h`).
Every string value supports `_ENV_NAME_` environment variable expansion.

## Non-configurable constants

The following are contracts, not tunables, and stay hardcoded on purpose:

- API group/version (`vald.vdaas.org/v1`) and controller name (`valdoperatorrelease`).
- Nodepool label suffixes `namespace` / `type` / `role` (the prefix is
  configurable via `node_pool.label_prefix`).
- `app.kubernetes.io/managed-by` + `<group>/managed-resource` sub-resource labels.
- Kubernetes DNS label length limit (63) for generated resource names.
- Discoverer ClusterRole/ClusterRoleBinding names (derived from the CR namespace).
