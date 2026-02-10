// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
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
	"reflect"
	"testing"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/k8s"
	"github.com/vdaas/vald/internal/k8s/vald/mirror/target"
	"github.com/vdaas/vald/internal/test/goleak"
	k8smock "github.com/vdaas/vald/internal/test/mock/k8s"
)

func Test_discovery_startSync(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx  context.Context
		prev map[string]target.Target
	}
	type fields struct {
		ctrl k8s.Controller
		mirr Mirror
	}
	type want struct {
		want map[string]target.Target
		err  error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, map[string]target.Target, error) error
		beforeFunc func(*testing.T, *discovery, args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func(w want, got map[string]target.Target, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			prev := make(map[string]target.Target)
			current := map[string]target.Target{
				"mirror-1": {
					Host: "192.168.1.2",
					Port: 8081,
				},
			}
			return test{
				name: "Succeeded detecting the created resource and connecting to the new target",
				args: args{
					ctx:  context.Background(),
					prev: prev,
				},
				fields: fields{
					mirr: &MirrorMock{
						ConnectFunc: func(_ context.Context, _ ...*payload.Mirror_Target) error {
							return nil
						},
					},
					ctrl: &k8smock.ControllerMock{
						GetManagerFunc: k8smock.NewDefaultManagerMock,
					},
				},
				beforeFunc: func(t *testing.T, d *discovery, _ args) {
					t.Helper()
					d.onReconcile(context.Background(), current)
				},
				want: want{
					want: current,
				},
			}
		}(),
		func() test {
			prev := make(map[string]target.Target)
			current := map[string]target.Target{
				"mirror-1": {
					Host: "192.168.1.2",
					Port: 8081,
				},
			}
			return test{
				name: "Succeeded detecting the created resource but failed to connect to the new target",
				args: args{
					ctx:  context.Background(),
					prev: prev,
				},
				fields: fields{
					mirr: &MirrorMock{
						ConnectFunc: func(_ context.Context, _ ...*payload.Mirror_Target) error {
							return errors.New("err")
						},
					},
					ctrl: &k8smock.ControllerMock{
						GetManagerFunc: k8smock.NewDefaultManagerMock,
					},
				},
				beforeFunc: func(t *testing.T, d *discovery, _ args) {
					t.Helper()
					d.onReconcile(context.Background(), current)
				},
				want: want{
					want: current,
					err:  errors.New("err"),
				},
			}
		}(),
		func() test {
			prev := map[string]target.Target{
				"mirror-1": {
					Host: "192.168.1.2",
					Port: 8081,
				},
			}
			current := make(map[string]target.Target)
			return test{
				name: "Succeeded detecting the deleted resource and disconnecting the target",
				args: args{
					ctx:  context.Background(),
					prev: prev,
				},
				fields: fields{
					mirr: &MirrorMock{
						DisconnectFunc: func(_ context.Context, _ ...*payload.Mirror_Target) error {
							return nil
						},
					},
					ctrl: &k8smock.ControllerMock{
						GetManagerFunc: k8smock.NewDefaultManagerMock,
					},
				},
				beforeFunc: func(t *testing.T, d *discovery, _ args) {
					t.Helper()
					d.onReconcile(context.Background(), current)
				},
				want: want{
					want: current,
				},
			}
		}(),
		func() test {
			prev := map[string]target.Target{
				"mirror-1": {
					Host: "192.168.1.2",
					Port: 8081,
				},
			}
			current := make(map[string]target.Target)
			return test{
				name: "Succeeded detecting the deleted resource but failed to disconnect the target",
				args: args{
					ctx:  context.Background(),
					prev: prev,
				},
				fields: fields{
					mirr: &MirrorMock{
						DisconnectFunc: func(_ context.Context, _ ...*payload.Mirror_Target) error {
							return errors.New("err")
						},
					},
					ctrl: &k8smock.ControllerMock{
						GetManagerFunc: k8smock.NewDefaultManagerMock,
					},
				},
				beforeFunc: func(t *testing.T, d *discovery, _ args) {
					t.Helper()
					d.onReconcile(context.Background(), current)
				},
				want: want{
					want: current,
					err:  errors.New("err"),
				},
			}
		}(),
		func() test {
			prev := map[string]target.Target{
				"mirror-1": {
					Host: "192.168.1.2",
					Port: 8081,
				},
			}
			current := map[string]target.Target{
				"mirror-1": {
					Host: "192.168.1.3",
					Port: 8081,
				},
			}
			return test{
				name: "Succeeded detecting the updated resource and updating the target connection",
				args: args{
					ctx:  context.Background(),
					prev: prev,
				},
				fields: fields{
					mirr: &MirrorMock{
						ConnectFunc: func(_ context.Context, _ ...*payload.Mirror_Target) error {
							return nil
						},
						DisconnectFunc: func(_ context.Context, _ ...*payload.Mirror_Target) error {
							return nil
						},
					},
					ctrl: &k8smock.ControllerMock{
						GetManagerFunc: k8smock.NewDefaultManagerMock,
					},
				},
				beforeFunc: func(t *testing.T, d *discovery, _ args) {
					t.Helper()
					d.onReconcile(context.Background(), current)
				},
				want: want{
					want: current,
				},
			}
		}(),
		func() test {
			prev := map[string]target.Target{
				"mirror-1": {
					Host: "192.168.1.2",
					Port: 8081,
				},
			}
			current := map[string]target.Target{
				"mirror-1": {
					Host: "192.168.1.3",
					Port: 8081,
				},
			}
			return test{
				name: "Succeeded detecting the updated resource and failed to update the target connection",
				args: args{
					ctx:  context.Background(),
					prev: prev,
				},
				fields: fields{
					mirr: &MirrorMock{
						ConnectFunc: func(_ context.Context, _ ...*payload.Mirror_Target) error {
							return nil
						},
						DisconnectFunc: func(_ context.Context, _ ...*payload.Mirror_Target) error {
							return errors.New("err")
						},
					},
					ctrl: &k8smock.ControllerMock{
						GetManagerFunc: k8smock.NewDefaultManagerMock,
					},
				},
				beforeFunc: func(t *testing.T, d *discovery, _ args) {
					t.Helper()
					d.onReconcile(context.Background(), current)
				},
				want: want{
					want: current,
					err:  errors.New("err"),
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())

			d, err := NewDiscovery(
				WithDiscoveryController(test.fields.ctrl),
				WithDiscoveryMirror(test.fields.mirr),
			)
			if err != nil {
				t.Fatal(err)
			}

			if dis, ok := d.(*discovery); ok {
				if test.beforeFunc != nil {
					test.beforeFunc(tt, dis, test.args)
				}
				if test.afterFunc != nil {
					defer test.afterFunc(tt, test.args)
				}
				checkFunc := test.checkFunc
				if test.checkFunc == nil {
					checkFunc = defaultCheckFunc
				}

				got, err := dis.startSync(test.args.ctx, test.args.prev)
				if err := checkFunc(test.want, got, err); err != nil {
					tt.Errorf("error = %v", err)
				}
			}
		})
	}
}

func Test_discovery_syncWithAddr(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx      context.Context
		current  map[string]target.Target
		curAddrs map[string]string
	}
	type fields struct {
		ctrl k8s.Controller
		mirr Mirror
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
		beforeFunc func(*testing.T, *discovery, args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func(w want, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		return nil
	}
	tests := []test{
		func() test {
			current := map[string]target.Target{
				"mirror-1": {
					Host:  "192.168.1.2",
					Port:  8081,
					Phase: target.MirrorTargetPhaseConnected,
				},
			}
			curAddrs := map[string]string{
				"192.168.1.2:8081": "mirror-1",
			}

			return test{
				name: "Succeeded to change to disconnected phase when curAddrs are not connected",
				args: args{
					ctx:      context.Background(),
					current:  current,
					curAddrs: curAddrs,
				},
				fields: fields{
					mirr: &MirrorMock{
						IsConnectedFunc: func(_ context.Context, _ string) bool {
							return false
						},
						RangeMirrorAddrFunc: func(_ func(addr string, _ any) bool) {
							// There is no connection.
						},
					},
					ctrl: &k8smock.ControllerMock{
						GetManagerFunc: k8smock.NewDefaultManagerMock,
					},
				},
				beforeFunc: func(t *testing.T, d *discovery, _ args) {
					t.Helper()
					d.onReconcile(context.Background(), current)
				},
			}
		}(),
		func() test {
			current := map[string]target.Target{
				"mirror-1": {
					Host:  "192.168.1.2",
					Port:  8081,
					Phase: target.MirrorTargetPhaseDisconnected,
				},
			}
			curAddrs := map[string]string{
				"192.168.1.2:8081": "mirror-1",
			}

			return test{
				name: "Succeeded to change to connected phase when curAddrs are connected",
				args: args{
					ctx:      context.Background(),
					current:  current,
					curAddrs: curAddrs,
				},
				fields: fields{
					mirr: &MirrorMock{
						IsConnectedFunc: func(_ context.Context, _ string) bool {
							return true
						},
						RangeMirrorAddrFunc: func(f func(addr string, _ any) bool) {
							for addr := range curAddrs {
								f(addr, struct{}{})
							}
						},
					},
					ctrl: &k8smock.ControllerMock{
						GetManagerFunc: k8smock.NewDefaultManagerMock,
					},
				},
			}
		}(),
		func() test {
			current := map[string]target.Target{}
			curAddrs := map[string]string{}
			newAddrs := []string{
				"192.168.1.2:8081",
			}

			return test{
				name: "Succeeded to create new resource when there is a new connection",
				args: args{
					ctx:      context.Background(),
					current:  current,
					curAddrs: curAddrs,
				},
				fields: fields{
					mirr: &MirrorMock{
						IsConnectedFunc: func(_ context.Context, _ string) bool {
							return true
						},
						RangeMirrorAddrFunc: func(f func(addr string, _ any) bool) {
							for _, addr := range newAddrs {
								f(addr, struct{}{})
							}
						},
					},
					ctrl: &k8smock.ControllerMock{
						GetManagerFunc: k8smock.NewDefaultManagerMock,
					},
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())

			d, err := NewDiscovery(
				WithDiscoveryController(test.fields.ctrl),
				WithDiscoveryMirror(test.fields.mirr),
			)
			if err != nil {
				t.Fatal(err)
			}

			if dis, ok := d.(*discovery); ok {
				if test.beforeFunc != nil {
					test.beforeFunc(tt, dis, test.args)
				}
				if test.afterFunc != nil {
					defer test.afterFunc(tt, test.args)
				}
				checkFunc := test.checkFunc
				if test.checkFunc == nil {
					checkFunc = defaultCheckFunc
				}

				err := dis.syncWithAddr(test.args.ctx, test.args.current, test.args.curAddrs)
				if err := checkFunc(test.want, err); err != nil {
					tt.Errorf("error = %v", err)
				}
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
//
// func TestNewDiscovery(t *testing.T) {
// 	type args struct {
// 		opts []DiscoveryOption
// 	}
// 	type want struct {
// 		wantDsc Discovery
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, Discovery, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotDsc Discovery, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotDsc, w.wantDsc) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotDsc, w.wantDsc)
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
// 			gotDsc, err := NewDiscovery(test.args.opts...)
// 			if err := checkFunc(test.want, gotDsc, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_discovery_onReconcile(t *testing.T) {
// 	type args struct {
// 		in0  context.Context
// 		list map[string]target.Target
// 	}
// 	type fields struct {
// 		namespace       string
// 		labels          map[string]string
// 		colocation      string
// 		der             net.Dialer
// 		targetsByName   atomic.Pointer[map[string]target.Target]
// 		ctrl            k8s.Controller
// 		dur             time.Duration
// 		selfMirrAddrs   []string
// 		selfMirrAddrStr string
// 		mirr            Mirror
// 		eg              errgroup.Group
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
// 		           in0:nil,
// 		           list:nil,
// 		       },
// 		       fields: fields {
// 		           namespace:"",
// 		           labels:nil,
// 		           colocation:"",
// 		           der:nil,
// 		           targetsByName:nil,
// 		           ctrl:nil,
// 		           dur:nil,
// 		           selfMirrAddrs:nil,
// 		           selfMirrAddrStr:"",
// 		           mirr:nil,
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
// 		           in0:nil,
// 		           list:nil,
// 		           },
// 		           fields: fields {
// 		           namespace:"",
// 		           labels:nil,
// 		           colocation:"",
// 		           der:nil,
// 		           targetsByName:nil,
// 		           ctrl:nil,
// 		           dur:nil,
// 		           selfMirrAddrs:nil,
// 		           selfMirrAddrStr:"",
// 		           mirr:nil,
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
// 			d := &discovery{
// 				namespace:       test.fields.namespace,
// 				labels:          test.fields.labels,
// 				colocation:      test.fields.colocation,
// 				der:             test.fields.der,
// 				targetsByName:   test.fields.targetsByName,
// 				ctrl:            test.fields.ctrl,
// 				dur:             test.fields.dur,
// 				selfMirrAddrs:   test.fields.selfMirrAddrs,
// 				selfMirrAddrStr: test.fields.selfMirrAddrStr,
// 				mirr:            test.fields.mirr,
// 				eg:              test.fields.eg,
// 			}
//
// 			d.onReconcile(test.args.in0, test.args.list)
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_discovery_Start(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		namespace       string
// 		labels          map[string]string
// 		colocation      string
// 		der             net.Dialer
// 		targetsByName   atomic.Pointer[map[string]target.Target]
// 		ctrl            k8s.Controller
// 		dur             time.Duration
// 		selfMirrAddrs   []string
// 		selfMirrAddrStr string
// 		mirr            Mirror
// 		eg              errgroup.Group
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
// 		           namespace:"",
// 		           labels:nil,
// 		           colocation:"",
// 		           der:nil,
// 		           targetsByName:nil,
// 		           ctrl:nil,
// 		           dur:nil,
// 		           selfMirrAddrs:nil,
// 		           selfMirrAddrStr:"",
// 		           mirr:nil,
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
// 		           namespace:"",
// 		           labels:nil,
// 		           colocation:"",
// 		           der:nil,
// 		           targetsByName:nil,
// 		           ctrl:nil,
// 		           dur:nil,
// 		           selfMirrAddrs:nil,
// 		           selfMirrAddrStr:"",
// 		           mirr:nil,
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
// 			d := &discovery{
// 				namespace:       test.fields.namespace,
// 				labels:          test.fields.labels,
// 				colocation:      test.fields.colocation,
// 				der:             test.fields.der,
// 				targetsByName:   test.fields.targetsByName,
// 				ctrl:            test.fields.ctrl,
// 				dur:             test.fields.dur,
// 				selfMirrAddrs:   test.fields.selfMirrAddrs,
// 				selfMirrAddrStr: test.fields.selfMirrAddrStr,
// 				mirr:            test.fields.mirr,
// 				eg:              test.fields.eg,
// 			}
//
// 			got, err := d.Start(test.args.ctx)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_discovery_loadTargets(t *testing.T) {
// 	type fields struct {
// 		namespace       string
// 		labels          map[string]string
// 		colocation      string
// 		der             net.Dialer
// 		targetsByName   atomic.Pointer[map[string]target.Target]
// 		ctrl            k8s.Controller
// 		dur             time.Duration
// 		selfMirrAddrs   []string
// 		selfMirrAddrStr string
// 		mirr            Mirror
// 		eg              errgroup.Group
// 	}
// 	type want struct {
// 		want map[string]target.Target
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, map[string]target.Target) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got map[string]target.Target) error {
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
// 		           namespace:"",
// 		           labels:nil,
// 		           colocation:"",
// 		           der:nil,
// 		           targetsByName:nil,
// 		           ctrl:nil,
// 		           dur:nil,
// 		           selfMirrAddrs:nil,
// 		           selfMirrAddrStr:"",
// 		           mirr:nil,
// 		           eg:nil,
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
// 		           namespace:"",
// 		           labels:nil,
// 		           colocation:"",
// 		           der:nil,
// 		           targetsByName:nil,
// 		           ctrl:nil,
// 		           dur:nil,
// 		           selfMirrAddrs:nil,
// 		           selfMirrAddrStr:"",
// 		           mirr:nil,
// 		           eg:nil,
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
// 			d := &discovery{
// 				namespace:       test.fields.namespace,
// 				labels:          test.fields.labels,
// 				colocation:      test.fields.colocation,
// 				der:             test.fields.der,
// 				targetsByName:   test.fields.targetsByName,
// 				ctrl:            test.fields.ctrl,
// 				dur:             test.fields.dur,
// 				selfMirrAddrs:   test.fields.selfMirrAddrs,
// 				selfMirrAddrStr: test.fields.selfMirrAddrStr,
// 				mirr:            test.fields.mirr,
// 				eg:              test.fields.eg,
// 			}
//
// 			got := d.loadTargets()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_discovery_connectTarget(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		req map[string]*createdTarget
// 	}
// 	type fields struct {
// 		namespace       string
// 		labels          map[string]string
// 		colocation      string
// 		der             net.Dialer
// 		targetsByName   atomic.Pointer[map[string]target.Target]
// 		ctrl            k8s.Controller
// 		dur             time.Duration
// 		selfMirrAddrs   []string
// 		selfMirrAddrStr string
// 		mirr            Mirror
// 		eg              errgroup.Group
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
// 		           req:nil,
// 		       },
// 		       fields: fields {
// 		           namespace:"",
// 		           labels:nil,
// 		           colocation:"",
// 		           der:nil,
// 		           targetsByName:nil,
// 		           ctrl:nil,
// 		           dur:nil,
// 		           selfMirrAddrs:nil,
// 		           selfMirrAddrStr:"",
// 		           mirr:nil,
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
// 		           req:nil,
// 		           },
// 		           fields: fields {
// 		           namespace:"",
// 		           labels:nil,
// 		           colocation:"",
// 		           der:nil,
// 		           targetsByName:nil,
// 		           ctrl:nil,
// 		           dur:nil,
// 		           selfMirrAddrs:nil,
// 		           selfMirrAddrStr:"",
// 		           mirr:nil,
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
// 			d := &discovery{
// 				namespace:       test.fields.namespace,
// 				labels:          test.fields.labels,
// 				colocation:      test.fields.colocation,
// 				der:             test.fields.der,
// 				targetsByName:   test.fields.targetsByName,
// 				ctrl:            test.fields.ctrl,
// 				dur:             test.fields.dur,
// 				selfMirrAddrs:   test.fields.selfMirrAddrs,
// 				selfMirrAddrStr: test.fields.selfMirrAddrStr,
// 				mirr:            test.fields.mirr,
// 				eg:              test.fields.eg,
// 			}
//
// 			err := d.connectTarget(test.args.ctx, test.args.req)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_discovery_createMirrorTargetResource(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		name string
// 		host string
// 		port int
// 	}
// 	type fields struct {
// 		namespace       string
// 		labels          map[string]string
// 		colocation      string
// 		der             net.Dialer
// 		targetsByName   atomic.Pointer[map[string]target.Target]
// 		ctrl            k8s.Controller
// 		dur             time.Duration
// 		selfMirrAddrs   []string
// 		selfMirrAddrStr string
// 		mirr            Mirror
// 		eg              errgroup.Group
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
// 		           name:"",
// 		           host:"",
// 		           port:0,
// 		       },
// 		       fields: fields {
// 		           namespace:"",
// 		           labels:nil,
// 		           colocation:"",
// 		           der:nil,
// 		           targetsByName:nil,
// 		           ctrl:nil,
// 		           dur:nil,
// 		           selfMirrAddrs:nil,
// 		           selfMirrAddrStr:"",
// 		           mirr:nil,
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
// 		           name:"",
// 		           host:"",
// 		           port:0,
// 		           },
// 		           fields: fields {
// 		           namespace:"",
// 		           labels:nil,
// 		           colocation:"",
// 		           der:nil,
// 		           targetsByName:nil,
// 		           ctrl:nil,
// 		           dur:nil,
// 		           selfMirrAddrs:nil,
// 		           selfMirrAddrStr:"",
// 		           mirr:nil,
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
// 			d := &discovery{
// 				namespace:       test.fields.namespace,
// 				labels:          test.fields.labels,
// 				colocation:      test.fields.colocation,
// 				der:             test.fields.der,
// 				targetsByName:   test.fields.targetsByName,
// 				ctrl:            test.fields.ctrl,
// 				dur:             test.fields.dur,
// 				selfMirrAddrs:   test.fields.selfMirrAddrs,
// 				selfMirrAddrStr: test.fields.selfMirrAddrStr,
// 				mirr:            test.fields.mirr,
// 				eg:              test.fields.eg,
// 			}
//
// 			err := d.createMirrorTargetResource(test.args.ctx, test.args.name, test.args.host, test.args.port)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_discovery_disconnectTarget(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		req map[string]*deletedTarget
// 	}
// 	type fields struct {
// 		namespace       string
// 		labels          map[string]string
// 		colocation      string
// 		der             net.Dialer
// 		targetsByName   atomic.Pointer[map[string]target.Target]
// 		ctrl            k8s.Controller
// 		dur             time.Duration
// 		selfMirrAddrs   []string
// 		selfMirrAddrStr string
// 		mirr            Mirror
// 		eg              errgroup.Group
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
// 		           req:nil,
// 		       },
// 		       fields: fields {
// 		           namespace:"",
// 		           labels:nil,
// 		           colocation:"",
// 		           der:nil,
// 		           targetsByName:nil,
// 		           ctrl:nil,
// 		           dur:nil,
// 		           selfMirrAddrs:nil,
// 		           selfMirrAddrStr:"",
// 		           mirr:nil,
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
// 		           req:nil,
// 		           },
// 		           fields: fields {
// 		           namespace:"",
// 		           labels:nil,
// 		           colocation:"",
// 		           der:nil,
// 		           targetsByName:nil,
// 		           ctrl:nil,
// 		           dur:nil,
// 		           selfMirrAddrs:nil,
// 		           selfMirrAddrStr:"",
// 		           mirr:nil,
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
// 			d := &discovery{
// 				namespace:       test.fields.namespace,
// 				labels:          test.fields.labels,
// 				colocation:      test.fields.colocation,
// 				der:             test.fields.der,
// 				targetsByName:   test.fields.targetsByName,
// 				ctrl:            test.fields.ctrl,
// 				dur:             test.fields.dur,
// 				selfMirrAddrs:   test.fields.selfMirrAddrs,
// 				selfMirrAddrStr: test.fields.selfMirrAddrStr,
// 				mirr:            test.fields.mirr,
// 				eg:              test.fields.eg,
// 			}
//
// 			err := d.disconnectTarget(test.args.ctx, test.args.req)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_discovery_updateMirrorTargetPhase(t *testing.T) {
// 	type args struct {
// 		ctx   context.Context
// 		name  string
// 		phase target.MirrorTargetPhase
// 	}
// 	type fields struct {
// 		namespace       string
// 		labels          map[string]string
// 		colocation      string
// 		der             net.Dialer
// 		targetsByName   atomic.Pointer[map[string]target.Target]
// 		ctrl            k8s.Controller
// 		dur             time.Duration
// 		selfMirrAddrs   []string
// 		selfMirrAddrStr string
// 		mirr            Mirror
// 		eg              errgroup.Group
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
// 		           name:"",
// 		           phase:nil,
// 		       },
// 		       fields: fields {
// 		           namespace:"",
// 		           labels:nil,
// 		           colocation:"",
// 		           der:nil,
// 		           targetsByName:nil,
// 		           ctrl:nil,
// 		           dur:nil,
// 		           selfMirrAddrs:nil,
// 		           selfMirrAddrStr:"",
// 		           mirr:nil,
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
// 		           name:"",
// 		           phase:nil,
// 		           },
// 		           fields: fields {
// 		           namespace:"",
// 		           labels:nil,
// 		           colocation:"",
// 		           der:nil,
// 		           targetsByName:nil,
// 		           ctrl:nil,
// 		           dur:nil,
// 		           selfMirrAddrs:nil,
// 		           selfMirrAddrStr:"",
// 		           mirr:nil,
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
// 			d := &discovery{
// 				namespace:       test.fields.namespace,
// 				labels:          test.fields.labels,
// 				colocation:      test.fields.colocation,
// 				der:             test.fields.der,
// 				targetsByName:   test.fields.targetsByName,
// 				ctrl:            test.fields.ctrl,
// 				dur:             test.fields.dur,
// 				selfMirrAddrs:   test.fields.selfMirrAddrs,
// 				selfMirrAddrStr: test.fields.selfMirrAddrStr,
// 				mirr:            test.fields.mirr,
// 				eg:              test.fields.eg,
// 			}
//
// 			err := d.updateMirrorTargetPhase(test.args.ctx, test.args.name, test.args.phase)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_discovery_updateTarget(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		req map[string]*updatedTarget
// 	}
// 	type fields struct {
// 		namespace       string
// 		labels          map[string]string
// 		colocation      string
// 		der             net.Dialer
// 		targetsByName   atomic.Pointer[map[string]target.Target]
// 		ctrl            k8s.Controller
// 		dur             time.Duration
// 		selfMirrAddrs   []string
// 		selfMirrAddrStr string
// 		mirr            Mirror
// 		eg              errgroup.Group
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
// 		           req:nil,
// 		       },
// 		       fields: fields {
// 		           namespace:"",
// 		           labels:nil,
// 		           colocation:"",
// 		           der:nil,
// 		           targetsByName:nil,
// 		           ctrl:nil,
// 		           dur:nil,
// 		           selfMirrAddrs:nil,
// 		           selfMirrAddrStr:"",
// 		           mirr:nil,
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
// 		           req:nil,
// 		           },
// 		           fields: fields {
// 		           namespace:"",
// 		           labels:nil,
// 		           colocation:"",
// 		           der:nil,
// 		           targetsByName:nil,
// 		           ctrl:nil,
// 		           dur:nil,
// 		           selfMirrAddrs:nil,
// 		           selfMirrAddrStr:"",
// 		           mirr:nil,
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
// 			d := &discovery{
// 				namespace:       test.fields.namespace,
// 				labels:          test.fields.labels,
// 				colocation:      test.fields.colocation,
// 				der:             test.fields.der,
// 				targetsByName:   test.fields.targetsByName,
// 				ctrl:            test.fields.ctrl,
// 				dur:             test.fields.dur,
// 				selfMirrAddrs:   test.fields.selfMirrAddrs,
// 				selfMirrAddrStr: test.fields.selfMirrAddrStr,
// 				mirr:            test.fields.mirr,
// 				eg:              test.fields.eg,
// 			}
//
// 			err := d.updateTarget(test.args.ctx, test.args.req)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_isConnectedPhase(t *testing.T) {
// 	type args struct {
// 		phase target.MirrorTargetPhase
// 	}
// 	type want struct {
// 		want bool
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, bool) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got bool) error {
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
// 		           phase:nil,
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
// 		           phase:nil,
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
// 			got := isConnectedPhase(test.args.phase)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_isDisconnectedPhase(t *testing.T) {
// 	type args struct {
// 		phase target.MirrorTargetPhase
// 	}
// 	type want struct {
// 		want bool
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, bool) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got bool) error {
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
// 		           phase:nil,
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
// 		           phase:nil,
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
// 			got := isDisconnectedPhase(test.args.phase)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
