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

// Package gache provides implementation of cache using gache
package gache

import (
	"context"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		opts []Option
	}
	tests := []struct {
		name  string
		args  args
		wantC *cache
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotC := New(tt.args.opts...); !reflect.DeepEqual(gotC, tt.wantC) {
				t.Errorf("New() = %v, want %v", gotC, tt.wantC)
			}
		})
	}
}

func Test_cache_Start(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		c    *cache
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.Start(tt.args.ctx)
		})
	}
}

func Test_cache_Get(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name  string
		c     *cache
		args  args
		want  interface{}
		want1 bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.c.Get(tt.args.key)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("cache.Get() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("cache.Get() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_cache_Set(t *testing.T) {
	type args struct {
		key string
		val interface{}
	}
	tests := []struct {
		name string
		c    *cache
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.Set(tt.args.key, tt.args.val)
		})
	}
}

func Test_cache_Delete(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		c    *cache
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.Delete(tt.args.key)
		})
	}
}

func Test_cache_GetAndDelete(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name  string
		c     *cache
		args  args
		want  interface{}
		want1 bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.c.GetAndDelete(tt.args.key)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("cache.GetAndDelete() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("cache.GetAndDelete() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
