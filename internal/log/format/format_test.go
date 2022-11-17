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
package format

import (
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/test/goleak"
)

func TestFormat_String(t *testing.T) {
	type want struct {
		want string
	}
	type test struct {
		name       string
		f          Format
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
			name: "returns raw when f is RAW",
			f:    RAW,
			want: want{
				want: "raw",
			},
		},

		{
			name: "returns json when f is JSON",
			f:    JSON,
			want: want{
				want: "json",
			},
		},

		// {
		// 	name: "returns ltsv when f is LTSV",
		// 	f:    LTSV,
		// 	want: want{
		// 		want: "ltsv",
		// 	},
		// },

		{
			name: "returns unknown when f is 100",
			f:    Format(100),
			want: want{
				want: "unknown",
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
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

			got := test.f.String()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestAtof(t *testing.T) {
	type want struct {
		want Format
	}
	type test struct {
		name       string
		str        string
		want       want
		checkFunc  func(want, Format) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got Format) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}

	tests := []test{
		{
			name: "returns RAW when str is `raw`",
			str:  "raw",
			want: want{
				want: RAW,
			},
		},

		{
			name: "returns RAW when str is `RAw`",
			str:  "RAw",
			want: want{
				want: RAW,
			},
		},

		{
			name: "returns JSON when str is `json`",
			str:  "json",
			want: want{
				want: JSON,
			},
		},

		{
			name: "returns JSON when str is `JSOn`",
			str:  "JSOn",
			want: want{
				want: JSON,
			},
		},

		// {
		// 	name: "returns LTSV when str is `ltsv`",
		// 	str:  "ltsv",
		// 	want: want{
		// 		want: LTSV,
		// 	},
		// },
		//
		// {
		// 	name: "returns LTSV when str is `LTSv`",
		// 	str:  "LTSv",
		// 	want: want{
		// 		want: LTSV,
		// 	},
		// },

		{
			name: "returns Unknown when str is `Vald`",
			str:  "Vald",
			want: want{
				want: Unknown,
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
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
			got := Atof(test.str)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
