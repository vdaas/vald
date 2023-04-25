# How to use this example
After launching the k8s cluster, do the following steps at the root directory.

1. Deploy ValdRelease with filter-gateway settings
```
// TODO
```

2. Deploy an ingress filter server and an egress filter server
```
$ make docker/build/example/client/gateway/filter/ingress-filter/server
$ make k8s/example/client/gateway/filter/ingress-filter/server/deploy

$ make docker/build/example/client/gateway/filter/egress-filter/server
$ make k8s/example/client/gateway/filter/egress-filter/server/deploy
```

3. Run test
```
// Please change the argument according to your environment.
$ go run ./example/client/gateway/filter/main.go -addr "localhost:8081" -ingresshost "vald-ingress-filter.default.svc.cluster.local" -ingressport 8082 -egresshost "vald-egress-filter.default.svc.cluster.local" -egressport 8083
```