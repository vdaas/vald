# Capacity Planning

## What is capacity planning for the Vald cluster?

Capacity planning is essential before deploying the Vald cluster to the cloud service.
There are three viewpoints: Vald cluster view, Kubernetes view, and Component view.
Let's see each view.

<div class="notice">
When introducing production, we recommend that you actually measure how many resources are required for verification.
</div>

## Vald cluster view

The essential point at the Vald cluster view is the hardware specification, especially RAM.
The Vald cluster, especially Vald Agent components, requires much RAM capacity because the vector index is stored in memory.

The minimum required memory for each vector (bit) is:

```bash
// minimum required bits of vector
{ oid (64bit) + timestamp (64bit) + uuid (user defined) } * 2 + { dimension * 64 } + { the creation edge size + the search edge size } * 8
```

Considering the `index size` and `index_replica`, it is easy to figure out the minimum required RAM capacity by the following formula.

```bash
{ minimum required bits of vector } * { the index size } * { index_replica }
```

For example, you want to insert 1 million vectors with 900 dimensions with 32 byte (256 bit) UUID, the index replica is 3, `creation edge size` is 20, and `search edge size` is 10, the minimum required RAM capacity is:

```bash
{(64 + 64 + 256) × 2 + (900 × 64) + (20 + 10) × 8 } × 1,000,000 × 3 = 175,824,000,000 (bit) = 21.978 (GB)
```

It is just the minimum required RAM for indexing.
Considering the margin of RAM capacity, the minimum RAM capacity should be less than 60% of the actual RAM capacity.
Therefore, the actual minimum RAM capacity will be:

```bash
8,7168,000,000 (bit) / 0.6 = 145,280,000,000 (bit) = 36.63 (GB)
```

<div class="warn">
In the production usage, memory usage may be not enough in the minimum required RAM.<BR>
Because for example, there are a noisy problem, high memory usage for createIndex (indexing on memory), high traffic needs more memory, etc.
</div>

## Kubernetes cluster view

### Pod priority & QoS

When the Node capacity (e.g., RAM, CPU) reaches the limit, Kubernetes will decide to kill some Pods according to QoS and Pod priority.
Kubernetes performs pod scheduling with pods Priority Class as the priority and QoS as the second priority.

#### Pod priority

Pod priority has the integer value, and the higher value, the higher priority.

Each Vald component has the default priority value:

- Agent: 1000000000
- Discoverer: 1000000
- Filter Gateway: 1000000
- LB Gateway: 1000000
- Index Manager: 1000000

Therefore, the order of priority is as follows:

```bash
Agent > Discoverer = Filter Gateway = LB Gateway = Index Manger
```

Those values will be helpful when the Pods other than the Vald component are in the same Node.

It is easy to change by editing your `values.yaml`.

```yaml
# e.g. LB Gateway podPriority settings.
...
  gateway:
    lb:
    ...
      podPriority:
        enabled: true
        value: {new values}
      ...
```

#### QoS

QoS value can be either Guaranteed, Burstable, or BestEffort.
And, QoS priority is higher in the order of Guaranteed, Burstable, BestEffort, and Kubernetes will kill Pods in ascending order of importance.

Resource request and limit determine QoS.

The below table shows the condition for each QoS.

|    QoS     |   request CPU   | request Memory  |    limit CPU    | request Memory  | Sup.                                |
| :--------: | :-------------: | :-------------: | :-------------: | :-------------: | :---------------------------------- |
| Guaranteed |       :o:       |       :o:       |       :o:       |       :o:       | All settings are required.          |
| Burstable  | :o: (:warning:) | :o: (:warning:) | :o: (:warning:) | :o: (:warning:) | One to three settings are required. |
| BestEffort |       :x:       |       :x:       |       :x:       |       :x:       | No setting is required.             |

Vald requires many RAM resources because of on-memory indexing, so we highly recommend that you do not specify a limit, especially for the Vald Agent.
In this case, QoS will be Burstable.

**Throttling**

The CPU throttling affects the pod performance.

If it occurs, the Vald cluster operator must consider each component's CPU resource request and limit.
It is easy to change by editing the `values.yaml` file and applying it.

```yaml
# e.g. LB Gateway resources settings.
...
  gateway:
    lb:
    ...
      resources:
        requests:
          cpu: 200m
          memory: 150Mi
        limits:
          cpu: 2000m
          memory: 700Mi
    ...
```

<div class="warning">
Please take care of pod priority and QoS.
</div>

### Node & Pod affinity

Kubernetes scheduler often applies Pods on Node based on resource availability.

In production usage, other components sometimes work on the Kubernetes cluster where the Vald cluster runs.
Depending on the situation, you may want to deploy to a different Node: e.g., when running a machine learning component that requires high memory on an independent Node.

In this situation, we recommend you to set the affinity/anti-affinity configuration for each Vald component.
It is easy to change by editing each component setting on your `values.yaml`.

<div class="warning">
The affinity setting for Vald Agent is significant for the Vald cluster.<BR>
Please DO NOT remove the default settings.
</div>

```yaml
# e.g. Agent's affinity settings
...
  agent:
  ...
    affinity:
      nodeAffinity:
        preferredDuringSchedulingIgnoredDuringExecution: []
        requiredDuringSchedulingIgnoredDuringExecution:
          nodeSelectorTerms: []
      podAffinity:
        preferredDuringSchedulingIgnoredDuringExecution: []
        requiredDuringSchedulingIgnoredDuringExecution: []
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
        requiredDuringSchedulingIgnoredDuringExecution: []
  ...
```

For more information about Kubernetes affinity, please refer to [here](https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/#affinity-and-anti-affinity)

## Component view

Depending on the customization of each component for each user, there are some points to be aware of.

### Index Manager

If the `saveIndex` is executed frequently, the backup data per unit time will increase, which consumes bandwidth.

Similarly, as the saveIndex concurrency increases, the backup data per unit time increases.
