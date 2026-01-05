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
use bincode::config::{self, Configuration};
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

/// A performance-oriented bincode configuration.
const BINCODE_CONFIG: Configuration = config::standard()
    .with_little_endian()
    .with_variable_int_encoding();

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
    Encode(#[from] bincode::error::EncodeError),
    /// Error returned for deserialization failures.
    #[error("Codec decode error: {0}")]
    Decode(#[from] bincode::error::DecodeError),
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
                let (vec, _) = bincode::decode_from_slice(&val, BINCODE_CONFIG)?;
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
        let value = bincode::encode_to_vec(&vector, BINCODE_CONFIG)?;
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
}
