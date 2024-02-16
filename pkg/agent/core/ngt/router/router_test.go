//
// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
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

// Package router provides implementation of Go API for routing http Handler wrapped by rest.Func
package router

import (
	"net/http"
	"testing"

	"github.com/gorilla/mux"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/http/routing"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/internal/test/goleak"
	"github.com/vdaas/vald/pkg/agent/core/ngt/handler/rest"
)

func TestNew(t *testing.T) {
	t.Parallel()
	type args struct {
		opts []Option
	}
	type want struct {
		routes []routing.Route
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, http.Handler) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got http.Handler) error {
		gh, ok := got.(*mux.Router)
		if !ok {
			return errors.New("type cast got failed")
		}
		for _, r := range w.routes {
			gotR := gh.Get(r.Name)
			if gotR == nil {
				return errors.Errorf("route not found: %s", r.Name)
			}

			if gotR.GetHandler() == nil {
				return errors.Errorf("handler not found: %s", r.Name)
			}

			gotP, err := gotR.GetPathRegexp()
			if err != nil {
				return err
			}
			if gotP == "" {
				return errors.Errorf("pattern is empty: %s", r.Name)
			}
		}
		return nil
	}
	tests := []test{
		func() test {
			h := rest.New()
			eg := errgroup.Get()

			return test{
				name: "return handlers",
				args: args{
					opts: []Option{
						WithHandler(h),
						WithErrGroup(eg),
					},
				},
				want: want{
					routes: []routing.Route{
						{
							Name: "Index",
							Methods: []string{
								http.MethodGet,
							},
							Pattern:     "/",
							HandlerFunc: h.Index,
						},
						{
							Name: "Search",
							Methods: []string{
								http.MethodPost,
							},
							Pattern:     "/search",
							HandlerFunc: h.Search,
						},
						{
							Name: "Search By ID",
							Methods: []string{
								http.MethodPost,
							},
							Pattern:     "/id/search",
							HandlerFunc: h.SearchByID,
						},
						{
							Name: "Insert",
							Methods: []string{
								http.MethodPost,
							},
							Pattern:     "/insert",
							HandlerFunc: h.Insert,
						},
						{
							Name: "Multiple Insert",
							Methods: []string{
								http.MethodPost,
							},
							Pattern:     "/insert/multi",
							HandlerFunc: h.MultiInsert,
						},
						{
							Name: "Update",
							Methods: []string{
								http.MethodPost,
								http.MethodPatch,
								http.MethodPut,
							},
							Pattern:     "/update",
							HandlerFunc: h.Update,
						},
						{
							Name: "Multiple Update",
							Methods: []string{
								http.MethodPost,
								http.MethodPatch,
								http.MethodPut,
							},
							Pattern:     "/update/multi",
							HandlerFunc: h.MultiUpdate,
						},
						{
							Name: "Remove",
							Methods: []string{
								http.MethodDelete,
							},
							Pattern:     "/delete",
							HandlerFunc: h.Remove,
						},
						{
							Name: "Multiple Remove",
							Methods: []string{
								http.MethodDelete,
								http.MethodPost,
							},
							Pattern:     "/delete/multi",
							HandlerFunc: h.MultiRemove,
						},
						{
							Name: "Create Index",
							Methods: []string{
								http.MethodPost,
							},
							Pattern:     "/index/create",
							HandlerFunc: h.CreateIndex,
						},
						{
							Name: "Save Index",
							Methods: []string{
								http.MethodGet,
							},
							Pattern:     "/index/save",
							HandlerFunc: h.SaveIndex,
						},
						{
							Name: "GetObject",
							Methods: []string{
								http.MethodGet,
							},
							Pattern:     "/object/{id}",
							HandlerFunc: h.GetObject,
						},
					},
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			got := New(test.args.opts...)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
