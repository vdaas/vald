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
	stderrs "errors"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/log/logger"
	"github.com/vdaas/vald/internal/test/goleak"
)

func TestMain(m *testing.M) {
	log.Init(log.WithLoggerType(logger.NOP.String()))
	os.Exit(m.Run())
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
	defaultAfterFunc := func(args) {
		os.Args = []string{
			"test",
		}
	}
	tests := []test{
		{
			name: "returns nil when option is nil and version option is set",
			beforeFunc: func(args) {
				os.Args = []string{
					"test", "-version",
				}
			},
			want: want{
				err: nil,
			},
		},

		{
			name: "returns error when option is nil and params.Parse returns error",
			beforeFunc: func(args) {
				os.Args = []string{
					"test", "-team=set",
				}
			},
			want: want{
				err: errors.ErrArgumentParseFailed(stderrs.New("flag provided but not defined: -team")),
			},
		},

		{
			name: "returns error when option is not nil and r.loadConfig returns error",
			args: args{
				opts: []Option{
					WithConfigLoader(func(string) (interface{}, *config.GlobalConfig, error) {
						return nil, nil, errors.New("err")
					}),
				},
			},
			beforeFunc: func(args) {
				os.Args = []string{
					"test", "-c=./runner.go",
				}
			},
			want: want{
				err: errors.New("err"),
			},
		},

		{
			name: "returns error when option is not nil and ver.Check returns error",
			args: args{
				opts: []Option{
					WithVersion("v1.1.7", "v1.1.5", "v1.1.0"),
					WithConfigLoader(func(string) (interface{}, *config.GlobalConfig, error) {
						return nil, &config.GlobalConfig{
							Logging: &config.Logging{
								Logger: "glg",
								Level:  "info",
								Format: "json",
							},
							Version: "v1.1.7",
						}, nil
					}),
				},
			},
			beforeFunc: func(args) {
				os.Args = []string{
					"test", "-c=./runner.go",
				}
			},
			want: want{
				err: errors.ErrInvalidConfigVersion("1.1.7", ">= v1.1.0, <= v1.1.5"),
			},
		},

		{
			name: "returns error when option is not nil and r.initializeDaemon returns error",
			args: args{
				opts: []Option{
					WithVersion("v1.1.2", "v1.1.5", "v1.1.0"),
					WithConfigLoader(func(string) (interface{}, *config.GlobalConfig, error) {
						return nil, &config.GlobalConfig{
							Logging: &config.Logging{
								Logger: "glg",
								Level:  "info",
								Format: "json",
							},
							Version: "v1.1.2",
						}, nil
					}),
					WithDaemonInitializer(func(interface{}) (Runner, error) {
						return nil, errors.New("err")
					}),
				},
			},
			beforeFunc: func(args) {
				os.Args = []string{
					"test", "-c=./runner.go",
				}
			},
			want: want{
				err: errors.New("err"),
			},
		},

		{
			name: "returns nil when option is not nil and Run returns nil",
			args: args{
				opts: []Option{
					WithVersion("v1.1.2", "v1.1.5", "v1.1.0"),
					WithConfigLoader(func(string) (interface{}, *config.GlobalConfig, error) {
						return nil, &config.GlobalConfig{
							Logging: &config.Logging{
								Logger: "glg",
								Level:  "info",
								Format: "json",
							},
							Version: "v1.1.2",
						}, nil
					}),
					WithDaemonInitializer(func(interface{}) (Runner, error) {
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
						}, nil
					}),
				},
			},
			beforeFunc: func(args) {
				os.Args = []string{
					"test", "-c=./runner.go",
				}
				go func() {
					time.Sleep(2 * time.Second)
					syscall.Kill(syscall.Getpid(), syscall.SIGINT)
				}()
			},
			want: want{
				err: nil,
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			afterFunc := defaultAfterFunc
			if test.afterFunc != nil {
				afterFunc = test.afterFunc
			}
			defer afterFunc(test.args)

			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			err := Do(ctx, test.args.opts...)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
