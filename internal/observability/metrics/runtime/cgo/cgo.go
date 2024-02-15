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
package cgo

import (
	"context"
	"runtime"

	"github.com/vdaas/vald/internal/observability/metrics"
	api "go.opentelemetry.io/otel/metric"
	view "go.opentelemetry.io/otel/sdk/metric"
)

const (
	metricsName        = "cgo_call_count"
	metricsDescription = "Number of cgo call"
)

type cgo struct{}

func New() metrics.Metric {
	return &cgo{}
}

func (*cgo) View() ([]metrics.View, error) {
	return []metrics.View{
		view.NewView(
			view.Instrument{
				Name:        metricsName,
				Description: metricsDescription,
			},
			view.Stream{
<<<<<<< HEAD
				Aggregation: view.AggregationLastValue{},
=======
				Aggregation: meric.AggregationLastValue{},
>>>>>>> feature/agent/qbg
			},
		),
	}, nil
}

func (*cgo) Register(m metrics.Meter) error {
	count, err := m.Int64ObservableGauge(
		metricsName,
		metrics.WithDescription(metricsDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}
	_, err = m.RegisterCallback(
		func(_ context.Context, o api.Observer) error {
			o.ObserveInt64(count, runtime.NumCgoCall())
			return nil
		}, count,
	)
	return err
}
