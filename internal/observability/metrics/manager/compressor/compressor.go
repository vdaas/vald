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
	compressorBufferCount       metrics.Int64Measure
	registererBufferCount       metrics.Int64Measure
	registererWorkerBufferCount metrics.Int64Measure
}

func New(c service.Compressor, r service.Registerer) metrics.Metric {
	return &compressorMetrics{
		compressor: c,
		registerer: r,
		compressorBufferCount: *metrics.Int64(
			metrics.ValdOrg+"/manager/compressor/compressor_buffer_count",
			"Compressor compressor buffer count",
			metrics.UnitDimensionless),
		registererBufferCount: *metrics.Int64(
			metrics.ValdOrg+"/manager/compressor/registerer_buffer_count",
			"Compressor registerer buffer count",
			metrics.UnitDimensionless),
		registererWorkerBufferCount: *metrics.Int64(
			metrics.ValdOrg+"/manager/compressor/registerer_worker_buffer_count",
			"Compressor registerer worker buffer count",
			metrics.UnitDimensionless),
	}
}

func (c *compressorMetrics) Measurement(ctx context.Context) ([]metrics.Measurement, error) {
	return []metrics.Measurement{
		c.compressorBufferCount.M(int64(c.compressor.WorkerLen())),
		c.registererBufferCount.M(int64(c.registerer.Len())),
		c.registererWorkerBufferCount.M(int64(c.registerer.WorkerLen())),
	}, nil
}

func (c *compressorMetrics) MeasurementWithTags(ctx context.Context) ([]metrics.MeasurementWithTags, error) {
	return []metrics.MeasurementWithTags{}, nil
}

func (c *compressorMetrics) View() []*metrics.View {
	return []*metrics.View{
		&metrics.View{
			Name:        "compressor_compressor_buffer_count",
			Description: "Compressor compressor buffer count",
			Measure:     &c.compressorBufferCount,
			Aggregation: metrics.LastValue(),
		},
		&metrics.View{
			Name:        "compressor_registerer_buffer_count",
			Description: "Compressor registerer buffer count",
			Measure:     &c.registererBufferCount,
			Aggregation: metrics.LastValue(),
		},
		&metrics.View{
			Name:        "compressor_registerer_worker_buffer_count",
			Description: "Compressor registerer worker buffer count",
			Measure:     &c.registererWorkerBufferCount,
			Aggregation: metrics.LastValue(),
		},
	}
}
