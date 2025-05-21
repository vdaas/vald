//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
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

package grpc

import (
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/encoding"
	_ "google.golang.org/grpc/encoding/proto"
	"google.golang.org/grpc/keepalive"
)

func init() {
	encoding.RegisterCodec(Codec{})
}

type (
	// Server represents a gRPC server to serve RPC requests.
	Server = grpc.Server

	// ServerOption represents a gRPC server option.
	ServerOption = grpc.ServerOption
)

// ErrServerStopped indicates that the operation is now illegal because of
// the server being stopped.
var ErrServerStopped = grpc.ErrServerStopped

// NewServer returns the gRPC server.
func NewServer(opts ...ServerOption) *Server {
	// skipcq: GO-S0902
	return grpc.NewServer(opts...)
}

// Creds is a alias of grpc.Creds that sets credentials for server connections.
func Creds(c credentials.TransportCredentials) ServerOption {
	return grpc.Creds(c)
}

// KeepaliveParams is a alias of grpc.KeepaliveParams that sets keepalive and max-age parameters for the server.
func KeepaliveParams(kp keepalive.ServerParameters) ServerOption {
	return grpc.KeepaliveParams(kp)
}

// KeepaliveEnforcementPolicy is a alias of grpc.KeepaliveEnforcementPolicy that sets keepalive enforcement policy for the server.
func KeepaliveEnforcementPolicy(kep keepalive.EnforcementPolicy) ServerOption {
	return grpc.KeepaliveEnforcementPolicy(kep)
}

// MaxRecvMsgSize is a alias of grpc.MaxRecvMsgSize to set the max message size in bytes the server can receive.
func MaxRecvMsgSize(size int) ServerOption {
	return grpc.MaxRecvMsgSize(size)
}

// MaxSendMsgSize is a alias of grpc.MaxSendMsgSize to set the max message size in bytes the server can send.
func MaxSendMsgSize(size int) ServerOption {
	return grpc.MaxSendMsgSize(size)
}

// InitialWindowSize is a alias of grpc.InitialWindowSize that sets window size for stream.
func InitialWindowSize(size int32) ServerOption {
	return grpc.InitialWindowSize(size)
}

// InitialConnWindowSize is a alias of grpc.InitialConnWindowSize that sets window size for a connection.
func InitialConnWindowSize(size int32) ServerOption {
	return grpc.InitialConnWindowSize(size)
}

// ReadBufferSize is a alias of grpc.ReadBufferSize that lets you set the size of read buffer.
func ReadBufferSize(size int) ServerOption {
	return grpc.ReadBufferSize(size)
}

// WriteBufferSize is a alias of grpc.WriteBufferSize to determines how much data can be batched before doing a write on the wire.
func WriteBufferSize(size int) ServerOption {
	return grpc.WriteBufferSize(size)
}

// ConnectionTimeout is a alias of grpc.ConnectionTimeout that sets the timeout for
// connection establishment (up to and including HTTP/2 handshaking) for all
// new connections.
func ConnectionTimeout(d time.Duration) ServerOption {
	return grpc.ConnectionTimeout(d)
}

// MaxHeaderListSize is a alias of grpc.MaxHeaderListSize that sets the max (uncompressed) size
// of header list that the server is prepared to accept.
func MaxHeaderListSize(size uint32) ServerOption {
	return grpc.MaxHeaderListSize(size)
}

// HeaderTableSize is a alias of grpc.HeaderTableSize that sets the size of dynamic
// header table for stream.
func HeaderTableSize(size uint32) ServerOption {
	return grpc.HeaderTableSize(size)
}

// MaxConcurrentStreams returns a ServerOption that will apply a limit on the number of concurrent streams to each ServerTransport.
func MaxConcurrentStreams(n uint32) ServerOption {
	return grpc.MaxConcurrentStreams(n)
}

// NumStreamWorkers returns a ServerOption that sets the number of worker goroutines that should be used to process incoming streams. Setting this to zero
// (default) will disable workers and spawn a new goroutine for each stream.
func NumStreamWorkers(n uint32) ServerOption {
	return grpc.NumStreamWorkers(n)
}

// SharedWriteBuffer allows reusing per-connection transport write buffer. If this option is set to true every connection will release the buffer after flushing
// the data on the wire.
func SharedWriteBuffer(val bool) ServerOption {
	return grpc.SharedWriteBuffer(val)
}

// WaitForHandlers cause Stop to wait until all outstanding method handlers have exited before returning. If false, Stop will return as soon as all connections
// have closed, but method handlers may still be running. By default, Stop does not wait for method handlers to return.
func WaitForHandlers(val bool) ServerOption {
	return grpc.WaitForHandlers(val)
}

/*
API References https://pkg.go.dev/google.golang.org/grpc#ServerOption

1. Already Implemented APIs
- func ConnectionTimeout(d time.Duration) ServerOption
- func Creds(c credentials.TransportCredentials) ServerOption
- func HeaderTableSize(s uint32) ServerOption
- func InitialConnWindowSize(s int32) ServerOption
- func InitialWindowSize(s int32) ServerOption
- func KeepaliveEnforcementPolicy(kep keepalive.EnforcementPolicy) ServerOption
- func KeepaliveParams(kp keepalive.ServerParameters) ServerOption
- func MaxConcurrentStreams(n uint32) ServerOption
- func MaxHeaderListSize(s uint32) ServerOption
- func MaxRecvMsgSize(m int) ServerOption
- func MaxSendMsgSize(m int) ServerOption
- func NumStreamWorkers(numServerWorkers uint32) ServerOption
- func ReadBufferSize(s int) ServerOption
- func SharedWriteBuffer(val bool) ServerOption
- func WaitForHandlers(w bool) ServerOption
- func WriteBufferSize(s int) ServerOption

2. Unnecessary for this package APIs
- func ChainStreamInterceptor(interceptors ...StreamServerInterceptor) ServerOption
- func ChainUnaryInterceptor(interceptors ...UnaryServerInterceptor) ServerOption
- func StreamInterceptor(i StreamServerInterceptor) ServerOption
- func UnaryInterceptor(i UnaryServerInterceptor) ServerOption

3. Experimental APIs
- func ForceServerCodec(codec encoding.Codec) ServerOption
- func ForceServerCodecV2(codecV2 encoding.CodecV2) ServerOption
- func InTapHandle(h tap.ServerInHandle) ServerOption
- func StatsHandler(h stats.Handler) ServerOption
- func UnknownServiceHandler(streamHandler StreamHandler) ServerOption

4. Deprecated APIs
- func CustomCodec(codec Codec) ServerOption
- func MaxMsgSize(m int) ServerOption
- func RPCCompressor(cp Compressor) ServerOption
- func RPCDecompressor(dc Decompressor) ServerOption
*/
