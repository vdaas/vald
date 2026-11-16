// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
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
	"slices"
	"testing"

	agent "github.com/vdaas/vald/apis/grpc/v1/agent/core"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/client/v1/client/discoverer"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/pool"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/internal/test/goleak"
	clientmock "github.com/vdaas/vald/internal/test/mock/client"
	grpcmock "github.com/vdaas/vald/internal/test/mock/grpc"
)

const testAgentAddr = "127.0.0.1:8080"

func Test_index_Start(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context //nolint:containedctx // the args struct carrying ctx is the repo-wide table-driven test pattern
	}
	type fields struct {
		client      discoverer.Client
		targetAddrs []string
		concurrency int
	}
	type want struct {
		err error
	}
	type test struct {
		args       args
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
		name       string
		fields     fields
	}
	wantAnyError := func(_ want, err error) error {
		if err == nil {
			return errors.New("got_error: nil, want: an error for the unprocessed target agents")
		}
		return nil
	}
	defaultCheckFunc := func(w want, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%s\",\n\t\t\t\twant: \"%s\"", err.Error(), w.err.Error())
		}
		return nil
	}
	tests := []test{
		func() test {
			addrs := []string{
				testAgentAddr,
			}
			return test{
				name: "Fail: when OrderedRangeConcurrent never runs the save indexing request for the target agent",
				// the callback never ran, so no agent was processed and Start must fail.
				checkFunc: wantAnyError,
				args: args{
					ctx: t.Context(),
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
			return test{
				name: "Fail: when a target addresses (targetAddrs) is given and OrderedRangeConcurrent never runs the save indexing request for them",
				// the callback never ran, so no agent was processed and Start must fail.
				checkFunc: wantAnyError,
				args: args{
					ctx: t.Context(),
				},
				fields: fields{
					targetAddrs: []string{
						testAgentAddr,
						"127.0.0.1:8081",
						"127.0.0.1:8083",
					},
					client: &clientmock.DiscovererClientMock{
						GetAddrsFunc: func(_ context.Context) []string {
							return nil
						},
						GetClientFunc: func() grpc.Client {
							return &grpcmock.GRPCClientMock{
								OrderedRangeConcurrentFunc: func(_ context.Context, _ []string, _ int,
									_ func(_ context.Context, _ string, _ *grpc.ClientConn, _ ...grpc.CallOption) error,
								) error {
									return nil
								},
								ConnectFunc: func(_ context.Context, _ string, _ ...grpc.DialOption) (pool.Conn, error) {
									return nil, nil
								},
							}
						},
					},
				},
			}
		}(),
		func() test {
			addrs := []string{
				testAgentAddr,
			}
			return test{
				name: "Fail: when there is an error wrapped with gRPC status in the save indexing request process",
				args: args{
					ctx: t.Context(),
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
				testAgentAddr,
			}
			return test{
				name: "Fail: When the OrderedRangeConcurrent method returns a gRPC client conn not found error",
				args: args{
					ctx: t.Context(),
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
				client:      test.fields.client,
				targetAddrs: test.fields.targetAddrs,
				concurrency: test.fields.concurrency,
			}

			err := idx.Start(test.args.ctx)
			if cerr := checkFunc(test.want, err); cerr != nil {
				tt.Errorf("error = %v", cerr)
			}
		})
	}
}

// Test_index_doSaveIndex_retriedTarget pins that only the last attempt per target is
// reported: the client backoff re-runs the callback for the same target, and a
// retry that succeeds must clear the earlier failure.
func Test_index_doSaveIndex_retriedTarget(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		wantSubstr string
		notSubstr  string
		attempts   []error
		wantErr    bool
	}{
		{
			name:     "a failure followed by a successful retry is not an error",
			attempts: []error{status.Error(codes.Unavailable, "first failure"), nil},
		},
		{
			name:       "only the last failure of a target is reported",
			attempts:   []error{status.Error(codes.Unavailable, "first failure"), status.Error(codes.Internal, "final failure")},
			wantErr:    true,
			wantSubstr: "final failure",
			notSubstr:  "first failure",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(tt *testing.T) {
			tt.Parallel()
			idx := &index{
				client: &clientmock.DiscovererClientMock{
					GetAddrsFunc: func(_ context.Context) []string {
						return []string{testAgentAddr}
					},
					GetClientFunc: func() grpc.Client {
						return &grpcmock.GRPCClientMock{
							// re-run a failing callback like the client backoff does and swallow
							// the final error like the real OrderedRangeConcurrent.
							OrderedRangeConcurrentFunc: func(ctx context.Context, order []string, _ int,
								f func(_ context.Context, _ string, _ *grpc.ClientConn, _ ...grpc.CallOption) error,
							) error {
								for _, addr := range order {
									for range tc.attempts {
										if f(ctx, addr, nil) == nil {
											break
										}
									}
								}
								return nil
							},
						}
					},
				},
			}
			var calls int
			err := idx.doSaveIndex(tt.Context(),
				func(context.Context, agent.AgentClient, ...grpc.CallOption) (*payload.Empty, error) {
					if calls >= len(tc.attempts) {
						return nil, errors.New("unexpected extra attempt")
					}
					aerr := tc.attempts[calls]
					calls++
					return nil, aerr
				})
			if !tc.wantErr {
				if err != nil {
					tt.Errorf("doSaveIndex() error = %v, want nil", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tc.wantSubstr) {
				tt.Fatalf("doSaveIndex() error = %v, want it to contain %q", err, tc.wantSubstr)
			}
			if strings.Contains(err.Error(), tc.notSubstr) {
				tt.Errorf("doSaveIndex() error = %v, must not contain the retried failure %q", err, tc.notSubstr)
			}
		})
	}
}

// Test_index_doSaveIndex_unprocessedTarget pins that a target agent the client never ran the
// callback for (no connection pool, connection lookup failure, open circuit breaker, done
// ctx) fails the job: OrderedRangeConcurrent swallows those errors and returns nil.
func Test_index_doSaveIndex_unprocessedTarget(t *testing.T) {
	t.Parallel()
	targets := []string{"127.0.0.1:8080", "127.0.0.1:8081", "127.0.0.1:8082"}
	tests := []struct {
		name      string
		result    error
		reached   []string
		unreached []string
	}{
		{
			name:    "every target processed is not an error",
			reached: targets,
		},
		{
			name:      "a target the callback never ran for is an error",
			reached:   []string{targets[0], targets[2]},
			unreached: []string{targets[1]},
		},
		{
			name:      "no target reached is an error for every target",
			unreached: targets,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(tt *testing.T) {
			tt.Parallel()
			idx := &index{
				client: &clientmock.DiscovererClientMock{
					GetAddrsFunc: func(_ context.Context) []string {
						return targets
					},
					GetClientFunc: func() grpc.Client {
						return &grpcmock.GRPCClientMock{
							// skip the unreached targets and swallow that like the real OrderedRangeConcurrent.
							OrderedRangeConcurrentFunc: func(ctx context.Context, order []string, _ int,
								f func(_ context.Context, _ string, _ *grpc.ClientConn, _ ...grpc.CallOption) error,
							) error {
								for _, addr := range order {
									if slices.Contains(tc.reached, addr) {
										if err := f(ctx, addr, nil); err != nil {
											tt.Logf("callback for %s returned %v", addr, err)
										}
									}
								}
								return nil
							},
						}
					},
				},
			}
			err := idx.doSaveIndex(tt.Context(),
				func(context.Context, agent.AgentClient, ...grpc.CallOption) (*payload.Empty, error) {
					return nil, tc.result
				})
			if len(tc.unreached) == 0 {
				if err != nil {
					tt.Errorf("doSaveIndex() error = %v, want nil", err)
				}
				return
			}
			if err == nil {
				tt.Fatalf("doSaveIndex() error = nil, want an error for %v", tc.unreached)
			}
			for _, addr := range targets {
				want := slices.Contains(tc.unreached, addr)
				if got := errors.Is(err, errors.ErrGRPCClientConnNotFound(addr)); got != want {
					tt.Errorf("doSaveIndex() error = %v, reports %s as unprocessed: %t, want %t", err, addr, got, want)
				}
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
// 		})
// 	}
// }
//
// func Test_index_StartClient(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		client      discoverer.Client
// 		targetAddrs []string
// 		concurrency int
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
// 				client:      test.fields.client,
// 				targetAddrs: test.fields.targetAddrs,
// 				concurrency: test.fields.concurrency,
// 			}
//
// 			got, err := idx.StartClient(test.args.ctx)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_index_doSaveIndex(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		fn  func(_ context.Context, _ agent.AgentClient, _ ...grpc.CallOption) (*payload.Empty, error)
// 	}
// 	type fields struct {
// 		client      discoverer.Client
// 		targetAddrs []string
// 		concurrency int
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
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
// 		           ctx:nil,
// 		           fn:nil,
// 		       },
// 		       fields: fields {
// 		           client:nil,
// 		           targetAddrs:nil,
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
// 		           fn:nil,
// 		           },
// 		           fields: fields {
// 		           client:nil,
// 		           targetAddrs:nil,
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
// 				client:      test.fields.client,
// 				targetAddrs: test.fields.targetAddrs,
// 				concurrency: test.fields.concurrency,
// 			}
//
// 			err := idx.doSaveIndex(test.args.ctx, test.args.fn)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
