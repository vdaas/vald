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

use futures::{Stream, StreamExt};
use serde::{Serialize, de::DeserializeOwned};
use sled::{Db, Tree, transaction::TransactionError};
use std::borrow::Borrow;
use std::sync::Arc;
use std::sync::atomic::{AtomicUsize, Ordering};
use tokio::sync::mpsc;
use tokio_stream::wrappers::ReceiverStream;
use tracing::instrument;
use wincode::{SchemaRead, SchemaWrite};

use crate::map::{
    codec::Codec,
    error::Error,
    types::{KeyType, ValueType},
};

/// A trait that defines the core functionality of a key-value map.
///
/// This trait provides an abstraction over different map implementations,
/// such as unidirectional and bidirectional maps. It defines the essential
/// operations required for a key-value store, including getting, setting,
/// and deleting key-value pairs, as well as iterating over the map's contents.
pub trait MapBase: Sized + Sync + Send + 'static {
    /// The key type for the map.
    type K: KeyType;
    /// The value type for the map.
    type V: ValueType;
    /// The codec used for serializing and deserializing keys and values.
    type C: Codec;

    /// Creates a new map instance.
    ///
    /// # Arguments
    ///
    /// * `db` - The Sled database instance.
    /// * `scan_on_startup` - A boolean indicating whether to scan the database on startup.
    /// * `codec` - The codec to use for serialization and deserialization.
    fn new(db: Db, scan_on_startup: bool, codec: Self::C) -> Result<Self, Error>;

    /// Returns a reference to the underlying Sled database instance.
    fn _db(&self) -> &Arc<Db>;

    /// Returns a reference to the underlying Sled tree.
    fn _tree(&self) -> &sled::Tree;

    /// Returns a reference to the atomic counter for the map's length.
    fn _len(&self) -> &AtomicUsize;

    /// Returns a reference to the codec used by the map.
    fn _codec(&self) -> &Arc<Self::C>;

    /// Retrieves the value and timestamp associated with a given key.
    ///
    /// # Arguments
    ///
    /// * `key` - The key to retrieve the value for.
    fn get<Q>(&self, key: &Q) -> impl Future<Output = Result<(Self::V, u128), Error>> + Send
    where
        Self::K: Borrow<Q>,
        Q: Serialize + SchemaWrite<Src = Q> + ?Sized + Sync;

    /// Inserts or updates a key-value pair with a specified timestamp.
    ///
    /// # Arguments
    ///
    /// * `key` - The key to insert or update.
    /// * `value` - The value to associate with the key.
    /// * `timestamp` - The timestamp for the key-value pair.
    fn set(
        &self,
        key: Self::K,
        value: Self::V,
        timestamp: u128,
    ) -> impl Future<Output = Result<(), Error>> + Send;

    /// Deletes a pair by its key and returns the associated value.
    ///
    /// # Arguments
    ///
    /// * `key` - The key of the pair to delete.
    fn delete<Q>(&self, key: &Q) -> impl Future<Output = Result<Self::V, Error>> + Send
    where
        Self::K: Borrow<Q>,
        Q: Serialize + SchemaWrite<Src = Q> + ?Sized + Sync;

    /// Iterates over all key-value pairs using a callback function.
    ///
    /// The callback function `f` receives the key, value, and timestamp for each pair.
    /// If the callback returns `Ok(false)`, the iteration stops.
    ///
    /// # Arguments
    ///
    /// * `f` - The callback function to apply to each key-value pair.
    #[instrument(skip(self, f))]
    fn range<F>(&self, mut f: F) -> impl Future<Output = Result<(), Error>> + Send
    where
        F: FnMut(&Self::K, &Self::V, u128) -> Result<bool, Error> + Send + 'static,
    {
        let mut stream = Box::pin(self.range_stream());
        async move {
            while let Some(item) = stream.next().await {
                let (k, v, ts) = item?;
                if !f(&k, &v, ts)? {
                    break;
                }
            }
            Ok(())
        }
    }

    /// Returns a stream over all key-value pairs in the map.
    #[instrument(skip(self))]
    fn range_stream(&self) -> impl Stream<Item = Result<(Self::K, Self::V, u128), Error>> + Send {
        let codec = self._codec().clone();
        let tree = self._tree().clone();
        let (tx, rx) = mpsc::channel(128);

        tokio::task::spawn_blocking(move || {
            for item in tree.iter() {
                let result = (|| {
                    let (k_ivec, payload_ivec) = item?;
                    let (v_b, ts): (Vec<u8>, u128) = wincode::deserialize(&payload_ivec)
                        .map_err(|e| Error::Codec {
                            source: Box::new(e),
                        })?;
                    let k: Self::K = codec.decode(&k_ivec)?;
                    let v: Self::V = codec.decode(&v_b)?;
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
    #[instrument(skip(self))]
    fn len(&self) -> usize {
        self._len().load(Ordering::Relaxed)
    }

    /// Flushes all pending writes to the disk, ensuring durability.
    #[instrument(skip(self))]
    fn flush(&self) -> impl Future<Output = Result<(), Error>> + Send {
        let db = self._db().clone();
        async {
            tokio::task::spawn_blocking(move || db.flush()).await??;
            Ok(())
        }
    }

    // --- Private Helper Methods ---

    /// Internal helper to get a value from db given a key.
    fn perform_get<Input, Output>(
        &self,
        key: &Input,
        tree: &Tree,
    ) -> impl Future<Output = Result<(Output, u128), Error>> + Send
    where
        Input: Serialize + SchemaWrite<Src = Input> + ?Sized + Sync,
        Output: DeserializeOwned + for<'de> SchemaRead<'de, Dst = Output> + Send + 'static,
    {
        let tree = tree.clone();
        let codec = self._codec().clone();

        async move {
            let encoded_input = codec.encode(key)?;
            tokio::task::spawn_blocking(move || -> Result<(Output, u128), Error> {
                let payload_ivec = tree.get(encoded_input)?.ok_or(Error::NotFound)?;

                let (output_bytes, ts): (Vec<u8>, u128) = wincode::deserialize(&payload_ivec)
                    .map_err(|e| Error::Codec {
                        source: Box::new(e),
                    })?;
                let output = codec.decode(&output_bytes)?;
                Ok((output, ts))
            })
            .await?
        }
    }

    /// Internal helper to set an entry to db.
    fn perform_set<F>(
        &self,
        key: Self::K,
        value: Self::V,
        timestamp: u128,
        f: F,
    ) -> impl Future<Output = Result<(), Error>> + Send
    where
        F: FnOnce(Vec<u8>, Vec<u8>, u128) -> Result<bool, TransactionError<Error>> + Send + 'static,
    {
        let codec = self._codec().clone();

        async move {
            let key_bytes = codec.encode(&key)?;
            let val_bytes = codec.encode(&value)?;
            let was_inserted = tokio::task::spawn_blocking(move || -> Result<bool, Error> {
                let transaction_result = f(key_bytes, val_bytes, timestamp);

                transaction_result.map_err(|e: TransactionError<Error>| Error::SledTransaction {
                    source: Box::new(e),
                })
            })
            .await??;

            if was_inserted {
                self._len().fetch_add(1, Ordering::SeqCst);
            }

            Ok(())
        }
    }

    /// Internal helper to delete an entry from db
    fn perform_delete<Input, Output, F>(
        &self,
        key: &Input,
        f: F,
    ) -> impl Future<Output = Result<Output, Error>> + Send
    where
        Input: Serialize + SchemaWrite<Src = Input> + ?Sized + Sync,
        Output: DeserializeOwned + for<'de> SchemaRead<'de, Dst = Output> + Send + 'static,
        F: FnOnce(Vec<u8>) -> Result<Option<Vec<u8>>, TransactionError<Error>> + Send + 'static,
    {
        let codec = self._codec().clone();

        async move {
            let encoded_input = codec.encode(key)?;
            let deleted_bytes = tokio::task::spawn_blocking(move || {
                let transaction_result = f(encoded_input);
                transaction_result.map_err(|e: TransactionError<Error>| Error::SledTransaction {
                    source: Box::new(e),
                })
            })
            .await??;

            if let Some(bytes_to_decode) = deleted_bytes {
                self._len().fetch_sub(1, Ordering::SeqCst);
                codec.decode(&bytes_to_decode)
            } else {
                Err(Error::NotFound)
            }
        }
    }
}
