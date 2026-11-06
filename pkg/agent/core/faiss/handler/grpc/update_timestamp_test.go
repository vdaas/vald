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
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
)

// updateErrFaiss fails every timed update with err; other methods panic
// through the nil embedded interface.
type updateErrFaiss struct {
	stubFaiss
}

func (f *updateErrFaiss) UpdateWithTime(string, []float32, int64) error { return f.err }

// Test_server_Update_invalidTimestamp pins that the service's rejection of a
// MaxInt64 timestamp reaches the client as InvalidArgument, as the LB gateway
// reports it, instead of Internal.
func Test_server_Update_invalidTimestamp(t *testing.T) {
	t.Parallel()
	const ts = math.MaxInt64
	vec := []float32{1, 2, 3}
	s := &server{faiss: &updateErrFaiss{err: errors.ErrInvalidTimestamp(ts), dim: len(vec)}}
	_, err := s.Update(t.Context(), &payload.Update_Request{
		Vector: &payload.Object_Vector{Id: "a", Vector: vec},
		Config: &payload.Update_Config{Timestamp: ts},
	})
	st, ok := status.FromError(err)
	if err == nil || !ok || st.Code() != codes.InvalidArgument {
		t.Errorf("got error %v, want code %v", err, codes.InvalidArgument)
	}
}
