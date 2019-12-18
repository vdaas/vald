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
	"errors"
	"fmt"
	"testing"

	"github.com/vdaas/vald/internal/log"
)

func TestDo(t *testing.T) {
	type args struct {
		ctx  context.Context
		opts []Option
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Do(tt.args.ctx, tt.args.opts...); (err != nil) != tt.wantErr {
				t.Errorf("Do() error = %v, wantErr %v", err, tt.wantErr)
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
			log.Init(log.DefaultGlg())
			ctx, cancel := context.WithCancel(context.Background())
			run := &runnerMock{}

			return test{
				name: "run success",
				args: args{
					ctx: ctx,
					run: run,
				},
				checkFunc: func(args args) error {
					cancel()
					return Run(args.ctx, args.run, args.name)
				},
			}
		}(),
		func() test {
			log.Init(log.DefaultGlg())
			ctx, cancel := context.WithCancel(context.Background())
			run := &runnerMock{
				PreStartFunc: func(ctx context.Context) error {
					return errors.New("prestart error")
				},
			}

			return test{
				name: "run with prestart error",
				args: args{
					ctx: ctx,
					run: run,
				},
				checkFunc: func(args args) error {
					cancel()
					err := Run(args.ctx, args.run, args.name)
					if err.Error() != "prestart error" {
						return errors.New("prestart error should be thrown")
					}
					return nil
				},
			}
		}(),
		func() test {
			log.Init(log.DefaultGlg())
			ctx, cancel := context.WithCancel(context.Background())
			run := &runnerMock{
				StartFunc: func(ctx context.Context) <-chan error {
					errChan := make(chan error, 1)
					errChan <- errors.New("start error")
					return errChan
				},
			}

			return test{
				name: "run with start error",
				args: args{
					ctx: ctx,
					run: run,
				},
				checkFunc: func(args args) error {
					cancel()
					err := Run(args.ctx, args.run, args.name)
					if err.Error() != "error:\tstart error\tcount:\t0" {
						return fmt.Errorf("start error should be thrown, got: %s", err)
					}
					return nil
				},
			}
		}(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.checkFunc(tt.args); err != nil {
				t.Errorf("Run() error = %v", err)
			}
		})
	}
}
