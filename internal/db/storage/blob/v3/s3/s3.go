// Package s3 provides s3 functions.
package s3

import (
	"context"
	"io"
	"net/http"
	"reflect"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go/logging"

	"github.com/vdaas/vald/internal/db/storage/blob"
	"github.com/vdaas/vald/internal/db/storage/blob/v3/s3/downloader"
	"github.com/vdaas/vald/internal/db/storage/blob/v3/s3/logger"
	"github.com/vdaas/vald/internal/db/storage/blob/v3/s3/uploader"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
)

type client struct {
	eg errgroup.Group

	s3client *s3.Client
	client   *http.Client
	dclient  downloader.Client
	uclient  uploader.Client

	endpoint        string
	bucket          string
	region          string
	accessKey       string
	secretAccessKey string
	token           string

	maxRetries              int
	forcePathStyle          bool
	useAccelerate           bool
	useARNRegion            bool
	useDualStack            bool
	enableSSL               bool
	enableEndpointDiscovery bool

	maxPartSize int64
	concurrency int
	logger      logging.Logger
}

// New returns blob.Bucket implementation if no error occurs.
func New(opts ...Option) (b blob.Bucket, err error) {
	c := new(client)
	for _, opt := range append(defaultOptions, opts...) {
		if err := opt(c); err != nil {
			werr := errors.ErrOptionFailed(err, reflect.ValueOf(opt))

			e := new(errors.ErrCriticalOption)
			if errors.As(err, &e) {
				log.Error(werr)
				return nil, werr
			}
			log.Warn(werr)
		}
	}
	c.logger = logger.New()
	return c, nil
}

// Open loads the configuration and creates s3 client.
func (c *client) Open(ctx context.Context) (err error) {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithDefaultRegion(c.region),
		config.WithRegion(c.region),
		config.WithHTTPClient(c.client),
		config.WithLogger(c.logger),
		config.WithLogConfigurationWarnings(true),
		config.WithRetryer(func() aws.Retryer {
			if c.maxRetries < 1 {
				return aws.NopRetryer{}
			}
			return retry.AddWithMaxAttempts(retry.NewStandard(), c.maxRetries)
		}),
		func(lo *config.LoadOptions) error {
			lo.EnableEndpointDiscovery = aws.EndpointDiscoveryDisabled
			if c.enableEndpointDiscovery {
				lo.EnableEndpointDiscovery = aws.EndpointDiscoveryEnabled
			}
			return nil
		},
	)
	if err != nil {
		return err
	}

	c.s3client = s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = c.forcePathStyle
		o.UseAccelerate = c.useAccelerate
		o.UseARNRegion = c.useARNRegion

		if len(c.endpoint) != 0 {
			o.EndpointResolver = s3.EndpointResolverFromURL(c.endpoint)
		}
		if c.useDualStack {
			o.EndpointOptions.UseDualStackEndpoint = aws.DualStackEndpointStateEnabled
		}
		if !c.enableSSL {
			o.EndpointOptions.DisableHTTPS = true
		}

		if len(c.accessKey) != 0 && len(c.secretAccessKey) != 0 {
			o.Credentials = credentials.NewStaticCredentialsProvider(
				c.accessKey,
				c.secretAccessKey,
				c.token,
			)
		}
	})

	if c.dclient == nil {
		c.dclient, err = downloader.New(
			downloader.WithAPIClient(c.s3client),
			downloader.WithBucket(c.bucket),
			downloader.WithConcurrency(c.concurrency),
			downloader.WithMaxPartSize(c.maxPartSize),
		)
		if err != nil {
			return err
		}
	}

	if c.uclient == nil {
		c.uclient, err = uploader.New(
			uploader.WithErrGroup(c.eg),
			uploader.WithAPIClient(c.s3client),
			uploader.WithBucket(c.bucket),
			uploader.WithConcurrency(c.concurrency),
			uploader.WithMaxPartSize(c.maxPartSize),
		)
		if err != nil {
			return err
		}
	}
	return nil
}

// Close does nothing. Always returns nil.
func (c *client) Close() error {
	return nil
}

// Reader creates reader.Reader implementation and returns it.
// An error will be returned when an error occurs in downloader.Download.
func (c *client) Reader(ctx context.Context, key string) (rc io.ReadCloser, err error) {
	return c.dclient.Download(ctx, key)
}

// Writer creates writer.Writer implementation and returns it.
// An error will be returned when the uploader initialization fails or an error occurs in uploader.Open.
func (c *client) Writer(ctx context.Context, key string) (wc io.WriteCloser, err error) {
	err = c.uclient.Open(ctx, key)
	if err != nil {
		return nil, err
	}
	return c.uclient, nil
}
