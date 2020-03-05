# Get Started

For the one who is interested in Vald, this article will show you a way to deploy and run Vald on k8s.
This article uses Scylla DB as the backend data store for metadata-management and backup-manager.
Fashion-mnist is used as an example of a dataset.

1. [About](#About)
    1. [Main Features](#Main-Features)
    2. [Requirements](#Requirements)
2. [Starting Vald on k8s cluster](#Starting-Vald-on-k8s-cluster)
    1. [Deploy](#Deploy)
    2. [Run](#Run)
3. [Advanced](#Advanced)

## About

Vald is a distributed highly scalable and fast approximate nearest neighbor dense vector search engine.<br>
It provides the search result of any input (e.g. word, sentence, image and etc.) which is coverted to a multi-dimensional vector when searching the neighbor.<br>
Vald is designed base on Cloud Native.
It uses the fastest ANN Algorithm [NGT](https://github.com/yahoojapan/NGT) to search neighbors.
(If you are interested in ANN benchmarks, please refer to [the official website](http://ann-benchmarks.com/).)

### Main Features

- Auto Indexing
    - Normally, when changing the Graph Index, the Graph must be locked, but Vald uses distributed index graph, it is extremely difficult for the user to reconstruct the Graph.
    - Therefore, Vald automatically indexes distributed Graphs sequentially.

- Ingress/Egress Filltering
    - Vald has flexible customizability, which is the Ingress/Egress filter.
    - These can be freely configured by the user to fit the Vald grpc interface.
        - Ingress Filter: Ability to Vectorize through filter on request.
        - Egress Filter: A function to rerank or filter the Search response with your own algorithm.

- Horizontal Scaling
    - Vald is a cloud-native vector search engine running on Kubernetes, which enables horizontal scalling of memory and cpu for billion scale of vector data.

- Auto Index Backup
    - Vald has auto index backup feature using MySQL + Redis or Cassndora which enables disaster recovery.

- Distributed Indexing
    - Vald indexes vecotr to distributed multiple agent. It means each agent has different graph index.

- Index Replication
    - Vald stores each index in multiple agents which enables index replicas.
    - Rebalance replica when some pods go down.

- Easy to use
    - Vald can be easily installed in a few steps and Vald is highly customizable.

### Requirements

- k8s: 
- go: 
- helm: 
- hdf5: 

If helm/hdf5 is not installed, please install helm (see below details or [here](https://helm.sh/docs/intro/install))/hdf5 (see below details or [here](https://www.hdfgroup.org/)).

<details><summary>optional installation</summary><br>
install helm
<pre>
curl https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3 | bash
</pre>
install hdf5
<pre>
yum install -y hdf5-devel
apt-get install libhdf5-serial-dev
brew install hdf5
</pre></details>

## Starting Vald on k8s cluster

### Deploy

This section shows how to deploy Vald with Scylla, which is used as a datastore for index-metadata and backup-manager.
If you want to learn about Scylla, please refer to [the official website](https://www.scylladb.com/).

1. Confirm the right cluster to deploy

    ```bash
    kubectl cluster-info
    ```

2. Prepare Scylla DB and k8s metrics-server

    At first, you have to apply scylla as backup database.

    ```bash
    cd path/to/vald/root 
    
    kubectl apply -f k8s/jobs/db/initialize/cassandra/configmap.yaml
    kubectl apply -f k8s/external/scylla
    ```

    Apply k8s metrics-server

    ```bash
    kubectl apply -f k8s/metrics/metrics-server
    ```

3. Deployment from helm template command

    ```bash
    helm repo add vald https://vald.vdaas.org/charts
    helm install vald vald/vald --values example/helm/values-scylla.yaml
    ```

4. Verify

    ```bash
    kubectl get pods
    ```

<details><summary>output example</summary><br>
If the deployment is successful, you can check the following information.
<pre>
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
</pre></details>

### Run

This chapter shows the procudure of run Vald with fashion-mnist dataset.

1. Port Forward

    ```bash
    kubectl port-forward deployment/vald-gateway 8081:8081
    ```

2. Download dataset

    In this tutorial. we use [fashion-mnist](https://github.com/zalandoresearch/fashion-mnist) as a dataset for indexing and search query.

    ```bash
    # move to working directory
    cd example/client
    
    # get fashion-mnist
    wget http://ann-benchmarks.com/fashion-mnist-784-euclidean.hdf5
    ```

3. Running example

    Vald provides clients that support multiple langurages such as Java, Node.js, Python and so on.<br>
    In this case, we use [vald-client-go](https://github.com/vdaas/vald-client-go) which is written by golang.

    We use `example/client/main.go` for running example.
    This will execute 4 steps.
    1. init
    - Recognition paremters.
    2. load
    - Loading from fashion-mnist dataset and set id for each vector that is loaded. This step will return the training dataset, test dataset, and ids list of ids when loading is completed with success.
    3. insert
    - Insert and Indexing training dataset to Vald agent.
    4. search
    - Seach neighbor vector fot test vector and return list of neighbor vector.

    ```bash
    # run example
    go run main.go
    ```

## Advanced

### Customize Vald

Vald is highly customizable.
For example you can configure the number of vector dimension, the number of replica and etc.
You can customize Vald by creating/editing `values.yaml`.
We will publish the instructions of `values.yaml` soon.

### Another way to deploy Vald

In the `Get Started` section, we'll show you how to deploy Vald with Scylla DB.
Vald can be deployed in another way, which is used by Helm or `kind` command. (Compressor datastore is required, for example mysql + redis or casandora ).

