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
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"gocloud.dev/blob/s3blob"
)

type URLOpener = s3blob.URLOpener

type sess struct {
	endpoint        string
	region          string
	accessKey       string
	secretAccessKey string
	token           string
}

type Session interface {
	URLOpener() (*URLOpener, error)
}

func NewSession(opts ...Option) Session {
	s := new(sess)
	for _, opt := range append(defaultOpts, opts...) {
		opt(s)
	}

	return s
}

func (s *sess) URLOpener() (*URLOpener, error) {
	session, err := session.NewSession(
		aws.NewConfig().WithEndpoint(
			s.endpoint,
		).WithRegion(
			s.region,
		).WithCredentials(
			credentials.NewStaticCredentials(
				s.accessKey,
				s.secretAccessKey,
				s.token,
			),
		),
	)

	if err != nil {
		return nil, err
	}

	return &URLOpener{
		ConfigProvider: session,
		Options: s3blob.Options{
			UseLegacyList: false,
		},
	}, nil
}
