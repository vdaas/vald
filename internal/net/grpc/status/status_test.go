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

// Package status provides statuses and errors returned by grpc handler functions
package status

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/info"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/errdetails"
	"github.com/vdaas/vald/internal/test/goleak"
)

func TestMain(m *testing.M) {
	info.Init("")
	goleak.VerifyTestMain(m)
}

func TestParseError(t *testing.T) {
	t.Parallel()
	type args struct {
		err         error
		defaultCode codes.Code
		defaultMsg  string
		details     []interface{}
	}
	type want struct {
		wantSt  codes.Code
		wantMsg string
		err     error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *Status, string, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotSt *Status, gotMsg string, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotSt, w.wantSt) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotSt, w.wantSt)
		}
		if !reflect.DeepEqual(gotMsg, w.wantMsg) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotMsg, w.wantMsg)
		}
		return nil
	}
	tests := []test{
		func() test {
			return test{
				name: "test_case_1",
				args: args{
					err: WrapWithNotFound(errors.ErrEmptySearchResult.Error(), errors.ErrEmptySearchResult,
						&errdetails.RequestInfo{
							RequestId: "sample request ID",
						},
						&errdetails.ResourceInfo{
							ResourceType: "sample resource type",
							ResourceName: "sample resource name",
						}, info.Get()),
					defaultCode: codes.Internal,
					defaultMsg:  "failed to parse Search gRPC error response",
					details:     nil,
				},
				want: want{
					wantSt: codes.NotFound,
				},
				checkFunc: func(w want, gotSt *Status, gotMsg string, err error) error {
					if w.wantSt != gotSt.Code() {
						b, err := json.MarshalIndent(gotSt.Details(), "", "\t")
						t.Log(gotSt.String(), string(b), err)
						return errors.ErrEmptySearchResult
					}
					return nil
				},
			}
		}(),
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           err: nil,
		           defaultCode: nil,
		           defaultMsg: "",
		           details: nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           err: nil,
		           defaultCode: nil,
		           defaultMsg: "",
		           details: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
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

			gotSt, gotMsg, err := ParseError(test.args.err, test.args.defaultCode, test.args.defaultMsg, test.args.details...)
			if err := checkFunc(test.want, gotSt, gotMsg, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
