//
// Copyright (C) 2019 Vdaas.org Vald team ( kpango, kmrmt, rinx )
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

	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
)

func TestDo(t *testing.T) {
	type args struct {
		ctx  context.Context
		opts []Option
	}

	type test struct {
		name     string
		testArgs []string
		args     args
		wantErr  bool
	}

	tests := []test{
		func() test {
			loadConfig := func(string) (interface{}, *config.GlobalConfig, error) {
				return nil, &config.GlobalConfig{
					Version: "v0.0.1",
					Logging: &config.Logging{},
				}, nil
			}

			initializeDeamon := func(interface{}) (Runner, error) {
				return &runnerMock{
					PreStartFunc: func(ctx context.Context) error {
						return nil
					},
					StartFunc: func(ctx context.Context) (<-chan error, error) {
						return make(chan error), nil
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
			}

			ctx, cancel := context.WithCancel(context.Background())
			cancel()

			return test{
				name: "returns nil when no error occurs internally",
				testArgs: []string{
					"dummyCommand",
				},
				args: args{
					opts: []Option{
						WithDaemonInitializer(initializeDeamon),
						WithConfigLoader(loadConfig),
						WithVersion("v0.0.0", "v0.0.10", "v0.0.1"),
					},
					ctx: ctx,
				},
			}
		}(),

		{
			name: "returns error when parsing of params fails",
			testArgs: []string{
				"dummyCommand",
				"-test.Args=true",
			},
			wantErr: true,
		},

		{
			name: "returns nil when help option is set",
			testArgs: []string{
				"dummyCommand",
				"--help",
			},
		},

		func() test {
			loadConfig := func(string) (interface{}, *config.GlobalConfig, error) {
				return nil, &config.GlobalConfig{
					Version: "v0.0.1",
					Logging: new(config.Logging),
				}, nil
			}

			initializeDeamon := func(interface{}) (Runner, error) {
				return nil, errors.New("fail")
			}

			return test{
				name: "returns error when initialization of the daemon fails",
				testArgs: []string{
					"dummyCommand",
				},
				args: args{
					opts: []Option{
						WithDaemonInitializer(initializeDeamon),
						WithConfigLoader(loadConfig),
						WithVersion("v0.0.0", "v0.0.10", "v0.0.1"),
					},
				},
				wantErr: true,
			}
		}(),

		func() test {
			loadConfig := func(string) (interface{}, *config.GlobalConfig, error) {
				return nil, new(config.GlobalConfig), nil
			}

			return test{
				name: "returns nil when version options is set",
				testArgs: []string{
					"dummyCommand",
					"--version",
				},
				args: args{
					opts: []Option{
						WithConfigLoader(loadConfig),
					},
				},
			}
		}(),

		func() test {
			loadConfig := func(string) (interface{}, *config.GlobalConfig, error) {
				return nil, &config.GlobalConfig{
					Version: "v0.1.0",
					Logging: new(config.Logging),
				}, nil
			}

			initializeDeamon := func(interface{}) (Runner, error) {
				return nil, errors.New("fail")
			}

			return test{
				name: "returns error when version check fails",
				testArgs: []string{
					"dummyCommand",
				},
				args: args{
					opts: []Option{
						WithDaemonInitializer(initializeDeamon),
						WithConfigLoader(loadConfig),
						WithVersion("v0.0.0", "v0.0.10", "v0.0.1"),
					},
				},
				wantErr: true,
			}
		}(),

		func() test {
			loadConfig := func(string) (interface{}, *config.GlobalConfig, error) {
				return nil, nil, errors.New("fail")
			}

			return test{
				name: "returns error when loading the config fails",
				testArgs: []string{
					"dummyCommand",
				},
				args: args{
					opts: []Option{
						WithConfigLoader(loadConfig),
					},
				},
				wantErr: true,
			}
		}(),

		func() test {
			loadConfig := func(string) (interface{}, *config.GlobalConfig, error) {
				return nil, &config.GlobalConfig{
					Version: "v0.1.0",
				}, errors.New("fail")
			}

			return test{
				name: "returns error incase of invalid version",
				testArgs: []string{
					"dummyCommand",
				},
				args: args{
					opts: []Option{
						WithConfigLoader(loadConfig),
						WithVersion("v0.0.0", "v0.0.10", "v0.0.1"),
					},
				},
				wantErr: true,
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = tt.testArgs
			if err := Do(tt.args.ctx, tt.args.opts...); (err != nil) != tt.wantErr {
				t.Errorf("not equals. err: %v", err)
			}
		})
	}
}

func TestRun(t *testing.T) {
	type args struct {
		ctx  context.Context
		run  Runner
		name string
	}

	type test struct {
		name      string
		args      args
		checkFunc func(args) error
	}

	tests := []test{
		func() test {
			run := &runnerMock{
				PreStartFunc: func(ctx context.Context) error {
					return nil
				},
				StartFunc: func(ctx context.Context) (<-chan error, error) {
					ch := make(chan error, 1)
					defer close(ch)
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

			ctx, cancel := context.WithCancel(context.Background())

			return test{
				name: "returns nil when no error occurs internally",
				args: args{
					ctx: ctx,
					run: run,
				},
				checkFunc: func(args args) error {
					cancel()
					got := Run(args.ctx, args.run, args.name)
					if got != nil {
						return errors.Errorf("err is not nil. err: %v", got)
					}
					return nil
				},
			}
		}(),

		func() test {
			err := errors.New("fail")
			run := &runnerMock{
				PreStartFunc: func(ctx context.Context) error {
					return err
				},
			}

			return test{
				name: "returns error when prestart function returns error",
				args: args{
					ctx: context.Background(),
					run: run,
				},
				checkFunc: func(args args) error {
					got := Run(args.ctx, args.run, args.name)
					if !errors.Is(got, err) {
						return errors.Errorf("err not equals. want: %v, but got: %v", err, got)
					}
					return nil
				},
			}
		}(),

		func() test {
			err := errors.New("fail")
			run := &runnerMock{
				PreStartFunc: func(ctx context.Context) error {
					return nil
				},
				StartFunc: func(ctx context.Context) (<-chan error, error) {
					return nil, err
				},
			}

			return test{
				name: "returns error when start function returns error",
				args: args{
					ctx: context.Background(),
					run: run,
				},
				checkFunc: func(args args) error {
					got := Run(args.ctx, args.run, args.name)
					if !errors.Is(got, err) {
						return errors.Errorf("err not equals. want: %v, but got: %v", err, got)
					}
					return nil
				},
			}
		}(),

		func() test {
			name := "vald"

			ctx, cancel := context.WithCancel(context.Background())

			var (
				startErr          = errors.New("start error")
				preStopErr  error = nil
				stopErr           = errors.New("stop error")
				postStopErr error = nil
			)

			run := &runnerMock{
				PreStartFunc: func(ctx context.Context) error {
					return nil
				},
				StartFunc: func(ctx context.Context) (<-chan error, error) {
					ch := make(chan error, 1)
					go func() {
						defer close(ch)
						ch <- startErr
						cancel()
					}()
					return ch, nil
				},
				PreStopFunc: func(ctx context.Context) error {
					return preStopErr
				},
				StopFunc: func(ctx context.Context) error {
					return stopErr
				},
				PostStopFunc: func(ctx context.Context) error {
					return postStopErr
				},
			}

			return test{
				name: "returns nil when start and stop function returns error",
				args: args{
					ctx:  ctx,
					run:  run,
					name: name,
				},
				checkFunc: func(args args) error {
					var want error
					for _, err := range []error{
						errors.ErrStartFunc(name, startErr),
						errors.ErrStopFunc(name, stopErr),
					} {
						want = errors.Wrapf(want, "error:\t%s\tcount:\t%d", err.Error(), 0)
					}
					want = errors.ErrDaemonStopFailed(want)

					got := Run(args.ctx, args.run, args.name)
					if got == nil {
						return errors.New("err is nil")
					} else if got.Error() != want.Error() {
						return errors.Errorf("not equals. want: %v, but got: %v", want, got)
					}
					return nil
				},
			}
		}(),

		func() test {
			name := "vald"

			ctx, cancel := context.WithCancel(context.Background())

			var (
				preStopErr  = errors.New("prestop error")
				stopErr     = errors.New("stop error")
				postStopErr = errors.New("poststop error")
				startErr    = errors.New("start error")
			)

			run := &runnerMock{
				PreStartFunc: func(ctx context.Context) error {
					return nil
				},
				StartFunc: func(ctx context.Context) (<-chan error, error) {
					ch := make(chan error, 1)
					go func() {
						defer close(ch)
						ch <- startErr
						cancel()
					}()
					return ch, nil
				},
				PreStopFunc: func(ctx context.Context) error {
					return preStopErr
				},
				StopFunc: func(ctx context.Context) error {
					return stopErr
				},
				PostStopFunc: func(ctx context.Context) error {
					return postStopErr
				},
			}

			return test{
				name: "returns nil when all(start, prestop, stop, poststop) function returns error",
				args: args{
					ctx:  ctx,
					run:  run,
					name: name,
				},
				checkFunc: func(args args) error {
					var want error
					for _, err := range []error{
						errors.ErrStartFunc(name, startErr),
						errors.ErrPreStopFunc(name, preStopErr),
						errors.ErrStopFunc(name, stopErr),
						errors.ErrPostStopFunc(name, postStopErr),
					} {
						want = errors.Wrapf(want, "error:\t%s\tcount:\t%d", err.Error(), 0)
					}
					want = errors.ErrDaemonStopFailed(want)

					got := Run(args.ctx, args.run, args.name)
					if got == nil {
						return errors.New("err is nil")
					} else if got.Error() != want.Error() {
						return errors.Errorf("not equals. want: %v, but got: %v", want, got)
					}
					return nil
				},
			}
		}(),
	}

	log.Init()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.checkFunc(tt.args); err != nil {
				t.Errorf("Run() error = %v", err)
			}
		})
	}
}
