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

// Package metrics provides metrics.
package metrics

import (
	"context"
	"testing"
)

func TestRegisterView(t *testing.T) {
	type args struct {
		views []*View
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RegisterView(tt.args.views...); (err != nil) != tt.wantErr {
				t.Errorf("RegisterView() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRecord(t *testing.T) {
	type args struct {
		ctx context.Context
		ms  []Measurement
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Record(tt.args.ctx, tt.args.ms...)
		})
	}
}

func TestRecordWithTags(t *testing.T) {
	type args struct {
		ctx  context.Context
		mwts []MeasurementWithTags
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RecordWithTags(tt.args.ctx, tt.args.mwts...); (err != nil) != tt.wantErr {
				t.Errorf("RecordWithTags() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMeasurementsCount(t *testing.T) {
	type args struct {
		m Metric
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MeasurementsCount(tt.args.m); got != tt.want {
				t.Errorf("MeasurementsCount() = %v, want %v", got, tt.want)
			}
		})
	}
}
