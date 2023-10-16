//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
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

	// ErrGRPCClientConnectionClose represents a function to generate an error that the gRPC connection couldn't close.
	ErrGRPCClientConnectionClose = func(name string, err error) error {
		return Wrapf(err, "%s's gRPC connection close error", name)
	}

	// ErrInvalidGRPCPort represents a function to generate an error that the gRPC port is invalid.
	ErrInvalidGRPCPort = func(addr, host string, port uint16) error {
		return Errorf("invalid gRPC client connection port to addr: %s,\thost: %s\t port: %d", addr, host, port)
	}

	// ErrInvalidGRPCClientConn represents a function to generate an error that the vald internal gRPC connection is invalid.
	ErrInvalidGRPCClientConn = func(addr string) error {
		return Errorf("invalid gRPC client connection to %s", addr)
	}

	// ErrGRPCLookupIPAddrNotFound represents a function to generate an error that the vald internal gRPC client couldn't find IP address.
	ErrGRPCLookupIPAddrNotFound = func(host string) error {
		return Errorf("vald internal gRPC client could not find ip addrs for %s", host)
	}

	// ErrGRPCClientNotFound represents an error that the vald internal gRPC client couldn't find.
	ErrGRPCClientNotFound = New("vald internal gRPC client not found")

	// ErrGRPCPoolConnectionNotFound represents an error that the vald internal gRPC client pool connection couldn't find.
	ErrGRPCPoolConnectionNotFound = New("vald internal gRPC pool connection not found")

	// ErrGRPCClientConnNotFound represents a function to generate an error that the gRPC client connection couldn't find.
	ErrGRPCClientConnNotFound = func(addr string) error {
		return Errorf("gRPC client connection not found in %s", addr)
	}

	// ErrGRPCClientStreamNotFound represents an error that the vald internal gRPC client couldn't find any gRPC client stream connection.
	ErrGRPCClientStreamNotFound = New("vald internal gRPC client gRPC client stream not found")

	// ErrRPCCallFailed represents a function to generate an error that the RPC call failed.
	ErrRPCCallFailed = func(addr string, err error) error {
		return Wrapf(err, "addr: %s", addr)
	}

	// ErrGRPCTargetAddrNotFound represents an error that the gRPC target address couldn't find.
	ErrGRPCTargetAddrNotFound = New("gRPC connection target not found")

	// ErrGRPCUnexpectedStatusError represents an error that the gRPC status code is undefined.
	ErrGRPCUnexpectedStatusError = func(code string, err error) error {
		return Wrapf(err, "unexcepted error detected: code %s", code)
	}

	// ErrInvalidProtoMessageType represents an error that the gRPC protocol buffers message type is invalid.
	ErrInvalidProtoMessageType = func(v interface{}) error {
		return Errorf("failed to marshal/unmarshal proto message, message type is %T (missing vtprotobuf/protobuf helpers)", v)
	}

	// ErrServerStreamClientRecv represents a function to generate an error that the gRPC client couldn't receive from stream.
	ErrServerStreamClientRecv = func(err error) error {
		return Wrap(err, "gRPC client failed to receive from stream")
	}

	// ErrServerStreamClientSend represents a function to generate an error that the gRPC server couldn't send to stream.
	ErrServerStreamServerSend = func(err error) error {
		return Wrap(err, "gRPC server failed to send to stream")
	}
)
