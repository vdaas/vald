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
package service

import (
	"context"
	"slices"
	"testing"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/strings"
	clientmock "github.com/vdaas/vald/internal/test/mock/client"
	grpcmock "github.com/vdaas/vald/internal/test/mock/grpc"
)

const testAgentAddr = "127.0.0.1:8080"

// Test_index_doDeleteIndex_retriedTarget pins that only the last attempt per target is
// reported: the client backoff re-runs the callback for the same target, and a
// retry that succeeds must clear the earlier failure.
func Test_index_doDeleteIndex_retriedTarget(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		wantSubstr string
		notSubstr  string
		attempts   []error
		wantErr    bool
	}{
		{
			name:     "a failure followed by a successful retry is not an error",
			attempts: []error{status.Error(codes.Unavailable, "first failure"), nil},
		},
		{
			name:     "a failure followed by a skipped retry is not an error",
			attempts: []error{status.Error(codes.Unavailable, "first failure"), status.Error(codes.FailedPrecondition, "skip")},
		},
		{
			name:     "NotFound of a retried Remove that already succeeded is not an error",
			attempts: []error{status.Error(codes.Unavailable, "first failure"), status.Error(codes.NotFound, "not found")},
		},
		{
			name:       "only the last failure of a target is reported",
			attempts:   []error{status.Error(codes.Unavailable, "first failure"), status.Error(codes.Internal, "final failure")},
			wantErr:    true,
			wantSubstr: "final failure",
			notSubstr:  "first failure",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(tt *testing.T) {
			tt.Parallel()
			idx := &index{
				client: &clientmock.DiscovererClientMock{
					GetAddrsFunc: func(_ context.Context) []string {
						return []string{testAgentAddr}
					},
					GetClientFunc: func() grpc.Client {
						return &grpcmock.GRPCClientMock{
							// re-run a failing callback like the client backoff does and swallow
							// the final error like the real OrderedRangeConcurrent.
							OrderedRangeConcurrentFunc: func(ctx context.Context, order []string, _ int,
								f func(_ context.Context, _ string, _ *grpc.ClientConn, _ ...grpc.CallOption) error,
							) error {
								for _, addr := range order {
									for range tc.attempts {
										if f(ctx, addr, nil) == nil {
											break
										}
									}
								}
								return nil
							},
						}
					},
				},
			}
			var calls int
			err := idx.doDeleteIndex(tt.Context(),
				func(context.Context, vald.RemoveClient, ...grpc.CallOption) (*payload.Object_Location, error) {
					if calls >= len(tc.attempts) {
						return nil, errors.New("unexpected extra attempt")
					}
					aerr := tc.attempts[calls]
					calls++
					return nil, aerr
				})
			if !tc.wantErr {
				if err != nil {
					tt.Errorf("doDeleteIndex() error = %v, want nil", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tc.wantSubstr) {
				tt.Fatalf("doDeleteIndex() error = %v, want it to contain %q", err, tc.wantSubstr)
			}
			if strings.Contains(err.Error(), tc.notSubstr) {
				tt.Errorf("doDeleteIndex() error = %v, must not contain the retried failure %q", err, tc.notSubstr)
			}
		})
	}
}

// Test_index_doDeleteIndex_unprocessedTarget pins that a target agent the client never ran the
// callback for (no connection pool, connection lookup failure, open circuit breaker, done
// ctx) fails the job: OrderedRangeConcurrent swallows those errors and returns nil.
func Test_index_doDeleteIndex_unprocessedTarget(t *testing.T) {
	t.Parallel()
	targets := []string{"127.0.0.1:8080", "127.0.0.1:8081", "127.0.0.1:8082"}
	tests := []struct {
		name      string
		result    error
		reached   []string
		unreached []string
	}{
		{
			name:    "every target processed is not an error",
			reached: targets,
		},
		{
			name:    "a target skipped with NotFound is processed",
			result:  status.Error(codes.NotFound, "not found"),
			reached: targets,
		},
		{
			name:      "a target the callback never ran for is an error",
			reached:   []string{targets[0], targets[2]},
			unreached: []string{targets[1]},
		},
		{
			name:      "no target reached is an error for every target",
			unreached: targets,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(tt *testing.T) {
			tt.Parallel()
			idx := &index{
				client: &clientmock.DiscovererClientMock{
					GetAddrsFunc: func(_ context.Context) []string {
						return targets
					},
					GetClientFunc: func() grpc.Client {
						return &grpcmock.GRPCClientMock{
							// skip the unreached targets and swallow that like the real OrderedRangeConcurrent.
							OrderedRangeConcurrentFunc: func(ctx context.Context, order []string, _ int,
								f func(_ context.Context, _ string, _ *grpc.ClientConn, _ ...grpc.CallOption) error,
							) error {
								for _, addr := range order {
									if slices.Contains(tc.reached, addr) {
										if err := f(ctx, addr, nil); err != nil {
											tt.Logf("callback for %s returned %v", addr, err)
										}
									}
								}
								return nil
							},
						}
					},
				},
			}
			err := idx.doDeleteIndex(tt.Context(),
				func(context.Context, vald.RemoveClient, ...grpc.CallOption) (*payload.Object_Location, error) {
					return nil, tc.result
				})
			if len(tc.unreached) == 0 {
				if err != nil {
					tt.Errorf("doDeleteIndex() error = %v, want nil", err)
				}
				return
			}
			if err == nil {
				tt.Fatalf("doDeleteIndex() error = nil, want an error for %v", tc.unreached)
			}
			for _, addr := range targets {
				want := slices.Contains(tc.unreached, addr)
				if got := errors.Is(err, errors.ErrGRPCClientConnNotFound(addr)); got != want {
					tt.Errorf("doDeleteIndex() error = %v, reports %s as unprocessed: %t, want %t", err, addr, got, want)
				}
			}
		})
	}
}
