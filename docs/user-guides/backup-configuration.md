# Backup configuration

There are three types of options for the Vald cluster: backup, filtering, and core algorithms.

This page describes how enabled backup features for your Vald cluster.

## What is the backup

Vald's backup function is to save the index information indexed in each Vald Agent pod as file data in Persistent Volume or S3.
When the Vald Agent pod is restarted for some reason, the index state is restored from the saved index data.

## Backup configuration

This section shows the best practice backup configuration with PV, S3, or PV + S3.

Each sample configuration yaml is published on [here](https://github.com/vdaas/vald/tree/master/charts/vald/values).
Please try to reference.

### General

Regardless of the backup destination, the following Vald Agent settings must be set to enable backup.
`in-memory-mode=false` means storing index files in the local volume.
`index_path` is the location where those files are stored.

### Persistent Volume

#### requirement

You must prepare PV before deployment when using Kubernetes Persistent Volume (PV) for backup storage.
Please refer to the setup guide of the usage environment for the provisioning PV.

For example:
    - [GKE setup PV document](https://cloud.google.com/kubernetes-engine/docs/concepts/persistent-volumes)
    - [EKS storage document](https://docs.aws.amazon.com/eks/latest/userguide/storage.html)

#### configuration

After provisioning PV, the following parameters are needed to be set.

<!-- TODO: yaml -->

Each PV will be mounted on each Vald Agent Pod's `index_path`.

It is highly recommended to set `copy_on_write` (CoW) as `true`. 

The CoW is an option to update the backup file safely.
The backup file may be corrupted, and the Vald Agent pod may not restore from backup files when the Vald Agent pod terminates during saveIndex without CoW is not be enabled.
On the other hand, when CoW is enabled, the Vald Agent pod can restore the data from one generation ago.

### S3

#### requirement

Before deployment, you must provision the S3 object storage.
You can use any S3-compatible object storage.
    - [AWS S3](https://aws.amazon.com/s3/)
    - [Google Cloud Storage](https://cloud.google.com/storage/docs/)

#### configuration

After provisioning the object storage, the following parameters are needed to be set.
To enable the backup function with S3, the Vald Agent Sidecar should be enabled.

<!-- YAML -->

The Vald Agent Sidecar needs an access key and a secret access key to communicate with your object storage.
Before applying the helm chart, register each value in Kubernetes secrets with the following commands.

<!-- command -->

### Persistent Volume and S3

You can use both PV and S3 at the same time.

<!-- Example YAML -->

## Restore

Restoring from the backup file runs on initContainer when the config is set correctly and the backup file exists.

If you use both PV and S3, the backup file used for restoration will prioritize the file on PV.
If the backup file does not exist on the PV, the backup file will be retrieved from S3 via the Vald Agent Sidecar and restored.
