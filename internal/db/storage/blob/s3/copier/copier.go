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
	"net/url"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/vdaas/vald/internal/db/storage/blob/s3/sdk/s3/s3iface"
	"github.com/vdaas/vald/internal/errors"
)

type Copier interface {
	Copy(from, to string) error
}

type copier struct {
	service s3iface.S3API
	bucket  string
}

func New(opts ...Option) (Copier, error) {
	c := new(copier)
	for _, opt := range append(defaultOpts, opts...) {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	if c.service == nil {
		return nil, errors.NewErrInvalidOption("service", c.service)
	}

	if len(c.bucket) == 0 {
		return nil, errors.NewErrInvalidOption("bucket", c.bucket)
	}

	return c, nil
}

func (c *copier) Copy(from, to string) error {
	if len(from) == 0 {
		return errors.New("from file path is empty")
	}
	if len(to) == 0 {
		return errors.New("to file path is empty")
	}
	// copy object
	_, err := c.service.CopyObject(&s3.CopyObjectInput{
		Bucket:     aws.String(c.bucket),
		CopySource: aws.String(url.PathEscape(c.bucket + "/" + from)),
		Key:        aws.String(to),
	})
	if err != nil {
		return err
	}

	// Wait to see for check copied object
	err = c.service.WaitUntilObjectExists(&s3.HeadObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(to),
	})
	if err != nil {
		return err
	}

	return nil
}
