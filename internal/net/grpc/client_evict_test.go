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
	"slices"
	"testing"

	"github.com/vdaas/vald/internal/circuitbreaker"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/sync"
)

// evictRecordingDialer is a net.Dialer fake that records EvictAddr calls. The
// embedded Dialer is nil: these tests never dial.
type evictRecordingDialer struct {
	net.Dialer
	evicted []string
	mu      sync.Mutex
}

func (d *evictRecordingDialer) EvictAddr(addr string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.evicted = append(d.evicted, addr)
}

// Test_gRPCClient_Disconnect_evictsAddrState verifies that Disconnect drops the
// circuit breaker metrics and the dialer state of the removed address, and of
// no other address. They used to grow with every address ever connected.
func Test_gRPCClient_Disconnect_evictsAddrState(t *testing.T) {
	t.Parallel()
	const (
		method = "Test_gRPCClient_Disconnect_evictsAddrState"
		addr   = "192.0.2.166:8081"
		other  = "192.0.2.167:8081"
	)
	tests := []struct {
		name        string
		wantEvicted []string
		wantGone    []string
		wantKept    []string
		stored      bool
	}{
		{
			name:        "connected address",
			stored:      true,
			wantEvicted: []string{addr},
			wantGone:    []string{method + "/" + addr, addr},
			wantKept:    []string{method + "/" + other},
		},
		{
			name:     "unknown address",
			stored:   false,
			wantKept: []string{method + "/" + addr, addr, method + "/" + other},
		},
	}
	// Subtests run in sequence: the breaker metrics are package-global.
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := t.Context()
			cb, err := circuitbreaker.New()
			if err != nil {
				t.Fatal(err)
			}
			der := new(evictRecordingDialer)
			g, ok := New("evict-test", WithCircuitBreaker(cb), WithDialer("tcp", der)).(*gRPCClient)
			if !ok {
				t.Fatal("unexpected client type")
			}
			for _, key := range []string{method + "/" + addr, addr, method + "/" + other} {
				if _, derr := cb.Do(ctx, key, func(context.Context) (any, error) {
					return true, nil
				}); derr != nil {
					t.Fatal(derr)
				}
			}
			if test.stored {
				g.conns.Store(addr, &hookConn{})
				g.clientCount = 1
			}

			err = g.Disconnect(ctx, addr)
			if (err == nil) != test.stored {
				t.Fatalf("Disconnect returned error %v for stored=%v", err, test.stored)
			}

			ms := circuitbreaker.Metrics(ctx)
			for _, key := range test.wantGone {
				if _, ok := ms[key]; ok {
					t.Errorf("breaker metrics for %q remain after Disconnect", key)
				}
			}
			for _, key := range test.wantKept {
				if _, ok := ms[key]; !ok {
					t.Errorf("breaker metrics for %q were removed", key)
				}
			}
			der.mu.Lock()
			evicted := slices.Clone(der.evicted)
			der.mu.Unlock()
			// Disconnect evicts before and after closing the pool, so the same
			// address may be evicted more than once.
			slices.Sort(evicted)
			evicted = slices.Compact(evicted)
			if !slices.Equal(evicted, test.wantEvicted) {
				t.Errorf("dialer evicted %v, want %v", evicted, test.wantEvicted)
			}
		})
	}
}
