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

package resource

import (
	"testing"

	"github.com/vdaas/vald/internal/k8s/client"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/atomic"
	clientfake "k8s.io/client-go/kubernetes/fake"
)

const (
	testKindDeployment = "Deployment"
	testKindConfigMap  = "ConfigMap"
)

func newTestClients(t *testing.T) *Clients {
	t.Helper()
	return NewClients(&fakeClientSet{cs: clientfake.NewClientset()})
}

func TestClients_NamespacedKinds_CachePerNamespace(t *testing.T) {
	t.Parallel()

	tests := []struct {
		access func(cs *Clients, ns string) any
		name   string
	}{
		{name: "Pod", access: func(cs *Clients, ns string) any { return cs.Pod(ns) }},
		{name: testKindDeployment, access: func(cs *Clients, ns string) any { return cs.Deployment(ns) }},
		{name: "DaemonSet", access: func(cs *Clients, ns string) any { return cs.DaemonSet(ns) }},
		{name: "StatefulSet", access: func(cs *Clients, ns string) any { return cs.StatefulSet(ns) }},
		{name: "Job", access: func(cs *Clients, ns string) any { return cs.Job(ns) }},
		{name: "CronJob", access: func(cs *Clients, ns string) any { return cs.CronJob(ns) }},
		{name: "Service", access: func(cs *Clients, ns string) any { return cs.Service(ns) }},
		{name: "Secret", access: func(cs *Clients, ns string) any { return cs.Secret(ns) }},
		{name: testKindConfigMap, access: func(cs *Clients, ns string) any { return cs.ConfigMap(ns) }},
		{name: "Endpoints", access: func(cs *Clients, ns string) any { return cs.Endpoints(ns) }},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			cs := newTestClients(t)
			a := tc.access(cs, "default")
			b := tc.access(cs, "default")
			if a != b {
				t.Errorf("%s(%q) called twice returned different instances, want the same cached instance", tc.name, "default")
			}

			other := tc.access(cs, "other")
			if a == other {
				t.Errorf("%s(%q) and %s(%q) returned the same instance, want distinct per-namespace clients", tc.name, "default", tc.name, "other")
			}
		})
	}
}

func TestClients_ClusterScoped_CachesSingleton(t *testing.T) {
	t.Parallel()

	cs := newTestClients(t)
	a := cs.MutatingWebhookConfiguration()
	b := cs.MutatingWebhookConfiguration()
	if a != b {
		t.Error("MutatingWebhookConfiguration() called twice returned different instances, want the same cached instance")
	}

	va := cs.ValidatingWebhookConfiguration()
	vb := cs.ValidatingWebhookConfiguration()
	if va != vb {
		t.Error("ValidatingWebhookConfiguration() called twice returned different instances, want the same cached instance")
	}
}

// TestCached_ConcurrentConstructOnce exercises the shared double-checked-lock
// skeleton (cached) through both adapters under -race: every goroutine must
// observe the same instance and each construct must run exactly once.
func TestCached_ConcurrentConstructOnce(t *testing.T) {
	t.Parallel()

	cs := newTestClients(t)
	const goroutines = 32
	var (
		nsCalls, singletonCalls atomic.Int64
		wg                      sync.WaitGroup
	)
	nsResults := make([]PodClient, goroutines)
	singletonResults := make([]MutatingWebhookConfigurationClient, goroutines)
	for i := range goroutines {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			nsResults[i] = cachedNamespaced(cs, &cs.pods, "default",
				func(c client.Client, namespace string) PodClient {
					nsCalls.Add(1)
					return Pod(c, namespace)
				})
			singletonResults[i] = cachedSingleton(cs, &cs.mutatingWebhookConfiguration,
				func(c client.Client) MutatingWebhookConfigurationClient {
					singletonCalls.Add(1)
					return MutatingWebhookConfiguration(c)
				})
		}(i)
	}
	wg.Wait()

	if n := nsCalls.Load(); n != 1 {
		t.Errorf("namespaced construct ran %d times, want exactly once", n)
	}
	if n := singletonCalls.Load(); n != 1 {
		t.Errorf("singleton construct ran %d times, want exactly once", n)
	}
	for i := range goroutines {
		if nsResults[i] != nsResults[0] {
			t.Errorf("namespaced result[%d] differs from result[0], want the same cached instance", i)
		}
		if singletonResults[i] != singletonResults[0] {
			t.Errorf("singleton result[%d] differs from result[0], want the same cached instance", i)
		}
	}
}

func TestClients_ConcurrentAccess(t *testing.T) {
	t.Parallel()

	cs := newTestClients(t)
	var wg sync.WaitGroup
	results := make([]PodClient, 32)
	for i := range results {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			results[i] = cs.Pod("default")
		}(i)
	}
	wg.Wait()

	first := results[0]
	for i, r := range results {
		if r != first {
			t.Errorf("result[%d] = %v, want the same cached instance as result[0]", i, r)
		}
	}
}

// NOT IMPLEMENTED BELOW
//
// func TestNewClients(t *testing.T) {
// 	type args struct {
// 		c client.Client
// 	}
// 	type want struct {
// 		want *Clients
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, *Clients) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got *Clients) error {
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
// 		           c:nil,
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
// 		           c:nil,
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
// 			got := NewClients(test.args.c)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_cached(t *testing.T) {
// 	type args struct {
// 		mu     *sync.RWMutex
// 		lookup func() (T, bool)
// 		build  func() T
// 	}
// 	type want struct {
// 		want T
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, T) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got T) error {
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
// 		           mu:nil,
// 		           lookup:nil,
// 		           build:nil,
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
// 		           mu:nil,
// 		           lookup:nil,
// 		           build:nil,
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
// 			got := cached(test.args.mu, test.args.lookup, test.args.build)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_cachedNamespaced(t *testing.T) {
// 	type args struct {
// 		cs        *Clients
// 		cache     *map[string]T
// 		namespace string
// 		construct func(client.Client, string) T
// 	}
// 	type want struct {
// 		want T
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, T) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got T) error {
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
// 		           cs:Clients{},
// 		           cache:nil,
// 		           namespace:"",
// 		           construct:nil,
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
// 		           cs:Clients{},
// 		           cache:nil,
// 		           namespace:"",
// 		           construct:nil,
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
// 			got := cachedNamespaced(test.args.cs, test.args.cache, test.args.namespace, test.args.construct)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_cachedSingleton(t *testing.T) {
// 	type args struct {
// 		cs        *Clients
// 		cache     *T
// 		construct func(client.Client) T
// 	}
// 	type want struct {
// 		want T
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, T) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got T) error {
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
// 		           cs:Clients{},
// 		           cache:nil,
// 		           construct:nil,
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
// 		           cs:Clients{},
// 		           cache:nil,
// 		           construct:nil,
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
// 			got := cachedSingleton(test.args.cs, test.args.cache, test.args.construct)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestClients_Pod(t *testing.T) {
// 	type args struct {
// 		namespace string
// 	}
// 	type fields struct {
// 		c                              client.Client
// 		pods                           map[string]PodClient
// 		deployments                    map[string]DeploymentClient
// 		daemonSets                     map[string]DaemonSetClient
// 		statefulSets                   map[string]StatefulSetClient
// 		jobs                           map[string]JobClient
// 		cronJobs                       map[string]CronJobClient
// 		services                       map[string]ServiceClient
// 		secrets                        map[string]SecretClient
// 		configMaps                     map[string]ConfigMapClient
// 		endpoints                      map[string]EndpointClient
// 		mutatingWebhookConfiguration   MutatingWebhookConfigurationClient
// 		validatingWebhookConfiguration ValidatingWebhookConfigurationClient
// 	}
// 	type want struct {
// 		want PodClient
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, PodClient) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got PodClient) error {
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
// 		           namespace:"",
// 		       },
// 		       fields: fields {
// 		           c:nil,
// 		           pods:nil,
// 		           deployments:nil,
// 		           daemonSets:nil,
// 		           statefulSets:nil,
// 		           jobs:nil,
// 		           cronJobs:nil,
// 		           services:nil,
// 		           secrets:nil,
// 		           configMaps:nil,
// 		           endpoints:nil,
// 		           mutatingWebhookConfiguration:nil,
// 		           validatingWebhookConfiguration:nil,
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
// 		           namespace:"",
// 		           },
// 		           fields: fields {
// 		           c:nil,
// 		           pods:nil,
// 		           deployments:nil,
// 		           daemonSets:nil,
// 		           statefulSets:nil,
// 		           jobs:nil,
// 		           cronJobs:nil,
// 		           services:nil,
// 		           secrets:nil,
// 		           configMaps:nil,
// 		           endpoints:nil,
// 		           mutatingWebhookConfiguration:nil,
// 		           validatingWebhookConfiguration:nil,
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
// 			cs := &Clients{
// 				c:                              test.fields.c,
// 				pods:                           test.fields.pods,
// 				deployments:                    test.fields.deployments,
// 				daemonSets:                     test.fields.daemonSets,
// 				statefulSets:                   test.fields.statefulSets,
// 				jobs:                           test.fields.jobs,
// 				cronJobs:                       test.fields.cronJobs,
// 				services:                       test.fields.services,
// 				secrets:                        test.fields.secrets,
// 				configMaps:                     test.fields.configMaps,
// 				endpoints:                      test.fields.endpoints,
// 				mutatingWebhookConfiguration:   test.fields.mutatingWebhookConfiguration,
// 				validatingWebhookConfiguration: test.fields.validatingWebhookConfiguration,
// 			}
//
// 			got := cs.Pod(test.args.namespace)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestClients_Deployment(t *testing.T) {
// 	type args struct {
// 		namespace string
// 	}
// 	type fields struct {
// 		c                              client.Client
// 		pods                           map[string]PodClient
// 		deployments                    map[string]DeploymentClient
// 		daemonSets                     map[string]DaemonSetClient
// 		statefulSets                   map[string]StatefulSetClient
// 		jobs                           map[string]JobClient
// 		cronJobs                       map[string]CronJobClient
// 		services                       map[string]ServiceClient
// 		secrets                        map[string]SecretClient
// 		configMaps                     map[string]ConfigMapClient
// 		endpoints                      map[string]EndpointClient
// 		mutatingWebhookConfiguration   MutatingWebhookConfigurationClient
// 		validatingWebhookConfiguration ValidatingWebhookConfigurationClient
// 	}
// 	type want struct {
// 		want DeploymentClient
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, DeploymentClient) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got DeploymentClient) error {
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
// 		           namespace:"",
// 		       },
// 		       fields: fields {
// 		           c:nil,
// 		           pods:nil,
// 		           deployments:nil,
// 		           daemonSets:nil,
// 		           statefulSets:nil,
// 		           jobs:nil,
// 		           cronJobs:nil,
// 		           services:nil,
// 		           secrets:nil,
// 		           configMaps:nil,
// 		           endpoints:nil,
// 		           mutatingWebhookConfiguration:nil,
// 		           validatingWebhookConfiguration:nil,
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
// 		           namespace:"",
// 		           },
// 		           fields: fields {
// 		           c:nil,
// 		           pods:nil,
// 		           deployments:nil,
// 		           daemonSets:nil,
// 		           statefulSets:nil,
// 		           jobs:nil,
// 		           cronJobs:nil,
// 		           services:nil,
// 		           secrets:nil,
// 		           configMaps:nil,
// 		           endpoints:nil,
// 		           mutatingWebhookConfiguration:nil,
// 		           validatingWebhookConfiguration:nil,
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
// 			cs := &Clients{
// 				c:                              test.fields.c,
// 				pods:                           test.fields.pods,
// 				deployments:                    test.fields.deployments,
// 				daemonSets:                     test.fields.daemonSets,
// 				statefulSets:                   test.fields.statefulSets,
// 				jobs:                           test.fields.jobs,
// 				cronJobs:                       test.fields.cronJobs,
// 				services:                       test.fields.services,
// 				secrets:                        test.fields.secrets,
// 				configMaps:                     test.fields.configMaps,
// 				endpoints:                      test.fields.endpoints,
// 				mutatingWebhookConfiguration:   test.fields.mutatingWebhookConfiguration,
// 				validatingWebhookConfiguration: test.fields.validatingWebhookConfiguration,
// 			}
//
// 			got := cs.Deployment(test.args.namespace)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestClients_DaemonSet(t *testing.T) {
// 	type args struct {
// 		namespace string
// 	}
// 	type fields struct {
// 		c                              client.Client
// 		pods                           map[string]PodClient
// 		deployments                    map[string]DeploymentClient
// 		daemonSets                     map[string]DaemonSetClient
// 		statefulSets                   map[string]StatefulSetClient
// 		jobs                           map[string]JobClient
// 		cronJobs                       map[string]CronJobClient
// 		services                       map[string]ServiceClient
// 		secrets                        map[string]SecretClient
// 		configMaps                     map[string]ConfigMapClient
// 		endpoints                      map[string]EndpointClient
// 		mutatingWebhookConfiguration   MutatingWebhookConfigurationClient
// 		validatingWebhookConfiguration ValidatingWebhookConfigurationClient
// 	}
// 	type want struct {
// 		want DaemonSetClient
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, DaemonSetClient) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got DaemonSetClient) error {
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
// 		           namespace:"",
// 		       },
// 		       fields: fields {
// 		           c:nil,
// 		           pods:nil,
// 		           deployments:nil,
// 		           daemonSets:nil,
// 		           statefulSets:nil,
// 		           jobs:nil,
// 		           cronJobs:nil,
// 		           services:nil,
// 		           secrets:nil,
// 		           configMaps:nil,
// 		           endpoints:nil,
// 		           mutatingWebhookConfiguration:nil,
// 		           validatingWebhookConfiguration:nil,
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
// 		           namespace:"",
// 		           },
// 		           fields: fields {
// 		           c:nil,
// 		           pods:nil,
// 		           deployments:nil,
// 		           daemonSets:nil,
// 		           statefulSets:nil,
// 		           jobs:nil,
// 		           cronJobs:nil,
// 		           services:nil,
// 		           secrets:nil,
// 		           configMaps:nil,
// 		           endpoints:nil,
// 		           mutatingWebhookConfiguration:nil,
// 		           validatingWebhookConfiguration:nil,
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
// 			cs := &Clients{
// 				c:                              test.fields.c,
// 				pods:                           test.fields.pods,
// 				deployments:                    test.fields.deployments,
// 				daemonSets:                     test.fields.daemonSets,
// 				statefulSets:                   test.fields.statefulSets,
// 				jobs:                           test.fields.jobs,
// 				cronJobs:                       test.fields.cronJobs,
// 				services:                       test.fields.services,
// 				secrets:                        test.fields.secrets,
// 				configMaps:                     test.fields.configMaps,
// 				endpoints:                      test.fields.endpoints,
// 				mutatingWebhookConfiguration:   test.fields.mutatingWebhookConfiguration,
// 				validatingWebhookConfiguration: test.fields.validatingWebhookConfiguration,
// 			}
//
// 			got := cs.DaemonSet(test.args.namespace)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestClients_StatefulSet(t *testing.T) {
// 	type args struct {
// 		namespace string
// 	}
// 	type fields struct {
// 		c                              client.Client
// 		pods                           map[string]PodClient
// 		deployments                    map[string]DeploymentClient
// 		daemonSets                     map[string]DaemonSetClient
// 		statefulSets                   map[string]StatefulSetClient
// 		jobs                           map[string]JobClient
// 		cronJobs                       map[string]CronJobClient
// 		services                       map[string]ServiceClient
// 		secrets                        map[string]SecretClient
// 		configMaps                     map[string]ConfigMapClient
// 		endpoints                      map[string]EndpointClient
// 		mutatingWebhookConfiguration   MutatingWebhookConfigurationClient
// 		validatingWebhookConfiguration ValidatingWebhookConfigurationClient
// 	}
// 	type want struct {
// 		want StatefulSetClient
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, StatefulSetClient) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got StatefulSetClient) error {
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
// 		           namespace:"",
// 		       },
// 		       fields: fields {
// 		           c:nil,
// 		           pods:nil,
// 		           deployments:nil,
// 		           daemonSets:nil,
// 		           statefulSets:nil,
// 		           jobs:nil,
// 		           cronJobs:nil,
// 		           services:nil,
// 		           secrets:nil,
// 		           configMaps:nil,
// 		           endpoints:nil,
// 		           mutatingWebhookConfiguration:nil,
// 		           validatingWebhookConfiguration:nil,
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
// 		           namespace:"",
// 		           },
// 		           fields: fields {
// 		           c:nil,
// 		           pods:nil,
// 		           deployments:nil,
// 		           daemonSets:nil,
// 		           statefulSets:nil,
// 		           jobs:nil,
// 		           cronJobs:nil,
// 		           services:nil,
// 		           secrets:nil,
// 		           configMaps:nil,
// 		           endpoints:nil,
// 		           mutatingWebhookConfiguration:nil,
// 		           validatingWebhookConfiguration:nil,
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
// 			cs := &Clients{
// 				c:                              test.fields.c,
// 				pods:                           test.fields.pods,
// 				deployments:                    test.fields.deployments,
// 				daemonSets:                     test.fields.daemonSets,
// 				statefulSets:                   test.fields.statefulSets,
// 				jobs:                           test.fields.jobs,
// 				cronJobs:                       test.fields.cronJobs,
// 				services:                       test.fields.services,
// 				secrets:                        test.fields.secrets,
// 				configMaps:                     test.fields.configMaps,
// 				endpoints:                      test.fields.endpoints,
// 				mutatingWebhookConfiguration:   test.fields.mutatingWebhookConfiguration,
// 				validatingWebhookConfiguration: test.fields.validatingWebhookConfiguration,
// 			}
//
// 			got := cs.StatefulSet(test.args.namespace)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestClients_Job(t *testing.T) {
// 	type args struct {
// 		namespace string
// 	}
// 	type fields struct {
// 		c                              client.Client
// 		pods                           map[string]PodClient
// 		deployments                    map[string]DeploymentClient
// 		daemonSets                     map[string]DaemonSetClient
// 		statefulSets                   map[string]StatefulSetClient
// 		jobs                           map[string]JobClient
// 		cronJobs                       map[string]CronJobClient
// 		services                       map[string]ServiceClient
// 		secrets                        map[string]SecretClient
// 		configMaps                     map[string]ConfigMapClient
// 		endpoints                      map[string]EndpointClient
// 		mutatingWebhookConfiguration   MutatingWebhookConfigurationClient
// 		validatingWebhookConfiguration ValidatingWebhookConfigurationClient
// 	}
// 	type want struct {
// 		want JobClient
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, JobClient) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got JobClient) error {
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
// 		           namespace:"",
// 		       },
// 		       fields: fields {
// 		           c:nil,
// 		           pods:nil,
// 		           deployments:nil,
// 		           daemonSets:nil,
// 		           statefulSets:nil,
// 		           jobs:nil,
// 		           cronJobs:nil,
// 		           services:nil,
// 		           secrets:nil,
// 		           configMaps:nil,
// 		           endpoints:nil,
// 		           mutatingWebhookConfiguration:nil,
// 		           validatingWebhookConfiguration:nil,
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
// 		           namespace:"",
// 		           },
// 		           fields: fields {
// 		           c:nil,
// 		           pods:nil,
// 		           deployments:nil,
// 		           daemonSets:nil,
// 		           statefulSets:nil,
// 		           jobs:nil,
// 		           cronJobs:nil,
// 		           services:nil,
// 		           secrets:nil,
// 		           configMaps:nil,
// 		           endpoints:nil,
// 		           mutatingWebhookConfiguration:nil,
// 		           validatingWebhookConfiguration:nil,
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
// 			cs := &Clients{
// 				c:                              test.fields.c,
// 				pods:                           test.fields.pods,
// 				deployments:                    test.fields.deployments,
// 				daemonSets:                     test.fields.daemonSets,
// 				statefulSets:                   test.fields.statefulSets,
// 				jobs:                           test.fields.jobs,
// 				cronJobs:                       test.fields.cronJobs,
// 				services:                       test.fields.services,
// 				secrets:                        test.fields.secrets,
// 				configMaps:                     test.fields.configMaps,
// 				endpoints:                      test.fields.endpoints,
// 				mutatingWebhookConfiguration:   test.fields.mutatingWebhookConfiguration,
// 				validatingWebhookConfiguration: test.fields.validatingWebhookConfiguration,
// 			}
//
// 			got := cs.Job(test.args.namespace)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestClients_CronJob(t *testing.T) {
// 	type args struct {
// 		namespace string
// 	}
// 	type fields struct {
// 		c                              client.Client
// 		pods                           map[string]PodClient
// 		deployments                    map[string]DeploymentClient
// 		daemonSets                     map[string]DaemonSetClient
// 		statefulSets                   map[string]StatefulSetClient
// 		jobs                           map[string]JobClient
// 		cronJobs                       map[string]CronJobClient
// 		services                       map[string]ServiceClient
// 		secrets                        map[string]SecretClient
// 		configMaps                     map[string]ConfigMapClient
// 		endpoints                      map[string]EndpointClient
// 		mutatingWebhookConfiguration   MutatingWebhookConfigurationClient
// 		validatingWebhookConfiguration ValidatingWebhookConfigurationClient
// 	}
// 	type want struct {
// 		want CronJobClient
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, CronJobClient) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got CronJobClient) error {
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
// 		           namespace:"",
// 		       },
// 		       fields: fields {
// 		           c:nil,
// 		           pods:nil,
// 		           deployments:nil,
// 		           daemonSets:nil,
// 		           statefulSets:nil,
// 		           jobs:nil,
// 		           cronJobs:nil,
// 		           services:nil,
// 		           secrets:nil,
// 		           configMaps:nil,
// 		           endpoints:nil,
// 		           mutatingWebhookConfiguration:nil,
// 		           validatingWebhookConfiguration:nil,
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
// 		           namespace:"",
// 		           },
// 		           fields: fields {
// 		           c:nil,
// 		           pods:nil,
// 		           deployments:nil,
// 		           daemonSets:nil,
// 		           statefulSets:nil,
// 		           jobs:nil,
// 		           cronJobs:nil,
// 		           services:nil,
// 		           secrets:nil,
// 		           configMaps:nil,
// 		           endpoints:nil,
// 		           mutatingWebhookConfiguration:nil,
// 		           validatingWebhookConfiguration:nil,
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
// 			cs := &Clients{
// 				c:                              test.fields.c,
// 				pods:                           test.fields.pods,
// 				deployments:                    test.fields.deployments,
// 				daemonSets:                     test.fields.daemonSets,
// 				statefulSets:                   test.fields.statefulSets,
// 				jobs:                           test.fields.jobs,
// 				cronJobs:                       test.fields.cronJobs,
// 				services:                       test.fields.services,
// 				secrets:                        test.fields.secrets,
// 				configMaps:                     test.fields.configMaps,
// 				endpoints:                      test.fields.endpoints,
// 				mutatingWebhookConfiguration:   test.fields.mutatingWebhookConfiguration,
// 				validatingWebhookConfiguration: test.fields.validatingWebhookConfiguration,
// 			}
//
// 			got := cs.CronJob(test.args.namespace)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestClients_Service(t *testing.T) {
// 	type args struct {
// 		namespace string
// 	}
// 	type fields struct {
// 		c                              client.Client
// 		pods                           map[string]PodClient
// 		deployments                    map[string]DeploymentClient
// 		daemonSets                     map[string]DaemonSetClient
// 		statefulSets                   map[string]StatefulSetClient
// 		jobs                           map[string]JobClient
// 		cronJobs                       map[string]CronJobClient
// 		services                       map[string]ServiceClient
// 		secrets                        map[string]SecretClient
// 		configMaps                     map[string]ConfigMapClient
// 		endpoints                      map[string]EndpointClient
// 		mutatingWebhookConfiguration   MutatingWebhookConfigurationClient
// 		validatingWebhookConfiguration ValidatingWebhookConfigurationClient
// 	}
// 	type want struct {
// 		want ServiceClient
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, ServiceClient) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got ServiceClient) error {
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
// 		           namespace:"",
// 		       },
// 		       fields: fields {
// 		           c:nil,
// 		           pods:nil,
// 		           deployments:nil,
// 		           daemonSets:nil,
// 		           statefulSets:nil,
// 		           jobs:nil,
// 		           cronJobs:nil,
// 		           services:nil,
// 		           secrets:nil,
// 		           configMaps:nil,
// 		           endpoints:nil,
// 		           mutatingWebhookConfiguration:nil,
// 		           validatingWebhookConfiguration:nil,
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
// 		           namespace:"",
// 		           },
// 		           fields: fields {
// 		           c:nil,
// 		           pods:nil,
// 		           deployments:nil,
// 		           daemonSets:nil,
// 		           statefulSets:nil,
// 		           jobs:nil,
// 		           cronJobs:nil,
// 		           services:nil,
// 		           secrets:nil,
// 		           configMaps:nil,
// 		           endpoints:nil,
// 		           mutatingWebhookConfiguration:nil,
// 		           validatingWebhookConfiguration:nil,
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
// 			cs := &Clients{
// 				c:                              test.fields.c,
// 				pods:                           test.fields.pods,
// 				deployments:                    test.fields.deployments,
// 				daemonSets:                     test.fields.daemonSets,
// 				statefulSets:                   test.fields.statefulSets,
// 				jobs:                           test.fields.jobs,
// 				cronJobs:                       test.fields.cronJobs,
// 				services:                       test.fields.services,
// 				secrets:                        test.fields.secrets,
// 				configMaps:                     test.fields.configMaps,
// 				endpoints:                      test.fields.endpoints,
// 				mutatingWebhookConfiguration:   test.fields.mutatingWebhookConfiguration,
// 				validatingWebhookConfiguration: test.fields.validatingWebhookConfiguration,
// 			}
//
// 			got := cs.Service(test.args.namespace)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestClients_Secret(t *testing.T) {
// 	type args struct {
// 		namespace string
// 	}
// 	type fields struct {
// 		c                              client.Client
// 		pods                           map[string]PodClient
// 		deployments                    map[string]DeploymentClient
// 		daemonSets                     map[string]DaemonSetClient
// 		statefulSets                   map[string]StatefulSetClient
// 		jobs                           map[string]JobClient
// 		cronJobs                       map[string]CronJobClient
// 		services                       map[string]ServiceClient
// 		secrets                        map[string]SecretClient
// 		configMaps                     map[string]ConfigMapClient
// 		endpoints                      map[string]EndpointClient
// 		mutatingWebhookConfiguration   MutatingWebhookConfigurationClient
// 		validatingWebhookConfiguration ValidatingWebhookConfigurationClient
// 	}
// 	type want struct {
// 		want SecretClient
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, SecretClient) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got SecretClient) error {
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
// 		           namespace:"",
// 		       },
// 		       fields: fields {
// 		           c:nil,
// 		           pods:nil,
// 		           deployments:nil,
// 		           daemonSets:nil,
// 		           statefulSets:nil,
// 		           jobs:nil,
// 		           cronJobs:nil,
// 		           services:nil,
// 		           secrets:nil,
// 		           configMaps:nil,
// 		           endpoints:nil,
// 		           mutatingWebhookConfiguration:nil,
// 		           validatingWebhookConfiguration:nil,
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
// 		           namespace:"",
// 		           },
// 		           fields: fields {
// 		           c:nil,
// 		           pods:nil,
// 		           deployments:nil,
// 		           daemonSets:nil,
// 		           statefulSets:nil,
// 		           jobs:nil,
// 		           cronJobs:nil,
// 		           services:nil,
// 		           secrets:nil,
// 		           configMaps:nil,
// 		           endpoints:nil,
// 		           mutatingWebhookConfiguration:nil,
// 		           validatingWebhookConfiguration:nil,
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
// 			cs := &Clients{
// 				c:                              test.fields.c,
// 				pods:                           test.fields.pods,
// 				deployments:                    test.fields.deployments,
// 				daemonSets:                     test.fields.daemonSets,
// 				statefulSets:                   test.fields.statefulSets,
// 				jobs:                           test.fields.jobs,
// 				cronJobs:                       test.fields.cronJobs,
// 				services:                       test.fields.services,
// 				secrets:                        test.fields.secrets,
// 				configMaps:                     test.fields.configMaps,
// 				endpoints:                      test.fields.endpoints,
// 				mutatingWebhookConfiguration:   test.fields.mutatingWebhookConfiguration,
// 				validatingWebhookConfiguration: test.fields.validatingWebhookConfiguration,
// 			}
//
// 			got := cs.Secret(test.args.namespace)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestClients_ConfigMap(t *testing.T) {
// 	type args struct {
// 		namespace string
// 	}
// 	type fields struct {
// 		c                              client.Client
// 		pods                           map[string]PodClient
// 		deployments                    map[string]DeploymentClient
// 		daemonSets                     map[string]DaemonSetClient
// 		statefulSets                   map[string]StatefulSetClient
// 		jobs                           map[string]JobClient
// 		cronJobs                       map[string]CronJobClient
// 		services                       map[string]ServiceClient
// 		secrets                        map[string]SecretClient
// 		configMaps                     map[string]ConfigMapClient
// 		endpoints                      map[string]EndpointClient
// 		mutatingWebhookConfiguration   MutatingWebhookConfigurationClient
// 		validatingWebhookConfiguration ValidatingWebhookConfigurationClient
// 	}
// 	type want struct {
// 		want ConfigMapClient
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, ConfigMapClient) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got ConfigMapClient) error {
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
// 		           namespace:"",
// 		       },
// 		       fields: fields {
// 		           c:nil,
// 		           pods:nil,
// 		           deployments:nil,
// 		           daemonSets:nil,
// 		           statefulSets:nil,
// 		           jobs:nil,
// 		           cronJobs:nil,
// 		           services:nil,
// 		           secrets:nil,
// 		           configMaps:nil,
// 		           endpoints:nil,
// 		           mutatingWebhookConfiguration:nil,
// 		           validatingWebhookConfiguration:nil,
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
// 		           namespace:"",
// 		           },
// 		           fields: fields {
// 		           c:nil,
// 		           pods:nil,
// 		           deployments:nil,
// 		           daemonSets:nil,
// 		           statefulSets:nil,
// 		           jobs:nil,
// 		           cronJobs:nil,
// 		           services:nil,
// 		           secrets:nil,
// 		           configMaps:nil,
// 		           endpoints:nil,
// 		           mutatingWebhookConfiguration:nil,
// 		           validatingWebhookConfiguration:nil,
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
// 			cs := &Clients{
// 				c:                              test.fields.c,
// 				pods:                           test.fields.pods,
// 				deployments:                    test.fields.deployments,
// 				daemonSets:                     test.fields.daemonSets,
// 				statefulSets:                   test.fields.statefulSets,
// 				jobs:                           test.fields.jobs,
// 				cronJobs:                       test.fields.cronJobs,
// 				services:                       test.fields.services,
// 				secrets:                        test.fields.secrets,
// 				configMaps:                     test.fields.configMaps,
// 				endpoints:                      test.fields.endpoints,
// 				mutatingWebhookConfiguration:   test.fields.mutatingWebhookConfiguration,
// 				validatingWebhookConfiguration: test.fields.validatingWebhookConfiguration,
// 			}
//
// 			got := cs.ConfigMap(test.args.namespace)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestClients_Endpoints(t *testing.T) {
// 	type args struct {
// 		namespace string
// 	}
// 	type fields struct {
// 		c                              client.Client
// 		pods                           map[string]PodClient
// 		deployments                    map[string]DeploymentClient
// 		daemonSets                     map[string]DaemonSetClient
// 		statefulSets                   map[string]StatefulSetClient
// 		jobs                           map[string]JobClient
// 		cronJobs                       map[string]CronJobClient
// 		services                       map[string]ServiceClient
// 		secrets                        map[string]SecretClient
// 		configMaps                     map[string]ConfigMapClient
// 		endpoints                      map[string]EndpointClient
// 		mutatingWebhookConfiguration   MutatingWebhookConfigurationClient
// 		validatingWebhookConfiguration ValidatingWebhookConfigurationClient
// 	}
// 	type want struct {
// 		want EndpointClient
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, EndpointClient) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got EndpointClient) error {
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
// 		           namespace:"",
// 		       },
// 		       fields: fields {
// 		           c:nil,
// 		           pods:nil,
// 		           deployments:nil,
// 		           daemonSets:nil,
// 		           statefulSets:nil,
// 		           jobs:nil,
// 		           cronJobs:nil,
// 		           services:nil,
// 		           secrets:nil,
// 		           configMaps:nil,
// 		           endpoints:nil,
// 		           mutatingWebhookConfiguration:nil,
// 		           validatingWebhookConfiguration:nil,
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
// 		           namespace:"",
// 		           },
// 		           fields: fields {
// 		           c:nil,
// 		           pods:nil,
// 		           deployments:nil,
// 		           daemonSets:nil,
// 		           statefulSets:nil,
// 		           jobs:nil,
// 		           cronJobs:nil,
// 		           services:nil,
// 		           secrets:nil,
// 		           configMaps:nil,
// 		           endpoints:nil,
// 		           mutatingWebhookConfiguration:nil,
// 		           validatingWebhookConfiguration:nil,
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
// 			cs := &Clients{
// 				c:                              test.fields.c,
// 				pods:                           test.fields.pods,
// 				deployments:                    test.fields.deployments,
// 				daemonSets:                     test.fields.daemonSets,
// 				statefulSets:                   test.fields.statefulSets,
// 				jobs:                           test.fields.jobs,
// 				cronJobs:                       test.fields.cronJobs,
// 				services:                       test.fields.services,
// 				secrets:                        test.fields.secrets,
// 				configMaps:                     test.fields.configMaps,
// 				endpoints:                      test.fields.endpoints,
// 				mutatingWebhookConfiguration:   test.fields.mutatingWebhookConfiguration,
// 				validatingWebhookConfiguration: test.fields.validatingWebhookConfiguration,
// 			}
//
// 			got := cs.Endpoints(test.args.namespace)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestClients_MutatingWebhookConfiguration(t *testing.T) {
// 	type fields struct {
// 		c                              client.Client
// 		pods                           map[string]PodClient
// 		deployments                    map[string]DeploymentClient
// 		daemonSets                     map[string]DaemonSetClient
// 		statefulSets                   map[string]StatefulSetClient
// 		jobs                           map[string]JobClient
// 		cronJobs                       map[string]CronJobClient
// 		services                       map[string]ServiceClient
// 		secrets                        map[string]SecretClient
// 		configMaps                     map[string]ConfigMapClient
// 		endpoints                      map[string]EndpointClient
// 		mutatingWebhookConfiguration   MutatingWebhookConfigurationClient
// 		validatingWebhookConfiguration ValidatingWebhookConfigurationClient
// 	}
// 	type want struct {
// 		want MutatingWebhookConfigurationClient
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, MutatingWebhookConfigurationClient) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got MutatingWebhookConfigurationClient) error {
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
// 		       fields: fields {
// 		           c:nil,
// 		           pods:nil,
// 		           deployments:nil,
// 		           daemonSets:nil,
// 		           statefulSets:nil,
// 		           jobs:nil,
// 		           cronJobs:nil,
// 		           services:nil,
// 		           secrets:nil,
// 		           configMaps:nil,
// 		           endpoints:nil,
// 		           mutatingWebhookConfiguration:nil,
// 		           validatingWebhookConfiguration:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
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
// 		           fields: fields {
// 		           c:nil,
// 		           pods:nil,
// 		           deployments:nil,
// 		           daemonSets:nil,
// 		           statefulSets:nil,
// 		           jobs:nil,
// 		           cronJobs:nil,
// 		           services:nil,
// 		           secrets:nil,
// 		           configMaps:nil,
// 		           endpoints:nil,
// 		           mutatingWebhookConfiguration:nil,
// 		           validatingWebhookConfiguration:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
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
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			cs := &Clients{
// 				c:                              test.fields.c,
// 				pods:                           test.fields.pods,
// 				deployments:                    test.fields.deployments,
// 				daemonSets:                     test.fields.daemonSets,
// 				statefulSets:                   test.fields.statefulSets,
// 				jobs:                           test.fields.jobs,
// 				cronJobs:                       test.fields.cronJobs,
// 				services:                       test.fields.services,
// 				secrets:                        test.fields.secrets,
// 				configMaps:                     test.fields.configMaps,
// 				endpoints:                      test.fields.endpoints,
// 				mutatingWebhookConfiguration:   test.fields.mutatingWebhookConfiguration,
// 				validatingWebhookConfiguration: test.fields.validatingWebhookConfiguration,
// 			}
//
// 			got := cs.MutatingWebhookConfiguration()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestClients_ValidatingWebhookConfiguration(t *testing.T) {
// 	type fields struct {
// 		c                              client.Client
// 		pods                           map[string]PodClient
// 		deployments                    map[string]DeploymentClient
// 		daemonSets                     map[string]DaemonSetClient
// 		statefulSets                   map[string]StatefulSetClient
// 		jobs                           map[string]JobClient
// 		cronJobs                       map[string]CronJobClient
// 		services                       map[string]ServiceClient
// 		secrets                        map[string]SecretClient
// 		configMaps                     map[string]ConfigMapClient
// 		endpoints                      map[string]EndpointClient
// 		mutatingWebhookConfiguration   MutatingWebhookConfigurationClient
// 		validatingWebhookConfiguration ValidatingWebhookConfigurationClient
// 	}
// 	type want struct {
// 		want ValidatingWebhookConfigurationClient
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, ValidatingWebhookConfigurationClient) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got ValidatingWebhookConfigurationClient) error {
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
// 		       fields: fields {
// 		           c:nil,
// 		           pods:nil,
// 		           deployments:nil,
// 		           daemonSets:nil,
// 		           statefulSets:nil,
// 		           jobs:nil,
// 		           cronJobs:nil,
// 		           services:nil,
// 		           secrets:nil,
// 		           configMaps:nil,
// 		           endpoints:nil,
// 		           mutatingWebhookConfiguration:nil,
// 		           validatingWebhookConfiguration:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
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
// 		           fields: fields {
// 		           c:nil,
// 		           pods:nil,
// 		           deployments:nil,
// 		           daemonSets:nil,
// 		           statefulSets:nil,
// 		           jobs:nil,
// 		           cronJobs:nil,
// 		           services:nil,
// 		           secrets:nil,
// 		           configMaps:nil,
// 		           endpoints:nil,
// 		           mutatingWebhookConfiguration:nil,
// 		           validatingWebhookConfiguration:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
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
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			cs := &Clients{
// 				c:                              test.fields.c,
// 				pods:                           test.fields.pods,
// 				deployments:                    test.fields.deployments,
// 				daemonSets:                     test.fields.daemonSets,
// 				statefulSets:                   test.fields.statefulSets,
// 				jobs:                           test.fields.jobs,
// 				cronJobs:                       test.fields.cronJobs,
// 				services:                       test.fields.services,
// 				secrets:                        test.fields.secrets,
// 				configMaps:                     test.fields.configMaps,
// 				endpoints:                      test.fields.endpoints,
// 				mutatingWebhookConfiguration:   test.fields.mutatingWebhookConfiguration,
// 				validatingWebhookConfiguration: test.fields.validatingWebhookConfiguration,
// 			}
//
// 			got := cs.ValidatingWebhookConfiguration()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
