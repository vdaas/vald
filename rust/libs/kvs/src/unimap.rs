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

use bincode::{config::standard as bincode_standard, Decode, Encode};
use futures::{Stream, StreamExt};
use serde::{Serialize, de::DeserializeOwned};
use sled::{
    Db, IVec, Tree,
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

/// The internal struct holding the state and log of `UniMap`.
pub struct UniInner<K, V, C: Codec> {
    db: Arc<Db>,
    tree: Tree,
    len: AtomicUsize,
    codec: Arc<C>,
    _marker: std::marker::PhantomData<(K, V)>,
}

impl<K, V, C: Codec> UniInner<K, V, C>
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
        let tree = self.tree.clone();
        let codec = self.codec.clone();
        let encoded_input = self.codec.encode(&key)?;

        tokio::task::spawn_blocking(move || -> Result<(_, u128), Error> {
            let payload_ivec = tree.get(encoded_input)?.ok_or(Error::NotFound)?;

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

        /// Inserts or updates a key-value pair with a specified timestamp.
    #[instrument(skip(self, key, value))]
    pub async fn set(&self, key: K, value: V, timestamp: u128) -> Result<(), Error> {
        let tree = self.tree.clone();
        let codec = self.codec.clone();

        let was_inserted = tokio::task::spawn_blocking(move || -> Result<_, Error> {
            let key_bytes = codec.encode(&key)?;
            let val_bytes = codec.encode(&value)?;
            let encoded_payload =
                bincode::encode_to_vec((val_bytes.clone(), timestamp), bincode_standard())
                    .map_err(|e| Error::Codec {
                        source: Box::new(e),
                    })?;

            Ok(tree.insert(key_bytes.as_slice(), IVec::from(encoded_payload.clone())))
        })
        .await??;

        if let Ok(_) = was_inserted {
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
        let tree = self.tree.clone();
        let encoded_input = self.codec.encode(key)?;
        let result = tree.remove(encoded_input.as_slice())?;

        if let Some(result) = result {
            self.len.fetch_sub(1, Ordering::SeqCst);
            self.codec.decode(&result)
        } else {
            Err(Error::NotFound)
        }
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
        let tree = self.tree.clone();
        let (tx, rx) = mpsc::channel(128);

        tokio::task::spawn_blocking(move || {
            for item in tree.iter() {
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

    /// A public constructor method
    pub fn new(db: Db, tree: Tree, initial_len: usize, codec: C) -> UniInner<K, V, C> {
        UniInner {
            db: Arc::new(db),
            tree,
            len: AtomicUsize::new(initial_len),
            codec: Arc::new(codec),
            _marker: std::marker::PhantomData,
        }
    }
}