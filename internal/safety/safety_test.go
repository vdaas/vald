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
package safety

import (
	"fmt"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
)

func TestRecoverFunc(t *testing.T) {
	type test struct {
		name       string
		fn         func() error
		runtimeErr bool
		want       error
	}

	tests := []test{
		{
			name: "returns error when system paniced caused by runtime error",
			fn: func() error {
				_ = []string{}[10]
				return nil
			},
			runtimeErr: true,
			want:       errors.New("system paniced caused by runtime error: runtime error: index out of range [10] with length 0"),
		},

		{
			name: "returns error when system paniced caused by panic with string value",
			fn: func() error {
				panic("panic")
			},
			want: errors.New("panic recovered: panic"),
		},

		{
			name: "returns error when system paniced caused by panic with error",
			fn: func() error {
				panic(fmt.Errorf("error"))
			},
			want: errors.New("error"),
		},

		{
			name: "returns error when system paniced caused by panic with int value",
			fn: func() error {
				panic(10)
			},
			want: errors.New("panic recovered: 10"),
		},
	}

	log.Init()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if ok := tt.runtimeErr; ok {
					if want, got := tt.want, recover().(error); !errors.Is(got, want) {
						t.Errorf("not equals. want: %v, got: %v", want, got)
					}
				}
			}()

			got := RecoverFunc(tt.fn)()
			if !errors.Is(got, tt.want) {
				t.Errorf("not equals. want: %v, got: %v", tt.want, got)
			}
		})
	}
}
