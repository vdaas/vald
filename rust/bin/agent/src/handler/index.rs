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
use proto::{
    core::v1::agent_server,
    payload::v1::{control, info, object, Empty},
    vald::v1::index_server,
};

#[tonic::async_trait]
impl agent_server::Agent for super::Agent {
    async fn create_index(
        &self,
        request: tonic::Request<control::CreateIndexRequest>,
    ) -> std::result::Result<tonic::Response<Empty>, tonic::Status> {
        todo!()
    }

    async fn save_index(
        &self,
        request: tonic::Request<Empty>,
    ) -> std::result::Result<tonic::Response<Empty>, tonic::Status> {
        todo!()
    }

    #[doc = " Represent the creating and saving index RPC.\n"]
    async fn create_and_save_index(
        &self,
        request: tonic::Request<control::CreateIndexRequest>,
    ) -> std::result::Result<tonic::Response<Empty>, tonic::Status> {
        todo!()
    }
}

#[tonic::async_trait]
impl index_server::Index for super::Agent {
    #[doc = " Represent the RPC to get the agent index information.\n"]
    async fn index_info(
        &self,
        request: tonic::Request<Empty>,
    ) -> std::result::Result<tonic::Response<info::index::Count>, tonic::Status> {
        todo!()
    }

    #[doc = " Represent the RPC to get the agent index detailed information.\n"]
    async fn index_detail(
        &self,
        request: tonic::Request<Empty>,
    ) -> std::result::Result<tonic::Response<info::index::Detail>, tonic::Status> {
        todo!()
    }
    #[doc = " Represent the RPC to get the agent index statistics.\n"]
    async fn index_statistics(
        &self,
        request: tonic::Request<Empty>,
    ) -> std::result::Result<tonic::Response<info::index::Statistics>, tonic::Status> {
        todo!()
    }

    #[doc = " Represent the RPC to get the agent index detailed statistics.\n"]
    async fn index_statistics_detail(
        &self,
        request: tonic::Request<Empty>,
    ) -> std::result::Result<tonic::Response<info::index::StatisticsDetail>, tonic::Status> {
        todo!()
    }

    #[doc = " Represent the RPC to get the index property.\n"]
    async fn index_property(
        &self,
        request: tonic::Request<Empty>,
    ) -> std::result::Result<tonic::Response<info::index::PropertyDetail>, tonic::Status> {
        todo!()
    }
}
