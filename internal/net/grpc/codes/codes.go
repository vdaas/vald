//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
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

package codes

import "google.golang.org/grpc/codes"

type Code = codes.Code

var (
	OK                 = codes.OK
	Canceled           = codes.Canceled
	Unknown            = codes.Unknown
	InvalidArgument    = codes.InvalidArgument
	DeadlineExceeded   = codes.DeadlineExceeded
	NotFound           = codes.NotFound
	AlreadyExists      = codes.AlreadyExists
	PermissionDenied   = codes.PermissionDenied
	ResourceExhausted  = codes.ResourceExhausted
	FailedPrecondition = codes.FailedPrecondition
	Aborted            = codes.Aborted
	OutOfRange         = codes.OutOfRange
	Unimplemented      = codes.Unimplemented
	Internal           = codes.Internal
	Unavailable        = codes.Unavailable
	DataLoss           = codes.DataLoss
	Unauthenticated    = codes.Unauthenticated
)

type CodeType interface {
	int | int8 | int32 | int64 | uint | uint8 | uint32 | uint64 | Code
}

func ToString[T CodeType](c T) string {
	switch Code(c) {
	case OK:
		return "OK"
	case Canceled:
		return "Canceled"
	case Unknown:
		return "Unknown"
	case InvalidArgument:
		return "InvalidArgument"
	case DeadlineExceeded:
		return "DeadlineExceeded"
	case NotFound:
		return "NotFound"
	case AlreadyExists:
		return "AlreadyExists"
	case PermissionDenied:
		return "PermissionDenied"
	case ResourceExhausted:
		return "ResourceExhausted"
	case FailedPrecondition:
		return "FailedPrecondition"
	case Aborted:
		return "Aborted"
	case OutOfRange:
		return "OutOfRange"
	case Unimplemented:
		return "Unimplemented"
	case Internal:
		return "Internal"
	case Unavailable:
		return "Unavailable"
	case DataLoss:
		return "DataLoss"
	case Unauthenticated:
		return "Unauthenticated"
	default:
		return "InvalidStatus"
	}
}
