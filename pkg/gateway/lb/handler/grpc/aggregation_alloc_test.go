//go:build !race

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

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/net/grpc"
)

func Test_valdPairingHeapAggr_Allocations(t *testing.T) {
	ctx := t.Context()
	const n = 10
	resp1 := batchResponse(n)
	resp2 := batchResponse(n)

	// Warm pools.
	for range 5 {
		ag := newPairingHeap(n, n, 3)
		ag.Send(ctx, resp1)
		ag.Send(ctx, resp2)
		_ = ag.Result()
	}

	if allocs := testing.AllocsPerRun(50, func() {
		ag := newPairingHeap(n, n, 3)
		ag.Send(ctx, resp1)
		ag.Send(ctx, resp2)
		if got := len(ag.Result().GetResults()); got != n {
			t.Fatalf("got %d results, want %d", got, n)
		}
	}); allocs > 3 {
		t.Errorf("valdPairingHeapAggr round-trip allocated %v times, want <= 3", allocs)
	}

	ag := newPairingHeap(n, n, 1)
	ag.Send(ctx, resp1)
	if allocs := testing.AllocsPerRun(50, func() {
		ag.Send(ctx, resp1)
	}); allocs != 0 {
		t.Errorf("valdPairingHeapAggr duplicate Send allocated %v times, want 0", allocs)
	}
	_ = ag.Result()
}

func Test_valdPoolSliceAggr_Allocations(t *testing.T) {
	ctx := t.Context()
	const n = 10
	resp := batchResponse(n)

	// Warm pools.
	for range 5 {
		ag := newPoolSlice(n, n, 3)
		ag.Send(ctx, resp)
		_ = ag.Result()
	}

	if allocs := testing.AllocsPerRun(50, func() {
		ag := newPoolSlice(n, n, 3)
		ag.Send(ctx, resp)
		if got := len(ag.Result().GetResults()); got != n {
			t.Fatalf("got %d results, want %d", got, n)
		}
	}); allocs > 3 {
		t.Errorf("valdPoolSliceAggr round-trip allocated %v times, want <= 3", allocs)
	}
}

func Test_aggregationSearch_HealthyAllocs(t *testing.T) {
	ctx := t.Context()
	hit := newDistResponse(dist32("ok", 0.1))
	cli := &regClient{
		search: func(context.Context, *payload.Search_Request) (*payload.Search_Response, error) {
			return hit, nil
		},
	}
	gw := &regGateway{
		agents: 1,
		broadCast: func(ctx context.Context, f func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error) error {
			return f(ctx, "a0", cli)
		},
	}
	s := &server{gateway: gw, replica: 1, timeout: time.Second}
	vec := []float32{1, 2, 3}
	cfg := &payload.Search_Config{Num: 1}

	buildReq := func(fcfg *payload.Search_Config) *payload.Search_Request {
		return &payload.Search_Request{Vector: vec, Config: fcfg}
	}
	for range 5 {
		if _, _, err := doSearch(ctx, s, vald.SearchRPCName, cfg, buildReq, vald.Client.Search); err != nil {
			t.Fatalf("warm doSearch error = %v", err)
		}
	}

	allocs := testing.AllocsPerRun(50, func() {
		res, _, err := doSearch(ctx, s, vald.SearchRPCName, cfg, buildReq, vald.Client.Search)
		if err != nil || len(res.GetResults()) != 1 {
			t.Fatalf("doSearch = (%v, %v)", res, err)
		}
	})
	if allocs > 14 {
		t.Errorf("healthy doSearch allocated %v times, want <= 14", allocs)
	}
}
