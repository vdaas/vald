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

use dashmap::DashMap;
use std::cmp::Ordering;
use std::sync::atomic::{AtomicU64, Ordering as AtomicOrdering};
use std::sync::Arc;
use std::time::{SystemTime, UNIX_EPOCH};

/// Error type representing possible failures for queue operations.
#[derive(Debug)]
pub enum QueueError {
    /// The provided UUID is invalid or empty.
    InvalidUuid,
    /// The provided timestamp is older than an existing entry and cannot replace it.
    TimestampTooOld,
    /// Internal inconsistency error with a message.
    InternalError(String),
}

/// Trait definition for queue functionality. VQueue implements this interface.
pub trait Queue: Send + Sync {
    /// Inserts the specified UUID and vector into the insert queue.
    /// If timestamp is None, uses the current system time.
    fn push_insert(
        &self,
        uuid: String,
        vector: Arc<Vec<f32>>,
        timestamp: Option<i64>,
    ) -> Result<(), QueueError>;
    /// Inserts a delete request for the specified UUID into the delete queue.
    /// If timestamp is None, uses the current system time.
    fn push_delete(&self, uuid: String, timestamp: Option<i64>) -> Result<(), QueueError>;
    /// Pops and returns the vector and timestamp for the given UUID from the insert queue.
    /// Returns None if no entry exists.
    fn pop_insert(&self, uuid: &str) -> Result<Option<(Arc<Vec<f32>>, i64)>, QueueError>;
    /// Pops and returns the timestamp for the given UUID from the delete queue.
    /// Returns None if no entry exists.
    fn pop_delete(&self, uuid: &str) -> Result<Option<i64>, QueueError>;
    /// Retrieves the vector and timestamp for the given UUID without removing.
    /// Returns None if no valid vector exists.
    fn get_vector(&self, uuid: &str) -> Result<Option<(Arc<Vec<f32>>, i64)>, QueueError>;
    /// Retrieves the vector, insert timestamp, and delete timestamp for the given UUID.
    /// Returns None if no entry exists in either queue.
    fn get_vector_with_timestamp(
        &self,
        uuid: &str,
    ) -> Result<Option<(Arc<Vec<f32>>, i64, i64)>, QueueError>;
    /// Calls the provided closure for each insert entry in the queue.
    /// If the closure returns false, iteration stops.
    fn range<F>(&self, f: F) -> Result<(), QueueError>
    where
        F: FnMut(&str, &Arc<Vec<f32>>, i64) -> bool;
    /// Returns the current size of the insert queue.
    fn ivq_len(&self) -> usize;
    /// Returns the current size of the delete queue.
    fn dvq_len(&self) -> usize;
    /// Retrieves and sorts insert entries with timestamps <= now, calling the closure for each.
    fn range_pop_insert<F>(&self, now: i64, f: F) -> Result<(), QueueError>
    where
        F: FnMut(&str, &Arc<Vec<f32>>, i64) -> bool;
    /// Retrieves and sorts delete entries with timestamps <= now, calling the closure for each.
    fn range_pop_delete<F>(&self, now: i64, f: F) -> Result<(), QueueError>
    where
        F: FnMut(&str) -> bool;
    /// Returns insert timestamp if exists and newer than delete timestamp.
    fn iv_exists(&self, uuid: &str) -> Result<Option<i64>, QueueError>;
    /// Returns delete timestamp if exists and newer than insert timestamp.
    fn dv_exists(&self, uuid: &str) -> Result<Option<i64>, QueueError>;
}

/// Internal structure holding index information.
#[derive(Clone)]
struct Index {
    uuid: String,
    vector: Option<Arc<Vec<f32>>>,
    timestamp: i64,
}

/// VQueue implementation using DashMap and AtomicU64 for concurrent access.
pub struct VQueue {
    il: Arc<DashMap<String, Index>>, // Insert queue
    dl: Arc<DashMap<String, Index>>, // Delete queue
    ic: AtomicU64,
    dc: AtomicU64,
    /// Cancel flag for early termination of iteration
    cancel_flag: Arc<AtomicU64>,
}

impl VQueue {
    /// Creates a new VQueue.
    pub fn new() -> Self {
        VQueue {
            il: Arc::new(DashMap::new()),
            dl: Arc::new(DashMap::new()),
            ic: AtomicU64::new(0),
            dc: AtomicU64::new(0),
            cancel_flag: Arc::new(AtomicU64::new(0)),
        }
    }

    /// Returns the current time in nanoseconds.
    fn now_ns() -> i64 {
        let dur = SystemTime::now()
            .duration_since(UNIX_EPOCH)
            .expect("Time went backwards");
        dur.as_nanos() as i64
    }

    /// Loads the insert entry for the specified key.
    fn load_ivq(&self, uuid: &str) -> Option<Index> {
        self.il.get(uuid).map(|entry| entry.value().clone())
    }

    /// Loads the delete entry for the specified key.
    fn load_dvq(&self, uuid: &str) -> Option<Index> {
        self.dl.get(uuid).map(|entry| entry.value().clone())
    }

    /// Compares timestamps.
    #[inline]
    fn newer(ts1: i64, ts2: i64) -> bool {
        ts1 > ts2
    }

    /// Sets the cancel flag (call to interrupt iteration).
    pub fn cancel(&self) {
        self.cancel_flag.fetch_add(1, AtomicOrdering::SeqCst);
    }
}

impl Queue for VQueue {
    fn push_insert(
        &self,
        uuid: String,
        vector: Arc<Vec<f32>>,
        timestamp: Option<i64>,
    ) -> Result<(), QueueError> {
        if uuid.trim().is_empty() {
            return Err(QueueError::InvalidUuid);
        }
        let ts = timestamp.unwrap_or_else(|| VQueue::now_ns());
        if let Some(didx) = self.load_dvq(&uuid) {
            if VQueue::newer(didx.timestamp, ts) {
                return Err(QueueError::TimestampTooOld);
            }
        }
        let new_idx = Index {
            uuid: uuid.clone(),
            vector: Some(vector.clone()),
            timestamp: ts,
        };
        match self.il.entry(uuid.clone()) {
            dashmap::mapref::entry::Entry::Occupied(mut occ) => {
                let existing = occ.get();
                if VQueue::newer(ts, existing.timestamp) {
                    occ.insert(new_idx);
                } else {
                    return Err(QueueError::TimestampTooOld);
                }
            }
            dashmap::mapref::entry::Entry::Vacant(vac) => {
                vac.insert(new_idx);
                self.ic.fetch_add(1, AtomicOrdering::SeqCst);
            }
        }
        Ok(())
    }

    fn push_delete(&self, uuid: String, timestamp: Option<i64>) -> Result<(), QueueError> {
        if uuid.trim().is_empty() {
            return Err(QueueError::InvalidUuid);
        }
        let ts = timestamp.unwrap_or_else(|| VQueue::now_ns());
        let new_idx = Index {
            uuid: uuid.clone(),
            vector: None,
            timestamp: ts,
        };
        match self.dl.entry(uuid.clone()) {
            dashmap::mapref::entry::Entry::Occupied(mut occ) => {
                let existing = occ.get();
                if VQueue::newer(ts, existing.timestamp) {
                    occ.insert(new_idx);
                } else {
                    return Err(QueueError::TimestampTooOld);
                }
            }
            dashmap::mapref::entry::Entry::Vacant(vac) => {
                vac.insert(new_idx);
                self.dc.fetch_add(1, AtomicOrdering::SeqCst);
            }
        }
        Ok(())
    }

    fn pop_insert(&self, uuid: &str) -> Result<Option<(Arc<Vec<f32>>, i64)>, QueueError> {
        if uuid.trim().is_empty() {
            return Err(QueueError::InvalidUuid);
        }
        if let Some((_, idx)) = self.il.remove(uuid) {
            if idx.timestamp != 0 {
                self.ic.fetch_sub(1, AtomicOrdering::SeqCst);
                if let Some(vec_arc) = idx.vector {
                    return Ok(Some((vec_arc, idx.timestamp)));
                }
            }
        }
        Ok(None)
    }

    fn pop_delete(&self, uuid: &str) -> Result<Option<i64>, QueueError> {
        if uuid.trim().is_empty() {
            return Err(QueueError::InvalidUuid);
        }
        if let Some((_, idx)) = self.dl.remove(uuid) {
            if idx.timestamp != 0 {
                self.dc.fetch_sub(1, AtomicOrdering::SeqCst);
                return Ok(Some(idx.timestamp));
            }
        }
        Ok(None)
    }

    fn get_vector(&self, uuid: &str) -> Result<Option<(Arc<Vec<f32>>, i64)>, QueueError> {
        if uuid.trim().is_empty() {
            return Err(QueueError::InvalidUuid);
        }
        if let Some(idx) = self.load_ivq(uuid) {
            let its = idx.timestamp;
            if let Some(didx) = self.load_dvq(uuid) {
                if VQueue::newer(didx.timestamp, its) {
                    return Ok(None);
                }
            }
            if let Some(vec_arc) = idx.vector {
                return Ok(Some((vec_arc.clone(), its)));
            }
        }
        Ok(None)
    }

    fn get_vector_with_timestamp(
        &self,
        uuid: &str,
    ) -> Result<Option<(Arc<Vec<f32>>, i64, i64)>, QueueError> {
        if uuid.trim().is_empty() {
            return Err(QueueError::InvalidUuid);
        }
        let mut its = 0;
        let mut dts = 0;
        if let Some(idx) = self.load_ivq(uuid) {
            its = idx.timestamp;
        }
        if let Some(didx) = self.load_dvq(uuid) {
            dts = didx.timestamp;
        }
        if its == 0 && dts == 0 {
            return Ok(None);
        }
        if its == 0 {
            return Ok(Some((Arc::new(Vec::new()), 0, dts))); // no insert, only delete
        }
        if dts == 0 || VQueue::newer(its, dts) {
            if let Some(idx) = self.load_ivq(uuid) {
                if let Some(vec_arc) = idx.vector {
                    return Ok(Some((vec_arc.clone(), its, dts)));
                }
            }
            return Ok(None);
        }
        Ok(None)
    }

    fn range<F>(&self, mut f: F) -> Result<(), QueueError>
    where
        F: FnMut(&str, &Arc<Vec<f32>>, i64) -> bool,
    {
        for entry in self.il.iter() {
            if self.cancel_flag.load(AtomicOrdering::SeqCst) != 0 {
                break;
            }
            let idx = entry.value();
            let uuid = &idx.uuid;
            let its = idx.timestamp;
            if let Some(didx) = self.load_dvq(uuid) {
                if VQueue::newer(didx.timestamp, its) {
                    continue;
                }
            }
            if let Some(ref vec_arc) = idx.vector {
                if !f(uuid, vec_arc, its) {
                    break;
                }
            }
        }
        Ok(())
    }

    fn ivq_len(&self) -> usize {
        self.ic.load(AtomicOrdering::SeqCst) as usize
    }

    fn dvq_len(&self) -> usize {
        self.dc.load(AtomicOrdering::SeqCst) as usize
    }

    fn range_pop_insert<F>(&self, now: i64, mut f: F) -> Result<(), QueueError>
    where
        F: FnMut(&str, &Arc<Vec<f32>>, i64) -> bool,
    {
        let mut items: Vec<Index> = Vec::new();
        for entry in self.il.iter() {
            if self.cancel_flag.load(AtomicOrdering::SeqCst) != 0 {
                break;
            }
            let idx = entry.value().clone();
            if VQueue::newer(idx.timestamp, now) {
                continue;
            }
            if let Some(didx) = self.load_dvq(&idx.uuid) {
                if VQueue::newer(didx.timestamp, idx.timestamp) {
                    let _ = self.il.remove(&idx.uuid);
                    self.ic.fetch_sub(1, AtomicOrdering::SeqCst);
                    continue;
                }
            }
            items.push(idx);
        }
        items.sort_by(|a, b| b.timestamp.cmp(&a.timestamp));
        for idx in items {
            if self.cancel_flag.load(AtomicOrdering::SeqCst) != 0 {
                break;
            }
            if let Some(ref vec_arc) = idx.vector {
                if !f(&idx.uuid, vec_arc, idx.timestamp) {
                    return Ok(());
                }
            }
            if self.il.remove(&idx.uuid).is_some() {
                self.ic.fetch_sub(1, AtomicOrdering::SeqCst);
            }
        }
        Ok(())
    }

    fn range_pop_delete<F>(&self, now: i64, mut f: F) -> Result<(), QueueError>
    where
        F: FnMut(&str) -> bool,
    {
        let mut items: Vec<Index> = Vec::new();
        for entry in self.dl.iter() {
            if self.cancel_flag.load(AtomicOrdering::SeqCst) != 0 {
                break;
            }
            let idx = entry.value().clone();
            if VQueue::newer(idx.timestamp, now) {
                continue;
            }
            items.push(idx);
        }
        items.sort_by(|a, b| b.timestamp.cmp(&a.timestamp));
        for didx in items {
            if self.cancel_flag.load(AtomicOrdering::SeqCst) != 0 {
                break;
            }
            if !f(&didx.uuid) {
                return Ok(());
            }
            if self.dl.remove(&didx.uuid).is_some() {
                self.dc.fetch_sub(1, AtomicOrdering::SeqCst);
            }
            if let Some(iv) = self.load_ivq(&didx.uuid) {
                if VQueue::newer(didx.timestamp, iv.timestamp) {
                    if self.il.remove(&didx.uuid).is_some() {
                        self.ic.fetch_sub(1, AtomicOrdering::SeqCst);
                    }
                }
            }
        }
        Ok(())
    }

    fn iv_exists(&self, uuid: &str) -> Result<Option<i64>, QueueError> {
        if uuid.trim().is_empty() {
            return Err(QueueError::InvalidUuid);
        }
        if let Some((_, ts)) = self.get_vector(uuid)? {
            return Ok(Some(ts));
        }
        Ok(None)
    }

    fn dv_exists(&self, uuid: &str) -> Result<Option<i64>, QueueError> {
        if uuid.trim().is_empty() {
            return Err(QueueError::InvalidUuid);
        }
        if let Some((_vec, its)) = self.get_vector(uuid)? {
            if let Ok(Some(dts)) = self.pop_delete(uuid) {
                if VQueue::newer(dts, its) {
                    return Ok(Some(dts));
                }
            }
            return Ok(None);
        } else if let Some(didx) = self.load_dvq(uuid) {
            return Ok(Some(didx.timestamp));
        }
        Ok(None)
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use std::sync::Arc;

    #[test]
    fn test_push_pop() {
        let q = VQueue::new();
        let vec = Arc::new(vec![1.0, 2.0]);
        assert!(q.push_insert("key1".into(), vec.clone(), None).is_ok());
        assert_eq!(q.ivq_len(), 1);
        if let Ok(Some((vec_out, ts))) = q.pop_insert("key1") {
            assert_eq!(*vec_out, vec![1.0, 2.0]);
            assert!(ts > 0);
        } else {
            panic!("PopInsert failed");
        }
        assert_eq!(q.ivq_len(), 0);
    }

    #[test]
    fn test_delete() {
        let q = VQueue::new();
        assert!(q.push_delete("key1".into(), None).is_ok());
        assert_eq!(q.dvq_len(), 1);
        if let Ok(Some(ts)) = q.pop_delete("key1") {
            assert!(ts > 0);
        } else {
            panic!("PopDelete failed");
        }
        assert_eq!(q.dvq_len(), 0);
    }

    #[test]
    fn test_get_vector() {
        let q = VQueue::new();
        let vec = Arc::new(vec![3.0]);
        assert!(q.push_insert("key1".into(), vec.clone(), Some(100)).is_ok());
        assert!(q.push_delete("key1".into(), Some(50)).is_ok());
        if let Ok(Some((vec_out, ts))) = q.get_vector("key1") {
            assert_eq!(*vec_out, vec![3.0]);
            assert_eq!(ts, 100);
        } else {
            panic!("GetVector failed");
        }
        assert!(q.push_delete("key1".into(), Some(150)).is_ok());
        assert!(q.get_vector("key1").unwrap().is_none());
    }

    #[test]
    fn test_range() {
        let q = VQueue::new();
        let vec1 = Arc::new(vec![3.0]);
        let vec2 = Arc::new(vec![4.0]);
        assert!(q.push_insert("a".into(), vec1.clone(), Some(100)).is_ok());
        assert!(q.push_insert("b".into(), vec2.clone(), Some(200)).is_ok());
        let mut seen = Vec::new();
        q.range(|uuid, vec_arc, ts| {
            seen.push((uuid.to_string(), vec_arc.clone().clone(), ts));
            true
        })
        .unwrap();
        assert_eq!(seen.len(), 2);
    }
}
