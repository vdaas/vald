# Read Replica and Rotator

Read replica enhances the search QPS (Queries Per Second) of the Vald cluster by deploying read-only agents in addition to the regular agents and distributing the requests among them. Read replica is deployed as Kubernetes deployments, and depending on the number of replicas (N), QPS increases by approximately 1.7 to 1.8 times \* N.

<div class="notice">
The increase in QPS is possible with sufficient infrastructure (see <a href="#important-notes">Important notes</a>).
</div>

## How to deploy read replica

The read replica is managed with a separate chart from the Vald cluster and is deployed as an addon to the Vald cluster. The Vald cluster should be deployed first, followed by the deployment of the read replica.

> The reason Vald and Vald-readreplica are in separate charts is to avoid conflicts between the read replica's restart and the Helm operator's processes when Vald is managed by a helm operator. Therefore, the read replica will always be deployed using Helm commands.

### When you deploy Vald with Helm command

1. Edit `values.yaml` like below (Please refer to [deployment](../user-guides/deployment.md) for other fields.)

   ```yaml
   agent:
   ngt:
       export_index_info_to_k8s: true
   readreplica:
       enabled: true
       minReplicas: 1 # if you don't use hpa, this will be the replicas of the Deployment
       maxReplicas: 3
       hpa:
       enabled: true # if you prefer to use hpa
       targetCPUUtilizationPercentage: 80
   manager:
   index:
       operator:
       enabled: true
       rotation_job_concurrency: 2
   ```

1. Deploy Vald cluster

   ```bash
   helm install vald vald/vald --values values.yaml
   ```

1. Deploy `vald-readreplica` with the same `values.yaml`

   ```bash
   helm install vald-readreplica vald/vald-readreplica --values values.yaml
   ```

### When you deploy Vald cluster with `vald-helm-operator`

1. Edit `valdrelease.yaml` with the same fields as above

1. Deploy Vald cluster

   ```bash
   helm install vald-helm-operator-release vald/vald-helm-operator
   kubectl apply -f valdrelease.yaml
   ```

1. Deploy `vald-readreplica`

   ```bash
   helm install vald-readreplica vald/vald-readreplica --values <YOUR VALUES YAML FILE PATH>
   ```

## Architecture

Read replica mainly consists of the following four parts.

<img src="../../assets/docs/guides/read-replica-and-rotator/architecture.png" alt="Read Replica Architecture" />

### Read replica deployment

The deployment that generates Pods where the actual processing of read replica takes place. Read replica accepts read requests (search) and reads the index from the read replica PVC.

### Read replica PVC

The PVC for read replica Pods is used to read the index. It is generated based on the latest snapshot from the PVC of the regular agent. Unlike the agent PVC, it is generated as ROX, allowing it to be read from multiple Pods.

### Index operator

The operator handles the following processes:

1. Monitoring the time when the agent saved the index to the PVC and when the read replica performed index rotation
1. Generating [Read replica rotator](#read-replica-rotator) job when an index save occurs after the most recent rotation

> The Index operator also manages the timing of index create/save operations other than those mentioned above. Please refer to another document for details.

### Read replica rotator

The Kubernetes job handles the following processes:

1. Creating a snapshot from the agent's PVC
1. Generating a PVC for read replica from the snapshot
1. Rolling update of the read replica deployment to launch a group of read replica pods with the latest index

## Important notes

Result consistency is guaranteed

There is a time lag between index insertion, agent save, and the completion of read replica rotation. During this time, there may be inconsistencies between the index in the agent itself and the index in the read replica.

- Sufficient infrastructure is required for QPS scaling

  Even if read replicas are deployed, QPS will not scale if sufficient resources are not available in the Kubernetes cluster. Specifically, agent resources and read replica resources should be deployed on separate nodes. Vald sets `podAntiAffinity` to ensure that agent resources and read replica resources are deployed on separate nodes as much as possible.
