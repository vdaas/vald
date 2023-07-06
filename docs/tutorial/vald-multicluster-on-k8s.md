# Multiple Vald Clusters using Vald Mirror Gateway

This article shows how to deploy multiple Vald clusters on your Kubernetes cluster.

## Overview

Vald cluster consists of multiple micro-services.

In [Get Started](https://vald.vdaas.org/docs/tutorial/get-started), you may use 4 kinds of components to deploy the Vald cluster.

In this tutorial, you will deploy multiple Vald clusters with Vald Mirror Gateway, which connects another Vald cluster.

The below image shows the architecture image of this case.

<img src="../../assets/docs/tutorial/vald-multicluster-on-k8s.png">

## Requirements

- Vald: v1.8 ~
- Kubernetes: v1.27 ~
- Go: v1.20 ~
- Helm: v3 ~
- libhdf5 (_only required to get started_)

Helm is used to deploy Vald cluster on your Kubernetes and HDF5 is used to decode the sample data file to run the example code.

If Helm or HDF5 is not installed, please install [Helm](https://helm.sh/docs/intro/install) and [HDF5](https://www.hdfgroup.org/).

<details><summary>Installation command for Helm</summary><br>

```bash
curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash
```

</details>

## Prepare the Kubernetes Cluster

This tutorial requires the Kubernetes cluster.

Vald cluster runs on public Cloud Kubernetes Services such as GKE, EKS.
In the sense of trying to `Get Started`, [k3d](https://k3d.io/) or [kind](https://kind.sigs.k8s.io/) are easy Kubernetes tools to use.

This tutorial uses [Kubernetes Metrics Server](https://github.com/kubernetes-sigs/metrics-server) for running the Vald cluster.

Please make sure these functions are available.

The way to deploy Kubernetes Metrics Service is here:

```bash
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml && \
kubectl wait -n kube-system --for=condition=ready pod -l k8s-app=metrics-server --timeout=600s
````

Please prepare three Namespaces on the Kubernetes cluster.

```bash
kubectl create namespace vald-01 && \
kubectl create namespace vald-02 && \
kubectl create namespace vald-03
```

## Deploy Vald Clusters on each Kubernetes Namespace

This chapter shows how to deploy the multiple Vald clusters using Helm and run it on your Kubernetes cluster.
In this section, you will deploy three Vald clusters consisting of `vald-agent-ngt`, `vald-lb-gateway`, `vald-discoverer`, `vald-manager-index`, and `vald-mirror-gateway` using the basic configuration.

1. Clone the repository

   To use the `deployment` YAML for deployment, let’s clone `[vdaas/vald](https://github.com/vdaas/vald.git)` repository.

   ```bash
   git clone https://github.com/vdaas/vald.git &&　cd vald
   ```

2. Deploy on the `vald-01` Namespace using [dev-vald-01.yaml](https://github.com/vdaas/vald/blob/feature/mirror-gateway-definition/charts/vald/values/multi-vald/dev-vald-01.yaml) and [values.yaml](https://github.com/vdaas/vald/blob/main/example/helm/values.yaml)

   ```bash
   helm install vald-cluster-01 charts/vald \
     -f ./example/helm/values.yaml \
     -f ./charts/vald/values/multi-vald/dev-vald-01.yaml \
     -n vald-01
   ```

3. Deploy on the `vald-02` Namespace using [dev-vald-02.yaml](https://github.com/vdaas/vald/blob/feature/mirror-gateway-definition/charts/vald/values/multi-vald/dev-vald-02.yaml) and [values.yaml](https://github.com/vdaas/vald/blob/main/example/helm/values.yaml)

   ```bash
   helm install vald-cluster-02 charts/vald \
     -f ./example/helm/values.yaml \
     -f ./charts/vald/values/multi-vald/dev-vald-02.yaml \
     -n vald-02
   ```

4. Deploy on the `vald-03` Namespace using [dev-vald-03.yaml](https://github.com/vdaas/vald/blob/feature/mirror-gateway-definition/charts/vald/values/multi-vald/dev-vald-03.yaml) and [values.yaml](https://github.com/vdaas/vald/blob/main/example/helm/values.yaml)

   ```bash
   helm install vald-cluster-03 charts/vald \
     -f ./example/helm/values.yaml \
     -f ./charts/vald/values/multi-vald/dev-vald-03.yaml \
     -n vald-03
   ```

5. Verify

   If success deployment, the Vald cluster’s components should run on each Kubernetes Namespace.

    <details><summary>`vald-01` Namespace</summary><br>
    
     ```bash
     kubectl get pods -n vald-01
     NAME                                   READY   STATUS    RESTARTS   AGE
     vald-agent-ngt-0                       1/1     Running   0          2m41s
     vald-agent-ngt-2                       1/1     Running   0          2m41s
     vald-agent-ngt-3                       1/1     Running   0          2m41s
     vald-agent-ngt-4                       1/1     Running   0          2m41s
     vald-agent-ngt-5                       1/1     Running   0          2m41s
     vald-agent-ngt-1                       1/1     Running   0          2m41s
     vald-discoverer-77967c9697-brbsp       1/1     Running   0          2m41s
     vald-lb-gateway-587879d598-xmws7       1/1     Running   0          2m41s
     vald-lb-gateway-587879d598-dzn9c       1/1     Running   0          2m41s
     vald-manager-index-56d474c848-wkh6b    1/1     Running   0          2m41s
     vald-lb-gateway-587879d598-9wb5q       1/1     Running   0          2m41s
     vald-mirror-gateway-6df75cf7cf-gzcr4   1/1     Running   0          2m26s
     vald-mirror-gateway-6df75cf7cf-vjbqx   1/1     Running   0          2m26s
     vald-mirror-gateway-6df75cf7cf-c2g7t   1/1     Running   0          2m41s
     ```

    </details>

    <details><summary>`vald-02` Namespace</summary><br>

   ```bash
   kubectl get pods -n vald-02
   NAME                                  READY   STATUS    RESTARTS   AGE
   vald-agent-ngt-0                      1/1     Running   0          2m52s
   vald-agent-ngt-1                      1/1     Running   0          2m52s
   vald-agent-ngt-2                      1/1     Running   0          2m52s
   vald-agent-ngt-4                      1/1     Running   0          2m52s
   vald-agent-ngt-5                      1/1     Running   0          2m52s
   vald-agent-ngt-3                      1/1     Running   0          2m52s
   vald-discoverer-8cfcff76-vlmpg        1/1     Running   0          2m52s
   vald-lb-gateway-54896f9f49-wtlcv      1/1     Running   0          2m52s
   vald-lb-gateway-54896f9f49-hbklj      1/1     Running   0          2m52s
   vald-manager-index-676855f8d7-bb4wb   1/1     Running   0          2m52s
   vald-lb-gateway-54896f9f49-kgrdf      1/1     Running   0          2m52s
   vald-mirror-gateway-6598cf957-t2nz4   1/1     Running   0          2m37s
   vald-mirror-gateway-6598cf957-wr448   1/1     Running   0          2m52s
   vald-mirror-gateway-6598cf957-jdd6q   1/1     Running   0          2m37s
   ```

    </details>

    <details><summary>`vald-03` Namespace</summary><br>

   ```bash
   kubectl get pods -n vald-03
   NAME                                   READY   STATUS    RESTARTS   AGE
   vald-agent-ngt-0                       1/1     Running   0          2m46s
   vald-agent-ngt-1                       1/1     Running   0          2m46s
   vald-agent-ngt-2                       1/1     Running   0          2m46s
   vald-agent-ngt-3                       1/1     Running   0          2m46s
   vald-agent-ngt-4                       1/1     Running   0          2m46s
   vald-agent-ngt-5                       1/1     Running   0          2m46s
   vald-discoverer-879867b44-8m59h        1/1     Running   0          2m46s
   vald-lb-gateway-6c8c6b468d-ghlpx       1/1     Running   0          2m46s
   vald-lb-gateway-6c8c6b468d-rt688       1/1     Running   0          2m46s
   vald-lb-gateway-6c8c6b468d-jq7pl       1/1     Running   0          2m46s
   vald-manager-index-5596f89644-xfv4t    1/1     Running   0          2m46s
   vald-mirror-gateway-7b95956f8b-l57jz   1/1     Running   0          2m31s
   vald-mirror-gateway-7b95956f8b-xd9n5   1/1     Running   0          2m46s
   vald-mirror-gateway-7b95956f8b-dnxbb   1/1     Running   0          2m31s
   ```

    </details>

## Deploy ValdMirrorTarget Resource (Custom Resource)

It requires applying the `ValdMirrorTarget` Custom Resource to the one Namespace.

When applied successfully, the destination information is automatically created on other Namespaces when interconnected with each `vald-mirror-gateway`.

This tutorial will deploy the [ValdMirrorTarger](https://github.com/vdaas/vald/tree/main/charts/vald/values/mirror-target.yaml) Custom Resource to the `vald-03` Namespace with the following command.

```bash
kubectl apply -f ./charts/vald/values/multi-vald/mirror-target.yaml -n vald-03
```

The current connection status can be checked with the following command.

<details><summary>Example output</summary><br>

```bash
kubectl get vmt -A -o wide
NAMESPACE   NAME                                 HOST                                            PORT   STATUS      LAST TRANSITION TIME   AGE
vald-03     mirror-target-01                     vald-mirror-gateway.vald-01.svc.cluster.local   8081   Connected   2023-05-22T02:07:51Z   19m
vald-03     mirror-target-02                     vald-mirror-gateway.vald-02.svc.cluster.local   8081   Connected   2023-05-22T02:07:51Z   19m
vald-02     mirror-target-3296010438411762394    vald-mirror-gateway.vald-01.svc.cluster.local   8081   Connected   2023-05-22T02:07:53Z   19m
vald-02     mirror-target-12697587923462644654   vald-mirror-gateway.vald-03.svc.cluster.local   8081   Connected   2023-05-22T02:07:53Z   19m
vald-01     mirror-target-13698925675968803691   vald-mirror-gateway.vald-02.svc.cluster.local   8081   Connected   2023-05-22T02:07:53Z   19m
vald-01     mirror-target-17825710563723389324   vald-mirror-gateway.vald-03.svc.cluster.local   8081   Connected   2023-05-22T02:07:53Z   19m
```

</details>

## Run Example Code

In this chapter, you will execute insert, search, get, and delete vectors to your Vald clusters using the example code.

The [Fashion-MNIST](https://github.com/zalandoresearch/fashion-mnist) is used as a dataset for indexing and search query.

The example code is implemented in Go and uses [vald-client-go](https://github.com/vdaas/vald-client-go), one of the official Vald client libraries, for requesting to Vald cluster.

Vald provides multiple language client libraries such as Go, Java, Node.js, Python, etc.

If you are interested, please refer to [SDKs](https://vald.vdaas.org/docs/user-guides/sdks).

1. Port Forward(option)

   If you do NOT use Kubernetes Ingress, port forwarding is required to make requests from your local environment.

   ```bash
   kubectl port-forward svc/vald-mirror-gateway 8080:8081 -n vald-01 & \
     kubectl port-forward svc/vald-mirror-gateway 8081:8081 -n vald-02 & \
     kubectl port-forward svc/vald-mirror-gateway 8082:8081 -n vald-03
   ```

2. Download dataset

   Move to the working directory.

   ```bash
   cd example/client/mirror
   ```

   Download [Fashion-MNIST](https://github.com/zalandoresearch/fashion-mnist), which is used as a dataset for indexing and search query.

   ```bash
   wget https://ann-benchmarks.com/fashion-mnist-784-euclidean.hdf5
   ```

3. Run Example

   We use [example/client/mirror/main.go](https://github.com/vdaas/vald/blob/feature/mirror-gateway-example/example/client/mirror/main.go) to run the example.

   This example will insert and index 400 vectors into the Vald cluster from the Fashion-MNIST dataset via [gRPC](https://grpc.io/).
   And then, after waiting for indexing, it will request to search the nearest vector 10 times to all Vald clusters. You will get the 10 nearest neighbor vectors for each search query.
   And it will request to get vectors using inserted vector ID.

   Run example codes by executing the below command.

   ```bash
   go run main.go
   ```

   See [GetStarted](https://vald.vdaas.org/docs/tutorial/get-started/) for an explanation of the example code.

## Cleanup

Last, you can remove the deployed Vald cluster by executing the below command.

```bash
helm uninstall vald-cluster-01 -n vald-01 && \
helm uninstall vald-cluster-02 -n vald-02 && \
helm uninstall vald-cluster-03 -n vald-03
```

You can remove all `ValdMirrorTarget` resources (Custom Resource) with the below command.

```bash
kubectl delete --all ValdMirrorTargets -A
```

You can remove all created Namespaces with the below command.

```bash
kubectl delete namespaces vald-01 vald-02 vald-03
```

## Next Steps

Congratulation!
You completely entered the World of multiple Vald clusters!

For more information, we recommend you check the following:

- [Configuration](https://vald.vdaas.org/docs/user-guides/configuration/)
