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

func TestWithSearchParallel(t *testing.T) {
	type args struct {
		flag bool
	}
	tests := []struct {
		name string
		args args
		want SearchOption
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithSearchParallel(tt.args.flag); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithSearchParallel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithSearchConfig(t *testing.T) {
	type args struct {
		cfg *client.SearchConfig
	}
	tests := []struct {
		name string
		args args
		want SearchOption
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithSearchConfig(tt.args.cfg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithSearchConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
