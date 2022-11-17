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
package level

import (
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/test/goleak"
)

func TestAtol(t *testing.T) {
	type want struct {
		want Level
	}
	type test struct {
		name       string
		str        string
		checkFunc  func(want, Level) error
		beforeFunc func()
		afterFunc  func()
		want       want
	}

	defaultCheckFunc := func(w want, got Level) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}

	tests := []test{
		{
			name: "returns DEBUG when str is `debug`",
			str:  "debug",
			want: want{
				want: DEBUG,
			},
		},

		{
			name: "returns DEBUG when str is `deb`",
			str:  "deb",
			want: want{
				want: DEBUG,
			},
		},

		{
			name: "returns DEBUG when str is `DEBUg`",
			str:  "DEBUg",
			want: want{
				want: DEBUG,
			},
		},

		{
			name: "returns INFO when str is `info`",
			str:  "info",
			want: want{
				want: INFO,
			},
		},

		{
			name: "returns INFO when str is `INFo`",
			str:  "INFo",
			want: want{
				want: INFO,
			},
		},

		{
			name: "returns INFO when str is `INFos`",
			str:  "INFos",
			want: want{
				want: INFO,
			},
		},

		{
			name: "returns WARN when str is `warn`",
			str:  "warn",
			want: want{
				want: WARN,
			},
		},

		{
			name: "returns WARN when str is `WARn`",
			str:  "WARn",
			want: want{
				want: WARN,
			},
		},

		{
			name: "returns WARN when str is `WARns`",
			str:  "WARns",
			want: want{
				want: WARN,
			},
		},

		{
			name: "returns ERROR when str is `error`",
			str:  "error",
			want: want{
				want: ERROR,
			},
		},

		{
			name: "returns ERROR when str is `err`",
			str:  "err",
			want: want{
				want: ERROR,
			},
		},

		{
			name: "returns ERROR when str is `ERROr`",
			str:  "ERROr",
			want: want{
				want: ERROR,
			},
		},

		{
			name: "returns FATAL when str is `fatal`",
			str:  "fatal",
			want: want{
				want: FATAL,
			},
		},

		{
			name: "returns FATAL when str is `FATAl`",
			str:  "FATAl",
			want: want{
				want: FATAL,
			},
		},

		{
			name: "returns FATAL when str is `FATAls`",
			str:  "FATAls",
			want: want{
				want: FATAL,
			},
		},

		{
			name: "returns Unknown when str is `vald`",
			str:  "vald",
			want: want{
				want: Unknown,
			},
		},
		{
			name: "returns Unknown when str is `va`",
			str:  "va",
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
			got := Atol(test.str)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestLevel_String(t *testing.T) {
	type want struct {
		want string
	}
	type test struct {
		name       string
		l          Level
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
			name: "returns `DEBUG` when l is DEBUG",
			l:    DEBUG,
			want: want{
				want: "DEBUG",
			},
		},

		{
			name: "returns `INFO` when l is INFO",
			l:    INFO,
			want: want{
				want: "INFO",
			},
		},

		{
			name: "returns `WARN` when l is WARN",
			l:    WARN,
			want: want{
				want: "WARN",
			},
		},

		{
			name: "returns `ERROR` when l is ERROR",
			l:    ERROR,
			want: want{
				want: "ERROR",
			},
		},

		{
			name: "returns `FATAL` when l is FATAL",
			l:    FATAL,
			want: want{
				want: "FATAL",
			},
		},

		{
			name: "returns `Unknown` when l is 100",
			l:    Level(100),
			want: want{
				want: "Unknown",
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

			got := test.l.String()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
