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

// Package prometheus provides a prometheus exporter.
package prometheus

import (
	"context"
	"net/http"
	"reflect"
	"testing"

	"contrib.go.opencensus.io/exporter/prometheus"
)

func TestNew(t *testing.T) {
	type args struct {
		opts []PrometheusOption
	}
	tests := []struct {
		name    string
		args    args
		want    Prometheus
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.opts...)
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
		want *prometheus.Exporter
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

func Test_exporter_NewHTTPHandler(t *testing.T) {
	tests := []struct {
		name string
		e    *exporter
		want http.Handler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.NewHTTPHandler(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("exporter.NewHTTPHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExporter(t *testing.T) {
	tests := []struct {
		name    string
		want    Prometheus
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Exporter()
			if (err != nil) != tt.wantErr {
				t.Errorf("Exporter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Exporter() = %v, want %v", got, tt.want)
			}
		})
	}
}
