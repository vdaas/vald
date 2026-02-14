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

//! Daemon module for managing background tasks.
//!
//! This module provides functionality for running periodic background tasks such as:
//! - Auto indexing: Periodically creates indexes when vqueue reaches a threshold
//! - Auto save: Periodically saves indexes to disk
//! - Index limit: Force creates and saves index after a time limit

use std::sync::Arc;
use std::time::Duration;

use algorithm::{ANN, Error};
use tokio::sync::{RwLock, mpsc};
use tokio::time::{Instant, interval};
use tokio_util::sync::CancellationToken;
use tracing::{debug, info, warn};

/// Configuration for the daemon background tasks.
#[derive(Debug, Clone)]
pub struct DaemonConfig {
    /// Duration between auto indexing checks.
    /// If <= 0, auto indexing is effectively disabled (uses very long duration).
    pub auto_index_check_duration: Duration,

    /// Duration between auto save checks.
    /// If <= 0, auto save is effectively disabled.
    pub auto_save_index_duration: Duration,

    /// Time limit for forcing index creation and save.
    /// If <= 0, this limit is disabled.
    pub auto_index_limit: Duration,

    /// Minimum number of items in vqueue before triggering auto index.
    pub auto_index_length: usize,

    /// Pool size for create index operation.
    pub pool_size: u32,

    /// Initial delay before starting the daemon loop.
    pub initial_delay: Duration,

    /// Enable proactive garbage collection.
    pub enable_proactive_gc: bool,
}

impl Default for DaemonConfig {
    fn default() -> Self {
        Self {
            auto_index_check_duration: Duration::from_secs(1),
            auto_save_index_duration: Duration::from_secs(60),
            auto_index_limit: Duration::from_secs(3600), // 1 hour
            auto_index_length: 100,
            pool_size: 10000,
            initial_delay: Duration::ZERO,
            enable_proactive_gc: false,
        }
    }
}

impl DaemonConfig {
    /// Creates a new DaemonConfig from config settings.
    pub fn from_config(config: &crate::config::Daemon) -> Self {
        Self {
            auto_index_check_duration: Duration::from_millis(config.auto_index_check_duration_ms),
            auto_save_index_duration: Duration::from_millis(config.auto_save_index_duration_ms),
            auto_index_limit: Duration::from_millis(config.auto_index_limit_ms),
            auto_index_length: config.auto_index_length,
            pool_size: config.pool_size,
            initial_delay: Duration::from_millis(config.initial_delay_ms),
            enable_proactive_gc: config.enable_proactive_gc,
        }
    }
}

/// Handle for controlling the daemon.
#[derive(Clone)]
pub struct DaemonHandle {
    cancel_token: CancellationToken,
    /// Sender to notify when daemon has completed shutdown.
    shutdown_complete: Arc<tokio::sync::Notify>,
}

impl DaemonHandle {
    /// Signals the daemon to stop.
    pub fn stop(&self) {
        self.cancel_token.cancel();
    }

    /// Returns true if the daemon has been signaled to stop.
    pub fn is_cancelled(&self) -> bool {
        self.cancel_token.is_cancelled()
    }

    /// Waits for the daemon to complete shutdown.
    /// This should be called after stop() to ensure graceful shutdown.
    pub async fn wait(&self) {
        self.shutdown_complete.notified().await;
    }

    /// Stops the daemon and waits for it to complete.
    pub async fn stop_and_wait(&self) {
        self.stop();
        self.wait().await;
    }
}

/// Starts the daemon background tasks for the given ANN service.
///
/// This function spawns a background task that periodically:
/// 1. Checks if vqueue has enough items and creates an index if needed
/// 2. Saves the index to disk at regular intervals
/// 3. Forces index creation and save after a time limit
///
/// # Arguments
///
/// * `service` - Arc-wrapped RwLock of the ANN service
/// * `config` - Daemon configuration
///
/// # Returns
///
/// A tuple of (DaemonHandle, mpsc::Receiver<Error>):
/// - DaemonHandle: Used to control the daemon (stop it)
/// - Receiver: Receives any errors that occur during daemon operations
///
/// # Example
///
/// ```ignore
/// let service = Arc::new(RwLock::new(QBGService::new(settings).await));
/// let config = DaemonConfig::default();
/// let (handle, mut error_rx) = start(service.clone(), config).await;
///
/// // Handle errors in another task
/// tokio::spawn(async move {
///     while let Some(err) = error_rx.recv().await {
///         eprintln!("Daemon error: {:?}", err);
///     }
/// });
///
/// // Later, stop the daemon and wait for completion
/// handle.stop_and_wait().await;
/// ```
pub async fn start<T: ANN + 'static>(
    service: Arc<RwLock<T>>,
    config: DaemonConfig,
) -> (DaemonHandle, mpsc::Receiver<Error>) {
    let (error_tx, error_rx) = mpsc::channel::<Error>(16);
    let cancel_token = CancellationToken::new();
    let shutdown_complete = Arc::new(tokio::sync::Notify::new());
    let shutdown_complete_clone = shutdown_complete.clone();
    let handle = DaemonHandle {
        cancel_token: cancel_token.clone(),
        shutdown_complete,
    };

    let daemon_task = async move {
        // Apply initial delay if configured
        if !config.initial_delay.is_zero() {
            tokio::select! {
                _ = cancel_token.cancelled() => {
                    info!("Daemon cancelled during initial delay");
                    shutdown_complete_clone.notify_waiters();
                    return;
                }
                _ = tokio::time::sleep(config.initial_delay) => {}
            }
        }

        // Use very long intervals for disabled features
        let max_duration = Duration::from_secs(u64::MAX / 2);

        let index_check_interval = if config.auto_index_check_duration.is_zero() {
            max_duration
        } else {
            config.auto_index_check_duration
        };

        let save_interval = if config.auto_save_index_duration.is_zero() {
            max_duration
        } else {
            config.auto_save_index_duration
        };

        let limit_interval = if config.auto_index_limit.is_zero() {
            max_duration
        } else {
            config.auto_index_limit
        };

        let mut index_tick = interval(index_check_interval);
        let mut save_tick = interval(save_interval);
        let mut limit_tick = interval(limit_interval);

        // Skip immediate first tick
        index_tick.tick().await;
        save_tick.tick().await;
        limit_tick.tick().await;

        let start_time = Instant::now();

        loop {
            tokio::select! {
                _ = cancel_token.cancelled() => {
                    info!("Daemon shutdown requested, performing final index creation...");
                    // Perform final index creation before shutdown
                    let mut svc = service.write().await;
                    if let Err(e) = svc.create_index().await {
                        if !matches!(e, Error::UncommittedIndexNotFound {}) {
                            let _ = error_tx.send(e).await;
                        }
                    }
                    info!("Daemon shutdown complete");
                    shutdown_complete_clone.notify_waiters();
                    return;
                }

                _ = index_tick.tick() => {
                    let svc = service.read().await;
                    let ivq_len = svc.insert_vqueue_buffer_len() as usize;
                    let is_flushing = svc.is_flushing();
                    drop(svc);

                    if !is_flushing && ivq_len >= config.auto_index_length {
                        debug!("Auto index triggered: vqueue len {} >= threshold {}", ivq_len, config.auto_index_length);
                        let mut svc = service.write().await;
                        if let Err(e) = svc.create_index().await {
                            if !matches!(e, Error::UncommittedIndexNotFound {}) {
                                warn!("Auto index creation failed: {:?}", e);
                                let _ = error_tx.send(e).await;
                            }
                        }
                    }
                }

                _ = limit_tick.tick() => {
                    debug!("Index limit reached after {:?}, forcing create and save", start_time.elapsed());
                    let mut svc = service.write().await;
                    if let Err(e) = svc.create_and_save_index().await {
                        if !matches!(e, Error::UncommittedIndexNotFound {}) {
                            warn!("Forced create and save index failed: {:?}", e);
                            let _ = error_tx.send(e).await;
                        }
                    }
                }

                _ = save_tick.tick() => {
                    debug!("Auto save index triggered");
                    let mut svc = service.write().await;
                    if let Err(e) = svc.save_index().await {
                        warn!("Auto save index failed: {:?}", e);
                        let _ = error_tx.send(e).await;
                    }
                }
            }

            // Proactive GC if enabled (Rust doesn't have manual GC, but we can hint)
            if config.enable_proactive_gc {
                // In Rust, memory is managed automatically.
                // This is a placeholder for any custom memory management if needed.
                // For example, clearing caches or compacting data structures.
            }
        }
    };

    tokio::spawn(daemon_task);

    (handle, error_rx)
}

#[cfg(test)]
mod tests {
    use super::*;
    use proto::payload::v1::{info, search};
    use std::collections::HashMap;
    use std::sync::atomic::{AtomicU32, Ordering};

    /// Mock ANN service for testing
    struct MockANNService {
        ivq_len: AtomicU32,
        create_index_count: AtomicU32,
        save_index_count: AtomicU32,
        is_flushing: bool,
    }

    impl MockANNService {
        fn new() -> Self {
            Self {
                ivq_len: AtomicU32::new(0),
                create_index_count: AtomicU32::new(0),
                save_index_count: AtomicU32::new(0),
                is_flushing: false,
            }
        }

        fn set_ivq_len(&self, len: u32) {
            self.ivq_len.store(len, Ordering::SeqCst);
        }

        fn get_create_index_count(&self) -> u32 {
            self.create_index_count.load(Ordering::SeqCst)
        }

        fn get_save_index_count(&self) -> u32 {
            self.save_index_count.load(Ordering::SeqCst)
        }
    }

    impl ANN for MockANNService {
        async fn search(
            &self,
            _vector: Vec<f32>,
            _k: u32,
            _epsilon: f32,
            _radius: f32,
        ) -> Result<search::Response, Error> {
            Ok(search::Response::default())
        }

        async fn search_by_id(
            &self,
            _uuid: String,
            _k: u32,
            _epsilon: f32,
            _radius: f32,
        ) -> Result<search::Response, Error> {
            Ok(search::Response::default())
        }

        async fn linear_search(
            &self,
            _vector: Vec<f32>,
            _k: u32,
        ) -> Result<search::Response, Error> {
            Ok(search::Response::default())
        }

        async fn linear_search_by_id(
            &self,
            _uuid: String,
            _k: u32,
        ) -> Result<search::Response, Error> {
            Ok(search::Response::default())
        }

        async fn insert(&mut self, _uuid: String, _vector: Vec<f32>) -> Result<(), Error> {
            Ok(())
        }

        async fn insert_with_time(
            &mut self,
            _uuid: String,
            _vector: Vec<f32>,
            _t: i64,
        ) -> Result<(), Error> {
            Ok(())
        }

        async fn insert_multiple(
            &mut self,
            _vectors: HashMap<String, Vec<f32>>,
        ) -> Result<(), Error> {
            Ok(())
        }

        async fn insert_multiple_with_time(
            &mut self,
            _vectors: HashMap<String, Vec<f32>>,
            _t: i64,
        ) -> Result<(), Error> {
            Ok(())
        }

        async fn update(&mut self, _uuid: String, _vector: Vec<f32>) -> Result<(), Error> {
            Ok(())
        }

        async fn update_with_time(
            &mut self,
            _uuid: String,
            _vector: Vec<f32>,
            _t: i64,
        ) -> Result<(), Error> {
            Ok(())
        }

        async fn update_multiple(
            &mut self,
            _vectors: HashMap<String, Vec<f32>>,
        ) -> Result<(), Error> {
            Ok(())
        }

        async fn update_multiple_with_time(
            &mut self,
            _vectors: HashMap<String, Vec<f32>>,
            _t: i64,
        ) -> Result<(), Error> {
            Ok(())
        }

        async fn update_timestamp(
            &mut self,
            _uuid: String,
            _t: i64,
            _force: bool,
        ) -> Result<(), Error> {
            Ok(())
        }

        async fn remove(&mut self, _uuid: String) -> Result<(), Error> {
            Ok(())
        }

        async fn remove_with_time(&mut self, _uuid: String, _t: i64) -> Result<(), Error> {
            Ok(())
        }

        async fn remove_multiple(&mut self, _uuids: Vec<String>) -> Result<(), Error> {
            Ok(())
        }

        async fn remove_multiple_with_time(
            &mut self,
            _uuids: Vec<String>,
            _t: i64,
        ) -> Result<(), Error> {
            Ok(())
        }

        async fn regenerate_indexes(&mut self) -> Result<(), Error> {
            Ok(())
        }

        async fn create_index(&mut self) -> Result<(), Error> {
            self.create_index_count.fetch_add(1, Ordering::SeqCst);
            Ok(())
        }

        async fn save_index(&mut self) -> Result<(), Error> {
            self.save_index_count.fetch_add(1, Ordering::SeqCst);
            Ok(())
        }

        async fn create_and_save_index(&mut self) -> Result<(), Error> {
            self.create_index().await?;
            self.save_index().await
        }

        async fn get_object(&self, _uuid: String) -> Result<(Vec<f32>, i64), Error> {
            Err(Error::ObjectIDNotFound {
                uuid: "not found".to_string(),
            })
        }

        async fn exists(&self, _uuid: String) -> (usize, bool) {
            (0, false)
        }

        async fn uuids(&self) -> Vec<String> {
            vec![]
        }

        async fn list_object_func<F: FnMut(String, Vec<f32>, i64) -> bool + Send>(&self, _f: F) {}

        fn is_indexing(&self) -> bool {
            false
        }

        fn is_flushing(&self) -> bool {
            self.is_flushing
        }

        fn is_saving(&self) -> bool {
            false
        }

        fn len(&self) -> u32 {
            0
        }

        fn number_of_create_index_executions(&self) -> u64 {
            self.create_index_count.load(Ordering::SeqCst) as u64
        }

        fn insert_vqueue_buffer_len(&self) -> u32 {
            self.ivq_len.load(Ordering::SeqCst)
        }

        fn delete_vqueue_buffer_len(&self) -> u32 {
            0
        }

        fn get_dimension_size(&self) -> usize {
            128
        }

        fn broken_index_count(&self) -> u64 {
            0
        }

        fn is_statistics_enabled(&self) -> bool {
            false
        }

        fn index_statistics(&self) -> Result<info::index::Statistics, Error> {
            Ok(info::index::Statistics::default())
        }

        fn index_property(&self) -> Result<info::index::Property, Error> {
            Ok(info::index::Property::default())
        }

        async fn close(&mut self) -> Result<(), Error> {
            Ok(())
        }
    }

    #[tokio::test]
    async fn test_daemon_auto_index() {
        let service = Arc::new(RwLock::new(MockANNService::new()));

        let config = DaemonConfig {
            auto_index_check_duration: Duration::from_millis(50),
            auto_save_index_duration: Duration::from_secs(3600), // Disable save for this test
            auto_index_limit: Duration::from_secs(3600),         // Disable limit for this test
            auto_index_length: 10,
            pool_size: 100,
            initial_delay: Duration::ZERO,
            enable_proactive_gc: false,
        };

        // Set vqueue length above threshold
        service.read().await.set_ivq_len(15);

        let (handle, _error_rx) = start(service.clone(), config).await;

        // Wait for auto index to trigger
        tokio::time::sleep(Duration::from_millis(150)).await;

        // Check that create_index was called
        let create_count = service.read().await.get_create_index_count();
        assert!(
            create_count >= 1,
            "Expected at least 1 create_index call, got {}",
            create_count
        );

        handle.stop();
        tokio::time::sleep(Duration::from_millis(50)).await;
    }

    #[tokio::test]
    async fn test_daemon_auto_save() {
        let service = Arc::new(RwLock::new(MockANNService::new()));

        let config = DaemonConfig {
            auto_index_check_duration: Duration::from_secs(3600), // Disable index check
            auto_save_index_duration: Duration::from_millis(50),
            auto_index_limit: Duration::from_secs(3600), // Disable limit
            auto_index_length: 1000,
            pool_size: 100,
            initial_delay: Duration::ZERO,
            enable_proactive_gc: false,
        };

        let (handle, _error_rx) = start(service.clone(), config).await;

        // Wait for auto save to trigger
        tokio::time::sleep(Duration::from_millis(150)).await;

        // Check that save_index was called
        let save_count = service.read().await.get_save_index_count();
        assert!(
            save_count >= 1,
            "Expected at least 1 save_index call, got {}",
            save_count
        );

        handle.stop();
        tokio::time::sleep(Duration::from_millis(50)).await;
    }

    #[tokio::test]
    async fn test_daemon_handle_stop() {
        let service = Arc::new(RwLock::new(MockANNService::new()));

        let config = DaemonConfig::default();
        let (handle, _error_rx) = start(service.clone(), config).await;

        assert!(!handle.is_cancelled());
        handle.stop();
        assert!(handle.is_cancelled());

        // Give daemon time to shut down
        tokio::time::sleep(Duration::from_millis(50)).await;
    }

    #[tokio::test]
    async fn test_daemon_initial_delay() {
        let service = Arc::new(RwLock::new(MockANNService::new()));

        let config = DaemonConfig {
            auto_index_check_duration: Duration::from_millis(10),
            auto_save_index_duration: Duration::from_secs(3600),
            auto_index_limit: Duration::from_secs(3600),
            auto_index_length: 0, // Always trigger
            pool_size: 100,
            initial_delay: Duration::from_millis(100),
            enable_proactive_gc: false,
        };

        service.read().await.set_ivq_len(100);
        let (handle, _error_rx) = start(service.clone(), config).await;

        // Immediately after start, create_index should not have been called
        tokio::time::sleep(Duration::from_millis(20)).await;
        let count_before_delay = service.read().await.get_create_index_count();
        assert_eq!(
            count_before_delay, 0,
            "Should not have created index during initial delay"
        );

        // After initial delay passes
        tokio::time::sleep(Duration::from_millis(150)).await;
        let count_after_delay = service.read().await.get_create_index_count();
        assert!(
            count_after_delay >= 1,
            "Should have created index after initial delay"
        );

        handle.stop();
    }

    #[tokio::test]
    async fn test_daemon_skips_when_flushing() {
        let mut mock = MockANNService::new();
        mock.is_flushing = true;
        mock.ivq_len.store(1000, Ordering::SeqCst);
        let service = Arc::new(RwLock::new(mock));

        let config = DaemonConfig {
            auto_index_check_duration: Duration::from_millis(20),
            auto_save_index_duration: Duration::from_secs(3600),
            auto_index_limit: Duration::from_secs(3600),
            auto_index_length: 10,
            pool_size: 100,
            initial_delay: Duration::ZERO,
            enable_proactive_gc: false,
        };

        let (handle, _error_rx) = start(service.clone(), config).await;

        // Wait for several ticks
        tokio::time::sleep(Duration::from_millis(100)).await;

        // create_index should not have been called because is_flushing is true
        let count = service.read().await.get_create_index_count();
        assert_eq!(count, 0, "Should not have created index while flushing");

        handle.stop();
    }

    #[tokio::test]
    async fn test_daemon_shutdown_creates_final_index() {
        let service = Arc::new(RwLock::new(MockANNService::new()));

        let config = DaemonConfig {
            auto_index_check_duration: Duration::from_secs(3600), // Disable periodic
            auto_save_index_duration: Duration::from_secs(3600),
            auto_index_limit: Duration::from_secs(3600),
            auto_index_length: 1000,
            pool_size: 100,
            initial_delay: Duration::ZERO,
            enable_proactive_gc: false,
        };

        let (handle, _error_rx) = start(service.clone(), config).await;

        // No periodic index creation should have happened
        tokio::time::sleep(Duration::from_millis(50)).await;
        let count_before = service.read().await.get_create_index_count();
        assert_eq!(count_before, 0);

        // Stop the daemon - this should trigger final index creation
        handle.stop();
        tokio::time::sleep(Duration::from_millis(100)).await;

        let count_after = service.read().await.get_create_index_count();
        assert_eq!(
            count_after, 1,
            "Should have created final index on shutdown"
        );
    }

    // ========== Graceful Shutdown Tests ==========

    #[tokio::test]
    async fn test_daemon_stop_and_wait() {
        let service = Arc::new(RwLock::new(MockANNService::new()));

        let config = DaemonConfig {
            auto_index_check_duration: Duration::from_secs(3600),
            auto_save_index_duration: Duration::from_secs(3600),
            auto_index_limit: Duration::from_secs(3600),
            auto_index_length: 1000,
            pool_size: 100,
            initial_delay: Duration::ZERO,
            enable_proactive_gc: false,
        };

        let (handle, _error_rx) = start(service.clone(), config).await;

        // Give the daemon time to start
        tokio::time::sleep(Duration::from_millis(10)).await;

        // stop_and_wait should complete and create final index
        let start_time = std::time::Instant::now();
        handle.stop_and_wait().await;
        let elapsed = start_time.elapsed();

        // Should complete quickly (within 500ms for test)
        assert!(
            elapsed < Duration::from_millis(500),
            "stop_and_wait took too long: {:?}",
            elapsed
        );

        // Should have called create_index on shutdown
        let count = service.read().await.get_create_index_count();
        assert_eq!(count, 1, "Should have created final index on shutdown");
    }

    #[tokio::test]
    async fn test_daemon_wait_after_stop() {
        let service = Arc::new(RwLock::new(MockANNService::new()));

        let config = DaemonConfig {
            auto_index_check_duration: Duration::from_secs(3600),
            auto_save_index_duration: Duration::from_secs(3600),
            auto_index_limit: Duration::from_secs(3600),
            auto_index_length: 1000,
            pool_size: 100,
            initial_delay: Duration::ZERO,
            enable_proactive_gc: false,
        };

        let (handle, _error_rx) = start(service.clone(), config).await;

        // Stop first
        handle.stop();
        assert!(handle.is_cancelled());

        // Then wait - should complete immediately since stop already triggered
        let start_time = std::time::Instant::now();
        handle.wait().await;
        let elapsed = start_time.elapsed();

        assert!(
            elapsed < Duration::from_millis(200),
            "wait() took too long: {:?}",
            elapsed
        );
    }

    #[tokio::test]
    async fn test_daemon_multiple_wait_calls() {
        let service = Arc::new(RwLock::new(MockANNService::new()));

        let config = DaemonConfig::default();
        let (handle, _error_rx) = start(service.clone(), config).await;

        // Clone handle for multiple waiters
        let handle2 = handle.clone();

        // Spawn multiple tasks that wait
        let wait1 = tokio::spawn(async move {
            handle.stop_and_wait().await;
            "waiter1"
        });

        let wait2 = tokio::spawn(async move {
            handle2.wait().await;
            "waiter2"
        });

        // Both should complete
        let result1 = wait1.await.unwrap();
        let result2 = wait2.await.unwrap();

        assert_eq!(result1, "waiter1");
        assert_eq!(result2, "waiter2");
    }

    #[tokio::test]
    async fn test_daemon_shutdown_during_initial_delay() {
        let service = Arc::new(RwLock::new(MockANNService::new()));

        let config = DaemonConfig {
            auto_index_check_duration: Duration::from_millis(10),
            auto_save_index_duration: Duration::from_secs(3600),
            auto_index_limit: Duration::from_secs(3600),
            auto_index_length: 0,
            pool_size: 100,
            initial_delay: Duration::from_secs(10), // Very long initial delay
            enable_proactive_gc: false,
        };

        let (handle, _error_rx) = start(service.clone(), config).await;

        // Stop immediately (during initial delay)
        tokio::time::sleep(Duration::from_millis(10)).await;
        let start_time = std::time::Instant::now();
        handle.stop_and_wait().await;
        let elapsed = start_time.elapsed();

        // Should stop quickly, not wait for full initial delay
        assert!(
            elapsed < Duration::from_millis(500),
            "Shutdown should be fast: {:?}",
            elapsed
        );

        // No index creation should have happened (cancelled during initial delay)
        let count = service.read().await.get_create_index_count();
        assert_eq!(
            count, 0,
            "Should not create index when cancelled during initial delay"
        );
    }

    #[tokio::test]
    async fn test_daemon_graceful_shutdown_with_pending_operations() {
        let service = Arc::new(RwLock::new(MockANNService::new()));

        // Set high vqueue length to simulate pending operations
        service.read().await.set_ivq_len(1000);

        let config = DaemonConfig {
            auto_index_check_duration: Duration::from_millis(50),
            auto_save_index_duration: Duration::from_secs(3600),
            auto_index_limit: Duration::from_secs(3600),
            auto_index_length: 100, // Threshold lower than vqueue
            pool_size: 100,
            initial_delay: Duration::ZERO,
            enable_proactive_gc: false,
        };

        let (handle, _error_rx) = start(service.clone(), config).await;

        // Let it run for a bit and create some indexes
        tokio::time::sleep(Duration::from_millis(100)).await;

        let count_before = service.read().await.get_create_index_count();
        assert!(count_before >= 1, "Should have auto-indexed");

        // Now stop and wait
        handle.stop_and_wait().await;

        // Should have created one more final index
        let count_after = service.read().await.get_create_index_count();
        assert!(
            count_after > count_before,
            "Should have created final index on shutdown"
        );
    }
}
