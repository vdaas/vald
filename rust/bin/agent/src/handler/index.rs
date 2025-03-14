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
use log::{error, info};
use proto::{
    core::v1::agent_server,
    payload::v1::{control, info, Empty},
    vald::v1::index_server,
};
use std::collections::HashMap;
use tonic::{Code, Status};
use tonic_types::{ErrorDetails, PreconditionViolation, StatusExt};

#[tonic::async_trait]
impl agent_server::Agent for super::Agent {
    async fn create_index(
        &self,
        request: tonic::Request<control::CreateIndexRequest>,
    ) -> std::result::Result<tonic::Response<Empty>, tonic::Status> {
        info!("Recieved a request from {:?}", request.remote_addr());
        let req = request.get_ref();
        let pool_size = req.pool_size;
        let hostname = cargo::util::hostname()?;
        let domain = hostname.to_str().unwrap();
        let res = Empty {};
        {
            let mut s = self.s.write().await;
            let result = s.create_index();
            match result {
                Err(err) => {
                    let metadata = HashMap::new();
                    let resource_type = self.resource_type.clone() + "/qbg.CreateIndex";
                    let resource_name = format!("{}: {}({})", self.api_name, self.name, self.ip);
                    let status = match err {
                        Error::UncommittedIndexNotFound {} => {
                            let mut err_details = ErrorDetails::new();
                            err_details.set_error_info(err.to_string(), domain, metadata);
                            err_details.set_precondition_failure(vec![PreconditionViolation::new(
                                "uncommitted index is empty",
                                "failed to CreateIndex operation caused by empty uncommitted indices",
                                err.to_string(),
                            )]);
                            err_details.set_resource_info(resource_type, resource_name, "", "");
                            Status::with_error_details(
                                Code::FailedPrecondition,
                                format!("CreateIndex API failed to create indexes pool_size = {} due to the precondition failure, error: {}", pool_size, err.to_string()),
                                err_details,
                            )
                        }
                        Error::FlushingIsInProgress {} => {
                            let mut err_details = ErrorDetails::new();
                            err_details.set_error_info(err.to_string(), domain, metadata);
                            err_details.set_resource_info(resource_type, resource_name, "", "");
                            Status::with_error_details(
                                Code::Aborted,
                                "CreateIndex API aborted to process create indexes request due to flushing indices is in progress",
                                err_details,
                            )
                        }
                        _ => {
                            let mut err_details = ErrorDetails::new();
                            err_details.set_error_info(err.to_string(), domain, metadata);
                            err_details.set_resource_info(resource_type, resource_name, "", "");
                            let status = Status::with_error_details(
                                Code::Internal,
                                format!("CreateIndex API failed to create indexes pool_size = {}, error: {}", pool_size, err.to_string()),
                                err_details,
                            );
                            error!("{:?}", status);
                            status
                        }
                    };
                    Err(status)
                }
                Ok(()) => Ok(tonic::Response::new(res)),
            }
        }
    }

    async fn save_index(
        &self,
        request: tonic::Request<Empty>,
    ) -> std::result::Result<tonic::Response<Empty>, tonic::Status> {
        info!("Recieved a request from {:?}", request.remote_addr());
        let hostname = cargo::util::hostname()?;
        let domain = hostname.to_str().unwrap();
        let res = Empty {};
        {
            let mut s = self.s.write().await;
            let result = s.save_index();
            match result {
                Err(err) => {
                    error!("{:?}", err);
                    let metadata = HashMap::new();
                    let resource_type = self.resource_type.clone() + "/qbg.SaveIndex";
                    let resource_name = format!("{}: {}({})", self.api_name, self.name, self.ip);
                    let mut err_details = ErrorDetails::new();
                    err_details.set_error_info(err.to_string(), domain, metadata);
                    err_details.set_resource_info(resource_type, resource_name, "", "");
                    let status = Status::with_error_details(
                        Code::Internal,
                        "SaveIndex API failed to save indices",
                        err_details,
                    );
                    error!("{:?}", status);
                    Err(status)
                }
                Ok(()) => Ok(tonic::Response::new(res)),
            }
        }
    }

    #[doc = " Represent the creating and saving index RPC.\n"]
    async fn create_and_save_index(
        &self,
        _request: tonic::Request<control::CreateIndexRequest>,
    ) -> std::result::Result<tonic::Response<Empty>, tonic::Status> {
        todo!()
    }
}

#[tonic::async_trait]
impl index_server::Index for super::Agent {
    #[doc = " Represent the RPC to get the agent index information.\n"]
    async fn index_info(
        &self,
        request: tonic::Request<Empty>,
    ) -> std::result::Result<tonic::Response<info::index::Count>, tonic::Status> {
        info!("Recieved a request from {:?}", request.remote_addr());
        {
            let s = self.s.read().await;
            Ok(tonic::Response::new(info::index::Count {
                stored: s.len(),
                uncommitted: s.insert_vqueue_buffer_len() + s.delete_vqueue_buffer_len(),
                indexing: s.is_indexing(),
                saving: s.is_saving(),
            }))
        }
    }

    #[doc = " Represent the RPC to get the agent index detailed information.\n"]
    async fn index_detail(
        &self,
        _request: tonic::Request<Empty>,
    ) -> std::result::Result<tonic::Response<info::index::Detail>, tonic::Status> {
        todo!()
    }
    #[doc = " Represent the RPC to get the agent index statistics.\n"]
    async fn index_statistics(
        &self,
        _request: tonic::Request<Empty>,
    ) -> std::result::Result<tonic::Response<info::index::Statistics>, tonic::Status> {
        todo!()
    }

    #[doc = " Represent the RPC to get the agent index detailed statistics.\n"]
    async fn index_statistics_detail(
        &self,
        _request: tonic::Request<Empty>,
    ) -> std::result::Result<tonic::Response<info::index::StatisticsDetail>, tonic::Status> {
        todo!()
    }

    #[doc = " Represent the RPC to get the index property.\n"]
    async fn index_property(
        &self,
        _request: tonic::Request<Empty>,
    ) -> std::result::Result<tonic::Response<info::index::PropertyDetail>, tonic::Status> {
        todo!()
    }
}
