//
// Copyright (C) 2019 Vdaas.org Vald team ( kpango, kou-m, rinx )
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
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
	"encoding/json"
	"fmt"
	"os"

	"github.com/vdaas/vald/apis/grpc/errors"
	"github.com/vdaas/vald/internal/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func newStatus(code codes.Code, msg string, detail interface{}, err error) *status.Status {
	st := status.New(code, msg)

	data := errors.Errors_RPC{
		Type:   code.String(),
		Title:  msg,
		Detail: fmt.Sprintf("%#v", detail),
		Status: int64(code),
		Error:  err.Error(),
	}

	data.Instance, err = os.Hostname()
	if err != nil {
		ebody, _ := json.Marshal(data)
		log.Debugf("error body: %s, msg: %v", string(ebody), err)
		log.Error(err)
	}

	dst, err := st.WithDetails(&data)
	if err != nil {
		ebody, _ := json.Marshal(data)
		log.Debugf("error body: %s, msg: %v", string(ebody), err)
		log.Error(err)
		return st
	}

	return dst
}

func WrapWithCanceled(msg string, detail interface{}, err error) error {
	return newStatus(codes.Canceled, msg, detail, err).Err()
}

func WrapWithUnknown(msg string, detail interface{}, err error) error {
	return newStatus(codes.Unknown, msg, detail, err).Err()
}

func WrapWithInvalidArgument(msg string, detail interface{}, err error) error {
	return newStatus(codes.InvalidArgument, msg, detail, err).Err()
}

func WrapWithDeadlineExceeded(msg string, detail interface{}, err error) error {
	return newStatus(codes.DeadlineExceeded, msg, detail, err).Err()
}

func WrapWithNotFound(msg string, detail interface{}, err error) error {
	return newStatus(codes.NotFound, msg, detail, err).Err()
}

func WrapWithAlreadyExists(msg string, detail interface{}, err error) error {
	return newStatus(codes.AlreadyExists, msg, detail, err).Err()
}

func WrapWithPermissionDenied(msg string, detail interface{}, err error) error {
	return newStatus(codes.PermissionDenied, msg, detail, err).Err()
}

func WrapWithResourceExhausted(msg string, detail interface{}, err error) error {
	return newStatus(codes.ResourceExhausted, msg, detail, err).Err()
}

func WrapWithFailedPrecondition(msg string, detail interface{}, err error) error {
	return newStatus(codes.FailedPrecondition, msg, detail, err).Err()
}

func WrapWithAborted(msg string, detail interface{}, err error) error {
	return newStatus(codes.Aborted, msg, detail, err).Err()
}

func WrapWithOutOfRange(msg string, detail interface{}, err error) error {
	return newStatus(codes.OutOfRange, msg, detail, err).Err()
}

func WrapWithUnimplemented(msg string, detail interface{}, err error) error {
	return newStatus(codes.Unimplemented, msg, detail, err).Err()
}

func WrapWithInternal(msg string, detail interface{}, err error) error {
	return newStatus(codes.Internal, msg, detail, err).Err()
}

func WrapWithUnavailable(msg string, detail interface{}, err error) error {
	return newStatus(codes.Unavailable, msg, detail, err).Err()
}

func WrapWithDataLoss(msg string, detail interface{}, err error) error {
	return newStatus(codes.DataLoss, msg, detail, err).Err()
}

func WrapWithUnauthenticated(msg string, detail interface{}, err error) error {
	return newStatus(codes.Unauthenticated, msg, detail, err).Err()
}

func FromError(err error) *errors.Errors_RPC {
	for _, detail := range status.Convert(err).Details() {
		switch t := detail.(type) {
		case *errors.Errors_RPC:
			return t
		}
	}
	return nil
}
