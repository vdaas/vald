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
    payload::v1::{object, remove},
    vald::v1::remove_server,
};
use std::collections::HashMap;
use tonic::{Code, Status};
use tonic_types::{ErrorDetails, FieldViolation, StatusExt};

#[tonic::async_trait]
impl remove_server::Remove for super::Agent {
    async fn remove(
        &self,
        request: tonic::Request<remove::Request>,
    ) -> std::result::Result<tonic::Response<object::Location>, tonic::Status> {
        println!("Recieved a request from {:?}", request.remote_addr());
        let req = request.get_ref();
        let config = req.config.clone().unwrap();
        let id = req.id.clone().unwrap();
        let uuid = id.id;
        let hostname = cargo::util::hostname()?;
        let domain = hostname.to_str().unwrap();
        {
            let mut s = self.s.write().await;
            if uuid.len() == 0 {
                let err = Error::InvalidUUID { uuid: uuid.clone() };
                let metadata = HashMap::new();
                let resource_type = self.resource_type.clone() + "/qbg.Remove";
                let resource_name = format!("{}: {}({})", self.api_name, self.name, self.ip);
                let mut err_details = ErrorDetails::new();
                err_details.set_error_info(err.to_string(), domain, metadata);
                err_details.set_request_info(
                    uuid.clone(),
                    String::from_utf8(req.encode_to_vec()).unwrap(),
                );
                err_details.set_bad_request(vec![FieldViolation::new("uuid", err.to_string())]);
                err_details.set_resource_info(resource_type, resource_name, "", "");
                let status = Status::with_error_details(
                    Code::InvalidArgument,
                    format!("Remove API invalid argument for uuid \"{}\" detected", uuid),
                    err_details,
                );
                return Err(status);
            }
            let result = s.remove(uuid.clone(), config.timestamp);
            match result {
                Err(err) => {
                    let metadata = HashMap::new();
                    let resource_type = self.resource_type.clone() + "/qbg.Remove";
                    let resource_name = format!("{}: {}({})", self.api_name, self.name, self.ip);
                    let status = match err {
                        Error::FlushingIsInProgress {} => {
                            let mut err_details = ErrorDetails::new();
                            err_details.set_error_info(err.to_string(), domain, metadata);
                            err_details.set_request_info(
                                uuid,
                                String::from_utf8(req.encode_to_vec()).unwrap(),
                            );
                            err_details.set_resource_info(resource_type, resource_name, "", "");
                            Status::with_error_details(Code::Aborted, "Remove API aborted to process remove request due to flushing indices is in progress", err_details)
                        }
                        Error::ObjectIDNotFound { uuid: _ } => {
                            let mut err_details = ErrorDetails::new();
                            err_details.set_error_info(err.to_string(), domain, metadata);
                            err_details.set_request_info(
                                uuid.clone(),
                                String::from_utf8(req.encode_to_vec()).unwrap(),
                            );
                            err_details.set_resource_info(resource_type, resource_name, "", "");
                            Status::with_error_details(
                                Code::NotFound,
                                format!("Remove API uuid {} not found", uuid),
                                err_details,
                            )
                        }
                        Error::UUIDNotFound { id: _ } => {
                            let mut err_details = ErrorDetails::new();
                            err_details.set_error_info(err.to_string(), domain, metadata);
                            err_details.set_request_info(
                                uuid.clone(),
                                String::from_utf8(req.encode_to_vec()).unwrap(),
                            );
                            err_details.set_bad_request(vec![FieldViolation::new(
                                "uuid",
                                err.to_string(),
                            )]);
                            err_details.set_resource_info(resource_type, resource_name, "", "");
                            Status::with_error_details(
                                Code::InvalidArgument,
                                format!(
                                    "Remove API invalid argument for uuid \"{}\" detected",
                                    uuid
                                ),
                                err_details,
                            )
                        }
                        _ => {
                            let mut err_details = ErrorDetails::new();
                            err_details.set_error_info(err.to_string(), domain, metadata);
                            err_details.set_request_info(
                                uuid,
                                String::from_utf8(req.encode_to_vec()).unwrap(),
                            );
                            err_details.set_resource_info(resource_type, resource_name, "", "");
                            Status::with_error_details(
                                Code::Internal,
                                "Remove API failed",
                                err_details,
                            )
                        }
                    };
                    Err(status)
                }
                Ok(()) => Ok(tonic::Response::new(object::Location {
                    name: self.name.clone(),
                    uuid: uuid,
                    ips: vec![self.ip.clone()],
                })),
            }
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
        _request: tonic::Request<tonic::Streaming<remove::Request>>,
    ) -> std::result::Result<tonic::Response<Self::StreamRemoveStream>, tonic::Status> {
        todo!()
    }

    #[doc = " A method to remove multiple indexed vectors in a single request.\n"]
    async fn multi_remove(
        &self,
        _request: tonic::Request<remove::MultiRequest>,
    ) -> std::result::Result<tonic::Response<object::Locations>, tonic::Status> {
        todo!()
    }
}
