//
// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
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

// Package errdetails provides error detail for grpc status
package errdetails

import (
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/info"
	"github.com/vdaas/vald/internal/net/grpc/proto"
	"github.com/vdaas/vald/internal/net/grpc/types"
)

func Test_decodeDetails(t *testing.T) {
	t.Parallel()
	type args struct {
		objs []interface{}
	}
	tests := []struct {
		name        string
		args        args
		wantDetails []Detail
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotDetails := decodeDetails(tt.args.objs...); !reflect.DeepEqual(gotDetails, tt.wantDetails) {
				t.Errorf("decodeDetails() = %v, want %v", gotDetails, tt.wantDetails)
			}
		})
	}
}

func TestSerialize(t *testing.T) {
	t.Parallel()
	type args struct {
		objs []interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Serialize(tt.args.objs...); got != tt.want {
				t.Errorf("Serialize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnyToErrorDetail(t *testing.T) {
	t.Parallel()
	type args struct {
		a *types.Any
	}
	tests := []struct {
		name string
		args args
		want proto.Message
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AnyToErrorDetail(tt.args.a); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AnyToErrorDetail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDebugInfoFromInfoDetail(t *testing.T) {
	t.Parallel()
	type args struct {
		v *info.Detail
	}
	tests := []struct {
		name string
		args args
		want *DebugInfo
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DebugInfoFromInfoDetail(tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DebugInfoFromInfoDetail() = %v, want %v", got, tt.want)
			}
		})
	}
}
