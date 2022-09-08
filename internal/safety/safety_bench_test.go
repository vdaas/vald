// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package safety

import (
	"testing"

	"github.com/vdaas/vald/internal/errors"
)

func BenchmarkRecoverFunc(b *testing.B) {
	type args struct {
		fn func() error
	}
	type test struct {
		name string
		args args
	}
	tests := []test{
		{
			name: "fn return error func",
			args: args{
				fn: func() error {
					return errors.New("fn err")
				},
			},
		},
	}
	for _, tc := range tests {
		test := tc
		b.Run(test.name, func(b *testing.B) {
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					got := RecoverFunc(test.args.fn)
					if got == nil {
						b.Error("got is empty")
					}
				}
			})
		})
	}
}

func BenchmarkRecoverWithoutPanicFunc(b *testing.B) {
	type args struct {
		fn func() error
	}
	type test struct {
		name string
		args args
	}
	tests := []test{
		{
			name: "fn return error func",
			args: args{
				fn: func() error {
					return errors.New("fn err")
				},
			},
		},
		{
			name: "fn panic runtime error func",
			args: args{
				fn: func() error {
					_ = []string{}[10]
					return nil
				},
			},
		},
		{
			name: "fn panic string func",
			args: args{
				fn: func() error {
					panic("panic")
				},
			},
		},
		{
			name: "fn panic error func",
			args: args{
				fn: func() error {
					panic(errors.Errorf("error"))
				},
			},
		},
		{
			name: "fn panic int func",
			args: args{
				fn: func() error {
					panic(10)
				},
			},
		},
	}
	for _, tc := range tests {
		test := tc
		b.Run(test.name, func(b *testing.B) {
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					got := RecoverWithoutPanicFunc(test.args.fn)
					if got == nil {
						b.Error("got is empty")
					}
					err := got()
					if err == nil {
						b.Error("err is nil")
					}
				}
			})
		})
	}
}
