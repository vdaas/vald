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

pub mod config;
pub mod handler;
pub mod metrics;
pub mod middleware;
pub mod service;

use crate::config::AgentConfig;
use handler::Agent;
use observability::{TracingConfig, init_tracing, shutdown_tracing};
use service::QBGService;
use tracing::{error, info};

/// Starts the agent service with the given configuration.
pub async fn serve(config: AgentConfig) -> Result<(), Box<dyn std::error::Error>> {
    // Initialize tracing
    let tracing_config = TracingConfig::new()
        .enable_stdout(true)
        .enable_json(config.logging.json)
        .enable_otel(config.observability.tracer.enabled)
        .level(&config.logging.level)
        .service_name("vald-agent");

    // Build OpenTelemetry config if enabled
    let otel_config = if config.observability.enabled {
        Some(build_otel_config(&config))
    } else {
        None
    };

    let tracer_provider =
        init_tracing(&tracing_config, otel_config.as_ref()).expect("failed to initialize tracing");

    info!("starting vald-agent");

    let service = match config.service.type_.as_str() {
        "qbg" => QBGService::new(&config.qbg).await,
        _ => panic!("unsupported algorithm service"),
    };
    let mut agent = Agent::new(
        service,
        "agent-qbg",
        "127.0.0.1",
        "vald/internal/core/algorithm",
        "vald-agent",
        10,
    );

    // Start the daemon for automatic indexing and saving
    agent.start(&config).await;

    // Start health servers
    let health_servers = vec![
        &config.server_config.healths.liveness,
        &config.server_config.healths.readiness,
        &config.server_config.healths.startup,
    ];

    let mut bind_addrs = std::collections::HashSet::new();
    for s in health_servers {
        if s.enabled {
            let host = if s.host.is_empty() {
                "0.0.0.0"
            } else {
                &s.host
            };
            bind_addrs.insert(format!("{}:{}", host, s.port));
        }
    }

    for addr in bind_addrs {
        info!("Starting health server at {}", addr);
        let addr_clone = addr.clone();
        tokio::spawn(async move {
            match tokio::net::TcpListener::bind(&addr_clone).await {
                Ok(listener) => {
                    if let Err(e) = axum::serve(listener, handler::health::router()).await {
                        error!("Health server error on {}: {}", addr_clone, e);
                    }
                }
                Err(e) => {
                    error!("Failed to bind health server on {}: {}", addr_clone, e);
                }
            }
        });
    }

    // Register NGT metrics if metering is enabled
    if config.observability.enabled && config.observability.meter.enabled {
        if let Err(e) = metrics::register_metrics(agent.service()) {
            error!("failed to register metrics: {}", e);
        } else {
            info!("NGT metrics registered successfully");
        }
    }

    // Setup graceful shutdown
    let shutdown_agent = agent.clone();
    tokio::spawn(async move {
        match tokio::signal::ctrl_c().await {
            Ok(()) => {
                info!("Received shutdown signal, stopping daemon...");
                shutdown_agent.stop();
            }
            Err(e) => {
                error!("Failed to listen for shutdown signal: {}", e);
            }
        }
    });

    // Serve gRPC (blocks until server stops)
    let result = agent.serve_grpc(config).await;

    // Shutdown tracing
    if let Err(e) = shutdown_tracing(tracer_provider) {
        error!("failed to shutdown tracing: {}", e);
    }

    result
}

fn build_otel_config(config: &AgentConfig) -> observability::Config {
    use std::time::Duration;

    let endpoint = &config.observability.endpoint;
    let service_name = &config.observability.service_name;

    observability::Config::new()
        .enabled(config.observability.enabled)
        .endpoint(endpoint)
        .attribute(observability::observability::SERVICE_NAME, service_name)
        .tracer(observability::config::Tracer::new().enabled(config.observability.tracer.enabled))
        .meter(
            observability::config::Meter::new()
                .enabled(config.observability.meter.enabled)
                .export_duration(Duration::from_secs(
                    config.observability.meter.export_duration_secs,
                ))
                .export_timeout_duration(Duration::from_secs(
                    config.observability.meter.export_timeout_secs,
                )),
        )
}

#[cfg(test)]
mod tests {
    use super::*;

    /// Helper function to create test config
    fn create_test_config() -> AgentConfig {
        let config_str = r#"
logging:
  level: "info"
service:
  type: "qbg"
qbg:
  dimension: 128
  index_path: "/tmp/test_qbg_index"
server_config:
  servers:
    - name: grpc
      host: 0.0.0.0
      port: 8081
      grpc:
        max_receive_message_size: 4194304
        max_send_message_size: 4194304
        initial_window_size: 65535
        initial_conn_window_size: 65535
        max_header_list_size: 8192
        max_concurrent_streams: 100
        connection_timeout: 30s
        keepalive:
          max_conn_age: 300s
          time: 60s
          timeout: 20s
        interceptors:
          - accesslog
          - metric
"#;
        use ::config::FileFormat;
        let settings = ::config::Config::builder()
            .add_source(::config::File::from_str(config_str, FileFormat::Yaml))
            .build()
            .unwrap();

        let mut config: AgentConfig = settings.try_deserialize().unwrap();
        // Since deserialization might use defaults for missing fields, and `healths` might not be in the YAML,
        // it should be handled by `#[serde(default)]` in `config.rs`.
        // However, if we manually constructed AgentConfig in any test (which we didn't in this file), we'd need to fix it.
        // The `create_test_config` function uses `try_deserialize`, which respects `#[serde(default)]`.
        // So no manual change needed for `create_test_config` return value if `config.rs` has defaults.
        // But checking `config.rs`, `ServerConfig` derives `Default`.
        config
    }

    #[test]
    fn test_config_parsing() {
        let config = create_test_config();

        assert_eq!(config.logging.level, "info");
        assert_eq!(config.service.type_, "qbg");
        assert_eq!(config.qbg.dimension, 128);
    }

    #[test]
    fn test_config_grpc_settings() {
        let config = create_test_config();

        assert_eq!(config.server_config.servers.len(), 1);

        let server = &config.server_config.servers[0];
        assert_eq!(server.name, "grpc");
        assert_eq!(server.grpc.max_receive_message_size, 4194304);
    }

    #[test]
    fn test_unsupported_service_type() {
        let config_str = r#"
logging:
  level: "info"
service:
  type: "unsupported"
qbg:
  dimension: 128
  index_path: "/tmp/index"
"#;
        use ::config::FileFormat;
        let settings = ::config::Config::builder()
            .add_source(::config::File::from_str(config_str, FileFormat::Yaml))
            .build()
            .unwrap();

        let config: AgentConfig = settings.try_deserialize().unwrap();

        assert_eq!(config.service.type_, "unsupported");
    }
}
