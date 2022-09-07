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

package runner

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/log/logger"
	"github.com/vdaas/vald/internal/test/goleak"
)

func TestMain(m *testing.M) {
	log.Init(log.WithLoggerType(logger.NOP.String()))
	code := m.Run()
	os.Exit(code)
}

func TestRun(t *testing.T) {
	type args struct {
		ctx  context.Context
		run  Runner
		name string
	}
	type want struct {
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		return nil
	}
	tests := []test{
		func() test {
			ctx, cancel := context.WithCancel(context.Background())
			return test{
				name: "returns nil when internal functionally occurs no error",
				args: args{
					ctx: ctx,
					run: func() Runner {
						return &runnerMock{
							PreStartFunc: func(ctx context.Context) error {
								return nil
							},
							StartFunc: func(ctx context.Context) (<-chan error, error) {
								return make(chan error, 1), nil
							},
							PreStopFunc: func(ctx context.Context) error {
								return nil
							},
							StopFunc: func(ctx context.Context) error {
								return nil
							},
							PostStopFunc: func(ctx context.Context) error {
								return nil
							},
						}
					}(),
				},
				beforeFunc: func(args) {
					cancel()
				},
				want: want{
					err: nil,
				},
			}
		}(),

		func() test {
			ctx, cancel := context.WithCancel(context.Background())
			return test{
				name: "returns error when run.PreStop and run.Stop and run.PostStop returns error",
				args: args{
					ctx:  ctx,
					name: "vald",
					run: func() Runner {
						return &runnerMock{
							PreStartFunc: func(ctx context.Context) error {
								return nil
							},
							StartFunc: func(ctx context.Context) (<-chan error, error) {
								return make(chan error, 1), nil
							},
							PreStopFunc: func(ctx context.Context) error {
								return errors.New("err1")
							},
							StopFunc: func(ctx context.Context) error {
								return errors.New("err2")
							},
							PostStopFunc: func(ctx context.Context) error {
								return errors.New("err3")
							},
						}
					}(),
				},
				beforeFunc: func(args) {
					go func() {
						time.Sleep(2 * time.Second)
						cancel()
					}()
				},
				want: want{
					err: func() (err error) {
						details := []struct {
							err error
							cnt int
						}{
							{
								err: errors.New("err1"),
								cnt: 1,
							},
							{
								err: errors.New("err2"),
								cnt: 1,
							},
							{
								err: errors.New("err3"),
								cnt: 1,
							},
						}

						for _, det := range details {
							err = errors.Wrapf(err, "error:\t%s\tcount:\t%d", det.err.Error(), det.cnt)
						}

						return errors.ErrDaemonStopFailed(err)
					}(),
				},
			}
		}(),

		func() test {
			ctx, cancel := context.WithCancel(context.Background())
			return test{
				name: "returns error when channel of run.StartFunc contains error",
				args: args{
					ctx:  ctx,
					name: "vald",
					run: func() Runner {
						return &runnerMock{
							PreStartFunc: func(ctx context.Context) error {
								return nil
							},
							StartFunc: func(ctx context.Context) (<-chan error, error) {
								ch := make(chan error, 3)
								ch <- errors.New("err1")
								ch <- errors.New("err2")
								ch <- errors.New("err1")
								return ch, nil
							},
							PreStopFunc: func(ctx context.Context) error {
								return nil
							},
							StopFunc: func(ctx context.Context) error {
								return nil
							},
							PostStopFunc: func(ctx context.Context) error {
								return nil
							},
						}
					}(),
				},
				beforeFunc: func(args) {
					go func() {
						time.Sleep(2 * time.Second)
						cancel()
					}()
				},
				want: want{
					err: func() (err error) {
						details := []struct {
							err error
							cnt int
						}{
							{
								err: errors.New("err1"),
								cnt: 2,
							},
							{
								err: errors.New("err2"),
								cnt: 1,
							},
						}

						for _, detail := range details {
							err = errors.Wrapf(err, "error:\t%s\tcount:\t%d", detail.err.Error(), detail.cnt)
						}

						return errors.ErrDaemonStopFailed(err)
					}(),
				},
			}
		}(),

		{
			name: "returns error when run.PreStart returns error",
			args: args{
				ctx: context.Background(),
				run: func() Runner {
					return &runnerMock{
						PreStartFunc: func(context.Context) error {
							return errors.New("err")
						},
					}
				}(),
				name: "vald",
			},
			want: want{
				err: errors.New("err"),
			},
		},

		{
			name: "returns error when run.Start returns error",
			args: args{
				ctx: context.Background(),
				run: func() Runner {
					return &runnerMock{
						PreStartFunc: func(context.Context) error {
							return nil
						},
						StartFunc: func(context.Context) (<-chan error, error) {
							return nil, errors.New("err")
						},
					}
				}(),
				name: "vald",
			},
			want: want{
				err: errors.ErrDaemonStartFailed(errors.New("err")),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
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

			err := Run(test.args.ctx, test.args.run, test.args.name)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestDo(t *testing.T) {
	type args struct {
		ctx  context.Context
		opts []Option
	}
	type want struct {
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           ctx: nil,
		           opts: nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           ctx: nil,
		           opts: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
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

			err := Do(test.args.ctx, test.args.opts...)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
