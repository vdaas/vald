# OpenTelemetry Operator

This directory contains manifests to deploy OpenTelemetry Operator and Collector.

# Getting started

This section describes how to set up a monitoring system using OpenTelemetry, Jaeger and Prometheus and Grafana.

The following is an example structure.

```
                                        ===> Jaeger
                                       ||
Vald ===> OpenTelemetry Collector  ====
                                       ||
                                        ===> Prometheus <=== Grafana
```

1. To deploy the operator in a Kubernetes Cluster, we first need to deploy [cert-manager](https://cert-manager.io/docs/installation/).

   ```sh
   make k8s/external/cert-manager/deploy
   ```

2. Deploy Jaeger and Prometheus to export traces and metrics from OpenTelemetry Collector.

   - Deploy Jaeger

   ```sh
   make k8s/metrics/jaeger/deploy
   ```

   - Deploy Prometheus

   ```sh
   make k8s/metrics/prometheus/operator/deploy
   ```

3. Deploy Grafana and dashboard to visualize the metrics.

   ```sh
   make k8s/metrics/grafana/deploy
   ```

4. Deploy OpenTelemetry Operator and Collector.

   - Deploy OpenTelemetry Operator

   ```sh
   make k8s/otel/operator/install
   ```

   - Deploy OpenTelemetry Collector

   ```sh
   make k8s/otel/collector/install
   ```
