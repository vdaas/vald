package operation

import "github.com/vdaas/vald/internal/client/v1/client"

type Option func(*operation)

func WithClient(c client.Client) Option {
	return func(o *operation) {
		o.client = c
	}
}

func WithIndexer(c client.Indexer) Option {
	return func(o *operation) {
		o.indexerC = c
	}
}
