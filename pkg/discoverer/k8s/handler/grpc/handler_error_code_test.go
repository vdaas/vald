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

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/sync/singleflight"
)

// errDiscoverer is a service.Discoverer whose lookups all fail with err.
type errDiscoverer struct {
	err error
}

func (*errDiscoverer) Start(context.Context) (<-chan error, error) {
	ech := make(chan error)
	close(ech)
	return ech, nil
}

func (d *errDiscoverer) GetPods(*payload.Discoverer_Request) (*payload.Info_Pods, error) {
	return nil, d.err
}

func (d *errDiscoverer) GetNodes(*payload.Discoverer_Request) (*payload.Info_Nodes, error) {
	return nil, d.err
}

func (d *errDiscoverer) GetServices(*payload.Discoverer_Request) (*payload.Info_Services, error) {
	return nil, d.err
}

// Lookup misses must surface as NotFound (not retried, ignored by the client
// circuit breaker) and an unpopulated cache as Unavailable (still retried and
// counted by the breaker like Internal); anything else stays Internal.
func Test_server_discoverErrorCode(t *testing.T) {
	t.Parallel()
	req := &payload.Discoverer_Request{Name: "vald-agent", Namespace: "default", Node: "node-1"}
	calls := map[string]func(context.Context, *server) error{
		"Pods": func(ctx context.Context, s *server) error {
			_, err := s.Pods(ctx, req)
			return err
		},
		"Nodes": func(ctx context.Context, s *server) error {
			_, err := s.Nodes(ctx, req)
			return err
		},
		"Services": func(ctx context.Context, s *server) error {
			_, err := s.Services(ctx, req)
			return err
		},
	}
	tests := []struct {
		err      error
		name     string
		wantCode codes.Code
	}{
		{name: "node not found", err: errors.ErrNodeNotFound(req.GetNode()), wantCode: codes.NotFound},
		{name: "namespace not found", err: errors.ErrNamespaceNotFound(req.GetNamespace()), wantCode: codes.NotFound},
		{name: "pod name not found", err: errors.ErrPodNameNotFound(req.GetName()), wantCode: codes.NotFound},
		{name: "svc name not found", err: errors.ErrSvcNameNotFound(req.GetName()), wantCode: codes.NotFound},
		{name: "invalid discovery cache", err: errors.ErrInvalidDiscoveryCache, wantCode: codes.Unavailable},
		{name: "other error", err: errors.New("unexpected"), wantCode: codes.Internal},
	}
	for api, call := range calls {
		for _, tc := range tests {
			t.Run(api+"/"+tc.name, func(tt *testing.T) {
				tt.Parallel()
				s := &server{
					dsc:    &errDiscoverer{err: tc.err},
					pgroup: singleflight.New[*payload.Info_Pods](),
					ngroup: singleflight.New[*payload.Info_Nodes](),
					sgroup: singleflight.New[*payload.Info_Services](),
				}
				err := call(tt.Context(), s)
				if err == nil {
					tt.Fatal("expected error, got nil")
				}
				st, ok := status.FromError(err)
				if !ok || st == nil {
					tt.Fatalf("error is not a gRPC status: %v", err)
				}
				if st.Code() != tc.wantCode {
					tt.Errorf("code = %v, want %v (err: %v)", st.Code(), tc.wantCode, err)
				}
			})
		}
	}
}
