//
// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
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

// Package vald provides vald server interface
package vald

import grpc "google.golang.org/grpc"

type Server interface {
	FlushServer
	IndexServer
	InsertServer
	ObjectServer
	RemoveServer
	SearchServer
	UpdateServer
	UpsertServer
}

type ServerWithFilter interface {
	Server
	FilterServer
}

type UnimplementedValdServer struct {
	UnimplementedFlushServer
	UnimplementedIndexServer
	UnimplementedInsertServer
	UnimplementedObjectServer
	UnimplementedRemoveServer
	UnimplementedSearchServer
	UnimplementedUpdateServer
	UnimplementedUpsertServer
}

type UnimplementedValdServerWithFilter struct {
	UnimplementedValdServer
	UnimplementedFilterServer
}

type Client interface {
	FlushClient
	IndexClient
	InsertClient
	ObjectClient
	RemoveClient
	SearchClient
	UpdateClient
	UpsertClient
}

type ClientWithFilter interface {
	Client
	FilterClient
}

const PackageName = "vald.v1"

const (
	FilterRPCServiceName = "Filter"
	FlushRPCServiceName  = "Flush"
	IndexRPCServiceName  = "Index"
	InsertRPCServiceName = "Insert"
	ObjectRPCServiceName = "Object"
	RemoveRPCServiceName = "Remove"
	SearchRPCServiceName = "Search"
	UpdateRPCServiceName = "Update"
	UpsertRPCServiceName = "Upsert"
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

	RemoveRPCName            = "Remove"
	StreamRemoveRPCName      = "StreamRemove"
	MultiRemoveRPCName       = "MultiRemove"
	RemoveByTimestampRPCName = "RemoveByTimestamp"

	FlushRPCName = "Flush"

	ExistsRPCName           = "Exists"
	GetObjectRPCName        = "GetObject"
	GetTimestampRPCName     = "GetTimestamp"
	StreamGetObjectRPCName  = "StreamGetObject"
	StreamListObjectRPCName = "StreamListObject"

	IndexInfoRPCName   = "IndexInfo"
	IndexDetailRPCName = "IndexDetail"
)

type client struct {
	FlushClient
	IndexClient
	InsertClient
	ObjectClient
	RemoveClient
	SearchClient
	UpdateClient
	UpsertClient
}

func RegisterValdServer(s *grpc.Server, srv Server) {
	RegisterFlushServer(s, srv)
	RegisterIndexServer(s, srv)
	RegisterInsertServer(s, srv)
	RegisterObjectServer(s, srv)
	RegisterRemoveServer(s, srv)
	RegisterSearchServer(s, srv)
	RegisterUpdateServer(s, srv)
	RegisterUpsertServer(s, srv)
}

func RegisterValdServerWithFilter(s *grpc.Server, srv ServerWithFilter) {
	RegisterValdServer(s, srv)
	RegisterFilterServer(s, srv)
}

func NewValdClient(conn *grpc.ClientConn) Client {
	return &client{
		FlushClient:  NewFlushClient(conn),
		IndexClient:  NewIndexClient(conn),
		InsertClient: NewInsertClient(conn),
		ObjectClient: NewObjectClient(conn),
		RemoveClient: NewRemoveClient(conn),
		SearchClient: NewSearchClient(conn),
		UpdateClient: NewUpdateClient(conn),
		UpsertClient: NewUpsertClient(conn),
	}
}
