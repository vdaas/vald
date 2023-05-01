# How to use this example
After launching the k8s cluster, do the following steps at the root directory.

1. Deploy Vald cluster with filter-gateway
```
vim example/helm/values.yaml
---
...
gateway:
...
    filter:
        enabled: true
...
agent:
    ngt:
        dimension: 784
...

// deploy vald cluster
helm repo add vald https://vald.vdaas.org/charts
helm install vald vald/vald --values example/helm/values.yaml

// if you don't use the Kubernetes ingress, set the port forward
kubectl port-forward deployment/vald-filter-gateway 8081:8081
```

2. Deploy an ingress filter server and an egress filter server
```
// prepare for pushing images to your dockerhub
docker login
export YOUR_DOCKERHUB_ID=<YOUR_DOCKERHUB_ID>

// deploy ingress filter
make docker/build/example/client/gateway/filter/ingress-filter/server
docker push ${YOUR_DOCKERHUB_ID}/vald-ingress-filter:latest
make k8s/example/client/gateway/filter/ingress-filter/server/deploy

// deploy egress filter
make docker/build/example/client/gateway/filter/egress-filter/server
docker push ${YOUR_DOCKERHUB_ID}/vald-egress-filter:latest
make k8s/example/client/gateway/filter/egress-filter/server/deploy
```

3. Run test
```
// Please change the argument according to your environment.
go run ./example/client/gateway/filter/main.go -addr "localhost:8081" -ingresshost "vald-ingress-filter.default.svc.cluster.local" -ingressport 8082 -egresshost "vald-egress-filter.default.svc.cluster.local" -egressport 8083
```