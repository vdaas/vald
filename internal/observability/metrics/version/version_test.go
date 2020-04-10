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

// Package version provides version info metrics functions
package version

import (
	"context"
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/observability/metrics"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		want    metrics.Metric
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New()
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_labelKVs(t *testing.T) {
	tests := []struct {
		name    string
		want    map[metrics.Key]string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := labelKVs()
			if (err != nil) != tt.wantErr {
				t.Errorf("labelKVs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("labelKVs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_version_Measurement(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		v       *version
		args    args
		want    []metrics.Measurement
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.v.Measurement(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("version.Measurement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("version.Measurement() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_version_MeasurementWithTags(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		v       *version
		args    args
		want    []metrics.MeasurementWithTags
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.v.MeasurementWithTags(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("version.MeasurementWithTags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("version.MeasurementWithTags() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_version_View(t *testing.T) {
	tests := []struct {
		name string
		v    *version
		want []*metrics.View
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.View(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("version.View() = %v, want %v", got, tt.want)
			}
		})
	}
}
