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

//! # Persistent Vector Queue
//!
//! A persistent, async-safe, and concurrent vector queue implementation.
//! It is designed for high-concurrency environments, offloading blocking I/O operations
//! to a dedicated thread pool managed by `tokio`.
//!
//! The implementation uses `sled` as its underlying persistent storage engine to leverage
//! its robust transactional capabilities, ensuring data consistency.

use async_trait::async_trait;
use bincode::{Decode, Encode, config::standard as bincode_standard};
use serde::{Deserialize, Deserializer, Serialize, Serializer};
use sled::{
    Tree,
    transaction::{ConflictableTransactionError, TransactionError, Transactional},
};
use std::pin::Pin;
use std::sync::Arc;
use std::sync::atomic::{AtomicU64, Ordering as AtomicOrdering};
use std::time::{SystemTime, UNIX_EPOCH};
use thiserror::Error;
use tokio::sync::mpsc;
use tokio_stream::{Stream, wrappers::ReceiverStream};

/// Error type representing possible failures for queue operations.
#[derive(Debug, Error)]
pub enum QueueError {
    #[error("The provided UUID is invalid or empty.")]
    InvalidUuid,
    #[error("The provided timestamp is older than an existing entry and cannot replace it.")]
    TimestampTooOld,
    #[error("Sled database error")]
    Sled(#[from] sled::Error),
    #[error("Sled transaction error")]
    SledTransaction {
        #[source]
        source: Box<dyn std::error::Error + Send + Sync>,
    },
    #[error("Codec error")]
    Codec {
        #[source]
        source: Box<dyn std::error::Error + Send + Sync>,
    },
    #[error("Internal Tokio task error")]
    Internal(#[from] tokio::task::JoinError),
    #[error("UTF-8 conversion error")]
    Utf8(#[from] std::string::FromUtf8Error),
}

/// Trait definition for queue functionality.
#[async_trait]
pub trait Queue: Send + Sync {
    async fn push_insert(
        &self,
        uuid: String,
        vector: Arc<Vec<f32>>,
        timestamp: Option<i64>,
    ) -> Result<(), QueueError>;

    async fn push_delete(&self, uuid: String, timestamp: Option<i64>) -> Result<(), QueueError>;

    async fn pop_insert(&self, uuid: &str) -> Result<Option<(Arc<Vec<f32>>, i64)>, QueueError>;

    async fn pop_delete(&self, uuid: &str) -> Result<Option<i64>, QueueError>;

    async fn get_vector(&self, uuid: &str) -> Result<Option<(Arc<Vec<f32>>, i64)>, QueueError>;

    async fn get_vector_with_timestamp(
        &self,
        uuid: &str,
    ) -> Result<Option<(Arc<Vec<f32>>, i64, i64)>, QueueError>;

    fn range(
        &self,
    ) -> Pin<Box<dyn Stream<Item = Result<(String, Arc<Vec<f32>>, i64), QueueError>> + Send>>;

    fn ivq_len(&self) -> u64;

    fn dvq_len(&self) -> u64;

    fn range_pop_insert(
        &self,
        now: i64,
    ) -> Pin<Box<dyn Stream<Item = Result<(String, Arc<Vec<f32>>, i64), QueueError>> + Send>>;

    fn range_pop_delete(
        &self,
        now: i64,
    ) -> Pin<Box<dyn Stream<Item = Result<(String, i64), QueueError>> + Send>>;

    async fn iv_exists(&self, uuid: &str) -> Result<Option<i64>, QueueError>;

    async fn dv_exists(&self, uuid: &str) -> Result<Option<i64>, QueueError>;
}

/// Internal structure holding index information.
#[derive(Serialize, Deserialize, Encode, Decode, Clone, Debug)]
struct Index {
    #[serde(with = "arc_vec_serde")]
    vector: Option<Arc<Vec<f32>>>,
    timestamp: i64,
}

mod arc_vec_serde {
    use super::*;
    pub fn serialize<S>(vec: &Option<Arc<Vec<f32>>>, serializer: S) -> Result<S::Ok, S::Error>
    where
        S: Serializer,
    {
        let vec_ref: Option<&Vec<f32>> = vec.as_ref().map(|arc| &**arc);
        vec_ref.serialize(serializer)
    }

    pub fn deserialize<'de, D>(deserializer: D) -> Result<Option<Arc<Vec<f32>>>, D::Error>
    where
        D: Deserializer<'de>,
    {
        let vec: Option<Vec<f32>> = Option::deserialize(deserializer)?;
        Ok(vec.map(Arc::new))
    }
}

/// A persistent queue implementation using `sled`.
#[derive(Clone)]
pub struct PersistentQueue {
    insert_queue: Tree,
    delete_queue: Tree,
    insert_count: Arc<AtomicU64>,
    delete_count: Arc<AtomicU64>,
}

pub struct Builder {
    path: String,
}

impl Builder {
    pub fn new(path: impl AsRef<str>) -> Self {
        Self {
            path: path.as_ref().to_string(),
        }
    }

    pub async fn build(self) -> Result<PersistentQueue, QueueError> {
        let path = self.path;
        let db = tokio::task::spawn_blocking(move || sled::open(path)).await??;
        let insert_queue = db.open_tree("insert_queue")?;
        let delete_queue = db.open_tree("delete_queue")?;

        let insert_count = Arc::new(AtomicU64::new(insert_queue.len() as u64));
        let delete_count = Arc::new(AtomicU64::new(delete_queue.len() as u64));

        Ok(PersistentQueue {
            insert_queue,
            delete_queue,
            insert_count,
            delete_count,
        })
    }
}

impl PersistentQueue {
    fn now_ns() -> i64 {
        SystemTime::now()
            .duration_since(UNIX_EPOCH)
            .unwrap_or_default()
            .as_nanos() as i64
    }

    fn decode_index(bytes: &[u8]) -> Result<Index, QueueError> {
        bincode::decode_from_slice(bytes, bincode_standard())
            .map(|(i, _)| i)
            .map_err(|e| QueueError::Codec {
                source: Box::new(e),
            })
    }

    fn decode_index_in_transaction(
        bytes: &[u8],
    ) -> Result<Index, ConflictableTransactionError<QueueError>> {
        bincode::decode_from_slice(bytes, bincode_standard())
            .map(|(i, _)| i)
            .map_err(|e| {
                ConflictableTransactionError::Abort(QueueError::Codec {
                    source: Box::new(e),
                })
            })
    }
}

#[async_trait]
impl Queue for PersistentQueue {
    async fn push_insert(
        &self,
        uuid: String,
        vector: Arc<Vec<f32>>,
        timestamp: Option<i64>,
    ) -> Result<(), QueueError> {
        if uuid.trim().is_empty() {
            return Err(QueueError::InvalidUuid);
        }
        let ts = timestamp.unwrap_or_else(Self::now_ns);

        let iq = self.insert_queue.clone();
        let dq = self.delete_queue.clone();
        let ic = self.insert_count.clone();

        tokio::task::spawn_blocking(move || {
            let new_idx = Index {
                vector: Some(vector),
                timestamp: ts,
            };

            (&iq, &dq)
                .transaction(move |(iq_tx, dq_tx)| {
                    if let Some(d_bytes) = dq_tx.get(uuid.as_bytes())? {
                        let d_idx = Self::decode_index_in_transaction(&d_bytes)?;
                        if d_idx.timestamp > ts {
                            return Err(ConflictableTransactionError::Abort(
                                QueueError::TimestampTooOld,
                            ));
                        }
                    }

                    let is_new_entry;
                    if let Some(i_bytes) = iq_tx.get(uuid.as_bytes())? {
                        is_new_entry = false;
                        let i_idx = Self::decode_index_in_transaction(&i_bytes)?;
                        if i_idx.timestamp >= ts {
                            return Err(ConflictableTransactionError::Abort(
                                QueueError::TimestampTooOld,
                            ));
                        }
                    } else {
                        is_new_entry = true;
                    }

                    let new_idx_bytes = bincode::encode_to_vec(&new_idx, bincode_standard())
                        .map_err(|e| {
                            ConflictableTransactionError::Abort(QueueError::Codec {
                                source: Box::new(e),
                            })
                        })?;
                    iq_tx.insert(uuid.as_bytes(), new_idx_bytes)?;

                    if is_new_entry {
                        ic.fetch_add(1, AtomicOrdering::SeqCst);
                    }
                    Ok(())
                })
                .map_err(
                    |e: TransactionError<QueueError>| QueueError::SledTransaction {
                        source: Box::new(e),
                    },
                )
        })
        .await?
    }

    async fn push_delete(&self, uuid: String, timestamp: Option<i64>) -> Result<(), QueueError> {
        if uuid.trim().is_empty() {
            return Err(QueueError::InvalidUuid);
        }
        let ts = timestamp.unwrap_or_else(Self::now_ns);
        let dq = self.delete_queue.clone();
        let dc = self.delete_count.clone();

        tokio::task::spawn_blocking(move || {
            let new_idx = Index {
                vector: None,
                timestamp: ts,
            };

            (&dq)
                .transaction(move |dq_tx| {
                    let is_new_entry = dq_tx.get(uuid.as_bytes())?.is_none();
                    let new_idx_bytes = bincode::encode_to_vec(&new_idx, bincode_standard())
                        .map_err(|e| {
                            ConflictableTransactionError::Abort(QueueError::Codec {
                                source: Box::new(e),
                            })
                        })?;
                    dq_tx.insert(uuid.as_bytes(), new_idx_bytes)?;

                    if is_new_entry {
                        dc.fetch_add(1, AtomicOrdering::SeqCst);
                    }
                    Ok(())
                })
                .map_err(
                    |e: TransactionError<QueueError>| QueueError::SledTransaction {
                        source: Box::new(e),
                    },
                )
        })
        .await?
    }

    async fn pop_insert(&self, uuid: &str) -> Result<Option<(Arc<Vec<f32>>, i64)>, QueueError> {
        let iq = self.insert_queue.clone();
        let ic = self.insert_count.clone();
        let uuid = uuid.to_string();

        tokio::task::spawn_blocking(move || {
            if let Some(bytes) = iq.remove(uuid.as_bytes())? {
                ic.fetch_sub(1, AtomicOrdering::SeqCst);
                let idx = Self::decode_index(&bytes)?;
                if let Some(vector) = idx.vector {
                    return Ok(Some((vector, idx.timestamp)));
                }
            }
            Ok(None)
        })
        .await?
    }

    async fn pop_delete(&self, uuid: &str) -> Result<Option<i64>, QueueError> {
        let dq = self.delete_queue.clone();
        let iq = self.insert_queue.clone();
        let dc = self.delete_count.clone();
        let ic = self.insert_count.clone();
        let uuid = uuid.to_string();

        tokio::task::spawn_blocking(move || {
            let transaction_res = (&dq, &iq).transaction(|(dq_tx, iq_tx)| {
                if let Some(d_bytes) = dq_tx.remove(uuid.as_bytes())? {
                    dc.fetch_sub(1, AtomicOrdering::SeqCst);
                    let d_idx = Self::decode_index_in_transaction(&d_bytes)?;
                    if let Some(i_bytes) = iq_tx.get(uuid.as_bytes())? {
                        let i_idx = Self::decode_index_in_transaction(&i_bytes)?;
                        if d_idx.timestamp > i_idx.timestamp {
                            if iq_tx.remove(uuid.as_bytes())?.is_some() {
                                ic.fetch_sub(1, AtomicOrdering::SeqCst);
                            }
                        }
                    }
                    return Ok(Some(d_idx.timestamp));
                }
                Ok(None)
            });
            transaction_res.map_err(
                |e: TransactionError<QueueError>| QueueError::SledTransaction {
                    source: Box::new(e),
                },
            )
        })
        .await?
    }

    async fn get_vector(&self, uuid: &str) -> Result<Option<(Arc<Vec<f32>>, i64)>, QueueError> {
        if let Some((vec, ts, _)) = self.get_vector_with_timestamp(uuid).await? {
            return Ok(Some((vec, ts)));
        }
        Ok(None)
    }

    async fn get_vector_with_timestamp(
        &self,
        uuid: &str,
    ) -> Result<Option<(Arc<Vec<f32>>, i64, i64)>, QueueError> {
        let iq = self.insert_queue.clone();
        let (its_opt, dts_opt) = self.get_timestamps(uuid).await?;

        if let Some(its) = its_opt {
            if its > dts_opt.unwrap_or(0) {
                let iq_clone = iq.clone();
                let uuid = uuid.to_string();
                return tokio::task::spawn_blocking(move || {
                    if let Some(bytes) = iq_clone.get(uuid)? {
                        let idx = Self::decode_index(&bytes)?;
                        if let Some(vector) = idx.vector {
                            return Ok(Some((vector, its, dts_opt.unwrap_or(0))));
                        }
                    }
                    Ok(None)
                })
                .await?;
            }
        }
        Ok(None)
    }

    fn range(
        &self,
    ) -> Pin<Box<dyn Stream<Item = Result<(String, Arc<Vec<f32>>, i64), QueueError>> + Send>> {
        let (tx, rx) = mpsc::channel(128);
        let q = self.clone();

        tokio::spawn(async move {
            let iter = q.insert_queue.iter();
            for item in iter {
                let res = (|| -> Result<Option<(String, Arc<Vec<f32>>, i64)>, QueueError> {
                    let (uuid_bytes, i_bytes) = item?;
                    let uuid = String::from_utf8(uuid_bytes.to_vec())?;
                    let i_idx = Self::decode_index(&i_bytes)?;

                    if let Some(d_bytes) = q.delete_queue.get(&uuid_bytes)? {
                        let d_idx = Self::decode_index(&d_bytes)?;
                        if d_idx.timestamp > i_idx.timestamp {
                            return Ok(None);
                        }
                    }

                    if let Some(vector) = i_idx.vector {
                        Ok(Some((uuid, vector, i_idx.timestamp)))
                    } else {
                        Ok(None)
                    }
                })();

                match res {
                    Ok(Some(item)) => {
                        if tx.send(Ok(item)).await.is_err() {
                            break;
                        }
                    }
                    Err(e) => {
                        if tx.send(Err(e)).await.is_err() {
                            break;
                        }
                    }
                    Ok(None) => {}
                }
            }
        });

        Box::pin(ReceiverStream::new(rx))
    }

    fn ivq_len(&self) -> u64 {
        self.insert_count.load(AtomicOrdering::Relaxed)
    }
    fn dvq_len(&self) -> u64 {
        self.delete_count.load(AtomicOrdering::Relaxed)
    }

    fn range_pop_insert(
        &self,
        now: i64,
    ) -> Pin<Box<dyn Stream<Item = Result<(String, Arc<Vec<f32>>, i64), QueueError>> + Send>> {
        let q = self.clone();
        let (tx, rx) = mpsc::channel(128);

        tokio::spawn(async move {
            let keys = match q.collect_insert_keys_to_pop(now).await {
                Ok(keys) => keys,
                Err(e) => {
                    let _ = tx.send(Err(e)).await;
                    return;
                }
            };

            for uuid in keys {
                match q.pop_insert(&uuid).await {
                    Ok(Some(item)) => {
                        if tx.send(Ok((uuid, item.0, item.1))).await.is_err() {
                            break;
                        }
                    }
                    Ok(None) => { /* Already popped */ }
                    Err(e) => {
                        if tx.send(Err(e)).await.is_err() {
                            break;
                        }
                    }
                }
            }
        });

        Box::pin(ReceiverStream::new(rx))
    }

    fn range_pop_delete(
        &self,
        now: i64,
    ) -> Pin<Box<dyn Stream<Item = Result<(String, i64), QueueError>> + Send>> {
        let q = self.clone();
        let (tx, rx) = mpsc::channel(128);

        tokio::spawn(async move {
            let keys = match q.collect_delete_keys_to_pop(now).await {
                Ok(keys) => keys,
                Err(e) => {
                    let _ = tx.send(Err(e)).await;
                    return;
                }
            };

            for uuid in keys {
                match q.pop_delete(&uuid).await {
                    Ok(Some(ts)) => {
                        if tx.send(Ok((uuid, ts))).await.is_err() {
                            break;
                        }
                    }
                    Ok(None) => { /* Already popped */ }
                    Err(e) => {
                        if tx.send(Err(e)).await.is_err() {
                            break;
                        }
                    }
                }
            }
        });

        Box::pin(ReceiverStream::new(rx))
    }

    async fn iv_exists(&self, uuid: &str) -> Result<Option<i64>, QueueError> {
        let (its_opt, dts_opt) = self.get_timestamps(uuid).await?;
        if let Some(its) = its_opt {
            if its > dts_opt.unwrap_or(0) {
                return Ok(Some(its));
            }
        }
        Ok(None)
    }

    async fn dv_exists(&self, uuid: &str) -> Result<Option<i64>, QueueError> {
        let (its_opt, dts_opt) = self.get_timestamps(uuid).await?;
        if let Some(dts) = dts_opt {
            if dts > its_opt.unwrap_or(0) {
                return Ok(Some(dts));
            }
        }
        Ok(None)
    }
}

impl PersistentQueue {
    /// Helper to get both insert and delete timestamps atomically.
    async fn get_timestamps(&self, uuid: &str) -> Result<(Option<i64>, Option<i64>), QueueError> {
        let iq = self.insert_queue.clone();
        let dq = self.delete_queue.clone();
        let uuid = uuid.to_string();

        tokio::task::spawn_blocking(move || {
            let (i_idx_opt, d_idx_opt): (Option<Index>, Option<Index>) = (&iq, &dq)
                .transaction(|(iq_tx, dq_tx)| {
                    let i_bytes = iq_tx.get(uuid.as_bytes())?;
                    let d_bytes = dq_tx.get(uuid.as_bytes())?;
                    let i_idx_opt = i_bytes
                        .map(|b| Self::decode_index_in_transaction(&b))
                        .transpose()?;
                    let d_idx_opt = d_bytes
                        .map(|b| Self::decode_index_in_transaction(&b))
                        .transpose()?;
                    Ok((i_idx_opt, d_idx_opt))
                })
                .map_err(
                    |e: TransactionError<QueueError>| QueueError::SledTransaction {
                        source: Box::new(e),
                    },
                )?;

            Ok((
                i_idx_opt.map(|i| i.timestamp),
                d_idx_opt.map(|d| d.timestamp),
            ))
        })
        .await?
    }

    async fn collect_insert_keys_to_pop(&self, now: i64) -> Result<Vec<String>, QueueError> {
        let iq = self.insert_queue.clone();
        let dq = self.delete_queue.clone();
        let ic = self.insert_count.clone();

        tokio::task::spawn_blocking(move || {
            let mut keys = Vec::new();
            for item in iq.iter() {
                let (uuid_bytes, i_bytes) = item?;
                let i_idx = Self::decode_index(&i_bytes)?;

                if i_idx.timestamp > now {
                    continue;
                }

                if let Some(d_bytes) = dq.get(&uuid_bytes)? {
                    let d_idx = Self::decode_index(&d_bytes)?;
                    if d_idx.timestamp > i_idx.timestamp {
                        if iq.remove(&uuid_bytes)?.is_some() {
                            ic.fetch_sub(1, AtomicOrdering::SeqCst);
                        }
                        continue;
                    }
                }
                keys.push(String::from_utf8(uuid_bytes.to_vec())?);
            }
            Ok(keys)
        })
        .await?
    }

    async fn collect_delete_keys_to_pop(&self, now: i64) -> Result<Vec<String>, QueueError> {
        let dq = self.delete_queue.clone();
        tokio::task::spawn_blocking(move || {
            let mut keys = Vec::new();
            for item in dq.iter() {
                let (uuid_bytes, d_bytes) = item?;
                let d_idx = Self::decode_index(&d_bytes)?;
                if d_idx.timestamp <= now {
                    keys.push(String::from_utf8(uuid_bytes.to_vec())?);
                }
            }
            Ok(keys)
        })
        .await?
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
        let queue = Builder::new(&path).build().await.unwrap();
        (queue, guard)
    }

    #[tokio::test]
    async fn test_push_and_pop_insert() {
        let (q, _guard) = setup("push_pop_insert").await;
        let vec = Arc::new(vec![1.0, 2.0]);
        q.push_insert("key1".into(), vec.clone(), None)
            .await
            .unwrap();
        assert_eq!(q.ivq_len(), 1);
        let (vec_out, _) = q.pop_insert("key1").await.unwrap().unwrap();
        assert_eq!(*vec_out, vec![1.0, 2.0]);
        assert_eq!(q.ivq_len(), 0);
        assert!(q.pop_insert("key1").await.unwrap().is_none());
    }

    #[tokio::test]
    async fn test_pop_delete_is_fully_atomic() {
        let (q, _guard) = setup("pop_delete_atomic").await;
        q.push_insert("key1".to_string(), Arc::new(vec![1.0]), Some(100))
            .await
            .unwrap();
        q.push_delete("key1".to_string(), Some(200)).await.unwrap();

        assert_eq!(q.ivq_len(), 1);
        assert_eq!(q.dvq_len(), 1);

        let ts = q.pop_delete("key1").await.unwrap().unwrap();
        assert_eq!(ts, 200);

        assert_eq!(q.ivq_len(), 0);
        assert_eq!(q.dvq_len(), 0);
        assert!(q.get_vector("key1").await.unwrap().is_none());
    }

    #[tokio::test]
    async fn test_range_pop_insert_stream() {
        let (q, _guard) = setup("range_pop_insert_stream").await;
        q.push_insert("key1".to_string(), Arc::new(vec![1.0]), Some(100))
            .await
            .unwrap();
        q.push_insert("key2".to_string(), Arc::new(vec![2.0]), Some(200))
            .await
            .unwrap();
        q.push_insert("key3".to_string(), Arc::new(vec![3.0]), Some(300))
            .await
            .unwrap();
        q.push_delete("key2".to_string(), Some(250)).await.unwrap();

        let mut stream = q.range_pop_insert(280);
        let mut results = Vec::new();
        while let Some(item) = stream.next().await {
            results.push(item.unwrap());
        }

        assert_eq!(results.len(), 1);
        assert_eq!(results[0].0, "key1");
        assert_eq!(q.ivq_len(), 1);
        assert!(q.get_vector("key3").await.unwrap().is_some());
        assert!(q.get_vector("key1").await.unwrap().is_none());
    }

    #[tokio::test]
    async fn test_range_pop_delete_stream() {
        let (q, _guard) = setup("range_pop_delete_stream").await;
        q.push_insert("key1".to_string(), Arc::new(vec![1.0]), Some(100))
            .await
            .unwrap();
        q.push_insert("key2".to_string(), Arc::new(vec![2.0]), Some(300))
            .await
            .unwrap();
        q.push_delete("key1".to_string(), Some(200)).await.unwrap();
        q.push_delete("key2".to_string(), Some(250)).await.unwrap();
        q.push_delete("key3".to_string(), Some(400)).await.unwrap();

        let stream = q.range_pop_delete(500);
        let results: Vec<_> = stream.collect::<Vec<_>>().await;

        assert_eq!(results.len(), 3);
        assert_eq!(q.dvq_len(), 0);
        assert_eq!(q.ivq_len(), 1);
        assert!(q.get_vector("key1").await.unwrap().is_none());
        assert!(q.get_vector("key2").await.unwrap().is_some());
    }

    #[tokio::test]
    async fn test_concurrent_pushes() {
        let (q, _guard) = setup("concurrent_pushes").await;
        let queue = Arc::new(q);
        let mut tasks = JoinSet::new();
        let num_tasks = 100;

        for i in 0..num_tasks {
            let q_clone = queue.clone();
            tasks.spawn(async move {
                if i % 2 == 0 {
                    let vec = Arc::new(vec![i as f32]);
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
    }
}
