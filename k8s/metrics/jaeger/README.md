# Jaeger

This directory contains manifests to deploy jaeger-operator and jaeger.

The below is an example of integration with Vald.

```
# values.yaml
defaults:
  observability:
    enabled: true
    trace:
      enabled: true
      sampling_rate: 0.5
    jaeger:
      enabled: true
      agent_endpoint: jaeger-agent.default.svc.cluster.local:6831
      username: ""
      password: ""
      service_name: vald
      buffer_max_count: 10
```

And we can open webUI using port forwarding.

```
kubectl port-forward service/jaeger-query 16686:16686
```
