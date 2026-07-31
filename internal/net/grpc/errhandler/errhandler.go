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

// Package errhandler consolidates the span-recording and typed-error-return
// boilerplate that every Vald gRPC handler repeats on its error paths. Instead
// of copying the "record the error on the span with the attributes for its
// status code, mark the span errored, then return the zero payload and the
// error" tail into every RPC method, handlers call a single helper. It composes
// only the already-shared internal/observability/trace and
// internal/net/grpc/codes vocabularies, so it is reachable from every handler
// package (agent, gateway, discoverer, ...) without introducing a pkg-to-pkg
// dependency (Vald Law 5) and touches no generated code (Law 1).
package errhandler

import (
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/observability/trace"
)

// RecordSpanError records err on span with the tracing attributes that
// correspond to the given gRPC status code, then marks the span as errored. It
// is a no-op when span is nil, folding the ubiquitous
//
//	if span != nil {
//		span.RecordError(err)
//		span.SetAttributes(trace.FromGRPCStatus(code, err.Error())...)
//		span.SetStatus(trace.StatusError, err.Error())
//	}
//
// block into a single call. trace.FromGRPCStatus(code, msg) dispatches to the
// exact same trace.StatusCode<Code>(msg) helper the hand-written call sites
// used, so for a given code the recorded attributes are byte-for-byte identical
// — callers that previously wrote trace.StatusCodeInternal(...) pass
// codes.Internal here and observe no change. Pass the code that reproduces the
// call site's current attribute, not a "corrected" one, to preserve behavior.
func RecordSpanError(span trace.Span, code codes.Code, err error) {
	if span == nil {
		return
	}
	span.RecordError(err)
	span.SetAttributes(trace.FromGRPCStatus(code, err.Error())...)
	span.SetStatus(trace.StatusError, err.Error())
}

// RecordSpanStatus records err on span using the attributes derived from an
// already-parsed gRPC *status.Status (its Code and Message), then marks the span
// as errored. It is a no-op when st or span is nil, folding the stream/broadcast
// tail
//
//	st, _ := status.FromError(err)
//	if st != nil && span != nil {
//		span.RecordError(err)
//		span.SetAttributes(trace.FromGRPCStatus(st.Code(), st.Message())...)
//		span.SetStatus(trace.StatusError, err.Error())
//	}
//
// into a single call. Unlike RecordSpanError, the span attribute message is the
// status message (st.Message()), not err.Error(); these paths record the parsed
// status because the caller also builds a status-bearing stream response from
// the same st (e.g. st.Proto()), so st stays owned by the caller and is passed
// in rather than re-parsed here.
func RecordSpanStatus(span trace.Span, st *status.Status, err error) {
	if st == nil || span == nil {
		return
	}
	span.RecordError(err)
	span.SetAttributes(trace.FromGRPCStatus(st.Code(), st.Message())...)
	span.SetStatus(trace.StatusError, err.Error())
}

// RecordSpanAttrs records err on span with a caller-prepared attribute set, then
// marks the span as errored. It is a no-op when span is nil, folding the tail
//
//	if span != nil {
//		span.RecordError(err)
//		span.SetAttributes(attrs...)
//		span.SetStatus(trace.StatusError, err.Error())
//	}
//
// used by the broadcast paths that classify the error into attrs with their own
// (per-RPC, tolerable-vs-fatal) switch before recording it. Only the span
// recording is consolidated here; the caller keeps its classification and its
// return decision, so this changes no control flow.
func RecordSpanAttrs(span trace.Span, attrs trace.Attributes, err error) {
	if span == nil {
		return
	}
	span.RecordError(err)
	span.SetAttributes(attrs...)
	span.SetStatus(trace.StatusError, err.Error())
}

// HandleError records err on span (via RecordSpanError) and returns the zero
// value of the handler's response payload alongside err, collapsing the
//
//	if span != nil { ... }
//	return nil, err
//
// tail of a unary handler into a single statement:
//
//	return errhandler.HandleError[payload.Object_Location](span, code, err)
//
// The type parameter T is only ever used to produce a typed nil (*T)(nil) on
// the error path; the payload is never inspected or constructed, so this adds
// no cost beyond a nil pointer and cannot touch generated protobuf internals
// (Law 1). err must already be the classified/wrapped status error; HandleError
// deliberately does not classify, because the error-to-code mapping differs per
// core library and stays owned by the caller.
func HandleError[T any](span trace.Span, code codes.Code, err error) (*T, error) {
	RecordSpanError(span, code, err)
	return nil, err
}
