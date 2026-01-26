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

pub trait MultiError {
    fn new_uuid_already_exists(uuids: Vec<String>) -> Error;
    fn new_object_id_not_found(uuids: Vec<String>) -> Error;
    fn new_invalid_dimension_size(
        current: Vec<String>,
        limit: Vec<String>,
    ) -> Error;
    fn new_uuid_not_found(uuids: Vec<String>) -> Error;
    fn split_uuids(uuids: String) -> Vec<String>;
}

#[derive(thiserror::Error, Debug)]
pub enum Error {
    #[error("create indexing is in progress")]
    CreateIndexingIsInProgress {},
    #[error("search result is empty")]
    EmptySearchResult {},
    #[error("flush is in progress")]
    FlushingIsInProgress {},
    #[error("incompatible dimension size detected\trequested: {got},\tconfigured: {want}")]
    IncompatibleDimensionSize {
        got: usize,
        want: usize,
    },
    #[error("uuid {uuid} index already exists")]
    UUIDAlreadyExists {
        uuid: String,
    },
    #[error("object uuid{} not found", if uuid == "0" { "" } else { " {uuid}'s metadata" })]
    UUIDNotFound {
        uuid: String,
    },
    #[error("uncommitted indexes are not found")]
    UncommittedIndexNotFound {},
    #[error("uuid \"{uuid}\" is invalid")]
    InvalidUUID {
        uuid: String,
    },
    #[error("dimension size {} is invalid, the supporting dimension size must be {}", current, if limit == "0" { "bigger than 2" } else { "between 2 ~ {limit}" })]
    InvalidDimensionSize{
        current: String,
        limit: String,
    },
    #[error("uuid {uuid}'s object id not found")]
    ObjectIDNotFound {
        uuid: String,
    },
    #[error("write operation to read replica is not possible")]
    WriteOperationToReadReplica {},
    #[error("{method} is not supported for {algorithm}")]
    Unsupported {
        method: String,
        algorithm: String,
    },
    #[error("{0}")]
    Internal(#[from] Box<dyn std::error::Error + Send + Sync>),
    #[error("unknown error")]
    Unknown {},
}

impl MultiError for Error {
    fn new_uuid_already_exists(uuids: Vec<String>) -> Error {
        Error::UUIDAlreadyExists {
            uuid: uuids.join(","),
        }
    }

    fn new_object_id_not_found(uuids: Vec<String>) -> Error {
        Error::ObjectIDNotFound {
            uuid: uuids.join(","),
        }
    }

    fn new_invalid_dimension_size(
        current: Vec<String>,
        limit: Vec<String>,
    ) -> Error {
        Error::InvalidDimensionSize {
            current: current.join(","),
            limit: limit.join(","),
        }
    }

    fn new_uuid_not_found(uuids: Vec<String>) -> Error {
        Error::UUIDNotFound {
            uuid: uuids.join(","),
        }
    }

    fn split_uuids(uuids: String) -> Vec<String> {
        uuids.split(",").map(|x| x.to_string()).collect()
    }
}
