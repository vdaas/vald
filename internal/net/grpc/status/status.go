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

var Code = status.Code

func New(c codes.Code, msg string) *status.Status {
	return status.New(c, msg)
}

func newStatus(code codes.Code, msg string, err error, details ...interface{}) (st *status.Status) {
	st = New(code, msg)

	messages := make([]proto.Message, 0, 4)
	if len(details) != 0 {
		for _, detail := range details {
			switch v := detail.(type) {
			case info.Detail:
				debug := &errdetails.DebugInfo{
					Detail: fmt.Sprintf("Version: %s,Name: %s, GitCommit: %s, BuildTime: %s, NGT_Version: %s ,Go_Version: %s, GOARCH: %s, GOOS: %s, CGO_Enabled: %s, BuildCPUInfo: [%s]",
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
				if debug.StackEntries == nil {
					debug.StackEntries = make([]string, 0, len(v.StackTrace))
				}
				for i, stack := range v.StackTrace {
					debug.StackEntries = append(debug.StackEntries, fmt.Sprintf("id: %d stack_trace: %s", i, stack.String()))
				}
				messages = append(messages, debug)
			case proto.Message:
				messages = append(messages, v)
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
			dm, ok := detail.(proto.Message)
			if ok {
				messages = append(messages, dm)
			}
		}
	}

	st, err = st.WithDetails(messages...)
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

func Error(code codes.Code, msg string) error {
	return status.Error(code, msg)
}

func Errorf(code codes.Code, format string, args ...interface{}) error {
	return status.Errorf(code, format, args...)
}

func FromError(err error) (st *status.Status, ok bool) {
	if err == nil {
		return nil, false
	}
	root := err
	defer func() {
		if !ok {
			return
		}
		sst, ok := FromError(errors.Unwrap(err))
		if ok && sst != nil && sst.Err() != nil {
			pms := make([]proto.Message, 0, len(st.Details())+len(sst.Details())+1)
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
			ist, err := New(st.Code(), st.Message()).WithDetails(pms...)
			if err == nil {
				st = ist
			}
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
