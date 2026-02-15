//
// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

use agent::config::AgentConfig;
use agent::serve;
use config::Config;
use opentelemetry::global;
use opentelemetry_sdk::metrics::{InMemoryMetricExporter, PeriodicReader, SdkMeterProvider};
use proto::core::v1::agent_client::AgentClient;
use proto::payload::v1::{control, insert, object, remove, search, update, upsert};
use proto::vald::v1::insert_client::InsertClient;
use proto::vald::v1::object_client::ObjectClient;
use proto::vald::v1::remove_client::RemoveClient;
use proto::vald::v1::search_client::SearchClient;
use proto::vald::v1::update_client::UpdateClient;
use proto::vald::v1::upsert_client::UpsertClient;
use rand::prelude::*;
use rand_distr::{Distribution, Normal};
use std::time::Duration;
use tempfile::TempDir;
use tonic::transport::Channel;

#[tokio::test]
async fn test_agent_integration_qbg() -> Result<(), Box<dyn std::error::Error>> {
    // 1. Setup Observability
    let exporter = InMemoryMetricExporter::default();
    let reader = PeriodicReader::builder(exporter.clone()).build();
    let provider = SdkMeterProvider::builder().with_reader(reader).build();
    global::set_meter_provider(provider.clone());

    // 2. Setup Config
    let temp_dir = TempDir::new()?;
    let index_path = temp_dir.path().join("index");
    let mut rng = rand::rng();
    let port = 50051 + (rng.random::<u16>() % 1000); // Random port to avoid conflicts
    let addr = format!("http://127.0.0.1:{}", port);

    // We need to match AgentConfig structure. It might be complex to set via set_default for nested structs.
    // Better to use a YAML string source.
    let config_yaml = format!(
        r#"
logging:
  level: debug
service:
  type: qbg
qbg:
  dimension: 128
  index_path: "{}"
  distance_type: 1 # L2
  object_type: 1 # Float
observability:
  enabled: true
  meter:
    enabled: true
    export_duration_secs: 1
server_config:
  servers:
    - name: grpc
      host: 0.0.0.0
      port: {}
      grpc:
        max_receive_message_size: 4194304
        max_send_message_size: 4194304
        initial_window_size: 65535
        initial_conn_window_size: 65535
        max_header_list_size: 8192
        max_concurrent_streams: 100
        keepalive:
          time: 120s
          timeout: 20s
        interceptors:
          - MetricInterceptor
"#,
        index_path.to_str().unwrap(),
        port
    );

    let config_source = ::config::File::from_str(&config_yaml, ::config::FileFormat::Yaml);
    let settings = Config::builder().add_source(config_source).build()?;
    let mut agent_config: AgentConfig = settings.try_deserialize()?;
    agent_config.bind(); // Not strictly needed for tests but good practice

    // 3. Start Agent in background
    let config_clone = agent_config.clone();
    tokio::spawn(async move {
        if let Err(e) = serve(config_clone).await {
            eprintln!("Agent failed: {:?}", e);
        }
    });

    // 4. Wait for server readiness
    let mut client_endpoint = None;
    for _ in 0..40 {
        tokio::time::sleep(Duration::from_millis(500)).await;
        let channel = Channel::from_shared(addr.clone())?.connect().await;
        if channel.is_ok() {
            client_endpoint = Some(channel.unwrap());
            break;
        }
    }
    let channel = client_endpoint.expect("Failed to connect to agent");

    // 5. Create Clients
    let mut insert_client = InsertClient::new(channel.clone());
    let mut search_client = SearchClient::new(channel.clone());
    let mut update_client = UpdateClient::new(channel.clone());
    let mut remove_client = RemoveClient::new(channel.clone());
    let mut object_client = ObjectClient::new(channel.clone());
    let mut upsert_client = UpsertClient::new(channel.clone());
    let mut agent_client = AgentClient::new(channel.clone());

    // 6. Generate Data (Normal Distribution)
    let normal = Normal::new(0.0, 1.0)?;
    let generate_vector = |rng: &mut ThreadRng, dim: usize| -> Vec<f32> {
        (0..dim).map(|_| normal.sample(rng)).collect()
    };

    // 7. Test Insert
    let num_vectors = 1000;
    for i in 0..num_vectors {
        let req = insert::Request {
            vector: Some(object::Vector {
                id: format!("uuid-{}", i),
                vector: generate_vector(&mut rng, 128),
                timestamp: 0,
            }),
            config: Some(insert::Config::default()),
        };
        insert_client.insert(req).await?;
    }

    // 8. Create Index (Required for search in QBG usually)
    agent_client
        .create_index(control::CreateIndexRequest { pool_size: 10 })
        .await?;

    // 9. Test Search
    let search_vec = generate_vector(&mut rng, 128);
    let search_req = search::Request {
        vector: search_vec.clone(),
        config: Some(search::Config {
            num: 5,
            radius: -1.0,
            epsilon: 0.1,
            timeout: 3000000000, // 3s
            ..Default::default()
        }),
    };

    let res = search_client.search(search_req.clone()).await?;
    let response = res.into_inner();
    println!("Search results: {:?}", response.results);
    // Note: Search results might be empty depending on data distribution and clustering,
    // but the call should succeed.

    // 10. Test Update
    let update_id = "uuid-0";
    let new_vec = generate_vector(&mut rng, 128);
    let update_req = update::Request {
        vector: Some(object::Vector {
            id: update_id.to_string(),
            vector: new_vec.clone(),
            timestamp: 0,
        }),
        config: Some(update::Config::default()),
    };
    update_client.update(update_req).await?;

    // Verify update via GetObject
    let get_req = object::VectorRequest {
        id: Some(object::Id {
            id: update_id.to_string(),
        }),
        filters: None,
    };
    let res = object_client.get_object(get_req).await?;
    // Note: QBG might return approximate vector or quantized, but if using float32, it should be close.
    // QBG stores compressed vectors if configured? config default uses compression?
    // In our config: object_type: 1 (Float).
    // But QBG might not return exact vector if it's in the index?
    // GetObject checks KVStore (which we enabled).
    // KVS stores exact vector if not compressed?
    // Let's verify vector length at least.
    let retrieved_vec = res.into_inner().vector;
    assert_eq!(retrieved_vec.len(), 128);
    // Checking exact equality might fail if compression is enabled by default in AgentConfig for KVS.
    // `kvsdb.use_compression` defaults to true in `config.rs` memory.
    // If compressed, it might be lossy? No, KVS compression (zstd/lz4) is lossless.
    // So exact match should work.
    assert_eq!(retrieved_vec, new_vec);

    // 11. Test Upsert
    let upsert_id = "uuid-new";
    let upsert_vec = generate_vector(&mut rng, 128);
    let upsert_req = upsert::Request {
        vector: Some(object::Vector {
            id: upsert_id.to_string(),
            vector: upsert_vec.clone(),
            timestamp: 0,
        }),
        config: Some(upsert::Config::default()),
    };
    upsert_client.upsert(upsert_req).await?;

    // Verify Upsert
    let get_req_upsert = object::VectorRequest {
        id: Some(object::Id {
            id: upsert_id.to_string(),
        }),
        filters: None,
    };
    let res = object_client.get_object(get_req_upsert).await?;
    assert_eq!(res.into_inner().vector, upsert_vec);

    // 12. Test Remove
    let remove_id = "uuid-1";
    let remove_req = remove::Request {
        id: Some(object::Id {
            id: remove_id.to_string(),
        }),
        config: Some(remove::Config::default()),
    };
    remove_client.remove(remove_req).await?;

    // Verify Remove
    let get_req_remove = object::VectorRequest {
        id: Some(object::Id {
            id: remove_id.to_string(),
        }),
        filters: None,
    };
    let res = object_client.get_object(get_req_remove).await;
    // Should fail with NotFound
    assert!(res.is_err());

    // 13. Verify Observability
    // Force flush metrics to ensure collection happens
    if let Err(e) = provider.force_flush() {
        eprintln!("Failed to flush metrics: {:?}", e);
    }

    // Check exporter
    let metrics = exporter.get_finished_metrics()?;

    // We look for any metric starting with "agent_core_ngt_" or "vald_agent_"
    // The metric names in `metrics.rs` are `agent_core_ngt_index_count`, etc.
    let index_count_metric = metrics
        .iter()
        .flat_map(|rm| rm.scope_metrics())
        .flat_map(|sm| sm.metrics())
        .find(|m| m.name() == "agent_core_ngt_index_count");

    // Since we inserted and created index, index count should be tracked.
    assert!(
        index_count_metric.is_some(),
        "Metric agent_core_ngt_index_count not found. Available metrics: {:?}",
        metrics
            .iter()
            .flat_map(|rm| rm.scope_metrics())
            .flat_map(|sm| sm.metrics())
            .map(|m| m.name().to_string())
            .collect::<Vec<_>>()
    );

    // Also check request metrics if `MetricInterceptor` is working.
    // It should produce metrics like `vald_agent_grpc_server_handled_total`?
    // Middleware metrics depend on `observability` crate implementation.

    Ok(())
}
