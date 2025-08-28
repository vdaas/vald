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

//! # Map
//!
//! A generic, persistent, async-safe, and pluggable-codec `Map`.
//! It is designed for high-concurrency environments, offloading blocking I/O operations
//! to a dedicated thread pool managed by `tokio`. This version focuses on core
//! functionality without Time-To-Live (TTL) features.
//!
//! The implementation uses `sled` as its underlying persistent storage engine to leverage
//! its robust transactional capabilities, ensuring data consistency for bidirectional mappings.
use std::sync::Arc;

pub mod error;
use error::Error;

pub mod codec;
use codec::{BincodeCodec, Codec};

pub mod types;

pub mod map;
pub mod bidirectional_map;
pub mod unidirectional_map;

use crate::{bidirectional_map::BidirectionalMap, map::Map, unidirectional_map::UnidirectionalMap};

/// A builder for creating `Map` instances.
///
/// This builder allows for configuration of the map before it is created,
/// such as setting the path, codec, and startup behavior.
pub struct MapBuilder<M: Map, C: Codec = BincodeCodec> {
    path: String,
    codec: C,
    scan_on_startup: bool,
    _marker: std::marker::PhantomData<M>,
}

impl<M: Map<C = BincodeCodec>> MapBuilder<M, BincodeCodec> {
    /// Creates a new `MapBuilder` with the given path.
    ///
    /// By default, it uses `BincodeCodec` and scans the database on startup.
    ///
    /// # Arguments
    ///
    /// * `path` - The path to the Sled database file.
    pub fn new(path: impl AsRef<str>) -> Self {
        Self {
            path: path.as_ref().to_string(),
            codec: BincodeCodec::default(),
            scan_on_startup: true,
            _marker: std::marker::PhantomData,
        }
    }
}

impl<M: Map<C = C>, C: Codec> MapBuilder<M, C> {
    /// Sets a custom codec for the map, returning a new builder instance.
    ///
    /// # Arguments
    ///
    /// * `new_codec` - The new codec to use.
    pub fn codec<NewC: Codec>(self, new_codec: NewC) -> MapBuilder<M, NewC> {
        MapBuilder {
            path: self.path,
            codec: new_codec,
            scan_on_startup: self.scan_on_startup,
            _marker: std::marker::PhantomData,
        }
    }

    /// Disables the full database scan on startup.
    ///
    /// By default, the map scans the database to determine the number of items.
    /// Disabling this can speed up startup time for large databases.
    pub fn disable_scan_on_startup(mut self) -> Self {
        self.scan_on_startup = false;
        self
    }

    /// Builds the map, initializing the database.
    ///
    /// This method creates the database file and initializes the map.
    /// It returns an `Arc` wrapped map instance.
    pub async fn build(self) -> Result<Arc<M>, Error> {
        if let Some(dir) = std::path::Path::new(&self.path).parent() {
            tokio::fs::create_dir_all(dir).await?;
        }

        let path = self.path.clone();
        let db = tokio::task::spawn_blocking(move || sled::open(path)).await??;

        let inner = Arc::new(M::new(db, self.scan_on_startup, self.codec)?);

        Ok(inner)
    }
}

/// A type alias for a `MapBuilder` that creates a `BidirectionalMap`.
pub type BidirectionalMapBuilder<K, V, C> = MapBuilder<BidirectionalMap<K, V, C>, C>;

/// A type alias for a `MapBuilder` that creates a `UnidirectionalMap`.
pub type UnidirectionalMapBuilder<K, V, C> = MapBuilder<UnidirectionalMap<K, V, C>, C>;

#[cfg(test)]
mod integration_tests;
