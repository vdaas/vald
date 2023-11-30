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

// Package client provides vald component client interfaces
package client

import (
	"github.com/vdaas/vald/apis/grpc/v1/agent/core"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
)

type (
	ObjectID                  = payload.Object_ID
	ObjectIDs                 = payload.Object_IDs
	ObjectVector              = payload.Object_Vector
	ObjectVectors             = payload.Object_Vectors
	ObjectLocation            = payload.Object_Location
	ObjectLocations           = payload.Object_Locations
	SearchRequest             = payload.Search_Request
	SearchIDRequest           = payload.Search_IDRequest
	SearchResponse            = payload.Search_Response
	SearchResponses           = payload.Search_Responses
	InsertRequest             = payload.Insert_Request
	UpdateRequest             = payload.Update_Request
	UpsertRequest             = payload.Upsert_Request
	RemoveRequest             = payload.Remove_Request
	SearchMultiRequest        = payload.Search_MultiRequest
	SearchIDMultiRequest      = payload.Search_MultiIDRequest
	InsertMultiRequest        = payload.Insert_MultiRequest
	UpdateMultiRequest        = payload.Update_MultiRequest
	UpsertMultiRequest        = payload.Upsert_MultiRequest
	RemoveMultiRequest        = payload.Remove_MultiRequest
	ControlCreateIndexRequest = payload.Control_CreateIndexRequest
	InfoIndex                 = payload.Info_Index
	InfoIndexCount            = payload.Info_Index_Count
	Empty                     = payload.Empty
	SearchConfig              = payload.Search_Config
	ObjectDistance            = payload.Object_Distance

	Searcher     = vald.SearchClient
	Inserter     = vald.InsertClient
	Updater      = vald.UpdateClient
	Upsertor     = vald.UpsertClient
	Remover      = vald.RemoveClient
	ObjectReader = vald.ObjectClient
	Indexer      = core.AgentClient
)

type Client interface {
	Reader
	Writer
}

type Reader interface {
	Searcher
	ObjectReader
}

type Writer interface {
	Inserter
	Updater
	Upsertor
	Remover
}
