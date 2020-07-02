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
	"sync/atomic"
	"testing"
	"time"

	"github.com/vdaas/vald/apis/grpc/payload"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/k8s"
	"go.uber.org/goleak"
)

func TestNew(t *testing.T) {
	type args struct {
		opts []Option
	}
	type want struct {
		wantDsc Discoverer
		err     error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, Discoverer, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotDsc Discoverer, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(gotDsc, w.wantDsc) {
			return errors.Errorf("got = %v, want %v", gotDsc, w.wantDsc)
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

			gotDsc, err := New(test.args.opts...)
			if err := test.checkFunc(test.want, gotDsc, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_discoverer_Start(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type fields struct {
		maxPods         int
		nodes           nodeMap
		nodeMetrics     nodeMetricsMap
		pods            podsMap
		podMetrics      podMetricsMap
		podsByNode      atomic.Value
		podsByNamespace atomic.Value
		podsByName      atomic.Value
		nodeByName      atomic.Value
		ctrl            k8s.Controller
		namespace       string
		name            string
		csd             time.Duration
		eg              errgroup.Group
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
		           maxPods: 0,
		           nodes: nodeMap{},
		           nodeMetrics: nodeMetricsMap{},
		           pods: podsMap{},
		           podMetrics: podMetricsMap{},
		           podsByNode: nil,
		           podsByNamespace: nil,
		           podsByName: nil,
		           nodeByName: nil,
		           ctrl: nil,
		           namespace: "",
		           name: "",
		           csd: nil,
		           eg: nil,
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
		           maxPods: 0,
		           nodes: nodeMap{},
		           nodeMetrics: nodeMetricsMap{},
		           pods: podsMap{},
		           podMetrics: podMetricsMap{},
		           podsByNode: nil,
		           podsByNamespace: nil,
		           podsByName: nil,
		           nodeByName: nil,
		           ctrl: nil,
		           namespace: "",
		           name: "",
		           csd: nil,
		           eg: nil,
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
			d := &discoverer{
				maxPods:         test.fields.maxPods,
				nodes:           test.fields.nodes,
				nodeMetrics:     test.fields.nodeMetrics,
				pods:            test.fields.pods,
				podMetrics:      test.fields.podMetrics,
				podsByNode:      test.fields.podsByNode,
				podsByNamespace: test.fields.podsByNamespace,
				podsByName:      test.fields.podsByName,
				nodeByName:      test.fields.nodeByName,
				ctrl:            test.fields.ctrl,
				namespace:       test.fields.namespace,
				name:            test.fields.name,
				csd:             test.fields.csd,
				eg:              test.fields.eg,
			}

			got, err := d.Start(test.args.ctx)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_discoverer_GetPods(t *testing.T) {
	type args struct {
		req *payload.Discoverer_Request
	}
	type fields struct {
		maxPods         int
		nodes           nodeMap
		nodeMetrics     nodeMetricsMap
		pods            podsMap
		podMetrics      podMetricsMap
		podsByNode      atomic.Value
		podsByNamespace atomic.Value
		podsByName      atomic.Value
		nodeByName      atomic.Value
		ctrl            k8s.Controller
		namespace       string
		name            string
		csd             time.Duration
		eg              errgroup.Group
	}
	type want struct {
		wantPods *payload.Info_Pods
		err      error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *payload.Info_Pods, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotPods *payload.Info_Pods, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(gotPods, w.wantPods) {
			return errors.Errorf("got = %v, want %v", gotPods, w.wantPods)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           req: nil,
		       },
		       fields: fields {
		           maxPods: 0,
		           nodes: nodeMap{},
		           nodeMetrics: nodeMetricsMap{},
		           pods: podsMap{},
		           podMetrics: podMetricsMap{},
		           podsByNode: nil,
		           podsByNamespace: nil,
		           podsByName: nil,
		           nodeByName: nil,
		           ctrl: nil,
		           namespace: "",
		           name: "",
		           csd: nil,
		           eg: nil,
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
		           req: nil,
		           },
		           fields: fields {
		           maxPods: 0,
		           nodes: nodeMap{},
		           nodeMetrics: nodeMetricsMap{},
		           pods: podsMap{},
		           podMetrics: podMetricsMap{},
		           podsByNode: nil,
		           podsByNamespace: nil,
		           podsByName: nil,
		           nodeByName: nil,
		           ctrl: nil,
		           namespace: "",
		           name: "",
		           csd: nil,
		           eg: nil,
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
			d := &discoverer{
				maxPods:         test.fields.maxPods,
				nodes:           test.fields.nodes,
				nodeMetrics:     test.fields.nodeMetrics,
				pods:            test.fields.pods,
				podMetrics:      test.fields.podMetrics,
				podsByNode:      test.fields.podsByNode,
				podsByNamespace: test.fields.podsByNamespace,
				podsByName:      test.fields.podsByName,
				nodeByName:      test.fields.nodeByName,
				ctrl:            test.fields.ctrl,
				namespace:       test.fields.namespace,
				name:            test.fields.name,
				csd:             test.fields.csd,
				eg:              test.fields.eg,
			}

			gotPods, err := d.GetPods(test.args.req)
			if err := test.checkFunc(test.want, gotPods, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_discoverer_GetNodes(t *testing.T) {
	type args struct {
		req *payload.Discoverer_Request
	}
	type fields struct {
		maxPods         int
		nodes           nodeMap
		nodeMetrics     nodeMetricsMap
		pods            podsMap
		podMetrics      podMetricsMap
		podsByNode      atomic.Value
		podsByNamespace atomic.Value
		podsByName      atomic.Value
		nodeByName      atomic.Value
		ctrl            k8s.Controller
		namespace       string
		name            string
		csd             time.Duration
		eg              errgroup.Group
	}
	type want struct {
		wantNodes *payload.Info_Nodes
		err       error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *payload.Info_Nodes, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotNodes *payload.Info_Nodes, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(gotNodes, w.wantNodes) {
			return errors.Errorf("got = %v, want %v", gotNodes, w.wantNodes)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           req: nil,
		       },
		       fields: fields {
		           maxPods: 0,
		           nodes: nodeMap{},
		           nodeMetrics: nodeMetricsMap{},
		           pods: podsMap{},
		           podMetrics: podMetricsMap{},
		           podsByNode: nil,
		           podsByNamespace: nil,
		           podsByName: nil,
		           nodeByName: nil,
		           ctrl: nil,
		           namespace: "",
		           name: "",
		           csd: nil,
		           eg: nil,
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
		           req: nil,
		           },
		           fields: fields {
		           maxPods: 0,
		           nodes: nodeMap{},
		           nodeMetrics: nodeMetricsMap{},
		           pods: podsMap{},
		           podMetrics: podMetricsMap{},
		           podsByNode: nil,
		           podsByNamespace: nil,
		           podsByName: nil,
		           nodeByName: nil,
		           ctrl: nil,
		           namespace: "",
		           name: "",
		           csd: nil,
		           eg: nil,
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
			d := &discoverer{
				maxPods:         test.fields.maxPods,
				nodes:           test.fields.nodes,
				nodeMetrics:     test.fields.nodeMetrics,
				pods:            test.fields.pods,
				podMetrics:      test.fields.podMetrics,
				podsByNode:      test.fields.podsByNode,
				podsByNamespace: test.fields.podsByNamespace,
				podsByName:      test.fields.podsByName,
				nodeByName:      test.fields.nodeByName,
				ctrl:            test.fields.ctrl,
				namespace:       test.fields.namespace,
				name:            test.fields.name,
				csd:             test.fields.csd,
				eg:              test.fields.eg,
			}

			gotNodes, err := d.GetNodes(test.args.req)
			if err := test.checkFunc(test.want, gotNodes, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}
