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

/// Error types and helpers for ANN implementations.
pub mod error;
pub use error::{Error, MultiError};

use anyhow::Result;
use proto::payload::v1::{info, search};
use std::{collections::HashMap, future::Future, i64};

/// Trait for Approximate Nearest Neighbor (ANN) index implementations.
///
/// All methods that involve I/O or potentially blocking operations are async.
pub trait ANN: Send + Sync {
    // Search operations (async for potential I/O with vqueue/kvs)
    /// Searches for nearest neighbors by vector.
    fn search(
        &self,
        vector: Vec<f32>,
        k: u32,
        epsilon: f32,
        radius: f32,
    ) -> impl Future<Output = Result<search::Response, Error>> + Send;
    /// Searches for nearest neighbors by UUID.
    fn search_by_id(
        &self,
        uuid: String,
        k: u32,
        epsilon: f32,
        radius: f32,
    ) -> impl Future<Output = Result<search::Response, Error>> + Send;
    /// Performs a linear search by vector.
    fn linear_search(
        &self,
        vector: Vec<f32>,
        k: u32,
    ) -> impl Future<Output = Result<search::Response, Error>> + Send;
    /// Performs a linear search by UUID.
    fn linear_search_by_id(
        &self,
        uuid: String,
        k: u32,
    ) -> impl Future<Output = Result<search::Response, Error>> + Send;

    // Insert operations (async for vqueue push)
    /// Inserts a vector with a UUID.
    fn insert(
        &mut self,
        uuid: String,
        vector: Vec<f32>,
    ) -> impl Future<Output = Result<(), Error>> + Send;
    /// Inserts a vector with a UUID and timestamp.
    fn insert_with_time(
        &mut self,
        uuid: String,
        vector: Vec<f32>,
        t: i64,
    ) -> impl Future<Output = Result<(), Error>> + Send;
    /// Inserts multiple vectors.
    fn insert_multiple(
        &mut self,
        vectors: HashMap<String, Vec<f32>>,
    ) -> impl Future<Output = Result<(), Error>> + Send;
    /// Inserts multiple vectors with a shared timestamp.
    fn insert_multiple_with_time(
        &mut self,
        vectors: HashMap<String, Vec<f32>>,
        t: i64,
    ) -> impl Future<Output = Result<(), Error>> + Send;

    // Update operations (async for vqueue/kvs)
    /// Updates a vector by UUID.
    fn update(
        &mut self,
        uuid: String,
        vector: Vec<f32>,
    ) -> impl Future<Output = Result<(), Error>> + Send;
    /// Updates a vector by UUID with a timestamp.
    fn update_with_time(
        &mut self,
        uuid: String,
        vector: Vec<f32>,
        t: i64,
    ) -> impl Future<Output = Result<(), Error>> + Send;
    /// Updates multiple vectors.
    fn update_multiple(
        &mut self,
        vectors: HashMap<String, Vec<f32>>,
    ) -> impl Future<Output = Result<(), Error>> + Send;
    /// Updates multiple vectors with a shared timestamp.
    fn update_multiple_with_time(
        &mut self,
        vectors: HashMap<String, Vec<f32>>,
        t: i64,
    ) -> impl Future<Output = Result<(), Error>> + Send;
    /// Updates the timestamp for a UUID.
    fn update_timestamp(
        &mut self,
        uuid: String,
        t: i64,
        force: bool,
    ) -> impl Future<Output = Result<(), Error>> + Send;

    // Remove operations (async for vqueue push)
    /// Removes a vector by UUID.
    fn remove(&mut self, uuid: String) -> impl Future<Output = Result<(), Error>> + Send;
    /// Removes a vector by UUID with a timestamp.
    fn remove_with_time(
        &mut self,
        uuid: String,
        t: i64,
    ) -> impl Future<Output = Result<(), Error>> + Send;
    /// Removes multiple vectors.
    fn remove_multiple(
        &mut self,
        uuids: Vec<String>,
    ) -> impl Future<Output = Result<(), Error>> + Send;
    /// Removes multiple vectors with a shared timestamp.
    fn remove_multiple_with_time(
        &mut self,
        uuids: Vec<String>,
        t: i64,
    ) -> impl Future<Output = Result<(), Error>> + Send;

    // Index management (async for I/O)
    /// Regenerates indexes from persisted state.
    fn regenerate_indexes(&mut self) -> impl Future<Output = Result<(), Error>> + Send;
    /// Creates a new index from queued data.
    fn create_index(&mut self) -> impl Future<Output = Result<(), Error>> + Send;
    /// Saves the current index to storage.
    fn save_index(&mut self) -> impl Future<Output = Result<(), Error>> + Send;
    /// Creates and then saves an index.
    fn create_and_save_index(&mut self) -> impl Future<Output = Result<(), Error>> + Send;

    // Object retrieval (async for kvs/vqueue lookup)
    /// Returns an object by UUID.
    fn get_object(
        &self,
        uuid: String,
    ) -> impl Future<Output = Result<(Vec<f32>, i64), Error>> + Send;
    /// Returns whether a UUID exists and the associated object ID.
    fn exists(&self, uuid: String) -> impl Future<Output = (usize, bool)> + Send;
    /// Returns all UUIDs stored in the index.
    fn uuids(&self) -> impl Future<Output = Vec<String>> + Send;

    // List with callback (sync, but may need async variant in future)
    /// Iterates over objects, invoking a callback for each entry.
    fn list_object_func<F: FnMut(String, Vec<f32>, i64) -> bool + Send>(
        &self,
        f: F,
    ) -> impl Future<Output = ()> + Send;

    // Status queries (sync - these are typically fast in-memory checks)
    /// Returns true when indexing is in progress.
    fn is_indexing(&self) -> bool;
    /// Returns true when flushing is in progress.
    fn is_flushing(&self) -> bool;
    /// Returns true when saving is in progress.
    fn is_saving(&self) -> bool;
    /// Returns the number of indexed objects.
    fn len(&self) -> u32;
    /// Returns the total number of create-index executions.
    fn number_of_create_index_executions(&self) -> u64;
    /// Returns the insert vqueue buffer length.
    fn insert_vqueue_buffer_len(&self) -> u32;
    /// Returns the delete vqueue buffer length.
    fn delete_vqueue_buffer_len(&self) -> u32;
    /// Returns the configured dimension size.
    fn get_dimension_size(&self) -> usize;
    /// Returns the number of broken index backups.
    fn broken_index_count(&self) -> u64;
    /// Returns true if statistics collection is enabled.
    fn is_statistics_enabled(&self) -> bool;

    // Info queries (sync - typically fast)
    /// Returns index statistics.
    fn index_statistics(&self) -> Result<info::index::Statistics, Error>;
    /// Returns index property settings.
    fn index_property(&self) -> Result<info::index::Property, Error>;

    // Cleanup
    /// Closes the index and releases resources.
    fn close(&mut self) -> impl Future<Output = Result<(), Error>> + Send;
}
