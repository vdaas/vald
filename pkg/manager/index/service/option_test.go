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

// Package service
package service

import (
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/client/discoverer"
	"github.com/vdaas/vald/internal/errgroup"
)

func TestWithIndexingConcurrency(t *testing.T) {
	type args struct {
		c int
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
			if got := WithIndexingConcurrency(tt.args.c); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithIndexingConcurrency() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithIndexingDuration(t *testing.T) {
	type args struct {
		dur string
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
			if got := WithIndexingDuration(tt.args.dur); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithIndexingDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithIndexingDurationLimit(t *testing.T) {
	type args struct {
		dur string
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
			if got := WithIndexingDurationLimit(tt.args.dur); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithIndexingDurationLimit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithMinUncommitted(t *testing.T) {
	type args struct {
		n uint32
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
			if got := WithMinUncommitted(tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithMinUncommitted() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithCreationPoolSize(t *testing.T) {
	type args struct {
		size uint32
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
			if got := WithCreationPoolSize(tt.args.size); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithCreationPoolSize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithDiscoverer(t *testing.T) {
	type args struct {
		c discoverer.Client
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
			if got := WithDiscoverer(tt.args.c); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithDiscoverer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithErrGroup(t *testing.T) {
	type args struct {
		eg errgroup.Group
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
			if got := WithErrGroup(tt.args.eg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithErrGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}
