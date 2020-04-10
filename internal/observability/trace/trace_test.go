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

// Package trace provides trace functions.
package trace

import (
	"context"
	"reflect"
	"testing"

	"go.opencensus.io/trace"
)

func TestStartSpan(t *testing.T) {
	type args struct {
		ctx  context.Context
		name string
		opts []trace.StartOption
	}
	tests := []struct {
		name  string
		args  args
		want  context.Context
		want1 *Span
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := StartSpan(tt.args.ctx, tt.args.name, tt.args.opts...)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StartSpan() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("StartSpan() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		opts []TraceOption
	}
	tests := []struct {
		name string
		args args
		want Tracer
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_tracer_Start(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		t    *tracer
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.t.Start(tt.args.ctx)
		})
	}
}
