// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package service

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/k8s"
	"github.com/vdaas/vald/internal/sync/atomic"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/pkg/operator/vald/config"
)

func TestNew(t *testing.T) {
	t.Parallel()

	type test struct {
		wantErr error
		cfg     *config.Operator
		name    string
	}

	tests := []test{
		{
			name:    "returns ErrInvalidConfig when config is nil",
			cfg:     nil,
			wantErr: errors.ErrInvalidConfig,
		},
		{
			name: "returns error when the default vrs template cannot be loaded",
			cfg: &config.Operator{
				Name: "vald-operator",
				Vrs: &config.Vrs{
					DefaultVrsPath: filepath.Join(t.TempDir(), "missing.yaml"),
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_, err := New(tc.cfg)
			if err == nil {
				t.Fatal("New() error = nil, want error")
			}
			if tc.wantErr != nil && !errors.Is(err, tc.wantErr) {
				t.Errorf("New() error = %v, want %v", err, tc.wantErr)
			}
		})
	}
}

// ctrlMock is a minimal k8s.Controller stub whose Start returns a
// pre-populated error channel.
type ctrlMock struct {
	dech chan error
}

func (c *ctrlMock) Start(context.Context) (<-chan error, error) {
	return c.dech, nil
}

func (*ctrlMock) GetManager() k8s.Manager {
	return nil
}

// TestOperator_Start_CtxCancelUnblocksErrorForward pins the ctx guard on the
// ech forward: with no reader on ech (capacity 2), the third controller error
// blocks the send, and only a ctx-guarded select lets the goroutine exit on
// cancellation instead of leaking.
func TestOperator_Start_CtxCancelUnblocksErrorForward(t *testing.T) {
	t.Parallel()

	dech := make(chan error, 3)
	for range 3 {
		dech <- errors.New("controller error")
	}

	eg, egctx := errgroup.New(context.Background())
	ctx, cancel := context.WithCancel(egctx)
	defer cancel()

	o := &operator{ctrl: &ctrlMock{dech: dech}, eg: eg}
	ech, err := o.Start(ctx)
	if err != nil {
		t.Fatalf("Start() error = %v, want nil", err)
	}

	// Wait until the forwarding goroutine has filled ech and is blocked on
	// the third send.
	deadline := time.Now().Add(5 * time.Second)
	for len(ech) < cap(ech) {
		if time.Now().After(deadline) {
			t.Fatal("error forwarding did not fill the channel buffer")
		}
		time.Sleep(time.Millisecond)
	}

	cancel()

	done := make(chan struct{})
	go func() {
		defer close(done)
		_ = eg.Wait()
	}()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("forwarding goroutine did not exit on context cancel; the ech send must be ctx-guarded")
	}
}

// doneCountCtx counts Done() calls. A select loop evaluates ctx.Done() on every
// iteration, so the count is a deterministic measure of how often it spins.
type doneCountCtx struct {
	context.Context //nolint:containedctx // test wrapper that counts Done() calls
	n               atomic.Int64
}

func (c *doneCountCtx) Done() <-chan struct{} {
	c.n.Add(1)
	return c.Context.Done()
}

// TestOperator_Start_ClosedControllerChannelDoesNotSpin pins that the relay
// loop stops selecting the controller error channel once it is closed (the
// controller's manager returned while ctx is still live). A closed channel is
// always ready, so without nil-ing it the loop busy-spins; the spin shows up
// as an unbounded number of ctx.Done() evaluations.
func TestOperator_Start_ClosedControllerChannelDoesNotSpin(t *testing.T) {
	t.Parallel()

	dech := make(chan error)
	close(dech)

	eg, egctx := errgroup.New(context.Background())
	cctx, cancel := context.WithCancel(egctx)
	defer cancel()
	ctx := &doneCountCtx{Context: cctx}

	o := &operator{ctrl: &ctrlMock{dech: dech}, eg: eg}
	if _, err := o.Start(ctx); err != nil {
		t.Fatalf("Start() error = %v, want nil", err)
	}

	time.Sleep(100 * time.Millisecond)
	if n := ctx.n.Load(); n > 100 {
		t.Errorf("ctx.Done() evaluated %d times with a closed controller channel; the relay loop is busy-spinning", n)
	}

	cancel()
	done := make(chan struct{})
	go func() {
		defer close(done)
		if err := eg.Wait(); err != nil {
			t.Logf("eg.Wait() = %v", err)
		}
	}()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("relay goroutine did not exit on context cancel")
	}
}

// NOT IMPLEMENTED BELOW
//
// func Test_operator_Start(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		ctrl k8s.Controller
// 		eg   errgroup.Group
// 	}
// 	type want struct {
// 		want <-chan error
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, <-chan error, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got <-chan error, err error) error {
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
// 		           ctx:nil,
// 		       },
// 		       fields: fields {
// 		           ctrl:nil,
// 		           eg:nil,
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
// 		           ctx:nil,
// 		           },
// 		           fields: fields {
// 		           ctrl:nil,
// 		           eg:nil,
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
// 			o := &operator{
// 				ctrl: test.fields.ctrl,
// 				eg:   test.fields.eg,
// 			}
//
// 			got, err := o.Start(test.args.ctx)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
