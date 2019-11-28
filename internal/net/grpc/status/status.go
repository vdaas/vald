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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func WrapWithCanceled(err error) error {
	return status.New(codes.Canceled, err.Error()).Err()
}

func WrapWithUnknown(err error) error {
	return status.New(codes.Unknown, err.Error()).Err()
}

func WrapWithInvalidArgument(err error) error {
	return status.New(codes.InvalidArgument, err.Error()).Err()
}

func WrapWithDeadlineExceeded(err error) error {
	return status.New(codes.DeadlineExceeded, err.Error()).Err()
}

func WrapWithNotFound(err error) error {
	return status.New(codes.NotFound, err.Error()).Err()
}

func WrapWithAlreadyExists(err error) error {
	return status.New(codes.AlreadyExists, err.Error()).Err()
}

func WrapWithPermissionDenied(err error) error {
	return status.New(codes.PermissionDenied, err.Error()).Err()
}

func WrapWithResourceExhausted(err error) error {
	return status.New(codes.ResourceExhausted, err.Error()).Err()
}

func WrapWithFailedPrecondition(err error) error {
	return status.New(codes.FailedPrecondition, err.Error()).Err()
}

func WrapWithAborted(err error) error {
	return status.New(codes.Aborted, err.Error()).Err()
}

func WrapWithOutOfRange(err error) error {
	return status.New(codes.OutOfRange, err.Error()).Err()
}

func WrapWithUnimplemented(err error) error {
	return status.New(codes.Unimplemented, err.Error()).Err()
}

func WrapWithInternal(err error) error {
	return status.New(codes.Internal, err.Error()).Err()
}

func WrapWithUnavailable(err error) error {
	return status.New(codes.Unavailable, err.Error()).Err()
}

func WrapWithDataLoss(err error) error {
	return status.New(codes.DataLoss, err.Error()).Err()
}

func WrapWithUnauthenticated(err error) error {
	return status.New(codes.Unauthenticated, err.Error()).Err()
}
