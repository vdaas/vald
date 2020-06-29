# Deploy Vald Agent on Docker

Vald is designed and implemented based on Cloud-Native architecture.
However, there may be cases that want to use only Vald Agent without Kubernetes.

This article will show you how to deploy and run the Vald Agent on Docker.
Fashion-mnist is used as an example dataset, same as [Get Stared](../tutorial/get-started.md)

## Requirements

- docker: v19.0 ~
- go: v1.14 ~
- libhdf5 (_only required for this tutorial._)


Hdf5 is required for this tutorial. If hdf5 is not installed, please install [hdf5](https://www.hdfgroup.org/).
<details><summary>[Optional] Install hdf5</summary><br>

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

1. Create `config.yaml`

    The configuration of Vald agent for docker is set using `config.yaml`<br>
    You can also check [the sample](https://github.com/vdaas/vald/blob/master/cmd/agent/core/ngt/sample.yaml).

    ```bash
    cat << EOF > config.yaml
    ---
    version: v0.0.0
    time_zone: JST
    logging:
      logger: glg
      level: debug
      format: json
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
      # vector dimension
      dimension: 784
      # bulk insert chunk size
      bulk_insert_chunk_size: 10
      # distance_type, which should be `l1`, `l2`, `angle`, `hamming`, `cosine`, `normalizedangle` or `nomralizedcosine`
      distance_type: l2
      # object_type, which should be `float` or `uint8`
      object_type: float
      # creation edge size
      creation_edge_size: 20
      # search edge size
      search_edge_size: 10
      # in-memory mode enabled
      enable_in_memory_mode: true
      # The limit duration of automatic indexing 
      auto_index_duration_limit: 1m
      # Check duration of automatic indexing.
      auto_index_check_duration: 10s
      # The number of cache to trigger automatic indexing
      auto_index_length: 100
    EOF
    ```

1. Deploy Vald Agent on Docker

    To deploy Vald agent on docker with `config.yaml`, you can run below command.

    ```bash
    docker run -v $(pwd)/config.yaml:/etc/server/config.yaml -p 8081:8081 --rm -it vdaas/vald-agent-ngt
    ```

1. Verify

    If the deployment success, you can confirm the output will be similar to below.

    ```bash
    {"Date":"2020-06-24T05:49:53.615263Z","Level":"INFO","Detail":"maxprocs: Leaving GOMAXPROCS=4: CPU quota undefined"}
    {"Date":"2020-06-24T05:49:53.615263Z","Level":"INFO","Detail":"service agent ngt v0.0.0 starting..."}
    {"Date":"2020-06-24T05:49:53.615263Z","Level":"INFO","Detail":["daemon start"]}
    {"Date":"2020-06-24T05:49:53.620263Z","Level":"INFO","Detail":"server agent-grpc executing preStartFunc"}
    {"Date":"2020-06-24T05:49:53.620263Z","Level":"INFO","Detail":"gRPC server agent-grpc starting on 0.0.0.0:8081"}
    ```

## Run using example code

1. Download dataset

    In this tutorial. we use [fashion-mnist](https://github.com/zalandoresearch/fashion-mnist) as a dataset for indexing and search query.

    ```bash
    # move to working directory
    cd example/client/agent
    
    # download fashion-mnist testing dataset
    wget http://ann-benchmarks.com/fashion-mnist-784-euclidean.hdf5
    ```

1. Running example

    Vald provides multiple language client libraries such as Go, Java, Node.js, Python, and so on.<br>
    In this example, the fashion-mnist dataset will insert into the Vald and search using [vald-client-go](https://github.com/vdaas/vald-client-go).
    
    We use [`example/client/agent/main.go`](https://github.com/vdaas/vald/blob/master/example/client/agent/main.go) to run the example.
    The example code is the same as running an example only Vald agent on Kubernetes.
    If you want to learn the detail of running an example, please refer to the tutorial of [standalone Vald Agent on kubernetes](../tutorial/get-started.md/#run-using-example-code-1).

    ```bash
    # run example
    go run main.go
    ```
    Note:
      - We recommend you to run `CreateIndex()` after `Insert()` without waiting auto indexing.

1. Clean Up

    Remove the Vald Agent docker via `Ctrl+C`.

