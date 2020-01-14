//
// Copyright (C) 2019-2020 Vdaas.org Vald team ( kpango, rinx, kmrmt )
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
package rest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/vdaas/vald/internal/errors"
)

func TestHandlerToRestFunc(t *testing.T) {
	type test struct {
		name      string
		hfn       http.HandlerFunc
		checkFunc func(Func) error
	}

	tests := []test{
		func() test {
			cnt := 0

			hfn := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				cnt++
			})

			return test{
				name: "returns 200 status code",
				hfn:  hfn,
				checkFunc: func(fn Func) error {
					code, err := fn(httptest.NewRecorder(), new(http.Request))
					if err != nil {
						return errors.Errorf("err is not nil. err: %v", err)
					}

					if code != http.StatusOK {
						return errors.Errorf("status code is wrong. want: %v, got: %v", http.StatusOK, code)
					}

					if cnt != 1 {
						return errors.Errorf("called count is wrong. want: %v, got: %v", 1, cnt)
					}

					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fn := HandlerToRestFunc(tt.hfn)
			if err := tt.checkFunc(fn); err != nil {
				t.Error(err)
			}
		})
	}
}
