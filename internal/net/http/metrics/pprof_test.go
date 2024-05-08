// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
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
package metrics

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewPProfHandler(t *testing.T) {
	handler := NewPProfHandler()
	server := httptest.NewServer(handler)
	defer server.Close()

	for _, route := range GetProfileRoutes() {
		t.Run(route.Name, func(t *testing.T) {
			resp, err := http.Get(server.URL + route.Pattern)
			if err != nil {
				t.Errorf("Failed to make GET request for %s: %v", route.Name, err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				t.Errorf("Expected status code 200 for %s, got %d", route.Name, resp.StatusCode)
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
