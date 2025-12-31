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
use algorithm::{Error, MultiError};
use log::{error, info, warn};
use prost::Message;
use proto::{
    payload::v1::{insert, object},
    vald::v1::insert_server,
};
use tokio::sync::RwLock;
use std::{collections::HashMap, string::String, sync::Arc};
use tonic::{Code, Status};
use tonic_types::StatusExt;

use super::common::{build_error_details, bidirectional_stream};

pub(super) async fn insert(
    s: Arc<RwLock<dyn algorithm::ANN>>,
    resource_type: &str,
    api_name: &str,
    name: &str,
    ip: &str,
    request: &insert::Request,
) -> Result<object::Location, Status> {
    let config = match request.config.clone() {
        Some(cfg) => cfg,
        None => return Err(Status::invalid_argument("Missing configuration in request")),
    };
    let hostname = cargo::util::hostname()?;
    let domain = hostname.to_str().unwrap();
    {
        let mut s = s.write().await;
        let vec = match request.vector.clone() {
            Some(v) => v,
            None => return Err(Status::invalid_argument("Missing vector in request")),
        };
        if vec.vector.len() != s.get_dimension_size() {
            let err = Error::IncompatibleDimensionSize {
                got: vec.vector.len(),
                want: s.get_dimension_size(),
            };
            let resource_type = format!("{}/qbg.Insert", resource_type);
            let resource_name = format!("{}: {}({})", api_name, name, ip);
            let err_details = build_error_details(err, domain, &vec.id, request.encode_to_vec(), &resource_type, &resource_name, Some("vector dimension size"));
            let status = Status::with_error_details(
                Code::InvalidArgument,
                "Insert API Incombatible Dimension Size detedted",
                err_details,
            );
            warn!("{:?}", status);
            return Err(status);
        }
        let result = s.insert(vec.id.clone(), vec.vector.clone(), config.timestamp);
        match result {
            Err(err) => {
                let resource_type = format!("{}/qbg.Insert", resource_type);
                let resource_name = format!("{}: {}({})", api_name, name, ip);
                let request_bytes = request.encode_to_vec();
                let status = match err {
                    Error::FlushingIsInProgress {} => {
                        let err_details = build_error_details(err, domain, &vec.id, request_bytes, &resource_type, &resource_name, None);
                        let status = Status::with_error_details(Code::Aborted, "Insert API aborted to process insert request due to flushing indices is in progress", err_details);
                        warn!("{:?}", status);
                        status
                    }
                    Error::UUIDAlreadyExists { uuid: _ } => {
                        let err_details = build_error_details(err, domain, &vec.id, request_bytes, &resource_type, &resource_name, None);
                        let status = Status::with_error_details(
                            Code::AlreadyExists,
                            format!("Insert API uuid {} already exists", vec.id),
                            err_details,
                        );
                        warn!("{:?}", status);
                        status
                    }
                    Error::UUIDNotFound { uuid: _ } => {
                        let err_details = build_error_details(err, domain, &vec.id, request_bytes, &resource_type, &resource_name, Some("uuid"));
                        let status = Status::with_error_details(
                            Code::InvalidArgument,
                            format!(
                                "Insert API invalid id: \"{}\" or vector: {:?} was given",
                                vec.id, vec.vector
                            ),
                            err_details,
                        );
                        warn!("{:?}", status);
                        status
                    }
                    _ => {
                        let err_details = build_error_details(err, domain, &vec.id, request_bytes, &resource_type, &resource_name, None);
                        Status::with_error_details(
                            Code::Unknown,
                            "failed to parse Insert gRPC error response",
                            err_details,
                        )
                    }
                };
                Err(status)
            }
            Ok(()) => Ok(object::Location {
                name: name.to_owned(),
                uuid: vec.id,
                ips: vec![ip.to_owned()],
            }),
        }
    }
}

#[tonic::async_trait]
impl insert_server::Insert for super::Agent {
    async fn insert(
        &self,
        request: tonic::Request<insert::Request>,
    ) -> std::result::Result<tonic::Response<object::Location>, tonic::Status> {
        info!("Recieved a request from {:?}", request.remote_addr());
        let request = request.get_ref();
        let s = self.s.clone();
        let resource_type = self.resource_type.clone();
        let name = self.name.clone();
        let ip = self.ip.clone();
        let api_name = self.api_name.clone();
        match insert(s, &resource_type, &api_name, &name, &ip, request).await {
            Ok(location) => Ok(tonic::Response::new(location)),
            Err(e) => Err(e),
        }
    }

    #[doc = " Server streaming response type for the StreamInsert method."]
    type StreamInsertStream = crate::stream_type!(object::StreamLocation);

    #[doc = " A method to add new multiple vectors by bidirectional streaming.\n"]
    async fn stream_insert(
        &self,
        request: tonic::Request<tonic::Streaming<insert::Request>>,
    ) -> std::result::Result<tonic::Response<Self::StreamInsertStream>, tonic::Status> {
        info!("Received stream insert request from {:?}", request.remote_addr());
        let s = self.s.clone();
        let resource_type = format!("{}/qbg.StreamInsert", self.resource_type.clone());
        let name = self.name.clone();
        let ip = self.ip.clone();
        let api_name = self.api_name.clone();

        let process_fn = move |req: insert::Request| {
            let s = s.clone();
            let resource_type = resource_type.clone();
            let name = name.clone();
            let ip = ip.clone();
            let api_name = api_name.clone();
            async move {
                match insert(s, &resource_type, &api_name, &name, &ip, &req).await {
                    Ok(response) => {
                        Ok(object::StreamLocation {
                            payload: Some(object::stream_location::Payload::Location(response)),
                        })
                    }
                    Err(status) => Err(status),
                }
            }
        };

        bidirectional_stream(request, self.stream_concurrency, process_fn).await
    }

    #[doc = " A method to add new multiple vectors in a single request.\n"]
    async fn multi_insert(
        &self,
        request: tonic::Request<insert::MultiRequest>,
    ) -> std::result::Result<tonic::Response<object::Locations>, tonic::Status> {
        info!("Recieved a request from {:?}", request.remote_addr());
        let mreq = request.get_ref();
        let hostname = cargo::util::hostname()?;
        let domain = hostname.to_str().unwrap();
        let mut uuids: Vec<String> = Vec::new();
        let mut vmap = HashMap::new();
        {
            let mut s = self.s.write().await;
            for req in mreq.requests.clone() {
                let vec = match req.vector.clone() {
                    Some(v) => v,
                    None => return Err(Status::invalid_argument("Missing vector in request")),
                };
                if vec.vector.len() != s.get_dimension_size() {
                    let err = Error::IncompatibleDimensionSize {
                        got: vec.vector.len(),
                        want: s.get_dimension_size(),
                    };
                    let resource_type = format!("{}/qbg.MultiInsert", self.resource_type);
                    let resource_name = format!("{}: {}({})", self.api_name, self.name, self.ip);
                    let err_details = build_error_details(err, domain, &vec.id, mreq.encode_to_vec(), &resource_type, &resource_name, Some("vector dimension size"));
                    let status = Status::with_error_details(
                        Code::InvalidArgument,
                        "MultiInsert API Incombatible Dimension Size detedted",
                        err_details,
                    );
                    warn!("{:?}", status);
                    return Err(status);
                }
                uuids.push(vec.id.clone());
                vmap.insert(vec.id, vec.vector);
            }
            let result = s.insert_multiple(vmap);
            match result {
                Err(err) => {
                    let resource_type = format!("{}/qbg.MultiInsert", self.resource_type);
                    let resource_name = format!("{}: {}({})", self.api_name, self.name, self.ip);
                    let request_bytes = mreq.encode_to_vec();
                    let status = match err {
                        Error::FlushingIsInProgress {} => {
                            let err_details = build_error_details(err, domain, &uuids.join(", "), request_bytes, &resource_type, &resource_name, None);
                            let status = Status::with_error_details(Code::Aborted, "MultiInsert API aborted to process insert request due to flushing indices is in progress", err_details);
                            warn!("{:?}", status);
                            status
                        }
                        Error::UUIDAlreadyExists { ref uuid } => {
                            let err_details = build_error_details(&err, domain, uuid, request_bytes, &resource_type, &resource_name, None);
                            let uuids = Error::split_uuids(uuid.to_string());
                            let status = Status::with_error_details(
                                Code::AlreadyExists,
                                format!("MultiInsert API uuids {:?} already exists", uuids),
                                err_details,
                            );
                            warn!("{:?}", status);
                            status
                        }
                        Error::UUIDNotFound { uuid: _ } => {
                            let err_details = build_error_details(err, domain, &uuids.join(", "), request_bytes, &resource_type, &resource_name, Some("uuid"));
                            let status = Status::with_error_details(
                                Code::InvalidArgument,
                                format!("MultiInsert API invalid uuids \"{:?}\" detected", uuids),
                                err_details,
                            );
                            warn!("{:?}", status);
                            status
                        }
                        _ => {
                            let err_details = build_error_details(err, domain, &uuids.join(", "), request_bytes, &resource_type, &resource_name, None);
                            let status = Status::with_error_details(
                                Code::Internal,
                                "MultiInsert API failed",
                                err_details,
                            );
                            error!("{:?}", status);
                            status
                        }
                    };
                    Err(status)
                }
                Ok(()) => Ok(tonic::Response::new(object::Locations {
                    locations: uuids
                        .iter()
                        .map(|x| object::Location {
                            name: self.name.clone(),
                            uuid: x.to_string(),
                            ips: vec![self.ip.clone()],
                        })
                        .collect(),
                })),
            }
        }
    }
}
