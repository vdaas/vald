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
	"context"
	"io"
	"net/http"
	"reflect"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/vdaas/vald/internal/db/storage/blob"
	"github.com/vdaas/vald/internal/db/storage/blob/v2/s3/downloader"
	"github.com/vdaas/vald/internal/db/storage/blob/v2/s3/logger"
	"github.com/vdaas/vald/internal/db/storage/blob/v2/s3/uploader"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
)

type client struct {
	eg       errgroup.Group
	bucket   string
	s3client *s3.Client
	client   *http.Client
	logMode  aws.ClientLogMode
	dclient  downloader.Client
	uclient  *uploader.Client
	// uclient      *manager.New
	region       string
	maxPartSize  int64
	maxChunkSize int64
}

// New returns blob.Bucket implementation if no error occurs.
func New(opts ...Option) (b blob.Bucket, err error) {
	c := new(client)
	for _, opt := range append(defaultOptions, opts...) {
		if err := opt(c); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}

	return c, nil
}

// Open does nothing. Always returns nil.
func (c *client) Open(ctx context.Context) (err error) {
	l := logger.New()
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithHTTPClient(c.client),
		config.WithClientLogMode(c.logMode),
		config.WithLogger(l),
		config.WithDefaultRegion(c.region),
		config.WithRegion(c.region),
		config.WithRetryer(func() aws.Retryer {
			return aws.NopRetryer{}
		}),
		config.WithEndpointResolver(aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
			return aws.Endpoint{}, nil
		})),
		config.WithSharedConfigProfile(""),
		config.WithSharedConfigFiles(nil),
		config.WithSharedCredentialsFiles(nil),
		config.WithCustomCABundle(nil),
		config.WithCredentialsProvider(aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
			return aws.Credentials{
				AccessKeyID:     "",
				SecretAccessKey: "",
				SessionToken:    "",
				Source:          "",
				CanExpire:       false,
				Expires:         time.Now().AddDate(100, 0, 0),
			}, nil
		})),
	)
	if err != nil {
		return err
	}
	c.s3client = s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.Region = c.region
		o.Logger = l
		o.UseAccelerate = true
		o.UseARNRegion = true
		o.EndpointOptions.DisableHTTPS = false
	})
	c.dclient, err = downloader.New(
		downloader.WithAPIClient(c.s3client),
	)
	if err != nil {
		return err
	}
	c.uclient = uploader.New(c.s3client)
	return nil
}

// Close does nothing. Always returns nil.
func (c *client) Close() error {
	return nil
}

// Reader creates reader.Reader implementation and returns it.
// An error will be returned when the reader initialization fails or an error occurs in reader.Open.
func (c *client) Reader(ctx context.Context, key string) (rc io.ReadCloser, err error) {
	return c.dclient.Open(ctx, &s3.GetObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	})
}

// Writer creates writer.Writer implementation and returns it.
// An error will be returned when the writer initialization fails or an error occurs in writer.Open.
func (c *client) Writer(ctx context.Context, key string) (wc io.WriteCloser, err error) {
	pr, pw := io.Pipe()
	ul, err := c.uclient.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
		Body:   pr,
	}, func(u *uploader.Client) {
	})
	if err != nil {
		return nil, err
	}
	return pw, nil
}
