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

//! # Persistent Vector Queue
//!
//! A persistent, async-safe, and concurrent vector queue implementation.
//! It is designed for high-concurrency environments, offloading blocking I/O operations
//! to a dedicated thread pool managed by `tokio`.
//!
//! The implementation uses `sled` as its underlying persistent storage engine to leverage
//! its robust transactional capabilities, ensuring data consistency.
use async_trait::async_trait;
use sled::{
    Tree,
    transaction::{
        ConflictableTransactionError, TransactionError, Transactional, UnabortableTransactionError,
    },
};
use std::{
    collections::HashMap,
    convert::TryInto,
    pin::Pin,
    str,
    sync::{
        Arc,
        atomic::{AtomicU64, Ordering},
    },
    time::{SystemTime, UNIX_EPOCH},
};
use thiserror::Error;
use tokio::sync::mpsc;
use tokio_stream::{Stream, wrappers::ReceiverStream};

/// A tuple representing a queue item: (UUID, Vector, Timestamp).
pub type QueueItem = (String, Vec<f32>, i64);

type RangeStream = Pin<Box<dyn Stream<Item = Result<QueueItem, QueueError>> + Send>>;

/// Error type representing possible failures for queue operations.
#[derive(Debug, Error)]
pub enum QueueError {
    /// Error returned when the provided UUID is invalid or empty.
    #[error("The provided UUID is invalid or empty.")]
    InvalidUuid,
    /// Error returned for `sled` database-related failures.
    #[error("Sled database error")]
    Sled(#[from] sled::Error),
    /// Error returned for serialization failures.
    #[error("Codec encode error: {0}")]
    Encode(#[from] wincode::error::WriteError),
    /// Error returned for deserialization failures.
    #[error("Codec decode error: {0}")]
    Decode(#[from] wincode::error::ReadError),
    /// Error returned for internal Tokio task failures.
    #[error("Internal Tokio task error")]
    Internal(#[from] tokio::task::JoinError),
    /// Error returned for UTF-8 conversion failures.
    #[error("UTF-8 conversion error")]
    Utf8(#[from] str::Utf8Error),
    /// Error returned for key parsing failures.
    #[error("Key parsing error: {0}")]
    KeyParse(String),
    /// Error returned for `sled` unabortable transaction failures.
    #[error("Sled unabortable transaction error")]
    Unabortable(#[from] UnabortableTransactionError),
    /// Error returned when the requested UUID is not found in the queue.
    #[error("UUID not found in queue: {0}")]
    NotFound(String),
}

/// Represents an item drained from the queue.
#[derive(Debug, PartialEq, Clone)]
pub enum DrainItem {
    /// An insert operation with the UUID and vector.
    Insert(String, Vec<f32>),
    /// A delete operation with the UUID.
    Delete(String),
}

/// Trait definition for queue functionality.
#[async_trait]
pub trait Queue: Send + Sync {
    /// Pushes an insert/update operation for a vector into the queue.
    /// If an entry with the same UUID already exists, it will be replaced.
    ///
    /// # Arguments
    ///
    /// * `uuid` - The UUID of the vector.
    /// * `vector` - The vector to be inserted.
    /// * `timestamp` - An optional timestamp for the operation. If `None`, the current time is used.
    async fn push_insert(
        &self,
        uuid: impl AsRef<str> + Send,
        vector: Vec<f32>,
        timestamp: Option<i64>,
    ) -> Result<(), QueueError>;

    /// Pushes a delete operation for a vector into the queue.
    ///
    /// # Arguments
    ///
    /// * `uuid` - The UUID of the vector to be deleted.
    /// * `timestamp` - An optional timestamp for the operation. If `None`, the current time is used.
    async fn push_delete(
        &self,
        uuid: impl AsRef<str> + Send,
        timestamp: Option<i64>,
    ) -> Result<(), QueueError>;

    /// Pops and removes an insert operation from the queue by UUID.
    /// This is a destructive operation that removes the entry from the insert queue.
    ///
    /// # Arguments
    ///
    /// * `uuid` - The UUID of the vector to pop.
    ///
    /// # Returns
    ///
    /// A tuple of (vector, timestamp) if the UUID exists in the insert queue.
    async fn pop_insert(&self, uuid: impl AsRef<str> + Send)
    -> Result<(Vec<f32>, i64), QueueError>;

    /// Pops and removes a delete operation from the queue by UUID.
    /// This is a destructive operation that removes the entry from the delete queue.
    ///
    /// # Arguments
    ///
    /// * `uuid` - The UUID of the delete operation to pop.
    ///
    /// # Returns
    ///
    /// The timestamp of the delete operation if the UUID exists in the delete queue.
    async fn pop_delete(&self, uuid: impl AsRef<str> + Send) -> Result<i64, QueueError>;

    /// Checks if a UUID exists in the insert queue and returns its timestamp.
    /// This is a non-destructive read operation.
    ///
    /// # Arguments
    ///
    /// * `uuid` - The UUID to check.
    ///
    /// # Returns
    ///
    /// The insert timestamp if the UUID exists, or 0 if not found.
    async fn iv_exists(&self, uuid: impl AsRef<str> + Send) -> Result<i64, QueueError>;

    /// Checks if a UUID exists in the delete queue and returns its timestamp.
    /// This is a non-destructive read operation.
    ///
    /// # Arguments
    ///
    /// * `uuid` - The UUID to check.
    ///
    /// # Returns
    ///
    /// The delete timestamp if the UUID exists, or 0 if not found.
    async fn dv_exists(&self, uuid: impl AsRef<str> + Send) -> Result<i64, QueueError>;

    /// Returns the vector stored in the queue.
    /// If the same UUID exists in both the insert queue and the delete queue,
    /// the timestamp is compared and the vector is returned only if the insert timestamp is newer.
    ///
    /// # Arguments
    ///
    /// * `uuid` - The UUID of the vector to retrieve.
    ///
    /// # Returns
    ///
    /// A tuple of (vector, insert_timestamp, exists).
    async fn get_vector(&self, uuid: impl AsRef<str> + Send)
    -> Result<(Vec<f32>, i64), QueueError>;

    /// Returns the vector and both timestamps stored in the queue.
    /// This method returns both insert and delete timestamps, allowing the caller
    /// to determine the state of the vector.
    ///
    /// # Arguments
    ///
    /// * `uuid` - The UUID of the vector to retrieve.
    ///
    /// # Returns
    ///
    /// A tuple of (vector, insert_timestamp, delete_timestamp, exists).
    /// - `exists` is true if the vector is valid (insert timestamp > delete timestamp)
    /// - Even if `exists` is false, delete_timestamp may be non-zero if a delete is pending
    async fn get_vector_with_timestamp(
        &self,
        uuid: impl AsRef<str> + Send,
    ) -> Result<(Option<Vec<f32>>, i64, i64, bool), QueueError>;

    /// Returns a stream that drains both the insert and delete queues up to the given timestamp.
    ///
    /// It resolves conflicts between inserts and deletes, yielding a stream of `DrainItem`s.
    /// This is a destructive operation that removes the drained entries from the queues.
    ///
    /// # Arguments
    ///
    /// * `now` - The timestamp up to which the queues should be drained.
    /// * `batch_size` - The number of items to drain in each batch.
    fn drain_queues(
        &self,
        now: i64,
        batch_size: usize,
    ) -> Pin<Box<dyn Stream<Item = Result<DrainItem, QueueError>> + Send>>;

    /// Returns the number of vectors in the insert queue.
    fn ivq_len(&self) -> u64;

    /// Returns the number of vectors in the delete queue.
    fn dvq_len(&self) -> u64;

    /// Iterates over all items in the insert queue, filtering out items that have a newer delete.
    /// This is a non-destructive operation that does not modify the queue.
    /// Returns a stream of (uuid, vector, timestamp) tuples for each valid item.
    fn range(&self) -> RangeStream;
}

/// A persistent queue implementation using `sled`.
#[derive(Clone)]
pub struct PersistentQueue {
    /// The `sled` tree for the insert queue.
    insert_queue: Tree,
    /// The `sled` tree for the delete queue.
    delete_queue: Tree,
    /// The `sled` tree for the insert index.
    insert_index: Tree,
    /// The `sled` tree for the delete index.
    delete_index: Tree,
    /// The number of items in the insert queue.
    insert_count: Arc<AtomicU64>,
    /// The number of items in the delete queue.
    delete_count: Arc<AtomicU64>,
}

/// A builder for creating a `PersistentQueue`.
pub struct Builder {
    /// The path to the `sled` database.
    path: String,
    /// The configuration for the `sled` database.
    config: sled::Config,
}

impl Builder {
    /// Creates a new builder for a `PersistentQueue`.
    ///
    /// # Arguments
    ///
    /// * `path` - The path to the `sled` database.
    pub fn new(path: impl AsRef<str>) -> Self {
        Self {
            path: path.as_ref().to_string(),
            config: sled::Config::default(),
        }
    }

    /// Sets the cache capacity for the `sled` database.
    ///
    /// # Arguments
    ///
    /// * `capacity` - The cache capacity in bytes.
    pub fn cache_capacity(mut self, capacity: u64) -> Self {
        self.config = self.config.cache_capacity(capacity);
        self
    }

    /// Enables or disables compression for the `sled` database.
    ///
    /// # Arguments
    ///
    /// * `use_compression` - Whether to use compression.
    pub fn use_compression(mut self, use_compression: bool) -> Self {
        self.config = self.config.use_compression(use_compression);
        self
    }

    /// Builds the `PersistentQueue`.
    /// This method opens the `sled` database and its trees, and initializes the queue counts.
    pub async fn build(self) -> Result<PersistentQueue, QueueError> {
        // Configure the database path.
        let config = self.config.path(self.path);
        // Open the database in a blocking task to avoid blocking the async runtime.
        let db = tokio::task::spawn_blocking(move || config.open()).await??;
        // Open the required trees for the queues and index.
        let insert_queue = db.open_tree("insert_queue")?;
        let delete_queue = db.open_tree("delete_queue")?;
        let insert_index = db.open_tree("insert_index")?;
        let delete_index = db.open_tree("delete_index")?;
        // Initialize the atomic counters with the current number of items in the queues.
        let insert_count = Arc::new(AtomicU64::new(insert_queue.len() as u64));
        let delete_count = Arc::new(AtomicU64::new(delete_queue.len() as u64));

        Ok(PersistentQueue {
            insert_queue,
            delete_queue,
            insert_index,
            delete_index,
            insert_count,
            delete_count,
        })
    }
}

impl PersistentQueue {
    /// Returns the current time in nanoseconds since the UNIX epoch.
    fn now_ns() -> i64 {
        SystemTime::now()
            .duration_since(UNIX_EPOCH)
            .unwrap_or_default()
            .as_nanos() as i64
    }

    /// Creates a key for the `sled` database from a timestamp and a UUID.
    fn create_key(timestamp: i64, uuid: &str) -> Vec<u8> {
        let uuid_bytes = uuid.as_bytes();
        let mut buf = Vec::with_capacity(8 + uuid_bytes.len());
        buf.extend_from_slice(&timestamp.to_be_bytes());
        buf.extend_from_slice(uuid_bytes);
        buf
    }

    /// Parses a key from the `sled` database into a timestamp and a UUID.
    ///
    /// # Arguments
    ///
    /// * `key` - The key to parse.
    fn parse_key(key: &[u8]) -> Result<(i64, String), QueueError> {
        if key.len() < 8 {
            return Err(QueueError::KeyParse("Key is too short".to_string()));
        }
        let ts_bytes: [u8; 8] = key[0..8]
            .try_into()
            .map_err(|_| QueueError::KeyParse("Failed to parse timestamp bytes".to_string()))?;
        let ts = i64::from_be_bytes(ts_bytes);
        let uuid = str::from_utf8(&key[8..])?.to_string();
        Ok((ts, uuid))
    }

    /// Drains a batch of items from the insert and delete queues atomically.
    /// This method resolves conflicts, removes processed items from storage,
    /// and returns them sorted by timestamp.
    async fn drain_batch(&self, now: i64, batch_size: usize) -> Result<Vec<DrainItem>, QueueError> {
        // Clone Tree handles (cheap Arc clone) to move into the blocking task.
        let iq = self.insert_queue.clone();
        let dq = self.delete_queue.clone();
        let ii = self.insert_index.clone();
        let di = self.delete_index.clone();
        let ic = self.insert_count.clone();
        let dc = self.delete_count.clone();

        tokio::task::spawn_blocking(move || {
            let end_key = (now + 1).to_be_bytes();
            // Step 1: Collect all delete operations within the timestamp range.
            let mut deletes: HashMap<String, (i64, Vec<u8>)> = HashMap::new();
            for item in dq.range(..end_key.as_ref()).take(batch_size) {
                let (key, _) = item?;
                let (ts, uuid) = Self::parse_key(&key)?;
                deletes.insert(uuid, (ts, key.to_vec()));
            }
            // Step 2: Collect insert operations and resolve conflicts against the collected deletes.
            let mut final_items: Vec<(i64, DrainItem)> = Vec::new();
            let mut inserts_to_remove: Vec<(Vec<u8>, Vec<u8>)> = Vec::new();
            let mut deletes_to_remove_db: Vec<(Vec<u8>, Vec<u8>)> = Vec::new();

            for item in iq.range(..end_key.as_ref()).take(batch_size) {
                let (key, val) = item?;
                let (its, uuid) = Self::parse_key(&key)?;
                inserts_to_remove.push((key.to_vec(), uuid.as_bytes().to_vec()));

                if let Some((dts, delete_key)) = deletes.get(&uuid) {
                    if *dts >= its {
                        continue;
                    } else {
                        deletes_to_remove_db.push((delete_key.clone(), uuid.as_bytes().to_vec()));
                        deletes.remove(&uuid);
                    }
                }
                let vec = wincode::deserialize(&val)?;
                final_items.push((its, DrainItem::Insert(uuid, vec)));
            }
            // Step 3: Collect the remaining valid deletes.
            deletes_to_remove_db.extend(
                deletes
                    .into_iter()
                    .map(|(uuid, (ts, key))| {
                        final_items.push((ts, DrainItem::Delete(uuid.clone())));
                        (key, uuid.into_bytes())
                    })
                    .collect::<Vec<(Vec<u8>, Vec<u8>)>>(),
            );
            // Step 4: Atomically apply all removals in a single database transaction.
            if !inserts_to_remove.is_empty() || !deletes_to_remove_db.is_empty() {
                (&iq, &dq, &ii, &di)
                    .transaction(|(iq_tx, dq_tx, ii_tx, di_tx)| {
                        for (key, uuid_bytes) in &inserts_to_remove {
                            iq_tx.remove(key.as_slice())?;
                            ii_tx.remove(uuid_bytes.as_slice())?;
                        }
                        for (key, uuid_bytes) in &deletes_to_remove_db {
                            dq_tx.remove(key.as_slice())?;
                            di_tx.remove(uuid_bytes.as_slice())?;
                        }
                        Ok(())
                    })
                    .map_err(|e| match e {
                        TransactionError::Abort(e) => QueueError::Sled(e),
                        TransactionError::Storage(e) => QueueError::Sled(e),
                    })?;
            }
            // Step 5: Update the in-memory queue counts.
            ic.fetch_sub(inserts_to_remove.len() as u64, Ordering::Relaxed);
            dc.fetch_sub(deletes_to_remove_db.len() as u64, Ordering::Relaxed);
            // Step 6: Combine and sort the final list of operations by timestamp.
            final_items.sort_by_key(|item| item.0);
            Ok(final_items.into_iter().map(|item| item.1).collect())
        })
        .await?
    }

    /// Internal helper to push an operation to the queue transactionally.
    async fn push_internal(
        &self,
        uuid: &str,
        timestamp: Option<i64>,
        queue: &Tree,
        index: &Tree,
        counter: &Arc<AtomicU64>,
        value: Vec<u8>,
    ) -> Result<(), QueueError> {
        if uuid.trim().is_empty() {
            return Err(QueueError::InvalidUuid);
        }
        let ts = timestamp.unwrap_or_else(Self::now_ns);
        let key = Self::create_key(ts, uuid);
        let uuid_bytes = uuid.as_bytes().to_vec();

        let q = queue.clone();
        let i = index.clone();
        let c = counter.clone();

        tokio::task::spawn_blocking(move || {
            (&q, &i)
                .transaction(|(q_tx, i_tx)| {
                    let to_abortable = |e| ConflictableTransactionError::Abort(e);
                    if let Some(old_ts_bytes) =
                        i_tx.insert(uuid_bytes.as_slice(), &ts.to_be_bytes())?
                    {
                        let old_ts_bytes_arr: [u8; 8] = old_ts_bytes
                            .as_ref()
                            .try_into()
                            .map_err(|_| {
                                QueueError::KeyParse("Invalid timestamp in index".to_string())
                            })
                            .map_err(to_abortable)?;
                        let old_ts = i64::from_be_bytes(old_ts_bytes_arr);
                        let uuid_str = str::from_utf8(&uuid_bytes)
                            .map_err(QueueError::from)
                            .map_err(to_abortable)?;
                        let old_key = Self::create_key(old_ts, uuid_str);
                        q_tx.remove(old_key.as_slice())?;
                    } else {
                        c.fetch_add(1, Ordering::Relaxed);
                    }
                    q_tx.insert(key.as_slice(), value.as_slice())?;
                    Ok(())
                })
                .map_err(|e| match e {
                    TransactionError::Abort(qe) => qe,
                    TransactionError::Storage(sled_err) => QueueError::Sled(sled_err),
                })
        })
        .await?
    }

    /// Loads a vector from the insert queue without removing it.
    /// Returns (vector, timestamp) if found.
    async fn load_ivq(&self, uuid: &str) -> Result<(Vec<f32>, i64), QueueError> {
        let uuid_bytes = uuid.as_bytes().to_vec();
        let uuid_string = uuid.to_string();
        let index = self.insert_index.clone();
        let queue = self.insert_queue.clone();

        tokio::task::spawn_blocking(move || {
            // Get timestamp from index
            let ts_bytes = match index.get(&uuid_bytes)? {
                Some(bytes) => bytes,
                None => return Err(QueueError::NotFound(uuid_string)),
            };

            let ts_bytes_arr: [u8; 8] = ts_bytes
                .as_ref()
                .try_into()
                .map_err(|_| QueueError::KeyParse("Invalid timestamp in index".to_string()))?;
            let ts = i64::from_be_bytes(ts_bytes_arr);

            // Get vector from queue
            let uuid_str = str::from_utf8(&uuid_bytes)?;
            let key = Self::create_key(ts, uuid_str);
            let value = match queue.get(&key)? {
                Some(bytes) => bytes,
                None => return Err(QueueError::NotFound(uuid_string)),
            };

            let vec = wincode::deserialize(&value)?;
            Ok((vec, ts))
        })
        .await?
    }

    /// Loads a timestamp from the delete queue without removing it.
    /// Returns the timestamp if found.
    async fn load_dvq(&self, uuid: &str) -> Result<i64, QueueError> {
        let uuid_bytes = uuid.as_bytes().to_vec();
        let uuid_string = uuid.to_string();
        let index = self.delete_index.clone();

        tokio::task::spawn_blocking(move || match index.get(&uuid_bytes)? {
            Some(ts_bytes) => {
                let ts_bytes_arr: [u8; 8] = ts_bytes
                    .as_ref()
                    .try_into()
                    .map_err(|_| QueueError::KeyParse("Invalid timestamp in index".to_string()))?;
                Ok(i64::from_be_bytes(ts_bytes_arr))
            }
            None => Err(QueueError::NotFound(uuid_string)),
        })
        .await?
    }

    /// Internal implementation of get_vector with timestamp.
    /// If enable_delete_timestamp is false, delete timestamp information is not returned.
    async fn get_vector_internal(
        &self,
        uuid: &str,
        enable_delete_timestamp: bool,
    ) -> Result<(Option<Vec<f32>>, i64, i64, bool), QueueError> {
        // Try to load from insert queue
        let ivq_result = self.load_ivq(uuid).await;

        match ivq_result {
            Ok((vec, its)) => {
                // Vector exists in insert queue, check delete queue
                let dts = match self.load_dvq(uuid).await {
                    Ok(ts) => ts,
                    Err(QueueError::NotFound(_)) => 0,
                    Err(e) => return Err(e),
                };

                if dts == 0 {
                    // Not in delete queue, vector exists
                    Ok((Some(vec), its, 0, true))
                } else {
                    // Both queues have the UUID, compare timestamps
                    // Vector exists if insert timestamp is newer than delete timestamp
                    let exists = its > dts;
                    Ok((Some(vec), its, dts, exists))
                }
            }
            Err(QueueError::NotFound(_)) => {
                // Not in insert queue
                if !enable_delete_timestamp {
                    // Don't check delete queue, just return not found
                    return Ok((None, 0, 0, false));
                }

                // Check delete queue
                let dts = match self.load_dvq(uuid).await {
                    Ok(ts) => ts,
                    Err(QueueError::NotFound(_)) => {
                        // Not in either queue
                        return Ok((None, 0, 0, false));
                    }
                    Err(e) => return Err(e),
                };

                // In delete queue but not insert queue
                Ok((None, 0, dts, false))
            }
            Err(e) => Err(e),
        }
    }

    /// Internal helper to pop an item from a queue by UUID.
    /// Returns the value bytes and timestamp if found.
    async fn pop_internal(
        &self,
        uuid: &str,
        queue: &Tree,
        index: &Tree,
        counter: &Arc<AtomicU64>,
    ) -> Result<(Vec<u8>, i64), QueueError> {
        if uuid.trim().is_empty() {
            return Err(QueueError::InvalidUuid);
        }
        let uuid_bytes = uuid.as_bytes().to_vec();
        let uuid_string = uuid.to_string();

        let q = queue.clone();
        let i = index.clone();
        let c = counter.clone();

        tokio::task::spawn_blocking(move || {
            (&q, &i)
                .transaction(|(q_tx, i_tx)| {
                    let to_abortable = |e| ConflictableTransactionError::Abort(e);
                    // Get the timestamp from the index
                    let ts_bytes = i_tx
                        .remove(uuid_bytes.as_slice())?
                        .ok_or_else(|| QueueError::NotFound(uuid_string.clone()))
                        .map_err(to_abortable)?;

                    let ts_bytes_arr: [u8; 8] = ts_bytes
                        .as_ref()
                        .try_into()
                        .map_err(|_| QueueError::KeyParse("Invalid timestamp in index".to_string()))
                        .map_err(to_abortable)?;
                    let ts = i64::from_be_bytes(ts_bytes_arr);

                    // Create the key and remove from queue
                    let uuid_str = str::from_utf8(&uuid_bytes)
                        .map_err(QueueError::from)
                        .map_err(to_abortable)?;
                    let key = Self::create_key(ts, uuid_str);

                    let value = q_tx
                        .remove(key.as_slice())?
                        .ok_or_else(|| QueueError::NotFound(uuid_string.clone()))
                        .map_err(to_abortable)?;

                    c.fetch_sub(1, Ordering::Relaxed);

                    Ok((value.to_vec(), ts))
                })
                .map_err(|e| match e {
                    TransactionError::Abort(qe) => qe,
                    TransactionError::Storage(sled_err) => QueueError::Sled(sled_err),
                })
        })
        .await?
    }
}

#[async_trait]
impl Queue for PersistentQueue {
    /// Pushes an insert/update operation for a vector into the queue.
    async fn push_insert(
        &self,
        uuid: impl AsRef<str> + Send,
        vector: Vec<f32>,
        timestamp: Option<i64>,
    ) -> Result<(), QueueError> {
        let value = wincode::serialize(&vector)?;
        self.push_internal(
            uuid.as_ref(),
            timestamp,
            &self.insert_queue,
            &self.insert_index,
            &self.insert_count,
            value,
        )
        .await
    }

    /// Pushes a delete operation for a vector into the queue.
    async fn push_delete(
        &self,
        uuid: impl AsRef<str> + Send,
        timestamp: Option<i64>,
    ) -> Result<(), QueueError> {
        self.push_internal(
            uuid.as_ref(),
            timestamp,
            &self.delete_queue,
            &self.delete_index,
            &self.delete_count,
            Vec::new(),
        )
        .await
    }

    /// Returns a stream that drains both the insert and delete queues up to the given timestamp.
    fn drain_queues(
        &self,
        now: i64,
        batch_size: usize,
    ) -> Pin<Box<dyn Stream<Item = Result<DrainItem, QueueError>> + Send>> {
        let (tx, rx) = mpsc::channel(batch_size);
        let q = self.clone();

        tokio::spawn(async move {
            loop {
                match q.drain_batch(now, batch_size).await {
                    Ok(items) if items.is_empty() => break,
                    Ok(items) => {
                        for item in items {
                            if tx.send(Ok(item)).await.is_err() {
                                return;
                            }
                        }
                    }
                    Err(e) => {
                        let _ = tx.send(Err(e)).await;
                        return;
                    }
                }
            }
        });

        Box::pin(ReceiverStream::new(rx))
    }

    /// Returns the number of vectors in the insert queue.
    fn ivq_len(&self) -> u64 {
        self.insert_count.load(Ordering::Acquire)
    }

    /// Returns the number of vectors in the delete queue.
    fn dvq_len(&self) -> u64 {
        self.delete_count.load(Ordering::Acquire)
    }

    /// Iterates over all items in the insert queue, filtering out items that have a newer delete.
    fn range(&self) -> RangeStream {
        let (tx, rx) = mpsc::channel(64);
        let iq = self.insert_queue.clone();
        let di = self.delete_index.clone();

        tokio::spawn(async move {
            let result = tokio::task::spawn_blocking(move || {
                let mut items = Vec::new();
                for (key, val) in iq.iter().flatten() {
                    if let Ok((its, uuid)) = Self::parse_key(&key) {
                        // Check if there's a newer delete for this uuid
                        let skip = di
                            .get(uuid.as_bytes())
                            .ok()
                            .flatten()
                            .filter(|dts_bytes| dts_bytes.len() >= 8)
                            .map(|dts_bytes| {
                                let dts_arr: [u8; 8] =
                                    dts_bytes[0..8].try_into().unwrap_or_default();
                                let dts = i64::from_be_bytes(dts_arr);
                                dts >= its
                            })
                            .unwrap_or(false);

                        if skip {
                            continue;
                        }
                        // Decode the vector
                        if let Ok(vec) = wincode::deserialize(&val) {
                            items.push((uuid, vec, its));
                        }
                    }
                }
                items
            })
            .await;

            match result {
                Ok(items) => {
                    for item in items {
                        if tx.send(Ok(item)).await.is_err() {
                            break;
                        }
                    }
                }
                Err(e) => {
                    let _ = tx.send(Err(QueueError::Internal(e))).await;
                }
            }
        });

        Box::pin(ReceiverStream::new(rx))
    }

    /// Pops an insert operation from the queue by UUID.
    /// Returns the vector and timestamp if found.
    async fn pop_insert(
        &self,
        uuid: impl AsRef<str> + Send,
    ) -> Result<(Vec<f32>, i64), QueueError> {
        let (value_bytes, ts) = self
            .pop_internal(
                uuid.as_ref(),
                &self.insert_queue,
                &self.insert_index,
                &self.insert_count,
            )
            .await?;

        let vec = wincode::deserialize(&value_bytes)?;
        Ok((vec, ts))
    }

    /// Pops a delete operation from the queue by UUID.
    /// Returns the timestamp if found.
    async fn pop_delete(&self, uuid: impl AsRef<str> + Send) -> Result<i64, QueueError> {
        let (_, ts) = self
            .pop_internal(
                uuid.as_ref(),
                &self.delete_queue,
                &self.delete_index,
                &self.delete_count,
            )
            .await?;
        Ok(ts)
    }

    /// Checks if a UUID exists in the insert queue and returns its timestamp.
    async fn iv_exists(&self, uuid: impl AsRef<str> + Send) -> Result<i64, QueueError> {
        let uuid_bytes = uuid.as_ref().as_bytes().to_vec();
        let uuid_string = uuid.as_ref().to_string();
        let index = self.insert_index.clone();

        tokio::task::spawn_blocking(move || match index.get(&uuid_bytes)? {
            Some(ts_bytes) => {
                let ts_bytes_arr: [u8; 8] = ts_bytes
                    .as_ref()
                    .try_into()
                    .map_err(|_| QueueError::KeyParse("Invalid timestamp in index".to_string()))?;
                Ok(i64::from_be_bytes(ts_bytes_arr))
            }
            None => Err(QueueError::NotFound(uuid_string)),
        })
        .await?
    }

    /// Checks if a UUID exists in the delete queue and returns its timestamp.
    async fn dv_exists(&self, uuid: impl AsRef<str> + Send) -> Result<i64, QueueError> {
        let uuid_bytes = uuid.as_ref().as_bytes().to_vec();
        let uuid_string = uuid.as_ref().to_string();
        let index = self.delete_index.clone();

        tokio::task::spawn_blocking(move || match index.get(&uuid_bytes)? {
            Some(ts_bytes) => {
                let ts_bytes_arr: [u8; 8] = ts_bytes
                    .as_ref()
                    .try_into()
                    .map_err(|_| QueueError::KeyParse("Invalid timestamp in index".to_string()))?;
                Ok(i64::from_be_bytes(ts_bytes_arr))
            }
            None => Err(QueueError::NotFound(uuid_string)),
        })
        .await?
    }

    /// Returns the vector stored in the queue.
    /// If the same UUID exists in both the insert queue and the delete queue,
    /// the timestamp is compared and the vector is returned only if the insert timestamp is newer.
    async fn get_vector(
        &self,
        uuid: impl AsRef<str> + Send,
    ) -> Result<(Vec<f32>, i64), QueueError> {
        let (vec_opt, its, _dts, exists) = self.get_vector_internal(uuid.as_ref(), false).await?;

        if !exists {
            return Err(QueueError::NotFound(uuid.as_ref().to_string()));
        }

        match vec_opt {
            Some(vec) => Ok((vec, its)),
            None => Err(QueueError::NotFound(uuid.as_ref().to_string())),
        }
    }

    /// Returns the vector and both timestamps stored in the queue.
    /// This method returns both insert and delete timestamps, allowing the caller
    /// to determine the state of the vector.
    async fn get_vector_with_timestamp(
        &self,
        uuid: impl AsRef<str> + Send,
    ) -> Result<(Option<Vec<f32>>, i64, i64, bool), QueueError> {
        self.get_vector_internal(uuid.as_ref(), true).await
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use std::fs;
    use tokio::task::JoinSet;
    use tokio_stream::StreamExt;

    struct TestGuard {
        path: String,
    }

    impl Drop for TestGuard {
        fn drop(&mut self) {
            let _ = fs::remove_dir_all(&self.path);
        }
    }

    async fn setup(test_name: &str) -> (PersistentQueue, TestGuard) {
        let path = format!("./test_db_queue_{}", test_name);
        let _ = fs::remove_dir_all(&path);
        let guard = TestGuard { path: path.clone() };
        let queue = Builder::new(&path)
            .cache_capacity(1024 * 1024)
            .build()
            .await
            .unwrap();
        (queue, guard)
    }

    #[tokio::test]
    async fn test_push_and_drain() {
        let (q, _guard) = setup("push_and_drain").await;
        let vec = vec![1.0, 2.0];
        q.push_insert("key1", vec.clone(), Some(100)).await.unwrap();
        assert_eq!(q.ivq_len(), 1);

        let mut stream = q.drain_queues(200, 10);
        let item = stream.next().await.unwrap().unwrap();
        assert_eq!(item, DrainItem::Insert("key1".into(), vec));
        assert!(stream.next().await.is_none());
        assert_eq!(q.ivq_len(), 0);
    }

    #[tokio::test]
    async fn test_delete_newer_than_insert() {
        let (q, _guard) = setup("delete_newer_than_insert").await;
        q.push_insert("key1", vec![1.0], Some(100)).await.unwrap();
        q.push_delete("key1", Some(200)).await.unwrap();

        let mut stream = q.drain_queues(300, 10);
        let item = stream.next().await.unwrap().unwrap();
        assert_eq!(item, DrainItem::Delete("key1".into()));
        assert!(stream.next().await.is_none());
        assert_eq!(q.ivq_len(), 0);
        assert_eq!(q.dvq_len(), 0);
    }

    #[tokio::test]
    async fn test_insert_newer_than_delete() {
        let (q, _guard) = setup("insert_newer_than_delete").await;
        q.push_delete("key1", Some(100)).await.unwrap();
        q.push_insert("key1", vec![1.0], Some(200)).await.unwrap();

        let mut stream = q.drain_queues(300, 10);
        let item = stream.next().await.unwrap().unwrap();
        assert_eq!(item, DrainItem::Insert("key1".into(), vec![1.0]));
        assert!(stream.next().await.is_none());
        assert_eq!(q.ivq_len(), 0);
        assert_eq!(q.dvq_len(), 0);
    }

    #[tokio::test]
    async fn test_drain_in_batches() {
        let (q, _guard) = setup("drain_in_batches").await;
        q.push_insert("key1", vec![1.0], Some(100)).await.unwrap();
        q.push_insert("key2", vec![2.0], Some(200)).await.unwrap();
        q.push_delete("key3", Some(150)).await.unwrap();

        let mut stream = q.drain_queues(300, 1);
        let mut items = Vec::new();
        while let Some(item) = stream.next().await {
            items.push(item.unwrap());
        }

        assert_eq!(items.len(), 3);
        assert_eq!(items[0], DrainItem::Insert("key1".into(), vec![1.0]));
        assert_eq!(items[1], DrainItem::Delete("key3".into()));
        assert_eq!(items[2], DrainItem::Insert("key2".into(), vec![2.0]));
        assert_eq!(q.ivq_len(), 0);
        assert_eq!(q.dvq_len(), 0);
    }

    #[tokio::test(flavor = "multi_thread", worker_threads = 10)]
    async fn test_concurrent_pushes_and_drain() {
        let (q, _guard) = setup("concurrent_pushes_and_drain").await;
        let queue = Arc::new(q);
        let mut tasks = JoinSet::new();
        let num_tasks = 100;

        for i in 0..num_tasks {
            let q_clone = queue.clone();
            tasks.spawn(async move {
                if i % 2 == 0 {
                    let vec = vec![i as f32];
                    q_clone
                        .push_insert(format!("key{}", i), vec, Some(i as i64))
                        .await
                } else {
                    q_clone
                        .push_delete(format!("key{}", i), Some(i as i64))
                        .await
                }
            });
        }

        while let Some(res) = tasks.join_next().await {
            res.unwrap().unwrap();
        }

        assert_eq!(queue.ivq_len(), num_tasks / 2);
        assert_eq!(queue.dvq_len(), num_tasks / 2);

        let mut stream = queue.drain_queues(num_tasks as i64, 10);
        let mut items = Vec::new();
        while let Some(item) = stream.next().await {
            items.push(item.unwrap());
        }

        assert_eq!(items.len(), num_tasks as usize);
        assert_eq!(queue.ivq_len(), 0);
        assert_eq!(queue.dvq_len(), 0);
    }

    #[tokio::test]
    async fn test_invalid_uuid() {
        let (q, _guard) = setup("invalid_uuid").await;
        let res = q.push_insert("", vec![1.0], None).await;
        assert!(matches!(res, Err(QueueError::InvalidUuid)));
        let res = q.push_delete(" ", None).await;
        assert!(matches!(res, Err(QueueError::InvalidUuid)));
    }

    #[tokio::test]
    async fn test_parse_key_invalid() {
        // Key shorter than 8 bytes should error
        let short_key = vec![0, 1, 2, 3, 4, 5, 6];
        assert!(matches!(
            PersistentQueue::parse_key(&short_key),
            Err(QueueError::KeyParse(_))
        ));
    }

    #[tokio::test]
    async fn test_decode_error_on_corrupt_data() {
        let (q, _guard) = setup("corrupt_data").await;
        // Directly insert a bad bincode payload
        let key = PersistentQueue::create_key(100, "key_corrupt");
        let _ = q.insert_queue.insert(&key, &[0xff, 0xff, 0xff]); // invalid payload
        let mut stream = q.drain_queues(200, 10);
        // The decode error should surface
        let err = stream.next().await.unwrap().unwrap_err();
        assert!(matches!(err, QueueError::Decode(_)));
    }

    #[tokio::test]
    async fn test_replace_insert() {
        let (q, _guard) = setup("replace_insert").await;
        q.push_insert("key1", vec![1.0], Some(100)).await.unwrap();
        assert_eq!(q.ivq_len(), 1);
        q.push_insert("key1", vec![2.0], Some(200)).await.unwrap();
        assert_eq!(q.ivq_len(), 1);

        let mut stream = q.drain_queues(300, 10);
        let item = stream.next().await.unwrap().unwrap();
        assert_eq!(item, DrainItem::Insert("key1".into(), vec![2.0]));
        assert!(stream.next().await.is_none());
        assert_eq!(q.ivq_len(), 0);
    }

    #[tokio::test]
    async fn test_partial_drain() {
        let (q, _guard) = setup("partial_drain").await;
        q.push_insert("key1", vec![1.0], Some(100)).await.unwrap();
        q.push_insert("key2", vec![2.0], Some(300)).await.unwrap();

        let mut stream = q.drain_queues(200, 10);
        let item = stream.next().await.unwrap().unwrap();
        assert_eq!(item, DrainItem::Insert("key1".into(), vec![1.0]));
        assert!(stream.next().await.is_none());
        assert_eq!(q.ivq_len(), 1);
    }

    #[tokio::test]
    async fn test_persistence() {
        let path = "./test_db_queue_persistence";
        let _ = fs::remove_dir_all(path);
        let _guard = TestGuard { path: path.into() };

        let q1 = Builder::new(path).build().await.unwrap();
        q1.push_insert("key1", vec![1.0], Some(100)).await.unwrap();
        q1.push_delete("key2", Some(100)).await.unwrap();
        drop(q1);

        let q2 = Builder::new(path).build().await.unwrap();
        assert_eq!(q2.ivq_len(), 1);
        assert_eq!(q2.dvq_len(), 1);

        let mut stream = q2.drain_queues(200, 10);
        let mut items = Vec::new();
        while let Some(item) = stream.next().await {
            items.push(item.unwrap());
        }
        assert_eq!(items.len(), 2);
        assert!(items.contains(&DrainItem::Insert("key1".into(), vec![1.0])));
        assert!(items.contains(&DrainItem::Delete("key2".into())));
    }

    #[tokio::test(flavor = "multi_thread", worker_threads = 10)]
    async fn test_transactional_complex_ordering() {
        let (q, _guard) = setup("transactional_complex_ordering").await;

        // 1. Insert key4 at t=90
        q.push_insert("key4", vec![4.0], Some(90)).await.unwrap();
        // 2. Delete key5 at t=95
        q.push_delete("key5", Some(95)).await.unwrap();
        // 3. Insert key1 at t=100
        q.push_insert("key1", vec![1.0], Some(100)).await.unwrap();
        // 4. Insert key2 at t=110
        q.push_insert("key2", vec![2.0], Some(110)).await.unwrap();
        // 5. Delete key1 at t=120
        q.push_delete("key1", Some(120)).await.unwrap();
        // 6. Delete key3 at t=125 (delete before insert)
        q.push_delete("key3", Some(125)).await.unwrap();
        // 7. Insert key3 at t=130
        q.push_insert("key3", vec![3.0], Some(130)).await.unwrap();
        // 8. Insert key1 again at t=140 (re-insert, should override delete)
        q.push_insert("key1", vec![1.4], Some(140)).await.unwrap();
        // 9. Delete key2 at t=150 (should override insert)
        q.push_delete("key2", Some(150)).await.unwrap();

        // After all operations:
        // ivq should have 4 items: key1, key2, key3, key4.
        // dvq should have 4 items: key1, key2, key3, key5.
        assert_eq!(q.ivq_len(), 4);
        assert_eq!(q.dvq_len(), 4);

        let mut stream = q.drain_queues(200, 10);
        let mut items = Vec::new();
        while let Some(item) = stream.next().await {
            items.push(item.unwrap());
        }

        // Expected final state after draining and conflict resolution:
        // - key4: inserted at 90. Kept.
        // - key5: deleted at 95. Kept.
        // - key1: inserted at 100, deleted at 120, inserted at 140. Final: Insert(140).
        // - key2: inserted at 110, deleted at 150. Final: Delete(150).
        // - key3: deleted at 125, inserted at 130. Final: Insert(130).
        let expected_items = vec![
            DrainItem::Insert("key4".into(), vec![4.0]), // ts=90
            DrainItem::Delete("key5".into()),            // ts=95
            DrainItem::Insert("key3".into(), vec![3.0]), // ts=130
            DrainItem::Insert("key1".into(), vec![1.4]), // ts=140
            DrainItem::Delete("key2".into()),            // ts=150
        ];

        assert_eq!(items.len(), 5);
        assert_eq!(items, expected_items);

        assert_eq!(q.ivq_len(), 0);
        assert_eq!(q.dvq_len(), 0);
    }

    #[tokio::test]
    async fn test_pop_insert_basic() {
        let (q, _guard) = setup("pop_insert_basic").await;
        let vec = vec![1.0, 2.0, 3.0];
        q.push_insert("key1", vec.clone(), Some(100)).await.unwrap();
        assert_eq!(q.ivq_len(), 1);

        let (popped_vec, ts) = q.pop_insert("key1").await.unwrap();
        assert_eq!(popped_vec, vec);
        assert_eq!(ts, 100);
        assert_eq!(q.ivq_len(), 0);

        // Trying to pop again should return NotFound
        let res = q.pop_insert("key1").await;
        assert!(matches!(res, Err(QueueError::NotFound(_))));
    }

    #[tokio::test]
    async fn test_pop_delete_basic() {
        let (q, _guard) = setup("pop_delete_basic").await;
        q.push_delete("key1", Some(200)).await.unwrap();
        assert_eq!(q.dvq_len(), 1);

        let ts = q.pop_delete("key1").await.unwrap();
        assert_eq!(ts, 200);
        assert_eq!(q.dvq_len(), 0);

        // Trying to pop again should return NotFound
        let res = q.pop_delete("key1").await;
        assert!(matches!(res, Err(QueueError::NotFound(_))));
    }

    #[tokio::test]
    async fn test_pop_insert_not_found() {
        let (q, _guard) = setup("pop_insert_not_found").await;
        let res = q.pop_insert("nonexistent").await;
        assert!(matches!(res, Err(QueueError::NotFound(_))));
    }

    #[tokio::test]
    async fn test_pop_delete_not_found() {
        let (q, _guard) = setup("pop_delete_not_found").await;
        let res = q.pop_delete("nonexistent").await;
        assert!(matches!(res, Err(QueueError::NotFound(_))));
    }

    #[tokio::test]
    async fn test_pop_insert_invalid_uuid() {
        let (q, _guard) = setup("pop_insert_invalid_uuid").await;
        let res = q.pop_insert("").await;
        assert!(matches!(res, Err(QueueError::InvalidUuid)));
        let res = q.pop_insert("   ").await;
        assert!(matches!(res, Err(QueueError::InvalidUuid)));
    }

    #[tokio::test]
    async fn test_pop_delete_invalid_uuid() {
        let (q, _guard) = setup("pop_delete_invalid_uuid").await;
        let res = q.pop_delete("").await;
        assert!(matches!(res, Err(QueueError::InvalidUuid)));
        let res = q.pop_delete("   ").await;
        assert!(matches!(res, Err(QueueError::InvalidUuid)));
    }

    #[tokio::test]
    async fn test_pop_insert_after_update() {
        let (q, _guard) = setup("pop_insert_after_update").await;
        // Push initial vector
        q.push_insert("key1", vec![1.0], Some(100)).await.unwrap();
        // Update with new vector
        q.push_insert("key1", vec![2.0], Some(200)).await.unwrap();
        assert_eq!(q.ivq_len(), 1);

        // Pop should return the latest vector
        let (vec, ts) = q.pop_insert("key1").await.unwrap();
        assert_eq!(vec, vec![2.0]);
        assert_eq!(ts, 200);
        assert_eq!(q.ivq_len(), 0);
    }

    #[tokio::test]
    async fn test_iv_exists() {
        let (q, _guard) = setup("iv_exists").await;
        // Should not exist initially
        let res = q.iv_exists("key1").await;
        assert!(matches!(res, Err(QueueError::NotFound(_))));

        // After push, should exist
        q.push_insert("key1", vec![1.0], Some(100)).await.unwrap();
        let ts = q.iv_exists("key1").await.unwrap();
        assert_eq!(ts, 100);

        // After pop, should not exist
        q.pop_insert("key1").await.unwrap();
        let res = q.iv_exists("key1").await;
        assert!(matches!(res, Err(QueueError::NotFound(_))));
    }

    #[tokio::test]
    async fn test_dv_exists() {
        let (q, _guard) = setup("dv_exists").await;
        // Should not exist initially
        let res = q.dv_exists("key1").await;
        assert!(matches!(res, Err(QueueError::NotFound(_))));

        // After push, should exist
        q.push_delete("key1", Some(200)).await.unwrap();
        let ts = q.dv_exists("key1").await.unwrap();
        assert_eq!(ts, 200);

        // After pop, should not exist
        q.pop_delete("key1").await.unwrap();
        let res = q.dv_exists("key1").await;
        assert!(matches!(res, Err(QueueError::NotFound(_))));
    }

    #[tokio::test(flavor = "multi_thread", worker_threads = 4)]
    async fn test_pop_insert_concurrent() {
        let (q, _guard) = setup("pop_insert_concurrent").await;
        let queue = Arc::new(q);
        let num_items = 50;

        // Push multiple items
        for i in 0..num_items {
            queue
                .push_insert(format!("key{}", i), vec![i as f32], Some(i as i64))
                .await
                .unwrap();
        }
        assert_eq!(queue.ivq_len(), num_items);

        // Pop all items concurrently
        let mut tasks = JoinSet::new();
        for i in 0..num_items {
            let q_clone = queue.clone();
            tasks.spawn(async move { q_clone.pop_insert(format!("key{}", i)).await });
        }

        let mut success_count = 0;
        while let Some(res) = tasks.join_next().await {
            if res.unwrap().is_ok() {
                success_count += 1;
            }
        }

        assert_eq!(success_count, num_items as usize);
        assert_eq!(queue.ivq_len(), 0);
    }

    #[tokio::test(flavor = "multi_thread", worker_threads = 4)]
    async fn test_pop_delete_concurrent() {
        let (q, _guard) = setup("pop_delete_concurrent").await;
        let queue = Arc::new(q);
        let num_items = 50;

        // Push multiple delete items
        for i in 0..num_items {
            queue
                .push_delete(format!("key{}", i), Some(i as i64))
                .await
                .unwrap();
        }
        assert_eq!(queue.dvq_len(), num_items);

        // Pop all items concurrently
        let mut tasks = JoinSet::new();
        for i in 0..num_items {
            let q_clone = queue.clone();
            tasks.spawn(async move { q_clone.pop_delete(format!("key{}", i)).await });
        }

        let mut success_count = 0;
        while let Some(res) = tasks.join_next().await {
            if res.unwrap().is_ok() {
                success_count += 1;
            }
        }

        assert_eq!(success_count, num_items as usize);
        assert_eq!(queue.dvq_len(), 0);
    }

    #[tokio::test]
    async fn test_pop_insert_multiple_vectors() {
        let (q, _guard) = setup("pop_insert_multiple_vectors").await;

        q.push_insert("key1", vec![1.0, 1.1], Some(100))
            .await
            .unwrap();
        q.push_insert("key2", vec![2.0, 2.1, 2.2], Some(200))
            .await
            .unwrap();
        q.push_insert("key3", vec![3.0], Some(300)).await.unwrap();
        assert_eq!(q.ivq_len(), 3);

        let (vec2, ts2) = q.pop_insert("key2").await.unwrap();
        assert_eq!(vec2, vec![2.0, 2.1, 2.2]);
        assert_eq!(ts2, 200);
        assert_eq!(q.ivq_len(), 2);

        let (vec1, ts1) = q.pop_insert("key1").await.unwrap();
        assert_eq!(vec1, vec![1.0, 1.1]);
        assert_eq!(ts1, 100);
        assert_eq!(q.ivq_len(), 1);

        let (vec3, ts3) = q.pop_insert("key3").await.unwrap();
        assert_eq!(vec3, vec![3.0]);
        assert_eq!(ts3, 300);
        assert_eq!(q.ivq_len(), 0);
    }

    #[tokio::test]
    async fn test_get_vector_basic() {
        let (q, _guard) = setup("get_vector_basic").await;
        let vec = vec![1.0, 2.0, 3.0];
        q.push_insert("key1", vec.clone(), Some(100)).await.unwrap();

        // get_vector should return the vector without removing it
        let (got_vec, ts) = q.get_vector("key1").await.unwrap();
        assert_eq!(got_vec, vec);
        assert_eq!(ts, 100);

        // Queue length should remain unchanged
        assert_eq!(q.ivq_len(), 1);

        // Should still be able to pop
        let (popped_vec, _) = q.pop_insert("key1").await.unwrap();
        assert_eq!(popped_vec, vec);
        assert_eq!(q.ivq_len(), 0);
    }

    #[tokio::test]
    async fn test_get_vector_not_found() {
        let (q, _guard) = setup("get_vector_not_found").await;
        let res = q.get_vector("nonexistent").await;
        assert!(matches!(res, Err(QueueError::NotFound(_))));
    }

    #[tokio::test]
    async fn test_get_vector_with_delete_newer() {
        let (q, _guard) = setup("get_vector_with_delete_newer").await;
        // Insert at t=100, delete at t=200
        q.push_insert("key1", vec![1.0], Some(100)).await.unwrap();
        q.push_delete("key1", Some(200)).await.unwrap();

        // get_vector should return NotFound because delete is newer
        let res = q.get_vector("key1").await;
        assert!(matches!(res, Err(QueueError::NotFound(_))));
    }

    #[tokio::test]
    async fn test_get_vector_with_insert_newer() {
        let (q, _guard) = setup("get_vector_with_insert_newer").await;
        // Delete at t=100, insert at t=200
        q.push_delete("key1", Some(100)).await.unwrap();
        q.push_insert("key1", vec![1.0], Some(200)).await.unwrap();

        // get_vector should return the vector because insert is newer
        let (vec, ts) = q.get_vector("key1").await.unwrap();
        assert_eq!(vec, vec![1.0]);
        assert_eq!(ts, 200);
    }

    #[tokio::test]
    async fn test_get_vector_with_timestamp_basic() {
        let (q, _guard) = setup("get_vector_with_timestamp_basic").await;
        let vec = vec![1.0, 2.0];
        q.push_insert("key1", vec.clone(), Some(100)).await.unwrap();

        let (got_vec, its, dts, exists) = q.get_vector_with_timestamp("key1").await.unwrap();
        assert_eq!(got_vec, Some(vec));
        assert_eq!(its, 100);
        assert_eq!(dts, 0);
        assert!(exists);
    }

    #[tokio::test]
    async fn test_get_vector_with_timestamp_not_found() {
        let (q, _guard) = setup("get_vector_with_timestamp_not_found").await;
        let (vec, its, dts, exists) = q.get_vector_with_timestamp("nonexistent").await.unwrap();
        assert!(vec.is_none());
        assert_eq!(its, 0);
        assert_eq!(dts, 0);
        assert!(!exists);
    }

    #[tokio::test]
    async fn test_get_vector_with_timestamp_delete_only() {
        let (q, _guard) = setup("get_vector_with_timestamp_delete_only").await;
        q.push_delete("key1", Some(100)).await.unwrap();

        let (vec, its, dts, exists) = q.get_vector_with_timestamp("key1").await.unwrap();
        assert!(vec.is_none());
        assert_eq!(its, 0);
        assert_eq!(dts, 100);
        assert!(!exists);
    }

    #[tokio::test]
    async fn test_get_vector_with_timestamp_both_queues_insert_newer() {
        let (q, _guard) = setup("get_vector_with_timestamp_both_insert_newer").await;
        q.push_delete("key1", Some(100)).await.unwrap();
        q.push_insert("key1", vec![1.0], Some(200)).await.unwrap();

        let (vec, its, dts, exists) = q.get_vector_with_timestamp("key1").await.unwrap();
        assert_eq!(vec, Some(vec![1.0]));
        assert_eq!(its, 200);
        assert_eq!(dts, 100);
        assert!(exists); // insert is newer, so exists is true
    }

    #[tokio::test]
    async fn test_get_vector_with_timestamp_both_queues_delete_newer() {
        let (q, _guard) = setup("get_vector_with_timestamp_both_delete_newer").await;
        q.push_insert("key1", vec![1.0], Some(100)).await.unwrap();
        q.push_delete("key1", Some(200)).await.unwrap();

        let (vec, its, dts, exists) = q.get_vector_with_timestamp("key1").await.unwrap();
        assert_eq!(vec, Some(vec![1.0]));
        assert_eq!(its, 100);
        assert_eq!(dts, 200);
        assert!(!exists); // delete is newer, so exists is false
    }

    #[tokio::test]
    async fn test_get_vector_with_timestamp_same_timestamp() {
        let (q, _guard) = setup("get_vector_with_timestamp_same_ts").await;
        // Same timestamp for insert and delete (like update operation)
        q.push_insert("key1", vec![1.0], Some(100)).await.unwrap();
        q.push_delete("key1", Some(100)).await.unwrap();

        let (vec, its, dts, exists) = q.get_vector_with_timestamp("key1").await.unwrap();
        assert_eq!(vec, Some(vec![1.0]));
        assert_eq!(its, 100);
        assert_eq!(dts, 100);
        assert!(!exists); // same timestamp means not newer, so exists is false
    }

    #[tokio::test]
    async fn test_get_vector_does_not_modify_queue() {
        let (q, _guard) = setup("get_vector_no_modify").await;
        q.push_insert("key1", vec![1.0], Some(100)).await.unwrap();
        q.push_insert("key2", vec![2.0], Some(200)).await.unwrap();
        assert_eq!(q.ivq_len(), 2);

        // Multiple get_vector calls should not modify the queue
        for _ in 0..5 {
            let _ = q.get_vector("key1").await.unwrap();
            let _ = q.get_vector("key2").await.unwrap();
        }

        assert_eq!(q.ivq_len(), 2);

        // get_vector_with_timestamp should also not modify
        let _ = q.get_vector_with_timestamp("key1").await.unwrap();
        let _ = q.get_vector_with_timestamp("key2").await.unwrap();

        assert_eq!(q.ivq_len(), 2);
    }

    // ========== Range Tests ==========

    #[tokio::test]
    async fn test_range_empty_queue() {
        let (q, _guard) = setup("range_empty_queue").await;

        let stream = q.range();
        let items: Vec<_> = tokio_stream::StreamExt::collect(stream).await;
        assert!(items.is_empty());
    }

    #[tokio::test]
    async fn test_range_single_item() {
        let (q, _guard) = setup("range_single_item").await;

        q.push_insert("key1", vec![1.0, 2.0], Some(100))
            .await
            .unwrap();

        let stream = q.range();
        let items: Vec<_> = tokio_stream::StreamExt::collect(stream).await;

        assert_eq!(items.len(), 1);
        let (uuid, vec, ts) = items[0].as_ref().unwrap();
        assert_eq!(uuid, "key1");
        assert_eq!(vec, &vec![1.0, 2.0]);
        assert_eq!(*ts, 100);
    }

    #[tokio::test]
    async fn test_range_multiple_items() {
        let (q, _guard) = setup("range_multiple_items").await;

        q.push_insert("key1", vec![1.0], Some(100)).await.unwrap();
        q.push_insert("key2", vec![2.0], Some(200)).await.unwrap();
        q.push_insert("key3", vec![3.0], Some(300)).await.unwrap();

        let stream = q.range();
        let items: Vec<_> = tokio_stream::StreamExt::collect(stream).await;

        assert_eq!(items.len(), 3);

        // Collect all uuids
        let uuids: Vec<_> = items
            .iter()
            .filter_map(|r| r.as_ref().ok())
            .map(|(uuid, _, _)| uuid.clone())
            .collect();

        assert!(uuids.contains(&"key1".to_string()));
        assert!(uuids.contains(&"key2".to_string()));
        assert!(uuids.contains(&"key3".to_string()));
    }

    #[tokio::test]
    async fn test_range_filters_newer_delete() {
        let (q, _guard) = setup("range_filters_newer_delete").await;

        // Insert at t=100, delete at t=200 (delete is newer, should be filtered)
        q.push_insert("key1", vec![1.0], Some(100)).await.unwrap();
        q.push_delete("key1", Some(200)).await.unwrap();

        // Insert at t=300, delete at t=100 (insert is newer, should appear)
        q.push_insert("key2", vec![2.0], Some(300)).await.unwrap();
        q.push_delete("key2", Some(100)).await.unwrap();

        let stream = q.range();
        let items: Vec<_> = tokio_stream::StreamExt::collect(stream).await;

        // Only key2 should appear because key1 has a newer delete
        assert_eq!(items.len(), 1);
        let (uuid, vec, ts) = items[0].as_ref().unwrap();
        assert_eq!(uuid, "key2");
        assert_eq!(vec, &vec![2.0]);
        assert_eq!(*ts, 300);
    }

    #[tokio::test]
    async fn test_range_same_timestamp_filtered() {
        let (q, _guard) = setup("range_same_timestamp_filtered").await;

        // Insert and delete at same timestamp (delete >= insert, should be filtered)
        q.push_insert("key1", vec![1.0], Some(100)).await.unwrap();
        q.push_delete("key1", Some(100)).await.unwrap();

        let stream = q.range();
        let items: Vec<_> = tokio_stream::StreamExt::collect(stream).await;

        assert!(items.is_empty());
    }

    #[tokio::test]
    async fn test_range_does_not_modify_queue() {
        let (q, _guard) = setup("range_no_modify").await;

        q.push_insert("key1", vec![1.0], Some(100)).await.unwrap();
        q.push_insert("key2", vec![2.0], Some(200)).await.unwrap();

        assert_eq!(q.ivq_len(), 2);

        // Multiple range calls should not modify the queue
        for _ in 0..3 {
            let stream = q.range();
            let _: Vec<_> = tokio_stream::StreamExt::collect(stream).await;
        }

        assert_eq!(q.ivq_len(), 2);
    }

    #[tokio::test]
    async fn test_range_no_delete() {
        let (q, _guard) = setup("range_no_delete").await;

        // Items without any delete should all appear
        q.push_insert("key1", vec![1.0], Some(100)).await.unwrap();
        q.push_insert("key2", vec![2.0], Some(200)).await.unwrap();

        let stream = q.range();
        let items: Vec<_> = tokio_stream::StreamExt::collect(stream).await;

        assert_eq!(items.len(), 2);
    }
}
