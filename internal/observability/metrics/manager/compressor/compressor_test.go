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

// Package compressor provides functions for compressor stats
package compressor

import (
	"context"
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/observability/metrics"
	"github.com/vdaas/vald/pkg/manager/compressor/service"
)

func TestNew(t *testing.T) {
	type args struct {
		c service.Compressor
		r service.Registerer
	}
	tests := []struct {
		name string
		args args
		want metrics.Metric
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.c, tt.args.r); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_compressorMetrics_Measurement(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		c       *compressorMetrics
		args    args
		want    []metrics.Measurement
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.Measurement(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("compressorMetrics.Measurement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("compressorMetrics.Measurement() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_compressorMetrics_MeasurementWithTags(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		c       *compressorMetrics
		args    args
		want    []metrics.MeasurementWithTags
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.MeasurementWithTags(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("compressorMetrics.MeasurementWithTags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("compressorMetrics.MeasurementWithTags() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_compressorMetrics_View(t *testing.T) {
	tests := []struct {
		name string
		c    *compressorMetrics
		want []*metrics.View
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.View(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("compressorMetrics.View() = %v, want %v", got, tt.want)
			}
		})
	}
}
