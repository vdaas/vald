//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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

// Package grpc provides generic functionality for grpc
package grpc

import (
	"context"
	"reflect"
	"testing"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/io"
	"github.com/vdaas/vald/internal/test/data/request"
	"github.com/vdaas/vald/internal/test/data/vector"
	"github.com/vdaas/vald/internal/test/goleak"
	"github.com/vdaas/vald/internal/test/mock"
)

func TestBidirectionalStream(t *testing.T) {
	type args struct {
		concurrency int
		f           func(context.Context, *payload.Insert_Request) (*payload.Object_StreamLocation, error)
		insertReqs  []*payload.Insert_Request
	}
	type want struct {
		rpcResp []*payload.Object_StreamLocation
		err     error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func([]*payload.Object_StreamLocation, want, error) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
	}

	const (
		name        = "vald-agent-ngt-1" // agent name
		f32VecDim   = 3                  // float32 vector dimension
		ip          = "localhost"        // ip address
		concurrency = 10
	)

	defaultCallbackFn := func(ctx context.Context, i *payload.Insert_Request) (*payload.Object_StreamLocation, error) {
		return &payload.Object_StreamLocation{
			Payload: &payload.Object_StreamLocation_Location{
				Location: &payload.Object_Location{
					Name: name,
					Uuid: i.GetVector().GetId(),
					Ips:  []string{ip},
				},
			},
		}, nil
	}

	defaultCheckFunc := func(rpcResp []*payload.Object_StreamLocation, w want, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}

		// since the insert order is not guaranteed, check only the error count on the response
		sm := make(map[int32]int) // want status map
		for _, r := range w.rpcResp {
			sm[r.GetStatus().GetCode()]++
		}
		gsm := make(map[int32]int) // got status map
		for _, r := range rpcResp {
			gsm[r.GetStatus().GetCode()]++
		}

		if !reflect.DeepEqual(sm, gsm) {
			return errors.Errorf("status count is not correct, got: %v, want: %v", gsm, sm)
		}
		return nil
	}
	tests := []test{
		func() test {
			insertCnt := 1
			reqs, err := request.GenMultiInsertReq(request.Float, vector.Gaussian, insertCnt, f32VecDim, nil)
			if err != nil {
				t.Fatal(err)
			}
			return test{
				name: "success to receive 1 message from stream",
				args: args{
					insertReqs:  reqs.Requests,
					concurrency: concurrency,
					f:           defaultCallbackFn,
				},
				want: want{
					rpcResp: request.GenObjectStreamLocation(insertCnt, name, ip),
				},
			}
		}(),
		func() test {
			insertCnt := 10
			reqs, err := request.GenMultiInsertReq(request.Float, vector.Gaussian, insertCnt, f32VecDim, nil)
			if err != nil {
				t.Fatal(err)
			}
			return test{
				name: "success to receive 10 message from stream",
				args: args{
					insertReqs:  reqs.Requests,
					concurrency: concurrency,
					f:           defaultCallbackFn,
				},
				want: want{
					rpcResp: request.GenObjectStreamLocation(insertCnt, name, ip),
				},
			}
		}(),
		func() test {
			return test{
				name: "success to receive 0 message from stream",
				args: args{
					insertReqs:  []*payload.Insert_Request{},
					concurrency: concurrency,
					f:           defaultCallbackFn,
				},
				want: want{
					rpcResp: nil,
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			if test.beforeFunc != nil {
				test.beforeFunc(tt, test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			insertReqs := test.args.insertReqs
			rpcResp := make([]*payload.Object_StreamLocation, 0)
			recvIdx := 0

			stream := &mock.StreamInsertServerMock{
				ServerStream: &mock.ServerStreamMock{
					ContextFunc: func() context.Context {
						return ctx
					},
					RecvMsgFunc: func(i interface{}) error {
						if recvIdx >= len(insertReqs) {
							return io.EOF
						}

						obj := i.(*payload.Insert_Request)
						if insertReqs[recvIdx] != nil {
							obj.Vector = insertReqs[recvIdx].Vector
							obj.Config = insertReqs[recvIdx].Config
						}
						recvIdx++

						return nil
					},
					SendMsgFunc: func(i interface{}) error {
						rpcResp = append(rpcResp, i.(*payload.Object_StreamLocation))
						return nil
					},
				},
			}

			err := BidirectionalStream(ctx, stream, test.args.concurrency, test.args.f)
			if err := checkFunc(rpcResp, test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
//
// func TestBidirectionalStreamClient(t *testing.T) {
// 	type args struct {
// 		stream       ClientStream
// 		dataProvider func() interface{}
// 		newData      func() interface{}
// 		f            func(interface{}, error)
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           stream:nil,
// 		           dataProvider:nil,
// 		           newData:nil,
// 		           f:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
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
// 		           args: args {
// 		           stream:nil,
// 		           dataProvider:nil,
// 		           newData:nil,
// 		           f:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
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
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			err := BidirectionalStreamClient(test.args.stream, test.args.dataProvider, test.args.newData, test.args.f)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_removeDuplicates(t *testing.T) {
// 	type args struct {
// 		x    S
// 		less func(left, right E) int
// 	}
// 	type want struct {
// 		want S
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, S) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got S) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           x:nil,
// 		           less:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
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
// 		           args: args {
// 		           x:nil,
// 		           less:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
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
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			got := removeDuplicates(test.args.x, test.args.less)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
