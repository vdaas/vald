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
use algorithm::Error;
use log::{debug, error, info, warn};
use prost::Message;
use proto::{payload::v1::search, vald::v1::search_server};
use std::{collections::HashMap, string::String};
use tonic::{Code, Status};
use tonic_types::{ErrorDetails, FieldViolation, StatusExt};

#[tonic::async_trait]
impl search_server::Search for super::Agent {
    async fn search(
        &self,
        request: tonic::Request<search::Request>,
    ) -> Result<tonic::Response<search::Response>, Status> {
        info!("Recieved a request from {:?}", request.remote_addr());
        let req = request.get_ref();
        let config = match req.config.clone() {
            Some(cfg) => cfg,
            None => return Err(Status::invalid_argument("Missing configuration in request")),
        };
        let hostname = cargo::util::hostname()?;
        let domain = hostname.to_str().unwrap();
        {
            let s = self.s.read().await;
            if req.vector.len() != s.get_dimension_size() {
                let err = Error::IncompatibleDimensionSize {
                    got: req.vector.len(),
                    want: s.get_dimension_size(),
                };
                let metadata = HashMap::new();
                let resource_type = self.resource_type.clone() + "/qbg.Search";
                let resource_name = format!("{}: {}({})", self.api_name, self.name, self.ip);
                let mut err_details = ErrorDetails::new();
                err_details.set_error_info(err.to_string(), domain, metadata);
                err_details.set_request_info(
                    config.request_id,
                    String::from_utf8(req.encode_to_vec())
                        .unwrap_or_else(|_| "<invalid UTF-8>".to_string()),
                );
                err_details.set_bad_request(vec![FieldViolation::new(
                    "vector dimension size",
                    err.to_string(),
                )]);
                err_details.set_resource_info(resource_type, resource_name, "", "");
                let status = Status::with_error_details(
                    Code::InvalidArgument,
                    "Search API Incombatible Dimension Size detedted",
                    err_details,
                );
                warn!("{:?}", status);
                return Err(status);
            }
            let result = s.search(
                req.vector.clone(),
                config.num,
                config.epsilon,
                config.radius,
            );
            match result {
                Err(err) => {
                    let metadata = HashMap::new();
                    let resource_type = self.resource_type.clone() + "/qbg.Search";
                    let resource_name = format!("{}: {}({})", self.api_name, self.name, self.ip);
                    let status = match err {
                        Error::CreateIndexingIsInProgress {} => {
                            let mut err_details = ErrorDetails::new();
                            err_details.set_error_info(err.to_string(), domain, metadata);
                            err_details.set_request_info(
                                config.request_id,
                                String::from_utf8(req.encode_to_vec())
                                    .unwrap_or_else(|_| "<invalid UTF-8>".to_string()),
                            );
                            err_details.set_resource_info(resource_type, resource_name, "", "");
                            let status = Status::with_error_details(Code::Aborted, "Search API aborted to process search request due to creating indices is in progress", err_details);
                            debug!("{:?}", status);
                            status
                        }
                        Error::FlushingIsInProgress {} => {
                            let mut err_details = ErrorDetails::new();
                            err_details.set_error_info(err.to_string(), domain, metadata);
                            err_details.set_request_info(
                                config.request_id,
                                String::from_utf8(req.encode_to_vec())
                                    .unwrap_or_else(|_| "<invalid UTF-8>".to_string()),
                            );
                            err_details.set_resource_info(resource_type, resource_name, "", "");
                            let status = Status::with_error_details(Code::Aborted, "Search API aborted to process search request due to flushing indices is in progress", err_details);
                            debug!("{:?}", status);
                            status
                        }
                        Error::EmptySearchResult {} => {
                            let request_id = config.request_id;
                            let mut err_details = ErrorDetails::new();
                            err_details.set_error_info(err.to_string(), domain, metadata);
                            err_details.set_request_info(
                                &request_id,
                                String::from_utf8(req.encode_to_vec())
                                    .unwrap_or_else(|_| "<invalid UTF-8>".to_string()),
                            );
                            err_details.set_resource_info(resource_type, resource_name, "", "");
                            let status = Status::with_error_details(
                                Code::NotFound,
                                format!(
                                    "Search API requestID {}'s search result not found",
                                    &request_id
                                ),
                                err_details,
                            );
                            debug!("{:?}", status);
                            status
                        }
                        Error::IncompatibleDimensionSize { got: _, want: _ } => {
                            let mut err_details = ErrorDetails::new();
                            err_details.set_error_info(err.to_string(), domain, metadata);
                            err_details.set_request_info(
                                config.request_id,
                                String::from_utf8(req.encode_to_vec())
                                    .unwrap_or_else(|_| "<invalid UTF-8>".to_string()),
                            );
                            err_details.set_resource_info(resource_type, resource_name, "", "");
                            err_details.set_bad_request(vec![FieldViolation::new(
                                "vector dimension size",
                                err.to_string(),
                            )]);
                            let status = Status::with_error_details(
                                Code::InvalidArgument,
                                "Search API Incompatible Dimension Size detected",
                                err_details,
                            );
                            warn!("{:?}", status);
                            status
                        }
                        _ => {
                            let mut err_details = ErrorDetails::new();
                            err_details.set_error_info(err.to_string(), domain, metadata);
                            err_details.set_request_info(
                                config.request_id,
                                String::from_utf8(req.encode_to_vec())
                                    .unwrap_or_else(|_| "<invalid UTF-8>".to_string()),
                            );
                            err_details.set_resource_info(resource_type, resource_name, "", "");
                            let status = Status::with_error_details(
                                Code::Internal,
                                "Search API failed to process search request",
                                err_details,
                            );
                            error!("{:?}", status);
                            status
                        }
                    };
                    Err(status)
                }
                Ok(mut response) => {
                    response.request_id = config.request_id;
                    Ok(tonic::Response::new(response))
                }
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
