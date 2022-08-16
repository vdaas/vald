//
// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

// Package trace provides trace functions.
package trace

import (
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"go.opentelemetry.io/otel/attribute"
	ocodes "go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

const (
	msgAttributeKey = attribute.Key("rpc.grpc.message")

	StatusOK    = ocodes.Ok
	StatusError = ocodes.Error
)

func FromGRPCStatus(code codes.Code, msg string) []attribute.KeyValue {
	switch code {
	case codes.OK:
		return StatusCodeOK(msg)
	case codes.Canceled:
		return StatusCodeCancelled(msg)
	case codes.InvalidArgument:
		return StatusCodeInvalidArgument(msg)
	case codes.DeadlineExceeded:
		return StatusCodeDeadlineExceeded(msg)
	case codes.NotFound:
		return StatusCodeNotFound(msg)
	case codes.AlreadyExists:
		return StatusCodeAlreadyExists(msg)
	case codes.PermissionDenied:
		return StatusCodePermissionDenied(msg)
	case codes.ResourceExhausted:
		return StatusCodeResourceExhausted(msg)
	case codes.FailedPrecondition:
		return StatusCodeFailedPrecondition(msg)
	case codes.Aborted:
		return StatusCodeAborted(msg)
	case codes.OutOfRange:
		return StatusCodeOutOfRange(msg)
	case codes.Unimplemented:
		return StatusCodeUnimplemented(msg)
	case codes.Internal:
		return StatusCodeInternal(msg)
	case codes.Unavailable:
		return StatusCodeUnavailable(msg)
	case codes.DataLoss:
		return StatusCodeDataLoss(msg)
	case codes.Unauthenticated:
		return StatusCodeUnauthenticated(msg)
	}
	return StatusCodeUnknown(msg)
}

func StatusCodeOK(msg string) []attribute.KeyValue {
	return []attribute.KeyValue{
		semconv.RPCGRPCStatusCodeOk,
		msgAttributeKey.String(msg),
	}
}

func StatusCodeCancelled(msg string) []attribute.KeyValue {
	return []attribute.KeyValue{
		semconv.RPCGRPCStatusCodeCancelled,
		msgAttributeKey.String(msg),
	}
}

func StatusCodeUnknown(msg string) []attribute.KeyValue {
	return []attribute.KeyValue{
		semconv.RPCGRPCStatusCodeUnknown,
		msgAttributeKey.String(msg),
	}
}

func StatusCodeInvalidArgument(msg string) []attribute.KeyValue {
	return []attribute.KeyValue{
		semconv.RPCGRPCStatusCodeInvalidArgument,
		msgAttributeKey.String(msg),
	}
}

func StatusCodeDeadlineExceeded(msg string) []attribute.KeyValue {
	return []attribute.KeyValue{
		semconv.RPCGRPCStatusCodeDeadlineExceeded,
		msgAttributeKey.String(msg),
	}
}

func StatusCodeNotFound(msg string) []attribute.KeyValue {
	return []attribute.KeyValue{
		semconv.RPCGRPCStatusCodeNotFound,
		msgAttributeKey.String(msg),
	}
}

func StatusCodeAlreadyExists(msg string) []attribute.KeyValue {
	return []attribute.KeyValue{
		semconv.RPCGRPCStatusCodeAlreadyExists,
		msgAttributeKey.String(msg),
	}
}

func StatusCodePermissionDenied(msg string) []attribute.KeyValue {
	return []attribute.KeyValue{
		semconv.RPCGRPCStatusCodePermissionDenied,
		msgAttributeKey.String(msg),
	}
}

func StatusCodeResourceExhausted(msg string) []attribute.KeyValue {
	return []attribute.KeyValue{
		semconv.RPCGRPCStatusCodeResourceExhausted,
		msgAttributeKey.String(msg),
	}
}

func StatusCodeFailedPrecondition(msg string) []attribute.KeyValue {
	return []attribute.KeyValue{
		semconv.RPCGRPCStatusCodeFailedPrecondition,
		msgAttributeKey.String(msg),
	}
}

func StatusCodeAborted(msg string) []attribute.KeyValue {
	return []attribute.KeyValue{
		semconv.RPCGRPCStatusCodeAborted,
		msgAttributeKey.String(msg),
	}
}

func StatusCodeOutOfRange(msg string) []attribute.KeyValue {
	return []attribute.KeyValue{
		semconv.RPCGRPCStatusCodeOutOfRange,
		msgAttributeKey.String(msg),
	}
}

func StatusCodeUnimplemented(msg string) []attribute.KeyValue {
	return []attribute.KeyValue{
		semconv.RPCGRPCStatusCodeUnimplemented,
		msgAttributeKey.String(msg),
	}
}

func StatusCodeInternal(msg string) []attribute.KeyValue {
	return []attribute.KeyValue{
		semconv.RPCGRPCStatusCodeInternal,
		msgAttributeKey.String(msg),
	}
}

func StatusCodeUnavailable(msg string) []attribute.KeyValue {
	return []attribute.KeyValue{
		semconv.RPCGRPCStatusCodeUnavailable,
		msgAttributeKey.String(msg),
	}
}

func StatusCodeDataLoss(msg string) []attribute.KeyValue {
	return []attribute.KeyValue{
		semconv.RPCGRPCStatusCodeDataLoss,
		msgAttributeKey.String(msg),
	}
}

func StatusCodeUnauthenticated(msg string) []attribute.KeyValue {
	return []attribute.KeyValue{
		semconv.RPCGRPCStatusCodeUnauthenticated,
		msgAttributeKey.String(msg),
	}
}
