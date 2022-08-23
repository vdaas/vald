# Backup configuration

There are three types of options for the Vald cluster: backup, filtering, and core algorithms.

This page describes how to enable the backup feature on your Vald cluster.

## What is the backup

Vald's backup function is to save the index data in each Vald Agent pod as a data file to the Persistent Volume or S3.
When the Vald Agent pod is restarted for some reason, the index state is restored from the saved index data.

## Backup methods

You can choose one of three types of backup methods.

- PV (recommended)
- S3
- PV + S3

Please refer to the following tables and decide which method fit for your case.

|             | PV                                                                                                                                                          | S3                                                                                                    | PV+S3                                                                                                                                            |
| :---------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------- | :---------------------------------------------------------------------------------------------------- | :----------------------------------------------------------------------------------------------------------------------------------------------- |
| usecase     | Want to use backup with low cost<BR>Would not like to use some external storage for backup<BR>Want to backup with highly compatible storage with Kubernetes | Want to use the same backup file with several Vald clusters<BR>Want to access the backup files easily | Want to use backup with PV basically and access the backup files easily<BR>Want to prevent backup file failure due to Kubernetes cluster failure |
| merit :+1:  | Easy to use<BR>Highly compatible with Kubernetes<BR>Complete with internal network<BR>Safety backup if applying CoW                                         | Easy to access backup files<BR>It fan be shared and used by multiple clusters                         | The safest of these methods                                                                                                                      |
| demerit :x: | A bit hard to check backup files<BR>Can not share backup files for several Vald clusters                                                                    | Need to communicate with external network                                                             | Need to operate both storages<BR>The most expensive way                                                                                          |

## Backup configuration

This section shows the best practice for configuring backup features with PV, S3, or PV + S3.

Each sample configuration yaml is published on [here](https://github.com/vdaas/vald/tree/master/charts/vald/values).
Please refer it for more details.

### General

Regardless of the backup destination, the following Vald Agent settings must be set to enable backup.

- `agent.ngt.index_path` is the location where those files are stored.
- `agent.ngt.in-memory-mode=false` means storing index files in the local volume.

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
  terminationGracePeriodSeconds: 3600
  ...
  ngt:
    ...
    index_path: "/var/ngt/index"
    enable_in_memory_mode: false
    # copy on write function flag
    enable_copy_on_write: true
    ...
```

Each PV will be mounted on each Vald Agent Pod's `index_path`.

You can choose `copy_on_write` (CoW) function.

The CoW is an option to update the backup file safely.

The backup file may be corrupted, and the Vald Agent pod may not restore from backup files when the Vald Agent pod terminates during saveIndex without CoW is not be enabled.
On the other hand, if CoW is enabled, the Vald Agent pod can restore the data from one generation ago.

<div class="caution">
When CoW is enabled, PV temporarily has two backup files; new and old versions.<BR>
So, A double storage capacity is required if CoW is enabled, e.g., when set 2Gi as size without CoW, the size should be more than 4Gi with CoW.
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
  terminationGracePeriodSeconds: 3600
  ...
  ngt:
    ...
    index_path: "/var/ngt/index"
    enable_in_memory_mode: false
    ...
  sidecar:
    enabled: true
    # run sidecar with initContainerMode or not
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
Please refer to the before sections for provisioning storages.

```yaml
agent:
  ...
  terminationGracePeriodSeconds: 3600
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
  ngt:
    ...
    enable_in_memory_mode: false
    # copy on write function flag
    enable_copy_on_write: true
    ...
  sidecar:
    enabled: true
    # run sidecar with initContainerMode or not. If true, it allows restoration.
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
    ...
```

## Restore

Restoring from the backup file runs on start Pod when the config is set correctly, and the backup file exists.

### PV

In using the PV case, restoration starts when Pod starts.

If the configuration is correct, the backup file will be automatically mounted, loaded, and indexed when the Pod starts.

### S3

In using the S3 case, restoration runs only `initContainerMode`.
To enable restoration, you have to set `sidecar.initContainerMode` as `true`.

Agent Sidecar tries to get the backup file from S3, unpacks it, and starts indexing.

### PV + S3

In using both the PV and S3 case, the backup file used for restoration will prioritize the file on PV.
If the backup file does not exist on the PV, the backup file will be retrieved from S3 via the Vald Agent Sidecar and restored.
