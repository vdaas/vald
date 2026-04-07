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

/// Helper constructors for multi-object errors.
pub trait MultiError {
    /// Builds an error for UUIDs that already exist.
    fn new_uuid_already_exists(uuids: Vec<String>) -> Error;
    /// Builds an error for missing object IDs.
    fn new_object_id_not_found(uuids: Vec<String>) -> Error;
    /// Builds an error for invalid dimension sizes.
    fn new_invalid_dimension_size(current: Vec<String>, limit: Vec<String>) -> Error;
    /// Builds an error for missing UUIDs.
    fn new_uuid_not_found(uuids: Vec<String>) -> Error;
    /// Splits a comma-separated UUID list into a vector.
    fn split_uuids(uuids: String) -> Vec<String>;
}

/// Error types returned by ANN (Approximate Nearest Neighbor) operations.
///
/// This enum represents all possible error conditions that can occur during index construction,
/// search operations, and data management in the algorithm layer. Each variant corresponds to
/// a specific error condition with appropriate context information.
///
/// # Variants
///
/// * `CreateIndexingIsInProgress` - Index creation is currently running, operations must wait
/// * `EmptySearchResult` - Query returned no matching vectors
/// * `FlushingIsInProgress` - Flush operation is in progress, blocking concurrent operations
/// * `IncompatibleDimensionSize` - Query/insert vector dimension doesn't match index configuration
/// * `UUIDAlreadyExists` - Attempted to insert a vector with an already existing UUID
/// * `UUIDNotFound` - Requested UUID does not exist in the index
/// * `UncommittedIndexNotFound` - No uncommitted (pending) index operations found
/// * `InvalidUUID` - UUID format is invalid
/// * `InvalidDimensionSize` - Vector dimension size violates constraints
/// * `ObjectIDNotFound` - Object ID metadata lookup failed
/// * `WriteOperationToReadReplica` - Write operations are not allowed on read-only replicas
/// * `Unsupported` - Operation is not supported for the given algorithm
/// * `IndexNotFound` - Index does not exist or failed to load
/// * `InvalidTimestamp` - Timestamp value is invalid
/// * `NewerTimestampAlreadyExists` - UUID with a newer timestamp already exists (conflict)
/// * `Internal` - Wrapped internal error from underlying components
/// * `Unknown` - Unexpected error with no specific categorization
#[derive(thiserror::Error, Debug)]
pub enum Error {
    /// Index creation is currently in progress.
    ///
    /// Returned when attempting to perform operations that require an exclusive index lock
    /// while the index is being created.
    #[error("create indexing is in progress")]
    CreateIndexingIsInProgress {},

    /// Search operation returned no results.
    ///
    /// Indicates that the search completed successfully but found no matching vectors
    /// within the configured search parameters.
    #[error("search result is empty")]
    EmptySearchResult {},

    /// Flush operation is currently in progress.
    ///
    /// Returned when attempting operations that conflict with ongoing flush operations
    /// which persist pending changes to disk.
    #[error("flush is in progress")]
    FlushingIsInProgress {},

    /// Query vector dimension doesn't match the index configuration.
    ///
    /// Contains the actual dimension (`got`) and the expected dimension (`want`).
    #[error("incompatible dimension size detected\trequested: {got},\tconfigured: {want}")]
    IncompatibleDimensionSize { got: usize, want: usize },

    /// UUID already exists in the index.
    ///
    /// Attempted to insert or create an object with a UUID that is already indexed.
    #[error("uuid {uuid} index already exists")]
    UUIDAlreadyExists { uuid: String },

    /// UUID not found in the index.
    ///
    /// Requested UUID does not exist or has been deleted.
    #[error("object uuid{} not found", if uuid == "0" { "" } else { " {uuid}'s metadata" })]
    UUIDNotFound { uuid: String },

    /// No uncommitted indexes found.
    ///
    /// Returned when attempting to flush or finalize uncommitted changes but none exist.
    #[error("uncommitted indexes are not found")]
    UncommittedIndexNotFound {},

    /// UUID format is invalid.
    ///
    /// The provided UUID does not conform to the expected format.
    #[error("uuid \"{uuid}\" is invalid")]
    InvalidUUID { uuid: String },

    /// Vector dimension size is invalid.
    ///
    /// Dimension must be >= 2 and <= configured limit.
    /// Contains current dimension and the limit.
    #[error("dimension size {} is invalid, the supporting dimension size must be {}", current, if limit == "0" { "bigger than 2" } else { "between 2 ~ {limit}" })]
    InvalidDimensionSize { current: String, limit: String },

    /// Object ID not found in the index.
    ///
    /// The object metadata could not be retrieved.
    #[error("uuid {uuid}'s object id not found")]
    ObjectIDNotFound { uuid: String },

    /// Write operation attempted on a read-only replica.
    ///
    /// This instance is configured as a read replica and does not accept write operations.
    #[error("write operation to read replica is not possible")]
    WriteOperationToReadReplica {},

    /// Operation is not supported for the specified algorithm.
    ///
    /// Some operations may not be available for all algorithm implementations.
    /// Contains the operation method name and the algorithm name.
    #[error("{method} is not supported for {algorithm}")]
    Unsupported { method: String, algorithm: String },

    /// Index does not exist or could not be loaded.
    ///
    /// The requested index file is missing or corrupted.
    #[error("index not found")]
    IndexNotFound {},

    /// Timestamp value is invalid.
    ///
    /// The provided timestamp does not meet validity requirements.
    #[error("timestamp {timestamp} is invalid")]
    InvalidTimestamp { timestamp: i64 },

    /// UUID with a newer timestamp already exists.
    ///
    /// Conflict detected: an update attempt with an older timestamp for a UUID that already
    /// has a newer timestamp recorded.
    #[error("uuid {uuid}'s newer timestamp {timestamp} already exists")]
    NewerTimestampAlreadyExists { uuid: String, timestamp: i64 },

    /// Internal error from underlying components.
    ///
    /// Wraps errors from dependencies and internal subsystems.
    #[error("{0}")]
    Internal(#[from] Box<dyn std::error::Error + Send + Sync>),

    /// Unexpected error with no specific categorization.
    ///
    /// Indicates an error condition that doesn't fit other categories.
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

    fn new_invalid_dimension_size(current: Vec<String>, limit: Vec<String>) -> Error {
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
        uuids.split(',').map(|x| x.to_string()).collect()
    }
}
