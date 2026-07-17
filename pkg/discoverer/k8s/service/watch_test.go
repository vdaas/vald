//
// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
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
	"testing"

	"github.com/vdaas/vald/internal/k8s"
	"github.com/vdaas/vald/internal/sync"
)

const testNodeAName = "node-a"

func TestSyncOnReconcile(t *testing.T) {
	t.Parallel()

	tests := []struct {
		initial   map[string]string
		converted map[string]string
		want      map[string]string
		name      string
	}{
		{
			name:      "stores entries into an empty map",
			converted: map[string]string{"a": "1", "b": "2"},
			want:      map[string]string{"a": "1", "b": "2"},
		},
		{
			name:      "updates existing entries and prunes stale ones",
			initial:   map[string]string{"a": "old", "stale": "x"},
			converted: map[string]string{"a": "new", "b": "2"},
			want:      map[string]string{"a": "new", "b": "2"},
		},
		{
			name:    "empty snapshot prunes everything",
			initial: map[string]string{"a": "1", "b": "2"},
			want:    map[string]string{},
		},
		{
			name:      "identical snapshot keeps entries",
			initial:   map[string]string{"a": "1"},
			converted: map[string]string{"a": "1"},
			want:      map[string]string{"a": "1"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			var dst sync.Map[string, string]
			for k, v := range test.initial {
				dst.Store(k, v)
			}
			list := new(k8s.PodList)
			cb := syncOnReconcile(&dst, func(l *k8s.PodList) map[string]string {
				if l != list {
					t.Errorf("conv received list %p, want %p", l, list)
				}
				return test.converted
			})
			cb(context.Background(), list)

			got := make(map[string]string)
			dst.Range(func(key, value string) bool {
				got[key] = value
				return true
			})
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("synced map = %#v, want %#v", got, test.want)
			}
		})
	}
}

func TestDebugLogged(t *testing.T) {
	t.Parallel()

	in := new(k8s.NodeList)
	want := map[string]NodeMetrics{testNodeAName: {Name: testNodeAName, CPU: 1.0}}
	conv := debugLogged("test subject", func(l *k8s.NodeList) map[string]NodeMetrics {
		if l != in {
			t.Errorf("conv received list %p, want %p", l, in)
		}
		return want
	})
	if got := conv(in); !reflect.DeepEqual(got, want) {
		t.Errorf("debugLogged conv = %#v, want %#v", got, want)
	}
}

func TestDiscovererPodEntries(t *testing.T) {
	t.Parallel()

	t.Run("groups pods into distinct pointer entries and raises maxPods", func(t *testing.T) {
		t.Parallel()

		d := &discoverer{namespace: "default", maxPods: 1}
		list := &k8s.PodList{Items: []k8s.Pod{
			newRunningPod("agent-0", "vald-agent"),
			newRunningPod("agent-1", "vald-agent"),
			newRunningPod("gateway-0", "vald-lb-gateway"),
		}}

		got := d.podEntries(list)
		if len(got) != 2 {
			t.Fatalf("podEntries() groups = %d, want 2", len(got))
		}
		agents := got["vald-agent"]
		if agents == nil || len(*agents) != 2 {
			t.Fatalf("vald-agent entry = %v, want 2 pods", agents)
		}
		gateways := got["vald-lb-gateway"]
		if gateways == nil || len(*gateways) != 1 {
			t.Fatalf("vald-lb-gateway entry = %v, want 1 pod", gateways)
		}
		if agents == gateways {
			t.Error("entries share a single pointer")
		}
		if d.maxPods != 2 {
			t.Errorf("maxPods = %d, want 2", d.maxPods)
		}
	})

	t.Run("keeps maxPods when every group is smaller", func(t *testing.T) {
		t.Parallel()

		d := &discoverer{maxPods: 5}
		list := &k8s.PodList{Items: []k8s.Pod{
			newRunningPod("agent-0", "vald-agent"),
		}}

		if got := d.podEntries(list); len(got) != 1 {
			t.Fatalf("podEntries() groups = %d, want 1", len(got))
		}
		if d.maxPods != 5 {
			t.Errorf("maxPods = %d, want 5", d.maxPods)
		}
	})
}

func TestNodeEntries(t *testing.T) {
	t.Parallel()

	list := &k8s.NodeList{Items: []k8s.Node{
		{ObjectMeta: k8s.ObjectMeta{Name: testNodeAName}},
		{ObjectMeta: k8s.ObjectMeta{Name: "node-b"}},
	}}

	got := nodeEntries(list)
	if len(got) != 2 {
		t.Fatalf("nodeEntries() = %d entries, want 2", len(got))
	}
	a := got[testNodeAName]
	if a == nil || a.Name != testNodeAName {
		t.Fatalf("node-a entry = %+v, want converted node", a)
	}
	if got[testNodeAName] == got["node-b"] {
		t.Error("entries share a single pointer")
	}
}

// NOT IMPLEMENTED BELOW
//
// func Test_listSyncController(t *testing.T) {
// 	type args struct {
// 		name      string
// 		what      string
// 		obj       k8s.Object
// 		dst       *sync.Map[string, V]
// 		conv      func(list PL) map[string]V
// 		namespace string
// 		fields    map[string]string
// 		labels    map[string]string
// 		extra     []reconciler.ListOption
// 	}
// 	type want struct {
// 		want k8s.Option
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, k8s.Option) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got k8s.Option) error {
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
// 		           name:"",
// 		           what:"",
// 		           obj:nil,
// 		           dst:nil,
// 		           conv:nil,
// 		           namespace:"",
// 		           fields:nil,
// 		           labels:nil,
// 		           extra:nil,
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
// 		           name:"",
// 		           what:"",
// 		           obj:nil,
// 		           dst:nil,
// 		           conv:nil,
// 		           namespace:"",
// 		           fields:nil,
// 		           labels:nil,
// 		           extra:nil,
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
// 			got := listSyncController(test.args.name, test.args.what, test.args.obj, test.args.dst, test.args.conv, test.args.namespace, test.args.fields, test.args.labels, test.args.extra...)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_syncOnReconcile(t *testing.T) {
// 	type args struct {
// 		dst  *sync.Map[string, V]
// 		conv func(list L) map[string]V
// 	}
// 	type want struct {
// 		want func(ctx context.Context, list L)
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, func(ctx context.Context, list L)) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got func(ctx context.Context, list L)) error {
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
// 		           dst:nil,
// 		           conv:nil,
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
// 		           dst:nil,
// 		           conv:nil,
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
// 			got := syncOnReconcile(test.args.dst, test.args.conv)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_debugLogged(t *testing.T) {
// 	type args struct {
// 		what string
// 		conv func(I) O
// 	}
// 	type want struct {
// 		want func(I) O
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, func(I) O) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got func(I) O) error {
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
// 		           what:"",
// 		           conv:nil,
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
// 		           what:"",
// 		           conv:nil,
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
// 			got := debugLogged(test.args.what, test.args.conv)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_discoverer_podEntries(t *testing.T) {
// 	type args struct {
// 		list *k8s.PodList
// 	}
// 	type fields struct {
// 		eg              errgroup.Group
// 		der             net.Dialer
// 		ctrl            k8s.Controller
// 		podsByName      atomic.Pointer[map[string][]*payload.Info_Pod]
// 		svcsByName      atomic.Pointer[map[string]*payload.Info_Service]
// 		nodeByName      atomic.Pointer[map[string]*payload.Info_Node]
// 		podsByNode      atomic.Pointer[map[string]map[string]map[string][]*payload.Info_Pod]
// 		podsByNamespace atomic.Pointer[map[string]map[string][]*payload.Info_Pod]
// 		namespace       string
// 		name            string
// 		pods            sync.Map[string, *[]Pod]
// 		nodeMetrics     sync.Map[string, NodeMetrics]
// 		nodes           sync.Map[string, *Node]
// 		podMetrics      sync.Map[string, PodMetrics]
// 		maxPods         int
// 		csd             time.Duration
// 	}
// 	type want struct {
// 		want map[string]*[]Pod
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, map[string]*[]Pod) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got map[string]*[]Pod) error {
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
// 		           list:nil,
// 		       },
// 		       fields: fields {
// 		           eg:nil,
// 		           der:nil,
// 		           ctrl:nil,
// 		           podsByName:nil,
// 		           svcsByName:nil,
// 		           nodeByName:nil,
// 		           podsByNode:nil,
// 		           podsByNamespace:nil,
// 		           namespace:"",
// 		           name:"",
// 		           pods:nil,
// 		           nodeMetrics:nil,
// 		           nodes:nil,
// 		           podMetrics:nil,
// 		           maxPods:0,
// 		           csd:nil,
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
// 		           list:nil,
// 		           },
// 		           fields: fields {
// 		           eg:nil,
// 		           der:nil,
// 		           ctrl:nil,
// 		           podsByName:nil,
// 		           svcsByName:nil,
// 		           nodeByName:nil,
// 		           podsByNode:nil,
// 		           podsByNamespace:nil,
// 		           namespace:"",
// 		           name:"",
// 		           pods:nil,
// 		           nodeMetrics:nil,
// 		           nodes:nil,
// 		           podMetrics:nil,
// 		           maxPods:0,
// 		           csd:nil,
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
// 			d := &discoverer{
// 				eg:              test.fields.eg,
// 				der:             test.fields.der,
// 				ctrl:            test.fields.ctrl,
// 				podsByName:      test.fields.podsByName,
// 				svcsByName:      test.fields.svcsByName,
// 				nodeByName:      test.fields.nodeByName,
// 				podsByNode:      test.fields.podsByNode,
// 				podsByNamespace: test.fields.podsByNamespace,
// 				namespace:       test.fields.namespace,
// 				name:            test.fields.name,
// 				pods:            test.fields.pods,
// 				nodeMetrics:     test.fields.nodeMetrics,
// 				nodes:           test.fields.nodes,
// 				podMetrics:      test.fields.podMetrics,
// 				maxPods:         test.fields.maxPods,
// 				csd:             test.fields.csd,
// 			}
//
// 			got := d.podEntries(test.args.list)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_nodeEntries(t *testing.T) {
// 	type args struct {
// 		list *k8s.NodeList
// 	}
// 	type want struct {
// 		want map[string]*Node
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, map[string]*Node) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got map[string]*Node) error {
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
// 		           list:nil,
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
// 		           list:nil,
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
// 			got := nodeEntries(test.args.list)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
