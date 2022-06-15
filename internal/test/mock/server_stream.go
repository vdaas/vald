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
