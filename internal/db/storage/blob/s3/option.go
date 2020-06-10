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

package s3

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/unit"
)

type Option func(c *client) error

var (
	defaultOpts = []Option{
		WithErrGroup(errgroup.Get()),
	}
)

func WithErrGroup(eg errgroup.Group) Option {
	return func(c *client) error {
		if eg != nil {
			c.eg = eg
		}
		return nil
	}
}

func WithSession(sess *session.Session) Option {
	return func(c *client) error {
		if sess != nil {
			c.session = sess
		}
		return nil
	}
}

func WithBucket(bucket string) Option {
	return func(c *client) error {
		c.bucket = bucket
		return nil
	}
}

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
