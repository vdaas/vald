# Get Started

For who want to use or be interested in Vald, this article shows a way to deploy and run Vald.
In this article, Scylla DB is used as a backend datastore of meta and backup-manager.
Fashion-minist is used as an example for dataset.
Also, you can see basic info of Vald.

1. [About](#About)
    1. [Main Features](#Main-Features)
    2. [Requirements](#Requirements)
2. [Starting Vald on k8s cluster](#Starting-Vald-on-k8s-cluster)
    1. [Deploy](#Deploy)
    2. [Run](#Run)
3. [Advanced](#Advanced)

## About

Vald is distibuted high scalable and high-speed approximate nearest neighbor search engine.<br>
It provides search result lists for input (e.g. word, sentence, image ...) which is coverted as (also high-dimensional) vector when searching neighbor.<br>
Vald is designed as Cloud Native.
It uses World Fastest ANN Algorithm [NGT](https://github.com/yahoojapan/NGT) for search neighbors.

### Main Features

- Auto Indexing
    - ???

- Ingress Filltering
    - ???

- Horizonal Scaling
    - Vald, which is designed as Cloud Native search engine, can use cluster resource especially memory and cpu.

- Auto Index Backup
    - Vald has auto index backup with MySQL + Redis or Cassndora which is helpful for quick recovery.

- Distributed Indexing
    - Vald indexes insert vecotr to distibuted multiple NGT agent. 

- Index Replication
    - ???

- Easy to use
    - Vald can be easily installed in a few steps and accepts many kind of settings.

### Requirements

- k8s: 
- go:
- docker:
- helm: 

If helm is not installed, let see below details or [here](https://htlm.sh/docs/intro/install).

<details>
    <summary>optional installation</summray>
    <div>
    <pre>
        1. helm
            <code>
            curl https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3 | bash
            </code>
    </pre>
    </div>
</details>

## Starting Vald on k8s cluster

You have to install helm. Please check [here](https://github.com/helm/helm#install).

### Deploy

This section shows how to deploy Vald with Syclla, that is used as a datastore for meta and backup-manager.
If you want to learn about Scylla, please refer to [the official website](https://www.scylladb.com/).
Please confirm the right cluster to deploy it by running `kubectl cluster-info` before the steps below.

1. Prepare scylla database and k8s metrics-server

    At first, you have to apply scylla as backup database.

    ```bash
    cd {Vald_project_root}
    
    kubectl apply -f k8s/jobs/db/initialize/cassandra/configmap.yaml
    kubectl apply -f k8s/external/scylla
    ```

    Apply k8s metrics-server

    ```bash
    kubectl apply -f k8s/metrics/metrics-server
    ```

2. Deployment from helm template command

    ```
    helm install --values example/helm/values-scylla.yaml vald vald/vald
    ```

### Run

This chapter shows procudure of running Vald from localhost with fashion-mnist dataset.

1. Port Forward

    ```bash
    kubectl port-forward {Vald-gateway-pod-name} 8081:8081
    ```

2. Install dataset

    To search neighbors of input vector, Vald should have dataset.
    In this case, you use [fashion-mnist](https://github.com/zalandoresearch/fashion-mnist) as dataset.

    ```bash
    # move to working directory
    cd example/client
    
    # get fashion-mnist
    wget https://ann-benchmarks.com/fashion-mnist-784-euclidean.hdf5
    ```

3. Running example

    Insert, Search and Remove will be running as follows.
        - Insert: insert and indexing training dataset to Vald agent
        - Search: search neighbor vector of test vector.
        - Remove: remove indexing from Vald agent.

    ```bash
    go run main.go
    ```

## Advanced

### Customize values

Vald can be customized as desired, e.g. the number of vector dimension, the number of replica and so on.
If it needs to customized Vald, it can be applied by creating/editing `values.yaml`.
In near future, instructions of `values.yaml` will be published.

### Another way to deployment

For `Get Started`, we show the deployment procudure for Vald with scylla.
Vald can be deployed by another ways, which are used by helm or `kind` command. (There are needed compressor database, mysql + redis or casandora ).
Let try another way suited for you.

