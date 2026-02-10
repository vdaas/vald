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
const NODES_SKIPPED_FOR_INDEGREE_DISTANCE: &str = "agent_core_ngt_nodes_skipped_for_indegree_distance";
const NUMBER_OF_EDGES: &str = "agent_core_ngt_number_of_edges";
const NUMBER_OF_INDEXED_OBJECTS: &str = "agent_core_ngt_number_of_indexed_objects";
const NUMBER_OF_NODES: &str = "agent_core_ngt_number_of_nodes";
const NUMBER_OF_NODES_WITHOUT_EDGES: &str = "agent_core_ngt_number_of_nodes_without_edges";
const NUMBER_OF_NODES_WITHOUT_INDEGREE: &str = "agent_core_ngt_number_of_nodes_without_indegree";
const NUMBER_OF_OBJECTS: &str = "agent_core_ngt_number_of_objects";
const NUMBER_OF_REMOVED_OBJECTS: &str = "agent_core_ngt_number_of_removed_objects";
const SIZE_OF_OBJECT_REPOSITORY: &str = "agent_core_ngt_size_of_object_repository";
const SIZE_OF_REFINEMENT_OBJECT_REPOSITORY: &str = "agent_core_ngt_size_of_refinement_object_repository";
const VARIANCE_OF_INDEGREE: &str = "agent_core_ngt_variance_of_indegree";
const VARIANCE_OF_OUTDEGREE: &str = "agent_core_ngt_variance_of_outdegree";
const MEAN_EDGE_LENGTH: &str = "agent_core_ngt_mean_edge_length";
const MEAN_EDGE_LENGTH_FOR_10_EDGES: &str = "agent_core_ngt_mean_edge_length_for_10_edges";
const MEAN_INDEGREE_DISTANCE_FOR_10_EDGES: &str = "agent_core_ngt_mean_indegree_distance_for_10_edges";
const MEAN_NUMBER_OF_EDGES_PER_NODE: &str = "agent_core_ngt_mean_number_of_edges_per_node";
const C1_INDEGREE: &str = "agent_core_ngt_c1_indegree";
const C5_INDEGREE: &str = "agent_core_ngt_c5_indegree";
const C95_OUTDEGREE: &str = "agent_core_ngt_c95_outdegree";
const C99_OUTDEGREE: &str = "agent_core_ngt_c99_outdegree";

pub fn register_metrics<S>(service: Arc<RwLock<S>>) -> anyhow::Result<()>
where
    S: ANN + 'static,
{
    let meter = global::meter("vald-agent");
    let svc = Arc::downgrade(&service);

    // Basic Metrics
    let svc_index_count = svc.clone();
    let _index_count = meter
        .i64_observable_gauge(INDEX_COUNT)
        .with_description("Agent NGT index count")
        .with_callback(move |observer| {
            if let Some(service) = svc_index_count.upgrade() {
                if let Ok(s) = service.try_read() {
                    observer.observe(s.len() as i64, &[]);
                }
            }
        })
        .build();
    let svc_uncommitted_index_count = svc.clone();
    let _uncommitted_index_count = meter
        .i64_observable_gauge(UNCOMMITTED_INDEX_COUNT)
        .with_description("Agent NGT uncommitted index count")
        .with_callback(move |observer| {
            if let Some(service) = svc_uncommitted_index_count.upgrade() {
                if let Ok(s) = service.try_read() {
                    let total = s.insert_vqueue_buffer_len() + s.delete_vqueue_buffer_len();
                    observer.observe(total as i64, &[]);
                }
            }
        })
        .build();
    let svc_insert_vqueue_count = svc.clone();
    let _insert_vqueue_count = meter
        .i64_observable_gauge(INSERT_VQUEUE_COUNT)
        .with_description("Agent NGT insert vqueue count")
        .with_callback(move |observer| {
            if let Some(service) = svc_insert_vqueue_count.upgrade() {
                if let Ok(s) = service.try_read() {
                    observer.observe(s.insert_vqueue_buffer_len() as i64, &[]);
                }
            }
        })
        .build();
    let svc_delete_vqueue_count = svc.clone();
    let _delete_vqueue_count = meter
        .i64_observable_gauge(DELETE_VQUEUE_COUNT)
        .with_description("Agent NGT delete vqueue count")
        .with_callback(move |observer| {
            if let Some(service) = svc_delete_vqueue_count.upgrade() {
                if let Ok(s) = service.try_read() {
                    observer.observe(s.delete_vqueue_buffer_len() as i64, &[]);
                }
            }
        })
        .build();
    let svc_completed_create_index_total = svc.clone();
    let _completed_create_index_total = meter
        .i64_observable_gauge(COMPLETED_CREATE_INDEX_TOTAL)
        .with_description("The cumulative count of completed create index execution")
        .with_callback(move |observer| {
            if let Some(service) = svc_completed_create_index_total.upgrade() {
                if let Ok(s) = service.try_read() {
                    observer.observe(s.number_of_create_index_executions() as i64, &[]);
                }
            }
        })
        .build();
    let _executed_proactive_gc_total = meter
        .i64_observable_gauge(EXECUTED_PROACTIVE_GC_TOTAL)
        .with_description("The cumulative count of proactive GC execution")
        .with_callback(|observer| {
            observer.observe(0_i64, &[]);
        })
        .build();
    let svc_is_indexing = svc.clone();
    let _is_indexing = meter
        .i64_observable_gauge(IS_INDEXING)
        .with_description("Currently indexing or no")
        .with_callback(move |observer| {
            if let Some(service) = svc_is_indexing.upgrade() {
                if let Ok(s) = service.try_read() {
                    observer.observe(if s.is_indexing() { 1 } else { 0 }, &[]);
                }
            }
        })
        .build();
    let svc_is_saving = svc.clone();
    let _is_saving = meter
        .i64_observable_gauge(IS_SAVING)
        .with_description("Currently saving or not")
        .with_callback(move |observer| {
            if let Some(service) = svc_is_saving.upgrade() {
                if let Ok(s) = service.try_read() {
                    observer.observe(if s.is_saving() { 1 } else { 0 }, &[]);
                }
            }
        })
        .build();
    let svc_broken_index_store_count = svc.clone();
    let _broken_index_store_count = meter
        .i64_observable_gauge(BROKEN_INDEX_STORE_COUNT)
        .with_description("How many broken index generations have been stored")
        .with_callback(move |observer| {
            if let Some(service) = svc_broken_index_store_count.upgrade() {
                if let Ok(s) = service.try_read() {
                    observer.observe(s.broken_index_count() as i64, &[]);
                }
            }
        })
        .build();

    // Statistics Metrics (Int64)
    let svc_median_indegree = svc.clone();
    let _median_indegree = meter
        .i64_observable_gauge(MEDIAN_INDEGREE)
        .with_description("Median indegree of nodes")
        .with_callback(move |observer| {
            if let Some(service) = svc_median_indegree.upgrade() {
                if let Ok(s) = service.try_read() {
                    if s.is_statistics_enabled() {
                        if let Ok(stats) = s.index_statistics() {
                            observer.observe(stats.median_indegree as i64, &[]);
                        }
                    }
                }
            }
        })
        .build();
    let svc_median_outdegree = svc.clone();
    let _median_outdegree = meter
        .i64_observable_gauge(MEDIAN_OUTDEGREE)
        .with_description("Median outdegree of nodes")
        .with_callback(move |observer| {
            if let Some(service) = svc_median_outdegree.upgrade() {
                if let Ok(s) = service.try_read() {
                    if s.is_statistics_enabled() {
                        if let Ok(stats) = s.index_statistics() {
                            observer.observe(stats.median_outdegree as i64, &[]);
                        }
                    }
                }
            }
        })
        .build();
    let svc_max_number_of_indegree = svc.clone();
    let _max_number_of_indegree = meter
        .i64_observable_gauge(MAX_NUMBER_OF_INDEGREE)
        .with_description("Maximum number of indegree")
        .with_callback(move |observer| {
            if let Some(service) = svc_max_number_of_indegree.upgrade() {
                if let Ok(s) = service.try_read() {
                    if s.is_statistics_enabled() {
                        if let Ok(stats) = s.index_statistics() {
                            observer.observe(stats.max_number_of_indegree as i64, &[]);
                        }
                    }
                }
            }
        })
        .build();
    let svc_max_number_of_outdegree = svc.clone();
    let _max_number_of_outdegree = meter
        .i64_observable_gauge(MAX_NUMBER_OF_OUTDEGREE)
        .with_description("Maximum number of outdegree")
        .with_callback(move |observer| {
            if let Some(service) = svc_max_number_of_outdegree.upgrade() {
                if let Ok(s) = service.try_read() {
                    if s.is_statistics_enabled() {
                        if let Ok(stats) = s.index_statistics() {
                            observer.observe(stats.max_number_of_outdegree as i64, &[]);
                        }
                    }
                }
            }
        })
        .build();
    let svc_min_number_of_indegree = svc.clone();
    let _min_number_of_indegree = meter
        .i64_observable_gauge(MIN_NUMBER_OF_INDEGREE)
        .with_description("Minimum number of indegree")
        .with_callback(move |observer| {
            if let Some(service) = svc_min_number_of_indegree.upgrade() {
                if let Ok(s) = service.try_read() {
                    if s.is_statistics_enabled() {
                        if let Ok(stats) = s.index_statistics() {
                            observer.observe(stats.min_number_of_indegree as i64, &[]);
                        }
                    }
                }
            }
        })
        .build();
    let svc_min_number_of_outdegree = svc.clone();
    let _min_number_of_outdegree = meter
        .i64_observable_gauge(MIN_NUMBER_OF_OUTDEGREE)
        .with_description("Minimum number of outdegree")
        .with_callback(move |observer| {
            if let Some(service) = svc_min_number_of_outdegree.upgrade() {
                if let Ok(s) = service.try_read() {
                    if s.is_statistics_enabled() {
                        if let Ok(stats) = s.index_statistics() {
                            observer.observe(stats.min_number_of_outdegree as i64, &[]);
                        }
                    }
                }
            }
        })
        .build();
    let svc_mode_indegree = svc.clone();
    let _mode_indegree = meter
        .i64_observable_gauge(MODE_INDEGREE)
        .with_description("Mode of indegree")
        .with_callback(move |observer| {
            if let Some(service) = svc_mode_indegree.upgrade() {
                if let Ok(s) = service.try_read() {
                    if s.is_statistics_enabled() {
                        if let Ok(stats) = s.index_statistics() {
                            observer.observe(stats.mode_indegree as i64, &[]);
                        }
                    }
                }
            }
        })
        .build();
    let svc_mode_outdegree = svc.clone();
    let _mode_outdegree = meter
        .i64_observable_gauge(MODE_OUTDEGREE)
        .with_description("Mode of outdegree")
        .with_callback(move |observer| {
            if let Some(service) = svc_mode_outdegree.upgrade() {
                if let Ok(s) = service.try_read() {
                    if s.is_statistics_enabled() {
                        if let Ok(stats) = s.index_statistics() {
                            observer.observe(stats.mode_outdegree as i64, &[]);
                        }
                    }
                }
            }
        })
        .build();
    let svc_nodes_skipped_for_10_edges = svc.clone();
    let _nodes_skipped_for_10_edges = meter
        .i64_observable_gauge(NODES_SKIPPED_FOR_10_EDGES)
        .with_description("Nodes skipped for 10 edges")
        .with_callback(move |observer| {
            if let Some(service) = svc_nodes_skipped_for_10_edges.upgrade() {
                if let Ok(s) = service.try_read() {
                    if s.is_statistics_enabled() {
                        if let Ok(stats) = s.index_statistics() {
                            observer.observe(stats.nodes_skipped_for_10_edges as i64, &[]);
                        }
                    }
                }
            }
        })
        .build();
    let svc_nodes_skipped_for_indegree_distance = svc.clone();
    let _nodes_skipped_for_indegree_distance = meter
        .i64_observable_gauge(NODES_SKIPPED_FOR_INDEGREE_DISTANCE)
        .with_description("Nodes skipped for indegree distance")
        .with_callback(move |observer| {
            if let Some(service) = svc_nodes_skipped_for_indegree_distance.upgrade() {
                if let Ok(s) = service.try_read() {
                    if s.is_statistics_enabled() {
                        if let Ok(stats) = s.index_statistics() {
                            observer.observe(stats.nodes_skipped_for_indegree_distance as i64, &[]);
                        }
                    }
                }
            }
        })
        .build();
    let svc_number_of_edges = svc.clone();
    let _number_of_edges = meter
        .i64_observable_gauge(NUMBER_OF_EDGES)
        .with_description("Number of edges")
        .with_callback(move |observer| {
            if let Some(service) = svc_number_of_edges.upgrade() {
                if let Ok(s) = service.try_read() {
                    if s.is_statistics_enabled() {
                        if let Ok(stats) = s.index_statistics() {
                            observer.observe(stats.number_of_edges as i64, &[]);
                        }
                    }
                }
            }
        })
        .build();
    let svc_number_of_indexed_objects = svc.clone();
    let _number_of_indexed_objects = meter
        .i64_observable_gauge(NUMBER_OF_INDEXED_OBJECTS)
        .with_description("Number of indexed objects")
        .with_callback(move |observer| {
            if let Some(service) = svc_number_of_indexed_objects.upgrade() {
                if let Ok(s) = service.try_read() {
                    if s.is_statistics_enabled() {
                        if let Ok(stats) = s.index_statistics() {
                            observer.observe(stats.number_of_indexed_objects as i64, &[]);
                        }
                    }
                }
            }
        })
        .build();
    let svc_number_of_nodes = svc.clone();
    let _number_of_nodes = meter
        .i64_observable_gauge(NUMBER_OF_NODES)
        .with_description("Number of nodes")
        .with_callback(move |observer| {
            if let Some(service) = svc_number_of_nodes.upgrade() {
                if let Ok(s) = service.try_read() {
                    if s.is_statistics_enabled() {
                        if let Ok(stats) = s.index_statistics() {
                            observer.observe(stats.number_of_nodes as i64, &[]);
                        }
                    }
                }
            }
        })
        .build();
    let svc_number_of_nodes_without_edges = svc.clone();
    let _number_of_nodes_without_edges = meter
        .i64_observable_gauge(NUMBER_OF_NODES_WITHOUT_EDGES)
        .with_description("Number of nodes without edges")
        .with_callback(move |observer| {
            if let Some(service) = svc_number_of_nodes_without_edges.upgrade() {
                if let Ok(s) = service.try_read() {
                    if s.is_statistics_enabled() {
                        if let Ok(stats) = s.index_statistics() {
                            observer.observe(stats.number_of_nodes_without_edges as i64, &[]);
                        }
                    }
                }
            }
        })
        .build();
    let svc_number_of_nodes_without_indegree = svc.clone();
    let _number_of_nodes_without_indegree = meter
        .i64_observable_gauge(NUMBER_OF_NODES_WITHOUT_INDEGREE)
        .with_description("Number of nodes without indegree")
        .with_callback(move |observer| {
            if let Some(service) = svc_number_of_nodes_without_indegree.upgrade() {
                if let Ok(s) = service.try_read() {
                    if s.is_statistics_enabled() {
                        if let Ok(stats) = s.index_statistics() {
                            observer.observe(stats.number_of_nodes_without_indegree as i64, &[]);
                        }
                    }
                }
            }
        })
        .build();
    let svc_number_of_objects = svc.clone();
    let _number_of_objects = meter
        .i64_observable_gauge(NUMBER_OF_OBJECTS)
        .with_description("Number of objects")
        .with_callback(move |observer| {
            if let Some(service) = svc_number_of_objects.upgrade() {
                if let Ok(s) = service.try_read() {
                    if s.is_statistics_enabled() {
                        if let Ok(stats) = s.index_statistics() {
                            observer.observe(stats.number_of_objects as i64, &[]);
                        }
                    }
                }
            }
        })
        .build();
    let svc_number_of_removed_objects = svc.clone();
    let _number_of_removed_objects = meter
        .i64_observable_gauge(NUMBER_OF_REMOVED_OBJECTS)
        .with_description("Number of removed objects")
        .with_callback(move |observer| {
            if let Some(service) = svc_number_of_removed_objects.upgrade() {
                if let Ok(s) = service.try_read() {
                    if s.is_statistics_enabled() {
                        if let Ok(stats) = s.index_statistics() {
                            observer.observe(stats.number_of_removed_objects as i64, &[]);
                        }
                    }
                }
            }
        })
        .build();
    let svc_size_of_object_repository = svc.clone();
    let _size_of_object_repository = meter
        .i64_observable_gauge(SIZE_OF_OBJECT_REPOSITORY)
        .with_description("Size of object repository")
        .with_callback(move |observer| {
            if let Some(service) = svc_size_of_object_repository.upgrade() {
                if let Ok(s) = service.try_read() {
                    if s.is_statistics_enabled() {
                        if let Ok(stats) = s.index_statistics() {
                            observer.observe(stats.size_of_object_repository as i64, &[]);
                        }
                    }
                }
            }
        })
        .build();
    let svc_size_of_refinement_object_repository = svc.clone();
    let _size_of_refinement_object_repository = meter
        .i64_observable_gauge(SIZE_OF_REFINEMENT_OBJECT_REPOSITORY)
        .with_description("Size of refinement object repository")
        .with_callback(move |observer| {
            if let Some(service) = svc_size_of_refinement_object_repository.upgrade() {
                if let Ok(s) = service.try_read() {
                    if s.is_statistics_enabled() {
                        if let Ok(stats) = s.index_statistics() {
                            observer.observe(stats.size_of_refinement_object_repository as i64, &[]);
                        }
                    }
                }
            }
        })
        .build();

    // Statistics Metrics (Float64)
    let svc_variance_of_indegree = svc.clone();
    let _variance_of_indegree = meter
        .f64_observable_gauge(VARIANCE_OF_INDEGREE)
        .with_description("Variance of indegree")
        .with_callback(move |observer| {
            if let Some(service) = svc_variance_of_indegree.upgrade() {
                if let Ok(s) = service.try_read() {
                    if s.is_statistics_enabled() {
                        if let Ok(stats) = s.index_statistics() {
                            observer.observe(stats.variance_of_indegree as f64, &[]);
                        }
                    }
                }
            }
        })
        .build();
    let svc_variance_of_outdegree = svc.clone();
    let _variance_of_outdegree = meter
        .f64_observable_gauge(VARIANCE_OF_OUTDEGREE)
        .with_description("Variance of outdegree")
        .with_callback(move |observer| {
            if let Some(service) = svc_variance_of_outdegree.upgrade() {
                if let Ok(s) = service.try_read() {
                    if s.is_statistics_enabled() {
                        if let Ok(stats) = s.index_statistics() {
                            observer.observe(stats.variance_of_outdegree as f64, &[]);
                        }
                    }
                }
            }
        })
        .build();
    let svc_mean_edge_length = svc.clone();
    let _mean_edge_length = meter
        .f64_observable_gauge(MEAN_EDGE_LENGTH)
        .with_description("Mean edge length")
        .with_callback(move |observer| {
            if let Some(service) = svc_mean_edge_length.upgrade() {
                if let Ok(s) = service.try_read() {
                    if s.is_statistics_enabled() {
                        if let Ok(stats) = s.index_statistics() {
                            observer.observe(stats.mean_edge_length as f64, &[]);
                        }
                    }
                }
            }
        })
        .build();
    let svc_mean_edge_length_for_10_edges = svc.clone();
    let _mean_edge_length_for_10_edges = meter
        .f64_observable_gauge(MEAN_EDGE_LENGTH_FOR_10_EDGES)
        .with_description("Mean edge length for 10 edges")
        .with_callback(move |observer| {
            if let Some(service) = svc_mean_edge_length_for_10_edges.upgrade() {
                if let Ok(s) = service.try_read() {
                    if s.is_statistics_enabled() {
                        if let Ok(stats) = s.index_statistics() {
                            observer.observe(stats.mean_edge_length_for_10_edges as f64, &[]);
                        }
                    }
                }
            }
        })
        .build();
    let svc_mean_indegree_distance_for_10_edges = svc.clone();
    let _mean_indegree_distance_for_10_edges = meter
        .f64_observable_gauge(MEAN_INDEGREE_DISTANCE_FOR_10_EDGES)
        .with_description("Mean indegree distance for 10 edges")
        .with_callback(move |observer| {
            if let Some(service) = svc_mean_indegree_distance_for_10_edges.upgrade() {
                if let Ok(s) = service.try_read() {
                    if s.is_statistics_enabled() {
                        if let Ok(stats) = s.index_statistics() {
                            observer.observe(stats.mean_indegree_distance_for_10_edges as f64, &[]);
                        }
                    }
                }
            }
        })
        .build();
    let svc_mean_number_of_edges_per_node = svc.clone();
    let _mean_number_of_edges_per_node = meter
        .f64_observable_gauge(MEAN_NUMBER_OF_EDGES_PER_NODE)
        .with_description("Mean number of edges per node")
        .with_callback(move |observer| {
            if let Some(service) = svc_mean_number_of_edges_per_node.upgrade() {
                if let Ok(s) = service.try_read() {
                    if s.is_statistics_enabled() {
                        if let Ok(stats) = s.index_statistics() {
                            observer.observe(stats.mean_number_of_edges_per_node as f64, &[]);
                        }
                    }
                }
            }
        })
        .build();
    let svc_c1_indegree = svc.clone();
    let _c1_indegree = meter
        .f64_observable_gauge(C1_INDEGREE)
        .with_description("C1 indegree")
        .with_callback(move |observer| {
            if let Some(service) = svc_c1_indegree.upgrade() {
                if let Ok(s) = service.try_read() {
                    if s.is_statistics_enabled() {
                        if let Ok(stats) = s.index_statistics() {
                            observer.observe(stats.c1_indegree as f64, &[]);
                        }
                    }
                }
            }
        })
        .build();
    let svc_c5_indegree = svc.clone();
    let _c5_indegree = meter
        .f64_observable_gauge(C5_INDEGREE)
        .with_description("C5 indegree")
        .with_callback(move |observer| {
            if let Some(service) = svc_c5_indegree.upgrade() {
                if let Ok(s) = service.try_read() {
                    if s.is_statistics_enabled() {
                        if let Ok(stats) = s.index_statistics() {
                            observer.observe(stats.c5_indegree as f64, &[]);
                        }
                    }
                }
            }
        })
        .build();
    let svc_c95_outdegree = svc.clone();
    let _c95_outdegree = meter
        .f64_observable_gauge(C95_OUTDEGREE)
        .with_description("C95 outdegree")
        .with_callback(move |observer| {
            if let Some(service) = svc_c95_outdegree.upgrade() {
                if let Ok(s) = service.try_read() {
                    if s.is_statistics_enabled() {
                        if let Ok(stats) = s.index_statistics() {
                            observer.observe(stats.c95_outdegree as f64, &[]);
                        }
                    }
                }
            }
        })
        .build();
    let svc_c99_outdegree = svc;
    let _c99_outdegree = meter
        .f64_observable_gauge(C99_OUTDEGREE)
        .with_description("C99 outdegree")
        .with_callback(move |observer| {
            if let Some(service) = svc_c99_outdegree.upgrade() {
                if let Ok(s) = service.try_read() {
                    if s.is_statistics_enabled() {
                        if let Ok(stats) = s.index_statistics() {
                            observer.observe(stats.c99_outdegree as f64, &[]);
                        }
                    }
                }
            }
        })
        .build();

    Ok(())
}

#[cfg(test)]
mod tests {
    use super::*;
    use algorithm::{Error, ANN};
    use proto::payload::v1::{info, search};
    use std::collections::HashMap;
    use std::future::Future;
    use opentelemetry_sdk::metrics::{InMemoryMetricExporter, PeriodicReader, SdkMeterProvider};

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
        fn search(&self, _v: Vec<f32>, _k: u32, _e: f32, _r: f32) -> impl Future<Output = Result<search::Response, Error>> + Send { async { Ok(search::Response::default()) } }
        fn search_by_id(&self, _u: String, _k: u32, _e: f32, _r: f32) -> impl Future<Output = Result<search::Response, Error>> + Send { async { Ok(search::Response::default()) } }
        fn linear_search(&self, _v: Vec<f32>, _k: u32) -> impl Future<Output = Result<search::Response, Error>> + Send { async { Ok(search::Response::default()) } }
        fn linear_search_by_id(&self, _u: String, _k: u32) -> impl Future<Output = Result<search::Response, Error>> + Send { async { Ok(search::Response::default()) } }
        fn insert(&mut self, _u: String, _v: Vec<f32>) -> impl Future<Output = Result<(), Error>> + Send { async { Ok(()) } }
        fn insert_with_time(&mut self, _u: String, _v: Vec<f32>, _t: i64) -> impl Future<Output = Result<(), Error>> + Send { async { Ok(()) } }
        fn insert_multiple(&mut self, _vs: HashMap<String, Vec<f32>>) -> impl Future<Output = Result<(), Error>> + Send { async { Ok(()) } }
        fn insert_multiple_with_time(&mut self, _vs: HashMap<String, Vec<f32>>, _t: i64) -> impl Future<Output = Result<(), Error>> + Send { async { Ok(()) } }
        fn update(&mut self, _u: String, _v: Vec<f32>) -> impl Future<Output = Result<(), Error>> + Send { async { Ok(()) } }
        fn update_with_time(&mut self, _u: String, _v: Vec<f32>, _t: i64) -> impl Future<Output = Result<(), Error>> + Send { async { Ok(()) } }
        fn update_multiple(&mut self, _vs: HashMap<String, Vec<f32>>) -> impl Future<Output = Result<(), Error>> + Send { async { Ok(()) } }
        fn update_multiple_with_time(&mut self, _vs: HashMap<String, Vec<f32>>, _t: i64) -> impl Future<Output = Result<(), Error>> + Send { async { Ok(()) } }
        fn update_timestamp(&mut self, _u: String, _t: i64, _f: bool) -> impl Future<Output = Result<(), Error>> + Send { async { Ok(()) } }
        fn remove(&mut self, _u: String) -> impl Future<Output = Result<(), Error>> + Send { async { Ok(()) } }
        fn remove_with_time(&mut self, _u: String, _t: i64) -> impl Future<Output = Result<(), Error>> + Send { async { Ok(()) } }
        fn remove_multiple(&mut self, _us: Vec<String>) -> impl Future<Output = Result<(), Error>> + Send { async { Ok(()) } }
        fn remove_multiple_with_time(&mut self, _us: Vec<String>, _t: i64) -> impl Future<Output = Result<(), Error>> + Send { async { Ok(()) } }
        fn regenerate_indexes(&mut self) -> impl Future<Output = Result<(), Error>> + Send { async { Ok(()) } }
        fn create_index(&mut self) -> impl Future<Output = Result<(), Error>> + Send { async { Ok(()) } }
        fn save_index(&mut self) -> impl Future<Output = Result<(), Error>> + Send { async { Ok(()) } }
        fn create_and_save_index(&mut self) -> impl Future<Output = Result<(), Error>> + Send { async { Ok(()) } }
        fn get_object(&self, _u: String) -> impl Future<Output = Result<(Vec<f32>, i64), Error>> + Send { async { Ok((vec![], 0)) } }
        fn exists(&self, _u: String) -> impl Future<Output = (usize, bool)> + Send { async { (0, false) } }
        fn uuids(&self) -> impl Future<Output = Vec<String>> + Send { async { vec![] } }
        fn list_object_func<F: FnMut(String, Vec<f32>, i64) -> bool + Send>(&self, _f: F) -> impl Future<Output = ()> + Send { async {} }
        fn close(&mut self) -> impl Future<Output = Result<(), Error>> + Send { async { Ok(()) } }

        // Metrics methods
        fn is_indexing(&self) -> bool { self.indexing }
        fn is_flushing(&self) -> bool { false }
        fn is_saving(&self) -> bool { self.saving }
        fn len(&self) -> u32 { self.len }
        fn number_of_create_index_executions(&self) -> u64 { self.create_index_count }
        fn insert_vqueue_buffer_len(&self) -> u32 { self.insert_buffer }
        fn delete_vqueue_buffer_len(&self) -> u32 { self.delete_buffer }
        fn get_dimension_size(&self) -> usize { 128 }
        fn broken_index_count(&self) -> u64 { self.broken_count }
        fn is_statistics_enabled(&self) -> bool { self.stats_enabled }
        fn index_statistics(&self) -> Result<info::index::Statistics, Error> {
            Ok(info::index::Statistics {
                median_indegree: 10,
                median_outdegree: 20,
                ..Default::default()
            })
        }
        fn index_property(&self) -> Result<info::index::Property, Error> { Ok(info::index::Property::default()) }
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
