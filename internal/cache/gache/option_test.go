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
	"time"

	"github.com/kpango/gache"
)

func TestWithGache(t *testing.T) {
	type args struct {
		g gache.Gache
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithGache(tt.args.g); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithGache() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithExpiredHook(t *testing.T) {
	type args struct {
		f func(context.Context, string)
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithExpiredHook(tt.args.f); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithExpiredHook() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithExpireDuration(t *testing.T) {
	type args struct {
		dur time.Duration
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithExpireDuration(tt.args.dur); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithExpireDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithExpireCheckDuration(t *testing.T) {
	type args struct {
		dur time.Duration
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithExpireCheckDuration(tt.args.dur); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithExpireCheckDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}
