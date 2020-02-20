# Get Started 

This article is written about tutorial for users who wants to use vald.

1. [About](#About)
    1. [Main Features](#Main-Features)
    2. [Requirements](#Requirements)
2. [Starting Vald](#Starting-Vald)
    1. [Deploy via helm](#Deploy-via-helm)
    2. [Deploy from sourcecode](#Deploy-from-sourcecode)
    3. [Indexing](#Indexing)
        1. [Insert Index](#Insert-Index)
        2. [Search Index](#Search-Index)
3. [Advanced](#Advanced)

## About

Vald is distibuted high scalable & high-speed approximate nearest neighor seardch engine.
It provides search result lists for input (e.g. word, sentence, image ...) which is coverted as (also high-dimensional) vector when searching neighor.

### Main Features

```
- Auto Indexing
- Ingress Filltering
- Horizonal Scaling
- Easy to use
- Auto Index Backup


e.g. You should prepare mysql + redis or casandora, due to vald need to use for backup.
```

### Requirements

```
write requirements lists
 - k8s <version>
 - go <version>
 - docker
 - helm
 - etc
```

[helm](https://github.com/helm/helm#install)
docker(if you'd like to deploy on your local envirionment)

## Starting Vald

We provides two methods for deployment.
Before deployment, notice that Vald requires **docker** or **k8s** envirionment.
You can choose [deploy via helm](#Deploy-via-helm) or [deoloy from sourcecode](#Deploy-from-sourcecode).

### Deploy via helm

You have to install helm. Please check [here](https://github.com/helm/helm#install).
You can deploy as follow.

```bash
hele repo add vald https://vald.vdaas.org/charts

export KUBECONFIG="$(kind get kubeconfig-path --name="vald")"

hele install --generate-name vald/vald
```

### Deploy from sourcecode

You can also deploy by cloning sourcecode.
Let you see following.

```bash
git clone https://github.com/vdaas/vald.github

cd vald
# deploy local
make kind start
export KUBECONFIG="$(kind get kubeconfig-path --name="vald")"

make k8s/vald/deploy
```

### Indexing

indexing dataset is most important step for vald.

#### Insert Index

#### Search Index

## Advanced

Vald can be customized as you want, e.g. the number of vector dimension / replica and so on.
If it needs to customized Vald, let you see following.

### values.yaml(git base)
### helm version
    - kubectl delete namespace 
    - helm install -f your.name -generate-name vald/vald 

