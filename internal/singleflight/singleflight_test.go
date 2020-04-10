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

// Package singleflight represents zero time caching
package singleflight

import (
	"context"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		size int
	}
	tests := []struct {
		name string
		args args
		want Group
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.size); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_group_Do(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
		fn  func() (interface{}, error)
	}
	tests := []struct {
		name       string
		g          *group
		args       args
		wantV      interface{}
		wantShared bool
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotV, gotShared, err := tt.g.Do(tt.args.ctx, tt.args.key, tt.args.fn)
			if (err != nil) != tt.wantErr {
				t.Errorf("group.Do() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotV, tt.wantV) {
				t.Errorf("group.Do() gotV = %v, want %v", gotV, tt.wantV)
			}
			if gotShared != tt.wantShared {
				t.Errorf("group.Do() gotShared = %v, want %v", gotShared, tt.wantShared)
			}
		})
	}
}
