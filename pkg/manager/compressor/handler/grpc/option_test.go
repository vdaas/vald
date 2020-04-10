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

// Package grpc provides grpc server logic
package grpc

import (
	"reflect"
	"testing"

	"github.com/vdaas/vald/pkg/manager/compressor/service"
)

func TestWithCompressor(t *testing.T) {
	type args struct {
		c service.Compressor
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
			if got := WithCompressor(tt.args.c); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithCompressor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithBackup(t *testing.T) {
	type args struct {
		c service.Backup
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
			if got := WithBackup(tt.args.c); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithBackup() = %v, want %v", got, tt.want)
			}
		})
	}
}
