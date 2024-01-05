# Vald Read Replica

**THIS CHART IS A WORK IN PROGRESS AND IS NOT YET FUNCTIONAL**

This is a Helm chart to install Vald readreplica components.

Current chart version is `v1.7.10`

## Install

Add Vald Helm repository

    $ helm repo add vald https://vald.vdaas.org/charts

Install Vald cluster first.

    $ helm install vald vald/vald

Run the following command to install the chart,

    $ helm install vald-readreplica vald/vald-readreplica

## Configuration

### Overview

[`values.yaml`](https://github.com/vdaas/vald/blob/main/charts/vald-readreplica/values.yaml) of this chart is a symbolic link to the [`values.yaml`](https://github.com/vdaas/vald/blob/main/charts/vald/values.yaml) of the main vald chart
because all the configurations must be synced with the main vald cluster.
So please look at the document of the main vald chart for configurations.

When you deploy this chart with custom `values.yaml` on install, you should deploy the vald
cluster with the same `values.yaml` as well.
