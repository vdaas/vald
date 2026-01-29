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
use proto::{
    core::v1::agent_server,
    payload::v1::{Empty, control, info},
    vald::v1::index_server,
};
use std::collections::HashMap;
use tonic::{Code, Status};
use tonic_types::{PreconditionViolation, StatusExt};

use crate::handler::common::build_error_details;

#[tonic::async_trait]
impl<S: algorithm::ANN + 'static> agent_server::Agent for super::Agent<S> {
    async fn create_index(
        &self,
        request: tonic::Request<control::CreateIndexRequest>,
    ) -> std::result::Result<tonic::Response<Empty>, tonic::Status> {
        info!("Received a request from {:?}", request.remote_addr());
        let req = request.get_ref();
        let pool_size = req.pool_size;
        let res = Empty {};
        let mut s = self.s.write().await;
        let result = s.create_index().await;
        match result {
            Err(err) => {
                let resource_type = format!("{}/qbg.CreateIndex", self.resource_type);
                let resource_name = format!("{}: {}({})", self.api_name, self.name, self.ip);
                let status = match err {
                    Error::UncommittedIndexNotFound {} => {
                        let mut err_details = build_error_details(&err, "", vec![], &resource_type, &resource_name, None);
                        err_details.set_precondition_failure(vec![PreconditionViolation::new(
                            "uncommitted index is empty",
                            "failed to CreateIndex operation caused by empty uncommitted indices",
                            err.to_string(),
                        )]);
                        Status::with_error_details(
                            Code::FailedPrecondition,
                            format!("CreateIndex API failed to create indexes pool_size = {} due to the precondition failure, error: {}", pool_size, err),
                            err_details,
                        )
                    }
                    Error::FlushingIsInProgress {} => {
                        let err_details = build_error_details(&err, "", vec![], &resource_type, &resource_name, None);
                        Status::with_error_details(
                            Code::Aborted,
                            "CreateIndex API aborted to process create indexes request due to flushing indices is in progress",
                            err_details,
                        )
                    }
                    _ => {
                        let err_details = build_error_details(&err, "", vec![], &resource_type, &resource_name, None);
                        let status = Status::with_error_details(
                            Code::Internal,
                            format!("CreateIndex API failed to create indexes pool_size = {}, error: {}", pool_size, err),
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

    async fn save_index(
        &self,
        request: tonic::Request<Empty>,
    ) -> std::result::Result<tonic::Response<Empty>, tonic::Status> {
        info!("Received a request from {:?}", request.remote_addr());
        let res = Empty {};
        {
            let mut s = self.s.write().await;
            let result = s.save_index().await;
            match result {
                Err(err) => {
                    error!("{:?}", err);
                    let resource_type = format!("{}/qbg.SaveIndex", self.resource_type);
                    let resource_name = format!("{}: {}({})", self.api_name, self.name, self.ip);
                    let err_details = build_error_details(&err, "", vec![], &resource_type, &resource_name, None);
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
        request: tonic::Request<control::CreateIndexRequest>,
    ) -> std::result::Result<tonic::Response<Empty>, tonic::Status> {
        info!("Received a request from {:?}", request.remote_addr());
        let req = request.get_ref();
        let pool_size = req.pool_size;
        let res = Empty {};
        let mut s = self.s.write().await;
        let result = s.create_and_save_index().await;
        match result {
            Err(err) => {
                let resource_type = format!("{}/qbg.CreateAndSaveIndex", self.resource_type);
                let resource_name = format!("{}: {}({})", self.api_name, self.name, self.ip);
                let status = match err {
                    Error::UncommittedIndexNotFound {} => {
                        let mut err_details = build_error_details(&err, "", vec![], &resource_type, &resource_name, None);
                        err_details.set_precondition_failure(vec![PreconditionViolation::new(
                            "uncommitted index is empty",
                            "failed to CreateAndSaveIndex operation caused by empty uncommitted indices",
                            err.to_string(),
                        )]);
                        Status::with_error_details(
                            Code::FailedPrecondition,
                            format!("CreateAndSaveIndex API failed to create indexes pool_size = {} due to the precondition failure, error: {}", pool_size, err),
                            err_details,
                        )
                    }
                    Error::FlushingIsInProgress {} => {
                        let err_details = build_error_details(&err, "", vec![], &resource_type, &resource_name, None);
                        Status::with_error_details(
                            Code::Aborted,
                            "CreateAndSaveIndex API aborted to process create indexes request due to flushing indices is in progress",
                            err_details,
                        )
                    }
                    _ => {
                        let err_details = build_error_details(&err, "", vec![], &resource_type, &resource_name, None);
                        let status = Status::with_error_details(
                            Code::Internal,
                            format!("CreateAndSaveIndex API failed to create indexes pool_size = {}, error: {}", pool_size, err),
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

#[tonic::async_trait]
impl<S: algorithm::ANN + 'static> index_server::Index for super::Agent<S> {
    #[doc = " Represent the RPC to get the agent index information.\n"]
    async fn index_info(
        &self,
        request: tonic::Request<Empty>,
    ) -> std::result::Result<tonic::Response<info::index::Count>, tonic::Status> {
        info!("Received a request from {:?}", request.remote_addr());
        let s = self.s.read().await;
        Ok(tonic::Response::new(info::index::Count {
            stored: s.len(),
            uncommitted: s.insert_vqueue_buffer_len() + s.delete_vqueue_buffer_len(),
            indexing: s.is_indexing(),
            saving: s.is_saving(),
        }))
    }

    #[doc = " Represent the RPC to get the agent index detailed information.\n"]
    async fn index_detail(
        &self,
        request: tonic::Request<Empty>,
    ) -> std::result::Result<tonic::Response<info::index::Detail>, tonic::Status> {
        info!("Received a request from {:?}", request.remote_addr());
        let s = self.s.read().await;
        let mut counts = HashMap::new();
        counts.insert(
            self.name.clone(),
            info::index::Count {
                stored: s.len(),
                uncommitted: s.insert_vqueue_buffer_len() + s.delete_vqueue_buffer_len(),
                indexing: s.is_indexing(),
                saving: s.is_saving(),
            },
        );
        Ok(tonic::Response::new(info::index::Detail {
            counts,
            replica: 1,
            live_agents: 1,
        }))
    }

    #[doc = " Represent the RPC to get the agent index statistics.\n"]
    async fn index_statistics(
        &self,
        request: tonic::Request<Empty>,
    ) -> std::result::Result<tonic::Response<info::index::Statistics>, tonic::Status> {
        info!("Received a request from {:?}", request.remote_addr());
        let s = self.s.read().await;
        match s.index_statistics() {
            Ok(stats) => Ok(tonic::Response::new(stats)),
            Err(err) => {
                error!("IndexStatistics API failed: {:?}", err);
                let resource_type = format!("{}/qbg.IndexStatistics", self.resource_type);
                let resource_name = format!("{}: {}({})", self.api_name, self.name, self.ip);
                let err_details = build_error_details(&err, "", vec![], &resource_type, &resource_name, None);
                Err(Status::with_error_details(
                    Code::Internal,
                    format!("IndexStatistics API failed: {}", err),
                    err_details,
                ))
            }
        }
    }

    #[doc = " Represent the RPC to get the agent index detailed statistics.\n"]
    async fn index_statistics_detail(
        &self,
        request: tonic::Request<Empty>,
    ) -> std::result::Result<tonic::Response<info::index::StatisticsDetail>, tonic::Status> {
        info!("Received a request from {:?}", request.remote_addr());
        let s = self.s.read().await;
        match s.index_statistics() {
            Ok(stats) => {
                let mut details = HashMap::new();
                details.insert(self.name.clone(), stats);
                Ok(tonic::Response::new(info::index::StatisticsDetail { details }))
            }
            Err(err) => {
                error!("IndexStatisticsDetail API failed: {:?}", err);
                let resource_type = format!("{}/qbg.IndexStatisticsDetail", self.resource_type);
                let resource_name = format!("{}: {}({})", self.api_name, self.name, self.ip);
                let err_details = build_error_details(&err, "", vec![], &resource_type, &resource_name, None);
                Err(Status::with_error_details(
                    Code::Internal,
                    format!("IndexStatisticsDetail API failed: {}", err),
                    err_details,
                ))
            }
        }
    }

    #[doc = " Represent the RPC to get the index property.\n"]
    async fn index_property(
        &self,
        request: tonic::Request<Empty>,
    ) -> std::result::Result<tonic::Response<info::index::PropertyDetail>, tonic::Status> {
        info!("Received a request from {:?}", request.remote_addr());
        let s = self.s.read().await;
        match s.index_property() {
            Ok(prop) => {
                let mut details = HashMap::new();
                details.insert(self.name.clone(), prop);
                Ok(tonic::Response::new(info::index::PropertyDetail { details }))
            }
            Err(err) => {
                error!("IndexProperty API failed: {:?}", err);
                let resource_type = format!("{}/qbg.IndexProperty", self.resource_type);
                let resource_name = format!("{}: {}({})", self.api_name, self.name, self.ip);
                let err_details = build_error_details(&err, "", vec![], &resource_type, &resource_name, None);
                Err(Status::with_error_details(
                    Code::Internal,
                    format!("IndexProperty API failed: {}", err),
                    err_details,
                ))
            }
        }
    }
}
