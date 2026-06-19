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
package service

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/k8s"
	"github.com/vdaas/vald/internal/k8s/client"
	mock "github.com/vdaas/vald/internal/test/mock/k8s"
	"github.com/vdaas/vald/internal/test/testify"
)

func Test_getNewBaseName(t *testing.T) {
	type args struct {
		old string
	}
	type want struct {
		want string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "fist rotation just returns the name",
			args: args{
				old: "vald-agent-ngt-readreplica-0",
			},
			want: want{
				want: "vald-agent-ngt-readreplica-0-",
			},
		},
		{
			name: "successfully remove timestamp",
			args: args{
				old: "vald-agent-ngt-readreplica-0-20220101",
			},
			want: want{
				want: "vald-agent-ngt-readreplica-0-",
			},
		},
		{
			name: "successfully remove timestamp when the name has no dashes",
			args: args{
				old: "vald-1-20220101",
			},
			want: want{
				want: "vald-1-",
			},
		},
		{
			name: "no replica id basename returns empty string",
			args: args{
				old: "vald",
			},
			want: want{
				want: "",
			},
		},
	}
	for _, test := range tests {
		tt := test
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := getNewBaseName(tt.args.old); got != tt.want.want {
				t.Errorf("getNewBaseName() = %v, want %v", got, tt.want.want)
			}
		})
	}
}

func Test_parseReplicaID(t *testing.T) {
	labelKey := "foo"
	type args struct {
		c         client.Client
		replicaID string
	}
	type want struct {
		err error
		ids []string
	}
	type test struct {
		args args
		name string
		want want
	}
	tests := []test{
		{
			name: "single replicaID",
			args: args{
				replicaID: "0",
				c:         nil,
			},
			want: want{
				ids: []string{"0"},
				err: nil,
			},
		},
		{
			name: "multiple replicaIDs",
			args: args{
				replicaID: "0,1",
				c:         nil,
			},
			want: want{
				ids: []string{"0", "1"},
				err: nil,
			},
		},
		{
			name: "returns error when replicaID is empty",
			args: args{
				replicaID: "",
				c:         nil,
			},
			want: want{
				ids: nil,
				err: errors.ErrReadReplicaIDEmpty,
			},
		},
		func() test {
			wantID1 := "bar"
			wantID2 := "baz"
			mock := &mock.ValdK8sClientMock{}

			mock.On("LabelSelector", testify.Anything, testify.Anything, testify.Anything).Return(client.NewSelector(), nil)
			mock.On("List", testify.Anything, testify.Anything, testify.Anything).Run(func(args testify.Arguments) {
				if depList, ok := args.Get(1).(*k8s.DeploymentList); ok {
					depList.Items = []k8s.Deployment{
						{
							ObjectMeta: k8s.ObjectMeta{
								Labels: map[string]string{
									labelKey: wantID1,
								},
							},
						},
						{
							ObjectMeta: k8s.ObjectMeta{
								Labels: map[string]string{
									labelKey: wantID2,
								},
							},
						},
					}
				}
			}).Return(nil)
			return test{
				name: "returns all ids when rotate-all option is set",
				args: args{
					replicaID: rotateAllID,
					c:         mock,
				},
				want: want{
					ids: []string{wantID1, wantID2},
					err: nil,
				},
			}
		}(),
	}
	for _, test := range tests {
		tt := test
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := &rotator{
				readReplicaLabelKey: labelKey,
			}
			ids, err := r.parseReplicaID(tt.args.replicaID, tt.args.c)
			require.Equal(t, tt.want.ids, ids)
			require.Equal(t, tt.want.err, err)
		})
	}
}

// NOT IMPLEMENTED BELOW
//
// func TestNew(t *testing.T) {
// 	type args struct {
// 		replicaID string
// 		opts      []Option
// 	}
// 	type want struct {
// 		want Rotator
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, Rotator, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got Rotator, err error) error {
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
// 		           replicaID:"",
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
// 		           replicaID:"",
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
// 			got, err := New(test.args.replicaID, test.args.opts...)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_rotator_Start(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		namespace           string
// 		volumeName          string
// 		readReplicaLabelKey string
// 		subProcesses        []subProcess
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
// 		           namespace:"",
// 		           volumeName:"",
// 		           readReplicaLabelKey:"",
// 		           subProcesses:nil,
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
// 		           volumeName:"",
// 		           readReplicaLabelKey:"",
// 		           subProcesses:nil,
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
// 			r := &rotator{
// 				namespace:           test.fields.namespace,
// 				volumeName:          test.fields.volumeName,
// 				readReplicaLabelKey: test.fields.readReplicaLabelKey,
// 				subProcesses:        test.fields.subProcesses,
// 			}
//
// 			err := r.Start(test.args.ctx)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_rotator_newSubprocess(t *testing.T) {
// 	type args struct {
// 		c         client.Client
// 		replicaID string
// 	}
// 	type fields struct {
// 		namespace           string
// 		volumeName          string
// 		readReplicaLabelKey string
// 		subProcesses        []subProcess
// 	}
// 	type want struct {
// 		want subProcess
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, subProcess, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got subProcess, err error) error {
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
// 		           c:nil,
// 		           replicaID:"",
// 		       },
// 		       fields: fields {
// 		           namespace:"",
// 		           volumeName:"",
// 		           readReplicaLabelKey:"",
// 		           subProcesses:nil,
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
// 		           replicaID:"",
// 		           },
// 		           fields: fields {
// 		           namespace:"",
// 		           volumeName:"",
// 		           readReplicaLabelKey:"",
// 		           subProcesses:nil,
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
// 			r := &rotator{
// 				namespace:           test.fields.namespace,
// 				volumeName:          test.fields.volumeName,
// 				readReplicaLabelKey: test.fields.readReplicaLabelKey,
// 				subProcesses:        test.fields.subProcesses,
// 			}
//
// 			got, err := r.newSubprocess(test.args.c, test.args.replicaID)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_subProcess_rotate(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		listOpts   k8s.ListOptions
// 		client     client.Client
// 		volumeName string
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
// 		           listOpts:nil,
// 		           client:nil,
// 		           volumeName:"",
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
// 		           listOpts:nil,
// 		           client:nil,
// 		           volumeName:"",
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
// 			s := &subProcess{
// 				listOpts:   test.fields.listOpts,
// 				client:     test.fields.client,
// 				volumeName: test.fields.volumeName,
// 			}
//
// 			err := s.rotate(test.args.ctx)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_subProcess_createSnapshot(t *testing.T) {
// 	type args struct {
// 		ctx        context.Context
// 		deployment *k8s.Deployment
// 	}
// 	type fields struct {
// 		listOpts   k8s.ListOptions
// 		client     client.Client
// 		volumeName string
// 	}
// 	type want struct {
// 		wantNewSnap *k8s.VolumeSnapshot
// 		wantOldSnap *k8s.VolumeSnapshot
// 		err         error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *k8s.VolumeSnapshot, *k8s.VolumeSnapshot, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotNewSnap *k8s.VolumeSnapshot, gotOldSnap *k8s.VolumeSnapshot, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotNewSnap, w.wantNewSnap) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotNewSnap, w.wantNewSnap)
// 		}
// 		if !reflect.DeepEqual(gotOldSnap, w.wantOldSnap) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotOldSnap, w.wantOldSnap)
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
// 		           deployment:nil,
// 		       },
// 		       fields: fields {
// 		           listOpts:nil,
// 		           client:nil,
// 		           volumeName:"",
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
// 		           deployment:nil,
// 		           },
// 		           fields: fields {
// 		           listOpts:nil,
// 		           client:nil,
// 		           volumeName:"",
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
// 			s := &subProcess{
// 				listOpts:   test.fields.listOpts,
// 				client:     test.fields.client,
// 				volumeName: test.fields.volumeName,
// 			}
//
// 			gotNewSnap, gotOldSnap, err := s.createSnapshot(test.args.ctx, test.args.deployment)
// 			if err := checkFunc(test.want, gotNewSnap, gotOldSnap, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_subProcess_createPVC(t *testing.T) {
// 	type args struct {
// 		ctx         context.Context
// 		newSnapShot string
// 		deployment  *k8s.Deployment
// 	}
// 	type fields struct {
// 		listOpts   k8s.ListOptions
// 		client     client.Client
// 		volumeName string
// 	}
// 	type want struct {
// 		wantNewPvc *k8s.PersistentVolumeClaim
// 		wantOldPvc *k8s.PersistentVolumeClaim
// 		err        error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *k8s.PersistentVolumeClaim, *k8s.PersistentVolumeClaim, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotNewPvc *k8s.PersistentVolumeClaim, gotOldPvc *k8s.PersistentVolumeClaim, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotNewPvc, w.wantNewPvc) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotNewPvc, w.wantNewPvc)
// 		}
// 		if !reflect.DeepEqual(gotOldPvc, w.wantOldPvc) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotOldPvc, w.wantOldPvc)
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
// 		           newSnapShot:"",
// 		           deployment:nil,
// 		       },
// 		       fields: fields {
// 		           listOpts:nil,
// 		           client:nil,
// 		           volumeName:"",
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
// 		           newSnapShot:"",
// 		           deployment:nil,
// 		           },
// 		           fields: fields {
// 		           listOpts:nil,
// 		           client:nil,
// 		           volumeName:"",
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
// 			s := &subProcess{
// 				listOpts:   test.fields.listOpts,
// 				client:     test.fields.client,
// 				volumeName: test.fields.volumeName,
// 			}
//
// 			gotNewPvc, gotOldPvc, err := s.createPVC(test.args.ctx, test.args.newSnapShot, test.args.deployment)
// 			if err := checkFunc(test.want, gotNewPvc, gotOldPvc, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_subProcess_getDeployment(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		listOpts   k8s.ListOptions
// 		client     client.Client
// 		volumeName string
// 	}
// 	type want struct {
// 		want *k8s.Deployment
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *k8s.Deployment, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got *k8s.Deployment, err error) error {
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
// 		           listOpts:nil,
// 		           client:nil,
// 		           volumeName:"",
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
// 		           listOpts:nil,
// 		           client:nil,
// 		           volumeName:"",
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
// 			s := &subProcess{
// 				listOpts:   test.fields.listOpts,
// 				client:     test.fields.client,
// 				volumeName: test.fields.volumeName,
// 			}
//
// 			got, err := s.getDeployment(test.args.ctx)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_subProcess_updateDeployment(t *testing.T) {
// 	type args struct {
// 		ctx          context.Context
// 		newPVC       string
// 		deployment   *k8s.Deployment
// 		snapshotTime time.Time
// 	}
// 	type fields struct {
// 		listOpts   k8s.ListOptions
// 		client     client.Client
// 		volumeName string
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
// 		           newPVC:"",
// 		           deployment:nil,
// 		           snapshotTime:time.Time{},
// 		       },
// 		       fields: fields {
// 		           listOpts:nil,
// 		           client:nil,
// 		           volumeName:"",
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
// 		           newPVC:"",
// 		           deployment:nil,
// 		           snapshotTime:time.Time{},
// 		           },
// 		           fields: fields {
// 		           listOpts:nil,
// 		           client:nil,
// 		           volumeName:"",
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
// 			s := &subProcess{
// 				listOpts:   test.fields.listOpts,
// 				client:     test.fields.client,
// 				volumeName: test.fields.volumeName,
// 			}
//
// 			err := s.updateDeployment(test.args.ctx, test.args.newPVC, test.args.deployment, test.args.snapshotTime)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_subProcess_deleteSnapshot(t *testing.T) {
// 	type args struct {
// 		ctx      context.Context
// 		snapshot *k8s.VolumeSnapshot
// 	}
// 	type fields struct {
// 		listOpts   k8s.ListOptions
// 		client     client.Client
// 		volumeName string
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
// 		           snapshot:nil,
// 		       },
// 		       fields: fields {
// 		           listOpts:nil,
// 		           client:nil,
// 		           volumeName:"",
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
// 		           snapshot:nil,
// 		           },
// 		           fields: fields {
// 		           listOpts:nil,
// 		           client:nil,
// 		           volumeName:"",
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
// 			s := &subProcess{
// 				listOpts:   test.fields.listOpts,
// 				client:     test.fields.client,
// 				volumeName: test.fields.volumeName,
// 			}
//
// 			err := s.deleteSnapshot(test.args.ctx, test.args.snapshot)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_subProcess_deletePVC(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		pvc *k8s.PersistentVolumeClaim
// 	}
// 	type fields struct {
// 		listOpts   k8s.ListOptions
// 		client     client.Client
// 		volumeName string
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
// 		           pvc:nil,
// 		       },
// 		       fields: fields {
// 		           listOpts:nil,
// 		           client:nil,
// 		           volumeName:"",
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
// 		           pvc:nil,
// 		           },
// 		           fields: fields {
// 		           listOpts:nil,
// 		           client:nil,
// 		           volumeName:"",
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
// 			s := &subProcess{
// 				listOpts:   test.fields.listOpts,
// 				client:     test.fields.client,
// 				volumeName: test.fields.volumeName,
// 			}
//
// 			err := s.deletePVC(test.args.ctx, test.args.pvc)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_rotator_parseReplicaID(t *testing.T) {
// 	type args struct {
// 		replicaID string
// 		c         client.Client
// 	}
// 	type fields struct {
// 		namespace           string
// 		volumeName          string
// 		readReplicaLabelKey string
// 		subProcesses        []subProcess
// 	}
// 	type want struct {
// 		want []string
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, []string, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got []string, err error) error {
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
// 		           replicaID:"",
// 		           c:nil,
// 		       },
// 		       fields: fields {
// 		           namespace:"",
// 		           volumeName:"",
// 		           readReplicaLabelKey:"",
// 		           subProcesses:nil,
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
// 		           replicaID:"",
// 		           c:nil,
// 		           },
// 		           fields: fields {
// 		           namespace:"",
// 		           volumeName:"",
// 		           readReplicaLabelKey:"",
// 		           subProcesses:nil,
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
// 			r := &rotator{
// 				namespace:           test.fields.namespace,
// 				volumeName:          test.fields.volumeName,
// 				readReplicaLabelKey: test.fields.readReplicaLabelKey,
// 				subProcesses:        test.fields.subProcesses,
// 			}
//
// 			got, err := r.parseReplicaID(test.args.replicaID, test.args.c)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
