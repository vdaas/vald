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

use anyhow::{Context, Result};
use kube::{
    api::{Api, Patch, PatchParams},
    Client,
};
use k8s_openapi::api::core::v1::Pod;
use serde_json::json;
use std::collections::HashMap;
use tracing::{debug, error, info};

/// Annotation keys for exporting index metrics to pod annotations.
pub mod annotations {
    /// Annotation key for index count.
    pub const INDEX_COUNT: &str = "vald.vdaas.org/index-count";
    /// Annotation key for uncommitted entry count.
    pub const UNCOMMITTED_COUNT: &str = "vald.vdaas.org/uncommitted-entries";
    /// Annotation key for processed vqueue entries.
    pub const PROCESSED_VQ_COUNT: &str = "vald.vdaas.org/processed-vq-entries";
    /// Annotation key for last save index timestamp.
    pub const LAST_SAVE_TIMESTAMP: &str = "vald.vdaas.org/last-save-timestamp";
    /// Annotation key for unsaved create index execution count.
    pub const UNSAVED_CREATE_INDEX_EXEC: &str = "vald.vdaas.org/unsaved-create-index-execution";
}

/// Trait for applying annotations to Kubernetes resources.
#[async_trait::async_trait]
pub trait Patcher: Send + Sync {
    /// Apply annotations to the specified pod.
    async fn apply_pod_annotations(
        &self,
        name: &str,
        namespace: &str,
        annotations: HashMap<String, String>,
    ) -> Result<()>;
}

/// Kubernetes client for interacting with the Kubernetes API.
pub struct K8sClient {
    client: Client,
}

impl K8sClient {
    /// Create a new K8sClient.
    pub async fn new() -> Result<Self> {
        let client = Client::try_default()
            .await
            .context("failed to create Kubernetes client")?;
        Ok(Self { client })
    }

    /// Create a new K8sClient with a custom client.
    pub fn with_client(client: Client) -> Self {
        Self { client }
    }
}

#[async_trait::async_trait]
impl Patcher for K8sClient {
    async fn apply_pod_annotations(
        &self,
        name: &str,
        namespace: &str,
        annotations: HashMap<String, String>,
    ) -> Result<()> {
        if annotations.is_empty() {
            debug!("no annotations to apply, skipping");
            return Ok(());
        }

        let pods: Api<Pod> = Api::namespaced(self.client.clone(), namespace);

        // Build the patch using server-side apply
        let patch = json!({
            "apiVersion": "v1",
            "kind": "Pod",
            "metadata": {
                "name": name,
                "annotations": annotations
            }
        });

        let params = PatchParams::apply("vald-agent").force();

        match pods.patch(name, &params, &Patch::Apply(&patch)).await {
            Ok(_) => {
                debug!(
                    "successfully applied annotations to pod {}/{}: {:?}",
                    namespace, name, annotations
                );
                Ok(())
            }
            Err(e) => {
                error!(
                    "failed to apply annotations to pod {}/{}: {}",
                    namespace, name, e
                );
                Err(e.into())
            }
        }
    }
}

/// Index metrics for exporting to pod annotations.
#[derive(Debug, Clone, Default)]
pub struct IndexMetrics {
    /// Number of indexed vectors.
    pub index_count: Option<u64>,
    /// Number of uncommitted entries.
    pub uncommitted_count: Option<u64>,
    /// Number of processed vqueue entries.
    pub processed_vq_count: Option<u64>,
    /// Last save index timestamp in RFC3339 format.
    pub last_save_timestamp: Option<String>,
    /// Number of create index executions since last save.
    pub unsaved_create_index_exec: Option<u64>,
}

impl IndexMetrics {
    /// Convert metrics to annotation map.
    pub fn to_annotations(&self) -> HashMap<String, String> {
        let mut annotations = HashMap::new();

        if let Some(v) = self.index_count {
            annotations.insert(annotations::INDEX_COUNT.to_string(), v.to_string());
        }
        if let Some(v) = self.uncommitted_count {
            annotations.insert(annotations::UNCOMMITTED_COUNT.to_string(), v.to_string());
        }
        if let Some(v) = self.processed_vq_count {
            annotations.insert(annotations::PROCESSED_VQ_COUNT.to_string(), v.to_string());
        }
        if let Some(v) = &self.last_save_timestamp {
            annotations.insert(annotations::LAST_SAVE_TIMESTAMP.to_string(), v.clone());
        }
        if let Some(v) = self.unsaved_create_index_exec {
            annotations.insert(annotations::UNSAVED_CREATE_INDEX_EXEC.to_string(), v.to_string());
        }

        annotations
    }
}

/// Manager for exporting index metrics to Kubernetes pod annotations.
pub struct MetricsExporter {
    patcher: Box<dyn Patcher>,
    pod_name: String,
    pod_namespace: String,
    enabled: bool,
}

impl MetricsExporter {
    /// Create a new MetricsExporter.
    pub fn new(
        patcher: Box<dyn Patcher>,
        pod_name: String,
        pod_namespace: String,
        enabled: bool,
    ) -> Self {
        Self {
            patcher,
            pod_name,
            pod_namespace,
            enabled,
        }
    }

    /// Check if export is enabled.
    pub fn is_enabled(&self) -> bool {
        self.enabled
    }

    /// Export metrics for tick event.
    /// Exports: uncommitted_count, index_count
    pub async fn export_on_tick(&self, index_count: u64, uncommitted_count: u64) -> Result<()> {
        if !self.enabled {
            return Ok(());
        }

        let metrics = IndexMetrics {
            index_count: Some(index_count),
            uncommitted_count: Some(uncommitted_count),
            ..Default::default()
        };

        info!(
            "exporting tick metrics: index_count={}, uncommitted_count={}",
            index_count, uncommitted_count
        );

        self.patcher
            .apply_pod_annotations(&self.pod_name, &self.pod_namespace, metrics.to_annotations())
            .await
    }

    /// Export metrics after create_index operation.
    /// Exports: uncommitted_count, processed_vq_count, unsaved_create_index_exec, index_count
    pub async fn export_on_create_index(
        &self,
        index_count: u64,
        uncommitted_count: u64,
        processed_vq_count: u64,
        unsaved_create_index_exec: u64,
    ) -> Result<()> {
        if !self.enabled {
            return Ok(());
        }

        let metrics = IndexMetrics {
            index_count: Some(index_count),
            uncommitted_count: Some(uncommitted_count),
            processed_vq_count: Some(processed_vq_count),
            unsaved_create_index_exec: Some(unsaved_create_index_exec),
            ..Default::default()
        };

        info!(
            "exporting create_index metrics: index_count={}, uncommitted_count={}, processed_vq={}, unsaved_exec={}",
            index_count, uncommitted_count, processed_vq_count, unsaved_create_index_exec
        );

        self.patcher
            .apply_pod_annotations(&self.pod_name, &self.pod_namespace, metrics.to_annotations())
            .await
    }

    /// Export metrics after save_index operation.
    /// Exports: last_save_timestamp, unsaved_create_index_exec, processed_vq_count
    pub async fn export_on_save_index(
        &self,
        last_save_timestamp: String,
        processed_vq_count: u64,
    ) -> Result<()> {
        if !self.enabled {
            return Ok(());
        }

        let metrics = IndexMetrics {
            last_save_timestamp: Some(last_save_timestamp.clone()),
            processed_vq_count: Some(processed_vq_count),
            unsaved_create_index_exec: Some(0), // Reset after save
            ..Default::default()
        };

        info!(
            "exporting save_index metrics: timestamp={}, processed_vq={}",
            last_save_timestamp, processed_vq_count
        );

        self.patcher
            .apply_pod_annotations(&self.pod_name, &self.pod_namespace, metrics.to_annotations())
            .await
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use std::sync::Mutex;

    struct MockPatcher {
        applied: Mutex<Vec<HashMap<String, String>>>,
    }

    impl MockPatcher {
        fn new() -> Self {
            Self {
                applied: Mutex::new(Vec::new()),
            }
        }

        fn get_applied(&self) -> Vec<HashMap<String, String>> {
            self.applied.lock().unwrap().clone()
        }
    }

    #[async_trait::async_trait]
    impl Patcher for MockPatcher {
        async fn apply_pod_annotations(
            &self,
            _name: &str,
            _namespace: &str,
            annotations: HashMap<String, String>,
        ) -> Result<()> {
            self.applied.lock().unwrap().push(annotations);
            Ok(())
        }
    }

    #[tokio::test]
    async fn test_export_on_tick() {
        let patcher = Box::new(MockPatcher::new());
        let exporter = MetricsExporter::new(
            patcher,
            "test-pod".to_string(),
            "default".to_string(),
            true,
        );

        exporter.export_on_tick(100, 5).await.unwrap();

        // Note: Can't access mock directly due to Box, would need Arc<dyn Patcher>
    }

    #[test]
    fn test_index_metrics_to_annotations() {
        let metrics = IndexMetrics {
            index_count: Some(100),
            uncommitted_count: Some(5),
            processed_vq_count: Some(10),
            last_save_timestamp: Some("2024-01-01T00:00:00Z".to_string()),
            unsaved_create_index_exec: Some(2),
        };

        let annotations = metrics.to_annotations();
        assert_eq!(annotations.get(annotations::INDEX_COUNT), Some(&"100".to_string()));
        assert_eq!(annotations.get(annotations::UNCOMMITTED_COUNT), Some(&"5".to_string()));
        assert_eq!(annotations.get(annotations::PROCESSED_VQ_COUNT), Some(&"10".to_string()));
        assert_eq!(
            annotations.get(annotations::LAST_SAVE_TIMESTAMP),
            Some(&"2024-01-01T00:00:00Z".to_string())
        );
        assert_eq!(
            annotations.get(annotations::UNSAVED_CREATE_INDEX_EXEC),
            Some(&"2".to_string())
        );
    }

    #[test]
    fn test_index_metrics_partial() {
        let metrics = IndexMetrics {
            index_count: Some(50),
            ..Default::default()
        };

        let annotations = metrics.to_annotations();
        assert_eq!(annotations.len(), 1);
        assert_eq!(annotations.get(annotations::INDEX_COUNT), Some(&"50".to_string()));
    }
}
