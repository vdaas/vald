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

	"github.com/vdaas/vald/internal/errgroup"
)

func TestWithLimitation(t *testing.T) {
	type args struct {
		limit int
	}
	tests := []struct {
		name string
		args args
		want CompressorOption
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithLimitation(tt.args.limit); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithLimitation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithBuffer(t *testing.T) {
	type args struct {
		b int
	}
	tests := []struct {
		name string
		args args
		want CompressorOption
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithBuffer(tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithBuffer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithCompressAlgorithm(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want CompressorOption
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithCompressAlgorithm(tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithCompressAlgorithm() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithCompressionLevel(t *testing.T) {
	type args struct {
		level int
	}
	tests := []struct {
		name string
		args args
		want CompressorOption
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithCompressionLevel(tt.args.level); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithCompressionLevel() = %v, want %v", got, tt.want)
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
		want CompressorOption
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
