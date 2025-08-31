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
use std::{path::Path, sync::Arc};

pub mod codec;
pub mod error;
pub mod map;
pub mod types;

mod bidirectional_map;
mod unidirectional_map;

use crate::{
    codec::{BincodeCodec, Codec},
    error::Error,
    map::Map,
};

/// A builder for creating `Map` instances.
///
/// This builder allows for configuration of the map before it is created,
/// such as setting the path, codec, and startup behavior.
pub struct MapBuilder<M: Map, C: Codec = BincodeCodec> {
    path: String,
    codec: C,
    config: Config,
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
            config: Config::default(),
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
            config: self.config,
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

    /// https://docs.rs/sled/0.34.7/sled/struct.Config.html#method.cache_capacity
    pub fn cache_capacity(mut self, to: u64) -> Self {
        self.config = self.config.cache_capacity(to);
        self
    }

    /// https://docs.rs/sled/0.34.7/sled/struct.Config.html#method.mode
    pub fn mode(mut self, to: Mode) -> Self {
        self.config = self.config.mode(to);
        self
    }

    /// https://docs.rs/sled/0.34.7/sled/struct.Config.html#method.use_compression
    pub fn use_compression(mut self, to: bool) -> Self {
        self.config = self.config.use_compression(to);
        self
    }

    /// https://docs.rs/sled/0.34.7/sled/struct.Config.html#method.compression_factor
    pub fn compression_factor(mut self, to: i32) -> Self {
        self.config = self.config.compression_factor(to);
        self
    }

    /// https://docs.rs/sled/0.34.7/sled/struct.Config.html#method.print_profile_on_drop
    pub fn print_profile_on_drop(mut self, to: bool) -> Self {
        self.config = self.config.print_profile_on_drop(to);
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

        let db = tokio::task::spawn_blocking(move || {
            self.config.path(Path::new(&self.path)).open()
        }).await??;

        let map = Arc::new(M::new(db, self.scan_on_startup, self.codec)?);

        Ok(map)
    }
}

/// A type alias for a `sled::Config`
pub type Config = sled::Config;

/// A type alias for a `sled::Mode`
pub type Mode = sled::Mode;

/// A type alias for a `BidirectionalMap`
pub type BidirectionalMap<K, V, C> = bidirectional_map::BidirectionalMap<K, V, C>;

/// A type alias for a `MapBuilder` that creates a `BidirectionalMap`.
pub type BidirectionalMapBuilder<K, V, C> = MapBuilder<BidirectionalMap<K, V, C>, C>;

/// A type alias for a `UnidirectionalMap`
pub type UnidirectionalMap<K, V, C> = unidirectional_map::UnidirectionalMap<K, V, C>;

/// A type alias for a `MapBuilder` that creates a `UnidirectionalMap`.
pub type UnidirectionalMapBuilder<K, V, C> = MapBuilder<UnidirectionalMap<K, V, C>, C>;

#[cfg(test)]
mod integration_tests;
