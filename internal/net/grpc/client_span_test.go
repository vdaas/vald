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

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/observability/attribute"
	"github.com/vdaas/vald/internal/observability/trace"
	ocodes "go.opentelemetry.io/otel/codes"
	otrace "go.opentelemetry.io/otel/trace"
)

// recordingSpan records the calls recordSpanError makes. The embedded nil Span
// panics on any other method, so unexpected calls fail the test.
type recordingSpan struct {
	trace.Span
	desc   string
	errs   []error
	attrs  []attribute.KeyValue
	status int
	code   ocodes.Code
}

func (s *recordingSpan) RecordError(err error, _ ...otrace.EventOption) {
	s.errs = append(s.errs, err)
}

func (s *recordingSpan) SetAttributes(kv ...attribute.KeyValue) {
	s.attrs = append(s.attrs, kv...)
}

func (s *recordingSpan) SetStatus(code ocodes.Code, desc string) {
	s.status++
	s.code, s.desc = code, desc
}

func TestRecordSpanError(t *testing.T) {
	t.Parallel()
	grpcErr := status.Error(codes.Unavailable, "unavailable")
	plainErr := errors.New("plain")
	tests := []struct {
		err       error
		name      string
		wantAttrs []attribute.KeyValue
		wantCalls bool
	}{
		{
			name: "nil error records nothing",
		},
		{
			name:      "gRPC status error records error, status attributes and status",
			err:       grpcErr,
			wantCalls: true,
			wantAttrs: trace.FromGRPCStatus(codes.Unavailable, grpcErr.Error()),
		},
		{
			name:      "non-status error records error and status without attributes",
			err:       plainErr,
			wantCalls: true,
		},
	}
	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			span := new(recordingSpan)
			recordSpanError(span, test.err)
			if !test.wantCalls {
				if len(span.errs) != 0 || len(span.attrs) != 0 || span.status != 0 {
					tt.Errorf("span was modified: %+v", span)
				}
				return
			}
			if len(span.errs) != 1 || !errors.Is(span.errs[0], test.err) {
				tt.Errorf("RecordError calls = %v, want [%v]", span.errs, test.err)
			}
			if len(span.attrs) != len(test.wantAttrs) {
				tt.Fatalf("attributes = %v, want %v", span.attrs, test.wantAttrs)
			}
			for i := range span.attrs {
				if span.attrs[i] != test.wantAttrs[i] {
					tt.Errorf("attribute[%d] = %v, want %v", i, span.attrs[i], test.wantAttrs[i])
				}
			}
			if span.status != 1 || span.code != trace.StatusError || span.desc != test.err.Error() {
				tt.Errorf("SetStatus calls = %d (%v, %q), want 1 (%v, %q)",
					span.status, span.code, span.desc, trace.StatusError, test.err.Error())
			}
		})
	}
	// a nil span must not panic.
	recordSpanError(nil, plainErr)
}

// TestStartAddrSpan_TracingDisabled pins that, with tracing off, startAddrSpan
// returns the given ctx and a nil span without building the span name.
// Not parallel: AllocsPerRun changes GOMAXPROCS.
func TestStartAddrSpan_TracingDisabled(t *testing.T) {
	if trace.Enabled() {
		t.Skip("tracing is enabled in this process")
	}
	type ctxKey struct{}
	ctx := context.WithValue(t.Context(), ctxKey{}, "v")
	addr := string([]byte("10.0.0.1:8081")) // not a constant, like a real target
	gotCtx, span := startAddrSpan(ctx, apiName+"/Client.do/", addr)
	if gotCtx != ctx {
		t.Error("startAddrSpan returned a different ctx")
	}
	if span != nil {
		t.Errorf("startAddrSpan returned span %v, want nil", span)
	}
	if allocs := testing.AllocsPerRun(100, func() {
		_, span = startAddrSpan(ctx, apiName+"/Client.do/", addr)
	}); allocs != 0 {
		t.Errorf("startAddrSpan allocated %v times per call, want 0", allocs)
	}
}
