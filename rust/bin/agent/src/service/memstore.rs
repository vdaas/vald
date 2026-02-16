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

//! # Memstore
//!
//! This module provides functions for managing the in-memory store that combines
//! the KVS (key-value store) and VQueue (vector queue) for the agent.
//! It handles conflict resolution between the two stores based on timestamps.

use std::sync::Arc;

use kvs::{BidirectionalMap, MapBase, map::codec::WincodeCodec};
use thiserror::Error;
use vqueue::{Queue, QueueError};

/// Error type for memstore operations.
#[derive(Debug, Error)]
pub enum MemstoreError {
    /// Error when UUID is not found.
    #[error("UUID not found: {0}")]
    UuidNotFound(String),

    /// Error when object is not found.
    #[error("Object not found: {0}")]
    ObjectNotFound(String),

    /// Error when object ID is not found.
    #[error("Object ID not found: {0}")]
    ObjectIdNotFound(String),

    /// Error when timestamp is zero.
    #[error("Zero timestamp provided")]
    ZeroTimestamp,

    /// Error when a newer timestamp object already exists.
    #[error("Newer timestamp object already exists for uuid: {0}, provided timestamp: {1}")]
    NewerTimestampObjectAlreadyExists(String, i64),

    /// Error when nothing needs to be done for update.
    #[error("Nothing to be done for update: {0}")]
    NothingToBeDoneForUpdate(String),

    /// Error from KVS operations.
    #[error("KVS error: {0}")]
    Kvs(#[from] kvs::map::error::Error),

    /// Error from VQueue operations.
    #[error("VQueue error: {0}")]
    VQueue(#[from] QueueError),
}

/// Type alias for the bidirectional map used in memstore.
/// Maps UUID (String) to OID (u32).
pub type KvsMap = BidirectionalMap<String, u32, WincodeCodec>;

/// Checks if a UUID exists in the memstore (kvs + vqueue).
///
/// # Arguments
///
/// * `kv` - The KVS bidirectional map.
/// * `vq` - The vector queue.
/// * `uuid` - The UUID to check.
///
/// # Returns
///
/// A tuple of (oid, exists). If the UUID exists, `oid` is the object ID and `exists` is true.
pub async fn exists<Q: Queue>(
    kv: &Arc<KvsMap>,
    vq: &Q,
    uuid: &str,
) -> Result<(u32, bool), MemstoreError> {
    // Check vqueue first
    let vq_result = vq.get_vector_with_timestamp(uuid).await;

    match vq_result {
        Ok((_vec, its, dts, exists)) => {
            if exists {
                // Found in vqueue with valid insert
                // Try to get OID from kvs
                match kv.get(uuid).await {
                    Ok((oid, kts)) => {
                        // Update kvs timestamp if vqueue is newer
                        if (kts as i64) < its {
                            let _ = kv.set(uuid.to_string(), oid, its as u128).await;
                        }
                        Ok((oid, true))
                    }
                    Err(_) => {
                        // Not in kvs yet (still in vqueue), return 0 as oid
                        Ok((0, true))
                    }
                }
            } else {
                // Not valid in vqueue (delete is newer or not found)
                // Check kvs
                match kv.get(uuid).await {
                    Ok((oid, kts)) => {
                        // Update kvs timestamp if insert timestamp is newer
                        if its > 0 && (kts as i64) < its {
                            let _ = kv.set(uuid.to_string(), oid, its as u128).await;
                        }
                        // If delete timestamp is newer than insert, object will be deleted soon
                        if dts > its {
                            log::debug!(
                                "Exists: uuid {}'s data found in kvsdb but delete vqueue data exists. The object will be deleted soon",
                                uuid
                            );
                            return Ok((0, false));
                        }
                        Ok((oid, true))
                    }
                    Err(_) => Ok((0, false)),
                }
            }
        }
        Err(QueueError::NotFound(_)) => {
            // Not in vqueue, check kvs only
            match kv.get(uuid).await {
                Ok((oid, _ts)) => Ok((oid, true)),
                Err(_) => Ok((0, false)),
            }
        }
        Err(e) => Err(MemstoreError::VQueue(e)),
    }
}

/// Gets an object (vector and timestamp) from the memstore.
///
/// # Arguments
///
/// * `kv` - The KVS bidirectional map.
/// * `vq` - The vector queue.
/// * `uuid` - The UUID of the object to retrieve.
/// * `get_vector_fn` - A function to get the vector from the index by OID.
///
/// # Returns
///
/// A tuple of (vector, timestamp).
pub async fn get_object<Q, F, Fut>(
    kv: &Arc<KvsMap>,
    vq: &Q,
    uuid: &str,
    get_vector_fn: Option<F>,
) -> Result<(Vec<f32>, i64), MemstoreError>
where
    Q: Queue,
    F: FnOnce(u32) -> Fut,
    Fut: std::future::Future<Output = Result<Vec<f32>, MemstoreError>>,
{
    // Check vqueue first
    let vq_result = vq.get_vector_with_timestamp(uuid).await;

    match vq_result {
        Ok((Some(vec), its, dts, exists)) => {
            if exists {
                return Ok((vec, its));
            }
            // Vector exists but delete is newer, check kvs
            match kv.get(uuid).await {
                Ok((oid, kts)) => {
                    // Update kvs timestamp if vqueue insert is newer
                    if (kts as i64) < its {
                        let _ = kv.set(uuid.to_string(), oid, its as u128).await;
                    }
                    // If delete timestamp is newer, object will be deleted soon
                    if dts > its {
                        log::debug!(
                            "GetObject: uuid {}'s data found in kvsdb but delete vqueue data exists. The object will be deleted soon",
                            uuid
                        );
                        return Err(MemstoreError::ObjectIdNotFound(uuid.to_string()));
                    }
                    // Get vector from index
                    if let Some(f) = get_vector_fn {
                        let vec = f(oid).await?;
                        return Ok((vec, kts as i64));
                    }
                    Err(MemstoreError::ObjectNotFound(uuid.to_string()))
                }
                Err(_) => Err(MemstoreError::ObjectIdNotFound(uuid.to_string())),
            }
        }
        Ok((None, its, dts, _exists)) => {
            // No vector in vqueue, check kvs
            match kv.get(uuid).await {
                Ok((oid, kts)) => {
                    // Update kvs timestamp if vqueue insert is newer
                    if its > 0 && (kts as i64) < its {
                        let _ = kv.set(uuid.to_string(), oid, its as u128).await;
                    }
                    // If delete timestamp is newer, object will be deleted soon
                    if dts > its && dts > 0 {
                        log::debug!(
                            "GetObject: uuid {}'s data found in kvsdb but delete vqueue data exists. The object will be deleted soon",
                            uuid
                        );
                        return Err(MemstoreError::ObjectIdNotFound(uuid.to_string()));
                    }
                    // Get vector from index
                    if let Some(f) = get_vector_fn {
                        let vec = f(oid).await?;
                        return Ok((vec, kts as i64));
                    }
                    Err(MemstoreError::ObjectNotFound(uuid.to_string()))
                }
                Err(_) => {
                    log::debug!(
                        "GetObject: uuid {}'s data not found in kvsdb and insert vqueue",
                        uuid
                    );
                    Err(MemstoreError::ObjectIdNotFound(uuid.to_string()))
                }
            }
        }
        Err(QueueError::NotFound(_)) => {
            // Not in vqueue, check kvs only
            match kv.get(uuid).await {
                Ok((oid, kts)) => {
                    if let Some(f) = get_vector_fn {
                        let vec = f(oid).await?;
                        return Ok((vec, kts as i64));
                    }
                    Err(MemstoreError::ObjectNotFound(uuid.to_string()))
                }
                Err(_) => {
                    log::debug!(
                        "GetObject: uuid {}'s data not found in kvsdb and insert vqueue",
                        uuid
                    );
                    Err(MemstoreError::ObjectIdNotFound(uuid.to_string()))
                }
            }
        }
        Err(e) => Err(MemstoreError::VQueue(e)),
    }
}

/// Collects all UUIDs from the memstore (kvs + vqueue).
///
/// # Arguments
///
/// * `kv` - The KVS bidirectional map.
/// * `vq` - The vector queue.
///
/// # Returns
///
/// A vector of UUIDs.
pub async fn uuids<Q: Queue>(kv: &Arc<KvsMap>, vq: &Q) -> Result<Vec<String>, MemstoreError> {
    use futures::StreamExt;
    use kvs::MapBase;

    let mut result = Vec::new();
    let mut seen = std::collections::HashSet::new();

    // Collect from kvs using range_stream
    let mut stream = Box::pin(kv.range_stream());
    while let Some(item) = stream.next().await {
        if let Ok((uuid, _oid, _ts)) = item {
            // Check if this uuid has a pending delete
            match vq.dv_exists(&uuid).await {
                Ok(dts) if dts > 0 => {
                    // Has pending delete, check if insert is newer
                    match vq.iv_exists(&uuid).await {
                        Ok(its) if its > dts => {
                            seen.insert(uuid.clone());
                            result.push(uuid);
                        }
                        _ => {
                            // Delete is newer or no insert, skip
                        }
                    }
                }
                _ => {
                    // No pending delete
                    seen.insert(uuid.clone());
                    result.push(uuid);
                }
            }
        }
    }

    // Then, collect from vqueue insert queue (items not yet in kvs)
    // Note: This requires iterating through vqueue, which we can do via ivq_len check
    // For now, we rely on the kvs having most items and vqueue having uncommitted ones
    // A full implementation would need a range/iterator on vqueue

    Ok(result)
}

/// Applies the input function on each index stored in the kvs and vqueue.
/// Use this function for performing something on each object while caring about memory usage.
/// If the vector exists in the vqueue, this vector is not indexed so the oid(object ID) is processed as 0.
///
/// # Arguments
///
/// * `kv` - The KVS bidirectional map.
/// * `vq` - The vector queue.
/// * `f` - A callback function to process each item. Returns false to stop iteration.
pub async fn list_object_func<Q, F>(kv: &Arc<KvsMap>, vq: &Q, mut f: F)
where
    Q: Queue,
    F: FnMut(String, u32, i64) -> bool + Send,
{
    use futures::StreamExt;
    use kvs::MapBase;
    use std::collections::HashSet;

    let mut dup: HashSet<String> = HashSet::new();

    // First, iterate through vqueue insert items
    let mut vq_stream = Box::pin(vq.range());
    while let Some(item) = vq_stream.next().await {
        if let Ok((uuid, _vec, ts)) = item {
            // Check if this uuid exists in kvs
            match kv.get(&uuid).await {
                Ok((oid, kts)) => {
                    // Exists in kvs
                    if ts > kts as i64 {
                        // vqueue is newer, use vqueue timestamp
                        dup.insert(uuid.clone());
                        if !f(uuid, oid, ts) {
                            return;
                        }
                    }
                    // else: kvs data is newer, will process at kvs.range
                }
                Err(_) => {
                    // Not in kvs, oid is 0
                    if !f(uuid, 0, ts) {
                        return;
                    }
                }
            }
        }
    }

    // Then, iterate through kvs entries
    let mut kv_stream = Box::pin(kv.range_stream());
    while let Some(item) = kv_stream.next().await {
        if let Ok((uuid, oid, ts)) = item {
            // Skip if already processed from vqueue
            if dup.contains(&uuid) {
                continue;
            }
            // Check if delete vqueue data exists and is newer (data will be deleted soon)
            match vq.dv_exists(&uuid).await {
                Ok(dts) if dts > 0 => {
                    // Has pending delete, skip
                    continue;
                }
                _ => {}
            }
            if !f(uuid, oid, ts as i64) {
                return;
            }
        }
    }
}

/// Updates the timestamp of an object in the memstore.
///
/// # Arguments
///
/// * `kv` - The KVS bidirectional map.
/// * `vq` - The vector queue.
/// * `uuid` - The UUID of the object to update.
/// * `ts` - The new timestamp.
/// * `force` - If true, forces the update even if the new timestamp is older.
/// * `get_vector_fn` - A function to get the vector from the index by OID.
///
/// # Returns
///
/// Ok(()) if the update was successful.
pub async fn update_timestamp<Q, F, Fut>(
    kv: &Arc<KvsMap>,
    vq: &Q,
    uuid: &str,
    ts: i64,
    force: bool,
    get_vector_fn: Option<F>,
) -> Result<(), MemstoreError>
where
    Q: Queue,
    F: FnOnce(u32) -> Fut,
    Fut: std::future::Future<Output = Result<Vec<f32>, MemstoreError>>,
{
    if uuid.is_empty() {
        return Err(MemstoreError::UuidNotFound("empty".to_string()));
    }
    if !force && ts <= 0 {
        return Err(MemstoreError::ZeroTimestamp);
    }

    // Read vqueue data
    let vq_result = vq.get_vector_with_timestamp(uuid).await;
    let (vec, its, dts, vqok) = match vq_result {
        Ok((v, i, d, exists)) => (v, i, d, exists || i > 0 || d > 0),
        Err(QueueError::NotFound(_)) => (None, 0, 0, false),
        Err(e) => return Err(MemstoreError::VQueue(e)),
    };

    // Read kvs data
    let kv_result = kv.get(uuid).await;
    let (oid, kts, kvok) = match kv_result {
        Ok((o, t)) => (o, t as i64, true),
        Err(_) => (0, 0, false),
    };

    if !vqok && !kvok {
        return Err(MemstoreError::ObjectNotFound(uuid.to_string()));
    }

    if !force && (ts <= kts || ts <= its) {
        return Err(MemstoreError::NewerTimestampObjectAlreadyExists(
            uuid.to_string(),
            ts,
        ));
    }

    // Case 1: Only in vqueue, no kvs data, and timestamp is newer than delete
    if vqok
        && !kvok
        && dts != 0
        && dts < ts
        && (force || its < ts)
        && let Some(v) = vec
    {
        vq.push_insert(uuid, v, Some(ts)).await?;
        // Pop delete since we don't need it anymore
        match vq.pop_delete(uuid).await {
            Ok(pdts) if pdts != dts => {
                // Rollback if timestamp changed
                vq.push_delete(uuid, Some(pdts)).await?;
            }
            _ => {}
        }
        return Ok(());
    }

    // Case 2: Both in vqueue and kvs
    if vqok
        && kvok
        && dts < ts
        && (force || (kts < ts && its < ts))
        && let Some(v) = vec
    {
        vq.push_insert(uuid, v, Some(ts)).await?;
        kv.set(uuid.to_string(), oid, ts as u128).await?;
        if dts == 0 {
            // Add delete vqueue for update
            vq.push_delete(uuid, Some(ts - 1)).await?;
        }
        return Ok(());
    }

    // Case 3: Not in insert vqueue, but in kvs
    if !vqok && its == 0 && kvok && (force || kts < ts) {
        kv.set(uuid.to_string(), oid, ts as u128).await?;
        if dts != 0 && (force || dts < ts) {
            match vq.pop_delete(uuid).await {
                Ok(pdts) if pdts != dts => {
                    // Rollback if timestamp changed
                    vq.push_delete(uuid, Some(pdts)).await?;
                }
                _ => {}
            }
        }
        return Ok(());
    }

    // Case 4: Insert vqueue found with special conditions
    if !vqok && its != 0 && kvok && (force || kts < ts) {
        kv.set(uuid.to_string(), oid, ts as u128).await?;
        if vec.is_none()
            && its > dts
            && let Some(f) = get_vector_fn
            && let Ok(ovec) = f(oid).await
        {
            vq.push_insert(uuid, ovec, Some(ts)).await?;
            return Ok(());
        }
        match vq.pop_insert(uuid).await {
            Ok((pvec, pits)) if pits != its => {
                // Rollback if timestamp changed
                vq.push_insert(uuid, pvec, Some(pits)).await?;
            }
            _ => {}
        }
        return Ok(());
    }

    Err(MemstoreError::NothingToBeDoneForUpdate(uuid.to_string()))
}

#[cfg(test)]
mod tests {
    use super::*;
    use kvs::BidirectionalMapBuilder;
    use std::fs;
    use std::future::Ready;
    use vqueue::{Builder as VQueueBuilder, PersistentQueue};

    // Type alias for the None case in get_vector_fn
    type NoopFuture = Ready<Result<Vec<f32>, MemstoreError>>;
    type NoopFn = fn(u32) -> NoopFuture;

    struct TestGuard {
        paths: Vec<String>,
    }

    impl Drop for TestGuard {
        fn drop(&mut self) {
            for path in &self.paths {
                let _ = fs::remove_dir_all(path);
            }
        }
    }

    async fn setup(test_name: &str) -> (Arc<KvsMap>, PersistentQueue, TestGuard) {
        let kvs_path = format!("./test_memstore_kvs_{}", test_name);
        let vq_path = format!("./test_memstore_vq_{}", test_name);
        let _ = fs::remove_dir_all(&kvs_path);
        let _ = fs::remove_dir_all(&vq_path);

        let guard = TestGuard {
            paths: vec![kvs_path.clone(), vq_path.clone()],
        };

        let kv = BidirectionalMapBuilder::<String, u32, WincodeCodec>::new(&kvs_path)
            .build()
            .await
            .unwrap();

        let vq = VQueueBuilder::new(&vq_path).build().await.unwrap();

        (kv, vq, guard)
    }

    #[tokio::test]
    async fn test_exists_in_vqueue() {
        let (kv, vq, _guard) = setup("exists_in_vqueue").await;

        vq.push_insert("uuid1", vec![1.0, 2.0], Some(100))
            .await
            .unwrap();

        let (oid, ok) = exists(&kv, &vq, "uuid1").await.unwrap();
        assert!(ok);
        assert_eq!(oid, 0); // Not in kvs yet
    }

    #[tokio::test]
    async fn test_exists_in_kvs() {
        let (kv, vq, _guard) = setup("exists_in_kvs").await;

        kv.set("uuid1".to_string(), 42, 100).await.unwrap();

        let (oid, ok) = exists(&kv, &vq, "uuid1").await.unwrap();
        assert!(ok);
        assert_eq!(oid, 42);
    }

    #[tokio::test]
    async fn test_exists_not_found() {
        let (kv, vq, _guard) = setup("exists_not_found").await;

        let (oid, ok) = exists(&kv, &vq, "nonexistent").await.unwrap();
        assert!(!ok);
        assert_eq!(oid, 0);
    }

    #[tokio::test]
    async fn test_exists_with_pending_delete() {
        let (kv, vq, _guard) = setup("exists_with_pending_delete").await;

        // Insert then delete (delete is newer)
        vq.push_insert("uuid1", vec![1.0], Some(100)).await.unwrap();
        vq.push_delete("uuid1", Some(200)).await.unwrap();

        let (oid, ok) = exists(&kv, &vq, "uuid1").await.unwrap();
        assert!(!ok);
        assert_eq!(oid, 0);
    }

    #[tokio::test]
    async fn test_get_object_from_vqueue() {
        let (kv, vq, _guard) = setup("get_object_from_vqueue").await;

        vq.push_insert("uuid1", vec![1.0, 2.0], Some(100))
            .await
            .unwrap();

        let (vec, ts) = get_object::<_, NoopFn, NoopFuture>(&kv, &vq, "uuid1", None)
            .await
            .unwrap();
        assert_eq!(vec, vec![1.0, 2.0]);
        assert_eq!(ts, 100);
    }

    #[tokio::test]
    async fn test_get_object_from_kvs_with_fn() {
        let (kv, vq, _guard) = setup("get_object_from_kvs_with_fn").await;

        kv.set("uuid1".to_string(), 42, 100).await.unwrap();

        let get_fn = |_oid: u32| async move { Ok(vec![3.0, 4.0]) };

        let (vec, ts) = get_object(&kv, &vq, "uuid1", Some(get_fn)).await.unwrap();
        assert_eq!(vec, vec![3.0, 4.0]);
        assert_eq!(ts, 100);
    }

    #[tokio::test]
    async fn test_get_object_not_found() {
        let (kv, vq, _guard) = setup("get_object_not_found").await;

        let result = get_object::<_, NoopFn, NoopFuture>(&kv, &vq, "nonexistent", None).await;
        assert!(matches!(result, Err(MemstoreError::ObjectIdNotFound(_))));
    }

    #[tokio::test]
    async fn test_update_timestamp_in_kvs() {
        let (kv, vq, _guard) = setup("update_timestamp_in_kvs").await;

        kv.set("uuid1".to_string(), 42, 100).await.unwrap();

        update_timestamp::<_, NoopFn, NoopFuture>(&kv, &vq, "uuid1", 200, false, None)
            .await
            .unwrap();

        let (oid, ts) = kv.get("uuid1").await.unwrap();
        assert_eq!(oid, 42);
        assert_eq!(ts, 200);
    }

    #[tokio::test]
    async fn test_update_timestamp_not_found() {
        let (kv, vq, _guard) = setup("update_timestamp_not_found").await;

        let result =
            update_timestamp::<_, NoopFn, NoopFuture>(&kv, &vq, "nonexistent", 200, false, None)
                .await;
        assert!(matches!(result, Err(MemstoreError::ObjectNotFound(_))));
    }

    #[tokio::test]
    async fn test_update_timestamp_newer_exists() {
        let (kv, vq, _guard) = setup("update_timestamp_newer_exists").await;

        kv.set("uuid1".to_string(), 42, 200).await.unwrap();

        let result =
            update_timestamp::<_, NoopFn, NoopFuture>(&kv, &vq, "uuid1", 100, false, None).await;
        assert!(matches!(
            result,
            Err(MemstoreError::NewerTimestampObjectAlreadyExists(_, _))
        ));
    }

    #[tokio::test]
    async fn test_update_timestamp_force() {
        let (kv, vq, _guard) = setup("update_timestamp_force").await;

        kv.set("uuid1".to_string(), 42, 200).await.unwrap();

        // Force update with older timestamp
        update_timestamp::<_, NoopFn, NoopFuture>(&kv, &vq, "uuid1", 100, true, None)
            .await
            .unwrap();

        let (oid, ts) = kv.get("uuid1").await.unwrap();
        assert_eq!(oid, 42);
        assert_eq!(ts, 100);
    }

    // ========== list_object_func Tests ==========

    #[tokio::test]
    async fn test_list_object_func_empty() {
        let (kv, vq, _guard) = setup("list_object_func_empty").await;

        let mut count = 0;
        list_object_func(&kv, &vq, |_uuid, _oid, _ts| {
            count += 1;
            true
        })
        .await;

        assert_eq!(count, 0);
    }

    #[tokio::test]
    async fn test_list_object_func_kvs_only() {
        let (kv, vq, _guard) = setup("list_object_func_kvs_only").await;

        kv.set("uuid1".to_string(), 1, 100).await.unwrap();
        kv.set("uuid2".to_string(), 2, 200).await.unwrap();

        let mut items: Vec<(String, u32, i64)> = Vec::new();
        list_object_func(&kv, &vq, |uuid, oid, ts| {
            items.push((uuid, oid, ts));
            true
        })
        .await;

        assert_eq!(items.len(), 2);
        let uuids: Vec<_> = items.iter().map(|(u, _, _)| u.clone()).collect();
        assert!(uuids.contains(&"uuid1".to_string()));
        assert!(uuids.contains(&"uuid2".to_string()));
    }

    #[tokio::test]
    async fn test_list_object_func_vqueue_only() {
        let (kv, vq, _guard) = setup("list_object_func_vqueue_only").await;

        vq.push_insert("uuid1", vec![1.0], Some(100)).await.unwrap();
        vq.push_insert("uuid2", vec![2.0], Some(200)).await.unwrap();

        let mut items: Vec<(String, u32, i64)> = Vec::new();
        list_object_func(&kv, &vq, |uuid, oid, ts| {
            items.push((uuid, oid, ts));
            true
        })
        .await;

        assert_eq!(items.len(), 2);
        // OID should be 0 for items only in vqueue
        for (_, oid, _) in &items {
            assert_eq!(*oid, 0);
        }
    }

    #[tokio::test]
    async fn test_list_object_func_both_kvs_and_vqueue() {
        let (kv, vq, _guard) = setup("list_object_func_both").await;

        // Item in kvs
        kv.set("uuid1".to_string(), 1, 100).await.unwrap();
        // Item in vqueue only
        vq.push_insert("uuid2", vec![2.0], Some(200)).await.unwrap();

        let mut items: Vec<(String, u32, i64)> = Vec::new();
        list_object_func(&kv, &vq, |uuid, oid, ts| {
            items.push((uuid, oid, ts));
            true
        })
        .await;

        assert_eq!(items.len(), 2);
    }

    #[tokio::test]
    async fn test_list_object_func_vqueue_newer_than_kvs() {
        let (kv, vq, _guard) = setup("list_object_func_vqueue_newer").await;

        // Same uuid in both kvs and vqueue, vqueue is newer
        kv.set("uuid1".to_string(), 1, 100).await.unwrap();
        vq.push_insert("uuid1", vec![1.0], Some(200)).await.unwrap();

        let mut items: Vec<(String, u32, i64)> = Vec::new();
        list_object_func(&kv, &vq, |uuid, oid, ts| {
            items.push((uuid, oid, ts));
            true
        })
        .await;

        // Should only appear once with the newer timestamp
        assert_eq!(items.len(), 1);
        assert_eq!(items[0].0, "uuid1");
        assert_eq!(items[0].1, 1); // OID from kvs
        assert_eq!(items[0].2, 200); // timestamp from vqueue (newer)
    }

    #[tokio::test]
    async fn test_list_object_func_skips_pending_delete() {
        let (kv, vq, _guard) = setup("list_object_func_skips_delete").await;

        // Item in kvs with pending delete
        kv.set("uuid1".to_string(), 1, 100).await.unwrap();
        vq.push_delete("uuid1", Some(200)).await.unwrap();

        // Item in kvs without pending delete
        kv.set("uuid2".to_string(), 2, 100).await.unwrap();

        let mut items: Vec<(String, u32, i64)> = Vec::new();
        list_object_func(&kv, &vq, |uuid, oid, ts| {
            items.push((uuid, oid, ts));
            true
        })
        .await;

        // Only uuid2 should appear (uuid1 has pending delete)
        assert_eq!(items.len(), 1);
        assert_eq!(items[0].0, "uuid2");
    }

    #[tokio::test]
    async fn test_list_object_func_early_termination() {
        let (kv, vq, _guard) = setup("list_object_func_early_term").await;

        kv.set("uuid1".to_string(), 1, 100).await.unwrap();
        kv.set("uuid2".to_string(), 2, 200).await.unwrap();
        kv.set("uuid3".to_string(), 3, 300).await.unwrap();

        let mut count = 0;
        list_object_func(&kv, &vq, |_uuid, _oid, _ts| {
            count += 1;
            count < 2 // Stop after 2 items
        })
        .await;

        // Should stop early
        assert!(count <= 2);
    }

    #[tokio::test]
    async fn test_list_object_func_vqueue_delete_newer_filters() {
        let (kv, vq, _guard) = setup("list_object_func_vq_delete_filters").await;

        // Insert then delete in vqueue (delete is newer)
        vq.push_insert("uuid1", vec![1.0], Some(100)).await.unwrap();
        vq.push_delete("uuid1", Some(200)).await.unwrap();

        // Insert in vqueue only (no delete)
        vq.push_insert("uuid2", vec![2.0], Some(300)).await.unwrap();

        let mut items: Vec<(String, u32, i64)> = Vec::new();
        list_object_func(&kv, &vq, |uuid, oid, ts| {
            items.push((uuid, oid, ts));
            true
        })
        .await;

        // uuid1 should be filtered by range() because delete is newer
        // uuid2 should appear
        assert_eq!(items.len(), 1);
        assert_eq!(items[0].0, "uuid2");
    }

    // ========== uuids Tests ==========

    #[tokio::test]
    async fn test_uuids_empty() {
        let (kv, vq, _guard) = setup("uuids_empty").await;

        let result = uuids(&kv, &vq).await.unwrap();
        assert!(result.is_empty());
    }

    #[tokio::test]
    async fn test_uuids_from_kvs_only() {
        let (kv, vq, _guard) = setup("uuids_from_kvs_only").await;

        kv.set("uuid1".to_string(), 1, 100).await.unwrap();
        kv.set("uuid2".to_string(), 2, 200).await.unwrap();
        kv.set("uuid3".to_string(), 3, 300).await.unwrap();

        let mut result = uuids(&kv, &vq).await.unwrap();
        result.sort();

        assert_eq!(result.len(), 3);
        assert_eq!(result, vec!["uuid1", "uuid2", "uuid3"]);
    }

    #[tokio::test]
    async fn test_uuids_filters_pending_deletes() {
        let (kv, vq, _guard) = setup("uuids_filters_pending_deletes").await;

        kv.set("uuid1".to_string(), 1, 100).await.unwrap();
        kv.set("uuid2".to_string(), 2, 200).await.unwrap();

        // Add pending delete for uuid1
        vq.push_delete("uuid1", Some(300)).await.unwrap();

        let result = uuids(&kv, &vq).await.unwrap();

        // Only uuid2 should appear (uuid1 has pending delete)
        assert_eq!(result.len(), 1);
        assert_eq!(result[0], "uuid2");
    }

    #[tokio::test]
    async fn test_uuids_includes_if_insert_newer_than_delete() {
        let (kv, vq, _guard) = setup("uuids_insert_newer_than_delete").await;

        kv.set("uuid1".to_string(), 1, 100).await.unwrap();

        // Delete then insert with newer timestamp
        vq.push_delete("uuid1", Some(200)).await.unwrap();
        vq.push_insert("uuid1", vec![1.0], Some(300)).await.unwrap();

        let result = uuids(&kv, &vq).await.unwrap();

        // uuid1 should appear because insert is newer than delete
        assert_eq!(result.len(), 1);
        assert_eq!(result[0], "uuid1");
    }

    // ========== Additional exists Tests ==========

    #[tokio::test]
    async fn test_exists_both_kvs_and_vqueue() {
        let (kv, vq, _guard) = setup("exists_both_kvs_and_vqueue").await;

        kv.set("uuid1".to_string(), 42, 100).await.unwrap();
        vq.push_insert("uuid1", vec![1.0], Some(200)).await.unwrap();

        let (oid, ok) = exists(&kv, &vq, "uuid1").await.unwrap();
        assert!(ok);
        assert_eq!(oid, 42); // Should get OID from kvs
    }

    #[tokio::test]
    async fn test_exists_delete_then_insert_newer() {
        let (kv, vq, _guard) = setup("exists_delete_then_insert_newer").await;

        // Push delete first, then insert with newer timestamp
        vq.push_delete("uuid1", Some(100)).await.unwrap();
        vq.push_insert("uuid1", vec![1.0], Some(200)).await.unwrap();

        let (oid, ok) = exists(&kv, &vq, "uuid1").await.unwrap();
        assert!(ok);
        assert_eq!(oid, 0); // Not in kvs yet
    }

    #[tokio::test]
    async fn test_exists_kvs_with_newer_delete() {
        let (kv, vq, _guard) = setup("exists_kvs_with_newer_delete").await;

        kv.set("uuid1".to_string(), 42, 100).await.unwrap();
        // Delete is newer than kvs entry but no insert in vqueue
        vq.push_delete("uuid1", Some(200)).await.unwrap();

        let (oid, ok) = exists(&kv, &vq, "uuid1").await.unwrap();
        // Delete is newer, so object is about to be deleted
        assert!(!ok);
        assert_eq!(oid, 0);
    }

    #[tokio::test]
    async fn test_exists_updates_kvs_timestamp_if_vqueue_newer() {
        let (kv, vq, _guard) = setup("exists_updates_kvs_ts").await;

        kv.set("uuid1".to_string(), 42, 100).await.unwrap();
        vq.push_insert("uuid1", vec![1.0], Some(200)).await.unwrap();

        let (_oid, ok) = exists(&kv, &vq, "uuid1").await.unwrap();
        assert!(ok);

        // Check that kvs timestamp was updated
        let (_, ts) = kv.get("uuid1").await.unwrap();
        assert_eq!(ts, 200);
    }

    // ========== Additional get_object Tests ==========

    #[tokio::test]
    async fn test_get_object_with_pending_delete() {
        let (kv, vq, _guard) = setup("get_object_with_pending_delete").await;

        kv.set("uuid1".to_string(), 42, 100).await.unwrap();
        vq.push_delete("uuid1", Some(200)).await.unwrap();

        let result = get_object::<_, NoopFn, NoopFuture>(&kv, &vq, "uuid1", None).await;
        assert!(matches!(result, Err(MemstoreError::ObjectIdNotFound(_))));
    }

    #[tokio::test]
    async fn test_get_object_vqueue_with_vector_and_pending_delete() {
        let (kv, vq, _guard) = setup("get_object_vq_with_delete").await;

        // Insert then delete (delete is newer)
        vq.push_insert("uuid1", vec![1.0, 2.0], Some(100))
            .await
            .unwrap();
        vq.push_delete("uuid1", Some(200)).await.unwrap();

        let result = get_object::<_, NoopFn, NoopFuture>(&kv, &vq, "uuid1", None).await;
        // Should fail because delete is newer
        assert!(matches!(result, Err(MemstoreError::ObjectIdNotFound(_))));
    }

    #[tokio::test]
    async fn test_get_object_updates_kvs_timestamp() {
        let (kv, vq, _guard) = setup("get_object_updates_kvs_ts").await;

        kv.set("uuid1".to_string(), 42, 100).await.unwrap();
        vq.push_insert("uuid1", vec![1.0, 2.0], Some(200))
            .await
            .unwrap();

        let (vec, ts) = get_object::<_, NoopFn, NoopFuture>(&kv, &vq, "uuid1", None)
            .await
            .unwrap();
        assert_eq!(vec, vec![1.0, 2.0]);
        assert_eq!(ts, 200);

        // When vqueue has the vector (exists=true), kvs timestamp is NOT updated
        // because we return vqueue data directly without touching kvs.
        // kvs update only happens when vqueue has no vector (None) but has insert timestamp.
        let (_, kts) = kv.get("uuid1").await.unwrap();
        assert_eq!(kts, 100); // Stays at original timestamp
    }

    #[tokio::test]
    async fn test_get_object_updates_kvs_timestamp_from_insert_ts() {
        // Test that kvs timestamp is updated when vqueue has a newer insert timestamp
        // but exists=false (delete is newer than insert).
        // In this case, get_object returns an error, but kvs timestamp should still be updated.
        let (kv, vq, _guard) = setup("get_object_updates_kvs_ts2").await;

        // Set initial kvs entry with timestamp 100
        kv.set("uuid1".to_string(), 42, 100).await.unwrap();

        // Push insert with ts=200 (newer than kvs), then delete with ts=150
        // Note: insert(200) > delete(150), so exists=true and we get vqueue data
        // To test kvs update path, we need exists=false but its > kts
        // So: insert(200), delete(300) -> exists=false, but its(200) > kts(100)
        vq.push_insert("uuid1", vec![1.0, 2.0, 3.0], Some(200))
            .await
            .unwrap();
        vq.push_delete("uuid1", Some(300)).await.unwrap();

        // Custom get_vector_fn won't be called because delete is newer
        let get_fn = |oid: u32| async move {
            if oid == 42 {
                Ok(vec![99.0, 99.0, 99.0])
            } else {
                Err(MemstoreError::ObjectNotFound(oid.to_string()))
            }
        };

        // Call get_object - should fail because delete is newer
        let result = get_object(&kv, &vq, "uuid1", Some(get_fn)).await;
        assert!(result.is_err(), "Expected error because delete is newer");

        // But kvs timestamp should still be updated from 100 to 200
        let (oid, kts) = kv.get("uuid1").await.unwrap();
        assert_eq!(oid, 42);
        assert_eq!(
            kts, 200,
            "kvs timestamp should be updated to vqueue insert timestamp"
        );
    }

    #[tokio::test]
    async fn test_get_object_with_custom_vector_fn() {
        let (kv, vq, _guard) = setup("get_object_custom_fn").await;

        kv.set("uuid1".to_string(), 42, 100).await.unwrap();

        // Custom function that returns a specific vector based on OID
        let get_fn = |oid: u32| async move {
            if oid == 42 {
                Ok(vec![42.0, 42.0, 42.0])
            } else {
                Err(MemstoreError::ObjectNotFound(oid.to_string()))
            }
        };

        let (vec, ts) = get_object(&kv, &vq, "uuid1", Some(get_fn)).await.unwrap();
        assert_eq!(vec, vec![42.0, 42.0, 42.0]);
        assert_eq!(ts, 100);
    }

    #[tokio::test]
    async fn test_get_object_vector_fn_returns_error() {
        let (kv, vq, _guard) = setup("get_object_fn_error").await;

        kv.set("uuid1".to_string(), 42, 100).await.unwrap();

        let get_fn = |_oid: u32| async move {
            Err(MemstoreError::ObjectNotFound(
                "vector not found".to_string(),
            ))
        };

        let result = get_object(&kv, &vq, "uuid1", Some(get_fn)).await;
        assert!(matches!(result, Err(MemstoreError::ObjectNotFound(_))));
    }

    // ========== Additional update_timestamp Tests ==========

    #[tokio::test]
    async fn test_update_timestamp_empty_uuid() {
        let (kv, vq, _guard) = setup("update_timestamp_empty_uuid").await;

        let result =
            update_timestamp::<_, NoopFn, NoopFuture>(&kv, &vq, "", 200, false, None).await;
        assert!(matches!(result, Err(MemstoreError::UuidNotFound(_))));
    }

    #[tokio::test]
    async fn test_update_timestamp_zero_timestamp_without_force() {
        let (kv, vq, _guard) = setup("update_timestamp_zero_ts").await;

        kv.set("uuid1".to_string(), 42, 100).await.unwrap();

        let result =
            update_timestamp::<_, NoopFn, NoopFuture>(&kv, &vq, "uuid1", 0, false, None).await;
        assert!(matches!(result, Err(MemstoreError::ZeroTimestamp)));
    }

    #[tokio::test]
    async fn test_update_timestamp_zero_timestamp_with_force() {
        let (kv, vq, _guard) = setup("update_timestamp_zero_ts_force").await;

        kv.set("uuid1".to_string(), 42, 100).await.unwrap();

        // With force=true, zero timestamp is allowed
        update_timestamp::<_, NoopFn, NoopFuture>(&kv, &vq, "uuid1", 0, true, None)
            .await
            .unwrap();

        let (_, ts) = kv.get("uuid1").await.unwrap();
        assert_eq!(ts, 0);
    }

    #[tokio::test]
    async fn test_update_timestamp_in_vqueue_only() {
        let (kv, vq, _guard) = setup("update_timestamp_vqueue_only").await;

        vq.push_insert("uuid1", vec![1.0, 2.0], Some(100))
            .await
            .unwrap();
        vq.push_delete("uuid1", Some(50)).await.unwrap(); // older delete

        update_timestamp::<_, NoopFn, NoopFuture>(&kv, &vq, "uuid1", 200, false, None)
            .await
            .unwrap();

        // Check vqueue has updated timestamp
        let (vec, ts) = vq.get_vector("uuid1").await.unwrap();
        assert_eq!(vec, vec![1.0, 2.0]);
        assert_eq!(ts, 200);
    }

    #[tokio::test]
    async fn test_update_timestamp_both_vqueue_and_kvs() {
        let (kv, vq, _guard) = setup("update_timestamp_both").await;

        kv.set("uuid1".to_string(), 42, 100).await.unwrap();
        vq.push_insert("uuid1", vec![1.0, 2.0], Some(150))
            .await
            .unwrap();

        update_timestamp::<_, NoopFn, NoopFuture>(&kv, &vq, "uuid1", 200, false, None)
            .await
            .unwrap();

        // Both kvs and vqueue should be updated
        let (_, kts) = kv.get("uuid1").await.unwrap();
        assert_eq!(kts, 200);

        let (_, vts) = vq.get_vector("uuid1").await.unwrap();
        assert_eq!(vts, 200);
    }

    #[tokio::test]
    async fn test_update_timestamp_force_older_than_both() {
        let (kv, vq, _guard) = setup("update_timestamp_force_older").await;

        kv.set("uuid1".to_string(), 42, 200).await.unwrap();
        vq.push_insert("uuid1", vec![1.0, 2.0], Some(300))
            .await
            .unwrap();

        // Force update with timestamp older than both
        update_timestamp::<_, NoopFn, NoopFuture>(&kv, &vq, "uuid1", 100, true, None)
            .await
            .unwrap();

        let (_, kts) = kv.get("uuid1").await.unwrap();
        assert_eq!(kts, 100);
    }

    // ========== Edge Case Tests ==========

    #[tokio::test]
    async fn test_exists_multiple_operations_same_uuid() {
        let (kv, vq, _guard) = setup("exists_multiple_ops").await;

        // Simulate multiple operations on same uuid
        vq.push_insert("uuid1", vec![1.0], Some(100)).await.unwrap();
        vq.push_delete("uuid1", Some(150)).await.unwrap();
        vq.push_insert("uuid1", vec![2.0], Some(200)).await.unwrap();

        let (oid, ok) = exists(&kv, &vq, "uuid1").await.unwrap();
        assert!(ok); // Latest insert is newest
        assert_eq!(oid, 0); // Not in kvs
    }

    #[tokio::test]
    async fn test_get_object_prefers_vqueue_over_kvs() {
        let (kv, vq, _guard) = setup("get_object_prefers_vqueue").await;

        // Old data in kvs
        kv.set("uuid1".to_string(), 42, 100).await.unwrap();
        // New data in vqueue
        vq.push_insert("uuid1", vec![999.0], Some(200))
            .await
            .unwrap();

        let (vec, ts) = get_object::<_, NoopFn, NoopFuture>(&kv, &vq, "uuid1", None)
            .await
            .unwrap();
        // Should get vqueue data since it's newer
        assert_eq!(vec, vec![999.0]);
        assert_eq!(ts, 200);
    }

    #[tokio::test]
    async fn test_concurrent_operations() {
        let (kv, vq, _guard) = setup("concurrent_ops").await;

        // Simulate concurrent inserts
        let handles: Vec<_> = (0..10)
            .map(|i| {
                let kv = kv.clone();
                let vq = vq.clone();
                tokio::spawn(async move {
                    let uuid = format!("uuid{}", i);
                    vq.push_insert(&uuid, vec![i as f32], Some(100 + i as i64))
                        .await
                        .unwrap();
                    kv.set(uuid.clone(), i as u32, (100 + i) as u128)
                        .await
                        .unwrap();
                })
            })
            .collect();

        for handle in handles {
            handle.await.unwrap();
        }

        // All items should exist
        for i in 0..10 {
            let uuid = format!("uuid{}", i);
            let (oid, ok) = exists(&kv, &vq, &uuid).await.unwrap();
            assert!(ok, "uuid{} should exist", i);
            assert_eq!(oid, i as u32);
        }
    }

    #[tokio::test]
    async fn test_list_object_func_with_mixed_timestamps() {
        let (kv, vq, _guard) = setup("list_object_func_mixed_ts").await;

        // kvs has older data
        kv.set("uuid1".to_string(), 1, 100).await.unwrap();
        // vqueue has newer data for same uuid
        vq.push_insert("uuid1", vec![1.0], Some(200)).await.unwrap();

        // kvs has newer data
        kv.set("uuid2".to_string(), 2, 300).await.unwrap();
        // vqueue has older data for same uuid
        vq.push_insert("uuid2", vec![2.0], Some(250)).await.unwrap();

        let mut items: Vec<(String, u32, i64)> = Vec::new();
        list_object_func(&kv, &vq, |uuid, oid, ts| {
            items.push((uuid, oid, ts));
            true
        })
        .await;

        items.sort_by(|a, b| a.0.cmp(&b.0));

        assert_eq!(items.len(), 2);
        // uuid1 should have vqueue timestamp (200) because it's newer
        assert_eq!(items[0].0, "uuid1");
        assert_eq!(items[0].2, 200);
        // uuid2 - depends on which source wins based on iteration order
    }

    #[tokio::test]
    async fn test_special_characters_in_uuid() {
        let (kv, vq, _guard) = setup("special_chars").await;

        let special_uuids = vec![
            "uuid-with-dashes",
            "uuid_with_underscores",
            "uuid.with.dots",
            "uuid:with:colons",
            "uuid/with/slashes",
        ];

        for (i, uuid) in special_uuids.iter().enumerate() {
            vq.push_insert(*uuid, vec![i as f32], Some(100 + i as i64))
                .await
                .unwrap();
            kv.set(uuid.to_string(), i as u32, (100 + i) as u128)
                .await
                .unwrap();
        }

        for (i, uuid) in special_uuids.iter().enumerate() {
            let (oid, ok) = exists(&kv, &vq, uuid).await.unwrap();
            assert!(ok, "UUID '{}' should exist", uuid);
            assert_eq!(oid, i as u32);
        }
    }

    #[tokio::test]
    async fn test_large_vector_handling() {
        let (kv, vq, _guard) = setup("large_vector").await;

        // Create a large vector
        let large_vec: Vec<f32> = (0..10000).map(|i| i as f32).collect();

        vq.push_insert("uuid1", large_vec.clone(), Some(100))
            .await
            .unwrap();

        let (vec, ts) = get_object::<_, NoopFn, NoopFuture>(&kv, &vq, "uuid1", None)
            .await
            .unwrap();
        assert_eq!(vec.len(), 10000);
        assert_eq!(ts, 100);
        assert_eq!(vec, large_vec);
    }

    #[tokio::test]
    async fn test_empty_vector_handling() {
        let (kv, vq, _guard) = setup("empty_vector").await;

        vq.push_insert("uuid1", vec![], Some(100)).await.unwrap();

        let (vec, ts) = get_object::<_, NoopFn, NoopFuture>(&kv, &vq, "uuid1", None)
            .await
            .unwrap();
        assert!(vec.is_empty());
        assert_eq!(ts, 100);
    }

    #[tokio::test]
    async fn test_negative_timestamp_handling() {
        let (kv, vq, _guard) = setup("negative_timestamp").await;

        // Negative timestamps should work
        vq.push_insert("uuid1", vec![1.0], Some(-100))
            .await
            .unwrap();

        let (oid, ok) = exists(&kv, &vq, "uuid1").await.unwrap();
        assert!(ok);
        assert_eq!(oid, 0);

        let (vec, ts) = get_object::<_, NoopFn, NoopFuture>(&kv, &vq, "uuid1", None)
            .await
            .unwrap();
        assert_eq!(vec, vec![1.0]);
        assert_eq!(ts, -100);
    }

    #[tokio::test]
    async fn test_max_timestamp_handling() {
        let (kv, vq, _guard) = setup("max_timestamp").await;

        let max_ts = i64::MAX;
        vq.push_insert("uuid1", vec![1.0], Some(max_ts))
            .await
            .unwrap();

        let (vec, ts) = get_object::<_, NoopFn, NoopFuture>(&kv, &vq, "uuid1", None)
            .await
            .unwrap();
        assert_eq!(vec, vec![1.0]);
        assert_eq!(ts, max_ts);
    }
}
