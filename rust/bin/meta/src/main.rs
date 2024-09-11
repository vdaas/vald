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

mod handler;

use opentelemetry::global;
use opentelemetry::propagation::Extractor;
use tonic::transport::Server;
use tonic::Request;
use observability::{config::Config, observability::{Observability, ObservabilityImpl}};

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
            .collect::<Vec<_>>()    }
}

fn intercept(mut req: Request<()>) -> Result<Request<()>, tonic::Status> {
    let parent_cx = global::get_text_map_propagator(|prop| {
        prop.extract(&MetadataMap(req.metadata()))
    });
    req.extensions_mut().insert(parent_cx);
    Ok(req)
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    // TODO: yaml is given from outside
    // let config_yaml = r#"enabled: false
    // endpoint: ""
    // attributes: []
    // tracer:
    //     enabled: false
    // meter:
    //     enabled: false
    //     export_duration:
    //         secs: 60
    //         nanos: 0
    //     export_timeout_duration:
    //         secs: 30
    //         nanos: 0
    // "#;
    //
    // decode config yaml
    // let observability_cfg = serde_yaml::from_str(config_yaml).unwrap();
    let observability_cfg = Config::default();
    let mut observability = ObservabilityImpl::new(observability_cfg)?;

    let addr = "[::1]:8081".parse()?;
    let cfg_path = "/var/lib/meta/database"; // TODO: set the appropriate path
    let meta = handler::Meta::new(cfg_path)?;

    tonic::transport::Server::builder()
        .add_service(proto::meta::v1::meta_server::MetaServer::with_interceptor(meta, intercept))
        .serve(addr)
        .await?;

    observability.shutdown()?;
    Ok(())
}
