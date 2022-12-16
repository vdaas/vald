# Vald Agent Standalone on Kubernetes

This article will show you how to deploy a standalone Vald Agent using Helm and run it on your Kubernetes cluster.

## Overview

Vald is made up of multiple microservices.
In [Get Started](../tutorial/get-started.md), you may use 4 kinds of components to deploy Vald.
In this case, you use only 1 component, `Vald Agent` that is the core component for Vald named `vald-agent-ngt`, to deploy.
The below image shows the architecture image of this case.

<img src="../../assets/docs/tutorial/vald-agent-standalone-on-k8s.svg">

<div class="warning">
Using only Vald Agent, the auto indexing function is not in use.
</div>

The 5 steps to Vald Agent Standalone on Kubernetes with Vald:

1. [Check and Satisfy the Requirements](#Requirements)
1. [Prepare Kubernetes Cluster](#Prepare-the-Kubernetes-Cluster)
1. [Deploy Vald Agent Standalone on Kubernetes cluster](#Deploy-Vald-Agent-Standalone-on-Kubernetes-Cluster)
1. [Run Example Code](#Run-Example-Code)
1. [Cleanup](#Cleanup)

## Requirements

- Kubernetes: v1.19 ~
- Go: v1.15 ~
- Helm: v3 ~
- libhdf5 (_only required for tutorial_)

Helm is used to deploying Vald on your Kubernetes, and Hdf5 decodes the sample data.<br>
If Helm or HDF5 is not installed, please install [Helm](https://helm.sh/docs/intro/install) and [HDF5](https://www.hdfgroup.org/).

<details><summary>Installation command for Helm</summary><br>

```bash
curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash
```

</details>

<details><summary>Installation command for HDF5</summary><br>

```bash
# yum
yum install -y hdf5-devel

# apt
apt-get install libhdf5-serial-dev

# homebrew
brew install hdf5
```

</details>

## Prepare the Kubernetes Cluster

1. Prepare Kubernetes cluster

   To complete get started, the Kubernetes cluster is required.<br>
   Vald will run on Cloud Services such as GKE, AWS.
   In the sense of trying to "Get-Started", [k3d](https://k3d.io/) or [kind](https://kind.sigs.k8s.io/) are easy Kubernetes tools to use.

## Deploy Vald Agent Standalone on Kubernetes Cluster

This chapter will show you how to deploy a standalone Vald Agent using Helm and run it on your Kubernetes cluster. <br>
This chapter uses [NGT](https://github.com/yahoojapan/ngt) as Vald Agent to perform vector insertion operation, indexing, and searching operation.<br>

1. Clone the repository

   To use the `deployment yaml` for deployment, let's clone [`vdaas/vald`](https://github.com/vdaas/vald.git) repository.

   ```bash
   git clone https://github.com/vdaas/vald.git && \
   cd vald
   ```

1. Deploy Vald Agent Standalone using Helm

   There is the [values.yaml](https://github.com/vdaas/vald/blob/main/example/helm/values-standalone-agent-ngt.yaml) to deploy standalone Vald Agent.
   Each component can be disabled by setting the value `false` to the `[component].enabled` field.
   This is useful for deploying standalone Vald Agent NGT pods.

   ```bash
   helm repo add vald https://vald.vdaas.org/charts && \
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

## Run Example Code

1.  Port Forward

    At first, port-forward the vald-lb-gateway is required to make request from your local environment possible.

    ```bash
    kubectl port-forward service/vald-agent-ngt 8081:8081
    ```

1.  Download dataset

    Download [Fashion-MNIST](https://github.com/zalandoresearch/fashion-mnist) that is used as a dataset for indexing and search query.

    ```bash
    # move to the work directory
    cd example/client/agent
    ```

    ```bash
    # download Fashion-MNIST testing dataset
    wget http://ann-benchmarks.com/fashion-mnist-784-euclidean.hdf5
    ```

1.  Run Example

    We use [`example/client/agent/main.go`](https://github.com/vdaas/vald/blob/main/example/client/agent/main.go) to run the example.<br>
    This example will insert and index 400 vectors into the Vald from the Fashion-MNIST dataset via gRPC.
    And then after waiting for indexing, it will request for searching the nearest vector 10 times.
    You will get the 10 nearest neighbor vectors for each search query.<br>
    Run example codes by executing the below command.

    ```bash
    # run example
    go run main.go
    ```

    <details><summary>The detailed explanation of example code is here</summary><br>
    This will execute 6 steps.

    1.  init

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

    1.  load

        - Loading from Fashion-MNIST dataset and set id for each vector that is loaded. This step will return the training dataset, test dataset, and ids list of ids when loading is completed with success.
            <details><summary>example code</summary><br>

          ```go
          ids, train, test, err := load(datasetPath)
          if err != nil {
              glg.Fatal(err)
          }
          ```

            </details>

    1.  Create the gRPC connection and Vald client with gRPC connection.
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

    1.  Insert and Index

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
          In this example, you can create index manually using `CreateAndSaveIndex()` method in the client library
          <details><summary>example code</summary><br>

              ```go
              _, err = client.CreateAndSaveIndex(ctx, &payload.Control_CreateIndexRequest{
                  PoolSize: uint32(insertCount),
              })
              if err != nil {
                  glg.Fatal(err)
              }
              ```

              </details>

    1.  Search

        - Search 10 neighbor vectors for each 20 test datasets and return a list of neighbor vectors.

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
                      Epsilon: 0.1,
                      Timeout: 100000000,
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

    1.  Remove

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
            The removed vectors still exist in the NGT graph index before the SaveIndex (or CreateAndSaveIndex) API is called.
            If you run the below code, the indexes will be removed completely from the Vald Agent NGT graph and the Backup file.
                <details><summary>example code</summary><br>

                ```go
                _, err = client.SaveIndex(ctx, &payload.Empty{})
                if err != nil {
                    glg.Fatal(err)
                }
                ```

                </details>

        </details>

    <div class="caution">
    It would be best to run CreateIndex() after Insert() without waiting for auto-indexing in your client code, even you can wait for the finishing auto createIndex function, which sometimes takes a long time.
    The backup files (e.g., ngt-meta.kvsdb) will be in your mount directory when vald-agent-ngt finishes indexing.
    </div>
      
    <div class="warning">
    If you use Go(v1.16~) and catch the error like `missing go.sum entry to add it` when running `go run main.go`, please run `go mod tidy` and retry.
    This error comes from Go Command Changes of Go 1.16 Release Notes.(Please refer to https://golang.org/doc/go1.16#go-command for more details).
    </div>

## Cleanup

In the last, you can remove all deployed Vald pods by executing the below command.

```bash
helm uninstall vald-agent-ngt
```

## Recommended Documents

Congratulation! You achieved this tutorial!

For more information, we recommend you to check:

- [Architecture](../overview/architecture.md)
- [Configuration](../user-guides/configuration.md)
- [Operations](../user-guides/operations.md)
