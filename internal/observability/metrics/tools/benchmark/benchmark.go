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
package benchmark

import (
	"context"

	v1 "github.com/vdaas/vald/internal/k8s/vald/benchmark/api/v1"
	"github.com/vdaas/vald/internal/observability/metrics"
	"github.com/vdaas/vald/pkg/operator/benchmark/service"
	api "go.opentelemetry.io/otel/metric"
	view "go.opentelemetry.io/otel/sdk/metric"
)

const (
	AppliedScenarioCount            = "benchmark_operator_applied_scenario"
	AppliedScenarioCountDescription = "Benchmark Operator applied scenario count"

	RunningScenarioCount            = "benchmark_operator_running_scenario"
	RunningScenarioCountDescription = "Benchmark Operator running scenario count"

	CompleteScenarioCount            = "benchmark_operator_complete_scenario"
	CompleteScenarioCountDescription = "Benchmark Operator complete scenario count"

	AppliedBenchmarkJobCount            = "benchmark_operator_applied_benchmark_job"
	AppliedBenchmarkJobCountDescription = "Benchmark Operator applied benchmark job count"

	RunningBenchmarkJobCount            = "benchmark_operator_running_benchmark_job"
	RunningBenchmarkJobCountDescription = "Benchmark Operator running benchmark job count"

	CompleteBenchmarkJobCount            = "benchmark_operator_complete_benchmark_job"
	CompleteBenchmarkJobCountDescription = "Benchmark Operator complete benchmark job count"
)

type operatorMetrics struct {
	op service.Operator
}

func New(om service.Operator) metrics.Metric {
	return &operatorMetrics{
		op: om,
	}
}

// TODO: implement here.
func (om *operatorMetrics) View() ([]metrics.View, error) {
	return []metrics.View{
		view.NewView(
			view.Instrument{
				Name:        AppliedScenarioCount,
				Description: AppliedScenarioCountDescription,
			},
			view.Stream{
				Aggregation: view.AggregationLastValue{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        RunningScenarioCount,
				Description: RunningScenarioCountDescription,
			},
			view.Stream{
				Aggregation: view.AggregationLastValue{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        CompleteScenarioCount,
				Description: CompleteScenarioCountDescription,
			},
			view.Stream{
				Aggregation: view.AggregationLastValue{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        AppliedBenchmarkJobCount,
				Description: AppliedBenchmarkJobCountDescription,
			},
			view.Stream{
				Aggregation: view.AggregationLastValue{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        RunningBenchmarkJobCount,
				Description: RunningBenchmarkJobCountDescription,
			},
			view.Stream{
				Aggregation: view.AggregationLastValue{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        CompleteBenchmarkJobCount,
				Description: CompleteBenchmarkJobCountDescription,
			},
			view.Stream{
				Aggregation: view.AggregationLastValue{},
			},
		),
	}, nil
}

// TODO: implement here.
func (om *operatorMetrics) Register(m metrics.Meter) error {
	appliedScenarioCount, err := m.Int64ObservableGauge(
		AppliedScenarioCount,
		metrics.WithDescription(AppliedScenarioCountDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}
	runningScenarioCount, err := m.Int64ObservableGauge(
		RunningScenarioCount,
		metrics.WithDescription(RunningScenarioCountDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}
	completeScenarioCount, err := m.Int64ObservableGauge(
		CompleteScenarioCount,
		metrics.WithDescription(CompleteScenarioCountDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	appliedBenchJobCount, err := m.Int64ObservableGauge(
		AppliedBenchmarkJobCount,
		metrics.WithDescription(AppliedBenchmarkJobCountDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}
	runningBenchJobCount, err := m.Int64ObservableGauge(
		RunningBenchmarkJobCount,
		metrics.WithDescription(RunningBenchmarkJobCountDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}
	completeBenchJobCount, err := m.Int64ObservableGauge(
		CompleteBenchmarkJobCount,
		metrics.WithDescription(CompleteBenchmarkJobCountDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	_, err = m.RegisterCallback(
		func(_ context.Context, o api.Observer) error {
			// scenario status
			var sApplied, sRunning, sComplete int64
			for k, v := range om.op.GetScenarioStatus() {
				sApplied += v
				if k == v1.BenchmarkScenarioCompleted {
					sComplete += v
				} else {
					sRunning += v
				}
			}
			o.ObserveInt64(appliedScenarioCount, sApplied)
			o.ObserveInt64(runningScenarioCount, sRunning)
			o.ObserveInt64(completeScenarioCount, sComplete)

			// benchmark job status
			var bApplied, bRunning, bComplete int64
			for k, v := range om.op.GetBenchmarkJobStatus() {
				bApplied += v
				if k == v1.BenchmarkJobCompleted {
					bComplete += v
				} else {
					bRunning += v
				}
			}
			o.ObserveInt64(appliedBenchJobCount, bApplied)
			o.ObserveInt64(runningBenchJobCount, bRunning)
			o.ObserveInt64(completeBenchJobCount, bComplete)
			return nil
		},
		appliedScenarioCount,
		runningScenarioCount,
		completeScenarioCount,
		appliedBenchJobCount,
		runningBenchJobCount,
		completeBenchJobCount,
	)
	return err
}
