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

package recover

import (
	"context"
	"net"
	"runtime"
	"testing"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/log/logger"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/internal/test/goleak"
	ggrpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

func TestMain(m *testing.M) {
	log.Init(log.WithLoggerType(logger.NOP.String()))
	goleak.VerifyTestMain(m)
}

func TestHandlePanic(t *testing.T) {
	t.Parallel()
	var nilStr *string
	validStr := "valid string"

	tests := []struct {
		name     string
		panicVal any
		wantSub  string
	}{
		{
			name:     "typed nil string pointer",
			panicVal: nilStr,
			wantSub:  "panic recovered",
		},
		{
			name:     "valid string pointer",
			panicVal: &validStr,
			wantSub:  "valid string",
		},
		{
			name:     "string value",
			panicVal: "direct string",
			wantSub:  "direct string",
		},
		{
			name:     "error value",
			panicVal: errors.New("sentinel err"),
			wantSub:  "sentinel err",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := handlePanic(tt.panicVal, nil)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !strings.Contains(err.Error(), tt.wantSub) {
				t.Errorf("expected %q in error %q", tt.wantSub, err.Error())
			}
		})
	}
}

type customPanicStruct struct {
	Msg  string
	Code int
}

// TestInterceptor_arbitraryPanicTypes tests that Interceptor catches all
// arbitrary panic types without crashing, wrapping them in Vald error types.
func TestInterceptor_arbitraryPanicTypes(t *testing.T) {
	t.Parallel()

	dummyInfo := &grpc.UnaryServerInfo{
		FullMethod: "/vald.v1.Vald/Insert",
	}

	strPointerVal := "pointer string panic"
	targetErr := errors.New("custom sentinel error")

	tests := []struct {
		name      string
		panicVal  any
		checkFunc func(t *testing.T, err error)
	}{
		{
			name:     "panic with nil",
			panicVal: nil,
			checkFunc: func(t *testing.T, err error) {
				t.Helper()
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				// Under Go 1.21+, panic(nil) recovers as *runtime.PanicNilError
				// or gets wrapped as panic recovered
				errStr := err.Error()
				if !strings.Contains(errStr, "nil") && !strings.Contains(errStr, "panic") {
					t.Errorf("expected nil or panic mention in error, got: %v", err)
				}
			},
		},
		{
			name:     "panic with string",
			panicVal: "critical database failure",
			checkFunc: func(t *testing.T, err error) {
				t.Helper()
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if !strings.Contains(err.Error(), "critical database failure") {
					t.Errorf("expected string panic message in error, got: %v", err)
				}
			},
		},
		{
			name:     "panic with string pointer",
			panicVal: &strPointerVal,
			checkFunc: func(t *testing.T, err error) {
				t.Helper()
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if !strings.Contains(err.Error(), "pointer string panic") {
					t.Errorf("expected dereferenced string in error, got: %v", err)
				}
			},
		},
		{
			name:     "panic with error",
			panicVal: targetErr,
			checkFunc: func(t *testing.T, err error) {
				t.Helper()
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if !errors.Is(err, targetErr) {
					t.Errorf("expected error to wrap targetErr, got: %v", err)
				}
			},
		},
		{
			name:     "panic with custom struct",
			panicVal: customPanicStruct{Msg: "corrupt index", Code: 503},
			checkFunc: func(t *testing.T, err error) {
				t.Helper()
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				errStr := err.Error()
				if !strings.Contains(errStr, "corrupt index") || !strings.Contains(errStr, "503") {
					t.Errorf("expected struct contents in error message, got: %v", err)
				}
			},
		},
		{
			name: "panic with runtime.Error (nil pointer dereference)",
			panicVal: func() runtime.Error {
				var r runtime.Error
				defer func() {
					if rec := recover(); rec != nil {
						r = rec.(runtime.Error)
					}
				}()
				var p *int
				_ = *p
				return r
			}(),
			checkFunc: func(t *testing.T, err error) {
				t.Helper()
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if !strings.Contains(err.Error(), "runtime error") {
					t.Errorf("expected runtime error mention, got: %v", err)
				}
			},
		},
		{
			name: "panic with runtime.Error (slice index out of range)",
			panicVal: func() runtime.Error {
				var r runtime.Error
				defer func() {
					if rec := recover(); rec != nil {
						r = rec.(runtime.Error)
					}
				}()
				var s []int
				_ = s[10]
				return r
			}(),
			checkFunc: func(t *testing.T, err error) {
				t.Helper()
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if !strings.Contains(err.Error(), "runtime error") {
					t.Errorf("expected runtime error mention, got: %v", err)
				}
			},
		},
		{
			name:     "panic with byte slice",
			panicVal: []byte("arbitrary bytes"),
			checkFunc: func(t *testing.T, err error) {
				t.Helper()
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if !strings.Contains(err.Error(), "panic recovered") {
					t.Errorf("expected panic recovered in error, got: %v", err)
				}
			},
		},
		{
			name:     "panic with integer",
			panicVal: 404,
			checkFunc: func(t *testing.T, err error) {
				t.Helper()
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if !strings.Contains(err.Error(), "404") {
					t.Errorf("expected 404 in error, got: %v", err)
				}
			},
		},
		{
			name:     "panic with map",
			panicVal: map[string]int{"fail": 1},
			checkFunc: func(t *testing.T, err error) {
				t.Helper()
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if !strings.Contains(err.Error(), "fail") {
					t.Errorf("expected map key in error, got: %v", err)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			interceptor := Interceptor()

			handler := func(ctx context.Context, req any) (any, error) {
				panic(tt.panicVal)
			}

			// Interceptor must not crash or re-panic
			resp, err := interceptor(context.Background(), "test-req", dummyInfo, handler)
			if resp != nil {
				t.Errorf("expected nil resp, got: %v", resp)
			}
			tt.checkFunc(t, err)
		})
	}
}

// TestStreamInterceptor_arbitraryPanicTypes tests StreamInterceptor with arbitrary panic types.
func TestStreamInterceptor_arbitraryPanicTypes(t *testing.T) {
	t.Parallel()

	dummyInfo := &grpc.StreamServerInfo{
		FullMethod: "/vald.v1.Vald/StreamInsert",
	}

	tests := []struct {
		name     string
		panicVal any
	}{
		{name: "nil panic", panicVal: nil},
		{name: "string panic", panicVal: "stream aborted"},
		{name: "error panic", panicVal: errors.New("stream err")},
		{name: "struct panic", panicVal: customPanicStruct{Msg: "stream bad", Code: 500}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			interceptor := StreamInterceptor()

			handler := func(srv any, ss grpc.ServerStream) error {
				panic(tt.panicVal)
			}

			err := interceptor("dummy-srv", nil, dummyInfo, handler)
			if err == nil {
				t.Fatal("expected error from StreamInterceptor, got nil")
			}
		})
	}
}

// testPanicValdServer implements ValdServer to simulate arbitrary panics inside real gRPC RPC calls.
type testPanicValdServer struct {
	vald.UnimplementedValdServer
}

func (s *testPanicValdServer) Insert(
	ctx context.Context, req *payload.Insert_Request,
) (*payload.Object_Location, error) {
	if req == nil || req.Vector == nil {
		return &payload.Object_Location{Uuid: "normal-ok"}, nil
	}

	switch req.Vector.Id {
	case "panic-nil":
		panic(nil)
	case "panic-string":
		panic("string panic inside gRPC handler")
	case "panic-error":
		panic(errors.New("custom error inside gRPC handler"))
	case "panic-struct":
		panic(customPanicStruct{Msg: "struct panic", Code: 502})
	case "panic-nil-deref":
		var p *int
		_ = *p
	default:
		// Normal execution
		return &payload.Object_Location{Uuid: req.Vector.Id}, nil
	}
	return nil, nil
}

// TestInterceptor_grpcStatusError executes end-to-end gRPC RPC calls over bufconn
// verifying that panics are converted to gRPC status errors and client receives non-OK status codes,
// while the server process remains healthy for subsequent calls.
func TestInterceptor_grpcStatusError(t *testing.T) {
	t.Parallel()

	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() {
		lis.Close()
	})

	srv := ggrpc.NewServer(
		ggrpc.UnaryInterceptor(Interceptor()),
	)
	vald.RegisterValdServer(srv, &testPanicValdServer{})

	go func() {
		_ = srv.Serve(lis)
	}()
	t.Cleanup(func() {
		srv.GracefulStop()
	})

	ctx := t.Context()
	conn, err := ggrpc.NewClient("passthrough://bufnet",
		ggrpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}),
		ggrpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("failed to dial bufnet: %v", err)
	}
	t.Cleanup(func() {
		conn.Close()
	})

	client := vald.NewValdClient(conn)

	panicScenarios := []struct {
		name        string
		id          string
		expectedSub string
	}{
		{
			name:        "panic nil",
			id:          "panic-nil",
			expectedSub: "panic",
		},
		{
			name:        "panic string",
			id:          "panic-string",
			expectedSub: "string panic inside gRPC handler",
		},
		{
			name:        "panic error",
			id:          "panic-error",
			expectedSub: "custom error inside gRPC handler",
		},
		{
			name:        "panic struct",
			id:          "panic-struct",
			expectedSub: "struct panic",
		},
		{
			name:        "panic nil dereference",
			id:          "panic-nil-deref",
			expectedSub: "runtime error",
		},
	}

	for _, sc := range panicScenarios {
		t.Run(sc.name, func(t *testing.T) {
			req := &payload.Insert_Request{
				Vector: &payload.Object_Vector{
					Id: sc.id,
				},
			}

			resp, err := client.Insert(ctx, req)
			if err == nil {
				t.Fatalf("expected error from server panic, got resp: %v", resp)
			}

			st, ok := status.FromError(err)
			if !ok {
				t.Fatalf("expected gRPC status error, got: %T (%v)", err, err)
			}

			if st.Code() == codes.OK {
				t.Errorf("expected non-OK status code, got: %v", st.Code())
			}

			if !strings.Contains(st.Message(), sc.expectedSub) {
				t.Errorf("expected status message to contain %q, got: %q", sc.expectedSub, st.Message())
			}
		})
	}

	// Verify server is still alive and processes normal request successfully
	t.Run("post-panic normal request succeeds", func(t *testing.T) {
		req := &payload.Insert_Request{
			Vector: &payload.Object_Vector{
				Id: "healthy-subsequent-request",
			},
		}

		resp, err := client.Insert(ctx, req)
		if err != nil {
			t.Fatalf("expected normal request to succeed after recovered panics, got err: %v", err)
		}
		if resp == nil || resp.Uuid != "healthy-subsequent-request" {
			t.Fatalf("expected normal response, got: %v", resp)
		}
	})
}

// NOT IMPLEMENTED BELOW
//
// func TestInterceptor(t *testing.T) {
// 	type want struct {
// 		want grpc.UnaryServerInterceptor
// 	}
// 	type test struct {
// 		name       string
// 		want       want
// 		checkFunc  func(want, grpc.UnaryServerInterceptor) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got grpc.UnaryServerInterceptor) error {
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
// 			got := Interceptor()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestStreamInterceptor(t *testing.T) {
// 	type want struct {
// 		want grpc.StreamServerInterceptor
// 	}
// 	type test struct {
// 		name       string
// 		want       want
// 		checkFunc  func(want, grpc.StreamServerInterceptor) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got grpc.StreamServerInterceptor) error {
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
// 			got := RecoverStreamInterceptor()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
