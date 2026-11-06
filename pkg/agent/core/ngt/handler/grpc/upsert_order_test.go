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
	"math"
	"slices"
	"testing"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/sync/errgroup"
)

// upsertNGT reports the IDs in existing as present, fails every update with
// err and accepts every insert; other methods panic through the nil embedded
// interface.
type upsertNGT struct {
	existing map[string]bool
	stubNGT
}

func (n *upsertNGT) Exists(uuid string) (uint32, bool) { return 0, n.existing[uuid] }

func (n *upsertNGT) UpdateWithTime(string, []float32, int64) error { return n.err }

func (n *upsertNGT) UpdateMultiple(map[string][]float32) error { return n.err }

func (n *upsertNGT) UpdateMultipleWithTime(map[string][]float32, int64) error { return n.err }

func (*upsertNGT) InsertWithTime(string, []float32, int64) error { return nil }

func (*upsertNGT) InsertMultiple(map[string][]float32) error { return nil }

func (*upsertNGT) InsertMultipleWithTime(map[string][]float32, int64) error { return nil }

// Test_server_update_invalidTimestamp pins that the service's rejection of a
// MaxInt64 timestamp reaches the client as InvalidArgument, as the LB gateway
// and UpdateTimestamp report it, instead of Internal, through every update path.
func Test_server_update_invalidTimestamp(t *testing.T) {
	t.Parallel()
	const ts = math.MaxInt64
	vec := []float32{1, 2, 3}
	ovec := func(id string) *payload.Object_Vector { return &payload.Object_Vector{Id: id, Vector: vec} }
	upserts := func(ids ...string) *payload.Upsert_MultiRequest {
		reqs := &payload.Upsert_MultiRequest{Requests: make([]*payload.Upsert_Request, 0, len(ids))}
		for _, id := range ids {
			reqs.Requests = append(reqs.Requests, &payload.Upsert_Request{Vector: ovec(id), Config: &payload.Upsert_Config{Timestamp: ts}})
		}
		return reqs
	}
	tests := []struct {
		call func(context.Context, *server) error
		name string
	}{
		{
			name: "Update",
			call: func(ctx context.Context, s *server) error {
				_, err := s.Update(ctx, &payload.Update_Request{Vector: ovec("a"), Config: &payload.Update_Config{Timestamp: ts}})
				return err
			},
		},
		{
			name: "MultiUpdate",
			call: func(ctx context.Context, s *server) error {
				_, err := s.MultiUpdate(ctx, &payload.Update_MultiRequest{Requests: []*payload.Update_Request{
					{Vector: ovec("a"), Config: &payload.Update_Config{Timestamp: ts}},
					{Vector: ovec("b"), Config: &payload.Update_Config{Timestamp: ts}},
				}})
				return err
			},
		},
		{
			name: "Upsert of an existing ID",
			call: func(ctx context.Context, s *server) error {
				_, err := s.Upsert(ctx, &payload.Upsert_Request{Vector: ovec("a"), Config: &payload.Upsert_Config{Timestamp: ts}})
				return err
			},
		},
		{
			name: "MultiUpsert of existing IDs",
			call: func(ctx context.Context, s *server) error {
				_, err := s.MultiUpsert(ctx, upserts("a", "b"))
				return err
			},
		},
		{
			name: "MultiUpsert of existing and new IDs",
			call: func(ctx context.Context, s *server) error {
				_, err := s.MultiUpsert(ctx, upserts("a", "new"))
				return err
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			eg, ctx := errgroup.New(t.Context())
			s := &server{eg: eg, ngt: &upsertNGT{
				existing: map[string]bool{"a": true, "b": true},
				err:      errors.ErrInvalidTimestamp(ts), dim: len(vec),
			}}
			err := test.call(ctx, s)
			st, ok := status.FromError(err)
			if err == nil || !ok || st.Code() != codes.InvalidArgument {
				t.Errorf("got error %v, want code %v", err, codes.InvalidArgument)
			}
		})
	}
}

// Test_server_MultiUpsert_requestOrder pins that MultiUpsert returns the
// locations in request order, like the other Multi* handlers, also when the
// batch mixes existing (updated) and new (inserted) IDs.
func Test_server_MultiUpsert_requestOrder(t *testing.T) {
	t.Parallel()
	const (
		new1 = "new-1"
		new2 = "new-2"
		old1 = "old-1"
		old2 = "old-2"
	)
	vec := []float32{1, 2, 3}
	tests := []struct {
		name string
		ids  []string
	}{
		{name: "new ID first", ids: []string{new1, old1, new2, old2}},
		{name: "existing ID first", ids: []string{old1, new1, old2, new2}},
		{name: "existing IDs only", ids: []string{old2, old1}},
		{name: "new IDs only", ids: []string{new2, new1}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			eg, ctx := errgroup.New(t.Context())
			s := &server{eg: eg, ngt: &upsertNGT{
				existing: map[string]bool{old1: true, old2: true},
				dim:      len(vec),
			}}
			reqs := &payload.Upsert_MultiRequest{Requests: make([]*payload.Upsert_Request, 0, len(test.ids))}
			for _, id := range test.ids {
				reqs.Requests = append(reqs.Requests, &payload.Upsert_Request{Vector: &payload.Object_Vector{Id: id, Vector: vec}})
			}
			res, err := s.MultiUpsert(ctx, reqs)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			got := make([]string, 0, len(res.GetLocations()))
			for _, loc := range res.GetLocations() {
				got = append(got, loc.GetUuid())
			}
			if !slices.Equal(got, test.ids) {
				t.Errorf("location uuids = %v, want %v", got, test.ids)
			}
		})
	}
}
