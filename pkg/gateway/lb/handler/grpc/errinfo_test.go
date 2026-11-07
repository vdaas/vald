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
	"fmt"
	"slices"
	"testing"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/errdetails"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/sync/errgroup"
)

// extClient adds the agent calls okClient does not stub.
type extClient struct {
	*okClient
	noIps bool // answer UpdateTimestamp with a location that carries no ips
}

func (c extClient) UpdateTimestamp(
	context.Context, *payload.Update_TimestampRequest, ...grpc.CallOption,
) (*payload.Object_Location, error) {
	if c.noIps {
		return &payload.Object_Location{}, nil
	}
	return &payload.Object_Location{Uuid: fixtureUUID, Name: "agent-ts", Ips: []string{"10.0.0.9"}}, nil
}

func (c extClient) LinearSearch(
	context.Context, *payload.Search_Request, ...grpc.CallOption,
) (*payload.Search_Response, error) {
	return &payload.Search_Response{Results: []*payload.Object_Distance{{Id: "id", Distance: 0.1}}}, nil
}

func newExtClient(vec []float32) extClient {
	c := extClient{okClient: newOKClient(vec)}
	c.search = func(context.Context, *payload.Search_Request) (*payload.Search_Response, error) {
		return &payload.Search_Response{Results: []*payload.Object_Distance{{Id: "id", Distance: 0.1}}}, nil
	}
	return c
}

// The error details (full request JSON, resource names, gateway addresses) are only
// needed when an RPC fails. A successful RPC must not build them: gateway.Addrs is
// queried only while building them, so it must not be called at all.
func Test_server_successDoesNotBuildErrorDetails(t *testing.T) {
	t.Parallel()
	const uuid = "uuid"
	vec := []float32{1, 2, 3}
	tests := []struct {
		call func(ctx context.Context, s *server) error
		name string
	}{
		{
			name: vald.InsertRPCName,
			call: func(ctx context.Context, s *server) error {
				_, err := s.Insert(ctx, &payload.Insert_Request{Vector: &payload.Object_Vector{Id: uuid, Vector: vec}})
				return err
			},
		},
		{
			name: vald.UpsertRPCName,
			call: func(ctx context.Context, s *server) error {
				_, err := s.Upsert(ctx, &payload.Upsert_Request{
					Vector: &payload.Object_Vector{Id: uuid, Vector: vec},
					Config: &payload.Upsert_Config{SkipStrictExistCheck: true},
				})
				return err
			},
		},
		{
			name: vald.UpdateRPCName,
			call: func(ctx context.Context, s *server) error {
				_, err := s.Update(ctx, &payload.Update_Request{Vector: &payload.Object_Vector{Id: uuid, Vector: vec}})
				return err
			},
		},
		{
			name: vald.UpdateTimestampRPCName,
			call: func(ctx context.Context, s *server) error {
				_, err := s.UpdateTimestamp(ctx, &payload.Update_TimestampRequest{Id: uuid, Timestamp: 1})
				return err
			},
		},
		{
			name: vald.RemoveRPCName,
			call: func(ctx context.Context, s *server) error {
				_, err := s.Remove(ctx, &payload.Remove_Request{
					Id:     &payload.Object_ID{Id: uuid},
					Config: &payload.Remove_Config{SkipStrictExistCheck: true},
				})
				return err
			},
		},
		{
			name: vald.SearchByIDRPCName,
			call: func(ctx context.Context, s *server) error {
				_, err := s.SearchByID(ctx, &payload.Search_IDRequest{Id: uuid, Config: &payload.Search_Config{Num: 3}})
				return err
			},
		},
		{
			name: vald.LinearSearchByIDRPCName,
			call: func(ctx context.Context, s *server) error {
				_, err := s.LinearSearchByID(ctx, &payload.Search_IDRequest{Id: uuid, Config: &payload.Search_Config{Num: 3}})
				return err
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			// The stored vector differs from the requested one so Update proceeds to remove+insert,
			// and Exists finds nothing so Insert may proceed.
			c := newExtClient([]float32{4, 5, 6})
			c.exists = func(context.Context, *payload.Object_ID) (*payload.Object_ID, error) {
				return &payload.Object_ID{}, nil
			}
			gw := syncGateway(c)
			if test.name == "Upsert" || test.name == "Update" {
				// Upsert/Update look the object up first; the stored vector must be found.
				c.exists = func(context.Context, *payload.Object_ID) (*payload.Object_ID, error) {
					return &payload.Object_ID{Id: uuid}, nil
				}
			}
			eg, _ := errgroup.New(t.Context())
			s := &server{gateway: gw, replica: 1, eg: eg, timeout: time.Second}
			if err := test.call(t.Context(), s); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if n := gw.addrsCalls.Load(); n != 0 {
				t.Errorf("a successful %s built the error details: gateway.Addrs called %d times, want 0", test.name, n)
			}
		})
	}
}

// statusSig reduces an error to its status code, message and details. The details are
// sorted because the status builder emits them in map order.
func statusSig(t *testing.T, err error) string {
	t.Helper()
	st, ok := status.FromError(err)
	if !ok || st == nil {
		t.Fatalf("not a status error: %v", err)
	}
	ds := make([]string, 0, len(st.Proto().GetDetails()))
	for _, d := range st.Proto().GetDetails() {
		ds = append(ds, d.GetTypeUrl()+"|"+string(d.GetValue()))
	}
	slices.Sort(ds)
	return fmt.Sprintf("%v|%s|%q", st.Code(), st.Message(), ds)
}

// The details of every error the lazily built RequestInfo/ResourceInfo feed into must be
// exactly what the handlers used to build up front: expected values are assembled here
// from the same inputs the old code used (request JSON, resource type, name, gateway addrs).
func Test_server_errorDetailsUnchanged(t *testing.T) {
	t.Parallel()
	const (
		uuid = "uuid"
		name = "lb-name"
		ip   = "10.0.0.1"
	)
	addrs := []string{"10.0.0.2:8081", "10.0.0.3:8081"}
	vec := []float32{1, 2, 3}
	invalidID := &errdetails.BadRequest{
		FieldViolations: []*errdetails.BadRequestFieldViolation{
			{Field: "invalid id", Description: errors.ErrInvalidMetaDataConfig.Error()},
		},
	}
	insertReq := &payload.Insert_Request{Vector: &payload.Object_Vector{Vector: vec}}
	upsertReq := &payload.Upsert_Request{Vector: &payload.Object_Vector{Vector: vec}}
	updateReq := &payload.Update_Request{Vector: &payload.Object_Vector{Vector: vec}}
	tsReq := &payload.Update_TimestampRequest{Timestamp: 1}
	tsNotFoundReq := &payload.Update_TimestampRequest{Id: uuid, Timestamp: 1}
	searchByIDReq := &payload.Search_IDRequest{Config: &payload.Search_Config{Num: 3}}
	removeReq := &payload.Remove_Request{
		Id:     &payload.Object_ID{Id: uuid},
		Config: &payload.Remove_Config{SkipStrictExistCheck: true},
	}
	tests := []struct {
		call         func(ctx context.Context, s *server) error
		req          any
		wrap         func(msg string, err error, details ...any) error
		cause        error
		name         string
		msg          string
		uuid         string
		resourceType string
		extra        []any
		noIps        bool
	}{
		{
			name: "Insert invalid uuid", req: insertReq, wrap: status.WrapWithInvalidArgument,
			call: func(ctx context.Context, s *server) error { _, err := s.Insert(ctx, insertReq); return err },
			msg:  vald.InsertRPCName + " API invalid uuid", cause: errors.ErrInvalidMetaDataConfig,
			resourceType: "/vald.v1." + vald.InsertRPCName, extra: []any{invalidID},
		},
		{
			name: "Upsert invalid uuid", req: upsertReq, wrap: status.WrapWithInvalidArgument,
			call: func(ctx context.Context, s *server) error { _, err := s.Upsert(ctx, upsertReq); return err },
			msg:  vald.UpsertRPCName + " API invalid uuid", cause: errors.ErrInvalidMetaDataConfig,
			resourceType: "/vald.v1." + vald.UpsertRPCName, extra: []any{invalidID},
		},
		{
			name: "Update invalid uuid", req: updateReq, wrap: status.WrapWithInvalidArgument,
			call: func(ctx context.Context, s *server) error { _, err := s.Update(ctx, updateReq); return err },
			msg:  vald.UpdateRPCName + " API invalid uuid", cause: errors.ErrInvalidMetaDataConfig,
			resourceType: "/vald.v1." + vald.UpdateRPCName + "." + vald.GetObjectRPCName, extra: []any{invalidID},
		},
		{
			name: "UpdateTimestamp invalid uuid", req: tsReq, wrap: status.WrapWithInvalidArgument,
			call: func(ctx context.Context, s *server) error { _, err := s.UpdateTimestamp(ctx, tsReq); return err },
			msg:  vald.UpdateTimestampRPCName + " API invalid uuid", cause: errors.ErrInvalidMetaDataConfig,
			resourceType: "/vald.v1." + vald.UpdateTimestampRPCName + "." + vald.GetObjectRPCName, extra: []any{invalidID},
		},
		{
			name: "SearchByID invalid uuid", req: searchByIDReq, wrap: status.WrapWithInvalidArgument,
			call: func(ctx context.Context, s *server) error { _, err := s.SearchByID(ctx, searchByIDReq); return err },
			msg:  vald.SearchByIDRPCName + " API invalid uuid", cause: errors.ErrInvalidMetaDataConfig,
			resourceType: "/vald.v1." + vald.SearchByIDRPCName, extra: []any{invalidID},
		},
		{
			name: "LinearSearchByID invalid uuid", req: searchByIDReq, wrap: status.WrapWithInvalidArgument,
			call: func(ctx context.Context, s *server) error {
				_, err := s.LinearSearchByID(ctx, searchByIDReq)
				return err
			},
			msg: vald.LinearSearchByIDRPCName + " API invalid uuid", cause: errors.ErrInvalidMetaDataConfig,
			resourceType: "/vald.v1." + vald.LinearSearchByIDRPCName, extra: []any{invalidID},
		},
		{
			name: "Remove target not found", req: removeReq, uuid: uuid, wrap: status.WrapWithNotFound,
			call: func(ctx context.Context, s *server) error { _, err := s.Remove(ctx, removeReq); return err },
			msg:  vald.RemoveRPCName + " API remove target not found", cause: errors.ErrIndexNotFound,
			resourceType: "/vald.v1." + vald.RemoveRPCName,
		},
		{
			name: "UpdateTimestamp target not found", req: tsNotFoundReq, uuid: uuid, wrap: status.WrapWithNotFound, noIps: true,
			call: func(ctx context.Context, s *server) error {
				_, err := s.UpdateTimestamp(ctx, tsNotFoundReq)
				return err
			},
			msg: vald.UpdateTimestampRPCName + " API update target not found", cause: errors.ErrIndexNotFound,
			resourceType: "/vald.v1." + vald.UpdateTimestampRPCName + "." + vald.GetObjectRPCName,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			c := newExtClient(vec)
			c.noIps = test.noIps
			c.remove = func(context.Context, *payload.Remove_Request) (*payload.Object_Location, error) {
				return &payload.Object_Location{}, nil
			}
			gw := syncGateway(c)
			gw.addrs = addrs
			s := &server{gateway: gw, name: name, ip: ip, replica: 1, timeout: time.Second}

			err := test.call(t.Context(), s)
			if err == nil {
				t.Fatal("expected an error, got nil")
			}
			details := append([]any{
				&errdetails.RequestInfo{RequestId: test.uuid, ServingData: errdetails.Serialize(test.req)},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + test.resourceType,
					ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, name, ip, addrs),
				},
			}, test.extra...)
			want := test.wrap(test.msg, test.cause, details...)
			if got, wantSig := statusSig(t, err), statusSig(t, want); got != wantSig {
				t.Errorf("error status changed:\n got: %s\nwant: %s", got, wantSig)
			}
			// The details are built exactly once, on the error path.
			if n := gw.addrsCalls.Load(); n != 1 {
				t.Errorf("gateway.Addrs called %d times, want 1", n)
			}
		})
	}
}
