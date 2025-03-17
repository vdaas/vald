//go:build e2e

//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
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
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc/proto"
	"github.com/vdaas/vald/tests/v2/e2e/config"
)

func indexCallBack[R proto.Message](res R, err error) bool {
	if err != nil {
		return true
	}
	log.Infof("response: %v", res)
	return true
}

func (r *runner) processIndex(t *testing.T, ctx context.Context, plan *config.Execution) {
	t.Helper()
	if plan == nil {
		t.Fatalf("index operation plan is nil")
		return
	}
	switch plan.Type {
	case config.OpIndexInfo:
		single(t, ctx, 0, plan, new(payload.Empty), r.client.IndexInfo)
	case config.OpIndexDetail:
		single(t, ctx, 0, plan, new(payload.Empty), r.client.IndexDetail)
	case config.OpIndexStatistics:
		single(t, ctx, 0, plan, new(payload.Empty), r.client.IndexStatistics)
	case config.OpIndexStatisticsDetail:
		single(t, ctx, 0, plan, new(payload.Empty), r.client.IndexStatisticsDetail)
	case config.OpIndexProperty:
		single(t, ctx, 0, plan, new(payload.Empty), r.client.IndexProperty)
	case config.OpFlush:
		single(t, ctx, 0, plan, new(payload.Flush_Request), r.client.Flush)
	default:
		t.Fatalf("unsupported index operation: %s", plan.Type)
	}
}
