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

mod common;
/// Flush RPC handlers.
pub mod flush;
/// Index RPC handlers.
pub mod index;
/// Insert RPC handlers.
pub mod insert;
/// Object RPC handlers.
pub mod object;
/// Remove RPC handlers.
pub mod remove;
/// Search RPC handlers.
pub mod search;
/// Update RPC handlers.
pub mod update;
/// Upsert RPC handlers.
pub mod upsert;

use crate::config::AgentConfig;
use crate::middleware;
use crate::service::{DaemonConfig, DaemonHandle, start_daemon};
use proto::{
    core::v1::agent_server,
    vald::v1::{
        flush_server, index_server, insert_server, object_server, remove_server, search_server,
        update_server, upsert_server,
    },
};
use std::sync::Arc;
use std::time::Duration;
use tokio::sync::{RwLock, mpsc};

/// Agent service wrapper for running the ANN implementation and gRPC server.
pub struct Agent<S: algorithm::ANN + 'static> {
    s: Arc<RwLock<S>>,
    name: String,
    ip: String,
    resource_type: String,
    api_name: String,
    stream_concurrency: usize,
    daemon_handle: Option<DaemonHandle>,
    error_rx: Option<mpsc::Receiver<algorithm::Error>>,
}

impl<S: algorithm::ANN + 'static> Agent<S> {
    /// Creates a new agent instance with its service and identity settings.
    pub fn new(
        s: S,
        name: &str,
        ip: &str,
        resource_type: &str,
        api_name: &str,
        stream_concurrency: usize,
    ) -> Self {
        Self {
            s: Arc::new(RwLock::new(s)),
            name: name.to_string(),
            ip: ip.to_string(),
            resource_type: resource_type.to_string(),
            api_name: api_name.to_string(),
            stream_concurrency,
            daemon_handle: None,
            error_rx: None,
        }
    }

    /// Starts the daemon for automatic indexing and saving.
    /// This should be called before serve_grpc.
    pub async fn start(&mut self, config: &AgentConfig) {
        let daemon_config = DaemonConfig::from_config(&config.daemon);
        log::info!("Starting daemon with config: {:?}", daemon_config);

        let (handle, error_rx) = start_daemon(self.s.clone(), daemon_config).await;
        self.daemon_handle = Some(handle);
        self.error_rx = Some(error_rx);

        log::info!("Daemon started successfully");
    }

    /// Stops the daemon gracefully.
    pub fn stop(&self) {
        if let Some(ref handle) = self.daemon_handle {
            log::info!("Stopping daemon...");
            handle.stop();
            log::info!("Daemon stop signal sent");
        }
    }

    /// Performs a graceful shutdown of the agent.
    ///
    /// This method:
    /// 1. Stops the daemon and waits for it to complete final index creation
    /// 2. Calls close() on the underlying service to:
    ///    - Create and save any uncommitted index changes
    ///    - Close the QBG index
    ///    - Flush and close KVS
    ///
    /// This should be called when the application is shutting down to ensure
    /// all data is persisted correctly.
    pub async fn shutdown(&self) -> Result<(), algorithm::Error> {
        log::info!("Agent shutdown initiated...");

        // Stop daemon and wait for it to complete
        if let Some(ref handle) = self.daemon_handle {
            log::info!("Waiting for daemon to complete shutdown...");
            handle.stop_and_wait().await;
            log::info!("Daemon shutdown complete");
        }

        // Close the service
        log::info!("Closing service...");
        let mut service = self.s.write().await;
        let result = service.close().await;

        match &result {
            Ok(()) => log::info!("Agent shutdown complete"),
            Err(e) => log::error!("Agent shutdown completed with errors: {:?}", e),
        }

        result
    }

    /// Returns the service wrapped in Arc<RwLock<S>> for external access.
    pub fn service(&self) -> Arc<RwLock<S>> {
        self.s.clone()
    }

    /// Starts the gRPC server with all registered services.
    pub async fn serve_grpc(self, config: AgentConfig) -> Result<(), Box<dyn std::error::Error>> {
        let addr = "0.0.0.0:8081".parse()?;

        let grpc_server_config = config
            .server_config
            .servers
            .iter()
            .find(|s| s.name == "grpc")
            .map(|s| &s.grpc)
            .ok_or_else(|| {
                std::io::Error::new(std::io::ErrorKind::NotFound, "grpc server config not found")
            })?;

        let mut builder = tonic::transport::Server::builder();
        if let Some(duration) =
            parse_duration_from_string(&grpc_server_config.keepalive.max_conn_age)
        {
            builder = builder.max_connection_age(duration);
        }
        if let Some(duration) = parse_duration_from_string(&grpc_server_config.connection_timeout) {
            builder = builder.timeout(duration);
        }

        let mut accessloginterceptor: Option<()> = None;
        let mut metricinterceptor: Option<()> = None;
        for name in &grpc_server_config.interceptors {
            match name.to_lowercase().as_str() {
                "accessloginterceptor" | "accesslog" => accessloginterceptor = Some(()),
                "metricinterceptor" | "metric" => metricinterceptor = Some(()),
                _ => {}
            }
        }

        let layer = tower::ServiceBuilder::new()
            .option_layer(
                accessloginterceptor.map(|_| middleware::AccessLogMiddlewareLayer::default()),
            )
            .option_layer(metricinterceptor.map(|_| middleware::MetricMiddlewareLayer::default()))
            .into_inner();

        let max_recv_size = grpc_server_config.max_receive_message_size;
        let max_send_size = grpc_server_config.max_send_message_size;

        builder
            .initial_stream_window_size(Some(grpc_server_config.initial_window_size))
            .initial_connection_window_size(Some(grpc_server_config.initial_conn_window_size))
            .http2_keepalive_interval(parse_duration_from_string(
                &grpc_server_config.keepalive.time,
            ))
            .http2_keepalive_timeout(parse_duration_from_string(
                &grpc_server_config.keepalive.timeout,
            ))
            .http2_max_header_list_size(Some(grpc_server_config.max_header_list_size))
            .max_concurrent_streams(Some(grpc_server_config.max_concurrent_streams))
            .layer(layer)
            .add_service(
                agent_server::AgentServer::new(self.clone())
                    .max_decoding_message_size(max_recv_size)
                    .max_encoding_message_size(max_send_size),
            )
            .add_service(
                search_server::SearchServer::new(self.clone())
                    .max_decoding_message_size(max_recv_size)
                    .max_encoding_message_size(max_send_size),
            )
            .add_service(
                insert_server::InsertServer::new(self.clone())
                    .max_decoding_message_size(max_recv_size)
                    .max_encoding_message_size(max_send_size),
            )
            .add_service(
                update_server::UpdateServer::new(self.clone())
                    .max_decoding_message_size(max_recv_size)
                    .max_encoding_message_size(max_send_size),
            )
            .add_service(
                upsert_server::UpsertServer::new(self.clone())
                    .max_decoding_message_size(max_recv_size)
                    .max_encoding_message_size(max_send_size),
            )
            .add_service(
                remove_server::RemoveServer::new(self.clone())
                    .max_decoding_message_size(max_recv_size)
                    .max_encoding_message_size(max_send_size),
            )
            .add_service(
                object_server::ObjectServer::new(self.clone())
                    .max_decoding_message_size(max_recv_size)
                    .max_encoding_message_size(max_send_size),
            )
            .add_service(
                index_server::IndexServer::new(self.clone())
                    .max_decoding_message_size(max_recv_size)
                    .max_encoding_message_size(max_send_size),
            )
            .add_service(
                flush_server::FlushServer::new(self.clone())
                    .max_decoding_message_size(max_recv_size)
                    .max_encoding_message_size(max_send_size),
            )
            .serve(addr)
            .await?;

        Ok(())
    }
}

impl<S: algorithm::ANN + 'static> Clone for Agent<S> {
    fn clone(&self) -> Self {
        Self {
            s: self.s.clone(),
            name: self.name.clone(),
            ip: self.ip.clone(),
            resource_type: self.resource_type.clone(),
            api_name: self.api_name.clone(),
            stream_concurrency: self.stream_concurrency,
            daemon_handle: self.daemon_handle.clone(),
            error_rx: None, // error_rx is not cloneable, only main instance handles errors
        }
    }
}

impl<S: algorithm::ANN + 'static> Drop for Agent<S> {
    fn drop(&mut self) {
        self.stop();
    }
}

/// Parses a duration string like "30s", "5m", "1h" into a Duration.
fn parse_duration_from_string(input: &str) -> Option<Duration> {
    if input.len() < 2 {
        return None;
    }
    let last_char = match input.chars().last() {
        Some(c) => c,
        None => return None,
    };
    if last_char.is_numeric() {
        return None;
    }

    let (value, unit) = input.split_at(input.len() - 1);
    let num: u64 = match value.parse() {
        Ok(n) => n,
        Err(_) => return None,
    };
    match unit {
        "s" => Some(Duration::from_secs(num)),
        "m" => Some(Duration::from_secs(num * 60)),
        "h" => Some(Duration::from_secs(num * 60 * 60)),
        _ => None,
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use algorithm::{ANN, Error};
    use proto::payload::v1::{info, insert, object, remove, search, update, upsert};
    use proto::vald::v1::{
        insert_server::Insert, object_server::Object, remove_server::Remove, search_server::Search,
    };
    use std::collections::HashMap;

    /// Minimal mock ANN service for handler testing.
    /// Returns fixed responses without business logic.
    struct MockANNService {
        dimension: usize,
    }

    impl MockANNService {
        fn new(dimension: usize) -> Self {
            Self { dimension }
        }
    }

    impl ANN for MockANNService {
        fn get_dimension_size(&self) -> usize {
            self.dimension
        }

        fn search(
            &self,
            _vector: Vec<f32>,
            num: u32,
            _epsilon: f32,
            _radius: f32,
        ) -> impl std::future::Future<Output = Result<search::Response, Error>> + Send {
            async move {
                Ok(search::Response {
                    request_id: String::new(),
                    results: (0..num)
                        .map(|i| object::Distance {
                            id: format!("result-{}", i),
                            distance: 0.1 * i as f32,
                        })
                        .collect(),
                })
            }
        }

        fn search_by_id(
            &self,
            _uuid: String,
            num: u32,
            _epsilon: f32,
            _radius: f32,
        ) -> impl std::future::Future<Output = Result<search::Response, Error>> + Send {
            async move {
                Ok(search::Response {
                    request_id: String::new(),
                    results: (0..num)
                        .map(|i| object::Distance {
                            id: format!("result-{}", i),
                            distance: 0.1 * i as f32,
                        })
                        .collect(),
                })
            }
        }

        fn linear_search(
            &self,
            _v: Vec<f32>,
            _n: u32,
        ) -> impl std::future::Future<Output = Result<search::Response, Error>> + Send {
            async {
                Err(Error::Unsupported {
                    method: "linear_search".into(),
                    algorithm: "Mock".into(),
                })
            }
        }

        fn linear_search_by_id(
            &self,
            _u: String,
            _n: u32,
        ) -> impl std::future::Future<Output = Result<search::Response, Error>> + Send {
            async {
                Err(Error::Unsupported {
                    method: "linear_search_by_id".into(),
                    algorithm: "Mock".into(),
                })
            }
        }

        fn insert(
            &mut self,
            _u: String,
            _v: Vec<f32>,
        ) -> impl std::future::Future<Output = Result<(), Error>> + Send {
            async { Ok(()) }
        }
        fn insert_with_time(
            &mut self,
            _u: String,
            _v: Vec<f32>,
            _t: i64,
        ) -> impl std::future::Future<Output = Result<(), Error>> + Send {
            async { Ok(()) }
        }
        fn insert_multiple(
            &mut self,
            _vs: HashMap<String, Vec<f32>>,
        ) -> impl std::future::Future<Output = Result<(), Error>> + Send {
            async { Ok(()) }
        }
        fn insert_multiple_with_time(
            &mut self,
            _vs: HashMap<String, Vec<f32>>,
            _t: i64,
        ) -> impl std::future::Future<Output = Result<(), Error>> + Send {
            async { Ok(()) }
        }
        fn update(
            &mut self,
            _u: String,
            _v: Vec<f32>,
        ) -> impl std::future::Future<Output = Result<(), Error>> + Send {
            async { Ok(()) }
        }
        fn update_with_time(
            &mut self,
            _u: String,
            _v: Vec<f32>,
            _t: i64,
        ) -> impl std::future::Future<Output = Result<(), Error>> + Send {
            async { Ok(()) }
        }
        fn update_multiple(
            &mut self,
            _vs: HashMap<String, Vec<f32>>,
        ) -> impl std::future::Future<Output = Result<(), Error>> + Send {
            async { Ok(()) }
        }
        fn update_multiple_with_time(
            &mut self,
            _vs: HashMap<String, Vec<f32>>,
            _t: i64,
        ) -> impl std::future::Future<Output = Result<(), Error>> + Send {
            async { Ok(()) }
        }
        fn update_timestamp(
            &mut self,
            _u: String,
            _t: i64,
            _f: bool,
        ) -> impl std::future::Future<Output = Result<(), Error>> + Send {
            async { Ok(()) }
        }
        fn remove(
            &mut self,
            _u: String,
        ) -> impl std::future::Future<Output = Result<(), Error>> + Send {
            async { Ok(()) }
        }
        fn remove_with_time(
            &mut self,
            _u: String,
            _t: i64,
        ) -> impl std::future::Future<Output = Result<(), Error>> + Send {
            async { Ok(()) }
        }
        fn remove_multiple(
            &mut self,
            _us: Vec<String>,
        ) -> impl std::future::Future<Output = Result<(), Error>> + Send {
            async { Ok(()) }
        }
        fn remove_multiple_with_time(
            &mut self,
            _us: Vec<String>,
            _t: i64,
        ) -> impl std::future::Future<Output = Result<(), Error>> + Send {
            async { Ok(()) }
        }

        fn get_object(
            &self,
            _uuid: String,
        ) -> impl std::future::Future<Output = Result<(Vec<f32>, i64), Error>> + Send {
            let dim = self.dimension;
            async move { Ok((vec![0.0; dim], 12345)) }
        }

        fn exists(&self, _uuid: String) -> impl std::future::Future<Output = (usize, bool)> + Send {
            async { (1, true) }
        }
        fn uuids(&self) -> impl std::future::Future<Output = Vec<String>> + Send {
            async { vec!["uuid-1".into()] }
        }
        fn list_object_func<F: FnMut(String, Vec<f32>, i64) -> bool + Send>(
            &self,
            _f: F,
        ) -> impl std::future::Future<Output = ()> + Send {
            async {}
        }
        fn create_index(&mut self) -> impl std::future::Future<Output = Result<(), Error>> + Send {
            async { Ok(()) }
        }
        fn save_index(&mut self) -> impl std::future::Future<Output = Result<(), Error>> + Send {
            async { Ok(()) }
        }
        fn create_and_save_index(
            &mut self,
        ) -> impl std::future::Future<Output = Result<(), Error>> + Send {
            async { Ok(()) }
        }
        fn regenerate_indexes(
            &mut self,
        ) -> impl std::future::Future<Output = Result<(), Error>> + Send {
            async { Ok(()) }
        }
        fn len(&self) -> u32 {
            100
        }
        fn insert_vqueue_buffer_len(&self) -> u32 {
            5
        }
        fn delete_vqueue_buffer_len(&self) -> u32 {
            2
        }
        fn is_indexing(&self) -> bool {
            false
        }
        fn is_flushing(&self) -> bool {
            false
        }
        fn is_saving(&self) -> bool {
            false
        }
        fn number_of_create_index_executions(&self) -> u64 {
            10
        }
        fn broken_index_count(&self) -> u64 {
            0
        }
        fn is_statistics_enabled(&self) -> bool {
            false
        }
        fn index_statistics(&self) -> Result<info::index::Statistics, Error> {
            Ok(info::index::Statistics::default())
        }
        fn index_property(&self) -> Result<info::index::Property, Error> {
            Ok(info::index::Property::default())
        }
        fn close(&mut self) -> impl std::future::Future<Output = Result<(), Error>> + Send {
            async { Ok(()) }
        }
    }

    fn create_test_agent(dimension: usize) -> Agent<MockANNService> {
        Agent::new(
            MockANNService::new(dimension),
            "test-agent",
            "127.0.0.1",
            "vald.v1",
            "vald-agent",
            10,
        )
    }

    fn gen_vector(dim: usize, seed: u64) -> Vec<f32> {
        let mut state = seed;
        (0..dim)
            .map(|i| {
                state = state
                    .wrapping_mul(6364136223846793005)
                    .wrapping_add(i as u64);
                ((state >> 33) as f32 / u32::MAX as f32) * 2.0 - 1.0
            })
            .collect()
    }

    // ==================== Insert Handler Tests ====================

    #[tokio::test]
    async fn test_insert_handler_success() {
        let agent = create_test_agent(128);

        let request = tonic::Request::new(insert::Request {
            vector: Some(object::Vector {
                id: "test-uuid-1".to_string(),
                vector: gen_vector(128, 1),
                timestamp: 0,
            }),
            config: Some(insert::Config {
                skip_strict_exist_check: false,
                timestamp: 0,
                filters: None,
            }),
        });

        let result = agent.insert(request).await;
        assert!(result.is_ok());

        let response = result.unwrap().into_inner();
        assert_eq!(response.uuid, "test-uuid-1");
        assert_eq!(response.name, "test-agent");
    }

    #[tokio::test]
    async fn test_insert_handler_duplicate_uuid() {
        let agent = create_test_agent(128);

        let vector = gen_vector(128, 1);

        // First insert
        let request1 = tonic::Request::new(insert::Request {
            vector: Some(object::Vector {
                id: "duplicate-uuid".to_string(),
                vector: vector.clone(),
                timestamp: 0,
            }),
            config: Some(insert::Config::default()),
        });
        let _ = agent.insert(request1).await.unwrap();

        // Second insert with same UUID - Mock always succeeds, so we just verify handler doesn't crash
        let request2 = tonic::Request::new(insert::Request {
            vector: Some(object::Vector {
                id: "duplicate-uuid".to_string(),
                vector: vector,
                timestamp: 0,
            }),
            config: Some(insert::Config::default()),
        });

        // With simplified mock, this succeeds (no duplicate check)
        let result = agent.insert(request2).await;
        assert!(result.is_ok());
    }

    #[tokio::test]
    async fn test_insert_handler_invalid_dimension() {
        let agent = create_test_agent(128);

        let request = tonic::Request::new(insert::Request {
            vector: Some(object::Vector {
                id: "test-uuid".to_string(),
                vector: gen_vector(64, 1), // Wrong dimension
                timestamp: 0,
            }),
            config: Some(insert::Config::default()),
        });

        let result = agent.insert(request).await;
        assert!(result.is_err());

        let status = result.unwrap_err();
        assert_eq!(status.code(), tonic::Code::InvalidArgument);
    }

    #[tokio::test]
    async fn test_insert_handler_missing_config() {
        let agent = create_test_agent(128);

        let request = tonic::Request::new(insert::Request {
            vector: Some(object::Vector {
                id: "test-uuid".to_string(),
                vector: gen_vector(128, 1),
                timestamp: 0,
            }),
            config: None, // Missing config
        });

        let result = agent.insert(request).await;
        assert!(result.is_err());

        let status = result.unwrap_err();
        assert_eq!(status.code(), tonic::Code::InvalidArgument);
    }

    // ==================== Search Handler Tests ====================

    #[tokio::test]
    async fn test_search_handler_success() {
        let agent = create_test_agent(128);

        // Insert some vectors first
        for i in 0..5 {
            let request = tonic::Request::new(insert::Request {
                vector: Some(object::Vector {
                    id: format!("vec-{}", i),
                    vector: gen_vector(128, i),
                    timestamp: 0,
                }),
                config: Some(insert::Config::default()),
            });
            agent.insert(request).await.unwrap();
        }

        // Search
        let search_request = tonic::Request::new(search::Request {
            vector: gen_vector(128, 100),
            config: Some(search::Config {
                request_id: "req-1".to_string(),
                num: 3,
                radius: -1.0,
                epsilon: 0.1,
                timeout: 0,
                ingress_filters: None,
                egress_filters: None,
                min_num: 0,
                aggregation_algorithm: 0,
                ratio: None,
                nprobe: 0,
            }),
        });

        let result = agent.search(search_request).await;
        assert!(result.is_ok());

        let response = result.unwrap().into_inner();
        assert!(!response.results.is_empty());
        assert!(response.results.len() <= 3);
    }

    #[tokio::test]
    async fn test_search_handler_invalid_dimension() {
        let agent = create_test_agent(128);

        let request = tonic::Request::new(search::Request {
            vector: gen_vector(64, 1), // Wrong dimension
            config: Some(search::Config {
                request_id: "req-1".to_string(),
                num: 3,
                radius: -1.0,
                epsilon: 0.1,
                timeout: 0,
                ingress_filters: None,
                egress_filters: None,
                min_num: 0,
                aggregation_algorithm: 0,
                ratio: None,
                nprobe: 0,
            }),
        });

        let result = agent.search(request).await;
        assert!(result.is_err());

        let status = result.unwrap_err();
        assert_eq!(status.code(), tonic::Code::InvalidArgument);
    }

    #[tokio::test]
    async fn test_search_handler_empty_index() {
        let agent = create_test_agent(128);

        let request = tonic::Request::new(search::Request {
            vector: gen_vector(128, 1),
            config: Some(search::Config {
                request_id: "req-1".to_string(),
                num: 3,
                radius: -1.0,
                epsilon: 0.1,
                timeout: 0,
                ingress_filters: None,
                egress_filters: None,
                min_num: 0,
                aggregation_algorithm: 0,
                ratio: None,
                nprobe: 0,
            }),
        });

        // Mock always returns results, so this succeeds
        let result = agent.search(request).await;
        assert!(result.is_ok());
    }

    // ==================== Remove Handler Tests ====================

    #[tokio::test]
    async fn test_remove_handler_success() {
        let agent = create_test_agent(128);

        // Insert a vector first
        let insert_request = tonic::Request::new(insert::Request {
            vector: Some(object::Vector {
                id: "to-remove".to_string(),
                vector: gen_vector(128, 1),
                timestamp: 0,
            }),
            config: Some(insert::Config::default()),
        });
        agent.insert(insert_request).await.unwrap();

        // Remove
        let remove_request = tonic::Request::new(remove::Request {
            id: Some(object::Id {
                id: "to-remove".to_string(),
            }),
            config: Some(remove::Config {
                skip_strict_exist_check: false,
                timestamp: 0,
            }),
        });

        let result = agent.remove(remove_request).await;
        assert!(result.is_ok());

        let response = result.unwrap().into_inner();
        assert_eq!(response.uuid, "to-remove");
    }

    #[tokio::test]
    async fn test_remove_handler_not_found() {
        let agent = create_test_agent(128);

        let request = tonic::Request::new(remove::Request {
            id: Some(object::Id {
                id: "nonexistent".to_string(),
            }),
            config: Some(remove::Config::default()),
        });

        // Mock always succeeds
        let result = agent.remove(request).await;
        assert!(result.is_ok());
    }

    #[tokio::test]
    async fn test_remove_handler_empty_uuid() {
        let agent = create_test_agent(128);

        let request = tonic::Request::new(remove::Request {
            id: Some(object::Id {
                id: "".to_string(), // Empty UUID
            }),
            config: Some(remove::Config::default()),
        });

        let result = agent.remove(request).await;
        assert!(result.is_err());

        let status = result.unwrap_err();
        assert_eq!(status.code(), tonic::Code::InvalidArgument);
    }

    // ==================== Object Handler Tests ====================

    #[tokio::test]
    async fn test_get_object_handler_success() {
        let agent = create_test_agent(128);

        // Get object - Mock returns fixed values
        let get_request = tonic::Request::new(object::VectorRequest {
            id: Some(object::Id {
                id: "get-object-test".to_string(),
            }),
            filters: None,
        });

        let result = agent.get_object(get_request).await;
        assert!(result.is_ok());

        let response = result.unwrap().into_inner();
        assert_eq!(response.id, "get-object-test");
        assert_eq!(response.vector.len(), 128); // Mock returns vec![0.0; 128]
    }

    #[tokio::test]
    async fn test_get_object_handler_not_found() {
        let agent = create_test_agent(128);

        let request = tonic::Request::new(object::VectorRequest {
            id: Some(object::Id {
                id: "nonexistent".to_string(),
            }),
            filters: None,
        });

        // Mock always returns success
        let result = agent.get_object(request).await;
        assert!(result.is_ok());
    }

    #[tokio::test]
    async fn test_get_object_handler_empty_uuid() {
        let agent = create_test_agent(128);

        let request = tonic::Request::new(object::VectorRequest {
            id: Some(object::Id {
                id: "".to_string(), // Empty UUID
            }),
            filters: None,
        });

        let result = agent.get_object(request).await;
        assert!(result.is_err());

        let status = result.unwrap_err();
        assert_eq!(status.code(), tonic::Code::InvalidArgument);
    }

    // ==================== Multi-operation Tests ====================

    #[tokio::test]
    async fn test_multi_insert_handler_success() {
        use proto::vald::v1::insert_server::Insert;

        let agent = create_test_agent(128);

        let requests: Vec<insert::Request> = (0..5)
            .map(|i| insert::Request {
                vector: Some(object::Vector {
                    id: format!("multi-{}", i),
                    vector: gen_vector(128, i),
                    timestamp: 0,
                }),
                config: Some(insert::Config::default()),
            })
            .collect();

        let request = tonic::Request::new(insert::MultiRequest { requests });
        let result = agent.multi_insert(request).await;

        assert!(result.is_ok());
        let response = result.unwrap().into_inner();
        assert_eq!(response.locations.len(), 5);
    }

    #[tokio::test]
    async fn test_multi_search_handler_success() {
        use proto::vald::v1::search_server::Search;

        let agent = create_test_agent(128);

        // Insert vectors first
        for i in 0..10 {
            let request = tonic::Request::new(insert::Request {
                vector: Some(object::Vector {
                    id: format!("vec-{}", i),
                    vector: gen_vector(128, i),
                    timestamp: 0,
                }),
                config: Some(insert::Config::default()),
            });
            agent.insert(request).await.unwrap();
        }

        let requests: Vec<search::Request> = (0..3)
            .map(|i| search::Request {
                vector: gen_vector(128, i + 100),
                config: Some(search::Config {
                    request_id: format!("req-{}", i),
                    num: 2,
                    radius: -1.0,
                    epsilon: 0.1,
                    timeout: 0,
                    ingress_filters: None,
                    egress_filters: None,
                    min_num: 0,
                    aggregation_algorithm: 0,
                    ratio: None,
                    nprobe: 0,
                }),
            })
            .collect();

        let request = tonic::Request::new(search::MultiRequest { requests });
        let result = agent.multi_search(request).await;

        assert!(result.is_ok());
        let response = result.unwrap().into_inner();
        assert_eq!(response.responses.len(), 3);
    }

    // ==================== Parse Duration Tests ====================

    #[test]
    fn test_parse_duration_seconds() {
        assert_eq!(
            parse_duration_from_string("30s"),
            Some(Duration::from_secs(30))
        );
        assert_eq!(
            parse_duration_from_string("1s"),
            Some(Duration::from_secs(1))
        );
        assert_eq!(
            parse_duration_from_string("0s"),
            Some(Duration::from_secs(0))
        );
    }

    #[test]
    fn test_parse_duration_minutes() {
        assert_eq!(
            parse_duration_from_string("5m"),
            Some(Duration::from_secs(300))
        );
        assert_eq!(
            parse_duration_from_string("1m"),
            Some(Duration::from_secs(60))
        );
    }

    #[test]
    fn test_parse_duration_hours() {
        assert_eq!(
            parse_duration_from_string("1h"),
            Some(Duration::from_secs(3600))
        );
        assert_eq!(
            parse_duration_from_string("2h"),
            Some(Duration::from_secs(7200))
        );
    }

    #[test]
    fn test_parse_duration_invalid() {
        assert_eq!(parse_duration_from_string(""), None);
        assert_eq!(parse_duration_from_string("30"), None);
        assert_eq!(parse_duration_from_string("abc"), None);
        assert_eq!(parse_duration_from_string("s"), None);
    }

    // ==================== Update Handler Tests ====================

    #[tokio::test]
    async fn test_update_handler_success() {
        use proto::vald::v1::update_server::Update;

        let agent = create_test_agent(128);

        // Update the vector - Mock always succeeds
        let new_vector = gen_vector(128, 100);
        let update_request = tonic::Request::new(update::Request {
            vector: Some(object::Vector {
                id: "update-test".to_string(),
                vector: new_vector.clone(),
                timestamp: 0,
            }),
            config: Some(update::Config::default()),
        });

        let result = agent.update(update_request).await;
        assert!(result.is_ok());
        assert_eq!(result.unwrap().into_inner().uuid, "update-test");
    }

    #[tokio::test]
    async fn test_update_handler_not_found() {
        use proto::vald::v1::update_server::Update;

        let agent = create_test_agent(128);

        let request = tonic::Request::new(update::Request {
            vector: Some(object::Vector {
                id: "nonexistent".to_string(),
                vector: gen_vector(128, 1),
                timestamp: 0,
            }),
            config: Some(update::Config::default()),
        });

        // Mock always succeeds
        let result = agent.update(request).await;
        assert!(result.is_ok());
    }

    #[tokio::test]
    async fn test_update_handler_invalid_dimension() {
        use proto::vald::v1::update_server::Update;

        let agent = create_test_agent(128);

        // Insert first
        let insert_request = tonic::Request::new(insert::Request {
            vector: Some(object::Vector {
                id: "update-dim-test".to_string(),
                vector: gen_vector(128, 1),
                timestamp: 0,
            }),
            config: Some(insert::Config::default()),
        });
        agent.insert(insert_request).await.unwrap();

        // Try to update with wrong dimension
        let request = tonic::Request::new(update::Request {
            vector: Some(object::Vector {
                id: "update-dim-test".to_string(),
                vector: gen_vector(64, 1), // Wrong dimension
                timestamp: 0,
            }),
            config: Some(update::Config::default()),
        });

        let result = agent.update(request).await;
        assert!(result.is_err());
        assert_eq!(result.unwrap_err().code(), tonic::Code::InvalidArgument);
    }

    // ==================== Upsert Handler Tests ====================

    #[tokio::test]
    async fn test_upsert_handler_insert_new() {
        use proto::vald::v1::upsert_server::Upsert;

        let agent = create_test_agent(128);

        let vector = gen_vector(128, 1);
        let request = tonic::Request::new(upsert::Request {
            vector: Some(object::Vector {
                id: "upsert-new".to_string(),
                vector: vector.clone(),
                timestamp: 0,
            }),
            config: Some(upsert::Config::default()),
        });

        let result = agent.upsert(request).await;
        assert!(result.is_ok());
        assert_eq!(result.unwrap().into_inner().uuid, "upsert-new");
    }

    #[tokio::test]
    async fn test_upsert_handler_update_existing() {
        use proto::vald::v1::upsert_server::Upsert;

        let agent = create_test_agent(128);

        // Upsert (update) with new vector - Mock always reports exists=true
        let new_vector = gen_vector(128, 100);
        let request = tonic::Request::new(upsert::Request {
            vector: Some(object::Vector {
                id: "upsert-update".to_string(),
                vector: new_vector.clone(),
                timestamp: 0,
            }),
            config: Some(upsert::Config::default()),
        });

        let result = agent.upsert(request).await;
        assert!(result.is_ok());
        assert_eq!(result.unwrap().into_inner().uuid, "upsert-update");
    }

    // ==================== Exists Handler Tests ====================

    #[tokio::test]
    async fn test_exists_handler_found() {
        use proto::vald::v1::object_server::Object;

        let agent = create_test_agent(128);

        // Insert a vector
        let insert_request = tonic::Request::new(insert::Request {
            vector: Some(object::Vector {
                id: "exists-test".to_string(),
                vector: gen_vector(128, 1),
                timestamp: 0,
            }),
            config: Some(insert::Config::default()),
        });
        agent.insert(insert_request).await.unwrap();

        // Check exists
        let request = tonic::Request::new(object::Id {
            id: "exists-test".to_string(),
        });

        let result = agent.exists(request).await;
        assert!(result.is_ok());
        assert_eq!(result.unwrap().into_inner().id, "exists-test");
    }

    #[tokio::test]
    async fn test_exists_handler_not_found() {
        use proto::vald::v1::object_server::Object;

        let agent = create_test_agent(128);

        let request = tonic::Request::new(object::Id {
            id: "nonexistent".to_string(),
        });

        // Mock always returns exists=true
        let result = agent.exists(request).await;
        assert!(result.is_ok());
        assert_eq!(result.unwrap().into_inner().id, "nonexistent");
    }

    #[tokio::test]
    async fn test_exists_handler_empty_uuid() {
        use proto::vald::v1::object_server::Object;

        let agent = create_test_agent(128);

        let request = tonic::Request::new(object::Id { id: "".to_string() });

        let result = agent.exists(request).await;
        assert!(result.is_err());
        assert_eq!(result.unwrap_err().code(), tonic::Code::InvalidArgument);
    }

    // ==================== Index Handler Tests ====================

    #[tokio::test]
    async fn test_create_index_handler() {
        use proto::core::v1::agent_server::Agent as AgentServer;
        use proto::payload::v1::control;

        let agent = create_test_agent(128);

        let request = tonic::Request::new(control::CreateIndexRequest { pool_size: 10 });
        let result = agent.create_index(request).await;

        assert!(result.is_ok());
    }

    #[tokio::test]
    async fn test_save_index_handler() {
        use proto::core::v1::agent_server::Agent as AgentServer;
        use proto::payload::v1::Empty;

        let agent = create_test_agent(128);

        let request = tonic::Request::new(Empty {});
        let result = agent.save_index(request).await;

        assert!(result.is_ok());
    }

    #[tokio::test]
    async fn test_create_and_save_index_handler() {
        use proto::core::v1::agent_server::Agent as AgentServer;
        use proto::payload::v1::control;

        let agent = create_test_agent(128);

        let request = tonic::Request::new(control::CreateIndexRequest { pool_size: 10 });
        let result = agent.create_and_save_index(request).await;

        assert!(result.is_ok());
    }

    #[tokio::test]
    async fn test_index_info_handler() {
        use proto::payload::v1::Empty;
        use proto::vald::v1::index_server::Index;

        let agent = create_test_agent(128);

        let request = tonic::Request::new(Empty {});
        let result = agent.index_info(request).await;

        assert!(result.is_ok());
        let response = result.unwrap().into_inner();
        assert!(!response.indexing);
        assert!(!response.saving);
    }

    #[tokio::test]
    async fn test_index_detail_handler() {
        use proto::payload::v1::Empty;
        use proto::vald::v1::index_server::Index;

        let agent = create_test_agent(128);

        let request = tonic::Request::new(Empty {});
        let result = agent.index_detail(request).await;

        assert!(result.is_ok());
        let response = result.unwrap().into_inner();
        assert_eq!(response.replica, 1);
        assert_eq!(response.live_agents, 1);
        assert!(response.counts.contains_key("test-agent"));
    }

    #[tokio::test]
    async fn test_index_statistics_handler() {
        use proto::payload::v1::Empty;
        use proto::vald::v1::index_server::Index;

        let agent = create_test_agent(128);

        let request = tonic::Request::new(Empty {});
        let result = agent.index_statistics(request).await;

        assert!(result.is_ok());
    }

    #[tokio::test]
    async fn test_index_property_handler() {
        use proto::payload::v1::Empty;
        use proto::vald::v1::index_server::Index;

        let agent = create_test_agent(128);

        let request = tonic::Request::new(Empty {});
        let result = agent.index_property(request).await;

        assert!(result.is_ok());
        let response = result.unwrap().into_inner();
        assert!(response.details.contains_key("test-agent"));
    }

    // ==================== Flush Handler Tests ====================

    #[tokio::test]
    async fn test_flush_handler() {
        use proto::payload::v1::flush;
        use proto::vald::v1::flush_server::Flush;

        let agent = create_test_agent(128);

        let request = tonic::Request::new(flush::Request {});
        let result = agent.flush(request).await;

        assert!(result.is_ok());
        let response = result.unwrap().into_inner();
        assert!(!response.indexing);
        assert!(!response.saving);
    }

    // ==================== Search By ID Handler Tests ====================

    #[tokio::test]
    async fn test_search_by_id_handler_success() {
        use proto::vald::v1::search_server::Search;

        let agent = create_test_agent(128);

        // Insert vectors first
        for i in 0..10 {
            let request = tonic::Request::new(insert::Request {
                vector: Some(object::Vector {
                    id: format!("search-id-{}", i),
                    vector: gen_vector(128, i),
                    timestamp: 0,
                }),
                config: Some(insert::Config::default()),
            });
            agent.insert(request).await.unwrap();
        }

        let request = tonic::Request::new(search::IdRequest {
            id: "search-id-0".to_string(),
            config: Some(search::Config {
                request_id: "req-1".to_string(),
                num: 5,
                radius: -1.0,
                epsilon: 0.1,
                timeout: 0,
                ingress_filters: None,
                egress_filters: None,
                min_num: 0,
                aggregation_algorithm: 0,
                ratio: None,
                nprobe: 0,
            }),
        });

        let result = agent.search_by_id(request).await;
        assert!(result.is_ok());
    }

    #[tokio::test]
    async fn test_search_by_id_handler_empty_uuid() {
        use proto::vald::v1::search_server::Search;

        let agent = create_test_agent(128);

        let request = tonic::Request::new(search::IdRequest {
            id: "".to_string(),
            config: Some(search::Config::default()),
        });

        let result = agent.search_by_id(request).await;
        assert!(result.is_err());
        assert_eq!(result.unwrap_err().code(), tonic::Code::InvalidArgument);
    }

    #[tokio::test]
    async fn test_search_by_id_handler_not_found() {
        use proto::vald::v1::search_server::Search;

        let agent = create_test_agent(128);

        let request = tonic::Request::new(search::IdRequest {
            id: "nonexistent".to_string(),
            config: Some(search::Config {
                request_id: "req-1".to_string(),
                num: 5,
                radius: -1.0,
                epsilon: 0.1,
                timeout: 0,
                ingress_filters: None,
                egress_filters: None,
                min_num: 0,
                aggregation_algorithm: 0,
                ratio: None,
                nprobe: 0,
            }),
        });

        // Mock always returns results
        let result = agent.search_by_id(request).await;
        assert!(result.is_ok());
    }

    // ==================== Linear Search Handler Tests ====================

    #[tokio::test]
    async fn test_linear_search_handler_unsupported() {
        use proto::vald::v1::search_server::Search;

        let agent = create_test_agent(128);

        let request = tonic::Request::new(search::Request {
            vector: gen_vector(128, 1),
            config: Some(search::Config {
                request_id: "req-1".to_string(),
                num: 5,
                radius: -1.0,
                epsilon: 0.1,
                timeout: 0,
                ingress_filters: None,
                egress_filters: None,
                min_num: 0,
                aggregation_algorithm: 0,
                ratio: None,
                nprobe: 0,
            }),
        });

        let result = agent.linear_search(request).await;
        // MockANNService returns Unsupported error for linear_search
        assert!(result.is_err());
        assert_eq!(result.unwrap_err().code(), tonic::Code::Unimplemented);
    }

    #[tokio::test]
    async fn test_linear_search_by_id_handler_unsupported() {
        use proto::vald::v1::search_server::Search;

        let agent = create_test_agent(128);

        let request = tonic::Request::new(search::IdRequest {
            id: "test-uuid".to_string(),
            config: Some(search::Config {
                request_id: "req-1".to_string(),
                num: 5,
                radius: -1.0,
                epsilon: 0.1,
                timeout: 0,
                ingress_filters: None,
                egress_filters: None,
                min_num: 0,
                aggregation_algorithm: 0,
                ratio: None,
                nprobe: 0,
            }),
        });

        let result = agent.linear_search_by_id(request).await;
        // MockANNService returns Unsupported error for linear_search_by_id
        assert!(result.is_err());
        assert_eq!(result.unwrap_err().code(), tonic::Code::Unimplemented);
    }

    // ==================== Multi Remove Handler Tests ====================

    #[tokio::test]
    async fn test_multi_remove_handler_success() {
        use proto::vald::v1::remove_server::Remove;

        let agent = create_test_agent(128);

        // Insert vectors first
        for i in 0..5 {
            let request = tonic::Request::new(insert::Request {
                vector: Some(object::Vector {
                    id: format!("multi-remove-{}", i),
                    vector: gen_vector(128, i),
                    timestamp: 0,
                }),
                config: Some(insert::Config::default()),
            });
            agent.insert(request).await.unwrap();
        }

        let requests: Vec<remove::Request> = (0..5)
            .map(|i| remove::Request {
                id: Some(object::Id {
                    id: format!("multi-remove-{}", i),
                }),
                config: None,
            })
            .collect();

        let request = tonic::Request::new(remove::MultiRequest { requests });
        let result = agent.multi_remove(request).await;

        assert!(result.is_ok());
        let response = result.unwrap().into_inner();
        assert_eq!(response.locations.len(), 5);
    }

    // ==================== Multi Update Handler Tests ====================

    #[tokio::test]
    async fn test_multi_update_handler_success() {
        use proto::vald::v1::update_server::Update;

        let agent = create_test_agent(128);

        // Insert vectors first
        for i in 0..3 {
            let request = tonic::Request::new(insert::Request {
                vector: Some(object::Vector {
                    id: format!("multi-update-{}", i),
                    vector: gen_vector(128, i),
                    timestamp: 0,
                }),
                config: Some(insert::Config::default()),
            });
            agent.insert(request).await.unwrap();
        }

        let requests: Vec<update::Request> = (0..3)
            .map(|i| update::Request {
                vector: Some(object::Vector {
                    id: format!("multi-update-{}", i),
                    vector: gen_vector(128, i + 100),
                    timestamp: 0,
                }),
                config: Some(update::Config::default()),
            })
            .collect();

        let request = tonic::Request::new(update::MultiRequest { requests });
        let result = agent.multi_update(request).await;

        assert!(result.is_ok());
        let response = result.unwrap().into_inner();
        assert_eq!(response.locations.len(), 3);
    }

    // ==================== Multi Upsert Handler Tests ====================

    #[tokio::test]
    async fn test_multi_upsert_handler_success() {
        use proto::vald::v1::upsert_server::Upsert;

        let agent = create_test_agent(128);

        let requests: Vec<upsert::Request> = (0..5)
            .map(|i| upsert::Request {
                vector: Some(object::Vector {
                    id: format!("multi-upsert-{}", i),
                    vector: gen_vector(128, i),
                    timestamp: 0,
                }),
                config: Some(upsert::Config::default()),
            })
            .collect();

        let request = tonic::Request::new(upsert::MultiRequest { requests });
        let result = agent.multi_upsert(request).await;

        assert!(result.is_ok());
        let response = result.unwrap().into_inner();
        assert_eq!(response.locations.len(), 5);
    }

    // ==================== Graceful Shutdown Tests ====================

    /// Mock ANN service with shutdown tracking for testing graceful shutdown
    struct MockShutdownService {
        dimension: usize,
        close_called: std::sync::atomic::AtomicBool,
        create_index_count: std::sync::atomic::AtomicU32,
        save_index_count: std::sync::atomic::AtomicU32,
    }

    impl MockShutdownService {
        fn new(dimension: usize) -> Self {
            Self {
                dimension,
                close_called: std::sync::atomic::AtomicBool::new(false),
                create_index_count: std::sync::atomic::AtomicU32::new(0),
                save_index_count: std::sync::atomic::AtomicU32::new(0),
            }
        }

        fn is_close_called(&self) -> bool {
            self.close_called.load(std::sync::atomic::Ordering::SeqCst)
        }

        fn get_create_index_count(&self) -> u32 {
            self.create_index_count
                .load(std::sync::atomic::Ordering::SeqCst)
        }

        fn get_save_index_count(&self) -> u32 {
            self.save_index_count
                .load(std::sync::atomic::Ordering::SeqCst)
        }
    }

    impl ANN for MockShutdownService {
        fn get_dimension_size(&self) -> usize {
            self.dimension
        }

        fn search(
            &self,
            _v: Vec<f32>,
            _n: u32,
            _e: f32,
            _r: f32,
        ) -> impl std::future::Future<Output = Result<search::Response, Error>> + Send {
            async { Ok(search::Response::default()) }
        }
        fn search_by_id(
            &self,
            _u: String,
            _n: u32,
            _e: f32,
            _r: f32,
        ) -> impl std::future::Future<Output = Result<search::Response, Error>> + Send {
            async { Ok(search::Response::default()) }
        }
        fn linear_search(
            &self,
            _v: Vec<f32>,
            _n: u32,
        ) -> impl std::future::Future<Output = Result<search::Response, Error>> + Send {
            async { Ok(search::Response::default()) }
        }
        fn linear_search_by_id(
            &self,
            _u: String,
            _n: u32,
        ) -> impl std::future::Future<Output = Result<search::Response, Error>> + Send {
            async { Ok(search::Response::default()) }
        }
        fn insert(
            &mut self,
            _u: String,
            _v: Vec<f32>,
        ) -> impl std::future::Future<Output = Result<(), Error>> + Send {
            async { Ok(()) }
        }
        fn insert_with_time(
            &mut self,
            _u: String,
            _v: Vec<f32>,
            _t: i64,
        ) -> impl std::future::Future<Output = Result<(), Error>> + Send {
            async { Ok(()) }
        }
        fn insert_multiple(
            &mut self,
            _vs: HashMap<String, Vec<f32>>,
        ) -> impl std::future::Future<Output = Result<(), Error>> + Send {
            async { Ok(()) }
        }
        fn insert_multiple_with_time(
            &mut self,
            _vs: HashMap<String, Vec<f32>>,
            _t: i64,
        ) -> impl std::future::Future<Output = Result<(), Error>> + Send {
            async { Ok(()) }
        }
        fn update(
            &mut self,
            _u: String,
            _v: Vec<f32>,
        ) -> impl std::future::Future<Output = Result<(), Error>> + Send {
            async { Ok(()) }
        }
        fn update_with_time(
            &mut self,
            _u: String,
            _v: Vec<f32>,
            _t: i64,
        ) -> impl std::future::Future<Output = Result<(), Error>> + Send {
            async { Ok(()) }
        }
        fn update_multiple(
            &mut self,
            _vs: HashMap<String, Vec<f32>>,
        ) -> impl std::future::Future<Output = Result<(), Error>> + Send {
            async { Ok(()) }
        }
        fn update_multiple_with_time(
            &mut self,
            _vs: HashMap<String, Vec<f32>>,
            _t: i64,
        ) -> impl std::future::Future<Output = Result<(), Error>> + Send {
            async { Ok(()) }
        }
        fn update_timestamp(
            &mut self,
            _u: String,
            _t: i64,
            _f: bool,
        ) -> impl std::future::Future<Output = Result<(), Error>> + Send {
            async { Ok(()) }
        }
        fn remove(
            &mut self,
            _u: String,
        ) -> impl std::future::Future<Output = Result<(), Error>> + Send {
            async { Ok(()) }
        }
        fn remove_with_time(
            &mut self,
            _u: String,
            _t: i64,
        ) -> impl std::future::Future<Output = Result<(), Error>> + Send {
            async { Ok(()) }
        }
        fn remove_multiple(
            &mut self,
            _us: Vec<String>,
        ) -> impl std::future::Future<Output = Result<(), Error>> + Send {
            async { Ok(()) }
        }
        fn remove_multiple_with_time(
            &mut self,
            _us: Vec<String>,
            _t: i64,
        ) -> impl std::future::Future<Output = Result<(), Error>> + Send {
            async { Ok(()) }
        }

        fn get_object(
            &self,
            _uuid: String,
        ) -> impl std::future::Future<Output = Result<(Vec<f32>, i64), Error>> + Send {
            let dim = self.dimension;
            async move { Ok((vec![0.0; dim], 12345)) }
        }

        fn exists(&self, _uuid: String) -> impl std::future::Future<Output = (usize, bool)> + Send {
            async { (1, true) }
        }
        fn uuids(&self) -> impl std::future::Future<Output = Vec<String>> + Send {
            async { vec![] }
        }
        fn list_object_func<F: FnMut(String, Vec<f32>, i64) -> bool + Send>(
            &self,
            _f: F,
        ) -> impl std::future::Future<Output = ()> + Send {
            async {}
        }

        fn create_index(&mut self) -> impl std::future::Future<Output = Result<(), Error>> + Send {
            self.create_index_count
                .fetch_add(1, std::sync::atomic::Ordering::SeqCst);
            async { Ok(()) }
        }
        fn save_index(&mut self) -> impl std::future::Future<Output = Result<(), Error>> + Send {
            self.save_index_count
                .fetch_add(1, std::sync::atomic::Ordering::SeqCst);
            async { Ok(()) }
        }
        fn create_and_save_index(
            &mut self,
        ) -> impl std::future::Future<Output = Result<(), Error>> + Send {
            self.create_index_count
                .fetch_add(1, std::sync::atomic::Ordering::SeqCst);
            self.save_index_count
                .fetch_add(1, std::sync::atomic::Ordering::SeqCst);
            async { Ok(()) }
        }
        fn regenerate_indexes(
            &mut self,
        ) -> impl std::future::Future<Output = Result<(), Error>> + Send {
            async { Ok(()) }
        }
        fn len(&self) -> u32 {
            100
        }
        fn insert_vqueue_buffer_len(&self) -> u32 {
            0
        }
        fn delete_vqueue_buffer_len(&self) -> u32 {
            0
        }
        fn is_indexing(&self) -> bool {
            false
        }
        fn is_flushing(&self) -> bool {
            false
        }
        fn is_saving(&self) -> bool {
            false
        }
        fn number_of_create_index_executions(&self) -> u64 {
            0
        }
        fn broken_index_count(&self) -> u64 {
            0
        }
        fn is_statistics_enabled(&self) -> bool {
            false
        }
        fn index_statistics(&self) -> Result<info::index::Statistics, Error> {
            Ok(info::index::Statistics::default())
        }
        fn index_property(&self) -> Result<info::index::Property, Error> {
            Ok(info::index::Property::default())
        }

        fn close(&mut self) -> impl std::future::Future<Output = Result<(), Error>> + Send {
            self.close_called
                .store(true, std::sync::atomic::Ordering::SeqCst);
            async { Ok(()) }
        }
    }

    #[tokio::test]
    async fn test_agent_shutdown_without_daemon() {
        // Test shutdown when daemon is not started
        let agent = Agent::new(
            MockShutdownService::new(128),
            "test",
            "127.0.0.1",
            "vald.v1",
            "vald-agent",
            10,
        );

        // Shutdown should succeed even without daemon
        let result = agent.shutdown().await;
        assert!(result.is_ok(), "Shutdown should succeed without daemon");

        // Verify close was called
        let service = agent.service();
        let svc = service.read().await;
        assert!(
            svc.is_close_called(),
            "close() should be called during shutdown"
        );
    }

    #[tokio::test]
    async fn test_agent_shutdown_with_daemon() {
        use crate::service::{DaemonConfig, start_daemon};

        let service = MockShutdownService::new(128);
        let service_arc = Arc::new(RwLock::new(service));

        // Create daemon manually
        let daemon_config = DaemonConfig {
            auto_index_check_duration: std::time::Duration::from_secs(3600),
            auto_save_index_duration: std::time::Duration::from_secs(3600),
            auto_index_limit: std::time::Duration::from_secs(3600),
            auto_index_length: 1000,
            pool_size: 100,
            initial_delay: std::time::Duration::ZERO,
            enable_proactive_gc: false,
        };

        let (handle, error_rx) = start_daemon(service_arc.clone(), daemon_config).await;

        // Create agent with daemon
        let agent = Agent {
            s: service_arc.clone(),
            name: "test".to_string(),
            ip: "127.0.0.1".to_string(),
            resource_type: "vald.v1".to_string(),
            api_name: "vald-agent".to_string(),
            stream_concurrency: 10,
            daemon_handle: Some(handle),
            error_rx: Some(error_rx),
        };

        // Let daemon start
        tokio::time::sleep(std::time::Duration::from_millis(50)).await;

        // Shutdown should complete and call close
        let start = std::time::Instant::now();
        let result = agent.shutdown().await;
        let elapsed = start.elapsed();

        assert!(result.is_ok(), "Shutdown should succeed");
        assert!(
            elapsed < std::time::Duration::from_secs(1),
            "Shutdown should be fast"
        );

        // Verify close was called
        let svc = service_arc.read().await;
        assert!(
            svc.is_close_called(),
            "close() should be called during shutdown"
        );

        // Verify final index was created (daemon shutdown creates index)
        assert!(
            svc.get_create_index_count() >= 1,
            "create_index should be called on shutdown"
        );
    }

    #[tokio::test]
    async fn test_agent_stop_signals_daemon() {
        use crate::service::{DaemonConfig, start_daemon};

        let service = MockShutdownService::new(128);
        let service_arc = Arc::new(RwLock::new(service));

        let daemon_config = DaemonConfig::default();
        let (handle, error_rx) = start_daemon(service_arc.clone(), daemon_config).await;

        let agent = Agent {
            s: service_arc.clone(),
            name: "test".to_string(),
            ip: "127.0.0.1".to_string(),
            resource_type: "vald.v1".to_string(),
            api_name: "vald-agent".to_string(),
            stream_concurrency: 10,
            daemon_handle: Some(handle.clone()),
            error_rx: Some(error_rx),
        };

        // Verify daemon is not cancelled yet
        assert!(
            !handle.is_cancelled(),
            "Daemon should not be cancelled initially"
        );

        // Stop should signal daemon
        agent.stop();

        assert!(
            handle.is_cancelled(),
            "Daemon should be cancelled after stop()"
        );
    }

    #[tokio::test]
    async fn test_agent_shutdown_is_idempotent() {
        let agent = Agent::new(
            MockShutdownService::new(128),
            "test",
            "127.0.0.1",
            "vald.v1",
            "vald-agent",
            10,
        );

        // First shutdown
        let result1 = agent.shutdown().await;
        assert!(result1.is_ok());

        // Second shutdown should also succeed (idempotent)
        let result2 = agent.shutdown().await;
        assert!(result2.is_ok());
    }
}
