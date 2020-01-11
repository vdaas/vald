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

	"github.com/vdaas/vald/internal/log"
)

func TestRecoverFunc(t *testing.T) {
	type args struct {
		fn func() error
	}

	type test struct {
		name       string
		args       args
		runtimeErr bool
		want       error
	}

	tests := []test{
		{
			name: "runtime error",
			args: args{
				fn: func() error {
					_ = []string{}[10]
					return nil
				},
			},
			runtimeErr: true,
			want:       fmt.Errorf("system paniced caused by runtime error: runtime error: index out of range [10] with length 0"),
		},

		{
			name: "panic string",
			args: args{
				fn: func() error {
					panic("panic")
				},
			},
			want: fmt.Errorf("panic recovered: panic"),
		},

		{
			name: "panic error",
			args: args{
				fn: func() error {
					panic(fmt.Errorf("error"))
				},
			},
			want: fmt.Errorf("error"),
		},

		{
			name: "default case panic",
			args: args{
				fn: func() error {
					panic(10)
				},
			},
			want: fmt.Errorf("panic recovered: 10"),
		},
	}

	log.Init(log.DefaultGlg())

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if ok := tt.runtimeErr; ok {
					if want, got := tt.want, recover().(error); want.Error() != got.Error() {
						t.Errorf("want: %v, got: %v", want, got)
					}
				}
			}()

			got := RecoverFunc(tt.args.fn)()

			if tt.want == nil && got != nil {
				t.Errorf("RecoverFunc return error: %v", got)
			} else if tt.want != nil {
				if tt.want.Error() != got.Error() {
					t.Errorf("RecoverFunc is wrong, want: %v, got: %v", tt.want, got)
				}
			}
		})
	}
}
