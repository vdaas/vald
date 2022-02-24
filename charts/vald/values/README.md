# Vald

This is the sample manifests to deploy Vald.

## Deploy an only vald-agent-ngt

```
helm install --values agent-ngt-standalone.yaml vald-cluster vald/vald
```

## Deploy an vald-agent-ngt with backup features

### Backup to the persistent volume

```
helm install --values vald-backup-via-pv.yaml vald-cluster vald/vald
```

### Backup to the Amazon S3

```
helm install --values vald-backup-via-s3.yaml vald-cluster vald/vald
```

### Backup to the persistent volume and the Amazon S3

```
helm install --values vald-backup-via-pv-and-s3.yaml vald-cluster vald/vald
```
