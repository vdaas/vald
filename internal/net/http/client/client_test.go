//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
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

package client

import (
	"context"
	"crypto/tls"
	"net/http"
	"net/url"
	"reflect"
	"sync/atomic"
	"testing"

	"github.com/vdaas/vald/internal/backoff"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/log/logger"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/test/comparator"
	"github.com/vdaas/vald/internal/test/goleak"
)

var (
	// Goroutine leak is detected by `fastime`, but it should be ignored in the test because it is an external package.
	goleakIgnoreOptions = []goleak.Option{
		goleak.IgnoreTopFunction("github.com/kpango/fastime.(*fastime).StartTimerD.func1"),
	}

	transportComparator = []comparator.Option{
		comparator.AllowUnexported(transport{}),
		comparator.AllowUnexported(http.Transport{}),
		comparator.IgnoreFields(http.Transport{}, "idleLRU", "altProto", "TLSNextProto", "dialsInProgress"),
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
		// skipcq: VET-V0008
		comparator.Comparer(func(x, y sync.Mutex) bool {
			// skipcq: VET-V0008
			return reflect.DeepEqual(x, y)
		}),
		comparator.Comparer(func(x, y atomic.Value) bool {
			return reflect.DeepEqual(x.Load(), y.Load())
		}),
		// skipcq: VET-V0008
		comparator.Comparer(func(x, y sync.Once) bool {
			// skipcq: VET-V0008
			return reflect.DeepEqual(x, y)
		}),
		comparator.Comparer(func(x, y *tls.Config) bool {
			return reflect.DeepEqual(x, y)
		}),
		// skipcq: VET-V0008
		comparator.Comparer(func(x, y sync.WaitGroup) bool {
			// skipcq: VET-V0008
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
				want: func() (c *http.Client) {
					c, _ = New()
					return c
				}(),
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
				want: func() (c *http.Client) {
					c, _ = New()
					return c
				}(),
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

// NOT IMPLEMENTED BELOW
//
// func TestNewWithTransport(t *testing.T) {
// 	type args struct {
// 		rt   http.RoundTripper
// 		opts []Option
// 	}
// 	type want struct {
// 		want *http.Client
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, *http.Client, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got *http.Client, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
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
// 		           rt:nil,
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
// 		           rt:nil,
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
// 			got, err := NewWithTransport(test.args.rt, test.args.opts...)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
