//go:build e2e

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

package config

import (
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/tests/v2/e2e/metrics"
)

func waitExecution(name string) *Execution {
	return &Execution{Name: name, Type: OpWait}
}

// TestStrategy_Bind_NilReceiver guards against a nil dereference when building
// the error message for a nil Strategy.
func TestStrategy_Bind_NilReceiver(t *testing.T) {
	var s *Strategy
	got, err := s.Bind(nil)
	if got != nil {
		t.Errorf("expected nil strategy, got %v", got)
	}
	if !errors.Is(err, errors.ErrInvalidConfig) {
		t.Errorf("expected ErrInvalidConfig, got %v", err)
	}
}

// TestOperation_Bind_NilReceiver guards against a nil dereference when building
// the error message for a nil Operation.
func TestOperation_Bind_NilReceiver(t *testing.T) {
	var o *Operation
	got, err := o.Bind("strategy", nil)
	if got != nil {
		t.Errorf("expected nil operation, got %v", got)
	}
	if !errors.Is(err, errors.ErrInvalidConfig) {
		t.Errorf("expected ErrInvalidConfig, got %v", err)
	}
}

// TestOperation_Bind_SkipsNilExecutions ensures nil entries are dropped and the
// remaining executions are not duplicated by the compaction.
func TestOperation_Bind_SkipsNilExecutions(t *testing.T) {
	b := waitExecution("b")
	o := &Operation{
		Name:       "op",
		Executions: []*Execution{nil, b, nil},
	}
	got, err := o.Bind("strategy", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Executions) != 1 {
		t.Fatalf("expected 1 execution, got %d: %v", len(got.Executions), got.Executions)
	}
	if got.Executions[0] != b {
		t.Errorf("expected execution %q, got %v", b.Name, got.Executions[0])
	}
	if b.Strategy != "strategy" || b.Operation != "op" {
		t.Errorf("expected strategy/operation to be propagated, got %q/%q", b.Strategy, b.Operation)
	}
}

// TestStrategy_Bind_SkipsNilOperations ensures [nil, B] does not become [B, B].
func TestStrategy_Bind_SkipsNilOperations(t *testing.T) {
	op := &Operation{
		Name:       "op",
		Executions: []*Execution{waitExecution("e")},
	}
	s := &Strategy{
		Name:       "strategy",
		Operations: []*Operation{nil, op},
	}
	got, err := s.Bind(nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Operations) != 1 {
		t.Fatalf("expected 1 operation, got %d: %v", len(got.Operations), got.Operations)
	}
	if got.Operations[0] != op {
		t.Errorf("expected operation %q, got %v", op.Name, got.Operations[0])
	}
}

// TestPortForward_Bind_ExpandsPorts ensures the expanded port values are stored
// back on the PortForward instead of being discarded.
func TestPortForward_Bind_ExpandsPorts(t *testing.T) {
	t.Setenv("E2E_TEST_TARGET_PORT", "9090")
	t.Setenv("E2E_TEST_LOCAL_PORT", "9091")
	pf := &PortForward{
		Enabled:     true,
		Namespace:   "default",
		ServiceName: "svc",
		TargetPort:  "${E2E_TEST_TARGET_PORT}",
		LocalPort:   "${E2E_TEST_LOCAL_PORT}",
	}
	got, err := pf.Bind()
	if err != nil {
		t.Fatal(err)
	}
	if got.TargetPort != "9090" {
		t.Errorf("expected expanded target port 9090, got %q", got.TargetPort)
	}
	if got.LocalPort != "9091" {
		t.Errorf("expected expanded local port 9091, got %q", got.LocalPort)
	}
}

// TestOperator_Bind_Normalizes guards that comparison operators are returned
// in their canonical lowercase form, which is what assertion callers switch on.
func TestOperator_Bind_Normalizes(t *testing.T) {
	t.Parallel()
	tests := []struct {
		in   Operator
		want Operator
	}{
		{in: "eq", want: Eq},
		{in: "GE", want: Ge},
		{in: "Lt", want: Lt},
		{in: "NE", want: Ne},
		{in: "Contains", want: Contains},
		{in: "", want: Eq},
	}
	for _, tt := range tests {
		t.Run(string(tt.in), func(t *testing.T) {
			t.Parallel()
			got, err := tt.in.Bind()
			if err != nil {
				t.Fatalf("Bind(%q) error: %v", tt.in, err)
			}
			if got != tt.want {
				t.Errorf("Bind(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

// TestMetrics_Opts_TDigestPartialConfig guards that unset TDigest fields fall
// back to the metrics defaults instead of being rejected as zero values.
func TestMetrics_Opts_TDigestPartialConfig(t *testing.T) {
	t.Parallel()
	tests := []struct {
		m    *Metrics
		name string
	}{
		{name: "bare tdigest block", m: &Metrics{Enabled: true, TDigest: &TDigest{}}},
		{name: "latency quantiles only", m: &Metrics{Enabled: true, LatencyTDigest: &TDigest{Quantiles: []float64{0.5, 0.99}}}},
		{name: "queue wait shards only", m: &Metrics{Enabled: true, QueueWaitTDigest: &TDigest{NumShards: 2}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			m, err := tt.m.Bind()
			if err != nil {
				t.Fatalf("Bind() error: %v", err)
			}
			if m == nil {
				t.Fatal("Bind() returned nil metrics config")
			}
			if _, err := metrics.NewCollector(m.Opts()...); err != nil {
				t.Errorf("NewCollector(Opts()) error: %v", err)
			}
		})
	}
}
