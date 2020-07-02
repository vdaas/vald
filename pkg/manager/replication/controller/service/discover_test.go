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

// Package service manages the main logic of server.
package service

import (
	"context"
	"reflect"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/k8s"
	"github.com/vdaas/vald/internal/net/grpc"
	"go.uber.org/goleak"
)

func TestNew(t *testing.T) {
	type args struct {
		opts []Option
	}
	type want struct {
		wantRp Replicator
		err    error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, Replicator, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotRp Replicator, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(gotRp, w.wantRp) {
			return errors.Errorf("got = %v, want %v", gotRp, w.wantRp)
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

			gotRp, err := New(test.args.opts...)
			if err := test.checkFunc(test.want, gotRp, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_replicator_Start(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type fields struct {
		pods      atomic.Value
		ctrl      k8s.Controller
		namespace string
		name      string
		eg        errgroup.Group
		rdur      time.Duration
		rpods     sync.Map
		client    grpc.Client
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
		           pods: nil,
		           ctrl: nil,
		           namespace: "",
		           name: "",
		           eg: nil,
		           rdur: nil,
		           rpods: sync.Map{},
		           client: nil,
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
		           pods: nil,
		           ctrl: nil,
		           namespace: "",
		           name: "",
		           eg: nil,
		           rdur: nil,
		           rpods: sync.Map{},
		           client: nil,
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
			r := &replicator{
				pods:      test.fields.pods,
				ctrl:      test.fields.ctrl,
				namespace: test.fields.namespace,
				name:      test.fields.name,
				eg:        test.fields.eg,
				rdur:      test.fields.rdur,
				rpods:     test.fields.rpods,
				client:    test.fields.client,
			}

			got, err := r.Start(test.args.ctx)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_replicator_GetCurrentPodIPs(t *testing.T) {
	type fields struct {
		pods      atomic.Value
		ctrl      k8s.Controller
		namespace string
		name      string
		eg        errgroup.Group
		rdur      time.Duration
		rpods     sync.Map
		client    grpc.Client
	}
	type want struct {
		want  []string
		want1 bool
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, []string, bool) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got []string, got1 bool) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got = %v, want %v", got, w.want)
		}
		if !reflect.DeepEqual(got1, w.want1) {
			return errors.Errorf("got = %v, want %v", got1, w.want1)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       fields: fields {
		           pods: nil,
		           ctrl: nil,
		           namespace: "",
		           name: "",
		           eg: nil,
		           rdur: nil,
		           rpods: sync.Map{},
		           client: nil,
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
		           pods: nil,
		           ctrl: nil,
		           namespace: "",
		           name: "",
		           eg: nil,
		           rdur: nil,
		           rpods: sync.Map{},
		           client: nil,
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
			r := &replicator{
				pods:      test.fields.pods,
				ctrl:      test.fields.ctrl,
				namespace: test.fields.namespace,
				name:      test.fields.name,
				eg:        test.fields.eg,
				rdur:      test.fields.rdur,
				rpods:     test.fields.rpods,
				client:    test.fields.client,
			}

			got, got1 := r.GetCurrentPodIPs()
			if err := test.checkFunc(test.want, got, got1); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_replicator_SendRecoveryRequest(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type fields struct {
		pods      atomic.Value
		ctrl      k8s.Controller
		namespace string
		name      string
		eg        errgroup.Group
		rdur      time.Duration
		rpods     sync.Map
		client    grpc.Client
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
		           pods: nil,
		           ctrl: nil,
		           namespace: "",
		           name: "",
		           eg: nil,
		           rdur: nil,
		           rpods: sync.Map{},
		           client: nil,
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
		           pods: nil,
		           ctrl: nil,
		           namespace: "",
		           name: "",
		           eg: nil,
		           rdur: nil,
		           rpods: sync.Map{},
		           client: nil,
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
			r := &replicator{
				pods:      test.fields.pods,
				ctrl:      test.fields.ctrl,
				namespace: test.fields.namespace,
				name:      test.fields.name,
				eg:        test.fields.eg,
				rdur:      test.fields.rdur,
				rpods:     test.fields.rpods,
				client:    test.fields.client,
			}

			err := r.SendRecoveryRequest(test.args.ctx)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}
