# How to use this example
After launching the k8s cluster, do the following steps at the root directory.

1. Deploy Vald cluster with filter-gateway

      ```bash
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
              distance_type: l2
      ...
      
      # deploy vald cluster
      helm repo add vald https://vald.vdaas.org/charts
      helm install vald vald/vald --values example/helm/values.yaml
      ```

2. Build and publish example ingress filter and egress filter docker image

```bash
# login to docker if needed, and setup your DockerHub ID
docker login
export DOCKERHUB_ID=<DOCKERHUB_ID>

# build and publish ingress filter image
docker build \
    -f example/manifest/filter/ingress/Dockerfile \
    -t $DOCKERHUB_ID/vald-ingress-filter:latest . \
    --build-arg GO_VERSION=$(make version/go)

docker push ${DOCKERHUB_ID}/vald-ingress-filter:latest

# build and publish egress filter image
docker build \
    -f example/manifest/filter/egress/Dockerfile \
    -t $DOCKERHUB_ID/vald-egress-filter:latest . \
    --build-arg GO_VERSION=$(make version/go)

docker push ${DOCKERHUB_ID}/vald-egress-filter:latest
```

3. Deploy ingress filter server and egress filter server
```bash
# deploy ingress filter
sed -e "s/DOCKERHUB_ID/${DOCKERHUB_ID}/g" example/manifest/filter/egress/deployment.yaml | kubectl apply -f - \
&& kubectl apply -f example/manifest/filter/egress/service.yaml

# deploy egress filter
sed -e "s/DOCKERHUB_ID/${DOCKERHUB_ID}/g" example/manifest/filter/ingress/deployment.yaml | kubectl apply -f - \
&& kubectl apply -f example/manifest/filter/ingress/service.yaml
```

4. Run test
```bash
# if you don't use the Kubernetes ingress, set the port forward
kubectl port-forward deployment/vald-filter-gateway 8081:8081

# Please change the argument according to your environment.
go run ./example/client/gateway/filter/main.go -addr "localhost:8081" -ingresshost "vald-ingress-filter.default.svc.cluster.local" -ingressport 8082 -egresshost "vald-egress-filter.default.svc.cluster.local" -egressport 8083
```

5. Cleanup
```bash
helm uninstall vald
kubectl delete -f ./example/manifest/filter/egress/deployment.yaml -f ./example/manifest/filter/egress/service.yaml -f ./example/manifest/filter/ingress/deployment.yaml -f ./example/manifest/filter/ingress/service.yaml
```