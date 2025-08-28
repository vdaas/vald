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

use bincode::{config::standard as bincode_standard, Encode};
use serde::Serialize;
use sled::{
    Db, IVec, Tree,
    transaction::{ConflictableTransactionError, TransactionError},
};
use tracing::instrument;
use std::{borrow::Borrow, sync::Arc};
use std::sync::atomic::AtomicUsize;

use crate::{codec::Codec, error::Error, map::Map, types::{KeyType, ValueType}};


const TREE_NAME: &str = "tree";
/// A map implementation that stores key-value pairs in a single Sled tree.
///
/// This is a basic key-value store where each key is associated with a single value.
/// It uses the `Map` trait for its core functionality.
pub struct UnidirectionalMap<K: KeyType, V: ValueType, C: Codec> {
    db: Arc<Db>,
    tree: Tree,
    len: AtomicUsize,
    codec: Arc<C>,
    _marker: std::marker::PhantomData<(K, V)>,
}

impl<K: KeyType, V: ValueType, C: Codec> Map for UnidirectionalMap<K, V, C> {
    type K = K;
    type V = V;
    type C = C;

    fn _db(&self) -> &Arc<Db> {
        &self.db
    }
    
    fn _tree(&self) -> &sled::Tree {
        &self.tree
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
        Q: Serialize + Encode + ?Sized + Sync
    {
        self.perform_get(key, &self.tree)
    }

    #[instrument(skip(self, key, value))]
    fn set(&self, key: Self::K, value: Self::V, timestamp: u128) -> impl Future<Output = Result<(), Error>> + Send {
        let t = self.tree.clone();
        let f = set_transaction_func(t);
        self.perform_set(key, value, timestamp, f)
    }

    #[instrument(skip(self, key))]
    fn delete<Q>(&self, key: &Q) -> impl Future<Output = Result<Self::V, Error>> + Send
    where
        Self::K: Borrow<Q>,
        Q: Serialize + Encode + ?Sized + Sync
    {
        let t = self.tree.clone();
        let f = delete_transaction_func(t);
        self.perform_delete(key, f)
    }

    fn new(db: Db, scan_on_startup: bool, codec: C) -> Result<UnidirectionalMap<K, V, C>, Error> {
        let tree = db.open_tree(TREE_NAME)?;
        let initial_len = if scan_on_startup { tree.len() } else { 0 };

        Ok(UnidirectionalMap {
            db: Arc::new(db),
            tree: tree,
            len: AtomicUsize::new(initial_len),
            codec: Arc::new(codec),
            _marker: std::marker::PhantomData
        })
    }
}

fn set_transaction_func(t: Tree) -> impl FnOnce(Vec<u8>, Vec<u8>, u128) -> Result<bool, TransactionError<Error>> + Send + 'static {
    move |key: Vec<u8>, value: Vec<u8>, timestamp: u128| -> Result<bool, TransactionError<Error>> {
        let encoded_payload = bincode::encode_to_vec((value.clone(), timestamp), bincode_standard())
            .map_err(|e| TransactionError::Abort(Error::Codec {
                source: Box::new(e),
            }))?;
        (&t).transaction(move |tx| {
            let is_new= !tx.get(key.as_slice())?.is_some();
            tx.insert(key.as_slice(), IVec::from(encoded_payload.clone()))?;
            
            Ok(is_new)
        })
    }
}

fn delete_transaction_func(t: Tree) -> impl FnOnce(Vec<u8>) -> Result<Option<Vec<u8>>, TransactionError<Error>> + Send + 'static {
    move |key: Vec<u8>| -> Result<Option<Vec<u8>>, TransactionError<Error>> {
        (&t).transaction(move |tx| {
            if let Some(payload_ivec) = tx.remove(key.as_slice())? {
                let (inverse_key_bytes, _): (Vec<u8>, u128) =
                    bincode::decode_from_slice(&payload_ivec, bincode_standard())
                        .map(|(decoded, _)| decoded)
                        .map_err(|e| {
                            ConflictableTransactionError::Abort(Error::Codec {
                                source: Box::new(e),
                            })
                        })?;
                Ok(Some(inverse_key_bytes))
            } else {
                Ok(None)
            }
        })
    }
}
