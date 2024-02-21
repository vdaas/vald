# Continuous Benchmark Tool

## What is the Continuous Benchmark Tool?

Continuous Benchmark Tool allows you to get benchmark of Vald cluster in 24/365.

Assumed use case is:

- Verification with workload close to production environment
- Verification before service installation when Vald version up

## Architecture

Continuous Benchmark Tool has following 2 components:

- Benchmark Operator: Manages benchmark jobs
- Benchmark Job: Executes CRUDs request to the target Vald cluster

## Benchmark component and its feature

### Benchmark Operator

- Manages benchmark jobs according to applied manifest.
- Apply method:
  - Scenario method: one manifest with multiple benchmark jobs
  - Job method: one manifest with one benchmark job

### Benchmark Job

- Executes CRUD request to the target Vald cluster based on defined config.
- Execute steps are:
  1. Load dataset (valid only for HDF5 format )
  1. Execute request with load dataset

## Benchmark CRD

Benchmark workload can be set by applying the Kubernetes Custom Resources(CRDs), `ValdBenchmarkScenarioResource` or `ValdBenchmarkJobResource`.
Benchmark Operator manages benchmark job according to the applied manifest.

### ValdBenchmarkJob

[`ValdBenchmarkJob`](https://github.com/vdaas/vald/blob/main/charts/vald-benchmark-operator/crds/valdbenchmarkjob.yaml) is used for executing single benchmark job.

And, Benchmark Operator also applies it to the Kubernetes cluster based on `ValdBenchmarkScenarioResource`.

**main properties**

| Name                       | mandatory | Description                                                                                                           | type                                                                     | sample                                                                                       |
| :------------------------- | :-------- | :-------------------------------------------------------------------------------------------------------------------- | :----------------------------------------------------------------------- | :------------------------------------------------------------------------------------------- |
| target                     | \*        | target Vald cluster                                                                                                   | object                                                                   | ref: [target](#target-prop)                                                                  |
| dataset                    | \*        | dataset information                                                                                                   | object                                                                   | ref: [dataset](#dataset-prop)                                                                |
| job_type                   | \*        | execute job type                                                                                                      | string enum: [insert, update, upsert, remove, search, getobject, exists] | search                                                                                       |
| repetition                 |           | the number of job repetitions<BR>default: `1`                                                                         | integer                                                                  | 1                                                                                            |
| replica                    |           | the number of job concurrent job executions<BR>default: `1`                                                           | integer                                                                  | 2                                                                                            |
| rps                        |           | designed request per sec to the target cluster<BR>default: `1000`                                                     | integer                                                                  | 1000                                                                                         |
| concurrency_limit          |           | goroutine count limit for rps adjustment<BR>default: `200`                                                            | integer                                                                  | 20                                                                                           |
| ttl_seconds_after_finished |           | time until deletion of Pod after job end<BR>default: `600`                                                            | integer                                                                  | 120                                                                                          |
| insert_config              |           | request config for insert job                                                                                         | object                                                                   | ref: [config](#insert-cfg-props)                                                             |
| update_config              |           | request config for update job                                                                                         | object                                                                   | ref: [config](#update-cfg-props)                                                             |
| upsert_config              |           | request config for upsert job                                                                                         | object                                                                   | ref: [config](#upsert-cfg-props)                                                             |
| search_config              |           | request config for search job                                                                                         | object                                                                   | ref: [config](#search-cfg-props)                                                             |
| remove_config              |           | request config for remove job                                                                                         | object                                                                   | ref: [config](#remove-cfg-props)                                                             |
| object_config              |           | request config for object job                                                                                         | object                                                                   | ref: [config](#object-cfg-props)                                                             |
| client_config              |           | gRPC client config for running benchmark job<BR>Tune if can not getting the expected performance with default config. | object                                                                   | ref: [defaults.grpc](https://github.com/vdaas/vald/blob/main/charts/vald/README.md)          |
| server_config              |           | server config for benchmark job pod<BR>Tune if can not getting the expected performance with default config.          | object                                                                   | ref: [defaults.server_config](https://github.com/vdaas/vald/blob/main/charts/vald/README.md) |

<a id="target-prop" />

**target**

- target Vald cluster information
- type: object

| property | mandatory | description           | type    | sample    |
| :------- | :-------- | :-------------------- | :------ | :-------- |
| host     | \*        | target cluster's host | string  | localhost |
| port     | \*        | target cluster's port | integer | 8081      |

<a id="dataset-prop" />

**dataset**

- dataset which is used for executing job operation
- type: object

| property    | mandatory | description                                                                                          | type                                   | sample        |
| :---------- | :-------- | :--------------------------------------------------------------------------------------------------- | :------------------------------------- | :------------ |
| name        | \*        | dataset name                                                                                         | string enum: [fashion-mnist, original] | fashion-mnist |
| group       | \*        | group name                                                                                           | string enum: [train, test, neighbors]  | train         |
| indexes     | \*        | amount of index size                                                                                 | integer                                | 1000000       |
| range       | \*        | range of indexes to be used (if there are many indexes, the range will be corrected on the job side) | object                                 | -             |
| range.start | \*        | start of range                                                                                       | integer                                | 1             |
| range.end   | \*        | end of range                                                                                         | integer                                | 1000000       |
| url         |           | the dataset url. It should be set when set `name` as `original`                                      | string                                 |               |

<a id="insert-cfg-props" />

**insert_config**

- rpc config for insert request
- type: object

| property                | mandatory | description                                                                                                  | type   | sample     |
| :---------------------- | :-------- | :----------------------------------------------------------------------------------------------------------- | :----- | :--------- |
| skip_strict_exist_check |           | Check whether the same vector is already inserted or not.<br>The ID should be unique if the value is `true`. | bool   | false      |
| timestamp               |           | The timestamp of the vector inserted.<br>If it is N/A, the current time will be used.                        | string | 1707272658 |

<a id="update-cfg-props" />

**update_config**

- rpc config for update request
- type: object

| property                | mandatory | description                                                                                                  | type   | sample     |
| :---------------------- | :-------- | :----------------------------------------------------------------------------------------------------------- | :----- | :--------- |
| skip_strict_exist_check |           | Check whether the same vector is already inserted or not.<br>The ID should be unique if the value is `true`. | bool   | false      |
| timestamp               |           | The timestamp of the vector inserted.<br>If it is N/A, the current time will be used.                        | string | 1707272658 |
| disable_balanced_update |           | A flag to disable balanced update (split remove -&gt; insert operation) during update operation.             | bool   | false      |

<a id="upsert-cfg-props" />

**upsert_config**

- rpc config for upsert request
- type: object

| property                | mandatory | description                                                                                                  | type   | sample     |
| :---------------------- | :-------- | :----------------------------------------------------------------------------------------------------------- | :----- | :--------- |
| skip_strict_exist_check |           | Check whether the same vector is already inserted or not.<br>The ID should be unique if the value is `true`. | bool   | false      |
| timestamp               |           | The timestamp of the vector inserted.<br>If it is N/A, the current time will be used.                        | string | 1707272658 |
| disable_balanced_update |           | A flag to disable balanced update (split remove -&gt; insert operation) during update operation.             | bool   | false      |

<a id="search-cfg-props" />

**upsert_config**

- rpc config for search request
- type: object

| property              | mandatory | description                                                                                                                                     | type                                                                                     | sample |
| :-------------------- | :-------- | :---------------------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- | :----- |
| radius                |           | The search radius.<BR>default: `-1`                                                                                                             | number                                                                                   | -1     |
| epsilon               |           | The search coefficient.<BR>default: `0.05`                                                                                                      | number                                                                                   | 0.05   |
| num                   | \*        | The maximum number of results to be returned.                                                                                                   | integer                                                                                  | 10     |
| min_num               |           | The minimum number of results to be returned.                                                                                                   | integer                                                                                  | 5      |
| timeout               |           | Search timeout in nanoseconds<BR>default: `10s`                                                                                                 | string                                                                                   | 3s     |
| enable_linear_search  |           | A flag to enable linear search operation for estimating search recall.<BR>If it is `true`, search operation with linear operation will execute. | bool                                                                                     | false  |
| aggregation_algorithm |           | The search aggregation algorithm option.<BR>default: `Unknown`                                                                                  | string enum: ["Unknown", "ConcurrentQueue", "SortSlice", "SortPoolSlice", "PairingHeap"] |        |

<a id="remove-cfg-props" />

**remove_config**

- rpc config for remove request
- type: object

| property                | mandatory | description                                                                                                  | type   | sample     |
| :---------------------- | :-------- | :----------------------------------------------------------------------------------------------------------- | :----- | :--------- |
| skip_strict_exist_check |           | Check whether the same vector is already inserted or not.<br>The ID should be unique if the value is `true`. | bool   | false      |
| timestamp               |           | The timestamp of the vector inserted.<br>If it is N/A, the current time will be used.                        | string | 1707272658 |

<a id="object-cfg-props" />

**object_config**

- rpc config for get object request
- type: object

| property              | mandatory | description                                                 | type     | sample |
| :-------------------- | :-------- | :---------------------------------------------------------- | :------- | :----- |
| filter_config.targets |           | filter target host and port for bypassing filter component. | []object |        |

### ValdBenchmarkScenario

[`ValdBenchmarkScenario`](https://github.com/vdaas/vald/blob/main/charts/vald-benchmark-operator/crds/valdbenchmarkscenario.yaml) is used for executing single or multiple benchmark job.

Benchmark Operator decomposes manifest and creates benchmark resources one by one.
The `target` and `dataset` property are the global config for scenario, they can be overwritten when each job has own config.

**main properties**

| property | mandatory | description                                                                            | type   | sample                                  |
| :------- | :-------- | :------------------------------------------------------------------------------------- | :----- | :-------------------------------------- |
| target   | \*        | target Vald cluster information<BR>It will be overwritten when each job has own config | object | ref: [target](#target-prop)             |
| dataset  | \*        | dataset information<BR>It will be overwritten when each job has own config             | object | ref: [dataset](#dataset-prop)           |
| jobs     | \*        | benchmark job config<BR>The jobs written above will be executed in order.              | object | ref: [benchmark job](#valdbenchmarkjob) |

## Deploy Benchmark Operator

Continuous benchmark operator can be applied with `Helm` same as Vald cluster.

It requires `ValdBenchmarkOperatorRelease` for deploying `vald-benchmark-operator`.

It is not must to apply, so please edit and apply as necessary.

<details><summary>Sample ValdBenchmarkOperatorRelease YAML</summary>

```yaml
# @schema {"name": "name", "type": "string"}
# name -- name of the deployment
name: vald-benchmark-operator
# @schema {"name": "time_zone", "type": "string"}
# time_zone -- time_zone
time_zone: ""
# @schema {"name": "image", "type": "object"}
image:
  # @schema {"name": "image.repository", "type": "string"}
  # image.repository -- image repository
  repository: vdaas/vald-benchmark-operator
  # @schema {"name": "image.tag", "type": "string"}
  # image.tag -- image tag
  tag: v1.7.5
  # @schema {"name": "image.pullPolicy", "type": "string", "enum": ["Always", "Never", "IfNotPresent"]}
  # image.pullPolicy -- image pull policy
  pullPolicy: Always
# @schema {"name": "job_image", "type": "object"}
job_image:
  # @schema {"name": "job_image.repository", "type": "string"}
  # image.repository -- job image repository
  repository: vdaas/vald-benchmark-job
  # @schema {"name": "job_image.tag", "type": "string"}
  # image.tag -- image tag for job docker image
  tag: v1.7.5
  # @schema {"name": "job_image.pullPolicy", "type": "string", "enum": ["Always", "Never", "IfNotPresent"]}
  # image.pullPolicy -- image pull policy
  pullPolicy: Always
# @schema {"name": "resources", "type": "object"}
# resources -- kubernetes resources of pod
resources:
  # @schema {"name": "resources.limits", "type": "object"}
  limits:
    cpu: 300m
    memory: 300Mi
  # @schema {"name": "resources.requests", "type": "object"}
  requests:
    cpu: 200m
    memory: 200Mi
# @schema {"name": "logging", "type": "object"}
logging:
  # @schema {"name": "logging.logger", "type": "string", "enum": ["glg", "zap"]}
  # logging.logger -- logger name.
  logger: glg
  # @schema {"name": "logging.level", "type": "string", "enum": ["debug", "info", "warn", "error", "fatal"]}
  # logging.level -- logging level.
  level: debug
  # @schema {"name": "logging.format", "type": "string", "enum": ["raw", "json"]}
  # logging.format -- logging format.
  format: raw
```

</details>

For more details of the configuration of `vald-benchmark-operator-release`, please refer to [here](https://github.com/vdaas/vald/blob/main/charts/vald-benchmark-operator/values.yaml)

1. Add Vald repo into the helm repo

   ```bash
   helm repo add vald https://vdaas.vald.org
   ```

1. Deploy `vald-benchmark-operator-release`

   ```bash
   helm install vald-benchmark-operator-release vald/vald-benchmark-operator
   ```

1. Apply `vbor.yaml` (optional)

   ```bash
   kubectl apply -f vbor.yaml
   ```

## Running Continuous Benchmarks

After deploy the benchmark operator, you can execute continuous benchmark by applying `ValdBenchmarkScenarioRelease` or `ValdBenchmarkJobRelease`.

Please configure designed benchmark and apply by `kubectl` command.

The sample manifests are [here](https://github.com/vdaas/vald/tree/main/example/helm/benchmark).
