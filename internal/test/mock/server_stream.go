//
// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
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
package mock

import (
	"context"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/net/grpc"
	"google.golang.org/grpc/metadata"
)

type StreamInsertServerMock struct {
	SendFunc func(*payload.Object_StreamLocation) error
	RecvFunc func() (*payload.Insert_Request, error)
	grpc.ServerStream
}

func (m *StreamInsertServerMock) Send(l *payload.Object_StreamLocation) error {
	return m.SendFunc(l)
}

func (m *StreamInsertServerMock) Recv() (*payload.Insert_Request, error) {
	return m.RecvFunc()
}

// ServerStreamMock implements grpc.ServerStream mock implementation.
type ServerStreamMock struct {
	SetHeaderFunc  func(metadata.MD) error
	SendHeaderFunc func(metadata.MD) error
	SetTrailerFunc func(metadata.MD)
	ContextFunc    func() context.Context
	SendMsgFunc    func(interface{}) error
	RecvMsgFunc    func(interface{}) error
}

func (m *ServerStreamMock) SetHeader(md metadata.MD) error {
	return m.SetHeaderFunc(md)
}

func (m *ServerStreamMock) SendHeader(md metadata.MD) error {
	return m.SendHeaderFunc(md)
}

func (m *ServerStreamMock) SetTrailer(md metadata.MD) {
	m.SetTrailerFunc(md)
}

func (m *ServerStreamMock) Context() context.Context {
	return m.ContextFunc()
}

func (m *ServerStreamMock) SendMsg(msg interface{}) error {
	return m.SendMsgFunc(msg)
}

func (m *ServerStreamMock) RecvMsg(msg interface{}) error {
	return m.RecvMsgFunc(msg)
}
