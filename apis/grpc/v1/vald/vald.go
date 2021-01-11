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

// Package vald provides vald server interface
package vald

import (
	grpc "google.golang.org/grpc"
)

type Server interface {
	InsertServer
	UpdateServer
	UpsertServer
	SearchServer
	RemoveServer
	ObjectServer
}

type ServerWithFilter interface {
	Server
	FilterServer
}

type Client interface {
	InsertClient
	UpdateClient
	UpsertClient
	SearchClient
	RemoveClient
	ObjectClient
}

type ClientWithFilter interface {
	Client
	FilterClient
}

type client struct {
	InsertClient
	UpdateClient
	UpsertClient
	SearchClient
	RemoveClient
	ObjectClient
}

func RegisterValdServer(s *grpc.Server, srv Server) {
	RegisterInsertServer(s, srv)
	RegisterUpdateServer(s, srv)
	RegisterUpsertServer(s, srv)
	RegisterSearchServer(s, srv)
	RegisterRemoveServer(s, srv)
	RegisterObjectServer(s, srv)
}

func RegisterValdServerWithFilter(s *grpc.Server, srv ServerWithFilter) {
	RegisterValdServer(s, srv)
	RegisterFilterServer(s, srv)
}

func NewValdClient(conn *grpc.ClientConn) Client {
	return &client{
		NewInsertClient(conn),
		NewUpdateClient(conn),
		NewUpsertClient(conn),
		NewSearchClient(conn),
		NewRemoveClient(conn),
		NewObjectClient(conn),
	}
}
