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

#[tokio::test]
async fn test_bidimap_crud_and_len() {
    let (path, _guard) = setup("bidimap_crud_and_len");
    let map = BidiBuilder::new(&path).build().await.unwrap();
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
    assert!(map.get_inverse("one").await.is_err());

    let removed = map.delete("alpha").await.unwrap();
    assert_eq!(removed, "uno");
    assert_eq!(map.len(), 0);
    assert!(map.get("alpha").await.is_err());
}

#[tokio::test]
async fn test_bidimap_delete_inverse() {
    let (path, _guard) = setup("bidimap_delete_inverse");
    let map = BidiBuilder::new(&path).build().await.unwrap();
    map.set("a".to_string(), "1".to_string(), 1).await.unwrap();
    assert_eq!(map.len(), 1);

    let removed_key = map.delete_inverse("1").await.unwrap();
    assert_eq!(removed_key, "a");
    assert_eq!(map.len(), 0);
    assert!(map.get("a").await.is_err());
    assert!(map.get_inverse("1").await.is_err());
}

#[tokio::test]
async fn test_bidimap_range_callback() {
    let (path, _guard) = setup("bidimap_range_callback");
    let map = BidiBuilder::new(&path).build().await.unwrap();
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
async fn test_bidimap_range_stream() {
    let (path, _guard) = setup("bidimap_range_stream");
    let map = BidiBuilder::new(&path).build().await.unwrap();
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
async fn test_bidimap_disable_scan_on_startup() {
    let (path, _guard) = setup("bidimap_disable_scan_on_startup");
    {
        let map = BidiBuilder::<String, String>::new(&path)
            .build()
            .await
            .unwrap();
        map.set("a".to_string(), "1".to_string(), 1).await.unwrap();
        map.set("b".to_string(), "2".to_string(), 2).await.unwrap();
        map.flush().await.unwrap();
    }

    let map = BidiBuilder::<String, String>::new(&path)
        .disable_scan_on_startup()
        .build()
        .await
        .unwrap();

    assert_eq!(map.len(), 0);

    let (v, _) = map.get("a").await.unwrap();
    assert_eq!(v, "1");
}

#[tokio::test]
async fn test_bidimap_concurrent_access() {
    let (path, _guard) = setup("bidimap_concurrent_access");
    let map = BidiBuilder::new(&path).build().await.unwrap();

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
            let (read_v, read_ts) = map.get(k.as_str()).await.unwrap();
            assert_eq!(read_v, v);
            assert_eq!(read_ts, ts);

            let (read_k, read_ts_inv) = map.get_inverse(v.as_str()).await.unwrap();
            assert_eq!(read_k, k);
            assert_eq!(read_ts_inv, ts);
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
            if i % 2 == 0 {
                let deleted_v = map.delete(k.as_str()).await.unwrap();
                assert_eq!(deleted_v, v);
            } else {
                let deleted_k = map.delete_inverse(v.as_str()).await.unwrap();
                assert_eq!(deleted_k, k);
            }
        });
    }
    while let Some(res) = set.join_next().await {
        res.unwrap();
    }

    assert_eq!(map.len(), 0);
}

#[tokio::test]
async fn test_unidimap_crud_and_len() {
    let (path, _guard) = setup("unidimap_crud_and_len");
    let map = BidiBuilder::new(&path).build().await.unwrap();
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
    assert!(map.get_inverse("one").await.is_err());

    let removed = map.delete("alpha").await.unwrap();
    assert_eq!(removed, "uno");
    assert_eq!(map.len(), 0);
    assert!(map.get("alpha").await.is_err());
}

#[tokio::test]
async fn test_unidimap_delete_inverse() {
    let (path, _guard) = setup("unidimap_delete_inverse");
    let map = BidiBuilder::new(&path).build().await.unwrap();
    map.set("a".to_string(), "1".to_string(), 1).await.unwrap();
    assert_eq!(map.len(), 1);

    let removed_key = map.delete_inverse("1").await.unwrap();
    assert_eq!(removed_key, "a");
    assert_eq!(map.len(), 0);
    assert!(map.get("a").await.is_err());
    assert!(map.get_inverse("1").await.is_err());
}

#[tokio::test]
async fn test_unidimap_range_callback() {
    let (path, _guard) = setup("unidimap_range_callback");
    let map = BidiBuilder::new(&path).build().await.unwrap();
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
async fn test_unidimap_range_stream() {
    let (path, _guard) = setup("unidimap_range_stream");
    let map = BidiBuilder::new(&path).build().await.unwrap();
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
async fn test_unidimap_disable_scan_on_startup() {
    let (path, _guard) = setup("unidimap_disable_scan_on_startup");
    {
        let map = BidiBuilder::<String, String>::new(&path)
            .build()
            .await
            .unwrap();
        map.set("a".to_string(), "1".to_string(), 1).await.unwrap();
        map.set("b".to_string(), "2".to_string(), 2).await.unwrap();
        map.flush().await.unwrap();
    }

    let map = BidiBuilder::<String, String>::new(&path)
        .disable_scan_on_startup()
        .build()
        .await
        .unwrap();

    assert_eq!(map.len(), 0);

    let (v, _) = map.get("a").await.unwrap();
    assert_eq!(v, "1");
}

#[tokio::test]
async fn test_unidimap_concurrent_access() {
    let (path, _guard) = setup("unidimap_concurrent_access");
    let map = BidiBuilder::new(&path).build().await.unwrap();

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
            let (read_v, read_ts) = map.get(k.as_str()).await.unwrap();
            assert_eq!(read_v, v);
            assert_eq!(read_ts, ts);

            let (read_k, read_ts_inv) = map.get_inverse(v.as_str()).await.unwrap();
            assert_eq!(read_k, k);
            assert_eq!(read_ts_inv, ts);
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
            if i % 2 == 0 {
                let deleted_v = map.delete(k.as_str()).await.unwrap();
                assert_eq!(deleted_v, v);
            } else {
                let deleted_k = map.delete_inverse(v.as_str()).await.unwrap();
                assert_eq!(deleted_k, k);
            }
        });
    }
    while let Some(res) = set.join_next().await {
        res.unwrap();
    }

    assert_eq!(map.len(), 0);
}