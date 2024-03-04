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
	api "go.opentelemetry.io/otel/metric"
	view "go.opentelemetry.io/otel/sdk/metric"
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

func (*correctionMetrics) View() ([]metrics.View, error) {
	return []metrics.View{
		view.NewView(
			view.Instrument{
				Name:        checkedIndexCount,
				Description: checkedIndexCountDesc,
			},
			view.Stream{
				Aggregation: view.AggregationLastValue{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        correctedOldIndexCount,
				Description: correctedOldIndexCountDesc,
			},
			view.Stream{
				Aggregation: view.AggregationLastValue{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        correctedReplicationCount,
				Description: correctedReplicationCountDesc,
			},
			view.Stream{
				Aggregation: view.AggregationLastValue{},
			},
		),
	}, nil
}

func (c *correctionMetrics) Register(m metrics.Meter) error {
	checkedIndexCount, err := m.Int64ObservableGauge(
		checkedIndexCount,
		metrics.WithDescription(checkedIndexCountDesc),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	oldIndexCount, err := m.Int64ObservableGauge(
		correctedOldIndexCount,
		metrics.WithDescription(correctedOldIndexCountDesc),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	replicationCount, err := m.Int64ObservableGauge(
		correctedReplicationCount,
		metrics.WithDescription(correctedReplicationCountDesc),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	_, err = m.RegisterCallback(
		func(_ context.Context, o api.Observer) error {
			o.ObserveInt64(checkedIndexCount, int64(c.correction.NumberOfCheckedIndex()))
			o.ObserveInt64(oldIndexCount, int64(c.correction.NumberOfCorrectedOldIndex()))
			o.ObserveInt64(replicationCount, int64(c.correction.NumberOfCorrectedReplication()))
			return nil
		},
		checkedIndexCount,
		oldIndexCount,
		replicationCount,
	)
	return err
}
