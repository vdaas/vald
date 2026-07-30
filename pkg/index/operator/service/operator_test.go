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
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/vdaas/vald/internal/k8s"
	"github.com/vdaas/vald/internal/k8s/resource"
	"github.com/vdaas/vald/internal/k8s/vald"
	mock "github.com/vdaas/vald/internal/test/mock/k8s"
)

//nolint:maintidx // table-driven test with multiple literal test-case constructors; complexity is inherent to the pattern, not tangled logic
func Test_operator_podOnReconcile(t *testing.T) {
	t.Parallel()

	const (
		readReplicaLabelKey = "app"
		rotatorName         = "rotator"
		testDeploymentName  = "deploymentName"
		testRunningJob1     = "already running job1"
	)

	type want struct {
		res          k8s.Result
		err          error
		createCalled bool
	}
	type test struct {
		want                   want
		agentPod               *k8s.Pod
		readReplicaDeployment  *k8s.Deployment
		name                   string
		runningJobs            []k8s.Job
		rotationJobConcurrency uint
		readReplicaEnabled     bool
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
						Name:   testDeploymentName,
						Labels: map[string]string{readReplicaLabelKey: "0"},
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
						Name:   testDeploymentName,
						Labels: map[string]string{readReplicaLabelKey: "0"},
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
						Name:   testDeploymentName,
						Labels: map[string]string{readReplicaLabelKey: "0"},
						Annotations: map[string]string{
							vald.LastTimeSnapshotTimestampAnnotationsKey: rotateTime.Format(vald.TimeFormat),
						},
					},
				},
				runningJobs: []k8s.Job{
					{
						ObjectMeta: k8s.ObjectMeta{
							Name:   testRunningJob1,
							Labels: map[string]string{readReplicaLabelKey: rotatorName},
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
						Name:   testDeploymentName,
						Labels: map[string]string{readReplicaLabelKey: "0"},
						Annotations: map[string]string{
							vald.LastTimeSnapshotTimestampAnnotationsKey: rotateTime.Format(vald.TimeFormat),
						},
					},
				},
				runningJobs: []k8s.Job{
					{
						ObjectMeta: k8s.ObjectMeta{
							Name:   testRunningJob1,
							Labels: map[string]string{readReplicaLabelKey: rotatorName},
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
						Name:   testDeploymentName,
						Labels: map[string]string{readReplicaLabelKey: "0"},
						Annotations: map[string]string{
							vald.LastTimeSnapshotTimestampAnnotationsKey: rotateTime.Format(vald.TimeFormat),
						},
					},
				},
				runningJobs: []k8s.Job{
					{
						ObjectMeta: k8s.ObjectMeta{
							Name:   testRunningJob1,
							Labels: map[string]string{readReplicaLabelKey: rotatorName},
						},
						Status: k8s.JobStatus{
							Active: 1,
						},
					},
					{
						ObjectMeta: k8s.ObjectMeta{
							Name:   "already running job2",
							Labels: map[string]string{readReplicaLabelKey: rotatorName},
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

			mockClient := &mock.ValdK8sClientMock{}

			scheme := k8s.NewScheme()
			require.NoError(t, k8s.AddClientGoScheme(scheme))
			var initObjs []k8s.Object
			if test.readReplicaDeployment != nil {
				initObjs = append(initObjs, test.readReplicaDeployment)
			}
			for i := range test.runningJobs {
				initObjs = append(initObjs, &test.runningJobs[i])
			}

			// testify/mock does not accept to set Times(0), so the call count is
			// tracked manually via an interceptor instead of a mock expectation.
			var createCalls int
			base := mock.NewFakeClientBuilder().WithScheme(scheme).WithObjects(initObjs...).Build()
			fakeClient := mock.NewInterceptorClient(base, mock.InterceptorFuncs{
				Create: func(
					ctx context.Context, c k8s.WithWatch, obj k8s.Object, opts ...k8s.CreateOption,
				) error {
					createCalls++
					return c.Create(ctx, obj, opts...)
				},
			})

			concurrency := uint(1)
			if test.rotationJobConcurrency != 0 {
				concurrency = test.rotationJobConcurrency
			}
			op := operator{
				client:                 mockClient,
				deployments:            resource.NewClient(fakeClient, new(k8s.Deployment), new(k8s.DeploymentList)),
				jobs:                   resource.NewClient(fakeClient, new(k8s.Job), new(k8s.JobList)),
				readReplicaLabelKey:    readReplicaLabelKey,
				rotatorName:            rotatorName,
				readReplicaEnabled:     test.readReplicaEnabled,
				rotationJobConcurrency: concurrency,
			}

			op.rotatorJob = &k8s.Job{
				ObjectMeta: k8s.ObjectMeta{
					Name: "foo job",
				},
			}

			res, err := op.podOnReconcile(tt.Context(), test.agentPod)
			require.Equal(t, test.want.err, err)
			require.Equal(t, test.want.res, res)
			wantCreateCalls := 0
			if test.want.createCalled {
				wantCreateCalls = 1
			}
			require.Equal(t, wantCreateCalls, createCalls)
		})
	}
}

func Test_operator_ensureJobConcurrency(t *testing.T) {
	t.Parallel()

	const (
		rotatorName = "rotator"
		namespace   = "default"
		targetKey   = "vald.vdaas.org/target-read-replica-id"
		podIdx      = "0"
		otherPodIdx = "1"
	)

	// newJob builds a rotation job fixture that matches the label selector
	// (app=rotatorName) and namespace used by ensureJobConcurrency. An empty
	// idx leaves the pod template annotations nil.
	newJob := func(name string, status k8s.JobStatus, idx string) k8s.Job {
		job := k8s.Job{
			ObjectMeta: k8s.ObjectMeta{
				Name:      name,
				Namespace: namespace,
				Labels:    map[string]string{"app": rotatorName},
			},
			Status: status,
		}
		if idx != "" {
			job.Spec.Template.Annotations = map[string]string{targetKey: idx}
		}
		return job
	}

	type test struct {
		name                   string
		jobs                   []k8s.Job
		rotationJobConcurrency uint
		want                   jobReconcileResult
	}

	tests := []test{
		{
			name:                   "returns createRequired when no jobs exist",
			rotationJobConcurrency: 1,
			want:                   createRequired,
		},
		{
			// Regression: a just-created job has Active==0 with zero
			// Succeeded/Failed until the Job controller schedules its pod.
			// It must stay in the accounting so that a second reconcile in
			// that window dedups instead of creating a duplicate rotation job.
			name: "returns createSkipped when a brand-new job (all zero status) has the same podIdx",
			jobs: []k8s.Job{
				newJob("brand-new", k8s.JobStatus{}, podIdx),
			},
			rotationJobConcurrency: 2,
			want:                   createSkipped,
		},
		{
			name: "returns requeueRequired when a brand-new job (all zero status) fills the concurrency limit",
			jobs: []k8s.Job{
				newJob("brand-new", k8s.JobStatus{}, otherPodIdx),
			},
			rotationJobConcurrency: 1,
			want:                   requeueRequired,
		},
		{
			name: "returns createRequired when a succeeded job is dropped from accounting",
			jobs: []k8s.Job{
				newJob("succeeded", k8s.JobStatus{Succeeded: 1}, podIdx),
			},
			rotationJobConcurrency: 1,
			want:                   createRequired,
		},
		{
			name: "returns createRequired when a failed job is dropped from accounting",
			jobs: []k8s.Job{
				newJob("failed", k8s.JobStatus{Failed: 1}, podIdx),
			},
			rotationJobConcurrency: 1,
			want:                   createRequired,
		},
		{
			name: "returns requeueRequired when a running job with a different podIdx fills the concurrency limit",
			jobs: []k8s.Job{
				newJob("running", k8s.JobStatus{Active: 1}, otherPodIdx),
			},
			rotationJobConcurrency: 1,
			want:                   requeueRequired,
		},
		{
			name: "returns createSkipped when a running job has the same podIdx",
			jobs: []k8s.Job{
				newJob("running", k8s.JobStatus{Active: 1}, podIdx),
			},
			rotationJobConcurrency: 2,
			want:                   createSkipped,
		},
		{
			name: "returns createSkipped when a retrying job (Active>0 with Failed>0) has the same podIdx",
			jobs: []k8s.Job{
				newJob("retrying", k8s.JobStatus{Active: 1, Failed: 1}, podIdx),
			},
			rotationJobConcurrency: 2,
			want:                   createSkipped,
		},
		{
			name: "returns createRequired when a running job has no target annotation",
			jobs: []k8s.Job{
				newJob("running-no-annotation", k8s.JobStatus{Active: 1}, ""),
			},
			rotationJobConcurrency: 2,
			want:                   createRequired,
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()

			scheme := k8s.NewScheme()
			require.NoError(tt, k8s.AddClientGoScheme(scheme))
			var initObjs []k8s.Object
			for i := range test.jobs {
				initObjs = append(initObjs, &test.jobs[i])
			}
			fakeClient := mock.NewFakeClientBuilder().WithScheme(scheme).WithObjects(initObjs...).Build()

			op := operator{
				jobs:                              resource.NewClient(fakeClient, new(k8s.Job), new(k8s.JobList)),
				namespace:                         namespace,
				rotatorName:                       rotatorName,
				rotationJobConcurrency:            test.rotationJobConcurrency,
				targetReadReplicaIDAnnotationsKey: targetKey,
			}

			got, err := op.ensureJobConcurrency(tt.Context(), podIdx)
			require.NoError(tt, err)
			require.Equal(tt, test.want, got)
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
// 		want Operator
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, Operator, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got Operator, err error) error {
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
// 			got, err := New(test.args.namespace, test.args.agentName, test.args.rotatorName, test.args.targetReadReplicaIDKey, test.args.rotatorJob, test.args.opts...)
// 			if err := checkFunc(test.want, got, err); err != nil {
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
// 		client                            client.Client
// 		deployments                       *resource.Client[*k8s.Deployment, *k8s.DeploymentList]
// 		jobs                              *resource.Client[*k8s.Job, *k8s.JobList]
// 		rotatorJob                        *k8s.Job
// 		namespace                         string
// 		rotatorName                       string
// 		targetReadReplicaIDAnnotationsKey string
// 		readReplicaLabelKey               string
// 		rotationJobConcurrency            uint
// 		readReplicaEnabled                bool
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
// 		           client:nil,
// 		           deployments:nil,
// 		           jobs:nil,
// 		           rotatorJob:nil,
// 		           namespace:"",
// 		           rotatorName:"",
// 		           targetReadReplicaIDAnnotationsKey:"",
// 		           readReplicaLabelKey:"",
// 		           rotationJobConcurrency:0,
// 		           readReplicaEnabled:false,
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
// 		           client:nil,
// 		           deployments:nil,
// 		           jobs:nil,
// 		           rotatorJob:nil,
// 		           namespace:"",
// 		           rotatorName:"",
// 		           targetReadReplicaIDAnnotationsKey:"",
// 		           readReplicaLabelKey:"",
// 		           rotationJobConcurrency:0,
// 		           readReplicaEnabled:false,
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
// 				client:                            test.fields.client,
// 				deployments:                       test.fields.deployments,
// 				jobs:                              test.fields.jobs,
// 				rotatorJob:                        test.fields.rotatorJob,
// 				namespace:                         test.fields.namespace,
// 				rotatorName:                       test.fields.rotatorName,
// 				targetReadReplicaIDAnnotationsKey: test.fields.targetReadReplicaIDAnnotationsKey,
// 				readReplicaLabelKey:               test.fields.readReplicaLabelKey,
// 				rotationJobConcurrency:            test.fields.rotationJobConcurrency,
// 				readReplicaEnabled:                test.fields.readReplicaEnabled,
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
// 		client                            client.Client
// 		deployments                       *resource.Client[*k8s.Deployment, *k8s.DeploymentList]
// 		jobs                              *resource.Client[*k8s.Job, *k8s.JobList]
// 		rotatorJob                        *k8s.Job
// 		namespace                         string
// 		rotatorName                       string
// 		targetReadReplicaIDAnnotationsKey string
// 		readReplicaLabelKey               string
// 		rotationJobConcurrency            uint
// 		readReplicaEnabled                bool
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
// 		           client:nil,
// 		           deployments:nil,
// 		           jobs:nil,
// 		           rotatorJob:nil,
// 		           namespace:"",
// 		           rotatorName:"",
// 		           targetReadReplicaIDAnnotationsKey:"",
// 		           readReplicaLabelKey:"",
// 		           rotationJobConcurrency:0,
// 		           readReplicaEnabled:false,
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
// 		           client:nil,
// 		           deployments:nil,
// 		           jobs:nil,
// 		           rotatorJob:nil,
// 		           namespace:"",
// 		           rotatorName:"",
// 		           targetReadReplicaIDAnnotationsKey:"",
// 		           readReplicaLabelKey:"",
// 		           rotationJobConcurrency:0,
// 		           readReplicaEnabled:false,
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
// 				client:                            test.fields.client,
// 				deployments:                       test.fields.deployments,
// 				jobs:                              test.fields.jobs,
// 				rotatorJob:                        test.fields.rotatorJob,
// 				namespace:                         test.fields.namespace,
// 				rotatorName:                       test.fields.rotatorName,
// 				targetReadReplicaIDAnnotationsKey: test.fields.targetReadReplicaIDAnnotationsKey,
// 				readReplicaLabelKey:               test.fields.readReplicaLabelKey,
// 				rotationJobConcurrency:            test.fields.rotationJobConcurrency,
// 				readReplicaEnabled:                test.fields.readReplicaEnabled,
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
// 		client                            client.Client
// 		deployments                       *resource.Client[*k8s.Deployment, *k8s.DeploymentList]
// 		jobs                              *resource.Client[*k8s.Job, *k8s.JobList]
// 		rotatorJob                        *k8s.Job
// 		namespace                         string
// 		rotatorName                       string
// 		targetReadReplicaIDAnnotationsKey string
// 		readReplicaLabelKey               string
// 		rotationJobConcurrency            uint
// 		readReplicaEnabled                bool
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
// 		           client:nil,
// 		           deployments:nil,
// 		           jobs:nil,
// 		           rotatorJob:nil,
// 		           namespace:"",
// 		           rotatorName:"",
// 		           targetReadReplicaIDAnnotationsKey:"",
// 		           readReplicaLabelKey:"",
// 		           rotationJobConcurrency:0,
// 		           readReplicaEnabled:false,
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
// 		           client:nil,
// 		           deployments:nil,
// 		           jobs:nil,
// 		           rotatorJob:nil,
// 		           namespace:"",
// 		           rotatorName:"",
// 		           targetReadReplicaIDAnnotationsKey:"",
// 		           readReplicaLabelKey:"",
// 		           rotationJobConcurrency:0,
// 		           readReplicaEnabled:false,
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
// 				client:                            test.fields.client,
// 				deployments:                       test.fields.deployments,
// 				jobs:                              test.fields.jobs,
// 				rotatorJob:                        test.fields.rotatorJob,
// 				namespace:                         test.fields.namespace,
// 				rotatorName:                       test.fields.rotatorName,
// 				targetReadReplicaIDAnnotationsKey: test.fields.targetReadReplicaIDAnnotationsKey,
// 				readReplicaLabelKey:               test.fields.readReplicaLabelKey,
// 				rotationJobConcurrency:            test.fields.rotationJobConcurrency,
// 				readReplicaEnabled:                test.fields.readReplicaEnabled,
// 			}
//
// 			gotRq, err := o.createRotationJobOrRequeue(test.args.ctx, test.args.podIdx)
// 			if err := checkFunc(test.want, gotRq, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
