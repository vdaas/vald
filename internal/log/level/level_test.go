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
package level

import (
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"go.uber.org/goleak"
)

func TestString(t *testing.T) {
	type test struct {
		name  string
		level Level
		want  string
	}

	tests := []test{
		{
			name:  "returns DEBUG",
			level: DEBUG,
			want:  "DEBUG",
		},

		{
			name:  "returns INFO",
			level: INFO,
			want:  "INFO",
		},

		{
			name:  "returns WARN",
			level: WARN,
			want:  "WARN",
		},

		{
			name:  "returns ERROR",
			level: ERROR,
			want:  "ERROR",
		},

		{
			name:  "returns FATAL",
			level: FATAL,
			want:  "FATAL",
		},

		{
			name:  "returns Unknown",
			level: Level(100),
			want:  "Unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.level.String()
			if got != tt.want {
				t.Errorf("not equals. want: %v, but got: %v", tt.want, got)
			}
		})
	}
}

func TestAtol(t *testing.T) {
	type test struct {
		name string
		str  string
		want Level
	}

	tests := []test{
		{
			name: "returns DEBUG when str is debug",
			str:  "debug",
			want: DEBUG,
		},

		{
			name: "returns DEBUG when str is deb",
			str:  "deb",
			want: DEBUG,
		},

		{
			name: "returns DEBUG when str is DEBUg",
			str:  "DEBUg",
			want: DEBUG,
		},

		{
			name: "returns INFO when str is info",
			str:  "info",
			want: INFO,
		},

		{
			name: "returns INFO when str is INFo",
			str:  "INFo",
			want: INFO,
		},

		{
			name: "returns WARN when str is warn",
			str:  "warn",
			want: WARN,
		},

		{
			name: "returns WARN when str is WARn",
			str:  "WARn",
			want: WARN,
		},

		{
			name: "returns ERROR when str is error",
			str:  "error",
			want: ERROR,
		},

		{
			name: "returns ERROR when str is err",
			str:  "err",
			want: ERROR,
		},

		{
			name: "returns ERROR when str is ERROr",
			str:  "ERROr",
			want: ERROR,
		},

		{
			name: "returns FATAL when str is fatal",
			str:  "fatal",
			want: FATAL,
		},

		{
			name: "returns FATAL when str is FATAl",
			str:  "FATAl",
			want: FATAL,
		},

		{
			name: "returns Unknown when str is vald",
			str:  "vald",
			want: Unknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Atol(tt.str)
			if got != tt.want {
				t.Errorf("not equals. want: %v, but got: %v", tt.want, got)
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
			return errors.Errorf("got = %v, want %v", got, w.want)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       want: want{},
		       checkFunc: defaultCheckFunc,
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := test.l.String()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}
