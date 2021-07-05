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
package copier

import (
	"github.com/vdaas/vald/internal/db/storage/blob/s3/sdk/s3"
	"github.com/vdaas/vald/internal/errors"
)

type Option func(c *copier) error

var defaultOpts = []Option{}

// WithService returns the option to set s for copier.
func WithService(s *s3.S3) Option {
	return func(c *copier) error {
		if s == nil {
			return errors.NewErrInvalidOption("service", s)
		}
		c.service = s
		return nil
	}
}

// WithBucket returns the option to set bucket for copier.
func WithBucket(bucket string) Option {
	return func(c *copier) error {
		if len(bucket) == 0 {
			return errors.NewErrInvalidOption("bucket", bucket)
		}
		c.bucket = bucket
		return nil
	}
}
