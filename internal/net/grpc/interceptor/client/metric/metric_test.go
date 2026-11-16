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
package metric

import (
	"context"
	"maps"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/io"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"google.golang.org/grpc"
)

// fakeClientStream returns the queued RecvMsg errors in order (io.EOF once
// exhausted) and sendErr from every SendMsg.
type fakeClientStream struct {
	grpc.ClientStream
	sendErr error
	recv    []error
}

func (f *fakeClientStream) RecvMsg(any) error {
	if len(f.recv) == 0 {
		return io.EOF
	}
	err := f.recv[0]
	f.recv = f.recv[1:]
	return err
}

func (f *fakeClientStream) SendMsg(any) error { return f.sendErr }

func (*fakeClientStream) CloseSend() error { return nil }

func TestClientInterceptors_StreamRecordsAtCompletion(t *testing.T) {
	t.Parallel()
	const method = "/vald.v1.Search/StreamSearch"
	errInternal := status.Error(codes.Internal, "internal")
	errUnavailable := status.Error(codes.Unavailable, "unavailable")
	errExhausted := status.Error(codes.ResourceExhausted, "exhausted")

	type collectFunc func() map[string]int64
	tests := []struct {
		streamerErr    error
		fake           *fakeClientStream
		wantAfterSetup map[string]int64
		run            func(t *testing.T, cs grpc.ClientStream, cancel context.CancelFunc, collect collectFunc)
		want           map[string]int64
		name           string
		serverStreams  bool
	}{
		{
			name:           "setup is not recorded and the final status error is",
			serverStreams:  true,
			fake:           &fakeClientStream{recv: []error{nil, errInternal}},
			wantAfterSetup: map[string]int64{},
			run: func(t *testing.T, cs grpc.ClientStream, _ context.CancelFunc, collect collectFunc) {
				t.Helper()
				if err := cs.RecvMsg(nil); err != nil {
					t.Fatalf("RecvMsg() error = %v", err)
				}
				if got := collect(); len(got) != 0 {
					t.Errorf("recorded before stream end: %v", got)
				}
				if err := cs.RecvMsg(nil); !errors.Is(err, errInternal) {
					t.Errorf("RecvMsg() error = %v, want %v", err, errInternal)
				}
			},
			want: map[string]int64{codes.Internal.String(): 1},
		},
		{
			name:           "io.EOF records OK exactly once",
			serverStreams:  true,
			fake:           &fakeClientStream{},
			wantAfterSetup: map[string]int64{},
			run: func(t *testing.T, cs grpc.ClientStream, _ context.CancelFunc, _ collectFunc) {
				t.Helper()
				for range 2 {
					if err := cs.RecvMsg(nil); !errors.Is(err, io.EOF) {
						t.Errorf("RecvMsg() error = %v, want io.EOF", err)
					}
				}
			},
			want: map[string]int64{codes.OK.String(): 1},
		},
		{
			name:           "non-server-streaming RPC records OK on its single response",
			serverStreams:  false,
			fake:           &fakeClientStream{recv: []error{nil}},
			wantAfterSetup: map[string]int64{},
			run: func(t *testing.T, cs grpc.ClientStream, _ context.CancelFunc, _ collectFunc) {
				t.Helper()
				if err := cs.RecvMsg(nil); err != nil {
					t.Errorf("RecvMsg() error = %v", err)
				}
			},
			want: map[string]int64{codes.OK.String(): 1},
		},
		{
			name:           "CloseSend is not recorded",
			serverStreams:  true,
			fake:           &fakeClientStream{recv: []error{errUnavailable}},
			wantAfterSetup: map[string]int64{},
			run: func(t *testing.T, cs grpc.ClientStream, _ context.CancelFunc, collect collectFunc) {
				t.Helper()
				if err := cs.CloseSend(); err != nil {
					t.Fatalf("CloseSend() error = %v", err)
				}
				if got := collect(); len(got) != 0 {
					t.Errorf("recorded on CloseSend: %v", got)
				}
				if err := cs.RecvMsg(nil); !errors.Is(err, errUnavailable) {
					t.Errorf("RecvMsg() error = %v, want %v", err, errUnavailable)
				}
			},
			want: map[string]int64{codes.Unavailable.String(): 1},
		},
		{
			name:           "SendMsg non-EOF error ends the RPC and is recorded once",
			serverStreams:  true,
			fake:           &fakeClientStream{sendErr: errExhausted},
			wantAfterSetup: map[string]int64{},
			run: func(t *testing.T, cs grpc.ClientStream, _ context.CancelFunc, _ collectFunc) {
				t.Helper()
				if err := cs.SendMsg(nil); !errors.Is(err, errExhausted) {
					t.Errorf("SendMsg() error = %v, want %v", err, errExhausted)
				}
				if err := cs.RecvMsg(nil); !errors.Is(err, io.EOF) {
					t.Errorf("RecvMsg() error = %v, want io.EOF", err)
				}
			},
			want: map[string]int64{codes.ResourceExhausted.String(): 1},
		},
		{
			name:           "SendMsg io.EOF defers to the RecvMsg status",
			serverStreams:  true,
			fake:           &fakeClientStream{sendErr: io.EOF, recv: []error{errUnavailable}},
			wantAfterSetup: map[string]int64{},
			run: func(t *testing.T, cs grpc.ClientStream, _ context.CancelFunc, _ collectFunc) {
				t.Helper()
				if err := cs.SendMsg(nil); !errors.Is(err, io.EOF) {
					t.Errorf("SendMsg() error = %v, want io.EOF", err)
				}
				if err := cs.RecvMsg(nil); !errors.Is(err, errUnavailable) {
					t.Errorf("RecvMsg() error = %v, want %v", err, errUnavailable)
				}
			},
			want: map[string]int64{codes.Unavailable.String(): 1},
		},
		{
			name:           "streamer error is recorded immediately",
			serverStreams:  true,
			streamerErr:    errUnavailable,
			wantAfterSetup: map[string]int64{codes.Unavailable.String(): 1},
			want:           map[string]int64{codes.Unavailable.String(): 1},
		},
		{
			name:           "stream abandoned by ctx cancel is recorded as Canceled",
			serverStreams:  true,
			fake:           &fakeClientStream{},
			wantAfterSetup: map[string]int64{},
			run: func(t *testing.T, _ grpc.ClientStream, cancel context.CancelFunc, collect collectFunc) {
				t.Helper()
				cancel()
				// context.AfterFunc records from its own goroutine: poll, bounded.
				deadline := time.Now().Add(5 * time.Second)
				// record() adds the latency sample before the counter, so once
				// the counter is visible both are; the consistency check in
				// collectCounts runs only after this poll.
				for len(collect()) == 0 && time.Now().Before(deadline) {
					time.Sleep(10 * time.Millisecond)
				}
			},
			want: map[string]int64{codes.Canceled.String(): 1},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			reader := sdkmetric.NewManualReader()
			mp := sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader))
			t.Cleanup(func() {
				if err := mp.Shutdown(context.WithoutCancel(t.Context())); err != nil {
					t.Errorf("Shutdown() error = %v", err)
				}
			})
			collect := func() map[string]int64 {
				counts, _, _ := rawCounts(t, reader)
				return counts
			}

			_, streamInterceptor, err := clientInterceptors(mp.Meter("test"))
			if err != nil {
				t.Fatalf("clientInterceptors() error = %v", err)
			}
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()
			streamer := func(context.Context, *grpc.StreamDesc, *grpc.ClientConn, string, ...grpc.CallOption) (grpc.ClientStream, error) {
				if tc.streamerErr != nil {
					return nil, tc.streamerErr
				}
				return tc.fake, nil
			}
			cs, err := streamInterceptor(ctx, &grpc.StreamDesc{ServerStreams: tc.serverStreams}, nil, method, streamer)
			if !errors.Is(err, tc.streamerErr) {
				t.Fatalf("stream interceptor error = %v, want %v", err, tc.streamerErr)
			}
			if got := collect(); !maps.Equal(got, tc.wantAfterSetup) {
				t.Errorf("after setup: recorded %v, want %v", got, tc.wantAfterSetup)
			}
			if tc.run != nil {
				tc.run(t, cs, cancel, collect)
			}
			if got := collectCounts(t, reader); !maps.Equal(got, tc.want) {
				t.Errorf("recorded %v, want %v", got, tc.want)
			}
		})
	}
}

// collectCounts returns client_completed_rpcs per grpc_client_status and
// checks that every completed RPC has exactly one latency sample.
func collectCounts(t *testing.T, reader *sdkmetric.ManualReader) map[string]int64 {
	t.Helper()
	counts, total, latencies := rawCounts(t, reader)
	if uint64(total) != latencies { //nolint:gosec // total sums a monotonic counter that only ever adds 1, so it is never negative
		t.Errorf("completed RPCs = %d but latency samples = %d", total, latencies)
	}
	return counts
}

// rawCounts collects without asserting that the counter and the histogram
// agree: record() updates them in two steps, so a Collect that overlaps a
// record running on another goroutine may observe only one of them.
func rawCounts(
	t *testing.T, reader *sdkmetric.ManualReader,
) (counts map[string]int64, total int64, latencies uint64) {
	t.Helper()
	var rm metricdata.ResourceMetrics
	if err := reader.Collect(context.WithoutCancel(t.Context()), &rm); err != nil {
		t.Fatalf("Collect() error = %v", err)
	}
	counts = make(map[string]int64)
	for _, sm := range rm.ScopeMetrics {
		for _, m := range sm.Metrics {
			switch d := m.Data.(type) {
			case metricdata.Sum[int64]:
				if m.Name != completedRPCsMetricsName {
					continue
				}
				for _, dp := range d.DataPoints {
					v, _ := dp.Attributes.Value(gRPCStatus)
					counts[v.AsString()] += dp.Value
					total += dp.Value
				}
			case metricdata.Histogram[float64]:
				if m.Name != latencyMetricsName {
					continue
				}
				for _, dp := range d.DataPoints {
					latencies += dp.Count
				}
			}
		}
	}
	return counts, total, latencies
}

// NOT IMPLEMENTED BELOW
//
// func TestClientInterceptors(t *testing.T) {
// 	type want struct {
// 		want  grpc.UnaryClientInterceptor
// 		want1 grpc.StreamClientInterceptor
// 		err   error
// 	}
// 	type test struct {
// 		name       string
// 		want       want
// 		checkFunc  func(want, grpc.UnaryClientInterceptor, grpc.StreamClientInterceptor, error) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got grpc.UnaryClientInterceptor, got1 grpc.StreamClientInterceptor, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		if !reflect.DeepEqual(got1, w.want1) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got1, w.want1)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
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
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
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
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			got, got1, err := ClientInterceptors()
// 			if err := checkFunc(test.want, got, got1, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_attributesFromError(t *testing.T) {
// 	type args struct {
// 		method string
// 		err    error
// 	}
// 	type want struct {
// 		want []attribute.KeyValue
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, []attribute.KeyValue) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got []attribute.KeyValue) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           method:"",
// 		           err:nil,
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
// 		           method:"",
// 		           err:nil,
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
//
// 			got := attributesFromError(test.args.method, test.args.err)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
