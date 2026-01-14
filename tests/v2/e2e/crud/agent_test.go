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
	"github.com/vdaas/vald/tests/v2/e2e/config"
)

func (r *runner) processAgent(t *testing.T, ctx context.Context, plan *config.Execution) error {
	t.Helper()
	if plan == nil {
		t.Fatalf("index operation plan is nil")
		return nil
	}
	switch plan.Type {
	case config.OpCreateIndex:
		return single(t, ctx, 0, plan, &payload.Control_CreateIndexRequest{
			PoolSize: plan.Agent.PoolSize,
		}, r.aclient.CreateIndex, emptyCallback[*payload.Empty](plan.Name))
	case config.OpSaveIndex:
		return single(t, ctx, 0, plan, new(payload.Empty), r.aclient.SaveIndex, emptyCallback[*payload.Empty](plan.Name))
	case config.OpCreateAndSaveIndex:
		return single(t, ctx, 0, plan, &payload.Control_CreateIndexRequest{
			PoolSize: plan.Agent.PoolSize,
		}, r.aclient.CreateAndSaveIndex, emptyCallback[*payload.Empty](plan.Name))
	default:
		t.Fatalf("unsupported agent operation: %s", plan.Type)
	}
	return nil
}
