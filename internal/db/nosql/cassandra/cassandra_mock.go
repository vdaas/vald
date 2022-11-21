// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package cassandra

import (
	"context"
	"net"

	"github.com/gocql/gocql"
)

type MockClusterConfig struct {
	CreateSessionFunc func() (*gocql.Session, error)
}

func (m *MockClusterConfig) CreateSession() (*gocql.Session, error) {
	return m.CreateSessionFunc()
}

type DialerMock struct {
	DialContextFunc func(ctx context.Context, network, addr string) (net.Conn, error)
}

func (dm *DialerMock) DialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	return dm.DialContextFunc(ctx, network, addr)
}

func (dm *DialerMock) GetDialer() func(ctx context.Context, network, addr string) (net.Conn, error) {
	return dm.DialContextFunc
}
func (*DialerMock) StartDialerCache(_ context.Context) {}
