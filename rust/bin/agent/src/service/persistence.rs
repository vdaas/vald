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

use anyhow::Result;
use std::fs;
use std::path::{Path, PathBuf};

#[derive(Debug)]
pub struct IndexPaths {
    pub primary_path: PathBuf,
    pub secondary_path: PathBuf,
    pub next_path: PathBuf,
    pub temporary_path: PathBuf,
    pub history_path: PathBuf,
    pub broken_path: PathBuf,
}

#[derive(Debug, Clone)]
pub struct PersistenceConfig {
    pub enable_copy_on_write: bool,
    pub broken_index_history_limit: usize,
}

pub struct PersistenceManager {
    #[allow(dead_code)]
    path: PathBuf,
    config: PersistenceConfig,
    paths: IndexPaths,
}

impl PersistenceManager {
    pub fn new<P: AsRef<Path>>(path: P, config: PersistenceConfig) -> Self {
        let path = path.as_ref().to_path_buf();
        let paths = IndexPaths {
            primary_path: path.clone(),
            secondary_path: path.with_file_name(format!(
                "{}_secondary",
                path.file_name().unwrap().to_string_lossy()
            )),
            next_path: path.with_file_name(format!(
                "{}_next",
                path.file_name().unwrap().to_string_lossy()
            )),
            temporary_path: path.with_file_name(format!(
                "{}_tmp",
                path.file_name().unwrap().to_string_lossy()
            )),
            history_path: path.with_file_name(format!(
                "{}_history",
                path.file_name().unwrap().to_string_lossy()
            )),
            broken_path: path.with_file_name(format!(
                "{}_broken",
                path.file_name().unwrap().to_string_lossy()
            )),
        };

        Self {
            path,
            config,
            paths,
        }
    }

    pub fn prepare_folders(&self) -> Result<()> {
        if self.config.enable_copy_on_write {
            // Ensure temporary directory is clean
            if self.paths.temporary_path.exists() {
                fs::remove_dir_all(&self.paths.temporary_path)?;
            }
            fs::create_dir_all(&self.paths.temporary_path)?;
        } else {
            // Ensure parent directory of primary path exists
            if let Some(parent) = self.paths.primary_path.parent() {
                if !parent.exists() {
                    fs::create_dir_all(parent)?;
                }
            }
        }
        Ok(())
    }

    pub fn index_exists(&self) -> bool {
        self.paths.primary_path.join("metadata").exists()
    }

    pub fn broken_index_count(&self) -> u64 {
        // Implementation to count broken indices in history/broken folder
        // For now, return 0 as placeholder
        0
    }

    pub fn paths(&self) -> &IndexPaths {
        &self.paths
    }

    pub fn needs_backup(path: &Path) -> bool {
        // Check if index at path might be broken
        // Simplified check: if folder exists but metadata is missing
        path.exists() && !path.join("metadata").exists()
    }

    pub fn backup_broken(&self) -> Result<()> {
        // Move primary path to broken/history path
        // Simplified implementation
        Ok(())
    }

    pub fn mktmp(&self) -> Result<()> {
        if self.paths.temporary_path.exists() {
            fs::remove_dir_all(&self.paths.temporary_path)?;
        }
        fs::create_dir_all(&self.paths.temporary_path)?;
        Ok(())
    }

    pub fn get_save_path(&self) -> PathBuf {
        if self.config.enable_copy_on_write {
            self.paths.temporary_path.clone()
        } else {
            self.paths.primary_path.clone()
        }
    }

    pub fn save_metadata<T: serde::Serialize>(&self, metadata: &T) -> Result<()> {
        let metadata_path = self.paths.primary_path.join("metadata");
        let file = fs::File::create(metadata_path)?;
        serde_json::to_writer(file, metadata)?;
        Ok(())
    }

    pub fn save_metadata_to_save_path<T: serde::Serialize>(&self, metadata: &T) -> Result<()> {
        let path = self.get_save_path();
        let metadata_path = path.join("metadata");
        let file = fs::File::create(metadata_path)?;
        serde_json::to_writer(file, metadata)?;
        Ok(())
    }

    pub fn is_copy_on_write_enabled(&self) -> bool {
        self.config.enable_copy_on_write
    }

    pub fn move_and_switch_saved_data(&self) -> Result<()> {
        if !self.config.enable_copy_on_write {
            return Ok(());
        }

        // Logic to swap temporary and primary paths atomically (or as close as possible)
        // 1. Rename primary to secondary (backup)
        if self.paths.primary_path.exists() {
            if self.paths.secondary_path.exists() {
                fs::remove_dir_all(&self.paths.secondary_path)?;
            }
            fs::rename(&self.paths.primary_path, &self.paths.secondary_path)?;
        }

        // 2. Rename temporary to primary
        fs::rename(&self.paths.temporary_path, &self.paths.primary_path)?;

        // 3. Remove secondary
        if self.paths.secondary_path.exists() {
            fs::remove_dir_all(&self.paths.secondary_path)?;
        }

        Ok(())
    }
}

#[allow(dead_code)]
fn copy_dir<P: AsRef<Path>, Q: AsRef<Path>>(src: P, dst: Q) -> Result<()> {
    let src = src.as_ref();
    let dst = dst.as_ref();
    fs::create_dir_all(dst)?;
    for entry in fs::read_dir(src)? {
        let entry = entry?;
        let ty = entry.file_type()?;
        if ty.is_dir() {
            copy_dir(entry.path(), dst.join(entry.file_name()))?;
        } else {
            fs::copy(entry.path(), dst.join(entry.file_name()))?;
        }
    }
    Ok(())
}
