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
	"testing"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/k8s"
	"github.com/vdaas/vald/internal/k8s/vald/mirror/target"
	"github.com/vdaas/vald/internal/sync/atomic"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/internal/test/goleak"
	k8smock "github.com/vdaas/vald/internal/test/mock/k8s"
)

const (
	testMirrorName1  = "mirror-1"
	testMirrorHost1  = "192.168.1.2"
	testMirrorHost2  = "192.168.1.3"
	testMirrorAddr1  = "192.168.1.2:8081"
	testMirrorAddr2  = "192.168.1.3:8081"
	testGatewayHost1 = "192.168.2.2"
	testGatewayHost2 = "192.168.3.2"
)

//nolint:maintidx // table-driven test with 6 literal test-case constructors; complexity is inherent to the pattern, not tangled logic
func Test_discovery_startSync(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx  context.Context //nolint:containedctx // table-driven test args struct, ctx passed through to the function under test
		prev map[string]target.Endpoint
	}
	type fields struct {
		ctrl k8s.Controller
		mirr Mirror
	}
	type want struct {
		want map[string]target.Endpoint
		err  error
	}
	type test struct {
		fields     fields
		args       args
		want       want
		checkFunc  func(want, map[string]target.Endpoint, error) error
		beforeFunc func(*testing.T, *discovery, args)
		afterFunc  func(*testing.T, args)
		name       string
	}
	defaultCheckFunc := func(w want, got map[string]target.Endpoint, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			prev := make(map[string]target.Endpoint)
			current := map[string]target.Endpoint{
				testMirrorName1: {
					Host: testMirrorHost1,
					Port: 8081,
				},
			}
			return test{
				name: "Succeeded detecting the created resource and connecting to the new target",
				args: args{
					ctx:  t.Context(),
					prev: prev,
				},
				fields: fields{
					mirr: &MirrorMock{
						ConnectFunc: func(_ context.Context, _ ...*payload.Mirror_Target) error {
							return nil
						},
					},
					ctrl: &k8smock.ControllerMock{
						GetManagerFunc: k8smock.NewDefaultManagerMock,
					},
				},
				beforeFunc: func(t *testing.T, d *discovery, _ args) {
					t.Helper()
					d.onReconcile(t.Context(), current)
				},
				want: want{
					want: current,
				},
			}
		}(),
		func() test {
			prev := make(map[string]target.Endpoint)
			current := map[string]target.Endpoint{
				testMirrorName1: {
					Host: testMirrorHost1,
					Port: 8081,
				},
			}
			return test{
				name: "Succeeded detecting the created resource but failed to connect to the new target",
				args: args{
					ctx:  t.Context(),
					prev: prev,
				},
				fields: fields{
					mirr: &MirrorMock{
						ConnectFunc: func(_ context.Context, _ ...*payload.Mirror_Target) error {
							return errors.New("err")
						},
					},
					ctrl: &k8smock.ControllerMock{
						GetManagerFunc: k8smock.NewDefaultManagerMock,
					},
				},
				beforeFunc: func(t *testing.T, d *discovery, _ args) {
					t.Helper()
					d.onReconcile(t.Context(), current)
				},
				want: want{
					want: current,
					err:  errors.New("err"),
				},
			}
		}(),
		func() test {
			prev := map[string]target.Endpoint{
				testMirrorName1: {
					Host: testMirrorHost1,
					Port: 8081,
				},
			}
			current := make(map[string]target.Endpoint)
			return test{
				name: "Succeeded detecting the deleted resource and disconnecting the target",
				args: args{
					ctx:  t.Context(),
					prev: prev,
				},
				fields: fields{
					mirr: &MirrorMock{
						DisconnectFunc: func(_ context.Context, _ ...*payload.Mirror_Target) error {
							return nil
						},
					},
					ctrl: &k8smock.ControllerMock{
						GetManagerFunc: k8smock.NewDefaultManagerMock,
					},
				},
				beforeFunc: func(t *testing.T, d *discovery, _ args) {
					t.Helper()
					d.onReconcile(t.Context(), current)
				},
				want: want{
					want: current,
				},
			}
		}(),
		func() test {
			prev := map[string]target.Endpoint{
				testMirrorName1: {
					Host: testMirrorHost1,
					Port: 8081,
				},
			}
			current := make(map[string]target.Endpoint)
			return test{
				name: "Succeeded detecting the deleted resource but failed to disconnect the target",
				args: args{
					ctx:  t.Context(),
					prev: prev,
				},
				fields: fields{
					mirr: &MirrorMock{
						DisconnectFunc: func(_ context.Context, _ ...*payload.Mirror_Target) error {
							return errors.New("err")
						},
					},
					ctrl: &k8smock.ControllerMock{
						GetManagerFunc: k8smock.NewDefaultManagerMock,
					},
				},
				beforeFunc: func(t *testing.T, d *discovery, _ args) {
					t.Helper()
					d.onReconcile(t.Context(), current)
				},
				want: want{
					want: current,
					err:  errors.New("err"),
				},
			}
		}(),
		func() test {
			prev := map[string]target.Endpoint{
				testMirrorName1: {
					Host: testMirrorHost1,
					Port: 8081,
				},
			}
			current := map[string]target.Endpoint{
				testMirrorName1: {
					Host: testMirrorHost2,
					Port: 8081,
				},
			}
			return test{
				name: "Succeeded detecting the updated resource and updating the target connection",
				args: args{
					ctx:  t.Context(),
					prev: prev,
				},
				fields: fields{
					mirr: &MirrorMock{
						ConnectFunc: func(_ context.Context, _ ...*payload.Mirror_Target) error {
							return nil
						},
						DisconnectFunc: func(_ context.Context, _ ...*payload.Mirror_Target) error {
							return nil
						},
					},
					ctrl: &k8smock.ControllerMock{
						GetManagerFunc: k8smock.NewDefaultManagerMock,
					},
				},
				beforeFunc: func(t *testing.T, d *discovery, _ args) {
					t.Helper()
					d.onReconcile(t.Context(), current)
				},
				want: want{
					want: current,
				},
			}
		}(),
		func() test {
			prev := map[string]target.Endpoint{
				testMirrorName1: {
					Host: testMirrorHost1,
					Port: 8081,
				},
			}
			current := map[string]target.Endpoint{
				testMirrorName1: {
					Host: testMirrorHost2,
					Port: 8081,
				},
			}
			return test{
				name: "Succeeded detecting the updated resource and failed to update the target connection",
				args: args{
					ctx:  t.Context(),
					prev: prev,
				},
				fields: fields{
					mirr: &MirrorMock{
						ConnectFunc: func(_ context.Context, _ ...*payload.Mirror_Target) error {
							return nil
						},
						DisconnectFunc: func(_ context.Context, _ ...*payload.Mirror_Target) error {
							return errors.New("err")
						},
					},
					ctrl: &k8smock.ControllerMock{
						GetManagerFunc: k8smock.NewDefaultManagerMock,
					},
				},
				beforeFunc: func(t *testing.T, d *discovery, _ args) {
					t.Helper()
					d.onReconcile(t.Context(), current)
				},
				want: want{
					want: current,
					err:  errors.New("err"),
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())

			d, err := NewDiscovery(
				WithDiscoveryController(test.fields.ctrl),
				WithDiscoveryMirror(test.fields.mirr),
			)
			if err != nil {
				t.Fatal(err)
			}

			if dis, ok := d.(*discovery); ok {
				if test.beforeFunc != nil {
					test.beforeFunc(tt, dis, test.args)
				}
				if test.afterFunc != nil {
					defer test.afterFunc(tt, test.args)
				}
				checkFunc := test.checkFunc
				if test.checkFunc == nil {
					checkFunc = defaultCheckFunc
				}

				got, err := dis.startSync(test.args.ctx, test.args.prev)
				if cerr := checkFunc(test.want, got, err); cerr != nil {
					tt.Errorf("error = %v", cerr)
				}
			}
		})
	}
}

func Test_discovery_syncWithAddr(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx      context.Context //nolint:containedctx // table-driven test args struct, ctx passed through to the function under test
		current  map[string]target.Endpoint
		curAddrs map[string]string
	}
	type fields struct {
		ctrl k8s.Controller
		mirr Mirror
	}
	type want struct {
		err error
	}
	type test struct {
		args       args
		fields     fields
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(*testing.T, *discovery, args)
		afterFunc  func(*testing.T, args)
		name       string
	}
	defaultCheckFunc := func(w want, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		return nil
	}
	tests := []test{
		func() test {
			current := map[string]target.Endpoint{
				testMirrorName1: {
					Host:  testMirrorHost1,
					Port:  8081,
					Phase: target.MirrorTargetPhaseConnected,
				},
			}
			curAddrs := map[string]string{
				testMirrorAddr1: testMirrorName1,
			}

			return test{
				name: "Succeeded to change to disconnected phase when curAddrs are not connected",
				args: args{
					ctx:      t.Context(),
					current:  current,
					curAddrs: curAddrs,
				},
				fields: fields{
					mirr: &MirrorMock{
						IsConnectedFunc: func(_ context.Context, _ string) bool {
							return false
						},
						RangeMirrorAddrFunc: func(_ func(addr string, _ any) bool) {
							// There is no connection.
						},
					},
					ctrl: &k8smock.ControllerMock{
						GetManagerFunc: k8smock.NewDefaultManagerMock,
					},
				},
				beforeFunc: func(t *testing.T, d *discovery, _ args) {
					t.Helper()
					d.onReconcile(t.Context(), current)
				},
			}
		}(),
		func() test {
			current := map[string]target.Endpoint{
				testMirrorName1: {
					Host:  testMirrorHost1,
					Port:  8081,
					Phase: target.MirrorTargetPhaseDisconnected,
				},
			}
			curAddrs := map[string]string{
				testMirrorAddr1: testMirrorName1,
			}

			return test{
				name: "Succeeded to change to connected phase when curAddrs are connected",
				args: args{
					ctx:      t.Context(),
					current:  current,
					curAddrs: curAddrs,
				},
				fields: fields{
					mirr: &MirrorMock{
						IsConnectedFunc: func(_ context.Context, _ string) bool {
							return true
						},
						RangeMirrorAddrFunc: func(f func(addr string, _ any) bool) {
							for addr := range curAddrs {
								f(addr, struct{}{})
							}
						},
					},
					ctrl: &k8smock.ControllerMock{
						GetManagerFunc: k8smock.NewDefaultManagerMock,
					},
				},
			}
		}(),
		func() test {
			current := map[string]target.Endpoint{}
			curAddrs := map[string]string{}
			newAddrs := []string{
				testMirrorAddr1,
			}

			return test{
				name: "Succeeded to create new resource when there is a new connection",
				args: args{
					ctx:      t.Context(),
					current:  current,
					curAddrs: curAddrs,
				},
				fields: fields{
					mirr: &MirrorMock{
						IsConnectedFunc: func(_ context.Context, _ string) bool {
							return true
						},
						RangeMirrorAddrFunc: func(f func(addr string, _ any) bool) {
							for _, addr := range newAddrs {
								f(addr, struct{}{})
							}
						},
					},
					ctrl: &k8smock.ControllerMock{
						GetManagerFunc: k8smock.NewDefaultManagerMock,
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

			d, err := NewDiscovery(
				WithDiscoveryController(test.fields.ctrl),
				WithDiscoveryMirror(test.fields.mirr),
			)
			if err != nil {
				t.Fatal(err)
			}

			if dis, ok := d.(*discovery); ok {
				if test.beforeFunc != nil {
					test.beforeFunc(tt, dis, test.args)
				}
				if test.afterFunc != nil {
					defer test.afterFunc(tt, test.args)
				}
				checkFunc := test.checkFunc
				if test.checkFunc == nil {
					checkFunc = defaultCheckFunc
				}

				err := dis.syncWithAddr(test.args.ctx, test.args.current, test.args.curAddrs)
				if cerr := checkFunc(test.want, err); cerr != nil {
					tt.Errorf("error = %v", cerr)
				}
			}
		})
	}
}

// Test_discovery_syncWithAddr_updatesPhaseOncePerSync guards against the second
// (RangeMirrorAddr) loop of syncWithAddr repeating a phase update the first loop
// already made for the same target.
func Test_discovery_syncWithAddr_updatesPhaseOncePerSync(t *testing.T) {
	tests := []struct {
		name      string
		phase     target.MirrorTargetPhase
		connected bool
		want      int64
	}{
		{name: "Disconnected -> Connected", phase: target.MirrorTargetPhaseDisconnected, connected: true, want: 1},
		{name: "Connected -> Disconnected", phase: target.MirrorTargetPhaseConnected, connected: false, want: 1},
		{name: "Pending -> Disconnected", phase: target.MirrorTargetPhasePending, connected: false, want: 1},
		{name: "Connected stays Connected", phase: target.MirrorTargetPhaseConnected, connected: true, want: 0},
		{name: "Disconnected stays Disconnected", phase: target.MirrorTargetPhaseDisconnected, connected: false, want: 0},
	}
	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			var gets atomic.Int64
			mgr := func() k8s.Manager {
				m := k8smock.NewDefaultManagerMock().(*k8smock.ManagerMock)
				getClient := m.GetClientFunc
				m.GetClientFunc = func() k8s.Client {
					c := getClient().(*k8smock.ClientMock)
					c.GetFunc = func(context.Context, k8s.ObjectKey, k8s.Object, ...k8s.GetOption) error {
						gets.Add(1)
						return nil
					}
					return c
				}
				return m
			}
			d, err := NewDiscovery(
				WithDiscoveryController(&k8smock.ControllerMock{GetManagerFunc: mgr}),
				WithDiscoveryMirror(&MirrorMock{
					IsConnectedFunc: func(context.Context, string) bool { return test.connected },
					RangeMirrorAddrFunc: func(f func(addr string, _ any) bool) {
						f(testMirrorAddr1, struct{}{})
					},
				}),
			)
			if err != nil {
				tt.Fatal(err)
			}
			current := map[string]target.Endpoint{
				testMirrorName1: {Host: testMirrorHost1, Port: 8081, Phase: test.phase},
			}
			if err := d.(*discovery).syncWithAddr(tt.Context(), current, map[string]string{testMirrorAddr1: testMirrorName1}); err != nil {
				tt.Fatalf("syncWithAddr() error = %v", err)
			}
			if got := gets.Load(); got != test.want {
				tt.Errorf("phase updates (mirror target Gets): got %d, want %d", got, test.want)
			}
		})
	}
}

// Test_discovery_startSync_retriesFailedConnect pins that a MirrorTarget whose
// first Connect failed is reconnected by later syncs, while a target that was
// connected and then dropped (Mirror.registers disconnects a peer whose
// Register returns Unimplemented) is not redialed every tick.
func Test_discovery_startSync_retriesFailedConnect(t *testing.T) {
	t.Parallel()
	var (
		connects  atomic.Int64
		connected atomic.Bool
		failDial  atomic.Bool
	)
	failDial.Store(true)
	d, err := NewDiscovery(
		WithDiscoveryController(&k8smock.ControllerMock{GetManagerFunc: k8smock.NewDefaultManagerMock}),
		WithDiscoveryMirror(&MirrorMock{
			ConnectFunc: func(context.Context, ...*payload.Mirror_Target) error {
				connects.Add(1)
				if failDial.Load() {
					return errors.New("dial failed")
				}
				connected.Store(true)
				return nil
			},
			IsConnectedFunc: func(context.Context, string) bool { return connected.Load() },
			RangeMirrorAddrFunc: func(func(addr string, _ any) bool) {
				// The failed Connect removed the address from the mirror's addrs.
			},
		}),
	)
	if err != nil {
		t.Fatal(err)
	}
	dis := d.(*discovery)
	ctx := t.Context()
	pending := map[string]target.Endpoint{
		testMirrorName1: {Host: testMirrorHost1, Port: 8081, Phase: target.MirrorTargetPhasePending},
	}
	dis.onReconcile(ctx, pending)

	// 1st sync: the target is created and its Connect fails.
	prev, err := dis.startSync(ctx, map[string]target.Endpoint{})
	if err == nil {
		t.Fatal("1st sync: want the Connect error")
	}
	if got := connects.Load(); got != 1 {
		t.Fatalf("1st sync: Connect calls = %d, want 1", got)
	}

	// 2nd sync: no CR change and the peer is still unreachable; retried and failing.
	if prev, err = dis.startSync(ctx, prev); err == nil {
		t.Fatal("2nd sync: want the retried Connect error")
	}
	if got := connects.Load(); got != 2 {
		t.Fatalf("2nd sync: Connect calls = %d, want 2 (failed target must be retried)", got)
	}

	// 3rd sync: the peer became reachable; the retry connects it.
	failDial.Store(false)
	if prev, err = dis.startSync(ctx, prev); err != nil {
		t.Fatalf("3rd sync: %v", err)
	}
	if got := connects.Load(); got != 3 {
		t.Fatalf("3rd sync: Connect calls = %d, want 3", got)
	}
	if !connected.Load() {
		t.Fatal("3rd sync: target not connected")
	}

	// 4th sync: the connected target was dropped by registers (Unimplemented);
	// it must only be marked Disconnected, not redialed.
	connected.Store(false)
	dis.onReconcile(ctx, map[string]target.Endpoint{
		testMirrorName1: {Host: testMirrorHost1, Port: 8081, Phase: target.MirrorTargetPhaseConnected},
	})
	if _, err = dis.startSync(ctx, prev); err != nil {
		t.Fatalf("4th sync: %v", err)
	}
	if got := connects.Load(); got != 3 {
		t.Fatalf("4th sync: Connect calls = %d, want 3 (deliberately disconnected target must not be redialed)", got)
	}
}

// Test_discovery_startSync_retriesFailedDisconnect pins that a replaced or
// removed MirrorTarget whose Disconnect failed is disconnected by later syncs,
// that a replaced target's new endpoint is connected even when the old
// endpoint's Disconnect failed, and that the still-open old connection is not
// turned into an auto-created MirrorTarget meanwhile.
func Test_discovery_startSync_retriesFailedDisconnect(t *testing.T) {
	t.Parallel()
	tests := []struct {
		prev    map[string]target.Endpoint
		current map[string]target.Endpoint
		name    string
		// wantConnects is the number of Connect calls made by the 1st sync.
		wantConnects int64
	}{
		{
			name:         "updated target",
			prev:         map[string]target.Endpoint{testMirrorName1: {Host: testMirrorHost1, Port: 8081}},
			current:      map[string]target.Endpoint{testMirrorName1: {Host: testMirrorHost2, Port: 8081, Phase: target.MirrorTargetPhaseConnected}},
			wantConnects: 1,
		},
		{
			name:         "deleted target",
			prev:         map[string]target.Endpoint{testMirrorName1: {Host: testMirrorHost1, Port: 8081}},
			current:      map[string]target.Endpoint{},
			wantConnects: 0,
		},
	}
	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			var (
				connects    atomic.Int64
				creates     atomic.Int64
				failDisconn atomic.Bool
				oldOpen     atomic.Bool
			)
			failDisconn.Store(true)
			oldOpen.Store(true)
			var disconnected []string
			mgr := func() k8s.Manager {
				m := k8smock.NewDefaultManagerMock().(*k8smock.ManagerMock)
				getClient := m.GetClientFunc
				m.GetClientFunc = func() k8s.Client {
					c := getClient().(*k8smock.ClientMock)
					c.CreateFunc = func(context.Context, k8s.Object, ...k8s.CreateOption) error {
						creates.Add(1)
						return nil
					}
					return c
				}
				return m
			}
			d, err := NewDiscovery(
				WithDiscoveryController(&k8smock.ControllerMock{GetManagerFunc: mgr}),
				WithDiscoveryMirror(&MirrorMock{
					ConnectFunc: func(context.Context, ...*payload.Mirror_Target) error {
						connects.Add(1)
						return nil
					},
					DisconnectFunc: func(_ context.Context, tgts ...*payload.Mirror_Target) error {
						for _, tgt := range tgts {
							disconnected = append(disconnected, tgt.GetHost())
						}
						if failDisconn.Load() {
							return errors.New("disconnect failed")
						}
						oldOpen.Store(false)
						return nil
					},
					IsConnectedFunc: func(_ context.Context, addr string) bool {
						if addr == testMirrorAddr1 {
							return oldOpen.Load()
						}
						return true
					},
					RangeMirrorAddrFunc: func(f func(addr string, _ any) bool) {
						if oldOpen.Load() {
							f(testMirrorAddr1, struct{}{})
						}
					},
				}),
			)
			if err != nil {
				tt.Fatal(err)
			}
			dis := d.(*discovery)
			ctx := tt.Context()
			dis.onReconcile(ctx, test.current)

			// 1st sync: the old endpoint's Disconnect fails.
			prev, err := dis.startSync(ctx, test.prev)
			if err == nil {
				tt.Fatal("1st sync: want the Disconnect error")
			}
			if got := connects.Load(); got != test.wantConnects {
				tt.Fatalf("1st sync: Connect calls = %d, want %d (new endpoint must be connected despite the failed Disconnect)", got, test.wantConnects)
			}

			// 2nd sync: still failing; retried, and the open old connection must
			// not be turned into an auto-created MirrorTarget.
			if prev, err = dis.startSync(ctx, prev); err == nil {
				tt.Fatal("2nd sync: want the retried Disconnect error")
			}
			if got := len(disconnected); got != 2 {
				tt.Fatalf("2nd sync: Disconnect calls = %d, want 2 (failed Disconnect must be retried)", got)
			}
			if got := creates.Load(); got != 0 {
				tt.Fatalf("2nd sync: MirrorTarget creates = %d, want 0 (stale endpoint must not be re-created)", got)
			}

			// 3rd sync: the retry succeeds.
			failDisconn.Store(false)
			if prev, err = dis.startSync(ctx, prev); err != nil {
				tt.Fatalf("3rd sync: %v", err)
			}
			if got := len(disconnected); got != 3 {
				tt.Fatalf("3rd sync: Disconnect calls = %d, want 3", got)
			}

			// 4th sync: nothing left to retry.
			if _, err = dis.startSync(ctx, prev); err != nil {
				tt.Fatalf("4th sync: %v", err)
			}
			if got := len(disconnected); got != 3 {
				tt.Fatalf("4th sync: Disconnect calls = %d, want 3 (no retry after success)", got)
			}
			for i, host := range disconnected {
				if host != testMirrorHost1 {
					tt.Errorf("Disconnect call %d: host = %q, want the old endpoint %q", i, host, testMirrorHost1)
				}
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
//
// func TestNewDiscovery(t *testing.T) {
// 	type args struct {
// 		opts []DiscoveryOption
// 	}
// 	type want struct {
// 		want Discovery
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, Discovery, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got Discovery, err error) error {
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
// 			got, err := NewDiscovery(test.args.opts...)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_discovery_onReconcile(t *testing.T) {
// 	type args struct {
// 		in0  context.Context
// 		list map[string]target.Endpoint
// 	}
// 	type fields struct {
// 		der             net.Dialer
// 		ctrl            k8s.Controller
// 		mirr            Mirror
// 		eg              errgroup.Group
// 		labels          map[string]string
// 		targetsByName   atomic.Pointer[map[string]target.Endpoint]
// 		mirrorTargets   *resource.Client[*target.MirrorTarget, *mirrv1.ValdMirrorTargetList]
// 		namespace       string
// 		colocation      string
// 		selfMirrAddrStr string
// 		selfMirrAddrs   []string
// 		dur             time.Duration
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
// 		           in0:nil,
// 		           list:nil,
// 		       },
// 		       fields: fields {
// 		           der:nil,
// 		           ctrl:nil,
// 		           mirr:nil,
// 		           eg:nil,
// 		           labels:nil,
// 		           targetsByName:nil,
// 		           mirrorTargets:nil,
// 		           namespace:"",
// 		           colocation:"",
// 		           selfMirrAddrStr:"",
// 		           selfMirrAddrs:nil,
// 		           dur:nil,
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
// 		           in0:nil,
// 		           list:nil,
// 		           },
// 		           fields: fields {
// 		           der:nil,
// 		           ctrl:nil,
// 		           mirr:nil,
// 		           eg:nil,
// 		           labels:nil,
// 		           targetsByName:nil,
// 		           mirrorTargets:nil,
// 		           namespace:"",
// 		           colocation:"",
// 		           selfMirrAddrStr:"",
// 		           selfMirrAddrs:nil,
// 		           dur:nil,
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
// 			d := &discovery{
// 				der:             test.fields.der,
// 				ctrl:            test.fields.ctrl,
// 				mirr:            test.fields.mirr,
// 				eg:              test.fields.eg,
// 				labels:          test.fields.labels,
// 				targetsByName:   test.fields.targetsByName,
// 				mirrorTargets:   test.fields.mirrorTargets,
// 				namespace:       test.fields.namespace,
// 				colocation:      test.fields.colocation,
// 				selfMirrAddrStr: test.fields.selfMirrAddrStr,
// 				selfMirrAddrs:   test.fields.selfMirrAddrs,
// 				dur:             test.fields.dur,
// 			}
//
// 			d.onReconcile(test.args.in0, test.args.list)
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_discovery_Start(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		der             net.Dialer
// 		ctrl            k8s.Controller
// 		mirr            Mirror
// 		eg              errgroup.Group
// 		labels          map[string]string
// 		targetsByName   atomic.Pointer[map[string]target.Endpoint]
// 		mirrorTargets   *resource.Client[*target.MirrorTarget, *mirrv1.ValdMirrorTargetList]
// 		namespace       string
// 		colocation      string
// 		selfMirrAddrStr string
// 		selfMirrAddrs   []string
// 		dur             time.Duration
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
// 		           der:nil,
// 		           ctrl:nil,
// 		           mirr:nil,
// 		           eg:nil,
// 		           labels:nil,
// 		           targetsByName:nil,
// 		           mirrorTargets:nil,
// 		           namespace:"",
// 		           colocation:"",
// 		           selfMirrAddrStr:"",
// 		           selfMirrAddrs:nil,
// 		           dur:nil,
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
// 		           der:nil,
// 		           ctrl:nil,
// 		           mirr:nil,
// 		           eg:nil,
// 		           labels:nil,
// 		           targetsByName:nil,
// 		           mirrorTargets:nil,
// 		           namespace:"",
// 		           colocation:"",
// 		           selfMirrAddrStr:"",
// 		           selfMirrAddrs:nil,
// 		           dur:nil,
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
// 			d := &discovery{
// 				der:             test.fields.der,
// 				ctrl:            test.fields.ctrl,
// 				mirr:            test.fields.mirr,
// 				eg:              test.fields.eg,
// 				labels:          test.fields.labels,
// 				targetsByName:   test.fields.targetsByName,
// 				mirrorTargets:   test.fields.mirrorTargets,
// 				namespace:       test.fields.namespace,
// 				colocation:      test.fields.colocation,
// 				selfMirrAddrStr: test.fields.selfMirrAddrStr,
// 				selfMirrAddrs:   test.fields.selfMirrAddrs,
// 				dur:             test.fields.dur,
// 			}
//
// 			got, err := d.Start(test.args.ctx)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_discovery_loadTargets(t *testing.T) {
// 	type fields struct {
// 		der             net.Dialer
// 		ctrl            k8s.Controller
// 		mirr            Mirror
// 		eg              errgroup.Group
// 		labels          map[string]string
// 		targetsByName   atomic.Pointer[map[string]target.Endpoint]
// 		mirrorTargets   *resource.Client[*target.MirrorTarget, *mirrv1.ValdMirrorTargetList]
// 		namespace       string
// 		colocation      string
// 		selfMirrAddrStr string
// 		selfMirrAddrs   []string
// 		dur             time.Duration
// 	}
// 	type want struct {
// 		want map[string]target.Endpoint
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, map[string]target.Endpoint) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got map[string]target.Endpoint) error {
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
// 		           der:nil,
// 		           ctrl:nil,
// 		           mirr:nil,
// 		           eg:nil,
// 		           labels:nil,
// 		           targetsByName:nil,
// 		           mirrorTargets:nil,
// 		           namespace:"",
// 		           colocation:"",
// 		           selfMirrAddrStr:"",
// 		           selfMirrAddrs:nil,
// 		           dur:nil,
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
// 		           der:nil,
// 		           ctrl:nil,
// 		           mirr:nil,
// 		           eg:nil,
// 		           labels:nil,
// 		           targetsByName:nil,
// 		           mirrorTargets:nil,
// 		           namespace:"",
// 		           colocation:"",
// 		           selfMirrAddrStr:"",
// 		           selfMirrAddrs:nil,
// 		           dur:nil,
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
// 			d := &discovery{
// 				der:             test.fields.der,
// 				ctrl:            test.fields.ctrl,
// 				mirr:            test.fields.mirr,
// 				eg:              test.fields.eg,
// 				labels:          test.fields.labels,
// 				targetsByName:   test.fields.targetsByName,
// 				mirrorTargets:   test.fields.mirrorTargets,
// 				namespace:       test.fields.namespace,
// 				colocation:      test.fields.colocation,
// 				selfMirrAddrStr: test.fields.selfMirrAddrStr,
// 				selfMirrAddrs:   test.fields.selfMirrAddrs,
// 				dur:             test.fields.dur,
// 			}
//
// 			got := d.loadTargets()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_newMirrorTarget(t *testing.T) {
// 	type args struct {
// 		host string
// 		port int
// 	}
// 	type want struct {
// 		want *payload.Mirror_Target
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, *payload.Mirror_Target) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got *payload.Mirror_Target) error {
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
// 		           host:"",
// 		           port:0,
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
// 		           host:"",
// 		           port:0,
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
// 			got := newMirrorTarget(test.args.host, test.args.port)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_discovery_connectTarget(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		req map[string]*createdTarget
// 	}
// 	type fields struct {
// 		der             net.Dialer
// 		ctrl            k8s.Controller
// 		mirr            Mirror
// 		eg              errgroup.Group
// 		labels          map[string]string
// 		targetsByName   atomic.Pointer[map[string]target.Endpoint]
// 		mirrorTargets   *resource.Client[*target.MirrorTarget, *mirrv1.ValdMirrorTargetList]
// 		namespace       string
// 		colocation      string
// 		selfMirrAddrStr string
// 		selfMirrAddrs   []string
// 		dur             time.Duration
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
// 		           req:nil,
// 		       },
// 		       fields: fields {
// 		           der:nil,
// 		           ctrl:nil,
// 		           mirr:nil,
// 		           eg:nil,
// 		           labels:nil,
// 		           targetsByName:nil,
// 		           mirrorTargets:nil,
// 		           namespace:"",
// 		           colocation:"",
// 		           selfMirrAddrStr:"",
// 		           selfMirrAddrs:nil,
// 		           dur:nil,
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
// 		           req:nil,
// 		           },
// 		           fields: fields {
// 		           der:nil,
// 		           ctrl:nil,
// 		           mirr:nil,
// 		           eg:nil,
// 		           labels:nil,
// 		           targetsByName:nil,
// 		           mirrorTargets:nil,
// 		           namespace:"",
// 		           colocation:"",
// 		           selfMirrAddrStr:"",
// 		           selfMirrAddrs:nil,
// 		           dur:nil,
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
// 			d := &discovery{
// 				der:             test.fields.der,
// 				ctrl:            test.fields.ctrl,
// 				mirr:            test.fields.mirr,
// 				eg:              test.fields.eg,
// 				labels:          test.fields.labels,
// 				targetsByName:   test.fields.targetsByName,
// 				mirrorTargets:   test.fields.mirrorTargets,
// 				namespace:       test.fields.namespace,
// 				colocation:      test.fields.colocation,
// 				selfMirrAddrStr: test.fields.selfMirrAddrStr,
// 				selfMirrAddrs:   test.fields.selfMirrAddrs,
// 				dur:             test.fields.dur,
// 			}
//
// 			err := d.connectTarget(test.args.ctx, test.args.req)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_discovery_createMirrorTargetResource(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		name string
// 		host string
// 		port int
// 	}
// 	type fields struct {
// 		der             net.Dialer
// 		ctrl            k8s.Controller
// 		mirr            Mirror
// 		eg              errgroup.Group
// 		labels          map[string]string
// 		targetsByName   atomic.Pointer[map[string]target.Endpoint]
// 		mirrorTargets   *resource.Client[*target.MirrorTarget, *mirrv1.ValdMirrorTargetList]
// 		namespace       string
// 		colocation      string
// 		selfMirrAddrStr string
// 		selfMirrAddrs   []string
// 		dur             time.Duration
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
// 		           name:"",
// 		           host:"",
// 		           port:0,
// 		       },
// 		       fields: fields {
// 		           der:nil,
// 		           ctrl:nil,
// 		           mirr:nil,
// 		           eg:nil,
// 		           labels:nil,
// 		           targetsByName:nil,
// 		           mirrorTargets:nil,
// 		           namespace:"",
// 		           colocation:"",
// 		           selfMirrAddrStr:"",
// 		           selfMirrAddrs:nil,
// 		           dur:nil,
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
// 		           name:"",
// 		           host:"",
// 		           port:0,
// 		           },
// 		           fields: fields {
// 		           der:nil,
// 		           ctrl:nil,
// 		           mirr:nil,
// 		           eg:nil,
// 		           labels:nil,
// 		           targetsByName:nil,
// 		           mirrorTargets:nil,
// 		           namespace:"",
// 		           colocation:"",
// 		           selfMirrAddrStr:"",
// 		           selfMirrAddrs:nil,
// 		           dur:nil,
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
// 			d := &discovery{
// 				der:             test.fields.der,
// 				ctrl:            test.fields.ctrl,
// 				mirr:            test.fields.mirr,
// 				eg:              test.fields.eg,
// 				labels:          test.fields.labels,
// 				targetsByName:   test.fields.targetsByName,
// 				mirrorTargets:   test.fields.mirrorTargets,
// 				namespace:       test.fields.namespace,
// 				colocation:      test.fields.colocation,
// 				selfMirrAddrStr: test.fields.selfMirrAddrStr,
// 				selfMirrAddrs:   test.fields.selfMirrAddrs,
// 				dur:             test.fields.dur,
// 			}
//
// 			err := d.createMirrorTargetResource(test.args.ctx, test.args.name, test.args.host, test.args.port)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_discovery_disconnectTarget(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		req map[string]*deletedTarget
// 	}
// 	type fields struct {
// 		der             net.Dialer
// 		ctrl            k8s.Controller
// 		mirr            Mirror
// 		eg              errgroup.Group
// 		labels          map[string]string
// 		targetsByName   atomic.Pointer[map[string]target.Endpoint]
// 		mirrorTargets   *resource.Client[*target.MirrorTarget, *mirrv1.ValdMirrorTargetList]
// 		namespace       string
// 		colocation      string
// 		selfMirrAddrStr string
// 		selfMirrAddrs   []string
// 		dur             time.Duration
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
// 		           req:nil,
// 		       },
// 		       fields: fields {
// 		           der:nil,
// 		           ctrl:nil,
// 		           mirr:nil,
// 		           eg:nil,
// 		           labels:nil,
// 		           targetsByName:nil,
// 		           mirrorTargets:nil,
// 		           namespace:"",
// 		           colocation:"",
// 		           selfMirrAddrStr:"",
// 		           selfMirrAddrs:nil,
// 		           dur:nil,
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
// 		           req:nil,
// 		           },
// 		           fields: fields {
// 		           der:nil,
// 		           ctrl:nil,
// 		           mirr:nil,
// 		           eg:nil,
// 		           labels:nil,
// 		           targetsByName:nil,
// 		           mirrorTargets:nil,
// 		           namespace:"",
// 		           colocation:"",
// 		           selfMirrAddrStr:"",
// 		           selfMirrAddrs:nil,
// 		           dur:nil,
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
// 			d := &discovery{
// 				der:             test.fields.der,
// 				ctrl:            test.fields.ctrl,
// 				mirr:            test.fields.mirr,
// 				eg:              test.fields.eg,
// 				labels:          test.fields.labels,
// 				targetsByName:   test.fields.targetsByName,
// 				mirrorTargets:   test.fields.mirrorTargets,
// 				namespace:       test.fields.namespace,
// 				colocation:      test.fields.colocation,
// 				selfMirrAddrStr: test.fields.selfMirrAddrStr,
// 				selfMirrAddrs:   test.fields.selfMirrAddrs,
// 				dur:             test.fields.dur,
// 			}
//
// 			err := d.disconnectTarget(test.args.ctx, test.args.req)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_discovery_updateMirrorTargetPhase(t *testing.T) {
// 	type args struct {
// 		ctx   context.Context
// 		name  string
// 		phase target.MirrorTargetPhase
// 	}
// 	type fields struct {
// 		der             net.Dialer
// 		ctrl            k8s.Controller
// 		mirr            Mirror
// 		eg              errgroup.Group
// 		labels          map[string]string
// 		targetsByName   atomic.Pointer[map[string]target.Endpoint]
// 		mirrorTargets   *resource.Client[*target.MirrorTarget, *mirrv1.ValdMirrorTargetList]
// 		namespace       string
// 		colocation      string
// 		selfMirrAddrStr string
// 		selfMirrAddrs   []string
// 		dur             time.Duration
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
// 		           name:"",
// 		           phase:nil,
// 		       },
// 		       fields: fields {
// 		           der:nil,
// 		           ctrl:nil,
// 		           mirr:nil,
// 		           eg:nil,
// 		           labels:nil,
// 		           targetsByName:nil,
// 		           mirrorTargets:nil,
// 		           namespace:"",
// 		           colocation:"",
// 		           selfMirrAddrStr:"",
// 		           selfMirrAddrs:nil,
// 		           dur:nil,
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
// 		           name:"",
// 		           phase:nil,
// 		           },
// 		           fields: fields {
// 		           der:nil,
// 		           ctrl:nil,
// 		           mirr:nil,
// 		           eg:nil,
// 		           labels:nil,
// 		           targetsByName:nil,
// 		           mirrorTargets:nil,
// 		           namespace:"",
// 		           colocation:"",
// 		           selfMirrAddrStr:"",
// 		           selfMirrAddrs:nil,
// 		           dur:nil,
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
// 			d := &discovery{
// 				der:             test.fields.der,
// 				ctrl:            test.fields.ctrl,
// 				mirr:            test.fields.mirr,
// 				eg:              test.fields.eg,
// 				labels:          test.fields.labels,
// 				targetsByName:   test.fields.targetsByName,
// 				mirrorTargets:   test.fields.mirrorTargets,
// 				namespace:       test.fields.namespace,
// 				colocation:      test.fields.colocation,
// 				selfMirrAddrStr: test.fields.selfMirrAddrStr,
// 				selfMirrAddrs:   test.fields.selfMirrAddrs,
// 				dur:             test.fields.dur,
// 			}
//
// 			err := d.updateMirrorTargetPhase(test.args.ctx, test.args.name, test.args.phase)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_discovery_updateTarget(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		req map[string]*updatedTarget
// 	}
// 	type fields struct {
// 		der             net.Dialer
// 		ctrl            k8s.Controller
// 		mirr            Mirror
// 		eg              errgroup.Group
// 		labels          map[string]string
// 		targetsByName   atomic.Pointer[map[string]target.Endpoint]
// 		mirrorTargets   *resource.Client[*target.MirrorTarget, *mirrv1.ValdMirrorTargetList]
// 		namespace       string
// 		colocation      string
// 		selfMirrAddrStr string
// 		selfMirrAddrs   []string
// 		dur             time.Duration
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
// 		           req:nil,
// 		       },
// 		       fields: fields {
// 		           der:nil,
// 		           ctrl:nil,
// 		           mirr:nil,
// 		           eg:nil,
// 		           labels:nil,
// 		           targetsByName:nil,
// 		           mirrorTargets:nil,
// 		           namespace:"",
// 		           colocation:"",
// 		           selfMirrAddrStr:"",
// 		           selfMirrAddrs:nil,
// 		           dur:nil,
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
// 		           req:nil,
// 		           },
// 		           fields: fields {
// 		           der:nil,
// 		           ctrl:nil,
// 		           mirr:nil,
// 		           eg:nil,
// 		           labels:nil,
// 		           targetsByName:nil,
// 		           mirrorTargets:nil,
// 		           namespace:"",
// 		           colocation:"",
// 		           selfMirrAddrStr:"",
// 		           selfMirrAddrs:nil,
// 		           dur:nil,
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
// 			d := &discovery{
// 				der:             test.fields.der,
// 				ctrl:            test.fields.ctrl,
// 				mirr:            test.fields.mirr,
// 				eg:              test.fields.eg,
// 				labels:          test.fields.labels,
// 				targetsByName:   test.fields.targetsByName,
// 				mirrorTargets:   test.fields.mirrorTargets,
// 				namespace:       test.fields.namespace,
// 				colocation:      test.fields.colocation,
// 				selfMirrAddrStr: test.fields.selfMirrAddrStr,
// 				selfMirrAddrs:   test.fields.selfMirrAddrs,
// 				dur:             test.fields.dur,
// 			}
//
// 			err := d.updateTarget(test.args.ctx, test.args.req)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_isConnectedPhase(t *testing.T) {
// 	type args struct {
// 		phase target.MirrorTargetPhase
// 	}
// 	type want struct {
// 		want bool
// 	}
// 	type test struct {
// 		name       string
// 		args       args
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
// 		           phase:nil,
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
// 		           phase:nil,
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
// 			got := isConnectedPhase(test.args.phase)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_isDisconnectedPhase(t *testing.T) {
// 	type args struct {
// 		phase target.MirrorTargetPhase
// 	}
// 	type want struct {
// 		want bool
// 	}
// 	type test struct {
// 		name       string
// 		args       args
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
// 		           phase:nil,
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
// 		           phase:nil,
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
// 			got := isDisconnectedPhase(test.args.phase)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }

// doneCountCtx counts Done() calls. A select loop evaluates ctx.Done() on every
// iteration, so the count is a deterministic measure of how often it spins.
type doneCountCtx struct {
	context.Context //nolint:containedctx // test wrapper that counts Done() calls
	n               atomic.Int64
}

func (c *doneCountCtx) Done() <-chan struct{} {
	c.n.Add(1)
	return c.Context.Done()
}

// TestDiscovery_Start_ClosedControllerChannelDoesNotSpin pins that the relay
// loop stops selecting the controller error channel once it is closed (the
// controller's manager returned while ctx is still live). A closed channel is
// always ready, so without nil-ing it the loop busy-spins; the spin shows up
// as an unbounded number of ctx.Done() evaluations.
func TestDiscovery_Start_ClosedControllerChannelDoesNotSpin(t *testing.T) {
	t.Parallel()

	ch := make(chan error)
	close(ch)

	eg, egctx := errgroup.New(t.Context())
	cctx, cancel := context.WithCancel(egctx)
	defer cancel()
	ctx := &doneCountCtx{Context: cctx}

	d := &discovery{
		eg: eg,
		// the sync tick must not fire during the test
		dur: time.Hour,
		ctrl: &k8smock.ControllerMock{
			StartFunc: func(context.Context) (<-chan error, error) { return ch, nil },
		},
	}
	if _, err := d.Start(ctx); err != nil {
		t.Fatalf("Start() error = %v, want nil", err)
	}

	time.Sleep(100 * time.Millisecond)
	if n := ctx.n.Load(); n > 100 {
		t.Errorf("ctx.Done() evaluated %d times with a closed controller channel; the relay loop is busy-spinning", n)
	}

	cancel()
	done := make(chan struct{})
	go func() {
		defer close(done)
		if err := eg.Wait(); err != nil {
			t.Logf("eg.Wait() = %v", err)
		}
	}()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("relay goroutine did not exit on context cancel")
	}
}

// TestDiscovery_Start_ErrorRelayStopsOnCancelWhenConsumerGone pins that the
// relay goroutine returns when ctx is cancelled while it is blocked sending a
// controller error to a consumer that stopped reading.
func TestDiscovery_Start_ErrorRelayStopsOnCancelWhenConsumerGone(t *testing.T) {
	t.Parallel()

	const sendTimeout = 5 * time.Second

	dech := make(chan error)
	eg, egctx := errgroup.New(t.Context())
	ctx, cancel := context.WithCancel(egctx)
	defer cancel()

	d := &discovery{
		eg:  eg,
		dur: time.Hour,
		ctrl: &k8smock.ControllerMock{
			StartFunc: func(context.Context) (<-chan error, error) { return dech, nil },
		},
	}
	// the returned channel is intentionally never read.
	if _, err := d.Start(ctx); err != nil {
		t.Fatalf("Start() error = %v, want nil", err)
	}

	// one more error than the channel buffer: the last one cannot be buffered.
	for range errChanBufferSize + 1 {
		select {
		case dech <- errors.New("controller error"):
		case <-time.After(sendTimeout):
			t.Fatal("relay goroutine did not receive the controller error")
		}
	}

	cancel()
	done := make(chan struct{})
	go func() {
		defer close(done)
		if err := eg.Wait(); err != nil {
			t.Logf("eg.Wait() = %v", err)
		}
	}()
	select {
	case <-done:
	case <-time.After(sendTimeout):
		t.Fatal("relay goroutine stayed blocked on the error channel after ctx cancel")
	}
}
