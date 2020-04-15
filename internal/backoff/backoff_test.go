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

package backoff

import (
	"context"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
)

func TestNew(t *testing.T) {
	type test struct {
		name      string
		opts      []Option
		want      *backoff
		checkFunc func(got, want *backoff) error
	}

	tests := []test{
		{
			name: "returns backoff instance",
			opts: []Option{
				WithBackOffFactor(0.5),
			},
			want: &backoff{
				initialDuration:  float64(10 * time.Millisecond),
				backoffTimeLimit: 5 * time.Minute,
				maxDuration:      float64(time.Hour),
				jitterLimit:      float64(time.Minute),
				backoffFactor:    1.1,
				maxRetryCount:    50,
				errLog:           true,
				durationLimit:    float64(time.Hour) / 1.1,
			},
			checkFunc: func(got *backoff, want *backoff) error {
				got.jittedInitialDuration, want.jittedInitialDuration = 1, 1
				if !reflect.DeepEqual(got, want) {
					return errors.Errorf("not equals. want: %v, got: %v", got, want)
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.opts...)
			if err := tt.checkFunc(got.(*backoff), tt.want); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestDo(t *testing.T) {
	type args struct {
		fn   func() (interface{}, error)
		opts []Option
	}

	type test struct {
		name      string
		args      args
		ctxFn     func() (context.Context, context.CancelFunc)
		checkFunc func(got, want error) error
		want      error
	}

	tests := []test{
		func() test {
			cnt := 0
			fn := func() (interface{}, error) {
				cnt++
				return nil, nil
			}

			return test{
				name: "returns response and nil when function return not nil",
				args: args{
					fn: fn,
					opts: []Option{
						WithDisableErrorLog(),
					},
				},
				ctxFn: func() (context.Context, context.CancelFunc) {
					return context.WithCancel(context.Background())
				},
				checkFunc: func(got, want error) error {
					if cnt != 1 {
						return errors.Errorf("error count is wrong, want: %v, got: %v", 2, cnt)
					}

					if !errors.Is(want, got) {
						return errors.Errorf("not equals. want: %v, got: %v", want, got)
					}

					return nil
				},
				want: nil,
			}
		}(),

		func() test {
			cnt := 0
			fn := func() (interface{}, error) {
				cnt++
				if cnt == 2 {
					return nil, nil
				}
				return nil, errors.Errorf("error (%d)", cnt)
			}

			return test{
				name: "returns response and nil when retried twice and did not return an error",
				args: args{
					fn: fn,
					opts: []Option{
						WithDisableErrorLog(),
						WithRetryCount(6),
					},
				},
				ctxFn: func() (context.Context, context.CancelFunc) {
					return context.WithCancel(context.Background())
				},
				checkFunc: func(got, want error) error {
					if cnt != 2 {
						return errors.Errorf("error count is wrong, want: %v, got: %v", 2, cnt)
					}

					if !errors.Is(want, got) {
						return errors.Errorf("not equals. want: %v, got: %v", want, got)
					}

					return nil
				},
				want: nil,
			}
		}(),

		func() test {
			cnt := 0
			fn := func() (interface{}, error) {
				cnt++
				return nil, errors.Errorf("error (%d)", cnt)
			}

			return test{
				name: "returns error when retrying the maximum number of times",
				args: args{
					fn: fn,
					opts: []Option{
						WithRetryCount(6),
					},
				},
				ctxFn: func() (context.Context, context.CancelFunc) {
					return context.WithCancel(context.Background())
				},
				checkFunc: func(got, want error) error {
					if cnt != 7 {
						return errors.Errorf("error count is wrong, want: %v, got: %v", 7, cnt)
					}

					if want.Error() != got.Error() {
						return errors.Errorf("not equals. want: %v, got: %v", want, got)
					}

					return nil
				},
				want: errors.New("error (7)"),
			}
		}(),

		func() test {
			ctx, cancel := context.WithCancel(context.Background())

			cnt := 0
			fn := func() (interface{}, error) {
				cnt++
				if cnt == 2 {
					cancel()
				}
				return nil, errors.Errorf("error (%d)", cnt)
			}

			return test{
				name: "return response and error when context canceled",
				args: args{
					fn: fn,
					opts: []Option{
						WithDisableErrorLog(),
						WithRetryCount(6),
					},
				},
				ctxFn: func() (context.Context, context.CancelFunc) {
					return ctx, cancel
				},
				checkFunc: func(got, want error) error {
					if cnt != 2 {
						return errors.Errorf("error count is wrong, want: %v, got: %v", 2, cnt)
					}

					if got.Error() != want.Error() {
						return errors.Errorf("not equals. want: %v, got: %v", want, got)
					}

					return nil
				},
				want: errors.Wrap(errors.New("error (2)"), context.Canceled.Error()),
			}
		}(),

		func() test {
			err := errors.New("error")
			fn := func() (interface{}, error) {
				return nil, err
			}

			return test{
				name: "return response and error when backoff timeout",
				args: args{
					fn: fn,
					opts: []Option{
						WithDisableErrorLog(),
						WithRetryCount(6),
						WithBackOffTimeLimit("0s"),
					},
				},
				ctxFn: func() (context.Context, context.CancelFunc) {
					return context.WithCancel(context.Background())
				},
				checkFunc: func(got, want error) error {
					if !errors.Is(got, want) {
						return errors.Errorf("not equals. want: %v, got: %v", want, got)
					}
					return nil
				},
				want: err,
			}
		}(),
	}

	log.Init()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := tt.ctxFn()
			defer cancel()

			_, err := New(tt.args.opts...).Do(ctx, tt.args.fn)
			if tt.want == nil && err != nil {
				t.Errorf("Do return err: %v", err)
			}

			if err := tt.checkFunc(err, tt.want); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestClose(t *testing.T) {
	type test struct {
		name string
		bo   *backoff
	}

	tests := []test{
		{
			name: "processing is successes",
			bo: &backoff{
				wg: sync.WaitGroup{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.bo.Close()
		})
	}
}
