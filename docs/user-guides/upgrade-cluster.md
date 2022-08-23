# Upgrade Vald cluster

Before this document, we recommend reading [the configuration document](../user-guides/configuration.md) first.

## Update Configuration

When you'd like to update the configuration for the Vald cluster, you need to edit your configuration file and apply it according to your deployment method.
Please refer to one of the sections below according to your deployment method.

<div class="caution">
If each Vald Agent already has the indexes and uses the backup feature, please DO NOT edit the dimension and the distance type to avoid initContainer error.<BR>
(Each Vald Agent will restore from backup data when performing a rolling update)
</div>

<div class="warn">
The indexes will be clear when performing a rolling update if each Vald Agent already has the indexes and does not use the backup.
</div>

### Using Helm for deployment

The configuration file is `values.yaml`. (https://helm.sh/docs/chart_template_guide/values_files/)


You can edit your `values.yaml` and update your Vald cluster by following steps:
1. Edit `values.yaml`.
1. Update with Helm command. 

    ```bash
    helm upgrade <NAME> vald/vald --values <YOUR VALUES FILE PATH>
    # e.g., helm upgrade vald vald/vald --values valeus.yaml
    ```

### Using vald-helm-operator

There are two configuration files:
- `vhor.yaml` for vald-helm-operator
- `vr.yaml` for the Vald cluster

If you don't need to update the `vald-helm-operator`, please skip step.1 and step.2.
The Vald cluster will be updated by the following steps:
1. Edit `vhor.yaml`.
1. Apply `vhor.yaml`.

    ```bash
    kubectl apply -f <VHOR YAML FILE PATH>
    # e.g., kubectl apply -f vhor.yaml 
    ```

1. Edit `vr.yaml`.
1. Apply `vr.yaml`.

    ```bash
    kubectl apply -f <VR YAML FILE PATH>
    # e.g., kubectl apply -f vr.yaml 
    ```

## Upgrade Vald version

We recommend using the latest Vald version.

If your Vald cluster uses an old version and you'd like to upgrade the Vald version, you can upgrade by the following sections.
Please refer to one of the sections below according to your deployment method.

<div class="caution">
When the upgrade request is applied, Kubernetes will perform a rolling update.
If your Vald cluster does not apply backup, the index will be cleared by a rolling update.<BR>
In addition, if each Vald Agent already has the indexes and uses the backup feature, please DO NOT edit the dimension and the distance type to avoid initContainer error.
(Each Vald Agent will restore from backup data when backup when performing a rolling update)
</div>

### Using Helm for deployment

Using the `Helm` command for the deployment of the Vald cluster, the upgrading steps are below.

1. Update Helm repo

    ```bash
    helm repo update <NAME> https://vald.vdaas.org/charts
    # e.g., helm repo update vald https://vald.vdaas.org/charts
    ```

1. Edit `values.yaml`
    - We recommend setting a specific version as an image tag.

        ```yaml
        defaults:
        ...
          image:
            tag: v1.5.6 # set a new version
        ...
        ```

1. Apply `values.yaml`

    ```bash
    helm upgrade <NAME> vald/vald --values <YOUR VALUES FILE PATH>
    # e.g., helm upgrade vald vald/vald --values valeus.yaml
    ```

### Using vald-helm-operator

The upgrading steps are below if you use `vald-helm-operator` for the deployment.

1. Upgrade CRDs

    ```bash
    kubectl replace -f https://raw.githubusercontent.com/vdaas/vald/<VERSION>/charts/vald-helm-operator/crds/valdrelease.yaml
    kubectl replace -f https://raw.githubusercontent.com/vdaas/vald/<VERSION>/charts/vald-helm-operator/crds/valdhelmoperatorrelease.yaml
    ```

1. Update `vhor`

    ```bash
    kubectl patch vhor vhor-release -p '{"spec":{"image":{"tag":"<VERSION>"}}}'
    ```

    - Also, you can upgrade manually by editing `vhor.yaml` and applying it.

1. Edit `vr.yaml`
    - We recommend setting a specific version as an image tag.

        ```yaml
        ...
        spec:
          defaults:
            ...
            image:
              tag: v1.5.6 # set a new version
            ...
        ```

1. Apply `vr.yaml`

    ```bash
    kubectl apply -f <VR YAML FILE PATH>
    # e.g., kubectl apply -f vr.yaml 
    ```
