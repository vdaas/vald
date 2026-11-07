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
	"testing"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/k8s"
	"github.com/vdaas/vald/internal/k8s/vald/mirror/target"
	"github.com/vdaas/vald/internal/sync/atomic"
	k8smock "github.com/vdaas/vald/internal/test/mock/k8s"
)

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
