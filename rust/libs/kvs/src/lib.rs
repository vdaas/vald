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

pub mod map;

use crate::map::{
    base::MapBase,
    codec::{BincodeCodec, Codec},
    error::Error,
};

/// A builder for creating `Map` instances.
///
/// This builder allows for configuration of the map before it is created,
/// such as setting the path, codec, and startup behavior.
pub struct MapBuilder<M: MapBase, C: Codec = BincodeCodec> {
    path: String,
    codec: C,
    config: Config,
    scan_on_startup: bool,
    _marker: std::marker::PhantomData<M>,
}

impl<M: MapBase<C = BincodeCodec>> MapBuilder<M, BincodeCodec> {
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

impl<M: MapBase<C = C>, C: Codec> MapBuilder<M, C> {
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

    /// Sets the cache capacity in bytes for the database.
    ///
    /// See: https://docs.rs/sled/0.34.7/sled/struct.Config.html#method.cache_capacity
    pub fn cache_capacity(mut self, to: u64) -> Self {
        self.config = self.config.cache_capacity(to);
        self
    }

    /// Sets the database access mode.
    ///
    /// See: https://docs.rs/sled/0.34.7/sled/struct.Config.html#method.mode
    pub fn mode(mut self, to: Mode) -> Self {
        self.config = self.config.mode(to);
        self
    }

    /// Enables or disables transparent zstd compression.
    ///
    /// See: https://docs.rs/sled/0.34.7/sled/struct.Config.html#method.use_compression
    pub fn use_compression(mut self, to: bool) -> Self {
        self.config = self.config.use_compression(to);
        self
    }

    /// Sets the zstd compression factor.
    ///
    /// See: https://docs.rs/sled/0.34.7/sled/struct.Config.html#method.compression_factor
    pub fn compression_factor(mut self, to: i32) -> Self {
        self.config = self.config.compression_factor(to);
        self
    }

    /// If set to `true`, prints performance statistics when the `Db` is dropped.
    ///
    /// See: https://docs.rs/sled/0.34.7/sled/struct.Config.html#method.print_profile_on_drop
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

        let db =
            tokio::task::spawn_blocking(move || self.config.path(Path::new(&self.path)).open())
                .await??;

        let map = Arc::new(M::new(db, self.scan_on_startup, self.codec)?);

        Ok(map)
    }
}

/// Re-exports `sled::Config` to allow for configuration of the `Map``.
pub use sled::Config;
/// Re-exports `sled::Mode` to allow for configuration of the `Map`.
pub use sled::Mode;

/// Re-exports the `BidirectionalMap` implementation.
///
/// A map that allows for lookups of both keys and values.
pub use crate::map::bidirectional_map::BidirectionalMap;

/// Re-exports the `UnidirectionalMap` implementation.
///
/// A standard key-value map.
pub use crate::map::unidirectional_map::UnidirectionalMap;

/// A type alias for a `MapBuilder` that creates a `BidirectionalMap`.
pub type BidirectionalMapBuilder<K, V, C> = MapBuilder<BidirectionalMap<K, V, C>, C>;

/// A type alias for a `MapBuilder` that creates a `UnidirectionalMap`.
pub type UnidirectionalMapBuilder<K, V, C> = MapBuilder<UnidirectionalMap<K, V, C>, C>;

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

    async fn test_crud_and_len<M: MapBase<K = String, V = String, C = BincodeCodec>>(
        path: &str,
    ) -> Arc<M> {
        let map = MapBuilder::<M>::new(path).build().await.unwrap();
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

        let removed = map.delete("alpha").await.unwrap();
        assert_eq!(removed, "uno");
        assert_eq!(map.len(), 0);
        assert!(map.get("alpha").await.is_err());

        map
    }

    #[tokio::test]
    async fn test_bidirectional_map_crud_and_len() {
        let (path, _guard) = setup("bidirectional_map_crud_and_len");

        let map = test_crud_and_len::<BidirectionalMap<String, String, BincodeCodec>>(&path).await;

        assert!(map.get_inverse("one").await.is_err());
    }

    #[tokio::test]
    async fn test_unidirectional_map_crud_and_len() {
        let (path, _guard) = setup("unidirectional_map_crud_and_len");

        test_crud_and_len::<UnidirectionalMap<String, String, BincodeCodec>>(&path).await;
    }

    #[tokio::test]
    async fn test_bidirectional_map_delete_inverse() {
        let (path, _guard) = setup("bidirectional_map_delete_inverse");
        let map = MapBuilder::<BidirectionalMap<String, String, BincodeCodec>>::new(&path)
            .build()
            .await
            .unwrap();
        map.set("a".to_string(), "1".to_string(), 1).await.unwrap();
        assert_eq!(map.len(), 1);

        let removed_key = map.delete_inverse("1").await.unwrap();
        assert_eq!(removed_key, "a");
        assert_eq!(map.len(), 0);
        assert!(map.get("a").await.is_err());
        assert!(map.get_inverse("1").await.is_err());
    }

    async fn test_range_callback<M: MapBase<K = String, V = String, C = BincodeCodec>>(path: &str) {
        let map = MapBuilder::<M>::new(&path).build().await.unwrap();
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
    async fn test_bidirectional_map_range_callback() {
        let (path, _guard) = setup("bidirectional_map_range_callback");
        test_range_callback::<BidirectionalMap<String, String, BincodeCodec>>(&path).await;
    }

    #[tokio::test]
    async fn test_unidirectional_map_range_callback() {
        let (path, _guard) = setup("unidirectional_map_range_callback");
        test_range_callback::<UnidirectionalMap<String, String, BincodeCodec>>(&path).await;
    }

    async fn test_range_stream<M: MapBase<K = String, V = String, C = BincodeCodec>>(path: &str) {
        let map = MapBuilder::<M>::new(path).build().await.unwrap();
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
    async fn test_bidirectional_map_range_stream() {
        let (path, _guard) = setup("bidirectional_map_range_stream");
        test_range_stream::<BidirectionalMap<String, String, BincodeCodec>>(&path).await;
    }

    #[tokio::test]
    async fn test_unidirectional_map_range_stream() {
        let (path, _guard) = setup("unidirectional_map_range_stream");
        test_range_stream::<UnidirectionalMap<String, String, BincodeCodec>>(&path).await;
    }

    async fn test_disable_scan_on_startup<M: MapBase<K = String, V = String, C = BincodeCodec>>(
        path: &str,
    ) {
        {
            let map = MapBuilder::<M>::new(&path).build().await.unwrap();
            map.set("a".to_string(), "1".to_string(), 1).await.unwrap();
            map.set("b".to_string(), "2".to_string(), 2).await.unwrap();
            map.flush().await.unwrap();
        }

        let map = MapBuilder::<M>::new(&path)
            .disable_scan_on_startup()
            .build()
            .await
            .unwrap();

        assert_eq!(map.len(), 0);

        let (v, _) = map.get("a").await.unwrap();
        assert_eq!(v, "1");
    }

    #[tokio::test]
    async fn test_bidirectional_map_disable_scan_on_startup() {
        let (path, _guard) = setup("bidirectional_map_disable_scan_on_startup");
        test_disable_scan_on_startup::<BidirectionalMap<String, String, BincodeCodec>>(&path).await;
    }

    #[tokio::test]
    async fn test_unidirectional_map_disable_scan_on_startup() {
        let (path, _guard) = setup("unidirectional_map_disable_scan_on_startup");
        test_disable_scan_on_startup::<UnidirectionalMap<String, String, BincodeCodec>>(&path)
            .await;
    }

    async fn test_concurrent_access<M, F1, F2, Fut1, Fut2>(path: &str, f1: F1, f2: F2)
    where
        M: MapBase<K = String, V = String, C = BincodeCodec>,
        F1: Fn(Arc<M>, String, String, u128) -> Fut1 + Send + Sync + Copy + 'static,
        F2: Fn(Arc<M>, String, String, usize) -> Fut2 + Send + Sync + Copy + 'static,
        Fut1: Future<Output = ()> + Send,
        Fut2: Future<Output = ()> + Send,
    {
        let map = MapBuilder::<M>::new(&path).build().await.unwrap();

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
                f1(map, k, v, ts).await;
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
                f2(map, k, v, i).await;
            });
        }
        while let Some(res) = set.join_next().await {
            res.unwrap();
        }

        assert_eq!(map.len(), 0);
    }

    #[tokio::test]
    async fn test_bidirectional_map_concurrent_access() {
        let (path, _guard) = setup("bidirectional_map_concurrent_access");
        let f1 = |map: Arc<BidirectionalMap<String, String, BincodeCodec>>,
                  key: String,
                  value: String,
                  timestamp: u128| async move {
            let (read_v, read_ts) = map.get(key.as_str()).await.unwrap();
            assert_eq!(read_v, value);
            assert_eq!(read_ts, timestamp);
            let (read_k, read_ts_inv) = map.get_inverse(value.as_str()).await.unwrap();
            assert_eq!(read_k, key);
            assert_eq!(read_ts_inv, timestamp);
        };
        let f2 = |map: Arc<BidirectionalMap<String, String, BincodeCodec>>,
                  key: String,
                  value: String,
                  i: usize| async move {
            if i % 2 == 0 {
                let deleted_v = map.delete(key.as_str()).await.unwrap();
                assert_eq!(deleted_v, value);
            } else {
                let deleted_k = map.delete_inverse(value.as_str()).await.unwrap();
                assert_eq!(deleted_k, key);
            }
        };
        test_concurrent_access(&path, f1, f2).await;
    }

    #[tokio::test]
    async fn test_unidirectional_map_concurrent_access() {
        let (path, _guard) = setup("unidirectional_map_concurrent_access");
        let f1 = |map: Arc<UnidirectionalMap<String, String, BincodeCodec>>,
                  key: String,
                  value: String,
                  timestamp: u128| async move {
            let (read_v, read_ts) = map.get(key.as_str()).await.unwrap();
            assert_eq!(read_v, value);
            assert_eq!(read_ts, timestamp);
        };
        let f2 = |map: Arc<UnidirectionalMap<String, String, BincodeCodec>>,
                  key: String,
                  value: String,
                  _: usize| async move {
            let deleted_v = map.delete(key.as_str()).await.unwrap();
            assert_eq!(deleted_v, value);
        };
        test_concurrent_access(&path, f1, f2).await;
    }
}
