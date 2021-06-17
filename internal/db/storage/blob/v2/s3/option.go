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

package s3

import (
	"github.com/vdaas/vald/internal/errgroup"
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

// WithBucket returns the option to set bucket.
func WithBucket(bucket string) Option {
	return func(c *client) error {
		if len(bucket) != 0 {
			c.bucket = bucket
		}
		return nil
	}
}
