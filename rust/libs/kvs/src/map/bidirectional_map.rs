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

use serde::Serialize;
use sled::{
    Db, IVec, Tree,
    transaction::{ConflictableTransactionError, TransactionError, Transactional},
};
use std::{
    borrow::Borrow,
    sync::{Arc, atomic::AtomicUsize},
};
use tracing::instrument;
use wincode::SchemaWrite;

use crate::map::{
    base::MapBase,
    codec::Codec,
    error::Error,
    types::{KeyType, ValueType},
};

const PRIMARY_TREE_NAME: &str = "uo";
const SECONDARY_TREE_NAME: &str = "ou";
/// A map implementation that allows for efficient lookups of both keys and values.
///
/// This is achieved by maintaining two separate Sled trees: a primary tree for key-to-value mappings
/// and a secondary tree for value-to-key mappings. This allows for efficient inverse lookups.
/// It uses the `Map` trait for its core functionality and adds methods for inverse operations.
pub struct BidirectionalMap<K: KeyType, V: ValueType, C: Codec> {
    db: Arc<Db>,
    primary_tree: Tree,
    secondary_tree: Tree,
    len: AtomicUsize,
    codec: Arc<C>,
    _marker: std::marker::PhantomData<(K, V)>,
}

impl<K: KeyType, V: ValueType, C: Codec> MapBase for BidirectionalMap<K, V, C> {
    type K = K;
    type V = V;
    type C = C;

    fn _db(&self) -> &Arc<Db> {
        &self.db
    }

    fn _tree(&self) -> &sled::Tree {
        &self.primary_tree
    }

    fn _len(&self) -> &AtomicUsize {
        &self.len
    }

    fn _codec(&self) -> &Arc<C> {
        &self.codec
    }

    #[instrument(skip(self, key))]
    fn get<Q>(&self, key: &Q) -> impl Future<Output = Result<(Self::V, u128), Error>> + Send
    where
        Self::K: Borrow<Q>,
        Q: Serialize + SchemaWrite<Src = Q> + ?Sized + Sync,
    {
        self.perform_get(key, &self.primary_tree)
    }

    #[instrument(skip(self, key, value, timestamp))]
    fn set(
        &self,
        key: Self::K,
        value: Self::V,
        timestamp: u128,
    ) -> impl Future<Output = Result<(), Error>> + Send {
        let pt = self.primary_tree.clone();
        let st = self.secondary_tree.clone();
        let f = set_transaction_func(pt, st);
        self.perform_set(key, value, timestamp, f)
    }

    #[instrument(skip(self, key))]
    fn delete<Q>(&self, key: &Q) -> impl Future<Output = Result<Self::V, Error>> + Send
    where
        Self::K: Borrow<Q>,
        Q: Serialize + SchemaWrite<Src = Q> + ?Sized + Sync,
    {
        let pt = self.primary_tree.clone();
        let st = self.secondary_tree.clone();
        let f = delete_transaction_func(pt, st);
        self.perform_delete(key, f)
    }

    fn new(db: Db, scan_on_startup: bool, codec: C) -> Result<BidirectionalMap<K, V, C>, Error> {
        let pt = db.open_tree(PRIMARY_TREE_NAME)?;
        let st = db.open_tree(SECONDARY_TREE_NAME)?;
        let initial_len = if scan_on_startup { pt.len() } else { 0 };

        Ok(BidirectionalMap {
            db: Arc::new(db),
            primary_tree: pt,
            secondary_tree: st,
            len: AtomicUsize::new(initial_len),
            codec: Arc::new(codec),
            _marker: std::marker::PhantomData,
        })
    }
}

impl<K: KeyType, V: ValueType, C: Codec> BidirectionalMap<K, V, C> {
    /// Retrieves the key and timestamp associated with a given value.
    #[instrument(skip(self, value))]
    pub fn get_inverse<Q>(&self, value: &Q) -> impl Future<Output = Result<(K, u128), Error>> + Send
    where
        V: Borrow<Q>,
        Q: Serialize + SchemaWrite<Src = Q> + ?Sized + Sync,
    {
        self.perform_get(value, &self.secondary_tree)
    }

    /// Deletes a pair by its value and returns the associated key.
    #[instrument(skip(self, value))]
    pub fn delete_inverse<Q>(&self, value: &Q) -> impl Future<Output = Result<K, Error>> + Send
    where
        V: Borrow<Q>,
        Q: Serialize + wincode::SchemaWrite<Src = Q> + ?Sized + Sync,
    {
        let pt = self.primary_tree.clone();
        let st = self.secondary_tree.clone();
        let f = delete_transaction_func(st, pt);
        self.perform_delete(value, f)
    }
}

// --- Private Helper Methods ---

/// function generator for set transaction
fn set_transaction_func(
    t1: Tree,
    t2: Tree,
) -> impl FnOnce(Vec<u8>, Vec<u8>, u128) -> Result<bool, TransactionError<Error>> + Send + 'static {
    move |key: Vec<u8>, value: Vec<u8>, timestamp: u128| -> Result<bool, TransactionError<Error>> {
        let encoded_key_payload = wincode::serialize(&(key.clone(), timestamp)).map_err(|e| {
            TransactionError::Abort(Error::Codec {
                source: Box::new(e),
            })
        })?;
        let encoded_val_payload = wincode::serialize(&(value.clone(), timestamp)).map_err(|e| {
            TransactionError::Abort(Error::Codec {
                source: Box::new(e),
            })
        })?;
        (&t1, &t2).transaction(move |(tx1, tx2)| {
            let is_new;
            if let Some(old_payload_ivec) = tx1.get(key.as_slice())? {
                is_new = false;
                let (old_val_bytes, _): (Vec<u8>, u128) = wincode::deserialize(&old_payload_ivec)
                    .map(|decoded| decoded)
                    .map_err(|e| {
                        ConflictableTransactionError::Abort(Error::Codec {
                            source: Box::new(e),
                        })
                    })?;
                tx2.remove(old_val_bytes.as_slice())?;
            } else {
                is_new = true;
            }

            tx1.insert(key.as_slice(), IVec::from(encoded_val_payload.clone()))?;
            tx2.insert(value.as_slice(), IVec::from(encoded_key_payload.clone()))?;

            Ok(is_new)
        })
    }
}

/// function generator for delete transaction
fn delete_transaction_func(
    t1: Tree,
    t2: Tree,
) -> impl FnOnce(Vec<u8>) -> Result<Option<Vec<u8>>, TransactionError<Error>> + Send + 'static {
    move |key: Vec<u8>| -> Result<Option<Vec<u8>>, TransactionError<Error>> {
        (&t1, &t2).transaction(move |(tx1, tx2)| {
            if let Some(payload_ivec) = tx1.remove(key.as_slice())? {
                let (inverse_key_bytes, _): (Vec<u8>, u128) = wincode::deserialize(&payload_ivec)
                    .map(|decoded| decoded)
                    .map_err(|e| {
                        ConflictableTransactionError::Abort(Error::Codec {
                            source: Box::new(e),
                        })
                    })?;
                tx2.remove(inverse_key_bytes.as_slice())?;
                Ok(Some(inverse_key_bytes))
            } else {
                Ok(None)
            }
        })
    }
}
