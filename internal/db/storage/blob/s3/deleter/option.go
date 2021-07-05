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
package deleter

import (
	"github.com/vdaas/vald/internal/db/storage/blob/s3/sdk/s3"
	"github.com/vdaas/vald/internal/errors"
)

type Option func(d *deleter) error

var defaultOpts = []Option{}

// WithService returns the option to set s for deleter.
func WithService(s *s3.S3) Option {
	return func(d *deleter) error {
		if s == nil {
			return errors.NewErrInvalidOption("service", s)
		}
		d.service = s
		return nil
	}
}

// WithBucket returns the option to set bucket for deleter.
func WithBucket(bucket string) Option {
	return func(d *deleter) error {
		if len(bucket) == 0 {
			return errors.NewErrInvalidOption("bucket", bucket)
		}
		d.bucket = bucket
		return nil
	}
}
