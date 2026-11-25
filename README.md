# Vald E2E V2 Benchmarks

This orphan branch stores the latest daily E2E V2 benchmark metrics and pprof
baseline profiles for [`vdaas/vald`](https://github.com/vdaas/vald).

Each daily scheduled run (`0 21 * * *` UTC / `06:00` JST) overwrites the files
under `e2e/v2/latest/` and `e2e/v2/nightly/` with the newest verified results.

## Directory Layout

```text
README.md
e2e/v2/
  latest/
    <scenario>-<deployment>-<cluster>-<environment>.json
    profiles/
      <service>-<profile_type>.pb
  nightly/
    <scenario>-<deployment>-<cluster>-<environment>.json
    profiles/
      <service>-<profile_type>.pb
```

## Metrics JSON Schema & Naming

- Filename pattern: `<scenario>-<deployment>-<cluster>-<environment>.json`
- Example: `stream_crud-helm-chart-k3d-null.json`
- Top-level `metrics.Report` (`schema_version: "v1"`):
  - `schema_version`, `scenario`, `timestamp`
  - `summary`: aggregated `*metrics.GlobalSnapshot`
  - `strategies[]`: per-strategy `name`, `summary`, `operations[]`
    - `operations[]`: per-operation `name`, `summary`, `executions[]`
      - `executions[]`: per-RPC `index`, `name`, `type`, `mode`, `summary`
- `*metrics.GlobalSnapshot` fields (`summary`):
  - `start_time`, `last_updated` (elapsed duration and throughput are derived from `last_updated - start_time` and `total`)
  - `total`, `errors`, `codes`, `error_details`
  - `latencies` (`HistogramSnapshot`, nanoseconds): `total`, `min`, `mean`, `max`, `std_dev`, `sum`, `sum_sq`, `m2`, `bounds`, `counts`
  - `lat_percentiles` (`tdigest`, nanoseconds): `centroids`, `quantiles`, `compression`, `compression_trigger_factor`, `count`
  - `queue_waits`, `qw_percentiles`
  - `custom_counters` (e.g., `recall_sum_ppm`, `recall_samples`)
  - `exemplars`, `exemplar_details`, `invariants_ok`, `sketch_kind`, `bounds_hash`

## Profile Baselines

Profiles are collected via `profefe` during the `profile` environment run and
stored at `e2e/v2/<tag>/profiles/<service>-<profile_type>.pb`:

- `<tag>`: `latest`, `nightly`
- `<service>`: `vald-agent-ngt`, `vald-lb-gateway`, `vald-discoverer`, `vald-manager-index`
- `<profile_type>`: `heap`, `cpu`, `goroutine`, `threadcreate`
