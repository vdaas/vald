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

// Package trace provides trace functions.
package trace

import (
	"reflect"
	"testing"

	"go.opencensus.io/trace"
)

func TestStatusCodeOK(t *testing.T) {
	type args struct {
		msg string
	}
	tests := []struct {
		name string
		args args
		want trace.Status
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StatusCodeOK(tt.args.msg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StatusCodeOK() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStatusCodeCancelled(t *testing.T) {
	type args struct {
		msg string
	}
	tests := []struct {
		name string
		args args
		want trace.Status
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StatusCodeCancelled(tt.args.msg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StatusCodeCancelled() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStatusCodeUnknown(t *testing.T) {
	type args struct {
		msg string
	}
	tests := []struct {
		name string
		args args
		want trace.Status
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StatusCodeUnknown(tt.args.msg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StatusCodeUnknown() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStatusCodeInvalidArgument(t *testing.T) {
	type args struct {
		msg string
	}
	tests := []struct {
		name string
		args args
		want trace.Status
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StatusCodeInvalidArgument(tt.args.msg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StatusCodeInvalidArgument() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStatusCodeDeadlineExceeded(t *testing.T) {
	type args struct {
		msg string
	}
	tests := []struct {
		name string
		args args
		want trace.Status
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StatusCodeDeadlineExceeded(tt.args.msg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StatusCodeDeadlineExceeded() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStatusCodeNotFound(t *testing.T) {
	type args struct {
		msg string
	}
	tests := []struct {
		name string
		args args
		want trace.Status
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StatusCodeNotFound(tt.args.msg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StatusCodeNotFound() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStatusCodeAlreadyExists(t *testing.T) {
	type args struct {
		msg string
	}
	tests := []struct {
		name string
		args args
		want trace.Status
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StatusCodeAlreadyExists(tt.args.msg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StatusCodeAlreadyExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStatusCodePermissionDenied(t *testing.T) {
	type args struct {
		msg string
	}
	tests := []struct {
		name string
		args args
		want trace.Status
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StatusCodePermissionDenied(tt.args.msg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StatusCodePermissionDenied() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStatusCodeResourceExhausted(t *testing.T) {
	type args struct {
		msg string
	}
	tests := []struct {
		name string
		args args
		want trace.Status
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StatusCodeResourceExhausted(tt.args.msg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StatusCodeResourceExhausted() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStatusCodeFailedPrecondition(t *testing.T) {
	type args struct {
		msg string
	}
	tests := []struct {
		name string
		args args
		want trace.Status
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StatusCodeFailedPrecondition(tt.args.msg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StatusCodeFailedPrecondition() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStatusCodeAborted(t *testing.T) {
	type args struct {
		msg string
	}
	tests := []struct {
		name string
		args args
		want trace.Status
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StatusCodeAborted(tt.args.msg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StatusCodeAborted() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStatusCodeOutOfRange(t *testing.T) {
	type args struct {
		msg string
	}
	tests := []struct {
		name string
		args args
		want trace.Status
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StatusCodeOutOfRange(tt.args.msg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StatusCodeOutOfRange() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStatusCodeUnimplemented(t *testing.T) {
	type args struct {
		msg string
	}
	tests := []struct {
		name string
		args args
		want trace.Status
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StatusCodeUnimplemented(tt.args.msg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StatusCodeUnimplemented() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStatusCodeInternal(t *testing.T) {
	type args struct {
		msg string
	}
	tests := []struct {
		name string
		args args
		want trace.Status
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StatusCodeInternal(tt.args.msg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StatusCodeInternal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStatusCodeUnavailable(t *testing.T) {
	type args struct {
		msg string
	}
	tests := []struct {
		name string
		args args
		want trace.Status
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StatusCodeUnavailable(tt.args.msg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StatusCodeUnavailable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStatusCodeDataLoss(t *testing.T) {
	type args struct {
		msg string
	}
	tests := []struct {
		name string
		args args
		want trace.Status
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StatusCodeDataLoss(tt.args.msg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StatusCodeDataLoss() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStatusCodeUnauthenticated(t *testing.T) {
	type args struct {
		msg string
	}
	tests := []struct {
		name string
		args args
		want trace.Status
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StatusCodeUnauthenticated(tt.args.msg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StatusCodeUnauthenticated() = %v, want %v", got, tt.want)
			}
		})
	}
}
