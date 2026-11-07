// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
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
package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/vdaas/vald/pkg/gateway/lb/handler/rest"
)

// routeRecorder records which handler method served the request.
type routeRecorder struct {
	rest.Handler
	called string
}

func (r *routeRecorder) MultiRemove(http.ResponseWriter, *http.Request) (int, error) {
	r.called = "MultiRemove"
	return http.StatusOK, nil
}

func (r *routeRecorder) MultiSearch(http.ResponseWriter, *http.Request) (int, error) {
	r.called = "MultiSearch"
	return http.StatusOK, nil
}

func (r *routeRecorder) MultiSearchByID(http.ResponseWriter, *http.Request) (int, error) {
	r.called = "MultiSearchByID"
	return http.StatusOK, nil
}

func (r *routeRecorder) Remove(http.ResponseWriter, *http.Request) (int, error) {
	r.called = "Remove"
	return http.StatusOK, nil
}

func (r *routeRecorder) SearchByID(http.ResponseWriter, *http.Request) (int, error) {
	r.called = "SearchByID"
	return http.StatusOK, nil
}

func (r *routeRecorder) LinearSearch(http.ResponseWriter, *http.Request) (int, error) {
	r.called = "LinearSearch"
	return http.StatusOK, nil
}

func (r *routeRecorder) LinearSearchByID(http.ResponseWriter, *http.Request) (int, error) {
	r.called = "LinearSearchByID"
	return http.StatusOK, nil
}

func (r *routeRecorder) MultiLinearSearch(http.ResponseWriter, *http.Request) (int, error) {
	r.called = "MultiLinearSearch"
	return http.StatusOK, nil
}

func (r *routeRecorder) MultiLinearSearchByID(http.ResponseWriter, *http.Request) (int, error) {
	r.called = "MultiLinearSearchByID"
	return http.StatusOK, nil
}

func (r *routeRecorder) Exists(http.ResponseWriter, *http.Request) (int, error) {
	r.called = "Exists"
	return http.StatusOK, nil
}

func (r *routeRecorder) GetObject(http.ResponseWriter, *http.Request) (int, error) {
	r.called = "GetObject"
	return http.StatusOK, nil
}

func (r *routeRecorder) GetTimestamp(http.ResponseWriter, *http.Request) (int, error) {
	r.called = "GetTimestamp"
	return http.StatusOK, nil
}

func (r *routeRecorder) UpdateTimestamp(http.ResponseWriter, *http.Request) (int, error) {
	r.called = "UpdateTimestamp"
	return http.StatusOK, nil
}

func (r *routeRecorder) RemoveByTimestamp(http.ResponseWriter, *http.Request) (int, error) {
	r.called = "RemoveByTimestamp"
	return http.StatusOK, nil
}

func (r *routeRecorder) IndexInfo(http.ResponseWriter, *http.Request) (int, error) {
	r.called = "IndexInfo"
	return http.StatusOK, nil
}

func (r *routeRecorder) IndexDetail(http.ResponseWriter, *http.Request) (int, error) {
	r.called = "IndexDetail"
	return http.StatusOK, nil
}

func (r *routeRecorder) IndexStatistics(http.ResponseWriter, *http.Request) (int, error) {
	r.called = "IndexStatistics"
	return http.StatusOK, nil
}

func (r *routeRecorder) IndexStatisticsDetail(http.ResponseWriter, *http.Request) (int, error) {
	r.called = "IndexStatisticsDetail"
	return http.StatusOK, nil
}

func (r *routeRecorder) IndexProperty(http.ResponseWriter, *http.Request) (int, error) {
	r.called = "IndexProperty"
	return http.StatusOK, nil
}

func TestNew_staticRoutesNotShadowedByID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		path    string
		want    string
		methods []string
	}{
		{name: "multi search", methods: []string{http.MethodPost}, path: "/search/multi", want: "MultiSearch"},
		{name: "search by id", methods: []string{http.MethodGet}, path: "/search/uuid", want: "SearchByID"},
		{name: "multi search by id", methods: []string{http.MethodGet}, path: "/search/multi/uuid", want: "MultiSearchByID"},
		{name: "multi remove", methods: []string{http.MethodDelete, http.MethodPost}, path: "/delete/multi", want: "MultiRemove"},
		{name: "remove by id", methods: []string{http.MethodDelete}, path: "/delete/uuid", want: "Remove"},
		{name: "linear search", methods: []string{http.MethodPost}, path: "/linearsearch", want: "LinearSearch"},
		{name: "multi linear search", methods: []string{http.MethodPost}, path: "/linearsearch/multi", want: "MultiLinearSearch"},
		{name: "linear search by id", methods: []string{http.MethodGet}, path: "/linearsearch/uuid", want: "LinearSearchByID"},
		{name: "multi linear search by id", methods: []string{http.MethodGet}, path: "/linearsearch/multi/uuid", want: "MultiLinearSearchByID"},
		{name: "exists", methods: []string{http.MethodGet}, path: "/exists/uuid", want: "Exists"},
		{name: "get object", methods: []string{http.MethodGet}, path: "/object/uuid", want: "GetObject"},
		{name: "get timestamp", methods: []string{http.MethodGet}, path: "/object/meta/uuid", want: "GetTimestamp"},
		{name: "update timestamp", methods: []string{http.MethodPost}, path: "/update/timestamp", want: "UpdateTimestamp"},
		{name: "remove by timestamp", methods: []string{http.MethodPost}, path: "/remove/timestamp", want: "RemoveByTimestamp"},
		{name: "index info", methods: []string{http.MethodGet}, path: "/index/info", want: "IndexInfo"},
		{name: "index detail", methods: []string{http.MethodGet}, path: "/index/detail", want: "IndexDetail"},
		{name: "index statistics", methods: []string{http.MethodGet}, path: "/index/statistics", want: "IndexStatistics"},
		{name: "index statistics detail", methods: []string{http.MethodGet}, path: "/index/statistics/detail", want: "IndexStatisticsDetail"},
		{name: "index property", methods: []string{http.MethodGet}, path: "/index/property", want: "IndexProperty"},
	}
	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			for _, method := range test.methods {
				// embed a real handler so that every other method value is bound.
				rec := &routeRecorder{Handler: rest.New()}
				h := New(WithHandler(rec))
				req := httptest.NewRequestWithContext(tt.Context(), method, test.path, nil)
				h.ServeHTTP(httptest.NewRecorder(), req)
				if rec.called != test.want {
					tt.Errorf("%s %s served by %q, want %q", method, test.path, rec.called, test.want)
				}
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
//
// func TestNew(t *testing.T) {
// 	type args struct {
// 		opts []Option
// 	}
// 	type want struct {
// 		want http.Handler
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, http.Handler) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got http.Handler) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           opts:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           opts:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			got := New(test.args.opts...)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
