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
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/atomic"
	"github.com/vdaas/vald/internal/sync/errgroup"
)

// seqGateway broadcasts to the given clients one after another and joins their errors.
func seqGateway(clients ...vald.Client) *regGateway {
	return &regGateway{
		agents: len(clients),
		broadCast: func(ctx context.Context, f func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error) error {
			errs := make([]error, 0, len(clients))
			for _, c := range clients {
				errs = append(errs, f(ctx, "127.0.0.1:8081", c))
			}
			return errors.Join(errs...)
		},
	}
}

func Test_broadcastFirst(t *testing.T) {
	t.Parallel()
	errBroadCast := errors.New("broadcast failed")
	miss := func(context.Context, *payload.Object_ID) (*payload.Object_ID, error) {
		return &payload.Object_ID{}, nil
	}
	hit := func(id string) func(context.Context, *payload.Object_ID) (*payload.Object_ID, error) {
		return func(context.Context, *payload.Object_ID) (*payload.Object_ID, error) {
			return &payload.Object_ID{Id: id}, nil
		}
	}
	tests := []struct {
		name    string
		gw      func() *regGateway
		wantErr error
		wantID  string
		wantOK  bool
	}{
		{
			name:   "first hit is returned",
			gw:     func() *regGateway { return seqGateway(&regClient{exists: miss}, &regClient{exists: hit("a")}) },
			wantID: "a",
			wantOK: true,
		},
		{
			name: "no hit is not an error",
			gw:   func() *regGateway { return seqGateway(&regClient{exists: miss}, &regClient{exists: miss}) },
		},
		{
			name: "broadcast error is reported when nothing was found",
			gw: func() *regGateway {
				return &regGateway{
					broadCast: func(context.Context, func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error) error {
						return errBroadCast
					},
				}
			},
			wantErr: errBroadCast,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
			defer cancel()
			eg, _ := errgroup.New(ctx)
			eg.SetLimit(1)
			s := &server{gateway: test.gw(), eg: eg}
			req := &payload.Object_ID{Id: "a"}
			got, ok, err := broadcastFirst(ctx, s, "exists",
				func(oid *payload.Object_ID) bool { return oid.GetId() != "" },
				func(ctx context.Context, _ trace.Span, _ string, vc vald.Client, copts ...grpc.CallOption) (*payload.Object_ID, error) {
					return vc.Exists(ctx, req, copts...)
				})
			if (test.wantErr == nil && err != nil) || (test.wantErr != nil && !errors.Is(err, test.wantErr)) {
				t.Errorf("err = %v, want %v", err, test.wantErr)
			}
			if ok != test.wantOK || got.GetId() != test.wantID {
				t.Errorf("got (%v, %t), want (%q, %t)", got, ok, test.wantID, test.wantOK)
			}
		})
	}
}

// exists/getObject/GetTimestamp build the read-only agent request once per call
// and share it across the broadcast targets instead of allocating one per target.
func Test_server_objectReads_shareRequestAcrossTargets(t *testing.T) {
	t.Parallel()
	const (
		uuid    = "uuid"
		targets = 3
	)
	tests := []struct {
		call func(ctx context.Context, record func(any)) error
		name string
	}{
		{
			name: "exists",
			call: func(ctx context.Context, record func(any)) error {
				c := &regClient{exists: func(_ context.Context, in *payload.Object_ID) (*payload.Object_ID, error) {
					record(in)
					return &payload.Object_ID{}, nil
				}}
				s := &server{gateway: seqGateway(c, c, c)}
				s.eg, _ = errgroup.New(ctx)
				_, err := s.exists(ctx, uuid)
				return err
			},
		},
		{
			name: "getObject",
			call: func(ctx context.Context, record func(any)) error {
				c := &regClient{getObject: func(_ context.Context, in *payload.Object_VectorRequest) (*payload.Object_Vector, error) {
					record(in)
					return &payload.Object_Vector{}, nil
				}}
				s := &server{gateway: seqGateway(c, c, c)}
				s.eg, _ = errgroup.New(ctx)
				_, err := s.getObject(ctx, uuid)
				return err
			},
		},
		{
			name: "GetTimestamp",
			call: func(ctx context.Context, record func(any)) error {
				c := &regClient{getTS: func(_ context.Context, in *payload.Object_TimestampRequest) (*payload.Object_Timestamp, error) {
					record(in)
					return &payload.Object_Timestamp{}, nil
				}}
				s := &server{gateway: seqGateway(c, c, c)}
				s.eg, _ = errgroup.New(ctx)
				_, err := s.GetTimestamp(ctx, &payload.Object_TimestampRequest{Id: &payload.Object_ID{Id: uuid}})
				return err
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
			defer cancel()
			var (
				mu   sync.Mutex
				seen []any
			)
			record := func(in any) {
				mu.Lock()
				defer mu.Unlock()
				seen = append(seen, in)
			}
			// every target misses (empty answer), so the call reports NotFound after all of them ran.
			if err := test.call(ctx, record); err == nil {
				t.Fatal("expected a not-found error, got nil")
			}
			if ctx.Err() != nil {
				t.Fatalf("call did not finish before the test deadline: %v", ctx.Err())
			}
			mu.Lock()
			defer mu.Unlock()
			if len(seen) != targets {
				t.Fatalf("agents called %d times, want %d", len(seen), targets)
			}
			for i, in := range seen[1:] {
				if in != seen[0] {
					t.Errorf("target %d got request %p, want the shared request %p", i+1, in, seen[0])
				}
			}
		})
	}
}

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

// exists/getObject used to return a spurious NotFound when the select picked the
// nil broadcast result from ech while the found object was still buffered in the
// result channel. The race is random, so repeat the call and require every one to succeed.
func Test_server_existsAndGetObject_foundObjectWinsOverNilBroadCast(t *testing.T) {
	t.Parallel()
	const (
		uuid  = "uuid"
		calls = 200
	)
	eg, _ := errgroup.New(t.Context())
	// limit 1 runs eg.Go inline, so the found result, the canceled ctx and the broadcast result on ech
	// are all ready before the select runs; the old code then picks the ech case at random and fails.
	eg.SetLimit(1)
	s := &server{gateway: syncGateway(newOKClient([]float32{1, 2, 3})), replica: 1, eg: eg}
	t.Run("exists", func(t *testing.T) {
		t.Parallel()
		for i := range calls {
			ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
			id, err := s.exists(ctx, uuid)
			cancel()
			if err != nil || id.GetId() != uuid {
				t.Fatalf("call %d: got (%v, %v), want id %q", i, id, err, uuid)
			}
		}
	})
	t.Run("getObject", func(t *testing.T) {
		t.Parallel()
		for i := range calls {
			ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
			vec, err := s.getObject(ctx, uuid)
			cancel()
			if err != nil || vec.GetId() != uuid {
				t.Fatalf("call %d: got (%v, %v), want id %q", i, vec, err, uuid)
			}
		}
	})
	t.Run("GetTimestamp", func(t *testing.T) {
		t.Parallel()
		c := newOKClient([]float32{1, 2, 3})
		c.getTS = func(context.Context, *payload.Object_TimestampRequest) (*payload.Object_Timestamp, error) {
			return &payload.Object_Timestamp{Id: uuid, Timestamp: 1}, nil
		}
		eg, _ := errgroup.New(t.Context())
		eg.SetLimit(1)
		ts := &server{gateway: syncGateway(c), replica: 1, eg: eg}
		req := &payload.Object_TimestampRequest{Id: &payload.Object_ID{Id: uuid}}
		for i := range calls {
			ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
			got, err := ts.GetTimestamp(ctx, req)
			cancel()
			if err != nil || got.GetId() != uuid {
				t.Fatalf("call %d: got (%v, %v), want id %q", i, got, err, uuid)
			}
		}
	})
}

// When nothing was found, exists/getObject/GetTimestamp used to pick the closed result channel
// or the broadcast error at random, turning a connection error into NotFound about half the time.
func Test_server_existsAndGetObject_broadCastErrorWinsOverClosedResult(t *testing.T) {
	t.Parallel()
	const (
		uuid  = "uuid"
		calls = 200
	)
	newServer := func() *server {
		eg, _ := errgroup.New(t.Context())
		// limit 1 runs eg.Go inline, so the closed result channel and the buffered broadcast
		// error are both ready before the select runs.
		eg.SetLimit(1)
		return &server{
			gateway: &regGateway{
				broadCast: func(context.Context, func(context.Context, string, vald.Client, ...grpc.CallOption) error) error {
					return errors.ErrGRPCClientConnNotFound("*")
				},
			},
			replica: 1,
			eg:      eg,
		}
	}
	t.Run("exists", func(t *testing.T) {
		t.Parallel()
		s := newServer()
		for i := range calls {
			if _, err := s.exists(t.Context(), uuid); !errors.Is(err, errors.ErrGRPCClientConnNotFound("*")) {
				t.Fatalf("call %d: got %v, want connection not found", i, err)
			}
		}
	})
	t.Run("getObject", func(t *testing.T) {
		t.Parallel()
		s := newServer()
		for i := range calls {
			if _, err := s.getObject(t.Context(), uuid); !errors.Is(err, errors.ErrGRPCClientConnNotFound("*")) {
				t.Fatalf("call %d: got %v, want connection not found", i, err)
			}
		}
	})
	t.Run("GetTimestamp", func(t *testing.T) {
		t.Parallel()
		s := newServer()
		req := &payload.Object_TimestampRequest{Id: &payload.Object_ID{Id: uuid}}
		for i := range calls {
			_, err := s.GetTimestamp(t.Context(), req)
			if st, _ := status.FromError(err); st.Code() != codes.Internal {
				t.Fatalf("call %d: got %v, want code Internal", i, err)
			}
		}
	})
}

// exists/getObject/GetTimestamp derive a cancelable context that was only canceled once an object was found;
// every other return path left it registered with the parent. The ctx handed to BroadCast must be
// canceled by the time they return.
func Test_server_existsAndGetObject_cancelDerivedContext(t *testing.T) {
	t.Parallel()
	const uuid = "uuid"
	run := func(t *testing.T, call func(s *server, ctx context.Context) error) {
		t.Helper()
		var bctx context.Context
		c := newOKClient(nil)
		// nothing is found, so the call returns through the broadcast-result path
		c.exists = func(context.Context, *payload.Object_ID) (*payload.Object_ID, error) {
			return &payload.Object_ID{}, nil
		}
		c.getObject = func(context.Context, *payload.Object_VectorRequest) (*payload.Object_Vector, error) {
			return &payload.Object_Vector{}, nil
		}
		c.getTS = func(context.Context, *payload.Object_TimestampRequest) (*payload.Object_Timestamp, error) {
			return &payload.Object_Timestamp{}, nil
		}
		gw := &regGateway{
			agents: 1,
			broadCast: func(ctx context.Context, f func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error) error {
				bctx = ctx
				return f(ctx, "127.0.0.1:8081", c)
			},
		}
		eg, _ := errgroup.New(t.Context())
		s := &server{gateway: gw, replica: 1, eg: eg}

		ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
		defer cancel()
		if err := call(s, ctx); err == nil {
			t.Fatal("expected a not found error")
		}
		if bctx == nil {
			t.Fatal("BroadCast was not called")
		}
		if bctx.Err() == nil {
			t.Error("the derived context was not canceled when the call returned")
		}
	}
	t.Run("exists", func(t *testing.T) {
		t.Parallel()
		run(t, func(s *server, ctx context.Context) error {
			_, err := s.exists(ctx, uuid)
			return err
		})
	})
	t.Run("getObject", func(t *testing.T) {
		t.Parallel()
		run(t, func(s *server, ctx context.Context) error {
			_, err := s.getObject(ctx, uuid)
			return err
		})
	})
	t.Run("GetTimestamp", func(t *testing.T) {
		t.Parallel()
		run(t, func(s *server, ctx context.Context) error {
			_, err := s.GetTimestamp(ctx, &payload.Object_TimestampRequest{Id: &payload.Object_ID{Id: uuid}})
			return err
		})
	})
}

// NOT IMPLEMENTED BELOW
//
// func Test_server_exists(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		uuid string
// 	}
// 	type fields struct {
// 		UnimplementedValdServer vald.UnimplementedValdServer
// 		eg                      errgroup.Group
// 		gateway                 service.Gateway
// 		name                    string
// 		ip                      string
// 		timeout                 time.Duration
// 		replica                 int
// 		streamConcurrency       int
// 		multiConcurrency        int
// 	}
// 	type want struct {
// 		wantId *payload.Object_ID
// 		err    error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Object_ID, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotId *payload.Object_ID, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotId, w.wantId) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotId, w.wantId)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           uuid:"",
// 		       },
// 		       fields: fields {
// 		           UnimplementedValdServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           name:"",
// 		           ip:"",
// 		           timeout:nil,
// 		           replica:0,
// 		           streamConcurrency:0,
// 		           multiConcurrency:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           uuid:"",
// 		           },
// 		           fields: fields {
// 		           UnimplementedValdServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           name:"",
// 		           ip:"",
// 		           timeout:nil,
// 		           replica:0,
// 		           streamConcurrency:0,
// 		           multiConcurrency:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			s := &server{
// 				UnimplementedValdServer: test.fields.UnimplementedValdServer,
// 				eg:                      test.fields.eg,
// 				gateway:                 test.fields.gateway,
// 				name:                    test.fields.name,
// 				ip:                      test.fields.ip,
// 				timeout:                 test.fields.timeout,
// 				replica:                 test.fields.replica,
// 				streamConcurrency:       test.fields.streamConcurrency,
// 				multiConcurrency:        test.fields.multiConcurrency,
// 			}
//
// 			gotId, err := s.exists(test.args.ctx, test.args.uuid)
// 			if err := checkFunc(test.want, gotId, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_Exists(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		meta *payload.Object_ID
// 	}
// 	type fields struct {
// 		UnimplementedValdServer vald.UnimplementedValdServer
// 		eg                      errgroup.Group
// 		gateway                 service.Gateway
// 		name                    string
// 		ip                      string
// 		timeout                 time.Duration
// 		replica                 int
// 		streamConcurrency       int
// 		multiConcurrency        int
// 	}
// 	type want struct {
// 		wantId *payload.Object_ID
// 		err    error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Object_ID, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotId *payload.Object_ID, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotId, w.wantId) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotId, w.wantId)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           meta:nil,
// 		       },
// 		       fields: fields {
// 		           UnimplementedValdServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           name:"",
// 		           ip:"",
// 		           timeout:nil,
// 		           replica:0,
// 		           streamConcurrency:0,
// 		           multiConcurrency:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           meta:nil,
// 		           },
// 		           fields: fields {
// 		           UnimplementedValdServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           name:"",
// 		           ip:"",
// 		           timeout:nil,
// 		           replica:0,
// 		           streamConcurrency:0,
// 		           multiConcurrency:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			s := &server{
// 				UnimplementedValdServer: test.fields.UnimplementedValdServer,
// 				eg:                      test.fields.eg,
// 				gateway:                 test.fields.gateway,
// 				name:                    test.fields.name,
// 				ip:                      test.fields.ip,
// 				timeout:                 test.fields.timeout,
// 				replica:                 test.fields.replica,
// 				streamConcurrency:       test.fields.streamConcurrency,
// 				multiConcurrency:        test.fields.multiConcurrency,
// 			}
//
// 			gotId, err := s.Exists(test.args.ctx, test.args.meta)
// 			if err := checkFunc(test.want, gotId, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_getObject(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		uuid string
// 	}
// 	type fields struct {
// 		UnimplementedValdServer vald.UnimplementedValdServer
// 		eg                      errgroup.Group
// 		gateway                 service.Gateway
// 		name                    string
// 		ip                      string
// 		timeout                 time.Duration
// 		replica                 int
// 		streamConcurrency       int
// 		multiConcurrency        int
// 	}
// 	type want struct {
// 		wantVec *payload.Object_Vector
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Object_Vector, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotVec *payload.Object_Vector, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotVec, w.wantVec) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotVec, w.wantVec)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           uuid:"",
// 		       },
// 		       fields: fields {
// 		           UnimplementedValdServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           name:"",
// 		           ip:"",
// 		           timeout:nil,
// 		           replica:0,
// 		           streamConcurrency:0,
// 		           multiConcurrency:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           uuid:"",
// 		           },
// 		           fields: fields {
// 		           UnimplementedValdServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           name:"",
// 		           ip:"",
// 		           timeout:nil,
// 		           replica:0,
// 		           streamConcurrency:0,
// 		           multiConcurrency:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			s := &server{
// 				UnimplementedValdServer: test.fields.UnimplementedValdServer,
// 				eg:                      test.fields.eg,
// 				gateway:                 test.fields.gateway,
// 				name:                    test.fields.name,
// 				ip:                      test.fields.ip,
// 				timeout:                 test.fields.timeout,
// 				replica:                 test.fields.replica,
// 				streamConcurrency:       test.fields.streamConcurrency,
// 				multiConcurrency:        test.fields.multiConcurrency,
// 			}
//
// 			gotVec, err := s.getObject(test.args.ctx, test.args.uuid)
// 			if err := checkFunc(test.want, gotVec, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_GetObject(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		req *payload.Object_VectorRequest
// 	}
// 	type fields struct {
// 		UnimplementedValdServer vald.UnimplementedValdServer
// 		eg                      errgroup.Group
// 		gateway                 service.Gateway
// 		name                    string
// 		ip                      string
// 		timeout                 time.Duration
// 		replica                 int
// 		streamConcurrency       int
// 		multiConcurrency        int
// 	}
// 	type want struct {
// 		wantVec *payload.Object_Vector
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Object_Vector, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotVec *payload.Object_Vector, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotVec, w.wantVec) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotVec, w.wantVec)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           req:nil,
// 		       },
// 		       fields: fields {
// 		           UnimplementedValdServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           name:"",
// 		           ip:"",
// 		           timeout:nil,
// 		           replica:0,
// 		           streamConcurrency:0,
// 		           multiConcurrency:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           req:nil,
// 		           },
// 		           fields: fields {
// 		           UnimplementedValdServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           name:"",
// 		           ip:"",
// 		           timeout:nil,
// 		           replica:0,
// 		           streamConcurrency:0,
// 		           multiConcurrency:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			s := &server{
// 				UnimplementedValdServer: test.fields.UnimplementedValdServer,
// 				eg:                      test.fields.eg,
// 				gateway:                 test.fields.gateway,
// 				name:                    test.fields.name,
// 				ip:                      test.fields.ip,
// 				timeout:                 test.fields.timeout,
// 				replica:                 test.fields.replica,
// 				streamConcurrency:       test.fields.streamConcurrency,
// 				multiConcurrency:        test.fields.multiConcurrency,
// 			}
//
// 			gotVec, err := s.GetObject(test.args.ctx, test.args.req)
// 			if err := checkFunc(test.want, gotVec, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_StreamGetObject(t *testing.T) {
// 	type args struct {
// 		stream vald.Object_StreamGetObjectServer
// 	}
// 	type fields struct {
// 		UnimplementedValdServer vald.UnimplementedValdServer
// 		eg                      errgroup.Group
// 		gateway                 service.Gateway
// 		name                    string
// 		ip                      string
// 		timeout                 time.Duration
// 		replica                 int
// 		streamConcurrency       int
// 		multiConcurrency        int
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           stream:nil,
// 		       },
// 		       fields: fields {
// 		           UnimplementedValdServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           name:"",
// 		           ip:"",
// 		           timeout:nil,
// 		           replica:0,
// 		           streamConcurrency:0,
// 		           multiConcurrency:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           stream:nil,
// 		           },
// 		           fields: fields {
// 		           UnimplementedValdServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           name:"",
// 		           ip:"",
// 		           timeout:nil,
// 		           replica:0,
// 		           streamConcurrency:0,
// 		           multiConcurrency:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			s := &server{
// 				UnimplementedValdServer: test.fields.UnimplementedValdServer,
// 				eg:                      test.fields.eg,
// 				gateway:                 test.fields.gateway,
// 				name:                    test.fields.name,
// 				ip:                      test.fields.ip,
// 				timeout:                 test.fields.timeout,
// 				replica:                 test.fields.replica,
// 				streamConcurrency:       test.fields.streamConcurrency,
// 				multiConcurrency:        test.fields.multiConcurrency,
// 			}
//
// 			err := s.StreamGetObject(test.args.stream)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_StreamListObject(t *testing.T) {
// 	type args struct {
// 		req    *payload.Object_List_Request
// 		stream vald.Object_StreamListObjectServer
// 	}
// 	type fields struct {
// 		UnimplementedValdServer vald.UnimplementedValdServer
// 		eg                      errgroup.Group
// 		gateway                 service.Gateway
// 		name                    string
// 		ip                      string
// 		timeout                 time.Duration
// 		replica                 int
// 		streamConcurrency       int
// 		multiConcurrency        int
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           req:nil,
// 		           stream:nil,
// 		       },
// 		       fields: fields {
// 		           UnimplementedValdServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           name:"",
// 		           ip:"",
// 		           timeout:nil,
// 		           replica:0,
// 		           streamConcurrency:0,
// 		           multiConcurrency:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           req:nil,
// 		           stream:nil,
// 		           },
// 		           fields: fields {
// 		           UnimplementedValdServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           name:"",
// 		           ip:"",
// 		           timeout:nil,
// 		           replica:0,
// 		           streamConcurrency:0,
// 		           multiConcurrency:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			s := &server{
// 				UnimplementedValdServer: test.fields.UnimplementedValdServer,
// 				eg:                      test.fields.eg,
// 				gateway:                 test.fields.gateway,
// 				name:                    test.fields.name,
// 				ip:                      test.fields.ip,
// 				timeout:                 test.fields.timeout,
// 				replica:                 test.fields.replica,
// 				streamConcurrency:       test.fields.streamConcurrency,
// 				multiConcurrency:        test.fields.multiConcurrency,
// 			}
//
// 			err := s.StreamListObject(test.args.req, test.args.stream)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_GetTimestamp(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		req *payload.Object_TimestampRequest
// 	}
// 	type fields struct {
// 		UnimplementedValdServer vald.UnimplementedValdServer
// 		eg                      errgroup.Group
// 		gateway                 service.Gateway
// 		name                    string
// 		ip                      string
// 		timeout                 time.Duration
// 		replica                 int
// 		streamConcurrency       int
// 		multiConcurrency        int
// 	}
// 	type want struct {
// 		wantTs *payload.Object_Timestamp
// 		err    error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Object_Timestamp, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotTs *payload.Object_Timestamp, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotTs, w.wantTs) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotTs, w.wantTs)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           req:nil,
// 		       },
// 		       fields: fields {
// 		           UnimplementedValdServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           name:"",
// 		           ip:"",
// 		           timeout:nil,
// 		           replica:0,
// 		           streamConcurrency:0,
// 		           multiConcurrency:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           req:nil,
// 		           },
// 		           fields: fields {
// 		           UnimplementedValdServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           name:"",
// 		           ip:"",
// 		           timeout:nil,
// 		           replica:0,
// 		           streamConcurrency:0,
// 		           multiConcurrency:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			s := &server{
// 				UnimplementedValdServer: test.fields.UnimplementedValdServer,
// 				eg:                      test.fields.eg,
// 				gateway:                 test.fields.gateway,
// 				name:                    test.fields.name,
// 				ip:                      test.fields.ip,
// 				timeout:                 test.fields.timeout,
// 				replica:                 test.fields.replica,
// 				streamConcurrency:       test.fields.streamConcurrency,
// 				multiConcurrency:        test.fields.multiConcurrency,
// 			}
//
// 			gotTs, err := s.GetTimestamp(test.args.ctx, test.args.req)
// 			if err := checkFunc(test.want, gotTs, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
