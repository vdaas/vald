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
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
)

// An agent that rejects the search (wrong dimension: InvalidArgument; indexing or
// flushing: Aborted) used to be treated as a miss, so when no agent had results
// the client got NotFound instead of the agents' code. Results from the other
// agents are still served.
func Test_server_doSearch_agentErrorCodes(t *testing.T) {
	t.Parallel()
	agentA, agentB := new(regClient), new(regClient)
	hit := &payload.Search_Response{
		Results: []*payload.Object_Distance{{Id: "hit", Distance: 0.1}},
	}
	tests := []struct {
		errA     error
		errB     error
		resB     *payload.Search_Response
		name     string
		wantCode codes.Code // codes.OK: a result is expected
	}{
		{
			name:     "all agents InvalidArgument",
			errA:     status.Error(codes.InvalidArgument, "Search API Incompatible Dimension Size detected"),
			errB:     status.Error(codes.InvalidArgument, "Search API Incompatible Dimension Size detected"),
			wantCode: codes.InvalidArgument,
		},
		{
			name:     "all agents Aborted",
			errA:     status.Error(codes.Aborted, "indexing"),
			errB:     status.Error(codes.Aborted, "flushing"),
			wantCode: codes.Aborted,
		},
		{
			name:     "all agents NotFound stays NotFound",
			errA:     status.Error(codes.NotFound, "nf"),
			errB:     status.Error(codes.NotFound, "nf"),
			wantCode: codes.NotFound,
		},
		{
			name:     "Aborted agent does not hide the other agent's results",
			errA:     status.Error(codes.Aborted, "indexing"),
			resB:     hit,
			wantCode: codes.OK,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			gw := &regGateway{
				agents: 2,
				broadCast: func(ctx context.Context, f func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error) error {
					// models service.Gateway.BroadCast: every callback error is returned.
					return errors.Join(f(ctx, "a:8081", agentA), f(ctx, "b:8081", agentB))
				},
			}
			s := &server{gateway: gw, replica: 1, timeout: 10 * time.Second}
			res, _, err := s.doSearch(t.Context(), vald.SearchRPCName, &payload.Search_Config{Num: 1},
				func(_ context.Context, _ *payload.Search_Config, vc vald.Client, _ ...grpc.CallOption) (*payload.Search_Response, error) {
					if vc == agentA {
						return nil, test.errA
					}
					return test.resB, test.errB
				})
			if test.wantCode == codes.OK {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if got := len(res.GetResults()); got != 1 {
					t.Fatalf("got %d results, want 1", got)
				}
				return
			}
			st, ok := status.FromError(err)
			if !ok || st == nil || st.Code() != test.wantCode {
				t.Fatalf("got error: %v, want code %v", err, test.wantCode)
			}
		})
	}
}
