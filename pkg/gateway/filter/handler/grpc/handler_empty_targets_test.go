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

	egressapi "github.com/vdaas/vald/apis/grpc/v1/filter/egress"
	ingressapi "github.com/vdaas/vald/apis/grpc/v1/filter/ingress"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/client/v1/client/filter/egress"
	"github.com/vdaas/vald/internal/client/v1/client/filter/ingress"
	client "github.com/vdaas/vald/internal/client/v1/client/vald"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
)

// okGatewayStub is an LB gateway whose search/GetObject calls all succeed.
type okGatewayStub struct {
	client.Client
}

func okSearchResponse() *payload.Search_Response {
	return &payload.Search_Response{Results: []*payload.Object_Distance{{Id: testStrictExistUUID, Distance: 0.1}}}
}

func (*okGatewayStub) Search(
	context.Context, *payload.Search_Request, ...grpc.CallOption,
) (*payload.Search_Response, error) {
	return okSearchResponse(), nil
}

func (*okGatewayStub) SearchByID(
	context.Context, *payload.Search_IDRequest, ...grpc.CallOption,
) (*payload.Search_Response, error) {
	return okSearchResponse(), nil
}

func (*okGatewayStub) LinearSearch(
	context.Context, *payload.Search_Request, ...grpc.CallOption,
) (*payload.Search_Response, error) {
	return okSearchResponse(), nil
}

func (*okGatewayStub) LinearSearchByID(
	context.Context, *payload.Search_IDRequest, ...grpc.CallOption,
) (*payload.Search_Response, error) {
	return okSearchResponse(), nil
}

func (*okGatewayStub) GetObject(
	context.Context, *payload.Object_VectorRequest, ...grpc.CallOption,
) (*payload.Object_Vector, error) {
	return &payload.Object_Vector{Id: testStrictExistUUID, Vector: []float32{0.1}}, nil
}

// noTargetIngressStub and noTargetEgressStub fail Target like the real
// clients do when no address is given (and on any other call, which these
// tests never expect).
type noTargetIngressStub struct {
	ingress.Client
}

func (*noTargetIngressStub) Release(context.Context, ...string) {}

func (*noTargetIngressStub) Target(
	_ context.Context, addrs ...string,
) (ingressapi.FilterClient, error) {
	if len(addrs) == 0 {
		return nil, errors.ErrTargetNotFound
	}
	return nil, errors.Errorf("unexpected ingress targets %v", addrs)
}

type noTargetEgressStub struct {
	egress.Client
}

func (*noTargetEgressStub) Release(context.Context, ...string) {}

func (*noTargetEgressStub) Target(
	_ context.Context, addrs ...string,
) (egressapi.FilterClient, error) {
	if len(addrs) == 0 {
		return nil, errors.ErrTargetNotFound
	}
	return nil, errors.Errorf("unexpected egress targets %v", addrs)
}

// An empty (non-nil) targets list, as decoded from REST `"targets": []`, and
// empty configured filters must skip filtering instead of failing with
// "filter targets [] not found".
func Test_server_emptyFilterTargetsSkipped(t *testing.T) {
	t.Parallel()
	empty := func() *payload.Filter_Config {
		return &payload.Filter_Config{Targets: []*payload.Filter_Target{}}
	}
	searchCfg := func() *payload.Search_Config {
		return &payload.Search_Config{IngressFilters: empty(), EgressFilters: empty()}
	}
	tests := []struct {
		call func(context.Context, *server) error
		name string
	}{
		{
			name: "Search",
			call: func(ctx context.Context, s *server) error {
				_, err := s.Search(ctx, &payload.Search_Request{Vector: []float32{0.1}, Config: searchCfg()})
				return err
			},
		},
		{
			name: "SearchByID",
			call: func(ctx context.Context, s *server) error {
				_, err := s.SearchByID(ctx, &payload.Search_IDRequest{Id: testStrictExistUUID, Config: searchCfg()})
				return err
			},
		},
		{
			name: "LinearSearch",
			call: func(ctx context.Context, s *server) error {
				_, err := s.LinearSearch(ctx, &payload.Search_Request{Vector: []float32{0.1}, Config: searchCfg()})
				return err
			},
		},
		{
			name: "LinearSearchByID",
			call: func(ctx context.Context, s *server) error {
				_, err := s.LinearSearchByID(ctx, &payload.Search_IDRequest{Id: testStrictExistUUID, Config: searchCfg()})
				return err
			},
		},
		{
			name: "GetObject",
			call: func(ctx context.Context, s *server) error {
				_, err := s.GetObject(ctx, &payload.Object_VectorRequest{
					Id:      &payload.Object_ID{Id: testStrictExistUUID},
					Filters: empty(),
				})
				return err
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(tt *testing.T) {
			tt.Parallel()
			s := &server{
				gateway:         &okGatewayStub{},
				ingress:         &noTargetIngressStub{},
				egress:          &noTargetEgressStub{},
				SearchFilters:   []string{},
				DistanceFilters: []string{},
				ObjectFilters:   []string{},
			}
			if err := tc.call(tt.Context(), s); err != nil {
				tt.Errorf("%s() error = %v, want nil", tc.name, err)
			}
		})
	}
}
