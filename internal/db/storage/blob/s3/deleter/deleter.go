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
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/vdaas/vald/internal/db/storage/blob/s3/sdk/s3/s3iface"
	"github.com/vdaas/vald/internal/errors"
)

type Deleter interface {
	Delete(key string) error
}

type deleter struct {
	service s3iface.S3API
	bucket  string
}

func New(opts ...Option) (Deleter, error) {
	d := new(deleter)
	for _, opt := range append(defaultOpts, opts...) {
		if err := opt(d); err != nil {
			return nil, err
		}
	}

	if d.service == nil {
		return nil, errors.NewErrInvalidOption("service", d.service)
	}
	if len(d.bucket) == 0 {
		return nil, errors.NewErrInvalidOption("bucket", d.bucket)
	}

	return d, nil
}

func (d *deleter) Delete(key string) error {
	if len(key) == 0 {
		// TODO: use defined error
		return errors.New("key is empty")
	}

	_, err := d.service.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(d.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return err
	}

	err = d.service.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(d.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return err
	}

	return nil
}
