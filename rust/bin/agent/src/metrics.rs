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

use algorithm::ANN;
use opentelemetry::global;
use std::sync::Arc;
use tokio::sync::RwLock;

// Metric names
const INDEX_COUNT: &str = "agent_core_ngt_index_count";
const UNCOMMITTED_INDEX_COUNT: &str = "agent_core_ngt_uncommitted_index_count";
const INSERT_VQUEUE_COUNT: &str = "agent_core_ngt_insert_vqueue_count";
const DELETE_VQUEUE_COUNT: &str = "agent_core_ngt_delete_vqueue_count";
const COMPLETED_CREATE_INDEX_TOTAL: &str = "agent_core_ngt_completed_create_index_total";
const EXECUTED_PROACTIVE_GC_TOTAL: &str = "agent_core_ngt_executed_proactive_gc_total";
const IS_INDEXING: &str = "agent_core_ngt_is_indexing";
const IS_SAVING: &str = "agent_core_ngt_is_saving";
const BROKEN_INDEX_STORE_COUNT: &str = "agent_core_ngt_broken_index_store_count";

// Statistic metric names
const MEDIAN_INDEGREE: &str = "agent_core_ngt_median_indegree";
const MEDIAN_OUTDEGREE: &str = "agent_core_ngt_median_outdegree";
const MAX_NUMBER_OF_INDEGREE: &str = "agent_core_ngt_max_number_of_indegree";
const MAX_NUMBER_OF_OUTDEGREE: &str = "agent_core_ngt_max_number_of_outdegree";
const MIN_NUMBER_OF_INDEGREE: &str = "agent_core_ngt_min_number_of_indegree";
const MIN_NUMBER_OF_OUTDEGREE: &str = "agent_core_ngt_min_number_of_outdegree";
const MODE_INDEGREE: &str = "agent_core_ngt_mode_indegree";
const MODE_OUTDEGREE: &str = "agent_core_ngt_mode_outdegree";
const NODES_SKIPPED_FOR_10_EDGES: &str = "agent_core_ngt_nodes_skipped_for_10_edges";
const NODES_SKIPPED_FOR_INDEGREE_DISTANCE: &str =
    "agent_core_ngt_nodes_skipped_for_indegree_distance";
const NUMBER_OF_EDGES: &str = "agent_core_ngt_number_of_edges";
const NUMBER_OF_INDEXED_OBJECTS: &str = "agent_core_ngt_number_of_indexed_objects";
const NUMBER_OF_NODES: &str = "agent_core_ngt_number_of_nodes";
const NUMBER_OF_NODES_WITHOUT_EDGES: &str = "agent_core_ngt_number_of_nodes_without_edges";
const NUMBER_OF_NODES_WITHOUT_INDEGREE: &str = "agent_core_ngt_number_of_nodes_without_indegree";
const NUMBER_OF_OBJECTS: &str = "agent_core_ngt_number_of_objects";
const NUMBER_OF_REMOVED_OBJECTS: &str = "agent_core_ngt_number_of_removed_objects";
const SIZE_OF_OBJECT_REPOSITORY: &str = "agent_core_ngt_size_of_object_repository";
const SIZE_OF_REFINEMENT_OBJECT_REPOSITORY: &str =
    "agent_core_ngt_size_of_refinement_object_repository";
const VARIANCE_OF_INDEGREE: &str = "agent_core_ngt_variance_of_indegree";
const VARIANCE_OF_OUTDEGREE: &str = "agent_core_ngt_variance_of_outdegree";
const MEAN_EDGE_LENGTH: &str = "agent_core_ngt_mean_edge_length";
const MEAN_EDGE_LENGTH_FOR_10_EDGES: &str = "agent_core_ngt_mean_edge_length_for_10_edges";
const MEAN_INDEGREE_DISTANCE_FOR_10_EDGES: &str =
    "agent_core_ngt_mean_indegree_distance_for_10_edges";
const MEAN_NUMBER_OF_EDGES_PER_NODE: &str = "agent_core_ngt_mean_number_of_edges_per_node";
const C1_INDEGREE: &str = "agent_core_ngt_c1_indegree";
const C5_INDEGREE: &str = "agent_core_ngt_c5_indegree";
const C95_OUTDEGREE: &str = "agent_core_ngt_c95_outdegree";
const C99_OUTDEGREE: &str = "agent_core_ngt_c99_outdegree";

/// Registers an i64 observable gauge that reads a value from the ANN service.
macro_rules! register_basic_gauge {
    ($meter:expr, $svc:expr, $name:expr, $desc:expr, |$s:ident| $value:expr) => {{
        let svc = $svc.clone();
        $meter
            .i64_observable_gauge($name)
            .with_description($desc)
            .with_callback(move |observer| {
                if let Some(service) = svc.upgrade()
                    && let Ok($s) = service.try_read()
                {
                    observer.observe($value, &[]);
                }
            })
            .build();
    }};
}

/// Registers an i64 observable gauge backed by a field from `index_statistics()`.
macro_rules! register_stats_gauge_i64 {
    ($meter:expr, $svc:expr, $name:expr, $desc:expr, $field:ident) => {{
        let svc = $svc.clone();
        $meter
            .i64_observable_gauge($name)
            .with_description($desc)
            .with_callback(move |observer| {
                if let Some(service) = svc.upgrade()
                    && let Ok(s) = service.try_read()
                    && s.is_statistics_enabled()
                    && let Ok(stats) = s.index_statistics()
                {
                    observer.observe(stats.$field as i64, &[]);
                }
            })
            .build();
    }};
}

/// Registers an f64 observable gauge backed by a field from `index_statistics()`.
macro_rules! register_stats_gauge_f64 {
    ($meter:expr, $svc:expr, $name:expr, $desc:expr, $field:ident) => {{
        let svc = $svc.clone();
        $meter
            .f64_observable_gauge($name)
            .with_description($desc)
            .with_callback(move |observer| {
                if let Some(service) = svc.upgrade()
                    && let Ok(s) = service.try_read()
                    && s.is_statistics_enabled()
                    && let Ok(stats) = s.index_statistics()
                {
                    observer.observe(stats.$field, &[]);
                }
            })
            .build();
    }};
}

/// Registers OpenTelemetry metrics backed by the ANN service state.
pub fn register_metrics<S>(service: Arc<RwLock<S>>) -> anyhow::Result<()>
where
    S: ANN + 'static,
{
    let meter = global::meter("vald-agent");
    let svc = Arc::downgrade(&service);

    // Basic Metrics
    register_basic_gauge!(
        meter,
        svc,
        INDEX_COUNT,
        "Agent NGT index count",
        |s| s.len() as i64
    );
    register_basic_gauge!(
        meter,
        svc,
        UNCOMMITTED_INDEX_COUNT,
        "Agent NGT uncommitted index count",
        |s| { (s.insert_vqueue_buffer_len() + s.delete_vqueue_buffer_len()) as i64 }
    );
    register_basic_gauge!(
        meter,
        svc,
        INSERT_VQUEUE_COUNT,
        "Agent NGT insert vqueue count",
        |s| s.insert_vqueue_buffer_len() as i64
    );
    register_basic_gauge!(
        meter,
        svc,
        DELETE_VQUEUE_COUNT,
        "Agent NGT delete vqueue count",
        |s| s.delete_vqueue_buffer_len() as i64
    );
    register_basic_gauge!(
        meter,
        svc,
        COMPLETED_CREATE_INDEX_TOTAL,
        "The cumulative count of completed create index execution",
        |s| s.number_of_create_index_executions() as i64
    );
    meter
        .i64_observable_gauge(EXECUTED_PROACTIVE_GC_TOTAL)
        .with_description("The cumulative count of proactive GC execution")
        .with_callback(|observer| {
            observer.observe(0_i64, &[]);
        })
        .build();
    register_basic_gauge!(
        meter,
        svc,
        IS_INDEXING,
        "Currently indexing or no",
        |s| if s.is_indexing() { 1 } else { 0 }
    );
    register_basic_gauge!(
        meter,
        svc,
        IS_SAVING,
        "Currently saving or not",
        |s| if s.is_saving() { 1 } else { 0 }
    );
    register_basic_gauge!(
        meter,
        svc,
        BROKEN_INDEX_STORE_COUNT,
        "How many broken index generations have been stored",
        |s| s.broken_index_count() as i64
    );

    // Statistics Metrics (Int64)
    register_stats_gauge_i64!(
        meter,
        svc,
        MEDIAN_INDEGREE,
        "Median indegree of nodes",
        median_indegree
    );
    register_stats_gauge_i64!(
        meter,
        svc,
        MEDIAN_OUTDEGREE,
        "Median outdegree of nodes",
        median_outdegree
    );
    register_stats_gauge_i64!(
        meter,
        svc,
        MAX_NUMBER_OF_INDEGREE,
        "Maximum number of indegree",
        max_number_of_indegree
    );
    register_stats_gauge_i64!(
        meter,
        svc,
        MAX_NUMBER_OF_OUTDEGREE,
        "Maximum number of outdegree",
        max_number_of_outdegree
    );
    register_stats_gauge_i64!(
        meter,
        svc,
        MIN_NUMBER_OF_INDEGREE,
        "Minimum number of indegree",
        min_number_of_indegree
    );
    register_stats_gauge_i64!(
        meter,
        svc,
        MIN_NUMBER_OF_OUTDEGREE,
        "Minimum number of outdegree",
        min_number_of_outdegree
    );
    register_stats_gauge_i64!(meter, svc, MODE_INDEGREE, "Mode of indegree", mode_indegree);
    register_stats_gauge_i64!(
        meter,
        svc,
        MODE_OUTDEGREE,
        "Mode of outdegree",
        mode_outdegree
    );
    register_stats_gauge_i64!(
        meter,
        svc,
        NODES_SKIPPED_FOR_10_EDGES,
        "Nodes skipped for 10 edges",
        nodes_skipped_for_10_edges
    );
    register_stats_gauge_i64!(
        meter,
        svc,
        NODES_SKIPPED_FOR_INDEGREE_DISTANCE,
        "Nodes skipped for indegree distance",
        nodes_skipped_for_indegree_distance
    );
    register_stats_gauge_i64!(
        meter,
        svc,
        NUMBER_OF_EDGES,
        "Number of edges",
        number_of_edges
    );
    register_stats_gauge_i64!(
        meter,
        svc,
        NUMBER_OF_INDEXED_OBJECTS,
        "Number of indexed objects",
        number_of_indexed_objects
    );
    register_stats_gauge_i64!(
        meter,
        svc,
        NUMBER_OF_NODES,
        "Number of nodes",
        number_of_nodes
    );
    register_stats_gauge_i64!(
        meter,
        svc,
        NUMBER_OF_NODES_WITHOUT_EDGES,
        "Number of nodes without edges",
        number_of_nodes_without_edges
    );
    register_stats_gauge_i64!(
        meter,
        svc,
        NUMBER_OF_NODES_WITHOUT_INDEGREE,
        "Number of nodes without indegree",
        number_of_nodes_without_indegree
    );
    register_stats_gauge_i64!(
        meter,
        svc,
        NUMBER_OF_OBJECTS,
        "Number of objects",
        number_of_objects
    );
    register_stats_gauge_i64!(
        meter,
        svc,
        NUMBER_OF_REMOVED_OBJECTS,
        "Number of removed objects",
        number_of_removed_objects
    );
    register_stats_gauge_i64!(
        meter,
        svc,
        SIZE_OF_OBJECT_REPOSITORY,
        "Size of object repository",
        size_of_object_repository
    );
    register_stats_gauge_i64!(
        meter,
        svc,
        SIZE_OF_REFINEMENT_OBJECT_REPOSITORY,
        "Size of refinement object repository",
        size_of_refinement_object_repository
    );

    // Statistics Metrics (Float64)
    register_stats_gauge_f64!(
        meter,
        svc,
        VARIANCE_OF_INDEGREE,
        "Variance of indegree",
        variance_of_indegree
    );
    register_stats_gauge_f64!(
        meter,
        svc,
        VARIANCE_OF_OUTDEGREE,
        "Variance of outdegree",
        variance_of_outdegree
    );
    register_stats_gauge_f64!(
        meter,
        svc,
        MEAN_EDGE_LENGTH,
        "Mean edge length",
        mean_edge_length
    );
    register_stats_gauge_f64!(
        meter,
        svc,
        MEAN_EDGE_LENGTH_FOR_10_EDGES,
        "Mean edge length for 10 edges",
        mean_edge_length_for_10_edges
    );
    register_stats_gauge_f64!(
        meter,
        svc,
        MEAN_INDEGREE_DISTANCE_FOR_10_EDGES,
        "Mean indegree distance for 10 edges",
        mean_indegree_distance_for_10_edges
    );
    register_stats_gauge_f64!(
        meter,
        svc,
        MEAN_NUMBER_OF_EDGES_PER_NODE,
        "Mean number of edges per node",
        mean_number_of_edges_per_node
    );
    register_stats_gauge_f64!(meter, svc, C1_INDEGREE, "C1 indegree", c1_indegree);
    register_stats_gauge_f64!(meter, svc, C5_INDEGREE, "C5 indegree", c5_indegree);
    register_stats_gauge_f64!(meter, svc, C95_OUTDEGREE, "C95 outdegree", c95_outdegree);
    register_stats_gauge_f64!(meter, svc, C99_OUTDEGREE, "C99 outdegree", c99_outdegree);

    Ok(())
}

#[cfg(test)]
mod tests {
    use super::*;
    use algorithm::{ANN, Error};
    use opentelemetry_sdk::metrics::{InMemoryMetricExporter, PeriodicReader, SdkMeterProvider};
    use proto::payload::v1::{info, search};
    use std::collections::HashMap;
    use std::future::Future;

    #[derive(Clone)]
    struct MockANN {
        len: u32,
        insert_buffer: u32,
        delete_buffer: u32,
        create_index_count: u64,
        indexing: bool,
        saving: bool,
        broken_count: u64,
        stats_enabled: bool,
    }

    impl MockANN {
        fn new() -> Self {
            Self {
                len: 100,
                insert_buffer: 10,
                delete_buffer: 5,
                create_index_count: 3,
                indexing: true,
                saving: false,
                broken_count: 1,
                stats_enabled: true,
            }
        }
    }

    impl ANN for MockANN {
        fn search(
            &self,
            _v: Vec<f32>,
            _k: u32,
            _e: f32,
            _r: f32,
        ) -> impl Future<Output = Result<search::Response, Error>> + Send {
            async { Ok(search::Response::default()) }
        }
        fn search_by_id(
            &self,
            _u: String,
            _k: u32,
            _e: f32,
            _r: f32,
        ) -> impl Future<Output = Result<search::Response, Error>> + Send {
            async { Ok(search::Response::default()) }
        }
        fn linear_search(
            &self,
            _v: Vec<f32>,
            _k: u32,
        ) -> impl Future<Output = Result<search::Response, Error>> + Send {
            async { Ok(search::Response::default()) }
        }
        fn linear_search_by_id(
            &self,
            _u: String,
            _k: u32,
        ) -> impl Future<Output = Result<search::Response, Error>> + Send {
            async { Ok(search::Response::default()) }
        }
        fn insert(
            &mut self,
            _u: String,
            _v: Vec<f32>,
        ) -> impl Future<Output = Result<(), Error>> + Send {
            async { Ok(()) }
        }
        fn insert_with_time(
            &mut self,
            _u: String,
            _v: Vec<f32>,
            _t: i64,
        ) -> impl Future<Output = Result<(), Error>> + Send {
            async { Ok(()) }
        }
        fn insert_multiple(
            &mut self,
            _vs: HashMap<String, Vec<f32>>,
        ) -> impl Future<Output = Result<(), Error>> + Send {
            async { Ok(()) }
        }
        fn insert_multiple_with_time(
            &mut self,
            _vs: HashMap<String, Vec<f32>>,
            _t: i64,
        ) -> impl Future<Output = Result<(), Error>> + Send {
            async { Ok(()) }
        }
        fn update(
            &mut self,
            _u: String,
            _v: Vec<f32>,
        ) -> impl Future<Output = Result<(), Error>> + Send {
            async { Ok(()) }
        }
        fn update_with_time(
            &mut self,
            _u: String,
            _v: Vec<f32>,
            _t: i64,
        ) -> impl Future<Output = Result<(), Error>> + Send {
            async { Ok(()) }
        }
        fn update_multiple(
            &mut self,
            _vs: HashMap<String, Vec<f32>>,
        ) -> impl Future<Output = Result<(), Error>> + Send {
            async { Ok(()) }
        }
        fn update_multiple_with_time(
            &mut self,
            _vs: HashMap<String, Vec<f32>>,
            _t: i64,
        ) -> impl Future<Output = Result<(), Error>> + Send {
            async { Ok(()) }
        }
        fn update_timestamp(
            &mut self,
            _u: String,
            _t: i64,
            _f: bool,
        ) -> impl Future<Output = Result<(), Error>> + Send {
            async { Ok(()) }
        }
        fn remove(&mut self, _u: String) -> impl Future<Output = Result<(), Error>> + Send {
            async { Ok(()) }
        }
        fn remove_with_time(
            &mut self,
            _u: String,
            _t: i64,
        ) -> impl Future<Output = Result<(), Error>> + Send {
            async { Ok(()) }
        }
        fn remove_multiple(
            &mut self,
            _us: Vec<String>,
        ) -> impl Future<Output = Result<(), Error>> + Send {
            async { Ok(()) }
        }
        fn remove_multiple_with_time(
            &mut self,
            _us: Vec<String>,
            _t: i64,
        ) -> impl Future<Output = Result<(), Error>> + Send {
            async { Ok(()) }
        }
        fn regenerate_indexes(&mut self) -> impl Future<Output = Result<(), Error>> + Send {
            async { Ok(()) }
        }
        fn create_index(&mut self) -> impl Future<Output = Result<(), Error>> + Send {
            async { Ok(()) }
        }
        fn save_index(&mut self) -> impl Future<Output = Result<(), Error>> + Send {
            async { Ok(()) }
        }
        fn create_and_save_index(&mut self) -> impl Future<Output = Result<(), Error>> + Send {
            async { Ok(()) }
        }
        fn get_object(
            &self,
            _u: String,
        ) -> impl Future<Output = Result<(Vec<f32>, i64), Error>> + Send {
            async { Ok((vec![], 0)) }
        }
        fn exists(&self, _u: String) -> impl Future<Output = (usize, bool)> + Send {
            async { (0, false) }
        }
        fn uuids(&self) -> impl Future<Output = Vec<String>> + Send {
            async { vec![] }
        }
        fn list_object_func<F: FnMut(String, Vec<f32>, i64) -> bool + Send>(
            &self,
            _f: F,
        ) -> impl Future<Output = ()> + Send {
            async {}
        }
        fn close(&mut self) -> impl Future<Output = Result<(), Error>> + Send {
            async { Ok(()) }
        }

        // Metrics methods
        fn is_indexing(&self) -> bool {
            self.indexing
        }
        fn is_flushing(&self) -> bool {
            false
        }
        fn is_saving(&self) -> bool {
            self.saving
        }
        fn len(&self) -> u32 {
            self.len
        }
        fn number_of_create_index_executions(&self) -> u64 {
            self.create_index_count
        }
        fn insert_vqueue_buffer_len(&self) -> u32 {
            self.insert_buffer
        }
        fn delete_vqueue_buffer_len(&self) -> u32 {
            self.delete_buffer
        }
        fn get_dimension_size(&self) -> usize {
            128
        }
        fn broken_index_count(&self) -> u64 {
            self.broken_count
        }
        fn is_statistics_enabled(&self) -> bool {
            self.stats_enabled
        }
        fn index_statistics(&self) -> Result<info::index::Statistics, Error> {
            Ok(info::index::Statistics {
                median_indegree: 10,
                median_outdegree: 20,
                ..Default::default()
            })
        }
        fn index_property(&self) -> Result<info::index::Property, Error> {
            Ok(info::index::Property::default())
        }
    }

    #[test]
    fn test_metrics_integration() {
        let exporter = InMemoryMetricExporter::default();
        let reader = PeriodicReader::builder(exporter).build();
        let provider = SdkMeterProvider::builder().with_reader(reader).build();
        global::set_meter_provider(provider);

        let mock_ann = MockANN::new();
        let service = Arc::new(RwLock::new(mock_ann));

        register_metrics(service).unwrap();
    }
}
