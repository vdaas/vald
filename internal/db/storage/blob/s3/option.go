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

package s3

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/vdaas/vald/internal/backoff"
	"github.com/vdaas/vald/internal/db/storage/blob/s3/reader"
	"github.com/vdaas/vald/internal/db/storage/blob/s3/writer"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/unit"
)

// Option represents the functional option for client.
type Option func(c *client) error

var defaultOptions = []Option{
	WithErrGroup(errgroup.Get()),
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

// WithSession returns the option to set the session.
func WithSession(sess *session.Session) Option {
	return func(c *client) error {
		if sess != nil {
			c.session = sess
		}
		return nil
	}
}

// WithBucket returns the option to set bucket.
func WithBucket(bucket string) Option {
	return func(c *client) error {
		if len(bucket) != 0 {
			c.bucket = bucket
		}
		return nil
	}
}

// WithReader returns the option to set the reader.
func WithReader(r reader.Reader) Option {
	return func(c *client) error {
		if r != nil {
			c.reader = r
		}
		return nil
	}
}

// WithWriter returns the option to set the reader.
func WithWriter(w writer.Writer) Option {
	return func(c *client) error {
		if w != nil {
			c.writer = w
		}
		return nil
	}
}

// WithMaxPartSize returns the option to set maxPartSize.
// The value is set when the value exceeds the s3manager.MinUploadPartSize(1024 * 1024 * 5).
func WithMaxPartSize(size string) Option {
	return func(c *client) error {
		b, err := unit.ParseBytes(size)
		if err != nil {
			return err
		}

		if int64(b) >= s3manager.MinUploadPartSize {
			c.maxPartSize = int64(b)
		}

		return nil
	}
}

// WithMaxChunkSize returns the option to set maxChunkSize.
// The value is set when the value exceeds the s3manager.MinUploadPartSize(1024 * 1024 * 5).
func WithMaxChunkSize(size string) Option {
	return func(c *client) error {
		b, err := unit.ParseBytes(size)
		if err != nil {
			return err
		}

		if int64(b) >= s3manager.MinUploadPartSize {
			c.maxChunkSize = int64(b)
		}

		return nil
	}
}

// WithReaderBackoff returns the option to set readerBackoffEnabled.
func WithReaderBackoff(enabled bool) Option {
	return func(c *client) error {
		c.readerBackoffEnabled = enabled
		return nil
	}
}

// WithReaderBackoffOpts returns the option to set readerBackoffOpts.
func WithReaderBackoffOpts(opts ...backoff.Option) Option {
	return func(c *client) error {
		if opts == nil {
			return nil
		}
		if c.readerBackoffOpts != nil {
			c.readerBackoffOpts = append(c.readerBackoffOpts, opts...)
			return nil
		}

		c.readerBackoffOpts = opts

		return nil
	}
}
