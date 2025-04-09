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
mod common;
pub mod index;
pub mod insert;
pub mod object;
pub mod remove;
pub mod search;
pub mod update;
pub mod upsert;
use std::sync::Arc;
use tokio::sync::RwLock;

pub struct Agent {
    s: Arc<RwLock<dyn algorithm::ANN>>,
    name: String,
    ip: String,
    resource_type: String,
    api_name: String,
    stream_concurrency: usize,
}

impl Agent {
    pub fn new(
        s: impl algorithm::ANN + 'static,
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
            stream_concurrency: stream_concurrency
        }
    }
}
