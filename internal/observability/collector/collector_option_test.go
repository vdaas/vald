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

// Package collector provides metrics collector
package collector

import (
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/observability/metrics"
)

func TestWithErrGroup(t *testing.T) {
	type args struct {
		eg errgroup.Group
	}
	tests := []struct {
		name string
		args args
		want CollectorOption
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

func TestWithDuration(t *testing.T) {
	type args struct {
		dur string
	}
	tests := []struct {
		name string
		args args
		want CollectorOption
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithDuration(tt.args.dur); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithMetrics(t *testing.T) {
	type args struct {
		metrics []metrics.Metric
	}
	tests := []struct {
		name string
		args args
		want CollectorOption
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithMetrics(tt.args.metrics...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithMetrics() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithVersionInfo(t *testing.T) {
	type args struct {
		enabled bool
	}
	tests := []struct {
		name string
		args args
		want CollectorOption
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithVersionInfo(tt.args.enabled); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithVersionInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithMemoryMetrics(t *testing.T) {
	type args struct {
		enabled bool
	}
	tests := []struct {
		name string
		args args
		want CollectorOption
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithMemoryMetrics(tt.args.enabled); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithMemoryMetrics() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithGoroutineMetrics(t *testing.T) {
	type args struct {
		enabled bool
	}
	tests := []struct {
		name string
		args args
		want CollectorOption
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithGoroutineMetrics(tt.args.enabled); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithGoroutineMetrics() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithCGOMetrics(t *testing.T) {
	type args struct {
		enabled bool
	}
	tests := []struct {
		name string
		args args
		want CollectorOption
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithCGOMetrics(tt.args.enabled); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithCGOMetrics() = %v, want %v", got, tt.want)
			}
		})
	}
}
