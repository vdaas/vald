# Design Overview

This document explains what the `valdoperatorrelease` controller does and how it works.
It is grounded in the current code; see the referenced files for the source of truth.

## What it is

A Kubernetes controller that manages the lifecycle of [Vald](https://vald.vdaas.org/)
clusters through a single custom resource, `ValdOperatorRelease` (short name `vor`,
API group `vald.vdaas.org/v1`).

From a small, high-level input the controller generates one or more
[`ValdRelease`](https://github.com/vdaas/vald-helm-operator) (VRS) resources — the
manifests that the Vald Helm Operator (VHO) reconciles into a running Vald cluster.

```
ValdOperatorRelease (vor)  ──reconcile──▶  ValdRelease (VRS) × (active infra × clusters)
   minimal input                       consumed by vald-helm-operator
```

The controller only *generates* VRS definitions. It does not run Vald itself; a VHO
running in the target cluster turns each VRS into Vald pods. This separation lets a
management cluster emit VRS definitions and distribute them to other clusters (for
example via an external workflow engine).

## The `ValdOperatorRelease` (VOR) resource

The goal of VOR is to collapse the large `ValdRelease` configuration surface into the
minimum a user must supply. There are two input groups.

1. **Infrastructure / node-pool information** — `spec.infrastructure[]`
2. **Minimal Vald settings** — `spec.vectorEngine.vald`

Source of truth: `api/v1/valdoperatorrelease_types.go`.

### `spec.infrastructure[]`

A list of infrastructure entries. Each entry produces VRS objects for its clusters.

| Field | Notes |
|-------|-------|
| `role` | Free-form role label (e.g. `hot`, `standby`, `blue`, `green`). Copied to VRS labels. |
| `type` | Cluster type label. |
| `active` | When `false`, the entry is skipped during VRS generation. |
| `clusters[]` | `{ id, name }`. `id` is typically filled in by an external system once the target cluster exists; `name` must be set. |
| `nodePools` | Map keyed by pool type: `general` (required) and `agent` (optional). Each pool carries `name`, `replicas`, and `machineResource{ cpu, memory, storage }`. |

The `machineResource` values drive per-component resource and replica calculations for
the generated VRS.

### `spec.vectorEngine.vald`

The minimal Vald configuration reflected into each VRS:

- `defaults.logLevel`
- `agent.ngt` — `dimension` (required, ≥2), `creationEdgeSize`, `searchEdgeSize`, `distanceType`, `objectType`
- `agent.persistentVolume` — `{ enabled, storageClass, accessMode }` (storageClass/accessMode fall back to env defaults when omitted)
- `indexer` — `indexSchedule`, `saveSchedule`, `concurrency`, `manager`, suspend flags, durations
- `gateway` — `indexReplica`, `serviceType`, `ingress{ enabled, host }`
- `discoverer.kind` — `DaemonSet` or `Deployment`
- `overlay` — an arbitrary JSON patch merged onto the generated VRS (see *Overlay* below)

## Lifecycle

Reconciliation is modeled as an ordered flow of phases. Each phase carries a `Condition`
plus an optional `Builder` (creates/updates resources) and `Checker` (reports readiness).

Phases (`internal/pkg/lifecycle/condition.go`, `internal/pkg/domain/valdoperatorrelease/lifecycle.go`):

| Phase (`status.phase`) | Builder | Checker | Purpose |
|------------------------|---------|---------|---------|
| `WaitingClusterCreate` | – | yes | Validate infra config; wait until every cluster has `id` and `name`. |
| `WaitingCreateVrs` | yes | yes | Build the VRS objects and wait until they are ready. |
| `Completed` | – | – | Terminal phase. |

### Accumulating, self-healing conditions

`status.conditions` **accumulates** as phases progress (conditions are never removed),
and every reconcile re-evaluates **all** conditions. This is intentional: if a
previously-`True` condition later breaks (e.g. a generated VRS is deleted, or the spec
changes), the controller detects it and restarts work from that phase. `status.phase`
tracks the condition currently being processed.

See `internal/controller/valdoperatorrelease_controller.go` (`reconcileRoutine` /
`reconcileCondition`) and `internal/pkg/domain/valdoperatorrelease/phase.go`
(`AdvanceToNextPhase`).

### Readiness result states

A `Checker` returns one of four `desired.Result` states
(`internal/pkg/lifecycle/desired/result.go`):

| Result | Condition status | Meaning |
|--------|------------------|---------|
| `Progressing("")` | Unknown / Progressing | The controller is actively working. |
| `Pending("msg")` | Unknown / Pending | Waiting on something external (e.g. external system assigning a cluster id). |
| `Succeeded()` | True / Succeeded | Done. |
| `Failed(err)` | False / Failed | Misconfiguration or unrecoverable error. |

For the cluster-create check: a missing `cluster.id`/`cluster.name` yields **Pending**
(external system not done yet), while empty `clusters` or empty `infrastructure` yields
**Failed** (misconfiguration). See
`internal/pkg/domain/valdoperatorrelease/condition_wait_for_cluster_create.go`.

## VRS generation flow

Implemented in `internal/pkg/lifecycle/builder/vald/`. `VrsBuilder.Build` is a pure
function of `(CR, Config, NodePoolCapability)` — it makes no Kubernetes API calls.

1. **Validate** the CR (`infrastructure` non-empty, each cluster has a `name`).
2. **Iterate `infrastructure[]`**, skipping entries where `active == false` (and, when
   node-pool matching is enabled, entries with no matching general pool).
3. **Assemble the VRS spec** from the CR:
   - `buildDefaults` — log level etc. (`defaults.go`)
   - `buildGateway` — `indexReplica`, service type, ingress (`gateway.go`)
   - `buildAgent` — NGT settings (`agent.go`)
   - `buildManager` — Manager mode when `indexer.manager == true`, otherwise
     Creator/Saver mode (`manager.go`)
   - `buildDiscoverer` — kind + namespaced RBAC names (`discoverer.go`)
4. **Resolve resources from node pools** — `SetRelationalResources` derives replicas and
   per-component CPU/memory from the resolved agent node pool. The general-pool fallback
   rule lives in the domain layer (`Rules.ResolveAgentNodePool`).
5. **Optional settings** — persistent volume (`persistent_volume.go`) and node
   affinities (`common.go`).
6. **Per cluster** — for each `infra.clusters[]`, set `name = <namespace>-<clusterName>`
   (truncated to 63 chars), apply labels (`namespace`, `type`, `role`), then merge the
   overlay. One VRS is produced per cluster, so a single VOR can yield many VRS objects.

### Applying generated resources (`ResourceSyncer`)

The controller delegates the write side to `ResourceSyncer`
(`internal/controller/resource_syncer.go`):

- `Build` the desired objects, then `CreateOrUpdate` each with an owner reference and a
  `managed-generation` label set to the owner's `.metadata.generation`.
- **Prune**: list owned resources of the same kind and delete any that the current
  `Build` no longer produces (stale generations).

### Overlay

`spec.vectorEngine.vald.overlay` is a JSON patch merged onto the generated VRS, layered
on top of the default VRS template loaded at startup (`DEFAULT_VRS_PATH`). See
`internal/pkg/lifecycle/builder/vald/kustomize.go`.

## Multi-cluster distribution and node-pool matching

The same VOR can be distributed to multiple clusters, generating a VRS only where a
matching node pool exists.

- With `REQUIRE_NODEPOOL_MATCH=true`, the controller lists Nodes and resolves a
  `NodePoolCapability`; a VRS is generated only when a `general` pool is present
  (`internal/pkg/lifecycle/builder/vald/capability.go`).
- Node match is by labels `namespace=<vor namespace>` and `type=general` (required);
  `type=agent` is optional — when present, agent pods are scheduled onto `type=agent`,
  otherwise they fall back to `type=general`.
- Label keys default to `namespace` / `type`; set `NODEPOOL_LABEL_PREFIX`
  (e.g. `vald.vdaas.org`) to namespace them.
- With `REQUIRE_NODEPOOL_MATCH=false` (default) the controller skips node listing and
  treats every infra entry as schedulable (`AlwaysAvailable`).

## Configuration (environment variables)

Loaded once at startup into `config.Config` (`internal/infrastructure/config/`).

| Env var | Default | Purpose |
|---------|---------|---------|
| `DEFAULT_VRS_PATH` | `/opt/valdoperatorrelease/config/vrs.yaml` | Default VRS template merged with the overlay. |
| `REQUIRE_NODEPOOL_MATCH` | `false` | Only generate VRS where matching node pools exist. |
| `NODEPOOL_LABEL_PREFIX` | `""` | Prefix for the `namespace`/`type`/`role` node labels. |
| `AGENT_PODS_PER_NODE` | `2` | Agent pods packed per node when computing replicas. |
| `DEFAULT_STORAGE_CLASS` | `standard` | PV storage class fallback. |
| `DEFAULT_ACCESS_MODE` | `ReadWriteOnce` | PV access mode fallback. |
| `PV_BUFFER_RATIO` | `1.5` | PV size = `max(memoryRequest * ratio, min)`. |
| `PV_MIN_SIZE_BYTES` | `1073741824` | Minimum PV size in bytes. |
| `ENABLE_INGRESS` | `true` | Enable gateway ingress generation. |
| `GATEWAY_INGRESS_ANNOTATIONS` | `""` | YAML map of annotations applied to the gateway ingress. |
| `GATEWAY_SERVICE_TYPE` | `NodePort` | Gateway service type (`NodePort` / `ClusterIP` / `LoadBalancer`). |
| `VRS_LOG_LEVEL` | `warn` | Log level passed through to the generated VRS. |
| `DISCOVERER_DS_MAX_SURGE` | `30%` | Discoverer DaemonSet rolling-update `maxSurge`. |
| `DISCOVERER_DS_MAX_UNAVAILABLE` | `0%` | Discoverer DaemonSet rolling-update `maxUnavailable`. |
| `INTERNAL_HOST_DOMAIN` | `""` | Optional internal host domain (empty by default for OSS). |

## Source map

| Concern | Location |
|---------|----------|
| CRD types | `api/v1/valdoperatorrelease_types.go` |
| Reconciler | `internal/controller/valdoperatorrelease_controller.go` |
| Resource apply / prune | `internal/controller/resource_syncer.go` |
| Lifecycle / conditions | `internal/pkg/lifecycle/`, `internal/pkg/domain/valdoperatorrelease/` |
| VRS builder | `internal/pkg/lifecycle/builder/vald/` |
| VRS API model | `internal/pkg/api/valdrelease/` |
| Config / env | `internal/infrastructure/config/`, `internal/infrastructure/env/` |

> Note: there is currently no admission webhook implementation in this repository. The
> webhook server scaffolding in `cmd/main.go` and the `webhooks` entry in `PROJECT` are
> kubebuilder leftovers; webhooks are tracked as future work.
