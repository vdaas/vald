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

// Package service
package service

import (
	"context"
	"strings"
	"testing"

	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/client/v1/client/discoverer"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/sync"
)

// rangeClient models gRPCClient.RangeConcurrent: f runs per target, up to
// attempts times while it fails (the client backoff), and every per-target
// error that is not a local ctx error is dropped from the return value.
type rangeClient struct {
	grpc.Client
	// onRange, when set, runs when the broadcast starts, to change the agent list
	// that the discoverer reports while the write is in flight.
	onRange  func()
	addrs    []string
	attempts int
}

func (c *rangeClient) RangeConcurrent(
	ctx context.Context,
	_ int,
	f func(ctx context.Context, addr string, conn *grpc.ClientConn, copts ...grpc.CallOption) error,
) error {
	if c.onRange != nil {
		c.onRange()
	}
	if len(c.addrs) == 0 {
		return errors.ErrGRPCClientConnNotFound("*")
	}
	for _, addr := range c.addrs {
		for range max(c.attempts, 1) {
			if f(ctx, addr, nil) == nil {
				break
			}
		}
	}
	return nil
}

type rangeDiscoverer struct {
	discoverer.Client
	c grpc.Client
	// addrs are the agents the discoverer knows, reached or not.
	addrs []string
}

func (d *rangeDiscoverer) GetClient() grpc.Client { return d.c }

func (d *rangeDiscoverer) GetReadClient() grpc.Client { return d.c }

func (d *rangeDiscoverer) GetAddrs(context.Context) []string { return d.addrs }

// BroadCast returned whatever RangeConcurrent returned, so an agent failure
// (Unavailable/Internal/Aborted/...) reached the lb handlers as nil and they
// answered NotFound or OK instead of the agent's code.
func Test_gateway_BroadCast_propagatesCallbackErrors(t *testing.T) {
	t.Parallel()
	unavailable := status.Error(codes.Unavailable, "agent unavailable")
	aborted := status.Error(codes.Aborted, "flush in progress")
	internal := status.Error(codes.Internal, "agent internal")
	const addrA, addrB = "a:8081", "b:8081"
	tests := []struct {
		// errs maps target to the error of each successive attempt; a missing
		// entry or an exhausted list means success.
		errs     map[string][]error
		name     string
		addrs    []string
		wantIs   []error
		attempts int
		wantCode codes.Code // codes.OK: nil error expected
	}{
		{
			name:  "all agents succeed",
			addrs: []string{addrA, addrB},
		},
		{
			name:     "one agent unavailable",
			addrs:    []string{addrA, addrB},
			errs:     map[string][]error{addrB: {unavailable}},
			wantCode: codes.Unavailable,
			wantIs:   []error{unavailable},
		},
		{
			name:     "all agents aborted",
			addrs:    []string{addrA, addrB},
			errs:     map[string][]error{addrA: {aborted}, addrB: {aborted}},
			wantCode: codes.Aborted,
			wantIs:   []error{aborted},
		},
		{
			name:     "mixed codes are all kept and the first target's code leads",
			addrs:    []string{addrB, addrA},
			errs:     map[string][]error{addrA: {internal}, addrB: {aborted}},
			wantCode: codes.Internal,
			wantIs:   []error{internal, aborted},
		},
		{
			name:     "a retry that succeeds clears the earlier error",
			addrs:    []string{addrA},
			attempts: 2,
			errs:     map[string][]error{addrA: {unavailable}},
		},
		{
			name:     "a retry that fails again keeps the last error",
			addrs:    []string{addrA},
			attempts: 2,
			errs:     map[string][]error{addrA: {unavailable, internal}},
			wantCode: codes.Internal,
			wantIs:   []error{internal},
		},
	}
	for _, test := range tests {
		for _, kind := range []BroadCastKind{READ, WRITE} {
			t.Run(test.name, func(t *testing.T) {
				t.Parallel()
				g := &gateway{client: &rangeDiscoverer{c: &rangeClient{addrs: test.addrs, attempts: test.attempts}}}
				calls := make(map[string]int) // the model calls f sequentially
				err := g.BroadCast(t.Context(), kind, func(_ context.Context, target string, _ vald.Client, _ ...grpc.CallOption) error {
					n := calls[target]
					calls[target]++
					if errs := test.errs[target]; n < len(errs) {
						return errs[n]
					}
					return nil
				})
				if test.wantCode == codes.OK {
					if err != nil {
						t.Fatalf("unexpected error: %v", err)
					}
					return
				}
				st, ok := status.FromError(err)
				if !ok || st == nil || st.Code() != test.wantCode {
					t.Fatalf("got error: %v, want code %v", err, test.wantCode)
				}
				for _, want := range test.wantIs {
					if !errors.Is(err, want) {
						t.Errorf("error %v does not contain %v", err, want)
					}
				}
			})
		}
	}
}

// No connection at all is still reported as ErrGRPCClientConnNotFound, which the
// handlers map to Internal.
func Test_gateway_BroadCast_connNotFound(t *testing.T) {
	t.Parallel()
	g := &gateway{client: &rangeDiscoverer{c: new(rangeClient)}}
	err := g.BroadCast(t.Context(), WRITE, func(context.Context, string, vald.Client, ...grpc.CallOption) error {
		t.Error("callback must not run without connections")
		return nil
	})
	if !errors.Is(err, errors.ErrGRPCClientConnNotFound("*")) {
		t.Fatalf("got %v, want connection not found", err)
	}
}

// RangeConcurrent returns nil for a target whose connection could not be acquired
// (or whose circuit breaker is open): the callback never runs and the target's
// error is dropped. A write BroadCast must still report such a target, otherwise
// Remove/Flush look successful while that agent keeps the object. The agents that
// are reachable must be written regardless.
func Test_gateway_BroadCast_unreachedTargets(t *testing.T) {
	t.Parallel()
	const addrA, addrB, addrC = "a:8081", "b:8081", "c:8081"
	known := []string{addrA, addrB, addrC}
	tests := []struct {
		cbErr     map[string]error
		name      string
		reached   []string // targets RangeConcurrent runs the callback for
		wantInErr []string // substrings that the returned error must contain
		kind      BroadCastKind
		wantCode  codes.Code
	}{
		{
			name:     "write: an unreached agent is reported as Unavailable",
			kind:     WRITE,
			reached:  []string{addrA, addrC},
			wantCode: codes.Unavailable, wantInErr: []string{addrB},
		},
		{
			name:     "write: callback failures and unreached agents are both reported",
			kind:     WRITE,
			reached:  []string{addrA, addrB},
			cbErr:    map[string]error{addrB: status.Error(codes.Internal, "agent internal")},
			wantCode: codes.Internal, wantInErr: []string{"agent internal", addrC},
		},
		{
			name:    "write: all agents reached",
			kind:    WRITE,
			reached: known,
		},
		{
			name:    "read: unreached agents keep being ignored (partial search results)",
			kind:    READ,
			reached: []string{addrA},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			g := &gateway{client: &rangeDiscoverer{c: &rangeClient{addrs: test.reached}, addrs: known}}
			var mu sync.Mutex
			called := map[string]bool{}
			err := g.BroadCast(t.Context(), test.kind, func(_ context.Context, target string, _ vald.Client, _ ...grpc.CallOption) error {
				mu.Lock()
				called[target] = true
				mu.Unlock()
				return test.cbErr[target]
			})
			for _, addr := range test.reached {
				if !called[addr] {
					t.Errorf("reachable agent %s was not written", addr)
				}
			}
			if test.wantCode == codes.OK {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			if err == nil {
				t.Fatal("got nil error, want the unreached agent reported")
			}
			if st, ok := status.FromError(err); !ok || st == nil || st.Code() != test.wantCode {
				t.Errorf("got error %v, want code %v", err, test.wantCode)
			}
			for _, sub := range test.wantInErr {
				if !strings.Contains(err.Error(), sub) {
					t.Errorf("error %q does not mention %q", err, sub)
				}
			}
		})
	}
}

// The agents a write has to reach are the ones known when it starts. An agent that
// falls out of the discoverer's list during the broadcast and was not reached must
// still be reported: it can rejoin with the object the write was meant to remove.
// An agent that joins the list during the broadcast is not a target of this write.
func Test_gateway_BroadCast_agentListChangesDuringBroadcast(t *testing.T) {
	t.Parallel()
	const addrA, addrB, addrC = "a:8081", "b:8081", "c:8081"
	tests := []struct {
		name      string
		wantInErr string   // empty: nil error expected
		before    []string // the discoverer's list when the broadcast starts
		after     []string // the list after the broadcast has started
		reached   []string // targets RangeConcurrent runs the callback for
	}{
		{
			name:      "an unreached agent that leaves the list is still reported",
			before:    []string{addrA, addrB},
			after:     []string{addrA},
			reached:   []string{addrA},
			wantInErr: addrB,
		},
		{
			name:    "an unreached agent that joins the list is not a target of this write",
			before:  []string{addrA},
			after:   []string{addrA, addrC},
			reached: []string{addrA},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			d := &rangeDiscoverer{addrs: test.before}
			d.c = &rangeClient{
				addrs:   test.reached,
				onRange: func() { d.addrs = test.after },
			}
			g := &gateway{client: d}
			err := g.BroadCast(t.Context(), WRITE, func(context.Context, string, vald.Client, ...grpc.CallOption) error { return nil })
			if test.wantInErr == "" {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			if err == nil {
				t.Fatal("got nil error, want the agent that left the list reported as Unavailable")
			}
			if st, ok := status.FromError(err); !ok || st == nil || st.Code() != codes.Unavailable {
				t.Errorf("got error %v, want code %v", err, codes.Unavailable)
			}
			if !strings.Contains(err.Error(), test.wantInErr) {
				t.Errorf("error %q does not mention %q", err, test.wantInErr)
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
//
// func TestNewGateway(t *testing.T) {
// 	type args struct {
// 		opts []Option
// 	}
// 	type want struct {
// 		wantGw Gateway
// 		err    error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, Gateway, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotGw Gateway, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotGw, w.wantGw) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotGw, w.wantGw)
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
// 			gotGw, err := NewGateway(test.args.opts...)
// 			if err := checkFunc(test.want, gotGw, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_gateway_Start(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		client discoverer.Client
// 		eg     errgroup.Group
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
// 		           eg:nil,
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
// 		           eg:nil,
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
// 			g := &gateway{
// 				client: test.fields.client,
// 				eg:     test.fields.eg,
// 			}
//
// 			got, err := g.Start(test.args.ctx)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_gateway_BroadCast(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		kind BroadCastKind
// 		f    func(ctx context.Context, target string, ac vald.Client, copts ...grpc.CallOption) error
// 	}
// 	type fields struct {
// 		client discoverer.Client
// 		eg     errgroup.Group
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
// 		           kind:nil,
// 		           f:nil,
// 		       },
// 		       fields: fields {
// 		           client:nil,
// 		           eg:nil,
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
// 		           kind:nil,
// 		           f:nil,
// 		           },
// 		           fields: fields {
// 		           client:nil,
// 		           eg:nil,
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
// 			g := &gateway{
// 				client: test.fields.client,
// 				eg:     test.fields.eg,
// 			}
//
// 			err := g.BroadCast(test.args.ctx, test.args.kind, test.args.f)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_gateway_DoMulti(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		num int
// 		f   func(ctx context.Context, target string, ac vald.Client, copts ...grpc.CallOption) error
// 	}
// 	type fields struct {
// 		client discoverer.Client
// 		eg     errgroup.Group
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
// 		           num:0,
// 		           f:nil,
// 		       },
// 		       fields: fields {
// 		           client:nil,
// 		           eg:nil,
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
// 		           num:0,
// 		           f:nil,
// 		           },
// 		           fields: fields {
// 		           client:nil,
// 		           eg:nil,
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
// 			g := &gateway{
// 				client: test.fields.client,
// 				eg:     test.fields.eg,
// 			}
//
// 			err := g.DoMulti(test.args.ctx, test.args.num, test.args.f)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_gateway_GetAgentCount(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		client discoverer.Client
// 		eg     errgroup.Group
// 	}
// 	type want struct {
// 		want int
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, int) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got int) error {
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
// 		           eg:nil,
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
// 		           eg:nil,
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
// 			g := &gateway{
// 				client: test.fields.client,
// 				eg:     test.fields.eg,
// 			}
//
// 			got := g.GetAgentCount(test.args.ctx)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_gateway_Addrs(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		client discoverer.Client
// 		eg     errgroup.Group
// 	}
// 	type want struct {
// 		want []string
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, []string) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got []string) error {
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
// 		           eg:nil,
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
// 		           eg:nil,
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
// 			g := &gateway{
// 				client: test.fields.client,
// 				eg:     test.fields.eg,
// 			}
//
// 			got := g.Addrs(test.args.ctx)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
