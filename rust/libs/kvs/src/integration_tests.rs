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

use crate::unidirectional_map::UnidirectionalMap;

use super::*;
use futures::stream::StreamExt;
use std::collections::HashMap;
use std::fs;
use tokio::task::JoinSet;
use bidirectional_map::BidirectionalMap;

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

async fn test_crud_and_len<M: Map<K = String, V = String, C = BincodeCodec>>(path: &str) -> Arc<M> {
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
    let map = MapBuilder::<BidirectionalMap<String, String, BincodeCodec>>::new(&path).build().await.unwrap();
    map.set("a".to_string(), "1".to_string(), 1).await.unwrap();
    assert_eq!(map.len(), 1);

    let removed_key = map.delete_inverse("1").await.unwrap();
    assert_eq!(removed_key, "a");
    assert_eq!(map.len(), 0);
    assert!(map.get("a").await.is_err());
    assert!(map.get_inverse("1").await.is_err());
}

async fn test_range_callback<M: Map<K = String, V = String, C = BincodeCodec>>(path: &str) {
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
    let _ = test_range_callback::<BidirectionalMap<String, String, BincodeCodec>>(&path);
}

#[tokio::test]
async fn test_unidirectional_map_range_callback() {
    let (path, _guard) = setup("unidirectional_map_range_callback");
    let _ = test_range_callback::<UnidirectionalMap<String, String, BincodeCodec>>(&path);
}

async fn test_range_stream<M: Map<K = String, V = String, C = BincodeCodec>>(path: &str) {
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
    let _ = test_range_stream::<BidirectionalMap<String, String, BincodeCodec>>(&path);
}

#[tokio::test]
async fn test_unidirectional_map_range_stream() {
    let (path, _guard) = setup("unidirectional_map_range_stream");
    let _ = test_range_stream::<UnidirectionalMap<String, String, BincodeCodec>>(&path);
}

async fn test_disable_scan_on_startup<M: Map<K = String, V = String, C = BincodeCodec>>(path: &str) {
    {
        let map = MapBuilder::<M>::new(&path)
            .build()
            .await
            .unwrap();
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
    let _ = test_disable_scan_on_startup::<BidirectionalMap<String, String, BincodeCodec>>(&path);
}

#[tokio::test]
async fn test_unidirectional_map_disable_scan_on_startup() {
    let (path, _guard) = setup("unidirectional_map_disable_scan_on_startup");
    let _ = test_disable_scan_on_startup::<UnidirectionalMap<String, String, BincodeCodec>>(&path);
}

async fn test_concurrent_access<M, F1, F2>(path: &str, mut f1: F1, mut f2: F2)
where
    M: Map<K = String, V = String, C = BincodeCodec>,
    F1: FnMut(Arc<M>, String, String, u128) -> () + Send + Copy + 'static,
    F2: FnMut(Arc<M>, String, String, usize) -> () + Send + Copy + 'static,
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
            let _ = map.set(k, v, ts).await;
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
            f1(map, k, v, ts)
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
            f2(map, k, v, i)
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
    let f1 = |map: Arc<BidirectionalMap<String, String, BincodeCodec>>, key: String, value: String, timestamp: u128|{
        let _ = async {
            let (read_v, read_ts) = map.get(key.as_str()).await.unwrap();
            assert_eq!(read_v, value);
            assert_eq!(read_ts, timestamp);
            let (read_k, read_ts_inv) = map.get_inverse(value.as_str()).await.unwrap();
            assert_eq!(read_k, key);
            assert_eq!(read_ts_inv, timestamp);
        };
    };
    let f2 = |map: Arc<BidirectionalMap<String, String, BincodeCodec>>, key: String, value: String, i: usize| {
        let _ = async {
            if i % 2 == 0 {
                let deleted_v = map.delete(key.as_str()).await.unwrap();
                assert_eq!(deleted_v, value);
            } else {
                let deleted_k = map.delete_inverse(value.as_str()).await.unwrap();
                assert_eq!(deleted_k, key);
            }
        };
    };
    let _ = test_concurrent_access(&path, f1, f2);
}

#[tokio::test]
async fn test_unidirectional_map_concurrent_access() {
    let (path, _guard) = setup("unidirectional_map_concurrent_access");
    let f1 = |map: Arc<UnidirectionalMap<String, String, BincodeCodec>>, key: String, value: String, timestamp: u128| {
        let _ = async {
            let (read_v, read_ts) = map.get(key.as_str()).await.unwrap();
            assert_eq!(read_v, value);
            assert_eq!(read_ts, timestamp);
        };
    };
    let f2 = |map: Arc<UnidirectionalMap<String, String, BincodeCodec>>, key: String, value: String, _: usize| {
        let _ = async {
            let deleted_v = map.delete(key.as_str()).await.unwrap();
            assert_eq!(deleted_v, value);
        };
    };
    let _ = test_concurrent_access(&path, f1, f2);
}
