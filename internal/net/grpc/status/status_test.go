//
// Copyright (C) 2019-2020 Vdaas.org Vald team ( kpango, rinx, kmrmt )
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
	"reflect"
	"testing"

	"github.com/vdaas/vald/apis/grpc/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Test_newStatus(t *testing.T) {
	type args struct {
		code    codes.Code
		msg     string
		err     error
		details []interface{}
	}
	tests := []struct {
		name   string
		args   args
		wantSt *status.Status
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotSt := newStatus(tt.args.code, tt.args.msg, tt.args.err, tt.args.details...); !reflect.DeepEqual(gotSt, tt.wantSt) {
				t.Errorf("newStatus() = %v, want %v", gotSt, tt.wantSt)
			}
		})
	}
}

func TestWrapWithCanceled(t *testing.T) {
	type args struct {
		msg     string
		err     error
		details []interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := WrapWithCanceled(tt.args.msg, tt.args.err, tt.args.details...); (err != nil) != tt.wantErr {
				t.Errorf("WrapWithCanceled() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWrapWithUnknown(t *testing.T) {
	type args struct {
		msg     string
		err     error
		details []interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := WrapWithUnknown(tt.args.msg, tt.args.err, tt.args.details...); (err != nil) != tt.wantErr {
				t.Errorf("WrapWithUnknown() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWrapWithInvalidArgument(t *testing.T) {
	type args struct {
		msg     string
		err     error
		details []interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := WrapWithInvalidArgument(tt.args.msg, tt.args.err, tt.args.details...); (err != nil) != tt.wantErr {
				t.Errorf("WrapWithInvalidArgument() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWrapWithDeadlineExceeded(t *testing.T) {
	type args struct {
		msg     string
		err     error
		details []interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := WrapWithDeadlineExceeded(tt.args.msg, tt.args.err, tt.args.details...); (err != nil) != tt.wantErr {
				t.Errorf("WrapWithDeadlineExceeded() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWrapWithNotFound(t *testing.T) {
	type args struct {
		msg     string
		err     error
		details []interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := WrapWithNotFound(tt.args.msg, tt.args.err, tt.args.details...); (err != nil) != tt.wantErr {
				t.Errorf("WrapWithNotFound() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWrapWithAlreadyExists(t *testing.T) {
	type args struct {
		msg     string
		err     error
		details []interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := WrapWithAlreadyExists(tt.args.msg, tt.args.err, tt.args.details...); (err != nil) != tt.wantErr {
				t.Errorf("WrapWithAlreadyExists() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWrapWithPermissionDenied(t *testing.T) {
	type args struct {
		msg     string
		err     error
		details []interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := WrapWithPermissionDenied(tt.args.msg, tt.args.err, tt.args.details...); (err != nil) != tt.wantErr {
				t.Errorf("WrapWithPermissionDenied() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWrapWithResourceExhausted(t *testing.T) {
	type args struct {
		msg     string
		err     error
		details []interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := WrapWithResourceExhausted(tt.args.msg, tt.args.err, tt.args.details...); (err != nil) != tt.wantErr {
				t.Errorf("WrapWithResourceExhausted() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWrapWithFailedPrecondition(t *testing.T) {
	type args struct {
		msg     string
		err     error
		details []interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := WrapWithFailedPrecondition(tt.args.msg, tt.args.err, tt.args.details...); (err != nil) != tt.wantErr {
				t.Errorf("WrapWithFailedPrecondition() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWrapWithAborted(t *testing.T) {
	type args struct {
		msg     string
		err     error
		details []interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := WrapWithAborted(tt.args.msg, tt.args.err, tt.args.details...); (err != nil) != tt.wantErr {
				t.Errorf("WrapWithAborted() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWrapWithOutOfRange(t *testing.T) {
	type args struct {
		msg     string
		err     error
		details []interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := WrapWithOutOfRange(tt.args.msg, tt.args.err, tt.args.details...); (err != nil) != tt.wantErr {
				t.Errorf("WrapWithOutOfRange() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWrapWithUnimplemented(t *testing.T) {
	type args struct {
		msg     string
		err     error
		details []interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := WrapWithUnimplemented(tt.args.msg, tt.args.err, tt.args.details...); (err != nil) != tt.wantErr {
				t.Errorf("WrapWithUnimplemented() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWrapWithInternal(t *testing.T) {
	type args struct {
		msg     string
		err     error
		details []interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := WrapWithInternal(tt.args.msg, tt.args.err, tt.args.details...); (err != nil) != tt.wantErr {
				t.Errorf("WrapWithInternal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWrapWithUnavailable(t *testing.T) {
	type args struct {
		msg     string
		err     error
		details []interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := WrapWithUnavailable(tt.args.msg, tt.args.err, tt.args.details...); (err != nil) != tt.wantErr {
				t.Errorf("WrapWithUnavailable() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWrapWithDataLoss(t *testing.T) {
	type args struct {
		msg     string
		err     error
		details []interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := WrapWithDataLoss(tt.args.msg, tt.args.err, tt.args.details...); (err != nil) != tt.wantErr {
				t.Errorf("WrapWithDataLoss() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWrapWithUnauthenticated(t *testing.T) {
	type args struct {
		msg     string
		err     error
		details []interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := WrapWithUnauthenticated(tt.args.msg, tt.args.err, tt.args.details...); (err != nil) != tt.wantErr {
				t.Errorf("WrapWithUnauthenticated() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFromError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want *errors.Errors_RPC
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FromError(tt.args.err); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromError() = %v, want %v", got, tt.want)
			}
		})
	}
}
