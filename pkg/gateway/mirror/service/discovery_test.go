// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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

func Test_discovery_Start(t *testing.T) {
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
						GetManagerFunc: func() k8s.Manager {
							return k8smock.NewDefaultManagerMock()
						},
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
						GetManagerFunc: func() k8s.Manager {
							return k8smock.NewDefaultManagerMock()
						},
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
						GetManagerFunc: func() k8s.Manager {
							return k8smock.NewDefaultManagerMock()
						},
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
						GetManagerFunc: func() k8s.Manager {
							return k8smock.NewDefaultManagerMock()
						},
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
						GetManagerFunc: func() k8s.Manager {
							return k8smock.NewDefaultManagerMock()
						},
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
						GetManagerFunc: func() k8s.Manager {
							return k8smock.NewDefaultManagerMock()
						},
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
//
// 		})
// 	}
// }
//
