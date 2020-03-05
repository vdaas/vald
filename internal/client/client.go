package client

import (
	"context"

	"github.com/vdaas/vald/apis/grpc/payload"
)

type (
	ObjectID                  = payload.Object_ID
	ObjectIDs                 = payload.Object_IDs
	ObjectVector              = payload.Object_Vector
	ObjectVectors             = payload.Object_Vectors
	SearchRequest             = payload.Search_Request
	SearchIDRequest           = payload.Search_IDRequest
	SearchResponse            = payload.Search_Response
	ControlCreateIndexRequest = payload.Control_CreateIndexRequest
	InfoIndex                 = payload.Info_Index
	MetaObject                = payload.Backup_MetaVector
	Empty                     = payload.Empty
	SearchConfig              = payload.Search_Config
	ObjectDistance            = payload.Object_Distance
	BackupMetaVector          = payload.Backup_MetaVector
)

type Client interface {
	Reader
	Writer
}

type Reader interface {
	Exists(context.Context, *ObjectID) (*ObjectID, error)
	Search(context.Context, *SearchRequest) (*SearchResponse, error)
	SearchByID(context.Context, *SearchIDRequest) (*SearchResponse, error)
	StreamSearch(context.Context, func() *SearchRequest, func(*SearchResponse, error)) error
	StreamSearchByID(context.Context, func() *SearchIDRequest, func(*SearchResponse, error)) error
}

type Writer interface {
	Insert(context.Context, *ObjectVector) error
	StreamInsert(context.Context, func() *ObjectVector, func(error)) error
	MultiInsert(context.Context, *ObjectVectors) error
	Update(context.Context, *ObjectVector) error
	StreamUpdate(context.Context, func() *ObjectVector, func(error)) error
	MultiUpdate(context.Context, *ObjectVectors) error
	Remove(context.Context, *ObjectID) error
	StreamRemove(context.Context, func() *ObjectID, func(error)) error
	MultiRemove(context.Context, *ObjectIDs) error
}

type Upserter interface {
	Upsert(context.Context, *ObjectVector) error
	MultiUpsert(context.Context, *ObjectVectors) error
	StreamUpsert(context.Context, func() *ObjectVector, func(error)) error
}

type ObjectReader interface {
	GetObject(context.Context, *ObjectID) (*ObjectVector, error)
	StreamGetObject(context.Context, func() *ObjectID, func(*ObjectVector, error)) error
}

type MetaObjectReader interface {
	GetObject(context.Context, *ObjectID) (*MetaObject, error)
	StreamGetObject(context.Context, func() *ObjectID, func(*MetaObject, error)) error
}

type Indexer interface {
	CreateIndex(context.Context, *ControlCreateIndexRequest) error
	SaveIndex(context.Context) error
	CreateAndSaveIndex(context.Context, *ControlCreateIndexRequest) error
	IndexInfo(context.Context) (*InfoIndex, error)
}
