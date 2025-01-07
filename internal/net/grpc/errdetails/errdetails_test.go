//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
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

// Package errdetails provides error detail for gRPC status
package errdetails

import (
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/info"
	"github.com/vdaas/vald/internal/net/grpc/proto"
	"github.com/vdaas/vald/internal/net/grpc/types"
)

func Test_decodeDetails(t *testing.T) {
	t.Parallel()
	type args struct {
		objs []any
	}
	tests := []struct {
		name        string
		args        args
		wantDetails []Detail
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotDetails := decodeDetails(tt.args.objs...); !reflect.DeepEqual(gotDetails, tt.wantDetails) {
				t.Errorf("decodeDetails() = %v, want %v", gotDetails, tt.wantDetails)
			}
		})
	}
}

func TestSerialize(t *testing.T) {
	t.Parallel()
	type args struct {
		objs []any
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Serialize(tt.args.objs...); got != tt.want {
				t.Errorf("Serialize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnyToErrorDetail(t *testing.T) {
	t.Parallel()
	type args struct {
		a *types.Any
	}
	tests := []struct {
		name string
		args args
		want proto.Message
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AnyToErrorDetail(tt.args.a); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AnyToErrorDetail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDebugInfoFromInfoDetail(t *testing.T) {
	t.Parallel()
	type args struct {
		v *info.Detail
	}
	tests := []struct {
		name string
		args args
		want *DebugInfo
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DebugInfoFromInfoDetail(tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DebugInfoFromInfoDetail() = %v, want %v", got, tt.want)
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
//
// func TestDetail_MarshalJSON(t *testing.T) {
// 	type fields struct {
// 		TypeURL string
// 		Message proto.Message
// 	}
// 	type want struct {
// 		wantBody []byte
// 		err      error
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, []byte, error) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, gotBody []byte, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotBody, w.wantBody) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotBody, w.wantBody)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           TypeURL:"",
// 		           Message:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           fields: fields {
// 		           TypeURL:"",
// 		           Message:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			d := &Detail{
// 				TypeURL: test.fields.TypeURL,
// 				Message: test.fields.Message,
// 			}
//
// 			gotBody, err := d.MarshalJSON()
// 			if err := checkFunc(test.want, gotBody, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
