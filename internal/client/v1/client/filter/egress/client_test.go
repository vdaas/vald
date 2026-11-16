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

package egress

import (
	"context"
	"slices"
	"testing"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/filter/egress"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/pool"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/sync"
	grpcmock "github.com/vdaas/vald/internal/test/mock/grpc"
	ggrpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

// fakeFilter adds tag to the distance or appends it to the vector, or fails
// with err when set.
type fakeFilter struct {
	egress.UnimplementedFilterServer
	err    error
	mu     sync.Mutex
	called int
	tag    float32
}

func (f *fakeFilter) call() error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.called++
	return f.err
}

func (f *fakeFilter) calls() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.called
}

func (f *fakeFilter) FilterDistance(
	_ context.Context, in *payload.Object_Distance,
) (*payload.Object_Distance, error) {
	if err := f.call(); err != nil {
		return nil, err
	}
	return &payload.Object_Distance{Id: in.GetId(), Distance: in.GetDistance() + f.tag}, nil
}

func (f *fakeFilter) FilterVector(
	_ context.Context, in *payload.Object_Vector,
) (*payload.Object_Vector, error) {
	if err := f.call(); err != nil {
		return nil, err
	}
	return &payload.Object_Vector{Id: in.GetId(), Vector: append(slices.Clone(in.GetVector()), f.tag)}, nil
}

// fakeGRPCClient runs callbacks on bufconn connections. OrderedRange mirrors
// gRPCClient.OrderedRange, which drops per-address errors.
type fakeGRPCClient struct {
	grpc.Client
	conns map[string]*grpc.ClientConn
}

func (f *fakeGRPCClient) Do(
	ctx context.Context,
	addr string,
	fn func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (any, error),
) (any, error) {
	conn, ok := f.conns[addr]
	if !ok {
		return nil, errors.ErrGRPCClientConnNotFound(addr)
	}
	return fn(ctx, conn)
}

func (f *fakeGRPCClient) OrderedRange(
	ctx context.Context,
	order []string,
	fn func(ctx context.Context, addr string, conn *grpc.ClientConn, copts ...grpc.CallOption) error,
) error {
	for _, addr := range order {
		if conn, ok := f.conns[addr]; ok {
			if err := fn(ctx, addr, conn); err != nil {
				// gRPCClient.OrderedRange joins the error and still returns nil.
				continue
			}
		}
	}
	return nil
}

func newFakeGRPCClient(
	t *testing.T, order []string, filters map[string]*fakeFilter,
) *fakeGRPCClient {
	t.Helper()
	conns := make(map[string]*grpc.ClientConn, len(order))
	for _, addr := range order {
		lis := bufconn.Listen(1 << 20)
		srv := ggrpc.NewServer()
		egress.RegisterFilterServer(srv, filters[addr])
		served := make(chan struct{})
		go func() {
			defer close(served)
			if err := srv.Serve(lis); err != nil {
				t.Logf("fake filter %s stopped: %v", addr, err)
			}
		}()
		conn, err := ggrpc.NewClient("passthrough:///"+addr,
			ggrpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) {
				return lis.DialContext(ctx)
			}),
			ggrpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			t.Fatalf("failed to create client for %s: %v", addr, err)
		}
		t.Cleanup(func() {
			if err := conn.Close(); err != nil {
				t.Logf("failed to close conn to %s: %v", addr, err)
			}
			srv.Stop()
			select {
			case <-served:
			case <-time.After(10 * time.Second):
				t.Errorf("fake filter %s did not stop", addr)
			}
		})
		conns[addr] = conn
	}
	return &fakeGRPCClient{conns: conns}
}

func Test_multipleAddrsClient_Chain(t *testing.T) {
	t.Parallel()
	failure := status.Error(codes.InvalidArgument, "filter failed")
	tests := []struct {
		errs      map[string]error
		wantCalls map[string]int
		name      string
		wantErr   bool
	}{
		{
			name:      "applies every filter in order",
			wantCalls: map[string]int{"a": 1, "b": 1},
		},
		{
			name:      "returns the error of a failing middle filter and stops the chain",
			errs:      map[string]error{"a": failure},
			wantCalls: map[string]int{"a": 1, "b": 0},
			wantErr:   true,
		},
		{
			name:      "returns the error of a failing last filter instead of (nil, nil)",
			errs:      map[string]error{"b": failure},
			wantCalls: map[string]int{"a": 1, "b": 1},
			wantErr:   true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			order := []string{"a", "b"}
			newClient := func() (*multipleAddrsClient, map[string]*fakeFilter) {
				filters := map[string]*fakeFilter{
					"a": {tag: 1, err: tc.errs["a"]},
					"b": {tag: 2, err: tc.errs["b"]},
				}
				return &multipleAddrsClient{addrs: order, c: newFakeGRPCClient(t, order, filters)}, filters
			}
			checkCalls := func(t *testing.T, filters map[string]*fakeFilter) {
				t.Helper()
				for addr, want := range tc.wantCalls {
					if c := filters[addr].calls(); c != want {
						t.Errorf("filter %s called %d times, want %d", addr, c, want)
					}
				}
			}
			ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
			defer cancel()

			m, filters := newClient()
			dist, err := m.FilterDistance(ctx, &payload.Object_Distance{Id: "x", Distance: 0.5})
			if tc.wantErr {
				if err == nil {
					t.Errorf("FilterDistance() error = nil, want error (got %v)", dist)
				} else if dist != nil {
					t.Errorf("FilterDistance() got = %v, want nil on error", dist)
				}
			} else if err != nil {
				t.Errorf("FilterDistance() error = %v, want nil", err)
			} else if got := dist.GetDistance(); got != 3.5 {
				t.Errorf("FilterDistance() distance = %v, want 3.5", got)
			}
			checkCalls(t, filters)

			m, filters = newClient()
			vec, err := m.FilterVector(ctx, &payload.Object_Vector{Id: "x", Vector: []float32{0}})
			if tc.wantErr {
				if err == nil {
					t.Errorf("FilterVector() error = nil, want error (got %v)", vec)
				} else if vec != nil {
					t.Errorf("FilterVector() got = %v, want nil on error", vec)
				}
			} else if err != nil {
				t.Errorf("FilterVector() error = %v, want nil", err)
			} else if want := []float32{0, 1, 2}; !slices.Equal(vec.GetVector(), want) {
				t.Errorf("FilterVector() vector = %v, want %v", vec.GetVector(), want)
			}
			checkCalls(t, filters)
		})
	}
}

// connTracker is a fake gRPC client connection table.
type connTracker struct {
	conns       map[string]bool
	unreachable map[string]bool
	connects    map[string]int
	disconnects map[string]int
	mu          sync.Mutex
}

func newConnTracker(unreachable ...string) *connTracker {
	ct := &connTracker{
		conns:       make(map[string]bool),
		unreachable: make(map[string]bool),
		connects:    make(map[string]int),
		disconnects: make(map[string]int),
	}
	for _, addr := range unreachable {
		ct.unreachable[addr] = true
	}
	return ct
}

func (ct *connTracker) mock() *grpcmock.GRPCClientMock {
	return &grpcmock.GRPCClientMock{
		ConnectFunc: func(_ context.Context, addr string, _ ...grpc.DialOption) (pool.Conn, error) {
			ct.mu.Lock()
			defer ct.mu.Unlock()
			if ct.unreachable[addr] {
				return nil, errors.New("unreachable")
			}
			ct.connects[addr]++
			ct.conns[addr] = true
			return nil, nil
		},
		DisconnectFunc: func(_ context.Context, addr string) error {
			ct.mu.Lock()
			defer ct.mu.Unlock()
			ct.disconnects[addr]++
			if !ct.conns[addr] {
				return errors.ErrGRPCClientConnNotFound(addr)
			}
			delete(ct.conns, addr)
			return nil
		},
		IsConnectedFunc: func(_ context.Context, addr string) bool {
			ct.mu.Lock()
			defer ct.mu.Unlock()
			// a pool left behind for an addr that went down is not connected.
			return ct.conns[addr] && !ct.unreachable[addr]
		},
	}
}

func (ct *connTracker) state(addr string) (connected bool, disconnects int) {
	ct.mu.Lock()
	defer ct.mu.Unlock()
	return ct.conns[addr], ct.disconnects[addr]
}

// waitFor polls cond until it holds; it fails the test after a bounded time.
func waitFor(t *testing.T, what string, cond func() bool) {
	t.Helper()
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()
	tick := time.NewTicker(time.Millisecond)
	defer tick.Stop()
	for !cond() {
		select {
		case <-ctx.Done():
			t.Fatalf("timed out waiting for %s", what)
		case <-tick.C:
		}
	}
}

// Test_client_Target_Release pins that Target and Release go through the
// request-derived target pool with the configured idle timeout: a released
// request-derived target is disconnected after it, configured targets stay
// connected, and the targets acquired before a failed connect are released.
// The pool itself is tested in package dyntarget.
func Test_client_Target_Release(t *testing.T) {
	t.Parallel()
	const (
		configured = "configured:8081"
		dynamic    = "dynamic:8081"
		down       = "down:8081"
	)
	ct := newConnTracker(down)
	c, err := New(WithAddrs(configured), WithClient(ct.mock()), WithIdleTimeout("1ms"))
	if err != nil {
		t.Fatal(err)
	}
	cl, ok := c.(*client)
	if !ok {
		t.Fatalf("New returned %T, want *client", c)
	}
	defer cl.targets.Stop()
	ctx := t.Context()
	if _, err = c.Target(ctx); !errors.Is(err, errors.ErrTargetNotFound) {
		t.Errorf("Target() error = %v, want %v", err, errors.ErrTargetNotFound)
	}
	fc, err := c.Target(ctx, configured, dynamic)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := fc.(*multipleAddrsClient); !ok {
		t.Errorf("Target returned %T, want *multipleAddrsClient", fc)
	}
	c.Release(ctx, configured, dynamic)
	waitFor(t, "the released request-derived target to be disconnected", func() bool {
		connected, disconnects := ct.state(dynamic)
		return !connected && disconnects == 1
	})
	if connected, disconnects := ct.state(configured); !connected || disconnects != 0 {
		t.Errorf("configured connected = %v, Disconnect calls = %d, want true, 0", connected, disconnects)
	}
	if _, err = c.Target(ctx, dynamic, down); err == nil {
		t.Fatal("Target to an unreachable addr succeeded")
	}
	waitFor(t, "the target acquired before the failed connect to be disconnected", func() bool {
		connected, disconnects := ct.state(dynamic)
		return !connected && disconnects == 2
	})
}

// lifecycleClient adds the connection monitor and Close that the client's
// Start and Stop need to the mock.
type lifecycleClient struct {
	*grpcmock.GRPCClientMock
}

func (lifecycleClient) StartConnectionMonitor(context.Context) (<-chan error, error) {
	ech := make(chan error)
	close(ech)
	return ech, nil
}

func (lifecycleClient) Close(context.Context) error { return nil }

// releaseKeepsConnection reports whether a released request-derived target is
// still connected, with no further Disconnect, well after the idle timeout.
func releaseKeepsConnection(
	ctx context.Context, t *testing.T, c Client, ct *connTracker, addr string,
) bool {
	t.Helper()
	_, before := ct.state(addr)
	if _, err := c.Target(ctx, addr); err != nil {
		t.Fatal(err)
	}
	c.Release(ctx, addr)
	time.Sleep(50 * time.Millisecond) // 50x the 1ms idle timeout
	connected, after := ct.state(addr)
	return connected && after == before
}

// Test_client_StartStop_wireIdleDisconnects pins that the client's Start ties
// the request-derived target pool to its ctx and that Stop stops it: once
// either has happened, a released request-derived target is no longer
// disconnected after the idle timeout.
func Test_client_StartStop_wireIdleDisconnects(t *testing.T) {
	t.Parallel()
	const dynamic = "dynamic:8081"
	newClient := func(t *testing.T) (Client, *connTracker) {
		t.Helper()
		ct := newConnTracker()
		c, err := New(WithAddrs("configured:8081"),
			WithClient(lifecycleClient{ct.mock()}), WithIdleTimeout("1ms"))
		if err != nil {
			t.Fatal(err)
		}
		return c, ct
	}
	t.Run("Start ctx done", func(t *testing.T) {
		t.Parallel()
		c, ct := newClient(t)
		ctx, cancel := context.WithCancel(t.Context())
		if _, err := c.Start(ctx); err != nil {
			t.Fatal(err)
		}
		cancel()
		// Start stops the pool asynchronously (context.AfterFunc) once ctx is done.
		waitFor(t, "a released target to stay connected after the Start ctx is done", func() bool {
			return releaseKeepsConnection(t.Context(), t, c, ct, dynamic)
		})
	})
	t.Run("Stop", func(t *testing.T) {
		t.Parallel()
		c, ct := newClient(t)
		if err := c.Stop(t.Context()); err != nil {
			t.Fatal(err)
		}
		if !releaseKeepsConnection(t.Context(), t, c, ct, dynamic) {
			t.Error("a released target was disconnected after Stop")
		}
	})
}

// NOT IMPLEMENTED BELOW
//
// func TestNew(t *testing.T) {
// 	type args struct {
// 		opts []Option
// 	}
// 	type want struct {
// 		want Client
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, Client, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got Client, err error) error {
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
// func Test_client_Start(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		c     grpc.Client
// 		cl    sync.Map[string, any]
// 		addrs []string
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
// 		           c:nil,
// 		           cl:nil,
// 		           addrs:nil,
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
// 		           c:nil,
// 		           cl:nil,
// 		           addrs:nil,
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
// 			c := &client{
// 				c:     test.fields.c,
// 				cl:    test.fields.cl,
// 				addrs: test.fields.addrs,
// 			}
//
// 			got, err := c.Start(test.args.ctx)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_client_Stop(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		c     grpc.Client
// 		cl    sync.Map[string, any]
// 		addrs []string
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
// 		       },
// 		       fields: fields {
// 		           c:nil,
// 		           cl:nil,
// 		           addrs:nil,
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
// 		           c:nil,
// 		           cl:nil,
// 		           addrs:nil,
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
// 			c := &client{
// 				c:     test.fields.c,
// 				cl:    test.fields.cl,
// 				addrs: test.fields.addrs,
// 			}
//
// 			err := c.Stop(test.args.ctx)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_client_GRPCClient(t *testing.T) {
// 	type fields struct {
// 		c     grpc.Client
// 		cl    sync.Map[string, any]
// 		addrs []string
// 	}
// 	type want struct {
// 		want grpc.Client
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, grpc.Client) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got grpc.Client) error {
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
// 		       fields: fields {
// 		           c:nil,
// 		           cl:nil,
// 		           addrs:nil,
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
// 		           c:nil,
// 		           cl:nil,
// 		           addrs:nil,
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
// 			c := &client{
// 				c:     test.fields.c,
// 				cl:    test.fields.cl,
// 				addrs: test.fields.addrs,
// 			}
//
// 			got := c.GRPCClient()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_client_Target(t *testing.T) {
// 	type args struct {
// 		ctx     context.Context
// 		targets []string
// 	}
// 	type fields struct {
// 		c     grpc.Client
// 		cl    sync.Map[string, any]
// 		addrs []string
// 	}
// 	type want struct {
// 		want egress.FilterClient
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, egress.FilterClient, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got egress.FilterClient, err error) error {
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
// 		           targets:nil,
// 		       },
// 		       fields: fields {
// 		           c:nil,
// 		           cl:nil,
// 		           addrs:nil,
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
// 		           targets:nil,
// 		           },
// 		           fields: fields {
// 		           c:nil,
// 		           cl:nil,
// 		           addrs:nil,
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
// 			c := &client{
// 				c:     test.fields.c,
// 				cl:    test.fields.cl,
// 				addrs: test.fields.addrs,
// 			}
//
// 			got, err := c.Target(test.args.ctx, test.args.targets...)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_client_FilterDistance(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		in   *payload.Object_Distance
// 		opts []grpc.CallOption
// 	}
// 	type fields struct {
// 		c     grpc.Client
// 		cl    sync.Map[string, any]
// 		addrs []string
// 	}
// 	type want struct {
// 		wantRes *payload.Object_Distance
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Object_Distance, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotRes *payload.Object_Distance, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotRes, w.wantRes) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
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
// 		           in:nil,
// 		           opts:nil,
// 		       },
// 		       fields: fields {
// 		           c:nil,
// 		           cl:nil,
// 		           addrs:nil,
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
// 		           in:nil,
// 		           opts:nil,
// 		           },
// 		           fields: fields {
// 		           c:nil,
// 		           cl:nil,
// 		           addrs:nil,
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
// 			c := &client{
// 				c:     test.fields.c,
// 				cl:    test.fields.cl,
// 				addrs: test.fields.addrs,
// 			}
//
// 			gotRes, err := c.FilterDistance(test.args.ctx, test.args.in, test.args.opts...)
// 			if err := checkFunc(test.want, gotRes, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_specificAddrClient_FilterDistance(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		in   *payload.Object_Distance
// 		opts []grpc.CallOption
// 	}
// 	type fields struct {
// 		c    grpc.Client
// 		addr string
// 	}
// 	type want struct {
// 		wantRes *payload.Object_Distance
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Object_Distance, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotRes *payload.Object_Distance, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotRes, w.wantRes) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
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
// 		           in:nil,
// 		           opts:nil,
// 		       },
// 		       fields: fields {
// 		           c:nil,
// 		           addr:"",
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
// 		           in:nil,
// 		           opts:nil,
// 		           },
// 		           fields: fields {
// 		           c:nil,
// 		           addr:"",
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
// 			s := &specificAddrClient{
// 				c:    test.fields.c,
// 				addr: test.fields.addr,
// 			}
//
// 			gotRes, err := s.FilterDistance(test.args.ctx, test.args.in, test.args.opts...)
// 			if err := checkFunc(test.want, gotRes, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_multipleAddrsClient_FilterDistance(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		in   *payload.Object_Distance
// 		opts []grpc.CallOption
// 	}
// 	type fields struct {
// 		c     grpc.Client
// 		addrs []string
// 	}
// 	type want struct {
// 		wantRes *payload.Object_Distance
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Object_Distance, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotRes *payload.Object_Distance, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotRes, w.wantRes) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
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
// 		           in:nil,
// 		           opts:nil,
// 		       },
// 		       fields: fields {
// 		           c:nil,
// 		           addrs:nil,
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
// 		           in:nil,
// 		           opts:nil,
// 		           },
// 		           fields: fields {
// 		           c:nil,
// 		           addrs:nil,
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
// 			m := &multipleAddrsClient{
// 				c:     test.fields.c,
// 				addrs: test.fields.addrs,
// 			}
//
// 			gotRes, err := m.FilterDistance(test.args.ctx, test.args.in, test.args.opts...)
// 			if err := checkFunc(test.want, gotRes, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_client_FilterVector(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		in   *payload.Object_Vector
// 		opts []grpc.CallOption
// 	}
// 	type fields struct {
// 		c     grpc.Client
// 		cl    sync.Map[string, any]
// 		addrs []string
// 	}
// 	type want struct {
// 		wantRes *payload.Object_Vector
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Object_Vector, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotRes *payload.Object_Vector, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotRes, w.wantRes) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
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
// 		           in:nil,
// 		           opts:nil,
// 		       },
// 		       fields: fields {
// 		           c:nil,
// 		           cl:nil,
// 		           addrs:nil,
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
// 		           in:nil,
// 		           opts:nil,
// 		           },
// 		           fields: fields {
// 		           c:nil,
// 		           cl:nil,
// 		           addrs:nil,
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
// 			c := &client{
// 				c:     test.fields.c,
// 				cl:    test.fields.cl,
// 				addrs: test.fields.addrs,
// 			}
//
// 			gotRes, err := c.FilterVector(test.args.ctx, test.args.in, test.args.opts...)
// 			if err := checkFunc(test.want, gotRes, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_specificAddrClient_FilterVector(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		in   *payload.Object_Vector
// 		opts []grpc.CallOption
// 	}
// 	type fields struct {
// 		c    grpc.Client
// 		addr string
// 	}
// 	type want struct {
// 		wantRes *payload.Object_Vector
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Object_Vector, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotRes *payload.Object_Vector, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotRes, w.wantRes) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
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
// 		           in:nil,
// 		           opts:nil,
// 		       },
// 		       fields: fields {
// 		           c:nil,
// 		           addr:"",
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
// 		           in:nil,
// 		           opts:nil,
// 		           },
// 		           fields: fields {
// 		           c:nil,
// 		           addr:"",
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
// 			s := &specificAddrClient{
// 				c:    test.fields.c,
// 				addr: test.fields.addr,
// 			}
//
// 			gotRes, err := s.FilterVector(test.args.ctx, test.args.in, test.args.opts...)
// 			if err := checkFunc(test.want, gotRes, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_multipleAddrsClient_FilterVector(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		in   *payload.Object_Vector
// 		opts []grpc.CallOption
// 	}
// 	type fields struct {
// 		c     grpc.Client
// 		addrs []string
// 	}
// 	type want struct {
// 		wantRes *payload.Object_Vector
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Object_Vector, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotRes *payload.Object_Vector, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotRes, w.wantRes) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
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
// 		           in:nil,
// 		           opts:nil,
// 		       },
// 		       fields: fields {
// 		           c:nil,
// 		           addrs:nil,
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
// 		           in:nil,
// 		           opts:nil,
// 		           },
// 		           fields: fields {
// 		           c:nil,
// 		           addrs:nil,
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
// 			m := &multipleAddrsClient{
// 				c:     test.fields.c,
// 				addrs: test.fields.addrs,
// 			}
//
// 			gotRes, err := m.FilterVector(test.args.ctx, test.args.in, test.args.opts...)
// 			if err := checkFunc(test.want, gotRes, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
