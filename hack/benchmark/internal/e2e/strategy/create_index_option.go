package strategy

import "github.com/vdaas/vald/internal/client"

type CreateIndexOption func(*createIndex)

var (
	defaultCreateIndexOptions = []CreateIndexOption{
		WithCreateIndexPoolSize(10000),
	}
)

func WithCreateIndexPoolSize(size int) CreateIndexOption {
	return func(ci *createIndex) {
		if size > 0 {
			ci.poolSize = uint32(size)
		}
	}
}

func WithCreateIndexClient(idxc client.Indexer) CreateIndexOption {
	return func(ci *createIndex) {
		if idxc != nil {
			ci.idxc = idxc
		}
	}
}
