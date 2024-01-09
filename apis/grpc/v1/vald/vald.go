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

import (
	"github.com/vdaas/vald/apis/grpc/v1/mirror"
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

type ServerWithMirror interface {
	Server
	mirror.MirrorServer
}

type UnimplementedValdServer struct {
	UnimplementedInsertServer
	UnimplementedUpdateServer
	UnimplementedUpsertServer
	UnimplementedSearchServer
	UnimplementedRemoveServer
	UnimplementedObjectServer
}

type UnimplementedValdServerWithFilter struct {
	UnimplementedValdServer
	UnimplementedFilterServer
}

type UnimplementedValdServerWithMirror struct {
	UnimplementedValdServer
	mirror.UnimplementedMirrorServer
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

type ClientWithMirror interface {
	Client
	mirror.MirrorClient
}

const PackageName = "vald.v1"

const (
	InsertRPCServiceName = "Insert"
	UpdateRPCServiceName = "Update"
	UpsertRPCServiceName = "Upsert"
	SearchRPCServiceName = "Search"
	RemoveRPCServiceName = "Remove"
	ObjectRPCServiceName = "Object"
	FilterRPCServiceName = "Filter"
	MirrorRPCServiceName = "Mirror"
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

	ExistsRPCName           = "Exists"
	GetObjectRPCName        = "GetObject"
	GetTimestampRPCName     = "GetTimestamp"
	StreamGetObjectRPCName  = "StreamGetObject"
	StreamListObjectRPCName = "StreamListObject"

	RegisterRPCName = "Register"
)

type client struct {
	InsertClient
	UpdateClient
	UpsertClient
	SearchClient
	RemoveClient
	ObjectClient
}

type clientWithMirror struct {
	Client
	mirror.MirrorClient
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

func RegisterValdServerWithMirror(s *grpc.Server, srv ServerWithMirror) {
	RegisterValdServer(s, srv)
	mirror.RegisterMirrorServer(s, srv)
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

func NewValdClientWithMirror(conn *grpc.ClientConn) ClientWithMirror {
	return &clientWithMirror{
		Client:       NewValdClient(conn),
		MirrorClient: mirror.NewMirrorClient(conn),
	}
}
