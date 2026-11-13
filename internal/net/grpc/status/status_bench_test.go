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
package status

import (
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/errdetails"
)

// The error helpers run on every failed RPC (NotFound and InvalidArgument are
// common in normal operation), so their per-call cost is on the request path.

func BenchmarkParseError_gRPCStatus(b *testing.B) {
	err := Error(codes.NotFound, "object not found")
	b.ReportAllocs()
	for b.Loop() {
		if _, _, err := ParseError(err, codes.Internal, "failed to parse gRPC error response"); err == nil {
			b.Fatal("expected an error")
		}
	}
}

func BenchmarkParseError_plainError(b *testing.B) {
	err := errors.New("plain error")
	b.ReportAllocs()
	for b.Loop() {
		if _, _, err := ParseError(err, codes.Internal, "failed to parse gRPC error response"); err == nil {
			b.Fatal("expected an error")
		}
	}
}

func BenchmarkParseError_withDetails(b *testing.B) {
	err := WrapWithInvalidArgument("API invalid argument", errors.New("invalid argument"),
		&errdetails.RequestInfo{RequestId: "uuid-0001", ServingData: "serving data"},
		&errdetails.BadRequest{FieldViolations: []*errdetails.BadRequestFieldViolation{{Field: "vector", Description: "invalid"}}},
	)
	b.ReportAllocs()
	for b.Loop() {
		if _, _, err := ParseError(err, codes.Internal, "failed to parse gRPC error response"); err == nil {
			b.Fatal("expected an error")
		}
	}
}

func BenchmarkWrapWithInvalidArgument_withDetails(b *testing.B) {
	base := errors.New("invalid argument")
	b.ReportAllocs()
	for b.Loop() {
		if err := WrapWithInvalidArgument("API invalid argument", base,
			&errdetails.RequestInfo{RequestId: "uuid-0001", ServingData: "serving data"},
			&errdetails.BadRequest{FieldViolations: []*errdetails.BadRequestFieldViolation{{Field: "vector", Description: "invalid"}}},
		); err == nil {
			b.Fatal("expected an error")
		}
	}
}

func BenchmarkWrapWithNotFound(b *testing.B) {
	base := errors.New("not found")
	b.ReportAllocs()
	for b.Loop() {
		if err := WrapWithNotFound("API not found", base); err == nil {
			b.Fatal("expected an error")
		}
	}
}
