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

	ingressapi "github.com/vdaas/vald/apis/grpc/v1/filter/ingress"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/client/v1/client/filter/ingress"
	client "github.com/vdaas/vald/internal/client/v1/client/vald"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
)

// lbErrGatewayStub is an LB gateway whose write/GetObject calls all fail with err.
type lbErrGatewayStub struct {
	client.Client
	err error
}

func (g *lbErrGatewayStub) Insert(
	context.Context, *payload.Insert_Request, ...grpc.CallOption,
) (*payload.Object_Location, error) {
	return nil, g.err
}

func (g *lbErrGatewayStub) Update(
	context.Context, *payload.Update_Request, ...grpc.CallOption,
) (*payload.Object_Location, error) {
	return nil, g.err
}

func (g *lbErrGatewayStub) Upsert(
	context.Context, *payload.Upsert_Request, ...grpc.CallOption,
) (*payload.Object_Location, error) {
	return nil, g.err
}

func (g *lbErrGatewayStub) GetObject(
	context.Context, *payload.Object_VectorRequest, ...grpc.CallOption,
) (*payload.Object_Vector, error) {
	return nil, g.err
}

// passIngressStub is an ingress filter that returns the vector unchanged.
type passIngressStub struct {
	ingress.Client
}

func (*passIngressStub) Release(context.Context, ...string) {}

func (p *passIngressStub) Target(context.Context, ...string) (ingressapi.FilterClient, error) {
	return p, nil
}

func (*passIngressStub) FilterVector(
	_ context.Context, in *payload.Object_Vector, _ ...grpc.CallOption,
) (*payload.Object_Vector, error) {
	return in, nil
}

// The filter gateway must preserve the LB gateway's status code after ingress
// filtering (and on GetObject) instead of rewriting it to Internal / NotFound.
func Test_server_lbErrorCodePreserved(t *testing.T) {
	t.Parallel()
	vec := func() *payload.Object_Vector {
		return &payload.Object_Vector{Id: "uuid-1", Vector: []float32{0.1, 0.2, 0.3, 0.4}}
	}
	calls := []struct {
		call func(context.Context, *server) error
		name string
	}{
		{
			name: "Insert",
			call: func(ctx context.Context, s *server) error {
				_, err := s.Insert(ctx, &payload.Insert_Request{
					Vector: vec(),
					Config: &payload.Insert_Config{SkipStrictExistCheck: true},
				})
				return err
			},
		},
		{
			name: "Update",
			call: func(ctx context.Context, s *server) error {
				_, err := s.Update(ctx, &payload.Update_Request{Vector: vec()})
				return err
			},
		},
		{
			name: "Upsert",
			call: func(ctx context.Context, s *server) error {
				_, err := s.Upsert(ctx, &payload.Upsert_Request{Vector: vec()})
				return err
			},
		},
		{
			name: "GetObject",
			call: func(ctx context.Context, s *server) error {
				_, err := s.GetObject(ctx, &payload.Object_VectorRequest{Id: &payload.Object_ID{Id: "uuid-1"}})
				return err
			},
		},
	}
	errs := []struct {
		lbErr    error
		name     string
		wantCode codes.Code
	}{
		{name: "NotFound", lbErr: status.Error(codes.NotFound, "not found"), wantCode: codes.NotFound},
		{name: "AlreadyExists", lbErr: status.Error(codes.AlreadyExists, "exists"), wantCode: codes.AlreadyExists},
		{name: "Unavailable", lbErr: status.Error(codes.Unavailable, "lb unavailable"), wantCode: codes.Unavailable},
		{name: "non-status error", lbErr: errors.New("connection not found"), wantCode: codes.Internal},
	}
	for _, c := range calls {
		for _, e := range errs {
			t.Run(c.name+"/"+e.name, func(tt *testing.T) {
				tt.Parallel()
				filters := []string{"ingress:8081"}
				s := &server{
					gateway:       &lbErrGatewayStub{err: e.lbErr},
					ingress:       &passIngressStub{},
					InsertFilters: filters,
					UpdateFilters: filters,
					UpsertFilters: filters,
				}
				err := c.call(tt.Context(), s)
				if err == nil {
					tt.Fatal("expected error, got nil")
				}
				st, ok := status.FromError(err)
				if !ok || st == nil {
					tt.Fatalf("error is not a gRPC status: %v", err)
				}
				if st.Code() != e.wantCode {
					tt.Errorf("code = %v, want %v (err: %v)", st.Code(), e.wantCode, err)
				}
			})
		}
	}
}
