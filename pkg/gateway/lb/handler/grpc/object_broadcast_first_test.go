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
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/errgroup"
)

// seqGateway broadcasts to the given clients one after another and joins their errors.
func seqGateway(clients ...vald.Client) *regGateway {
	return &regGateway{
		agents: len(clients),
		broadCast: func(ctx context.Context, f func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error) error {
			errs := make([]error, 0, len(clients))
			for _, c := range clients {
				errs = append(errs, f(ctx, "127.0.0.1:8081", c))
			}
			return errors.Join(errs...)
		},
	}
}

func Test_broadcastFirst(t *testing.T) {
	t.Parallel()
	errBroadCast := errors.New("broadcast failed")
	miss := func(context.Context, *payload.Object_ID) (*payload.Object_ID, error) {
		return &payload.Object_ID{}, nil
	}
	hit := func(id string) func(context.Context, *payload.Object_ID) (*payload.Object_ID, error) {
		return func(context.Context, *payload.Object_ID) (*payload.Object_ID, error) {
			return &payload.Object_ID{Id: id}, nil
		}
	}
	tests := []struct {
		name    string
		gw      func() *regGateway
		wantErr error
		wantID  string
		wantOK  bool
	}{
		{
			name:   "first hit is returned",
			gw:     func() *regGateway { return seqGateway(&regClient{exists: miss}, &regClient{exists: hit("a")}) },
			wantID: "a",
			wantOK: true,
		},
		{
			name: "no hit is not an error",
			gw:   func() *regGateway { return seqGateway(&regClient{exists: miss}, &regClient{exists: miss}) },
		},
		{
			name: "broadcast error is reported when nothing was found",
			gw: func() *regGateway {
				return &regGateway{
					broadCast: func(context.Context, func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error) error {
						return errBroadCast
					},
				}
			},
			wantErr: errBroadCast,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
			defer cancel()
			eg, _ := errgroup.New(ctx)
			eg.SetLimit(1)
			s := &server{gateway: test.gw(), eg: eg}
			req := &payload.Object_ID{Id: "a"}
			got, ok, err := broadcastFirst(ctx, s, "exists",
				func(oid *payload.Object_ID) bool { return oid.GetId() != "" },
				func(ctx context.Context, _ trace.Span, _ string, vc vald.Client, copts ...grpc.CallOption) (*payload.Object_ID, error) {
					return vc.Exists(ctx, req, copts...)
				})
			if (test.wantErr == nil && err != nil) || (test.wantErr != nil && !errors.Is(err, test.wantErr)) {
				t.Errorf("err = %v, want %v", err, test.wantErr)
			}
			if ok != test.wantOK || got.GetId() != test.wantID {
				t.Errorf("got (%v, %t), want (%q, %t)", got, ok, test.wantID, test.wantOK)
			}
		})
	}
}

// exists/getObject/GetTimestamp build the read-only agent request once per call
// and share it across the broadcast targets instead of allocating one per target.
func Test_server_objectReads_shareRequestAcrossTargets(t *testing.T) {
	t.Parallel()
	const (
		uuid    = "uuid"
		targets = 3
	)
	tests := []struct {
		call func(ctx context.Context, record func(any)) error
		name string
	}{
		{
			name: "exists",
			call: func(ctx context.Context, record func(any)) error {
				c := &regClient{exists: func(_ context.Context, in *payload.Object_ID) (*payload.Object_ID, error) {
					record(in)
					return &payload.Object_ID{}, nil
				}}
				s := &server{gateway: seqGateway(c, c, c)}
				s.eg, _ = errgroup.New(ctx)
				_, err := s.exists(ctx, uuid)
				return err
			},
		},
		{
			name: "getObject",
			call: func(ctx context.Context, record func(any)) error {
				c := &regClient{getObject: func(_ context.Context, in *payload.Object_VectorRequest) (*payload.Object_Vector, error) {
					record(in)
					return &payload.Object_Vector{}, nil
				}}
				s := &server{gateway: seqGateway(c, c, c)}
				s.eg, _ = errgroup.New(ctx)
				_, err := s.getObject(ctx, uuid)
				return err
			},
		},
		{
			name: "GetTimestamp",
			call: func(ctx context.Context, record func(any)) error {
				c := &regClient{getTS: func(_ context.Context, in *payload.Object_TimestampRequest) (*payload.Object_Timestamp, error) {
					record(in)
					return &payload.Object_Timestamp{}, nil
				}}
				s := &server{gateway: seqGateway(c, c, c)}
				s.eg, _ = errgroup.New(ctx)
				_, err := s.GetTimestamp(ctx, &payload.Object_TimestampRequest{Id: &payload.Object_ID{Id: uuid}})
				return err
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
			defer cancel()
			var (
				mu   sync.Mutex
				seen []any
			)
			record := func(in any) {
				mu.Lock()
				defer mu.Unlock()
				seen = append(seen, in)
			}
			// every target misses (empty answer), so the call reports NotFound after all of them ran.
			if err := test.call(ctx, record); err == nil {
				t.Fatal("expected a not-found error, got nil")
			}
			if ctx.Err() != nil {
				t.Fatalf("call did not finish before the test deadline: %v", ctx.Err())
			}
			mu.Lock()
			defer mu.Unlock()
			if len(seen) != targets {
				t.Fatalf("agents called %d times, want %d", len(seen), targets)
			}
			for i, in := range seen[1:] {
				if in != seen[0] {
					t.Errorf("target %d got request %p, want the shared request %p", i+1, in, seen[0])
				}
			}
		})
	}
}
