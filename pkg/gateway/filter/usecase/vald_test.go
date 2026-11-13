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

package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/client/v1/client/filter/egress"
	"github.com/vdaas/vald/internal/client/v1/client/filter/ingress"
	client "github.com/vdaas/vald/internal/client/v1/client/vald"
	"github.com/vdaas/vald/internal/sync/errgroup"
)

// clientLifecycle records the context a fake outbound client was started with
// and whether it was stopped.
type clientLifecycle struct {
	ctx     context.Context //nolint:containedctx // records the Start context under test
	stopped bool
}

func (l *clientLifecycle) start(ctx context.Context) (<-chan error, error) {
	l.ctx = ctx
	// never closed: the error forwarding loop of Start receives from it.
	return make(chan error), nil
}

func (l *clientLifecycle) stop(context.Context) error {
	l.stopped = true
	return nil
}

type fakeValdClient struct {
	client.Client
	l clientLifecycle
}

func (f *fakeValdClient) Start(ctx context.Context) (<-chan error, error) { return f.l.start(ctx) }
func (f *fakeValdClient) Stop(ctx context.Context) error                  { return f.l.stop(ctx) }

type fakeIngressClient struct {
	ingress.Client
	l clientLifecycle
}

func (f *fakeIngressClient) Start(ctx context.Context) (<-chan error, error) { return f.l.start(ctx) }

func (f *fakeIngressClient) Stop(ctx context.Context) error { return f.l.stop(ctx) }

type fakeEgressClient struct {
	egress.Client
	l clientLifecycle
}

func (f *fakeEgressClient) Start(ctx context.Context) (<-chan error, error) { return f.l.start(ctx) }
func (f *fakeEgressClient) Stop(ctx context.Context) error                  { return f.l.stop(ctx) }

type fakeServer struct{}

func (fakeServer) ListenAndServe(context.Context) <-chan error { return nil }
func (fakeServer) Shutdown(context.Context) error              { return nil }

// Test_run_outboundClientsOutliveShutdownSignal pins that the shutdown signal
// (cancellation of the Start context) does not stop the outbound clients, so
// the server can still drain requests that use them; PostStop stops them.
func Test_run_outboundClientsOutliveShutdownSignal(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()
	eg, egctx := errgroup.New(ctx)
	vc, ic, ec := new(fakeValdClient), new(fakeIngressClient), new(fakeEgressClient)
	r := &run{
		eg:      eg,
		server:  fakeServer{},
		client:  vc,
		ingress: ic,
		egress:  ec,
	}
	sigctx, signal := context.WithCancel(egctx)
	if _, err := r.Start(sigctx); err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	signal()
	lifecycles := map[string]*clientLifecycle{"vald": &vc.l, "ingress": &ic.l, "egress": &ec.l}
	for name, l := range lifecycles {
		if l.ctx == nil {
			t.Fatalf("%s client was not started", name)
		}
		if err := l.ctx.Err(); err != nil {
			t.Errorf("%s client context is done after the shutdown signal: %v", name, err)
		}
	}
	if err := r.Stop(ctx); err != nil {
		t.Fatalf("Stop() error = %v", err)
	}
	for name, l := range lifecycles {
		if err := l.ctx.Err(); err != nil {
			t.Errorf("%s client context is done before PostStop: %v", name, err)
		}
	}
	if err := r.PostStop(ctx); err != nil {
		t.Fatalf("PostStop() error = %v", err)
	}
	for name, l := range lifecycles {
		if !l.stopped {
			t.Errorf("%s client was not stopped by PostStop", name)
		}
		if l.ctx.Err() == nil {
			t.Errorf("%s client context is not cancelled after PostStop", name)
		}
	}
	done := make(chan error, 1)
	go func() { done <- eg.Wait() }()
	select {
	case <-done:
	case <-ctx.Done():
		t.Fatal("timed out waiting for the error forwarding goroutine")
	}
}

// NOT IMPLEMENTED BELOW
//
// func TestNew(t *testing.T) {
// 	type args struct {
// 		cfg *config.Data
// 	}
// 	type want struct {
// 		wantR runner.Interface
// 		err   error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, runner.Interface, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotR runner.Interface, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotR, w.wantR) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotR, w.wantR)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           cfg:nil,
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
// 		           cfg:nil,
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
// 			gotR, err := New(test.args.cfg)
// 			if err := checkFunc(test.want, gotR, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_run_PreStart(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		eg            errgroup.Group
// 		cfg           *config.Data
// 		server        starter.Server
// 		observability observability.Observability
// 		client        client.Client
// 		ingress       ingress.Client
// 		egress        egress.Client
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
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
// 		           eg:nil,
// 		           cfg:nil,
// 		           server:nil,
// 		           observability:nil,
// 		           client:nil,
// 		           ingress:nil,
// 		           egress:nil,
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
// 		           eg:nil,
// 		           cfg:nil,
// 		           server:nil,
// 		           observability:nil,
// 		           client:nil,
// 		           ingress:nil,
// 		           egress:nil,
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
// 			r := &run{
// 				eg:            test.fields.eg,
// 				cfg:           test.fields.cfg,
// 				server:        test.fields.server,
// 				observability: test.fields.observability,
// 				client:        test.fields.client,
// 				ingress:       test.fields.ingress,
// 				egress:        test.fields.egress,
// 			}
//
// 			err := r.PreStart(test.args.ctx)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_run_Start(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		eg            errgroup.Group
// 		cfg           *config.Data
// 		server        starter.Server
// 		observability observability.Observability
// 		client        client.Client
// 		ingress       ingress.Client
// 		egress        egress.Client
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
// 		           eg:nil,
// 		           cfg:nil,
// 		           server:nil,
// 		           observability:nil,
// 		           client:nil,
// 		           ingress:nil,
// 		           egress:nil,
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
// 		           eg:nil,
// 		           cfg:nil,
// 		           server:nil,
// 		           observability:nil,
// 		           client:nil,
// 		           ingress:nil,
// 		           egress:nil,
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
// 			r := &run{
// 				eg:            test.fields.eg,
// 				cfg:           test.fields.cfg,
// 				server:        test.fields.server,
// 				observability: test.fields.observability,
// 				client:        test.fields.client,
// 				ingress:       test.fields.ingress,
// 				egress:        test.fields.egress,
// 			}
//
// 			got, err := r.Start(test.args.ctx)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_run_PreStop(t *testing.T) {
// 	type args struct {
// 		in0 context.Context
// 	}
// 	type fields struct {
// 		eg            errgroup.Group
// 		cfg           *config.Data
// 		server        starter.Server
// 		observability observability.Observability
// 		client        client.Client
// 		ingress       ingress.Client
// 		egress        egress.Client
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           in0:nil,
// 		       },
// 		       fields: fields {
// 		           eg:nil,
// 		           cfg:nil,
// 		           server:nil,
// 		           observability:nil,
// 		           client:nil,
// 		           ingress:nil,
// 		           egress:nil,
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
// 		           in0:nil,
// 		           },
// 		           fields: fields {
// 		           eg:nil,
// 		           cfg:nil,
// 		           server:nil,
// 		           observability:nil,
// 		           client:nil,
// 		           ingress:nil,
// 		           egress:nil,
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
// 			r := &run{
// 				eg:            test.fields.eg,
// 				cfg:           test.fields.cfg,
// 				server:        test.fields.server,
// 				observability: test.fields.observability,
// 				client:        test.fields.client,
// 				ingress:       test.fields.ingress,
// 				egress:        test.fields.egress,
// 			}
//
// 			err := r.PreStop(test.args.in0)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_run_Stop(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		eg            errgroup.Group
// 		cfg           *config.Data
// 		server        starter.Server
// 		observability observability.Observability
// 		client        client.Client
// 		ingress       ingress.Client
// 		egress        egress.Client
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
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
// 		           eg:nil,
// 		           cfg:nil,
// 		           server:nil,
// 		           observability:nil,
// 		           client:nil,
// 		           ingress:nil,
// 		           egress:nil,
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
// 		           eg:nil,
// 		           cfg:nil,
// 		           server:nil,
// 		           observability:nil,
// 		           client:nil,
// 		           ingress:nil,
// 		           egress:nil,
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
// 			r := &run{
// 				eg:            test.fields.eg,
// 				cfg:           test.fields.cfg,
// 				server:        test.fields.server,
// 				observability: test.fields.observability,
// 				client:        test.fields.client,
// 				ingress:       test.fields.ingress,
// 				egress:        test.fields.egress,
// 			}
//
// 			err := r.Stop(test.args.ctx)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_run_PostStop(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		eg            errgroup.Group
// 		cfg           *config.Data
// 		server        starter.Server
// 		observability observability.Observability
// 		client        client.Client
// 		ingress       ingress.Client
// 		egress        egress.Client
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
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
// 		           eg:nil,
// 		           cfg:nil,
// 		           server:nil,
// 		           observability:nil,
// 		           client:nil,
// 		           ingress:nil,
// 		           egress:nil,
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
// 		           eg:nil,
// 		           cfg:nil,
// 		           server:nil,
// 		           observability:nil,
// 		           client:nil,
// 		           ingress:nil,
// 		           egress:nil,
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
// 			r := &run{
// 				eg:            test.fields.eg,
// 				cfg:           test.fields.cfg,
// 				server:        test.fields.server,
// 				observability: test.fields.observability,
// 				client:        test.fields.client,
// 				ingress:       test.fields.ingress,
// 				egress:        test.fields.egress,
// 			}
//
// 			err := r.PostStop(test.args.ctx)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
