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

// Package config providers configuration type and load configuration logic
package config

import (
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errors"
)

func TestLogging_Bind(t *testing.T) {
	type fields struct {
		Logger string
		Level  string
		Format string
	}
	type want struct {
		want *Logging
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, *Logging) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got *Logging) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "returns Logging when all fields contain no prefix/suffix symbol",
			fields: fields{
				Logger: "logger",
				Level:  "info",
				Format: "json",
			},
			want: want{
				want: &Logging{
					Logger: "logger",
					Level:  "info",
					Format: "json",
				},
			},
		},
		{
			name: "returns Logging with environment variable when it contains `_` prefix and suffix",
			fields: fields{
				Logger: "_LOGGING_BIND_LOGGER_",
				Level:  "_LOGGING_BIND_LEVEL_",
				Format: "_LOGGING_BIND_FORMAT_",
			},
			beforeFunc: func() {
				t.Setenv("LOGGING_BIND_LOGGER", "glg")
				t.Setenv("LOGGING_BIND_LEVEL", "info")
				t.Setenv("LOGGING_BIND_FORMAT", "json")
			},
			want: want{
				want: &Logging{
					Logger: "glg",
					Level:  "info",
					Format: "json",
				},
			},
		},
		{
			name: "returns Logging when all fields are empty",
			want: want{
				want: new(Logging),
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
			l := &Logging{
				Logger: test.fields.Logger,
				Level:  test.fields.Level,
				Format: test.fields.Format,
			}

			got := l.Bind()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
