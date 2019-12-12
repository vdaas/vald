//
// Copyright (C) 2019 Vdaas.org Vald team ( kpango, kmrmt, rinx )
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

// "github.com/pkg/errors"

var (

	// gRPC

	ErrAgentClientNotConnected = New("agent client not connected")

	ErrgRPCClientConnectionClose = func(name string, err error) error {
		return Wrapf(err, "%s's gRPC connection close error", name)
	}

	ErrInvalidGRPCClientConn = func(addr string) error {
		return Errorf("invalid gRPC client connection to %s", addr)
	}

	ErrGRPCClientConnNotFound = func(addr string) error {
		return Errorf("gRPC client connection not found on %s", addr)
	}

	ErrRPCCallFailed = func(addr string, err error) error {
		return Wrapf(err, "addr: %s", addr)
	}
)
