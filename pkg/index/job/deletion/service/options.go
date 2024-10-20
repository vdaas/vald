package service

import (
	"github.com/vdaas/vald/internal/client/v1/client/discoverer"
	"github.com/vdaas/vald/internal/errors"
)

type Option func(_ *index) error

// WithDiscoverer returns Option that sets discoverer client.
func WithDiscoverer(client discoverer.Client) Option {
	return func(idx *index) error {
		if client == nil {
			return errors.NewErrCriticalOption("discoverer", client)
		}
		idx.client = client
		return nil
	}
}

// WithIndexingConcurrency returns Option that sets indexing concurrency.
func WithIndexingConcurrency(num int) Option {
	return func(idx *index) error {
		if num <= 0 {
			return errors.NewErrInvalidOption("indexingConcurrency", num)
		}
		idx.concurrency = num
		return nil
	}
}

// WithTargetAddrs returns Option that sets indexing target addresses.
func WithTargetAddrs(addrs ...string) Option {
	return func(idx *index) error {
		if len(addrs) != 0 {
			idx.targetAddrs = append(idx.targetAddrs, addrs...)
		}
		return nil
	}
}

// WithTargetIndexID returns Option that sets target deleting index ID.
func WithTargetIndexID(indexID string) Option {
	return func(idx *index) error {
		idx.targetIndexID = indexID
		return nil
	}
}
