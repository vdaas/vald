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
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/atomic"
	"github.com/vdaas/vald/pkg/gateway/mirror/service"
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
// vector, as a real vald gateway does) on every StreamListObject call.
type listAgentClient struct {
	service.MirrorClient
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

func Test_server_doStreamListObject_noReplayAfterPartialDelivery(t *testing.T) {
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
			s := &server{streamConcurrency: 1}
			ss := &listServerStream{ctx: ctx}
			var forwarded atomic.Bool

			// doStreamListObject is the function gateway.Do re-invokes on retry.
			err := retryLikeExecuteRPC(maxAttempts, func() error {
				return s.doStreamListObject(ctx, ac, &payload.Object_List_Request{}, ss, &forwarded)
			})
			if ctx.Err() != nil {
				tt.Fatalf("doStreamListObject did not finish before the test deadline: %v", ctx.Err())
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
