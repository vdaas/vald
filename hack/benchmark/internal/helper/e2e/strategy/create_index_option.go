package strategy

import "github.com/vdaas/vald/internal/client"

type CreateIndexOption func(*createIndex)

var (
	defaultCreateIndexOptions = []CreateIndexOption{}
)

func WithCreateIndexPoolSize(size int) CreateIndexOption {
	return func(ci *createIndex) {
		if size > 0 {
			ci.poolSize = uint32(size)
		}
	}
}

func WithCreateIndexClient(cidx client.Indexer) CreateIndexOption {
	return func(ci *createIndex) {
		if cidx != nil {
			ci.cidx = cidx
		}
	}
}
