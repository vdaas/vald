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
	"context"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		opts []Option
	}
	tests := []struct {
		name    string
		args    args
		wantIdx Indexer
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIdx, err := New(tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotIdx, tt.wantIdx) {
				t.Errorf("New() = %v, want %v", gotIdx, tt.wantIdx)
			}
		})
	}
}

func Test_index_Start(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		idx     *index
		args    args
		want    <-chan error
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.idx.Start(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("index.Start() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("index.Start() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_index_execute(t *testing.T) {
	type args struct {
		ctx                context.Context
		enableLowIndexSkip bool
	}
	tests := []struct {
		name    string
		idx     *index
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.idx.execute(tt.args.ctx, tt.args.enableLowIndexSkip); (err != nil) != tt.wantErr {
				t.Errorf("index.execute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_index_loadInfos(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		idx     *index
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.idx.loadInfos(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("index.loadInfos() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_index_IsIndexing(t *testing.T) {
	tests := []struct {
		name string
		idx  *index
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.idx.IsIndexing(); got != tt.want {
				t.Errorf("index.IsIndexing() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_index_NumberOfUUIDs(t *testing.T) {
	tests := []struct {
		name string
		idx  *index
		want uint32
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.idx.NumberOfUUIDs(); got != tt.want {
				t.Errorf("index.NumberOfUUIDs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_index_NumberOfUncommittedUUIDs(t *testing.T) {
	tests := []struct {
		name string
		idx  *index
		want uint32
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.idx.NumberOfUncommittedUUIDs(); got != tt.want {
				t.Errorf("index.NumberOfUncommittedUUIDs() = %v, want %v", got, tt.want)
			}
		})
	}
}
