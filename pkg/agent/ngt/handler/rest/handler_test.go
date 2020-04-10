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

func Test_handler_Search(t *testing.T) {
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
			gotCode, err := tt.h.Search(tt.args.w, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("handler.Search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotCode != tt.wantCode {
				t.Errorf("handler.Search() = %v, want %v", gotCode, tt.wantCode)
			}
		})
	}
}

func Test_handler_SearchByID(t *testing.T) {
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
			gotCode, err := tt.h.SearchByID(tt.args.w, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("handler.SearchByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotCode != tt.wantCode {
				t.Errorf("handler.SearchByID() = %v, want %v", gotCode, tt.wantCode)
			}
		})
	}
}

func Test_handler_Insert(t *testing.T) {
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
			gotCode, err := tt.h.Insert(tt.args.w, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("handler.Insert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotCode != tt.wantCode {
				t.Errorf("handler.Insert() = %v, want %v", gotCode, tt.wantCode)
			}
		})
	}
}

func Test_handler_MultiInsert(t *testing.T) {
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
			gotCode, err := tt.h.MultiInsert(tt.args.w, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("handler.MultiInsert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotCode != tt.wantCode {
				t.Errorf("handler.MultiInsert() = %v, want %v", gotCode, tt.wantCode)
			}
		})
	}
}

func Test_handler_Update(t *testing.T) {
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
			gotCode, err := tt.h.Update(tt.args.w, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("handler.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotCode != tt.wantCode {
				t.Errorf("handler.Update() = %v, want %v", gotCode, tt.wantCode)
			}
		})
	}
}

func Test_handler_MultiUpdate(t *testing.T) {
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
			gotCode, err := tt.h.MultiUpdate(tt.args.w, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("handler.MultiUpdate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotCode != tt.wantCode {
				t.Errorf("handler.MultiUpdate() = %v, want %v", gotCode, tt.wantCode)
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
			gotCode, err := tt.h.Remove(tt.args.w, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("handler.Remove() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotCode != tt.wantCode {
				t.Errorf("handler.Remove() = %v, want %v", gotCode, tt.wantCode)
			}
		})
	}
}

func Test_handler_MultiRemove(t *testing.T) {
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
			gotCode, err := tt.h.MultiRemove(tt.args.w, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("handler.MultiRemove() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotCode != tt.wantCode {
				t.Errorf("handler.MultiRemove() = %v, want %v", gotCode, tt.wantCode)
			}
		})
	}
}

func Test_handler_CreateIndex(t *testing.T) {
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
			gotCode, err := tt.h.CreateIndex(tt.args.w, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("handler.CreateIndex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotCode != tt.wantCode {
				t.Errorf("handler.CreateIndex() = %v, want %v", gotCode, tt.wantCode)
			}
		})
	}
}

func Test_handler_SaveIndex(t *testing.T) {
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
			gotCode, err := tt.h.SaveIndex(tt.args.w, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("handler.SaveIndex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotCode != tt.wantCode {
				t.Errorf("handler.SaveIndex() = %v, want %v", gotCode, tt.wantCode)
			}
		})
	}
}

func Test_handler_CreateAndSaveIndex(t *testing.T) {
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
			gotCode, err := tt.h.CreateAndSaveIndex(tt.args.w, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("handler.CreateAndSaveIndex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotCode != tt.wantCode {
				t.Errorf("handler.CreateAndSaveIndex() = %v, want %v", gotCode, tt.wantCode)
			}
		})
	}
}

func Test_handler_GetObject(t *testing.T) {
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
			gotCode, err := tt.h.GetObject(tt.args.w, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("handler.GetObject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotCode != tt.wantCode {
				t.Errorf("handler.GetObject() = %v, want %v", gotCode, tt.wantCode)
			}
		})
	}
}

func Test_handler_Exists(t *testing.T) {
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
			gotCode, err := tt.h.Exists(tt.args.w, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("handler.Exists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotCode != tt.wantCode {
				t.Errorf("handler.Exists() = %v, want %v", gotCode, tt.wantCode)
			}
		})
	}
}
