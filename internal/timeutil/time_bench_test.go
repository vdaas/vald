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
package timeutil

import (
	"testing"
	"time"
)

func BenchmarkParse(b *testing.B) {
	type args struct {
		t string
	}
	type test struct {
		name string
		args args
	}
	tests := []test{
		{
			name: "t 1ns",
			args: args{
				t: "1ns",
			},
		},
		{
			name: "t 1µs",
			args: args{
				t: "1µs",
			},
		},
		{
			name: "t 1ms",
			args: args{
				t: "1ms",
			},
		},
		{
			name: "t 1s",
			args: args{
				t: "1s",
			},
		},
		{
			name: "t 1m",
			args: args{
				t: "1m",
			},
		},
		{
			name: "t 1h",
			args: args{
				t: "1h",
			},
		},
	}
	for _, test := range tests {
		b.Run(test.name, func(b *testing.B) {
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					got, err := Parse(test.args.t)
					if err != nil {
						b.Error(err)
						return
					}
					if got == 0 {
						b.Error("got is 0")
					}
				}
			})
		})
	}
}

func BenchmarkParseWithDefault(b *testing.B) {
	type args struct {
		t string
		d time.Duration
	}
	type test struct {
		name string
		args args
	}
	tests := []test{
		{
			name: "t 1ns and d 50 millisecond",
			args: args{
				t: "1ns",
				d: 50 * time.Millisecond,
			},
		},
		{
			name: "t 1µs and d 50 millisecond",
			args: args{
				t: "1µs",
				d: 50 * time.Millisecond,
			},
		},
		{
			name: "t 1ms and d 50 millisecond",
			args: args{
				t: "1ms",
				d: 50 * time.Millisecond,
			},
		},
		{
			name: "t 1s and d 50 millisecond",
			args: args{
				t: "1s",
				d: 50 * time.Millisecond,
			},
		},
		{
			name: "t 1m and d 50 millisecond",
			args: args{
				t: "1m",
				d: 50 * time.Millisecond,
			},
		},
		{
			name: "t 1h and d 50 millisecond",
			args: args{
				t: "1h",
				d: 50 * time.Millisecond,
			},
		},

		// The following are error pattern.
		// So default value is returned.

		// The length of t is 5.
		{
			name: "t 1days and d 50 millisecond",
			args: args{
				t: "1days",
				d: 50 * time.Millisecond,
			},
		},
		// The length of t is 6.
		{
			name: "t 1month and d 50 millisecond",
			args: args{
				t: "1month",
				d: 50 * time.Millisecond,
			},
		},
		// The length of t is 7.
		{
			name: "t 10month and d 50 millisecond",
			args: args{
				t: "10month",
				d: 50 * time.Millisecond,
			},
		},
		// The length of t is 8.
		{
			name: "t 100month and d 50 millisecond",
			args: args{
				t: "100month",
				d: 50 * time.Millisecond,
			},
		},
		// The length of t is 9.
		{
			name: "t 1000month and d 50 millisecond",
			args: args{
				t: "1000month",
				d: 50 * time.Millisecond,
			},
		},
		// The length of t is 10.
		{
			name: "t 10000month and d 50 millisecond",
			args: args{
				t: "10000month",
				d: 50 * time.Millisecond,
			},
		},
	}
	for _, test := range tests {
		b.Run(test.name, func(b *testing.B) {
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					got := ParseWithDefault(test.args.t, test.args.d)
					if got == 0 {
						b.Error("got is 0")
					}
				}
			})
		})
	}
}
