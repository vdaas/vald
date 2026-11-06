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
	client "github.com/vdaas/vald/internal/client/v1/client/vald"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
)

// strictExistGatewayStub records what the filter gateway forwards to the LB gateway.
type strictExistGatewayStub struct {
	client.Client
	existsID    *payload.Object_ID
	existsErr   error
	insertReq   *payload.Insert_Request
	updateReq   *payload.Update_Request
	upsertReq   *payload.Upsert_Request
	existsCalls int
}

func (g *strictExistGatewayStub) Exists(
	_ context.Context, _ *payload.Object_ID, _ ...grpc.CallOption,
) (*payload.Object_ID, error) {
	g.existsCalls++
	return g.existsID, g.existsErr
}

func (g *strictExistGatewayStub) Insert(
	_ context.Context, req *payload.Insert_Request, _ ...grpc.CallOption,
) (*payload.Object_Location, error) {
	g.insertReq = req
	return &payload.Object_Location{Uuid: req.GetVector().GetId()}, nil
}

func (g *strictExistGatewayStub) Update(
	_ context.Context, req *payload.Update_Request, _ ...grpc.CallOption,
) (*payload.Object_Location, error) {
	g.updateReq = req
	return &payload.Object_Location{Uuid: req.GetVector().GetId()}, nil
}

func (g *strictExistGatewayStub) Upsert(
	_ context.Context, req *payload.Upsert_Request, _ ...grpc.CallOption,
) (*payload.Object_Location, error) {
	g.upsertReq = req
	return &payload.Object_Location{Uuid: req.GetVector().GetId()}, nil
}

const testStrictExistUUID = "uuid-1"

func testStrictExistVector() *payload.Object_Vector {
	return &payload.Object_Vector{Id: testStrictExistUUID, Vector: []float32{0.1, 0.2, 0.3, 0.4}}
}

func Test_server_Insert_strictExistCheck(t *testing.T) {
	t.Parallel()
	tests := []struct {
		existsErr   error
		existsID    *payload.Object_ID
		name        string
		wantCode    codes.Code
		wantForward bool
	}{
		{
			name:        "not found proceeds to insert",
			existsErr:   status.Error(codes.NotFound, "not found"),
			wantCode:    codes.OK,
			wantForward: true,
		},
		{
			name:     "existing id returns AlreadyExists",
			existsID: &payload.Object_ID{Id: testStrictExistUUID},
			wantCode: codes.AlreadyExists,
		},
		{
			name:      "unavailable Exists error is returned, not bypassed",
			existsErr: status.Error(codes.Unavailable, "lb gateway unavailable"),
			wantCode:  codes.Unavailable,
		},
		{
			name:      "deadline exceeded Exists error is returned, not bypassed",
			existsErr: context.DeadlineExceeded,
			wantCode:  codes.DeadlineExceeded,
		},
		{
			name:      "non-status Exists error is returned as Internal",
			existsErr: errors.New("connection not found"),
			wantCode:  codes.Internal,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(tt *testing.T) {
			tt.Parallel()
			gw := &strictExistGatewayStub{existsID: tc.existsID, existsErr: tc.existsErr}
			s := &server{gateway: gw}
			_, err := s.Insert(tt.Context(), &payload.Insert_Request{Vector: testStrictExistVector()})
			gotCode := codes.OK
			if err != nil {
				st, ok := status.FromError(err)
				if !ok || st == nil {
					tt.Fatalf("error is not a gRPC status: %v", err)
				}
				gotCode = st.Code()
			}
			if gotCode != tc.wantCode {
				tt.Errorf("code = %v, want %v (err: %v)", gotCode, tc.wantCode, err)
			}
			if forwarded := gw.insertReq != nil; forwarded != tc.wantForward {
				tt.Errorf("forwarded = %v, want %v", forwarded, tc.wantForward)
			}
		})
	}
}

// Update/Upsert of an existing ID must reach the LB gateway, which performs
// the strict exist check (same-vector comparison) itself.
func Test_server_UpdateUpsert_existingIDForwarded(t *testing.T) {
	t.Parallel()
	tests := []struct {
		call func(*server, context.Context, *strictExistGatewayStub) (forwarded, skip bool, err error)
		name string
	}{
		{
			name: "Update",
			call: func(s *server, ctx context.Context, gw *strictExistGatewayStub) (bool, bool, error) {
				_, err := s.Update(ctx, &payload.Update_Request{Vector: testStrictExistVector()})
				return gw.updateReq != nil, gw.updateReq.GetConfig().GetSkipStrictExistCheck(), err
			},
		},
		{
			name: "Upsert",
			call: func(s *server, ctx context.Context, gw *strictExistGatewayStub) (bool, bool, error) {
				_, err := s.Upsert(ctx, &payload.Upsert_Request{Vector: testStrictExistVector()})
				return gw.upsertReq != nil, gw.upsertReq.GetConfig().GetSkipStrictExistCheck(), err
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(tt *testing.T) {
			tt.Parallel()
			gw := &strictExistGatewayStub{existsID: &payload.Object_ID{Id: testStrictExistUUID}}
			s := &server{gateway: gw}
			forwarded, skip, err := tc.call(s, tt.Context(), gw)
			if err != nil {
				tt.Fatalf("unexpected error: %v", err)
			}
			if !forwarded {
				tt.Fatal("request was not forwarded to the LB gateway")
			}
			if skip {
				tt.Error("SkipStrictExistCheck was forced to true; LB strict check is bypassed")
			}
		})
	}
}
