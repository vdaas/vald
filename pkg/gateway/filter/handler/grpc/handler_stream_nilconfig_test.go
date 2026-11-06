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
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	client "github.com/vdaas/vald/internal/client/v1/client/vald"
	"github.com/vdaas/vald/internal/io"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/sync"
)

// linearSearchGatewayStub answers LinearSearch with a fixed response.
type linearSearchGatewayStub struct {
	client.Client
	res *payload.Search_Response
}

func (g *linearSearchGatewayStub) LinearSearch(
	context.Context, *payload.Search_Request, ...grpc.CallOption,
) (*payload.Search_Response, error) {
	return g.res, nil
}

// linearSearchStreamStub feeds reqs to Recv (then io.EOF) and records Send.
type linearSearchStreamStub struct {
	vald.Search_StreamLinearSearchServer
	ctx  func() context.Context
	reqs []*payload.Search_Request
	sent []*payload.Search_StreamResponse
	mu   sync.Mutex
}

func (s *linearSearchStreamStub) Context() context.Context { return s.ctx() }

func (s *linearSearchStreamStub) Recv() (*payload.Search_Request, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.reqs) == 0 {
		return nil, io.EOF
	}
	req := s.reqs[0]
	s.reqs = s.reqs[1:]
	return req, nil
}

func (s *linearSearchStreamStub) Send(res *payload.Search_StreamResponse) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sent = append(s.sent, res)
	return nil
}

// TestStreamLinearSearch_nilConfig verifies that a stream message without a
// Config (legal in proto3) is served instead of nil-dereferencing
// req.Config.RequestId, which panicked inside the stream worker and ended the
// whole stream with an error.
func TestStreamLinearSearch_nilConfig(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()
	want := &payload.Search_Response{Results: []*payload.Object_Distance{{Id: "a", Distance: 0.1}}}
	s := &server{gateway: &linearSearchGatewayStub{res: want}}
	stream := &linearSearchStreamStub{
		ctx:  func() context.Context { return ctx },
		reqs: []*payload.Search_Request{{Vector: []float32{1, 2}}},
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- s.StreamLinearSearch(stream)
	}()
	select {
	case err := <-errCh:
		if err != nil {
			t.Fatalf("StreamLinearSearch() error = %v, want nil", err)
		}
	case <-ctx.Done():
		t.Fatal("timed out waiting for StreamLinearSearch")
	}

	stream.mu.Lock()
	defer stream.mu.Unlock()
	if len(stream.sent) != 1 {
		t.Fatalf("sent %d responses, want 1", len(stream.sent))
	}
	if got := stream.sent[0].GetResponse(); got != want {
		t.Errorf("response = %v, want %v", got, want)
	}
}
