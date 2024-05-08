# Load Testing

This document describes how to perform the load testing in the Vald cluster or Vald Agent using Vald Load Testing Tools.

## Overview

We develop our Load Testing Tools to perform the load testing in the Vald cluster or the Vald Agent.
You can perform the load test on the Vald cluster to test the insert and search performance in Vald.

## Prerequisite

- Vald Cluster or Vald Agent

  Please refer to [here](https://vald.vdaas.org/docs/tutorial/get-started/) for the installation guide of Vald Cluster or Vald Agent.

- Docker: v19.03 ~

## Install load test tools

We provide a docker image for the load testing.
Please download the image using the below command.

```bash
docker pull vdaas/vald-loadtest
```

## Configure load test tools

Please refer to the [Sample configuration file](https://github.com/vdaas/vald/blob/main/cmd/tools/cli/loadtest/sample.yaml) to configure the Vald Load Testing Tools.

Here are the important configurations and following the explanation.

| Name        | Description                                                                                                                                                             | Example                                     |
| :---------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------- | :------------------------------------------ |
| service     | By setting the value to the gateway, you can test the Vald cluster through the Vald gateway. <br />By setting the value to agent, you can test the Vald agent directly. | gateway or agent                            |
| operation   | To perform the specific load test action to the Vald cluster.                                                                                                           | insert, streaminsert, search or steamsearch |
| dataset     | The dataset is used in insert and stream insert operation.                                                                                                              |                                             |
| concurrency | The number of concurrent execution of the load test.                                                                                                                    |                                             |
| batch_size  | The batch size of the dataset.                                                                                                                                          |                                             |
| addr        | The cluster you want to test                                                                                                                                            |                                             |

## Execute load test

Download the sample configuration file.

```bash
wget https://raw.githubusercontent.com/vdaas/vald/main/cmd/tools/cli/loadtest/sample.yaml
mv sample.yaml config.yaml
```

Please refer to the [configure load test tools](#configure-load-test-tools) section to configure the load testing tool.

```bash
vi config.yaml
```

Execute the load test in docker.

```bash
docker run -it --rm -v $(pwd):/etc/server vald-load-test
```

## Result

After the Vald Load Testing tools finished, the following output will be displayed.

```bash
2020-07-13 08:10:45	[INFO]:	maxprocs: Leaving GOMAXPROCS=8: CPU quota undefined
2020-07-13 08:10:45	[INFO]:	service load test v0.0.0 starting...
2020-07-13 08:10:45	[INFO]:	executing daemon pre-start function
2020-07-13 08:10:45	[DEBG]:	start loading: random-786-100000
2020-07-13 17:10:49	[INFO]:	executing daemon start function
2020-07-13 17:10:49	[INFO]:	start load test(Gateway, Insert)
2020-07-13 17:10:52	[INFO]:	progress 177 requests, 5899.977806[vps], error: 0
2020-07-13 17:10:55	[INFO]:	progress 370 requests, 6166.631312[vps], error: 0
2020-07-13 17:10:58	[INFO]:	progress 558 requests, 6199.932718[vps], error: 0
2020-07-13 17:11:01	[INFO]:	progress 727 requests, 6058.303243[vps], error: 0
2020-07-13 17:11:04	[INFO]:	progress 910 requests, 6066.658704[vps], error: 0
2020-07-13 17:11:05	[INFO]:	result:Gateway	32	100	6060.743248
2020-07-13 17:11:05	[WARN]:	terminated signal received daemon will stopping soon...
2020-07-13 17:11:05	[INFO]:	executing daemon pre-stop function
2020-07-13 17:11:05	[INFO]:	executing daemon stop function
2020-07-13 17:11:05	[INFO]:	executing daemon post-stop function
2020-07-13 17:11:05	[ERR]:	error occurred in runner.Wait at load test: context canceled
2020-07-13 17:11:05	[WARN]:	daemon stopped
```

These lines show the result of the load test.

```bash
2020-07-13 08:10:45 [DEBG]: start loading: random-786-100000
...
2020-07-13 17:10:49 [INFO]: start load test(Gateway, Insert)
...
2020-07-13 17:11:05	[INFO]:	result:Gateway	32	100	6060.743248
```

This means that the `Insert Gateway` mode is used, which means the dataset `random-786-100000` will insert into the Vald Cluster.

The line `result:Gateway 32 100 6060.743248` means that the Gateway mode is used, with `32` concurrent execution and `100` batch size, with the VPS (Vector Per Sec.) of `6060.743248`.

It means that is performed the insertion of 6060 vectors into the Vald Cluster per second from the Vald Load Testing tools.

The result includes all the network latency and filtering latency.
