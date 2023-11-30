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
