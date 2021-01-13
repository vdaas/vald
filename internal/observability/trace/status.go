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
	"go.opencensus.io/trace"
)

func StatusCodeOK(msg string) trace.Status {
	return trace.Status{
		Code:    trace.StatusCodeOK,
		Message: msg,
	}
}

func StatusCodeCancelled(msg string) trace.Status {
	return trace.Status{
		Code:    trace.StatusCodeCancelled,
		Message: msg,
	}
}

func StatusCodeUnknown(msg string) trace.Status {
	return trace.Status{
		Code:    trace.StatusCodeUnknown,
		Message: msg,
	}
}

func StatusCodeInvalidArgument(msg string) trace.Status {
	return trace.Status{
		Code:    trace.StatusCodeInvalidArgument,
		Message: msg,
	}
}

func StatusCodeDeadlineExceeded(msg string) trace.Status {
	return trace.Status{
		Code:    trace.StatusCodeDeadlineExceeded,
		Message: msg,
	}
}

func StatusCodeNotFound(msg string) trace.Status {
	return trace.Status{
		Code:    trace.StatusCodeNotFound,
		Message: msg,
	}
}

func StatusCodeAlreadyExists(msg string) trace.Status {
	return trace.Status{
		Code:    trace.StatusCodeAlreadyExists,
		Message: msg,
	}
}

func StatusCodePermissionDenied(msg string) trace.Status {
	return trace.Status{
		Code:    trace.StatusCodePermissionDenied,
		Message: msg,
	}
}

func StatusCodeResourceExhausted(msg string) trace.Status {
	return trace.Status{
		Code:    trace.StatusCodeResourceExhausted,
		Message: msg,
	}
}

func StatusCodeFailedPrecondition(msg string) trace.Status {
	return trace.Status{
		Code:    trace.StatusCodeFailedPrecondition,
		Message: msg,
	}
}

func StatusCodeAborted(msg string) trace.Status {
	return trace.Status{
		Code:    trace.StatusCodeAborted,
		Message: msg,
	}
}

func StatusCodeOutOfRange(msg string) trace.Status {
	return trace.Status{
		Code:    trace.StatusCodeOutOfRange,
		Message: msg,
	}
}

func StatusCodeUnimplemented(msg string) trace.Status {
	return trace.Status{
		Code:    trace.StatusCodeUnimplemented,
		Message: msg,
	}
}

func StatusCodeInternal(msg string) trace.Status {
	return trace.Status{
		Code:    trace.StatusCodeInternal,
		Message: msg,
	}
}

func StatusCodeUnavailable(msg string) trace.Status {
	return trace.Status{
		Code:    trace.StatusCodeUnavailable,
		Message: msg,
	}
}

func StatusCodeDataLoss(msg string) trace.Status {
	return trace.Status{
		Code:    trace.StatusCodeDataLoss,
		Message: msg,
	}
}

func StatusCodeUnauthenticated(msg string) trace.Status {
	return trace.Status{
		Code:    trace.StatusCodeUnauthenticated,
		Message: msg,
	}
}
