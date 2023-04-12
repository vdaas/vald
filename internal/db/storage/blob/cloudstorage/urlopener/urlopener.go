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
package urlopener

import (
	"context"
	"net/http"
	"reflect"

	"cloud.google.com/go/storage"
	"github.com/vdaas/vald/internal/conv"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/file"
	"gocloud.dev/blob/gcsblob"
	"gocloud.dev/gcp"
	"golang.org/x/oauth2/google"
)

// URLOpener returns gcsblob.URLOpener and error.
type URLOpener interface {
	URLOpener(context.Context) (*gcsblob.URLOpener, error)
}

type urlOpener struct {
	credentialsFilePath string
	credentialsJSON     string

	client *http.Client
}

// New returns URLOpener implementation.
func New(opts ...Option) (URLOpener, error) {
	uo := new(urlOpener)
	for _, opt := range opts {
		if err := opt(uo); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}

	if uo.client == nil {
		return nil, errors.NewErrInvalidOption("client", uo.client)
	}

	return uo, nil
}

func (uo *urlOpener) URLOpener(ctx context.Context) (guo *gcsblob.URLOpener, err error) {
	var creds *google.Credentials
	scope := storage.ScopeReadWrite

	switch {
	case len(uo.credentialsFilePath) != 0:
		data, err := file.ReadFile(uo.credentialsFilePath)
		if err != nil || data == nil {
			return nil, err
		}
		creds, err = google.CredentialsFromJSON(ctx, data, scope)
		if err != nil {
			return nil, err
		}
	case len(uo.credentialsJSON) != 0:
		data := conv.Atob(uo.credentialsJSON)
		creds, err = google.CredentialsFromJSON(ctx, data, scope)
		if err != nil {
			return nil, err
		}
	default:
		creds, err = google.FindDefaultCredentials(ctx, scope)
		if err != nil {
			return nil, err
		}
	}

	cfg, err := google.JWTConfigFromJSON(creds.JSON, scope)
	if err != nil {
		return nil, err
	}

	client, err := gcp.NewHTTPClient(
		uo.client.Transport,
		gcp.CredentialsTokenSource(creds),
	)
	if err != nil {
		return nil, err
	}

	return &gcsblob.URLOpener{
		Client: client,
		Options: gcsblob.Options{
			GoogleAccessID: cfg.Email,
			PrivateKey:     cfg.PrivateKey,
		},
	}, nil
}
