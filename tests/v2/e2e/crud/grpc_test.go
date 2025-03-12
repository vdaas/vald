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

// Package crud provides end-to-end tests using ann-benchmarks datasets.
package crud

import (
	"context"
	"testing"

	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/proto"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/tests/v2/e2e/config"
)

// handleGRPCCallError centralizes the gRPC error handling and logging.
// It compares the error's status code with the expected codes from the plan.
// If the error is expected, it logs a message; otherwise, it logs an error.
func handleGRPCCallError(t *testing.T, err error, plan *config.Execution) {
	t.Helper()
	if err != nil {
		if st, ok := status.FromError(err); ok && st != nil {
			if len(plan.ExpectedStatusCodes) != 0 && !plan.ExpectedStatusCodes.Equals(st.Code().String()) {
				t.Errorf("unexpected error: %v", st)
			}
			return
		}
		t.Errorf("failed to search vector: %v", err)
	}
}

// Type aliases for generic search functions.
type (
	// grpcCall is a generic function type for making gRPC calls.
	grpcCall[Q, R proto.Message] func(ctx context.Context, query Q, opts ...grpc.CallOption) (response R, err error)
	// newMultiRequest is a generic type for functions that build bulk search requests.
	newMultiRequest[R, S proto.Message] func([]R) S
	// newStream is a generic type for functions that create a new gRPC stream.
	newStream[S grpc.ClientStream] func(ctx context.Context, opts ...grpc.CallOption) (S, error)
)
