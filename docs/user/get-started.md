# Get Started

This article will show you how to deploy and run Vald on your k8s cluster.
This article uses Scylla DB as the backend data store for metadata-management and backup-manager.
Fashion-mnist is used as an example of a dataset.

- [Get Started](#get-started)
  - [About](#about)
    - [Requirements](#requirements)
  - [Deploy and Run Vald on K8s cluster](#deploy-and-run-vald-on-k8s-cluster)
    - [Deploy](#deploy)
    - [Example](#example)
  - [Another way to deploy Vald](#another-way-to-deploy-vald)

## About


### Requirements

- k8s:  v1.17 ~
- go:   v1.14 ~
- helm: v3 ~
- libhdf5 (_only required for this tutorial._)

Helm and hdf5 is required for this tutorial. If helm or hdf5 is not installed, please install [helm](https://helm.sh/docs/intro/install) and [hdf5](https://www.hdfgroup.org/).

<details><summary>[Optional] Install helm</summary><br>

```bash
curl https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3 | bash
```
</details>

<details><summary>[Optional] Install hdf5</summary><br>


```bash
# yum
yum install -y hdf5-devel

# apt
apt-get install libhdf5-serial-dev

# homebrew
brew install hdf5
```
</details>

## Deploy and Run Vald on K8s cluster

### Deploy

This section shows how to deploy Vald with Scylla using Helm. Scylla is used as a datastore for index-metadata and backup-manager.
If you want to learn about Scylla, please refer to [the official website](https://www.scylladb.com/).

1. Confirm which cluster to deploy

    ```bash
    kubectl cluster-info
    ```

2. Prepare Scylla DB and k8s metrics-server

    Deploy Scylla as a backup database.

    ```bash
    kubectl apply -f k8s/jobs/db/initialize/cassandra/configmap.yaml
    kubectl apply -f k8s/external/scylla
    ```

    Apply k8s metrics-server

    ```bash
    kubectl apply -f k8s/metrics/metrics-server
    ```

3. Deploy Vald using helm

    ```bash
    helm repo add vald https://vald.vdaas.org/charts
    helm install vald vald/vald --values example/helm/values-scylla.yaml
    ```

4. Verify

    ```bash
    kubectl get pods
    ```

    <details><summary>Example output</summary><br>
    If the deployment is successful, all Vald components should be running.

    ```bash
    NAME                                       READY   STATUS    RESTARTS   AGE
    scylla-0                                   1/1     Running   0          13m
    scylla-1                                   1/1     Running   0          12m
    scylla-2                                   1/1     Running   0          10m
    vald-agent-ngt-0                           1/1     Running   0          5m49s
    vald-agent-ngt-1                           1/1     Running   0          5m49s
    vald-agent-ngt-2                           1/1     Running   0          5m49s
    vald-agent-ngt-3                           1/1     Running   0          5m49s
    vald-agent-ngt-4                           1/1     Runnnig   0          5m49s
    vald-discoverer-97c88678b-wj6xn            1/1     Running   0          5m49s
    vald-gateway-5bf95f8d97-2v76g              1/1     Running   0          5m49s
    vald-gateway-5bf95f8d97-5wtb2              1/1     Running   0          78s
    vald-gateway-5bf95f8d97-7d6j7              1/1     Running   0          78s
    vald-gateway-5bf95f8d97-gx45c              1/1     Running   0          5m49s
    vald-gateway-5bf95f8d97-kx2c5              1/1     Running   0          78s
    vald-gateway-5bf95f8d97-np2lc              1/1     Running   0          5m49s
    vald-manager-backup-6c9695b69b-9xngp       1/1     Running   0          5m49s
    vald-manager-backup-6c9695b69b-jvwft       1/1     Running   0          5m49s
    vald-manager-backup-6c9695b69b-mjs2r       1/1     Running   0          5m49s
    vald-manager-compressor-6c95bdbfb5-m5t7t   1/1     Running   0          5m49s
    vald-manager-compressor-6c95bdbfb5-q8hc6   1/1     Running   0          5m49s
    vald-manager-compressor-6c95bdbfb5-zp8hb   1/1     Running   0          5m49s
    vald-manager-index-59676f54bb-nzfwt        1/1     Running   0          5m49s
    vald-meta-559744db-bcrdw                   1/1     Running   0          5m49s
    vald-meta-559744db-hz7gd                   1/1     Running   0          5m49s
    ```
    </details>

### Example

This chapter shows how to perform a search action in Vald with fashion-mnist dataset.

1. Port Forward

    ```bash
    kubectl port-forward deployment/vald-gateway 8081:8081
    ```

2. Download dataset

    In this tutorial. we use [fashion-mnist](https://github.com/zalandoresearch/fashion-mnist) as a dataset for indexing and search query.

    ```bash
    # move to working directory
    cd example/client
    
    # download fashion-mnist testing dataset
    wget http://ann-benchmarks.com/fashion-mnist-784-euclidean.hdf5
    ```

3. Perform a search action in Vald

    In this example, the fashion-mnist dataset will insert into the vald cluster and perform a search using [vald-client-go](https://github.com/vdaas/vald-client-go).
    Full example code is [`here`](../../example/client/main.go) and the explaination is [here](https://github.com/vdaas/vald-client-go/docs/users/code-explanation.md).

    ```bash
    # run example
    go run main.go
    ```

## Another way to deploy Vald

In the `Get Started` section, you learnt how to deploy Vald with Scylla DB using Helm.
Vald can be deployed by Helm or `kind` command. (Compressor datastore is required, for example mysql + redis or casandora ).
We will publish the document about the setup procedure and configuration in the future.
