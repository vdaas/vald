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

use proto::{meta::v1::meta_server, payload::v1::{meta, Empty}};

#[tonic::async_trait]
impl meta_server::Meta for super::Meta {
    async fn get(
        &self,
        request: tonic::Request<meta::Key>,
    ) -> std::result::Result<tonic::Response<meta::Value>, tonic::Status> {
        todo!()
    }
    async fn set(
        &self,
        request: tonic::Request<meta::KeyValue>,
    ) -> std::result::Result<tonic::Response<Empty>, tonic::Status> {
        todo!()
    }
    
    async fn delete(
        &self,
        request: tonic::Request<meta::Key>,
    ) -> std::result::Result<tonic::Response<Empty>, tonic::Status> {
        todo!()
    }
}
