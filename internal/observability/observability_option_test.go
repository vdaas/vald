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

// Package observability provides observability functions
package observability

import (
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/observability/collector"
	"github.com/vdaas/vald/internal/observability/exporter/jaeger"
	"github.com/vdaas/vald/internal/observability/exporter/prometheus"
	"github.com/vdaas/vald/internal/observability/trace"
)

func TestWithErrGroup(t *testing.T) {
	type args struct {
		eg errgroup.Group
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithErrGroup(tt.args.eg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithErrGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithCollector(t *testing.T) {
	type args struct {
		c collector.Collector
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithCollector(tt.args.c); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithCollector() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithTracer(t *testing.T) {
	type args struct {
		t trace.Tracer
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithTracer(tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithTracer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithPrometheus(t *testing.T) {
	type args struct {
		p prometheus.Prometheus
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithPrometheus(tt.args.p); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithPrometheus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithJaeger(t *testing.T) {
	type args struct {
		j jaeger.Jaeger
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithJaeger(tt.args.j); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithJaeger() = %v, want %v", got, tt.want)
			}
		})
	}
}
