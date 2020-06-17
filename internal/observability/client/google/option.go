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

// Package google provides a google API client options.
package google

import (
	"google.golang.org/api/option"
)

type Option = option.ClientOption

var (
	WithAPIKey             = option.WithAPIKey
	WithAudiences          = option.WithAudiences
	WithCredentialsFile    = option.WithCredentialsFile
	WithEndpoint           = option.WithEndpoint
	WithQuotaProject       = option.WithQuotaProject
	WithRequestReason      = option.WithRequestReason
	WithScopes             = option.WithScopes
	WithServiceAccountFile = option.WithServiceAccountFile
	WithUserAgent          = option.WithUserAgent

	// WithClientCertSource(s ClientCertSource) ClientOption
	// WithCredentials(creds *google.Credentials) ClientOption
	// WithGRPCConn(conn *grpc.ClientConn) ClientOption
	// WithGRPCConnectionPool(size int) ClientOption
	// WithGRPCDialOption(opt grpc.DialOption) ClientOption
	// WithHTTPClient(client *http.Client) ClientOption
	// WithTokenSource(s oauth2.TokenSource) ClientOption
)

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
