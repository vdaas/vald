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

	"github.com/vdaas/vald/apis/grpc/payload"
)

func TestNewFilter(t *testing.T) {
	type args struct {
		opts []FilterOption
	}
	tests := []struct {
		name    string
		args    args
		wantEf  Filter
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotEf, err := NewFilter(tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewFilter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotEf, tt.wantEf) {
				t.Errorf("NewFilter() = %v, want %v", gotEf, tt.wantEf)
			}
		})
	}
}

func Test_filter_Start(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		f       *filter
		args    args
		want    <-chan error
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.f.Start(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("filter.Start() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("filter.Start() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_filter_FilterSearch(t *testing.T) {
	type args struct {
		ctx context.Context
		res *payload.Search_Response
	}
	tests := []struct {
		name    string
		f       *filter
		args    args
		want    *payload.Search_Response
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.f.FilterSearch(tt.args.ctx, tt.args.res)
			if (err != nil) != tt.wantErr {
				t.Errorf("filter.FilterSearch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("filter.FilterSearch() = %v, want %v", got, tt.want)
			}
		})
	}
}
