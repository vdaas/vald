//
// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
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

// Package router provides implementation of Go API for routing http Handler wrapped by rest.Func
package router

import (
	"net/http"
	"testing"

	"github.com/gorilla/mux"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/http/routing"
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
							"Index",
							[]string{
								http.MethodGet,
							},
							"/",
							h.Index,
						},
						{
							"Search",
							[]string{
								http.MethodPost,
							},
							"/search",
							h.Search,
						},
						{
							"Search By ID",
							[]string{
								http.MethodPost,
							},
							"/id/search",
							h.SearchByID,
						},
						{
							"Insert",
							[]string{
								http.MethodPost,
							},
							"/insert",
							h.Insert,
						},
						{
							"Multiple Insert",
							[]string{
								http.MethodPost,
							},
							"/insert/multi",
							h.MultiInsert,
						},
						{
							"Update",
							[]string{
								http.MethodPost,
								http.MethodPatch,
								http.MethodPut,
							},
							"/update",
							h.Update,
						},
						{
							"Multiple Update",
							[]string{
								http.MethodPost,
								http.MethodPatch,
								http.MethodPut,
							},
							"/update/multi",
							h.MultiUpdate,
						},
						{
							"Remove",
							[]string{
								http.MethodDelete,
							},
							"/delete",
							h.Remove,
						},
						{
							"Multiple Remove",
							[]string{
								http.MethodDelete,
								http.MethodPost,
							},
							"/delete/multi",
							h.MultiRemove,
						},
						{
							"Create Index",
							[]string{
								http.MethodPost,
							},
							"/index/create",
							h.CreateIndex,
						},
						{
							"Save Index",
							[]string{
								http.MethodGet,
							},
							"/index/save",
							h.SaveIndex,
						},
						{
							"GetObject",
							[]string{
								http.MethodGet,
							},
							"/object/{id}",
							h.GetObject,
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
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := New(test.args.opts...)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
