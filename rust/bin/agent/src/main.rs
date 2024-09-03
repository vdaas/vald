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
use std::{thread::sleep, time::Duration};

use algorithm::Error;
use anyhow::Result;
use proto::{payload::v1::search, vald::v1::{search_client::SearchClient, search_server::SearchServer}};
use tonic::transport::Server;

mod handler;

#[derive(Debug)]
struct MockService {
    dim: usize
}

impl algorithm::ANN for MockService {
    fn get_dimension_size(&self) -> usize {
        self.dim
    }

    fn search(&self, _vector: Vec<f32>, dim: u32, _epsilon: f32, _radius: f32) -> Result<tonic::Response<search::Response>, Error> {
        Err(Error::IncompatibleDimensionSize{got: dim as usize, want: self.dim}.into())
    }
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let addr = "[::1]:8081".parse().unwrap();

    let service = MockService{ dim: 42 };
    let agent = handler::Agent::new(service, "agent-ngt", "127.0.0.1", "vald/internal/core/algorithm", "vald-agent");

    tokio::spawn(async move {
        Server::builder()
            .add_service(SearchServer::new(agent))
            .serve(addr)
            .await
    });

    sleep(Duration::from_secs(3));
    
    let mut client = SearchClient::connect("http://[::1]:8081").await?;

    let cfg = search::Config::default();
    let request = tonic::Request::new(search::Request { vector: vec![0.0], config: Some(cfg) });

    let response = client.search(request).await?;

    println!("RESPONSE={:?}", response);

    Ok(())
}
