# Capacity Planning

## What is capacity planning for the Vald cluster?

Capacity planning is essential before deploying the Vald cluster to the cloud service.
There are three viewpoints: Vald cluster view, Kubernetes view, and Component view.
Let's see each view.

## Vald cluster view

The essential point at the Vald cluster view is the hardware specification, especially RAM.
The Vald cluster, especially Vald Agent components, requires much RAM capacity because the vector index is stored in memory.

It is easy to figure out the minimum required RAM capacity by the following formula.

```bash
( { the dimension vector } × { bit number of vector } + { the bit of vectors ID string } ) × { the maximum number of the vector } × { the index replica }
```

For example, if you want to insert 1 million vectors with 900 dimensions and the object type is 32-bit with 32 byte (256 bit) ID, and the index replica is 3, the minimum required RAM capacity is:

```bash
(900 × 32 + 256 ) × 1,000,000 × 3 = 8,7168,000,000 (bit) = 10.896 (GB)
```

It is just the minimum required RAM.
Considering the margin of RAM capacity, the minimum RAM capacity should be less than 60% of the actual RAM capacity.
Therefore, the actual minimum RAM capacity will be:

```bash
8,7168,000,000 (bit) / 0.6 = 145,280,000,000 (bit) = 18.16 (GB)
```

<div class="warn">
Memory usage is not enough in the minimum required RAM.<BR>
E.g., there is a noisy problem, high memory usage for createIndex, high traffic needs more memory, etc...
There may be A noisy neighbor problem, High memory usage, 
</div>


## Kubernetes cluster view

### Pod priority & QoS

When the Node capacity (e.g., RAM, CPU) reaches the limit, Kubernetes will decide to kill some Pods according to QoS and Pod priority.
Kubernetes performs pod scheduling with pods Priority Class as the priority and QoS as the second priority.

**Pod priority**

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

**QoS**

QoS value can be either Guaranteed, Burstable, or BestEffort.
And, QoS priority is higher in the order of Guaranteed, Burstable, BestEffort, and Kubernetes will kill Pods in ascending order of importance.

Resource request and limit determine QoS.

The below table shows the condition for each QoS.

|    QoS     |   request CPU   | request Memory |    limit CPU    | request Memory  | Sup.                                |
| :--------: | :-------------: | :------------: | :-------------: | :-------------: | :---------------------------------- |
| Guaranteed |       :o:       |      :o:       |       :o:       |       :o:       | All settings are required.          |
| Guaranteed | :o: (:warning:) | :o: (:warning) | :o: (:warning:) | :o: (:warning:) | One to three settings are required. |
| BestEffort |       :x:       |      :x:       |       :x:       |       :x:       | No setting is required.             |

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

### Component view

Depending on the customization of each component for each user, there are some points to be aware of.

**Index Manager**

If the `saveIndex` is executed frequently, the backup data per unit time will increase, which consumes bandwidth.

Similarly, as the saveIndex concurrency increases, the backup data per unit time increases.
