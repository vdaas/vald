# Get Started

This article will show you how to deploy and run Vald on your k8s cluster.
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

Vald is a distributed highly scalable and fast approximate nearest neighbor dense feature vector search engine.<br>
Vald can handle large amount of specific objects, such as image, text, video, and etc.<br>
Vald is designed base on Cloud Native.
It uses the fastest ANN Algorithm [NGT](https://github.com/yahoojapan/NGT) to search neighbors.
(If you are interested in ANN benchmarks, please refer to [the official website](http://ann-benchmarks.com/).)

### Main Features

- Auto Indexing
    - Normally, when changing the Graph Index, the Graph must be locked, but Vald uses distributed index graph, it is extremely difficult for the user to reconstruct the Graph.
    - Therefore, Vald automatically indexes distributed Graphs sequentially.

- Ingress/Egress Filltering
    - Vald implements it's own highly customizable Ingress/Egress filter.
    - Which can be configured by the user to fit the gRPC interface.
        - Ingress Filter: Ability to Vectorize through filter on request.
        - Egress Filter: rerank or filter the searching result with your own algorithm.

- Horizontal Scalable
    - Vald is a cloud-native vector search engine running on Kubernetes, which enables horizontal scalling of memory and cpu for billion scale of vector data.

- Auto Indexing Backup
    - Vald has auto index backup feature using MySQL + Redis or Cassndora which enables disaster recovery.

- Distributed Indexing
    - Vald distribute vector index to multiple agent, each agent stores different index.

- Index Replication
    - Vald stores each index in multiple agents which enables index replicas.
    - Automatically rebalance the replica when some Vald agent goes down.

- Easy to use
    - Vald can be easily installed in a few steps and Vald is highly customizable.

### Requirements

- k8s:  v1.17 ~
- go:   v1.14 ~
- helm: v3 ~
- libhdf5 (_only required for this tutorial._)

Helm and hdf5 is required for this tutorial. If helm or hdf5 is not installed, please install [helm](https://helm.sh/docs/intro/install)) and [hdf5](https://www.hdfgroup.org/)).

<details><summary>optional installation</summary><br>
install helm

```bash
curl https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3 | bash
```
install hdf5
```bash
## yum
yum install -y hdf5-devel

## apt
apt-get install libhdf5-serial-dev

## homebrew
brew install hdf5
```
</details>

## Deploy and Run Vald on K8s cluster

### Deploy

This section shows how to deploy Vald with Scylla, which is used as a datastore for index-metadata and backup-manager.
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

### Run

This chapter shows how to run Vald with fashion-mnist dataset.

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

    Vald provides multiple langurages client library such as golang, Java, Node.js, Python and so on.<br>
    In this case, we use [vald-client-go](https://github.com/vdaas/vald-client-go) which is written by golang.

    We use [`example/client/main.go`](../../example/client/main.go) to run the example.
    This will execute 4 steps.
    1. init
    - Import packages
        <details><summary>example code</summary><br>

        ```go
        package main

        import (
            "context"
            "encoding/json"
            "flag"
            "time"

            "github.com/kpango/fuid"
            "github.com/kpango/glg"
            "github.com/vdaas/vald-client-go/gateway/vald"
            "github.com/vdaas/vald-client-go/payload"

            "gonum.org/v1/hdf5"
            "google.golang.org/grpc"
        )
        ```
        </details>
    - Set variables
        - The constant number of training datasets and test datasets.
            <details><summary>example code</summary><br>

            ```go
            const (
                insertCount = 400
                testCount = 20
            )
            ```
            </details>

        - The variables for configuration.
            <details><summary>example code</summary><br>

            ```go
            const (
                datasetPath         string
                grpcServerAddr      string
                indexingWaitSeconds uint
            )
            ```
            </details>
    - Recognition paremters.
        <details><summary>example code</summary><br>

        ```go
        func init() {
            flag.StringVar(&datasetPath, "path", "fashion-mnist-784-euclidean.hdf5", "set dataset path")
            flag.StringVar(&grpcServerAddr, "addr", "127.0.0.1:8081", "set gRPC server address")
            flag.UintVar(&indexingWaitSeconds, "wait", 60, "set indexing wait seconds")
            flag.Parse()
        }
        ```
        </details>
    2. load
    - Loading from fashion-mnist dataset and set id for each vector that is loaded. This step will return the training dataset, test dataset, and ids list of ids when loading is completed with success.
        <details><summary>example code</summary><br>

        ```go
        ids, train, test, err := load(datasetPath)
        if err != nil {
            glg.Fatal(err)
        }
        ```
        </details>
    3. Create the gRPC connection and Vald client with gRPC connection.
        <details><summary>example code</summary><br>

        ```go
        ctx := context.Background()

        conn, err := grpc.DialContext(ctx, grpcServerAddr, grpc.WithInsecure())
        if err != nil {
            glg.Fatal(err)
        }

        client := vald.NewValdClient(conn)
        ```
        </details>
    4. Insert and Index
    - Insert and Indexing 400 training datasets to the Vald agent.
        <details><summary>example code</summary><br>

        ```go
        for i := range ids [:insertCount] {
            if i%10 == 0 {
                glg.Infof("Inserted %d", i)
            }
            _, err := client.Insert(ctx, &payload.Object_Vector{
                Id: ids[i],
                Vector: train[i],
            })
            if err != nil {
                glg.Fatal(err)
            }
        }
        ```
        </details>
    - Wait until indexing finish.
        <details><summary>example code</summary><br>

        ```go
        glg.Info("Wait for indexing to finish")
        time.Sleep(time.Duration(indexingWaitSeconds) * time.Second)
        ```
        </details>
    5. Search
    - Search 10 neighbor vectors for each 20 test datasets and return list of neighbor vector.
    - When getting approximate vectors, the Vald client sends search config and vector to the server via gRPC.
        <details><summary>example code</summary><br>

        ```go
        glg.Infof("Start search %d times", testCount)
        for i, vec := range test[:testCount] {
            res, err := client.Seach(ctx, &payload.Search_Request){
                Vector: vec,
                Config: &payload.Search_Config{
                    Num: 10,
                    Radius: -1,
                    Epsilon: 0.01,
                }
            }
            if err != nil {
                glg.Fatal(err)
            }

            b, _ := json.MarshalIndent(res.GetResults(), "", " ")
            glg.Infof("%d - Results : %s\n\n", i+1, string(b))
            time.Sleep(1 * time.Second)
        }
        ```
        </details>

    ```bash
    # run example
    go run main.go
    ```

## Advanced

### Customize Vald

Vald is highly customizable.
For example you can configure the number of vector dimension, the number of replica and etc.
You can customize Vald by creating/editing [`values.yaml`](../../charts/vald/values.yaml).
We will publish the description of `values.yaml` soon.

### Another way to deploy Vald

In the `Get Started` section, you learnt how to deploy Vald with Scylla DB using Helm.
Vald can be deployed by Helm or `kind` command. (Compressor datastore is required, for example mysql + redis or casandora ).
We will publish the document about the setup procedure and configuration in the future.
