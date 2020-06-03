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

package session

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

type sess struct {
	endpoint        string
	region          string
	accessKey       string
	secretAccessKey string
	token           string
}

type Session interface {
	Session() (*session.Session, error)
}

func New(opts ...Option) Session {
	s := new(sess)
	for _, opt := range append(defaultOpts, opts...) {
		opt(s)
	}

	return s
}

func (s *sess) Session() (*session.Session, error) {
	cfg := aws.NewConfig()

	if s.endpoint != "" {
		cfg = cfg.WithEndpoint(s.endpoint)
	}

	if s.region != "" {
		cfg = cfg.WithRegion(s.region)
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

	return session.NewSession(cfg)
}
