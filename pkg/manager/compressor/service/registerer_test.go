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

package service

import (
	"context"
	"reflect"
	"sync"
	"testing"

	"github.com/vdaas/vald/apis/grpc/payload"
	client "github.com/vdaas/vald/internal/client/compressor"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/worker"
	"go.uber.org/goleak"
)

func TestNewRegisterer(t *testing.T) {
	type args struct {
		opts []RegistererOption
	}
	type want struct {
		want Registerer
		err  error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, Registerer, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got Registerer, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got = %v, want %v", got, w.want)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
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
		           opts: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got, err := NewRegisterer(test.args.opts...)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_registerer_PreStart(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type fields struct {
		worker     worker.Worker
		workerOpts []worker.WorkerOption
		eg         errgroup.Group
		backup     Backup
		compressor Compressor
		client     client.Client
		metas      map[string]*payload.Backup_MetaVector
		metasMux   sync.Mutex
	}
	type want struct {
		err error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
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
		       },
		       fields: fields {
		           worker: nil,
		           workerOpts: nil,
		           eg: nil,
		           backup: nil,
		           compressor: nil,
		           client: nil,
		           metas: nil,
		           metasMux: sync.Mutex{},
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
		           },
		           fields: fields {
		           worker: nil,
		           workerOpts: nil,
		           eg: nil,
		           backup: nil,
		           compressor: nil,
		           client: nil,
		           metas: nil,
		           metasMux: sync.Mutex{},
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			r := &registerer{
				worker:     test.fields.worker,
				workerOpts: test.fields.workerOpts,
				eg:         test.fields.eg,
				backup:     test.fields.backup,
				compressor: test.fields.compressor,
				client:     test.fields.client,
				metas:      test.fields.metas,
				metasMux:   test.fields.metasMux,
			}

			err := r.PreStart(test.args.ctx)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_registerer_Start(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type fields struct {
		worker     worker.Worker
		workerOpts []worker.WorkerOption
		eg         errgroup.Group
		backup     Backup
		compressor Compressor
		client     client.Client
		metas      map[string]*payload.Backup_MetaVector
		metasMux   sync.Mutex
	}
	type want struct {
		want <-chan error
		err  error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, <-chan error, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got <-chan error, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got = %v, want %v", got, w.want)
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
		       },
		       fields: fields {
		           worker: nil,
		           workerOpts: nil,
		           eg: nil,
		           backup: nil,
		           compressor: nil,
		           client: nil,
		           metas: nil,
		           metasMux: sync.Mutex{},
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
		           },
		           fields: fields {
		           worker: nil,
		           workerOpts: nil,
		           eg: nil,
		           backup: nil,
		           compressor: nil,
		           client: nil,
		           metas: nil,
		           metasMux: sync.Mutex{},
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			r := &registerer{
				worker:     test.fields.worker,
				workerOpts: test.fields.workerOpts,
				eg:         test.fields.eg,
				backup:     test.fields.backup,
				compressor: test.fields.compressor,
				client:     test.fields.client,
				metas:      test.fields.metas,
				metasMux:   test.fields.metasMux,
			}

			got, err := r.Start(test.args.ctx)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_registerer_PostStop(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type fields struct {
		worker     worker.Worker
		workerOpts []worker.WorkerOption
		eg         errgroup.Group
		backup     Backup
		compressor Compressor
		client     client.Client
		metas      map[string]*payload.Backup_MetaVector
		metasMux   sync.Mutex
	}
	type want struct {
		err error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
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
		       },
		       fields: fields {
		           worker: nil,
		           workerOpts: nil,
		           eg: nil,
		           backup: nil,
		           compressor: nil,
		           client: nil,
		           metas: nil,
		           metasMux: sync.Mutex{},
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
		           },
		           fields: fields {
		           worker: nil,
		           workerOpts: nil,
		           eg: nil,
		           backup: nil,
		           compressor: nil,
		           client: nil,
		           metas: nil,
		           metasMux: sync.Mutex{},
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			r := &registerer{
				worker:     test.fields.worker,
				workerOpts: test.fields.workerOpts,
				eg:         test.fields.eg,
				backup:     test.fields.backup,
				compressor: test.fields.compressor,
				client:     test.fields.client,
				metas:      test.fields.metas,
				metasMux:   test.fields.metasMux,
			}

			err := r.PostStop(test.args.ctx)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_registerer_Register(t *testing.T) {
	type args struct {
		ctx  context.Context
		meta *payload.Backup_MetaVector
	}
	type fields struct {
		worker     worker.Worker
		workerOpts []worker.WorkerOption
		eg         errgroup.Group
		backup     Backup
		compressor Compressor
		client     client.Client
		metas      map[string]*payload.Backup_MetaVector
		metasMux   sync.Mutex
	}
	type want struct {
		err error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
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
		           meta: nil,
		       },
		       fields: fields {
		           worker: nil,
		           workerOpts: nil,
		           eg: nil,
		           backup: nil,
		           compressor: nil,
		           client: nil,
		           metas: nil,
		           metasMux: sync.Mutex{},
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
		           meta: nil,
		           },
		           fields: fields {
		           worker: nil,
		           workerOpts: nil,
		           eg: nil,
		           backup: nil,
		           compressor: nil,
		           client: nil,
		           metas: nil,
		           metasMux: sync.Mutex{},
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			r := &registerer{
				worker:     test.fields.worker,
				workerOpts: test.fields.workerOpts,
				eg:         test.fields.eg,
				backup:     test.fields.backup,
				compressor: test.fields.compressor,
				client:     test.fields.client,
				metas:      test.fields.metas,
				metasMux:   test.fields.metasMux,
			}

			err := r.Register(test.args.ctx, test.args.meta)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_registerer_RegisterMulti(t *testing.T) {
	type args struct {
		ctx   context.Context
		metas *payload.Backup_MetaVectors
	}
	type fields struct {
		worker     worker.Worker
		workerOpts []worker.WorkerOption
		eg         errgroup.Group
		backup     Backup
		compressor Compressor
		client     client.Client
		metas      map[string]*payload.Backup_MetaVector
		metasMux   sync.Mutex
	}
	type want struct {
		err error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
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
		           metas: nil,
		       },
		       fields: fields {
		           worker: nil,
		           workerOpts: nil,
		           eg: nil,
		           backup: nil,
		           compressor: nil,
		           client: nil,
		           metas: nil,
		           metasMux: sync.Mutex{},
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
		           metas: nil,
		           },
		           fields: fields {
		           worker: nil,
		           workerOpts: nil,
		           eg: nil,
		           backup: nil,
		           compressor: nil,
		           client: nil,
		           metas: nil,
		           metasMux: sync.Mutex{},
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			r := &registerer{
				worker:     test.fields.worker,
				workerOpts: test.fields.workerOpts,
				eg:         test.fields.eg,
				backup:     test.fields.backup,
				compressor: test.fields.compressor,
				client:     test.fields.client,
				metas:      test.fields.metas,
				metasMux:   test.fields.metasMux,
			}

			err := r.RegisterMulti(test.args.ctx, test.args.metas)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_registerer_Len(t *testing.T) {
	type fields struct {
		worker     worker.Worker
		workerOpts []worker.WorkerOption
		eg         errgroup.Group
		backup     Backup
		compressor Compressor
		client     client.Client
		metas      map[string]*payload.Backup_MetaVector
		metasMux   sync.Mutex
	}
	type want struct {
		want uint64
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, uint64) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got uint64) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got = %v, want %v", got, w.want)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       fields: fields {
		           worker: nil,
		           workerOpts: nil,
		           eg: nil,
		           backup: nil,
		           compressor: nil,
		           client: nil,
		           metas: nil,
		           metasMux: sync.Mutex{},
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
		           fields: fields {
		           worker: nil,
		           workerOpts: nil,
		           eg: nil,
		           backup: nil,
		           compressor: nil,
		           client: nil,
		           metas: nil,
		           metasMux: sync.Mutex{},
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			r := &registerer{
				worker:     test.fields.worker,
				workerOpts: test.fields.workerOpts,
				eg:         test.fields.eg,
				backup:     test.fields.backup,
				compressor: test.fields.compressor,
				client:     test.fields.client,
				metas:      test.fields.metas,
				metasMux:   test.fields.metasMux,
			}

			got := r.Len()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_registerer_TotalRequested(t *testing.T) {
	type fields struct {
		worker     worker.Worker
		workerOpts []worker.WorkerOption
		eg         errgroup.Group
		backup     Backup
		compressor Compressor
		client     client.Client
		metas      map[string]*payload.Backup_MetaVector
		metasMux   sync.Mutex
	}
	type want struct {
		want uint64
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, uint64) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got uint64) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got = %v, want %v", got, w.want)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       fields: fields {
		           worker: nil,
		           workerOpts: nil,
		           eg: nil,
		           backup: nil,
		           compressor: nil,
		           client: nil,
		           metas: nil,
		           metasMux: sync.Mutex{},
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
		           fields: fields {
		           worker: nil,
		           workerOpts: nil,
		           eg: nil,
		           backup: nil,
		           compressor: nil,
		           client: nil,
		           metas: nil,
		           metasMux: sync.Mutex{},
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			r := &registerer{
				worker:     test.fields.worker,
				workerOpts: test.fields.workerOpts,
				eg:         test.fields.eg,
				backup:     test.fields.backup,
				compressor: test.fields.compressor,
				client:     test.fields.client,
				metas:      test.fields.metas,
				metasMux:   test.fields.metasMux,
			}

			got := r.TotalRequested()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_registerer_TotalCompleted(t *testing.T) {
	type fields struct {
		worker     worker.Worker
		workerOpts []worker.WorkerOption
		eg         errgroup.Group
		backup     Backup
		compressor Compressor
		client     client.Client
		metas      map[string]*payload.Backup_MetaVector
		metasMux   sync.Mutex
	}
	type want struct {
		want uint64
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, uint64) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got uint64) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got = %v, want %v", got, w.want)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       fields: fields {
		           worker: nil,
		           workerOpts: nil,
		           eg: nil,
		           backup: nil,
		           compressor: nil,
		           client: nil,
		           metas: nil,
		           metasMux: sync.Mutex{},
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
		           fields: fields {
		           worker: nil,
		           workerOpts: nil,
		           eg: nil,
		           backup: nil,
		           compressor: nil,
		           client: nil,
		           metas: nil,
		           metasMux: sync.Mutex{},
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			r := &registerer{
				worker:     test.fields.worker,
				workerOpts: test.fields.workerOpts,
				eg:         test.fields.eg,
				backup:     test.fields.backup,
				compressor: test.fields.compressor,
				client:     test.fields.client,
				metas:      test.fields.metas,
				metasMux:   test.fields.metasMux,
			}

			got := r.TotalCompleted()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_registerer_dispatch(t *testing.T) {
	type args struct {
		ctx  context.Context
		meta *payload.Backup_MetaVector
	}
	type fields struct {
		worker     worker.Worker
		workerOpts []worker.WorkerOption
		eg         errgroup.Group
		backup     Backup
		compressor Compressor
		client     client.Client
		metas      map[string]*payload.Backup_MetaVector
		metasMux   sync.Mutex
	}
	type want struct {
		err error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
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
		           meta: nil,
		       },
		       fields: fields {
		           worker: nil,
		           workerOpts: nil,
		           eg: nil,
		           backup: nil,
		           compressor: nil,
		           client: nil,
		           metas: nil,
		           metasMux: sync.Mutex{},
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
		           meta: nil,
		           },
		           fields: fields {
		           worker: nil,
		           workerOpts: nil,
		           eg: nil,
		           backup: nil,
		           compressor: nil,
		           client: nil,
		           metas: nil,
		           metasMux: sync.Mutex{},
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			r := &registerer{
				worker:     test.fields.worker,
				workerOpts: test.fields.workerOpts,
				eg:         test.fields.eg,
				backup:     test.fields.backup,
				compressor: test.fields.compressor,
				client:     test.fields.client,
				metas:      test.fields.metas,
				metasMux:   test.fields.metasMux,
			}

			err := r.dispatch(test.args.ctx, test.args.meta)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_registerer_registerProcessFunc(t *testing.T) {
	type args struct {
		meta *payload.Backup_MetaVector
	}
	type fields struct {
		worker     worker.Worker
		workerOpts []worker.WorkerOption
		eg         errgroup.Group
		backup     Backup
		compressor Compressor
		client     client.Client
		metas      map[string]*payload.Backup_MetaVector
		metasMux   sync.Mutex
	}
	type want struct {
		want worker.JobFunc
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, worker.JobFunc) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got worker.JobFunc) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got = %v, want %v", got, w.want)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           meta: nil,
		       },
		       fields: fields {
		           worker: nil,
		           workerOpts: nil,
		           eg: nil,
		           backup: nil,
		           compressor: nil,
		           client: nil,
		           metas: nil,
		           metasMux: sync.Mutex{},
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
		           meta: nil,
		           },
		           fields: fields {
		           worker: nil,
		           workerOpts: nil,
		           eg: nil,
		           backup: nil,
		           compressor: nil,
		           client: nil,
		           metas: nil,
		           metasMux: sync.Mutex{},
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			r := &registerer{
				worker:     test.fields.worker,
				workerOpts: test.fields.workerOpts,
				eg:         test.fields.eg,
				backup:     test.fields.backup,
				compressor: test.fields.compressor,
				client:     test.fields.client,
				metas:      test.fields.metas,
				metasMux:   test.fields.metasMux,
			}

			got := r.registerProcessFunc(test.args.meta)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_registerer_forwardMetas(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type fields struct {
		worker     worker.Worker
		workerOpts []worker.WorkerOption
		eg         errgroup.Group
		backup     Backup
		compressor Compressor
		client     client.Client
		metas      map[string]*payload.Backup_MetaVector
		metasMux   sync.Mutex
	}
	type want struct {
		err error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
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
		       },
		       fields: fields {
		           worker: nil,
		           workerOpts: nil,
		           eg: nil,
		           backup: nil,
		           compressor: nil,
		           client: nil,
		           metas: nil,
		           metasMux: sync.Mutex{},
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
		           },
		           fields: fields {
		           worker: nil,
		           workerOpts: nil,
		           eg: nil,
		           backup: nil,
		           compressor: nil,
		           client: nil,
		           metas: nil,
		           metasMux: sync.Mutex{},
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			r := &registerer{
				worker:     test.fields.worker,
				workerOpts: test.fields.workerOpts,
				eg:         test.fields.eg,
				backup:     test.fields.backup,
				compressor: test.fields.compressor,
				client:     test.fields.client,
				metas:      test.fields.metas,
				metasMux:   test.fields.metasMux,
			}

			err := r.forwardMetas(test.args.ctx)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}
