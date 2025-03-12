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

	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc/proto"
	"github.com/vdaas/vald/tests/v2/e2e/config"
)

// unaryModify performs unary modification operations.
// It sends a single request to the server and checks the response.
func (r *runner) processIndex(t *testing.T, ctx context.Context, plan *config.Execution) {
	t.Helper()
	if plan == nil {
		t.Fatalf("index operation plan is nil")
		return
	}
	switch plan.Type {
	case config.OpIndexInfo:
		unaryIndex(t, ctx, plan, r.client.IndexInfo)
	case config.OpIndexDetail:
		unaryIndex(t, ctx, plan, r.client.IndexDetail)
	case config.OpIndexStatistics:
		unaryIndex(t, ctx, plan, r.client.IndexStatistics)
	case config.OpIndexStatisticsDetail:
		unaryIndex(t, ctx, plan, r.client.IndexStatisticsDetail)
	case config.OpIndexProperty:
		unaryIndex(t, ctx, plan, r.client.IndexProperty)
	case config.OpFlush:
		unaryIndex(t, ctx, plan, r.client.Flush)
	default:
		t.Fatalf("unsupported index operation: %s", plan.Type)
	}
}

// unaryIndex handles unary modification requests. It iterates over the data and for each vector,
// it sends a modification request with the specified operation and configuration.
// The function logs the result of each modification.
// The function is used for insert, update, upsert, and remove operations.
func unaryIndex[Q, R proto.Message](
	t *testing.T, ctx context.Context, plan *config.Execution, call grpcCall[Q, R],
) {
	t.Helper()
	// Launch the index modify request in a goroutine.
	var req Q
	// Execute the modify gRPC call.
	res, err := call(ctx, req)
	if err != nil {
		// Handle the error using the centralized error handler.
		handleGRPCCallError(t, err, plan)
		return
	}
	log.Infof("Index operation %s: returned %v", plan.Type, res)
	return
}
