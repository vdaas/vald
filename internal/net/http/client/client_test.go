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

package client

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/backoff"
	"github.com/vdaas/vald/internal/errors"
	htr "github.com/vdaas/vald/internal/net/http/transport"
	"github.com/vdaas/vald/internal/test/comparator"
	"go.uber.org/goleak"
	"golang.org/x/net/http2"
)

var (
	// Goroutine leak is detected by `fastime`, but it should be ignored in the test because it is an external package.
	goleakIgnoreOptions = []goleak.Option{
		goleak.IgnoreTopFunction("github.com/kpango/fastime.(*Fastime).StartTimerD.func1"),
		goleak.IgnoreTopFunction("internal/poll.runtime_pollWait"),
	}
)

func TestNew(t *testing.T) {
	type args struct {
		opts []Option
	}
	type want struct {
		want *http.Client
		err  error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *http.Client, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got *http.Client, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}

		opts := []comparator.Option{
			comparator.AllowUnexported(http.Client{}),

			comparator.Comparer(func(x, y http.RoundTripper) bool {
				return reflect.DeepEqual(x, y)
				//return comparator.Diff(x, y) == ""
			}),
		}
		if diff := comparator.Diff(got, w.want, opts...); diff != "" {
			return errors.New(diff)
		}
		return nil
	}
	tests := []test{
		{
			name: "return default http client success",
			args: args{
				opts: nil,
			},
			want: want{
				want: &http.Client{
					Transport: func() http.RoundTripper {
						t := &http.Transport{
							Proxy:              http.ProxyFromEnvironment,
							DisableKeepAlives:  false,
							DisableCompression: false,
						}
						_ = http2.ConfigureTransport(t)

						return htr.NewExpBackoff(
							htr.WithRoundTripper(t),

							htr.WithBackoff(
								backoff.New(),
							),
						)
					}(),
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got, err := New(test.args.opts...)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}
