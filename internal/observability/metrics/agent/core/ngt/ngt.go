// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
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
	IndexCountMetricsName        = "agent_core_ngt_index_count"
	IndexCountMetricsDescription = "Agent NGT index count"

	UncommittedIndexCountMetricsName        = "agent_core_ngt_uncommitted_index_count"
	UncommittedIndexCountMetricsDescription = "Agent NGT uncommitted index count"

	InsertVQueueCountMetricsName        = "agent_core_ngt_insert_vqueue_count"
	InsertVQueueCountMetricsDescription = "Agent NGT insert vqueue count"

	DeleteVQueueCountMetricsName        = "agent_core_ngt_delete_vqueue_count"
	DeleteVQueueCountMetricsDescription = "Agent NGT delete vqueue count"

	CompletedCreateIndexTotalMetricsName        = "agent_core_ngt_completed_create_index_total"
	CompletedCreateIndexTotalMetricsDescription = "The cumulative count of completed create index execution"

	ExecutedProactiveGCTotalMetricsName        = "agent_core_ngt_executed_proactive_gc_total"
	ExecutedProactiveGCTotalMetricsDescription = "The cumulative count of proactive GC execution"

	IsIndexingMetricsName        = "agent_core_ngt_is_indexing"
	IsIndexingMetricsDescription = "Currently indexing or no"

	IsSavingMetricsName        = "agent_core_ngt_is_saving"
	IsSavingMetricsDescription = "Currently saving or not"

	BrokenIndexStoreCountMetricsName        = "agent_core_ngt_broken_index_store_count"
	BrokenIndexStoreCountMetricsDescription = "How many broken index generations have been stored"

	MedianIndegreeMetricsName        = "agent_core_ngt_median_indegree"
	MedianIndegreeMetricsDescription = "Median indegree of nodes"

	MedianOutdegreeMetricsName        = "agent_core_ngt_median_outdegree"
	MedianOutdegreeMetricsDescription = "Median outdegree of nodes"

	MaxNumberOfIndegreeMetricsName        = "agent_core_ngt_max_number_of_indegree"
	MaxNumberOfIndegreeMetricsDescription = "Maximum number of indegree"

	MaxNumberOfOutdegreeMetricsName        = "agent_core_ngt_max_number_of_outdegree"
	MaxNumberOfOutdegreeMetricsDescription = "Maximum number of outdegree"

	MinNumberOfIndegreeMetricsName        = "agent_core_ngt_min_number_of_indegree"
	MinNumberOfIndegreeMetricsDescription = "Minimum number of indegree"

	MinNumberOfOutdegreeMetricsName        = "agent_core_ngt_min_number_of_outdegree"
	MinNumberOfOutdegreeMetricsDescription = "Minimum number of outdegree"

	ModeIndegreeMetricsName        = "agent_core_ngt_mode_indegree"
	ModeIndegreeMetricsDescription = "Mode of indegree"

	ModeOutdegreeMetricsName        = "agent_core_ngt_mode_outdegree"
	ModeOutdegreeMetricsDescription = "Mode of outdegree"

	NodesSkippedFor10EdgesMetricsName        = "agent_core_ngt_nodes_skipped_for_10_edges"
	NodesSkippedFor10EdgesMetricsDescription = "Nodes skipped for 10 edges"

	NodesSkippedForIndegreeDistanceMetricsName        = "agent_core_ngt_nodes_skipped_for_indegree_distance"
	NodesSkippedForIndegreeDistanceMetricsDescription = "Nodes skipped for indegree distance"

	NumberOfEdgesMetricsName        = "agent_core_ngt_number_of_edges"
	NumberOfEdgesMetricsDescription = "Number of edges"

	NumberOfIndexedObjectsMetricsName        = "agent_core_ngt_number_of_indexed_objects"
	NumberOfIndexedObjectsMetricsDescription = "Number of indexed objects"

	NumberOfNodesMetricsName        = "agent_core_ngt_number_of_nodes"
	NumberOfNodesMetricsDescription = "Number of nodes"

	NumberOfNodesWithoutEdgesMetricsName        = "agent_core_ngt_number_of_nodes_without_edges"
	NumberOfNodesWithoutEdgesMetricsDescription = "Number of nodes without edges"

	NumberOfNodesWithoutIndegreeMetricsName        = "agent_core_ngt_number_of_nodes_without_indegree"
	NumberOfNodesWithoutIndegreeMetricsDescription = "Number of nodes without indegree"

	NumberOfObjectsMetricsName        = "agent_core_ngt_number_of_objects"
	NumberOfObjectsMetricsDescription = "Number of objects"

	NumberOfRemovedObjectsMetricsName        = "agent_core_ngt_number_of_removed_objects"
	NumberOfRemovedObjectsMetricsDescription = "Number of removed objects"

	SizeOfObjectRepositoryMetricsName        = "agent_core_ngt_size_of_object_repository"
	SizeOfObjectRepositoryMetricsDescription = "Size of object repository"

	SizeOfRefinementObjectRepositoryMetricsName        = "agent_core_ngt_size_of_refinement_object_repository"
	SizeOfRefinementObjectRepositoryMetricsDescription = "Size of refinement object repository"

	VarianceOfIndegreeMetricsName        = "agent_core_ngt_variance_of_indegree"
	VarianceOfIndegreeMetricsDescription = "Variance of indegree"

	VarianceOfOutdegreeMetricsName        = "agent_core_ngt_variance_of_outdegree"
	VarianceOfOutdegreeMetricsDescription = "Variance of outdegree"

	MeanEdgeLengthMetricsName        = "agent_core_ngt_mean_edge_length"
	MeanEdgeLengthMetricsDescription = "Mean edge length"

	MeanEdgeLengthFor10EdgesMetricsName        = "agent_core_ngt_mean_edge_length_for_10_edges"
	MeanEdgeLengthFor10EdgesMetricsDescription = "Mean edge length for 10 edges"

	MeanIndegreeDistanceFor10EdgesMetricsName        = "agent_core_ngt_mean_indegree_distance_for_10_edges"
	MeanIndegreeDistanceFor10EdgesMetricsDescription = "Mean indegree distance for 10 edges"

	MeanNumberOfEdgesPerNodeMetricsName        = "agent_core_ngt_mean_number_of_edges_per_node"
	MeanNumberOfEdgesPerNodeMetricsDescription = "Mean number of edges per node"

	C1IndegreeMetricsName        = "agent_core_ngt_c1_indegree"
	C1IndegreeMetricsDescription = "C1 indegree"

	C5IndegreeMetricsName        = "agent_core_ngt_c5_indegree"
	C5IndegreeMetricsDescription = "C5 indegree"

	C95OutdegreeMetricsName        = "agent_core_ngt_c95_outdegree"
	C95OutdegreeMetricsDescription = "C95 outdegree"

	C99OutdegreeMetricsName        = "agent_core_ngt_c99_outdegree"
	C99OutdegreeMetricsDescription = "C99 outdegree"

	IndegreeCountMetricsName        = "agent_core_ngt_indegree_count"
	IndegreeCountMetricsDescription = "Indegree count"

	OutdegreeHistogramMetricsName        = "agent_core_ngt_outdegree_histogram"
	OutdegreeHistogramMetricsDescription = "Outdegree histogram"

	IndegreeHistogramMetricsName        = "agent_core_ngt_indegree_histogram"
	IndegreeHistogramMetricsDescription = "Indegree histogram"
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
				Name:        IndexCountMetricsName,
				Description: IndexCountMetricsDescription,
			},
			view.Stream{
				Aggregation: view.AggregationLastValue{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        UncommittedIndexCountMetricsName,
				Description: UncommittedIndexCountMetricsDescription,
			},
			view.Stream{
				Aggregation: view.AggregationLastValue{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        InsertVQueueCountMetricsName,
				Description: InsertVQueueCountMetricsDescription,
			},
			view.Stream{
				Aggregation: view.AggregationLastValue{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        DeleteVQueueCountMetricsName,
				Description: DeleteVQueueCountMetricsDescription,
			},
			view.Stream{
				Aggregation: view.AggregationLastValue{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        CompletedCreateIndexTotalMetricsName,
				Description: CompletedCreateIndexTotalMetricsDescription,
			},
			view.Stream{
				Aggregation: view.AggregationLastValue{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        ExecutedProactiveGCTotalMetricsName,
				Description: ExecutedProactiveGCTotalMetricsDescription,
			},
			view.Stream{
				Aggregation: view.AggregationLastValue{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        IsIndexingMetricsName,
				Description: IsIndexingMetricsDescription,
			},
			view.Stream{
				Aggregation: view.AggregationLastValue{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        IsSavingMetricsName,
				Description: IsSavingMetricsDescription,
			},
			view.Stream{
				Aggregation: view.AggregationLastValue{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        BrokenIndexStoreCountMetricsName,
				Description: BrokenIndexStoreCountMetricsDescription,
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
					Name:        MedianIndegreeMetricsName,
					Description: MedianIndegreeMetricsDescription,
				},
				view.Stream{
					Aggregation: view.AggregationLastValue{},
				},
			),
			view.NewView(
				view.Instrument{
					Name:        MedianOutdegreeMetricsName,
					Description: MedianOutdegreeMetricsDescription,
				},
				view.Stream{
					Aggregation: view.AggregationLastValue{},
				},
			),
			view.NewView(
				view.Instrument{
					Name:        MaxNumberOfIndegreeMetricsName,
					Description: MaxNumberOfIndegreeMetricsDescription,
				},
				view.Stream{
					Aggregation: view.AggregationLastValue{},
				},
			),
			view.NewView(
				view.Instrument{
					Name:        MaxNumberOfOutdegreeMetricsName,
					Description: MaxNumberOfOutdegreeMetricsDescription,
				},
				view.Stream{
					Aggregation: view.AggregationLastValue{},
				},
			),
			view.NewView(
				view.Instrument{
					Name:        MinNumberOfIndegreeMetricsName,
					Description: MinNumberOfIndegreeMetricsDescription,
				},
				view.Stream{
					Aggregation: view.AggregationLastValue{},
				},
			),
			view.NewView(
				view.Instrument{
					Name:        MinNumberOfOutdegreeMetricsName,
					Description: MinNumberOfOutdegreeMetricsDescription,
				},
				view.Stream{
					Aggregation: view.AggregationLastValue{},
				},
			),
			view.NewView(
				view.Instrument{
					Name:        ModeIndegreeMetricsName,
					Description: ModeIndegreeMetricsDescription,
				},
				view.Stream{
					Aggregation: view.AggregationLastValue{},
				},
			),
			view.NewView(
				view.Instrument{
					Name:        ModeOutdegreeMetricsName,
					Description: ModeOutdegreeMetricsDescription,
				},
				view.Stream{
					Aggregation: view.AggregationLastValue{},
				},
			),
			view.NewView(
				view.Instrument{
					Name:        NodesSkippedFor10EdgesMetricsName,
					Description: NodesSkippedFor10EdgesMetricsDescription,
				},
				view.Stream{
					Aggregation: view.AggregationLastValue{},
				},
			),
			view.NewView(
				view.Instrument{
					Name:        NodesSkippedForIndegreeDistanceMetricsName,
					Description: NodesSkippedForIndegreeDistanceMetricsDescription,
				},
				view.Stream{
					Aggregation: view.AggregationLastValue{},
				},
			),
			view.NewView(
				view.Instrument{
					Name:        NumberOfEdgesMetricsName,
					Description: NumberOfEdgesMetricsDescription,
				},
				view.Stream{
					Aggregation: view.AggregationLastValue{},
				},
			),
			view.NewView(
				view.Instrument{
					Name:        NumberOfIndexedObjectsMetricsName,
					Description: NumberOfIndexedObjectsMetricsDescription,
				},
				view.Stream{
					Aggregation: view.AggregationLastValue{},
				},
			),
			view.NewView(
				view.Instrument{
					Name:        NumberOfNodesMetricsName,
					Description: NumberOfNodesMetricsDescription,
				},
				view.Stream{
					Aggregation: view.AggregationLastValue{},
				},
			),
			view.NewView(
				view.Instrument{
					Name:        NumberOfNodesWithoutEdgesMetricsName,
					Description: NumberOfNodesWithoutEdgesMetricsDescription,
				},
				view.Stream{
					Aggregation: view.AggregationLastValue{},
				},
			),
			view.NewView(
				view.Instrument{
					Name:        NumberOfNodesWithoutIndegreeMetricsName,
					Description: NumberOfNodesWithoutIndegreeMetricsDescription,
				},
				view.Stream{
					Aggregation: view.AggregationLastValue{},
				},
			),
			view.NewView(
				view.Instrument{
					Name:        NumberOfObjectsMetricsName,
					Description: NumberOfObjectsMetricsDescription,
				},
				view.Stream{
					Aggregation: view.AggregationLastValue{},
				},
			),
			view.NewView(
				view.Instrument{
					Name:        NumberOfRemovedObjectsMetricsName,
					Description: NumberOfRemovedObjectsMetricsDescription,
				},
				view.Stream{
					Aggregation: view.AggregationLastValue{},
				},
			),
			view.NewView(
				view.Instrument{
					Name:        SizeOfObjectRepositoryMetricsName,
					Description: SizeOfObjectRepositoryMetricsDescription,
				},
				view.Stream{
					Aggregation: view.AggregationLastValue{},
				},
			),
			view.NewView(
				view.Instrument{
					Name:        SizeOfRefinementObjectRepositoryMetricsName,
					Description: SizeOfRefinementObjectRepositoryMetricsDescription,
				},
				view.Stream{
					Aggregation: view.AggregationLastValue{},
				},
			),
			view.NewView(
				view.Instrument{
					Name:        VarianceOfIndegreeMetricsName,
					Description: VarianceOfIndegreeMetricsDescription,
				},
				view.Stream{
					Aggregation: view.AggregationLastValue{},
				},
			),
			view.NewView(
				view.Instrument{
					Name:        VarianceOfOutdegreeMetricsName,
					Description: VarianceOfOutdegreeMetricsDescription,
				},
				view.Stream{
					Aggregation: view.AggregationLastValue{},
				},
			),
			view.NewView(
				view.Instrument{
					Name:        MeanEdgeLengthMetricsName,
					Description: MeanEdgeLengthMetricsDescription,
				},
				view.Stream{
					Aggregation: view.AggregationLastValue{},
				},
			),
			view.NewView(
				view.Instrument{
					Name:        MeanEdgeLengthFor10EdgesMetricsName,
					Description: MeanEdgeLengthFor10EdgesMetricsDescription,
				},
				view.Stream{
					Aggregation: view.AggregationLastValue{},
				},
			),
			view.NewView(
				view.Instrument{
					Name:        MeanIndegreeDistanceFor10EdgesMetricsName,
					Description: MeanIndegreeDistanceFor10EdgesMetricsDescription,
				},
				view.Stream{
					Aggregation: view.AggregationLastValue{},
				},
			),
			view.NewView(
				view.Instrument{
					Name:        MeanNumberOfEdgesPerNodeMetricsName,
					Description: MeanNumberOfEdgesPerNodeMetricsDescription,
				},
				view.Stream{
					Aggregation: view.AggregationLastValue{},
				},
			),
			view.NewView(
				view.Instrument{
					Name:        C1IndegreeMetricsName,
					Description: C1IndegreeMetricsDescription,
				},
				view.Stream{
					Aggregation: view.AggregationLastValue{},
				},
			),
			view.NewView(
				view.Instrument{
					Name:        C5IndegreeMetricsName,
					Description: C5IndegreeMetricsDescription,
				},
				view.Stream{
					Aggregation: view.AggregationLastValue{},
				},
			),
			view.NewView(
				view.Instrument{
					Name:        C95OutdegreeMetricsName,
					Description: C95OutdegreeMetricsDescription,
				},
				view.Stream{
					Aggregation: view.AggregationLastValue{},
				},
			),
			view.NewView(
				view.Instrument{
					Name:        C99OutdegreeMetricsName,
					Description: C99OutdegreeMetricsDescription,
				},
				view.Stream{
					Aggregation: view.AggregationLastValue{},
				},
			),
			view.NewView(
				view.Instrument{
					Name:        IndegreeCountMetricsName,
					Description: IndegreeCountMetricsDescription,
				},
				view.Stream{
					Aggregation: view.AggregationLastValue{},
				},
			),
			view.NewView(
				view.Instrument{
					Name:        OutdegreeHistogramMetricsName,
					Description: OutdegreeHistogramMetricsDescription,
				},
				view.Stream{
					Aggregation: view.AggregationLastValue{},
				},
			),
			view.NewView(
				view.Instrument{
					Name:        IndegreeHistogramMetricsName,
					Description: IndegreeHistogramMetricsDescription,
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
		IndexCountMetricsName,
		metrics.WithDescription(IndexCountMetricsDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	uncommittedIndexCount, err = m.Int64ObservableGauge(
		UncommittedIndexCountMetricsName,
		metrics.WithDescription(UncommittedIndexCountMetricsDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	insertVQueueCount, err = m.Int64ObservableGauge(
		InsertVQueueCountMetricsName,
		metrics.WithDescription(InsertVQueueCountMetricsDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	deleteVQueueCount, err = m.Int64ObservableGauge(
		DeleteVQueueCountMetricsName,
		metrics.WithDescription(DeleteVQueueCountMetricsDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	completedCreateIndexTotal, err = m.Int64ObservableGauge(
		CompletedCreateIndexTotalMetricsName,
		metrics.WithDescription(CompletedCreateIndexTotalMetricsDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	executedProactiveGCTotal, err = m.Int64ObservableGauge(
		ExecutedProactiveGCTotalMetricsName,
		metrics.WithDescription(ExecutedProactiveGCTotalMetricsDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	isIndexing, err = m.Int64ObservableGauge(
		IsIndexingMetricsName,
		metrics.WithDescription(IsIndexingMetricsDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	isSaving, err = m.Int64ObservableGauge(
		IsSavingMetricsName,
		metrics.WithDescription(IsSavingMetricsDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	brokenIndexCount, err = m.Int64ObservableGauge(
		BrokenIndexStoreCountMetricsName,
		metrics.WithDescription(BrokenIndexStoreCountMetricsDescription),
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
			MedianIndegreeMetricsName,
			metrics.WithDescription(MedianIndegreeMetricsDescription),
			metrics.WithUnit(metrics.Dimensionless),
		)
		if err != nil {
			return err
		}

		medianOutdegree, err = m.Int64ObservableGauge(
			MedianOutdegreeMetricsName,
			metrics.WithDescription(MedianOutdegreeMetricsDescription),
			metrics.WithUnit(metrics.Dimensionless),
		)
		if err != nil {
			return err
		}

		maxNumberOfIndegree, err = m.Int64ObservableGauge(
			MaxNumberOfIndegreeMetricsName,
			metrics.WithDescription(MaxNumberOfIndegreeMetricsDescription),
			metrics.WithUnit(metrics.Dimensionless),
		)
		if err != nil {
			return err
		}

		maxNumberOfOutdegree, err = m.Int64ObservableGauge(
			MaxNumberOfOutdegreeMetricsName,
			metrics.WithDescription(MaxNumberOfOutdegreeMetricsDescription),
			metrics.WithUnit(metrics.Dimensionless),
		)
		if err != nil {
			return err
		}

		minNumberOfIndegree, err = m.Int64ObservableGauge(
			MinNumberOfIndegreeMetricsName,
			metrics.WithDescription(MinNumberOfIndegreeMetricsDescription),
			metrics.WithUnit(metrics.Dimensionless),
		)
		if err != nil {
			return err
		}

		minNumberOfOutdegree, err = m.Int64ObservableGauge(
			MinNumberOfOutdegreeMetricsName,
			metrics.WithDescription(MinNumberOfOutdegreeMetricsDescription),
			metrics.WithUnit(metrics.Dimensionless),
		)
		if err != nil {
			return err
		}

		modeIndegree, err = m.Int64ObservableGauge(
			ModeIndegreeMetricsName,
			metrics.WithDescription(ModeIndegreeMetricsDescription),
			metrics.WithUnit(metrics.Dimensionless),
		)
		if err != nil {
			return err
		}

		modeOutdegree, err = m.Int64ObservableGauge(
			ModeOutdegreeMetricsName,
			metrics.WithDescription(ModeOutdegreeMetricsDescription),
			metrics.WithUnit(metrics.Dimensionless),
		)
		if err != nil {
			return err
		}

		nodesSkippedFor10Edges, err = m.Int64ObservableGauge(
			NodesSkippedFor10EdgesMetricsName,
			metrics.WithDescription(NodesSkippedFor10EdgesMetricsDescription),
			metrics.WithUnit(metrics.Dimensionless),
		)
		if err != nil {
			return err
		}

		nodesSkippedForIndegreeDistance, err = m.Int64ObservableGauge(
			NodesSkippedForIndegreeDistanceMetricsName,
			metrics.WithDescription(NodesSkippedForIndegreeDistanceMetricsDescription),
			metrics.WithUnit(metrics.Dimensionless),
		)
		if err != nil {
			return err
		}

		numberOfEdges, err = m.Int64ObservableGauge(
			NumberOfEdgesMetricsName,
			metrics.WithDescription(NumberOfEdgesMetricsDescription),
			metrics.WithUnit(metrics.Dimensionless),
		)
		if err != nil {
			return err
		}

		numberOfIndexedObjects, err = m.Int64ObservableGauge(
			NumberOfIndexedObjectsMetricsName,
			metrics.WithDescription(NumberOfIndexedObjectsMetricsDescription),
			metrics.WithUnit(metrics.Dimensionless),
		)
		if err != nil {
			return err
		}

		numberOfNodes, err = m.Int64ObservableGauge(
			NumberOfNodesMetricsName,
			metrics.WithDescription(NumberOfNodesMetricsDescription),
			metrics.WithUnit(metrics.Dimensionless),
		)
		if err != nil {
			return err
		}

		numberOfNodesWithoutEdges, err = m.Int64ObservableGauge(
			NumberOfNodesWithoutEdgesMetricsName,
			metrics.WithDescription(NumberOfNodesWithoutEdgesMetricsDescription),
			metrics.WithUnit(metrics.Dimensionless),
		)
		if err != nil {
			return err
		}

		numberOfNodesWithoutIndegree, err = m.Int64ObservableGauge(
			NumberOfNodesWithoutIndegreeMetricsName,
			metrics.WithDescription(NumberOfNodesWithoutIndegreeMetricsDescription),
			metrics.WithUnit(metrics.Dimensionless),
		)
		if err != nil {
			return err
		}

		numberOfObjects, err = m.Int64ObservableGauge(
			NumberOfObjectsMetricsName,
			metrics.WithDescription(NumberOfObjectsMetricsDescription),
			metrics.WithUnit(metrics.Dimensionless),
		)
		if err != nil {
			return err
		}

		numberOfRemovedObjects, err = m.Int64ObservableGauge(
			NumberOfRemovedObjectsMetricsName,
			metrics.WithDescription(NumberOfRemovedObjectsMetricsDescription),
			metrics.WithUnit(metrics.Dimensionless),
		)
		if err != nil {
			return err
		}

		sizeOfObjectRepository, err = m.Int64ObservableGauge(
			SizeOfObjectRepositoryMetricsName,
			metrics.WithDescription(SizeOfObjectRepositoryMetricsDescription),
			metrics.WithUnit(metrics.Dimensionless),
		)
		if err != nil {
			return err
		}

		sizeOfRefinementObjectRepository, err = m.Int64ObservableGauge(
			SizeOfRefinementObjectRepositoryMetricsName,
			metrics.WithDescription(SizeOfRefinementObjectRepositoryMetricsDescription),
			metrics.WithUnit(metrics.Dimensionless),
		)
		if err != nil {
			return err
		}

		varianceOfIndegree, err = m.Float64ObservableGauge(
			VarianceOfIndegreeMetricsName,
			metrics.WithDescription(VarianceOfIndegreeMetricsDescription),
			metrics.WithUnit(metrics.Dimensionless),
		)
		if err != nil {
			return err
		}

		varianceOfOutdegree, err = m.Float64ObservableGauge(
			VarianceOfOutdegreeMetricsName,
			metrics.WithDescription(VarianceOfOutdegreeMetricsDescription),
			metrics.WithUnit(metrics.Dimensionless),
		)
		if err != nil {
			return err
		}

		meanEdgeLength, err = m.Float64ObservableGauge(
			MeanEdgeLengthMetricsName,
			metrics.WithDescription(MeanEdgeLengthMetricsDescription),
			metrics.WithUnit(metrics.Dimensionless),
		)
		if err != nil {
			return err
		}

		meanEdgeLengthFor10Edges, err = m.Float64ObservableGauge(
			MeanEdgeLengthFor10EdgesMetricsName,
			metrics.WithDescription(MeanEdgeLengthFor10EdgesMetricsDescription),
			metrics.WithUnit(metrics.Dimensionless),
		)
		if err != nil {
			return err
		}

		meanIndegreeDistanceFor10Edges, err = m.Float64ObservableGauge(
			MeanIndegreeDistanceFor10EdgesMetricsName,
			metrics.WithDescription(MeanIndegreeDistanceFor10EdgesMetricsDescription),
			metrics.WithUnit(metrics.Dimensionless),
		)
		if err != nil {
			return err
		}

		meanNumberOfEdgesPerNode, err = m.Float64ObservableGauge(
			MeanNumberOfEdgesPerNodeMetricsName,
			metrics.WithDescription(MeanNumberOfEdgesPerNodeMetricsDescription),
			metrics.WithUnit(metrics.Dimensionless),
		)
		if err != nil {
			return err
		}

		c1Indegree, err = m.Float64ObservableGauge(
			C1IndegreeMetricsName,
			metrics.WithDescription(C1IndegreeMetricsDescription),
			metrics.WithUnit(metrics.Dimensionless),
		)
		if err != nil {
			return err
		}

		c5Indegree, err = m.Float64ObservableGauge(
			C5IndegreeMetricsName,
			metrics.WithDescription(C5IndegreeMetricsDescription),
			metrics.WithUnit(metrics.Dimensionless),
		)
		if err != nil {
			return err
		}

		c95Outdegree, err = m.Float64ObservableGauge(
			C95OutdegreeMetricsName,
			metrics.WithDescription(C95OutdegreeMetricsDescription),
			metrics.WithUnit(metrics.Dimensionless),
		)
		if err != nil {
			return err
		}

		c99Outdegree, err = m.Float64ObservableGauge(
			C99OutdegreeMetricsName,
			metrics.WithDescription(C99OutdegreeMetricsDescription),
			metrics.WithUnit(metrics.Dimensionless),
		)
		if err != nil {
			return err
		}

		indegreeCount, err = m.Float64ObservableGauge(
			IndegreeCountMetricsName,
			metrics.WithDescription(IndegreeCountMetricsDescription),
			metrics.WithUnit(metrics.Dimensionless),
		)
		if err != nil {
			return err
		}

		outdegreeHistogram, err = m.Float64ObservableGauge(
			OutdegreeHistogramMetricsName,
			metrics.WithDescription(OutdegreeHistogramMetricsDescription),
			metrics.WithUnit(metrics.Dimensionless),
		)
		if err != nil {
			return err
		}

		indegreeHistogram, err = m.Float64ObservableGauge(
			IndegreeHistogramMetricsName,
			metrics.WithDescription(IndegreeHistogramMetricsDescription),
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
