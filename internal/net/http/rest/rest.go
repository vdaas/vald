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

// Package rest provides REST API common logic & variable
package rest

import (
	"net/http"
)

type Func func(http.ResponseWriter, *http.Request) (code int, err error)

const (
	// ContentType represents a HTTP header name "Content-Type".
	ContentType = "Content-Type"

	// ApplicationJSON represents a HTTP content type "application/json".
	ApplicationJSON = "application/json"

	// ProblemJSON represents a HTTP content type "application/problem+json".
	ProblemJSON = "application/problem+json"

	// TextPlain represents a HTTP content type "text/plain".
	TextPlain = "text/plain"

	// CharsetUTF8 represents a UTF-8 charset for HTTP response "charset=UTF-8".
	CharsetUTF8 = "charset=UTF-8"
)

func HandlerToRestFunc(f http.HandlerFunc) Func {
	return func(w http.ResponseWriter, r *http.Request) (code int, _ error) {
		f(w, r)
		return http.StatusOK, nil
	}
}
