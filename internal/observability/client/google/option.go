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

// Package google provides a google API client options.
package google

import (
	"google.golang.org/api/option"
)

type Option = option.ClientOption

// WithClientCertSource(s ClientCertSource) ClientOption
// WithCredentials(creds *google.Credentials) ClientOption
// WithGRPCConn(conn *grpc.ClientConn) ClientOption
// WithGRPCConnectionPool(size int) ClientOption
// WithGRPCDialOption(opt grpc.DialOption) ClientOption
// WithHTTPClient(client *http.Client) ClientOption
// WithTokenSource(s oauth2.TokenSource) ClientOption

func WithAPIKey(apiKey string) Option {
	if apiKey == "" {
		return nil
	}
	return option.WithAPIKey(apiKey)
}

func WithAudiences(audiences ...string) Option {
	if len(audiences) == 0 {
		return nil
	}

	return option.WithAudiences(audiences...)
}

func WithCredentialsFile(path string) Option {
	if path == "" {
		return nil
	}

	return option.WithCredentialsFile(path)
}

func WithEndpoint(endpoint string) Option {
	if endpoint == "" {
		return nil
	}

	return option.WithEndpoint(endpoint)
}

func WithQuotaProject(qp string) Option {
	if qp == "" {
		return nil
	}

	return option.WithQuotaProject(qp)
}

func WithRequestReason(rr string) Option {
	if rr == "" {
		return nil
	}

	return option.WithRequestReason(rr)
}

func WithScopes(scopes ...string) Option {
	if len(scopes) == 0 {
		return nil
	}

	return option.WithScopes(scopes...)
}

func WithUserAgent(ua string) Option {
	if ua == "" {
		return nil
	}

	return option.WithUserAgent(ua)
}

func WithCredentialsJSON(json string) Option {
	if json != "" {
		return option.WithCredentialsJSON([]byte(json))
	}

	return nil
}

func WithTelemetry(enabled bool) Option {
	if !enabled {
		return option.WithTelemetryDisabled()
	}

	return nil
}

func WithAuthentication(enabled bool) Option {
	if !enabled {
		return option.WithoutAuthentication()
	}

	return nil
}
