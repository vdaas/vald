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
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/gogo/status"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/info"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/errdetails"
	"github.com/vdaas/vald/internal/net/grpc/proto"
)

type Status = status.Status

var Code = status.Code

func New(c codes.Code, msg string) *Status {
	return status.New(c, msg)
}

func newStatus(code codes.Code, msg string, err error, details ...interface{}) (st *Status) {
	st = New(code, msg)

	messages := make([]proto.Message, 0, 4)
	debugFunc := func(v *info.Detail) *errdetails.DebugInfo {
		debug := &errdetails.DebugInfo{
			Detail: fmt.Sprintf(
				"Version: %s,Name: %s, GitCommit: %s, BuildTime: %s, NGT_Version: %s ,Go_Version: %s, GOARCH: %s, GOOS: %s, CGO_Enabled: %s, BuildCPUInfo: [%s]",
				v.Version,
				v.ServerName,
				v.GitCommit,
				v.BuildTime,
				v.NGTVersion,
				v.GoVersion,
				v.GoArch,
				v.GoOS,
				v.CGOEnabled,
				strings.Join(v.BuildCPUInfoFlags, ", "),
			),
		}
		if debug.GetStackEntries() == nil {
			debug.StackEntries = make([]string, 0, len(v.StackTrace))
		}
		for i, stack := range v.StackTrace {
			debug.StackEntries = append(
				debug.GetStackEntries(),
				fmt.Sprintf("id: %d stack_trace: %s", i, stack.String()),
			)
		}
		return debug
	}
	if len(details) != 0 {
		for _, detail := range details {
			switch v := detail.(type) {
			case *info.Detail:
				messages = append(messages, debugFunc(v))
			case info.Detail:
				messages = append(messages, debugFunc(&v))
			case proto.Message:
				messages = append(messages, v)
			case *proto.Message:
				messages = append(messages, *v)
			}
		}
	}

	if err != nil {
		messages = append(messages, &errdetails.ErrorInfo{
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

	prevSt, ok := FromError(err)
	if ok {
		for _, detail := range prevSt.Details() {
			switch v := detail.(type) {
			case *info.Detail:
				messages = append(messages, debugFunc(v))
			case info.Detail:
				messages = append(messages, debugFunc(&v))
			case proto.Message:
				messages = append(messages, v)
			case *proto.Message:
				messages = append(messages, *v)
			}
		}
	}

	st, err = st.WithDetails(messages...)
	if err != nil {
		log.Warn("failed to set error details:", err)
	}

	err = st.Err()
	if err != nil {
		switch st.Code() {
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
			log.Warn(err.Error())
		}
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

func Error(code codes.Code, msg string) error {
	return status.Error(code, msg)
}

func Errorf(code codes.Code, format string, args ...interface{}) error {
	return status.Errorf(code, format, args...)
}

func ParseError(
	err error,
	defaultCode codes.Code,
	defaultMsg string,
	details ...interface{},
) (st *Status, msg string, rerr error) {
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
		rerr = errors.Wrap(st.Err(), err.Error())
		if st == nil || st.Message() == "" {
			msg = rerr.Error()
		} else {
			msg = st.Message()
		}
		return st, msg, errors.Wrap(st.Err(), err.Error())
	}

	pms := make([]proto.Message, 0, len(details))
	for _, detail := range details {
		pm, ok := detail.(proto.Message)
		if ok {
			pms = append(pms, pm)
		}
	}
	sst, err := st.WithDetails(pms...)
	if err == nil {
		st = sst
		err = st.Err()
	} else {
		log.Error(err)
	}
	err = st.Err()
	if err == nil {
		msg = st.Message()
	} else {
		msg = st.Err().Error()
		rerr = err
		switch st.Code() {
		case codes.Internal:
			log.Error(rerr.Error())
		case codes.Unavailable,
			codes.ResourceExhausted:
			log.Warn(rerr.Error())
		}
	}
	return st, msg, rerr
}

func FromError(err error) (st *Status, ok bool) {
	if err == nil {
		return nil, false
	}
	root := err
	defer func() {
		if !ok {
			return
		}
		ierr := errors.Unwrap(err)
		sst, ok := FromError(ierr)
		if ok && sst != nil && sst.Err() != nil {
			pms := make([]interface{}, 0, len(st.Details())+len(sst.Details())+1)
			for _, detail := range append(st.Details(), sst.Details()) {
				pm, ok := detail.(proto.Message)
				if ok {
					pms = append(pms, pm)
				}
			}
			pms = append(pms, &errdetails.ErrorInfo{
				Domain: fmt.Sprintf("code: %d, message: %s", sst.Code(), sst.Message()),
				Reason: sst.Err().Error(),
			})
			st = newStatus(st.Code(), st.Message(), st.Err(), pms...)
		}
	}()

	for {
		if st, ok = status.FromError(err); ok && st != nil {
			return st, true
		}
		if uerr := errors.Unwrap(err); uerr != nil {
			err = uerr
		} else {
			err = root
			for {
				switch err {
				case context.DeadlineExceeded:
					st = newStatus(codes.DeadlineExceeded, root.Error(), errors.Unwrap(err))
					return st, true
				case context.Canceled:
					st = newStatus(codes.Canceled, root.Error(), errors.Unwrap(err))
					return st, true
				case nil:
					st = New(codes.Unknown, root.Error())
					return st, false
				}
				if uerr := errors.Unwrap(err); uerr == nil {
					return New(codes.Unknown, root.Error()), false
				} else {
					err = uerr
				}
			}
		}
	}
}
