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
package grpc

import (
	"context"
	"testing"
	"time"
)

// TestGRPCClient_StartConnectionMonitor_initialConnectionFailure verifies that a failed
// initial connection does not leave the monitor marked as running, so that
// following calls retry and report the error instead of returning (nil, nil).
func TestGRPCClient_StartConnectionMonitor_initialConnectionFailure(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		calls int
	}{
		{name: "first call fails", calls: 1},
		{name: "every retry reports the failure again", calls: 3},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
			defer cancel()
			// an empty address is skipped on the initial connection, so no connection is
			// established (no dial, no port scan) and the monitor start fails.
			c := New("monitor-failure-test", WithAddrs(""))
			g, ok := c.(*gRPCClient)
			if !ok {
				t.Fatalf("unexpected client type: %T", c)
			}
			for i := range test.calls {
				ech, err := g.StartConnectionMonitor(ctx)
				if err == nil {
					t.Fatalf("call %d: StartConnectionMonitor must return an error, got ech: %v", i, ech)
				}
				if ech != nil {
					t.Errorf("call %d: error channel must be nil on error", i)
				}
				if g.monitorRunning.Load() {
					t.Errorf("call %d: monitorRunning must be reset after the failure", i)
				}
			}
		})
	}
}

// TestGRPCClient_StartConnectionMonitor_alreadyRunning verifies that a call while the monitor
// is running does not start another one and returns the existing channel.
func TestGRPCClient_StartConnectionMonitor_alreadyRunning(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
	defer cancel()
	g := &gRPCClient{}
	ech := make(chan error)
	g.ech = ech
	g.monitorRunning.Store(true)

	got, err := g.StartConnectionMonitor(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if got != (<-chan error)(ech) {
		t.Errorf("got channel %v, want %v", got, ech)
	}
	if !g.monitorRunning.Load() {
		t.Error("monitorRunning must stay true")
	}
}

// TestGRPCClient_StartConnectionMonitor_concurrentFailure verifies that concurrent callers
// never observe (nil, nil) while the initial connection of the starting caller fails.
func TestGRPCClient_StartConnectionMonitor_concurrentFailure(t *testing.T) {
	t.Parallel()
	const callers = 8
	ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
	defer cancel()
	c := New("monitor-concurrent-failure-test", WithAddrs(""))
	g, ok := c.(*gRPCClient)
	if !ok {
		t.Fatalf("unexpected client type: %T", c)
	}
	type result struct {
		ech <-chan error
		err error
	}
	results := make(chan result, callers)
	for range callers {
		go func() {
			ech, err := g.StartConnectionMonitor(ctx)
			results <- result{ech: ech, err: err}
		}()
	}
	for range callers {
		select {
		case r := <-results:
			if r.err == nil {
				t.Errorf("StartConnectionMonitor must return an error, got ech: %v", r.ech)
			}
			if r.ech != nil {
				t.Error("error channel must be nil on error")
			}
		case <-ctx.Done():
			t.Fatal("timed out waiting for StartConnectionMonitor callers")
		}
	}
	if g.monitorRunning.Load() {
		t.Error("monitorRunning must be reset after the failure")
	}
}
