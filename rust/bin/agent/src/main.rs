//
// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
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

use anyhow::Result;
use proto::payload::v1::search;

mod handler;

#[derive(Debug)]
struct MockService {}

impl algorithm::ANN for MockService {
    fn get_dimension_size(&self) -> usize {
        42
    }

    fn search(&self, vector: Vec<f32>, dim: usize, epsilon: f64, radius: f64) -> Result<tonic::Response<search::Response>> {
        Err(handler::search::IncompatibleDimensionSize::new(dim, 42))
    }
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let addr = "[::1]:8081".parse()?;
    let service = MockService{};
    let agent = handler::Agent::new(service, "agent-ngt", "127.0.0.1", "vald/internal/core/algorithm", "vald-agent");

    tonic::transport::Server::builder()
        .add_service(proto::core::v1::agent_server::AgentServer::new(agent))
        .serve(addr)
        .await?;

    Ok(())
}
