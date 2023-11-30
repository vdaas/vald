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

package reader

import (
	"github.com/vdaas/vald/internal/backoff"
	"github.com/vdaas/vald/internal/db/storage/blob/s3/reader/io"
	"github.com/vdaas/vald/internal/db/storage/blob/s3/sdk/s3/s3iface"
	"github.com/vdaas/vald/internal/sync/errgroup"
)

// Option represents the functional option for reader.
type Option func(r *reader)

var defaultOptions = []Option{
	WithErrGroup(errgroup.Get()),
	WithMaxChunkSize(512 * 1024 * 1024),
	WithBackoff(false),
	func(r *reader) {
		r.ctxio = io.New()
	},
}

// WithErrGroup returns the option to set the eg.
func WithErrGroup(eg errgroup.Group) Option {
	return func(r *reader) {
		if eg != nil {
			r.eg = eg
		}
	}
}

// WithService returns the option to set the service.
func WithService(s s3iface.S3API) Option {
	return func(r *reader) {
		if s != nil {
			r.service = s
		}
	}
}

// WithBucket returns the option to set the bucket.
func WithBucket(bucket string) Option {
	return func(r *reader) {
		r.bucket = bucket
	}
}

// WithMaxChunkSize retunrs the option to set the maxChunkSize.
func WithMaxChunkSize(size int64) Option {
	return func(r *reader) {
		r.maxChunkSize = size
	}
}

// WithBackoff returns the option to set the backoffEnabled.
func WithBackoff(enabled bool) Option {
	return func(r *reader) {
		r.backoffEnabled = enabled
	}
}

// WithBackoffOpts returns the option to set the backoffOpts.
func WithBackoffOpts(opts ...backoff.Option) Option {
	return func(r *reader) {
		if opts == nil {
			return
		}
		if r.backoffOpts == nil {
			r.backoffOpts = opts
			return
		}

		r.backoffOpts = append(r.backoffOpts, opts...)
	}
}
