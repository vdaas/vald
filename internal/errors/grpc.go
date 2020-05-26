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

var (

	// gRPC

	ErrgRPCClientConnectionClose = func(name string, err error) error {
		return Wrapf(err, "%s's gRPC connection close error", name)
	}

	ErrInvalidGRPCPort = func(addr, host string, port uint16) error {
		return Errorf("invalid gRPC client connection port to addr: %s,\thost: %s\t port: %d", addr, host, port)
	}

	ErrInvalidGRPCClientConn = func(addr string) error {
		return Errorf("invalid gRPC client connection to %s", addr)
	}

	ErrGRPCLookupIPAddrNotFound = func(host string) error {
		return Errorf("vald internal gRPC client could not find ip addrs for %s", host)
	}

	ErrGRPCClientNotFound = New("vald internal gRPC client not found")

	ErrGRPCClientConnNotFound = func(addr string) error {
		return Errorf("gRPC client connection not found in %s", addr)
	}

	ErrRPCCallFailed = func(addr string, err error) error {
		return Wrapf(err, "addr: %s", addr)
	}

	ErrGRPCTargetAddrNotFound = New("grpc connection target not found")
)
