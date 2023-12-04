// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package service

import (
	"github.com/vdaas/vald/internal/client/v1/client/discoverer"
	"github.com/vdaas/vald/internal/errors"
)

// Option represents the functional option for index corrector.
type Option func(*correct) error

var defaultOpts = []Option{
	WithStreamListConcurrency(200),     //nolint:gomnd
	WithKvsAsyncWriteConcurrency(2048), //nolint:gomnd
}

// WithIndexReplica returns Option that sets index replica.
func WithIndexReplica(num int) Option {
	return func(c *correct) error {
		if num <= 1 {
			return errors.NewErrCriticalOption("indexReplica", num, errors.ErrIndexReplicaOne)
		}
		c.indexReplica = num
		return nil
	}
}

// WithDiscoverer returns Option that sets discoverer client.
func WithDiscoverer(client discoverer.Client) Option {
	return func(c *correct) error {
		if client == nil {
			return errors.NewErrCriticalOption("discoverer", client)
		}
		c.discoverer = client
		return nil
	}
}

// WithStreamListConcurrency returns Option that sets concurrency for StreamList field value.
func WithStreamListConcurrency(num int) Option {
	return func(c *correct) error {
		if num <= 0 {
			return errors.NewErrInvalidOption("streamListConcurrency", num)
		}
		c.streamListConcurrency = num
		return nil
	}
}

// WithKvsAsyncWriteConcurrency returns Option that sets concurrency for kvs async write.
func WithKvsAsyncWriteConcurrency(num int) Option {
	return func(c *correct) error {
		if num <= 0 {
			return errors.NewErrInvalidOption("kvsAsyncWriteConcurrency", num)
		}
		c.bboltAsyncWriteConcurrency = num
		return nil
	}
}
