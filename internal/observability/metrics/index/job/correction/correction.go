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
package correction

import (
	"context"

	"github.com/vdaas/vald/internal/observability/metrics"
	"github.com/vdaas/vald/pkg/index/job/correction/service"
	"go.opentelemetry.io/otel/sdk/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/view"
)

const (
	checkedIndexCount     = "index_job_correction_checked_index_count"
	checkedIndexCountDesc = "The number of checked indexes while index correction job"

	correctedOldIndexCount     = "index_job_correction_corrected_old_index_count"
	correctedOldIndexCountDesc = "The number of corrected old indexes while index correction job"

	correctedReplicationCount     = "index_job_correction_corrected_replication_count"
	correctedReplicationCountDesc = "The number of operation happened to correct replication number while index correction job"
)

type correctionMetrics struct {
	correction service.Corrector
}

func New(c service.Corrector) metrics.Metric {
	return &correctionMetrics{
		correction: c,
	}
}

func (*correctionMetrics) View() ([]*metrics.View, error) {
	checkedIndexCount, err := view.New(
		view.MatchInstrumentName(checkedIndexCount),
		view.WithSetDescription(checkedIndexCountDesc),
		view.WithSetAggregation(aggregation.LastValue{}),
	)
	if err != nil {
		return nil, err
	}

	oldIndexCount, err := view.New(
		view.MatchInstrumentName(correctedOldIndexCount),
		view.WithSetDescription(correctedOldIndexCountDesc),
		view.WithSetAggregation(aggregation.LastValue{}),
	)
	if err != nil {
		return nil, err
	}

	replicationCount, err := view.New(
		view.MatchInstrumentName(correctedReplicationCount),
		view.WithSetDescription(correctedReplicationCountDesc),
		view.WithSetAggregation(aggregation.LastValue{}),
	)
	if err != nil {
		return nil, err
	}

	return []*metrics.View{
		&checkedIndexCount,
		&oldIndexCount,
		&replicationCount,
	}, nil
}

func (c *correctionMetrics) Register(m metrics.Meter) error {
	checkedIndexCount, err := m.AsyncInt64().Gauge(
		checkedIndexCount,
		metrics.WithDescription(checkedIndexCountDesc),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	oldIndexCount, err := m.AsyncInt64().Gauge(
		correctedOldIndexCount,
		metrics.WithDescription(correctedOldIndexCountDesc),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	replicationCount, err := m.AsyncInt64().Gauge(
		correctedReplicationCount,
		metrics.WithDescription(correctedReplicationCountDesc),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	return m.RegisterCallback(
		[]metrics.AsynchronousInstrument{
			checkedIndexCount,
			oldIndexCount,
			replicationCount,
		},
		func(ctx context.Context) {
			checkedIndexCount.Observe(ctx, int64(c.correction.NumberOfCheckedIndex()))
			oldIndexCount.Observe(ctx, int64(c.correction.NumberOfCorrectedOldIndex()))
			replicationCount.Observe(ctx, int64(c.correction.NumberOfCorrectedReplication()))
		},
	)
}
