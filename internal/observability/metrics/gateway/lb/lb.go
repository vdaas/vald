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
package lb

import (
	"context"

	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/observability/attribute"
	"github.com/vdaas/vald/internal/observability/metrics"
	"github.com/vdaas/vald/internal/sync"
	api "go.opentelemetry.io/otel/metric"
	view "go.opentelemetry.io/otel/sdk/metric"
)

const (
	SearchAgentFailureMetricsName        = "gateway_lb_search_agent_failure_total"
	SearchAgentFailureMetricsDescription = "Cumulative count of agents whose search broadcast call failed"

	SearchPartialSuccessMetricsName        = "gateway_lb_search_partial_success_total"
	SearchPartialSuccessMetricsDescription = "Cumulative count of search requests that returned results although some agents failed"

	WriteAgentFailureMetricsName        = "gateway_lb_write_agent_failure_total"
	WriteAgentFailureMetricsDescription = "Cumulative count of agents whose write call failed in a write request that succeeded on other agents"

	WritePartialSuccessMetricsName        = "gateway_lb_write_partial_success_total"
	WritePartialSuccessMetricsDescription = "Cumulative count of write requests that succeeded although they were applied to fewer agents than intended"

	rpcKey        = "rpc"
	grpcCodeKey   = "grpc_code"
	recallRiskKey = "recall_risk"
)

type failureKey struct {
	rpc  string
	code codes.Code
}

type partialKey struct {
	rpc        string
	recallRisk bool
}

// SearchMetrics counts the per-agent failures of the lb gateway search
// broadcast and of the write requests that succeeded partially. The counts are cumulative because the instruments are observable
// counters, which must report the running total. The agent address is
// deliberately not an attribute to keep the cardinality bounded.
// A nil *SearchMetrics is valid and records nothing.
type SearchMetrics struct {
	failures      map[failureKey]int64
	partials      map[partialKey]int64
	writeFailures map[failureKey]int64
	writePartials map[string]int64
	mu            sync.Mutex
}

func New() *SearchMetrics {
	return &SearchMetrics{
		failures:      make(map[failureKey]int64),
		partials:      make(map[partialKey]int64),
		writeFailures: make(map[failureKey]int64),
		writePartials: make(map[string]int64),
	}
}

// RecordAgentFailure counts one agent whose search call for rpc failed with code.
func (sm *SearchMetrics) RecordAgentFailure(rpc string, code codes.Code) {
	if sm == nil {
		return
	}
	sm.mu.Lock()
	sm.failures[failureKey{rpc: rpc, code: code}]++
	sm.mu.Unlock()
}

// RecordPartialSuccess counts one rpc request that returned results although
// some agents failed. recallRisk reports whether the failed agents were at least
// the index replica count, i.e. some vectors may have been unreachable.
func (sm *SearchMetrics) RecordPartialSuccess(rpc string, recallRisk bool) {
	if sm == nil {
		return
	}
	sm.mu.Lock()
	sm.partials[partialKey{rpc: rpc, recallRisk: recallRisk}]++
	sm.mu.Unlock()
}

// RecordWriteAgentFailure counts one agent whose write call for rpc failed with
// code in a write request that succeeded on other agents.
func (sm *SearchMetrics) RecordWriteAgentFailure(rpc string, code codes.Code) {
	if sm == nil {
		return
	}
	sm.mu.Lock()
	sm.writeFailures[failureKey{rpc: rpc, code: code}]++
	sm.mu.Unlock()
}

// RecordWritePartialSuccess counts one rpc write request that succeeded although
// it was applied to fewer agents than intended.
func (sm *SearchMetrics) RecordWritePartialSuccess(rpc string) {
	if sm == nil {
		return
	}
	sm.mu.Lock()
	sm.writePartials[rpc]++
	sm.mu.Unlock()
}

func (*SearchMetrics) View() ([]metrics.View, error) {
	return []metrics.View{
		view.NewView(
			view.Instrument{
				Name:        SearchAgentFailureMetricsName,
				Description: SearchAgentFailureMetricsDescription,
			},
			view.Stream{
				Aggregation: view.AggregationSum{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        SearchPartialSuccessMetricsName,
				Description: SearchPartialSuccessMetricsDescription,
			},
			view.Stream{
				Aggregation: view.AggregationSum{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        WriteAgentFailureMetricsName,
				Description: WriteAgentFailureMetricsDescription,
			},
			view.Stream{
				Aggregation: view.AggregationSum{},
			},
		),
		view.NewView(
			view.Instrument{
				Name:        WritePartialSuccessMetricsName,
				Description: WritePartialSuccessMetricsDescription,
			},
			view.Stream{
				Aggregation: view.AggregationSum{},
			},
		),
	}, nil
}

func (sm *SearchMetrics) Register(m metrics.Meter) error {
	failureTotal, err := m.Int64ObservableCounter(
		SearchAgentFailureMetricsName,
		metrics.WithDescription(SearchAgentFailureMetricsDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}
	partialTotal, err := m.Int64ObservableCounter(
		SearchPartialSuccessMetricsName,
		metrics.WithDescription(SearchPartialSuccessMetricsDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}
	writeFailureTotal, err := m.Int64ObservableCounter(
		WriteAgentFailureMetricsName,
		metrics.WithDescription(WriteAgentFailureMetricsDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}
	writePartialTotal, err := m.Int64ObservableCounter(
		WritePartialSuccessMetricsName,
		metrics.WithDescription(WritePartialSuccessMetricsDescription),
		metrics.WithUnit(metrics.Dimensionless),
	)
	if err != nil {
		return err
	}

	_, err = m.RegisterCallback(
		func(_ context.Context, o api.Observer) error {
			sm.mu.Lock()
			defer sm.mu.Unlock()
			for k, n := range sm.failures {
				o.ObserveInt64(failureTotal, n, api.WithAttributes(
					attribute.String(rpcKey, k.rpc),
					attribute.String(grpcCodeKey, k.code.String()),
				))
			}
			for k, n := range sm.partials {
				o.ObserveInt64(partialTotal, n, api.WithAttributes(
					attribute.String(rpcKey, k.rpc),
					attribute.Bool(recallRiskKey, k.recallRisk),
				))
			}
			for k, n := range sm.writeFailures {
				o.ObserveInt64(writeFailureTotal, n, api.WithAttributes(
					attribute.String(rpcKey, k.rpc),
					attribute.String(grpcCodeKey, k.code.String()),
				))
			}
			for rpc, n := range sm.writePartials {
				o.ObserveInt64(writePartialTotal, n, api.WithAttributes(
					attribute.String(rpcKey, rpc),
				))
			}
			return nil
		},
		failureTotal,
		partialTotal,
		writeFailureTotal,
		writePartialTotal,
	)
	return err
}
