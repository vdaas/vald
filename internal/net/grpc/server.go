//
// Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
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

// Package grpc provides generic functionality for grpc
package grpc

import (
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

type (
	Server       = grpc.Server
	ServerOption = grpc.ServerOption
)

var ErrServerStopped = grpc.ErrServerStopped

func NewServer(opts ...ServerOption) *Server {
	return grpc.NewServer(opts...)
}

func Creds(c credentials.TransportCredentials) ServerOption {
	return grpc.Creds(c)
}

func KeepaliveParams(kp keepalive.ServerParameters) ServerOption {
	return grpc.KeepaliveParams(kp)
}

func MaxRecvMsgSize(size int) ServerOption {
	return grpc.MaxRecvMsgSize(size)
}

func MaxSendMsgSize(size int) ServerOption {
	return grpc.MaxSendMsgSize(size)
}

func InitialWindowSize(size int32) ServerOption {
	return grpc.InitialWindowSize(size)
}

func InitialConnWindowSize(size int32) ServerOption {
	return grpc.InitialConnWindowSize(size)
}

func ReadBufferSize(size int) ServerOption {
	return grpc.ReadBufferSize(size)
}

func WriteBufferSize(size int) ServerOption {
	return grpc.WriteBufferSize(size)
}

func ConnectionTimeout(d time.Duration) ServerOption {
	return grpc.ConnectionTimeout(d)
}

func MaxHeaderListSize(size uint32) ServerOption {
	return grpc.MaxHeaderListSize(size)
}

func HeaderTableSize(size uint32) ServerOption {
	return grpc.HeaderTableSize(size)
}
