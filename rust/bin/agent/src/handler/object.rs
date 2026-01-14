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
use log::{info, warn};
use prost::Message;
use proto::{payload::v1::object, vald::v1::object_server};
use std::sync::Arc;
use tokio::sync::RwLock;
use tonic::{Code, Status};
use tonic_types::StatusExt;

use super::common::{bidirectional_stream, build_error_details};

async fn get_object(
    s: Arc<RwLock<dyn algorithm::ANN>>,
    resource_type: &str,
    api_name: &str,
    name: &str,
    ip: &str,
    request: &object::VectorRequest,
) -> Result<object::Vector, Status> {
    let id = match request.id.clone() {
        Some(id) => id,
        None => return Err(Status::invalid_argument("Missing ID in request")),
    };
    let uuid = id.id;
    let hostname = cargo::util::hostname()?;
    let domain = hostname.to_str().unwrap();
    {
        let s = s.read().await;
        if uuid.len() == 0 {
            let err = Error::InvalidUUID { uuid: uuid.clone() };
            let resource_type = format!("{}/qbg.GetObject", resource_type);
            let resource_name = format!("{}: {}({})", api_name, name, ip);
            let err_details = build_error_details(
                err,
                domain,
                &uuid,
                request.encode_to_vec(),
                &resource_type,
                &resource_name,
                Some("uuid"),
            );
            let status = Status::with_error_details(
                Code::InvalidArgument,
                format!(
                    "GetObject API invalid argument for uuid \"{}\" detected",
                    uuid
                ),
                err_details,
            );
            warn!("{:?}", status);
            return Err(status);
        }
        let result = s.get_object(uuid.clone());
        match result {
            Err(_err) => {
                let status =
                    Status::new(Code::NotFound, format!("uuid {}'s object not found", uuid));
                Err(status)
            }
            Ok((vec, ts)) => Ok(object::Vector {
                id: uuid,
                vector: vec,
                timestamp: ts,
            }),
        }
    }
}

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
        info!("Recieved a request from {:?}", request.remote_addr());
        let request = request.get_ref();
        let s = self.s.clone();
        let resource_type = self.resource_type.clone();
        let name = self.name.clone();
        let ip = self.ip.clone();
        let api_name = self.api_name.clone();
        match get_object(s, &resource_type, &api_name, &name, &ip, request).await {
            Ok(vector) => Ok(tonic::Response::new(vector)),
            Err(e) => Err(e),
        }
    }

    type StreamGetObjectStream = crate::stream_type!(object::StreamVector);

    async fn stream_get_object(
        &self,
        request: tonic::Request<tonic::Streaming<object::VectorRequest>>,
    ) -> std::result::Result<tonic::Response<Self::StreamGetObjectStream>, tonic::Status> {
        info!(
            "Received stream get object request from {:?}",
            request.remote_addr()
        );

        let s = self.s.clone();
        let resource_type = self.resource_type.clone() + "/qbg.StreamGetObject";
        let name = self.name.clone();
        let ip = self.ip.clone();
        let api_name = self.api_name.clone();

        let process_fn = move |req: object::VectorRequest| {
            let s = s.clone();
            let resource_type = resource_type.clone();
            let name = name.clone();
            let ip = ip.clone();
            let api_name = api_name.clone();
            async move {
                match get_object(s, &resource_type, &api_name, &name, &ip, &req).await {
                    Ok(vector) => Ok(object::StreamVector {
                        payload: Some(object::stream_vector::Payload::Vector(vector)),
                    }),
                    Err(status) => Err(status),
                }
            }
        };

        bidirectional_stream(request, self.stream_concurrency, process_fn).await
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
