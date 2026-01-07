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

package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/vdaas/vald/internal/k8s"
	"github.com/vdaas/vald/internal/k8s/client"
	"github.com/vdaas/vald/internal/k8s/vald"
	mock "github.com/vdaas/vald/internal/test/mock/k8s"
	"github.com/vdaas/vald/internal/test/testify"
)

func Test_operator_podOnReconcile(t *testing.T) {
	t.Parallel()

	type want struct {
		res          k8s.Result
		createCalled bool
		err          error
	}
	type test struct {
		name                   string
		agentPod               *k8s.Pod
		readReplicaEnabled     bool
		readReplicaDeployment  *k8s.Deployment
		runningJobs            []k8s.Job
		rotationJobConcurrency uint
		want                   want
	}

	tests := []test{
		{
			name:               "returns client.Result{} when read replica is not enabled",
			readReplicaEnabled: false,
			want: want{
				res:          k8s.Result{},
				createCalled: false,
				err:          nil,
			},
		},
		{
			name:               "returns client.Result{} when pod is not a statefulset",
			readReplicaEnabled: true,
			agentPod:           &k8s.Pod{},
			want: want{
				res:          k8s.Result{},
				createCalled: false,
				err:          nil,
			},
		},
		func() test {
			saveTime := time.Now()
			rotateTime := saveTime.Add(1 * time.Second)
			return test{
				name:               "returns requeue: false when last snapshot time is after the last save time",
				readReplicaEnabled: true,
				agentPod: &k8s.Pod{
					ObjectMeta: k8s.ObjectMeta{
						Labels: map[string]string{
							k8s.PodIndexLabel: "0",
						},
						Annotations: map[string]string{
							vald.LastTimeSaveIndexTimestampAnnotationsKey: saveTime.Format(vald.TimeFormat),
						},
					},
				},
				readReplicaDeployment: &k8s.Deployment{
					ObjectMeta: k8s.ObjectMeta{
						Name: "deploymentName",
						Annotations: map[string]string{
							vald.LastTimeSnapshotTimestampAnnotationsKey: rotateTime.Format(vald.TimeFormat),
						},
					},
				},
				want: want{
					res: k8s.Result{
						Requeue: false,
					},
					createCalled: false,
					err:          nil,
				},
			}
		}(),
		func() test {
			saveTime := time.Now()
			rotateTime := saveTime.Add(-1 * time.Second)
			return test{
				name:               "returns requeue: false and calls client.Create once when last snapshot time is before the last save time",
				readReplicaEnabled: true,
				agentPod: &k8s.Pod{
					ObjectMeta: k8s.ObjectMeta{
						Labels: map[string]string{
							k8s.PodIndexLabel: "0",
						},
						Annotations: map[string]string{
							vald.LastTimeSaveIndexTimestampAnnotationsKey: saveTime.Format(vald.TimeFormat),
						},
					},
				},
				readReplicaDeployment: &k8s.Deployment{
					ObjectMeta: k8s.ObjectMeta{
						Name: "deploymentName",
						Annotations: map[string]string{
							vald.LastTimeSnapshotTimestampAnnotationsKey: rotateTime.Format(vald.TimeFormat),
						},
					},
				},
				want: want{
					res: k8s.Result{
						Requeue: false,
					},
					createCalled: true,
					err:          nil,
				},
			}
		}(),
		func() test {
			saveTime := time.Now()
			rotateTime := saveTime.Add(-1 * time.Second)
			return test{
				name:               "returns requeue: true when there is already one running job when rotation job concurrency is 1",
				readReplicaEnabled: true,
				agentPod: &k8s.Pod{
					ObjectMeta: k8s.ObjectMeta{
						Labels: map[string]string{
							k8s.PodIndexLabel: "0",
						},
						Annotations: map[string]string{
							vald.LastTimeSaveIndexTimestampAnnotationsKey: saveTime.Format(vald.TimeFormat),
						},
					},
				},
				readReplicaDeployment: &k8s.Deployment{
					ObjectMeta: k8s.ObjectMeta{
						Name: "deploymentName",
						Annotations: map[string]string{
							vald.LastTimeSnapshotTimestampAnnotationsKey: rotateTime.Format(vald.TimeFormat),
						},
					},
				},
				runningJobs: []k8s.Job{
					{
						ObjectMeta: k8s.ObjectMeta{
							Name: "already running job1",
						},
						Status: k8s.JobStatus{
							Active: 1,
						},
					},
				},
				rotationJobConcurrency: 1,
				want: want{
					res: k8s.Result{
						Requeue: true,
					},
					createCalled: false,
					err:          nil,
				},
			}
		}(),
		func() test {
			saveTime := time.Now()
			rotateTime := saveTime.Add(-1 * time.Second)
			return test{
				name:               "returns requeue: false and create job when there is one running job when rotation job concurrency is 2",
				readReplicaEnabled: true,
				agentPod: &k8s.Pod{
					ObjectMeta: k8s.ObjectMeta{
						Labels: map[string]string{
							k8s.PodIndexLabel: "0",
						},
						Annotations: map[string]string{
							vald.LastTimeSaveIndexTimestampAnnotationsKey: saveTime.Format(vald.TimeFormat),
						},
					},
				},
				readReplicaDeployment: &k8s.Deployment{
					ObjectMeta: k8s.ObjectMeta{
						Name: "deploymentName",
						Annotations: map[string]string{
							vald.LastTimeSnapshotTimestampAnnotationsKey: rotateTime.Format(vald.TimeFormat),
						},
					},
				},
				runningJobs: []k8s.Job{
					{
						ObjectMeta: k8s.ObjectMeta{
							Name: "already running job1",
						},
						Status: k8s.JobStatus{
							Active: 1,
						},
					},
				},
				rotationJobConcurrency: 2,
				want: want{
					res: k8s.Result{
						Requeue: false,
					},
					createCalled: true,
					err:          nil,
				},
			}
		}(),
		func() test {
			saveTime := time.Now()
			rotateTime := saveTime.Add(-1 * time.Second)
			return test{
				name:               "returns requeue: true when there are two running jobs when rotation job concurrency is 2",
				readReplicaEnabled: true,
				agentPod: &k8s.Pod{
					ObjectMeta: k8s.ObjectMeta{
						Labels: map[string]string{
							k8s.PodIndexLabel: "0",
						},
						Annotations: map[string]string{
							vald.LastTimeSaveIndexTimestampAnnotationsKey: saveTime.Format(vald.TimeFormat),
						},
					},
				},
				readReplicaDeployment: &k8s.Deployment{
					ObjectMeta: k8s.ObjectMeta{
						Name: "deploymentName",
						Annotations: map[string]string{
							vald.LastTimeSnapshotTimestampAnnotationsKey: rotateTime.Format(vald.TimeFormat),
						},
					},
				},
				runningJobs: []k8s.Job{
					{
						ObjectMeta: k8s.ObjectMeta{
							Name: "already running job1",
						},
						Status: k8s.JobStatus{
							Active: 1,
						},
					},
					{
						ObjectMeta: k8s.ObjectMeta{
							Name: "already running job2",
						},
						Status: k8s.JobStatus{
							Active: 1,
						},
					},
				},
				rotationJobConcurrency: 2,
				want: want{
					res: k8s.Result{
						Requeue: true,
					},
					createCalled: false,
					err:          nil,
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()

			mock := &mock.ValdK8sClientMock{}
			mock.On("LabelSelector", testify.Anything, testify.Anything, testify.Anything).Return(client.NewSelector(), nil).Maybe()
			mock.On("List", testify.Anything, testify.AnythingOfType("*v1.DeploymentList"), testify.Anything).Run(func(args testify.Arguments) {
				arg, ok := args.Get(1).(*k8s.DeploymentList)
				require.True(t, ok)

				arg.Items = []k8s.Deployment{*test.readReplicaDeployment}
			}).Return(nil).Maybe()

			mock.On("List", testify.Anything, testify.AnythingOfType("*v1.JobList"), testify.Anything).Run(func(args testify.Arguments) {
				arg, ok := args.Get(1).(*k8s.JobList)
				require.True(t, ok)

				arg.Items = test.runningJobs
			}).Return(nil).Maybe()

			// testify/mock does not accept to set Times(0) so you cannot do things like .Return(nil).Once(calledTimes)
			// ref: https://github.com/stretchr/testify/issues/566
			if test.want.createCalled {
				mock.On("Create", testify.Anything, testify.Anything, testify.Anything).Return(nil).Once()
			}
			defer mock.AssertExpectations(tt)

			concurrency := uint(1)
			if test.rotationJobConcurrency != 0 {
				concurrency = test.rotationJobConcurrency
			}
			op := operator{
				client:                 mock,
				readReplicaEnabled:     test.readReplicaEnabled,
				rotationJobConcurrency: concurrency,
			}

			op.rotatorJob = &k8s.Job{
				ObjectMeta: k8s.ObjectMeta{
					Name: "foo job",
				},
			}

			res, err := op.podOnReconcile(context.Background(), test.agentPod)
			require.Equal(t, test.want.err, err)
			require.Equal(t, test.want.res, res)
		})
	}
}

// NOT IMPLEMENTED BELOW
//
// func TestNew(t *testing.T) {
// 	type args struct {
// 		namespace              string
// 		agentName              string
// 		rotatorName            string
// 		targetReadReplicaIDKey string
// 		rotatorJob             *k8s.Job
// 		opts                   []Option
// 	}
// 	type want struct {
// 		wantO Operator
// 		err   error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, Operator, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotO Operator, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotO, w.wantO) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotO, w.wantO)
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
// 		           agentName:"",
// 		           rotatorName:"",
// 		           targetReadReplicaIDKey:"",
// 		           rotatorJob:nil,
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
// 		           namespace:"",
// 		           agentName:"",
// 		           rotatorName:"",
// 		           targetReadReplicaIDKey:"",
// 		           rotatorJob:nil,
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
// 			gotO, err := New(test.args.namespace, test.args.agentName, test.args.rotatorName, test.args.targetReadReplicaIDKey, test.args.rotatorJob, test.args.opts...)
// 			if err := checkFunc(test.want, gotO, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_operator_Start(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		ctrl                              k8s.Controller
// 		eg                                errgroup.Group
// 		namespace                         string
// 		client                            client.Client
// 		rotatorName                       string
// 		targetReadReplicaIDAnnotationsKey string
// 		readReplicaEnabled                bool
// 		readReplicaLabelKey               string
// 		rotationJobConcurrency            uint
// 		rotatorJob                        *k8s.Job
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
// 		           ctrl:nil,
// 		           eg:nil,
// 		           namespace:"",
// 		           client:nil,
// 		           rotatorName:"",
// 		           targetReadReplicaIDAnnotationsKey:"",
// 		           readReplicaEnabled:false,
// 		           readReplicaLabelKey:"",
// 		           rotationJobConcurrency:0,
// 		           rotatorJob:nil,
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
// 		           ctrl:nil,
// 		           eg:nil,
// 		           namespace:"",
// 		           client:nil,
// 		           rotatorName:"",
// 		           targetReadReplicaIDAnnotationsKey:"",
// 		           readReplicaEnabled:false,
// 		           readReplicaLabelKey:"",
// 		           rotationJobConcurrency:0,
// 		           rotatorJob:nil,
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
// 			o := &operator{
// 				ctrl:                              test.fields.ctrl,
// 				eg:                                test.fields.eg,
// 				namespace:                         test.fields.namespace,
// 				client:                            test.fields.client,
// 				rotatorName:                       test.fields.rotatorName,
// 				targetReadReplicaIDAnnotationsKey: test.fields.targetReadReplicaIDAnnotationsKey,
// 				readReplicaEnabled:                test.fields.readReplicaEnabled,
// 				readReplicaLabelKey:               test.fields.readReplicaLabelKey,
// 				rotationJobConcurrency:            test.fields.rotationJobConcurrency,
// 				rotatorJob:                        test.fields.rotatorJob,
// 			}
//
// 			got, err := o.Start(test.args.ctx)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_operator_reconcileRotatorJob(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		pod *k8s.Pod
// 	}
// 	type fields struct {
// 		ctrl                              k8s.Controller
// 		eg                                errgroup.Group
// 		namespace                         string
// 		client                            client.Client
// 		rotatorName                       string
// 		targetReadReplicaIDAnnotationsKey string
// 		readReplicaEnabled                bool
// 		readReplicaLabelKey               string
// 		rotationJobConcurrency            uint
// 		rotatorJob                        *k8s.Job
// 	}
// 	type want struct {
// 		wantRequeue bool
// 		err         error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, bool, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotRequeue bool, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotRequeue, w.wantRequeue) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRequeue, w.wantRequeue)
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
// 		           pod:nil,
// 		       },
// 		       fields: fields {
// 		           ctrl:nil,
// 		           eg:nil,
// 		           namespace:"",
// 		           client:nil,
// 		           rotatorName:"",
// 		           targetReadReplicaIDAnnotationsKey:"",
// 		           readReplicaEnabled:false,
// 		           readReplicaLabelKey:"",
// 		           rotationJobConcurrency:0,
// 		           rotatorJob:nil,
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
// 		           pod:nil,
// 		           },
// 		           fields: fields {
// 		           ctrl:nil,
// 		           eg:nil,
// 		           namespace:"",
// 		           client:nil,
// 		           rotatorName:"",
// 		           targetReadReplicaIDAnnotationsKey:"",
// 		           readReplicaEnabled:false,
// 		           readReplicaLabelKey:"",
// 		           rotationJobConcurrency:0,
// 		           rotatorJob:nil,
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
// 			o := &operator{
// 				ctrl:                              test.fields.ctrl,
// 				eg:                                test.fields.eg,
// 				namespace:                         test.fields.namespace,
// 				client:                            test.fields.client,
// 				rotatorName:                       test.fields.rotatorName,
// 				targetReadReplicaIDAnnotationsKey: test.fields.targetReadReplicaIDAnnotationsKey,
// 				readReplicaEnabled:                test.fields.readReplicaEnabled,
// 				readReplicaLabelKey:               test.fields.readReplicaLabelKey,
// 				rotationJobConcurrency:            test.fields.rotationJobConcurrency,
// 				rotatorJob:                        test.fields.rotatorJob,
// 			}
//
// 			gotRequeue, err := o.reconcileRotatorJob(test.args.ctx, test.args.pod)
// 			if err := checkFunc(test.want, gotRequeue, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_needsRotation(t *testing.T) {
// 	type args struct {
// 		agentAnnotations       map[string]string
// 		readReplicaAnnotations map[string]string
// 	}
// 	type want struct {
// 		want bool
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, bool, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got bool, err error) error {
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
// 		           agentAnnotations:nil,
// 		           readReplicaAnnotations:nil,
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
// 		           agentAnnotations:nil,
// 		           readReplicaAnnotations:nil,
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
// 			got, err := needsRotation(test.args.agentAnnotations, test.args.readReplicaAnnotations)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_operator_createRotationJobOrRequeue(t *testing.T) {
// 	type args struct {
// 		ctx    context.Context
// 		podIdx string
// 	}
// 	type fields struct {
// 		ctrl                              k8s.Controller
// 		eg                                errgroup.Group
// 		namespace                         string
// 		client                            client.Client
// 		rotatorName                       string
// 		targetReadReplicaIDAnnotationsKey string
// 		readReplicaEnabled                bool
// 		readReplicaLabelKey               string
// 		rotationJobConcurrency            uint
// 		rotatorJob                        *k8s.Job
// 	}
// 	type want struct {
// 		wantRq bool
// 		err    error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, bool, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotRq bool, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotRq, w.wantRq) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRq, w.wantRq)
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
// 		           podIdx:"",
// 		       },
// 		       fields: fields {
// 		           ctrl:nil,
// 		           eg:nil,
// 		           namespace:"",
// 		           client:nil,
// 		           rotatorName:"",
// 		           targetReadReplicaIDAnnotationsKey:"",
// 		           readReplicaEnabled:false,
// 		           readReplicaLabelKey:"",
// 		           rotationJobConcurrency:0,
// 		           rotatorJob:nil,
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
// 		           podIdx:"",
// 		           },
// 		           fields: fields {
// 		           ctrl:nil,
// 		           eg:nil,
// 		           namespace:"",
// 		           client:nil,
// 		           rotatorName:"",
// 		           targetReadReplicaIDAnnotationsKey:"",
// 		           readReplicaEnabled:false,
// 		           readReplicaLabelKey:"",
// 		           rotationJobConcurrency:0,
// 		           rotatorJob:nil,
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
// 			o := &operator{
// 				ctrl:                              test.fields.ctrl,
// 				eg:                                test.fields.eg,
// 				namespace:                         test.fields.namespace,
// 				client:                            test.fields.client,
// 				rotatorName:                       test.fields.rotatorName,
// 				targetReadReplicaIDAnnotationsKey: test.fields.targetReadReplicaIDAnnotationsKey,
// 				readReplicaEnabled:                test.fields.readReplicaEnabled,
// 				readReplicaLabelKey:               test.fields.readReplicaLabelKey,
// 				rotationJobConcurrency:            test.fields.rotationJobConcurrency,
// 				rotatorJob:                        test.fields.rotatorJob,
// 			}
//
// 			gotRq, err := o.createRotationJobOrRequeue(test.args.ctx, test.args.podIdx)
// 			if err := checkFunc(test.want, gotRq, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_operator_ensureJobConcurrency(t *testing.T) {
// 	type args struct {
// 		ctx    context.Context
// 		podIdx string
// 	}
// 	type fields struct {
// 		ctrl                              k8s.Controller
// 		eg                                errgroup.Group
// 		namespace                         string
// 		client                            client.Client
// 		rotatorName                       string
// 		targetReadReplicaIDAnnotationsKey string
// 		readReplicaEnabled                bool
// 		readReplicaLabelKey               string
// 		rotationJobConcurrency            uint
// 		rotatorJob                        *k8s.Job
// 	}
// 	type want struct {
// 		want jobReconcileResult
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, jobReconcileResult, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got jobReconcileResult, err error) error {
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
// 		           podIdx:"",
// 		       },
// 		       fields: fields {
// 		           ctrl:nil,
// 		           eg:nil,
// 		           namespace:"",
// 		           client:nil,
// 		           rotatorName:"",
// 		           targetReadReplicaIDAnnotationsKey:"",
// 		           readReplicaEnabled:false,
// 		           readReplicaLabelKey:"",
// 		           rotationJobConcurrency:0,
// 		           rotatorJob:nil,
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
// 		           podIdx:"",
// 		           },
// 		           fields: fields {
// 		           ctrl:nil,
// 		           eg:nil,
// 		           namespace:"",
// 		           client:nil,
// 		           rotatorName:"",
// 		           targetReadReplicaIDAnnotationsKey:"",
// 		           readReplicaEnabled:false,
// 		           readReplicaLabelKey:"",
// 		           rotationJobConcurrency:0,
// 		           rotatorJob:nil,
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
// 			o := &operator{
// 				ctrl:                              test.fields.ctrl,
// 				eg:                                test.fields.eg,
// 				namespace:                         test.fields.namespace,
// 				client:                            test.fields.client,
// 				rotatorName:                       test.fields.rotatorName,
// 				targetReadReplicaIDAnnotationsKey: test.fields.targetReadReplicaIDAnnotationsKey,
// 				readReplicaEnabled:                test.fields.readReplicaEnabled,
// 				readReplicaLabelKey:               test.fields.readReplicaLabelKey,
// 				rotationJobConcurrency:            test.fields.rotationJobConcurrency,
// 				rotatorJob:                        test.fields.rotatorJob,
// 			}
//
// 			got, err := o.ensureJobConcurrency(test.args.ctx, test.args.podIdx)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
