// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package ngt

import (
	"context"

	"github.com/vdaas/vald/internal/observability/attribute"
	"github.com/vdaas/vald/internal/observability/metrics"
	"github.com/vdaas/vald/pkg/agent/core/ngt/service"
	api "go.opentelemetry.io/otel/metric"
	view "go.opentelemetry.io/otel/sdk/metric"
)

const (
	indexCountMetricsName        = "agent_core_ngt_index_count"
	indexCountMetricsDescription = "Agent NGT index count"

	uncommittedIndexCountMetricsName        = "agent_core_ngt_uncommitted_index_count"
	uncommittedIndexCountMetricsDescription = "Agent NGT index count"

	insertVQueueCountMetricsName        = "agent_core_ngt_insert_vqueue_count"
	insertVQueueCountMetricsDescription = "Agent NGT insert vqueue count"

	deleteVQueueCountMetricsName        = "agent_core_ngt_delete_vqueue_count"
	deleteVQueueCountMetricsDescription = "Agent NGT delete vqueue count"

	completedCreateIndexTotalMetricsName        = "agent_core_ngt_completed_create_index_total"
	completedCreateIndexTotalMetricsDescription = "The cumulative count of completed create index execution"

	executedProactiveGCTotalMetricsName        = "agent_core_ngt_executed_proactive_gc_total"
	executedProactiveGCTotalMetricsDescription = "The cumulative count of proactive GC execution"

	isIndexingMetricsName        = "agent_core_ngt_is_indexing"
	isIndexingMetricsDescription = "Currently indexing or no"

	isSavingMetricsName        = "agent_core_ngt_is_saving"
	isSavingMetricsDescription = "Currently saving or not"

	brokenIndexStoreCountMetricsName        = "agent_core_ngt_broken_index_store_count"
	brokenIndexStoreCountMetricsDescription = "How many broken index generations have been stored"

	medianIndegreeMetricsName        = "agent_core_ngt_median_indegree"
	medianIndegreeMetricsDescription = "Median indegree of nodes"

	medianOutdegreeMetricsName        = "agent_core_ngt_median_outdegree"
	medianOutdegreeMetricsDescription = "Median outdegree of nodes"

	maxNumberOfIndegreeMetricsName        = "agent_core_ngt_max_number_of_indegree"
	maxNumberOfIndegreeMetricsDescription = "Maximum number of indegree"

	maxNumberOfOutdegreeMetricsName        = "agent_core_ngt_max_number_of_outdegree"
	maxNumberOfOutdegreeMetricsDescription = "Maximum number of outdegree"

	minNumberOfIndegreeMetricsName        = "agent_core_ngt_min_number_of_indegree"
	minNumberOfIndegreeMetricsDescription = "Minimum number of indegree"

	minNumberOfOutdegreeMetricsName        = "agent_core_ngt_min_number_of_outdegree"
	minNumberOfOutdegreeMetricsDescription = "Minimum number of outdegree"

	modeIndegreeMetricsName        = "agent_core_ngt_mode_indegree"
	modeIndegreeMetricsDescription = "Mode of indegree"

	modeOutdegreeMetricsName        = "agent_core_ngt_mode_outdegree"
	modeOutdegreeMetricsDescription = "Mode of outdegree"

	nodesSkippedFor10EdgesMetricsName        = "agent_core_ngt_nodes_skipped_for_10_edges"
	nodesSkippedFor10EdgesMetricsDescription = "Nodes skipped for 10 edges"

	nodesSkippedForIndegreeDistanceMetricsName        = "agent_core_ngt_nodes_skipped_for_indegree_distance"
	nodesSkippedForIndegreeDistanceMetricsDescription = "Nodes skipped for indegree distance"

	numberOfEdgesMetricsName        = "agent_core_ngt_number_of_edges"
	numberOfEdgesMetricsDescription = "Number of edges"

	numberOfIndexedObjectsMetricsName        = "agent_core_ngt_number_of_indexed_objects"
	numberOfIndexedObjectsMetricsDescription = "Number of indexed objects"

	numberOfNodesMetricsName        = "agent_core_ngt_number_of_nodes"
	numberOfNodesMetricsDescription = "Number of nodes"

	numberOfNodesWithoutEdgesMetricsName        = "agent_core_ngt_number_of_nodes_without_edges"
	numberOfNodesWithoutEdgesMetricsDescription = "Number of nodes without edges"

	numberOfNodesWithoutIndegreeMetricsName        = "agent_core_ngt_number_of_nodes_without_indegree"
	numberOfNodesWithoutIndegreeMetricsDescription = "Number of nodes without indegree"

	numberOfObjectsMetricsName        = "agent_core_ngt_number_of_objects"
	numberOfObjectsMetricsDescription = "Number of objects"

	numberOfRemovedObjectsMetricsName        = "agent_core_ngt_number_of_removed_objects"
	numberOfRemovedObjectsMetricsDescription = "Number of removed objects"

	sizeOfObjectRepositoryMetricsName        = "agent_core_ngt_size_of_object_repository"
	sizeOfObjectRepositoryMetricsDescription = "Size of object repository"

	sizeOfRefinementObjectRepositoryMetricsName        = "agent_core_ngt_size_of_refinement_object_repository"
	sizeOfRefinementObjectRepositoryMetricsDescription = "Size of refinement object repository"

	varianceOfIndegreeMetricsName        = "agent_core_ngt_variance_of_indegree"
	varianceOfIndegreeMetricsDescription = "Variance of indegree"

	varianceOfOutdegreeMetricsName        = "agent_core_ngt_variance_of_outdegree"
	varianceOfOutdegreeMetricsDescription = "Variance of outdegree"

	meanEdgeLengthMetricsName        = "agent_core_ngt_mean_edge_length"
	meanEdgeLengthMetricsDescription = "Mean edge length"

	meanEdgeLengthFor10EdgesMetricsName        = "agent_core_ngt_mean_edge_length_for_10_edges"
	meanEdgeLengthFor10EdgesMetricsDescription = "Mean edge length for 10 edges"

	meanIndegreeDistanceFor10EdgesMetricsName        = "agent_core_ngt_mean_indegree_distance_for_10_edges"
	meanIndegreeDistanceFor10EdgesMetricsDescription = "Mean indegree distance for 10 edges"

	meanNumberOfEdgesPerNodeMetricsName        = "agent_core_ngt_mean_number_of_edges_per_node"
	meanNumberOfEdgesPerNodeMetricsDescription = "Mean number of edges per node"

	c1IndegreeMetricsName        = "agent_core_ngt_c1_indegree"
	c1IndegreeMetricsDescription = "C1 indegree"

	c5IndegreeMetricsName        = "agent_core_ngt_c5_indegree"
	c5IndegreeMetricsDescription = "C5 indegree"

	c95OutdegreeMetricsName        = "agent_core_ngt_c95_outdegree"
	c95OutdegreeMetricsDescription = "C95 outdegree"

	c99OutdegreeMetricsName        = "agent_core_ngt_c99_outdegree"
	c99OutdegreeMetricsDescription = "C99 outdegree"

	indegreeCountMetricsName        = "agent_core_ngt_indegree_count"
	indegreeCountMetricsDescription = "Indegree count"

	outdegreeHistogramMetricsName        = "agent_core_ngt_outdegree_histogram"
	outdegreeHistogramMetricsDescription = "Outdegree histogram"

	indegreeHistogramMetricsName        = "agent_core_ngt_indegree_histogram"
	indegreeHistogramMetricsDescription = "Indegree histogram"
)

type ngtMetrics struct {
	ngt service.NGT
}

func New(n service.NGT) metrics.Metric {
	return &ngtMetrics{
		ngt: n,
	}
}

func (n *ngtMetrics) View() (mv []metrics.View, err error) {
	mv = []metrics.View{
		view.NewView(
			view.Instrument{
				Name:        indexCountMetricsName,
				Description: indexCountMetricsDescription,
			},
			view.Stream{
				Aggregation: view.AggregationLastValue{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        uncommittedIndexCountMetricsName,
				Description: uncommittedIndexCountMetricsDescription,
			},
			view.Stream{
				Aggregation: view.AggregationLastValue{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        insertVQueueCountMetricsName,
				Description: insertVQueueCountMetricsDescription,
			},
			view.Stream{
				Aggregation: view.AggregationLastValue{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        deleteVQueueCountMetricsName,
				Description: deleteVQueueCountMetricsDescription,
			},
			view.Stream{
				Aggregation: view.AggregationLastValue{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        completedCreateIndexTotalMetricsName,
				Description: completedCreateIndexTotalMetricsDescription,
			},
			view.Stream{
				Aggregation: view.AggregationLastValue{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        executedProactiveGCTotalMetricsName,
				Description: executedProactiveGCTotalMetricsDescription,
			},
			view.Stream{
				Aggregation: view.AggregationLastValue{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        isIndexingMetricsName,
				Description: isIndexingMetricsDescription,
			},
			view.Stream{
				Aggregation: view.AggregationLastValue{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        isSavingMetricsName,
				Description: isSavingMetricsDescription,
			},
			view.Stream{
				Aggregation: view.AggregationLastValue{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        brokenIndexStoreCountMetricsName,
				Description: brokenIndexStoreCountMetricsDescription,
			},
			view.Stream{
				Aggregation: view.AggregationLastValue{},
			},
		),
	}

	if n.ngt.IsStatisticsEnabled() {
		mv = append(mv,
			view.NewView(
				view.Instrument{
					Name:        medianIndegreeMetricsName,
					Description: medianIndegreeMetricsDescription,
				},
				view.Stream{
					Aggregation: view.AggregationLastValue{},
				},
			),
			view.NewView(
				view.Instrument{
					Name:        medianOutdegreeMetricsName,
					Description: medianOutdegreeMetricsDescription,
				},
				view.Stream{
					Aggregation: view.AggregationLastValue{},
				},
			),
			view.NewView(
				view.Instrument{
					Name:        maxNumberOfIndegreeMetricsName,
					Description: maxNumberOfIndegreeMetricsDescription,
				},
				view.Stream{
					Aggregation: view.AggregationLastValue{},
				},
			),
			view.NewView(
				view.Instrument{
					Name:        maxNumberOfOutdegreeMetricsName,
					Description: maxNumberOfOutdegreeMetricsDescription,
				},
				view.Stream{
					Aggregation: view.AggregationLastValue{},
				},
			),
			view.NewView(
				view.Instrument{
					Name:        minNumberOfIndegreeMetricsName,
					Description: minNumberOfIndegreeMetricsDescription,
				},
				view.Stream{
					Aggregation: view.AggregationLastValue{},
				},
			),
			view.NewView(
				view.Instrument{
					Name:        minNumberOfOutdegreeMetricsName,
					Description: minNumberOfOutdegreeMetricsDescription,
				},
				view.Stream{
					Aggregation: view.AggregationLastValue{},
				},
			),
			view.NewView(
				view.Instrument{
					Name:        modeIndegreeMetricsName,
					Description: modeIndegreeMetricsDescription,
				},
				view.Stream{
					Aggregation: view.AggregationLastValue{},
				},
			),
			view.NewView(
				view.Instrument{
					Name:        modeOutdegreeMetricsName,
					Description: modeOutdegreeMetricsDescription,
				},
				view.Stream{
					Aggregation: view.AggregationLastValue{},
				},
			),
			view.NewView(
				view.Instrument{
					Name:        nodesSkippedFor10EdgesMetricsName,
					Description: nodesSkippedFor10EdgesMetricsDescription,
				},
				view.Stream{
					Aggregation: view.AggregationLastValue{},
				},
			),
			view.NewView(
				view.Instrument{
					Name:        nodesSkippedForIndegreeDistanceMetricsName,
					Description: nodesSkippedForIndegreeDistanceMetricsDescription,
				},
				view.Stream{
					Aggregation: view.AggregationLastValue{},
				},
			),
			view.NewView(
				view.Instrument{
					Name:        numberOfEdgesMetricsName,
					Description: numberOfEdgesMetricsDescription,
				},
				view.Stream{
					Aggregation: view.AggregationLastValue{},
				},
			),
			view.NewView(
				view.Instrument{
					Name:        numberOfIndexedObjectsMetricsName,
					Description: numberOfIndexedObjectsMetricsDescription,
				},
				view.Stream{
					Aggregation: view.AggregationLastValue{},
				},
			),
			view.NewView(
				view.Instrument{
					Name:        numberOfNodesMetricsName,
					Description: numberOfNodesMetricsDescription,
				},
				view.Stream{
					Aggregation: view.AggregationLastValue{},
				},
			),
			view.NewView(
				view.Instrument{
					Name:        numberOfNodesWithoutEdgesMetricsName,
					Description: numberOfNodesWithoutEdgesMetricsDescription,
				},
				view.Stream{
					Aggregation: view.AggregationLastValue{},
				},
			),
			view.NewView(
				view.Instrument{
					Name:        numberOfNodesWithoutIndegreeMetricsName,
					Description: numberOfNodesWithoutIndegreeMetricsDescription,
				},
				view.Stream{
					Aggregation: view.AggregationLastValue{},
				},
			),
			view.NewView(
				view.Instrument{
					Name:        numberOfObjectsMetricsName,
					Description: numberOfObjectsMetricsDescription,
				},
				view.Stream{
					Aggregation: view.AggregationLastValue{},
				},
			),
			view.NewView(
				view.Instrument{
					Name:        numberOfRemovedObjectsMetricsName,
					Description: numberOfRemovedObjectsMetricsDescription,
				},
				view.Stream{
					Aggregation: view.AggregationLastValue{},
				},
			),
			view.NewView(
				view.Instrument{
					Name:        sizeOfObjectRepositoryMetricsName,
					Description: sizeOfObjectRepositoryMetricsDescription,
				},
				view.Stream{
					Aggregation: view.AggregationLastValue{},
				},
			),
			view.NewView(
				view.Instrument{
					Name:        sizeOfRefinementObjectRepositoryMetricsName,
					Description: sizeOfRefinementObjectRepositoryMetricsDescription,
				},
				view.Stream{
					Aggregation: view.AggregationLastValue{},
				},
			),
			view.NewView(
				view.Instrument{
					Name:        varianceOfIndegreeMetricsName,
					Description: varianceOfIndegreeMetricsDescription,
				},
				view.Stream{
					Aggregation: view.AggregationLastValue{},
				},
			),
			view.NewView(
				view.Instrument{
					Name:        varianceOfOutdegreeMetricsName,
					Description: varianceOfOutdegreeMetricsDescription,
				},
				view.Stream{
					Aggregation: view.AggregationLastValue{},
				},
			),
			view.NewView(
				view.Instrument{
					Name:        meanEdgeLengthMetricsName,
					Description: meanEdgeLengthMetricsDescription,
				},
				view.Stream{
					Aggregation: view.AggregationLastValue{},
				},
			),
			view.NewView(
				view.Instrument{
					Name:        meanEdgeLengthFor10EdgesMetricsName,
					Description: meanEdgeLengthFor10EdgesMetricsDescription,
				},
				view.Stream{
					Aggregation: view.AggregationLastValue{},
				},
			),
			view.NewView(
				view.Instrument{
					Name:        meanIndegreeDistanceFor10EdgesMetricsName,
					Description: meanIndegreeDistanceFor10EdgesMetricsDescription,
				},
				view.Stream{
					Aggregation: view.AggregationLastValue{},
				},
			),
			view.NewView(
				view.Instrument{
					Name:        meanNumberOfEdgesPerNodeMetricsName,
					Description: meanNumberOfEdgesPerNodeMetricsDescription,
				},
				view.Stream{
					Aggregation: view.AggregationLastValue{},
				},
			),
			view.NewView(
				view.Instrument{
					Name:        c1IndegreeMetricsName,
					Description: c1IndegreeMetricsDescription,
				},
				view.Stream{
					Aggregation: view.AggregationLastValue{},
				},
			),
			view.NewView(
				view.Instrument{
					Name:        c5IndegreeMetricsName,
					Description: c5IndegreeMetricsDescription,
				},
				view.Stream{
					Aggregation: view.AggregationLastValue{},
				},
			),
			view.NewView(
				view.Instrument{
					Name:        c95OutdegreeMetricsName,
					Description: c95OutdegreeMetricsDescription,
				},
				view.Stream{
					Aggregation: view.AggregationLastValue{},
				},
			),
			view.NewView(
				view.Instrument{
					Name:        c99OutdegreeMetricsName,
					Description: c99OutdegreeMetricsDescription,
				},
				view.Stream{
					Aggregation: view.AggregationLastValue{},
				},
			),
			view.NewView(
				view.Instrument{
					Name:        indegreeCountMetricsName,
					Description: indegreeCountMetricsDescription,
				},
				view.Stream{
					Aggregation: view.AggregationLastValue{},
				},
			),
			view.NewView(
				view.Instrument{
					Name:        outdegreeHistogramMetricsName,
					Description: outdegreeHistogramMetricsDescription,
				},
				view.Stream{
					Aggregation: view.AggregationLastValue{},
				},
			),
			view.NewView(
				view.Instrument{
					Name:        indegreeHistogramMetricsName,
					Description: indegreeHistogramMetricsDescription,
				},
				view.Stream{
					Aggregation: view.AggregationLastValue{},
				},
			))
	}
	return mv, nil
}

func (n *ngtMetrics) Register(m metrics.Meter) (err error) {
	var indexCount,
		uncommittedIndexCount,
		insertVQueueCount,
		deleteVQueueCount,
		completedCreateIndexTotal,
		executedProactiveGCTotal,
		isIndexing,
		isSaving,
		brokenIndexCount metrics.Int64ObservableGauge

	indexCount, err = m.Int64ObservableGauge(
		indexCountMetricsName,
		metrics.WithDescription(indexCountMetricsDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	uncommittedIndexCount, err = m.Int64ObservableGauge(
		uncommittedIndexCountMetricsName,
		metrics.WithDescription(uncommittedIndexCountMetricsDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	insertVQueueCount, err = m.Int64ObservableGauge(
		insertVQueueCountMetricsName,
		metrics.WithDescription(insertVQueueCountMetricsDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	deleteVQueueCount, err = m.Int64ObservableGauge(
		deleteVQueueCountMetricsName,
		metrics.WithDescription(deleteVQueueCountMetricsDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	completedCreateIndexTotal, err = m.Int64ObservableGauge(
		completedCreateIndexTotalMetricsName,
		metrics.WithDescription(completedCreateIndexTotalMetricsDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	executedProactiveGCTotal, err = m.Int64ObservableGauge(
		executedProactiveGCTotalMetricsName,
		metrics.WithDescription(executedProactiveGCTotalMetricsDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	isIndexing, err = m.Int64ObservableGauge(
		isIndexingMetricsName,
		metrics.WithDescription(isIndexingMetricsDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	isSaving, err = m.Int64ObservableGauge(
		isSavingMetricsName,
		metrics.WithDescription(isSavingMetricsDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	brokenIndexCount, err = m.Int64ObservableGauge(
		brokenIndexStoreCountMetricsName,
		metrics.WithDescription(brokenIndexStoreCountMetricsDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	instruments := []api.Observable{
		indexCount,
		uncommittedIndexCount,
		insertVQueueCount,
		deleteVQueueCount,
		completedCreateIndexTotal,
		executedProactiveGCTotal,
		isIndexing,
		isSaving,
		brokenIndexCount,
	}
	var (
		medianIndegree,
		medianOutdegree,
		maxNumberOfIndegree,
		maxNumberOfOutdegree,
		minNumberOfIndegree,
		minNumberOfOutdegree,
		modeIndegree,
		modeOutdegree,
		nodesSkippedFor10Edges,
		nodesSkippedForIndegreeDistance,
		numberOfEdges,
		numberOfIndexedObjects,
		numberOfNodes,
		numberOfNodesWithoutEdges,
		numberOfNodesWithoutIndegree,
		numberOfObjects,
		numberOfRemovedObjects,
		sizeOfObjectRepository,
		sizeOfRefinementObjectRepository metrics.Int64ObservableGauge

		varianceOfIndegree,
		varianceOfOutdegree,
		meanEdgeLength,
		meanEdgeLengthFor10Edges,
		meanIndegreeDistanceFor10Edges,
		meanNumberOfEdgesPerNode,
		c1Indegree,
		c5Indegree,
		c95Outdegree,
		c99Outdegree,
		indegreeCount,
		outdegreeHistogram,
		indegreeHistogram metrics.Float64ObservableGauge
	)

	if n.ngt.IsStatisticsEnabled() {
		medianIndegree, err = m.Int64ObservableGauge(
			medianIndegreeMetricsName,
			metrics.WithDescription(medianIndegreeMetricsDescription),
			metrics.WithUnit(metrics.Dimensionless),
		)
		if err != nil {
			return err
		}

		medianOutdegree, err = m.Int64ObservableGauge(
			medianOutdegreeMetricsName,
			metrics.WithDescription(medianOutdegreeMetricsDescription),
			metrics.WithUnit(metrics.Dimensionless),
		)
		if err != nil {
			return err
		}

		maxNumberOfIndegree, err = m.Int64ObservableGauge(
			maxNumberOfIndegreeMetricsName,
			metrics.WithDescription(maxNumberOfIndegreeMetricsDescription),
			metrics.WithUnit(metrics.Dimensionless),
		)
		if err != nil {
			return err
		}

		maxNumberOfOutdegree, err = m.Int64ObservableGauge(
			maxNumberOfOutdegreeMetricsName,
			metrics.WithDescription(maxNumberOfOutdegreeMetricsDescription),
			metrics.WithUnit(metrics.Dimensionless),
		)
		if err != nil {
			return err
		}

		minNumberOfIndegree, err = m.Int64ObservableGauge(
			minNumberOfIndegreeMetricsName,
			metrics.WithDescription(minNumberOfIndegreeMetricsDescription),
			metrics.WithUnit(metrics.Dimensionless),
		)
		if err != nil {
			return err
		}

		minNumberOfOutdegree, err = m.Int64ObservableGauge(
			minNumberOfOutdegreeMetricsName,
			metrics.WithDescription(minNumberOfOutdegreeMetricsDescription),
			metrics.WithUnit(metrics.Dimensionless),
		)
		if err != nil {
			return err
		}

		modeIndegree, err = m.Int64ObservableGauge(
			modeIndegreeMetricsName,
			metrics.WithDescription(modeIndegreeMetricsDescription),
			metrics.WithUnit(metrics.Dimensionless),
		)
		if err != nil {
			return err
		}

		modeOutdegree, err = m.Int64ObservableGauge(
			modeOutdegreeMetricsName,
			metrics.WithDescription(modeOutdegreeMetricsDescription),
			metrics.WithUnit(metrics.Dimensionless),
		)
		if err != nil {
			return err
		}

		nodesSkippedFor10Edges, err = m.Int64ObservableGauge(
			nodesSkippedFor10EdgesMetricsName,
			metrics.WithDescription(nodesSkippedFor10EdgesMetricsDescription),
			metrics.WithUnit(metrics.Dimensionless),
		)
		if err != nil {
			return err
		}

		nodesSkippedForIndegreeDistance, err = m.Int64ObservableGauge(
			nodesSkippedForIndegreeDistanceMetricsName,
			metrics.WithDescription(nodesSkippedForIndegreeDistanceMetricsDescription),
			metrics.WithUnit(metrics.Dimensionless),
		)
		if err != nil {
			return err
		}

		numberOfEdges, err = m.Int64ObservableGauge(
			numberOfEdgesMetricsName,
			metrics.WithDescription(numberOfEdgesMetricsDescription),
			metrics.WithUnit(metrics.Dimensionless),
		)
		if err != nil {
			return err
		}

		numberOfIndexedObjects, err = m.Int64ObservableGauge(
			numberOfIndexedObjectsMetricsName,
			metrics.WithDescription(numberOfIndexedObjectsMetricsDescription),
			metrics.WithUnit(metrics.Dimensionless),
		)
		if err != nil {
			return err
		}

		numberOfNodes, err = m.Int64ObservableGauge(
			numberOfNodesMetricsName,
			metrics.WithDescription(numberOfNodesMetricsDescription),
			metrics.WithUnit(metrics.Dimensionless),
		)
		if err != nil {
			return err
		}

		numberOfNodesWithoutEdges, err = m.Int64ObservableGauge(
			numberOfNodesWithoutEdgesMetricsName,
			metrics.WithDescription(numberOfNodesWithoutEdgesMetricsDescription),
			metrics.WithUnit(metrics.Dimensionless),
		)
		if err != nil {
			return err
		}

		numberOfNodesWithoutIndegree, err = m.Int64ObservableGauge(
			numberOfNodesWithoutIndegreeMetricsName,
			metrics.WithDescription(numberOfNodesWithoutIndegreeMetricsDescription),
			metrics.WithUnit(metrics.Dimensionless),
		)
		if err != nil {
			return err
		}

		numberOfObjects, err = m.Int64ObservableGauge(
			numberOfObjectsMetricsName,
			metrics.WithDescription(numberOfObjectsMetricsDescription),
			metrics.WithUnit(metrics.Dimensionless),
		)
		if err != nil {
			return err
		}

		numberOfRemovedObjects, err = m.Int64ObservableGauge(
			numberOfRemovedObjectsMetricsName,
			metrics.WithDescription(numberOfRemovedObjectsMetricsDescription),
			metrics.WithUnit(metrics.Dimensionless),
		)
		if err != nil {
			return err
		}

		sizeOfObjectRepository, err = m.Int64ObservableGauge(
			sizeOfObjectRepositoryMetricsName,
			metrics.WithDescription(sizeOfObjectRepositoryMetricsDescription),
			metrics.WithUnit(metrics.Dimensionless),
		)
		if err != nil {
			return err
		}

		sizeOfRefinementObjectRepository, err = m.Int64ObservableGauge(
			sizeOfRefinementObjectRepositoryMetricsName,
			metrics.WithDescription(sizeOfRefinementObjectRepositoryMetricsDescription),
			metrics.WithUnit(metrics.Dimensionless),
		)
		if err != nil {
			return err
		}

		varianceOfIndegree, err = m.Float64ObservableGauge(
			varianceOfIndegreeMetricsName,
			metrics.WithDescription(varianceOfIndegreeMetricsDescription),
			metrics.WithUnit(metrics.Dimensionless),
		)
		if err != nil {
			return err
		}

		varianceOfOutdegree, err = m.Float64ObservableGauge(
			varianceOfOutdegreeMetricsName,
			metrics.WithDescription(varianceOfOutdegreeMetricsDescription),
			metrics.WithUnit(metrics.Dimensionless),
		)
		if err != nil {
			return err
		}

		meanEdgeLength, err = m.Float64ObservableGauge(
			meanEdgeLengthMetricsName,
			metrics.WithDescription(meanEdgeLengthMetricsDescription),
			metrics.WithUnit(metrics.Dimensionless),
		)
		if err != nil {
			return err
		}

		meanEdgeLengthFor10Edges, err = m.Float64ObservableGauge(
			meanEdgeLengthFor10EdgesMetricsName,
			metrics.WithDescription(meanEdgeLengthFor10EdgesMetricsDescription),
			metrics.WithUnit(metrics.Dimensionless),
		)
		if err != nil {
			return err
		}

		meanIndegreeDistanceFor10Edges, err = m.Float64ObservableGauge(
			meanIndegreeDistanceFor10EdgesMetricsName,
			metrics.WithDescription(meanIndegreeDistanceFor10EdgesMetricsDescription),
			metrics.WithUnit(metrics.Dimensionless),
		)
		if err != nil {
			return err
		}

		meanNumberOfEdgesPerNode, err = m.Float64ObservableGauge(
			meanNumberOfEdgesPerNodeMetricsName,
			metrics.WithDescription(meanNumberOfEdgesPerNodeMetricsDescription),
			metrics.WithUnit(metrics.Dimensionless),
		)
		if err != nil {
			return err
		}

		c1Indegree, err = m.Float64ObservableGauge(
			c1IndegreeMetricsName,
			metrics.WithDescription(c1IndegreeMetricsDescription),
			metrics.WithUnit(metrics.Dimensionless),
		)
		if err != nil {
			return err
		}

		c5Indegree, err = m.Float64ObservableGauge(
			c5IndegreeMetricsName,
			metrics.WithDescription(c5IndegreeMetricsDescription),
			metrics.WithUnit(metrics.Dimensionless),
		)
		if err != nil {
			return err
		}

		c95Outdegree, err = m.Float64ObservableGauge(
			c95OutdegreeMetricsName,
			metrics.WithDescription(c95OutdegreeMetricsDescription),
			metrics.WithUnit(metrics.Dimensionless),
		)
		if err != nil {
			return err
		}

		c99Outdegree, err = m.Float64ObservableGauge(
			c99OutdegreeMetricsName,
			metrics.WithDescription(c99OutdegreeMetricsDescription),
			metrics.WithUnit(metrics.Dimensionless),
		)
		if err != nil {
			return err
		}

		indegreeCount, err = m.Float64ObservableGauge(
			indegreeCountMetricsName,
			metrics.WithDescription(indegreeCountMetricsDescription),
			metrics.WithUnit(metrics.Dimensionless),
		)
		if err != nil {
			return err
		}

		outdegreeHistogram, err = m.Float64ObservableGauge(
			outdegreeHistogramMetricsName,
			metrics.WithDescription(outdegreeHistogramMetricsDescription),
			metrics.WithUnit(metrics.Dimensionless),
		)
		if err != nil {
			return err
		}

		indegreeHistogram, err = m.Float64ObservableGauge(
			indegreeHistogramMetricsName,
			metrics.WithDescription(indegreeHistogramMetricsDescription),
			metrics.WithUnit(metrics.Dimensionless),
		)
		if err != nil {
			return err
		}

		instruments = append(instruments,
			medianIndegree,
			medianOutdegree,
			maxNumberOfIndegree,
			maxNumberOfOutdegree,
			minNumberOfIndegree,
			minNumberOfOutdegree,
			modeIndegree,
			modeOutdegree,
			nodesSkippedFor10Edges,
			nodesSkippedForIndegreeDistance,
			numberOfEdges,
			numberOfIndexedObjects,
			numberOfNodes,
			numberOfNodesWithoutEdges,
			numberOfNodesWithoutIndegree,
			numberOfObjects,
			numberOfRemovedObjects,
			sizeOfObjectRepository,
			sizeOfRefinementObjectRepository,
			varianceOfIndegree,
			varianceOfOutdegree,
			meanEdgeLength,
			meanEdgeLengthFor10Edges,
			meanIndegreeDistanceFor10Edges,
			meanNumberOfEdgesPerNode,
			c1Indegree,
			c5Indegree,
			c95Outdegree,
			c99Outdegree,
			indegreeCount,
			outdegreeHistogram,
			indegreeHistogram,
		)
	}

	_, err = m.RegisterCallback(
		func(_ context.Context, o api.Observer) error {
			var indexing int64
			if n.ngt.IsIndexing() {
				indexing = 1
			}
			var saving int64
			if n.ngt.IsSaving() {
				saving = 1
			}
			o.ObserveInt64(indexCount, int64(n.ngt.Len()))
			o.ObserveInt64(uncommittedIndexCount, int64(n.ngt.InsertVQueueBufferLen()+n.ngt.DeleteVQueueBufferLen()))
			o.ObserveInt64(insertVQueueCount, int64(n.ngt.InsertVQueueBufferLen()))
			o.ObserveInt64(deleteVQueueCount, int64(int64(n.ngt.DeleteVQueueBufferLen())))
			o.ObserveInt64(completedCreateIndexTotal, int64(n.ngt.NumberOfCreateIndexExecution()))
			o.ObserveInt64(executedProactiveGCTotal, int64(n.ngt.NumberOfProactiveGCExecution()))
			o.ObserveInt64(isIndexing, int64(indexing))
			o.ObserveInt64(isSaving, int64(saving))
			o.ObserveInt64(brokenIndexCount, int64(n.ngt.BrokenIndexCount()))

			if n.ngt.IsStatisticsEnabled() {
				stats, err := n.ngt.IndexStatistics()
				if err == nil {
					o.ObserveInt64(medianIndegree, int64(stats.GetMedianIndegree()))
					o.ObserveInt64(medianOutdegree, int64(stats.GetMedianOutdegree()))
					o.ObserveInt64(maxNumberOfIndegree, int64(stats.GetMaxNumberOfIndegree()))
					o.ObserveInt64(maxNumberOfOutdegree, int64(stats.GetMaxNumberOfOutdegree()))
					o.ObserveInt64(minNumberOfIndegree, int64(stats.GetMinNumberOfIndegree()))
					o.ObserveInt64(minNumberOfOutdegree, int64(stats.GetMinNumberOfOutdegree()))
					o.ObserveInt64(modeIndegree, int64(stats.GetModeIndegree()))
					o.ObserveInt64(modeOutdegree, int64(stats.GetModeOutdegree()))
					o.ObserveInt64(nodesSkippedFor10Edges, int64(stats.GetNodesSkippedFor10Edges()))
					o.ObserveInt64(nodesSkippedForIndegreeDistance, int64(stats.GetNodesSkippedForIndegreeDistance()))
					o.ObserveInt64(numberOfEdges, int64(stats.GetNumberOfEdges()))
					o.ObserveInt64(numberOfIndexedObjects, int64(stats.GetNumberOfIndexedObjects()))
					o.ObserveInt64(numberOfNodes, int64(stats.GetNumberOfNodes()))
					o.ObserveInt64(numberOfNodesWithoutEdges, int64(stats.GetNumberOfNodesWithoutEdges()))
					o.ObserveInt64(numberOfNodesWithoutIndegree, int64(stats.GetNumberOfNodesWithoutIndegree()))
					o.ObserveInt64(numberOfObjects, int64(stats.GetNumberOfObjects()))
					o.ObserveInt64(numberOfRemovedObjects, int64(stats.GetNumberOfRemovedObjects()))
					o.ObserveInt64(sizeOfObjectRepository, int64(stats.GetSizeOfObjectRepository()))
					o.ObserveInt64(sizeOfRefinementObjectRepository, int64(stats.GetSizeOfRefinementObjectRepository()))
					o.ObserveFloat64(varianceOfIndegree, stats.GetVarianceOfIndegree())
					o.ObserveFloat64(varianceOfOutdegree, stats.GetVarianceOfOutdegree())
					o.ObserveFloat64(meanEdgeLength, stats.GetMeanEdgeLength())
					o.ObserveFloat64(meanEdgeLengthFor10Edges, stats.GetMeanEdgeLengthFor10Edges())
					o.ObserveFloat64(meanIndegreeDistanceFor10Edges, stats.GetMeanIndegreeDistanceFor10Edges())
					o.ObserveFloat64(meanNumberOfEdgesPerNode, stats.GetMeanNumberOfEdgesPerNode())
					o.ObserveFloat64(c1Indegree, stats.GetC1Indegree())
					o.ObserveFloat64(c5Indegree, stats.GetC5Indegree())
					o.ObserveFloat64(c95Outdegree, stats.GetC95Outdegree())
					o.ObserveFloat64(c99Outdegree, stats.GetC99Outdegree())
					// Calculate and observe the average of indegree count
					indegreeCounts := stats.GetIndegreeCount()
					if len(indegreeCounts) > 0 {
						totalIndegreeCount := int64(0)
						for _, count := range indegreeCounts {
							totalIndegreeCount += count
						}
						avgIndegreeCount := float64(totalIndegreeCount) / float64(len(indegreeCounts))
						o.ObserveFloat64(indegreeCount, avgIndegreeCount)
					}

					// Calculate and observe outdegree histogram
					outdegreeHist := stats.GetOutdegreeHistogram()
					for i, count := range outdegreeHist {
						o.ObserveFloat64(outdegreeHistogram, float64(count), api.WithAttributes(attribute.Int64("outdegree", int64(i))))
					}

					// Calculate and observe indegree histogram
					indegreeHist := stats.GetIndegreeHistogram()
					for i, count := range indegreeHist {
						o.ObserveFloat64(indegreeHistogram, float64(count), api.WithAttributes(attribute.Int64("indegree", int64(i))))
					}
				}
			}
			return nil
		},
		instruments...,
	)
	return err
}
