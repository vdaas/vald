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

pub mod error;
pub use error::{Error, MultiError};

use anyhow::Result;
use proto::payload::v1::{info, search};
<<<<<<< HEAD
use std::{collections::HashMap, error, fmt, i64};

pub trait MultiError {
    fn new_uuid_already_exists(uuids: Vec<String>) -> Error;
    fn new_object_id_not_found(uuids: Vec<String>) -> Error;
    fn new_invalid_dimension_size(
        uuids: Vec<String>,
        current: Vec<String>,
        limit: Vec<String>,
    ) -> Error;
    fn new_uuid_not_found(uuids: Vec<String>) -> Error;
    fn split_uuids(uuids: String) -> Vec<String>;
}

#[derive(Debug)]
pub enum Error {
    CreateIndexingIsInProgress {},
    FlushingIsInProgress {},
    EmptySearchResult {},
    IncompatibleDimensionSize {
        got: usize,
        want: usize,
    },
    UUIDAlreadyExists {
        uuid: String,
    },
    UUIDNotFound {
        uuid: String,
    },
    UncommittedIndexNotFound {},
    InvalidUUID {
        uuid: String,
    },
    InvalidDimensionSize {
        uuid: String,
        current: String,
        limit: String,
    },
    ObjectIDNotFound {
        uuid: String,
    },
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
        uuids: Vec<String>,
        current: Vec<String>,
        limit: Vec<String>,
    ) -> Error {
        Error::InvalidDimensionSize {
            uuid: uuids.join(","),
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

impl error::Error for Error {}

impl fmt::Display for Error {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        match self {
            Error::CreateIndexingIsInProgress {} => write!(f, "create indexing is in progress"),
            Error::FlushingIsInProgress {} => write!(f, "flush is in progress"),
            Error::EmptySearchResult {} => write!(f, "search result is empty"),
            Error::IncompatibleDimensionSize { got, want } => write!(
                f,
                "incompatible dimension size detected\trequested: {},\tconfigured: {}",
                got, want
            ),
            Error::UUIDAlreadyExists { uuid } => write!(f, "uuid {} index already exists", uuid),
            Error::UUIDNotFound { uuid } => {
                if *uuid == "0" {
                    write!(f, "object uuid not found")
                } else {
                    write!(f, "object uuid {}'s metadata not found", uuid)
                }
            }
            Error::UncommittedIndexNotFound {} => write!(f, "uncommitted indexes are not found"),
            Error::InvalidUUID { uuid } => write!(f, "uuid \"{}\" is invalid", uuid),
            Error::InvalidDimensionSize {
                uuid: _,
                current,
                limit,
            } => {
                if *limit == "0" {
                    write!(
                        f,
                        "dimension size {} is invalid, the supporting dimension size must be bigger than 2",
                        current
                    )
                } else {
                    write!(
                        f,
                        "dimension size {} is invalid, the supporting dimension size must be between 2 ~ {}",
                        current, limit
                    )
                }
            }
            Error::ObjectIDNotFound { uuid } => write!(f, "uuid {}'s object id not found", uuid),
            Error::Unknown {} => write!(f, "unknown error"),
        }
    }
}
||||||| parent of 5831713ed (fix)
use std::{collections::HashMap, error, fmt, i64};

pub trait MultiError {
    fn new_uuid_already_exists(uuids: Vec<String>) -> Error;
    fn new_object_id_not_found(uuids: Vec<String>) -> Error;
    fn new_invalid_dimension_size(
        uuids: Vec<String>,
        current: Vec<String>,
        limit: Vec<String>,
    ) -> Error;
    fn new_uuid_not_found(uuids: Vec<String>) -> Error;
    fn split_uuids(uuids: String) -> Vec<String>;
}

#[derive(Debug)]
pub enum Error {
    CreateIndexingIsInProgress {},
    FlushingIsInProgress {},
    EmptySearchResult {},
    IncompatibleDimensionSize {
        got: usize,
        want: usize,
    },
    UUIDAlreadyExists {
        uuid: String,
    },
    UUIDNotFound {
        uuid: String,
    },
    UncommittedIndexNotFound {},
    InvalidUUID {
        uuid: String,
    },
    InvalidDimensionSize {
        uuid: String,
        current: String,
        limit: String,
    },
    ObjectIDNotFound {
        uuid: String,
    },
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
        uuids: Vec<String>,
        current: Vec<String>,
        limit: Vec<String>,
    ) -> Error {
        Error::InvalidDimensionSize {
            uuid: uuids.join(","),
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

impl error::Error for Error {}

impl fmt::Display for Error {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        match self {
            Error::CreateIndexingIsInProgress {} => write!(f, "create indexing is in progress"),
            Error::FlushingIsInProgress {} => write!(f, "flush is in progress"),
            Error::EmptySearchResult {} => write!(f, "search result is empty"),
            Error::IncompatibleDimensionSize { got, want } => write!(
                f,
                "incompatible dimension size detected\trequested: {},\tconfigured: {}",
                got, want
            ),
            Error::UUIDAlreadyExists { uuid } => write!(f, "uuid {} index already exists", uuid),
            Error::UUIDNotFound { uuid } => {
                if *uuid == "0" {
                    write!(f, "object uuid not found")
                } else {
                    write!(f, "object uuid {}'s metadata not found", uuid)
                }
            }
            Error::UncommittedIndexNotFound {} => write!(f, "uncommitted indexes are not found"),
            Error::InvalidUUID { uuid } => write!(f, "uuid \"{}\" is invalid", uuid),
            Error::InvalidDimensionSize {
                uuid: _,
                current,
                limit,
            } => {
                if *limit == "0" {
                    write!(f, "dimension size {} is invalid, the supporting dimension size must be bigger than 2", current)
                } else {
                    write!(f, "dimension size {} is invalid, the supporting dimension size must be between 2 ~ {}", current, limit)
                }
            }
            Error::ObjectIDNotFound { uuid } => write!(f, "uuid {}'s object id not found", uuid),
            Error::Unknown {} => write!(f, "unknown error"),
        }
    }
}
=======
use std::{collections::HashMap, future::Future, i64};
>>>>>>> 5831713ed (fix)

/// Trait for Approximate Nearest Neighbor (ANN) index implementations.
///
/// All methods that involve I/O or potentially blocking operations are async.
pub trait ANN: Send + Sync {
    // Search operations (async for potential I/O with vqueue/kvs)
    fn search(&self, vector: Vec<f32>, k: u32, epsilon: f32, radius: f32) -> impl Future<Output = Result<search::Response, Error>> + Send;
    fn search_by_id(&self, uuid: String, k: u32, epsilon: f32, radius: f32) -> impl Future<Output = Result<search::Response, Error>> + Send;
    fn linear_search(&self, vector: Vec<f32>, k: u32) -> impl Future<Output = Result<search::Response, Error>> + Send;
    fn linear_search_by_id(&self, uuid: String, k: u32) -> impl Future<Output = Result<search::Response, Error>> + Send;

    // Insert operations (async for vqueue push)
    fn insert(&mut self, uuid: String, vector: Vec<f32>) -> impl Future<Output = Result<(), Error>> + Send;
    fn insert_with_time(&mut self, uuid: String, vector: Vec<f32>, t: i64) -> impl Future<Output = Result<(), Error>> + Send;
    fn insert_multiple(&mut self, vectors: HashMap<String, Vec<f32>>) -> impl Future<Output = Result<(), Error>> + Send;
    fn insert_multiple_with_time(&mut self, vectors: HashMap<String, Vec<f32>>, t: i64) -> impl Future<Output = Result<(), Error>> + Send;

    // Update operations (async for vqueue/kvs)
    fn update(&mut self, uuid: String, vector: Vec<f32>) -> impl Future<Output = Result<(), Error>> + Send;
    fn update_with_time(&mut self, uuid: String, vector: Vec<f32>, t: i64) -> impl Future<Output = Result<(), Error>> + Send;
    fn update_multiple(&mut self, vectors: HashMap<String, Vec<f32>>) -> impl Future<Output = Result<(), Error>> + Send;
    fn update_multiple_with_time(&mut self, vectors: HashMap<String, Vec<f32>>, t: i64) -> impl Future<Output = Result<(), Error>> + Send;
    fn update_timestamp(&mut self, uuid: String, t: i64, force: bool) -> impl Future<Output = Result<(), Error>> + Send;

    // Remove operations (async for vqueue push)
    fn remove(&mut self, uuid: String) -> impl Future<Output = Result<(), Error>> + Send;
    fn remove_with_time(&mut self, uuid: String, t: i64) -> impl Future<Output = Result<(), Error>> + Send;
    fn remove_multiple(&mut self, uuids: Vec<String>) -> impl Future<Output = Result<(), Error>> + Send;
    fn remove_multiple_with_time(&mut self, uuids: Vec<String>, t: i64) -> impl Future<Output = Result<(), Error>> + Send;

    // Index management (async for I/O)
    fn regenerate_indexes(&mut self) -> impl Future<Output = Result<(), Error>> + Send;
    fn create_index(&mut self) -> impl Future<Output = Result<(), Error>> + Send;
    fn save_index(&mut self) -> impl Future<Output = Result<(), Error>> + Send;
    fn create_and_save_index(&mut self) -> impl Future<Output = Result<(), Error>> + Send;

    // Object retrieval (async for kvs/vqueue lookup)
    fn get_object(&self, uuid: String) -> impl Future<Output = Result<(Vec<f32>, i64), Error>> + Send;
    fn exists(&self, uuid: String) -> impl Future<Output = (usize, bool)> + Send;
    fn uuids(&self) -> impl Future<Output = Vec<String>> + Send;

    // List with callback (sync, but may need async variant in future)
    fn list_object_func<F: FnMut(String, Vec<f32>, i64) -> bool + Send>(&self, f: F) -> impl Future<Output = ()> + Send;

    // Status queries (sync - these are typically fast in-memory checks)
    fn is_indexing(&self) -> bool;
    fn is_flushing(&self) -> bool;
    fn is_saving(&self) -> bool;
    fn len(&self) -> u32;
    fn number_of_create_index_executions(&self) -> u64;
    fn insert_vqueue_buffer_len(&self) -> u32;
    fn delete_vqueue_buffer_len(&self) -> u32;
    fn get_dimension_size(&self) -> usize;
    fn broken_index_count(&self) -> u64;
    fn is_statistics_enabled(&self) -> bool;

    // Info queries (sync - typically fast)
    fn index_statistics(&self) -> Result<info::index::Statistics, Error>;
    fn index_property(&self) -> Result<info::index::Property, Error>;

    // Cleanup
    fn close(&mut self) -> impl Future<Output = Result<(), Error>> + Send;
}
