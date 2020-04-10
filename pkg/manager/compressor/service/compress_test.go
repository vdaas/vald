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

package service

import (
	"context"
	"reflect"
	"testing"
)

func TestNewCompressor(t *testing.T) {
	type args struct {
		opts []CompressorOption
	}
	tests := []struct {
		name    string
		args    args
		want    Compressor
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewCompressor(tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCompressor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCompressor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_compressor_PreStart(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		c       *compressor
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.PreStart(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("compressor.PreStart() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_compressor_Start(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		c    *compressor
		args args
		want <-chan error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.Start(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("compressor.Start() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_compressor_dispatchCompress(t *testing.T) {
	type args struct {
		ctx     context.Context
		vectors [][]float32
	}
	tests := []struct {
		name        string
		c           *compressor
		args        args
		wantResults [][]byte
		wantErr     bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResults, err := tt.c.dispatchCompress(tt.args.ctx, tt.args.vectors...)
			if (err != nil) != tt.wantErr {
				t.Errorf("compressor.dispatchCompress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResults, tt.wantResults) {
				t.Errorf("compressor.dispatchCompress() = %v, want %v", gotResults, tt.wantResults)
			}
		})
	}
}

func Test_compressor_dispatchDecompress(t *testing.T) {
	type args struct {
		ctx    context.Context
		bytess [][]byte
	}
	tests := []struct {
		name        string
		c           *compressor
		args        args
		wantResults [][]float32
		wantErr     bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResults, err := tt.c.dispatchDecompress(tt.args.ctx, tt.args.bytess...)
			if (err != nil) != tt.wantErr {
				t.Errorf("compressor.dispatchDecompress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResults, tt.wantResults) {
				t.Errorf("compressor.dispatchDecompress() = %v, want %v", gotResults, tt.wantResults)
			}
		})
	}
}

func Test_compressor_Compress(t *testing.T) {
	type args struct {
		ctx    context.Context
		vector []float32
	}
	tests := []struct {
		name    string
		c       *compressor
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.Compress(tt.args.ctx, tt.args.vector)
			if (err != nil) != tt.wantErr {
				t.Errorf("compressor.Compress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("compressor.Compress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_compressor_Decompress(t *testing.T) {
	type args struct {
		ctx   context.Context
		bytes []byte
	}
	tests := []struct {
		name    string
		c       *compressor
		args    args
		want    []float32
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.Decompress(tt.args.ctx, tt.args.bytes)
			if (err != nil) != tt.wantErr {
				t.Errorf("compressor.Decompress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("compressor.Decompress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_compressor_MultiCompress(t *testing.T) {
	type args struct {
		ctx     context.Context
		vectors [][]float32
	}
	tests := []struct {
		name    string
		c       *compressor
		args    args
		want    [][]byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.MultiCompress(tt.args.ctx, tt.args.vectors)
			if (err != nil) != tt.wantErr {
				t.Errorf("compressor.MultiCompress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("compressor.MultiCompress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_compressor_MultiDecompress(t *testing.T) {
	type args struct {
		ctx    context.Context
		bytess [][]byte
	}
	tests := []struct {
		name    string
		c       *compressor
		args    args
		want    [][]float32
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.MultiDecompress(tt.args.ctx, tt.args.bytess)
			if (err != nil) != tt.wantErr {
				t.Errorf("compressor.MultiDecompress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("compressor.MultiDecompress() = %v, want %v", got, tt.want)
			}
		})
	}
}
