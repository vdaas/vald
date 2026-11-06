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
	"testing"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/info"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/log/logger"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/os"
	"github.com/vdaas/vald/internal/test/goleak"
	"github.com/vdaas/vald/pkg/agent/core/faiss/service"
)

// TestMain initializes the global logger and info provider once, before any parallel test runs,
// and checks the package for leaked goroutines.
func TestMain(m *testing.M) {
	log.Init(log.WithLoggerType(logger.NOP.String()))
	if err := info.Init(""); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	goleak.VerifyTestMain(m)
}

// stubFaiss returns a fixed error from Search; every other method panics through the nil embedded interface.
type stubFaiss struct {
	service.Faiss
	err error
	dim int
}

func (f *stubFaiss) GetDimensionSize() int { return f.dim }

func (f *stubFaiss) Search(_, _, _ uint32, _ []float32) (*payload.Search_Response, error) {
	return nil, f.err
}

// Test_server_Search_concurrentCoreError calls Search concurrently with a faiss core error.
// The errors.As target used to be a package-level variable, which is a data race under -race.
func Test_server_Search_concurrentCoreError(t *testing.T) {
	t.Parallel()

	const callers = 32
	tests := []struct {
		err     error
		name    string
		wantMsg string
	}{
		// the core returns FaissError by value (errors.NewFaissError), which the As target must match.
		{
			name:    "faiss core error takes the core error branch",
			err:     errors.NewFaissError("boom"),
			wantMsg: "Search API failed to process search request due to faiss core process returned error",
		},
		{
			name:    "other error takes the default branch",
			err:     errors.New("boom"),
			wantMsg: "Search API failed to process search request",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			vec := []float32{0, 1}
			s := &server{faiss: &stubFaiss{err: test.err, dim: len(vec)}}
			req := &payload.Search_Request{Vector: vec}

			start := make(chan struct{})
			errs := make(chan error, callers)
			for range callers {
				go func() {
					<-start
					_, err := s.Search(context.Background(), req)
					errs <- err
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
					t.Fatal("timed out waiting for Search callers")
				}
			}
		})
	}
}
