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

//! Agent metadata management for index persistence.
//!
//! This module provides functionality to load and store agent metadata
//! which tracks the state of the index (e.g., index count, validity).

use std::fs::{self, File};
use std::io::{BufReader, BufWriter};
use std::path::Path;

use serde::{Deserialize, Serialize};
use thiserror::Error;

/// The filename for agent metadata.
pub const AGENT_METADATA_FILENAME: &str = "metadata.json";

/// Errors that can occur during metadata operations.
#[derive(Debug, Error)]
pub enum MetadataError {
    #[error("metadata file not found: {0}")]
    FileNotFound(String),
    
    #[error("metadata file is empty: {0}")]
    FileEmpty(String),
    
    #[error("failed to read metadata: {0}")]
    ReadError(#[from] std::io::Error),
    
    #[error("failed to parse metadata: {0}")]
    ParseError(#[from] serde_json::Error),
    
    #[error("invalid metadata: {0}")]
    Invalid(String),
}

/// NGT-specific metadata.
#[derive(Debug, Clone, Serialize, Deserialize, Default, PartialEq)]
pub struct NgtMetadata {
    /// The number of indexed vectors.
    pub index_count: u64,
}

/// QBG-specific metadata (same structure as NGT for now).
#[derive(Debug, Clone, Serialize, Deserialize, Default, PartialEq)]
pub struct QbgMetadata {
    /// The number of indexed vectors.
    pub index_count: u64,
}

/// Agent metadata stored alongside the index.
#[derive(Debug, Clone, Serialize, Deserialize, Default, PartialEq)]
pub struct Metadata {
    /// Whether this index is marked as invalid.
    #[serde(default)]
    pub is_invalid: bool,
    
    /// NGT-specific metadata.
    #[serde(skip_serializing_if = "Option::is_none")]
    pub ngt: Option<NgtMetadata>,
    
    /// QBG-specific metadata.
    #[serde(skip_serializing_if = "Option::is_none")]
    pub qbg: Option<QbgMetadata>,
}

impl Metadata {
    /// Creates a new metadata instance for QBG with the given index count.
    pub fn new_qbg(index_count: u64) -> Self {
        Metadata {
            is_invalid: false,
            ngt: None,
            qbg: Some(QbgMetadata { index_count }),
        }
    }
    
    /// Creates a new metadata instance for NGT with the given index count.
    pub fn new_ngt(index_count: u64) -> Self {
        Metadata {
            is_invalid: false,
            ngt: Some(NgtMetadata { index_count }),
            qbg: None,
        }
    }
    
    /// Creates a metadata instance marked as invalid.
    pub fn invalid() -> Self {
        Metadata {
            is_invalid: true,
            ngt: None,
            qbg: None,
        }
    }
    
    /// Returns the index count from either NGT or QBG metadata.
    pub fn index_count(&self) -> u64 {
        self.qbg.as_ref().map(|q| q.index_count)
            .or_else(|| self.ngt.as_ref().map(|n| n.index_count))
            .unwrap_or(0)
    }
}

/// Loads metadata from the specified path.
///
/// # Arguments
/// * `path` - Path to the metadata file (e.g., "index/metadata.json")
///
/// # Returns
/// The loaded metadata or an error if the file cannot be read.
pub fn load<P: AsRef<Path>>(path: P) -> Result<Metadata, MetadataError> {
    let path = path.as_ref();
    
    // Check if file exists
    if !path.exists() {
        return Err(MetadataError::FileNotFound(path.display().to_string()));
    }
    
    // Check if file is empty
    let file_metadata = fs::metadata(path)?;
    if file_metadata.len() == 0 {
        return Err(MetadataError::FileEmpty(path.display().to_string()));
    }
    
    // Open and read the file
    let file = File::open(path)?;
    let reader = BufReader::new(file);
    
    let metadata: Metadata = serde_json::from_reader(reader)?;
    
    Ok(metadata)
}

/// Stores metadata to the specified path.
///
/// # Arguments
/// * `path` - Path to store the metadata file
/// * `metadata` - The metadata to store
///
/// # Returns
/// Ok(()) on success, or an error if the file cannot be written.
pub fn store<P: AsRef<Path>>(path: P, metadata: &Metadata) -> Result<(), MetadataError> {
    let path = path.as_ref();
    
    // Ensure parent directory exists
    if let Some(parent) = path.parent() {
        fs::create_dir_all(parent)?;
    }
    
    // Open file for writing (create or truncate)
    let file = File::create(path)?;
    let writer = BufWriter::new(file);
    
    // Write metadata as JSON
    serde_json::to_writer_pretty(writer, metadata)?;
    
    Ok(())
}

/// Returns the metadata file path for a given index directory.
pub fn metadata_path<P: AsRef<Path>>(index_dir: P) -> std::path::PathBuf {
    index_dir.as_ref().join(AGENT_METADATA_FILENAME)
}

#[cfg(test)]
mod tests {
    use super::*;
    use tempfile::tempdir;

    #[test]
    fn test_metadata_new_qbg() {
        let meta = Metadata::new_qbg(1000);
        assert!(!meta.is_invalid);
        assert!(meta.ngt.is_none());
        assert_eq!(meta.qbg.as_ref().unwrap().index_count, 1000);
        assert_eq!(meta.index_count(), 1000);
    }

    #[test]
    fn test_metadata_new_ngt() {
        let meta = Metadata::new_ngt(500);
        assert!(!meta.is_invalid);
        assert!(meta.qbg.is_none());
        assert_eq!(meta.ngt.as_ref().unwrap().index_count, 500);
        assert_eq!(meta.index_count(), 500);
    }

    #[test]
    fn test_metadata_invalid() {
        let meta = Metadata::invalid();
        assert!(meta.is_invalid);
        assert_eq!(meta.index_count(), 0);
    }

    #[test]
    fn test_store_and_load() {
        let dir = tempdir().unwrap();
        let path = dir.path().join("metadata.json");
        
        let original = Metadata::new_qbg(12345);
        store(&path, &original).unwrap();
        
        let loaded = load(&path).unwrap();
        assert_eq!(original, loaded);
    }

    #[test]
    fn test_load_nonexistent() {
        let result = load("/nonexistent/path/metadata.json");
        assert!(matches!(result, Err(MetadataError::FileNotFound(_))));
    }

    #[test]
    fn test_load_empty_file() {
        let dir = tempdir().unwrap();
        let path = dir.path().join("empty.json");
        
        // Create empty file
        File::create(&path).unwrap();
        
        let result = load(&path);
        assert!(matches!(result, Err(MetadataError::FileEmpty(_))));
    }

    #[test]
    fn test_json_serialization() {
        let meta = Metadata::new_qbg(100);
        let json = serde_json::to_string_pretty(&meta).unwrap();
        
        // Verify it matches the Go format
        assert!(json.contains("\"is_invalid\": false"));
        assert!(json.contains("\"index_count\": 100"));
    }

    #[test]
    fn test_metadata_path() {
        let path = metadata_path("/data/index");
        assert_eq!(path.to_str().unwrap(), "/data/index/metadata.json");
    }
}
