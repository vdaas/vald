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

package logging

import (
	"context"
	"fmt"
	"testing"

	"github.com/vdaas/vald/internal/encoding/json"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/sync/atomic"
	"github.com/vdaas/vald/internal/test/goleak"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

// countingError counts how many times its message is requested.
// AccessLogEntity.String reads the message of entity.Error, so the count tells
// whether the access log entity has been rendered.
type countingError struct {
	n *atomic.Int64
}

func (e countingError) Error() string {
	e.n.Add(1)
	return "counting error"
}

type stubServerStream struct {
	grpc.ServerStream
	ctx context.Context //nolint:containedctx // test stub that must return a fixed context
}

func (s stubServerStream) Context() context.Context {
	return s.ctx
}

// TestAccessLogInterceptor_LazyRender verifies that the access log entity is not
// rendered while the log level that would emit it is disabled.
// log.Init is effective only once per process, and no other test in this package
// initializes it, so the error level set here applies to the whole package.
func TestAccessLogInterceptor_LazyRender(t *testing.T) {
	log.Init(log.WithLoggerType("glg"), log.WithLevel("error"))

	t.Run("unary failure does not render the entity when warn is disabled", func(tt *testing.T) {
		var n atomic.Int64
		handlerErr := countingError{n: &n}
		_, err := AccessLogInterceptor()(
			tt.Context(),
			nil,
			&grpc.UnaryServerInfo{FullMethod: "/vald.v1.Search/Search"},
			func(context.Context, any) (any, error) { return nil, handlerErr },
		)
		if !errors.Is(err, handlerErr) {
			tt.Fatalf("handler error must be returned as is, got: %v", err)
		}
		if got := n.Load(); got != 0 {
			tt.Errorf("entity was rendered %d times while warn is disabled", got)
		}
	})

	t.Run("stream failure does not render the entity when warn is disabled", func(tt *testing.T) {
		var n atomic.Int64
		handlerErr := countingError{n: &n}
		err := AccessLogStreamInterceptor()(
			nil,
			stubServerStream{ctx: tt.Context()},
			&grpc.StreamServerInfo{FullMethod: "/vald.v1.Search/StreamSearch"},
			func(any, grpc.ServerStream) error { return handlerErr },
		)
		if !errors.Is(err, handlerErr) {
			tt.Fatalf("handler error must be returned as is, got: %v", err)
		}
		if got := n.Load(); got != 0 {
			tt.Errorf("entity was rendered %d times while warn is disabled", got)
		}
	})
}

// TestLazyAccessLog verifies that the lazily rendered value produces the same
// output as the pre-rendered string for the fmt based backends (glg RAW, zap) and
// for the glg JSON format.
func TestLazyAccessLog(t *testing.T) {
	var n atomic.Int64
	tests := []struct {
		entity *AccessLogEntity
		name   string
	}{
		{
			name: "success entity",
			entity: &AccessLogEntity{
				GRPC:      &AccessLogGRPCEntity{Kind: grpcKindUnary, Service: "vald.v1.Search", Method: "Search"},
				TraceID:   "0123456789abcdef",
				StartTime: 1,
				EndTime:   3,
				Latency:   2,
			},
		},
		{
			name: "failure entity with an error message",
			entity: &AccessLogEntity{
				GRPC:  &AccessLogGRPCEntity{Kind: grpcKindStream, Service: "vald.v1.Insert", Method: "StreamInsert"},
				Error: countingError{n: &n},
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			want := test.entity.String()
			lazy := (*lazyAccessLog)(test.entity)

			if got := lazy.String(); got != want {
				tt.Errorf("String() = %q, want %q", got, want)
			}
			// glg RAW format ("%v %v") and zap (fmt.Sprint)
			if got, w := fmt.Sprintf("%v %v", rpcCompletedMessage, lazy), fmt.Sprintf("%v %v", rpcCompletedMessage, want); got != w {
				tt.Errorf("Sprintf = %q, want %q", got, w)
			}
			if got, w := fmt.Sprint(rpcCompletedMessage, lazy), fmt.Sprint(rpcCompletedMessage, want); got != w {
				tt.Errorf("Sprint = %q, want %q", got, w)
			}
			// glg JSON format encodes the []any of arguments
			got, err := json.Marshal([]any{rpcFailedMessage, lazy})
			if err != nil {
				tt.Fatal(err)
			}
			w, err := json.Marshal([]any{rpcFailedMessage, want})
			if err != nil {
				tt.Fatal(err)
			}
			if string(got) != string(w) {
				tt.Errorf("json.Marshal = %s, want %s", got, w)
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
//
// func TestAccessLogEntity_String(t *testing.T) {
// 	type fields struct {
// 		Error     error
// 		GRPC      *AccessLogGRPCEntity
// 		TraceID   string
// 		StartTime int64
// 		EndTime   int64
// 		Latency   int64
// 	}
// 	type want struct {
// 		wantStr string
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, string) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, gotStr string) error {
// 		if !reflect.DeepEqual(gotStr, w.wantStr) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotStr, w.wantStr)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           Error:nil,
// 		           GRPC:AccessLogGRPCEntity{},
// 		           TraceID:"",
// 		           StartTime:0,
// 		           EndTime:0,
// 		           Latency:0,
// 		       },
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
// 		           fields: fields {
// 		           Error:nil,
// 		           GRPC:AccessLogGRPCEntity{},
// 		           TraceID:"",
// 		           StartTime:0,
// 		           EndTime:0,
// 		           Latency:0,
// 		           },
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
// 			e := AccessLogEntity{
// 				Error:     test.fields.Error,
// 				GRPC:      test.fields.GRPC,
// 				TraceID:   test.fields.TraceID,
// 				StartTime: test.fields.StartTime,
// 				EndTime:   test.fields.EndTime,
// 				Latency:   test.fields.Latency,
// 			}
//
// 			gotStr := e.String()
// 			if err := checkFunc(test.want, gotStr); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestAccessLogInterceptor(t *testing.T) {
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
// 			got := AccessLogInterceptor()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestAccessLogStreamInterceptor(t *testing.T) {
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
// 			got := AccessLogStreamInterceptor()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_parseMethod(t *testing.T) {
// 	type args struct {
// 		fullMethod string
// 	}
// 	type want struct {
// 		wantService string
// 		wantMethod  string
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, string, string) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotService string, gotMethod string) error {
// 		if !reflect.DeepEqual(gotService, w.wantService) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotService, w.wantService)
// 		}
// 		if !reflect.DeepEqual(gotMethod, w.wantMethod) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotMethod, w.wantMethod)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           fullMethod:"",
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
// 		           fullMethod:"",
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
// 			gotService, gotMethod := parseMethod(test.args.fullMethod)
// 			if err := checkFunc(test.want, gotService, gotMethod); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
