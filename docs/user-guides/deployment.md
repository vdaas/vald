# Deployment

There are two major ways for the deployment of the Vald cluster ways: using the helm command with `values.yaml` or vald-helm-operator.



using helm command
easy to deploy
need to helm command when upgrade configuration
vald-helm-operator
monitoring the Vald deployments and Vald custom resource
Use `kubectl` command to upgrade the Vald cluster configuration
Automate manage the Vald cluster based on CRD (Custom Resource Definitions) 

## Requirement

- Helm: v3 ~

If Helm or HDF5 is not installed, please install [Helm](https://helm.sh/docs/intro/install).

<details><summary>Installation command for Helm</summary><br>

```bash
curl https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3 | bash
```

</details>

## Before Deployment

Both ways need the configuration file when you need to overwrite the default configurations.

Please refer to [the basic configuration document](../user-guides/configuration.md).

In addition, as you need, please refer to the followings:
- [Backup Configuration](../user-guides/backup-configuration.md)
- [Filter Configuration](../user-guides/filter-configuration.md)

## Deploy with Helm command

When deploying with Helm command, it requires `value.yaml` to override the default values.

<details><summary>Sample `values.yaml`</summary><br>

```yaml
defaults:
  logging:
    level: debug
  image:
    tag: "v1.5.6"
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
  
  ## vald-agent settings
  agent:
    minReplicas: 3
    maxReplicas: 3
    podManagementPolicy: Parallel
    ngt:
      dimension: 784
      distance_type: l2
      object_type: float
      auto_index_check_duration_limit: 720h
      auto_index_duration_limit: 24h
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

## Deploy with vald-helm-operator

When deploying with vald-helm-operator, it requires `vr.yaml` file for applying `ValdRelease`.
`vald-helm-operator` manages the Vald cluster based on configuration of `vr.yaml`.

<details><summary>Sample `vr.yaml`</summary><br>

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
     tag: "v1.5.6"
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
   
   ## vald-agent settings
   agent:
     minReplicas: 3
     maxReplicas: 3
     podManagementPolicy: Parallel
     ngt:
       dimension: 784
       distance_type: l2
       object_type: float
       auto_index_check_duration_limit: 720h
       auto_index_duration_limit: 24h
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


After create `vr.yaml`, you can deploy by the following steps.

1. Add vald repo into the helm repo

    ```bash
    helm repo add vald https://vdaas.vald.org
     ```

1. Deploy `vald-helm-operator`

    ```bash
    helm install vald-helm-operator-release vald/vald-helm-operator
    ```

1. Apply `vr.yaml`

    ```bash
    kubectl apply -f vr.yaml
    ```

If you need to auto managing to vald-helm-operator, you must apply `vhor.yaml` for applying `ValdHelmOpearterRelease`.

<details><summary>Sample `vhor.yaml`</summary><br>

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

For more details of the configuration of vald-helm-operator-release, please refer to [here](https://github.com/vdaas/vald/tree/master/charts/vald-helm-operator#configuration).

1. Apply `vhor.yaml`

    ```bash
    kubectl apply -f valdHelmOperatorRelease yaml
    ```
