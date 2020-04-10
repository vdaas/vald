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

// Package errors provides error types and function
package errors

import "testing"

func TestErrMySQLNotFoundIdentity_Error(t *testing.T) {
	tests := []struct {
		name string
		e    *ErrMySQLNotFoundIdentity
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.Error(); got != tt.want {
				t.Errorf("ErrMySQLNotFoundIdentity.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsErrMySQLNotFound(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsErrMySQLNotFound(tt.args.err); got != tt.want {
				t.Errorf("IsErrMySQLNotFound() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrMySQLInvalidArgumentIdentity_Error(t *testing.T) {
	tests := []struct {
		name string
		e    *ErrMySQLInvalidArgumentIdentity
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.Error(); got != tt.want {
				t.Errorf("ErrMySQLInvalidArgumentIdentity.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsErrMySQLInvalidArgument(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsErrMySQLInvalidArgument(tt.args.err); got != tt.want {
				t.Errorf("IsErrMySQLInvalidArgument() = %v, want %v", got, tt.want)
			}
		})
	}
}
