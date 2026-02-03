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
    payload::v1::{object, remove},
    vald::v1::remove_server,
};
use std::sync::Arc;
use tokio::sync::RwLock;
use tonic::{Code, Status};
use tonic_types::StatusExt;

use super::common::{bidirectional_stream, build_error_details};

async fn remove<S: algorithm::ANN>(
    s: Arc<RwLock<S>>,
    resource_type: &str,
    api_name: &str,
    name: &str,
    ip: &str,
    request: &remove::Request,
) -> Result<object::Location, Status> {
    let _config = match request.config.clone() {
        Some(cfg) => cfg,
        None => return Err(Status::invalid_argument("Missing configuration in request")),
    };
    let id = match request.id.clone() {
        Some(id) => id,
        None => return Err(Status::invalid_argument("Missing ID in request")),
    };
    let uuid = id.id;
    {
        let mut s = s.write().await;
        if uuid.len() == 0 {
            let err = Error::InvalidUUID { uuid: uuid.clone() };
            let resource_type = format!("{}/qbg.Remove", resource_type);
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
                format!("Remove API invalid argument for uuid \"{}\" detected", uuid),
                err_details,
            );
            warn!("{:?}", status);
            return Err(status);
        }
        let result = s.remove(uuid.clone()).await;
        match result {
            Err(err) => {
                let resource_type = format!("{}/qbg.Remove", resource_type);
                let resource_name = format!("{}: {}({})", api_name, name, ip);
                let err_msg = err.to_string();
                let mut err_details = build_error_details(
                    err_msg.clone(),
                    &uuid,
                    request.encode_to_vec(),
                    &resource_type,
                    &resource_name,
                    None,
                );
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
                        let status = Status::with_error_details(
                            Code::Aborted,
                            "Remove API aborted to process remove request due to flushing indices is in progress",
                            err_details,
                        );
                        warn!("{:?}", status);
                        status
                    }
                    Error::ObjectIDNotFound { uuid: _ } => {
                        let status = Status::with_error_details(
                            Code::NotFound,
                            format!("Remove API uuid {} not found", uuid),
                            err_details,
                        );
                        warn!("{:?}", status);
                        status
                    }
                    Error::UUIDNotFound { uuid: _ } => {
                        err_details
                            .set_bad_request(vec![tonic_types::FieldViolation::new("id", err_msg)]);
                        let status = Status::with_error_details(
                            Code::InvalidArgument,
                            format!("Remove API invalid argument for uuid \"{}\" detected", uuid),
                            err_details,
                        );
                        warn!("{:?}", status);
                        status
                    }
                    _ => {
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
impl<S: algorithm::ANN + 'static> remove_server::Remove for super::Agent<S> {
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
        request: tonic::Request<remove::TimestampRequest>,
    ) -> std::result::Result<tonic::Response<object::Locations>, tonic::Status> {
        info!("Recieved a request from {:?}", request.remote_addr());
        let req = request.get_ref();
        let timestamps = &req.timestamps;

        let mut locations: Vec<object::Location> = Vec::new();
        let mut errors: Vec<Status> = Vec::new();

        // Build timestamp filter function
        let timestamp_filter = |obj_ts: i64| -> bool {
            for ts in timestamps {
                let op = remove::timestamp::Operator::try_from(ts.operator)
                    .unwrap_or(remove::timestamp::Operator::Eq);
                let matches = match op {
                    remove::timestamp::Operator::Eq => obj_ts == ts.timestamp,
                    remove::timestamp::Operator::Ne => obj_ts != ts.timestamp,
                    remove::timestamp::Operator::Ge => obj_ts >= ts.timestamp,
                    remove::timestamp::Operator::Gt => obj_ts > ts.timestamp,
                    remove::timestamp::Operator::Le => obj_ts <= ts.timestamp,
                    remove::timestamp::Operator::Lt => obj_ts < ts.timestamp,
                };
                if !matches {
                    return false;
                }
            }
            true
        };

        // Collect UUIDs to remove based on timestamp filter
        let uuids_to_remove: Vec<String>;
        {
            let s = self.s.read().await;
            let mut matching_uuids = Vec::new();
            s.list_object_func(|uuid, _vec, ts| {
                if timestamp_filter(ts) {
                    matching_uuids.push(uuid);
                }
                true
            })
            .await;
            uuids_to_remove = matching_uuids;
        }

        // Remove each matching object
        for uuid in uuids_to_remove {
            let remove_req = remove::Request {
                id: Some(object::Id { id: uuid.clone() }),
                config: None,
            };
            match remove(
                self.s.clone(),
                &self.resource_type,
                &self.api_name,
                &self.name,
                &self.ip,
                &remove_req,
            )
            .await
            {
                Ok(loc) => locations.push(loc),
                Err(e) => errors.push(e),
            }
        }

        if !errors.is_empty() && locations.is_empty() {
            // All removals failed
            return Err(errors.into_iter().next().unwrap());
        }

        if locations.is_empty() {
            let resource_type = format!("{}/qbg.RemoveByTimestamp", self.resource_type);
            let resource_name = format!("{}: {}({})", self.api_name, self.name, self.ip);
            let err = Error::IndexNotFound {};
            let err_details = build_error_details(
                err,
                "",
                req.encode_to_vec(),
                &resource_type,
                &resource_name,
                None,
            );
            let status = Status::with_error_details(
                Code::NotFound,
                "RemoveByTimestamp API remove target not found",
                err_details,
            );
            error!("{:?}", status);
            return Err(status);
        }

        Ok(tonic::Response::new(object::Locations { locations }))
    }

    #[doc = " Server streaming response type for the StreamRemove method."]
    type StreamRemoveStream = crate::stream_type!(object::StreamLocation);

    #[doc = " A method to remove multiple indexed vectors by bidirectional streaming.\n"]
    async fn stream_remove(
        &self,
        request: tonic::Request<tonic::Streaming<remove::Request>>,
    ) -> std::result::Result<tonic::Response<Self::StreamRemoveStream>, tonic::Status> {
        info!(
            "Received stream remove request from {:?}",
            request.remote_addr()
        );

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
                    Ok(location) => Ok(object::StreamLocation {
                        payload: Some(object::stream_location::Payload::Location(location)),
                    }),
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
            let result = s.remove_multiple(uuids.clone()).await;
            match result {
                Err(err) => {
                    let resource_type = self.resource_type.clone() + "/qbg.MultiRemove";
                    let resource_name = format!("{}: {}({})", self.api_name, self.name, self.ip);
                    let request_bytes = mreq.encode_to_vec();
                    let status = match err {
                        Error::FlushingIsInProgress {} => {
                            let err_details = build_error_details(
                                err,
                                &uuids.join(","),
                                request_bytes,
                                &resource_type,
                                &resource_name,
                                None,
                            );
                            let status = Status::with_error_details(
                                Code::Aborted,
                                "MultiRemove API aborted to process remove request due to flushing indices is in progress",
                                err_details,
                            );
                            warn!("{:?}", status);
                            status
                        }
                        Error::ObjectIDNotFound { ref uuid } => {
                            let err_details = build_error_details(
                                &err,
                                uuid,
                                request_bytes,
                                &resource_type,
                                &resource_name,
                                None,
                            );
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
                            let err_details = build_error_details(
                                err,
                                &uuids.join(","),
                                request_bytes,
                                &resource_type,
                                &resource_name,
                                Some("uuid"),
                            );
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
                            let err_details = build_error_details(
                                err,
                                &uuids.join(","),
                                request_bytes,
                                &resource_type,
                                &resource_name,
                                None,
                            );
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
