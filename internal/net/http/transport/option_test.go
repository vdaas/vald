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
package transport

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/backoff"
	"github.com/vdaas/vald/internal/errors"
)

func TestWithRoundTripper(t *testing.T) {
	type test struct {
		name      string
		tr        http.RoundTripper
		checkFunc func(Option) error
	}

	tests := []test{
		func() test {
			tr := new(roundTripMock)

			return test{
				name: "set success",
				tr:   tr,
				checkFunc: func(opt Option) error {
					got := new(ert)
					opt(got)

					if got, want := got.transport, tr; !reflect.DeepEqual(got, want) {
						return errors.New("invalid param was set")
					}
					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithRoundTripper(tt.tr)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithBackoff(t *testing.T) {
	type test struct {
		name      string
		bo        backoff.Backoff
		checkFunc func(Option) error
	}

	tests := []test{
		func() test {
			bo := new(backoffMock)

			return test{
				name: "set success",
				bo:   bo,
				checkFunc: func(opt Option) error {
					got := new(ert)
					opt(got)

					if got, want := got.bo, bo; !reflect.DeepEqual(got, want) {
						return errors.New("invalid param was set")
					}
					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithBackoff(tt.bo)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}
