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

package client

import (
	"context"
	"crypto/tls"
	"net/http"
	"net/url"
	"reflect"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/vdaas/vald/internal/backoff"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/log/logger"
	"github.com/vdaas/vald/internal/net"
	htr "github.com/vdaas/vald/internal/net/http/transport"
	"github.com/vdaas/vald/internal/test/comparator"
	"github.com/vdaas/vald/internal/test/goleak"
	"golang.org/x/net/http2"
)

var (
	// Goroutine leak is detected by `fastime`, but it should be ignored in the test because it is an external package.
	goleakIgnoreOptions = []goleak.Option{
		goleak.IgnoreTopFunction("github.com/kpango/fastime.(*fastime).StartTimerD.func1"),
	}

	transportComparator = []comparator.Option{
		comparator.AllowUnexported(transport{}),
		comparator.AllowUnexported(http.Transport{}),
		comparator.IgnoreFields(http.Transport{}, "idleLRU", "altProto", "TLSNextProto"),
		comparator.Exporter(func(t reflect.Type) bool {
			if t.Name() == "ert" || t.Name() == "backoff" {
				return true
			}
			return false
		}),

		comparator.Comparer(func(x, y backoff.Option) bool {
			return reflect.ValueOf(x).Pointer() == reflect.ValueOf(y).Pointer()
		}),
		comparator.Comparer(func(x, y func(*http.Request) (*url.URL, error)) bool {
			return reflect.ValueOf(x).Pointer() == reflect.ValueOf(y).Pointer()
		}),
		comparator.Comparer(func(x, y func(ctx context.Context, network, addr string) (net.Conn, error)) bool {
			return reflect.ValueOf(x).Pointer() == reflect.ValueOf(y).Pointer()
		}),
		comparator.Comparer(func(x, y sync.Mutex) bool {
			return reflect.DeepEqual(x, y)
		}),
		comparator.Comparer(func(x, y atomic.Value) bool {
			return reflect.DeepEqual(x.Load(), y.Load())
		}),
		comparator.Comparer(func(x, y sync.Once) bool {
			return reflect.DeepEqual(x, y)
		}),
		comparator.Comparer(func(x, y *tls.Config) bool {
			return reflect.DeepEqual(x, y)
		}),
		comparator.Comparer(func(x, y sync.WaitGroup) bool {
			return reflect.DeepEqual(x, y)
		}),
	}

	clientComparator = append(transportComparator,
		comparator.AllowUnexported(http.Client{}),
		comparator.FilterPath(func(p comparator.Path) bool {
			return p.String() == "Transport.bo.jittedInitialDuration"
		}, comparator.Ignore()),
	)
)

func TestMain(m *testing.M) {
	log.Init(log.WithLoggerType(logger.NOP.String()))
	goleak.VerifyTestMain(m)
}

func TestNew(t *testing.T) {
	t.Parallel()
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
		if diff := comparator.Diff(got, w.want, clientComparator...); diff != "" {
			return errors.New(diff)
		}
		return nil
	}
	tests := []test{
		{
			name: "initialize success with no option",
			args: args{
				opts: nil,
			},
			want: want{
				want: &http.Client{
					Transport: htr.NewExpBackoff(
						htr.WithRoundTripper(func() *http.Transport {
							t := new(http.Transport)
							t.Proxy = http.ProxyFromEnvironment
							_ = http2.ConfigureTransport(t)

							return t
						}()),
						htr.WithBackoff(
							backoff.New(),
						),
					),
				},
			},
		},
		{
			name: "fails and log invalid option error",
			args: args{
				opts: []Option{
					func(*transport) error {
						return errors.NewErrInvalidOption("dum", 1)
					},
				},
			},
			want: want{
				want: &http.Client{
					Transport: htr.NewExpBackoff(
						htr.WithRoundTripper(func() *http.Transport {
							t := new(http.Transport)
							t.Proxy = http.ProxyFromEnvironment
							_ = http2.ConfigureTransport(t)

							return t
						}()),
						htr.WithBackoff(
							backoff.New(),
						),
					),
				},
			},
		},
		{
			name: "fails with critical option error",
			args: args{
				opts: []Option{
					func(*transport) error {
						return errors.NewErrCriticalOption("dum", 1)
					},
				},
			},
			want: want{
				err: errors.NewErrCriticalOption("dum", 1),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
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

			got, err := New(test.args.opts...)
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
