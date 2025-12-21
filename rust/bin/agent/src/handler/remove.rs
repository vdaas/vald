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
use algorithm::{Error, MultiError};
use log::{error, info, warn};
use prost::Message;
use proto::{
    payload::v1::{object, remove},
    vald::v1::remove_server,
};
use tokio::sync::RwLock;
use std::sync::Arc;
use tonic::{Code, Status};
use tonic_types::StatusExt;

use super::common::{bidirectional_stream, build_error_details};

async fn remove(
    s: Arc<RwLock<dyn algorithm::ANN>>,
    resource_type: &str,
    api_name: &str,
    name: &str,
    ip: &str,
    request: &remove::Request
) -> Result<object::Location, Status> {
    let config = match request.config.clone() {
        Some(cfg) => cfg,
        None => return Err(Status::invalid_argument("Missing configuration in request")),
    };
    let id = match request.id.clone() {
        Some(id) => id,
        None => return Err(Status::invalid_argument("Missing ID in request")),
    };
    let uuid = id.id;
    let hostname = cargo::util::hostname()?;
    let domain = hostname.to_str().unwrap();
    {
        let mut s = s.write().await;
        if uuid.len() == 0 {
            let err = Error::InvalidUUID { uuid: uuid.clone() };
            let resource_type = format!("{}/qbg.Remove", resource_type);
            let resource_name = format!("{}: {}({})", api_name, name, ip);
            let err_details = build_error_details(err, domain, &uuid, request.encode_to_vec(), &resource_type, &resource_name, Some("uuid"));
            let status = Status::with_error_details(
                Code::InvalidArgument,
                format!("Remove API invalid argument for uuid \"{}\" detected", uuid),
                err_details,
            );
            warn!("{:?}", status);
            return Err(status);
        }
        let result = s.remove(uuid.clone(), config.timestamp);
        match result {
            Err(err) => {
                let resource_type = format!("{}/qbg.Remove", resource_type);
                let resource_name = format!("{}: {}({})", api_name, name, ip);
                let request_bytes = request.encode_to_vec();
                let status = match err {
                    Error::FlushingIsInProgress {} => {
                        let err_details = build_error_details(err, domain, &uuid, request_bytes, &resource_type, &resource_name, None);
                        let status = Status::with_error_details(Code::Aborted, "Remove API aborted to process remove request due to flushing indices is in progress", err_details);
                        warn!("{:?}", status);
                        status
                    }
                    Error::ObjectIDNotFound { uuid: _ } => {
                        let err_details = build_error_details(err, domain, &uuid, request_bytes, &resource_type, &resource_name, None);
                        let status = Status::with_error_details(
                            Code::NotFound,
                            format!("Remove API uuid {} not found", uuid),
                            err_details,
                        );
                        warn!("{:?}", status);
                        status
                    }
                    Error::UUIDNotFound { uuid: _ } => {
                        let err_details = build_error_details(err, domain, &uuid, request_bytes, &resource_type, &resource_name, Some("uuid"));
                        let status = Status::with_error_details(
                            Code::InvalidArgument,
                            format!(
                                "Remove API invalid argument for uuid \"{}\" detected",
                                uuid
                            ),
                            err_details,
                        );
                        warn!("{:?}", status);
                        status
                    }
                    _ => {
                        let err_details = build_error_details(err, domain, &uuid, request_bytes, &resource_type, &resource_name, None);
                        let status = Status::with_error_details(
                            Code::Internal,
                            "Remove API failed",
                            err_details,
                        );
                        error!("{:?}", status);
                        status
                    }
                };
                Err(status)
            }
            Ok(()) => Ok(object::Location {
                name: name.to_owned(),
                uuid: uuid,
                ips: vec![ip.to_owned()],
            }),
        }
    }
}

#[tonic::async_trait]
impl remove_server::Remove for super::Agent {
    async fn remove(
        &self,
        request: tonic::Request<remove::Request>,
    ) -> std::result::Result<tonic::Response<object::Location>, tonic::Status> {
        info!("Recieved a request from {:?}", request.remote_addr());
        let request = request.get_ref();
        let s = self.s.clone();
        let resource_type = self.resource_type.clone();
        let name = self.name.clone();
        let ip = self.ip.clone();
        let api_name = self.api_name.clone();
        match remove(s, &resource_type, &api_name, &name, &ip, request).await {
            Ok(location) => Ok(tonic::Response::new(location)),
            Err(e) => Err(e),
        }
    }

    #[doc = " A method to remove an indexed vector based on timestamp.\n"]
    async fn remove_by_timestamp(
        &self,
        _request: tonic::Request<remove::TimestampRequest>,
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
        info!("Received stream remove request from {:?}", request.remote_addr());

        let s = self.s.clone();
        let resource_type = self.resource_type.clone() + "/qbg.StreamRemove";
        let name = self.name.clone();
        let ip = self.ip.clone();
        let api_name = self.api_name.clone();

        let process_fn = move |req: remove::Request| {
            let s = s.clone();
            let resource_type = resource_type.clone();
            let name = name.clone();
            let ip = ip.clone();
            let api_name = api_name.clone();
            async move {
                match remove(s, &resource_type, &api_name, &name, &ip, &req).await {
                    Ok(location) => {
                        Ok(object::StreamLocation {
                            payload: Some(object::stream_location::Payload::Location(location)),
                        })
                    }
                    Err(status) => Err(status),
                }
            }
        };

        bidirectional_stream(request, self.stream_concurrency, process_fn).await
    }

    #[doc = " A method to remove multiple indexed vectors in a single request.\n"]
    async fn multi_remove(
        &self,
        request: tonic::Request<remove::MultiRequest>,
    ) -> std::result::Result<tonic::Response<object::Locations>, tonic::Status> {
        info!("Recieved a request from {:?}", request.remote_addr());
        let mreq = request.get_ref();
        let hostname = cargo::util::hostname()?;
        let domain = hostname.to_str().unwrap();
        let uuids: Vec<String> = mreq
            .requests
            .clone()
            .into_iter()
            .filter_map(|x| match x.id {
                Some(id) => Some(id.id),
                None => None,
            })
            .collect();
        {
            let mut s = self.s.write().await;
            let result = s.remove_multiple(uuids.clone());
            match result {
                Err(err) => {
                    let resource_type = self.resource_type.clone() + "/qbg.MultiRemove";
                    let resource_name = format!("{}: {}({})", self.api_name, self.name, self.ip);
                    let request_bytes = mreq.encode_to_vec();
                    let status = match err {
                        Error::FlushingIsInProgress {} => {
                            let err_details = build_error_details(err, domain, &uuids.join(","), request_bytes, &resource_type, &resource_name, None);
                            let status = Status::with_error_details(Code::Aborted, "MultiRemove API aborted to process remove request due to flushing indices is in progress", err_details);
                            warn!("{:?}", status);
                            status
                        }
                        Error::ObjectIDNotFound { ref uuid } => {
                            let err_details = build_error_details(&err, domain, uuid, request_bytes, &resource_type, &resource_name, None);
                            let uuids = Error::split_uuids(uuid.to_string());
                            let status = Status::with_error_details(
                                Code::NotFound,
                                format!("MultiRemove API uuids {:?} not found", uuids),
                                err_details,
                            );
                            warn!("{:?}", status);
                            status
                        }
                        Error::UUIDNotFound { uuid: _ } => {
                            let err_details = build_error_details(err, domain, &uuids.join(","), request_bytes, &resource_type, &resource_name, Some("uuid"));
                            let status = Status::with_error_details(
                                Code::InvalidArgument,
                                format!(
                                    "MultiRemove API invalid argument for uuids \"{:?}\" detected",
                                    uuids
                                ),
                                err_details,
                            );
                            warn!("{:?}", status);
                            status
                        }
                        _ => {
                            let err_details = build_error_details(err, domain, &uuids.join(","), request_bytes, &resource_type, &resource_name, None);
                            let status = Status::with_error_details(
                                Code::Internal,
                                "MultiRemove API failed",
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
