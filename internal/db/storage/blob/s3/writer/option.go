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
)

// Option represents the functional option for writer.
type Option func(w *writer)

var (
	defaultOpts = []Option{
		WithErrGroup(errgroup.Get()),
		WithContentType("application/octet-stream"),
		WithMaxPartSize(64 * 1024 * 1024),
	}
)

// WithErrGroup returns the option to set eg for writer.
func WithErrGroup(eg errgroup.Group) Option {
	return func(w *writer) {
		if eg != nil {
			w.eg = eg
		}
	}
}

// WithService returns the option to set s for writer.
func WithService(s *s3.S3) Option {
	return func(w *writer) {
		if s != nil {
			w.service = s
		}
	}
}

// WithBucket returns the option to set bucket for writer.
func WithBucket(bucket string) Option {
	return func(w *writer) {
		if len(bucket) != 0 {
			w.bucket = bucket
		}
	}
}

// WithKey returns the option to set key for writer.
func WithKey(key string) Option {
	return func(w *writer) {
		if len(key) != 0 {
			w.key = key
		}
	}
}

// WithContentType returns the option to set ct for writer.
func WithContentType(ct string) Option {
	return func(w *writer) {
		if len(ct) != 0 {
			w.contentType = ct
		}
	}
}

// WithMaxPartSize returns the option to set max for writer.
func WithMaxPartSize(max int64) Option {
	return func(w *writer) {
		if max > s3manager.DefaultUploadPartSize {
			w.maxPartSize = max
		}
	}
}
