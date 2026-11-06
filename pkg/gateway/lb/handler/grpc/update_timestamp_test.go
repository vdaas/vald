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

// Update removes at ts and inserts at ts+1; ts == MaxInt64 used to wrap the
// insert timestamp to MinInt64, so the agent dropped the insert and Update
// still reported success. It must be rejected before anything is sent.
func Test_server_Update_timestamp(t *testing.T) {
	t.Parallel()
	const uuid = "uuid"
	tests := []struct {
		name       string
		ts         int64
		wantCode   codes.Code
		wantRemove int64
		wantInsert int64
	}{
		{name: "max int64 is rejected", ts: math.MaxInt64, wantCode: codes.InvalidArgument},
		{name: "normal timestamp", ts: 7, wantCode: codes.OK, wantRemove: 7, wantInsert: 8},
		{name: "max int64 - 1 still fits", ts: math.MaxInt64 - 1, wantCode: codes.OK, wantRemove: math.MaxInt64 - 1, wantInsert: math.MaxInt64},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			// The stored vector differs from the requested one, so Update proceeds to remove+insert.
			c := newOKClient([]float32{4, 5, 6})
			eg, _ := errgroup.New(t.Context())
			s := &server{gateway: syncGateway(c), replica: 1, eg: eg}

			_, err := s.Update(t.Context(), &payload.Update_Request{
				Vector: &payload.Object_Vector{Id: uuid, Vector: []float32{1, 2, 3}},
				Config: &payload.Update_Config{SkipStrictExistCheck: true, Timestamp: test.ts},
			})
			if test.wantCode != codes.OK {
				st, ok := status.FromError(err)
				if err == nil || !ok || st == nil || st.Code() != test.wantCode {
					t.Fatalf("got error %v, want code %v", err, test.wantCode)
				}
				if c.removeIn != nil || c.insertIn != nil {
					t.Errorf("a rejected Update reached the agent: remove=%v insert=%v", c.removeIn, c.insertIn)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got := c.removeIn.GetConfig().GetTimestamp(); got != test.wantRemove {
				t.Errorf("remove timestamp = %d, want %d", got, test.wantRemove)
			}
			if got := c.insertIn.GetConfig().GetTimestamp(); got != test.wantInsert {
				t.Errorf("insert timestamp = %d, want %d", got, test.wantInsert)
			}
		})
	}
}
