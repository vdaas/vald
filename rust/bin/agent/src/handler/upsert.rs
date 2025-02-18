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
    payload::v1::{insert, object, update, upsert},
    vald::v1::{insert_server::Insert, update_server::Update, upsert_server},
};
use std::collections::HashMap;
use tonic::{Code, Status};
use tonic_types::{ErrorDetails, FieldViolation, StatusExt};

#[tonic::async_trait]
impl upsert_server::Upsert for super::Agent {
    async fn upsert(
        &self,
        request: tonic::Request<upsert::Request>,
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
            let s = self.s.read().await;
            let vec = match req.vector.clone() {
                Some(v) => v,
                None => return Err(Status::invalid_argument("Missing vector in request")),
            };
            let uuid = vec.id.clone();
            if vec.vector.len() != s.get_dimension_size() {
                let err = Error::IncompatibleDimensionSize {
                    got: vec.vector.len(),
                    want: s.get_dimension_size(),
                };
                let mut err_details = ErrorDetails::new();
                let metadata = HashMap::new();
                let resource_type = self.resource_type.clone() + "/qbg.Upsert";
                let resource_name = format!("{}: {}({})", self.api_name, self.name, self.ip);
                err_details.set_error_info(err.to_string(), domain, metadata);
                err_details.set_request_info(
                    uuid,
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
                    "Upsert API Incompatible Dimension Size detected",
                    err_details,
                );
                return Err(status);
            }
            if uuid.len() == 0 {
                let err = Error::InvalidUUID { uuid: uuid.clone() };
                let mut err_details = ErrorDetails::new();
                let metadata = HashMap::new();
                let resource_type = self.resource_type.clone() + "/qbg.Upsert";
                let resource_name = format!("{}: {}({})", self.api_name, self.name, self.ip);
                err_details.set_error_info(err.to_string(), domain, metadata);
                err_details.set_request_info(
                    uuid.clone(),
                    String::from_utf8(req.encode_to_vec())
                        .unwrap_or_else(|_| "<invalid UTF-8>".to_string()),
                );
                err_details.set_bad_request(vec![FieldViolation::new("uuid", err.to_string())]);
                err_details.set_resource_info(resource_type, resource_name, "", "");
                let status = Status::with_error_details(
                    Code::InvalidArgument,
                    format!("Upsert API invalid argument for uuid \"{}\" detected", uuid),
                    err_details,
                );
                return Err(status);
            }
            let rt_name;
            let result;
            let exists = s.exists(uuid.clone());
            if exists {
                result = self
                    .update(tonic::Request::new(update::Request {
                        vector: req.vector.clone(),
                        config: Some(update::Config {
                            skip_strict_exist_check: true,
                            filters: config.filters,
                            timestamp: config.timestamp,
                            disable_balanced_update: config.disable_balanced_update,
                        }),
                    }))
                    .await;
                rt_name = format!("{}{}", "/qbg.Upsert", "/qbg.Update");
            } else {
                result = self
                    .insert(tonic::Request::new(insert::Request {
                        vector: req.vector.clone(),
                        config: Some(insert::Config {
                            skip_strict_exist_check: true,
                            filters: config.filters,
                            timestamp: config.timestamp,
                        }),
                    }))
                    .await;
                rt_name = format!("{}{}", "/qbg.Upsert", "/qbg.Insert");
            }
            match result {
                Err(st) => {
                    let status = match st.code() {
                        Code::Aborted
                        | Code::Cancelled
                        | Code::DeadlineExceeded
                        | Code::AlreadyExists
                        | Code::NotFound
                        | Code::Ok
                        | Code::Unimplemented => return Err(st),
                        _ => {
                            let metadata = HashMap::new();
                            let resource_type =
                                format!("{}{}", self.resource_type.clone(), rt_name);
                            let resource_name =
                                format!("{}: {}({})", self.api_name, self.name, self.ip);
                            let mut err_details = ErrorDetails::new();
                            err_details.set_error_info(
                                st.get_details_error_info().unwrap().reason,
                                domain,
                                metadata,
                            );
                            err_details.set_request_info(
                                uuid,
                                String::from_utf8(req.encode_to_vec())
                                    .unwrap_or_else(|_| "<invalid UTF-8>".to_string()),
                            );
                            err_details.set_resource_info(resource_type, resource_name, "", "");
                            Status::with_error_details(st.code(), st.message(), err_details)
                        }
                    };
                    Err(status)
                }
                Ok(res) => Ok(res),
            }
        }
    }

    #[doc = " Server streaming response type for the StreamUpsert method."]
    type StreamUpsertStream = crate::stream_type!(object::StreamLocation);

    #[doc = " A method to insert/update multiple vectors by bidirectional streaming.\n"]
    async fn stream_upsert(
        &self,
        _request: tonic::Request<tonic::Streaming<upsert::Request>>,
    ) -> std::result::Result<tonic::Response<Self::StreamUpsertStream>, tonic::Status> {
        todo!()
    }

    #[doc = " A method to insert/update multiple vectors in a single request.\n"]
    async fn multi_upsert(
        &self,
        _request: tonic::Request<upsert::MultiRequest>,
    ) -> std::result::Result<tonic::Response<object::Locations>, tonic::Status> {
        todo!()
    }
}
