# Load Testing

This document explain why do we need to perform load test in Vald, and how to perform load test in Vald.

## Overview

Vald develop our own Load Testing Tools to perform the load test in Vald cluster.

## Prerequisite

- Vald Cluster or Vald Agent

    You can setup and deploy Vald cluster using helm.
    Please refer to [here](https://vald.vdaas.org/docs/tutorial/get-started/) for Vald Cluster or agent installation.

- Docker

## Install load test tools

We suggest to run the load test tools in Docker.

```
docker pull vdaas/vald-loadtest
```

## Configure load test tools

Please refer to the [Sample configuration file](https://github.com/vdaas/vald/blob/master/cmd/tools/cli/loadtest/sample.yaml) to configure the Vald Load Testing Tools. 

Here is the description of the important configurations.

| Name        | Description                                                                                                                                                      | Example                                             |
|-------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------|-----------------------------------------------------|
| service     | By setting the value to `gateway`, you can test the Vald cluster though the Vald gateway. By setting the value to `agent`, you can test the Vald agent directly. | `gateway` or `agent`                                |
| operation   | To perform the specific load test action to the Vald cluster.                                                                                                    | `insert`, `streaminsert`, `search` or `steamsearch` |
| dataset     | The dataset used in `insert` and `streaminsert` operation.                                                                                                       |                                                     |
| concurrency | The number of concurrent execution of the load test.                                                                                                             |                                                     |
| batch_size  | The batch size of the dataset.                                                                                                                                   |                                                     |
| addr        | The cluster you want to test                                                                                                                                     |                                                     |

## Execute load test

## Result
    explain how to read and interpret the result,
