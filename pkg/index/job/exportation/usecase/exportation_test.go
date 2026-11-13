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
	"reflect"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/file"
	"github.com/vdaas/vald/internal/os"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/internal/test/goleak"
	jconfig "github.com/vdaas/vald/pkg/index/job/exportation/config"
	"github.com/vdaas/vald/pkg/index/job/exportation/service"
)

// TestNew covers the two reachable branches of New without requiring a full
// vald.Client mock: the happy path (valid IndexPath, so pogreb.New and the
// starter/gateway construction all succeed) and a deterministic failure
// path (an IndexPath whose parent is a regular file, so the underlying
// service.New's file.MkdirAll fails and New propagates the error).
//
// New wires a real grpc gateway client, an errgroup and a starter server,
// none of which dial or bind anything during construction, so this test
// exercises the actual (non-mocked) construction path end-to-end against a
// temp directory instead of stubbing internal collaborators.
//
// Other error branches of New (grpc option parsing, starter TLS setup,
// observability construction) are not independently exercised here:
// reaching them deterministically would require either invalid low-level
// config.GRPCClient/TLS fixtures unrelated to this task's scope, or a full
// vald.Client mock for the gateway (see service/exporter_test.go for the
// same scoping rationale), which is disproportionate scaffolding for this
// bugfix + test-coverage task.
func TestNew(t *testing.T) {
	t.Parallel()

	t.Run("returns a Runner and nil when the config is valid", func(tt *testing.T) {
		tt.Parallel()
		defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())

		cfg := &jconfig.Data{
			Server:        &config.Servers{},
			Observability: &config.Observability{Enabled: false},
			Exporter: &config.IndexExporter{
				IndexPath:   tt.TempDir(),
				Concurrency: 1,
			},
		}

		r, err := New(cfg)
		if err != nil {
			tt.Fatalf("err: %v", err)
		}
		if r == nil {
			tt.Fatal("returned Runner must not be nil")
		}
		if err := r.PreStop(tt.Context()); err != nil {
			tt.Errorf("PreStop failed: %v", err)
		}
	})

	t.Run("applies the configured KVS background sync and compaction intervals", func(tt *testing.T) {
		tt.Parallel()
		defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())

		const (
			syncInterval       = 7 * time.Second
			compactionInterval = 11 * time.Second
		)
		cfg := &jconfig.Data{
			Server:        &config.Servers{},
			Observability: &config.Observability{Enabled: false},
			Exporter: &config.IndexExporter{
				IndexPath:                       tt.TempDir(),
				Concurrency:                     1,
				KVSBackgroundSyncInterval:       syncInterval.String(),
				KVSBackgroundCompactionInterval: compactionInterval.String(),
			},
		}

		r, err := New(cfg)
		if err != nil {
			tt.Fatalf("err: %v", err)
		}
		defer func() {
			if err := r.PreStop(tt.Context()); err != nil {
				tt.Errorf("PreStop failed: %v", err)
			}
		}()
		rr, ok := r.(*run)
		if !ok {
			tt.Fatalf("New returned %T, want *run", r)
		}
		// the exporter implementation is unexported in package service, so its
		// interval fields are read through reflection.
		ev := reflect.ValueOf(rr.exporter)
		if ev.Kind() != reflect.Pointer || ev.Elem().Kind() != reflect.Struct {
			tt.Fatalf("exporter is %T, want a pointer to a struct", rr.exporter)
		}
		for field, want := range map[string]time.Duration{
			"backgroundSyncInterval":       syncInterval,
			"backgroundCompactionInterval": compactionInterval,
		} {
			fv := ev.Elem().FieldByName(field)
			if !fv.IsValid() {
				tt.Fatalf("exporter has no field %s", field)
			}
			if got := time.Duration(fv.Int()); got != want {
				tt.Errorf("exporter %s = %v, want %v from the config", field, got, want)
			}
		}
	})

	t.Run("returns an error when the index path cannot be created", func(tt *testing.T) {
		tt.Parallel()
		defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())

		tmpDir := tt.TempDir()
		blocker := file.Join(tmpDir, "blocker")
		f, ferr := os.Create(blocker)
		if ferr != nil {
			tt.Fatal(ferr)
		}
		if ferr := f.Close(); ferr != nil {
			tt.Fatal(ferr)
		}

		cfg := &jconfig.Data{
			Server:        &config.Servers{},
			Observability: &config.Observability{Enabled: false},
			Exporter: &config.IndexExporter{
				IndexPath:   file.Join(blocker, "sub"),
				Concurrency: 1,
			},
		}

		r, err := New(cfg)
		if err == nil {
			tt.Error("expected an error when the index path parent is a regular file")
		}
		if r != nil {
			tt.Errorf("returned Runner must be nil, got: %#v", r)
		}
	})
}

type fakeListener struct{}

func (*fakeListener) ListenAndServe(context.Context) <-chan error {
	return nil
}

func (*fakeListener) Shutdown(context.Context) error {
	return nil
}

// fakeJob embeds the service interface so that only the methods used by
// run.Start need to be implemented.
type fakeJob struct {
	service.Exporter
}

func (*fakeJob) StartClient(context.Context) (<-chan error, error) {
	return make(chan error), nil
}

func (*fakeJob) Start(context.Context) error {
	return nil
}

// TestRun_Start_FindProcessError verifies that Start does not terminate the
// process (log.Fatal) when the job cannot find its own process to SIGTERM.
// With the old log.Fatalf this test binary would exit instead of reaching the
// assertions below.
func TestRun_Start_FindProcessError(t *testing.T) {
	// not parallel: os.FindProcess is a package level variable.
	called := make(chan struct{}, 1)
	orig := os.FindProcess
	t.Cleanup(func() { os.FindProcess = orig })
	os.FindProcess = func(int) (*os.Process, error) {
		select {
		case called <- struct{}{}:
		default:
		}
		return nil, errors.New("find process failed")
	}

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	eg, ctx := errgroup.New(ctx)
	r := &run{
		eg:       eg,
		server:   new(fakeListener),
		exporter: new(fakeJob),
	}
	if _, err := r.Start(ctx); err != nil {
		t.Fatalf("Start returned an error: %v", err)
	}

	select {
	case <-called:
	case <-time.After(10 * time.Second):
		t.Fatal("os.FindProcess was not called by the job goroutine")
	}

	cancel()
	done := make(chan error, 1)
	go func() {
		done <- eg.Wait()
	}()
	select {
	case err := <-done:
		if err != nil && !errors.Is(err, context.Canceled) {
			t.Errorf("unexpected error from errgroup: %v", err)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("goroutines started by Start did not finish")
	}
}

// NOT IMPLEMENTED BELOW
//
// func Test_run_PreStart(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		eg            errgroup.Group
// 		cfg           *config.Data
// 		observability observability.Observability
// 		server        starter.Server
// 		exporter      service.Exporter
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
// 		           observability:nil,
// 		           server:nil,
// 		           exporter:nil,
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
// 		           observability:nil,
// 		           server:nil,
// 		           exporter:nil,
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
// 				observability: test.fields.observability,
// 				server:        test.fields.server,
// 				exporter:      test.fields.exporter,
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
// 		observability observability.Observability
// 		server        starter.Server
// 		exporter      service.Exporter
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
// 		           observability:nil,
// 		           server:nil,
// 		           exporter:nil,
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
// 		           observability:nil,
// 		           server:nil,
// 		           exporter:nil,
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
// 				observability: test.fields.observability,
// 				server:        test.fields.server,
// 				exporter:      test.fields.exporter,
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
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		eg            errgroup.Group
// 		cfg           *config.Data
// 		observability observability.Observability
// 		server        starter.Server
// 		exporter      service.Exporter
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
// 		           observability:nil,
// 		           server:nil,
// 		           exporter:nil,
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
// 		           observability:nil,
// 		           server:nil,
// 		           exporter:nil,
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
// 				observability: test.fields.observability,
// 				server:        test.fields.server,
// 				exporter:      test.fields.exporter,
// 			}
//
// 			err := r.PreStop(test.args.ctx)
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
// 		observability observability.Observability
// 		server        starter.Server
// 		exporter      service.Exporter
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
// 		           observability:nil,
// 		           server:nil,
// 		           exporter:nil,
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
// 		           observability:nil,
// 		           server:nil,
// 		           exporter:nil,
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
// 				observability: test.fields.observability,
// 				server:        test.fields.server,
// 				exporter:      test.fields.exporter,
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
// 		in0 context.Context
// 	}
// 	type fields struct {
// 		eg            errgroup.Group
// 		cfg           *config.Data
// 		observability observability.Observability
// 		server        starter.Server
// 		exporter      service.Exporter
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
// 		           observability:nil,
// 		           server:nil,
// 		           exporter:nil,
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
// 		           observability:nil,
// 		           server:nil,
// 		           exporter:nil,
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
// 				observability: test.fields.observability,
// 				server:        test.fields.server,
// 				exporter:      test.fields.exporter,
// 			}
//
// 			err := r.PostStop(test.args.in0)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
