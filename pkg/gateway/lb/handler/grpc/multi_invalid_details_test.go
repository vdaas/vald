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
)

// The invalid-argument details of a Multi* batch (the serialized batch and the
// agent addresses) are the same for every invalid request, so they are built
// once per batch rather than once per invalid request.
func Test_server_Multi_buildsInvalidDetailsOncePerBatch(t *testing.T) {
	t.Parallel()
	const invalid = 3
	tests := []struct {
		call func(ctx context.Context, s *server) error
		name string
	}{
		{
			name: "MultiInsert",
			call: func(ctx context.Context, s *server) error {
				reqs := &payload.Insert_MultiRequest{}
				for range invalid {
					reqs.Requests = append(reqs.Requests, &payload.Insert_Request{Vector: &payload.Object_Vector{}})
				}
				_, err := s.MultiInsert(ctx, reqs)
				return err
			},
		},
		{
			name: "MultiUpdate",
			call: func(ctx context.Context, s *server) error {
				reqs := &payload.Update_MultiRequest{}
				for range invalid {
					reqs.Requests = append(reqs.Requests, &payload.Update_Request{Vector: &payload.Object_Vector{}})
				}
				_, err := s.MultiUpdate(ctx, reqs)
				return err
			},
		},
		{
			name: "MultiUpsert",
			call: func(ctx context.Context, s *server) error {
				reqs := &payload.Upsert_MultiRequest{}
				for range invalid {
					reqs.Requests = append(reqs.Requests, &payload.Upsert_Request{Vector: &payload.Object_Vector{}})
				}
				_, err := s.MultiUpsert(ctx, reqs)
				return err
			},
		},
		{
			name: "MultiRemove",
			call: func(ctx context.Context, s *server) error {
				reqs := &payload.Remove_MultiRequest{}
				for range invalid {
					reqs.Requests = append(reqs.Requests, &payload.Remove_Request{Id: &payload.Object_ID{}})
				}
				_, err := s.MultiRemove(ctx, reqs)
				return err
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			gw := &regGateway{addrs: []string{"127.0.0.1:8081"}}
			s := &server{gateway: gw, name: "lb", ip: "192.0.2.1", multiConcurrency: 1}
			if err := test.call(t.Context(), s); err == nil {
				t.Fatal("expected the invalid requests to be reported")
			}
			if n := gw.addrsCalls.Load(); n != 1 {
				t.Errorf("gateway.Addrs called %d times for %d invalid requests, want 1", n, invalid)
			}
		})
	}
}
