//
// Copyright (C) 2019 kpango (Yusuke Kato)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

// Package rest provides REST API common logic & variable
package rest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/vdaas/vald/internal/errors"
)

type Func func(http.ResponseWriter, *http.Request) (code int, err error)

const (
	// ContentType represents a HTTP header name "Content-Type"
	ContentType = "Content-Type"

	// ApplicationJSON represents a HTTP content type "application/json"
	ApplicationJSON = "application/json"

	// ProblemJSON represents a HTTP content type "application/problem+json"
	ProblemJSON = "application/problem+json"

	// TextPlain represents a HTTP content type "text/plain"
	TextPlain = "text/plain"

	// CharsetUTF8 represents a UTF-8 charset for HTTP response "charset=UTF-8"
	CharsetUTF8 = "charset=UTF-8"
)

func HandlerToRestFunc(f http.HandlerFunc) Func {
	return func(w http.ResponseWriter, r *http.Request) (code int, _ error) {
		f(w, r)
		return http.StatusOK, nil
	}
}

func IndexHandler(values map[string]interface{}, w http.ResponseWriter, r *http.Request) (cote int, err error) {
	if r == nil {
		return http.StatusBadRequest, errors.ErrInvalidRequest
	}
	var body []byte
	if r.Body != nil {
		body, err = ioutil.ReadAll(r.Body)
		if err != nil {
			return http.StatusInternalServerError, errors.ErrInvalidRequest
		}
		r.Body.Close()
	}
	w.WriteHeader(http.StatusOK)
	body, err = json.MarshalIndent(struct {
		Host             string                 `json:"host"`
		URI              string                 `json:"uri"`
		URL              string                 `json:"url"`
		Method           string                 `json:"method"`
		Proto            string                 `json:"proto"`
		Header           http.Header            `json:"header"`
		TransferEncoding []string               `json:"transfer_encoding"`
		RemoteAddr       string                 `json:"remote_addr"`
		ContentLength    int64                  `json:"content_length"`
		Body             []byte                 `json:"body"`
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
	}, "", "\t")
	if err != nil {
		return http.StatusInternalServerError, err
	}
	_, err = w.Write(body)
	return http.StatusOK, nil
}
