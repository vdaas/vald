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
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/io"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/atomic"
)

// listAgentStream yields n vectors and then fails every further Recv with failErr.
type listAgentStream struct {
	vald.Object_StreamListObjectClient
	failErr error
	n       int
	recvd   int
}

func (s *listAgentStream) Recv() (*payload.Object_List_Response, error) {
	if s.recvd >= s.n {
		return nil, s.failErr
	}
	s.recvd++
	return &payload.Object_List_Response{
		Payload: &payload.Object_List_Response_Vector{
			Vector: &payload.Object_Vector{Id: fmt.Sprintf("uuid-%d", s.recvd)},
		},
	}, nil
}

// listAgentClient opens a fresh listAgentStream (restarting from the first
// vector, as a real agent does) on every StreamListObject call.
type listAgentClient struct {
	vald.Client
	failErr error
	n       int
	opens   atomic.Int64
}

func (c *listAgentClient) StreamListObject(
	context.Context, *payload.Object_List_Request, ...grpc.CallOption,
) (vald.Object_StreamListObjectClient, error) {
	c.opens.Add(1)
	return &listAgentStream{n: c.n, failErr: c.failErr}, nil
}

// listServerStream records the ids forwarded to the caller.
type listServerStream struct {
	vald.Object_StreamListObjectServer
	ctx context.Context //nolint:containedctx // a fake grpc.ServerStream must return its stream context
	ids []string
	mu  sync.Mutex
}

func (s *listServerStream) Context() context.Context { return s.ctx }

func (s *listServerStream) Send(res *payload.Object_List_Response) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ids = append(s.ids, res.GetVector().GetId())
	return nil
}

// retryLikeExecuteRPC re-invokes f the way the backoff in internal/net/grpc
// gRPCClient.do does: executeRPC treats context errors as final and
// Internal/Unavailable/ResourceExhausted or non-status errors as retryable
// (for a healthy pool).
func retryLikeExecuteRPC(maxAttempts int, f func() error) (err error) {
	for range maxAttempts {
		err = f()
		if err == nil {
			return nil
		}
		if errors.IsAny(err, context.Canceled, context.DeadlineExceeded) {
			return err
		}
		st, ok := status.FromError(err)
		if !ok || st == nil {
			continue
		}
		if c := st.Code(); c == codes.Internal || c == codes.Unavailable || c == codes.ResourceExhausted {
			continue
		}
		return err
	}
	return err
}

func Test_server_StreamListObject_noReplayAfterPartialDelivery(t *testing.T) {
	t.Parallel()
	const maxAttempts = 3
	tests := []struct {
		name      string
		n         int
		wantCode  codes.Code
		wantOpens int64
		wantSent  int
	}{
		{
			name:      "failure after partial delivery is aborted and not replayed",
			n:         3,
			wantCode:  codes.Aborted,
			wantOpens: 1,
			wantSent:  3,
		},
		{
			name:      "failure before the first message keeps its retryable code",
			n:         0,
			wantCode:  codes.Unavailable,
			wantOpens: maxAttempts,
			wantSent:  0,
		},
	}
	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			// bounds the relay loop even if the stream never terminates.
			ctx, cancel := context.WithTimeout(tt.Context(), 10*time.Second)
			defer cancel()
			ac := &listAgentClient{
				n:       test.n,
				failErr: status.New(codes.Unavailable, "agent stream broken").Err(),
			}
			gw := &regGateway{
				broadCast: func(ctx context.Context, f func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error) error {
					return retryLikeExecuteRPC(maxAttempts, func() error {
						return f(ctx, "agent-0", ac)
					})
				},
			}
			s := &server{gateway: gw, streamConcurrency: 1}
			ss := &listServerStream{ctx: ctx}

			err := s.StreamListObject(&payload.Object_List_Request{}, ss)
			if ctx.Err() != nil {
				tt.Fatalf("StreamListObject did not finish before the test deadline: %v", ctx.Err())
			}
			st, ok := status.FromError(err)
			if !ok || st == nil || st.Code() != test.wantCode {
				tt.Errorf("error code = %v (err: %v), want %v", st.Code(), err, test.wantCode)
			}
			if got := ac.opens.Load(); got != test.wantOpens {
				tt.Errorf("agent stream opened %d times, want %d", got, test.wantOpens)
			}
			ss.mu.Lock()
			defer ss.mu.Unlock()
			if len(ss.ids) != test.wantSent {
				tt.Errorf("forwarded %d messages %v, want %d (no duplicates)", len(ss.ids), ss.ids, test.wantSent)
			}
		})
	}
}

// gatedListStream yields one vector once release is closed, then io.EOF.
// entered (if set) is closed when Recv is first entered.
type gatedListStream struct {
	vald.Object_StreamListObjectClient
	ctx     context.Context //nolint:containedctx // a fake client stream is bound to its call context
	release <-chan struct{}
	entered chan struct{}
	id      string
	once    sync.Once
	sent    bool
}

func (s *gatedListStream) Recv() (*payload.Object_List_Response, error) {
	if s.sent {
		return nil, io.EOF
	}
	if s.entered != nil {
		s.once.Do(func() { close(s.entered) })
	}
	select {
	case <-s.release:
	case <-s.ctx.Done():
		return nil, s.ctx.Err()
	}
	s.sent = true
	return &payload.Object_List_Response{
		Payload: &payload.Object_List_Response_Vector{
			Vector: &payload.Object_Vector{Id: s.id},
		},
	}, nil
}

type gatedListClient struct {
	vald.Client
	newStream func(ctx context.Context) *gatedListStream
}

func (c *gatedListClient) StreamListObject(
	ctx context.Context, _ *payload.Object_List_Request, _ ...grpc.CallOption,
) (vald.Object_StreamListObjectClient, error) {
	return c.newStream(ctx), nil
}

// notifyServerStream records the forwarded ids and closes onID's channel once the id is sent.
type notifyServerStream struct {
	vald.Object_StreamListObjectServer
	ctx  context.Context //nolint:containedctx // a fake grpc.ServerStream must return its stream context
	onID map[string]func()
	ids  []string
	mu   sync.Mutex
}

func (s *notifyServerStream) Context() context.Context { return s.ctx }

func (s *notifyServerStream) Send(res *payload.Object_List_Response) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	id := res.GetVector().GetId()
	s.ids = append(s.ids, id)
	if f := s.onID[id]; f != nil {
		f()
	}
	return nil
}

// A Recv blocked on one agent must not stop the other agents' streams. Agent a
// only yields once agent b's vector reached the caller, and agent b only yields
// once agent a is inside Recv; a Recv lock shared across targets deadlocks this
// until the deadline.
func Test_server_StreamListObject_targetsReceiveIndependently(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()

	aEntered := make(chan struct{})
	bSent := make(chan struct{})
	var bOnce sync.Once
	agentA := &gatedListClient{newStream: func(ctx context.Context) *gatedListStream {
		return &gatedListStream{ctx: ctx, id: "a", release: bSent, entered: aEntered}
	}}
	agentB := &gatedListClient{newStream: func(ctx context.Context) *gatedListStream {
		return &gatedListStream{ctx: ctx, id: "b", release: aEntered}
	}}
	gw := &regGateway{
		broadCast: func(ctx context.Context, f func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error) error {
			ech := make(chan error, 2)
			go func() { ech <- f(ctx, "agent-a", agentA) }()
			go func() { ech <- f(ctx, "agent-b", agentB) }()
			errs := make([]error, 0, 2)
			for range 2 {
				select {
				case err := <-ech:
					errs = append(errs, err)
				case <-time.After(30 * time.Second):
					return errors.New("broadcast callback did not return")
				}
			}
			return errors.Join(errs...)
		},
	}
	s := &server{gateway: gw, streamConcurrency: 2}
	ss := &notifyServerStream{
		ctx:  ctx,
		onID: map[string]func(){"b": func() { bOnce.Do(func() { close(bSent) }) }},
	}

	err := s.StreamListObject(&payload.Object_List_Request{}, ss)
	if ctx.Err() != nil {
		t.Fatalf("StreamListObject was blocked until the test deadline: %v", ctx.Err())
	}
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	ss.mu.Lock()
	defer ss.mu.Unlock()
	if len(ss.ids) != 2 || ss.ids[0] != "b" || ss.ids[1] != "a" {
		t.Errorf("forwarded %v, want [b a]", ss.ids)
	}
}
