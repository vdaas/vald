//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
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
	"github.com/vdaas/vald/internal/observability/attribute"
	ocodes "go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

type Attributes = []attribute.KeyValue

const (
	grpcMsgAttributeKey = attribute.Key("rpc.grpc.message")
	StatusOK            = ocodes.Ok
	StatusError         = ocodes.Error
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
		grpcMsgAttributeKey.String(msg),
	}
}

func StatusCodeCancelled(msg string) []attribute.KeyValue {
	return []attribute.KeyValue{
		semconv.RPCGRPCStatusCodeCancelled,
		grpcMsgAttributeKey.String(msg),
	}
}

func StatusCodeUnknown(msg string) []attribute.KeyValue {
	return []attribute.KeyValue{
		semconv.RPCGRPCStatusCodeUnknown,
		grpcMsgAttributeKey.String(msg),
	}
}

func StatusCodeInvalidArgument(msg string) []attribute.KeyValue {
	return []attribute.KeyValue{
		semconv.RPCGRPCStatusCodeInvalidArgument,
		grpcMsgAttributeKey.String(msg),
	}
}

func StatusCodeDeadlineExceeded(msg string) []attribute.KeyValue {
	return []attribute.KeyValue{
		semconv.RPCGRPCStatusCodeDeadlineExceeded,
		grpcMsgAttributeKey.String(msg),
	}
}

func StatusCodeNotFound(msg string) []attribute.KeyValue {
	return []attribute.KeyValue{
		semconv.RPCGRPCStatusCodeNotFound,
		grpcMsgAttributeKey.String(msg),
	}
}

func StatusCodeAlreadyExists(msg string) []attribute.KeyValue {
	return []attribute.KeyValue{
		semconv.RPCGRPCStatusCodeAlreadyExists,
		grpcMsgAttributeKey.String(msg),
	}
}

func StatusCodePermissionDenied(msg string) []attribute.KeyValue {
	return []attribute.KeyValue{
		semconv.RPCGRPCStatusCodePermissionDenied,
		grpcMsgAttributeKey.String(msg),
	}
}

func StatusCodeResourceExhausted(msg string) []attribute.KeyValue {
	return []attribute.KeyValue{
		semconv.RPCGRPCStatusCodeResourceExhausted,
		grpcMsgAttributeKey.String(msg),
	}
}

func StatusCodeFailedPrecondition(msg string) []attribute.KeyValue {
	return []attribute.KeyValue{
		semconv.RPCGRPCStatusCodeFailedPrecondition,
		grpcMsgAttributeKey.String(msg),
	}
}

func StatusCodeAborted(msg string) []attribute.KeyValue {
	return []attribute.KeyValue{
		semconv.RPCGRPCStatusCodeAborted,
		grpcMsgAttributeKey.String(msg),
	}
}

func StatusCodeOutOfRange(msg string) []attribute.KeyValue {
	return []attribute.KeyValue{
		semconv.RPCGRPCStatusCodeOutOfRange,
		grpcMsgAttributeKey.String(msg),
	}
}

func StatusCodeUnimplemented(msg string) []attribute.KeyValue {
	return []attribute.KeyValue{
		semconv.RPCGRPCStatusCodeUnimplemented,
		grpcMsgAttributeKey.String(msg),
	}
}

func StatusCodeInternal(msg string) []attribute.KeyValue {
	return []attribute.KeyValue{
		semconv.RPCGRPCStatusCodeInternal,
		grpcMsgAttributeKey.String(msg),
	}
}

func StatusCodeUnavailable(msg string) []attribute.KeyValue {
	return []attribute.KeyValue{
		semconv.RPCGRPCStatusCodeUnavailable,
		grpcMsgAttributeKey.String(msg),
	}
}

func StatusCodeDataLoss(msg string) []attribute.KeyValue {
	return []attribute.KeyValue{
		semconv.RPCGRPCStatusCodeDataLoss,
		grpcMsgAttributeKey.String(msg),
	}
}

func StatusCodeUnauthenticated(msg string) []attribute.KeyValue {
	return []attribute.KeyValue{
		semconv.RPCGRPCStatusCodeUnauthenticated,
		grpcMsgAttributeKey.String(msg),
	}
}
