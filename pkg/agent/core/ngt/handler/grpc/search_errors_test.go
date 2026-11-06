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
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/pkg/agent/core/ngt/service"
)

const searchErrorsTestID = "uuid"

// stubNGT returns a fixed error from every search method; other methods panic through the nil embedded interface.
type stubNGT struct {
	service.NGT
	err error
	dim int
}

func (n *stubNGT) GetDimensionSize() int { return n.dim }

func (n *stubNGT) Search(
	context.Context, []float32, uint32, float32, float32, int32,
) (*payload.Search_Response, error) {
	return nil, n.err
}

func (n *stubNGT) SearchByID(
	context.Context, string, uint32, float32, float32, int32,
) ([]float32, *payload.Search_Response, error) {
	return nil, nil, n.err
}

func (n *stubNGT) LinearSearch(
	context.Context, []float32, uint32,
) (*payload.Search_Response, error) {
	return nil, n.err
}

func (n *stubNGT) LinearSearchByID(
	context.Context, string, uint32,
) ([]float32, *payload.Search_Response, error) {
	return nil, nil, n.err
}

// Test_server_search_concurrentCoreError calls each search handler concurrently with an NGT core error.
// The errors.As target used to be a package-level variable, which is a data race under -race.
func Test_server_search_concurrentCoreError(t *testing.T) {
	t.Parallel()

	const callers = 32
	vec := []float32{0, 1}
	tests := []struct {
		call    func(context.Context, *server) error
		name    string
		wantMsg string
	}{
		{
			name: "Search",
			call: func(ctx context.Context, s *server) error {
				_, err := s.Search(ctx, &payload.Search_Request{Vector: vec})
				return err
			},
			wantMsg: "Search API failed to process search request due to ngt core process returned error",
		},
		{
			name: "SearchByID",
			call: func(ctx context.Context, s *server) error {
				_, err := s.SearchByID(ctx, &payload.Search_IDRequest{Id: searchErrorsTestID})
				return err
			},
			wantMsg: "SearchByID API failed to process search request due to ngt core process returned error",
		},
		{
			name: "LinearSearch",
			call: func(ctx context.Context, s *server) error {
				_, err := s.LinearSearch(ctx, &payload.Search_Request{Vector: vec})
				return err
			},
			wantMsg: "LinearSearch API failed to process search request due to ngt core process returned error",
		},
		{
			name: "LinearSearchByID",
			call: func(ctx context.Context, s *server) error {
				_, err := s.LinearSearchByID(ctx, &payload.Search_IDRequest{Id: searchErrorsTestID})
				return err
			},
			wantMsg: "LinearSearchByID API failed to process search request due to ngt core process returned error",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			// the core returns NGTError by value (errors.NewNGTError), which the As target must match.
			s := &server{ngt: &stubNGT{err: errors.NewNGTError("boom"), dim: len(vec)}}

			start := make(chan struct{})
			errs := make(chan error, callers)
			for range callers {
				go func() {
					<-start
					errs <- test.call(context.Background(), s)
				}()
			}
			close(start)

			timeout := time.After(10 * time.Second)
			for range callers {
				select {
				case err := <-errs:
					st, ok := status.FromError(err)
					if !ok {
						t.Errorf("got non-status error: %v", err)
						continue
					}
					if st.Code() != codes.Internal {
						t.Errorf("code = %v, want %v", st.Code(), codes.Internal)
					}
					if st.Message() != test.wantMsg {
						t.Errorf("message = %q, want %q", st.Message(), test.wantMsg)
					}
				case <-timeout:
					t.Fatal("timed out waiting for search callers")
				}
			}
		})
	}
}
