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
use bincode::{Decode, Encode};
use serde::{Serialize, de::DeserializeOwned};

use std::fmt::Debug;
use std::hash::Hash;
use std::sync::Arc;

pub mod error;
use error::Error;

pub mod codec;
use codec::{BincodeCodec, Codec};

mod bidimap;
use bidimap::BidiInner;

mod unimap;
use unimap::UniInner;

/// A builder for constructing a `BidiMap` instance with custom configurations.
pub struct BidiBuilder<K, V, C: Codec = BincodeCodec> {
    path: String,
    codec: C,
    scan_on_startup: bool,
    _marker: std::marker::PhantomData<(K, V)>,
}

impl<K, V> BidiBuilder<K, V, BincodeCodec>
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

impl<K, V, C: Codec> BidiBuilder<K, V, C>
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

        let inner = Arc::new(BidiInner::new(db, uo, ou, initial_len, self.codec));

        Ok(inner)
    }
}

/// A builder for constructing a `UnidiMap` instance with custom configurations.
pub struct UnidiBuilder<K, V, C: Codec = BincodeCodec> {
    path: String,
    codec: C,
    scan_on_startup: bool,
    _marker: std::marker::PhantomData<(K, V)>,
}

impl<K, V> UnidiBuilder<K, V, BincodeCodec>
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
    /// Creates a new `UnidiBuilder` with a specified database path and the default `BincodeCodec`.
    pub fn new(path: impl AsRef<str>) -> Self {
        Self {
            path: path.as_ref().to_string(),
            codec: BincodeCodec::default(),
            scan_on_startup: true,
            _marker: std::marker::PhantomData,
        }
    }
}

impl<K, V, C: Codec> UnidiBuilder<K, V, C>
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
    /// Sets a custom codec for the `UnidiMap`, returning a new builder instance.
    pub fn codec<NewC: Codec>(self, new_codec: NewC) -> UnidiBuilder<K, V, NewC> {
        UnidiBuilder {
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

    /// Builds the `UnidiMap`, initializing the database.
    pub async fn build(self) -> Result<Arc<UniInner<K, V, C>>, Error> {
        if let Some(dir) = std::path::Path::new(&self.path).parent() {
            tokio::fs::create_dir_all(dir).await?;
        }

        let path = self.path.clone();
        let db = tokio::task::spawn_blocking(move || sled::open(path)).await??;

        let tree = db.open_tree("tree")?;

        let initial_len = if self.scan_on_startup { tree.len() } else { 0 };

        let inner = Arc::new(UniInner::new(db, tree, initial_len, self.codec));

        Ok(inner)
    }
}

#[cfg(test)]
mod integration_tests;
