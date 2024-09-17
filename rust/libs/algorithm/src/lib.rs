//
// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
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
use std::{error, fmt};

use anyhow::Result;
use proto::payload::v1::search;

#[derive(Debug)]
pub enum Error {
    CreateIndexingIsInProgress{},
    FlushingIsInProgress{},
    EmptySearchResult{},
    IncompatibleDimensionSize{got: usize, want: usize},

    Unknown{},
}

impl error::Error for Error {}

impl fmt::Display for Error {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        match self {
            Error::CreateIndexingIsInProgress{} => write!(f, "create indexing is in progress"),
            Error::FlushingIsInProgress{} => write!(f, "flush is in progress"),
            Error::EmptySearchResult{} => write!(f, "search result is empty"),
            Error::IncompatibleDimensionSize { got, want } => write!(f, "incompatible dimension size detected\trequested: {},\tconfigured: {}", got, want),
            Error::Unknown {  } => write!(f, "unknown error")
        }
    }
}

pub trait ANN: Send + Sync {
    fn get_dimension_size(&self) -> usize;
    fn search(&self, vector: Vec<f32>, dim: u32, epsilon: f32, radius: f32) -> Result<tonic::Response<search::Response>, Error>;
}
