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
pub mod flush;
pub mod index;
pub mod insert;
pub mod object;
pub mod remove;
pub mod search;
pub mod update;
pub mod upsert;

use std::sync::Arc;
use std::time::Duration;
use tokio::sync::RwLock;
use config::Config;
use proto::{
    core::v1::agent_server,
    vald::v1::{
        flush_server, index_server, insert_server, object_server, remove_server, search_server, update_server, upsert_server
    }
};
use crate::middleware;

pub struct Agent<S: algorithm::ANN + 'static> {
    s: Arc<RwLock<S>>,
    name: String,
    ip: String,
    resource_type: String,
    api_name: String,
    stream_concurrency: usize,
}

impl<S: algorithm::ANN + 'static> Agent<S> {
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
            stream_concurrency: stream_concurrency,
        }
    }

    /// Starts the gRPC server with all registered services.
    pub async fn serve_grpc(self, settings: Config) -> Result<(), Box<dyn std::error::Error>> {
        let addr = "0.0.0.0:8081".parse()?;
        let mut grpc_key = String::new();
        for i in 0..settings.get_array("server_config.servers")?.len() {
            let name = settings.get::<String>(format!("server_config.servers[{i}].name").as_str())?;
            match name.as_str() {
                "grpc" => {
                    grpc_key = format!("server_config.servers[{i}]");
                }
                _ => {}
            }
        }

        let mut builder = tonic::transport::Server::builder();
        if let Some(duration) = parse_duration_from_string(
            settings
                .get::<String>(format!("{grpc_key}.grpc.keepalive.max_conn_age").as_str())?
                .as_str(),
        ) {
            builder = builder.max_connection_age(duration);
        }
        if let Some(duration) = parse_duration_from_string(
            settings
                .get::<String>(format!("{grpc_key}.grpc.connection_timeout").as_str())?
                .as_str(),
        ) {
            builder = builder.timeout(duration);
        }

        let mut accessloginterceptor: Option<()> = None;
        let mut metricinterceptor: Option<()> = None;
        for i in 0..settings
            .get_array(format!("{grpc_key}.grpc.interceptors").as_str())?
            .len()
        {
            let name = settings.get::<String>(format!("{grpc_key}.grpc.interceptors[{i}]").as_str())?;
            match name.to_lowercase().as_str() {
                "accessloginterceptor" | "accesslog" => accessloginterceptor = Some(()),
                "metricinterceptor" | "metric" => metricinterceptor = Some(()),
                _ => {}
            }
        }

        let layer = tower::ServiceBuilder::new()
            .option_layer(accessloginterceptor.map(|_| middleware::AccessLogMiddlewareLayer::default()))
            .option_layer(metricinterceptor.map(|_| middleware::MetricMiddlewareLayer::default()))
            .into_inner();

        let max_recv_size = settings.get::<usize>(format!("{grpc_key}.grpc.max_receive_message_size").as_str())?;
        let max_send_size = settings.get::<usize>(format!("{grpc_key}.grpc.max_send_message_size").as_str())?;

        builder
            .initial_stream_window_size(
                settings.get::<u32>(format!("{grpc_key}.grpc.initial_window_size").as_str())?,
            )
            .initial_connection_window_size(
                settings.get::<u32>(format!("{grpc_key}.grpc.initial_conn_window_size").as_str())?,
            )
            .http2_keepalive_interval(parse_duration_from_string(
                settings
                    .get::<String>(format!("{grpc_key}.grpc.keepalive.time").as_str())?
                    .as_str(),
            ))
            .http2_keepalive_timeout(parse_duration_from_string(
                settings
                    .get::<String>(format!("{grpc_key}.grpc.keepalive.timeout").as_str())?
                    .as_str(),
            ))
            .http2_max_header_list_size(
                settings.get::<u32>(format!("{grpc_key}.grpc.max_header_list_size").as_str())?,
            )
            .max_concurrent_streams(
                settings.get::<u32>(format!("{grpc_key}.grpc.max_concurrent_streams").as_str())?,
            )
            .layer(layer)
            .add_service(
                agent_server::AgentServer::new(self.clone())
                    .max_decoding_message_size(max_recv_size)
                    .max_encoding_message_size(max_send_size)
            )
            .add_service(
                search_server::SearchServer::new(self.clone())
                    .max_decoding_message_size(max_recv_size)
                    .max_encoding_message_size(max_send_size)
            )
            .add_service(
                insert_server::InsertServer::new(self.clone())
                    .max_decoding_message_size(max_recv_size)
                    .max_encoding_message_size(max_send_size)
            )
            .add_service(
                update_server::UpdateServer::new(self.clone())
                    .max_decoding_message_size(max_recv_size)
                    .max_encoding_message_size(max_send_size)
            )
            .add_service(
                upsert_server::UpsertServer::new(self.clone())
                    .max_decoding_message_size(max_recv_size)
                    .max_encoding_message_size(max_send_size)
            )
            .add_service(
                remove_server::RemoveServer::new(self.clone())
                    .max_decoding_message_size(max_recv_size)
                    .max_encoding_message_size(max_send_size)
            )
            .add_service(
                object_server::ObjectServer::new(self.clone())
                    .max_decoding_message_size(max_recv_size)
                    .max_encoding_message_size(max_send_size)
            )
            .add_service(
                index_server::IndexServer::new(self.clone())
                    .max_decoding_message_size(max_recv_size)
                    .max_encoding_message_size(max_send_size)
            )
            .add_service(
                flush_server::FlushServer::new(self.clone())
                    .max_decoding_message_size(max_recv_size)
                    .max_encoding_message_size(max_send_size)
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
        }
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
