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
use proto::{
    payload::v1::{object, remove},
    vald::v1::remove_server,
};

#[tonic::async_trait]
impl remove_server::Remove for super::Agent {
    async fn remove(
        &self,
        request: tonic::Request<remove::Request>,
    ) -> std::result::Result<tonic::Response<object::Location>, tonic::Status> {
        todo!()
    }

    #[doc = " A method to remove an indexed vector based on timestamp.\n"]
    async fn remove_by_timestamp(
        &self,
        request: tonic::Request<remove::TimestampRequest>,
    ) -> std::result::Result<tonic::Response<object::Locations>, tonic::Status> {
        todo!()
    }

    #[doc = " Server streaming response type for the StreamRemove method."]
    type StreamRemoveStream = crate::stream_type!(object::StreamLocation);

    #[doc = " A method to remove multiple indexed vectors by bidirectional streaming.\n"]
    async fn stream_remove(
        &self,
        request: tonic::Request<tonic::Streaming<remove::Request>>,
    ) -> std::result::Result<tonic::Response<Self::StreamRemoveStream>, tonic::Status> {
        todo!()
    }

    #[doc = " A method to remove multiple indexed vectors in a single request.\n"]
    async fn multi_remove(
        &self,
        request: tonic::Request<remove::MultiRequest>,
    ) -> std::result::Result<tonic::Response<object::Locations>, tonic::Status> {
        todo!()
    }
}
