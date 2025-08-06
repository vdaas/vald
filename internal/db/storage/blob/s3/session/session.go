//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
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

package session

import (
	"net/http"
	"reflect"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
)

type sess struct {
	// client is a http.Client.
	client                     *http.Client
	// endpoint is a s3 endpoint.
	endpoint                   string
	// region is a s3 region.
	region                     string
	// accessKey is a s3 access key.
	accessKey                  string
	// secretAccessKey is a s3 secret access key.
	secretAccessKey            string
	// token is a s3 token.
	token                      string
	// maxRetries is a max retries.
	maxRetries                 int
	// useARNRegion is a flag to use ARN region.
	useARNRegion               bool
	// useAccelerate is a flag to use accelerate.
	useAccelerate              bool
	// useDualStack is a flag to use dual stack.
	useDualStack               bool
	// enableSSL is a flag to enable SSL.
	enableSSL                  bool
	// enableParamValidation is a flag to enable param validation.
	enableParamValidation      bool
	// enable100Continue is a flag to enable 100 continue.
	enable100Continue          bool
	// enableContentMD5Validation is a flag to enable content MD5 validation.
	enableContentMD5Validation bool
	// enableEndpointDiscovery is a flag to enable endpoint discovery.
	enableEndpointDiscovery    bool
	// enableEndpointHostPrefix is a flag to enable endpoint host prefix.
	enableEndpointHostPrefix   bool
	// forcePathStyle is a flag to force path style.
	forcePathStyle             bool
}

// Session represents the interface to get AWS S3 session.
type Session interface {
	Session() (*session.Session, error)
}

// New returns the session implementation.
func New(opts ...Option) Session {
	s := new(sess)
	for _, opt := range append(defaultOptions, opts...) {
		if err := opt(s); err != nil {
			log.Warn(errors.ErrOptionFailed(err, reflect.ValueOf(opt)))
		}
	}

	return s
}

// Session returns the AWS S3 session or any error occurred.
func (s *sess) Session() (*session.Session, error) {
	cfg := aws.NewConfig().WithRegion(s.region)

	if s.endpoint != "" {
		cfg = cfg.WithEndpoint(s.endpoint)
	}

	if s.accessKey != "" && s.secretAccessKey != "" {
		creds := credentials.NewStaticCredentials(
			s.accessKey,
			s.secretAccessKey,
			s.token,
		)
		_, err := creds.Get()
		if err != nil {
			return nil, err
		}

		cfg = cfg.WithCredentials(creds)
	}

	if s.maxRetries != -1 {
		cfg = cfg.WithMaxRetries(s.maxRetries)
	}

	cfg = cfg.WithS3ForcePathStyle(s.forcePathStyle).
		WithS3UseAccelerate(s.useAccelerate).
		WithS3UseARNRegion(s.useARNRegion).
		WithUseDualStack(s.useDualStack).
		WithEndpointDiscovery(s.enableEndpointDiscovery)

	if !s.enableSSL {
		cfg = cfg.WithDisableSSL(true)
	}

	if !s.enableParamValidation {
		cfg = cfg.WithDisableParamValidation(true)
	}

	if !s.enable100Continue {
		cfg = cfg.WithS3Disable100Continue(true)
	}

	if !s.enableContentMD5Validation {
		cfg = cfg.WithS3DisableContentMD5Validation(true)
	}

	if !s.enableEndpointHostPrefix {
		cfg = cfg.WithDisableEndpointHostPrefix(true)
	}

	if s.client != nil {
		cfg = cfg.WithHTTPClient(s.client)
	}

	return session.NewSession(cfg)
}
