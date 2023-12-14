//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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
	endpoint        string
	region          string
	accessKey       string
	secretAccessKey string
	token           string

	maxRetries                 int
	forcePathStyle             bool
	useAccelerate              bool
	useARNRegion               bool
	useDualStack               bool
	enableSSL                  bool
	enableParamValidation      bool
	enable100Continue          bool
	enableContentMD5Validation bool
	enableEndpointDiscovery    bool
	enableEndpointHostPrefix   bool

	client *http.Client
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
