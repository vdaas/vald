//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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
	FlushServer
	ObjectServer
}

type ServerWithFilter interface {
	Server
	FilterServer
}

type UnimplementedValdServer struct {
	UnimplementedInsertServer
	UnimplementedUpdateServer
	UnimplementedUpsertServer
	UnimplementedSearchServer
	UnimplementedRemoveServer
	UnimplementedFlushServer
	UnimplementedObjectServer
}

type UnimplementedValdServerWithFilter struct {
	UnimplementedValdServer
	UnimplementedFilterServer
}

type Client interface {
	InsertClient
	UpdateClient
	UpsertClient
	SearchClient
	RemoveClient
	FlushClient
	ObjectClient
}

type ClientWithFilter interface {
	Client
	FilterClient
}

const PackageName = "vald.v1"

const (
	InsertRPCServiceName = "Insert"
	UpdateRPCServiceName = "Update"
	UpsertRPCServiceName = "Upsert"
	SearchRPCServiceName = "Search"
	RemoveRPCServiceName = "Remove"
	FlushRPCServiceName  = "Flush"
	ObjectRPCServiceName = "Object"
	FilterRPCServiceName = "Filter"
)

const (
	InsertRPCName             = "Insert"
	StreamInsertRPCName       = "StreamInsert"
	MultiInsertRPCName        = "MultiInsert"
	InsertObjectRPCName       = "InsertObject"
	StreamInsertObjectRPCName = "StreamInsertObject"
	MultiInsertObjectRPCName  = "MultiInsertObject"

	UpdateRPCName             = "Update"
	StreamUpdateRPCName       = "StreamUpdate"
	MultiUpdateRPCName        = "MultiUpdate"
	UpdateObjectRPCName       = "UpdateObject"
	StreamUpdateObjectRPCName = "StreamUpdateObject"
	MultiUpdateObjectRPCName  = "MultiUpdateObject"

	UpsertRPCName             = "Upsert"
	StreamUpsertRPCName       = "StreamUpsert"
	MultiUpsertRPCName        = "MultiUpsert"
	UpsertObjectRPCName       = "UpsertObject"
	StreamUpsertObjectRPCName = "StreamUpsertObject"
	MultiUpsertObjectRPCName  = "MultiUpsertObject"

	SearchRPCName                   = "Search"
	SearchByIDRPCName               = "SearchByID"
	StreamSearchRPCName             = "StreamSearch"
	StreamSearchByIDRPCName         = "StreamSearchByID"
	MultiSearchRPCName              = "MultiSearch"
	MultiSearchByIDRPCName          = "MultiSearchByID"
	LinearSearchRPCName             = "LinearSearch"
	LinearSearchByIDRPCName         = "LinearSearchByID"
	StreamLinearSearchRPCName       = "StreamLinearSearch"
	StreamLinearSearchByIDRPCName   = "StreamLinearSearchByID"
	MultiLinearSearchRPCName        = "MultiLinearSearch"
	MultiLinearSearchByIDRPCName    = "MultiLinearSearchByID"
	SearchObjectRPCName             = "SearchObject"
	MultiSearchObjectRPCName        = "MultiSearchObject"
	LinearSearchObjectRPCName       = "LinearSearchObject"
	MultiLinearSearchObjectRPCName  = "MultiLinearSearchObject"
	StreamLinearSearchObjectRPCName = "StreamLinearSearchObject"
	StreamSearchObjectRPCName       = "StreamSearchObject"

	RemoveRPCName       = "Remove"
	StreamRemoveRPCName = "StreamRemove"
	MultiRemoveRPCName  = "MultiRemove"

	FlushRPCName = "Flush"

	ExistsRPCName          = "Exists"
	GetObjectRPCName       = "GetObject"
	StreamGetObjectRPCName = "StreamGetObject"
)

type client struct {
	InsertClient
	UpdateClient
	UpsertClient
	SearchClient
	RemoveClient
	FlushClient
	ObjectClient
}

func RegisterValdServer(s *grpc.Server, srv Server) {
	RegisterInsertServer(s, srv)
	RegisterUpdateServer(s, srv)
	RegisterUpsertServer(s, srv)
	RegisterSearchServer(s, srv)
	RegisterRemoveServer(s, srv)
	RegisterFlushServer(s, srv)
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
		NewFlushClient(conn),
		NewObjectClient(conn),
	}
}
