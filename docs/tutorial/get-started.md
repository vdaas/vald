# Get Started

Vald is a highly scalable distributed fast approximate nearest neighbor dense vector search engine.<br>
Vald is designed and implemented based on Cloud-Native architecture.

This article will show you how to deploy and run the Vald components on your kubernetes cluster.
Fashion-mnist is used as an example of a dataset.

## Requirements

- kubernetes: v1.19 ~
- go: v1.15 ~
- helm: v3 ~
- libhdf5 (_only required for this tutorial._)

Helm is used for deploy Vald on your kubernetes and Hdf5 is used for running example code.<br>
If helm or hdf5 is not installed, please install [helm](https://helm.sh/docs/intro/install) and [hdf5](https://www.hdfgroup.org/).

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

## Deploy and Run Vald on kubernetes

This chapter shows the way to deploy Vald using Helm and to run on your kubernetes cluster.<br>
In this tutorial, you will deploy the minimum Vald that is consisted of vald-agent-ngt, vald-lb-gateway, vald-discoverer and vald-manager-index.<br>

### Deploy

1. Clone the vdaas/vald repository

   ```bash
   git clone https://github.com/vdaas/vald.git
   cd vald
   ```

1. Confirm which cluster to deploy

   ```bash
   kubectl cluster-info
   ```

   In the sense of trying to "Get-Started", [k3d](https://k3d.io/) or [kind](https://kind.sigs.k8s.io/) are easy kubernetes tools to use.

1. Apply kubernetes metrics server

   ```bash
   kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
   kubectl wait -n kube-system --for=condition=ready pod -l k8s-app=metrics-server --timeout=600s
   ```

1. Deploy Vald using helm

   ```bash
   helm repo add vald https://vald.vdaas.org/charts
   helm install vald vald/vald --values example/helm/values.yaml
   ```

1. Verify

   ```bash
   kubectl get pods
   ```

   <details><summary>Example output</summary><br>
   If the deployment is successful, all Vald components should be running.

   ```bash
   NAME                                       READY   STATUS      RESTARTS   AGE
   vald-agent-ngt-0                           1/1     Running     0          7m12s
   vald-agent-ngt-1                           1/1     Running     0          7m12s
   vald-agent-ngt-2                           1/1     Running     0          7m12s
   vald-agent-ngt-3                           1/1     Running     0          7m12s
   vald-agent-ngt-4                           1/1     Running     0          7m12s
   vald-discoverer-7f9f697dbb-q44qh           1/1     Running     0          7m11s
   vald-lb-gateway-6b7b9f6948-4z5md           1/1     Running     0          7m12s
   vald-lb-gateway-6b7b9f6948-68g94           1/1     Running     0          6m56s
   vald-lb-gateway-6b7b9f6948-cvspq           1/1     Running     0          6m56s
   vald-manager-index-74c7b5ddd6-jrnlw        1/1     Running     0          7m12s
   ```

   </details>

### Run using example code

This chapter shows how to perform a search action in Vald with fashion-mnist dataset.

1. Port Forward

   ```bash
   kubectl port-forward deployment/vald-lb-gateway 8081:8081
   ```

1. Download dataset

   In this tutorial. we use [fashion-mnist](https://github.com/zalandoresearch/fashion-mnist) as a dataset for indexing and search query.

   ```bash
   # move to working directory
   cd example/client

   # download fashion-mnist testing dataset
   wget http://ann-benchmarks.com/fashion-mnist-784-euclidean.hdf5
   ```

1. Running example

   Vald provides multiple language client libraries such as Go, Java, Node.js, Python, and so on.<br>
   In this example, the fashion-mnist dataset will insert into the Vald cluster and perform a search using [vald-client-go](https://github.com/vdaas/vald-client-go).

   We use [`example/client/main.go`](https://github.com/vdaas/vald/blob/master/example/client/main.go) to run the example.
   This will execute 6 steps.

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
            "github.com/vdaas/vald-client-go/v1/payload"
            "github.com/vdaas/vald-client-go/v1/vald"

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

      - Recognition parameters.
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

   1. load

      - Loading from fashion-mnist dataset and set id for each vector that is loaded. This step will return the training dataset, test dataset, and ids list of ids when loading is completed with success.
          <details><summary>example code</summary><br>

        ```go
        ids, train, test, err := load(datasetPath)
        if err != nil {
            glg.Fatal(err)
        }
        ```

          </details>

   1. Create the gRPC connection and Vald client with gRPC connection.

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

   1. Insert and Index

      - Insert and Indexing 400 training datasets to the Vald agent.
          <details><summary>example code</summary><br>

        ```go
        for i := range ids [:insertCount] {
            _, err := client.Insert(ctx, &payload.Insert_Request{
                Vector: &payload.Object_Vector{
                    Id: ids[i],
                    Vector: train[i],
                },
                Config: &payload.Insert_Config{
                    SkipStrictExistCheck: true,
                },
            })
            if err != nil {
                glg.Fatal(err)
            }
            if i%10 == 0 {
                glg.Infof("Inserted %d", i)
            }
        }
        ```

          </details>

      - Wait until indexing finish.
          <details><summary>example code</summary><br>

        ```go
        wt := time.Duration(indexingWaitSeconds) * time.Second
        glg.Infof("Wait %s for indexing to finish", wt)
        time.Sleep(wt)
        ```

          </details>

   1. Search

      - Search 10 neighbor vectors for each 20 test datasets and return list of neighbor vector.

      - When getting approximate vectors, the Vald client sends search config and vector to the server via gRPC.
          <details><summary>example code</summary><br>

        ```go
        glg.Infof("Start search %d times", testCount)
        for i, vec := range test[:testCount] {
            res, err := client.Search(ctx, &payload.Search_Request){
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

   1. Remove

      - Remove 400 indexed training datasets from the Vald agent.
          <details><summary>example code</summary><br>

        ```go
        for i := range ids [:insertCount] {
            _, err := client.Remove(ctx, &payload.Remove_Request{
                Id: &payload.Object_ID{
                    Id: ids[i],
                },
            })
            if err != nil {
                glg.Fatal(err)
            }
            if i%10 == 0 {
                glg.Infof("Removed %d", i)
            }
        }
        ```

          </details>

   ```bash
   # run example
   go run main.go
   ```

### Cleanup

1. Remove the Vald pods by executing:

   ```bash
   helm uninstall vald
   ```

## Deploy and Run full Vald on kubernetes

### Deploy

This chapter will show you how to deploy using Helm and run Vald on your kubernetes cluster.<br>
This chapter uses Scylla DB as a backend data store for indexing and data backup.<br>
If you want to learn about Scylla, please refer to [the official website](https://www.scylladb.com/).

1. Clone the vdaas/vald repository

   ```bash
   git clone https://github.com/vdaas/vald.git
   cd vald
   ```

1. Confirm which cluster to deploy

   ```bash
   kubectl cluster-info
   ```

   In the sense of trying to "Get-Started", [k3d](https://k3d.io/) or [kind](https://kind.sigs.k8s.io/) are easy kubernetes tools to use.

1. Prepare Scylla DB and kubernetes metrics-server

   Deploy Scylla as a backup database.

   ```bash
   make k8s/external/scylla/deploy
   ```

   In this make command, we are deploying a lightweight Cassandra-compatible scylladb using Operator.
   <details><summary>If you're interested in this make command, take a look here for more detail of make command</summary><br>

   1. Deploy cert-manager for ScyllaDB

   ```bash
   kubectl apply -f https://github.com/jetstack/cert-manager/releases/latest/download/cert-manager.yaml
   kubectl wait -n cert-manager --for=condition=ready pod -l app=cert-manager --timeout=60s
   kubectl wait -n cert-manager --for=condition=ready pod -l app=cainjector --timeout=60s
   kubectl wait -n cert-manager --for=condition=ready pod -l app=webhook --timeout=60s
   ```

   1. Deploy ScyllaDB Operator

   ```bash
   kubectl apply -f https://raw.githubusercontent.com/scylladb/scylla-operator/master/examples/common/operator.yaml
   kubectl wait -n scylla-operator-system --for=condition=ready pod -l statefulset.kubernetes.io/pod-name=scylla-operator-controller-manager-0 --timeout=600s
   ```

   1. Deploy ScyllaDB

   ```bash
   kubectl apply -f k8s/external/scylla/scyllacluster.yaml
   kubectl wait -n scylla --for=condition=ready pod -l statefulset.kubernetes.io/pod-name=vald-scylla-cluster-dc0-rack0-0 --timeout=600s
   kubectl -n scylla get pods
   ```

   1. Configure ScyllaDB

   ```bash

   kubectl apply -f example/manifest/scylla
   kubectl wait --for=condition=complete job/scylla-init --timeout=60s
   ```

   </details>

   For documentation on scylladb operator, please refer to [here](http://operator.docs.scylladb.com/master/generic)

   Apply kubernetes metrics-server

   ```bash
   kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
   kubectl wait -n kube-system --for=condition=ready pod -l k8s-app=metrics-server --timeout=600s
   ```

1. Deploy Vald using helm

   ```bash
   helm repo add vald https://vald.vdaas.org/charts
   helm install vald vald/vald --values example/helm/values-scylla.yaml
   ```

1. Verify

   ```bash
   kubectl get pods
   ```

   <details><summary>Example output</summary><br>
   If the deployment is successful, all Vald components should be running.

   ```bash
   NAME                                       READY   STATUS      RESTARTS   AGE
   scylla-init-vhdp5                          0/1     Completed   0          7m12s
   vald-agent-ngt-0                           1/1     Running     0          7m12s
   vald-agent-ngt-1                           1/1     Running     0          7m12s
   vald-agent-ngt-2                           1/1     Running     0          7m12s
   vald-agent-ngt-3                           1/1     Running     0          7m12s
   vald-agent-ngt-4                           1/1     Running     0          7m12s
   vald-agent-ngt-5                           1/1     Running     0          7m12s
   vald-backup-gateway-68c8b4ffd4-df8zp       1/1     Running     0          6m56s
   vald-backup-gateway-68c8b4ffd4-dmwrd       1/1     Running     0          6m56s
   vald-backup-gateway-68c8b4ffd4-nm8f7       1/1     Running     0          7m12s
   vald-discoverer-7f9f697dbb-q44qh           1/1     Running     0          7m11s
   vald-lb-gateway-6b7b9f6948-4z5md           1/1     Running     0          7m12s
   vald-lb-gateway-6b7b9f6948-68g94           1/1     Running     0          6m56s
   vald-lb-gateway-6b7b9f6948-cvspq           1/1     Running     0          6m56s
   vald-manager-backup-5fb5f8dc7-h22sv        1/1     Running     0          7m12s
   vald-manager-backup-5fb5f8dc7-ncrw4        1/1     Running     0          6m56s
   vald-manager-backup-5fb5f8dc7-nzbkh        1/1     Running     0          6m56s
   vald-manager-compressor-78bf64459f-27ckg   1/1     Running     0          6m56s
   vald-manager-compressor-78bf64459f-9kl9b   1/1     Running     0          7m12s
   vald-manager-compressor-78bf64459f-dkx24   1/1     Running     0          6m56s
   vald-manager-index-74c7b5ddd6-jrnlw        1/1     Running     0          7m12s
   vald-meta-747f757bbb-9v5xz                 1/1     Running     0          7m12s
   vald-meta-747f757bbb-mpwqp                 1/1     Running     0          6m56s
   vald-meta-gateway-8c5f55dd-8fsch           1/1     Running     0          6m56s
   vald-meta-gateway-8c5f55dd-sdd5q           1/1     Running     0          7m12s
   vald-meta-gateway-8c5f55dd-vfkn6           1/1     Running     0          6m56s
   ```

   </details>

### Run using example code

This chapter shows how to perform a search action in Vald with fashion-mnist dataset.

1. Port Forward

   ```bash
   kubectl port-forward deployment/vald-meta-gateway 8081:8081
   ```

1. Download dataset

   In this tutorial. we use [fashion-mnist](https://github.com/zalandoresearch/fashion-mnist) as a dataset for indexing and search query.

   ```bash
   # move to working directory
   cd example/client

   # download fashion-mnist testing dataset
   wget http://ann-benchmarks.com/fashion-mnist-784-euclidean.hdf5
   ```

1. Running example

   Vald provides multiple language client libraries such as Go, Java, Node.js, Python, and so on.<br>
   In this example, the fashion-mnist dataset will insert into the Vald cluster and perform a search using [vald-client-go](https://github.com/vdaas/vald-client-go).

   ```bash
   # run example
   go run main.go
   ```

1. Cleanup

   Remove the Vald pods by executing:

   ```bash
   helm uninstall vald
   ```

## Deploy and Run standalone Vald Agent on kubernetes

### Deploy

This chapter will show you how to deploy a standalone Vald Agent using Helm and run it on your kubernetes cluster. <br>
This chapter uses [NGT](https://github.com/yahoojapan/ngt) as Vald Agent to perform vector insertion operation, indexing and searching operations.<br>

1. Clone the vdaas/vald repository

   ```bash
   git clone https://github.com/vdaas/vald.git
   cd vald
   ```

1. Confirm which cluster to deploy

   ```bash
   kubectl cluster-info
   ```

   In the sense of trying to "Get-Started", [k3d](https://k3d.io/) or [kind](https://kind.sigs.k8s.io/) are easy kubernetes tools to use.

1. Deploy Vald Agent using helm

   There is the [values.yaml](https://github.com/vdaas/vald/blob/master/example/helm/values-standalone-agent-ngt.yaml) to deploy standalone Vald Agent.
   Each component can be disabled by setting the value `false` to the `[component].enabled` field.
   This is useful for deploying standalone Vald Agent NGT pods.

   ```bash
   helm repo add vald https://vald.vdaas.org/charts
   helm install vald-agent-ngt vald/vald --values example/helm/values-standalone-agent-ngt.yaml
   ```

1. Verify

   ```bash
   kubectl get pods
   ```

   <details><summary>Example output</summary><br>
   If the deployment is successful, Vald Agent component should be running.

   ```bash
   NAME               READY   STATUS    RESTARTS   AGE
   vald-agent-ngt-0   1/1     Running   0          20m
   vald-agent-ngt-1   1/1     Running   0          20m
   vald-agent-ngt-2   1/1     Running   0          20m
   vald-agent-ngt-3   1/1     Running   0          20m
   ```

   </details>

### Run using example code

1. Port Forward

   ```bash
   kubectl port-forward service/vald-agent-ngt 8081:8081
   ```

1. Download dataset

   In this tutorial. we use [fashion-mnist](https://github.com/zalandoresearch/fashion-mnist) as a dataset for indexing and search query.

   ```bash
   # move to working directory
   cd example/client/agent

   # download fashion-mnist testing dataset
   wget http://ann-benchmarks.com/fashion-mnist-784-euclidean.hdf5
   ```

1. Running example

   Vald provides multiple languages client library such as Go, Java, Node.js, Python and so on.<br>
   In this example, the fashion-mnist dataset will insert into the Vald and search using [vald-client-go](https://github.com/vdaas/vald-client-go).

   We use [`example/client/agent/main.go`](https://github.com/vdaas/vald/blob/master/example/client/agent/main.go) to run the example.
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
            agent "github.com/vdaas/vald-client-go/v1/agent/core"
            "github.com/vdaas/vald-client-go/v1/vald"
            "github.com/vdaas/vald-client-go/v1/payload"

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

      - Recognition parameters.
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

   1. load

      - Loading from fashion-mnist dataset and set id for each vector that is loaded. This step will return the training dataset, test dataset, and ids list of ids when loading is completed with success.
          <details><summary>example code</summary><br>

        ```go
        ids, train, test, err := load(datasetPath)
        if err != nil {
            glg.Fatal(err)
        }
        ```

          </details>

   1. Create the gRPC connection and Vald client with gRPC connection.
      <details><summary>example code</summary><br>

      ```go
      ctx := context.Background()

      conn, err := grpc.DialContext(ctx, grpcServerAddr, grpc.WithInsecure())
      if err != nil {
          glg.Fatal(err)
      }

      client := agent.NewAgentClient(conn)
      ```

      </details>

   1. Insert and Index

      - Insert and Indexing 400 training datasets to the Vald agent.
          <details><summary>example code</summary><br>

        ```go
        for i := range ids [:insertCount] {
            if i%10 == 0 {
                glg.Infof("Inserted %d", i)
            }
            _, err := client.Insert(ctx, &payload.Insert_Request{
                Vector: &payload.Object_Vector{
                    Id: ids[i],
                    Vector: train[i],
                },
                Config: &payload.Insert_Config{
                    SkipStrictExistCheck: true,
                },
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
        wt := time.Duration(indexingWaitSeconds) * time.Second
        glg.Infof("Wait %s for indexing to finish", wt)
        time.Sleep(wt)
        ```

          </details>

      - [Optional] Indexing manually instead of waiting for auto indexing
        You can set Agent NGT configuration `auto_index_duration_limit` and `auto_index_check_duration` for auto indexing.
        In this example, you can create index manually using `CreateAndSaveIndex()` mthod in the client library.
        <detail><summary>example code</summary><br>

        ```go
        _, err = client.CreateAndSaveIndex(ctx, &payload.Control_CreateIndexRequest{
        	PoolSize: uint32(insertCount),
        })
        if err != nil {
        	glg.Fatal(err)
        }
        ```

          </detail>

   1. Search

      - Search 10 neighbor vectors for each 20 test datasets and return list of neighbor vector.

      - When getting approximate vectors, the Vald client sends search config and vector to the server via gRPC.
          <details><summary>example code</summary><br>

        ```go
        glg.Infof("Start search %d times", testCount)
        for i, vec := range test[:testCount] {
            res, err := client.Search(ctx, &payload.Search_Request){
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

   1. Remove

      - Remove indexed 400 training datasets from the Vald agent.
          <details><summary>example code</summary><br>

        ```go
        for i := range ids [:insertCount] {
            _, err := client.Remove(ctx, &payload.Remove_Request{
                Id: &payload.Object_ID{
                    Id: ids[i],
                },
            })
            if err != nil {
                glg.Fatal(err)
            }
            if i%10 == 0 {
                glg.Infof("Removed %d", i)
            }
        }
        ```

          </details>


      - Remove from the index manually instead of waiting for auto indexing.
        The removed vectors are still exist in the NGT graph index before the SaveIndex (or CreateAndSaveIndex) API is called.
        If you run the below code, the indexes will be removed completely from the Vald Agent NGT graph and the Backup file.
        <detail><summary>example code</summary><br>

        ```go
        _, err = client.SaveIndex(ctx, &payload.Empty{})
        if err != nil {
            glg.Fatal(err)
        }
        ```

          </detail>

   ```bash
   # run example
   go run main.go
   ```

1. Cleanup

   Remove the Vald Agent pods by executing:

   ```bash
   helm uninstall vald-agent-ngt
   ```
