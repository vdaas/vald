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

use bincode::{Decode, Encode, config::standard as bincode_standard};
use futures::{Stream, StreamExt};
use serde::{Serialize, de::DeserializeOwned};
use sled::{
    Db, IVec, Tree,
    transaction::{ConflictableTransactionError, TransactionError, Transactional},
};
use std::borrow::Borrow;
use std::fmt::Debug;
use std::hash::Hash;
use std::sync::Arc;
use std::sync::atomic::{AtomicUsize, Ordering};
use tokio::sync::mpsc;
use tokio_stream::wrappers::ReceiverStream;
use tracing::instrument;

use crate::codec::Codec;
use crate::error::Error;

/// The internal struct holding the state and logic of the `BidiMap`.
pub struct BidiInner<K, V, C: Codec> {
    db: Arc<Db>,
    uo: Tree,
    ou: Tree,
    len: AtomicUsize,
    codec: Arc<C>,
    _marker: std::marker::PhantomData<(K, V)>,
}

impl<K, V, C: Codec> BidiInner<K, V, C>
where
    K: Serialize
        + DeserializeOwned
        + Encode
        + Decode<()>
        + Eq
        + Hash
        + Clone
        + Send
        + Sync
        + Debug
        + 'static,
    V: Serialize
        + DeserializeOwned
        + Encode
        + Decode<()>
        + Eq
        + Hash
        + Clone
        + Send
        + Sync
        + Debug
        + 'static,
{
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
    
    /// A public constructor method
    pub fn new(db: Db, uo: Tree, ou: Tree, initial_len: usize, codec: C) -> BidiInner<K, V, C> {
        BidiInner {
            db: Arc::new(db),
            uo,
            ou,
            len: AtomicUsize::new(initial_len),
            codec: Arc::new(codec),
            _marker: std::marker::PhantomData,
        }
    }
}

