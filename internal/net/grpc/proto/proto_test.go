//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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

// Package proto provides proto file logic
package proto

import (
	"reflect"
	"testing"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/test/comparator"
	"github.com/vdaas/vald/internal/test/goleak"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestMarshal(t *testing.T) {
	t.Parallel()
	type args struct {
		m Message
	}
	type want struct {
		want []byte
		err  error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, []byte, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got []byte, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return proto marshal result",
			args: args{
				m: &payload.Object_Vector{
					Id:     "1",
					Vector: []float32{1.0, 2.1},
				},
			},
			checkFunc: func(w want, b []byte, e error) error {
				if len(b) == 0 {
					return errors.New("return byte is empty")
				}
				return nil
			},
		},
		{
			name: "return empty byte when marshal a empty payload",
			args: args{
				m: &payload.Empty{},
			},
			want: want{
				want: []byte{},
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got, err := Marshal(test.args.m)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestUnmarshal(t *testing.T) {
	t.Parallel()
	type args struct {
		data []byte
		v    Message
	}
	type want struct {
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(test, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(test test, err error) error {
		if !errors.Is(err, test.want.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, test.want.err)
		}
		return nil
	}
	tests := []test{
		func() test {
			o := &payload.Object_Vector{
				Id:     "1",
				Vector: []float32{1.0, 2.1},
			}

			b, err := Marshal(o)
			if err != nil {
				t.Error(err)
			}

			return test{
				name: "return nil when unmarshal success",
				args: args{
					data: b,
					v:    &payload.Object_Vector{},
				},
				want: want{
					err: nil,
				},
				checkFunc: func(t test, e error) error {
					got, ok := t.args.v.(*payload.Object_Vector)
					if !ok {
						return errors.New("cannot cast result to object")
					}
					if got.Id != o.Id || !reflect.DeepEqual(got.Vector, o.Vector) {
						return errors.Errorf("got: %#v, want: %#v", got, o)
					}
					return nil
				},
			}
		}(),
		func() test {
			return test{
				name: "return error when unmarshal failed",
				args: args{
					data: []byte{1, 2, 3, 4},
					v:    &payload.Object_Vector{},
				},
				checkFunc: func(t test, e error) error {
					if e == nil {
						return errors.New("error should be returned")
					}
					return nil
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			err := Unmarshal(test.args.data, test.args.v)
			if err := test.checkFunc(test, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestClone(t *testing.T) {
	t.Parallel()
	type args struct {
		m Message
	}
	type want struct {
		want Message
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, Message) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got Message) error {
		diff := comparator.Diff(got, w.want, comparator.IgnoreUnexported(
			payload.Object_Vector{},
		))
		if diff != "" {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return the clone message",
			args: args{
				m: &payload.Object_Vector{
					Id:     "1",
					Vector: []float32{1.0, 2.1},
				},
			},
			want: want{
				want: &payload.Object_Vector{
					Id:     "1",
					Vector: []float32{1.0, 2.1},
				},
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
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

			got := Clone(test.args.m)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestToMessageV1(t *testing.T) {
	t.Parallel()
	type args struct {
		m Message
	}
	type want struct {
		want MessageV1
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, MessageV1) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got MessageV1) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return converted messageV1",
			args: args{
				m: &payload.Object_Vector{
					Id:     "1",
					Vector: []float32{1.0, 2.1},
				},
			},
			want: want{
				want: &payload.Object_Vector{
					Id:     "1",
					Vector: []float32{1.0, 2.1},
				},
			},
		},
		{
			name: "return nil when the message is nil",
			args: args{
				m: nil,
			},
			want: want{
				want: nil,
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := ToMessageV1(test.args.m)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
