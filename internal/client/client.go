package client

import (
	"context"
	
	"github.com/vdaas/vald/apis/grpc/payload"
)

type ObjectID = payload.Object_ID
type ObjectIDs = payload.Object_IDs
type ObjectVector = payload.Object_Vector
type ObjectVectors = payload.Object_Vectors
type SearchRequest = payload.Search_Request
type SearchIDRequest = payload.Search_IDRequest
type SearchResponse = payload.Search_Response
type ControlCreateIndexRequest = payload.Controll_CreateIndexRequest
type InfoIndex = payload.Info_Index
type MetaObject = payload.Backup_MetaVector

type Client interface {
	Reader
	Writer
}

type Reader interface {
	Exists(context.Context, *ObjectID) (*ObjectID, error)
	Search(context.Context, *SearchRequest) (*SearchResponse, error)
	SearchByID(context.Context, *SearchIDRequest) (*SearchResponse, error)
	StreamSearch(context.Context, func() *SearchRequest, func(*SearchResponse, error)) error
	StreamSearchByID(context.Context, func() *SearchRequest, func(*SearchResponse, error)) error
}

type Writer interface {
	Insert(context.Context, *ObjectVector) error
	StreamInsert(context.Context, func() *ObjectVector, func(error))
	MultiInsert(context.Context, *ObjectVectors) error
	Update(context.Context, *ObjectVector) error
	StreamUpdate(context.Context, func() *ObjectVector, func(error))
	MultiUpdate(context.Context, *ObjectVectors) error
	Remove(context.Context, *ObjectID) error
	StreamRemove(context.Context, func() *ObjectID, func(error))
	MultiRemove(context.Context, *ObjectIDs) error
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
