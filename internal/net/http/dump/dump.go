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

// Package dump provides http request/response dump logic
package dump

import (
	"net/http"

	"github.com/vdaas/vald/internal/errors"
)

func Request(values, body map[string]interface{}, r *http.Request) (res interface{}, err error) {
	if r == nil {
		return nil, errors.ErrInvalidRequest
	}
	return struct {
		Host             string                 `json:"host"`
		URI              string                 `json:"uri"`
		URL              string                 `json:"url"`
		Method           string                 `json:"method"`
		Proto            string                 `json:"proto"`
		Header           http.Header            `json:"header"`
		TransferEncoding []string               `json:"transfer_encoding"`
		RemoteAddr       string                 `json:"remote_addr"`
		ContentLength    int64                  `json:"content_length"`
		Body             map[string]interface{} `json:"body"`
		Values           map[string]interface{} `json:"values"`
	}{
		Host:             r.Host,
		URI:              r.RequestURI,
		URL:              r.URL.String(),
		Method:           r.Method,
		Proto:            r.Proto,
		Header:           r.Header,
		TransferEncoding: r.TransferEncoding,
		RemoteAddr:       r.RemoteAddr,
		ContentLength:    r.ContentLength,
		Body:             body,
		Values:           values,
	}, nil
}
