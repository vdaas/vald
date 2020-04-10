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

func TestErrCassandraNotFoundIdentity_Error(t *testing.T) {
	tests := []struct {
		name string
		e    *ErrCassandraNotFoundIdentity
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.Error(); got != tt.want {
				t.Errorf("ErrCassandraNotFoundIdentity.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsErrCassandraNotFound(t *testing.T) {
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
			if got := IsErrCassandraNotFound(tt.args.err); got != tt.want {
				t.Errorf("IsErrCassandraNotFound() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrCassandraUnavailableIdentity_Error(t *testing.T) {
	tests := []struct {
		name string
		e    *ErrCassandraUnavailableIdentity
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.Error(); got != tt.want {
				t.Errorf("ErrCassandraUnavailableIdentity.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsErrCassandraUnavailable(t *testing.T) {
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
			if got := IsErrCassandraUnavailable(tt.args.err); got != tt.want {
				t.Errorf("IsErrCassandraUnavailable() = %v, want %v", got, tt.want)
			}
		})
	}
}
