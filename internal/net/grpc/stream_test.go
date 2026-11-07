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
