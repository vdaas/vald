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

## Metrics JSON Naming

- Filename pattern: `<scenario>-<deployment>-<cluster>-<environment>.json`
- Example: `stream_crud-helm-chart-k3d-null.json`
- Recorded metrics per scenario, strategy, operation, and execution (`schema_version: "v1"`):
  - `total`, `errors`
  - `duration_ns`, `duration_sec`, `qps`
  - `latency`: `min_ns`, `max_ns`, `mean_ns`, `p99_ns`, `p95_ns`, `p90_ns`, `p50_ns`, `stddev_ns`

## Profile Baselines

Profiles are collected via `profefe` during the `profile` environment run and
stored at `e2e/v2/<tag>/profiles/<service>-<profile_type>.pb`:

- `<tag>`: `latest`, `nightly`
- `<service>`: `vald-agent-ngt`, `vald-lb-gateway`, `vald-discoverer`, `vald-manager-index`
- `<profile_type>`: `heap`, `cpu`, `goroutine`, `threadcreate`
