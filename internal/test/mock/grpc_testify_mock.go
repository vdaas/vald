// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package mock

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/backoff"
	"github.com/vdaas/vald/internal/net/grpc/pool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// ServerStreamTestifyMock is a testify mock struct for grpc.ServerStream.
// Define Send method on top of this for each specific usecases with your specific proto schema Go struct type.
type ServerStreamTestifyMock struct {
	mock.Mock
}

func (*ServerStreamTestifyMock) SendMsg(_ interface{}) error {
	return nil
}

func (*ServerStreamTestifyMock) SetHeader(metadata.MD) error {
	return nil
}

func (*ServerStreamTestifyMock) SendHeader(metadata.MD) error {
	return nil
}

func (*ServerStreamTestifyMock) SetTrailer(metadata.MD) {
}

func (*ServerStreamTestifyMock) Context() context.Context {
	return context.Background()
}

func (*ServerStreamTestifyMock) SendMsgWithContext(_ context.Context, _ interface{}) error {
	return nil
}

func (*ServerStreamTestifyMock) RecvMsg(_ interface{}) error {
	return nil
}

// ListObjectStreamMock is a testify mock struct for ListObjectStream based on ServerStreamTestifyMock
type ListObjectStreamMock struct {
	ServerStreamTestifyMock
}

func (losm *ListObjectStreamMock) Send(res *payload.Object_List_Response) error {
	args := losm.Called(res)
	return args.Error(0)
}

type ClientInternal struct {
	mock.Mock
}

type (
	CallOption = grpc.CallOption
	DialOption = pool.DialOption
	ClientConn = pool.ClientConn
)

func (c *ClientInternal) StartConnectionMonitor(ctx context.Context) (<-chan error, error) {
	args := c.Called(ctx)
	return args.Get(0).(<-chan error), args.Error(1)
}

func (c *ClientInternal) Connect(ctx context.Context, addr string, dopts ...DialOption) (pool.Conn, error) {
	args := c.Called(ctx, addr, dopts)
	return args.Get(0).(pool.Conn), args.Error(1)
}

func (c *ClientInternal) IsConnected(ctx context.Context, addr string) bool {
	args := c.Called(ctx, addr)
	return args.Bool(0)
}

func (c *ClientInternal) Disconnect(ctx context.Context, addr string) error {
	args := c.Called(ctx, addr)
	return args.Error(0)
}

func (c *ClientInternal) Range(ctx context.Context,
	f func(ctx context.Context,
		addr string,
		conn *ClientConn,
		copts ...CallOption) error) error {
	args := c.Called(ctx, f)
	return args.Error(0)
}

func (c *ClientInternal) RangeConcurrent(ctx context.Context,
	concurrency int,
	f func(ctx context.Context,
		addr string,
		conn *ClientConn,
		copts ...CallOption) error) error {
	args := c.Called(ctx, concurrency, f)
	return args.Error(0)
}

func (c *ClientInternal) OrderedRange(ctx context.Context,
	order []string,
	f func(ctx context.Context,
		addr string,
		conn *ClientConn,
		copts ...CallOption) error) error {
	args := c.Called(ctx, order, f)
	return args.Error(0)
}

func (c *ClientInternal) OrderedRangeConcurrent(ctx context.Context,
	order []string,
	concurrency int,
	f func(ctx context.Context,
		addr string,
		conn *ClientConn,
		copts ...CallOption) error) error {
	args := c.Called(ctx, order, concurrency, f)
	return args.Error(0)
}

func (c *ClientInternal) Do(ctx context.Context, addr string,
	f func(ctx context.Context,
		conn *ClientConn,
		copts ...CallOption) (interface{}, error)) (interface{}, error) {
	args := c.Called(ctx, addr, f)
	return args.Get(0), args.Error(1)
}

func (c *ClientInternal) RoundRobin(ctx context.Context, f func(ctx context.Context,
	conn *ClientConn,
	copts ...CallOption) (interface{}, error)) (interface{}, error) {
	args := c.Called(ctx, f)
	return args.Get(0), args.Error(1)
}

func (c *ClientInternal) GetDialOption() []DialOption {
	args := c.Called()
	return args.Get(0).([]DialOption)
}

func (c *ClientInternal) GetCallOption() []CallOption {
	args := c.Called()
	return args.Get(0).([]CallOption)
}

func (c *ClientInternal) GetBackoff() backoff.Backoff {
	args := c.Called()
	return args.Get(0).(backoff.Backoff)
}

func (c *ClientInternal) ConnectedAddrs() []string {
	args := c.Called()
	return args.Get(0).([]string)
}

func (c *ClientInternal) Close(ctx context.Context) error {
	args := c.Called(ctx)
	return args.Error(0)
}
