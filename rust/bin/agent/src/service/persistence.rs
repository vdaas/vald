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

//! Index persistence management for load/save/recovery operations.
//!
//! This module provides functionality to:
//! - Prepare index directories (origin, backup, broken)
//! - Load existing indexes from disk with fallback paths
//! - Save indexes atomically with concurrent writes
//! - Backup broken indexes with history limit
//! - Support Copy-on-Write (CoW) mode for safe updates

use std::fs;
use std::path::{Path, PathBuf};
use std::sync::atomic::{AtomicU64, Ordering};
use std::time::{SystemTime, UNIX_EPOCH};

use thiserror::Error;
use tracing::{debug, info, warn};

use super::metadata::{self, AGENT_METADATA_FILENAME, Metadata};

/// Directory name for backup index (Copy-on-Write mode).
const OLD_INDEX_DIR_NAME: &str = "backup";
/// Directory name for the origin/primary index.
const ORIGIN_INDEX_DIR_NAME: &str = "origin";
/// Directory name for broken index backups.
const BROKEN_INDEX_DIR_NAME: &str = "broken";

/// Errors that can occur during persistence operations.
#[derive(Debug, Error)]
pub enum PersistenceError {
    #[error("index file not found: {0}")]
    IndexFileNotFound(String),

    #[error("metadata file not found: {0}")]
    MetadataNotFound(String),

    #[error("invalid index: {0}")]
    InvalidIndex(String),

    #[error("index load timeout")]
    LoadTimeout,

    #[error("failed to prepare folders: {0}")]
    PrepareFoldersFailed(String),

    #[error("failed to backup broken index: {0}")]
    BackupFailed(String),

    #[error("failed to save index: {0}")]
    SaveFailed(String),

    #[error("io error: {0}")]
    IoError(#[from] std::io::Error),

    #[error("metadata error: {0}")]
    MetadataError(#[from] metadata::MetadataError),
}

/// Configuration for persistence operations.
#[derive(Debug, Clone)]
pub struct PersistenceConfig {
    /// Whether Copy-on-Write mode is enabled.
    pub enable_copy_on_write: bool,
    /// Maximum number of broken index generations to keep.
    pub broken_index_history_limit: usize,
}

impl Default for PersistenceConfig {
    fn default() -> Self {
        PersistenceConfig {
            enable_copy_on_write: false,
            broken_index_history_limit: 3,
        }
    }
}

/// Paths used for index persistence.
#[derive(Debug, Clone)]
pub struct IndexPaths {
    /// The base path (user-configured index path).
    pub base_path: PathBuf,
    /// The primary index path (base_path/origin).
    pub primary_path: PathBuf,
    /// The old/backup path for CoW mode (base_path/backup).
    pub old_path: PathBuf,
    /// The broken index backup path (base_path/broken).
    pub broken_path: PathBuf,
    /// Temporary path for atomic saves (only used in CoW mode).
    pub tmp_path: Option<PathBuf>,
}

impl IndexPaths {
    /// Creates a new IndexPaths from the base path.
    pub fn new<P: AsRef<Path>>(base_path: P) -> Self {
        let base = base_path.as_ref().to_path_buf();
        IndexPaths {
            primary_path: base.join(ORIGIN_INDEX_DIR_NAME),
            old_path: base.join(OLD_INDEX_DIR_NAME),
            broken_path: base.join(BROKEN_INDEX_DIR_NAME),
            base_path: base,
            tmp_path: None,
        }
    }

    /// Returns the metadata file path for the primary index.
    pub fn metadata_path(&self) -> PathBuf {
        self.primary_path.join(AGENT_METADATA_FILENAME)
    }
}

/// Manages index persistence state and filesystem paths.
pub struct PersistenceManager {
    config: PersistenceConfig,
    paths: IndexPaths,
    broken_index_count: AtomicU64,
    /// Temporary path for atomic saves in CoW mode.
    tmp_path: std::sync::RwLock<Option<PathBuf>>,
}

impl PersistenceManager {
    /// Creates a new PersistenceManager.
    pub fn new<P: AsRef<Path>>(base_path: P, config: PersistenceConfig) -> Self {
        PersistenceManager {
            paths: IndexPaths::new(base_path),
            config,
            broken_index_count: AtomicU64::new(0),
            tmp_path: std::sync::RwLock::new(None),
        }
    }

    /// Returns the paths managed by this instance.
    pub fn paths(&self) -> &IndexPaths {
        &self.paths
    }

    /// Returns the number of broken index backups.
    pub fn broken_index_count(&self) -> u64 {
        self.broken_index_count.load(Ordering::SeqCst)
    }

    /// Prepares the folder structure for index persistence.
    ///
    /// Creates the following directories if they don't exist:
    /// - base_path (for the index)
    /// - base_path/broken (broken index backups)
    /// - base_path/backup (if CoW is enabled)
    ///
    /// Note: base_path/origin is NOT created here because the index library (QBG/NGT)
    /// expects to create this directory itself during index initialization.
    pub fn prepare_folders(&self) -> Result<(), PersistenceError> {
        // Create base path if needed (parent of primary path)
        fs::create_dir_all(&self.paths.base_path).map_err(|e| {
            PersistenceError::PrepareFoldersFailed(format!(
                "failed to create base path {}: {}",
                self.paths.base_path.display(),
                e
            ))
        })?;
        debug!(
            "ensured base path exists: {}",
            self.paths.base_path.display()
        );

        // Create broken index backup directory
        fs::create_dir_all(&self.paths.broken_path).map_err(|e| {
            warn!("failed to create broken index directory: {}", e);
            PersistenceError::PrepareFoldersFailed(format!(
                "failed to create broken path {}: {}",
                self.paths.broken_path.display(),
                e
            ))
        })?;
        debug!(
            "created broken index directory: {}",
            self.paths.broken_path.display()
        );

        // Update broken index count
        if let Ok(entries) = fs::read_dir(&self.paths.broken_path) {
            let count = entries.filter_map(|e| e.ok()).count() as u64;
            self.broken_index_count.store(count, Ordering::SeqCst);
            debug!("broken index count: {}", count);
        }

        // Create old/backup directory if CoW is enabled
        if self.config.enable_copy_on_write {
            fs::create_dir_all(&self.paths.old_path).map_err(|e| {
                PersistenceError::PrepareFoldersFailed(format!(
                    "failed to create old path {}: {}",
                    self.paths.old_path.display(),
                    e
                ))
            })?;
            debug!(
                "created old/backup directory: {}",
                self.paths.old_path.display()
            );
        }

        Ok(())
    }

    /// Checks if the index at the given path needs to be backed up.
    ///
    /// Returns true if:
    /// - The path contains .json or .kvsdb files AND
    /// - metadata.json doesn't exist OR is invalid OR has index_count > 0
    pub fn needs_backup<P: AsRef<Path>>(path: P) -> bool {
        let path = path.as_ref();

        let entries = match fs::read_dir(path) {
            Ok(e) => e,
            Err(_) => return false,
        };

        let files: Vec<_> = entries
            .filter_map(|e| e.ok())
            .map(|e| e.file_name().to_string_lossy().to_string())
            .collect();

        if files.is_empty() {
            return false;
        }

        // Check if there are any .json or .kvsdb files (not initial state)
        let has_data_files = files
            .iter()
            .any(|f| f.ends_with(".json") || f.ends_with(".kvsdb"));
        if !has_data_files {
            return false;
        }

        // Check if metadata.json exists
        let metadata_path = path.join(AGENT_METADATA_FILENAME);
        if !metadata_path.exists() {
            return true;
        }

        // Check metadata content
        match metadata::load(&metadata_path) {
            Ok(meta) => meta.is_invalid || meta.index_count() > 0,
            Err(_) => false,
        }
    }

    /// Backs up a broken index to the broken directory.
    ///
    /// The backup directory is named with the current Unix nanosecond timestamp.
    /// If the history limit is exceeded, the oldest backup is removed.
    pub fn backup_broken(&self) -> Result<(), PersistenceError> {
        if self.config.broken_index_history_limit == 0 {
            return Ok(());
        }

        // Check if there's anything to backup
        let source_entries: Vec<_> = fs::read_dir(&self.paths.primary_path)
            .map_err(|e| PersistenceError::BackupFailed(e.to_string()))?
            .filter_map(|e| e.ok())
            .collect();

        if source_entries.is_empty() {
            debug!(
                "no files to backup in {}",
                self.paths.primary_path.display()
            );
            return Ok(());
        }

        // Check current backup count and remove oldest if at limit
        let mut backups: Vec<_> = fs::read_dir(&self.paths.broken_path)
            .map_err(|e| PersistenceError::BackupFailed(e.to_string()))?
            .filter_map(|e| e.ok())
            .map(|e| e.path())
            .collect();

        if backups.len() >= self.config.broken_index_history_limit {
            info!(
                "broken index history limit ({}) reached, removing oldest backup",
                self.config.broken_index_history_limit
            );
            backups.sort();
            if let Some(oldest) = backups.first() {
                fs::remove_dir_all(oldest).map_err(|e| {
                    PersistenceError::BackupFailed(format!(
                        "failed to remove oldest backup {}: {}",
                        oldest.display(),
                        e
                    ))
                })?;
            }
        }

        // Create new backup directory with timestamp
        let timestamp = SystemTime::now()
            .duration_since(UNIX_EPOCH)
            .unwrap()
            .as_nanos();
        let dest = self.paths.broken_path.join(timestamp.to_string());

        // Move the index to the backup directory
        info!("backing up broken index to {}", dest.display());
        move_dir(&self.paths.primary_path, &dest)?;

        // Update broken index count
        if let Ok(entries) = fs::read_dir(&self.paths.broken_path) {
            let count = entries.filter_map(|e| e.ok()).count() as u64;
            self.broken_index_count.store(count, Ordering::SeqCst);
            debug!("broken index count updated: {}", count);
        }

        // Recreate the primary path
        fs::create_dir_all(&self.paths.primary_path).map_err(|e| {
            PersistenceError::BackupFailed(format!(
                "failed to recreate primary path after backup: {}",
                e
            ))
        })?;

        Ok(())
    }

    /// Checks if an index exists at the primary path and is valid.
    ///
    /// Returns true if:
    /// - The primary path exists
    /// - metadata.json exists and is valid
    /// - index_count > 0
    pub fn index_exists(&self) -> bool {
        if !self.paths.primary_path.exists() {
            return false;
        }

        let metadata_path = self.paths.metadata_path();
        match metadata::load(&metadata_path) {
            Ok(meta) => !meta.is_invalid && meta.index_count() > 0,
            Err(_) => false,
        }
    }

    /// Loads metadata from the primary index path.
    pub fn load_metadata(&self) -> Result<Metadata, PersistenceError> {
        let metadata_path = self.paths.metadata_path();
        metadata::load(&metadata_path).map_err(|e| {
            PersistenceError::MetadataNotFound(format!("{}: {}", metadata_path.display(), e))
        })
    }

    /// Saves metadata to the primary index path.
    pub fn save_metadata(&self, metadata: &Metadata) -> Result<(), PersistenceError> {
        let metadata_path = self.paths.metadata_path();
        metadata::store(&metadata_path, metadata)?;
        Ok(())
    }

    /// Returns whether Copy-on-Write mode is enabled.
    pub fn is_copy_on_write_enabled(&self) -> bool {
        self.config.enable_copy_on_write
    }

    /// Creates a temporary directory for Copy-on-Write saves.
    ///
    /// This method creates a new temporary directory under the system temp directory
    /// and stores the path for later use by `get_save_path` and `move_and_switch_saved_data`.
    pub fn mktmp(&self) -> Result<(), PersistenceError> {
        if !self.config.enable_copy_on_write {
            return Ok(());
        }

        let vald_tmp_dir = std::env::temp_dir().join("vald");
        fs::create_dir_all(&vald_tmp_dir).map_err(|e| {
            PersistenceError::SaveFailed(format!(
                "failed to create vald temp directory {}: {}",
                vald_tmp_dir.display(),
                e
            ))
        })?;

        // Create a unique temp directory using timestamp and random suffix
        let timestamp = SystemTime::now()
            .duration_since(UNIX_EPOCH)
            .unwrap()
            .as_nanos();
        let tmp_name = format!("index-{}", timestamp);
        let tmp_path = vald_tmp_dir.join(&tmp_name);

        fs::create_dir_all(&tmp_path).map_err(|e| {
            PersistenceError::SaveFailed(format!(
                "failed to create temporary index directory {}: {}",
                tmp_path.display(),
                e
            ))
        })?;

        info!(
            "created temporary directory for CoW: {}",
            tmp_path.display()
        );

        let mut guard = self.tmp_path.write().unwrap();
        *guard = Some(tmp_path);

        Ok(())
    }

    /// Returns the path where the index should be saved.
    ///
    /// In Copy-on-Write mode, returns the temporary path.
    /// Otherwise, returns the primary path.
    pub fn get_save_path(&self) -> PathBuf {
        if self.config.enable_copy_on_write
            && let Some(tmp) = self.tmp_path.read().unwrap().as_ref()
        {
            return tmp.clone();
        }
        self.paths.primary_path.clone()
    }

    /// Saves metadata to the appropriate path (tmp for CoW, primary otherwise).
    pub fn save_metadata_to_save_path(&self, metadata: &Metadata) -> Result<(), PersistenceError> {
        let save_path = self.get_save_path();
        let metadata_path = save_path.join(AGENT_METADATA_FILENAME);
        metadata::store(&metadata_path, metadata)?;
        Ok(())
    }

    /// Moves and switches the saved data for Copy-on-Write mode.
    ///
    /// This performs an atomic switch of the index data:
    /// 1. Move `primary_path` (origin) → `old_path` (backup)
    /// 2. Move `tmp_path` → `primary_path` (origin)
    /// 3. Create a new temporary directory
    ///
    /// If step 2 fails, it attempts to rollback by moving backup back to primary.
    pub fn move_and_switch_saved_data(&self) -> Result<(), PersistenceError> {
        if !self.config.enable_copy_on_write {
            return Ok(());
        }

        let tmp_path = {
            let guard = self.tmp_path.read().unwrap();
            match guard.as_ref() {
                Some(p) => p.clone(),
                None => {
                    warn!("move_and_switch_saved_data called but no tmp_path is set");
                    return Ok(());
                }
            }
        };

        info!("starting move and switch saved data operation for copy on write");

        // Step 1: Move primary (origin) → old (backup)
        // First, remove old backup if it exists
        if self.paths.old_path.exists()
            && let Err(e) = fs::remove_dir_all(&self.paths.old_path)
        {
            warn!("failed to remove old backup directory: {}", e);
        }

        // Move primary to backup (only if primary exists and has content)
        if self.paths.primary_path.exists() {
            let has_content = fs::read_dir(&self.paths.primary_path)
                .map(|mut d| d.next().is_some())
                .unwrap_or(false);

            if has_content {
                if let Err(e) = move_dir(&self.paths.primary_path, &self.paths.old_path) {
                    warn!(
                        "failed to backup data from {} to {}: {}",
                        self.paths.primary_path.display(),
                        self.paths.old_path.display(),
                        e
                    );
                } else {
                    debug!(
                        "backed up primary to old: {} → {}",
                        self.paths.primary_path.display(),
                        self.paths.old_path.display()
                    );
                }
            }
        }

        // Step 2: Move tmp → primary (origin)
        if let Err(e) = move_dir(&tmp_path, &self.paths.primary_path) {
            warn!(
                "failed to move temporary index from {} to {}: {}, attempting rollback",
                tmp_path.display(),
                self.paths.primary_path.display(),
                e
            );
            // Rollback: move backup back to primary
            if self.paths.old_path.exists() {
                return move_dir(&self.paths.old_path, &self.paths.primary_path);
            }
            return Err(e);
        }

        info!(
            "successfully switched index: {} → {} → {}",
            tmp_path.display(),
            self.paths.primary_path.display(),
            self.paths.old_path.display()
        );

        // Step 3: Create new temporary directory
        self.mktmp()?;

        Ok(())
    }
}

/// Moves a directory from source to destination.
///
/// This function copies all contents from source to destination,
/// then removes the source directory.
fn move_dir<P: AsRef<Path>, Q: AsRef<Path>>(src: P, dst: Q) -> Result<(), PersistenceError> {
    let src = src.as_ref();
    let dst = dst.as_ref();

    // Create destination directory
    fs::create_dir_all(dst)?;

    // Copy all files/directories
    for entry in fs::read_dir(src)? {
        let entry = entry?;
        let src_path = entry.path();
        let dst_path = dst.join(entry.file_name());

        if src_path.is_dir() {
            move_dir(&src_path, &dst_path)?;
        } else {
            fs::copy(&src_path, &dst_path)?;
        }
    }

    // Remove source directory
    fs::remove_dir_all(src)?;

    Ok(())
}

/// Copies a directory from source to destination.
fn copy_dir<P: AsRef<Path>, Q: AsRef<Path>>(src: P, dst: Q) -> Result<(), PersistenceError> {
    let src = src.as_ref();
    let dst = dst.as_ref();

    // Create destination directory
    fs::create_dir_all(dst)?;

    // Copy all files/directories
    for entry in fs::read_dir(src)? {
        let entry = entry?;
        let src_path = entry.path();
        let dst_path = dst.join(entry.file_name());

        if src_path.is_dir() {
            copy_dir(&src_path, &dst_path)?;
        } else {
            fs::copy(&src_path, &dst_path)?;
        }
    }

    Ok(())
}

#[cfg(test)]
mod tests {
    use super::*;
    use tempfile::tempdir;

    #[test]
    fn test_index_paths_new() {
        let paths = IndexPaths::new("/data/index");
        assert_eq!(paths.base_path, PathBuf::from("/data/index"));
        assert_eq!(paths.primary_path, PathBuf::from("/data/index/origin"));
        assert_eq!(paths.old_path, PathBuf::from("/data/index/backup"));
        assert_eq!(paths.broken_path, PathBuf::from("/data/index/broken"));
    }

    #[test]
    fn test_persistence_manager_prepare_folders() {
        let dir = tempdir().unwrap();
        let manager = PersistenceManager::new(dir.path(), PersistenceConfig::default());

        manager.prepare_folders().unwrap();

        // base_path should exist (not primary_path, which is created by the index library)
        assert!(manager.paths.base_path.exists());
        assert!(manager.paths.broken_path.exists());
        // old_path not created when CoW is disabled
        assert!(!manager.paths.old_path.exists());
        // primary_path is NOT created by prepare_folders (index library creates it)
        assert!(!manager.paths.primary_path.exists());
    }

    #[test]
    fn test_persistence_manager_prepare_folders_cow() {
        let dir = tempdir().unwrap();
        let config = PersistenceConfig {
            enable_copy_on_write: true,
            ..Default::default()
        };
        let manager = PersistenceManager::new(dir.path(), config);

        manager.prepare_folders().unwrap();

        assert!(manager.paths.base_path.exists());
        assert!(manager.paths.broken_path.exists());
        assert!(manager.paths.old_path.exists());
    }

    #[test]
    fn test_needs_backup_empty_dir() {
        let dir = tempdir().unwrap();
        assert!(!PersistenceManager::needs_backup(dir.path()));
    }

    #[test]
    fn test_needs_backup_with_data_files() {
        let dir = tempdir().unwrap();

        // Create a .kvsdb file (indicates data exists)
        std::fs::write(dir.path().join("test.kvsdb"), b"data").unwrap();

        // No metadata.json -> needs backup
        assert!(PersistenceManager::needs_backup(dir.path()));
    }

    #[test]
    fn test_needs_backup_with_valid_metadata() {
        let dir = tempdir().unwrap();

        // Create data file
        std::fs::write(dir.path().join("test.kvsdb"), b"data").unwrap();

        // Create valid metadata with index_count > 0
        let meta = Metadata::new_qbg(100);
        metadata::store(dir.path().join(AGENT_METADATA_FILENAME), &meta).unwrap();

        // Has data with index_count > 0 -> needs backup
        assert!(PersistenceManager::needs_backup(dir.path()));
    }

    #[test]
    fn test_backup_broken() {
        let dir = tempdir().unwrap();
        let config = PersistenceConfig {
            broken_index_history_limit: 2,
            ..Default::default()
        };
        let manager = PersistenceManager::new(dir.path(), config);

        // Prepare folders first
        manager.prepare_folders().unwrap();

        // Manually create primary path (simulating index library behavior)
        fs::create_dir_all(&manager.paths.primary_path).unwrap();

        // Create some files in the primary path
        std::fs::write(manager.paths.primary_path.join("test.dat"), b"data").unwrap();

        // Backup
        manager.backup_broken().unwrap();

        // Primary path should be recreated but empty
        assert!(manager.paths.primary_path.exists());
        assert_eq!(
            fs::read_dir(&manager.paths.primary_path).unwrap().count(),
            0
        );

        // Broken path should have one backup
        assert_eq!(manager.broken_index_count(), 1);
    }

    #[test]
    fn test_backup_broken_history_limit() {
        let dir = tempdir().unwrap();
        let config = PersistenceConfig {
            broken_index_history_limit: 2,
            ..Default::default()
        };
        let manager = PersistenceManager::new(dir.path(), config);

        manager.prepare_folders().unwrap();

        // Create 3 backups
        for i in 0..3 {
            // Create primary path for each iteration (backup_broken moves it)
            fs::create_dir_all(&manager.paths.primary_path).unwrap();
            std::fs::write(
                manager.paths.primary_path.join(format!("test{}.dat", i)),
                format!("data{}", i).as_bytes(),
            )
            .unwrap();
            manager.backup_broken().unwrap();
            // Small delay to ensure unique timestamps
            std::thread::sleep(std::time::Duration::from_millis(10));
        }

        // Should only have 2 backups (history limit)
        assert_eq!(manager.broken_index_count(), 2);
    }

    #[test]
    fn test_needs_backup_invalid_metadata() {
        let dir = tempdir().unwrap();

        // Create data file
        std::fs::write(dir.path().join("test.kvsdb"), b"data").unwrap();

        // Create invalid metadata
        let meta = Metadata::invalid();
        metadata::store(dir.path().join(AGENT_METADATA_FILENAME), &meta).unwrap();

        // Invalid metadata -> needs backup
        assert!(PersistenceManager::needs_backup(dir.path()));
    }

    #[test]
    fn test_needs_backup_zero_index_count() {
        let dir = tempdir().unwrap();

        // Create data file
        std::fs::write(dir.path().join("test.kvsdb"), b"data").unwrap();

        // Create metadata with index_count = 0
        let meta = Metadata::new_qbg(0);
        metadata::store(dir.path().join(AGENT_METADATA_FILENAME), &meta).unwrap();

        // index_count == 0 -> does NOT need backup (clean state)
        assert!(!PersistenceManager::needs_backup(dir.path()));
    }

    #[test]
    fn test_needs_backup_initial_state_without_data_files() {
        let dir = tempdir().unwrap();

        // Create some non-data files (like grp, obj, prf, tre from NGT)
        std::fs::write(dir.path().join("grp"), b"grp data").unwrap();
        std::fs::write(dir.path().join("obj"), b"obj data").unwrap();

        // No .json or .kvsdb files -> initial state, does NOT need backup
        assert!(!PersistenceManager::needs_backup(dir.path()));
    }

    #[test]
    fn test_backup_broken_empty_primary() {
        let dir = tempdir().unwrap();
        let config = PersistenceConfig {
            broken_index_history_limit: 3,
            ..Default::default()
        };
        let manager = PersistenceManager::new(dir.path(), config);
        manager.prepare_folders().unwrap();

        // Create empty primary path
        fs::create_dir_all(&manager.paths.primary_path).unwrap();

        // Backup should succeed but not create any backup (nothing to backup)
        manager.backup_broken().unwrap();

        // No backups should exist
        assert_eq!(manager.broken_index_count(), 0);
    }

    #[test]
    fn test_backup_broken_history_limit_zero() {
        let dir = tempdir().unwrap();
        let config = PersistenceConfig {
            broken_index_history_limit: 0,
            ..Default::default()
        };
        let manager = PersistenceManager::new(dir.path(), config);
        manager.prepare_folders().unwrap();

        // Create primary path with data
        fs::create_dir_all(&manager.paths.primary_path).unwrap();
        std::fs::write(manager.paths.primary_path.join("test.dat"), b"data").unwrap();

        // Backup should return Ok immediately (history limit is 0)
        manager.backup_broken().unwrap();

        // Primary path should still have data (not moved)
        assert!(manager.paths.primary_path.join("test.dat").exists());

        // No backups should exist
        assert_eq!(manager.broken_index_count(), 0);
    }

    #[test]
    fn test_backup_broken_preserves_newest_backups() {
        let dir = tempdir().unwrap();
        let config = PersistenceConfig {
            broken_index_history_limit: 2,
            ..Default::default()
        };
        let manager = PersistenceManager::new(dir.path(), config);
        manager.prepare_folders().unwrap();

        // Create 3 backups with unique data
        for i in 0..3 {
            fs::create_dir_all(&manager.paths.primary_path).unwrap();
            std::fs::write(
                manager.paths.primary_path.join("data.txt"),
                format!("generation-{}", i),
            )
            .unwrap();
            manager.backup_broken().unwrap();
            std::thread::sleep(std::time::Duration::from_millis(10));
        }

        // Should have 2 backups (newest ones)
        assert_eq!(manager.broken_index_count(), 2);

        // Verify that the oldest backup (generation-0) was removed
        let backups: Vec<_> = fs::read_dir(&manager.paths.broken_path)
            .unwrap()
            .filter_map(|e| e.ok())
            .collect();

        for backup in backups {
            let content = fs::read_to_string(backup.path().join("data.txt")).unwrap();
            // Should NOT contain generation-0
            assert!(
                !content.contains("generation-0"),
                "oldest backup should have been removed"
            );
        }
    }

    #[test]
    fn test_backup_broken_recreates_primary_path() {
        let dir = tempdir().unwrap();
        let config = PersistenceConfig {
            broken_index_history_limit: 3,
            ..Default::default()
        };
        let manager = PersistenceManager::new(dir.path(), config);
        manager.prepare_folders().unwrap();

        // Create primary path with data
        fs::create_dir_all(&manager.paths.primary_path).unwrap();
        std::fs::write(manager.paths.primary_path.join("test.dat"), b"data").unwrap();

        // Backup
        manager.backup_broken().unwrap();

        // Primary path should be recreated (empty directory)
        assert!(manager.paths.primary_path.exists());
        assert!(manager.paths.primary_path.is_dir());
        assert_eq!(
            fs::read_dir(&manager.paths.primary_path).unwrap().count(),
            0
        );
    }

    #[test]
    fn test_index_exists() {
        let dir = tempdir().unwrap();
        let manager = PersistenceManager::new(dir.path(), PersistenceConfig::default());

        // No folder -> doesn't exist
        assert!(!manager.index_exists());

        manager.prepare_folders().unwrap();

        // No metadata -> doesn't exist
        assert!(!manager.index_exists());

        // Create valid metadata
        let meta = Metadata::new_qbg(100);
        manager.save_metadata(&meta).unwrap();

        // Now exists
        assert!(manager.index_exists());
    }

    #[test]
    fn test_index_exists_invalid_metadata() {
        let dir = tempdir().unwrap();
        let manager = PersistenceManager::new(dir.path(), PersistenceConfig::default());

        manager.prepare_folders().unwrap();

        // Create invalid metadata
        let meta = Metadata::invalid();
        manager.save_metadata(&meta).unwrap();

        // Invalid metadata -> doesn't exist
        assert!(!manager.index_exists());
    }

    #[test]
    fn test_load_save_metadata() {
        let dir = tempdir().unwrap();
        let manager = PersistenceManager::new(dir.path(), PersistenceConfig::default());

        manager.prepare_folders().unwrap();

        let original = Metadata::new_qbg(12345);
        manager.save_metadata(&original).unwrap();

        let loaded = manager.load_metadata().unwrap();
        assert_eq!(original, loaded);
    }

    #[test]
    fn test_mktmp_disabled() {
        let dir = tempdir().unwrap();
        let config = PersistenceConfig {
            enable_copy_on_write: false,
            ..Default::default()
        };
        let manager = PersistenceManager::new(dir.path(), config);

        // mktmp should succeed but not create a tmp path when CoW is disabled
        manager.mktmp().unwrap();

        let tmp = manager.tmp_path.read().unwrap();
        assert!(tmp.is_none());
    }

    #[test]
    fn test_mktmp_enabled() {
        let dir = tempdir().unwrap();
        let config = PersistenceConfig {
            enable_copy_on_write: true,
            ..Default::default()
        };
        let manager = PersistenceManager::new(dir.path(), config);

        manager.mktmp().unwrap();

        let tmp = manager.tmp_path.read().unwrap();
        assert!(tmp.is_some());
        let tmp_path = tmp.as_ref().unwrap();
        assert!(tmp_path.exists());
        assert!(tmp_path.starts_with(std::env::temp_dir().join("vald")));
    }

    #[test]
    fn test_mktmp_creates_unique_dirs() {
        let dir = tempdir().unwrap();
        let config = PersistenceConfig {
            enable_copy_on_write: true,
            ..Default::default()
        };
        let manager = PersistenceManager::new(dir.path(), config);

        manager.mktmp().unwrap();
        let first = manager.tmp_path.read().unwrap().clone().unwrap();

        // Small delay to ensure unique timestamp
        std::thread::sleep(std::time::Duration::from_millis(5));

        manager.mktmp().unwrap();
        let second = manager.tmp_path.read().unwrap().clone().unwrap();

        // Paths should be different
        assert_ne!(first, second);

        // Both should exist
        assert!(first.exists());
        assert!(second.exists());

        // Cleanup
        let _ = fs::remove_dir_all(first);
        let _ = fs::remove_dir_all(second);
    }

    #[test]
    fn test_get_save_path_cow_disabled() {
        let dir = tempdir().unwrap();
        let config = PersistenceConfig {
            enable_copy_on_write: false,
            ..Default::default()
        };
        let manager = PersistenceManager::new(dir.path(), config);

        // Should return primary path when CoW is disabled
        let save_path = manager.get_save_path();
        assert_eq!(save_path, manager.paths.primary_path);
    }

    #[test]
    fn test_get_save_path_cow_enabled_no_tmp() {
        let dir = tempdir().unwrap();
        let config = PersistenceConfig {
            enable_copy_on_write: true,
            ..Default::default()
        };
        let manager = PersistenceManager::new(dir.path(), config);

        // When CoW is enabled but mktmp hasn't been called, should return primary path
        let save_path = manager.get_save_path();
        assert_eq!(save_path, manager.paths.primary_path);
    }

    #[test]
    fn test_get_save_path_cow_enabled_with_tmp() {
        let dir = tempdir().unwrap();
        let config = PersistenceConfig {
            enable_copy_on_write: true,
            ..Default::default()
        };
        let manager = PersistenceManager::new(dir.path(), config);

        manager.mktmp().unwrap();

        let save_path = manager.get_save_path();
        let tmp_path = manager.tmp_path.read().unwrap().clone().unwrap();

        assert_eq!(save_path, tmp_path);
        assert_ne!(save_path, manager.paths.primary_path);

        // Cleanup
        let _ = fs::remove_dir_all(tmp_path);
    }

    #[test]
    fn test_save_metadata_to_save_path_cow_disabled() {
        let dir = tempdir().unwrap();
        let config = PersistenceConfig {
            enable_copy_on_write: false,
            ..Default::default()
        };
        let manager = PersistenceManager::new(dir.path(), config);
        manager.prepare_folders().unwrap();

        let meta = Metadata::new_qbg(100);
        manager.save_metadata_to_save_path(&meta).unwrap();

        // Should be saved to primary path
        let saved_path = manager.paths.primary_path.join(AGENT_METADATA_FILENAME);
        assert!(saved_path.exists());

        let loaded = metadata::load(&saved_path).unwrap();
        assert_eq!(meta, loaded);
    }

    #[test]
    fn test_save_metadata_to_save_path_cow_enabled() {
        let dir = tempdir().unwrap();
        let config = PersistenceConfig {
            enable_copy_on_write: true,
            ..Default::default()
        };
        let manager = PersistenceManager::new(dir.path(), config);
        manager.prepare_folders().unwrap();
        manager.mktmp().unwrap();

        let tmp_path = manager.tmp_path.read().unwrap().clone().unwrap();

        let meta = Metadata::new_qbg(200);
        manager.save_metadata_to_save_path(&meta).unwrap();

        // Should be saved to tmp path
        let saved_path = tmp_path.join(AGENT_METADATA_FILENAME);
        assert!(saved_path.exists());

        let loaded = metadata::load(&saved_path).unwrap();
        assert_eq!(meta, loaded);

        // Cleanup
        let _ = fs::remove_dir_all(tmp_path);
    }

    #[test]
    fn test_move_and_switch_saved_data_cow_disabled() {
        let dir = tempdir().unwrap();
        let config = PersistenceConfig {
            enable_copy_on_write: false,
            ..Default::default()
        };
        let manager = PersistenceManager::new(dir.path(), config);

        // Should succeed immediately when CoW is disabled
        manager.move_and_switch_saved_data().unwrap();
    }

    #[test]
    fn test_move_and_switch_saved_data_no_tmp() {
        let dir = tempdir().unwrap();
        let config = PersistenceConfig {
            enable_copy_on_write: true,
            ..Default::default()
        };
        let manager = PersistenceManager::new(dir.path(), config);

        // Should succeed with warning when no tmp path is set
        manager.move_and_switch_saved_data().unwrap();
    }

    #[test]
    fn test_move_and_switch_saved_data_full_cycle() {
        let dir = tempdir().unwrap();
        let config = PersistenceConfig {
            enable_copy_on_write: true,
            ..Default::default()
        };
        let manager = PersistenceManager::new(dir.path(), config);
        manager.prepare_folders().unwrap();

        // Create initial primary data
        fs::create_dir_all(&manager.paths.primary_path).unwrap();
        fs::write(
            manager.paths.primary_path.join("original.dat"),
            b"original data",
        )
        .unwrap();

        // Create temp directory and add new data
        manager.mktmp().unwrap();
        let tmp_path = manager.tmp_path.read().unwrap().clone().unwrap();
        fs::write(tmp_path.join("new.dat"), b"new data").unwrap();

        // Perform the switch
        manager.move_and_switch_saved_data().unwrap();

        // Verify: primary should now contain the new data
        assert!(manager.paths.primary_path.join("new.dat").exists());
        assert!(!manager.paths.primary_path.join("original.dat").exists());

        // Verify: old (backup) should contain the original data
        assert!(manager.paths.old_path.join("original.dat").exists());
        assert!(!manager.paths.old_path.join("new.dat").exists());

        // Verify: new tmp path should be created
        let new_tmp = manager.tmp_path.read().unwrap().clone().unwrap();
        assert!(new_tmp.exists());
        assert_ne!(new_tmp, tmp_path);

        // Cleanup
        let _ = fs::remove_dir_all(new_tmp);
    }

    #[test]
    fn test_move_and_switch_saved_data_empty_primary() {
        let dir = tempdir().unwrap();
        let config = PersistenceConfig {
            enable_copy_on_write: true,
            ..Default::default()
        };
        let manager = PersistenceManager::new(dir.path(), config);
        manager.prepare_folders().unwrap();

        // Create temp directory with data (primary is empty)
        manager.mktmp().unwrap();
        let tmp_path = manager.tmp_path.read().unwrap().clone().unwrap();
        fs::write(tmp_path.join("data.dat"), b"data").unwrap();

        // Perform the switch
        manager.move_and_switch_saved_data().unwrap();

        // Verify: primary should now contain the data
        assert!(manager.paths.primary_path.join("data.dat").exists());

        // Verify: old should be empty or not exist (nothing to backup)
        if manager.paths.old_path.exists() {
            let count = fs::read_dir(&manager.paths.old_path).unwrap().count();
            assert_eq!(count, 0);
        }

        // Cleanup
        let new_tmp = manager.tmp_path.read().unwrap().clone().unwrap();
        let _ = fs::remove_dir_all(new_tmp);
    }

    #[test]
    fn test_move_and_switch_saved_data_replaces_old_backup() {
        let dir = tempdir().unwrap();
        let config = PersistenceConfig {
            enable_copy_on_write: true,
            ..Default::default()
        };
        let manager = PersistenceManager::new(dir.path(), config);
        manager.prepare_folders().unwrap();

        // Create initial old backup
        fs::write(manager.paths.old_path.join("old_backup.dat"), b"old backup").unwrap();

        // Create primary data
        fs::create_dir_all(&manager.paths.primary_path).unwrap();
        fs::write(
            manager.paths.primary_path.join("primary.dat"),
            b"primary data",
        )
        .unwrap();

        // Create temp data
        manager.mktmp().unwrap();
        let tmp_path = manager.tmp_path.read().unwrap().clone().unwrap();
        fs::write(tmp_path.join("new.dat"), b"new data").unwrap();

        // Perform the switch
        manager.move_and_switch_saved_data().unwrap();

        // Verify: old backup should be replaced with primary data
        assert!(manager.paths.old_path.join("primary.dat").exists());
        assert!(!manager.paths.old_path.join("old_backup.dat").exists());

        // Cleanup
        let new_tmp = manager.tmp_path.read().unwrap().clone().unwrap();
        let _ = fs::remove_dir_all(new_tmp);
    }

    #[test]
    fn test_is_copy_on_write_enabled() {
        let dir = tempdir().unwrap();

        let disabled = PersistenceManager::new(
            dir.path(),
            PersistenceConfig {
                enable_copy_on_write: false,
                ..Default::default()
            },
        );
        assert!(!disabled.is_copy_on_write_enabled());

        let enabled = PersistenceManager::new(
            dir.path(),
            PersistenceConfig {
                enable_copy_on_write: true,
                ..Default::default()
            },
        );
        assert!(enabled.is_copy_on_write_enabled());
    }

    #[test]
    fn test_move_dir_helper() {
        let dir = tempdir().unwrap();
        let src = dir.path().join("src");
        let dst = dir.path().join("dst");

        // Create source with nested structure
        fs::create_dir_all(src.join("subdir")).unwrap();
        fs::write(src.join("file1.txt"), b"content1").unwrap();
        fs::write(src.join("subdir/file2.txt"), b"content2").unwrap();

        // Move
        move_dir(&src, &dst).unwrap();

        // Verify source is gone
        assert!(!src.exists());

        // Verify destination has all content
        assert!(dst.join("file1.txt").exists());
        assert!(dst.join("subdir/file2.txt").exists());
        assert_eq!(
            fs::read_to_string(dst.join("file1.txt")).unwrap(),
            "content1"
        );
        assert_eq!(
            fs::read_to_string(dst.join("subdir/file2.txt")).unwrap(),
            "content2"
        );
    }

    #[test]
    fn test_copy_dir_helper() {
        let dir = tempdir().unwrap();
        let src = dir.path().join("src");
        let dst = dir.path().join("dst");

        // Create source with nested structure
        fs::create_dir_all(src.join("subdir")).unwrap();
        fs::write(src.join("file1.txt"), b"content1").unwrap();
        fs::write(src.join("subdir/file2.txt"), b"content2").unwrap();

        // Copy
        copy_dir(&src, &dst).unwrap();

        // Verify source still exists
        assert!(src.exists());
        assert!(src.join("file1.txt").exists());

        // Verify destination has all content
        assert!(dst.join("file1.txt").exists());
        assert!(dst.join("subdir/file2.txt").exists());
        assert_eq!(
            fs::read_to_string(dst.join("file1.txt")).unwrap(),
            "content1"
        );
    }
}
