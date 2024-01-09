// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package service

import (
	"context"
	"testing"

	agent "github.com/vdaas/vald/apis/grpc/v1/agent/core"
	"github.com/vdaas/vald/internal/client/v1/client/discoverer"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/test/goleak"
	clientmock "github.com/vdaas/vald/internal/test/mock/client"
	grpcmock "github.com/vdaas/vald/internal/test/mock/grpc"
)

func Test_index_Start(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type fields struct {
		client         discoverer.Client
		targetAddrs    []string
		targetAddrList map[string]bool
		concurrency    int
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
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func(w want, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		return nil
	}
	tests := []test{
		func() test {
			addrs := []string{
				"127.0.0.1:8080",
			}
			return test{
				name: "Success: when there is no error in the save indexing request process",
				args: args{
					ctx: context.Background(),
				},

				fields: fields{
					client: &clientmock.DiscovererClientMock{
						GetAddrsFunc: func(_ context.Context) []string {
							return addrs
						},
						GetClientFunc: func() grpc.Client {
							return &grpcmock.GRPCClientMock{
								OrderedRangeConcurrentFunc: func(_ context.Context, _ []string, _ int,
									_ func(_ context.Context, _ string, _ *grpc.ClientConn, _ ...grpc.CallOption) error,
								) error {
									return nil
								},
							}
						},
					},
				},
			}
		}(),
		func() test {
			addrs := []string{
				"127.0.0.1:8080",
			}
			return test{
				name: "Fail: when there is an error wrapped with gRPC status in the save indexing request process",
				args: args{
					ctx: context.Background(),
				},
				fields: fields{
					client: &clientmock.DiscovererClientMock{
						GetAddrsFunc: func(_ context.Context) []string {
							return addrs
						},
						GetClientFunc: func() grpc.Client {
							return &grpcmock.GRPCClientMock{
								OrderedRangeConcurrentFunc: func(_ context.Context, _ []string, _ int,
									_ func(_ context.Context, _ string, _ *grpc.ClientConn, _ ...grpc.CallOption) error,
								) error {
									return status.WrapWithInternal(
										agent.SaveIndexRPCName+" API connection not found",
										errors.ErrGRPCClientConnNotFound("*"),
									)
								},
							}
						},
					},
				},
				want: want{
					err: status.Error(codes.Internal,
						agent.SaveIndexRPCName+" API connection not found"),
				},
			}
		}(),
		func() test {
			addrs := []string{
				"127.0.0.1:8080",
			}
			return test{
				name: "Fail: When the OrderedRangeConcurrent method returns a gRPC client conn not found error",
				args: args{
					ctx: context.Background(),
				},

				fields: fields{
					client: &clientmock.DiscovererClientMock{
						GetAddrsFunc: func(_ context.Context) []string {
							return addrs
						},
						GetClientFunc: func() grpc.Client {
							return &grpcmock.GRPCClientMock{
								OrderedRangeConcurrentFunc: func(_ context.Context, _ []string, _ int,
									_ func(_ context.Context, _ string, _ *grpc.ClientConn, _ ...grpc.CallOption) error,
								) error {
									return errors.ErrGRPCClientConnNotFound("*")
								},
							}
						},
					},
				},
				want: want{
					err: status.Error(codes.Internal,
						agent.SaveIndexRPCName+" API connection not found"),
				},
			}
		}(),
		func() test {
			targetAddrs := []string{
				"127.0.0.1:8080",
			}
			targetAddrList := map[string]bool{
				targetAddrs[0]: true,
			}
			return test{
				name: "Fail: when there is no address matching targetAddrList",
				args: args{
					ctx: context.Background(),
				},
				fields: fields{
					client: &clientmock.DiscovererClientMock{
						GetAddrsFunc: func(_ context.Context) []string {
							// NOTE: This function returns nil, meaning that the targetAddrs stored in the field are invalid values.
							return nil
						},
					},
					targetAddrs:    targetAddrs,
					targetAddrList: targetAddrList,
				},
				want: want{
					err: status.Error(codes.Internal,
						agent.SaveIndexRPCName+" API connection target address \"127.0.0.1:8080\" not found"),
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
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
			idx := &index{
				client:         test.fields.client,
				targetAddrs:    test.fields.targetAddrs,
				targetAddrList: test.fields.targetAddrList,
				concurrency:    test.fields.concurrency,
			}

			err := idx.Start(test.args.ctx)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
//
// func TestNew(t *testing.T) {
// 	type args struct {
// 		opts []Option
// 	}
// 	type want struct {
// 		want Indexer
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, Indexer, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got Indexer, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
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
// 		           opts:nil,
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
// 		           opts:nil,
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
// 			got, err := New(test.args.opts...)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
//
// func Test_index_StartClient(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		client         discoverer.Client
// 		targetAddrs    []string
// 		targetAddrList map[string]bool
// 		concurrency    int
// 	}
// 	type want struct {
// 		want <-chan error
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, <-chan error, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got <-chan error, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
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
// 		           ctx:nil,
// 		       },
// 		       fields: fields {
// 		           client:nil,
// 		           targetAddrs:nil,
// 		           targetAddrList:nil,
// 		           concurrency:0,
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
// 		           ctx:nil,
// 		           },
// 		           fields: fields {
// 		           client:nil,
// 		           targetAddrs:nil,
// 		           targetAddrList:nil,
// 		           concurrency:0,
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
// 			idx := &index{
// 				client:         test.fields.client,
// 				targetAddrs:    test.fields.targetAddrs,
// 				targetAddrList: test.fields.targetAddrList,
// 				concurrency:    test.fields.concurrency,
// 			}
//
// 			got, err := idx.StartClient(test.args.ctx)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
//
// 		})
// 	}
// }
