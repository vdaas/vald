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
)

type Option func(c *client)

var (
	defaultOpts = []Option{
		WithErrGroup(errgroup.Get()),
	}
)

func WithErrGroup(eg errgroup.Group) Option {
	return func(c *client) {
		if eg != nil {
			c.eg = eg
		}
	}
}

func WithSession(sess *session.Session) Option {
	return func(c *client) {
		if sess != nil {
			c.session = sess
		}
	}
}

func WithBucket(bucket string) Option {
	return func(c *client) {
		c.bucket = bucket
	}
}

func WithMaxPartSize(size int64) Option {
	return func(c *client) {
		if size >= s3manager.MinUploadPartSize {
			c.maxPartSize = size
		}
	}
}

func WithMaxPartSizeKB(kb int) Option {
	return WithMaxPartSize(int64(kb) * 1024)
}

func WithMaxPartSizeMB(mb int) Option {
	return WithMaxPartSize(int64(mb) * 1024 * 1024)
}

func WithMaxPartSizeGB(gb int) Option {
	return WithMaxPartSize(int64(gb) * 1024 * 1024 * 1024)
}
