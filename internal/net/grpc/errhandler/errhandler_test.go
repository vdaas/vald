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

package errhandler

import (
	"context"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/test/goleak"
	ocodes "go.opentelemetry.io/otel/codes"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

// recordingSpan starts a real recording span backed by an in-memory recorder so
// tests can assert exactly what RecordSpanError wrote to it.
func recordingSpan(t *testing.T) (trace.Span, *tracetest.SpanRecorder) {
	t.Helper()
	sr := tracetest.NewSpanRecorder()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(sr))
	_, span := tp.Tracer("errhandler-test").Start(context.Background(), "op")
	return span, sr
}

func TestRecordSpanError(t *testing.T) {
	type args struct {
		err  error
		code codes.Code
	}
	type want struct {
		wantStatusError bool
		wantEvent       bool
	}
	type test struct {
		name string
		args args
		want want
		// nilSpan drives the span==nil no-op path instead of a recording span.
		nilSpan bool
	}
	tests := []test{
		{
			name:    "nil span is a no-op and does not panic",
			args:    args{code: codes.Internal, err: errors.New("boom")},
			nilSpan: true,
			want:    want{wantStatusError: false, wantEvent: false},
		},
		{
			name: "records the error and sets error status for Internal",
			args: args{code: codes.Internal, err: errors.New("internal boom")},
			want: want{wantStatusError: true, wantEvent: true},
		},
		{
			name: "records the error and sets error status for Canceled",
			args: args{code: codes.Canceled, err: errors.New("canceled")},
			want: want{wantStatusError: true, wantEvent: true},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if tc.nilSpan {
				// must not panic
				RecordSpanError(nil, tc.args.code, tc.args.err)
				return
			}
			span, sr := recordingSpan(tt)
			RecordSpanError(span, tc.args.code, tc.args.err)
			span.End()
			ended := sr.Ended()
			if len(ended) != 1 {
				tt.Fatalf("expected exactly 1 ended span, got %d", len(ended))
			}
			ro := ended[0]
			if got := ro.Status().Code == ocodes.Error; got != tc.want.wantStatusError {
				tt.Errorf("span status error = %v, want %v (status=%v)", got, tc.want.wantStatusError, ro.Status())
			}
			hasExceptionEvent := false
			for _, ev := range ro.Events() {
				if ev.Name == "exception" {
					hasExceptionEvent = true
					break
				}
			}
			if hasExceptionEvent != tc.want.wantEvent {
				tt.Errorf("recorded exception event = %v, want %v", hasExceptionEvent, tc.want.wantEvent)
			}
			// The attributes attached must be exactly what trace.FromGRPCStatus
			// produces for the code, i.e. what the hand-written call sites wrote.
			wantAttrs := trace.FromGRPCStatus(tc.args.code, tc.args.err.Error())
			got := ro.Attributes()
			for _, wa := range wantAttrs {
				found := false
				for _, ga := range got {
					if ga.Key == wa.Key && ga.Value == wa.Value {
						found = true
						break
					}
				}
				if !found {
					tt.Errorf("missing expected span attribute %v=%v", wa.Key, wa.Value)
				}
			}
		})
	}
}

func TestRecordSpanStatus(t *testing.T) {
	type test struct {
		st      *status.Status
		err     error
		name    string
		nilSpan bool
		want    bool // whether the span should end up errored
	}
	tests := []test{
		{name: "nil status is a no-op", st: nil, err: errors.New("x"), want: false},
		{name: "nil span is a no-op", st: status.New(codes.NotFound, "missing"), err: errors.New("x"), nilSpan: true, want: false},
		{name: "records the status attributes for NotFound", st: status.New(codes.NotFound, "missing"), err: errors.New("boom"), want: true},
		{name: "records the status attributes for Internal", st: status.New(codes.Internal, "kaboom"), err: errors.New("boom"), want: true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if tc.nilSpan {
				RecordSpanStatus(nil, tc.st, tc.err)
				return
			}
			span, sr := recordingSpan(tt)
			RecordSpanStatus(span, tc.st, tc.err)
			span.End()
			ended := sr.Ended()
			if len(ended) != 1 {
				tt.Fatalf("expected exactly 1 ended span, got %d", len(ended))
			}
			ro := ended[0]
			if got := ro.Status().Code == ocodes.Error; got != tc.want {
				tt.Errorf("span errored = %v, want %v", got, tc.want)
			}
			if tc.want {
				wantAttrs := trace.FromGRPCStatus(tc.st.Code(), tc.st.Message())
				got := ro.Attributes()
				for _, wa := range wantAttrs {
					found := false
					for _, ga := range got {
						if ga.Key == wa.Key && ga.Value == wa.Value {
							found = true
							break
						}
					}
					if !found {
						tt.Errorf("missing expected span attribute %v=%v", wa.Key, wa.Value)
					}
				}
			}
		})
	}
}

func TestRecordSpanAttrs(t *testing.T) {
	t.Run("nil span is a no-op", func(tt *testing.T) {
		defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
		RecordSpanAttrs(nil, trace.StatusCodeInternal("x"), errors.New("x"))
	})
	t.Run("records the caller-prepared attributes and errors the span", func(tt *testing.T) {
		defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
		err := errors.New("boom")
		attrs := trace.FromGRPCStatus(codes.ResourceExhausted, "quota")
		span, sr := recordingSpan(tt)
		RecordSpanAttrs(span, attrs, err)
		span.End()
		ended := sr.Ended()
		if len(ended) != 1 {
			tt.Fatalf("expected exactly 1 ended span, got %d", len(ended))
		}
		ro := ended[0]
		if ro.Status().Code != ocodes.Error {
			tt.Errorf("span not errored: %v", ro.Status())
		}
		got := ro.Attributes()
		for _, wa := range attrs {
			found := false
			for _, ga := range got {
				if ga.Key == wa.Key && ga.Value == wa.Value {
					found = true
					break
				}
			}
			if !found {
				tt.Errorf("missing expected span attribute %v=%v", wa.Key, wa.Value)
			}
		}
	})
}

func TestHandleError(t *testing.T) {
	defer goleak.VerifyNone(t, goleak.IgnoreCurrent())
	testErr := errors.New("handler error")

	// The generic parameter T governs only the typed nil returned on the error
	// path; each instantiation must return a nil *T and the same error, and must
	// have recorded the error on a non-nil span.
	t.Run("returns typed nil and the same error, recording on the span", func(tt *testing.T) {
		span, sr := recordingSpan(tt)
		type payloadA struct{ X int }
		got, err := HandleError[payloadA](span, codes.Internal, testErr)
		span.End()
		if got != nil {
			tt.Errorf("expected nil *payloadA, got %#v", got)
		}
		if !errors.Is(err, testErr) {
			tt.Errorf("expected the input error back, got %v", err)
		}
		if ended := sr.Ended(); len(ended) != 1 || ended[0].Status().Code != ocodes.Error {
			tt.Errorf("expected the span to be marked errored, ended=%d", len(ended))
		}
	})

	t.Run("nil span still returns typed nil and the error", func(tt *testing.T) {
		type payloadB struct{ Y string }
		got, err := HandleError[payloadB](nil, codes.NotFound, testErr)
		if got != nil {
			tt.Errorf("expected nil *payloadB, got %#v", got)
		}
		if !errors.Is(err, testErr) {
			tt.Errorf("expected the input error back, got %v", err)
		}
	})
}
