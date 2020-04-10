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

// Package jaeger provides a jaeger exporter.
package jaeger

import (
	"context"
	"reflect"
	"testing"

	"contrib.go.opencensus.io/exporter/jaeger"
)

func TestNew(t *testing.T) {
	type args struct {
		opts []JaegerOption
	}
	tests := []struct {
		name    string
		args    args
		wantJ   Jaeger
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotJ, err := New(tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotJ, tt.wantJ) {
				t.Errorf("New() = %v, want %v", gotJ, tt.wantJ)
			}
		})
	}
}

func Test_exporter_Start(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		e       *exporter
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.Start(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("exporter.Start() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_exporter_Stop(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		e    *exporter
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.e.Stop(tt.args.ctx)
		})
	}
}

func Test_exporter_Exporter(t *testing.T) {
	tests := []struct {
		name string
		e    *exporter
		want *jaeger.Exporter
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.Exporter(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("exporter.Exporter() = %v, want %v", got, tt.want)
			}
		})
	}
}
