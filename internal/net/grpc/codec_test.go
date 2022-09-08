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

// Package grpc provides generic functionality for grpc
package grpc

import (
	"reflect"
	"testing"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/errors"
)

func TestCodec_Marshal(t *testing.T) {
	t.Parallel()
	type args struct {
		v interface{}
	}
	type want struct {
		want []byte
		err  error
	}
	type test struct {
		name       string
		args       args
		c          Codec
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
			name: "return marshal result when val is vtproto message",
			args: args{
				v: &payload.Object_Vector{
					Id:     "1",
					Vector: []float32{1.0, 2.1},
				},
			},
			checkFunc: func(w want, b []byte, e error) error {
				if e != nil {
					return e
				}
				if len(b) == 0 {
					return errors.New("return byte is empty")
				}
				return nil
			},
		},
		{
			name: "return marshal result when val is empty proto message",
			args: args{
				v: &payload.Empty{},
			},
			want: want{
				want: []byte{},
				err:  nil,
			},
		},
		{
			name: "return error when val is not proto message",
			args: args{
				v: []int{1},
			},
			want: want{
				want: nil,
				err:  errors.ErrInvalidProtoMessageType([]int{1}),
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
			c := Codec{}

			got, err := c.Marshal(test.args.v)
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestCodec_Unmarshal(t *testing.T) {
	t.Parallel()
	type args struct {
		data []byte
		v    interface{}
	}
	type want struct {
		err error
	}
	type test struct {
		name       string
		args       args
		c          Codec
		want       want
		checkFunc  func(test, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(t test, err error) error {
		if !errors.Is(err, t.want.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, t.want.err)
		}
		return nil
	}
	tests := []test{
		{
			name: "unmarshal success to parse data into v",
			args: args{
				data: func() []byte {
					b, _ := Codec{}.Marshal(&payload.Object_Vector{
						Id:     "1",
						Vector: []float32{1.0, 2.1},
					})
					return b
				}(),
				v: &payload.Object_Vector{},
			},
			checkFunc: func(t test, e error) error {
				if !reflect.DeepEqual(t.args.v, &payload.Object_Vector{
					Id:     "1",
					Vector: []float32{1.0, 2.1},
				}) {
					return errors.New("unmarshal result is not correct")
				}
				return nil
			},
		},
		{
			name: "return error when data is invalid",
			args: args{
				data: []byte{0, 1, 2},
				v:    &payload.Object_Vector{},
			},
			want: want{
				err: errors.New("proto: Object_Vector: illegal tag 0 (wire type 0)"),
			},
		},
		{
			name: "return error when v is invalid",
			args: args{
				data: []byte{0, 1, 2},
				v:    Codec{},
			},
			want: want{
				err: errors.New("failed to marshal/unmarshal proto message, message type is grpc.Codec (missing vtprotobuf/protobuf helpers)"),
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
			c := Codec{}

			err := c.Unmarshal(test.args.data, test.args.v)
			if err := checkFunc(test, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestCodec_Name(t *testing.T) {
	t.Parallel()
	type want struct {
		want string
	}
	type test struct {
		name       string
		c          Codec
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
			name: "return codec name",
			want: want{
				want: "proto",
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
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
			c := Codec{}

			got := c.Name()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
