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

use thiserror::Error;

/// Custom error type for all operations within the `kvs`.
///
/// This enum consolidates errors from various sources into a single, well-defined
/// error type, making error handling for the library's users more straightforward.
#[derive(Error, Debug)]
pub enum Error {
    #[error("item not found")]
    NotFound,
    /// Errors from the underlying `sled` database.
    #[error("sled db error")]
    Sled(#[from] sled::Error),
    /// I/O errors.
    #[error("I/O error")]
    Io(#[from] std::io::Error),
    /// Errors that occur during a `sled` transaction, including conflicts or custom aborts.
    #[error("sled transaction error")]
    SledTransaction {
        #[source]
        source: Box<dyn std::error::Error + Send + Sync>,
    },
    /// Error during serialization or deserialization.
    #[error("codec error")]
    Codec {
        #[source]
        source: Box<dyn std::error::Error + Send + Sync>,
    },
    /// Error related to internal Tokio task management, typically from `spawn_blocking`.
    #[error("internal task error")]
    Internal(#[from] tokio::task::JoinError),
}
