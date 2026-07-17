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

package portforward

import (
	"context"
	"net/http"
	"slices"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/test/goleak"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/rest"
	crclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

// mockClientSet wires the client-go fake clientset into the kclient.Client
// interface used by the forwarder. The forwarder only ever calls
// GetClientSet and GetRESTConfig on it, so the embedded crclient.Client is
// left nil: it exists solely to satisfy client.Client's k8s.Client method set
// at compile time.
type mockClientSet struct {
	crclient.Client
	cs kubernetes.Interface
}

func (m *mockClientSet) GetClientSet() kubernetes.Interface {
	return m.cs
}

func (*mockClientSet) GetRESTConfig() *rest.Config {
	return &rest.Config{Host: "https://127.0.0.1:6443"}
}

func newTestForwarder(t *testing.T) Forwarder {
	t.Helper()
	pf, err := New(
		WithClient(&mockClientSet{cs: fake.NewClientset()}),
		WithNamespace("default"),
		WithServiceName("vald-lb-gateway"),
		WithAddress("127.0.0.1"),
		WithPorts(map[uint16]uint16{18081: 8081}),
	)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	return pf
}

// TestPortforwardExtended_DoesNotMutateSharedClient guards against the SPDY
// transport being assigned to the caller-provided (possibly process-shared)
// http.Client: the session must build its own client instead.
func TestPortforwardExtended_DoesNotMutateSharedClient(t *testing.T) {
	t.Parallel()

	shared := &http.Client{Timeout: 3 * time.Second}

	cancel, _, err := portforwardExtended(
		context.Background(),
		nil, // errgroup unused: the call fails on address validation before spawning
		&mockClientSet{cs: fake.NewClientset()},
		"default", "vald-agent-0",
		nil, // no addresses so the call fails before dialing
		map[uint16]uint16{18081: 8081},
		shared,
	)
	if cancel != nil {
		cancel()
	}
	if !errors.Is(err, errors.ErrPortForwardAddressNotFound) {
		t.Fatalf("portforwardExtended() error = %v, want %v", err, errors.ErrPortForwardAddressNotFound)
	}
	if shared.Transport != nil {
		t.Errorf("shared http.Client transport mutated to %T, want nil", shared.Transport)
	}
	if shared.Timeout != 3*time.Second {
		t.Errorf("shared http.Client timeout mutated to %v", shared.Timeout)
	}
}

// TestNormalize verifies the shared slice normalization helper used by both
// updateTargets and portforwardExtended: the result must be sorted, fully
// deduplicated (including non-adjacent duplicates), and clipped to length.
func TestNormalize(t *testing.T) {
	t.Parallel()

	const (
		podA = "pod-a"
		podB = "pod-b"
		podC = "pod-c"
	)
	tests := []struct {
		name string
		in   []string
		want []string
	}{
		{
			name: "nil input stays nil",
			in:   nil,
			want: nil,
		},
		{
			name: "empty input stays empty",
			in:   []string{},
			want: []string{},
		},
		{
			name: "unsorted input is sorted",
			in:   []string{podC, podA, podB},
			want: []string{podA, podB, podC},
		},
		{
			name: "adjacent duplicates are removed",
			in:   []string{podA, podA, podB},
			want: []string{podA, podB},
		},
		{
			name: "non-adjacent duplicates are removed",
			in:   []string{podB, podA, podB, podA},
			want: []string{podA, podB},
		},
		{
			name: "excess capacity is clipped",
			in:   append(make([]string, 0, 16), podB, podA, podA),
			want: []string{podA, podB},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got := normalize(test.in)
			if !slices.Equal(got, test.want) {
				t.Errorf("normalize() = %v, want %v", got, test.want)
			}
			if cap(got) != len(got) {
				t.Errorf("normalize() cap = %d, want %d (clipped to length)", cap(got), len(got))
			}
		})
	}
}

// TestPortForward_StopBeforeStart guards Stop against nil panics when Start
// was never called.
func TestPortForward_StopBeforeStart(t *testing.T) {
	t.Parallel()

	pf := newTestForwarder(t)
	if err := pf.Stop(); err != nil {
		t.Errorf("Stop() before Start error = %v, want nil", err)
	}
	if err := pf.Stop(); err != nil {
		t.Errorf("second Stop() error = %v, want nil", err)
	}
}

// TestPortForward_StartStopIdempotent verifies that Start returns once the
// context expires (no busy-wait hang) and that Stop is idempotent afterwards
// (no double close of the error channel).
func TestPortForward_StartStopIdempotent(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	pf := newTestForwarder(t)
	ech, err := pf.Start(ctx)
	if ech == nil {
		t.Fatal("Start() ech = nil, want non-nil channel")
	}
	if err == nil {
		t.Error("Start() error = nil, want context error because no pod ever becomes healthy")
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		// First Stop drains the goroutines and may surface the context
		// error from the daemon goroutines; the second must be a no-op.
		if serr := pf.Stop(); serr != nil {
			t.Logf("first Stop() error = %v", serr)
		}
		if serr := pf.Stop(); serr != nil {
			t.Errorf("second Stop() error = %v, want nil", serr)
		}
	}()
	select {
	case <-done:
	case <-time.After(30 * time.Second):
		t.Fatal("Stop() did not return: goroutines leaked or deadlocked")
	}
}

// NOT IMPLEMENTED BELOW
//
// func TestNew(t *testing.T) {
// 	type args struct {
// 		opts []Option
// 	}
// 	type want struct {
// 		want Forwarder
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, Forwarder, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got Forwarder, err error) error {
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
// 		           opts:nil,
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
// 		           opts:nil,
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
// 			got, err := New(test.args.opts...)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_normalize(t *testing.T) {
// 	type args struct {
// 		ss []string
// 	}
// 	type want struct {
// 		want []string
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, []string) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got []string) error {
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
// 		           ss:nil,
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
// 		           ss:nil,
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
// 			got := normalize(test.args.ss)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_portForward_updateTargets(t *testing.T) {
// 	type args struct {
// 		pods []string
// 	}
// 	type fields struct {
// 		client      kclient.Client
// 		eclient     k8s.EndpointClient
// 		backoff     backoff.Backoff
// 		eg          errgroup.Group
// 		ports       map[uint16]uint16
// 		httpClient  *http.Client
// 		cancel      context.CancelFunc
// 		ech         chan error
// 		namespace   string
// 		serviceName string
// 		addresses   []string
// 		targets     []string
// 		healthy     atomic.Bool
// 		current     atomic.Uint64
// 	}
// 	type want struct{}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want) error {
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           pods:nil,
// 		       },
// 		       fields: fields {
// 		           client:nil,
// 		           eclient:nil,
// 		           backoff:nil,
// 		           eg:nil,
// 		           ports:nil,
// 		           httpClient:nil,
// 		           cancel:nil,
// 		           ech:nil,
// 		           namespace:"",
// 		           serviceName:"",
// 		           addresses:nil,
// 		           targets:nil,
// 		           healthy:nil,
// 		           current:nil,
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
// 		           pods:nil,
// 		           },
// 		           fields: fields {
// 		           client:nil,
// 		           eclient:nil,
// 		           backoff:nil,
// 		           eg:nil,
// 		           ports:nil,
// 		           httpClient:nil,
// 		           cancel:nil,
// 		           ech:nil,
// 		           namespace:"",
// 		           serviceName:"",
// 		           addresses:nil,
// 		           targets:nil,
// 		           healthy:nil,
// 		           current:nil,
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
// 			pf := &portForward{
// 				client:      test.fields.client,
// 				eclient:     test.fields.eclient,
// 				backoff:     test.fields.backoff,
// 				eg:          test.fields.eg,
// 				ports:       test.fields.ports,
// 				httpClient:  test.fields.httpClient,
// 				cancel:      test.fields.cancel,
// 				ech:         test.fields.ech,
// 				namespace:   test.fields.namespace,
// 				serviceName: test.fields.serviceName,
// 				addresses:   test.fields.addresses,
// 				targets:     test.fields.targets,
// 				healthy:     test.fields.healthy,
// 				current:     test.fields.current,
// 			}
//
// 			pf.updateTargets(test.args.pods)
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_portForward_getNextPod(t *testing.T) {
// 	type fields struct {
// 		client      kclient.Client
// 		eclient     k8s.EndpointClient
// 		backoff     backoff.Backoff
// 		eg          errgroup.Group
// 		ports       map[uint16]uint16
// 		httpClient  *http.Client
// 		cancel      context.CancelFunc
// 		ech         chan error
// 		namespace   string
// 		serviceName string
// 		addresses   []string
// 		targets     []string
// 		healthy     atomic.Bool
// 		current     atomic.Uint64
// 	}
// 	type want struct {
// 		want string
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, string, error) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got string, err error) error {
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
// 		       fields: fields {
// 		           client:nil,
// 		           eclient:nil,
// 		           backoff:nil,
// 		           eg:nil,
// 		           ports:nil,
// 		           httpClient:nil,
// 		           cancel:nil,
// 		           ech:nil,
// 		           namespace:"",
// 		           serviceName:"",
// 		           addresses:nil,
// 		           targets:nil,
// 		           healthy:nil,
// 		           current:nil,
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
// 		           client:nil,
// 		           eclient:nil,
// 		           backoff:nil,
// 		           eg:nil,
// 		           ports:nil,
// 		           httpClient:nil,
// 		           cancel:nil,
// 		           ech:nil,
// 		           namespace:"",
// 		           serviceName:"",
// 		           addresses:nil,
// 		           targets:nil,
// 		           healthy:nil,
// 		           current:nil,
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
// 			pf := &portForward{
// 				client:      test.fields.client,
// 				eclient:     test.fields.eclient,
// 				backoff:     test.fields.backoff,
// 				eg:          test.fields.eg,
// 				ports:       test.fields.ports,
// 				httpClient:  test.fields.httpClient,
// 				cancel:      test.fields.cancel,
// 				ech:         test.fields.ech,
// 				namespace:   test.fields.namespace,
// 				serviceName: test.fields.serviceName,
// 				addresses:   test.fields.addresses,
// 				targets:     test.fields.targets,
// 				healthy:     test.fields.healthy,
// 				current:     test.fields.current,
// 			}
//
// 			got, err := pf.getNextPod()
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_portForward_Start(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		client      kclient.Client
// 		eclient     k8s.EndpointClient
// 		backoff     backoff.Backoff
// 		eg          errgroup.Group
// 		ports       map[uint16]uint16
// 		httpClient  *http.Client
// 		cancel      context.CancelFunc
// 		ech         chan error
// 		namespace   string
// 		serviceName string
// 		addresses   []string
// 		targets     []string
// 		healthy     atomic.Bool
// 		current     atomic.Uint64
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
// 		           client:nil,
// 		           eclient:nil,
// 		           backoff:nil,
// 		           eg:nil,
// 		           ports:nil,
// 		           httpClient:nil,
// 		           cancel:nil,
// 		           ech:nil,
// 		           namespace:"",
// 		           serviceName:"",
// 		           addresses:nil,
// 		           targets:nil,
// 		           healthy:nil,
// 		           current:nil,
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
// 		           client:nil,
// 		           eclient:nil,
// 		           backoff:nil,
// 		           eg:nil,
// 		           ports:nil,
// 		           httpClient:nil,
// 		           cancel:nil,
// 		           ech:nil,
// 		           namespace:"",
// 		           serviceName:"",
// 		           addresses:nil,
// 		           targets:nil,
// 		           healthy:nil,
// 		           current:nil,
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
// 			pf := &portForward{
// 				client:      test.fields.client,
// 				eclient:     test.fields.eclient,
// 				backoff:     test.fields.backoff,
// 				eg:          test.fields.eg,
// 				ports:       test.fields.ports,
// 				httpClient:  test.fields.httpClient,
// 				cancel:      test.fields.cancel,
// 				ech:         test.fields.ech,
// 				namespace:   test.fields.namespace,
// 				serviceName: test.fields.serviceName,
// 				addresses:   test.fields.addresses,
// 				targets:     test.fields.targets,
// 				healthy:     test.fields.healthy,
// 				current:     test.fields.current,
// 			}
//
// 			got, err := pf.Start(test.args.ctx)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_portForward_notifyError(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		err error
// 	}
// 	type fields struct {
// 		client      kclient.Client
// 		eclient     k8s.EndpointClient
// 		backoff     backoff.Backoff
// 		eg          errgroup.Group
// 		ports       map[uint16]uint16
// 		httpClient  *http.Client
// 		cancel      context.CancelFunc
// 		ech         chan error
// 		namespace   string
// 		serviceName string
// 		addresses   []string
// 		targets     []string
// 		healthy     atomic.Bool
// 		current     atomic.Uint64
// 	}
// 	type want struct{}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want) error {
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           err:nil,
// 		       },
// 		       fields: fields {
// 		           client:nil,
// 		           eclient:nil,
// 		           backoff:nil,
// 		           eg:nil,
// 		           ports:nil,
// 		           httpClient:nil,
// 		           cancel:nil,
// 		           ech:nil,
// 		           namespace:"",
// 		           serviceName:"",
// 		           addresses:nil,
// 		           targets:nil,
// 		           healthy:nil,
// 		           current:nil,
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
// 		           err:nil,
// 		           },
// 		           fields: fields {
// 		           client:nil,
// 		           eclient:nil,
// 		           backoff:nil,
// 		           eg:nil,
// 		           ports:nil,
// 		           httpClient:nil,
// 		           cancel:nil,
// 		           ech:nil,
// 		           namespace:"",
// 		           serviceName:"",
// 		           addresses:nil,
// 		           targets:nil,
// 		           healthy:nil,
// 		           current:nil,
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
// 			pf := &portForward{
// 				client:      test.fields.client,
// 				eclient:     test.fields.eclient,
// 				backoff:     test.fields.backoff,
// 				eg:          test.fields.eg,
// 				ports:       test.fields.ports,
// 				httpClient:  test.fields.httpClient,
// 				cancel:      test.fields.cancel,
// 				ech:         test.fields.ech,
// 				namespace:   test.fields.namespace,
// 				serviceName: test.fields.serviceName,
// 				addresses:   test.fields.addresses,
// 				targets:     test.fields.targets,
// 				healthy:     test.fields.healthy,
// 				current:     test.fields.current,
// 			}
//
// 			pf.notifyError(test.args.ctx, test.args.err)
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_portForward_Stop(t *testing.T) {
// 	type fields struct {
// 		client      kclient.Client
// 		eclient     k8s.EndpointClient
// 		backoff     backoff.Backoff
// 		eg          errgroup.Group
// 		ports       map[uint16]uint16
// 		httpClient  *http.Client
// 		cancel      context.CancelFunc
// 		ech         chan error
// 		namespace   string
// 		serviceName string
// 		addresses   []string
// 		targets     []string
// 		healthy     atomic.Bool
// 		current     atomic.Uint64
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
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
// 		       fields: fields {
// 		           client:nil,
// 		           eclient:nil,
// 		           backoff:nil,
// 		           eg:nil,
// 		           ports:nil,
// 		           httpClient:nil,
// 		           cancel:nil,
// 		           ech:nil,
// 		           namespace:"",
// 		           serviceName:"",
// 		           addresses:nil,
// 		           targets:nil,
// 		           healthy:nil,
// 		           current:nil,
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
// 		           client:nil,
// 		           eclient:nil,
// 		           backoff:nil,
// 		           eg:nil,
// 		           ports:nil,
// 		           httpClient:nil,
// 		           cancel:nil,
// 		           ech:nil,
// 		           namespace:"",
// 		           serviceName:"",
// 		           addresses:nil,
// 		           targets:nil,
// 		           healthy:nil,
// 		           current:nil,
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
// 			pf := &portForward{
// 				client:      test.fields.client,
// 				eclient:     test.fields.eclient,
// 				backoff:     test.fields.backoff,
// 				eg:          test.fields.eg,
// 				ports:       test.fields.ports,
// 				httpClient:  test.fields.httpClient,
// 				cancel:      test.fields.cancel,
// 				ech:         test.fields.ech,
// 				namespace:   test.fields.namespace,
// 				serviceName: test.fields.serviceName,
// 				addresses:   test.fields.addresses,
// 				targets:     test.fields.targets,
// 				healthy:     test.fields.healthy,
// 				current:     test.fields.current,
// 			}
//
// 			err := pf.Stop()
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_portForward_endpointsWatcher(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		client      kclient.Client
// 		eclient     k8s.EndpointClient
// 		backoff     backoff.Backoff
// 		eg          errgroup.Group
// 		ports       map[uint16]uint16
// 		httpClient  *http.Client
// 		cancel      context.CancelFunc
// 		ech         chan error
// 		namespace   string
// 		serviceName string
// 		addresses   []string
// 		targets     []string
// 		healthy     atomic.Bool
// 		current     atomic.Uint64
// 	}
// 	type want struct {
// 		wantW watch.Interface
// 		err   error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, watch.Interface, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotW watch.Interface, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotW, w.wantW) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotW, w.wantW)
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
// 		           client:nil,
// 		           eclient:nil,
// 		           backoff:nil,
// 		           eg:nil,
// 		           ports:nil,
// 		           httpClient:nil,
// 		           cancel:nil,
// 		           ech:nil,
// 		           namespace:"",
// 		           serviceName:"",
// 		           addresses:nil,
// 		           targets:nil,
// 		           healthy:nil,
// 		           current:nil,
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
// 		           client:nil,
// 		           eclient:nil,
// 		           backoff:nil,
// 		           eg:nil,
// 		           ports:nil,
// 		           httpClient:nil,
// 		           cancel:nil,
// 		           ech:nil,
// 		           namespace:"",
// 		           serviceName:"",
// 		           addresses:nil,
// 		           targets:nil,
// 		           healthy:nil,
// 		           current:nil,
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
// 			pf := &portForward{
// 				client:      test.fields.client,
// 				eclient:     test.fields.eclient,
// 				backoff:     test.fields.backoff,
// 				eg:          test.fields.eg,
// 				ports:       test.fields.ports,
// 				httpClient:  test.fields.httpClient,
// 				cancel:      test.fields.cancel,
// 				ech:         test.fields.ech,
// 				namespace:   test.fields.namespace,
// 				serviceName: test.fields.serviceName,
// 				addresses:   test.fields.addresses,
// 				targets:     test.fields.targets,
// 				healthy:     test.fields.healthy,
// 				current:     test.fields.current,
// 			}
//
// 			gotW, err := pf.endpointsWatcher(test.args.ctx)
// 			if err := checkFunc(test.want, gotW, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_portForward_loadTargets(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		client      kclient.Client
// 		eclient     k8s.EndpointClient
// 		backoff     backoff.Backoff
// 		eg          errgroup.Group
// 		ports       map[uint16]uint16
// 		httpClient  *http.Client
// 		cancel      context.CancelFunc
// 		ech         chan error
// 		namespace   string
// 		serviceName string
// 		addresses   []string
// 		targets     []string
// 		healthy     atomic.Bool
// 		current     atomic.Uint64
// 	}
// 	type want struct{}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want) error {
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
// 		           client:nil,
// 		           eclient:nil,
// 		           backoff:nil,
// 		           eg:nil,
// 		           ports:nil,
// 		           httpClient:nil,
// 		           cancel:nil,
// 		           ech:nil,
// 		           namespace:"",
// 		           serviceName:"",
// 		           addresses:nil,
// 		           targets:nil,
// 		           healthy:nil,
// 		           current:nil,
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
// 		           client:nil,
// 		           eclient:nil,
// 		           backoff:nil,
// 		           eg:nil,
// 		           ports:nil,
// 		           httpClient:nil,
// 		           cancel:nil,
// 		           ech:nil,
// 		           namespace:"",
// 		           serviceName:"",
// 		           addresses:nil,
// 		           targets:nil,
// 		           healthy:nil,
// 		           current:nil,
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
// 			pf := &portForward{
// 				client:      test.fields.client,
// 				eclient:     test.fields.eclient,
// 				backoff:     test.fields.backoff,
// 				eg:          test.fields.eg,
// 				ports:       test.fields.ports,
// 				httpClient:  test.fields.httpClient,
// 				cancel:      test.fields.cancel,
// 				ech:         test.fields.ech,
// 				namespace:   test.fields.namespace,
// 				serviceName: test.fields.serviceName,
// 				addresses:   test.fields.addresses,
// 				targets:     test.fields.targets,
// 				healthy:     test.fields.healthy,
// 				current:     test.fields.current,
// 			}
//
// 			pf.loadTargets(test.args.ctx)
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_portForward_portForwardToService(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		client      kclient.Client
// 		eclient     k8s.EndpointClient
// 		backoff     backoff.Backoff
// 		eg          errgroup.Group
// 		ports       map[uint16]uint16
// 		httpClient  *http.Client
// 		cancel      context.CancelFunc
// 		ech         chan error
// 		namespace   string
// 		serviceName string
// 		addresses   []string
// 		targets     []string
// 		healthy     atomic.Bool
// 		current     atomic.Uint64
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
// 		           client:nil,
// 		           eclient:nil,
// 		           backoff:nil,
// 		           eg:nil,
// 		           ports:nil,
// 		           httpClient:nil,
// 		           cancel:nil,
// 		           ech:nil,
// 		           namespace:"",
// 		           serviceName:"",
// 		           addresses:nil,
// 		           targets:nil,
// 		           healthy:nil,
// 		           current:nil,
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
// 		           client:nil,
// 		           eclient:nil,
// 		           backoff:nil,
// 		           eg:nil,
// 		           ports:nil,
// 		           httpClient:nil,
// 		           cancel:nil,
// 		           ech:nil,
// 		           namespace:"",
// 		           serviceName:"",
// 		           addresses:nil,
// 		           targets:nil,
// 		           healthy:nil,
// 		           current:nil,
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
// 			pf := &portForward{
// 				client:      test.fields.client,
// 				eclient:     test.fields.eclient,
// 				backoff:     test.fields.backoff,
// 				eg:          test.fields.eg,
// 				ports:       test.fields.ports,
// 				httpClient:  test.fields.httpClient,
// 				cancel:      test.fields.cancel,
// 				ech:         test.fields.ech,
// 				namespace:   test.fields.namespace,
// 				serviceName: test.fields.serviceName,
// 				addresses:   test.fields.addresses,
// 				targets:     test.fields.targets,
// 				healthy:     test.fields.healthy,
// 				current:     test.fields.current,
// 			}
//
// 			err := pf.portForwardToService(test.args.ctx)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_portforwardExtended(t *testing.T) {
// 	type args struct {
// 		ctx       context.Context
// 		eg        errgroup.Group
// 		c         kclient.Client
// 		namespace string
// 		podName   string
// 		addresses []string
// 		ports     map[uint16]uint16
// 		hc        *http.Client
// 	}
// 	type want struct {
// 		wantCancel    context.CancelFunc
// 		wantErrorChan <-chan error
// 		err           error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, context.CancelFunc, <-chan error, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotCancel context.CancelFunc, gotErrorChan <-chan error, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotCancel, w.wantCancel) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotCancel, w.wantCancel)
// 		}
// 		if !reflect.DeepEqual(gotErrorChan, w.wantErrorChan) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotErrorChan, w.wantErrorChan)
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
// 		           eg:nil,
// 		           c:nil,
// 		           namespace:"",
// 		           podName:"",
// 		           addresses:nil,
// 		           ports:nil,
// 		           hc:nil,
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
// 		           eg:nil,
// 		           c:nil,
// 		           namespace:"",
// 		           podName:"",
// 		           addresses:nil,
// 		           ports:nil,
// 		           hc:nil,
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
// 			gotCancel, gotErrorChan, err := portforwardExtended(test.args.ctx, test.args.eg, test.args.c, test.args.namespace, test.args.podName, test.args.addresses, test.args.ports, test.args.hc)
// 			if err := checkFunc(test.want, gotCancel, gotErrorChan, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
