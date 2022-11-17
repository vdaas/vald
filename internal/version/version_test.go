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
package version

import (
	"fmt"
	"testing"

	"github.com/vdaas/vald/internal/errors"
)

func TestCheck(t *testing.T) {
	type args struct {
		cur string
		max string
		min string
	}

	type test struct {
		name string
		args args
		want error
	}

	tests := []test{
		{
			name: "return nil",
			args: args{
				cur: "1.0.5",
				max: "1.0.10",
				min: "1.0.0",
			},
			want: nil,
		},

		{
			name: "return error when cur format is invalid",
			args: args{
				cur: "vald",
				max: "1.0.10",
				min: "1.0.0",
			},
			want: fmt.Errorf("Malformed version: vald"),
		},

		{
			name: "return error when min format is invalid",
			args: args{
				cur: "1.5.10",
				max: "vald",
				min: "1.0.0",
			},
			want: fmt.Errorf("Malformed constraint:  <= vald"),
		},

		{
			name: "return error when min format is invalid",
			args: args{
				cur: "1.0.10",
				max: "1.0.5",
				min: "1.0.0",
			},
			want: errors.ErrInvalidConfigVersion("1.0.10", ">= 1.0.0, <= 1.0.5"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Check(tt.args.cur, tt.args.max, tt.args.min)
			if !errors.Is(err, tt.want) {
				t.Errorf("got: %s, want: %s", err, tt.want)
			}
		})
	}
}
