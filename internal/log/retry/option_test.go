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
package retry

import (
	"errors"
	"reflect"
	"testing"
)

func TestWithError(t *testing.T) {
	type test struct {
		name      string
		fn        func(vals ...interface{})
		checkFunc func(Option) error
	}

	tests := []test{
		func() test {
			fn := func(vals ...interface{}) {}

			return test{
				name: "set success when fn is not nil",
				fn:   fn,
				checkFunc: func(opt Option) error {
					got := new(retry)
					opt(got)

					if reflect.ValueOf(fn).Pointer() != reflect.ValueOf(got.errorFn).Pointer() {
						return errors.New("invalid params was set")
					}
					return nil
				},
			}
		}(),

		func() test {
			fn := func(vals ...interface{}) {}

			return test{
				name: "returns nothing when fn is nil",
				fn:   nil,
				checkFunc: func(opt Option) error {
					got := &retry{
						errorFn: fn,
					}
					opt(got)

					if reflect.ValueOf(fn).Pointer() != reflect.ValueOf(got.errorFn).Pointer() {
						return errors.New("invalid params was set")
					}
					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithError(tt.fn)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithWarn(t *testing.T) {
	type test struct {
		name      string
		fn        func(vals ...interface{})
		checkFunc func(Option) error
	}

	tests := []test{
		func() test {
			fn := func(vals ...interface{}) {}

			return test{
				name: "set success when fn is not nil",
				fn:   fn,
				checkFunc: func(opt Option) error {
					got := new(retry)
					opt(got)

					if reflect.ValueOf(fn).Pointer() != reflect.ValueOf(got.warnFn).Pointer() {
						return errors.New("invalid params was set")
					}
					return nil
				},
			}
		}(),

		func() test {
			fn := func(vals ...interface{}) {}

			return test{
				name: "returns nothing when fn is nil",
				fn:   nil,
				checkFunc: func(opt Option) error {
					got := &retry{
						warnFn: fn,
					}
					opt(got)

					if reflect.ValueOf(fn).Pointer() != reflect.ValueOf(got.warnFn).Pointer() {
						return errors.New("invalid params was set")
					}
					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithWarn(tt.fn)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}
