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

// Package strategy provides strategy for e2e testing functions
package strategy

import (
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/client"
)

func TestWithCreateIndexPoolSize(t *testing.T) {
	type args struct {
		size int
	}
	tests := []struct {
		name string
		args args
		want CreateIndexOption
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithCreateIndexPoolSize(tt.args.size); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithCreateIndexPoolSize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithCreateIndexClient(t *testing.T) {
	type args struct {
		c client.Indexer
	}
	tests := []struct {
		name string
		args args
		want CreateIndexOption
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithCreateIndexClient(tt.args.c); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithCreateIndexClient() = %v, want %v", got, tt.want)
			}
		})
	}
}
