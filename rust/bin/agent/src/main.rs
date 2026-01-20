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

use algorithm::{ANN, Error, MultiError};
use anyhow::Result;
use chrono::{Local, Timelike};
use config::Config;
use proto::{
    core::v1::agent_server,
    payload::v1::{
        object::Distance,
        search,
        info,
    },
    vald::v1::{
        flush_server, index_server,insert_server, object_server, remove_server, search_server, update_server, upsert_server
    }
};
use service::qbg::QBGService;
use std::collections::HashMap;
use std::time::Duration;

mod handler;
mod middleware;

macro_rules! new_svc {
    ($server:ty, $agent:expr, $settings:expr, $grpc_key:expr) => {
        <$server>::new($agent.clone())
            .max_decoding_message_size(
                $settings.get::<usize>(format!("{}.grpc.max_receive_message_size", $grpc_key).as_str())?,
            )
            .max_encoding_message_size(
                $settings.get::<usize>(format!("{}.grpc.max_send_message_size", $grpc_key).as_str())?,
            )
    };
}

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

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let addr = "0.0.0.0:8081".parse()?;
    let settings = Config::builder()
        .add_source(config::File::with_name("/etc/server/config.yaml"))
        .build()
        .unwrap();
    let _logger =
        flexi_logger::Logger::try_with_str(settings.get::<String>("logging.level")?)?.start()?;
    let service = QBGService::new(settings.clone());
    let agent = handler::Agent::new(
        service,
        "agent-qbg",
        "127.0.0.1",
        "vald/internal/core/algorithm",
        "vald-agent",
        10,
    );

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
        .add_service(new_svc!(agent_server::AgentServer<handler::Agent>, agent, settings, grpc_key))
        .add_service(new_svc!(search_server::SearchServer<handler::Agent>, agent, settings, grpc_key))
        .add_service(new_svc!(insert_server::InsertServer<handler::Agent>, agent, settings, grpc_key))
        .add_service(new_svc!(update_server::UpdateServer<handler::Agent>, agent, settings, grpc_key))
        .add_service(new_svc!(upsert_server::UpsertServer<handler::Agent>, agent, settings, grpc_key))
        .add_service(new_svc!(remove_server::RemoveServer<handler::Agent>, agent, settings, grpc_key))
        .add_service(new_svc!(object_server::ObjectServer<handler::Agent>, agent, settings, grpc_key))
        .add_service(new_svc!(index_server::IndexServer<handler::Agent>, agent, settings, grpc_key))
        .add_service(new_svc!(flush_server::FlushServer<handler::Agent>, agent, settings, grpc_key))
        .serve(addr)
        .await?;

    Ok(())
}
