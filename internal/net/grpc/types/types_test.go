//
// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
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

// Package types provides alias of protobuf library types
package types

import (
	"testing"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/test/goleak"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/runtime/protoimpl"
	"google.golang.org/protobuf/types/known/anypb"
)

func TestUnmarshalAny(t *testing.T) {
	t.Parallel()
	type args struct {
		any *Any
		pb  proto.Message
	}
	type want struct {
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		return nil
	}
	tests := []test{
		{
			name: "return nil when success to unmarshal insert",
			args: args{
				any: func() *Any {
					anyVal, err := anypb.New(new(payload.Insert_Request))
					if err != nil {
						t.Error(err)
					}
					return anyVal
				}(),
				pb: &payload.Insert_Request{
					Vector: &payload.Object_Vector{
						Id:     "1",
						Vector: []float32{1.0, 2.1, 3.1},
					},
				},
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "return error when unmarshal type mismatch",
			args: args{
				any: func() *Any {
					anyVal, err := anypb.New(new(payload.Insert_Request))
					if err != nil {
						t.Error(err)
					}
					return anyVal
				}(),
				pb: &payload.Object_Vector{
					Id:     "1",
					Vector: []float32{1.0, 2.1, 3.1},
				},
			},
			want: want{
				err: protoimpl.X.NewError("mismatched message type: got \"payload.v1.Object.Vector\", want \"payload.v1.Insert.Request\""),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
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

			err := UnmarshalAny(test.args.any, test.args.pb)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
