//
// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
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

package router

import (
	"net/http"
	"testing"

	"github.com/gorilla/mux"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/internal/test/goleak"
	"github.com/vdaas/vald/pkg/tools/embedder/handler/rest"
)

func TestNew(t *testing.T) {
	t.Parallel()
	defer goleak.VerifyNone(t, goleak.IgnoreCurrent())

	h := rest.New()
	got := New(WithHandler(h), WithErrGroup(errgroup.Get()))
	gh, ok := got.(*mux.Router)
	if !ok {
		t.Fatal("type cast got failed")
	}

	routes := map[string]string{
		"Index":                "/",
		"Search":               "/search",
		"LinearSearch":         "/linearsearch",
		"Insert":               "/insert",
		"Insert With Metadata": "/insert/with-metadata",
		"Update":               "/update",
		"Update With Metadata": "/update/with-metadata",
		"Upsert":               "/upsert",
		"Upsert With Metadata": "/upsert/with-metadata",
		"Remove":               "/remove",
		"Remove With Metadata": "/remove/with-metadata",
		"Embedding":            "/embedding",
	}
	for name, pattern := range routes {
		route := gh.Get(name)
		if route == nil {
			t.Fatal(errors.Errorf("route not found: %s", name))
		}
		gotPattern, err := route.GetPathTemplate()
		if err != nil {
			t.Fatal(err)
		}
		if gotPattern != pattern {
			t.Errorf("pattern = %s, want %s", gotPattern, pattern)
		}
		req, err := http.NewRequest(http.MethodGet, pattern, http.NoBody)
		if err != nil {
			t.Fatal(err)
		}
		if name != "Index" {
			req.Method = http.MethodPost
		}
		match := new(mux.RouteMatch)
		if !route.Match(req, match) {
			t.Fatalf("route %s does not match %s %s", name, req.Method, pattern)
		}
	}
}
