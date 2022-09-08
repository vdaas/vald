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

// Package cache provides implementation of cache
package cache

import (
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/test/goleak"
)

// Goroutine leak is detected by `fastime`, but it should be ignored in the test because it is an external package.
var goleakIgnoreOptions = []goleak.Option{
	goleak.IgnoreTopFunction("github.com/kpango/fastime.(*fastime).StartTimerD.func1"),
}

func TestNew(t *testing.T) {
	type args struct {
		opts []Option
	}
	type want struct {
		wantCc Cache
		err    error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, Cache, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotCc Cache, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotCc, w.wantCc) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotCc, w.wantCc)
		}
		return nil
	}
	tests := []test{
		{
			name: "return gache cacher",
			args: args{
				opts: []Option{WithType("gache")},
			},
			checkFunc: func(w want, got Cache, err error) error {
				if err != nil {
					return err
				}
				if got == nil {
					return errors.New("got cache is nil")
				}

				return nil
			},
		},
		{
			name: "return unknown error when type is unknown",
			args: args{
				opts: []Option{WithType("unknown")},
			},
			want: want{
				err: errors.ErrInvalidCacherType,
			},
		},
		{
			name: "return cache when type is empty",
			args: args{
				opts: []Option{WithType("")},
			},
			checkFunc: func(w want, got Cache, err error) error {
				if err != nil {
					return err
				}
				if got == nil {
					return errors.New("got cache is nil")
				}

				return nil
			},
		},
		{
			name: "return unknown error when type is dummy string",
			args: args{
				opts: []Option{WithType("dummy")},
			},
			want: want{
				err: errors.ErrInvalidCacherType,
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
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

			gotCc, err := New(test.args.opts...)
			if err := checkFunc(test.want, gotCc, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
