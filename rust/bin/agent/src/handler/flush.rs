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
use log::{debug, error, info};
use prost::Message;
use proto::{payload::v1::info, vald::v1::flush_server};
use tonic::{Code, Status};
use tonic_types::StatusExt;

use crate::handler::common::build_error_details;

#[tonic::async_trait]
impl<S: algorithm::ANN + 'static> flush_server::Flush for super::Agent<S> {
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
                    let resource_type = self.resource_type.clone() + "/qbg.Flush";
                    let resource_name = format!("{}: {}({})", self.api_name, self.name, self.ip);
                    let err_details = build_error_details(
                        err.to_string(),
                        domain,
                        "",
                        request.get_ref().encode_to_vec(),
                        &resource_type,
                        &resource_name,
                        None,
                    );
                    let status = match err {
                        Error::FlushingIsInProgress {} => {
                            let status = Status::with_error_details(Code::Aborted, "Flush API aborted due to flushing indices is in progress", err_details);
                            debug!("{:?}", status);
                            status
                        }
                        Error::WriteOperationToReadReplica {} => {
                            let status = Status::with_error_details(
                                Code::Aborted,
                                "Flush API aborted due to agent is read only",
                                err_details,
                            );
                            debug!("{:?}", status);
                            status
                        }
                        _ => {
                            let status = Status::with_error_details(
                                Code::Internal,
                                "Flush API is failed",
                                err_details,
                            );
                            error!("{:?}", status);
                            status
                        }
                    };
                    Err(status)
                }
                Ok(()) => {
                    let res = info::index::Count {
                        stored: 0,
                        uncommitted: 0,
                        indexing: false,
                        saving: false,
                    };
                    Ok(tonic::Response::new(res))
                }
            }
        }
    }
}
