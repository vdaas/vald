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
package circuitbreaker

import (
	"context"
	"slices"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/sync"
)

// Do counts every call per breaker name and state; the counters are lock-free
// per key, so concurrent calls must still be counted exactly.
func TestBreakerManager_Do_countsMetrics(t *testing.T) {
	t.Parallel()
	cb, err := New()
	if err != nil {
		t.Fatal(err)
	}
	const (
		key   = "TestBreakerManager_Do_countsMetrics"
		calls = 200
	)
	var wg sync.WaitGroup
	for range calls {
		wg.Go(func() {
			if _, err := cb.Do(t.Context(), key, func(context.Context) (any, error) {
				return struct{}{}, nil
			}); err != nil {
				t.Error(err)
			}
		})
	}
	wg.Wait()
	if got := Metrics(t.Context())[key][StateClosed]; got != calls {
		t.Errorf("Metrics()[%q][closed] = %d, want %d", key, got, calls)
	}
}

// EvictAddr removes the breakers and counters of every key that targets the
// address ("<addr>" or "<method>/<addr>") and keeps every other key, including
// keys that only end with the same characters.
func TestBreakerManager_EvictAddr(t *testing.T) {
	t.Parallel()
	const (
		method = "TestBreakerManager_EvictAddr"
		addr   = "10.0.0.1:8081"
	)
	evicted := []string{addr, method + "/" + addr, method + "/192.0.2.1:8081/" + addr}
	kept := []string{
		method + "/1" + addr,
		method + "/" + addr + "/192.0.2.1:8081",
		method + "/192.0.2.1:8081",
	}
	cb, err := New()
	if err != nil {
		t.Fatal(err)
	}
	bm, ok := cb.(*breakerManager)
	if !ok {
		t.Fatalf("New returned %T, want *breakerManager", cb)
	}
	for _, key := range append(slices.Clone(evicted), kept...) {
		if _, err := cb.Do(t.Context(), key, func(context.Context) (any, error) {
			return true, nil
		}); err != nil {
			t.Fatal(err)
		}
	}

	bm.EvictAddr(addr)

	ms := Metrics(t.Context())
	for _, key := range evicted {
		if _, ok := bm.m.Load(key); ok {
			t.Errorf("breaker %q remains after EvictAddr", key)
		}
		if _, ok := ms[key]; ok {
			t.Errorf("metrics for %q remain after EvictAddr", key)
		}
	}
	for _, key := range kept {
		if _, ok := bm.m.Load(key); !ok {
			t.Errorf("breaker %q was evicted", key)
		}
		if _, ok := ms[key]; !ok {
			t.Errorf("metrics for %q were evicted", key)
		}
	}
}

// A call still in flight when EvictAddr runs must not re-create the evicted
// key's counters when it finishes; otherwise the state of removed addresses
// keeps piling up under constant traffic.
func TestBreakerManager_EvictAddr_inFlightCall(t *testing.T) {
	t.Parallel()
	const addr = "10.0.0.2:8081"
	key := "TestBreakerManager_EvictAddr_inFlightCall/" + addr
	cb, err := New()
	if err != nil {
		t.Fatal(err)
	}
	bm, ok := cb.(*breakerManager)
	if !ok {
		t.Fatalf("New returned %T, want *breakerManager", cb)
	}
	started, release := make(chan struct{}), make(chan struct{})
	done := make(chan error, 1)
	go func() {
		_, err := cb.Do(t.Context(), key, func(context.Context) (any, error) {
			close(started)
			<-release
			return true, nil
		})
		done <- err
	}()
	select {
	case <-started:
	case <-time.After(10 * time.Second):
		t.Fatal("call did not start")
	}
	bm.EvictAddr(addr)
	close(release)
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("call did not finish")
	}
	if _, ok := Metrics(t.Context())[key]; ok {
		t.Errorf("metrics for %q were re-created by an in-flight call after EvictAddr", key)
	}
	if _, ok := bm.m.Load(key); ok {
		t.Errorf("breaker %q remains after EvictAddr", key)
	}
}
