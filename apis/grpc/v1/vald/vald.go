//
// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
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

type UnimplementedValdServerWithMetadata struct {
	UnimplementedValdServer
	UnimplementedSearchWithMetadataServer
	UnimplementedInsertWithMetadataServer
	UnimplementedObjectWithMetadataServer
	UnimplementedRemoveWithMetadataServer
	UnimplementedUpdateWithMetadataServer
	UnimplementedUpsertWithMetadataServer
}

type UnimplementedValdServerWithFilter struct {
	UnimplementedValdServer
	UnimplementedFilterServer
}

type Client interface {
	FlushClient
	IndexClient
	InsertClient
	InsertWithMetadataClient
	ObjectClient
	ObjectWithMetadataClient
	RemoveClient
	RemoveWithMetadataClient
	SearchClient
	SearchWithMetadataClient
	UpdateClient
	UpdateWithMetadataClient
	UpsertClient
	UpsertWithMetadataClient
}

type ClientWithFilter interface {
	Client
	FilterClient
}

const PackageName = "vald.v1"

const MetadataSpanName = "WithMetadata"

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
	InsertRPCName                   = "Insert"
	InsertWithMetadataRPCName       = "InsertWithMetadata"
	StreamInsertRPCName             = "StreamInsert"
	StreamInsertWithMetadataRPCName = "StreamInsertWithMetadata"
	MultiInsertRPCName              = "MultiInsert"
	MultiInsertWithMetadataRPCName  = "MultiInsertWithMetadata"
	InsertObjectRPCName             = "InsertObject"
	StreamInsertObjectRPCName       = "StreamInsertObject"
	MultiInsertObjectRPCName        = "MultiInsertObject"

	UpdateRPCName                      = "Update"
	UpdateWithMetadataRPCName          = "UpdateWithMetadata"
	StreamUpdateRPCName                = "StreamUpdate"
	StreamUpdateWithMetadataRPCName    = "StreamUpdateWithMetadata"
	MultiUpdateRPCName                 = "MultiUpdate"
	MultiUpdateWithMetadataRPCName     = "MultiUpdateWithMetadata"
	UpdateObjectRPCName                = "UpdateObject"
	StreamUpdateObjectRPCName          = "StreamUpdateObject"
	MultiUpdateObjectRPCName           = "MultiUpdateObject"
	UpdateTimestampRPCName             = "UpdateTimestamp"
	UpdateTimestampWithMetadataRPCName = "UpdateTimestampWithMetadata"

	UpsertRPCName                   = "Upsert"
	UpsertWithMetadataRPCName       = "UpsertWithMetadata"
	StreamUpsertRPCName             = "StreamUpsert"
	StreamUpsertWithMetadataRPCName = "StreamUpsertWithMetadata"
	MultiUpsertRPCName              = "MultiUpsert"
	MultiUpsertWithMetadataRPCName  = "MultiUpsertWithMetadata"
	UpsertObjectRPCName             = "UpsertObject"
	StreamUpsertObjectRPCName       = "StreamUpsertObject"
	MultiUpsertObjectRPCName        = "MultiUpsertObject"

	SearchRPCName                             = "Search"
	SearchWithMetadataRPCName                 = "SearchWithMetadata"
	SearchByIDRPCName                         = "SearchByID"
	SearchByIDWithMetadataRPCName             = "SearchByIDWithMetadata"
	StreamSearchRPCName                       = "StreamSearch"
	StreamSearchWithMetadataRPCName           = "StreamSearchWithMetadata"
	StreamSearchByIDRPCName                   = "StreamSearchByID"
	StreamSearchByIDWithMetadataRPCName       = "StreamSearchByIDWithMetadata"
	MultiSearchRPCName                        = "MultiSearch"
	MultiSearchWithMetadataRPCName            = "MultiSearchWithMetadata"
	MultiSearchByIDRPCName                    = "MultiSearchByID"
	MultiSearchByIDWithMetadataRPCName        = "MultiSearchByIDWithMetadata"
	LinearSearchRPCName                       = "LinearSearch"
	LinearSearchWithMetadataRPCName           = "LinearSearchWithMetadata"
	LinearSearchByIDRPCName                   = "LinearSearchByID"
	LinearSearchByIDWithMetadataRPCName       = "LinearSearchByIDWithMetadata"
	StreamLinearSearchRPCName                 = "StreamLinearSearch"
	StreamLinearSearchWithMetadataRPCName     = "StreamLinearSearchWithMetadata"
	StreamLinearSearchByIDRPCName             = "StreamLinearSearchByID"
	StreamLinearSearchByIDWithMetadataRPCName = "StreamLinearSearchByIDWithMetadata"
	MultiLinearSearchRPCName                  = "MultiLinearSearch"
	MultiLinearSearchWithMetadataRPCName      = "MultiLinearSearchWithMetadata"
	MultiLinearSearchByIDRPCName              = "MultiLinearSearchByID"
	MultiLinearSearchByIDWithMetadataRPCName  = "MultiLinearSearchByIDWithMetadata"
	SearchObjectRPCName                       = "SearchObject"
	MultiSearchObjectRPCName                  = "MultiSearchObject"
	LinearSearchObjectRPCName                 = "LinearSearchObject"
	MultiLinearSearchObjectRPCName            = "MultiLinearSearchObject"
	StreamLinearSearchObjectRPCName           = "StreamLinearSearchObject"
	StreamSearchObjectRPCName                 = "StreamSearchObject"

	RemoveRPCName                        = "Remove"
	RemoveWithMetadataRPCName            = "RemoveWithMetadata"
	StreamRemoveRPCName                  = "StreamRemove"
	StreamRemoveWithMetadataRPCName      = "StreamRemoveWithMetadata"
	MultiRemoveRPCName                   = "MultiRemove"
	MultiRemoveWithMetadataRPCName       = "MultiRemoveWithMetadata"
	RemoveByTimestampRPCName             = "RemoveByTimestamp"
	RemoveByTimestampWithMetadataRPCName = "RemoveByTimestampWithMetadata"

	FlushRPCName = "Flush"

	ExistsRPCName                       = "Exists"
	GetObjectRPCName                    = "GetObject"
	GetObjectWithMetadataRPCName        = "GetObjectWithMetadata"
	GetTimestampRPCName                 = "GetTimestamp"
	StreamGetObjectRPCName              = "StreamGetObject"
	StreamGetObjectWithMetadataRPCName  = "StreamGetObjectWithMetadata"
	StreamListObjectRPCName             = "StreamListObject"
	StreamListObjectWithMetadataRPCName = "StreamListObjectWithMetadata"

	IndexInfoRPCName             = "IndexInfo"
	IndexDetailRPCName           = "IndexDetail"
	IndexStatisticsRPCName       = "IndexStatistics"
	IndexStatisticsDetailRPCName = "IndexStatisticsDetail"
	IndexPropertyRPCName         = "IndexProperty"
)

type client struct {
	FlushClient
	IndexClient
	InsertClient
	InsertWithMetadataClient
	ObjectClient
	ObjectWithMetadataClient
	RemoveClient
	RemoveWithMetadataClient
	SearchClient
	SearchWithMetadataClient
	UpdateClient
	UpdateWithMetadataClient
	UpsertClient
	UpsertWithMetadataClient
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
		InsertWithMetadataClient: NewInsertWithMetadataClient(conn),
		ObjectClient: NewObjectClient(conn),
		ObjectWithMetadataClient: NewObjectWithMetadataClient(conn),
		RemoveClient: NewRemoveClient(conn),
		RemoveWithMetadataClient: NewRemoveWithMetadataClient(conn),
		SearchClient: NewSearchClient(conn),
		SearchWithMetadataClient: NewSearchWithMetadataClient(conn),
		UpdateClient: NewUpdateClient(conn),
		UpdateWithMetadataClient: NewUpdateWithMetadataClient(conn),
		UpsertClient: NewUpsertClient(conn),
		UpsertWithMetadataClient: NewUpsertWithMetadataClient(conn),
	}
}
