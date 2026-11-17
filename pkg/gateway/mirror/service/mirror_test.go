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
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/pool"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/atomic"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/internal/test/goleak"
	grpcmock "github.com/vdaas/vald/internal/test/mock/grpc"
)

func Test_mirr_Connect(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx     context.Context //nolint:containedctx // table-driven test args struct, ctx passed through to the function under test
		targets []*payload.Mirror_Target
	}
	type fields struct {
		gateway      Gateway
		gatewayAddr  string
		selfMirrAddr string
	}
	type want struct {
		err error
	}
	type test struct {
		fields     fields
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
		name       string
		args       args
	}
	defaultCheckFunc := func(w want, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		return nil
	}
	tests := []test{
		func() test {
			gatewayAddr := testMirrorAddr1
			selfMirrorAddr := testMirrorAddr2
			return test{
				name: "Succeeded to connect to other mirror gateways",
				args: args{
					ctx: t.Context(),
					targets: []*payload.Mirror_Target{
						{
							Host: testGatewayHost1,
							Port: 8081,
						},
						{
							Host: testGatewayHost2,
							Port: 8081,
						},
					},
				},
				fields: fields{
					selfMirrAddr: selfMirrorAddr,
					gatewayAddr:  gatewayAddr,
					gateway: &GatewayMock{
						GRPCClientFunc: func() grpc.Client {
							return &grpcmock.GRPCClientMock{
								IsConnectedFunc: func(_ context.Context, _ string) bool {
									return false
								},
								ConnectFunc: func(_ context.Context, _ string, _ ...grpc.DialOption) (conn pool.Conn, err error) {
									return conn, err
								},
								SetDisableResolveDNSAddrFunc: func(addr string, disabled bool) {},
							}
						},
					},
				},
			}
		}(),
		func() test {
			gatewayAddr := testMirrorAddr1
			selfMirrorAddr := testMirrorAddr2
			return test{
				name: "Failed to connect to other mirror gateways due to an invalid address",
				args: args{
					ctx: t.Context(),
					targets: []*payload.Mirror_Target{
						{
							Host: testGatewayHost1,
						},
					},
				},
				fields: fields{
					selfMirrAddr: selfMirrorAddr,
					gatewayAddr:  gatewayAddr,
					gateway: &GatewayMock{
						GRPCClientFunc: func() grpc.Client {
							return &grpcmock.GRPCClientMock{
								IsConnectedFunc: func(_ context.Context, _ string) bool {
									return false
								},
								ConnectFunc: func(_ context.Context, _ string, _ ...grpc.DialOption) (pool.Conn, error) {
									return nil, errors.New("missing port in address")
								},
								SetDisableResolveDNSAddrFunc: func(addr string, disabled bool) {},
							}
						},
					},
				},
				want: want{
					err: errors.New("missing port in address"),
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
			m, err := NewMirror(
				WithSelfMirrorAddrs(test.fields.selfMirrAddr),
				WithGatewayAddrs(test.fields.gatewayAddr),
				WithGateway(test.fields.gateway),
			)
			if err != nil {
				t.Fatal(err)
			}

			err = m.Connect(test.args.ctx, test.args.targets...)
			if cerr := checkFunc(test.want, err); cerr != nil {
				tt.Errorf("error = %v", cerr)
			}
		})
	}
}

func Test_mirr_Disconnect(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx     context.Context //nolint:containedctx // table-driven test args struct, ctx passed through to the function under test
		targets []*payload.Mirror_Target
	}
	type fields struct {
		gateway      Gateway
		gatewayAddr  string
		selfMirrAddr string
	}
	type want struct {
		err error
	}
	type test struct {
		fields     fields
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
		name       string
		args       args
	}
	defaultCheckFunc := func(w want, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		return nil
	}
	tests := []test{
		func() test {
			gatewayAddr := testMirrorAddr1
			selfMirrorAddr := testMirrorAddr2
			return test{
				name: "Succeeded to disconnect to other mirror gateways",
				args: args{
					ctx: t.Context(),
					targets: []*payload.Mirror_Target{
						{
							Host: testGatewayHost1,
							Port: 8081,
						},
						{
							Host: testGatewayHost2,
							Port: 8081,
						},
					},
				},
				fields: fields{
					selfMirrAddr: selfMirrorAddr,
					gatewayAddr:  gatewayAddr,
					gateway: &GatewayMock{
						GRPCClientFunc: func() grpc.Client {
							return &grpcmock.GRPCClientMock{
								IsConnectedFunc: func(_ context.Context, _ string) bool {
									return true
								},
								DisconnectFunc: func(_ context.Context, _ string) error {
									return nil
								},
								SetDisableResolveDNSAddrFunc: func(addr string, disabled bool) {},
							}
						},
					},
				},
			}
		}(),
		func() test {
			gatewayAddr := testMirrorAddr1
			selfMirrorAddr := testMirrorAddr2
			return test{
				name: "Failed to connect to other mirror gateways due to an invalid address",
				args: args{
					ctx: t.Context(),
					targets: []*payload.Mirror_Target{
						{
							Host: testGatewayHost1,
						},
					},
				},
				fields: fields{
					selfMirrAddr: selfMirrorAddr,
					gatewayAddr:  gatewayAddr,
					gateway: &GatewayMock{
						GRPCClientFunc: func() grpc.Client {
							return &grpcmock.GRPCClientMock{
								IsConnectedFunc: func(_ context.Context, _ string) bool {
									return true
								},
								DisconnectFunc: func(_ context.Context, _ string) error {
									return errors.New("missing port in address")
								},
								SetDisableResolveDNSAddrFunc: func(addr string, disabled bool) {},
							}
						},
					},
				},
				want: want{
					err: errors.New("missing port in address"),
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
			m, err := NewMirror(
				WithSelfMirrorAddrs(test.fields.selfMirrAddr),
				WithGatewayAddrs(test.fields.gatewayAddr),
				WithGateway(test.fields.gateway),
			)
			if err != nil {
				t.Fatal(err)
			}

			err = m.Disconnect(test.args.ctx, test.args.targets...)
			if cerr := checkFunc(test.want, err); cerr != nil {
				tt.Errorf("error = %v", cerr)
			}
		})
	}
}

func Test_mirr_MirrorTargets(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context //nolint:containedctx // table-driven test args struct, ctx passed through to the function under test
	}
	type fields struct {
		gateway      Gateway
		gatewayAddr  string
		selfMirrAddr string
	}
	type want struct {
		err  error
		want []*payload.Mirror_Target
	}
	type test struct {
		args       args
		checkFunc  func(want, []*payload.Mirror_Target, error) error
		beforeFunc func(*testing.T, Mirror)
		afterFunc  func(*testing.T)
		fields     fields
		name       string
		want       want
	}
	defaultCheckFunc := func(w want, got []*payload.Mirror_Target, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if len(got) != len(w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		sort.Slice(got, func(i, j int) bool {
			return got[i].Host > got[j].Host
		})
		sort.Slice(w.want, func(i, j int) bool {
			return w.want[i].Host > w.want[j].Host
		})
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			gatewayAddr := testMirrorAddr1
			selfMirrorAddr := testMirrorAddr2
			connectTargets := []*payload.Mirror_Target{
				{
					Host: testMirrorHost1, // gateway addresses
					Port: 8081,
				},
				{
					Host: testGatewayHost1, // other mirror address
					Port: 8081,
				},
				{
					Host: testGatewayHost2, // other mirror address
					Port: 8081,
				},
			}
			connected := make(map[string]bool)
			return test{
				name: "returns only the addresses of the mirror gateways",
				args: args{
					ctx: t.Context(),
				},
				fields: fields{
					gatewayAddr:  gatewayAddr,
					selfMirrAddr: selfMirrorAddr,
					gateway: &GatewayMock{
						GRPCClientFunc: func() grpc.Client {
							return &grpcmock.GRPCClientMock{
								ConnectFunc: func(_ context.Context, addr string, _ ...grpc.DialOption) (conn pool.Conn, err error) {
									connected[addr] = true
									return conn, err
								},
								IsConnectedFunc: func(_ context.Context, addr string) bool {
									return connected[addr]
								},
								SetDisableResolveDNSAddrFunc: func(addr string, disabled bool) {},
							}
						},
					},
				},
				beforeFunc: func(t *testing.T, m Mirror) {
					t.Helper()
					if err := m.Connect(t.Context(), connectTargets...); err != nil {
						t.Fatal(err)
					}
				},
				want: want{
					want: []*payload.Mirror_Target{
						{
							Host: testGatewayHost1, // other mirror address
							Port: 8081,
						},
						{
							Host: testGatewayHost2, // other mirror address
							Port: 8081,
						},
						{
							Host: testMirrorHost2, // self mirror address
							Port: 8081,
						},
					},
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())

			m, err := NewMirror(
				WithSelfMirrorAddrs(test.fields.selfMirrAddr),
				WithGatewayAddrs(test.fields.gatewayAddr),
				WithGateway(test.fields.gateway),
			)
			if err != nil {
				t.Fatal(err)
			}

			if test.beforeFunc != nil {
				test.beforeFunc(tt, m)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			got, err := m.MirrorTargets(test.args.ctx)
			if cerr := checkFunc(test.want, got, err); cerr != nil {
				tt.Errorf("error = %v", cerr)
			}
		})
	}
}

func Test_mirr_connectedOtherMirrorAddrs(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context //nolint:containedctx // table-driven test args struct, ctx passed through to the function under test
	}
	type fields struct {
		gateway      Gateway
		gatewayAddr  string
		selfMirrAddr string
	}
	type want struct {
		want []string
	}
	type test struct {
		fields     fields
		args       args
		checkFunc  func(want, []string) error
		beforeFunc func(*testing.T, Mirror)
		afterFunc  func(*testing.T)
		name       string
		want       want
	}
	defaultCheckFunc := func(w want, got []string) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			gatewayAddr := testMirrorAddr1
			selfMirrorAddr := testMirrorAddr2
			connectTargets := []*payload.Mirror_Target{
				{
					Host: testMirrorHost1, // gateway addresses
					Port: 8081,
				},
				{
					Host: testGatewayHost1, // other mirror address
					Port: 8081,
				},
			}
			connected := make(map[string]bool)
			return test{
				name: "returns only the address of the other mirror gateway",
				args: args{
					ctx: t.Context(),
				},
				fields: fields{
					selfMirrAddr: selfMirrorAddr,
					gatewayAddr:  gatewayAddr,
					gateway: &GatewayMock{
						GRPCClientFunc: func() grpc.Client {
							return &grpcmock.GRPCClientMock{
								ConnectFunc: func(_ context.Context, addr string, _ ...grpc.DialOption) (conn pool.Conn, err error) {
									connected[addr] = true
									return conn, err
								},
								IsConnectedFunc: func(_ context.Context, addr string) bool {
									return connected[addr]
								},
								SetDisableResolveDNSAddrFunc: func(addr string, disabled bool) {},
							}
						},
					},
				},
				beforeFunc: func(t *testing.T, m Mirror) {
					t.Helper()
					if err := m.Connect(t.Context(), connectTargets...); err != nil {
						t.Fatal(err)
					}
				},
				want: want{
					want: []string{
						"192.168.2.2:8081",
					},
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())

			m, err := NewMirror(
				WithSelfMirrorAddrs(test.fields.selfMirrAddr),
				WithGatewayAddrs(test.fields.gatewayAddr),
				WithGateway(test.fields.gateway),
			)
			if err != nil {
				t.Fatal(err)
			}

			if test.beforeFunc != nil {
				test.beforeFunc(tt, m)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			if mirr, ok := m.(*mirr); ok {
				got := mirr.connectedOtherMirrorAddrs(test.args.ctx)
				if err := checkFunc(test.want, got); err != nil {
					tt.Errorf("error = %v", err)
				}
			}
		})
	}
}

type registerClientStub struct {
	MirrorClient
	res *payload.Mirror_Targets
}

func (c *registerClientStub) Register(
	_ context.Context, _ *payload.Mirror_Targets, _ ...grpc.CallOption,
) (*payload.Mirror_Targets, error) {
	return c.res, nil
}

// Test_mirr_registers_dedup verifies that every distinct target returned by the
// peers is collected exactly once, even if several peers report the same list.
func Test_mirr_registers_dedup(t *testing.T) {
	t.Parallel()
	peers := []string{"peer-1:8081", "peer-2:8081"}
	res := &payload.Mirror_Targets{
		Targets: []*payload.Mirror_Target{
			{Host: "192.168.0.1", Port: 8081},
			{Host: "192.168.0.2", Port: 8081},
		},
	}
	m := &mirr{
		gateway: &GatewayMock{
			DoMultiFunc: func(ctx context.Context, _ []string, f func(context.Context, string, MirrorClient, ...grpc.CallOption) error) error {
				for _, peer := range peers {
					if err := f(ctx, peer, &registerClientStub{res: res}); err != nil {
						return err
					}
				}
				return nil
			},
		},
	}

	got, err := m.registers(t.Context(), &payload.Mirror_Targets{})
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != len(res.GetTargets()) {
		t.Fatalf("got %d targets, want %d: %v", len(got), len(res.GetTargets()), got)
	}
	seen := make(map[string]int, len(got))
	for _, tgt := range got {
		seen[tgt.GetHost()]++
	}
	for _, want := range res.GetTargets() {
		if n := seen[want.GetHost()]; n != 1 {
			t.Errorf("target %s appeared %d times, want 1", want.GetHost(), n)
		}
	}
}

// One unreachable target must not stop Connect from connecting the rest.
func Test_mirr_Connect_continuesPastFailure(t *testing.T) {
	t.Parallel()
	const (
		badAddr  = "192.168.10.1:8081"
		goodAddr = "192.168.10.2:8081"
	)
	failErr := errors.New("dial failure")
	var (
		mu        sync.Mutex
		attempted []string
	)
	m := &mirr{
		gateway: &GatewayMock{
			GRPCClientFunc: func() grpc.Client {
				return &grpcmock.GRPCClientMock{
					IsConnectedFunc: func(context.Context, string) bool { return false },
					ConnectFunc: func(_ context.Context, addr string, _ ...grpc.DialOption) (pool.Conn, error) {
						mu.Lock()
						attempted = append(attempted, addr)
						mu.Unlock()
						if addr == badAddr {
							return nil, failErr
						}
						return nil, nil
					},
					SetDisableResolveDNSAddrFunc: func(string, bool) {},
				}
			},
		},
	}

	err := m.Connect(t.Context(),
		&payload.Mirror_Target{Host: "192.168.10.1", Port: 8081},
		&payload.Mirror_Target{Host: "192.168.10.2", Port: 8081},
	)
	if !errors.Is(err, failErr) {
		t.Errorf("error = %v, want %v", err, failErr)
	}
	mu.Lock()
	got := append([]string(nil), attempted...)
	mu.Unlock()
	if len(got) != 2 || got[0] != badAddr || got[1] != goodAddr {
		t.Errorf("attempted = %v, want [%s %s]", got, badAddr, goodAddr)
	}
	if _, ok := m.addrs.Load(goodAddr); !ok {
		t.Errorf("%s not stored after successful connect", goodAddr)
	}
	if _, ok := m.addrs.Load(badAddr); ok {
		t.Errorf("%s stored after failed connect", badAddr)
	}
}

// Start must not run the register broadcast while no peer is connected
// (it can only fail with ErrTargetNotFound), and must run it once one is.
func Test_mirr_Start_registerOnlyWithPeers(t *testing.T) {
	t.Parallel()
	const peerAddr = "192.168.20.1:8081"
	tests := []struct {
		name          string
		peerConnected bool
	}{
		{name: "no connected peer skips register", peerConnected: false},
		{name: "connected peer runs register", peerConnected: true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(tt *testing.T) {
			tt.Parallel()
			ctx, cancel := context.WithTimeout(tt.Context(), 10*time.Second)
			defer cancel()

			var (
				isConnectedCalls atomic.Int64
				doMultiCalls     atomic.Int64
				once             sync.Once
			)
			// Signals that the loop has ticked enough (or reached register).
			reached := make(chan struct{})
			signal := func() { once.Do(func() { close(reached) }) }

			eg, egctx := errgroup.New(ctx)
			m := &mirr{
				eg:          eg,
				registerDur: time.Millisecond,
				gateway: &GatewayMock{
					GRPCClientFunc: func() grpc.Client {
						return &grpcmock.GRPCClientMock{
							IsConnectedFunc: func(context.Context, string) bool {
								// MirrorTargets and connectedOtherMirrorAddrs each
								// check the peer once per tick: 6 calls = 3 ticks.
								if isConnectedCalls.Add(1) >= 6 && !tc.peerConnected {
									signal()
								}
								return tc.peerConnected
							},
							ConnectedAddrsFunc: func() []string { return nil },
						}
					},
					DoMultiFunc: func(_ context.Context, targets []string, _ func(context.Context, string, MirrorClient, ...grpc.CallOption) error) error {
						doMultiCalls.Add(1)
						signal()
						if len(targets) == 0 {
							return errors.ErrTargetNotFound
						}
						return nil
					},
				},
			}
			m.addrs.Store(peerAddr, struct{}{})

			ech := m.Start(egctx)
			select {
			case <-reached:
			case <-ctx.Done():
				tt.Fatal("mirror register loop did not progress before timeout")
			}
			cancel()

			var errs []error
			for done := false; !done; {
				select {
				case err, ok := <-ech:
					if !ok {
						done = true
						break
					}
					errs = append(errs, err)
				case <-time.After(10 * time.Second):
					tt.Fatal("error channel was not closed after cancel")
				}
			}
			if err := eg.Wait(); err != nil && !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded) {
				tt.Errorf("unexpected errgroup error: %v", err)
			}

			if tc.peerConnected {
				if doMultiCalls.Load() == 0 {
					tt.Error("register was not called although a peer is connected")
				}
				return
			}
			if n := doMultiCalls.Load(); n != 0 {
				tt.Errorf("register called %d times without connected peers", n)
			}
			if len(errs) != 0 {
				tt.Errorf("errors reported without connected peers: %v", errs)
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
//
// func TestNewMirrorClient(t *testing.T) {
// 	type args struct {
// 		conn *grpc.ClientConn
// 	}
// 	type want struct {
// 		want MirrorClient
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, MirrorClient) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got MirrorClient) error {
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
// 		           conn:nil,
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
// 		           conn:nil,
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
// 			got := NewMirrorClient(test.args.conn)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestNewMirror(t *testing.T) {
// 	type args struct {
// 		opts []MirrorOption
// 	}
// 	type want struct {
// 		want Mirror
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, Mirror, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got Mirror, err error) error {
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
// 			got, err := NewMirror(test.args.opts...)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_mirr_Start(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		eg            errgroup.Group
// 		gateway       Gateway
// 		selfMirrTgts  []*payload.Mirror_Target
// 		addrs         sync.Map[string, any]
// 		selfMirrAddrs sync.Map[string, any]
// 		gwAddrs       sync.Map[string, any]
// 		registerDur   time.Duration
// 	}
// 	type want struct {
// 		want <-chan error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, <-chan error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got <-chan error) error {
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
// 		           eg:nil,
// 		           gateway:nil,
// 		           selfMirrTgts:nil,
// 		           addrs:nil,
// 		           selfMirrAddrs:nil,
// 		           gwAddrs:nil,
// 		           registerDur:nil,
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
// 		           eg:nil,
// 		           gateway:nil,
// 		           selfMirrTgts:nil,
// 		           addrs:nil,
// 		           selfMirrAddrs:nil,
// 		           gwAddrs:nil,
// 		           registerDur:nil,
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
// 			m := &mirr{
// 				eg:            test.fields.eg,
// 				gateway:       test.fields.gateway,
// 				selfMirrTgts:  test.fields.selfMirrTgts,
// 				addrs:         test.fields.addrs,
// 				selfMirrAddrs: test.fields.selfMirrAddrs,
// 				gwAddrs:       test.fields.gwAddrs,
// 				registerDur:   test.fields.registerDur,
// 			}
//
// 			got := m.Start(test.args.ctx)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_mirr_registers(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		tgts *payload.Mirror_Targets
// 	}
// 	type fields struct {
// 		eg            errgroup.Group
// 		gateway       Gateway
// 		selfMirrTgts  []*payload.Mirror_Target
// 		addrs         sync.Map[string, any]
// 		selfMirrAddrs sync.Map[string, any]
// 		gwAddrs       sync.Map[string, any]
// 		registerDur   time.Duration
// 	}
// 	type want struct {
// 		want []*payload.Mirror_Target
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, []*payload.Mirror_Target, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got []*payload.Mirror_Target, err error) error {
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
// 		           tgts:nil,
// 		       },
// 		       fields: fields {
// 		           eg:nil,
// 		           gateway:nil,
// 		           selfMirrTgts:nil,
// 		           addrs:nil,
// 		           selfMirrAddrs:nil,
// 		           gwAddrs:nil,
// 		           registerDur:nil,
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
// 		           tgts:nil,
// 		           },
// 		           fields: fields {
// 		           eg:nil,
// 		           gateway:nil,
// 		           selfMirrTgts:nil,
// 		           addrs:nil,
// 		           selfMirrAddrs:nil,
// 		           gwAddrs:nil,
// 		           registerDur:nil,
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
// 			m := &mirr{
// 				eg:            test.fields.eg,
// 				gateway:       test.fields.gateway,
// 				selfMirrTgts:  test.fields.selfMirrTgts,
// 				addrs:         test.fields.addrs,
// 				selfMirrAddrs: test.fields.selfMirrAddrs,
// 				gwAddrs:       test.fields.gwAddrs,
// 				registerDur:   test.fields.registerDur,
// 			}
//
// 			got, err := m.registers(test.args.ctx, test.args.tgts)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_mirr_IsConnected(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		addr string
// 	}
// 	type fields struct {
// 		eg            errgroup.Group
// 		gateway       Gateway
// 		selfMirrTgts  []*payload.Mirror_Target
// 		addrs         sync.Map[string, any]
// 		selfMirrAddrs sync.Map[string, any]
// 		gwAddrs       sync.Map[string, any]
// 		registerDur   time.Duration
// 	}
// 	type want struct {
// 		want bool
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, bool) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got bool) error {
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
// 		           addr:"",
// 		       },
// 		       fields: fields {
// 		           eg:nil,
// 		           gateway:nil,
// 		           selfMirrTgts:nil,
// 		           addrs:nil,
// 		           selfMirrAddrs:nil,
// 		           gwAddrs:nil,
// 		           registerDur:nil,
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
// 		           addr:"",
// 		           },
// 		           fields: fields {
// 		           eg:nil,
// 		           gateway:nil,
// 		           selfMirrTgts:nil,
// 		           addrs:nil,
// 		           selfMirrAddrs:nil,
// 		           gwAddrs:nil,
// 		           registerDur:nil,
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
// 			m := &mirr{
// 				eg:            test.fields.eg,
// 				gateway:       test.fields.gateway,
// 				selfMirrTgts:  test.fields.selfMirrTgts,
// 				addrs:         test.fields.addrs,
// 				selfMirrAddrs: test.fields.selfMirrAddrs,
// 				gwAddrs:       test.fields.gwAddrs,
// 				registerDur:   test.fields.registerDur,
// 			}
//
// 			got := m.IsConnected(test.args.ctx, test.args.addr)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_mirr_isSelfMirrorAddr(t *testing.T) {
// 	type args struct {
// 		addr string
// 	}
// 	type fields struct {
// 		eg            errgroup.Group
// 		gateway       Gateway
// 		selfMirrTgts  []*payload.Mirror_Target
// 		addrs         sync.Map[string, any]
// 		selfMirrAddrs sync.Map[string, any]
// 		gwAddrs       sync.Map[string, any]
// 		registerDur   time.Duration
// 	}
// 	type want struct {
// 		want bool
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, bool) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got bool) error {
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
// 		           addr:"",
// 		       },
// 		       fields: fields {
// 		           eg:nil,
// 		           gateway:nil,
// 		           selfMirrTgts:nil,
// 		           addrs:nil,
// 		           selfMirrAddrs:nil,
// 		           gwAddrs:nil,
// 		           registerDur:nil,
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
// 		           addr:"",
// 		           },
// 		           fields: fields {
// 		           eg:nil,
// 		           gateway:nil,
// 		           selfMirrTgts:nil,
// 		           addrs:nil,
// 		           selfMirrAddrs:nil,
// 		           gwAddrs:nil,
// 		           registerDur:nil,
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
// 			m := &mirr{
// 				eg:            test.fields.eg,
// 				gateway:       test.fields.gateway,
// 				selfMirrTgts:  test.fields.selfMirrTgts,
// 				addrs:         test.fields.addrs,
// 				selfMirrAddrs: test.fields.selfMirrAddrs,
// 				gwAddrs:       test.fields.gwAddrs,
// 				registerDur:   test.fields.registerDur,
// 			}
//
// 			got := m.isSelfMirrorAddr(test.args.addr)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_mirr_isGatewayAddr(t *testing.T) {
// 	type args struct {
// 		addr string
// 	}
// 	type fields struct {
// 		eg            errgroup.Group
// 		gateway       Gateway
// 		selfMirrTgts  []*payload.Mirror_Target
// 		addrs         sync.Map[string, any]
// 		selfMirrAddrs sync.Map[string, any]
// 		gwAddrs       sync.Map[string, any]
// 		registerDur   time.Duration
// 	}
// 	type want struct {
// 		want bool
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, bool) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got bool) error {
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
// 		           addr:"",
// 		       },
// 		       fields: fields {
// 		           eg:nil,
// 		           gateway:nil,
// 		           selfMirrTgts:nil,
// 		           addrs:nil,
// 		           selfMirrAddrs:nil,
// 		           gwAddrs:nil,
// 		           registerDur:nil,
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
// 		           addr:"",
// 		           },
// 		           fields: fields {
// 		           eg:nil,
// 		           gateway:nil,
// 		           selfMirrTgts:nil,
// 		           addrs:nil,
// 		           selfMirrAddrs:nil,
// 		           gwAddrs:nil,
// 		           registerDur:nil,
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
// 			m := &mirr{
// 				eg:            test.fields.eg,
// 				gateway:       test.fields.gateway,
// 				selfMirrTgts:  test.fields.selfMirrTgts,
// 				addrs:         test.fields.addrs,
// 				selfMirrAddrs: test.fields.selfMirrAddrs,
// 				gwAddrs:       test.fields.gwAddrs,
// 				registerDur:   test.fields.registerDur,
// 			}
//
// 			got := m.isGatewayAddr(test.args.addr)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_mirr_RangeMirrorAddr(t *testing.T) {
// 	type args struct {
// 		f func(addr string, _ any) bool
// 	}
// 	type fields struct {
// 		eg            errgroup.Group
// 		gateway       Gateway
// 		selfMirrTgts  []*payload.Mirror_Target
// 		addrs         sync.Map[string, any]
// 		selfMirrAddrs sync.Map[string, any]
// 		gwAddrs       sync.Map[string, any]
// 		registerDur   time.Duration
// 	}
// 	type want struct{}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want) error {
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           f:nil,
// 		       },
// 		       fields: fields {
// 		           eg:nil,
// 		           gateway:nil,
// 		           selfMirrTgts:nil,
// 		           addrs:nil,
// 		           selfMirrAddrs:nil,
// 		           gwAddrs:nil,
// 		           registerDur:nil,
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
// 		           f:nil,
// 		           },
// 		           fields: fields {
// 		           eg:nil,
// 		           gateway:nil,
// 		           selfMirrTgts:nil,
// 		           addrs:nil,
// 		           selfMirrAddrs:nil,
// 		           gwAddrs:nil,
// 		           registerDur:nil,
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
// 			m := &mirr{
// 				eg:            test.fields.eg,
// 				gateway:       test.fields.gateway,
// 				selfMirrTgts:  test.fields.selfMirrTgts,
// 				addrs:         test.fields.addrs,
// 				selfMirrAddrs: test.fields.selfMirrAddrs,
// 				gwAddrs:       test.fields.gwAddrs,
// 				registerDur:   test.fields.registerDur,
// 			}
//
// 			m.RangeMirrorAddr(test.args.f)
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
