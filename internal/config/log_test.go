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

// Package config providers configuration type and load configuration logic
package config

import (
	"os"
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
				Logger: "_logger_",
				Level:  "_level_",
				Format: "_format_",
			},
			beforeFunc: func() {
				_ = os.Setenv("logger", "glg")
				_ = os.Setenv("level", "info")
				_ = os.Setenv("format", "json")
			},
			afterFunc: func() {
				_ = os.Unsetenv("logger")
				_ = os.Unsetenv("level")
				_ = os.Unsetenv("format")
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

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			l := &Logging{
				Logger: test.fields.Logger,
				Level:  test.fields.Level,
				Format: test.fields.Format,
			}

			got := l.Bind()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}
