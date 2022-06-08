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
package grpc

import (
	"context"
	"reflect"
	"testing"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/test/goleak"
	"github.com/vdaas/vald/pkg/agent/core/ngt/service"
)

func Test_server_Upsert(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
		req *payload.Upsert_Request
	}
	type fields struct {
		name              string
		ip                string
		ngt               service.NGT
		eg                errgroup.Group
		streamConcurrency int
	}
	type want struct {
		wantLoc *payload.Object_Location
		err     error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *payload.Object_Location, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotLoc *payload.Object_Location, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotLoc, w.wantLoc) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotLoc, w.wantLoc)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           ctx: nil,
		           req: nil,
		       },
		       fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
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
		           ctx: nil,
		           req: nil,
		           },
		           fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
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
			s := &server{
				name:              test.fields.name,
				ip:                test.fields.ip,
				ngt:               test.fields.ngt,
				eg:                test.fields.eg,
				streamConcurrency: test.fields.streamConcurrency,
			}

			gotLoc, err := s.Upsert(test.args.ctx, test.args.req)
			if err := checkFunc(test.want, gotLoc, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_StreamUpsert(t *testing.T) {
	t.Parallel()
	type args struct {
		stream vald.Upsert_StreamUpsertServer
	}
	type fields struct {
		name              string
		ip                string
		ngt               service.NGT
		eg                errgroup.Group
		streamConcurrency int
	}
	type want struct {
		err error
	}
	type test struct {
		name       string
		args       args
		fields     fields
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

	/*
		Upsert test cases (only test float32 unless otherwise specified):
		- Equivalence Class Testing ( 1000 vectors inserted before an upsert )
			- case 1.1: success upsert with new ID
			- case 1.2: success upsert with existent ID
			- case 2.1: fail upsert with one different dimension vector (type: uint8)
			- case 2.2: fail upsert with one different dimension vector (type: float32)
		- Boundary Value Testing ( 1000 vectors inserted before an upsert if not specified )
			- case 1.1: fail upsert with "" as ID
			- case 1.2: fail upsert to empty index with "" as ID
			- case 2.1: success upsert with ^@ as duplicated ID
			- case 2.2: success upsert with ^@ as new ID
			- case 2.3: success upsert to empty index with ^@ as new ID
			- case 2.4: success upsert with ^I as duplicated ID
			- case 2.5: success upsert with ^I as new ID
			- case 2.6: success upsert to empty index with ^I as new ID
			- case 2.7: success upsert with ^J as duplicated ID
			- case 2.8: success upsert with ^J as new ID
			- case 2.9: success upsert to empty index with ^J as new ID
			- case 2.10: success upsert with ^M as duplicated ID
			- case 2.11: success upsert with ^M as new ID
			- case 2.12: success upsert to empty index with ^M as new ID
			- case 2.13: success upsert with ^[ as duplicated ID
			- case 2.14: success upsert with ^[ as new ID
			- case 2.15: success upsert to empty index with ^[ as new ID
			- case 2.16: success upsert with ^? as duplicated ID
			- case 2.17: success upsert with ^? as new ID
			- case 2.18: success upsert to empty index with ^? as new ID
			- case 3.1: success upsert with utf-8 ID from utf-8 index
			- case 3.2: success upsert with utf-8 ID from s-jis index
			- case 3.3: success upsert with utf-8 ID from euc-jp index
			- case 3.4: success upsert with s-jis ID from utf-8 index
			- case 3.5: success upsert with s-jis ID from s-jis index
			- case 3.6: success upsert with s-jis ID from euc-jp index
			- case 3.4: success upsert with euc-jp ID from utf-8 index
			- case 3.5: success upsert with euc-jp ID from s-jis index
			- case 3.6: success upsert with euc-jp ID from euc-jp index
			- case 4.1: success upsert with 😀 as duplicated ID
			- case 4.2: success upsert with 😀 as new ID
			- case 4.3: success upsert to empty index with 😀 as new ID
			- case 5.1: success upsert with one 0 value vector with duplicated ID (type: uint8)
			- case 5.2: success upsert with one 0 value vector with new ID (type: uint8)
			- case 5.3: success upsert to empty index with one 0 value vector with new ID (type: uint8)
			- case 5.4: success upsert with one +0 value vector with duplicated ID (type: float32)
			- case 5.5: success upsert with one +0 value vector with new ID (type: float32)
			- case 5.6: success upsert to empty index with one +0 value vector with new ID (type: float32)
			- case 5.7: success upsert with one -0 value vector with duplicated ID (type: float32)
			- case 5.8: success upsert with one -0 value vector with new ID (type: float32)
			- case 5.9: success upsert to empty index with one -0 value vector with new ID (type: float32)
			- case 6.1: success upsert with one min value vector with duplicated ID (type: uint8)
			- case 6.2: success upsert with one min value vector with new ID (type: uint8)
			- case 6.3: success upsert to empty index with one min value vector with new ID (type: uint8)
			- case 6.4: success upsert with one min value vector with duplicated ID (type: float32)
			- case 6.5: success upsert with one min value vector with new ID (type: float32)
			- case 6.6: success upsert to empty index with one min value vector with new ID (type: float32)
			- case 7.1: success upsert with one max value vector with duplicated ID (type: uint8)
			- case 7.2: success upsert with one max value vector with new ID (type: uint8)
			- case 7.3: success upsert to empty index with one max value vector with new ID (type: uint8)
			- case 7.4: success upsert with one max value vector with duplicated ID (type: float32)
			- case 7.5: success upsert with one max value vector (type: float32)
			- case 7.6: success upsert to empty index with one max value vector (type: float32)
			- case 8.1: success upsert with one NaN value vector with duplicated ID (type: float32) // NOTE: To fix it, it is necessary to check all of vector value
			- case 8.2: success upsert with one NaN value vector with new ID (type: float32) // NOTE: To fix it, it is necessary to check all of vector value
			- case 8.3: success upsert to empty index with one NaN value vector with new ID (type: float32) // NOTE: To fix it, it is necessary to check all of vector value
			- case 9.1: success upsert with one +inf value vector with duplicated ID (type: float32)
			- case 9.2: success upsert with one +inf value vector with new ID (type: float32)
			- case 9.3: success upsert to empty index with one +inf value vector with new ID (type: float32)
			- case 9.4: success upsert with one -inf value vector with duplicated ID (type: float32)
			- case 9.5: success upsert with one -inf value vector with new ID (type: float32)
			- case 9.6: success upsert to empty index with one -inf value vector with new ID (type: float32)
			- case 10.1: fail upsert with one nil vector
			- case 10.2: fail upsert to empty index with one nil vector
			- case 11.1: fail upsert with one empty vector
			- case 11.2: fail upsert to empty index with one empty vector
		- Decision Table Testing
			- case 1.1: fail upsert with one duplicated vector, duplicated ID and SkipStrictExistCheck is true
			- case 1.2: success upsert with one different vector, duplicated ID and SkipStrictExistsCheck is true
			- case 1.3: success upsert with one duplicated vector, different ID and SkipStrictExistCheck is true
			- case 2.1: success upsert with one duplicated vector, new ID and SkipStrictExistCheck is true
			- case 2.2: success upsert with one different vector, new ID and SkipStrictExistsCheck is true
			- case 2.3: success upsert with one duplicated vector, new ID and SkipStrictExistCheck is true
			- case 3.1: fail upsert with one duplicated vector, duplicated ID and SkipStrictExistCheck is false
			- case 3.2: success upsert with one different vector, duplicated ID and SkipStrictExistsCheck is false
			- case 3.3: success upsert with one duplicated vector, different ID and SkipStrictExistCheck is false
			- case 4.1: success upsert with one duplicated vector, new ID and SkipStrictExistCheck is false
			- case 4.2: success upsert with one different vector, new ID and SkipStrictExistsCheck is false
			- case 4.3: success upsert with one duplicated vector, new ID and SkipStrictExistCheck is false
	*/
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           stream: nil,
		       },
		       fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
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
		           stream: nil,
		           },
		           fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
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
			s := &server{
				name:              test.fields.name,
				ip:                test.fields.ip,
				ngt:               test.fields.ngt,
				eg:                test.fields.eg,
				streamConcurrency: test.fields.streamConcurrency,
			}

			err := s.StreamUpsert(test.args.stream)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_server_MultiUpsert(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx  context.Context
		reqs *payload.Upsert_MultiRequest
	}
	type fields struct {
		name              string
		ip                string
		ngt               service.NGT
		eg                errgroup.Group
		streamConcurrency int
	}
	type want struct {
		wantRes *payload.Object_Locations
		err     error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *payload.Object_Locations, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotRes *payload.Object_Locations, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotRes, w.wantRes) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           ctx: nil,
		           reqs: nil,
		       },
		       fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
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
		           ctx: nil,
		           reqs: nil,
		           },
		           fields: fields {
		           name: "",
		           ip: "",
		           ngt: nil,
		           eg: nil,
		           streamConcurrency: 0,
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
			s := &server{
				name:              test.fields.name,
				ip:                test.fields.ip,
				ngt:               test.fields.ngt,
				eg:                test.fields.eg,
				streamConcurrency: test.fields.streamConcurrency,
			}

			gotRes, err := s.MultiUpsert(test.args.ctx, test.args.reqs)
			if err := checkFunc(test.want, gotRes, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
