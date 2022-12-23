# Get Started

Vald is a highly scalable distributed fast approximate nearest neighbor dense vector search engine.<br>
Vald is designed and implemented based on Cloud-Native architecture.

This tutorial shows how to deploy and run the Vald components on your Kubernetes cluster.
And, Fashion-MNIST is used as an example of a dataset.

## Overview

The below image shows the architecture image about the deployment result of `Get Started`.<br>
The 4 kinds of components, `Vald LB Gateway`, `Vald Discoverer`, `Vald Agent`, and `Vald Index Manager` will be deployed to the Kubernetes.
For more information about Vald's architecture, please refer to [Architecture](../overview/architecture.md).

<img src="../../assets/docs/tutorial/getstarted.svg" />

The 5 steps to Get Started with Vald:

1. [Check and Satisfy the Requirements](#Requirements)
1. [Prepare Kubernetes Cluster](#Prepare-the-Kubernetes-Cluster)
1. [Deploy Vald on Kubernetes Cluster](#Deploy-Vald-on-Kubernetes-Cluster)
1. [Run Example Code](#Run-Example-Code)
1. [Cleanup](#Cleanup)

## Requirements

- Kubernetes: v1.19 ~
- Go: v1.15 ~
- Helm: v3 ~
- libhdf5 (_only required for get started_)

Helm is used to deploying Vald on your Kubernetes and HDF5 is used to decode the sample data file to run the example.<br>
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

This tutorial requires the Kubernetes cluster.<br>
Vald runs on public Cloud Kubernetes Services such as GKE, EKS.
In the sense of trying to `Get Started`, [k3d](https://k3d.io/) or [kind](https://kind.sigs.k8s.io/) are easy Kubernetes tools to use.

This tutorial uses Kubernetes Ingress and [Kubernetes Metrics Server](https://github.com/kubernetes-sigs/metrics-server) for running Vald.<br>
Please make sure these functions are available.<br>

The configuration of Kubernetes Ingress is depended on your Kubernetes cluster's provider.
Please refer to on yourself.

The way to deploy Kubernetes Metrics Service is here:

```bash
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml && \
kubectl wait -n kube-system --for=condition=ready pod -l k8s-app=metrics-server --timeout=600s
```

## Deploy Vald on Kubernetes Cluster

This chapter shows how to deploy Vald using Helm and run it on your Kubernetes cluster.<br>
In this tutorial, you will deploy the basic configuration of Vald that is consisted of vald-agent-ngt, vald-lb-gateway, vald-discoverer, and vald-manager-index.<br>

1. Clone the repository

   ```bash
   git clone https://github.com/vdaas/vald.git && \
   cd vald
   ```

1. Confirm which cluster to deploy

   ```bash
   kubectl cluster-info
   ```

1. Edit Configurations

   Set the parameters for connecting to the vald-lb-gateway through Kubernetes ingress from the external network.
   Please set these parameters.

   ```bash
   vim example/helm/values.yaml
   ===
   ## vald-lb-gateway settings
   gateway:
     lb:
       ...
       ingress:
         enabled: true
         # TODO: Set your ingress host.
         host: localhost
         # TODO: Set annotations which you have to set for your k8s cluster.
         annotations:
           ...
   ```

   Note:<br>
   If you decided to use port-forward instead of ingress, please set `gateway.lb.ingress.enabled` to `false`.

1. Deploy Vald using Helm

   Add vald repo into the helm repo.

   ```bash
   helm repo add vald https://vald.vdaas.org/charts
   ```

   Deploy vald on your Kubernetes cluster.

   ```bash
   helm install vald vald/vald --values example/helm/values.yaml
   ```

1. Verify

   When finish deploying Vald, you can check the Vald's pods status following command.

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

   ```bash
   kubectl get ingress
   ```

   <details><summary>Example output</summary><br>

   ```bash
   NAME                      CLASS    HOSTS       ADDRESS        PORTS   AGE
   vald-lb-gateway-ingress   <none>   localhost   192.168.16.2   80      7m43s
   ```

   </details>

   ```bash
   kubectl get svc
   ```

   <details><summary>Example output</summary><br>

   ```bash
   NAME                 TYPE        CLUSTER-IP   EXTERNAL-IP   PORT(S)             AGE
   kubernetes           ClusterIP   10.43.0.1    <none>        443/TCP             9m29s
   vald-agent-ngt       ClusterIP   None         <none>        8081/TCP,3001/TCP   8m48s
   vald-discoverer      ClusterIP   None         <none>        8081/TCP,3001/TCP   8m48s
   vald-manager-index   ClusterIP   None         <none>        8081/TCP,3001/TCP   8m48s
   vald-lb-gateway      ClusterIP   None         <none>        8081/TCP,3001/TCP   8m48s
   ```

   </details>

## Run Example Code

In this chapter, you will execute insert, search, and delete vectors to your Vald cluster using the example code.<br>
The [Fashion-MNIST](https://github.com/zalandoresearch/fashion-mnist) is used as a dataset for indexing and search query.

The example code is implemented in Go and using [vald-client-go](https://github.com/vdaas/vald-client-go), one of the official Vald client libraries, for requesting to Vald cluster.
Vald provides multiple language client libraries such as Go, Java, Node.js, Python, etc.
If you are interested, please refer to [SDKs](../user-guides/sdks.md).<br>

1.  Port Forward(option)

    If you do not use Kubernetes Ingress, port-forward is required to make requests from your local environment.

    ```bash
    kubectl port-forward deployment/vald-lb-gateway 8081:8081
    ```

1.  Download dataset

    Download [Fashion-MNIST](https://github.com/zalandoresearch/fashion-mnist) that is used as a dataset for indexing and search query.

    Move to the working directory

    ```bash
    cd example/client
    ```

    Download Fashion-MNIST testing dataset

    ```bash
    wget http://ann-benchmarks.com/fashion-mnist-784-euclidean.hdf5
    ```

1.  Run Example

    We use [`example/client/main.go`](https://github.com/vdaas/vald/blob/main/example/client/main.go) to run the example.<br>
    This example will insert and index 400 vectors into the Vald from the Fashion-MNIST dataset via [gRPC](https://grpc.io/).
    And then after waiting for indexing, it will request for searching the nearest vector 10 times.
    You will get the 10 nearest neighbor vectors for each search query.<br>
    Run example codes by executing the below command.

    ```bash
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

        client := vald.NewValdClient(conn)
        ```

        </details>

    1.  Insert and Index

        - Insert and Indexing 400 training datasets to the Vald agent.
            <details><summary>example code</summary><br>

          ```go
          for i := range ids [:insertCount] {
              _, err := client.Insert(ctx, &payload.Insert_Request{
                  Vector: &payload.Object_Vector{
                      Id:     ids[i],
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

    1.  Search

        - Search 10 neighbor vectors for each 20 test datasets and return a list of the neighbor vectors.

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

## Cleanup

In the last, you can remove the deployed Vald Cluster by executing the below command.

```bash
helm uninstall vald
```

## Next Steps

Congratulation! You completely entered the Vald World!

If you want, you can try other tutorials such as:

- [Vald Agent Standalone on k8s](../tutorial/vald-agent-standalone-on-k8s.md)
- [Vald Agent on Docker](../tutorial/vald-agent-standalone-on-docker.md)

For more information, we recommend you to check:

- [Configuration](../user-guides/configuration.md)
- [Operations](../user-guides/operations.md)
