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
package dump

import (
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"github.com/vdaas/vald/internal/errors"
)

func TestRequest(t *testing.T) {
	type args struct {
		values map[string]interface{}
		body   map[string]interface{}
		r      *http.Request
	}

	type test struct {
		name      string
		args      args
		checkFunc func(res interface{}, err error) error
	}

	tests := []test{
		{
			name: "returns object converted to structure",
			args: args{
				r: &http.Request{
					Host:       "hoge",
					RequestURI: "uri",
					URL: &url.URL{
						Scheme: "http",
					},
					Method: http.MethodGet,
					Proto:  "proto",
					Header: http.Header{},
					TransferEncoding: []string{
						"trans1",
					},
					RemoteAddr:    "0.0.0.0",
					ContentLength: 1234,
				},
				body: map[string]interface{}{
					"name": "vald",
				},
				values: map[string]interface{}{
					"version": "1.0.0",
				},
			},
			checkFunc: func(res interface{}, err error) error {
				if err != nil {
					return errors.Errorf("err is not nil. err: %v", err)
				}

				b, err := json.Marshal(res)
				if err != nil {
					return err
				}

				str := `{"host":"hoge","uri":"uri","url":"http:","method":"GET","proto":"proto","header":{},"transfer_encoding":["trans1"],"remote_addr":"0.0.0.0","content_length":1234,"body":{"name":"vald"},"values":{"version":"1.0.0"}}`
				if got, want := string(b), str; got != want {
					return errors.Errorf("response not equals. want: %v, got: %v", want, got)
				}

				return nil
			},
		},
		{
			name: "returns nil and error",
			checkFunc: func(res interface{}, err error) error {
				if got, want := err, errors.ErrInvalidRequest; !errors.Is(got, want) {
					return errors.Errorf("err not equals. want: %v, got: %v", want, got)
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := Request(tt.args.values, tt.args.body, tt.args.r)
			if err := tt.checkFunc(res, err); err != nil {
				t.Error(err)
			}
		})
	}
}
