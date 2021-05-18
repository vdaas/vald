# Vald Agent Standalone on Docker

Vald is designed and implemented based on Cloud-Native architecture.
However, there may be cases that want to use only Vald Agent without Kubernetes.

This article will show you how to deploy and run the Vald Agent on Docker.
Fashion-mnist is used as an example dataset, same as [Get Started](../tutorial/get-started.md).

## Requirements

- Docker: v19.0 ~
- Go: v1.15 ~
- libhdf5 (_only required for this tutorial._)

HDF5 is required for this tutorial. If hdf5 is not installed, please install [HDF5](https://www.hdfgroup.org/).
<details><summary>[Optional] Install HDF5</summary><br>

```bash
# yum
yum install -y hdf5-devel

# apt
apt-get install libhdf5-serial-dev

# homebrew
brew install hdf5
```
</details>

## Deploy

This chapter will show you how to deploy Vald Agent on docker.<br>
This chapter will use NGT for the core engine of Vald Agent.

1. Clone the vdaas/vald repository

    ```bash
    git clone https://github.com/vdaas/vald.git
    ```

1. Create directory for setting deploy vald-agent-ngt

    ```bash
    cd vald
    mkdir -p tutorial && cd tutorial
    ```

1. Create `config.yaml`

    The configuration of Vald agent for docker is set using `config.yaml`<br>
    You can also check [the sample configuration](https://github.com/vdaas/vald/blob/master/cmd/agent/core/ngt/sample.yaml).

    ```bash
    cat << EOF > config.yaml
    ---
    version: v0.0.0
    time_zone: JST
    logging:
      logger: glg
      level: debug
      format: raw
    server_config:
      servers:
        - name: agent-grpc
          host: 0.0.0.0
          port: 8081
          mode: GRPC
          probe_wait_time: "3s"
          http:
            shutdown_duration: "5s"
            handler_timeout: ""
            idle_timeout: ""
            read_header_timeout: ""
            read_timeout: ""
            write_timeout: ""
      startup_strategy:
        - agent-grpc
      shutdown_strategy:
        - agent-grpc
      full_shutdown_duration: 600s
      tls:
        enabled: false
        # cert: /path/to/cert
        # key: /path/to/key
        # ca: /path/to/ca
    ngt:
      # path to index data
      index_path: "/etc/server/backup"
      # vector dimension
      dimension: 784
      # bulk insert chunk size
      bulk_insert_chunk_size: 10
      # distance_type, which should be "l1", "l2" "angle", "hamming", "cosine", "normalizedangle", "normalizedcosine" or "jaccard"
      distance_type: l2
      # object_type, which should be "float" or "uint8"
      object_type: float
      # creation edge size
      creation_edge_size: 20
      # search edge size
      search_edge_size: 10
      # The limit duration of automatic indexing
      # auto_index_duration_limit should be 30m-6h for production use. Below setting is a just example
      auto_index_duration_limit: 1m
      # Check duration of automatic indexing.
      # auto_index_check_duration be 10m-1h for production use. Below setting is a just example
      auto_index_check_duration: 10s
      # The number of cache to trigger automatic indexing
      auto_index_length: 100
      # The limit duration of auto saving indexing
      # auto_save_index_duration should be 30m-60m for production use. The below setting is a just example.
      auto_save_index_duration: 90s
      # The maximum limit duration for an initial delay
      # initial_delay_max_duration should be 3m-5m for production use. The below setting is a just example.
      initial_delay_max_duration: 60s
    EOF
    ```

1. Create a backup directory

    To avoid removing the indexing data due to any trouble after finishing indexing, we should prepare the path for auto backup.

    ```bash
    mkdir -p backup
    ```

1. Deploy Vald Agent on Docker

    To deploy Vald agent on docker with `config.yaml`, you can run the below command.
    Note: 
      - Please check whether there is `config.yaml` file and `backup` directory in your current directory.

    ```bash
    docker run -v $(pwd)/:/etc/server -p 8081:8081 --rm -it vdaas/vald-agent-ngt
    ```

1. Verify

    If the deployment success, you can confirm the output will be similar to below.

    ```bash
    2020-07-01 03:02:41	[INFO]:	maxprocs: Leaving GOMAXPROCS=4: CPU quota undefined
    2020-07-01 03:02:41	[INFO]:	service agent ngt v0.0.0 starting...
    2020-07-01 03:02:41	[INFO]:	executing daemon pre-start function
    2020-07-01 03:02:41	[INFO]:	executing daemon start function
    2020-07-01 03:02:41	[INFO]:	server agent-grpc executing preStartFunc
    2020-07-01 12:02:41	[INFO]:	gRPC server agent-grpc starting on 0.0.0.0:8081
    ```

## Run using example code

1. Download dataset

    In this tutorial. we use [fashion-mnist](https://github.com/zalandoresearch/fashion-mnist) as a dataset for indexing and search query.

    ```bash
    # move to the working directory
    cd example/client/agent
    
    # download fashion-mnist testing dataset
    wget http://ann-benchmarks.com/fashion-mnist-784-euclidean.hdf5
    ```

1. Running example

    Vald provides multiple language client libraries such as Go, Java, Node.js, Python, and so on.<br>
    This example will insert and index 400 vectors into the Vald from the fashion-mnist dataset via gRPC.
    And then after waiting for indexing, it will request for searching the nearest vector 10 times.
    You will get the 10 nearest neighbor vectors for each search query.<br>
    
    We use [`example/client/agent/main.go`](https://github.com/vdaas/vald/blob/master/example/client/agent/main.go) to run the example.
    The example code is the same as running an example only Vald agent on Kubernetes.
    If you want to learn the detail of running an example, please refer to the tutorial of [standalone Vald Agent on kubernetes](../tutorial/get-started.md/#run-using-example-code-1).

    ```bash
    # run example
    go run main.go
    ```
    Note:
      - We recommend you to run `CreateIndex()` after `Insert()` without waiting for auto indexing.
      - When finish indexing completely, the backup files (metadata.json and ngt-meta.kvsdb) can be confirmed in your mount directory.

1. Clean Up

    Stop the Vald Agent docker container via `Ctrl+C`.

## Recommended Documents

Congratulation! You achieved this tutorial!

If you want, you can try other tutorials such as:
- [Vald Agent Standalone on k8s](../tutorial/vald-agent-standalone-on-k8s.md)

For more information, we recommend you to check:
- [Architecture](../overview/architecture.md)
- [Configuration](../user-guides/configuration.md)
- [Operations](../user-guides/operations.md)