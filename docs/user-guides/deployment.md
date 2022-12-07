# Deployment

There are two major ways for the deployment of the Vald cluster ways: Using the Helm command with `values.yaml` or without Helm command using operator called `vald-helm-operator`.

- Using Helm command with `values.yaml`

  - Easy to deploy
  - Allow editing Vald configuration values when the user executes Helm command with inlining.
  - Need Helm command when applying configuration

- Using `vald-helm-operator`
  - Monitoring the Vald deployments and Vald custom resource by `vald-helm-operator`
  - Use `kubectl` command to set or update of the Vald cluster configuration
  - Automate manage the Vald cluster based on CRD (Custom Resource Definitions)

<div class="notice">
The benefit of using vald-helm-operator, it eliminates the need to check the status of the Vald cluster from the Vald cluster operator.
It manages the Vald cluster based on the designed configuration instead of you.
</div>

## Requirement

- Helm: v3 ~

If Helm is not installed, please install [Helm](https://helm.sh/docs/intro/install).

<details><summary>Installation command for Helm</summary><br>

```bash
curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash
```

</details>

## Before Deployment

Both ways need the configuration file when you need to overwrite the default configurations before the deployment.

Please refer to [the basic configuration document](../user-guides/configuration.md).

In addition, as you need, please refer to the followings:

- [Backup Configuration](../user-guides/backup-configuration.md)
- [Filter Configuration](../user-guides/filtering-configuration.md)

Moreover, we publish the example configuration files for each use case.
Please also refer to [here](https://github.com/vdaas/vald/tree/main/charts/vald/values).

## Use Helm command

### Deployment

When deploying with Helm command, it requires `value.yaml` to override the default values.

<details><summary>Sample values YAML</summary><br>

```yaml
defaults:
  logging:
    level: debug
  image:
    # Please set the specified version, e.g., v1.5.6, instead of latest
    tag: "latest"
  server_config:
    healths:
      liveness:
        livenessProbe:
          initialDelaySeconds: 60
      readiness:
        readinessProbe:
          initialDelaySeconds: 60

  ## vald-lb-gateway settings
  gateway:
    lb:
      minReplicas: 2
      maxReplicas: 2
      gateway_config:
        index_replica: 2

  ## vald-agent settings
  agent:
    minReplicas: 6
    maxReplicas: 6
    podManagementPolicy: Parallel
    ngt:
      dimension: 784
      distance_type: l2
      object_type: float
      # When auto_index_check_duration_limit is minus value, the agent auto indexing is effectively disabled.
      auto_index_check_duration_limit: "-1s"
      # When auto_index_duration_limit is minus value, the agent auto indexing is effectively disabled.
      auto_index_duration_limit: "-1s"
      auto_create_index_pool_size: 10000
      default_pool_size: 10000

  ## vald-discoverer settings
  discoverer:
    resources:
      requests:
        cpu: 150m
        memory: 50Mi

  ## vald-manager settings
  manager:
    index:
      resources:
        requests:
          cpu: 150m
          memory: 30Mi
      indexer:
        auto_index_duration_limit: 1m
        auto_index_check_duration: 40s
```

</details>

After create `values.yaml`, you can deploy by the following steps.

1. Add vald repo into the helm repo.

   ```bash
   helm repo add vald https://vdaas.vald.org/charts
   ```

1. Deploy with the `values.yaml`

   ```bash
   helm install vald vald/vald --values <YOUR VALUES YAML FILE PATH>
   ```

### Update Configuration

When you need to update the configuration, you can update by following command with your new `values.yaml`.

```bash
helm upgrate vald vald/vald --values <YOUR NEW VALUES FILE PATH>
```

### Cleanup

The Vald cluster can be removed by the following command.

```bash
helm uninstall vald
```

<div class="caution">
If using PV for the backup, PV won't delete automatically.
</div>

## Using with vald-helm-operator

### Deployment

When deploying with vald-helm-operator, it requires `vr.yaml` file for applying `ValdRelease`.
`vald-helm-operator` manages the Vald cluster based on configuration of `vr.yaml`.

<details><summary>Sample ValdRelease YAML</summary><br>

```yaml
apiVersion: vald.vdaas.org/v1
kind: ValdRelease
metadata:
  name: vald-cluster
# the values of Helm chart for Vald can be placed under the `spec` field.
spec:
  defaults:
    logging:
      level: debug
    image:
      # Please set the specified version, e.g., v1.5.6, instead of latest
      tag: "latest"
    server_config:
      healths:
        liveness:
          livenessProbe:
            initialDelaySeconds: 60
        readiness:
          readinessProbe:
            initialDelaySeconds: 60

    ## vald-lb-gateway settings
    gateway:
      lb:
        minReplicas: 2
        maxReplicas: 2
        gateway_config:
          index_replica: 2

    ## vald-agent settings
    agent:
      minReplicas: 6
      maxReplicas: 6
      podManagementPolicy: Parallel
      ngt:
        dimension: 784
        distance_type: l2
        object_type: float
        # When auto_index_check_duration_limit is minus value, the agent auto indexing is effectively disabled.
        auto_index_check_duration_limit: "-1s"
        # When auto_index_duration_limit is minus value, the agent auto indexing is effectively disabled.
        auto_index_duration_limit: "-1s"
        auto_create_index_pool_size: 10000
        default_pool_size: 10000

    ## vald-discoverer settings
    discoverer:
      resources:
        requests:
          cpu: 150m
          memory: 50Mi

    ## vald-manager settings
    manager:
      index:
        resources:
          requests:
            cpu: 150m
            memory: 30Mi
        indexer:
          auto_index_duration_limit: 1m
          auto_index_check_duration: 40s
```

</details>

In addition, you can use the operator for the `vald-helm-operator` by applying `vhor.yaml` for `ValdHelmOperatorRelease`.
It is possible to auto-managing vald-helm-operator.

<details><summary>Sample ValdHelmOperatorRelease YAML</summary><br>

```yaml
apiVersion: vald.vdaas.org/v1
kind: ValdHelmOperatorRelease
metadata:
  name: vald-helm-operator-release
# the values of Helm chart for vald-helm-operator can be placed under the `spec` field.
spec:
  watchNamespaces: "default"
```

</details>

For more details of the configuration of vald-helm-operator-release, please refer to [here](https://github.com/vdaas/vald/tree/main/charts/vald-helm-operator#configuration).

After setting `vr.yaml` (and `vhor.yaml`), you can deploy by the following steps.

<div class="warning">
If you need to deploy ValdHelmOperatorRelease, you should apply vhor.yaml before applying vr.yaml.
</div>

1. Add vald repo into the helm repo

   ```bash
   helm repo add vald https://vdaas.vald.org
   ```

1. Deploy `vald-helm-operator`

   ```bash
   helm install vald-helm-operator-release vald/vald-helm-operator
   ```

1. Apply `vhor.yaml` (optional)

   ```bash
   kubectl apply -f vhor.yaml
   ```

   After deployment success, the Kubernetes cluster runs the rolling update for vald-helm-operator components.

1. Apply `vr.yaml`

   ```bash
   kubectl apply -f vr.yaml
   ```

### Update Configuration

When you need to update the configuration, you can update by following command with your new `vr.yaml` or `vhor.yaml`.

```bash
kubectl apply -f <new vr.yaml or new vhor.yaml>
```

### Cleanup

The Vald cluster can be removed by the following steps.

1. Delete the `ValdRelease`

   ```bash
   kubectl delete -f vr.yaml
   ```

1. Delete the `ValdHelmOperatorReleases`

   ```bash
   kubectl delete -f vhor.yaml
   ```

1. Delete crd

   ```bash
   kubectl patch crd/valdhelmopratorreleases.vald.vdaas.org -p '{"metadata":{"finalizers":[]}}' -type=merge
   ```

<div class="notice">
If using PV for the backup, PV won't delete automatically.
</div>
