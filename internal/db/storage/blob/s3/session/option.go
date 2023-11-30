//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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
	"net/http"

	"github.com/vdaas/vald/internal/errors"
)

// Option represents the functional option for session.
type Option func(s *sess) error

var defaultOptions = []Option{
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

// WithEndpoint returns the option to set the endpoint.
func WithEndpoint(ep string) Option {
	return func(s *sess) error {
		if len(ep) == 0 {
			return errors.NewErrInvalidOption("endpoint", ep)
		}
		s.endpoint = ep
		return nil
	}
}

// WithRegion returns the option to set the region.
func WithRegion(rg string) Option {
	return func(s *sess) error {
		if len(rg) == 0 {
			return errors.NewErrInvalidOption("region", rg)
		}
		s.region = rg
		return nil
	}
}

// WithAccessKey returns the option to set the accessKey.
func WithAccessKey(ak string) Option {
	return func(s *sess) error {
		if len(ak) == 0 {
			return errors.NewErrInvalidOption("accessKey", ak)
		}
		s.accessKey = ak
		return nil
	}
}

// WithSecretAccessKey returns the option to set the secretAccessKey.
func WithSecretAccessKey(sak string) Option {
	return func(s *sess) error {
		if len(sak) == 0 {
			return errors.NewErrInvalidOption("secretAccessKey", sak)
		}
		s.secretAccessKey = sak
		return nil
	}
}

// WithToken returns the option to set the token.
func WithToken(tk string) Option {
	return func(s *sess) error {
		if len(tk) == 0 {
			return errors.NewErrInvalidOption("token", tk)
		}
		s.token = tk
		return nil
	}
}

// WithMaxRetries returns the option to set the maxRetries.
func WithMaxRetries(r int) Option {
	return func(s *sess) error {
		s.maxRetries = r
		return nil
	}
}

// WithForcePathStyle returns the option to set the forcePathStyle.
func WithForcePathStyle(enabled bool) Option {
	return func(s *sess) error {
		s.forcePathStyle = enabled
		return nil
	}
}

// WithUseAccelerate returns the option to set the useAccelerate.
func WithUseAccelerate(enabled bool) Option {
	return func(s *sess) error {
		s.useAccelerate = enabled
		return nil
	}
}

// WithUseARNRegion returns the option to set the useARNRegion.
func WithUseARNRegion(enabled bool) Option {
	return func(s *sess) error {
		s.useARNRegion = enabled
		return nil
	}
}

// WithUseDualStack returns the option to set the useDualStack.
func WithUseDualStack(enabled bool) Option {
	return func(s *sess) error {
		s.useDualStack = enabled
		return nil
	}
}

// WithEnableSSL returns the option to set the enableSSL.
func WithEnableSSL(enabled bool) Option {
	return func(s *sess) error {
		s.enableSSL = enabled
		return nil
	}
}

// WithEnableParamValidation returns the option to set the enableParamValidation.
func WithEnableParamValidation(enabled bool) Option {
	return func(s *sess) error {
		s.enableParamValidation = enabled
		return nil
	}
}

// WithEnable100Continue returns the option to set the enable100Continue.
func WithEnable100Continue(enabled bool) Option {
	return func(s *sess) error {
		s.enable100Continue = enabled
		return nil
	}
}

// WithEnableContentMD5Validation returns the option to set the enableContentMD5Validation.
func WithEnableContentMD5Validation(enabled bool) Option {
	return func(s *sess) error {
		s.enableContentMD5Validation = enabled
		return nil
	}
}

// WithEnableEndpointDiscovery returns the option to set the enableEndpointDiscovery.
func WithEnableEndpointDiscovery(enabled bool) Option {
	return func(s *sess) error {
		s.enableEndpointDiscovery = enabled
		return nil
	}
}

// WithEnableEndpointHostPrefix returns the option to set the enableEndpointHostPrefix.
func WithEnableEndpointHostPrefix(enabled bool) Option {
	return func(s *sess) error {
		s.enableEndpointHostPrefix = enabled
		return nil
	}
}

// WithHTTPClient returns the option to set the client.
func WithHTTPClient(client *http.Client) Option {
	return func(s *sess) error {
		if client == nil {
			return errors.NewErrInvalidOption("httpClient", client)
		}
		s.client = client
		return nil
	}
}
