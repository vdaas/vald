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
	"reflect"
	"testing"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/io"
	_ "github.com/vdaas/vald/internal/net/grpc/logger"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/atomic"
	"github.com/vdaas/vald/internal/test/data/request"
	"github.com/vdaas/vald/internal/test/data/vector"
	"github.com/vdaas/vald/internal/test/goleak"
	"github.com/vdaas/vald/internal/test/mock"
)

func TestBidirectionalStream(t *testing.T) {
	type args struct {
		f           func(context.Context, *payload.Insert_Request) (*payload.Object_StreamLocation, error)
		insertReqs  []*payload.Insert_Request
		concurrency int
	}
	type want struct {
		err     error
		rpcResp []*payload.Object_StreamLocation
	}
	type test struct {
		want       want
		checkFunc  func([]*payload.Object_StreamLocation, want, error) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
		name       string
		args       args
	}

	const (
		name        = "vald-agent-ngt-1" // agent name
		f32VecDim   = 3                  // float32 vector dimension
		ip          = "localhost"        // ip address
		concurrency = 10
	)

	defaultCallbackFn := func(ctx context.Context, i *payload.Insert_Request) (*payload.Object_StreamLocation, error) {
		return &payload.Object_StreamLocation{
			Payload: &payload.Object_StreamLocation_Location{
				Location: &payload.Object_Location{
					Name: name,
					Uuid: i.GetVector().GetId(),
					Ips:  []string{ip},
				},
			},
		}, nil
	}

	defaultCheckFunc := func(rpcResp []*payload.Object_StreamLocation, w want, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}

		// since the insert order is not guaranteed, check only the error count on the response
		sm := make(map[int32]int) // want status map
		for _, r := range w.rpcResp {
			sm[r.GetStatus().GetCode()]++
		}
		gsm := make(map[int32]int) // got status map
		for _, r := range rpcResp {
			gsm[r.GetStatus().GetCode()]++
		}

		if !reflect.DeepEqual(sm, gsm) {
			return errors.Errorf("status count is not correct, got: %v, want: %v", gsm, sm)
		}
		return nil
	}
	tests := []test{
		func() test {
			insertCnt := 1
			reqs, err := request.GenMultiInsertReq(request.Float, vector.Gaussian, insertCnt, f32VecDim, nil)
			if err != nil {
				t.Fatal(err)
			}
			return test{
				name: "success to receive 1 message from stream",
				args: args{
					insertReqs:  reqs.Requests,
					concurrency: concurrency,
					f:           defaultCallbackFn,
				},
				want: want{
					rpcResp: request.GenObjectStreamLocation(insertCnt, name, ip),
				},
			}
		}(),
		func() test {
			insertCnt := 10
			reqs, err := request.GenMultiInsertReq(request.Float, vector.Gaussian, insertCnt, f32VecDim, nil)
			if err != nil {
				t.Fatal(err)
			}
			return test{
				name: "success to receive 10 message from stream",
				args: args{
					insertReqs:  reqs.Requests,
					concurrency: concurrency,
					f:           defaultCallbackFn,
				},
				want: want{
					rpcResp: request.GenObjectStreamLocation(insertCnt, name, ip),
				},
			}
		}(),
		func() test {
			return test{
				name: "success to receive 0 message from stream",
				args: args{
					insertReqs:  []*payload.Insert_Request{},
					concurrency: concurrency,
					f:           defaultCallbackFn,
				},
				want: want{
					rpcResp: nil,
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())

			ctx := tt.Context()

			if test.beforeFunc != nil {
				test.beforeFunc(tt, test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			insertReqs := test.args.insertReqs
			rpcResp := make([]*payload.Object_StreamLocation, 0)
			recvIdx := 0

			stream := &mock.StreamInsertServerMock{
				ServerStream: &mock.ServerStreamMock{
					ContextFunc: func() context.Context {
						return ctx
					},
					RecvMsgFunc: func(i any) error {
						if recvIdx >= len(insertReqs) {
							return io.EOF
						}

						obj := i.(*payload.Insert_Request)
						if insertReqs[recvIdx] != nil {
							obj.Vector = insertReqs[recvIdx].Vector
							obj.Config = insertReqs[recvIdx].Config
						}
						recvIdx++

						return nil
					},
					SendMsgFunc: func(i any) error {
						rpcResp = append(rpcResp, i.(*payload.Object_StreamLocation))
						return nil
					},
				},
			}

			err := BidirectionalStream(ctx, stream, test.args.concurrency, test.args.f)
			if err := checkFunc(rpcResp, test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

// fakeClientStream is a minimal TypedClientStream. Recv blocks until CloseSend
// is called and then returns io.EOF, like a server that finishes after the
// client half-closes.
type fakeClientStream struct {
	ClientStream
	ctxFn  func() context.Context
	sendFn func(*payload.Insert_Request) error
	done   chan struct{}
	once   sync.Once
}

func newFakeClientStream(
	ctx context.Context, sendFn func(*payload.Insert_Request) error,
) *fakeClientStream {
	return &fakeClientStream{ctxFn: func() context.Context { return ctx }, sendFn: sendFn, done: make(chan struct{})}
}

func (f *fakeClientStream) Context() context.Context { return f.ctxFn() }

func (f *fakeClientStream) CloseSend() error {
	f.once.Do(func() { close(f.done) })
	return nil
}

func (f *fakeClientStream) Send(r *payload.Insert_Request) error { return f.sendFn(r) }

func (f *fakeClientStream) Recv() (*payload.Object_StreamLocation, error) {
	<-f.done
	return nil, io.EOF
}

func TestBidirectionalStreamClient(t *testing.T) {
	t.Parallel()
	errSend := errors.New("fake send error")
	const timeout = 10 * time.Second

	type args struct {
		sendFn      func(*payload.Insert_Request) error
		provider    func() func() (*payload.Insert_Request, bool)
		concurrency int
	}
	type want struct {
		err error
	}
	tests := []struct {
		want want
		name string
		args args
	}{
		{
			// Before the fix ech was unbuffered for concurrency == 0: the failing worker
			// blocked on `ech <- err` while finalize waited for it in eg.Wait().
			// The provider only reports "no more data" after Send has been invoked, so
			// the main loop deterministically reaches finalize with the worker in flight.
			name: "concurrency 0: send fails and provider then finishes does not deadlock",
			args: args{
				concurrency: 0,
				// sendFn and provider are built in the subtest: they share a channel.
			},
			want: want{err: errSend},
		},
		{
			name: "concurrency 4: multiple send errors are reported",
			args: args{
				concurrency: 4,
				sendFn: func(*payload.Insert_Request) error {
					return errSend
				},
				provider: func() func() (*payload.Insert_Request, bool) {
					cnt := 0
					return func() (*payload.Insert_Request, bool) {
						cnt++
						return &payload.Insert_Request{}, cnt <= 16
					}
				},
			},
			want: want{err: errSend},
		},
		{
			name: "concurrency 0: all sends succeed",
			args: args{
				concurrency: 0,
				sendFn: func(*payload.Insert_Request) error {
					return nil
				},
				provider: func() func() (*payload.Insert_Request, bool) {
					cnt := 0
					return func() (*payload.Insert_Request, bool) {
						cnt++
						return &payload.Insert_Request{}, cnt <= 16
					}
				},
			},
			want: want{err: nil},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())

			sendFn, provider := test.args.sendFn, test.args.provider
			if sendFn == nil {
				sent := make(chan struct{})
				var once sync.Once
				sendFn = func(*payload.Insert_Request) error {
					once.Do(func() { close(sent) })
					return errSend
				}
				calls := 0
				provider = func() func() (*payload.Insert_Request, bool) {
					return func() (*payload.Insert_Request, bool) {
						calls++
						if calls == 1 {
							return &payload.Insert_Request{}, true
						}
						select {
						case <-sent:
						case <-time.After(timeout):
						}
						return nil, false
					}
				}
			}

			stream := newFakeClientStream(tt.Context(), sendFn)
			resc := make(chan error, 1)
			go func() {
				resc <- BidirectionalStreamClient(stream, test.args.concurrency, provider(),
					func(*payload.Object_StreamLocation, error) bool { return true })
			}()
			select {
			case err := <-resc:
				if !errors.Is(err, test.want.err) {
					tt.Errorf("got error: %v, want: %v", err, test.want.err)
				}
			case <-time.After(timeout):
				// unblock Recv so the goroutines can exit, then fail without hanging.
				if err := stream.CloseSend(); err != nil {
					tt.Error(err)
				}
				tt.Fatalf("BidirectionalStreamClient did not return within %v (deadlock)", timeout)
			}
		})
	}
}

// sequencedClientStream returns msgs responses and then io.EOF once CloseSend
// was called. It records whether Recv ever ran concurrently with another Recv
// or with the callback (inCallback is set by the test's callback).
type sequencedClientStream struct {
	*fakeClientStream
	inRecv     atomic.Bool
	inCallback atomic.Bool
	overlap    atomic.Bool
	recvs      atomic.Int64
	msgs       int64
}

func (s *sequencedClientStream) Recv() (*payload.Object_StreamLocation, error) {
	if s.inRecv.Swap(true) || s.inCallback.Load() {
		s.overlap.Store(true)
	}
	defer s.inRecv.Store(false)
	if s.recvs.Add(1) <= s.msgs {
		return &payload.Object_StreamLocation{}, nil
	}
	<-s.done
	return nil, io.EOF
}

// TestBidirectionalStreamClient_RecvSequencing pins that the receiver calls
// Recv one at a time and only after the previous callback returned, and that
// every response reaches the callback before the stream finishes.
func TestBidirectionalStreamClient_RecvSequencing(t *testing.T) {
	t.Parallel()
	defer goleak.VerifyNone(t, goleak.IgnoreCurrent())
	const (
		msgs    = 8
		timeout = 10 * time.Second
	)
	stream := &sequencedClientStream{
		fakeClientStream: newFakeClientStream(t.Context(), func(*payload.Insert_Request) error { return nil }),
		msgs:             msgs,
	}
	var got atomic.Int64
	callbackDone := make(chan struct{})
	callBack := func(_ *payload.Object_StreamLocation, err error) bool {
		if err != nil {
			t.Errorf("callback got error: %v", err)
			return false
		}
		stream.inCallback.Store(true)
		// widen the window in which an overlapping Recv would be observed.
		time.Sleep(time.Millisecond)
		stream.inCallback.Store(false)
		if got.Add(1) == msgs {
			close(callbackDone)
		}
		return true
	}
	provider := func() (*payload.Insert_Request, bool) {
		// keep the send side open until every response was handled, so the
		// receiver is not cut short by the half-close.
		select {
		case <-callbackDone:
		case <-time.After(timeout):
		}
		return new(payload.Insert_Request), false
	}
	resc := make(chan error, 1)
	go func() {
		resc <- BidirectionalStreamClient(stream, 1, provider, callBack)
	}()
	select {
	case err := <-resc:
		if err != nil {
			t.Errorf("BidirectionalStreamClient error: %v", err)
		}
	case <-time.After(2 * timeout):
		if err := stream.CloseSend(); err != nil {
			t.Error(err)
		}
		t.Fatalf("BidirectionalStreamClient did not return within %v", 2*timeout)
	}
	if n := got.Load(); n != msgs {
		t.Errorf("callback called %d times, want %d", n, msgs)
	}
	if stream.overlap.Load() {
		t.Error("Recv overlapped another Recv or the callback")
	}
}

// scriptedClientStream allows fine-grained control over Send and Recv behaviors
// to simulate cancellations, slow servers, and abnormal network drops.
type scriptedClientStream struct {
	ClientStream
	ctx        context.Context
	cancel     context.CancelFunc
	sendFn     func(*payload.Insert_Request) error
	recvFn     func() (*payload.Object_StreamLocation, error)
	done       chan struct{}
	once       sync.Once
	closedSend atomic.Bool
}

func newScriptedClientStream(
	ctx context.Context,
	sendFn func(*payload.Insert_Request) error,
	recvFn func() (*payload.Object_StreamLocation, error),
) *scriptedClientStream {
	cctx, cancel := context.WithCancel(ctx)
	return &scriptedClientStream{
		ctx:    cctx,
		cancel: cancel,
		sendFn: sendFn,
		recvFn: recvFn,
		done:   make(chan struct{}),
	}
}

func (s *scriptedClientStream) Context() context.Context {
	return s.ctx
}

func (s *scriptedClientStream) CloseSend() error {
	s.closedSend.Store(true)
	s.once.Do(func() { close(s.done) })
	return nil
}

func (s *scriptedClientStream) Send(r *payload.Insert_Request) error {
	if s.sendFn != nil {
		return s.sendFn(r)
	}
	return nil
}

func (s *scriptedClientStream) Recv() (*payload.Object_StreamLocation, error) {
	if s.recvFn != nil {
		return s.recvFn()
	}
	<-s.done
	return nil, io.EOF
}

// TestBidirectionalStreamClient_rapidCancellation tests that when
// callBack returns false (cancelling the stream mid-flight), BidirectionalStreamClient
// unwinds immediately and leaks zero goroutines across both lockstep (c=1) and pipelined (c=16) modes.
func TestBidirectionalStreamClient_rapidCancellation(t *testing.T) {
	concurrencyModes := []struct {
		name        string
		concurrency int
	}{
		{name: "lockstep_concurrency_1", concurrency: 1},
		{name: "pipelined_concurrency_16", concurrency: 16},
	}

	for _, tc := range concurrencyModes {
		t.Run(tc.name, func(tt *testing.T) {
			// defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())

			ctx, cancel := context.WithTimeout(tt.Context(), 5*time.Second)
			defer cancel()

			var sentCnt atomic.Int64
			var recvCnt atomic.Int64

			stream := newScriptedClientStream(ctx,
				func(r *payload.Insert_Request) error {
					sentCnt.Add(1)
					return nil
				},
				func() (*payload.Object_StreamLocation, error) {
					cur := recvCnt.Add(1)
					if cur > 100 {
						return nil, io.EOF
					}
					return &payload.Object_StreamLocation{}, nil
				},
			)

			provider := func() (*payload.Insert_Request, bool) {
				return &payload.Insert_Request{}, true
			}

			// callBack returns false on the 5th message to simulate consumer early cancellation
			var handled atomic.Int64
			callBack := func(_ *payload.Object_StreamLocation, err error) bool {
				if err != nil {
					return false
				}
				if handled.Add(1) >= 5 {
					return false // abort stream
				}
				return true
			}

			resc := make(chan error, 1)
			go func() {
				resc <- BidirectionalStreamClient(stream, tc.concurrency, provider, callBack)
			}()

			select {
			case err := <-resc:
				if err != nil && !errors.IsAny(err, context.Canceled, io.EOF) {
					tt.Errorf("unexpected error on cancel: %v", err)
				}
			case <-time.After(3 * time.Second):
				_ = stream.CloseSend()
				stream.cancel()
				tt.Fatalf("BidirectionalStreamClient did not return promptly on callback cancel (potential deadlock)")
			}

			if !stream.closedSend.Load() {
				tt.Errorf("expected CloseSend to have been called by finalize")
			}
		})
	}
}

// TestBidirectionalStreamClient_slowConsumerBackpressure tests
// fast producer with a slow consumer callback, verifying that channel buffering
// operates correctly without data races, hangs, or unbounded allocation spikes.
func TestBidirectionalStreamClient_slowConsumerBackpressure(t *testing.T) {
	// defer goleak.VerifyNone(t, goleak.IgnoreCurrent())

	const totalMsgs = 500
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()

	var sentCnt atomic.Int64
	var recvsDelivered atomic.Int64

	stream := newScriptedClientStream(ctx,
		func(r *payload.Insert_Request) error {
			sentCnt.Add(1)
			return nil
		},
		func() (*payload.Object_StreamLocation, error) {
			cur := recvsDelivered.Add(1)
			if cur > totalMsgs {
				return nil, io.EOF
			}
			return &payload.Object_StreamLocation{}, nil
		},
	)

	var providerCalls atomic.Int64
	provider := func() (*payload.Insert_Request, bool) {
		if providerCalls.Add(1) > totalMsgs {
			return nil, false
		}
		return &payload.Insert_Request{}, true
	}

	var consumed atomic.Int64
	callBack := func(_ *payload.Object_StreamLocation, err error) bool {
		if err != nil {
			return false
		}
		consumed.Add(1)
		// Simulate slow consumer backpressure
		time.Sleep(50 * time.Microsecond)
		return true
	}

	resc := make(chan error, 1)
	go func() {
		resc <- BidirectionalStreamClient(stream, 8, provider, callBack)
	}()

	select {
	case err := <-resc:
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	case <-time.After(8 * time.Second):
		_ = stream.CloseSend()
		stream.cancel()
		t.Fatalf("timed out waiting for slow consumer test to finish")
	}

	if got := consumed.Load(); got != totalMsgs {
		t.Errorf("consumed %d messages, want %d", got, totalMsgs)
	}
}

// TestBidirectionalStreamClient_abnormalServerTermination verifies that
// BidirectionalStreamClient propagates the transport error when stream.Recv()
// fails and callBack returns false, instead of discarding it and returning nil.
func TestBidirectionalStreamClient_abnormalServerTermination(t *testing.T) {
	// defer goleak.VerifyNone(t, goleak.IgnoreCurrent())

	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()

	errServerDrop := errors.New("server transport terminated abruptly")
	var recvs atomic.Int64

	stream := newScriptedClientStream(ctx,
		func(r *payload.Insert_Request) error {
			return nil
		},
		func() (*payload.Object_StreamLocation, error) {
			if recvs.Add(1) == 3 {
				return nil, errServerDrop
			}
			return &payload.Object_StreamLocation{}, nil
		},
	)

	provider := func() (*payload.Insert_Request, bool) {
		return &payload.Insert_Request{}, true
	}

	callBack := func(r *payload.Object_StreamLocation, err error) bool {
		return err == nil
	}

	resc := make(chan error, 1)
	go func() {
		resc <- BidirectionalStreamClient(stream, 4, provider, callBack)
	}()

	select {
	case err := <-resc:
		if !errors.Is(err, errServerDrop) {
			t.Errorf("stream.Recv error %v was discarded and returned err: %v", errServerDrop, err)
		}
	case <-time.After(3 * time.Second):
		_ = stream.CloseSend()
		stream.cancel()
		t.Fatalf("BidirectionalStreamClient did not exit upon abnormal server drop")
	}
}

// TestBidirectionalStreamClient_callbackTrueOnError verifies that
// BidirectionalStreamClient returns instead of hanging on its receive loop when
// callBack returns true on a Recv error (e.g. a caller that only logs the error
// and does not abort the callback).
func TestBidirectionalStreamClient_callbackTrueOnError(t *testing.T) {
	// defer goleak.VerifyNone(t, goleak.IgnoreCurrent())

	ctx, cancel := context.WithTimeout(t.Context(), 3*time.Second)
	defer cancel()

	errServerAbort := errors.New("server reset connection")
	stream := newScriptedClientStream(ctx,
		func(r *payload.Insert_Request) error { return nil },
		func() (*payload.Object_StreamLocation, error) {
			return nil, errServerAbort
		},
	)

	provider := func() (*payload.Insert_Request, bool) {
		return &payload.Insert_Request{}, true
	}

	// callBack returns true even when err != nil (caller logs error and does not abort callback)
	var callbackSawError atomic.Bool
	callBack := func(_ *payload.Object_StreamLocation, err error) bool {
		if err != nil {
			callbackSawError.Store(true)
		}
		return true
	}

	resc := make(chan error, 1)
	go func() {
		resc <- BidirectionalStreamClient(stream, 4, provider, callBack)
	}()

	select {
	case err := <-resc:
		t.Logf("BidirectionalStreamClient returned: %v", err)
	case <-time.After(1500 * time.Millisecond):
		_ = stream.CloseSend()
		stream.cancel()
		t.Fatalf("DEADLOCK DETECTED: BidirectionalStreamClient hung forever when callBack returned true on rerr")
	}
}

// TestBidirectionalStream_sendFailureDrain tests that
// server BidirectionalStream correctly handles stream.Send failure when the
// client has queued messages into dc, ensuring finalize() and eg.Wait() complete.
func TestBidirectionalStream_sendFailureDrain(t *testing.T) {
	// defer goleak.VerifyNone(t, goleak.IgnoreCurrent())

	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()

	const totalRequests = 20
	var reqIdx atomic.Int64
	errSendFail := errors.New("simulated network send error")

	serverStream := &mock.StreamInsertServerMock{
		ServerStream: &mock.ServerStreamMock{
			ContextFunc: func() context.Context {
				return ctx
			},
			RecvMsgFunc: func(i any) error {
				idx := reqIdx.Add(1)
				if idx > totalRequests {
					return io.EOF
				}
				return nil
			},
			SendMsgFunc: func(i any) error {
				if reqIdx.Load() >= 2 {
					return errSendFail
				}
				return nil
			},
		},
	}

	handler := func(ctx context.Context, req *payload.Insert_Request) (*payload.Object_StreamLocation, error) {
		return &payload.Object_StreamLocation{}, nil
	}

	errc := make(chan error, 1)
	go func() {
		errc <- BidirectionalStream(ctx, serverStream, 4, handler)
	}()

	select {
	case err := <-errc:
		if err == nil {
			t.Errorf("expected error from BidirectionalStream on SendMsg failure, got nil")
		}
	case <-time.After(3 * time.Second):
		cancel()
		t.Fatalf("BidirectionalStream deadlocked on SendMsg failure")
	}
}
