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

import "github.com/aws/aws-sdk-go/service/s3"

type Option func(w *writer)

var (
	defaultOpts = []Option{
		WithMaxPartSize(5 * 1024 * 1024),
	}
)

func WithService(s *s3.S3) Option {
	return func(w *writer) {
		if s != nil {
			w.service = s
		}
	}
}

func WithBucket(bucket string) Option {
	return func(w *writer) {
		w.bucket = bucket
	}
}

func WithKey(key string) Option {
	return func(w *writer) {
		w.key = key
	}
}

func WithMaxPartSize(max int) Option {
	return func(w *writer) {
		w.maxPartSize = int64(max)
	}
}
