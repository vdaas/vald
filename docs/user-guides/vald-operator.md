# Vald Operator

The Vald Operator is a Kubernetes controller that manages the lifecycle of Vald clusters through a single custom resource.
It watches `Mvaldrelease` (short name `mvrs`, API group `vald.vdaas.org/v1`) resources and generates
[`ValdRelease`](https://github.com/vdaas/vald-helm-operator) (VRS) manifests consumed by the Vald Helm Operator (VHO).

```
Mvaldrelease (mvrs)  ──reconcile──▶  ValdRelease (VRS) × (active infra × clusters)
   minimal input                       consumed by vald-helm-operator
```

The controller only *generates* VRS definitions.
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
make docker-build docker-push IMG=<some-registry>/vald-operator:tag
```

### Install CRDs and deploy the controller

```sh
make install
make deploy IMG=<some-registry>/vald-operator:tag
```

### Apply a sample Mvaldrelease

```sh
kubectl apply -k config/samples/
```

### Uninstall

```sh
kubectl delete -k config/samples/
make uninstall
make undeploy
```

## CRD: Mvaldrelease (`mvrs`)

The goal of the `Mvaldrelease` resource is to collapse the large `ValdRelease` configuration surface into the minimum a user must supply.
There are two input groups:

1. **Infrastructure / node-pool information** — `spec.infrastructure[]`
2. **Minimal Vald settings** — `spec.vectorEngine.vald`

### Spec example

```yaml
apiVersion: vald.vdaas.org/v1
kind: Mvaldrelease
metadata:
  name: my-mvrs
  namespace: vald
spec:
  infrastructure:
    - role: green               # arbitrary role label (e.g. green/blue, hot/cold)
      type: kind                # cluster type hint (informational)
      active: true
      clusters:
        - id: "abc-123"         # cluster identifier (populated by external system)
          name: "cluster-a"     # human-readable cluster name
      nodePools:
        general:                # gateway / discoverer / manager
          name: general-pool
          replicas: 3
          machineResource:
            cpu: "4"
            memory: "16Gi"
            storage: "100Gi"
        agent:                  # optional dedicated agent pool
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
        kind: DaemonSet         # DaemonSet | Deployment
      overlay: {}               # arbitrary JSON merged onto the generated VRS
```

### `spec.infrastructure[]`

| Field | Notes |
|-------|-------|
| `role` | Free-form role label (e.g. `hot`, `standby`, `blue`, `green`). Copied to VRS labels. |
| `type` | Cluster type label. |
| `active` | When `false`, the entry is skipped during VRS generation. |
| `clusters[]` | `{ id, name }`. `id` is typically filled by an external system; `name` must be set. |
| `nodePools` | Map keyed by pool type: `general` (required) and `agent` (optional). Each pool carries `name`, `replicas`, and `machineResource{ cpu, memory, storage }`. |

#### Node pool types

| Key | Hosts |
|-----|-------|
| `general` | gateway-lb, discoverer, manager-index |
| `agent` | vald-agent (NGT) |

When no `agent` pool is defined, the `general` pool is used for agent resource sizing and replica count as well.

### `spec.vectorEngine.vald`

| Field | Description |
|-------|-------------|
| `defaults.logLevel` | Log level for all Vald components. |
| `agent.ngt` | NGT settings: `dimension` (required, >= 2), `creationEdgeSize`, `searchEdgeSize`, `distanceType`, `objectType`. |
| `agent.persistentVolume` | Optional. `{ enabled, storageClass, accessMode }`. Falls back to environment defaults when omitted. |
| `indexer` | `indexSchedule`, `saveSchedule`, `concurrency`, `manager`, suspend flags, durations. |
| `gateway` | `indexReplica`, `serviceType`, `ingress{ enabled, host }`. |
| `discoverer.kind` | `DaemonSet` or `Deployment`. |
| `overlay` | Arbitrary JSON patch merged onto the generated VRS (Helm-style patch). |

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

| Field | Description |
|-------|-------------|
| `phase` | The condition type currently being evaluated. |
| `progress.total` | Total lifecycle phase count. |
| `progress.completed` | Phases that have reached `True`. |
| `conditions` | Accumulated `metav1.Condition` array. Phases are never removed. |

## Lifecycle

Reconciliation is modeled as an ordered flow of phases.
Each phase carries a `Condition` plus an optional `Builder` (creates/updates resources) and `Checker` (reports readiness).

| Phase | Builder | Checker | Purpose |
|-------|---------|---------|---------|
| `WaitingClusterCreate` | -- | yes | Validate infra config; wait until every cluster has `id` and `name`. |
| `WaitingCreateVrs` | yes | yes | Build VRS objects and wait until they are ready. |
| `Completed` | -- | -- | Terminal phase. |

### Accumulating, self-healing conditions

`status.conditions` **accumulates** as phases progress (conditions are never removed), and every reconcile re-evaluates **all** conditions.
If a previously-`True` condition later breaks (e.g. a generated VRS is deleted), the controller detects it and restarts work from that phase.

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

| Result | Condition Status | Meaning |
|--------|------------------|---------|
| `Progressing` | Unknown | The controller is actively working. |
| `Pending` | Unknown | Waiting for an external event (e.g. cluster ID assignment). |
| `Succeeded` | True | Phase is complete. |
| `Failed` | False | Misconfiguration or unrecoverable error. |

## VRS Generation

`VrsBuilder.Build` is a pure function of `(CR, Config, NodePoolCapability)` — it makes no Kubernetes API calls.

1. **Validate** the CR (`infrastructure` non-empty, each cluster has a `name`).
2. **Iterate `infrastructure[]`**, skipping inactive entries (and, when node-pool matching is enabled, entries with no matching general pool).
3. **Assemble VRS spec**: defaults, gateway, agent (NGT), manager, discoverer.
4. **Resolve resources from node pools**: derive replicas and per-component CPU/memory from the agent node pool.
5. **Apply optional settings**: persistent volume, node affinities.
6. **Per cluster**: for each `infra.clusters[]`, set name, apply labels, then merge the overlay.
   One VRS is produced per cluster, so a single MVRS can yield many VRS objects.

### Overlay

`spec.vectorEngine.vald.overlay` is a JSON patch merged onto the generated VRS, layered on top of the default VRS template loaded at startup (`DEFAULT_VRS_PATH`).
This provides an escape hatch for any VRS field not directly exposed in the MVRS spec.

### Resource management

All generated resources have a controller owner reference pointing to the `Mvaldrelease` CR, enabling garbage collection when the CR is deleted.
Every resource is labelled `managed-generation: <mvrs.Generation>`.
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
Memory limit is intentionally absent because the NGT index uses mmap and a hard limit would trigger OOM kills.

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

| Component | CPU req | RAM req | CPU lim | RAM lim |
|-----------|---------|---------|---------|---------|
| gateway-lb | 200m | 150Mi | 2000m | 700Mi |
| discoverer | 200m | 65Mi | 600m | 200Mi |
| manager-index | 200m | 80Mi | 1000m | 500Mi |

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

The same MVRS can be distributed to multiple clusters, generating a VRS only where a matching node pool exists.

### Node-pool matching

When `REQUIRE_NODEPOOL_MATCH=true`, the controller lists Nodes and resolves a `NodePoolCapability`.
A VRS is generated only when a `general` pool is present.
Nodes are matched by labels:

```
<prefix>/namespace = <mvrs namespace>
<prefix>/type      = general | agent
```

`prefix` defaults to empty and can be configured via `NODEPOOL_LABEL_PREFIX`.

When `REQUIRE_NODEPOOL_MATCH=false` (default), the controller skips node listing and treats every infra entry as schedulable.

### Use cases

- **Single cluster**: Deploy MVRS once; VRS generates in that cluster.
- **Multi-cluster same config**: Deploy MVRS to all clusters; only generate VRS where node-pool labels match.
- **Blue/green deployments**: Multiple MVRS objects with different roles and endpoints.
- **VRS generation only**: Generate VRS manifests in a management cluster and distribute separately to workload clusters.

## Configuration

All configuration is done via environment variables, loaded once at startup.

| Environment Variable | Default | Description |
|---------------------|---------|-------------|
| `DEFAULT_VRS_PATH` | `/opt/mvaldrelease/config/vrs.yaml` | Default VRS template merged with the overlay. |
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
| `INTERNAL_HOST_DOMAIN` | `""` | Optional internal host domain. |

## Development

### Common tasks

```sh
# Build and run unit tests
make build
make test

# Regenerate CRDs / RBAC / deepcopy after changing api/v1 types
make manifests generate

# Lint
make lint

# Build and push the controller image
make docker-build docker-push IMG=<registry>/vald-operator:tag

# Install CRDs and deploy the controller
make install
make deploy IMG=<registry>/vald-operator:tag
```

Run `make help` for the full list of targets.

<div class="notice">
<code>config/</code> is kubebuilder-generated. Do not hand-edit generated manifests; change the kubebuilder markers in the Go sources and regenerate with <code>make manifests generate</code>.
</div>
