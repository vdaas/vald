package agent

import (
	"github.com/vdaas/vald/internal/client"
)

type Client interface {
	client.Client
	client.ObjectReader
	client.Indexer
}
