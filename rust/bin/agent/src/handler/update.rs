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
    payload::v1::{object, update},
    vald::v1::update_server,
};
use std::{collections::HashMap, sync::Arc};
use tokio::sync::RwLock;
use tonic::{Code, Status};
use tonic_types::StatusExt;

use super::common::{bidirectional_stream, build_error_details};

pub(crate) async fn update<S: algorithm::ANN>(
    s: Arc<RwLock<S>>,
    resource_type: &str,
    api_name: &str,
    name: &str,
    ip: &str,
    request: &update::Request,
) -> Result<object::Location, Status> {
    let _config = match request.config.clone() {
        Some(cfg) => cfg,
        None => return Err(Status::invalid_argument("Missing configuration in request")),
    };
    let mut s = s.write().await;
    let vec = match request.vector.clone() {
        Some(v) => v,
        None => return Err(Status::invalid_argument("Missing vector in request")),
    };
    let uuid = vec.id.clone();
    if vec.vector.len() != s.get_dimension_size() {
        let err = Error::IncompatibleDimensionSize {
            got: vec.vector.len(),
            want: s.get_dimension_size(),
        };
        let resource_type = format!("{}/qbg.Update", resource_type);
        let resource_name = format!("{}: {}({})", api_name, name, ip);
        let err_details = build_error_details(
            err,
            &uuid,
            request.encode_to_vec(),
            &resource_type,
            &resource_name,
            Some("vector dimension size"),
        );
        let status = Status::with_error_details(
            Code::InvalidArgument,
            "Update API Incompatible Dimension Size detected",
            err_details,
        );
        warn!("{:?}", status);
        return Err(status);
    }
    if uuid.len() == 0 {
        let err = Error::InvalidUUID { uuid: uuid.clone() };
        let resource_type = format!("{}/qbg.Update", resource_type);
        let resource_name = format!("{}: {}({})", api_name, name, ip);
        let err_details = build_error_details(
            err,
            &uuid,
            request.encode_to_vec(),
            &resource_type,
            &resource_name,
            Some("uuid"),
        );
        let status = Status::with_error_details(
            Code::InvalidArgument,
            format!("Update API invalid argument for uuid \"{}\" detected", uuid),
            err_details,
        );
        warn!("{:?}", status);
        return Err(status);
    }
    let result = s.update(uuid.clone(), vec.vector.clone()).await;
    match result {
        Err(err) => {
            let resource_type = format!("{}/qbg.Update", resource_type);
            let resource_name = format!("{}: {}({})", api_name, name, ip);
            let request_bytes = request.encode_to_vec();
            let status = match err {
                Error::FlushingIsInProgress {} => {
                    let err_details = build_error_details(
                        err,
                        &uuid,
                        request_bytes,
                        &resource_type,
                        &resource_name,
                        None,
                    );
                    let status = Status::with_error_details(Code::Aborted, "Update API aborted to process update request due to flushing indices is in progress", err_details);
                    warn!("{:?}", status);
                    status
                }
                Error::ObjectIDNotFound { uuid: _ } => {
                    let err_details = build_error_details(
                        err,
                        &uuid,
                        request_bytes,
                        &resource_type,
                        &resource_name,
                        None,
                    );
                    let status = Status::with_error_details(
                        Code::NotFound,
                        format!("Update API uuid {} not found", uuid),
                        err_details,
                    );
                    warn!("{:?}", status);
                    status
                }
                Error::UUIDNotFound { uuid: _ } => {
                    let err_details = build_error_details(
                        err,
                        &uuid,
                        request_bytes,
                        &resource_type,
                        &resource_name,
                        Some("uuid or vector"),
                    );
                    let status = Status::with_error_details(
                        Code::InvalidArgument,
                        format!(
                            "Update API invalid argument for uuid \"{}\" vec \"{:?}\" detected",
                            uuid, vec.vector
                        ),
                        err_details,
                    );
                    warn!("{:?}", status);
                    status
                }
                Error::UUIDAlreadyExists { uuid: _ } => {
                    let err_details = build_error_details(
                        err,
                        &uuid,
                        request_bytes,
                        &resource_type,
                        &resource_name,
                        None,
                    );
                    let status = Status::with_error_details(
                        Code::AlreadyExists,
                        format!("Update API uuid {}'s same data already exists", uuid),
                        err_details,
                    );
                    warn!("{:?}", status);
                    status
                }
                _ => {
                    let err_details = build_error_details(
                        err,
                        &uuid,
                        request_bytes,
                        &resource_type,
                        &resource_name,
                        None,
                    );
                    let status = Status::with_error_details(
                        Code::Internal,
                        "Update API failed",
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

#[tonic::async_trait]
impl<S: algorithm::ANN + 'static> update_server::Update for super::Agent<S> {
    async fn update(
        &self,
        request: tonic::Request<update::Request>,
    ) -> std::result::Result<tonic::Response<object::Location>, tonic::Status> {
        info!("Recieved a request from {:?}", request.remote_addr());
        let request = request.get_ref();
        let s = self.s.clone();
        let resource_type = self.resource_type.clone();
        let name = self.name.clone();
        let ip = self.ip.clone();
        let api_name = self.api_name.clone();
        match update(s, &resource_type, &api_name, &name, &ip, request).await {
            Ok(location) => Ok(tonic::Response::new(location)),
            Err(e) => Err(e),
        }
    }

    #[doc = " Server streaming response type for the StreamUpdate method."]
    type StreamUpdateStream = crate::stream_type!(object::StreamLocation);

    #[doc = " A method to update multiple indexed vectors by bidirectional streaming.\n"]
    async fn stream_update(
        &self,
        request: tonic::Request<tonic::Streaming<update::Request>>,
    ) -> std::result::Result<tonic::Response<Self::StreamUpdateStream>, tonic::Status> {
        info!(
            "Received stream update request from {:?}",
            request.remote_addr()
        );

        let s = self.s.clone();
        let resource_type = self.resource_type.clone() + "/qbg.StreamUpdate";
        let name = self.name.clone();
        let ip = self.ip.clone();
        let api_name = self.api_name.clone();

        let process_fn = move |req: update::Request| {
            let s = s.clone();
            let resource_type = resource_type.clone();
            let name = name.clone();
            let ip = ip.clone();
            let api_name = api_name.clone();
            async move {
                match update(s, &resource_type, &api_name, &name, &ip, &req).await {
                    Ok(location) => Ok(object::StreamLocation {
                        payload: Some(object::stream_location::Payload::Location(location)),
                    }),
                    Err(status) => Err(status),
                }
            }
        };

        bidirectional_stream(request, self.stream_concurrency, process_fn).await
    }

    #[doc = " A method to update multiple indexed vectors in a single request.\n"]
    async fn multi_update(
        &self,
        request: tonic::Request<update::MultiRequest>,
    ) -> std::result::Result<tonic::Response<object::Locations>, tonic::Status> {
        info!("Recieved a request from {:?}", request.remote_addr());
        let mreq = request.get_ref();
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
                    let resource_type = self.resource_type.clone() + "/qbg.MultiUpdate";
                    let resource_name = format!("{}: {}({})", self.api_name, self.name, self.ip);
                    let err_details = build_error_details(
                        err,
                        &vec.id,
                        mreq.encode_to_vec(),
                        &resource_type,
                        &resource_name,
                        Some("vector dimension size"),
                    );
                    let status = Status::with_error_details(
                        Code::InvalidArgument,
                        "MultiUpdate API Incombatible Dimension Size detedted",
                        err_details,
                    );
                    warn!("{:?}", status);
                    return Err(status);
                }
                uuids.push(vec.id.clone());
                vmap.insert(vec.id, vec.vector);
            }
            let result = s.update_multiple(vmap).await;
            match result {
                Err(err) => {
                    let resource_type = self.resource_type.clone() + "/qbg.MultiUpdate";
                    let resource_name = format!("{}: {}({})", self.api_name, self.name, self.ip);
                    let request_bytes = mreq.encode_to_vec();
                    let status = match err {
                        Error::FlushingIsInProgress {} => {
                            let err_details = build_error_details(
                                err,
                                &uuids.join(", "),
                                request_bytes,
                                &resource_type,
                                &resource_name,
                                None,
                            );
                            let status = Status::with_error_details(
                                Code::Aborted,
                                "MultiUpdate API aborted to process update request due to flushing indices is in progress",
                                err_details,
                            );
                            warn!("{:?}", status);
                            status
                        }
                        Error::ObjectIDNotFound { ref uuid } => {
                            let err_details = build_error_details(
                                &err,
                                &uuid,
                                request_bytes,
                                &resource_type,
                                &resource_name,
                                None,
                            );
                            let uuids = Error::split_uuids(uuid.to_string());
                            let status = Status::with_error_details(
                                Code::NotFound,
                                format!("MultiUpdate API uuids {:?} not found", uuids),
                                err_details,
                            );
                            warn!("{:?}", status);
                            status
                        }
                        Error::InvalidDimensionSize {
                            current: _,
                            limit: _,
                        } => {
                            let err_details = build_error_details(
                                &err,
                                &uuids.join(","),
                                request_bytes,
                                &resource_type,
                                &resource_name,
                                Some("vector dimension"),
                            );
                            let status = Status::with_error_details(
                                Code::InvalidArgument,
                                format!("MultiUpdate API invalid dimension size detected"),
                                err_details,
                            );
                            warn!("{:?}", status);
                            status
                        }
                        Error::UUIDNotFound { ref uuid } => {
                            let err_details = build_error_details(
                                &err,
                                &uuid,
                                request_bytes,
                                &resource_type,
                                &resource_name,
                                Some("uuid"),
                            );
                            let uuids = Error::split_uuids(uuid.to_string());
                            let status = Status::with_error_details(
                                Code::InvalidArgument,
                                format!(
                                    "MultiUpdate API invalid argument for uuids {:?} detected",
                                    uuids
                                ),
                                err_details,
                            );
                            warn!("{:?}", status);
                            status
                        }
                        Error::UUIDAlreadyExists { ref uuid } => {
                            let err_details = build_error_details(
                                &err,
                                &uuid,
                                request_bytes,
                                &resource_type,
                                &resource_name,
                                None,
                            );
                            let uuids = Error::split_uuids(uuid.to_string());
                            let status = Status::with_error_details(
                                Code::AlreadyExists,
                                format!("MultiUpdate API uuids {:?} already exists", uuids),
                                err_details,
                            );
                            warn!("{:?}", status);
                            status
                        }
                        _ => {
                            let err_details = build_error_details(
                                err,
                                &uuids.join(", "),
                                request_bytes,
                                &resource_type,
                                &resource_name,
                                None,
                            );
                            let status = Status::with_error_details(
                                Code::Internal,
                                "MultiUpdate API failed",
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

    #[doc = " A method to update timestamp indexed vectors in a single request.\n"]
    async fn update_timestamp(
        &self,
        request: tonic::Request<update::TimestampRequest>,
    ) -> std::result::Result<tonic::Response<object::Location>, tonic::Status> {
        info!("Recieved a request from {:?}", request.remote_addr());
        let req = request.get_ref();
        let uuid = &req.id;
        let ts = req.timestamp;
        let force = req.force;
        let resource_type = format!("{}/qbg.UpdateTimestamp", self.resource_type);
        let resource_name = format!("{}: {}({})", self.api_name, self.name, self.ip);

        if uuid.is_empty() {
            let err = Error::InvalidUUID { uuid: uuid.clone() };
            let err_details = build_error_details(
                err,
                uuid,
                req.encode_to_vec(),
                &resource_type,
                &resource_name,
                Some("uuid"),
            );
            let status = Status::with_error_details(
                Code::InvalidArgument,
                "UpdateTimestamp API invalid uuid",
                err_details,
            );
            warn!("{:?}", status);
            return Err(status);
        }

        if !force && ts < 0 {
            let err = Error::InvalidTimestamp { timestamp: ts };
            let err_details = build_error_details(
                err,
                uuid,
                req.encode_to_vec(),
                &resource_type,
                &resource_name,
                Some("timestamp"),
            );
            let status = Status::with_error_details(
                Code::InvalidArgument,
                "UpdateTimestamp API invalid vector argument",
                err_details,
            );
            warn!("{:?}", status);
            return Err(status);
        }

        let mut s = self.s.write().await;
        match s.update_timestamp(uuid.clone(), ts, force).await {
            Err(err) => {
                let status = match &err {
                    Error::FlushingIsInProgress {} => {
                        let err_details = build_error_details(
                            err,
                            uuid,
                            req.encode_to_vec(),
                            &resource_type,
                            &resource_name,
                            None,
                        );
                        let status = Status::with_error_details(
                            Code::Aborted,
                            "UpdateTimestamp API aborted to process update request due to flushing indices is in progress",
                            err_details,
                        );
                        warn!("{:?}", status);
                        status
                    }
                    Error::ObjectIDNotFound { uuid: _ } => {
                        let err_details = build_error_details(
                            err,
                            uuid,
                            req.encode_to_vec(),
                            &resource_type,
                            &resource_name,
                            None,
                        );
                        let status = Status::with_error_details(
                            Code::NotFound,
                            format!("UpdateTimestamp API uuid {}'s data not found", uuid),
                            err_details,
                        );
                        warn!("{:?}", status);
                        status
                    }
                    Error::NewerTimestampAlreadyExists { uuid: _, timestamp: _ } => {
                        let err_details = build_error_details(
                            err,
                            uuid,
                            req.encode_to_vec(),
                            &resource_type,
                            &resource_name,
                            None,
                        );
                        let status = Status::with_error_details(
                            Code::AlreadyExists,
                            format!("UpdateTimestamp API uuid {}'s newer timestamp already exists", uuid),
                            err_details,
                        );
                        warn!("{:?}", status);
                        status
                    }
                    _ => {
                        let err_details = build_error_details(
                            err,
                            uuid,
                            req.encode_to_vec(),
                            &resource_type,
                            &resource_name,
                            None,
                        );
                        let status = Status::with_error_details(
                            Code::Internal,
                            "UpdateTimestamp API failed",
                            err_details,
                        );
                        error!("{:?}", status);
                        status
                    }
                };
                Err(status)
            }
            Ok(()) => Ok(tonic::Response::new(object::Location {
                name: self.name.clone(),
                uuid: uuid.clone(),
                ips: vec![self.ip.clone()],
            })),
        }
    }
}
