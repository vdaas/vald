//
// Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
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
package downloader

import (
	"github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/vdaas/vald/internal/backoff"
	"github.com/vdaas/vald/internal/db/storage/blob/v2/s3/ua"
	"github.com/vdaas/vald/internal/errgroup"
)

// Option represents the functional option for client.
type Option func(c *client) error

const (
	defaultDownloadPartSize    = 1024 * 1024 * 5
	defaultDownloadConcurrency = 5
)

var defaultOptions = []Option{
	WithErrGroup(errgroup.Get()),
	WithPartSize(defaultDownloadPartSize),
	WithConcurrency(defaultDownloadConcurrency),
	WithClientOptions(func(o *s3.Options) {
		o.APIOptions = append(o.APIOptions,
			middleware.AddSDKAgentKey(middleware.FeatureMetadata, ua.Key))
	}),
}

// WithErrGroup returns the option to set the eg.
func WithErrGroup(eg errgroup.Group) Option {
	return func(c *client) error {
		if eg != nil {
			c.eg = eg
		}
		return nil
	}
}

func WithAPIClient(mdc manager.DownloadAPIClient) Option {
	return func(c *client) error {
		if mdc != nil {
			c.client = mdc
		}
		return nil
	}
}

func WithClientOptions(opts ...func(*s3.Options)) Option {
	return func(c *client) error {
		if opts == nil {
			return nil
		}
		if c.clientOptions == nil {
			c.clientOptions = opts
		} else {
			c.clientOptions = append(c.clientOptions, opts...)
		}
		return nil
	}
}

func WithPartSize(size int64) Option {
	return func(c *client) error {
		if size > 0 {
			c.partSize = size
		}
		return nil
	}
}

func WithBackoff(bo backoff.Backoff) Option {
	return func(c *client) error {
		if bo != nil {
			c.bo = bo
		}
		return nil
	}
}

func WithConcurrency(concurrency int) Option {
	return func(c *client) error {
		if concurrency > 0 {
			c.concurrency = concurrency
		}
		return nil
	}
}

func WithBufferProvider(bp manager.WriterReadFromProvider) Option {
	return func(c *client) error {
		if bp != nil {
			c.bufferProvider = bp
		}
		return nil
	}
}
