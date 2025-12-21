//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
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

mod config;
mod handler;

use clap::Parser;
use observability::observability::{Observability, ObservabilityImpl};
use opentelemetry::global;
use opentelemetry::propagation::Extractor;
use tonic::transport::Server;
use tonic::Request;

#[derive(Parser, Debug)]
#[command(author, version, about, long_about = None)]
struct Args {
    #[arg(short, long)]
    config: Option<String>,

    #[arg(short, long)]
    database_path: Option<String>,

    #[arg(short, long)]
    server_addr: Option<String>,
}

struct MetadataMap<'a>(&'a tonic::metadata::MetadataMap);

impl<'a> Extractor for MetadataMap<'a> {
    fn get(&self, key: &str) -> Option<&str> {
        self.0.get(key).and_then(|metadata| metadata.to_str().ok())
    }

    fn keys(&self) -> Vec<&str> {
        self.0
            .keys()
            .map(|key| match key {
                tonic::metadata::KeyRef::Ascii(v) => v.as_str(),
                tonic::metadata::KeyRef::Binary(v) => v.as_str(),
            })
            .collect::<Vec<_>>()
    }
}

fn intercept(mut req: Request<()>) -> Result<Request<()>, tonic::Status> {
    let parent_cx =
        global::get_text_map_propagator(|prop| prop.extract(&MetadataMap(req.metadata())));
    req.extensions_mut().insert(parent_cx);
    Ok(req)
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let args = Args::parse();
    let mut cfg = if let Some(config_path) = args.config {
        config::Config::load(&config_path)?
    } else {
        config::Config::default()
    };

    if let Some(db_path) = args.database_path {
        cfg.database_path = db_path;
    }
    if let Some(addr) = args.server_addr {
        cfg.server_addr = addr;
    }

    let observability_cfg = cfg.observability.to_observability_config();
    let mut observability = ObservabilityImpl::new(observability_cfg)?;

    let addr = cfg.server_addr.parse()?;
    let meta = handler::Meta::new(&cfg.database_path)?;

    // the interceptor given here is implicitly executed for each request
    Server::builder()
        .add_service(proto::meta::v1::meta_server::MetaServer::with_interceptor(
            meta, intercept,
        ))
        .serve(addr)
        .await?;

    observability.shutdown()?;
    Ok(())
}
