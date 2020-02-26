# Get Started

For who want to use or be interested in vald, this article shows way to deploy and run vald using scylla for backup and fashion-minist for dataset.
Also, you can see basic info of vald.

1. [About](#About)
    1. [Main Features](#Main-Features)
    2. [Requirements](#Requirements)
2. [Starting Vald on k8s cluster](#Starting-Vald-on-k8s-cluster)
    1. [Deploy](#Deploy)
    2. [Run](#Run)
3. [Advanced](#Advanced)

## About

Vald is distibuted high scalable & high-speed approximate nearest neighor seardch engine.<br>
It provides search result lists for input (e.g. word, sentence, image ...) which is coverted as (also high-dimensional) vector when searching neighor.<br>
Vald is designed as Cloud Native and also use World Fastest ANN Algorithm [NGT](https://github.com/yahoojapan/NGT) for search neighors.

### Main Features

```
- Auto Indexing
- Ingress Filltering
- Horizonal Scaling
    - memory , cpu
- Easy to use
- Auto Index Backup

- Distributed Indexing
- Index Replication
-
```

### Requirements

```bash
write requirements lists
 - k8s <version>
 - go <version>
 - docker
 - helm
 - etc
```

[helm](https://github.com/helm/helm#install)
docker(if you'd like to deploy on your local envirionment)


## Starting Vald on k8s cluster

<details>
    <summary>optional installation</summary>
    <div>kind</div>
</details>

You have to install helm. Please check [here](https://github.com/helm/helm#install).

### Deploy

There shows procudure of deployment foe vald with scylla, which is used for backup database.
If you want to learn about scylla more, please refer to [here](https://www.scylladb.com/).
Please confirm `kubectl cluster-info` before deployment.

1. Prepare scylla database and k8s metrics-server

At first, you have to apply scylla as backup database.

```bash
cd {vald_project_root}

kubectl apply -f k8s/jobs/db/casandora/configmap.yaml
kubectl apply -f k8s/external/scylla
```

Apply k8s metrics-server

```
kubectl apply -f k8s/metrics/metrics-server
```

2. Deployment from helm template command

```
helm install --values example/values-scylla.yaml
```

### Run

This chapter shows procudure of running vald from localhost with fashion-mnist dataset.

1. Port Forward

`kubectl port-forward {vald-gateway-pod-name} 8081:8081`

2. Install dataset

To search neighors of input vector, vald should have dataset.
In this case, you use [fashion-mnist](https://github.com/erikbern/ann-benchmarks) as dataset.

```bash
# move to working directory
cd example/client

# get fashion-mnist
wget https://ann-benchmarks.com/fashion-mnist-784-euclidean.hdf5
```

3. Running example

Insert, Search and Remove will be running as follows.
  - Insert: Insert and Index training dataset to vald agent
  - Search: Search neighor vector of test vector.
  - Remove: Remove indexing from vald agent.

```
go run main.go
```

## Advanced

Vald can be customized as you want, e.g. the number of vector dimension / replica and so on.
If it needs to customized Vald, let you see following.

### values.yaml(git base)
### helm version
    - kubectl delete namespace 
    - helm install -f your.name -generate-name vald/vald 
### Another way to deployment
For `Get Started`, we show the deployment procudure for vald with scylla.
We also provide another ways for development, which are used by helm or `kind` command. (There are needed compressor database )

