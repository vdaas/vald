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
	"math"
	"testing"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/sync/errgroup"
)

// Update removes the stored vector from every agent before it inserts the new one, so
// a NaN/Inf vector that only the agents' Insert rejects would leave the object deleted
// and the client with an error. The gateway must reject such a request before it
// touches any agent, for every write API that ends in an Update.
func Test_server_write_rejectsNaNOrInfBeforeRemove(t *testing.T) {
	t.Parallel()
	nan, inf := float32(math.NaN()), float32(math.Inf(1))
	for _, vec := range [][]float32{{0.4, nan, 0.6}, {0.4, inf, 0.6}, {0.4, -inf, 0.6}} {
		tests := []struct {
			call func(s *server) error
			name string
		}{
			{
				name: "Update",
				call: func(s *server) error {
					_, err := s.Update(t.Context(), &payload.Update_Request{
						Vector: &payload.Object_Vector{Id: fixtureUUID, Vector: vec},
						Config: &payload.Update_Config{SkipStrictExistCheck: true},
					})
					return err
				},
			},
			{
				name: "Update with exist check",
				call: func(s *server) error {
					_, err := s.Update(t.Context(), &payload.Update_Request{
						Vector: &payload.Object_Vector{Id: fixtureUUID, Vector: vec},
					})
					return err
				},
			},
			{
				name: "Upsert",
				call: func(s *server) error {
					_, err := s.Upsert(t.Context(), &payload.Upsert_Request{
						Vector: &payload.Object_Vector{Id: fixtureUUID, Vector: vec},
						Config: &payload.Upsert_Config{SkipStrictExistCheck: true},
					})
					return err
				},
			},
			{
				name: "MultiUpdate",
				call: func(s *server) error {
					_, err := s.MultiUpdate(t.Context(), &payload.Update_MultiRequest{
						Requests: []*payload.Update_Request{{
							Vector: &payload.Object_Vector{Id: fixtureUUID, Vector: vec},
						}},
					})
					return err
				},
			},
			{
				name: "MultiUpsert",
				call: func(s *server) error {
					_, err := s.MultiUpsert(t.Context(), &payload.Upsert_MultiRequest{
						Requests: []*payload.Upsert_Request{{
							Vector: &payload.Object_Vector{Id: fixtureUUID, Vector: vec},
						}},
					})
					return err
				},
			},
		}
		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				t.Parallel()
				// The stored vector differs from the requested one, so a balanced Update would remove it.
				a := newTSAgent("agent-0", loopbackIP, []float32{0.1, 0.2, 0.3}, 1)
				eg, _ := errgroup.New(t.Context())
				s := &server{gateway: syncGateway(a), replica: 1, eg: eg}

				err := test.call(s)
				if err == nil {
					t.Fatal("expected an InvalidArgument error, got nil")
				}
				if st, ok := status.FromError(err); !ok || st.Code() != codes.InvalidArgument {
					t.Errorf("error = %v, want code InvalidArgument", err)
				}
				stored, ts, delTS, inserts := a.state()
				if !stored || ts != 1 || delTS != 0 || inserts != 0 {
					t.Errorf("agent state changed by a rejected request: stored=%v ts=%d delTS=%d inserts=%d, want stored=true ts=1 delTS=0 inserts=0",
						stored, ts, delTS, inserts)
				}
			})
		}
	}
}
