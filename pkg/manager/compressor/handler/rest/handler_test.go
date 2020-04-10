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

func Test_handler_GetVector(t *testing.T) {
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
			got, err := tt.h.GetVector(tt.args.w, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("handler.GetVector() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("handler.GetVector() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_handler_Locations(t *testing.T) {
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
			got, err := tt.h.Locations(tt.args.w, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("handler.Locations() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("handler.Locations() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_handler_Register(t *testing.T) {
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
			got, err := tt.h.Register(tt.args.w, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("handler.Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("handler.Register() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_handler_RegisterMulti(t *testing.T) {
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
			got, err := tt.h.RegisterMulti(tt.args.w, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("handler.RegisterMulti() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("handler.RegisterMulti() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_handler_Remove(t *testing.T) {
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
			got, err := tt.h.Remove(tt.args.w, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("handler.Remove() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("handler.Remove() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_handler_RemoveMulti(t *testing.T) {
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
			got, err := tt.h.RemoveMulti(tt.args.w, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("handler.RemoveMulti() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("handler.RemoveMulti() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_handler_RegisterIPs(t *testing.T) {
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
			got, err := tt.h.RegisterIPs(tt.args.w, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("handler.RegisterIPs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("handler.RegisterIPs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_handler_RemoveIPs(t *testing.T) {
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
			got, err := tt.h.RemoveIPs(tt.args.w, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("handler.RemoveIPs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("handler.RemoveIPs() = %v, want %v", got, tt.want)
			}
		})
	}
}
