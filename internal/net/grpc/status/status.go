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

// Package status provides statuses and errors returned by grpc handler functions
package status

import (
	"fmt"
	"os"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/info"
	"github.com/vdaas/vald/internal/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
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

	Code = status.Code
)

func newStatus(code codes.Code, msg string, err error, details ...interface{}) (st *status.Status) {
	st = status.New(code, msg)

	data := errors.Errors_RPC{
		Type:   code.String(),
		Msg:    msg,
		Status: int64(code),
	}
	if err != nil {
		data.Error = err.Error()
	}

	if len(details) != 0 {
		data.Details = make([]string, 0, len(details))
		for _, detail := range details {
			switch v := detail.(type) {
			case *errors.Errors_RPC:
				data.Roots = append(data.Roots, v)
			case info.Detail:
				data.Details = append(data.Details, v.String())
			default:
				data.Details = append(data.Details, fmt.Sprintf("%#v", detail))
			}
		}
	}

	root := FromError(err)
	if root != nil {
		data.Roots = append(data.Roots, root)
	}

	data.Instance, err = os.Hostname()
	if err != nil {
		log.Warn("failed to fetch hostname:", err)
	}

	st, err = st.WithDetails(&data)
	if err != nil {
		log.Warn("failed to set error details:", err)
	}

	return st
}

func WrapWithCanceled(msg string, err error, details ...interface{}) error {
	return newStatus(codes.Canceled, msg, err, details...).Err()
}

func WrapWithUnknown(msg string, err error, details ...interface{}) error {
	return newStatus(codes.Unknown, msg, err, details...).Err()
}

func WrapWithInvalidArgument(msg string, err error, details ...interface{}) error {
	return newStatus(codes.InvalidArgument, msg, err, details...).Err()
}

func WrapWithDeadlineExceeded(msg string, err error, details ...interface{}) error {
	return newStatus(codes.DeadlineExceeded, msg, err, details...).Err()
}

func WrapWithNotFound(msg string, err error, details ...interface{}) error {
	return newStatus(codes.NotFound, msg, err, details...).Err()
}

func WrapWithAlreadyExists(msg string, err error, details ...interface{}) error {
	return newStatus(codes.AlreadyExists, msg, err, details...).Err()
}

func WrapWithPermissionDenied(msg string, err error, details ...interface{}) error {
	return newStatus(codes.PermissionDenied, msg, err, details...).Err()
}

func WrapWithResourceExhausted(msg string, err error, details ...interface{}) error {
	return newStatus(codes.ResourceExhausted, msg, err, details...).Err()
}

func WrapWithFailedPrecondition(msg string, err error, details ...interface{}) error {
	return newStatus(codes.FailedPrecondition, msg, err, details...).Err()
}

func WrapWithAborted(msg string, err error, details ...interface{}) error {
	return newStatus(codes.Aborted, msg, err, details...).Err()
}

func WrapWithOutOfRange(msg string, err error, details ...interface{}) error {
	return newStatus(codes.OutOfRange, msg, err, details...).Err()
}

func WrapWithUnimplemented(msg string, err error, details ...interface{}) error {
	return newStatus(codes.Unimplemented, msg, err, details...).Err()
}

func WrapWithInternal(msg string, err error, details ...interface{}) error {
	return newStatus(codes.Internal, msg, err, details...).Err()
}

func WrapWithUnavailable(msg string, err error, details ...interface{}) error {
	return newStatus(codes.Unavailable, msg, err, details...).Err()
}

func WrapWithDataLoss(msg string, err error, details ...interface{}) error {
	return newStatus(codes.DataLoss, msg, err, details...).Err()
}

func WrapWithUnauthenticated(msg string, err error, details ...interface{}) error {
	return newStatus(codes.Unauthenticated, msg, err, details...).Err()
}

func FromError(err error) *errors.Errors_RPC {
	for _, detail := range status.Convert(err).Details() {
		if err, ok := detail.(*errors.Errors_RPC); ok {
			return err
		}
	}
	return nil
}
