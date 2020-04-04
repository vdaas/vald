//
// Copyright (C) 2019-2020 Vdaas.org Vald team ( kpango, rinx, kmrmt )
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

// Package compressor provides functions for compressor stats
package compressor

import (
	"context"

	"github.com/vdaas/vald/internal/observability/metrics"
	"github.com/vdaas/vald/pkg/manager/compressor/service"
)

type compressorMetrics struct {
	compressor                  service.Compressor
	registerer                  service.Registerer
	compressorBuffer            metrics.Int64Measure
	compressorTotalRequestedJob metrics.Int64Measure
	compressorTotalCompletedJob metrics.Int64Measure
	registererBuffer            metrics.Int64Measure
	registererTotalRequestedJob metrics.Int64Measure
	registererTotalCompletedJob metrics.Int64Measure
}

func New(c service.Compressor, r service.Registerer) metrics.Metric {
	return &compressorMetrics{
		compressor: c,
		registerer: r,
		compressorBuffer: *metrics.Int64(
			metrics.ValdOrg+"/manager/compressor/compressor_buffer",
			"the current number of compressor compress worker buffer elements",
			metrics.UnitDimensionless),
		compressorTotalRequestedJob: *metrics.Int64(
			metrics.ValdOrg+"/manager/compressor/compressor_requested_jobs_total",
			"the cumulative count of compressor compress worker requested job",
			metrics.UnitDimensionless),
		compressorTotalCompletedJob: *metrics.Int64(
			metrics.ValdOrg+"/manager/compressor/compressor_completed_jobs_total",
			"the cumulative count of compressor compress worker completed job",
			metrics.UnitDimensionless),
		registererBuffer: *metrics.Int64(
			metrics.ValdOrg+"/manager/compressor/registerer_buffer",
			"the current number of compressor registerer worker buffer elements",
			metrics.UnitDimensionless),
		registererTotalRequestedJob: *metrics.Int64(
			metrics.ValdOrg+"/manager/compressor/registerer_requested_jobs_total",
			"the cumulative count of compressor registerer worker requested job",
			metrics.UnitDimensionless),
		registererTotalCompletedJob: *metrics.Int64(
			metrics.ValdOrg+"/manager/compressor/registerer_completed_jobs_total",
			"the cumulative count of compressor registerer worker completed job",
			metrics.UnitDimensionless),
	}
}

func (c *compressorMetrics) Measurement(ctx context.Context) ([]metrics.Measurement, error) {
	return []metrics.Measurement{
		c.compressorBuffer.M(int64(c.compressor.Len())),
		c.compressorTotalRequestedJob.M(int64(c.compressor.TotalRequested())),
		c.compressorTotalCompletedJob.M(int64(c.compressor.TotalCompleted())),
		c.registererBuffer.M(int64(c.registerer.Len())),
		c.registererTotalRequestedJob.M(int64(c.registerer.TotalRequested())),
		c.registererTotalCompletedJob.M(int64(c.registerer.TotalCompleted())),
	}, nil
}

func (c *compressorMetrics) MeasurementWithTags(ctx context.Context) ([]metrics.MeasurementWithTags, error) {
	return []metrics.MeasurementWithTags{}, nil
}

func (c *compressorMetrics) View() []*metrics.View {
	return []*metrics.View{
		&metrics.View{
			Name:        "compressor_compressor_buffer",
			Description: "the current number of compressor compress worker buffer elements",
			Measure:     &c.compressorBuffer,
			Aggregation: metrics.LastValue(),
		},
		&metrics.View{
			Name:        "compressor_compressor_requested_jobs_total",
			Description: "the cumulative count of compressor compress worker requested job",
			Measure:     &c.compressorTotalRequestedJob,
			Aggregation: metrics.Count(),
		},
		&metrics.View{
			Name:        "compressor_compressor_completed_jobs_total",
			Description: "the cumulative count of compressor compress worker completed job",
			Measure:     &c.compressorTotalCompletedJob,
			Aggregation: metrics.Count(),
		},
		&metrics.View{
			Name:        "compressor_registerer_buffer",
			Description: "the current number of compressor registerer worker buffer elements",
			Measure:     &c.registererBuffer,
			Aggregation: metrics.LastValue(),
		},
		&metrics.View{
			Name:        "compressor_registerer_requested_jobs_total",
			Description: "the cumulative count of compressor registerer worker requested job",
			Measure:     &c.registererTotalRequestedJob,
			Aggregation: metrics.Count(),
		},
		&metrics.View{
			Name:        "compressor_registerer_completed_jobs_total",
			Description: "the cumulative count of compressor registerer worker completed job",
			Measure:     &c.registererTotalCompletedJob,
			Aggregation: metrics.Count(),
		},
	}
}
