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

//! # Bidirectional Map
//!
//! A generic, persistent, async-safe, and pluggable-codec bidirectional map.
//! It is designed for high-concurrency environments, offloading blocking I/O operations
//! to a dedicated thread pool managed by `tokio`. This version focuses on core
//! functionality without Time-To-Live (TTL) features.
//!
//! The implementation uses `sled` as its underlying persistent storage engine to leverage
//! its robust transactional capabilities, ensuring data consistency for bidirectional mappings.
use bincode::{Decode, Encode, config::standard as bincode_standard};
use futures::{Stream, StreamExt};
use serde::{Serialize, de::DeserializeOwned};
use sled::{
    Db, IVec, Tree,
    transaction::{ConflictableTransactionError, TransactionError, Transactional},
};
use std::borrow::Borrow;
use std::sync::Arc;
use std::sync::atomic::{AtomicUsize, Ordering};
use tokio::sync::mpsc;
use tokio_stream::wrappers::ReceiverStream;
use tracing::instrument;

pub mod codec;
pub mod error;
pub mod types;

use crate::{
    codec::{Codec, BincodeCodec},
    error::Error,
    types::{KeyType, ValueType},
};

/// A builder for constructing a `BidiMap` instance with custom configurations.
pub struct BidiBuilder<K, V, C: Codec = BincodeCodec> {
    path: String,
    codec: C,
    scan_on_startup: bool,
    _marker: std::marker::PhantomData<(K, V)>,
}

impl<K:KeyType, V: ValueType> BidiBuilder<K, V, BincodeCodec> {
    /// Creates a new `BidiBuilder` with a specified database path and the default `BincodeCodec`.
    pub fn new(path: impl AsRef<str>) -> Self {
        Self {
            path: path.as_ref().to_string(),
            codec: BincodeCodec::default(),
            scan_on_startup: true,
            _marker: std::marker::PhantomData,
        }
    }
}

impl<K: KeyType, V: ValueType, C: Codec> BidiBuilder<K, V, C> {
    /// Sets a custom codec for the `BidiMap`, returning a new builder instance.
    pub fn codec<NewC: Codec>(self, new_codec: NewC) -> BidiBuilder<K, V, NewC> {
        BidiBuilder {
            path: self.path,
            codec: new_codec,
            scan_on_startup: self.scan_on_startup,
            _marker: std::marker::PhantomData,
        }
    }

    /// Disables the full database scan on startup.
    pub fn disable_scan_on_startup(mut self) -> Self {
        self.scan_on_startup = false;
        self
    }

    /// Builds the `BidiMap`, initializing the database.
    pub async fn build(self) -> Result<Arc<BidiInner<K, V, C>>, Error> {
        if let Some(dir) = std::path::Path::new(&self.path).parent() {
            tokio::fs::create_dir_all(dir).await?;
        }

        let path = self.path.clone();
        let db = tokio::task::spawn_blocking(move || sled::open(path)).await??;

        let uo = db.open_tree("uo")?;
        let ou = db.open_tree("ou")?;

        let initial_len = if self.scan_on_startup { uo.len() } else { 0 };

        let inner = Arc::new(BidiInner {
            db: Arc::new(db),
            uo,
            ou,
            len: AtomicUsize::new(initial_len),
            codec: Arc::new(self.codec),
            _marker: std::marker::PhantomData,
        });

        Ok(inner)
    }
}

/// The internal struct holding the state and logic of the `BidiMap`.
pub struct BidiInner<K, V, C: Codec> {
    db: Arc<Db>,
    uo: Tree,
    ou: Tree,
    len: AtomicUsize,
    codec: Arc<C>,
    _marker: std::marker::PhantomData<(K, V)>,
}

impl<K: KeyType, V: ValueType, C: Codec> BidiInner<K, V, C> {
    /// Retrieves the value and timestamp associated with a given key.
    #[instrument(skip(self, key))]
    pub async fn get<Q>(&self, key: &Q) -> Result<(V, u128), Error>
    where
        K: Borrow<Q>,
        Q: Serialize + Encode + ?Sized,
    {
        self.perform_get(key, &self.uo).await
    }

    /// Retrieves the key and timestamp associated with a given value.
    #[instrument(skip(self, value))]
    pub async fn get_inverse<Q>(&self, value: &Q) -> Result<(K, u128), Error>
    where
        V: Borrow<Q>,
        Q: Serialize + Encode + ?Sized,
    {
        self.perform_get(value, &self.ou).await
    }

    /// Inserts or updates a key-value pair with a specified timestamp.
    #[instrument(skip(self, key, value))]
    pub async fn set(&self, key: K, value: V, timestamp: u128) -> Result<(), Error> {
        let uo = self.uo.clone();
        let ou = self.ou.clone();
        let codec = self.codec.clone();

        let was_inserted = tokio::task::spawn_blocking(move || -> Result<bool, Error> {
            let key_bytes = codec.encode(&key)?;
            let val_bytes = codec.encode(&value)?;
            let encoded_payload =
                bincode::encode_to_vec((val_bytes.clone(), timestamp), bincode_standard())
                    .map_err(|e| Error::Codec {
                        source: Box::new(e),
                    })?;
            let encoded_inverse_payload =
                bincode::encode_to_vec((key_bytes.clone(), timestamp), bincode_standard())
                    .map_err(|e| Error::Codec {
                        source: Box::new(e),
                    })?;

            let transaction_result = (&uo, &ou).transaction(move |(uo_tx, ou_tx)| {
                let is_new;
                if let Some(old_payload_ivec) = uo_tx.get(key_bytes.as_slice())? {
                    is_new = false;
                    let (old_val_bytes, _): (Vec<u8>, u128) =
                        bincode::decode_from_slice(&old_payload_ivec, bincode_standard())
                            .map(|(decoded, _)| decoded)
                            .map_err(|e| {
                                ConflictableTransactionError::Abort(Error::Codec {
                                    source: Box::new(e),
                                })
                            })?;
                    ou_tx.remove(old_val_bytes.as_slice())?;
                } else {
                    is_new = true;
                }

                uo_tx.insert(key_bytes.as_slice(), IVec::from(encoded_payload.clone()))?;
                ou_tx.insert(
                    val_bytes.as_slice(),
                    IVec::from(encoded_inverse_payload.clone()),
                )?;

                Ok(is_new)
            });

            transaction_result.map_err(|e: TransactionError<Error>| {
                Error::SledTransaction {
                    source: Box::new(e),
                }
            })
        })
        .await??;

        if was_inserted {
            self.len.fetch_add(1, Ordering::SeqCst);
        }

        Ok(())
    }

    /// Deletes a pair by its key and returns the associated value.
    #[instrument(skip(self, key))]
    pub async fn delete<Q>(&self, key: &Q) -> Result<V, Error>
    where
        K: Borrow<Q>,
        Q: Serialize + Encode + ?Sized,
    {
        self.perform_delete(key, &self.uo, &self.ou).await
    }

    /// Deletes a pair by its value and returns the associated key.
    #[instrument(skip(self, value))]
    pub async fn delete_inverse<Q>(&self, value: &Q) -> Result<K, Error>
    where
        V: Borrow<Q>,
        Q: Serialize + Encode + ?Sized,
    {
        self.perform_delete(value, &self.ou, &self.uo).await
    }

    /// Iterates over all key-value pairs using a callback function.
    pub async fn range<F>(&self, mut f: F) -> Result<(), Error>
    where
        F: FnMut(&K, &V, u128) -> Result<bool, Error> + Send + 'static,
    {
        let mut stream = self.range_stream();
        while let Some(item) = stream.next().await {
            let (k, v, ts) = item?;
            // The callback now receives references, avoiding clones in the loop.
            if !f(&k, &v, ts)? {
                break;
            }
        }
        Ok(())
    }

    /// Returns a stream over all key-value pairs in the map.
    pub fn range_stream(&self) -> impl Stream<Item = Result<(K, V, u128), Error>> + Send {
        let codec = self.codec.clone();
        let uo = self.uo.clone();
        let (tx, rx) = mpsc::channel(128);

        tokio::task::spawn_blocking(move || {
            for item in uo.iter() {
                let result = (|| {
                    let (k_ivec, payload_ivec) = item?;
                    let (v_b, ts): (Vec<u8>, u128) =
                        bincode::decode_from_slice(&payload_ivec, bincode_standard())
                            .map(|(decoded, _)| decoded)
                            .map_err(|e| Error::Codec {
                                source: Box::new(e),
                            })?;
                    let k: K = codec.decode(&k_ivec)?;
                    let v: V = codec.decode(&v_b)?;
                    Ok((k, v, ts))
                })();

                if tx.blocking_send(result).is_err() {
                    break;
                }
            }
        });

        ReceiverStream::new(rx)
    }

    /// Returns the number of elements in the map.
    pub fn len(&self) -> usize {
        self.len.load(Ordering::Relaxed)
    }

    /// Flushes all pending writes to the disk, ensuring durability.
    pub async fn flush(&self) -> Result<(), Error> {
        let db = self.db.clone();
        tokio::task::spawn_blocking(move || db.flush()).await??;
        Ok(())
    }

    // --- Private Helper Methods ---

    /// Internal helper to get a value from a tree given a key.
    async fn perform_get<Input, Output>(
        &self,
        input: &Input,
        tree: &Tree,
    ) -> Result<(Output, u128), Error>
    where
        Input: Serialize + Encode + ?Sized,
        Output: DeserializeOwned + Decode<()> + Send + 'static,
    {
        let t = tree.clone();
        let codec = self.codec.clone();
        let encoded_input = self.codec.encode(input)?;

        tokio::task::spawn_blocking(move || -> Result<(Output, u128), Error> {
            let payload_ivec = t.get(encoded_input)?.ok_or(Error::NotFound)?;

            let (output_bytes, ts): (Vec<u8>, u128) =
                bincode::decode_from_slice(&payload_ivec, bincode_standard())
                    .map(|(decoded, _)| decoded)
                    .map_err(|e| Error::Codec {
                        source: Box::new(e),
                    })?;
            let output = codec.decode(&output_bytes)?;
            Ok((output, ts))
        })
        .await?
    }

    /// Internal helper to delete an entry from a primary tree and its corresponding
    /// entry from the inverse tree.
    async fn perform_delete<Input, Output>(
        &self,
        input: &Input,
        primary_tree: &Tree,
        inverse_tree: &Tree,
    ) -> Result<Output, Error>
    where
        Input: Serialize + Encode + ?Sized,
        Output: DeserializeOwned + Decode<()> + Send + 'static,
    {
        let pt = primary_tree.clone();
        let it = inverse_tree.clone();
        let codec = self.codec.clone();
        let encoded_input = self.codec.encode(input)?;

        let deleted_bytes =
            tokio::task::spawn_blocking(move || -> Result<Option<Vec<u8>>, Error> {
                let transaction_result = (&pt, &it).transaction(move |(primary_tx, inverse_tx)| {
                    if let Some(payload_ivec) = primary_tx.remove(encoded_input.as_slice())? {
                        let (inverse_key_bytes, _): (Vec<u8>, u128) =
                            bincode::decode_from_slice(&payload_ivec, bincode_standard())
                                .map(|(decoded, _)| decoded)
                                .map_err(|e| {
                                    ConflictableTransactionError::Abort(Error::Codec {
                                        source: Box::new(e),
                                    })
                                })?;
                        inverse_tx.remove(inverse_key_bytes.as_slice())?;
                        Ok(Some(inverse_key_bytes))
                    } else {
                        Ok(None)
                    }
                });

                transaction_result.map_err(|e: TransactionError<Error>| {
                    Error::SledTransaction {
                        source: Box::new(e),
                    }
                })
            })
            .await??;

        if let Some(bytes_to_decode) = deleted_bytes {
            self.len.fetch_sub(1, Ordering::SeqCst);
            codec.decode(&bytes_to_decode)
        } else {
            Err(Error::NotFound)
        }
    }
}

#[cfg(test)]
mod integration_tests {
    use super::*;
    use futures::stream::StreamExt;
    use std::collections::HashMap;
    use std::fs;
    use tokio::task::JoinSet;

    const TEST_DB_BASE_PATH: &str = "./test_db";

    struct TestGuard {
        path: String,
    }

    impl Drop for TestGuard {
        fn drop(&mut self) {
            let _ = fs::remove_dir_all(&self.path);
        }
    }

    fn setup(test_name: &str) -> (String, TestGuard) {
        let path = format!("{}/{}", TEST_DB_BASE_PATH, test_name);
        let _ = fs::remove_dir_all(&path);
        let guard = TestGuard { path: path.clone() };
        (path, guard)
    }

    #[tokio::test]
    async fn test_crud_and_len() {
        let (path, _guard) = setup("crud_and_len");
        let map = BidiBuilder::new(&path).build().await.unwrap();
        assert_eq!(map.len(), 0);

        map.set("alpha".to_string(), "one".to_string(), 123)
            .await
            .unwrap();
        assert_eq!(map.len(), 1);

        // Now we can get by &str
        let (v, ts) = map.get("alpha").await.unwrap();
        assert_eq!(v, "one");
        assert_eq!(ts, 123);

        map.set("alpha".to_string(), "uno".to_string(), 456)
            .await
            .unwrap();
        assert_eq!(map.len(), 1);
        let (v2, ts2) = map.get("alpha").await.unwrap();
        assert_eq!(v2, "uno");
        assert_eq!(ts2, 456);
        assert!(map.get_inverse("one").await.is_err());

        let removed = map.delete("alpha").await.unwrap();
        assert_eq!(removed, "uno");
        assert_eq!(map.len(), 0);
        assert!(map.get("alpha").await.is_err());
    }

    #[tokio::test]
    async fn test_delete_inverse() {
        let (path, _guard) = setup("delete_inverse");
        let map = BidiBuilder::new(&path).build().await.unwrap();
        map.set("a".to_string(), "1".to_string(), 1).await.unwrap();
        assert_eq!(map.len(), 1);

        let removed_key = map.delete_inverse("1").await.unwrap();
        assert_eq!(removed_key, "a");
        assert_eq!(map.len(), 0);
        assert!(map.get("a").await.is_err());
        assert!(map.get_inverse("1").await.is_err());
    }

    #[tokio::test]
    async fn test_range_callback() {
        let (path, _guard) = setup("range_callback");
        let map = BidiBuilder::new(&path).build().await.unwrap();
        let mut expected = HashMap::new();
        for i in 0..10 {
            let k = format!("key{}", i);
            let v = format!("val{}", i);
            map.set(k.clone(), v.clone(), i as u128).await.unwrap();
            expected.insert(k, (v, i as u128));
        }

        let found = Arc::new(parking_lot::Mutex::new(HashMap::new()));
        let found_clone = found.clone();

        map.range(move |k, v, ts| {
            // Callback receives references, so we clone here to store them.
            found_clone.lock().insert(k.clone(), (v.clone(), ts));
            Ok(true)
        })
        .await
        .unwrap();

        assert_eq!(*found.lock(), expected);
    }

    #[tokio::test]
    async fn test_range_stream() {
        let (path, _guard) = setup("range_stream");
        let map = BidiBuilder::new(&path).build().await.unwrap();
        let mut expected = HashMap::new();
        for i in 0..10 {
            let k = format!("key{}", i);
            let v = format!("val{}", i);
            map.set(k.clone(), v.clone(), i as u128).await.unwrap();
            expected.insert(k, (v, i as u128));
        }

        let mut count = 0;
        let mut stream = Box::pin(map.range_stream());
        while let Some(Ok((k, v, ts))) = stream.next().await {
            let (expected_v, expected_ts) = expected.remove(&k).unwrap();
            assert_eq!(v, expected_v);
            assert_eq!(ts, expected_ts);
            count += 1;
        }
        assert_eq!(count, 10);
        assert!(expected.is_empty());
    }

    #[tokio::test]
    async fn test_disable_scan_on_startup() {
        let (path, _guard) = setup("disable_scan_on_startup");
        {
            let map = BidiBuilder::<String, String>::new(&path)
                .build()
                .await
                .unwrap();
            map.set("a".to_string(), "1".to_string(), 1).await.unwrap();
            map.set("b".to_string(), "2".to_string(), 2).await.unwrap();
            map.flush().await.unwrap();
        }

        let map = BidiBuilder::<String, String>::new(&path)
            .disable_scan_on_startup()
            .build()
            .await
            .unwrap();

        assert_eq!(map.len(), 0);

        let (v, _) = map.get("a").await.unwrap();
        assert_eq!(v, "1");
    }

    #[tokio::test]
    async fn test_concurrent_access() {
        let (path, _guard) = setup("concurrent_access");
        let map = BidiBuilder::new(&path).build().await.unwrap();

        let num_items = 100;
        let items: Vec<_> = (0..num_items)
            .map(|i| (format!("key-{}", i), format!("val-{}", i), i as u128))
            .collect();

        let mut set = JoinSet::new();
        for (k, v, ts) in items.clone() {
            let map = map.clone();
            set.spawn(async move {
                map.set(k, v, ts).await.unwrap();
            });
        }
        while let Some(res) = set.join_next().await {
            res.unwrap();
        }
        assert_eq!(map.len(), num_items as usize);

        let mut set = JoinSet::new();
        for (k, v, ts) in items.clone() {
            let map = map.clone();
            set.spawn(async move {
                let (read_v, read_ts) = map.get(k.as_str()).await.unwrap();
                assert_eq!(read_v, v);
                assert_eq!(read_ts, ts);

                let (read_k, read_ts_inv) = map.get_inverse(v.as_str()).await.unwrap();
                assert_eq!(read_k, k);
                assert_eq!(read_ts_inv, ts);
            });
        }
        while let Some(res) = set.join_next().await {
            res.unwrap();
        }

        let mut set = JoinSet::new();
        for (i, (k, v, _)) in items.iter().enumerate() {
            let map = map.clone();
            let k = k.clone();
            let v = v.clone();
            set.spawn(async move {
                if i % 2 == 0 {
                    let deleted_v = map.delete(k.as_str()).await.unwrap();
                    assert_eq!(deleted_v, v);
                } else {
                    let deleted_k = map.delete_inverse(v.as_str()).await.unwrap();
                    assert_eq!(deleted_k, k);
                }
            });
        }
        while let Some(res) = set.join_next().await {
            res.unwrap();
        }

        assert_eq!(map.len(), 0);
    }
}
