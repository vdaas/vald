# Get Started

For the one who is interested in Vald, this article will help you a way to deploy and run Vald on k8s.
In this article, Scylla DB is used as a backend datastore for index-metadata and backup-manager.
Fashion-mnist is used as an example dataset.

1. [About](#About)
    1. [Main Features](#Main-Features)
    2. [Requirements](#Requirements)
2. [Starting Vald on k8s cluster](#Starting-Vald-on-k8s-cluster)
    1. [Deploy](#Deploy)
    2. [Run](#Run)
3. [Advanced](#Advanced)

## About

Vald is distributed high scalable and high-speed approximate nearest neighbor search engine.<br>
It provides the searching result of any input (e.g. word, sentence, image ...) which is coverted to a multi-dimensional vector when searching the neighbor.<br>
Vald is designed base on Cloud Native.
It uses the fastest ANN Algorithm [NGT](https://github.com/yahoojapan/NGT) to search neighbors.

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

If helm is not installed, please install helm (see below details) or [here](https://helm.sh/docs/intro/install).

<details>
    <summary>optional installation</summray>

```bash
# install helm
curl https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3 | bash

# install hdf5
yum install -y hdf5-devel
apt-get install libhdf5-serial-dev
brew install hdf5
```
</details>

## Starting Vald on k8s cluster

You have to install helm. Please check [here](https://github.com/helm/helm#install).

### Deploy

This section shows how to deploy Vald with Scylla, which is used as a datastore for index-metadata and backup-manager.
If you want to learn about Scylla, please refer to [the official website](https://www.scylladb.com/).

1. Confirm the right cluster to deploy

    ```bash
    kubectl cluster-info
    ```

2. Prepare scylla database and k8s metrics-server

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

    Insert, Search and Remove will be running as follows.<br>
    - Insert: insert and indexing training dataset to Vald agent
    - Search: search neighbor vector of test vector.
    - Remove: remove indexing from Vald agent.

    ```bash
    # if you don't install hdf5-devel, please install hdf5-devl at first.
    # yum install -y hdf5-devel
    # apt-get install libhdf5-serial-dev
    # brew install hdf5

    # run example
    go run main.go
    ```

## Advanced

### Customize Vald

Vald is highly customizable.
For example you can configure the number of vector dimension, the number of replica and etc.
You can customiz Vald by creating/editing `values.yaml`.
We will publish the instructions of `values.yaml` soon.

### Another way to deploy Vlad

In the `Get Started` section, we show how to deploy Vald with Scylla.
Vald can be deployed in another way, which is used by helm or `kind` command. (Compressor datastore is required, for example mysql + redis or casandora ).

