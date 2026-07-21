//go:build e2e

//
// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

// package crud provides end-to-end tests using ann-benchmarks datasets.
package crud

import (
	"context"
	"testing"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/rpc/stats"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/tests/v2/e2e/config"
)

func (r *runner) resourceStatsDetail(
	ctx context.Context, in *payload.Empty, opts ...grpc.CallOption,
) (*payload.Info_Stats_ResourceStatsDetail, error) {
	return grpc.RoundRobin(ctx, r.client.GRPCClient(), func(
		ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (*payload.Info_Stats_ResourceStatsDetail, error) {
		return stats.NewStatsDetailClient(conn).ResourceStatsDetail(ctx, in, append(copts, opts...)...)
	})
}

func (r *runner) processStats(t *testing.T, ctx context.Context, plan *config.Execution) error {
	t.Helper()
	if plan == nil {
		t.Fatalf("stats operation plan is nil")
		return nil
	}
	switch plan.Type {
	case config.OpResourceStatsDetail:
		return single(t, ctx, 0, plan, new(payload.Empty), r.resourceStatsDetail, printCallback[*payload.Info_Stats_ResourceStatsDetail](passThrough))
	default:
		t.Fatalf("unsupported stats operation: %s", plan.Type)
	}
	return nil
}
