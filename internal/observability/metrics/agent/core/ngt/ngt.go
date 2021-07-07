//
// Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
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

// Package ngt provides functions for ngt stats
package ngt

import (
	"context"

	"github.com/vdaas/vald/internal/observability/metrics"
	"github.com/vdaas/vald/pkg/agent/core/ngt/service"
)

type ngtMetrics struct {
	ngt                            service.NGT
	indexCount                     metrics.Int64Measure
	uncommittedIndexCount          metrics.Int64Measure
	insertVQueueCount              metrics.Int64Measure
	deleteVQueueCount              metrics.Int64Measure
	insertVQueueChannelBufferCount metrics.Int64Measure
	deleteVQueueChannelBufferCount metrics.Int64Measure
	completedCreateIndexTotal      metrics.Int64Measure
	executedProactiveGCTotal       metrics.Int64Measure
	isIndexing                     metrics.Int64Measure
	isSaving                       metrics.Int64Measure
}

func New(n service.NGT) metrics.Metric {
	return &ngtMetrics{
		ngt: n,
		indexCount: *metrics.Int64(
			metrics.ValdOrg+"/agent/core/ngt/index_count",
			"Agent NGT index count",
			metrics.UnitDimensionless),
		uncommittedIndexCount: *metrics.Int64(
			metrics.ValdOrg+"/agent/core/ngt/uncommitted_index_count",
			"Agent NGT uncommitted index count",
			metrics.UnitDimensionless),
		insertVQueueCount: *metrics.Int64(
			metrics.ValdOrg+"/agent/core/ngt/insert_vqueue_count",
			"Agent NGT insert vqueue count",
			metrics.UnitDimensionless),
		deleteVQueueCount: *metrics.Int64(
			metrics.ValdOrg+"/agent/core/ngt/delete_vqueue_count",
			"Agent NGT delete vqueue count",
			metrics.UnitDimensionless),
		insertVQueueChannelBufferCount: *metrics.Int64(
			metrics.ValdOrg+"/agent/core/ngt/insert_vqueue_channel_buffer_count",
			"Agent NGT insert vqueue channel buffer count",
			metrics.UnitDimensionless),
		deleteVQueueChannelBufferCount: *metrics.Int64(
			metrics.ValdOrg+"/agent/core/ngt/delete_vqueue_channel_buffer_count",
			"Agent NGT delete vqueue channel buffer count",
			metrics.UnitDimensionless),
		completedCreateIndexTotal: *metrics.Int64(
			metrics.ValdOrg+"/agent/core/ngt/completed_create_index_total",
			"the cumulative count of completed create index execution",
			metrics.UnitDimensionless),
		executedProactiveGCTotal: *metrics.Int64(
			metrics.ValdOrg+"/agent/core/ngt/executed_proactive_gc_total",
			"the cumulative count of proactive GC execution",
			metrics.UnitDimensionless),
		isIndexing: *metrics.Int64(
			metrics.ValdOrg+"/agent/core/ngt/is_indexing",
			"currently indexing or not",
			metrics.UnitDimensionless),
		isSaving: *metrics.Int64(
			metrics.ValdOrg+"/agent/core/ngt/is_saving",
			"currently saving or not",
			metrics.UnitDimensionless),
	}
}

func (n *ngtMetrics) Measurement(ctx context.Context) ([]metrics.Measurement, error) {
	var isIndexing int64
	if n.ngt.IsIndexing() {
		isIndexing = 1
	}

	var isSaving int64
	if n.ngt.IsSaving() {
		isSaving = 1
	}

	return []metrics.Measurement{
		n.indexCount.M(int64(n.ngt.Len())),
		n.uncommittedIndexCount.M(
			int64(n.ngt.InsertVQueueBufferLen() + n.ngt.DeleteVQueueBufferLen()),
		),
		n.insertVQueueCount.M(int64(n.ngt.InsertVQueueBufferLen())),
		n.deleteVQueueCount.M(int64(n.ngt.DeleteVQueueBufferLen())),
		n.insertVQueueChannelBufferCount.M(int64(n.ngt.InsertVQueueChannelLen())),
		n.deleteVQueueChannelBufferCount.M(int64(n.ngt.DeleteVQueueChannelLen())),
		n.completedCreateIndexTotal.M(int64(n.ngt.NumberOfCreateIndexExecution())),
		n.executedProactiveGCTotal.M(int64(n.ngt.NumberOfProactiveGCExecution())),
		n.isIndexing.M(isIndexing),
		n.isSaving.M(isSaving),
	}, nil
}

func (n *ngtMetrics) MeasurementWithTags(
	ctx context.Context,
) ([]metrics.MeasurementWithTags, error) {
	return []metrics.MeasurementWithTags{}, nil
}

func (n *ngtMetrics) View() []*metrics.View {
	return []*metrics.View{
		{
			Name:        "agent_core_ngt_index_count",
			Description: n.indexCount.Description(),
			Measure:     &n.indexCount,
			Aggregation: metrics.LastValue(),
		},
		{
			Name:        "agent_core_ngt_uncommitted_index_count",
			Description: n.uncommittedIndexCount.Description(),
			Measure:     &n.uncommittedIndexCount,
			Aggregation: metrics.LastValue(),
		},
		{
			Name:        "agent_core_ngt_insert_vqueue_count",
			Description: n.insertVQueueCount.Description(),
			Measure:     &n.insertVQueueCount,
			Aggregation: metrics.LastValue(),
		},
		{
			Name:        "agent_core_ngt_delete_vqueue_count",
			Description: n.deleteVQueueCount.Description(),
			Measure:     &n.deleteVQueueCount,
			Aggregation: metrics.LastValue(),
		},
		{
			Name:        "agent_core_ngt_insert_vqueue_channel_buffer_count",
			Description: n.insertVQueueChannelBufferCount.Description(),
			Measure:     &n.insertVQueueChannelBufferCount,
			Aggregation: metrics.LastValue(),
		},
		{
			Name:        "agent_core_ngt_delete_vqueue_channel_buffer_count",
			Description: n.deleteVQueueChannelBufferCount.Description(),
			Measure:     &n.deleteVQueueChannelBufferCount,
			Aggregation: metrics.LastValue(),
		},
		{
			Name:        "agent_core_ngt_completed_create_index_total",
			Description: n.completedCreateIndexTotal.Description(),
			Measure:     &n.completedCreateIndexTotal,
			Aggregation: metrics.LastValue(),
		},
		{
			Name:        "agent_core_ngt_executed_proactive_gc_total",
			Description: n.executedProactiveGCTotal.Description(),
			Measure:     &n.executedProactiveGCTotal,
			Aggregation: metrics.LastValue(),
		},
		{
			Name:        "agent_core_ngt_is_indexing",
			Description: n.isIndexing.Description(),
			Measure:     &n.isIndexing,
			Aggregation: metrics.LastValue(),
		},
		{
			Name:        "agent_core_ngt_is_saving",
			Description: n.isSaving.Description(),
			Measure:     &n.isSaving,
			Aggregation: metrics.LastValue(),
		},
	}
}
