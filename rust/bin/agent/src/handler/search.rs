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
use proto::{payload::v1::search, vald::v1::search_server};

#[tonic::async_trait]
impl search_server::Search for super::Agent {
    async fn search(
        &self,
        request: tonic::Request<search::Request>,
    ) -> Result<tonic::Response<search::Response>, tonic::Status> {
        todo!()
    }

    #[doc = " A method to search indexed vectors by ID.\n"]
    async fn search_by_id(
        &self,
        request: tonic::Request<search::IdRequest>,
    ) -> Result<tonic::Response<search::Response>, tonic::Status> {
        todo!()
    }

    #[doc = " Server streaming response type for the StreamSearch method."]
    type StreamSearchStream = crate::stream_type!(search::StreamResponse);

    #[doc = " A method to search indexed vectors by multiple vectors.\n"]
    async fn stream_search(
        &self,
        request: tonic::Request<tonic::Streaming<search::Request>>,
    ) -> std::result::Result<tonic::Response<Self::StreamSearchStream>, tonic::Status> {
        todo!()
    }

    #[doc = " Server streaming response type for the StreamSearchByID method."]
    type StreamSearchByIDStream = crate::stream_type!(search::StreamResponse);

    #[doc = " A method to search indexed vectors by multiple IDs.\n"]
    async fn stream_search_by_id(
        &self,
        request: tonic::Request<tonic::Streaming<search::IdRequest>>,
    ) -> std::result::Result<tonic::Response<Self::StreamSearchByIDStream>, tonic::Status> {
        todo!()
    }

    #[doc = " A method to search indexed vectors by multiple vectors in a single request.\n"]
    async fn multi_search(
        &self,
        request: tonic::Request<search::MultiRequest>,
    ) -> std::result::Result<tonic::Response<search::Responses>, tonic::Status> {
        todo!()
    }

    #[doc = " A method to search indexed vectors by multiple IDs in a single request.\n"]
    async fn multi_search_by_id(
        &self,
        request: tonic::Request<search::MultiIdRequest>,
    ) -> std::result::Result<tonic::Response<search::Responses>, tonic::Status> {
        todo!()
    }

    #[doc = " A method to linear search indexed vectors by a raw vector.\n"]
    async fn linear_search(
        &self,
        request: tonic::Request<search::Request>,
    ) -> std::result::Result<tonic::Response<search::Response>, tonic::Status> {
        todo!()
    }

    #[doc = " A method to linear search indexed vectors by ID.\n"]
    async fn linear_search_by_id(
        &self,
        request: tonic::Request<search::IdRequest>,
    ) -> std::result::Result<tonic::Response<search::Response>, tonic::Status> {
        todo!()
    }

    #[doc = " Server streaming response type for the StreamLinearSearch method."]
    type StreamLinearSearchStream = crate::stream_type!(search::StreamResponse);

    #[doc = " A method to linear search indexed vectors by multiple vectors.\n"]
    async fn stream_linear_search(
        &self,
        request: tonic::Request<tonic::Streaming<search::Request>>,
    ) -> std::result::Result<tonic::Response<Self::StreamLinearSearchStream>, tonic::Status> {
        todo!()
    }

    #[doc = " Server streaming response type for the StreamLinearSearchByID method."]
    type StreamLinearSearchByIDStream = crate::stream_type!(search::StreamResponse);

    #[doc = " A method to linear search indexed vectors by multiple IDs.\n"]
    async fn stream_linear_search_by_id(
        &self,
        request: tonic::Request<tonic::Streaming<search::IdRequest>>,
    ) -> std::result::Result<tonic::Response<Self::StreamLinearSearchByIDStream>, tonic::Status>
    {
        todo!()
    }

    #[doc = " A method to linear search indexed vectors by multiple vectors in a single\n request.\n"]
    async fn multi_linear_search(
        &self,
        request: tonic::Request<search::MultiRequest>,
    ) -> std::result::Result<tonic::Response<search::Responses>, tonic::Status> {
        todo!()
    }

    #[doc = " A method to linear search indexed vectors by multiple IDs in a single\n request.\n"]
    async fn multi_linear_search_by_id(
        &self,
        request: tonic::Request<search::MultiIdRequest>,
    ) -> std::result::Result<tonic::Response<search::Responses>, tonic::Status> {
        todo!()
    }
}
