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
	"testing"

	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/client/v1/client/discoverer"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
)

// rangeClient models gRPCClient.RangeConcurrent: f runs per target, up to
// attempts times while it fails (the client backoff), and every per-target
// error that is not a local ctx error is dropped from the return value.
type rangeClient struct {
	grpc.Client
	addrs    []string
	attempts int
}

func (c *rangeClient) RangeConcurrent(
	ctx context.Context,
	_ int,
	f func(ctx context.Context, addr string, conn *grpc.ClientConn, copts ...grpc.CallOption) error,
) error {
	if len(c.addrs) == 0 {
		return errors.ErrGRPCClientConnNotFound("*")
	}
	for _, addr := range c.addrs {
		for range max(c.attempts, 1) {
			if f(ctx, addr, nil) == nil {
				break
			}
		}
	}
	return nil
}

type rangeDiscoverer struct {
	discoverer.Client
	c grpc.Client
}

func (d *rangeDiscoverer) GetClient() grpc.Client     { return d.c }
func (d *rangeDiscoverer) GetReadClient() grpc.Client { return d.c }

// BroadCast returned whatever RangeConcurrent returned, so an agent failure
// (Unavailable/Internal/Aborted/...) reached the lb handlers as nil and they
// answered NotFound or OK instead of the agent's code.
func Test_gateway_BroadCast_propagatesCallbackErrors(t *testing.T) {
	t.Parallel()
	unavailable := status.Error(codes.Unavailable, "agent unavailable")
	aborted := status.Error(codes.Aborted, "flush in progress")
	internal := status.Error(codes.Internal, "agent internal")
	const addrA, addrB = "a:8081", "b:8081"
	tests := []struct {
		// errs maps target to the error of each successive attempt; a missing
		// entry or an exhausted list means success.
		errs     map[string][]error
		name     string
		addrs    []string
		wantIs   []error
		attempts int
		wantCode codes.Code // codes.OK: nil error expected
	}{
		{
			name:  "all agents succeed",
			addrs: []string{addrA, addrB},
		},
		{
			name:     "one agent unavailable",
			addrs:    []string{addrA, addrB},
			errs:     map[string][]error{addrB: {unavailable}},
			wantCode: codes.Unavailable,
			wantIs:   []error{unavailable},
		},
		{
			name:     "all agents aborted",
			addrs:    []string{addrA, addrB},
			errs:     map[string][]error{addrA: {aborted}, addrB: {aborted}},
			wantCode: codes.Aborted,
			wantIs:   []error{aborted},
		},
		{
			name:     "mixed codes are all kept and the first target's code leads",
			addrs:    []string{addrB, addrA},
			errs:     map[string][]error{addrA: {internal}, addrB: {aborted}},
			wantCode: codes.Internal,
			wantIs:   []error{internal, aborted},
		},
		{
			name:     "a retry that succeeds clears the earlier error",
			addrs:    []string{addrA},
			attempts: 2,
			errs:     map[string][]error{addrA: {unavailable}},
		},
		{
			name:     "a retry that fails again keeps the last error",
			addrs:    []string{addrA},
			attempts: 2,
			errs:     map[string][]error{addrA: {unavailable, internal}},
			wantCode: codes.Internal,
			wantIs:   []error{internal},
		},
	}
	for _, test := range tests {
		for _, kind := range []BroadCastKind{READ, WRITE} {
			t.Run(test.name, func(t *testing.T) {
				t.Parallel()
				g := &gateway{client: &rangeDiscoverer{c: &rangeClient{addrs: test.addrs, attempts: test.attempts}}}
				calls := make(map[string]int) // the model calls f sequentially
				err := g.BroadCast(t.Context(), kind, func(_ context.Context, target string, _ vald.Client, _ ...grpc.CallOption) error {
					n := calls[target]
					calls[target]++
					if errs := test.errs[target]; n < len(errs) {
						return errs[n]
					}
					return nil
				})
				if test.wantCode == codes.OK {
					if err != nil {
						t.Fatalf("unexpected error: %v", err)
					}
					return
				}
				st, ok := status.FromError(err)
				if !ok || st == nil || st.Code() != test.wantCode {
					t.Fatalf("got error: %v, want code %v", err, test.wantCode)
				}
				for _, want := range test.wantIs {
					if !errors.Is(err, want) {
						t.Errorf("error %v does not contain %v", err, want)
					}
				}
			})
		}
	}
}

// No connection at all is still reported as ErrGRPCClientConnNotFound, which the
// handlers map to Internal.
func Test_gateway_BroadCast_connNotFound(t *testing.T) {
	t.Parallel()
	g := &gateway{client: &rangeDiscoverer{c: new(rangeClient)}}
	err := g.BroadCast(t.Context(), WRITE, func(context.Context, string, vald.Client, ...grpc.CallOption) error {
		t.Error("callback must not run without connections")
		return nil
	})
	if !errors.Is(err, errors.ErrGRPCClientConnNotFound("*")) {
		t.Fatalf("got %v, want connection not found", err)
	}
}
