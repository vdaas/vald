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
use std::{collections::HashMap, error, fmt};

use proto::{payload::v1::search, vald::v1::search_server};
use tonic::{Code, Status};
use tonic_types::{ErrorDetails, FieldViolation, StatusExt};

#[derive(Debug, Clone)]
pub struct CreateIndexingIsInProgress {}
#[derive(Debug, Clone)]
pub struct FlushingIsInProgress {}
#[derive(Debug, Clone)]
pub struct EmptySearchResult {}
#[derive(Debug, Clone)]
pub struct IncompatibleDimensionSize {
    got: usize,
    want: usize,
}

impl fmt::Display for CreateIndexingIsInProgress {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        write!(f, "create indexing is in progress")
    }
}

impl error::Error for CreateIndexingIsInProgress {
    fn source(&self) -> Option<&(dyn error::Error + 'static)> {
        None
    }
}

impl fmt::Display for FlushingIsInProgress {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        write!(f, "flush is in progress")
    }
}

impl error::Error for FlushingIsInProgress {
    fn source(&self) -> Option<&(dyn error::Error + 'static)> {
        None
    }
}

impl fmt::Display for EmptySearchResult {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        write!(f, "search result is empty")
    }
}

impl error::Error for EmptySearchResult {
    fn source(&self) -> Option<&(dyn error::Error + 'static)> {
        None
    }
}
impl IncompatibleDimensionSize {
    pub fn new(got: usize, want: usize) -> IncompatibleDimensionSize {
        IncompatibleDimensionSize {
            got: got,
            want: want,
        }
    }
}

impl fmt::Display for IncompatibleDimensionSize {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        write!(f, "incompatible dimension size detected\trequested: {},\tconfigured: {}", self.got, self.want)
    }
}

impl error::Error for IncompatibleDimensionSize {
    fn source(&self) -> Option<&(dyn error::Error + 'static)> {
        None
    }
}

#[tonic::async_trait]
impl search_server::Search for super::Agent {
    async fn search(
        &self,
        request: tonic::Request<search::Request>,
    ) -> Result<tonic::Response<search::Response>, Status> {
        let req = request.get_ref();
        let config = req.config.unwrap();
        let domain = cargo::util::hostname().unwrap_or("".into()).to_str().unwrap_or("");
        if req.vector.len() != self.s.get_dimension_size() {
            let err = IncompatibleDimensionSize::new(req.vector.len(), self.s.get_dimension_size());
            let msg = err.to_string();
            let mut err_details = ErrorDetails::new();
            err_details.set_error_info(msg, domain, HashMap::new());
            err_details.set_request_info(config.request_id, req.into());
            err_details.set_bad_request(vec![FieldViolation::new("vector dimension size", msg)]);
            err_details.set_resource_info(self.resource_type + "/ngt.Search", "", "", "");
            let status = Status::with_error_details(Code::InvalidArgument, "Search API Incombatible Dimension Size detedted", err_details);
            return Err(status);
        }

        let result = self.s.search(req.vector, config.num, config.epsilon, config.radius);
        match result {
            Err(err) => {
                let metadata = HashMap::new();
                let resource_type = self.resource_type + "/ngt.Search";
                let resource_name = format!("{}: {}({})", self.api_name, self.name, self.ip);
                let status = match err {
                    CreateIndexingIsInProgress => {
                        let mut err_details = ErrorDetails::new();
                        err_details.set_error_info(err.to_string(), domain, metadata);
                        err_details.set_request_info(config.request_id, req.into());
                        err_details.set_resource_info(resource_type, resource_name, "", "");
                        Status::with_error_details(Code::Aborted, "Search API aborted to process search request due to creating indices is in progress", err_details)
                    }
                    FlushingIsInProgress => {
                        let mut err_details = ErrorDetails::new();
                        err_details.set_error_info(err.to_string(), domain, metadata);
                        err_details.set_request_info(config.request_id, req.into());
                        err_details.set_resource_info(resource_type, resource_name, "", "");
                        Status::with_error_details(Code::Aborted, "Search API aborted to process search request due to flushing indices is in progress", err_details)
                    }
                    EmptySearchResult => {
                        let mut err_details = ErrorDetails::new();
                        err_details.set_error_info(err.to_string(), domain, metadata);
                        err_details.set_request_info(config.request_id, req.into());
                        err_details.set_resource_info(resource_type, resource_name, "", "");
                        Status::with_error_details(Code::NotFound, format!("Search API requestID {}'s search result not found", config.request_id), err_details)
                    }
                    IncompatibleDimensionSize => {
                        let mut err_details = ErrorDetails::new();
                        err_details.set_error_info(err.to_string(), domain, metadata);
                        err_details.set_request_info(config.request_id, req.into());
                        err_details.set_resource_info(resource_type, resource_name, "", "");
                        err_details.set_bad_request(vec![FieldViolation::new("vector dimension size", err.to_string())]);
                        Status::with_error_details(Code::InvalidArgument, "Search API Incompatible Dimension Size detected", err_details)
                    }
                    _ => {
                        let mut err_details = ErrorDetails::new();
                        err_details.set_error_info(err.to_string(), domain, metadata);
                        err_details.set_request_info(config.request_id, req.into());
                        err_details.set_resource_info(resource_type, resource_name, "", "");
                        Status::with_error_details(Code::Internal, "Search API failed to process search request", err_details)
                    }
                };
                Err(status)
            }
            Ok(mut response) => {
                response.get_mut().request_id = config.request_id;
                Ok(response)
            }
        }
    }

    #[doc = " A method to search indexed vectors by ID.\n"]
    async fn search_by_id(
        &self,
        _request: tonic::Request<search::IdRequest>,
    ) -> Result<tonic::Response<search::Response>, tonic::Status> {
        todo!()
    }

    #[doc = " Server streaming response type for the StreamSearch method."]
    type StreamSearchStream = crate::stream_type!(search::StreamResponse);

    #[doc = " A method to search indexed vectors by multiple vectors.\n"]
    async fn stream_search(
        &self,
        _request: tonic::Request<tonic::Streaming<search::Request>>,
    ) -> std::result::Result<tonic::Response<Self::StreamSearchStream>, tonic::Status> {
        todo!()
    }

    #[doc = " Server streaming response type for the StreamSearchByID method."]
    type StreamSearchByIDStream = crate::stream_type!(search::StreamResponse);

    #[doc = " A method to search indexed vectors by multiple IDs.\n"]
    async fn stream_search_by_id(
        &self,
        _request: tonic::Request<tonic::Streaming<search::IdRequest>>,
    ) -> std::result::Result<tonic::Response<Self::StreamSearchByIDStream>, tonic::Status> {
        todo!()
    }

    #[doc = " A method to search indexed vectors by multiple vectors in a single request.\n"]
    async fn multi_search(
        &self,
        _request: tonic::Request<search::MultiRequest>,
    ) -> std::result::Result<tonic::Response<search::Responses>, tonic::Status> {
        todo!()
    }

    #[doc = " A method to search indexed vectors by multiple IDs in a single request.\n"]
    async fn multi_search_by_id(
        &self,
        _request: tonic::Request<search::MultiIdRequest>,
    ) -> std::result::Result<tonic::Response<search::Responses>, tonic::Status> {
        todo!()
    }

    #[doc = " A method to linear search indexed vectors by a raw vector.\n"]
    async fn linear_search(
        &self,
        _request: tonic::Request<search::Request>,
    ) -> std::result::Result<tonic::Response<search::Response>, tonic::Status> {
        todo!()
    }

    #[doc = " A method to linear search indexed vectors by ID.\n"]
    async fn linear_search_by_id(
        &self,
        _request: tonic::Request<search::IdRequest>,
    ) -> std::result::Result<tonic::Response<search::Response>, tonic::Status> {
        todo!()
    }

    #[doc = " Server streaming response type for the StreamLinearSearch method."]
    type StreamLinearSearchStream = crate::stream_type!(search::StreamResponse);

    #[doc = " A method to linear search indexed vectors by multiple vectors.\n"]
    async fn stream_linear_search(
        &self,
        _request: tonic::Request<tonic::Streaming<search::Request>>,
    ) -> std::result::Result<tonic::Response<Self::StreamLinearSearchStream>, tonic::Status> {
        todo!()
    }

    #[doc = " Server streaming response type for the StreamLinearSearchByID method."]
    type StreamLinearSearchByIDStream = crate::stream_type!(search::StreamResponse);

    #[doc = " A method to linear search indexed vectors by multiple IDs.\n"]
    async fn stream_linear_search_by_id(
        &self,
        _request: tonic::Request<tonic::Streaming<search::IdRequest>>,
    ) -> std::result::Result<tonic::Response<Self::StreamLinearSearchByIDStream>, tonic::Status>
    {
        todo!()
    }

    #[doc = " A method to linear search indexed vectors by multiple vectors in a single\n request.\n"]
    async fn multi_linear_search(
        &self,
        _request: tonic::Request<search::MultiRequest>,
    ) -> std::result::Result<tonic::Response<search::Responses>, tonic::Status> {
        todo!()
    }

    #[doc = " A method to linear search indexed vectors by multiple IDs in a single\n request.\n"]
    async fn multi_linear_search_by_id(
        &self,
        _request: tonic::Request<search::MultiIdRequest>,
    ) -> std::result::Result<tonic::Response<search::Responses>, tonic::Status> {
        todo!()
    }
}
