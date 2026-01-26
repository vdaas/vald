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

mod handler;
mod middleware;
mod service;

use config::Config;
use handler::Agent;
use service::QBGService;

async fn serve(settings: Config) -> Result<(), Box<dyn std::error::Error>> {
    let _logger =
        flexi_logger::Logger::try_with_str(settings.get::<String>("logging.level")?)?.start()?;
    let service = match settings.get_string("service.type")?.as_str() {
        "qbg" => QBGService::new(settings.clone()).await,
        _ => panic!("unsupported algorithm service"),
    };
    let agent = Agent::new(
        service,
        "agent-qbg",
        "127.0.0.1",
        "vald/internal/core/algorithm",
        "vald-agent",
        10,
    );

    agent.serve_grpc(settings).await
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let settings = Config::builder()
        .add_source(config::File::with_name("/etc/server/config.yaml"))
        .build()
        .unwrap();
    
    serve(settings).await
}

#[cfg(test)]
mod tests {
    use super::*;

    /// Helper function to create test config
    fn create_test_config() -> Config {
        let config_str = r#"
logging:
  level: "info"
service:
  type: "qbg"
  dimension: 128
  creation_edge_size: 10
  search_edge_size: 40
  object_type: "Float"
  distance_type: "L2"
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
        Config::builder()
            .add_source(config::File::from_str(config_str, config::FileFormat::Yaml))
            .build()
            .unwrap()
    }

    #[test]
    fn test_config_parsing() {
        let config = create_test_config();
        
        assert_eq!(config.get_string("logging.level").unwrap(), "info");
        assert_eq!(config.get_string("service.type").unwrap(), "qbg");
        assert_eq!(config.get::<usize>("service.dimension").unwrap(), 128);
    }

    #[test]
    fn test_config_grpc_settings() {
        let config = create_test_config();
        
        let servers = config.get_array("server_config.servers").unwrap();
        assert_eq!(servers.len(), 1);
        
        let grpc_name = config.get_string("server_config.servers[0].name").unwrap();
        assert_eq!(grpc_name, "grpc");
        
        let max_recv = config.get::<usize>("server_config.servers[0].grpc.max_receive_message_size").unwrap();
        assert_eq!(max_recv, 4194304);
    }

    #[test]
    fn test_unsupported_service_type() {
        let config_str = r#"
logging:
  level: "info"
service:
  type: "unsupported"
"#;
        let config = Config::builder()
            .add_source(config::File::from_str(config_str, config::FileFormat::Yaml))
            .build()
            .unwrap();
        
        assert_eq!(config.get_string("service.type").unwrap(), "unsupported");
    }
}
