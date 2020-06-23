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

import "net/http"

type Option func(s *sess)

var (
	defaultOpts = []Option{
		WithMaxRetries(-1),
		WithForcePathStyle(false),
		WithUseAccelerate(false),
		WithUseARNRegion(false),
		WithUseDualStack(false),
		WithEnableSSL(true),
		WithEnableParamValidation(true),
		WithEnable100Continue(true),
		WithEnableContentMD5Validation(true),
		WithEnableEndpointDiscovery(false),
		WithEnableEndpointHostPrefix(true),
	}
)

func WithEndpoint(ep string) Option {
	return func(s *sess) {
		s.endpoint = ep
	}
}

func WithRegion(rg string) Option {
	return func(s *sess) {
		s.region = rg
	}
}

func WithAccessKey(ak string) Option {
	return func(s *sess) {
		s.accessKey = ak
	}
}

func WithSecretAccessKey(sak string) Option {
	return func(s *sess) {
		s.secretAccessKey = sak
	}
}

func WithToken(tk string) Option {
	return func(s *sess) {
		s.token = tk
	}
}

func WithMaxRetries(r int) Option {
	return func(s *sess) {
		s.maxRetries = r
	}
}

func WithForcePathStyle(enabled bool) Option {
	return func(s *sess) {
		s.forcePathStyle = enabled
	}
}

func WithUseAccelerate(enabled bool) Option {
	return func(s *sess) {
		s.useAccelerate = enabled
	}
}

func WithUseARNRegion(enabled bool) Option {
	return func(s *sess) {
		s.useARNRegion = enabled
	}
}

func WithUseDualStack(enabled bool) Option {
	return func(s *sess) {
		s.useDualStack = enabled
	}
}

func WithEnableSSL(enabled bool) Option {
	return func(s *sess) {
		s.enableSSL = enabled
	}
}

func WithEnableParamValidation(enabled bool) Option {
	return func(s *sess) {
		s.enableParamValidation = enabled
	}
}

func WithEnable100Continue(enabled bool) Option {
	return func(s *sess) {
		s.enable100Continue = enabled
	}
}

func WithEnableContentMD5Validation(enabled bool) Option {
	return func(s *sess) {
		s.enableContentMD5Validation = enabled
	}
}

func WithEnableEndpointDiscovery(enabled bool) Option {
	return func(s *sess) {
		s.enableEndpointDiscovery = enabled
	}
}

func WithEnableEndpointHostPrefix(enabled bool) Option {
	return func(s *sess) {
		s.enableEndpointHostPrefix = enabled
	}
}

func WithHTTPClient(client *http.Client) Option {
	return func(s *sess) {
		if client != nil {
			s.client = client
		}
	}
}
