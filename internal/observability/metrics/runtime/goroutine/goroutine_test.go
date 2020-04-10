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

// Package goroutine provides functions for goroutine runtime stats
package goroutine

import (
	"context"
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/observability/metrics"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want metrics.Metric
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_goroutines_Measurement(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		g       *goroutines
		args    args
		want    []metrics.Measurement
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.g.Measurement(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("goroutines.Measurement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("goroutines.Measurement() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_goroutines_MeasurementWithTags(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		g       *goroutines
		args    args
		want    []metrics.MeasurementWithTags
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.g.MeasurementWithTags(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("goroutines.MeasurementWithTags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("goroutines.MeasurementWithTags() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_goroutines_View(t *testing.T) {
	tests := []struct {
		name string
		g    *goroutines
		want []*metrics.View
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.g.View(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("goroutines.View() = %v, want %v", got, tt.want)
			}
		})
	}
}
