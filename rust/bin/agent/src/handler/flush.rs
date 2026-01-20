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
use log::{error, info};
use proto::{payload::v1::Empty, vald::v1::flush_server};
use std::collections::HashMap;
use tonic::{Code, Status};
use tonic_types::StatusExt;

#[tonic::async_trait]
impl flush_server::Flush for super::Agent {
    async fn flush(
        &self,
        request: tonic::Request<proto::payload::v1::flush::Request>,
    ) -> std::result::Result<tonic::Response<proto::payload::v1::info::index::Count>, Status> {
        info!("Recieved a request from {:?}", request.remote_addr());
        let hostname = cargo::util::hostname()?;
        let domain = hostname.to_str().unwrap();
        {
            let mut s = self.s.write().await;
            let result = s.regenerate_indexes().await;
            match result {
                Err(err) => {
                    error!("{:?}", err);
                    let metadata = HashMap::new();
                    let resource_type = self.resource_type.clone() + "/qbg.Flush";
                    let resource_name = format!("{}: {}({})", self.api_name, self.name, self.ip);
                    let mut err_details = tonic_types::ErrorDetails::new();
                    err_details.set_error_info(err.to_string(), domain, metadata);
                    err_details.set_resource_info(resource_type, resource_name, "", "");
                    let status = match err {
                        Error::FlushInprocess {} => {
                            let err_details = build_error_details(
                                err,
                                domain,
                                &vec.id,
                                request_bytes,
                                &resource_type,
                                &resource_name,
                                None,
                            );
                            let status = Status::with_error_details(Code::Aborted, "Flush API aborted due to flushing indices is in progress", err_details);
                            warn!("{:?}", status);
                            status
                        }
                        _ => {
                            let err_details = build_error_details(
                                err,
                                domain,
                                &vec.id,
                                request_bytes,
                                &resource_type,
                                &resource_name,
                                None,
                            );
                            Status::with_error_details(
                                Code::Unknown,
                                "failed to parse Insert gRPC error response",
                                err_details,
                            )
                        }
                    };
                    Err(status)
                }
                Ok(()) => {
                    counts = info::index::Count {
                        
                    };
                    Ok(tonic::Response::new(res))
                }
            }
        }
    }
}
