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
use prost::Message;
use proto::{
    payload::v1::{insert, object},
    vald::v1::insert_server,
};
use std::{collections::HashMap, string::String};
use tonic::{Code, Status};
use tonic_types::{ErrorDetails, FieldViolation, StatusExt};

#[tonic::async_trait]
impl insert_server::Insert for super::Agent {
    async fn insert(
        &self,
        request: tonic::Request<insert::Request>,
    ) -> std::result::Result<tonic::Response<object::Location>, tonic::Status> {
        println!("Recieved a request from {:?}", request.remote_addr());
        let req = request.get_ref();
        let config = match req.config.clone() {
            Some(cfg) => cfg,
            None => return Err(Status::invalid_argument("Missing configuration in request")),
        };
        let hostname = cargo::util::hostname()?;
        let domain = hostname.to_str().unwrap();
        {
            let mut s = self.s.write().await;
            let vec = match req.vector.clone() {
                Some(v) => v,
                None => return Err(Status::invalid_argument("Missing vector in request")),
            };
            if vec.vector.len() != s.get_dimension_size() {
                let err = Error::IncompatibleDimensionSize {
                    got: vec.vector.len(),
                    want: s.get_dimension_size(),
                };
                let mut err_details = ErrorDetails::new();
                let metadata = HashMap::new();
                let resource_type = self.resource_type.clone() + "/qbg.Insert";
                let resource_name = format!("{}: {}({})", self.api_name, self.name, self.ip);
                err_details.set_error_info(err.to_string(), domain, metadata);
                err_details.set_request_info(
                    vec.id,
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
                    "Insert API Incombatible Dimension Size detedted",
                    err_details,
                );
                return Err(status);
            }
            let result = s.insert(vec.id.clone(), vec.vector.clone(), config.timestamp);
            match result {
                Err(err) => {
                    let metadata = HashMap::new();
                    let resource_type = self.resource_type.clone() + "/qbg.Insert";
                    let resource_name = format!("{}: {}({})", self.api_name, self.name, self.ip);
                    let status = match err {
                        Error::FlushingIsInProgress {} => {
                            let mut err_details = ErrorDetails::new();
                            err_details.set_error_info(err.to_string(), domain, metadata);
                            err_details.set_request_info(
                                vec.id,
                                String::from_utf8(req.encode_to_vec())
                                    .unwrap_or_else(|_| "<invalid UTF-8>".to_string()),
                            );
                            err_details.set_resource_info(resource_type, resource_name, "", "");
                            Status::with_error_details(Code::Aborted, "Insert API aborted to process insert request due to flushing indices is in progress", err_details)
                        }
                        Error::UUIDAlreadyExists { uuid: _ } => {
                            let mut err_details = ErrorDetails::new();
                            err_details.set_error_info(err.to_string(), domain, metadata);
                            err_details.set_request_info(
                                vec.id.clone(),
                                String::from_utf8(req.encode_to_vec())
                                    .unwrap_or_else(|_| "<invalid UTF-8>".to_string()),
                            );
                            err_details.set_resource_info(resource_type, resource_name, "", "");
                            Status::with_error_details(
                                Code::AlreadyExists,
                                format!("Insert API uuid {} already exists", vec.id),
                                err_details,
                            )
                        }
                        Error::UUIDNotFound { id: _ } => {
                            let mut err_details = ErrorDetails::new();
                            err_details.set_error_info(err.to_string(), domain, metadata);
                            err_details.set_request_info(
                                vec.id.clone(),
                                String::from_utf8(req.encode_to_vec())
                                    .unwrap_or_else(|_| "<invalid UTF-8>".to_string()),
                            );
                            err_details.set_bad_request(vec![FieldViolation::new(
                                "uuid",
                                err.to_string(),
                            )]);
                            err_details.set_resource_info(resource_type, resource_name, "", "");
                            Status::with_error_details(
                                Code::InvalidArgument,
                                format!(
                                    "Insert API invalid id: \"{}\" or vector: {:?} was given",
                                    vec.id, vec.vector
                                ),
                                err_details,
                            )
                        }
                        _ => {
                            let mut err_details = ErrorDetails::new();
                            err_details.set_error_info(err.to_string(), domain, metadata);
                            err_details.set_request_info(
                                vec.id,
                                String::from_utf8(req.encode_to_vec())
                                    .unwrap_or_else(|_| "<invalid UTF-8>".to_string()),
                            );
                            err_details.set_resource_info(resource_type, resource_name, "", "");
                            Status::with_error_details(
                                Code::Unknown,
                                "failed to parse Insert gRPC error response",
                                err_details,
                            )
                        }
                    };
                    Err(status)
                }
                Ok(()) => Ok(tonic::Response::new(object::Location {
                    name: self.name.clone(),
                    uuid: vec.id,
                    ips: vec![self.ip.clone()],
                })),
            }
        }
    }

    #[doc = " Server streaming response type for the StreamInsert method."]
    type StreamInsertStream = crate::stream_type!(object::StreamLocation);

    #[doc = " A method to add new multiple vectors by bidirectional streaming.\n"]
    async fn stream_insert(
        &self,
        _request: tonic::Request<tonic::Streaming<insert::Request>>,
    ) -> std::result::Result<tonic::Response<Self::StreamInsertStream>, tonic::Status> {
        todo!()
    }

    #[doc = " A method to add new multiple vectors in a single request.\n"]
    async fn multi_insert(
        &self,
        _request: tonic::Request<insert::MultiRequest>,
    ) -> std::result::Result<tonic::Response<object::Locations>, tonic::Status> {
        todo!()
    }
}
