# Resource Allocation and Topology Strategy

## Background

Each `ValdOperatorRelease` CR describes one or more node pools. A "general" pool hosts the
gateway, discoverer, and manager components. An optional "agent" (ValdAgent) pool
hosts the NGT vector index pods.

When there is no dedicated agent pool the controller falls back to the general pool
for agent placement (via `effectiveNodePoolType`). Before this fix that fallback was
not applied to resource and replica sizing, producing zero-replica agent deployments
in no-agent-pool configurations.

## Node Pool Fallback for Agent Resources

The agent node-pool selection rule lives in
`pkg/operator/vald/service/rules.go` (`resolveAgentNodePool`) and is called directly by
the builder. This keeps `vrsBuilder.Build` a pure function of `(CR, Config, Capability)`:

```go
// pkg/operator/vald/service/builder.go
agentPool := resolveAgentNodePool(infra)
row.SetRelationalResources(agentPool.NodeCount, agentPool.MachineResource, resourceParams)
```

`ResolveAgentNodePool` returns the dedicated `agent` pool when present, and otherwise
falls back to the `general` pool for both replica count and machine resources — so a
config with no agent pool no longer produces zero-replica agent deployments.

`effectiveNodePoolType` (`pkg/operator/vald/service/builder.go`) applies the same fallback
to NodeSelector/Tolerations, so placement and resource sizing stay consistent.
`resourceParams` carries `AgentPodsPerNode` and the discoverer DaemonSet rolling-update
values from `Config`.

## Agent Resource Calculation

Agent pods request a fraction of the node's resources so that
`AgentPodsPerNode` (default 2) pods fit on one node:

```
CPU request  = nodeCPU  × ResourceRatio / AgentPodsPerNode
RAM request  = nodeRAM  × ResourceRatio / AgentPodsPerNode
CPU limit    = nodeCPU  × ResourceRatio          (all pods on this node share the limit)
RAM limit    = (not set)
```

`ResourceRatio = 0.6` — 60% of the node is reserved for agent pods combined; the
remaining 40% covers OS, DaemonSets, and other system processes.

Memory limit is intentionally absent: the NGT index grows on disk after startup and
is then mmap'd into address space. A hard memory limit would trigger OOM kills as the
index expands.

## Agent Pod Count

```
agentPodCount = agentNodeCount × AgentPodsPerNode
```

`MinReplicas = MaxReplicas = agentPodCount` (agent uses a StatefulSet with no HPA).

## Gateway Replica Scaling

Gateway replicas track agent replica count with an HPA range:

```
minReplicas = max(agentPodCount / 2, 1)
maxReplicas = max(agentPodCount × 2, 1)
```

## Gateway / Discoverer / Manager Resources

These components are not node-resource-dependent; they run against vald upstream
defaults and are unaffected by node pool sizing:

| Component     | CPU req | RAM req | CPU lim | RAM lim |
| ------------- | ------- | ------- | ------- | ------- |
| gateway-lb    | 200m    | 150Mi   | 2000m   | 700Mi   |
| discoverer    | 200m    | 65Mi    | 600m    | 200Mi   |
| manager-index | 200m    | 80Mi    | 1000m   | 500Mi   |

## Topology Spread Constraints

All four components receive a `topologySpreadConstraints` entry so Kubernetes
spreads pods across nodes rather than relying on implicit scheduler behaviour:

```yaml
topologySpreadConstraints:
  - maxSkew: 1
    topologyKey: kubernetes.io/hostname
    whenUnsatisfiable: DoNotSchedule
    labelSelector:
      matchLabels:
        app.kubernetes.io/component: <agent|gateway-lb|discoverer|manager-index>
```

`maxSkew: 1` with `DoNotSchedule` ensures at most one extra pod per node compared
to the least-loaded node. This replaces the previous approach of inflating resource
requests to crowd out extra pods, which was fragile and hard to reason about.
