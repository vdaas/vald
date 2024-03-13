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
package benchmark

import (
	"context"

	v1 "github.com/vdaas/vald/internal/k8s/vald/benchmark/api/v1"
	"github.com/vdaas/vald/internal/observability/metrics"
	"github.com/vdaas/vald/pkg/tools/benchmark/operator/service"
	api "go.opentelemetry.io/otel/metric"
	view "go.opentelemetry.io/otel/sdk/metric"
)

const (
	appliedScenarioCount            = "benchmark_operator_applied_scenario"
	appliedScenarioCountDescription = "Benchmark Operator applied scenario count"

	runningScenarioCount            = "benchmark_operator_running_scenario"
	runningScenarioCountDescription = "Benchmark Operator running scenario count"

	completeScenarioCount            = "benchmark_operator_complete_scenario"
	completeScenarioCountDescription = "Benchmark Operator complete scenario count"

	appliedBenchmarkJobCount            = "benchmark_operator_applied_benchmark_job"
	appliedBenchmarkJobCountDescription = "Benchmark Operator applied benchmark job count"

	runningBenchmarkJobCount            = "benchmark_operator_running_benchmark_job"
	runningBenchmarkJobCountDescription = "Benchmark Operator running benchmark job count"

	completeBenchmarkJobCount            = "benchmark_operator_complete_benchmark_job"
	completeBenchmarkJobCountDescription = "Benchmark Operator complete benchmark job count"

	appliedJobCount            = "benchmark_operator_applied_job"
	appliedJobCountDescription = "Benchmark Operator applied job count"

	runningJobCount            = "benchmark_operator_running_job"
	runningJobCountDescription = "Benchmark Operator running job count"

	completeJobCount            = "benchmark_operator_complete_job"
	completeJobCountDescription = "Benchmark Operator complete job count"
)

const (
	applied  = "applied"
	running  = "running"
	complete = "complete"
)

type operatorMetrics struct {
	op service.Operator
}

func New(om service.Operator) metrics.Metric {
	return &operatorMetrics{
		op: om,
	}
}

// TODO: implement here
func (om *operatorMetrics) View() ([]metrics.View, error) {
	return []metrics.View{
		view.NewView(
			view.Instrument{
				Name:        appliedScenarioCount,
				Description: appliedScenarioCountDescription,
			},
			view.Stream{
				Aggregation: view.AggregationLastValue{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        runningScenarioCount,
				Description: runningScenarioCountDescription,
			},
			view.Stream{
				Aggregation: view.AggregationLastValue{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        completeScenarioCount,
				Description: completeScenarioCountDescription,
			},
			view.Stream{
				Aggregation: view.AggregationLastValue{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        appliedBenchmarkJobCount,
				Description: appliedBenchmarkJobCountDescription,
			},
			view.Stream{
				Aggregation: view.AggregationLastValue{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        runningBenchmarkJobCount,
				Description: runningBenchmarkJobCountDescription,
			},
			view.Stream{
				Aggregation: view.AggregationLastValue{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        completeBenchmarkJobCount,
				Description: completeBenchmarkJobCountDescription,
			},
			view.Stream{
				Aggregation: view.AggregationLastValue{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        appliedJobCount,
				Description: appliedJobCountDescription,
			},
			view.Stream{
				Aggregation: view.AggregationLastValue{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        runningJobCount,
				Description: runningJobCountDescription,
			},
			view.Stream{
				Aggregation: view.AggregationLastValue{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        completeJobCount,
				Description: completeJobCountDescription,
			},
			view.Stream{
				Aggregation: view.AggregationLastValue{},
			},
		),
	}, nil
}

// TODO: implement here
func (om *operatorMetrics) Register(m metrics.Meter) error {
	appliedScenarioCount, err := m.Int64ObservableCounter(
		appliedScenarioCount,
		metrics.WithDescription(appliedScenarioCountDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}
	runningScenarioCount, err := m.Int64ObservableCounter(
		runningScenarioCount,
		metrics.WithDescription(runningScenarioCountDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}
	completeScenarioCount, err := m.Int64ObservableCounter(
		completeScenarioCount,
		metrics.WithDescription(completeScenarioCountDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	appliedBenchJobCount, err := m.Int64ObservableCounter(
		appliedBenchmarkJobCount,
		metrics.WithDescription(appliedScenarioCountDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}
	runningBenchJobCount, err := m.Int64ObservableCounter(
		runningBenchmarkJobCount,
		metrics.WithDescription(runningScenarioCountDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}
	completeBenchJobCount, err := m.Int64ObservableCounter(
		completeBenchmarkJobCount,
		metrics.WithDescription(completeScenarioCountDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}
	appliedJobCount, err := m.Int64ObservableCounter(
		appliedJobCount,
		metrics.WithDescription(appliedJobCountDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}
	runningJobCount, err := m.Int64ObservableCounter(
		runningJobCount,
		metrics.WithDescription(runningJobCountDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}
	completeJobCount, err := m.Int64ObservableCounter(
		completeBenchmarkJobCount,
		metrics.WithDescription(completeScenarioCountDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}
	_, err = m.RegisterCallback(
		func(_ context.Context, o api.Observer) error {
			// scenario status
			sst := map[string]int64{
				applied:  0,
				running:  0,
				complete: 0,
			}
			for k, v := range om.op.GetScenarioStatus() {
				sst[applied] += v
				if k == v1.BenchmarkScenarioCompleted {
					sst[complete] += v
				} else {
					sst[running] += v
				}
			}
			o.ObserveInt64(appliedScenarioCount, sst[applied])
			o.ObserveInt64(runningScenarioCount, sst[running])
			o.ObserveInt64(completeScenarioCount, sst[complete])

			// benchmark job status
			bst := map[string]int64{
				applied:  0,
				running:  0,
				complete: 0,
			}
			for k, v := range om.op.GetBenchmarkJobStatus() {
				bst[applied] += v
				if k == v1.BenchmarkJobCompleted {
					bst[complete] += v
				} else {
					bst[running] += v
				}
			}
			o.ObserveInt64(appliedBenchJobCount, bst[applied])
			o.ObserveInt64(runningBenchJobCount, bst[running])
			o.ObserveInt64(completeBenchJobCount, bst[complete])
			return nil
		},
		appliedScenarioCount,
		runningScenarioCount,
		completeScenarioCount,
		appliedBenchJobCount,
		runningBenchJobCount,
		completeBenchJobCount,
		appliedJobCount,
		runningJobCount,
		completeJobCount,
	)
	return nil
}
