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

// Package rest provides rest api logic
package rest

import (
	"net/http"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		opts []Option
	}
	tests := []struct {
		name string
		args args
		want Handler
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

func Test_handler_Index(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name    string
		h       *handler
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.h.Index(tt.args.w, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("handler.Index() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("handler.Index() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_handler_Pods(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name     string
		h        *handler
		args     args
		wantCode int
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCode, err := tt.h.Pods(tt.args.w, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("handler.Pods() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotCode != tt.wantCode {
				t.Errorf("handler.Pods() = %v, want %v", gotCode, tt.wantCode)
			}
		})
	}
}

func Test_handler_Nodes(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name     string
		h        *handler
		args     args
		wantCode int
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCode, err := tt.h.Nodes(tt.args.w, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("handler.Nodes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotCode != tt.wantCode {
				t.Errorf("handler.Nodes() = %v, want %v", gotCode, tt.wantCode)
			}
		})
	}
}
