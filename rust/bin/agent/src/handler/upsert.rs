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
use proto::{
    payload::v1::{insert, object, update, upsert},
    vald::v1::{insert_server::Insert, update_server::Update, upsert_server},
};
use std::sync::Arc;
use tokio::sync::RwLock;
use tonic::{Code, Status};
use tonic_types::StatusExt;

use super::common::{bidirectional_stream, build_error_details};
use super::insert::insert as insert_fn;
use super::update::update as update_fn;

async fn upsert<S: algorithm::ANN>(
    s: Arc<RwLock<S>>,
    resource_type: &str,
    api_name: &str,
    name: &str,
    ip: &str,
    request: &upsert::Request,
) -> Result<object::Location, Status> {
    let config = match request.config.clone() {
        Some(cfg) => cfg,
        None => return Err(Status::invalid_argument("Missing configuration in request")),
    };
    let vec = match request.vector.clone() {
        Some(v) => v,
        None => return Err(Status::invalid_argument("Missing vector in request")),
    };
    let uuid = vec.id.clone();

    // Check dimension size with a short-lived read lock
    {
        let s_inner = s.read().await;
        if vec.vector.len() != s_inner.get_dimension_size() {
            let err = Error::IncompatibleDimensionSize {
                got: vec.vector.len(),
                want: s_inner.get_dimension_size(),
            };
            let resource_type = format!("{}/qbg.Upsert", resource_type);
            let resource_name = format!("{}: {}({})", api_name, name, ip);
            let err_details = build_error_details(
                err,
                &vec.id,
                request.encode_to_vec(),
                &resource_type,
                &resource_name,
                Some("vector dimension size"),
            );
            let status = Status::with_error_details(
                Code::InvalidArgument,
                "Upsert API Incompatible Dimension Size detected",
                err_details,
            );
            warn!("{:?}", status);
            return Err(status);
        }
    }

    if uuid.is_empty() {
        let err = Error::InvalidUUID { uuid: uuid.clone() };
        let resource_type = format!("{}/qbg.Upsert", resource_type);
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
            format!("Upsert API invalid argument for uuid \"{}\" detected", uuid),
            err_details,
        );
        warn!("{:?}", status);
        return Err(status);
    }
    let rt_name;
    let result;
    let exists = {
        let s_inner = s.read().await;
        let (_, exists) = s_inner.exists(uuid.clone()).await;
        exists
    }; // s_inner dropped here to release read lock
    if exists {
        result = update_fn(
            s.clone(),
            resource_type,
            api_name,
            name,
            ip,
            &update::Request {
                vector: Some(vec),
                config: Some(update::Config {
                    skip_strict_exist_check: true,
                    filters: config.filters,
                    timestamp: config.timestamp,
                    disable_balanced_update: config.disable_balanced_update,
                }),
            },
        )
        .await;
        rt_name = format!("{}{}", "/qbg.Upsert", "/qbg.Update");
    } else {
        result = insert_fn(
            s.clone(),
            resource_type,
            api_name,
            name,
            ip,
            &insert::Request {
                vector: Some(vec),
                config: Some(insert::Config {
                    skip_strict_exist_check: true,
                    filters: config.filters,
                    timestamp: config.timestamp,
                }),
            },
        )
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
                    let resource_type = format!("{}{}", resource_type, rt_name);
                    let resource_name = format!("{}: {}({})", api_name, name, ip);
                    let err_details = build_error_details(
                        st.get_details_error_info().unwrap().reason,
                        &uuid,
                        request.encode_to_vec(),
                        &resource_type,
                        &resource_name,
                        None,
                    );
                    Status::with_error_details(st.code(), st.message(), err_details)
                }
            };
            Err(status)
        }
        Ok(res) => Ok(res),
    }
}

#[tonic::async_trait]
impl<S: algorithm::ANN + 'static> upsert_server::Upsert for super::Agent<S> {
    async fn upsert(
        &self,
        request: tonic::Request<upsert::Request>,
    ) -> std::result::Result<tonic::Response<object::Location>, tonic::Status> {
        info!("Recieved a request from {:?}", request.remote_addr());
        let request = request.get_ref();
        let s = self.s.clone();
        let resource_type = self.resource_type.clone();
        let name = self.name.clone();
        let ip = self.ip.clone();
        let api_name = self.api_name.clone();
        match upsert(s, &resource_type, &api_name, &name, &ip, request).await {
            Ok(location) => Ok(tonic::Response::new(location)),
            Err(e) => Err(e),
        }
    }

    #[doc = " Server streaming response type for the StreamUpsert method."]
    type StreamUpsertStream = crate::stream_type!(object::StreamLocation);

    #[doc = " A method to insert/update multiple vectors by bidirectional streaming.\n"]
    async fn stream_upsert(
        &self,
        request: tonic::Request<tonic::Streaming<upsert::Request>>,
    ) -> std::result::Result<tonic::Response<Self::StreamUpsertStream>, tonic::Status> {
        info!(
            "Received stream upsert request from {:?}",
            request.remote_addr()
        );
        let s = self.s.clone();
        let resource_type = self.resource_type.clone() + "/qbg.StreamUpsert";
        let name = self.name.clone();
        let ip = self.ip.clone();
        let api_name = self.api_name.clone();

        let process_fn = move |req: upsert::Request| {
            let s = s.clone();
            let resource_type = resource_type.clone();
            let name = name.clone();
            let ip = ip.clone();
            let api_name = api_name.clone();
            async move {
                match upsert(s, &resource_type, &api_name, &name, &ip, &req).await {
                    Ok(location) => Ok(object::StreamLocation {
                        payload: Some(object::stream_location::Payload::Location(location)),
                    }),
                    Err(status) => Err(status),
                }
            }
        };

        bidirectional_stream(request, self.stream_concurrency, process_fn).await
    }

    #[doc = " A method to insert/update multiple vectors in a single request.\n"]
    async fn multi_upsert(
        &self,
        request: tonic::Request<upsert::MultiRequest>,
    ) -> std::result::Result<tonic::Response<object::Locations>, tonic::Status> {
        info!("Received a request from {:?}", request.remote_addr());
        let mreq = request.get_ref();
        let mut ireqs = insert::MultiRequest { requests: vec![] };
        let mut ureqs = update::MultiRequest { requests: vec![] };
        let mut ids = vec![];

        // Use a block scope to release read lock before calling multi_insert/multi_update
        {
            let s = self.s.read().await;
            for req in mreq.requests.clone() {
                let vec = match req.vector.clone() {
                    Some(v) => v,
                    None => return Err(Status::invalid_argument("Missing vector in request")),
                };
                let config = match req.config.clone() {
                    Some(c) => c,
                    None => return Err(Status::invalid_argument("Missing config in request")),
                };
                if vec.vector.len() != s.get_dimension_size() {
                    let err = Error::IncompatibleDimensionSize {
                        got: vec.vector.len(),
                        want: s.get_dimension_size(),
                    };
                    let resource_type = self.resource_type.clone() + "/qbg.MultiUpsert";
                    let resource_name = format!("{}: {}({})", self.api_name, self.name, self.ip);
                    let err_details = build_error_details(
                        err,
                        &vec.id,
                        req.encode_to_vec(),
                        &resource_type,
                        &resource_name,
                        Some("vector dimension size"),
                    );
                    let status = Status::with_error_details(
                        Code::InvalidArgument,
                        "Upsert API Incompatible Dimension Size detected",
                        err_details,
                    );
                    warn!("{:?}", status);
                    return Err(status);
                }
                ids.push(vec.id.clone());
                let (_, exists) = s.exists(vec.id.clone()).await;
                if exists {
                    ureqs.requests.push(update::Request {
                        vector: Some(vec),
                        config: Some(update::Config {
                            skip_strict_exist_check: true,
                            filters: config.filters,
                            timestamp: config.timestamp,
                            disable_balanced_update: config.disable_balanced_update,
                        }),
                    });
                } else {
                    ireqs.requests.push(insert::Request {
                        vector: Some(vec),
                        config: Some(insert::Config {
                            skip_strict_exist_check: true,
                            filters: config.filters,
                            timestamp: config.timestamp,
                        }),
                    });
                }
            }
        } // read lock released here

        if ireqs.requests.is_empty() {
            let res = self.multi_update(tonic::Request::new(ureqs)).await?;
            return Ok(res);
        } else if ureqs.requests.is_empty() {
            let res = self.multi_insert(tonic::Request::new(ireqs)).await?;
            return Ok(res);
        } else {
            let ures = self.multi_update(tonic::Request::new(ureqs)).await?;
            let ires = self.multi_insert(tonic::Request::new(ireqs)).await?;

            let mut locs = object::Locations { locations: vec![] };
            let ilocs = ires.into_inner().locations;
            let ulocs = ures.into_inner().locations;
            if ulocs.is_empty() {
                locs.locations = ilocs;
            } else if ilocs.is_empty() {
                locs.locations = ulocs;
            } else {
                locs.locations = [ilocs, ulocs].concat();
            }
            return Ok(tonic::Response::new(locs));
        }
    }
}
