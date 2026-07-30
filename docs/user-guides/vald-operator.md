# Vald Operator

The Vald Operator is a Kubernetes controller that manages the lifecycle of Vald clusters through a single custom resource.
It watches `ValdOperatorRelease` (short name `vor`, API group `vald.vdaas.org/v1`) resources and generates
[`ValdRelease`](https://github.com/vdaas/vald-helm-operator) (VRS) manifests consumed by the Vald Helm Operator (VHO).

```
ValdOperatorRelease (vor)  ──reconcile──▶  ValdRelease (VRS) × (active infra × clusters)
   minimal input                              consumed by vald-helm-operator
```

The controller _generates_ VRS definitions.
It does not run Vald itself; a VHO running in the target cluster turns each VRS into Vald pods.
This separation lets a management cluster emit VRS definitions and distribute them to other clusters.

## Prerequisites

- Go 1.23+
- Docker 17.03+
- kubectl 1.11.3+ and access to a Kubernetes 1.11.3+ cluster
- A working [Vald Helm Operator (VHO)](https://vald.vdaas.org/docs/operator/vald-helm-operator/) in the target cluster

## Deployment

### Build and push the controller image

```sh
make docker/build/operator/vald
```

### Install CRDs and deploy the controller

```sh
make k8s/operator/vald/deploy
```

### Apply a sample ValdOperatorRelease

```sh
kubectl apply -f cmd/operator/vald/sample.yaml
```

### Uninstall

```sh
kubectl delete -f cmd/operator/vald/sample.yaml
make k8s/operator/vald/delete
```

## CRD: ValdOperatorRelease (`vor`)

The goal of the `ValdOperatorRelease` resource is to collapse the large `ValdRelease` configuration surface into what a user must supply.
There are two input groups:

1. **Infrastructure / node-pool information** — `spec.infrastructure[]`
2. **Minimal Vald settings** — `spec.vectorEngine.vald`

### Spec example

```yaml
apiVersion: vald.vdaas.org/v1
kind: ValdOperatorRelease
metadata:
  name: my-vor
  namespace: vald
spec:
  infrastructure:
    - role: green # arbitrary role label (e.g. green/blue, hot/cold)
      type: kind # cluster type hint (informational)
      active: true
      clusters:
        - id: "abc-123" # cluster ID (populated by external system)
          name: "cluster-a" # human-readable cluster name
      nodePools:
        general: # gateway / discoverer / manager
          name: general-pool
          replicas: 3
          machineResource:
            cpu: "4"
            memory: "16Gi"
            storage: "100Gi"
        agent: # optional dedicated agent pool
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
        persistentVolume:
          enabled: true
          storageClass: standard
          accessMode: ReadWriteOnce
      indexer:
        indexSchedule: "@every 1m"
        concurrency: 2
        manager: true
        indexDuration: 1h
        saveDuration: 1h
      gateway:
        indexReplica: 3
        serviceType: LoadBalancer
        ingress:
          enabled: true
          host: vald.example.com
      discoverer:
        kind: DaemonSet # DaemonSet | Deployment
      overlay: {} # arbitrary JSON merged onto the generated VRS
```

### `spec.infrastructure[]`

| Field        | Notes                                                                                                                                                     |
| ------------ | --------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `role`       | Free-form role label (e.g. `hot`, `standby`, `blue`, `green`). Copied to VRS labels.                                                                      |
| `type`       | Cluster type label.                                                                                                                                       |
| `active`     | When `false`, the entry is skipped during VRS generation.                                                                                                 |
| `clusters[]` | `{ id, name }`. The `id` field is typically filled by an external system; `name` must be set.                                                             |
| `nodePools`  | Map keyed by pool type: `general` (required) and `agent` (optional). Each pool carries `name`, `replicas`, and `machineResource{ cpu, memory, storage }`. |

#### Node pool types

| Key       | Hosts                                 |
| --------- | ------------------------------------- |
| `general` | gateway-lb, discoverer, manager-index |
| `agent`   | `vald-agent` (NGT)                    |

When no `agent` pool is defined, the `general` pool is used for agent resource sizing and replica count as well.

### `spec.vectorEngine.vald`

| Field                    | Description                                                                                                     |
| ------------------------ | --------------------------------------------------------------------------------------------------------------- |
| `defaults.logLevel`      | Log level for all Vald components.                                                                              |
| `agent.ngt`              | NGT settings: `dimension` (required, >= 2), `creationEdgeSize`, `searchEdgeSize`, `distanceType`, `objectType`. |
| `agent.persistentVolume` | Optional. `{ enabled, storageClass, accessMode }`. Falls back to environment defaults when omitted.             |
| `indexer`                | `indexSchedule`, `saveSchedule`, `concurrency`, `manager`, suspend flags, durations.                            |
| `gateway`                | `indexReplica`, `serviceType`, `ingress{ enabled, host }`.                                                      |
| `discoverer.kind`        | `DaemonSet` or `Deployment`.                                                                                    |
| `overlay`                | Arbitrary JSON patch merged onto the generated VRS (Helm-style patch).                                          |

### Status

```yaml
status:
  phase: WaitingCreateVrs
  progress:
    total: 3
    completed: 1
  conditions:
    - type: WaitingClusterCreate
      status: "True"
      reason: Succeeded
    - type: WaitingCreateVrs
      status: Unknown
      reason: Progressing
```

| Field                | Description                                                     |
| -------------------- | --------------------------------------------------------------- |
| `phase`              | The condition type currently being evaluated.                   |
| `progress.total`     | Total lifecycle phase count.                                    |
| `progress.completed` | Phases that have reached `True`.                                |
| `conditions`         | Accumulated `metav1.Condition` array. Phases are never removed. |

## Lifecycle

Reconciliation is modeled as an ordered flow of phases.
Each phase carries a `Condition` plus an optional `Builder` (creates/updates resources) and `Checker` (reports readiness).

| Phase                  | Builder | Checker | Purpose                                                              |
| ---------------------- | ------- | ------- | -------------------------------------------------------------------- |
| `WaitingClusterCreate` | --      | yes     | Validate infra config; wait until every cluster has `id` and `name`. |
| `WaitingCreateVrs`     | yes     | yes     | Build VRS objects and wait until they are ready.                     |
| `Completed`            | --      | --      | Terminal phase.                                                      |

### Accumulating, self-healing conditions

`status.conditions` **accumulates** as phases progress (conditions are never removed), and every reconcile re-evaluates **all** conditions.
If a `True` condition later breaks (e.g. a generated VRS is deleted), the controller detects it and restarts work from that phase.

```
Reconcile N:   [WaitingClusterCreate=True]
               → WaitingCreateVrs seeded as Progressing

Reconcile N+1: [WaitingClusterCreate=True, WaitingCreateVrs=Unknown]
               → VRS created, but Get returns NotFound → stay Progressing

Reconcile N+2: [WaitingClusterCreate=True, WaitingCreateVrs=True]
               → Completed seeded as Succeeded

--- ValdRelease deleted externally ---

Reconcile N+3: [WaitingClusterCreate=True, WaitingCreateVrs=False]
               → controller detects break → VRS is recreated → recovers
```

### Readiness states

| Result        | Condition Status | Meaning                                                     |
| ------------- | ---------------- | ----------------------------------------------------------- |
| `Progressing` | Unknown          | The controller is actively working.                         |
| `Pending`     | Unknown          | Waiting for an external event (e.g. cluster ID assignment). |
| `Succeeded`   | True             | Phase is complete.                                          |
| `Failed`      | False            | Misconfiguration or unrecoverable error.                    |

## VRS Generation

`VrsBuilder.Build` is a pure function of `(CR, Config, NodePoolCapability)` — it makes no Kubernetes API calls.

1. **Validate** the CR (`infrastructure` non-empty, each cluster has a `name`).
2. **Iterate `infrastructure[]`**, skipping inactive entries (and, when node-pool matching is enabled, entries with no matching general pool).
3. **Assemble VRS spec**: defaults, gateway, agent (NGT), manager, discoverer.
4. **Resolve resources from node pools**: derive replicas and per-component CPU/memory from the agent node pool.
5. **Apply optional settings**: persistent volume, node affinities.
6. **Per cluster**: for each `infra.clusters[]`, set name, apply labels, then merge the overlay.
   One VRS is produced for each cluster, so a single VOR can yield one VRS per target cluster.

### Overlay

`spec.vectorEngine.vald.overlay` is a JSON patch merged onto the generated VRS, layered on top of the default VRS template loaded at startup (`DEFAULT_VRS_PATH`).
This provides an escape hatch for any VRS field not directly exposed in the VOR spec.

### Resource management

All generated resources have a controller owner reference pointing to the `ValdOperatorRelease` CR, enabling garbage collection when the CR is deleted.
Every resource is labelled `managed-generation: <vor.Generation>`.
On each reconcile, resources owned by the CR but absent from the current build output are pruned (deleted).

## Resource Allocation and Topology

### Agent resource calculation

Agent pods request a fraction of the node's resources so that `AgentPodsPerNode` (default 2) pods fit on one node:

```
CPU request  = nodeCPU  * 0.6 / AgentPodsPerNode
RAM request  = nodeRAM  * 0.6 / AgentPodsPerNode
CPU limit    = nodeCPU  * 0.6
RAM limit    = (not set)
```

The 0.6 ratio reserves 60% of the node for agent pods; the remaining 40% covers OS, DaemonSets, and system processes.
Memory limit is intentionally absent because the NGT index uses memory-mapped files and a hard limit would trigger OOM kills.

### Agent pod count

```
agentPodCount = agentNodeCount * AgentPodsPerNode
```

`MinReplicas = MaxReplicas = agentPodCount` (agent uses a StatefulSet with no HPA).

### Gateway replica scaling

```
minReplicas = max(agentPodCount / 2, 1)
maxReplicas = max(agentPodCount * 2, 1)
```

### Gateway / discoverer / manager resources

These components use fixed defaults:

| Component     | CPU req | RAM req | CPU lim | RAM lim |
| ------------- | ------- | ------- | ------- | ------- |
| gateway-lb    | 200m    | 150Mi   | 2000m   | 700Mi   |
| discoverer    | 200m    | 65Mi    | 600m    | 200Mi   |
| manager-index | 200m    | 80Mi    | 1000m   | 500Mi   |

### Topology spread constraints

All components receive a `topologySpreadConstraints` entry to spread pods across nodes:

```yaml
topologySpreadConstraints:
  - maxSkew: 1
    topologyKey: kubernetes.io/hostname
    whenUnsatisfiable: DoNotSchedule
    labelSelector:
      matchLabels:
        app.kubernetes.io/component: <component>
```

## Multi-Cluster Distribution

The same VOR can be distributed across clusters and generates a VRS where a matching node pool exists.

### Node-pool matching

When `REQUIRE_NODEPOOL_MATCH=true`, the controller lists Nodes and resolves a `NodePoolCapability`.
A VRS is generated when a `general` pool is present.
Nodes are matched by labels:

```
<prefix>/namespace = <vor namespace>
<prefix>/type      = general | agent
```

`prefix` defaults to empty and can be configured via `NODEPOOL_LABEL_PREFIX`.

When `REQUIRE_NODEPOOL_MATCH=false` (default), the controller skips node listing and treats every infra entry as schedulable.

### Use cases

- **Single cluster**: Deploy VOR once; VRS generates in that cluster.
- **Multi-cluster same config**: Deploy VOR to all clusters; generate VRS where node-pool labels match.
- **Blue/green deployments**: Use separate VOR objects with different roles and endpoints.
- **VRS generation**: Generate VRS manifests in a management cluster and distribute them to workload clusters.

## Configuration

All configuration is loaded once at startup from a YAML config file, following the same convention as every other Vald component.
See [`cmd/operator/vald/sample.yaml`](https://github.com/vdaas/vald/blob/main/cmd/operator/vald/sample.yaml) for a complete example; when deploying via the Helm chart, every key is exposed through `charts/operator/vald/values.yaml`.

The reconciler-facing settings live under the `operator` key:

| Config Key                                                 | Default                                    | Description                                                       |
| ---------------------------------------------------------- | ------------------------------------------ | ----------------------------------------------------------------- |
| `operator.vrs.default_vrs_path`                            | `/opt/valdoperatorrelease/config/vrs.yaml` | Default VRS template merged with the overlay.                     |
| `operator.vrs.log_level`                                   | `warn`                                     | Log level passed through to the generated VRS.                    |
| `operator.node_pool.require_match`                         | `false`                                    | Only generate VRS where matching node pools exist.                |
| `operator.node_pool.label_prefix`                          | `""`                                       | Prefix for the `namespace`/`type`/`role` node labels.             |
| `operator.node_pool.agent_pods_per_node`                   | `2`                                        | Agent pods packed per node when computing replicas.               |
| `operator.persistent_volume.default_storage_class`         | `standard`                                 | PV storage class fallback.                                        |
| `operator.persistent_volume.default_access_mode`           | `ReadWriteOnce`                            | PV access mode fallback.                                          |
| `operator.persistent_volume.buffer_ratio`                  | `1.5`                                      | PV size = `max(memoryRequest * ratio, min)`.                      |
| `operator.persistent_volume.min_size_bytes`                | `1073741824`                               | Minimum PV size in bytes.                                         |
| `operator.networking.enable_ingress`                       | `true`                                     | Enable gateway ingress generation.                                |
| `operator.networking.gateway_ingress_annotations`          | `{}`                                       | Annotations applied to the gateway ingress.                       |
| `operator.networking.gateway_service_type`                 | `NodePort`                                 | Gateway service type (`NodePort` / `ClusterIP` / `LoadBalancer`). |
| `operator.networking.discoverer_daemonset_max_surge`       | `30%`                                      | Discoverer DaemonSet rolling-update `maxSurge`.                   |
| `operator.networking.discoverer_daemonset_max_unavailable` | `0%`                                       | Discoverer DaemonSet rolling-update `maxUnavailable`.             |

Controller-level settings (leader election, `requeue` intervals, reconcile concurrency, watched namespaces) live under `operator.controller`, and server/observability settings (`server_config`, `observability`) follow the standard Vald component layout.

## Development

### Common tasks

```sh
# Run unit tests
make test/operator/vald

# Regenerate k8s manifests after changing charts/operator/vald
make k8s/manifest/operator/vald/update

# Lint
make lint

# Build and push the controller image
make docker/build/operator/vald

# Deploy the controller to the current cluster
make k8s/operator/vald/deploy
```

Run `make help` for the full list of targets.

## Further Reading

| Document                                                                 | Description                                                                                 |
| ------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------- |
| [Design Overview](vald-operator/design.md)                               | Controller architecture, CRD spec, reconcile lifecycle, and environment-variable reference. |
| [Controller Specification](vald-operator/controller-spec.md)             | Detailed reconcile model, desired-state model, and resource-management rules.               |
| [Resource and Topology Strategy](vald-operator/resource-and-topology.md) | Agent resource calculation, replica scaling formulas, and topology spread constraints.      |
| [Test Strategy](vald-operator/test-strategy.md)                          | Test layers, golden file tests, and coverage philosophy.                                    |
