//
// Copyright (C) 2019-2020 Vdaas.org Vald team ( kpango, rinx, kmrmt )
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

package writer

import (
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
)

// Option represents the functional option for writer.
type Option func(w *writer) error

var (
	defaultOpts = []Option{
		WithErrGroup(errgroup.Get()),
		WithContentType("application/octet-stream"),
		WithMaxPartSize(64 * 1024 * 1024),
	}
)

// WithErrGroup returns the option to set eg for writer.
func WithErrGroup(eg errgroup.Group) Option {
	return func(w *writer) error {
		if eg == nil {
			return errors.ErrInvalidOption("errgroup", eg)
		}
		w.eg = eg
		return nil
	}
}

// WithService returns the option to set s for writer.
func WithService(s *s3.S3) Option {
	return func(w *writer) error {
		if s == nil {
			return errors.ErrInvalidOption("service", s)
		}
		w.service = s
		return nil
	}
}

// WithBucket returns the option to set bucket for writer.
func WithBucket(bucket string) Option {
	return func(w *writer) error {
		if len(bucket) == 0 {
			return errors.ErrInvalidOption("bucket", bucket)
		}
		w.bucket = bucket
		return nil
	}
}

// WithKey returns the option to set key for writer.
func WithKey(key string) Option {
	return func(w *writer) error {
		if len(key) == 0 {
			return errors.ErrInvalidOption("key", key)
		}
		w.key = key
		return nil
	}
}

// WithContentType returns the option to set ct for writer.
func WithContentType(ct string) Option {
	return func(w *writer) error {
		if len(ct) == 0 {
			return errors.ErrInvalidOption("contentType", ct)
		}
		w.contentType = ct
		return nil
	}
}

// WithMaxPartSize returns the option to set max for writer.
func WithMaxPartSize(max int64) Option {
	return func(w *writer) error {
		if max < s3manager.MinUploadPartSize {
			return errors.ErrInvalidOption("maxPartSize", max)
		}
		w.maxPartSize = max
		return nil
	}
}
