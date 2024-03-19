//
// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
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
	"github.com/vdaas/vald/internal/k8s/client"
	"github.com/vdaas/vald/internal/k8s/vald"
	"github.com/vdaas/vald/internal/test/mock/k8s"
	"github.com/vdaas/vald/internal/test/testify"
)

func Test_operator_podOnReconcile(t *testing.T) {
	t.Parallel()

	type want struct {
		res          client.Result
		createCalled bool
		err          error
	}
	type test struct {
		name                   string
		agentPod               *client.Pod
		readReplicaEnabled     bool
		readReplicaDeployment  *client.Deployment
		runningJobs            []client.Job
		rotationJobConcurrency uint
		want                   want
	}

	tests := []test{
		{
			name:               "returns client.Result{} when read replica is not enabled",
			readReplicaEnabled: false,
			want: want{
				res:          client.Result{},
				createCalled: false,
				err:          nil,
			},
		},
		{
			name:               "returns client.Result{} when pod is not a statefulset",
			readReplicaEnabled: true,
			agentPod:           &client.Pod{},
			want: want{
				res:          client.Result{},
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
				agentPod: &client.Pod{
					ObjectMeta: client.ObjectMeta{
						Labels: map[string]string{
							client.PodIndexLabel: "0",
						},
						Annotations: map[string]string{
							vald.LastTimeSaveIndexTimestampAnnotationsKey: saveTime.Format(vald.TimeFormat),
						},
					},
				},
				readReplicaDeployment: &client.Deployment{
					ObjectMeta: client.ObjectMeta{
						Name: "deploymentName",
						Annotations: map[string]string{
							vald.LastTimeSnapshotTimestampAnnotationsKey: rotateTime.Format(vald.TimeFormat),
						},
					},
				},
				want: want{
					res: client.Result{
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
				agentPod: &client.Pod{
					ObjectMeta: client.ObjectMeta{
						Labels: map[string]string{
							client.PodIndexLabel: "0",
						},
						Annotations: map[string]string{
							vald.LastTimeSaveIndexTimestampAnnotationsKey: saveTime.Format(vald.TimeFormat),
						},
					},
				},
				readReplicaDeployment: &client.Deployment{
					ObjectMeta: client.ObjectMeta{
						Name: "deploymentName",
						Annotations: map[string]string{
							vald.LastTimeSnapshotTimestampAnnotationsKey: rotateTime.Format(vald.TimeFormat),
						},
					},
				},
				want: want{
					res: client.Result{
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
				agentPod: &client.Pod{
					ObjectMeta: client.ObjectMeta{
						Labels: map[string]string{
							client.PodIndexLabel: "0",
						},
						Annotations: map[string]string{
							vald.LastTimeSaveIndexTimestampAnnotationsKey: saveTime.Format(vald.TimeFormat),
						},
					},
				},
				readReplicaDeployment: &client.Deployment{
					ObjectMeta: client.ObjectMeta{
						Name: "deploymentName",
						Annotations: map[string]string{
							vald.LastTimeSnapshotTimestampAnnotationsKey: rotateTime.Format(vald.TimeFormat),
						},
					},
				},
				runningJobs: []client.Job{
					{
						ObjectMeta: client.ObjectMeta{
							Name: "already running job1",
						},
						Status: client.JobStatus{
							Active: 1,
						},
					},
				},
				rotationJobConcurrency: 1,
				want: want{
					res: client.Result{
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
				agentPod: &client.Pod{
					ObjectMeta: client.ObjectMeta{
						Labels: map[string]string{
							client.PodIndexLabel: "0",
						},
						Annotations: map[string]string{
							vald.LastTimeSaveIndexTimestampAnnotationsKey: saveTime.Format(vald.TimeFormat),
						},
					},
				},
				readReplicaDeployment: &client.Deployment{
					ObjectMeta: client.ObjectMeta{
						Name: "deploymentName",
						Annotations: map[string]string{
							vald.LastTimeSnapshotTimestampAnnotationsKey: rotateTime.Format(vald.TimeFormat),
						},
					},
				},
				runningJobs: []client.Job{
					{
						ObjectMeta: client.ObjectMeta{
							Name: "already running job1",
						},
						Status: client.JobStatus{
							Active: 1,
						},
					},
				},
				rotationJobConcurrency: 2,
				want: want{
					res: client.Result{
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
				agentPod: &client.Pod{
					ObjectMeta: client.ObjectMeta{
						Labels: map[string]string{
							client.PodIndexLabel: "0",
						},
						Annotations: map[string]string{
							vald.LastTimeSaveIndexTimestampAnnotationsKey: saveTime.Format(vald.TimeFormat),
						},
					},
				},
				readReplicaDeployment: &client.Deployment{
					ObjectMeta: client.ObjectMeta{
						Name: "deploymentName",
						Annotations: map[string]string{
							vald.LastTimeSnapshotTimestampAnnotationsKey: rotateTime.Format(vald.TimeFormat),
						},
					},
				},
				runningJobs: []client.Job{
					{
						ObjectMeta: client.ObjectMeta{
							Name: "already running job1",
						},
						Status: client.JobStatus{
							Active: 1,
						},
					},
					{
						ObjectMeta: client.ObjectMeta{
							Name: "already running job2",
						},
						Status: client.JobStatus{
							Active: 1,
						},
					},
				},
				rotationJobConcurrency: 2,
				want: want{
					res: client.Result{
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

			mock := &k8s.ValdK8sClientMock{}
			mock.On("LabelSelector", testify.Anything, testify.Anything, testify.Anything).Return(client.NewSelector(), nil).Maybe()
			mock.On("List", testify.Anything, testify.AnythingOfType("*v1.DeploymentList"), testify.Anything).Run(func(args testify.Arguments) {
				arg, ok := args.Get(1).(*client.DeploymentList)
				require.True(t, ok)

				arg.Items = []client.Deployment{*test.readReplicaDeployment}
			}).Return(nil).Maybe()

			mock.On("List", testify.Anything, testify.AnythingOfType("*v1.JobList"), testify.Anything).Run(func(args testify.Arguments) {
				arg, ok := args.Get(1).(*client.JobList)
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

			op.rotatorJob = &client.Job{
				ObjectMeta: client.ObjectMeta{
					Name: "foo job",
				},
			}

			res, err := op.podOnReconcile(context.Background(), test.agentPod)
			require.Equal(t, test.want.err, err)
			require.Equal(t, test.want.res, res)
		})
	}
}
