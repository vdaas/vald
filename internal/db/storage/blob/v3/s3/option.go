package s3

import (
	"net/http"

	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/unit"
)

// Option represents the functional option for client.
type Option func(c *client) error

var defaultOptions = []Option{
	WithErrGroup(errgroup.Get()),
	WithMaxRetries(-1),
	WithForcePathStyle(false),
	WithUseAccelerate(false),
	WithUseARNRegion(false),
	WithUseDualStack(false),
	WithEnableSSL(true),
	WithEnableEndpointDiscovery(false),
	// WithEnableParamValidation(true),
	// WithEnable100Continue(true),
	// WithEnableContentMD5Validation(true),
	// WithEnableEndpointHostPrefix(true),
}

// WithErrGroup returns the option to set the eg.
func WithErrGroup(eg errgroup.Group) Option {
	return func(c *client) error {
		if eg == nil {
			return errors.NewErrInvalidOption("errgroup", eg)
		}
		c.eg = eg
		return nil
	}
}

// WithEndpoint returns the option to set the endpoint.
func WithEndpoint(ep string) Option {
	return func(c *client) error {
		if len(ep) == 0 {
			return errors.NewErrInvalidOption("endpoint", ep)
		}
		c.endpoint = ep
		return nil
	}
}

// WithBucket returns the option to set bucket.
func WithBucket(bucket string) Option {
	return func(c *client) error {
		if len(bucket) == 0 {
			return errors.NewErrInvalidOption("bucket", bucket)
		}
		c.bucket = bucket
		return nil
	}
}

// WithRegion returns the option to set the region.
func WithRegion(rg string) Option {
	return func(c *client) error {
		if len(rg) == 0 {
			return errors.NewErrInvalidOption("region", rg)
		}
		c.region = rg
		return nil
	}
}

// WithAccessKey returns the option to set the accessKey.
func WithAccessKey(ak string) Option {
	return func(c *client) error {
		if len(ak) == 0 {
			return errors.NewErrInvalidOption("accessKey", ak)
		}
		c.accessKey = ak
		return nil
	}
}

// WithSecretAccessKey returns the option to set the secretAccessKey.
func WithSecretAccessKey(sak string) Option {
	return func(c *client) error {
		if len(sak) == 0 {
			return errors.NewErrInvalidOption("secretAccessKey", sak)
		}
		c.secretAccessKey = sak
		return nil
	}
}

// WithToken returns the option to set the token.
func WithToken(tk string) Option {
	return func(c *client) error {
		if len(tk) == 0 {
			return errors.NewErrInvalidOption("token", tk)
		}
		c.token = tk
		return nil
	}
}

// WithMaxRetries returns the option to set the maxRetries.
func WithMaxRetries(r int) Option {
	return func(c *client) error {
		c.maxRetries = r
		return nil
	}
}

// WithForcePathStyle returns the option to set the forcePathStyle.
func WithForcePathStyle(enabled bool) Option {
	return func(c *client) error {
		c.forcePathStyle = enabled
		return nil
	}
}

// WithUseAccelerate returns the option to set the useAccelerate.
func WithUseAccelerate(enabled bool) Option {
	return func(c *client) error {
		c.useAccelerate = enabled
		return nil
	}
}

// WithUseARNRegion returns the option to set the useARNRegion.
func WithUseARNRegion(enabled bool) Option {
	return func(c *client) error {
		c.useARNRegion = enabled
		return nil
	}
}

// WithUseDualStack returns the option to set the useDualStack.
func WithUseDualStack(enabled bool) Option {
	return func(c *client) error {
		c.useDualStack = enabled
		return nil
	}
}

// WithEnableSSL returns the option to set the enableSSL.
func WithEnableSSL(enabled bool) Option {
	return func(c *client) error {
		c.enableSSL = enabled
		return nil
	}
}

// WithEnableEndpointDiscovery returns the option to set the enableEndpointDiscovery.
func WithEnableEndpointDiscovery(enabled bool) Option {
	return func(c *client) error {
		c.enableEndpointDiscovery = enabled
		return nil
	}
}

// WithHTTPClient returns the option to set the client.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *client) error {
		if hc == nil {
			return errors.NewErrInvalidOption("httpClient", hc)
		}
		c.client = hc
		return nil
	}
}

// WithMaxPartSize returns the option to set maxPartSize.
// The minimum allowed part size is 5MB, and if this value is set to zero,
// the DefaultUploadPartSize(DefaultDownloadPartSize) value will be used.
func WithMaxPartSize(size string) Option {
	return func(c *client) error {
		b, err := unit.ParseBytes(size)
		if err != nil {
			return err
		}

		if n := int64(b); n >= manager.DefaultUploadPartSize {
			c.maxPartSize = n
		}
		return nil
	}
}

func WithConcurrency(n int) Option {
	return func(c *client) error {
		if n >= manager.DefaultDownloadConcurrency {
			c.concurrency = n
		}
		return nil
	}
}
