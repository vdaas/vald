// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package urlopener

import (
	"net/http"
)

type Option func(*urlOpener) error

func WithCredentialsFile(path string) Option {
	return func(uo *urlOpener) error {
		if len(path) != 0 {
			uo.credentialsFilePath = path
		}
		return nil
	}
}

func WithCredentialsJSON(str string) Option {
	return func(uo *urlOpener) error {
		if len(str) != 0 {
			uo.credentialsJSON = str
		}
		return nil
	}
}

func WithHTTPClient(c *http.Client) Option {
	return func(uo *urlOpener) error {
		if c != nil {
			uo.client = c
		}
		return nil
	}
}
