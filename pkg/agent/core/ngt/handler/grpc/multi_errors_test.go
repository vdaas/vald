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
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/sync/errgroup"
)

// listNGT lists ids and fails DeleteWithTime for each of them with a per-id
// not-found error; other methods panic through the nil embedded interface.
type listNGT struct {
	stubNGT
	ids []string
}

func (n *listNGT) ListObjectFunc(
	_ context.Context, f func(uuid string, oid uint32, timestamp int64) bool,
) {
	for i, id := range n.ids {
		if !f(id, uint32(i), int64(i)) {
			return
		}
	}
}

func (*listNGT) DeleteWithTime(uuid string, _ int64) error {
	return errors.ErrObjectIDNotFound(uuid)
}

// Test_server_multi_joinedErrors pins the error the agent's Multi* search
// handlers and RemoveByTimestamp return now that per-item failures are
// collected in a slice and joined once after all items are processed: N
// distinct failures give one flat joined error with N members, in a status
// carrying the members' code.
func Test_server_multi_joinedErrors(t *testing.T) {
	t.Parallel()
	const n = 8
	vec := []float32{1, 2, 3}
	ids := make([]string, 0, n)
	for i := range n {
		ids = append(ids, "id-"+strconv.Itoa(i))
	}
	searchReqs := func() *payload.Search_MultiRequest {
		reqs := &payload.Search_MultiRequest{Requests: make([]*payload.Search_Request, 0, n)}
		for _, id := range ids {
			reqs.Requests = append(reqs.Requests, &payload.Search_Request{Vector: vec, Config: &payload.Search_Config{RequestId: id}})
		}
		return reqs
	}
	idReqs := func() *payload.Search_MultiIDRequest {
		reqs := &payload.Search_MultiIDRequest{Requests: make([]*payload.Search_IDRequest, 0, n)}
		for _, id := range ids {
			reqs.Requests = append(reqs.Requests, &payload.Search_IDRequest{Id: id})
		}
		return reqs
	}
	tests := []struct {
		call     func(context.Context, *server) error
		name     string
		wantCode codes.Code
	}{
		{
			name: "MultiSearch",
			call: func(ctx context.Context, s *server) error {
				_, err := s.MultiSearch(ctx, searchReqs())
				return err
			},
			wantCode: codes.NotFound,
		},
		{
			name: "MultiSearchByID",
			call: func(ctx context.Context, s *server) error {
				_, err := s.MultiSearchByID(ctx, idReqs())
				return err
			},
			wantCode: codes.NotFound,
		},
		{
			name: "MultiLinearSearch",
			call: func(ctx context.Context, s *server) error {
				_, err := s.MultiLinearSearch(ctx, searchReqs())
				return err
			},
			wantCode: codes.NotFound,
		},
		{
			name: "MultiLinearSearchByID",
			call: func(ctx context.Context, s *server) error {
				_, err := s.MultiLinearSearchByID(ctx, idReqs())
				return err
			},
			wantCode: codes.NotFound,
		},
		{
			name: "RemoveByTimestamp",
			call: func(ctx context.Context, s *server) error {
				_, err := s.RemoveByTimestamp(ctx, &payload.Remove_TimestampRequest{})
				return err
			},
			wantCode: codes.NotFound,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
			defer cancel()
			eg, ectx := errgroup.New(ctx)
			s := &server{
				ngt: &listNGT{err: errors.ErrEmptySearchResult, dim: len(vec), ids: ids},
				eg:  eg,
			}
			done := make(chan error, 1)
			go func() { done <- test.call(ectx, s) }()
			var err error
			select {
			case err = <-done:
			case <-ctx.Done():
				t.Fatal("timed out waiting for the handler")
			}
			if err == nil {
				t.Fatal("returned nil, want the per-item failures")
			}
			if st, ok := status.FromError(err); !ok || st == nil || st.Code() != test.wantCode {
				t.Errorf("status code = %v (ok=%t), want %v: %v", st.Code(), ok, test.wantCode, err)
			}
			var u interface{ Unwrap() []error }
			if !errors.As(err, &u) {
				t.Fatalf("got %T, want a joined error: %v", err, err)
			}
			if got := len(u.Unwrap()); got != n {
				t.Errorf("members = %d, want %d (flat): %v", got, n, err)
			}
			for _, id := range ids {
				if !strings.Contains(err.Error(), " "+id+"'s ") && !strings.Contains(err.Error(), " "+id+" ") {
					t.Errorf("error %q does not mention %q", err.Error(), id)
				}
			}
		})
	}
}
