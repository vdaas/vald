# Building gRPC proto

This page shows how to build gRPC proto file for calling API to your Vald Cluster.

<div class="notice">
Vald provides the official client libraries (see: <a href="https://vald.vdaas.org/docs/user-guides/sdks">SDK document</a>).
If you can use one of the SDKs we recommend using it.
</div>

## Target proto files

Vald defines the proto file for each API.
Let's check the below table for the details.

| API service name |                                          proto                                          |                                                                                             dependencies                                                                                              | usage                                                              |
| :--------------: | :-------------------------------------------------------------------------------------: | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------: | :----------------------------------------------------------------- |
|      Insert      | [insert.proto](https://github.com/vdaas/vald/blob/main/apis/proto/v1/vald/insert.proto) | [payload.proto](https://github.com/vdaas/vald/blob/main/apis/proto/v1/payload/payload.proto)<br>[annotation.proto](https://github.com/googleapis/googleapis/blob/master/google/api/annotations.proto) | Insert vectors into Vald Agent                                     |
|      Update      | [update.proto](https://github.com/vdaas/vald/blob/main/apis/proto/v1/vald/update.proto) | [payload.proto](https://github.com/vdaas/vald/blob/main/apis/proto/v1/payload/payload.proto)<br>[annotation.proto](https://github.com/googleapis/googleapis/blob/master/google/api/annotations.proto) | Update vectors stored in Vald Agent                                |
|      Upsert      | [upsert.proto](https://github.com/vdaas/vald/blob/main/apis/proto/v1/vald/upsert.proto) | [payload.proto](https://github.com/vdaas/vald/blob/main/apis/proto/v1/payload/payload.proto)<br>[annotation.proto](https://github.com/googleapis/googleapis/blob/master/google/api/annotations.proto) | Update vectors stored Vald Agent or Insert vectors into Vald Agent |
|      Search      | [search.proto](https://github.com/vdaas/vald/blob/main/apis/proto/v1/vald/search.proto) | [payload.proto](https://github.com/vdaas/vald/blob/main/apis/proto/v1/payload/payload.proto)<br>[annotation.proto](https://github.com/googleapis/googleapis/blob/master/google/api/annotations.proto) | Search similar vectors with query                                  |
|      Remove      | [remove.proto](https://github.com/vdaas/vald/blob/main/apis/proto/v1/vald/remove.proto) | [payload.proto](https://github.com/vdaas/vald/blob/main/apis/proto/v1/payload/payload.proto)<br>[annotation.proto](https://github.com/googleapis/googleapis/blob/master/google/api/annotations.proto) | Remove stored vectors from Vald Agent                              |
|      Object      | [object.proto](https://github.com/vdaas/vald/blob/main/apis/proto/v1/vald/object.proto) | [payload.proto](https://github.com/vdaas/vald/blob/main/apis/proto/v1/payload/payload.proto)<br>[annotation.proto](https://github.com/googleapis/googleapis/blob/master/google/api/annotations.proto) | Get object information of vector stored in Vald Agent              |

## How to build protobuf

### The way to build proto files

Let's build proto files using your favorite programming language.

There are 3 steps to building API proto:

1. Install gRPC tools

   - gRPC official document provides [the way to install for each language](https://grpc.io/docs/languages/).<br>
     If your favorite programming language is not there, you can find 3rd party tools for building.

1. Download Vald api proto files and external dependence

   - [vald api proto](https://github.com/vdaas/vald/tree/main/apis/proto/v1/vald)
   - [vald payload proto](https://github.com/vdaas/vald/tree/main/apis/proto/v1/payload)
   - [googleapis](https://github.com/googleapis/googleapis)
   - [PGV](https://github.com/envoyproxy/protoc-gen-validate)

1. Choose proto file(s) and build

### Example: Build proto files in Rust

This section shows the example steps for building proto files using Rust.

There are many tools for building proto in Rust, we use [tonic](https://github.com/hyperium/tonic) as an example.

1.  Check version

    This example runs in this environment.

    ```bash
    $ cargo version
    cargo 1.58.0 (7f08ace4f 2021-11-24)
    $ rustc -V
    rustc 1.58.0 (02072b482 2022-01-11)
    $ rustup -V
    rustup 1.24.3 (ce5817a94 2021-05-31)
    info: This is the version for the rustup toolchain manager, not the rustc compiler.
    info: The currently active `rustc` version is `rustc 1.57.0 (f1edd0429 2021-11-29)`
    $ rustup show
    Default host: x86_64-unknown-linux-gnu
    rustup home:  /home/user/.rustup

    stable-x86_64-unknown-linux-gnu (default)
    rustc 1.58.0 (02072b482 2022-01-11)
    ```

1.  Create project

    ```bash
    cargo new --lib vald-grpc
    ```

1.  Edit `Cargo.toml`

    ```bash
    cd vald-grpc && \
    vim Cargo.toml
    ---
    [package]
    name = "vald-grpc"
    version = "0.1.0"
    edition = "2021"

    # See more keys and their definitions at https://doc.rust-lang.org/cargo/reference/manifest.html

    [[bin]]
    name = "client"
    path = "src/client.rs"

    [dependencies]
    tonic = "0.6.2"
    tokio = "1.15"
    prost = "0.9"
    prost-types = "0.9"
    hdf5 = "0.8.1"
    chrono = "^0.4"

    [build-dependencies]
    tonic-build = "0.6.2"
    ```

1.  Download Vald proto files and dependence proto files

    ```bash
    // create proto file root dir
    mkdir -p proto
    ```

    ```bash
    // download vald proto files
    git clone https://github.com/vdaas/vald \
    && cp -R vald/apis/proto/v1/vald proto/vald \
    && cp -R vald/apis/proto/v1/payload proto/payload \
    && rm -rf vald
    ```

    ```bash
    // download googleapis
    git clone https://github.com/googleapis/googleapis \
    && cp -R googleapis/google proto/google \
    && rm -rf googleapis
    ```

    ```bash
    // download protoc-gen-validate
    git clone https://github.com/envoyproxy/protoc-gen-validate \
    && mv protoc-gen-validate proto/protoc-gen-validate
    ```

1.  Fixing import path

    ```bash
    find proto/vald -type f -name "*.proto" | xargs sed -i "s/apis\/proto\/v1\///g" && \
    find proto/vald -type f -name "*.proto" | xargs sed -i "s/github\.com\/googleapis\/googleapis\///g" && \
    find proto/payload -type f -name "*.proto" | xargs sed -i "s/github\.com\/googleapis\/googleapis\///g" && \
    find proto/payload -type f -name "*.proto" | xargs sed -i "s/github\.com\/envoyproxy\///g"
    ```

1.  Implement `build.rs` and Build proto

    1. `build.rs`

       ```rust
       fn main() -> Result<(), Box<dyn std::error::Error>> {
           let insert_proto = "./proto/vald/insert.proto";
           let update_proto = "./proto/vald/update.proto";
           let upsert_proto = "./proto/vald/upsert.proto";
           let search_proto = "./proto/vald/search.proto";
           let remove_proto = "./proto/vald/remove.proto";
           let object_proto = "./proto/vald/object.proto";

           tonic_build::configure()
               .build_client(true)
               .out_dir("./src/proto")
               .compile(
                   &[
                       insert_proto,
                       update_proto,
                       upsert_proto,
                       search_proto,
                       remove_proto,
                       object_proto,
                   ],
                   &["./proto"],
               )
               .unwrap_or_else(|e| panic!("protobuf compile error: {}", e));
           Ok(())
       }
       ```

    1. build proto

       ```bash
       cargo build
       ```

1.  Edit `Cargo.toml`

    ```bash
    cd vald-grpc && \
    vim Cargo.toml
    ---
    [package]
    ...

    [[bin]]
    name = "client"
    path = "src/client.rs"
    ...
    ```

1.  Implement code using client

    1.  `lib.rs`

        Import build proto in `src/lib.rs`

        ```rust
        pub mod vald {
            pub mod v1 {
                include!("./proto/vald.v1.rs");
            }
        }

        pub mod payload {
            pub mod v1 {
                include!("./proto/payload.v1.rs");
            }
        }

        pub mod google {
            pub mod rpc {
                include!("./proto/google.rpc.rs");
            }
            pub mod api {
                include!("./proto/google.api.rs");
            }
            pub mod protobuf {
                include!("./proto/google.protobuf.rs");
            }
        }
        ```

    1.  `src/clinet.rs`

        There are 4 steps in `src/clinet.rs`:

        1. Load dataset
        1. Insert vector to Vald cluster
        1. Search nearest neighbor vectors from Vald cluster after indexing finished
        1. Remove indexed vectors from Vald cluster

        The example is here:

        ```rust
        // import packages
        use chrono::Utc;
        use hdf5::{File, Result};
        use std::thread::sleep;
        use std::time;
        use tonic::transport::Endpoint;

        // import vald protos
        use vald_sample_rust_client::payload::v1::insert;
        use vald_sample_rust_client::payload::v1::object;
        use vald_sample_rust_client::payload::v1::remove;
        use vald_sample_rust_client::payload::v1::search;
        use vald_sample_rust_client::vald::v1::insert_client;
        use vald_sample_rust_client::vald::v1::remove_client;
        use vald_sample_rust_client::vald::v1::search_client;

        // Dataset file name
        static FILE: &str = "fashion-mnist-784-euclidean.hdf5";
        // Dataset name
        static DATASET: &str = "train";
        // set Vald cluster host
        static HOST: &str = "http://localhost:8080";
        // Time duration for waiting to finish `CreateIndex` and `SaveIndex`
        static DURATION: u64 = 15;

        // load data
        fn read_file() -> Result<Vec<Vec<f32>>> {
            let file = File::open(FILE).unwrap_or_else(|e| panic!("[ERR] failed to read file: {}", e));
            let data = file
                .dataset(DATASET)
                .unwrap_or_else(|e| panic!("[ERR] failed to get dataset: {}", e));
            let mut vector = Vec::new();
            for train in data.read_2d::<f32>()?.outer_iter() {
                let mut vec: Vec<f32> = Vec::new();
                vec.append(&mut train.to_vec());
                vector.push(vec);
                if vector.len() == 500 {
                    break;
                }
            }
            Ok(vector)
        }

        #[tokio::main(flavor = "current_thread")]
        async fn main() -> Result<(), Box<dyn std::error::Error>> {
            print!("[Start] Load {} file\n", FILE);
            let vec = read_file()?;
            print!("[End] Success to load {} file\n", FILE);

            print!("[Start] Insert phase\n");
            // create insert client
            let mut insert_client =
                insert_client::InsertClient::connect(Endpoint::from_static(HOST)).await?;
            let mut ids: Vec<String> = Vec::new();
            for v in vec.iter() {
                let id = Utc::now().timestamp_nanos().to_string();
                ids.push(id.to_string());
                // insert vector
                let _ = insert_client
                    .insert(insert::Request {
                        vector: Some(object::Vector {
                            id: id.to_string(),
                            vector: v.to_vec(),
                        }),
                        config: Some(insert::Config {
                            skip_strict_exist_check: true,
                            filters: None,
                            timestamp: Utc::now().timestamp(),
                        }),
                    })
                    .await?;
            }
            print!("[End] Finish Insert Phase\n");

            print!("[Sleep] Waiting SaveIndex is completed...\n");
            sleep(time::Duration::from_secs(DURATION));

            print!("[Start] Search phase\n");
            // create search client
            let mut search_client =
                search_client::SearchClient::connect(Endpoint::from_static(HOST)).await?;
            for id in ids.clone() {
                // search nearest neighbor vectors using searchById method
                let res = search_client
                    .search_by_id(search::IdRequest {
                        id: id.to_string(),
                        config: Some(search::Config {
                            request_id: id.to_string(),
                            num: 10,
                            radius: -1.0,
                            epsilon: -1.0,
                            timeout: 500,
                            ingress_filters: None,
                            egress_filters: None,
                        }),
                    })
                    .await?;
                print!("[Id]: {:?}\n", id.to_string());
                print!("[Result]: {:#?}\n", res.into_inner().results);
            }
            print!("[End] Finish Search Phase\n");

            print!("[Start] Remove phase\n");
            // create remove client
            let mut remove_client =
                remove_client::RemoveClient::connect(Endpoint::from_static(HOST)).await?;
            for id in ids.clone() {
                // remove vectors
                let _ = remove_client
                    .remove(remove::Request {
                        id: Some(object::Id { id: id.to_string() }),
                        config: None,
                    })
                    .await?;
            }
            print!("[End] Finish Remove phase\n");

            Ok(())
        }
        ```

1.  Build

    Build implemented codes before running example code.

    ```bash
    cargo build
    ```

1.  Running example

    After creating Vald cluster, you can run an example code by the following command.

    ```bash
    cargo run src/client.rs
    ```
