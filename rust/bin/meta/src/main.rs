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

mod config;
mod handler;

use anyhow::Result;
use clap::Parser;
use config::Config;
use observability::observability::{Observability, ObservabilityImpl};
use opentelemetry::global;
use opentelemetry::propagation::Extractor;
use tonic::Request;
use tonic::transport::Server;

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

#[derive(Parser, Debug)]
#[command(author, version, about, long_about = None)]
pub struct Args {
    /// Path to the configuration YAML file.
    #[arg(short, long, env = "META_CONFIG")]
    pub config: Option<String>,

    /// Path to the database directory. Overrides the value in the configuration file.
    #[arg(short, long, env = "META_DATABASE_PATH")]
    pub database_path: Option<String>,

    /// Address to bind the server to (e.g., "[::1]:8095"). Overrides the value in the configuration file.
    #[arg(short, long, env = "META_SERVER_ADDR")]
    pub server_addr: Option<String>,
}

fn apply_args_to_config(config: &mut Config, args: &Args) {
    if let Some(db_path) = &args.database_path {
        config.database_path = db_path.clone();
    }

    if let Some(addr) = &args.server_addr {
        config.server_addr = addr.clone();
    }
}

#[tokio::main]
async fn main() -> Result<()> {
    let args = Args::parse();

    let mut cfg = Config::load(args.config.as_deref())?;
    apply_args_to_config(&mut cfg, &args);

    cfg.validate()?;

    let observability_cfg = cfg.observability.to_observability_config()?;
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

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_args_override() {
        let mut cfg = Config::default();
        let default_db = cfg.database_path.clone();
        let default_addr = cfg.server_addr.clone();

        // No overrides
        let args = Args {
            config: None,
            database_path: None,
            server_addr: None,
        };
        apply_args_to_config(&mut cfg, &args);
        assert_eq!(cfg.database_path, default_db);
        assert_eq!(cfg.server_addr, default_addr);

        // With overrides
        let args = Args {
            config: None,
            database_path: Some("/new/db/path".to_string()),
            server_addr: Some("127.0.0.1:9090".to_string()),
        };
        apply_args_to_config(&mut cfg, &args);
        assert_eq!(cfg.database_path, "/new/db/path");
        assert_eq!(cfg.server_addr, "127.0.0.1:9090");
    }
}
