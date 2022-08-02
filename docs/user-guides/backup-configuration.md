# Backup configuration

There are three types of options for the Vald cluster: backup, filtering, and core algorithms.

This page describes how enabled backup features for your Vald cluster.

## What is the backup

Vald's backup function is to save the index data in each Vald Agent pod as a data file to the Persistent Volume or S3.
When the Vald Agent pod is restarted for some reason, the index state is restored from the saved index data.

## Backup configuration

This section shows the best practice backup configuration with PV, S3, or PV + S3.

Each sample configuration yaml is published on [here](https://github.com/vdaas/vald/tree/master/charts/vald/values).
Please try to reference.

### General

Regardless of the backup destination, the following Vald Agent settings must be set to enable backup.
- `index_path` is the backup file location.
- `in-memory-mode=false` means storing index files in the local volume.

In addition, `agent.terminationGracePeriodSeconds` value should be long enough to ensure the backup speed.

```yaml
agent:
  ...
  # We recommend setting this value long enough to ensure the backup speed of PV, since the Index is backed up at the end of the pod.
  terminationGracePeriodSeconds: 3600
  ngt:
    ...
    index_path: "/var/ngt/index"
    enable_in_memory_mode: false
    ...
```

### Persistent Volume

#### requirement

You must prepare PV before deployment when using Kubernetes Persistent Volume (PV) for backup storage.
Please refer to the setup guide of the usage environment for the provisioning PV.

For example:
- [GKE setup PV document](https://cloud.google.com/kubernetes-engine/docs/concepts/persistent-volumes)
- [EKS storage document](https://docs.aws.amazon.com/eks/latest/userguide/storage.html)

#### configuration

After provisioning PV, the following parameters are needed to be set.

```yaml
agent:
  ...
  persistentVolume:
    # use PV flag
    enabled: true
    # accessMode for PV (please verify your environment)
    accessMode: ReadWriteOncePod
    # storage class for PV (please verify your environment)
    storageClass: local-path
    # set enough size for backup
    size: 2Gi
  ...
```

Each PV will be mounted on each Vald Agent Pod's `index_path`.

It is highly recommended to set `copy_on_write` (CoW) as `true`. 

<div class="notice">
The CoW is an option to update the backup file safely.<BR>
The backup file may be corrupted, and the Vald Agent pod may not restore from backup files when the Vald Agent pod terminates during saveIndex without CoW is not be enabled.<BR>
On the other hand, when CoW is enabled, the Vald Agent pod can restore the data from one generation ago.
</div>

### S3

#### requirement

Before deployment, you must provision the S3 object storage.
You can use any S3-compatible object storage.

For example:
- [AWS S3](https://aws.amazon.com/s3/)
- [Google Cloud Storage](https://cloud.google.com/storage/docs/)

#### configuration

After provisioning the object storage, the following parameters are needed to be set.
To enable the backup function with S3, the Vald Agent Sidecar should be enabled.

```yaml
agent:
  ...
  sidecar:
    enabled: true
    initContainerEnabled: true
    # This is the Amazon S3 settings.
    # Please change it according to your environment.
    config:
      blob_storage:
        # storage type (default: s3)
        storage_type: "s3"
        # your bucket name
        bucket: "vald"
        s3:
          region: "us-central1"
    # If you enable sidecar, the following environment variables will be created automatically by default values.
    # So, please create the 'aws-secret' resource before deploying.
    # env:
    #   - name: AWS_ACCESS_KEY
    #     valueFrom:
    #       secretKeyRef:
    #         name: aws-secret
    #         key: access-key
    #   - name: AWS_SECRET_ACCESS_KEY
    #     valueFrom:
    #       secretKeyRef:
    #         name: aws-secret
    #         key: secret-access-key
```

The Vald Agent Sidecar needs an access key and a secret access key to communicate with your object storage.
Before applying the helm chart, register each value in Kubernetes secrets with the following commands.

```bash
kubectl create secret -n <Vald cluster namespace> aws-secret --access-key=<ACCESS KEY> --secret-access-key=<SECRET ACCESSS KEY>
```

### Persistent Volume and S3

You can use both PV and S3 at the same time.

```yaml
agent:
  minReplicas: 9
  maxReplicas: 9
  podManagementPolicy: Parallel
  resources:
    requests:
      cpu: 100m
      memory: 50Mi
  terminationGracePeriodSeconds: 3600
  persistentVolume:
    enabled: true
    accessMode: ReadWriteOncePod
    storageClass: local-path
    size: 2Gi
  ngt:
    dimension: 784
    index_path: "/var/ngt/index"
    enable_in_memory_mode: false
    auto_index_duration_limit: 730h
    auto_index_check_duration: 24h
    auto_index_length: 1000
    auto_save_index_duration: 365h
    auto_create_index_pool_size: 1000
  sidecar:
    enabled: true
    initContainerEnabled: true
    config:
      blob_storage:
        storage_type: "s3"
        bucket: "vald"
        s3:
          region: "us-central1"
    # If you enable sidecar, the following environment variables will be created automatically by default values.
    # So, please create the 'aws-secret' resource before deploying.
    # env:
    #   - name: AWS_ACCESS_KEY
    #     valueFrom:
    #       secretKeyRef:
    #         name: aws-secret
    #         key: access-key
    #   - name: AWS_SECRET_ACCESS_KEY
    #     valueFrom:
    #       secretKeyRef:
    #         name: aws-secret
    #         key: secret-access-key
```

## Restore

Restoring from the backup file runs on initContainer when the config is set correctly and the backup file exists.

If you use both PV and S3, the backup file used for restoration will prioritize the file on PV.
If the backup file does not exist on the PV, the backup file will be retrieved from S3 via the Vald Agent Sidecar and restored.
