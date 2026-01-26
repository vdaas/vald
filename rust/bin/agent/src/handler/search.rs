//
// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
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
use std::sync::Arc;
use tokio::sync::RwLock;
use tonic::{Code, Status};
use tonic_types::StatusExt;

use super::common::{bidirectional_stream, build_error_details};

async fn search<S: algorithm::ANN>(
    s: Arc<RwLock<S>>,
    resource_type: &str,
    api_name: &str,
    name: &str,
    ip: &str,
    request: &search::Request,
) -> Result<search::Response, Status> {
    let config = match request.config.clone() {
        Some(cfg) => cfg,
        None => return Err(Status::invalid_argument("Missing configuration in request")),
    };
    let hostname = cargo::util::hostname()?;
    let domain = hostname.to_str().unwrap();
    {
        let s = s.read().await;
        if request.vector.len() != s.get_dimension_size() {
            let err = Error::IncompatibleDimensionSize {
                got: request.vector.len(),
                want: s.get_dimension_size(),
            };
            let resource_type = format!("{}/qbg.Search", resource_type);
            let resource_name = format!("{}: {}({})", api_name, name, ip);
            let err_details = build_error_details(
                err,
                domain,
                &config.request_id,
                request.encode_to_vec(),
                &resource_type,
                &resource_name,
                Some("vector dimension size"),
            );
            let status = Status::with_error_details(
                Code::InvalidArgument,
                "Search API Incombatible Dimension Size detedted",
                err_details,
            );
            warn!("{:?}", status);
            return Err(status);
        }
        let result = s.search(
            request.vector.clone(),
            config.num,
            config.epsilon,
            config.radius,
        ).await;
        match result {
            Err(err) => {
                let resource_type = format!("{}/qbg.Search", resource_type);
                let resource_name = format!("{}: {}({})", api_name, name, ip);
                let request_bytes = request.encode_to_vec();
                let status = match err {
                    Error::CreateIndexingIsInProgress {} => {
                        let err_details = build_error_details(
                            err,
                            domain,
                            &config.request_id,
                            request_bytes,
                            &resource_type,
                            &resource_name,
                            None,
                        );
                        let status = Status::with_error_details(Code::Aborted, "Search API aborted to process search request due to creating indices is in progress", err_details);
                        debug!("{:?}", status);
                        status
                    }
                    Error::FlushingIsInProgress {} => {
                        let err_details = build_error_details(
                            err,
                            domain,
                            &config.request_id,
                            request_bytes,
                            &resource_type,
                            &resource_name,
                            None,
                        );
                        let status = Status::with_error_details(Code::Aborted, "Search API aborted to process search request due to flushing indices is in progress", err_details);
                        debug!("{:?}", status);
                        status
                    }
                    Error::EmptySearchResult {} => {
                        let err_details = build_error_details(
                            err,
                            domain,
                            &config.request_id,
                            request_bytes,
                            &resource_type,
                            &resource_name,
                            None,
                        );
                        let status = Status::with_error_details(
                            Code::NotFound,
                            format!(
                                "Search API requestID {}'s search result not found",
                                &config.request_id,
                            ),
                            err_details,
                        );
                        debug!("{:?}", status);
                        status
                    }
                    Error::IncompatibleDimensionSize { got: _, want: _ } => {
                        let err_details = build_error_details(
                            err,
                            domain,
                            &config.request_id,
                            request_bytes,
                            &resource_type,
                            &resource_name,
                            Some("vector dimension size"),
                        );
                        let status = Status::with_error_details(
                            Code::InvalidArgument,
                            "Search API Incompatible Dimension Size detected",
                            err_details,
                        );
                        warn!("{:?}", status);
                        status
                    }
                    _ => {
                        let err_details = build_error_details(
                            err,
                            domain,
                            &config.request_id,
                            request_bytes,
                            &resource_type,
                            &resource_name,
                            None,
                        );
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
                Ok(response)
            }
        }
    }
}

#[tonic::async_trait]
impl<S: algorithm::ANN + 'static> search_server::Search for super::Agent<S> {
    async fn search(
        &self,
        request: tonic::Request<search::Request>,
    ) -> Result<tonic::Response<search::Response>, Status> {
        info!("Recieved a request from {:?}", request.remote_addr());
        let request = request.get_ref();
        let s = self.s.clone();
        let resource_type = self.resource_type.clone();
        let name = self.name.clone();
        let ip = self.ip.clone();
        let api_name = self.api_name.clone();
        match search(s, &resource_type, &api_name, &name, &ip, request).await {
            Ok(response) => Ok(tonic::Response::new(response)),
            Err(e) => Err(e),
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
        request: tonic::Request<tonic::Streaming<search::Request>>,
    ) -> std::result::Result<tonic::Response<Self::StreamSearchStream>, tonic::Status> {
        info!(
            "Received stream search request from {:?}",
            request.remote_addr()
        );

        let s = self.s.clone();
        let resource_type = self.resource_type.clone() + "/qbg.StreamSearch";
        let name = self.name.clone();
        let ip = self.ip.clone();
        let api_name = self.api_name.clone();

        let process_fn = move |req: search::Request| {
            let s = s.clone();
            let resource_type = resource_type.clone();
            let name = name.clone();
            let ip = ip.clone();
            let api_name = api_name.clone();
            async move {
                match search(s, &resource_type, &api_name, &name, &ip, &req).await {
                    Ok(response) => Ok(search::StreamResponse {
                        payload: Some(search::stream_response::Payload::Response(response)),
                    }),
                    Err(status) => Err(status),
                }
            }
        };

        bidirectional_stream(request, self.stream_concurrency, process_fn).await
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
        request: tonic::Request<search::MultiRequest>,
    ) -> std::result::Result<tonic::Response<search::Responses>, tonic::Status> {
        info!("Recieved a request from {:?}", request.remote_addr());
        let mreq = request.get_ref();
        let hostname = cargo::util::hostname()?;
        let _domain = hostname.to_str().unwrap();
        let mut res = search::Responses { responses: vec![] };
        for req in mreq.requests.clone() {
            let response = self.search(tonic::Request::new(req)).await?;
            res.responses.push(response.into_inner());
        }
        Ok(tonic::Response::new(res))
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
