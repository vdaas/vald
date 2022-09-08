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

// Package cacher provides implementation of cache type definition
package cacher

import (
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errors"
)

func TestType_String(t *testing.T) {
	type want struct {
		want string
	}
	type test struct {
		name       string
		m          Type
		want       want
		checkFunc  func(want, string) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got string) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return `gache` when type is gache",
			m:    GACHE,
			want: want{
				want: "gache",
			},
		},
		{
			name: "return `unknown` when type is unknown",
			m:    Unknown,
			want: want{
				want: "unknown",
			},
		},
		{
			name: "return `unknown` when the type is invalid",
			m:    Type(100),
			want: want{
				want: "unknown",
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			got := test.m.String()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestToType(t *testing.T) {
	type args struct {
		str string
	}
	type want struct {
		want Type
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, Type) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got Type) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return GACHE type when the string is `gache`",
			args: args{
				str: "gache",
			},
			want: want{
				want: GACHE,
			},
		},
		{
			name: "return GACHE type when the string is `Gache`",
			args: args{
				str: "Gache",
			},
			want: want{
				want: GACHE,
			},
		},
		{
			name: "return GACHE type when the string is `GACHE`",
			args: args{
				str: "GACHE",
			},
			want: want{
				want: GACHE,
			},
		},
		{
			name: "return Unknown type when the string is invalid",
			args: args{
				str: "invalid",
			},
			want: want{
				want: Unknown,
			},
		},
		{
			name: "return Unknown type when the string is empty",
			args: args{
				str: "",
			},
			want: want{
				want: Unknown,
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			got := ToType(test.args.str)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
