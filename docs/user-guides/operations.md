# Operations

This page introduces best practices to operate a Vald cluster.

Table of Contents
---

- [Deployment](#deployment)
    - [Kubernetes cluster](#kubernetes-cluster)
    - [Multi tenant](#multi-tenant)
- [Monitoring](#monitoring)
    - [Observability features of Vald](#observability-features-of-vald)
    - [Enabling observability feature](#enabling-observability-feature)
    - [Use Prometheus and Grafana](#use-prometheus-and-grafana)
    - [Stackdriver (Cloud) Tracing, monitoring and Profiler](#stackdriver)
- [Upgrading](#upgrading)
    - [In case of manual deploy](#in-case-of-manual-deploy)
    - [In case of using Helm](#in-case-of-using-helm)
    - [In case of using Vald-Helm-Operator](#in-case-of-using-vald-helm-operator)
- [References](#references)


## Deployment

### Kubernetes cluster


### Multi tenant



## Monitoring

### Observability features of Vald


### Enabling observability feature

If observability features are enabled, the metrics will be collected periodically. the duration can be set on `observability.collector.duration`.

### Use Prometheus and Grafana


### <a name="stackdriver"></a> Stackdriver (Cloud) Tracing, Monitoring and Profiler


## Upgrading

### In case of manual deploy


### In case of using Helm


### In case of using Vald-Helm-Operator






### References

[vald-helm-chart]: https://github.com/vdaas/vald/tree/master/charts/vald
[vald-helm-operator-chart]: https://github.com/vdaas/vald/tree/master/charts/vald-helm-operator

[prometheus-io]: https://prometheus.io/
