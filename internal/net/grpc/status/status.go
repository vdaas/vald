//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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
	"context"
	"os"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/info"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/errdetails"
	"github.com/vdaas/vald/internal/net/grpc/proto"
	"github.com/vdaas/vald/internal/net/grpc/types"
	spb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/status"
)

type (
	Status = status.Status
	Code   = codes.Code
)

func New(c codes.Code, msg string) *Status {
	return status.New(c, msg)
}

func newStatus(code codes.Code, msg string, err error, details ...interface{}) (st *Status) {
	st = New(code, msg)
	return withDetails(st, err, details...)
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

func Error(code codes.Code, msg string) error {
	return status.Error(code, msg)
}

func Errorf(code codes.Code, format string, args ...interface{}) error {
	return status.Errorf(code, format, args...)
}

func ParseError(err error, defaultCode codes.Code, defaultMsg string, details ...interface{}) (st *Status, msg string, rerr error) {
	if err == nil {
		st = newStatus(codes.OK, "", nil, details...)
		msg = st.Message()
		return st, msg, nil
	}
	var ok bool
	st, ok = FromError(err)
	if !ok || st == nil {
		if defaultCode == 0 {
			defaultCode = codes.Internal
		}
		if len(defaultMsg) == 0 {
			defaultMsg = "failed to parse grpc status from error"
		}
		st = newStatus(defaultCode, defaultMsg, err, details...)
		if st == nil || st.Message() == "" {
			msg = st.Err().Error()
		} else {
			msg = st.Message()
		}
		return st, msg, st.Err()
	}

	st = withDetails(st, err, details...)
	msg = st.Message()
	return st, msg, err
}

func FromError(err error) (st *Status, ok bool) {
	if err == nil {
		return nil, false
	}
	root := err
	for {
		if st, ok = status.FromError(err); ok && st != nil {
			return st, true
		}
		if uerr := errors.Unwrap(err); uerr != nil {
			err = uerr
		} else {
			switch {
			case errors.Is(root, context.DeadlineExceeded):
				st = newStatus(codes.DeadlineExceeded, root.Error(), errors.Unwrap(root))
				return st, true
			case errors.Is(root, context.Canceled):
				st = newStatus(codes.Canceled, root.Error(), errors.Unwrap(root))
				return st, true
			default:
				st = newStatus(codes.Unknown, root.Error(), errors.Unwrap(root))
				return st, false
			}
		}
	}
}

func withDetails(st *Status, err error, details ...interface{}) *Status {
	msgs := make([]proto.MessageV1, 0, len(details)*2)
	if err != nil {
		msgs = append(msgs, &errdetails.ErrorInfo{
			Reason: err.Error(),
			Domain: func() (hostname string) {
				var err error
				hostname, err = os.Hostname()
				if err != nil {
					log.Warn("failed to fetch hostname:", err)
				}
				return hostname
			}(),
		})
	}
	for _, detail := range details {
		switch v := detail.(type) {
		case spb.Status:
			msgs = append(msgs, proto.ToMessageV1(&v))
		case *spb.Status:
			msgs = append(msgs, proto.ToMessageV1(v))
		case status.Status:
			msgs = append(msgs, proto.ToMessageV1(&spb.Status{
				Code:    v.Proto().GetCode(),
				Message: v.Message(),
			}))
			for _, d := range v.Proto().Details {
				msgs = append(msgs, proto.ToMessageV1(errdetails.AnyToErrorDetail(d)))
			}
		case *status.Status:
			msgs = append(msgs, proto.ToMessageV1(&spb.Status{
				Code:    v.Proto().GetCode(),
				Message: v.Message(),
			}))
			for _, d := range v.Proto().Details {
				msgs = append(msgs, proto.ToMessageV1(errdetails.AnyToErrorDetail(d)))
			}
		case info.Detail:
			msgs = append(msgs, errdetails.DebugInfoFromInfoDetail(&v))
		case *info.Detail:
			msgs = append(msgs, errdetails.DebugInfoFromInfoDetail(v))
		case proto.Message:
			msgs = append(msgs, proto.ToMessageV1(v))
		case *proto.Message:
			msgs = append(msgs, proto.ToMessageV1(*v))
		case proto.MessageV1:
			msgs = append(msgs, v)
		case *proto.MessageV1:
			msgs = append(msgs, *v)
		case types.Any:
			msgs = append(msgs, proto.ToMessageV1(errdetails.AnyToErrorDetail(&v)))
		}
	}

	if len(msgs) != 0 {
		sst, err := st.WithDetails(msgs...)
		if err == nil && sst != nil {
			st = sst
		} else {
			log.Warn("failed to set error details:", err)
		}
	}

	Log(st.Code(), st.Err())

	return st
}

func Log(code codes.Code, err error) {
	if err != nil {
		switch code {
		case codes.Internal,
			codes.DataLoss:
			log.Error(err.Error())
		case codes.Unavailable,
			codes.ResourceExhausted:
			log.Warn(err.Error())
		case codes.FailedPrecondition,
			codes.InvalidArgument,
			codes.OutOfRange,
			codes.Unauthenticated,
			codes.PermissionDenied,
			codes.Unknown:
			log.Debug(err.Error())
		case codes.Aborted,
			codes.Canceled,
			codes.DeadlineExceeded,
			codes.AlreadyExists,
			codes.NotFound,
			codes.OK,
			codes.Unimplemented:
		default:
			log.Warn(errors.ErrGRPCUnexpectedStatusError(code.String(), err))
		}
	}
}
