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

package vald

import (
	"strconv"
	"testing"

	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/credentials"
	"github.com/vdaas/vald/internal/sync"
)

func TestConnCache_NilAndCaching(t *testing.T) {
	t.Parallel()

	var cache ConnCache
	if got := cache.Get(nil); got == nil {
		t.Fatal("Get(nil) returned nil")
	} else if _, err := got.Search(t.Context(), nil); !errors.Is(err, errors.ErrGRPCClientConnNotFound("*")) {
		t.Fatalf("Get(nil).Search() error = %v, want %v", err, errors.ErrGRPCClientConnNotFound("*"))
	}
	if _, ok := cache.Load(nil); ok {
		t.Fatal("Load(nil) reported ok=true")
	}

	conn := new(grpc.ClientConn)
	c1 := cache.Get(conn)
	c2 := cache.Get(conn)
	if c1 == nil || c2 == nil {
		t.Fatal("expected non-nil cached client")
	}
	if c1 != c2 {
		t.Fatalf("expected cached instance reuse, got %p != %p", c1, c2)
	}
	if got, ok := cache.Load(conn); !ok || got != c1 {
		t.Fatalf("Load(conn) = (%v, %v), want (%p, true)", got, ok, c1)
	}
}

func TestConnCache_ConcurrentAndPrune(t *testing.T) {
	t.Parallel()

	var cache ConnCache
	conns := make([]*grpc.ClientConn, 12)
	for i := range conns {
		cc, err := grpc.NewClient("passthrough:///conn-cache-test-"+strconv.Itoa(i),
			grpc.WithTransportCredentials(credentials.NewInsecure()))
		if err != nil {
			t.Fatal(err)
		}
		conns[i] = cc
	}
	t.Cleanup(func() {
		for _, cc := range conns[1:] {
			if err := cc.Close(); err != nil {
				t.Errorf("close conn: %v", err)
			}
		}
	})

	const workers = 8
	var wg sync.WaitGroup
	results := make([]vald.Client, len(conns))
	for i, cc := range conns {
		results[i] = cache.Get(cc)
	}
	for range workers {
		wg.Go(func() {
			for i, cc := range conns {
				if got := cache.Get(cc); got != results[i] {
					t.Errorf("concurrent Get(%d) returned different instance", i)
				}
			}
		})
	}
	wg.Wait()

	if err := conns[0].Close(); err != nil {
		t.Fatal(err)
	}
	extra, err := grpc.NewClient("passthrough:///conn-cache-test-extra",
		grpc.WithTransportCredentials(credentials.NewInsecure()))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := extra.Close(); err != nil {
			t.Errorf("close extra conn: %v", err)
		}
	})

	_ = cache.Get(extra)
	if _, ok := cache.Load(extra); !ok {
		t.Fatal("expected extra conn in cache")
	}
	if _, ok := cache.Load(conns[0]); ok {
		t.Fatal("expected Shutdown conn to be pruned")
	}
}
