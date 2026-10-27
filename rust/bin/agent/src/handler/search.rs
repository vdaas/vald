// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
use algorithm::Error;
use log::{debug, error, info, warn};
use prost::Message;
use proto::{payload::v1::search, vald::v1::search_server};
use std::sync::Arc;
use tokio::sync::RwLock;
use tonic::{Code, Status};
use tonic_types::StatusExt;

use super::common::{bidirectional_stream, build_error_details};

async fn search(
    s: Arc<RwLock<Box<dyn algorithm::ANN>>>,
    resource_type: &str,
    api_name: &str,
    name: &str,
    ip: &str,
    request: &search::Request,
    linear: bool,
) -> Result<search::Response, Status> {
    let config = match request.config.clone() {
        Some(cfg) => cfg,
        None => return Err(Status::invalid_argument("Missing configuration in request")),
    };
    let hostname = super::common::get_hostname();
    let domain = hostname.as_str();
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
        let result = if linear {
            s.linear_search_with_options(request.vector.clone(), config.num, &config.options)
        } else {
            s.search_with_options(
                request.vector.clone(),
                config.num,
                config.epsilon,
                config.radius,
                &config.options,
            )
        };
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
                        let status = Status::with_error_details(
                            Code::Aborted,
                            "Search API aborted to process search request due to creating indices is in progress",
                            err_details,
                        );
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
                        let status = Status::with_error_details(
                            Code::Aborted,
                            "Search API aborted to process search request due to flushing indices is in progress",
                            err_details,
                        );
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

async fn resolve_id(
    s: &Arc<RwLock<Box<dyn algorithm::ANN>>>,
    resource_type: &str,
    api_name: &str,
    name: &str,
    ip: &str,
    request: search::IdRequest,
) -> Result<search::Request, Status> {
    let result = if request.id.is_empty() {
        Err(Error::InvalidUUID {
            uuid: request.id.clone(),
        })
    } else {
        s.read().await.get_object(request.id.clone())
    };
    let (vector, _) = result.map_err(|err| {
        let code = match &err {
            Error::InvalidUUID { .. } => Code::InvalidArgument,
            Error::ObjectIDNotFound { .. } | Error::UUIDNotFound { .. } => Code::NotFound,
            _ => Code::Internal,
        };
        let message = err.to_string();
        let details = build_error_details(
            err,
            &super::common::get_hostname(),
            &request.id,
            request.encode_to_vec(),
            &format!("{resource_type}/SearchByID"),
            &format!("{api_name}: {name}({ip})"),
            None,
        );
        Status::with_error_details(code, message, details)
    })?;
    Ok(search::Request {
        vector,
        config: request.config,
    })
}

#[tonic::async_trait]
impl search_server::Search for super::Agent {
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
        match search(s, &resource_type, &api_name, &name, &ip, request, false).await {
            Ok(response) => Ok(tonic::Response::new(response)),
            Err(e) => Err(e),
        }
    }

    #[doc = " A method to search indexed vectors by ID.\n"]
    async fn search_by_id(
        &self,
        request: tonic::Request<search::IdRequest>,
    ) -> Result<tonic::Response<search::Response>, Status> {
        let req = resolve_id(
            &self.s,
            &self.resource_type,
            &self.api_name,
            &self.name,
            &self.ip,
            request.into_inner(),
        )
        .await?;
        self.search(tonic::Request::new(req)).await
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
                match search(s, &resource_type, &api_name, &name, &ip, &req, false).await {
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
        request: tonic::Request<tonic::Streaming<search::IdRequest>>,
    ) -> Result<tonic::Response<Self::StreamSearchByIDStream>, Status> {
        let s = self.s.clone();
        let resource_type = self.resource_type.clone();
        let api_name = self.api_name.clone();
        let name = self.name.clone();
        let ip = self.ip.clone();
        let process_fn = move |req: search::IdRequest| {
            let (s, resource_type, api_name, name, ip) = (
                s.clone(),
                resource_type.clone(),
                api_name.clone(),
                name.clone(),
                ip.clone(),
            );
            async move {
                let req = resolve_id(&s, &resource_type, &api_name, &name, &ip, req).await?;
                let response =
                    search(s, &resource_type, &api_name, &name, &ip, &req, false).await?;
                Ok(search::StreamResponse {
                    payload: Some(search::stream_response::Payload::Response(response)),
                })
            }
        };
        bidirectional_stream(request, self.stream_concurrency, process_fn).await
    }

    #[doc = " A method to search indexed vectors by multiple vectors in a single request.\n"]
    async fn multi_search(
        &self,
        request: tonic::Request<search::MultiRequest>,
    ) -> std::result::Result<tonic::Response<search::Responses>, tonic::Status> {
        info!("Recieved a request from {:?}", request.remote_addr());
        let mreq = request.get_ref();
        let hostname = super::common::get_hostname();
        let _domain = hostname.as_str();
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
        request: tonic::Request<search::MultiIdRequest>,
    ) -> Result<tonic::Response<search::Responses>, Status> {
        let mut responses = Vec::new();
        for req in request.into_inner().requests {
            responses.push(
                self.search_by_id(tonic::Request::new(req))
                    .await?
                    .into_inner(),
            );
        }
        Ok(tonic::Response::new(search::Responses { responses }))
    }

    #[doc = " A method to linear search indexed vectors by a raw vector.\n"]
    async fn linear_search(
        &self,
        request: tonic::Request<search::Request>,
    ) -> Result<tonic::Response<search::Response>, Status> {
        search(
            self.s.clone(),
            &self.resource_type,
            &self.api_name,
            &self.name,
            &self.ip,
            request.get_ref(),
            true,
        )
        .await
        .map(tonic::Response::new)
    }

    #[doc = " A method to linear search indexed vectors by ID.\n"]
    async fn linear_search_by_id(
        &self,
        request: tonic::Request<search::IdRequest>,
    ) -> Result<tonic::Response<search::Response>, Status> {
        let req = resolve_id(
            &self.s,
            &self.resource_type,
            &self.api_name,
            &self.name,
            &self.ip,
            request.into_inner(),
        )
        .await?;
        self.linear_search(tonic::Request::new(req)).await
    }

    #[doc = " Server streaming response type for the StreamLinearSearch method."]
    type StreamLinearSearchStream = crate::stream_type!(search::StreamResponse);

    #[doc = " A method to linear search indexed vectors by multiple vectors.\n"]
    async fn stream_linear_search(
        &self,
        request: tonic::Request<tonic::Streaming<search::Request>>,
    ) -> Result<tonic::Response<Self::StreamLinearSearchStream>, Status> {
        let s = self.s.clone();
        let resource_type = self.resource_type.clone();
        let api_name = self.api_name.clone();
        let name = self.name.clone();
        let ip = self.ip.clone();
        let process_fn = move |req: search::Request| {
            let (s, resource_type, api_name, name, ip) = (
                s.clone(),
                resource_type.clone(),
                api_name.clone(),
                name.clone(),
                ip.clone(),
            );
            async move {
                let response = search(s, &resource_type, &api_name, &name, &ip, &req, true).await?;
                Ok(search::StreamResponse {
                    payload: Some(search::stream_response::Payload::Response(response)),
                })
            }
        };
        bidirectional_stream(request, self.stream_concurrency, process_fn).await
    }

    #[doc = " Server streaming response type for the StreamLinearSearchByID method."]
    type StreamLinearSearchByIDStream = crate::stream_type!(search::StreamResponse);

    #[doc = " A method to linear search indexed vectors by multiple IDs.\n"]
    async fn stream_linear_search_by_id(
        &self,
        request: tonic::Request<tonic::Streaming<search::IdRequest>>,
    ) -> Result<tonic::Response<Self::StreamLinearSearchByIDStream>, Status> {
        let s = self.s.clone();
        let resource_type = self.resource_type.clone();
        let api_name = self.api_name.clone();
        let name = self.name.clone();
        let ip = self.ip.clone();
        let process_fn = move |req: search::IdRequest| {
            let (s, resource_type, api_name, name, ip) = (
                s.clone(),
                resource_type.clone(),
                api_name.clone(),
                name.clone(),
                ip.clone(),
            );
            async move {
                let req = resolve_id(&s, &resource_type, &api_name, &name, &ip, req).await?;
                let response = search(s, &resource_type, &api_name, &name, &ip, &req, true).await?;
                Ok(search::StreamResponse {
                    payload: Some(search::stream_response::Payload::Response(response)),
                })
            }
        };
        bidirectional_stream(request, self.stream_concurrency, process_fn).await
    }

    #[doc = " A method to linear search indexed vectors by multiple vectors in a single\n request.\n"]
    async fn multi_linear_search(
        &self,
        request: tonic::Request<search::MultiRequest>,
    ) -> Result<tonic::Response<search::Responses>, Status> {
        let mut responses = Vec::new();
        for req in request.into_inner().requests {
            responses.push(
                self.linear_search(tonic::Request::new(req))
                    .await?
                    .into_inner(),
            );
        }
        Ok(tonic::Response::new(search::Responses { responses }))
    }

    #[doc = " A method to linear search indexed vectors by multiple IDs in a single\n request.\n"]
    async fn multi_linear_search_by_id(
        &self,
        request: tonic::Request<search::MultiIdRequest>,
    ) -> Result<tonic::Response<search::Responses>, Status> {
        let mut responses = Vec::new();
        for req in request.into_inner().requests {
            responses.push(
                self.linear_search_by_id(tonic::Request::new(req))
                    .await?
                    .into_inner(),
            );
        }
        Ok(tonic::Response::new(search::Responses { responses }))
    }
}

#[cfg(all(test, feature = "zvec"))]
mod tests {
    use super::*;
    use algorithm::ANN;
    use proto::payload::v1::{Empty, info};
    use proto::vald::v1::{index_server::Index, search_server::Search};

    #[tokio::test]
    async fn zvec_search_modes_and_property() {
        let path = std::env::temp_dir().join(format!(
            "vald-agent-search-{}-{}",
            std::process::id(),
            std::time::SystemTime::now()
                .duration_since(std::time::UNIX_EPOCH)
                .expect("clock")
                .as_nanos()
        ));
        let mut service = zvec::Service::new(path.to_str().expect("UTF-8 path"), 3, vec![], vec![])
            .expect("service");
        service
            .insert("first".into(), vec![0.0, 0.0, 0.0], 1)
            .expect("insert first");
        service
            .insert("second".into(), vec![1.0, 1.0, 1.0], 2)
            .expect("insert second");
        let agent = Arc::new(super::super::Agent::new(
            Box::new(service),
            "test-agent",
            "127.0.0.1",
            "test",
            "test",
            2,
        ));
        let property: info::index::PropertyDetail = agent
            .index_property(tonic::Request::new(Empty {}))
            .await
            .expect("property")
            .into_inner();
        assert_eq!(property.details["test-agent"].dimension, 3);
        assert_eq!(property.details["test-agent"].database_type, "Zvec");
        let listener = tokio::net::TcpListener::bind("127.0.0.1:0")
            .await
            .expect("listener");
        let address = listener.local_addr().expect("address");
        let (shutdown_tx, shutdown_rx) = tokio::sync::oneshot::channel();
        let server = tokio::spawn(
            tonic::transport::Server::builder()
                .add_service(proto::vald::v1::search_server::SearchServer::from_arc(
                    agent.clone(),
                ))
                .serve_with_incoming_shutdown(
                    tokio_stream::wrappers::TcpListenerStream::new(listener),
                    async {
                        shutdown_rx.await.expect("shutdown");
                    },
                ),
        );
        let mut client =
            proto::vald::v1::search_client::SearchClient::connect(format!("http://{address}"))
                .await
                .expect("connect");
        let config = Some(search::Config {
            num: 2,
            ..Default::default()
        });
        let request = search::Request {
            vector: vec![0.0, 0.0, 0.0],
            config: config.clone(),
        };
        let id_request = search::IdRequest {
            id: "first".into(),
            config,
        };
        let check = |response: search::Response| {
            assert_eq!(
                response
                    .results
                    .iter()
                    .map(|hit| hit.id.as_str())
                    .collect::<Vec<_>>(),
                vec!["first", "second"]
            );
        };
        check(
            client
                .linear_search(request.clone())
                .await
                .expect("unary linear")
                .into_inner(),
        );
        check(
            client
                .search_by_id(id_request.clone())
                .await
                .expect("unary by ID")
                .into_inner(),
        );
        check(
            client
                .linear_search_by_id(id_request.clone())
                .await
                .expect("unary linear by ID")
                .into_inner(),
        );
        for response in client
            .multi_linear_search(search::MultiRequest {
                requests: vec![request.clone(); 2],
            })
            .await
            .expect("multi linear")
            .into_inner()
            .responses
        {
            check(response);
        }
        for response in client
            .multi_search_by_id(search::MultiIdRequest {
                requests: vec![id_request.clone(); 2],
            })
            .await
            .expect("multi by ID")
            .into_inner()
            .responses
        {
            check(response);
        }
        for response in client
            .multi_linear_search_by_id(search::MultiIdRequest {
                requests: vec![id_request.clone(); 2],
            })
            .await
            .expect("multi linear by ID")
            .into_inner()
            .responses
        {
            check(response);
        }
        let mut streams = vec![
            client
                .stream_linear_search(tokio_stream::iter(vec![request; 2]))
                .await
                .expect("stream linear")
                .into_inner(),
            client
                .stream_search_by_id(tokio_stream::iter(vec![id_request.clone(); 2]))
                .await
                .expect("stream by ID")
                .into_inner(),
            client
                .stream_linear_search_by_id(tokio_stream::iter(vec![id_request; 2]))
                .await
                .expect("stream linear by ID")
                .into_inner(),
        ];
        for stream in &mut streams {
            let mut count = 0;
            while let Some(response) = stream.message().await.expect("stream message") {
                match response.payload.expect("payload") {
                    search::stream_response::Payload::Response(response) => check(response),
                    search::stream_response::Payload::Status(status) => {
                        panic!("unexpected stream status: {status:?}")
                    }
                }
                count += 1;
            }
            assert_eq!(count, 2);
        }
        shutdown_tx.send(()).expect("send shutdown");
        server.await.expect("server task").expect("server shutdown");
        drop(agent);
        std::fs::remove_dir_all(path).expect("cleanup");
    }
}
