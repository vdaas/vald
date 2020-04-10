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

func Test_handler_GetMeta(t *testing.T) {
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
			got, err := tt.h.GetMeta(tt.args.w, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("handler.GetMeta() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("handler.GetMeta() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_handler_GetMetas(t *testing.T) {
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
			got, err := tt.h.GetMetas(tt.args.w, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("handler.GetMetas() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("handler.GetMetas() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_handler_GetMetaInverse(t *testing.T) {
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
			got, err := tt.h.GetMetaInverse(tt.args.w, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("handler.GetMetaInverse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("handler.GetMetaInverse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_handler_GetMetasInverse(t *testing.T) {
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
			got, err := tt.h.GetMetasInverse(tt.args.w, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("handler.GetMetasInverse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("handler.GetMetasInverse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_handler_SetMeta(t *testing.T) {
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
			got, err := tt.h.SetMeta(tt.args.w, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("handler.SetMeta() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("handler.SetMeta() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_handler_SetMetas(t *testing.T) {
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
			got, err := tt.h.SetMetas(tt.args.w, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("handler.SetMetas() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("handler.SetMetas() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_handler_DeleteMeta(t *testing.T) {
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
			got, err := tt.h.DeleteMeta(tt.args.w, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("handler.DeleteMeta() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("handler.DeleteMeta() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_handler_DeleteMetas(t *testing.T) {
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
			got, err := tt.h.DeleteMetas(tt.args.w, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("handler.DeleteMetas() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("handler.DeleteMetas() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_handler_DeleteMetaInverse(t *testing.T) {
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
			got, err := tt.h.DeleteMetaInverse(tt.args.w, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("handler.DeleteMetaInverse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("handler.DeleteMetaInverse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_handler_DeleteMetasInverse(t *testing.T) {
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
			got, err := tt.h.DeleteMetasInverse(tt.args.w, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("handler.DeleteMetasInverse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("handler.DeleteMetasInverse() = %v, want %v", got, tt.want)
			}
		})
	}
}
