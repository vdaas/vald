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
use anyhow::Result;
use proto::payload::v1::search;
use std::{error, fmt, i64};

#[derive(Debug)]
pub enum Error {
    CreateIndexingIsInProgress {},
    FlushingIsInProgress {},
    EmptySearchResult {},
    IncompatibleDimensionSize { got: usize, want: usize },
    UUIDAlreadyExists { uuid: String },
    UUIDNotFound { id: usize },
    UncommittedIndexNotFound {},
    InvalidUUID { uuid: String },
    ObjectIDNotFound { uuid: String },
    Unknown {},
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
            Error::UUIDNotFound { id } => {
                if *id == (0 as usize) {
                    write!(f, "object uuid not found")
                } else {
                    write!(f, "object uuid {}'s metadata not found", id)
                }
            }
            Error::UncommittedIndexNotFound {} => write!(f, "uncommitted indexes are not found"),
            Error::InvalidUUID { uuid } => write!(f, "uuid \"{}\" is invalid", uuid),
            Error::ObjectIDNotFound { uuid } => write!(f, "uuid {}'s object id not found", uuid),
            Error::Unknown {} => write!(f, "unknown error"),
        }
    }
}

pub trait ANN: Send + Sync {
    fn exists(&self, uuid: String) -> bool;
    fn create_index(&mut self) -> Result<(), Error>;
    fn save_index(&mut self) -> Result<(), Error>;
    fn insert(&mut self, uuid: String, vector: Vec<f32>, ts: i64) -> Result<(), Error>;
    fn update(&mut self, uuid: String, vector: Vec<f32>, ts: i64) -> Result<(), Error>;
    fn remove(&mut self, uuid: String, ts: i64) -> Result<(), Error>;
    fn search(
        &self,
        vector: Vec<f32>,
        k: u32,
        epsilon: f32,
        radius: f32,
    ) -> Result<search::Response, Error>;
    fn get_object(&self, uuid: String) -> Result<(Vec<f32>, i64), Error>;
    fn get_dimension_size(&self) -> usize;
    fn len(&self) -> u32;
    fn insert_vqueue_buffer_len(&self) -> u32;
    fn delete_vqueue_buffer_len(&self) -> u32;
    fn is_indexing(&self) -> bool;
    fn is_saving(&self) -> bool;
}
