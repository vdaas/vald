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

// Package compress provides compress functions
package compress

import (
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"go.uber.org/goleak"
)

func TestWithLZ4Gob(t *testing.T) {
	type T = lz4Compressor
	type args struct {
		opts []GobOption
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got = %v, want %v", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		func() test {
			gobc, _ := NewGob()
			return test{
				name: "set success when opts is not nil.",
				args: args{
					opts: nil,
				},
				want: want{
					obj: &T{
						gobc: gobc,
					},
				},
			}
		}(),
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := WithLZ4Gob(test.args.opts...)
			obj := new(T)
			if err := test.checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithLZ4CompressionLevel(t *testing.T) {
	type T = lz4Compressor
	type args struct {
		level int
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got = %v, want %v", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		{
			name: "set success when level is less than 0.",
			args: args{
				level: -1,
			},
			want: want{
				obj: &T{
					compressionLevel: -1,
				},
			},
		},
		{
			name: "set success when level is nil.",
			want: want{
				obj: new(T),
			},
		},
		{
			name: "return error when level is more than 1.",
			args: args{
				level: 1,
			},
			want: want{
				obj: new(T),
				err: errors.ErrInvalidCompressionLevel(1),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}

			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := WithLZ4CompressionLevel(test.args.level)
			obj := new(T)
			if err := test.checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
