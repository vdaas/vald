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
use proto::{payload::v1::object, vald::v1::object_server};
use std::collections::HashMap;
use tonic::{Code, Status};
use tonic_types::{ErrorDetails, FieldViolation, StatusExt};

#[tonic::async_trait]
impl object_server::Object for super::Agent {
    async fn exists(
        &self,
        _request: tonic::Request<object::Id>,
    ) -> std::result::Result<tonic::Response<object::Id>, tonic::Status> {
        todo!()
    }
    async fn get_object(
        &self,
        request: tonic::Request<object::VectorRequest>,
    ) -> std::result::Result<tonic::Response<object::Vector>, tonic::Status> {
        println!("Recieved a request from {:?}", request.remote_addr());
        let req = request.get_ref();
        let id = match req.id.clone() {
            Some(id) => id,
            None => return Err(Status::invalid_argument("Missing ID in request")),
        };
        let uuid = id.id;
        let hostname = cargo::util::hostname()?;
        let domain = hostname.to_str().unwrap();
        {
            let s = self.s.read().await;
            if uuid.len() == 0 {
                let err = Error::InvalidUUID { uuid: uuid.clone() };
                let metadata = HashMap::new();
                let resource_type = self.resource_type.clone() + "/qbg.GetObject";
                let resource_name = format!("{}: {}({})", self.api_name, self.name, self.ip);
                let mut err_details = ErrorDetails::new();
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
                    format!(
                        "GetObject API invalid argument for uuid \"{}\" detected",
                        uuid
                    ),
                    err_details,
                );
                return Err(status);
            }
            let result = s.get_object(uuid.clone());
            match result {
                Err(_err) => {
                    let status =
                        Status::new(Code::NotFound, format!("uuid {}'s object not found", uuid));
                    Err(status)
                }
                Ok((vec, ts)) => Ok(tonic::Response::new(object::Vector {
                    id: uuid,
                    vector: vec,
                    timestamp: ts,
                })),
            }
        }
    }

    type StreamGetObjectStream = crate::stream_type!(object::StreamVector);

    async fn stream_get_object(
        &self,
        _request: tonic::Request<tonic::Streaming<object::VectorRequest>>,
    ) -> std::result::Result<tonic::Response<Self::StreamGetObjectStream>, tonic::Status> {
        todo!()
    }

    type StreamListObjectStream = crate::stream_type!(object::list::Response);

    async fn stream_list_object(
        &self,
        _request: tonic::Request<object::list::Request>,
    ) -> std::result::Result<tonic::Response<Self::StreamListObjectStream>, tonic::Status> {
        todo!()
    }

    async fn get_timestamp(
        &self,
        _request: tonic::Request<object::TimestampRequest>,
    ) -> std::result::Result<tonic::Response<object::Timestamp>, tonic::Status> {
        todo!()
    }
}
