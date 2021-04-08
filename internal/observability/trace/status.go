//
// Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
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
	"go.opencensus.io/trace"
)

type Status = trace.Status

func FromGRPCStatus(code codes.Code, msg string) Status {
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

func StatusCodeOK(msg string) Status {
	return trace.Status{
		Code:    trace.StatusCodeOK,
		Message: msg,
	}
}

func StatusCodeCancelled(msg string) Status {
	return trace.Status{
		Code:    trace.StatusCodeCancelled,
		Message: msg,
	}
}

func StatusCodeUnknown(msg string) Status {
	return trace.Status{
		Code:    trace.StatusCodeUnknown,
		Message: msg,
	}
}

func StatusCodeInvalidArgument(msg string) Status {
	return trace.Status{
		Code:    trace.StatusCodeInvalidArgument,
		Message: msg,
	}
}

func StatusCodeDeadlineExceeded(msg string) Status {
	return trace.Status{
		Code:    trace.StatusCodeDeadlineExceeded,
		Message: msg,
	}
}

func StatusCodeNotFound(msg string) Status {
	return trace.Status{
		Code:    trace.StatusCodeNotFound,
		Message: msg,
	}
}

func StatusCodeAlreadyExists(msg string) Status {
	return trace.Status{
		Code:    trace.StatusCodeAlreadyExists,
		Message: msg,
	}
}

func StatusCodePermissionDenied(msg string) Status {
	return trace.Status{
		Code:    trace.StatusCodePermissionDenied,
		Message: msg,
	}
}

func StatusCodeResourceExhausted(msg string) Status {
	return trace.Status{
		Code:    trace.StatusCodeResourceExhausted,
		Message: msg,
	}
}

func StatusCodeFailedPrecondition(msg string) Status {
	return trace.Status{
		Code:    trace.StatusCodeFailedPrecondition,
		Message: msg,
	}
}

func StatusCodeAborted(msg string) Status {
	return trace.Status{
		Code:    trace.StatusCodeAborted,
		Message: msg,
	}
}

func StatusCodeOutOfRange(msg string) Status {
	return trace.Status{
		Code:    trace.StatusCodeOutOfRange,
		Message: msg,
	}
}

func StatusCodeUnimplemented(msg string) Status {
	return trace.Status{
		Code:    trace.StatusCodeUnimplemented,
		Message: msg,
	}
}

func StatusCodeInternal(msg string) Status {
	return trace.Status{
		Code:    trace.StatusCodeInternal,
		Message: msg,
	}
}

func StatusCodeUnavailable(msg string) Status {
	return trace.Status{
		Code:    trace.StatusCodeUnavailable,
		Message: msg,
	}
}

func StatusCodeDataLoss(msg string) Status {
	return trace.Status{
		Code:    trace.StatusCodeDataLoss,
		Message: msg,
	}
}

func StatusCodeUnauthenticated(msg string) Status {
	return trace.Status{
		Code:    trace.StatusCodeUnauthenticated,
		Message: msg,
	}
}
