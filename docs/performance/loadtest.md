# Load Testing

This document describe how to perform the load testing in Vald cluster or Vald agent using Vald Load Testing Tools.

## Overview

We develop our own Load Testing Tools to perform the load testing in Vald cluster or Vald agent. You can perform the load test on Vald to test the insert and search performance in Vald.

## Prerequisite

- Vald Cluster or Vald Agent

    Please refer to [here](https://vald.vdaas.org/docs/tutorial/get-started/) for the installation guide of Vald Cluster or Vald Agent.

- Docker

## Install load test tools

We provide a docker image for the load testing. Please download the image using the below command.

```bash
docker pull vdaas/vald-loadtest
```

## Configure load test tools

Please refer to the [Sample configuration file](https://github.com/vdaas/vald/blob/master/cmd/tools/cli/loadtest/sample.yaml) to configure the Vald Load Testing Tools.

Here is the important configurations and following the explaination.

| Name        | Description                                                                                                                                                      | Example                                             |
|-------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------|-----------------------------------------------------|
| service     | By setting the value to `gateway`, you can test the Vald cluster though the Vald gateway. By setting the value to `agent`, you can test the Vald agent directly. | `gateway` or `agent`                                |
| operation   | To perform the specific load test action to the Vald cluster.                                                                                                    | `insert`, `streaminsert`, `search` or `steamsearch` |
| dataset     | The dataset used in `insert` and `streaminsert` operation.                                                                                                       |                                                     |
| concurrency | The number of concurrent execution of the load test.                                                                                                             |                                                     |
| batch_size  | The batch size of the dataset.                                                                                                                                   |                                                     |
| addr        | The cluster you want to test                                                                                                                                     |                                                     |

## Execute load test

Download the sample configuration file.

```bash
wget https://raw.githubusercontent.com/vdaas/vald/master/cmd/tools/cli/loadtest/sample.yaml
mv sample.yaml config.yaml
```

Please refer to the above section to configure the load testing tool.

```bash
vi config.yaml
```

Execute the load test in docker.

```bash
docker run -it --rm -v $(pwd):/etc/server vald-load-test
```

## Result
    explain how to read and interpret the result,
