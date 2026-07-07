# ValdOperatorRelease Controller Specification

## Overview

The ValdOperatorRelease controller manages the lifecycle of [ValdRelease](https://vald.vdaas.org) resources
on a Kubernetes cluster. Users describe their desired Vald configuration via a single `ValdOperatorRelease`
CR; the controller generates and maintains the corresponding `ValdRelease` objects, adjusting them
whenever the CR changes.

```
User applies ValdOperatorRelease CR
        │
        ▼
┌───────────────────────┐
│  ValdOperatorRelease          │
│  Controller            │
│                        │
│  1. Check infrastructure│
│  2. Build ValdRelease  │
│  3. Create / update    │
│  4. Prune stale VRS    │
└───────────────────────┘
        │
        ▼
  ValdRelease CR(s)
  (consumed by VHO)
```

ValdRelease objects are consumed by the [Vald Helm Operator (VHO)](https://vald.vdaas.org/docs/operator/vald-helm-operator/),
which deploys the actual Vald components.

---

## CRD: ValdOperatorRelease (`vor`)

### Spec

```yaml
spec:
  infrastructure:
    - role: green           # arbitrary role label (e.g. green/blue, hot/cold)
      type: kind            # cluster type hint (informational)
      active: true          # whether this infra entry is currently active
      clusters:
        - id: "abc-123"     # cluster identifier (populated by external system)
          name: "cluster-a" # human-readable cluster name
      nodePools:
        general:            # NodePoolTypeGeneral — gateway/discoverer/manager
          name: general-pool
          replicas: 3
          machineResource:
            cpu: "4"
            memory: "16Gi"
            storage: "100Gi"
        agent:              # NodePoolTypeValdAgent — optional dedicated agent pool
          name: agent-pool
          replicas: 6
          machineResource:
            cpu: "8"
            memory: "32Gi"
  vectorEngine:
    name: my-vald
    vald:
      defaults:
        logLevel: warn
      agent:
        ngt:
          dimension: 768
          distanceType: l2
          objectType: float
          creationEdgeSize: 20
          searchEdgeSize: 10
        persistentVolume:    # optional
          enabled: true
          storageClass: standard    # falls back to DEFAULT_STORAGE_CLASS
          accessMode: ReadWriteOnce  # falls back to DEFAULT_ACCESS_MODE
      indexer:
        indexSchedule: "@every 1m"
        concurrency: 2
        manager: true
        indexDuration: 1h
        saveDuration: 1h
      gateway:
        indexReplica: 3
        serviceType: LoadBalancer
        ingress:             # optional
          enabled: true
          host: vald.example.com
      discoverer:
        kind: DaemonSet      # DaemonSet | Deployment
      overlay: {}            # arbitrary JSON merged last (Helm-style patch)
```

#### `infrastructure[].clusters`

| Field | Required | Description |
|-------|----------|-------------|
| `id` | populated by external system | Cluster identifier. Empty while the cluster is being provisioned. |
| `name` | yes | Human-readable name used to derive the ValdRelease object name. |

#### `infrastructure[].nodePools`

Two pool types are recognized:

| Key | Constant | Hosts |
|-----|----------|-------|
| `general` | `NodePoolTypeGeneral` | gateway-lb, discoverer, manager-index |
| `agent` | `NodePoolTypeValdAgent` | vald-agent (NGT) |

When no `agent` pool is defined, the `general` pool is used for agent resource
sizing and replica count as well. See [Resource Allocation and Topology Strategy](resource-and-topology-strategy.md)
for the sizing formulas.

### Status

```yaml
status:
  phase: WaitingCreateVrs        # type of the condition currently being processed
  progress:
    total: 3                      # total number of lifecycle phases
    completed: 1                  # number of phases that have reached True
  conditions:
    - type: WaitingClusterCreate
      status: "True"
      reason: Succeeded
      message: "Waiting for Cluster Creation."
      lastTransitionTime: "2025-01-01T00:00:00Z"
    - type: WaitingCreateVrs
      status: Unknown
      reason: Progressing
      message: "Waiting for VRS creation."
      lastTransitionTime: "2025-01-01T00:00:01Z"
```

| Field | Description |
|-------|-------------|
| `phase` | The `type` of the condition currently being evaluated. Updated on every reconcile. |
| `progress.total` | Total lifecycle phase count. Set on every reconcile from the lifecycle definition. |
| `progress.completed` | Count of phases that have transitioned to `True`. Updated when a new phase is seeded. |
| `conditions` | Accumulated standard `metav1.Condition` array. Phases are never removed; each phase appears at most once. |

#### kubectl columns

```
NAME          PHASE              PROGRESS   STATUS   AGE
my-vor       WaitingCreateVrs   3          True     5m
```

> **Note**: The `STATUS` column shows `conditions[0].status`, which is always the
> first phase (`WaitingClusterCreate`). This reflects whether cluster provisioning
> has completed, not the overall readiness. Use `PHASE` for current processing state.

---

## Lifecycle Phases

The controller defines three phases in order:

| Order | Condition Type | Desired | Description |
|-------|---------------|---------|-------------|
| 0 | `WaitingClusterCreate` | `desired.Prop` | Validates that all clusters in `spec.infrastructure` have non-empty `id` and `name`. |
| 1 | `WaitingCreateVrs` | `desired.Resource` | Builds and creates `ValdRelease` objects; waits for them to exist in the API server. |
| 2 | `Completed` | `nil` | Terminal phase. No desired-state check; automatically marked `Succeeded` when reached. |

---

## Reconciliation Model

### Conditions Accumulation

Conditions are **never removed** from `Status.Conditions`. As phases progress, entries
are appended to the array. Every reconcile loop re-evaluates **all** accumulated conditions.

This provides self-healing: if a previously `True` condition later breaks (e.g. a
`ValdRelease` is deleted, or `spec.infrastructure` is modified), the controller detects
it on the next reconcile and halts at that phase until it recovers.

```
Reconcile N:   [WaitingClusterCreate=True]
               → WaitingCreateVrs seeded as Progressing

Reconcile N+1: [WaitingClusterCreate=True, WaitingCreateVrs=Unknown]
               → VRS created, but Get returns NotFound → stay Progressing

Reconcile N+2: [WaitingClusterCreate=True, WaitingCreateVrs=True]
               → Completed seeded as Succeeded

--- ValdRelease deleted externally ---

Reconcile N+3: [WaitingClusterCreate=True, WaitingCreateVrs=False, Completed=True]
               → WaitingCreateVrs fails → controller stops, returns error → requeued
               → VRS is recreated, WaitingCreateVrs recovers to True on next pass
```

### Reconcile Loop (simplified)

```go
// 1. First reconcile: no conditions yet → seed the first phase
if len(status.Conditions) == 0 {
    seed(phases[0], Progressing)
    return
}

// 2. Re-evaluate every accumulated condition
for _, cond := range status.Conditions {
    phase := phases.GetByType(cond.Type)
    if phase.Desired == nil { continue }

    createOrUpdateResources(phase)        // idempotent
    result := phase.Desired.IsReady(ctx)  // returns desired.Result
    updateCondition(cond.Type, result)

    if result != Succeeded {
        return error  // halt; controller-runtime requeues
    }
}

// 3. If all known phases passed, seed the next one
if next := phases.GetNext(); next exists {
    seed(next, Progressing or Succeeded)
}
```

### Phase Seeding

When a new phase is seeded:

- `Desired == nil` (e.g. `Completed`): seeded with `Succeeded` immediately.
- `Desired != nil`: seeded with `Progressing`; evaluated on the next reconcile.

---

## Desired State Model

Each lifecycle phase optionally holds a `Desired` implementation that answers two questions:

1. **`Build(ctx)`** — what Kubernetes resources should exist?
2. **`IsReady(ctx)`** — are those resources in the expected state?

### `desired.Result`

`IsReady` returns a `desired.Result` which maps directly to a `metav1.Condition`:

| Constructor | `Status` | `Reason` | Meaning |
|-------------|----------|----------|---------|
| `Progressing(msg)` | `Unknown` | `Progressing` | Controller is actively working (e.g. building resources). |
| `Pending(msg)` | `Unknown` | `Pending` | Waiting for an external event (e.g. external system setting `cluster.id`). |
| `Succeeded()` | `True` | `Succeeded` | Phase is complete. |
| `Failed(err)` | `False` | `Failed` | Unrecoverable error (e.g. missing required configuration). |

**`Pending` vs `Failed`**

`Pending` indicates an expected transient state — the system is waiting for something
outside its control. `Failed` indicates a configuration error that requires human intervention.

Example — `WaitingClusterCreate`:

| Condition | Result |
|-----------|--------|
| `spec.infrastructure` is empty | `Failed` — required field, user must fix |
| `infra.clusters` is empty | `Failed` — required field, user must fix |
| `cluster.id == ""` | `Pending` — external system has started but has not assigned an ID yet |
| `cluster.name == ""` | `Pending` — same reason |
| all clusters have `id` and `name` | `Succeeded` |

### `desired.Prop`

Used for pure in-memory checks (no Kubernetes API calls). Holds a `Check func() Result`.
`Build` is a no-op.

### `desired.Resource`

Used when checking Kubernetes resources via the API server. Holds a `Builder` (which produces
the list of expected objects) and a `Client`. `Build` calls the builder and stores the list;
`IsReady` iterates the list and calls `client.Get` for each object.

| Get result | `IsReady` returns |
|------------|-------------------|
| `NotFound` | `Progressing("")` — resources not created yet |
| other error | `Failed(err)` |
| found, `Check` returns non-True | returned as-is |
| all found and passing | `Succeeded()` |

---

## Resource Management

### NodePool Match (optional)

When the environment variable `REQUIRE_NODEPOOL_MATCH=true` is set, the controller
only generates `ValdRelease` objects for infra entries that have matching Kubernetes
`Node` objects. Nodes are matched by labels:

```
<prefix>/namespace = <vor namespace>
<prefix>/type      = general | agent
```

`prefix` defaults to empty (bare `namespace`/`type` keys) and can be configured via
`NODEPOOL_LABEL_PREFIX`.

### Generation Tracking

Every created/updated resource is labelled `managed-generation: <vor.Generation>`.
On each reconcile, resources owned by the `ValdOperatorRelease` CR but absent from the current
build output are **pruned** (deleted).

### Owner References

All generated resources have a controller owner reference pointing to the `ValdOperatorRelease` CR,
enabling garbage collection when the CR is deleted.

---

## Configuration

| Environment Variable | Default | Description |
|---------------------|---------|-------------|
| `REQUIRE_NODEPOOL_MATCH` | `false` | Only generate VRS when matching Node labels exist. |
| `NODEPOOL_LABEL_PREFIX` | `""` | Prefix for node pool label keys. |

See [explain.md](../../explain.md#configuration-environment-variables) for the full
environment-variable reference (PV sizing, ingress, gateway service type, discoverer
DaemonSet strategy, etc.).

---

## Open Items

- **Requeue behavior**: Pending/Progressing conditions currently return an `error`, causing
  exponential back-off and noisy logs. These should return `ctrl.Result{RequeueAfter: duration}`
  instead. The duration should be configurable.

- **`util.UpdateStatus`**: The custom condition upsert helper should be replaced with
  the standard `meta.SetStatusCondition` from `k8s.io/apimachinery/pkg/api/meta`.

- **`STATUS` printcolumn**: Shows `conditions[0].status` (always the first phase).
  Should be changed to reflect the current overall state more meaningfully.
