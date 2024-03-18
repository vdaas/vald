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

	tmock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/vdaas/vald/internal/k8s/client"
	"github.com/vdaas/vald/internal/k8s/vald"
	"github.com/vdaas/vald/internal/test/mock/k8s"
	"github.com/vdaas/vald/internal/test/testify"

	batchv1 "k8s.io/api/batch/v1"
)

func Test_operator_podOnReconcile(t *testing.T) {
	t.Parallel()

	type test struct {
		name     string
		testfunc func(t *testing.T)
	}

	tests := []test{
		{
			"returns client.Result{} when read replica is not enabled",
			func(t *testing.T) {
				mock := &k8s.ValdK8sClientMock{}
				o, err := New(
					"namespace",
					"agentName",
					"rotatorName",
					"targetReadReplicaIDEnvName",
					nil,
					WithK8sClient(mock),
					WithReadReplicaEnabled(false),
				)
				require.NoError(t, err)

				op, ok := o.(*operator)
				require.True(t, ok)

				res, err := op.podOnReconcile(context.Background(), nil)
				require.NoError(t, err)

				require.Equal(t, client.Result{}, res)
			},
		},
		{
			"returns client.Result{} when pod is not a statefulset",
			func(t *testing.T) {
				mock := &k8s.ValdK8sClientMock{}
				o, err := New(
					"namespace",
					"agentName",
					"rotatorName",
					"targetReadReplicaIDEnvName",
					nil,
					WithK8sClient(mock),
					WithReadReplicaEnabled(true),
				)
				require.NoError(t, err)

				op, ok := o.(*operator)
				require.True(t, ok)

				pod := &client.Pod{}
				res, err := op.podOnReconcile(context.Background(), pod)
				require.NoError(t, err)

				require.Equal(t, client.Result{}, res)
			},
		},
		{
			"returns requeue: false when last snapshot time is after the last save time",
			func(t *testing.T) {
				saveTime := time.Now()
				rotateTime := saveTime.Add(1 * time.Second)

				mock := &k8s.ValdK8sClientMock{}
				mock.On("LabelSelector", testify.Anything, testify.Anything, testify.Anything).Return(client.NewSelector(), nil)
				mock.On("List", testify.Anything, testify.AnythingOfType("*v1.DeploymentList"), testify.Anything).Run(func(args tmock.Arguments) {
					arg, ok := args.Get(1).(*client.DeploymentList)
					require.True(t, ok)

					arg.Items = []client.Deployment{
						{
							ObjectMeta: client.ObjectMeta{
								Name: "deploymentName",
								Annotations: map[string]string{
									vald.LastTimeSnapshotTimestampAnnotationsKey: rotateTime.Format(vald.TimeFormat),
								},
							},
						},
					}
				}).Return(nil)

				agentPod := &client.Pod{
					ObjectMeta: client.ObjectMeta{
						Labels: map[string]string{
							client.PodIndexLabel: "0",
						},
						Annotations: map[string]string{
							vald.LastTimeSaveIndexTimestampAnnotationsKey: saveTime.Format(vald.TimeFormat),
						},
					},
				}

				o, err := New(
					"namespace",
					"agentName",
					"rotatorName",
					"targetReadReplicaIDEnvName",
					nil,
					WithK8sClient(mock),
					WithReadReplicaEnabled(true),
				)
				require.NoError(t, err)

				op, ok := o.(*operator)
				require.True(t, ok)

				res, err := op.podOnReconcile(context.Background(), agentPod)
				require.NoError(t, err)

				require.Equal(t, client.Result{
					Requeue: false,
				}, res)
			},
		},
		{
			"returns requeue: false and calls client.Create once when last snapshot time is before the last save time",
			func(t *testing.T) {
				saveTime := time.Now()
				rotateTime := saveTime.Add(-1 * time.Second)

				mock := &k8s.ValdK8sClientMock{}
				mock.On("LabelSelector", testify.Anything, testify.Anything, testify.Anything).Return(client.NewSelector(), nil)
				mock.On("List", testify.Anything, testify.AnythingOfType("*v1.DeploymentList"), testify.Anything).Run(func(args tmock.Arguments) {
					arg, ok := args.Get(1).(*client.DeploymentList)
					require.True(t, ok)

					arg.Items = []client.Deployment{
						{
							ObjectMeta: client.ObjectMeta{
								Name: "deploymentName",
								Annotations: map[string]string{
									vald.LastTimeSnapshotTimestampAnnotationsKey: rotateTime.Format(vald.TimeFormat),
								},
							},
						},
					}
				}).Return(nil)

				mock.On("List", testify.Anything, testify.AnythingOfType("*v1.JobList"), testify.Anything).Run(func(args tmock.Arguments) {
					arg, ok := args.Get(1).(*client.JobList)
					require.True(t, ok)

					arg.Items = []client.Job{}
				}).Return(nil)

				mock.On("Create", testify.Anything, testify.Anything, testify.Anything).Return(nil).Once()

				agentPod := &client.Pod{
					ObjectMeta: client.ObjectMeta{
						Labels: map[string]string{
							client.PodIndexLabel: "0",
						},
						Annotations: map[string]string{
							vald.LastTimeSaveIndexTimestampAnnotationsKey: saveTime.Format(vald.TimeFormat),
						},
					},
				}

				o, err := New(
					"namespace",
					"agentName",
					"rotatorName",
					"targetReadReplicaIDEnvName",
					nil,
					WithK8sClient(mock),
					WithReadReplicaEnabled(true),
				)
				require.NoError(t, err)

				op, ok := o.(*operator)
				require.True(t, ok)

				op.rotatorJob = &client.Job{
					ObjectMeta: client.ObjectMeta{
						Name: "foo job",
					},
				}

				res, err := op.podOnReconcile(context.Background(), agentPod)
				require.NoError(t, err)

				require.Equal(t, client.Result{
					Requeue: false,
				}, res)
			},
		},
		{
			"returns requeue: true when there is already one running job when rotation job concurrency is 1",
			func(t *testing.T) {
				saveTime := time.Now()
				rotateTime := saveTime.Add(-1 * time.Second)

				mock := &k8s.ValdK8sClientMock{}
				mock.On("LabelSelector", testify.Anything, testify.Anything, testify.Anything).Return(client.NewSelector(), nil)
				mock.On("List", testify.Anything, testify.AnythingOfType("*v1.DeploymentList"), testify.Anything).Run(func(args tmock.Arguments) {
					arg, ok := args.Get(1).(*client.DeploymentList)
					require.True(t, ok)

					arg.Items = []client.Deployment{
						{
							ObjectMeta: client.ObjectMeta{
								Name: "deploymentName",
								Annotations: map[string]string{
									vald.LastTimeSnapshotTimestampAnnotationsKey: rotateTime.Format(vald.TimeFormat),
								},
							},
						},
					}
				}).Return(nil)

				mock.On("List", testify.Anything, testify.AnythingOfType("*v1.JobList"), testify.Anything).Run(func(args tmock.Arguments) {
					arg, ok := args.Get(1).(*client.JobList)
					require.True(t, ok)

					arg.Items = []client.Job{
						{
							ObjectMeta: client.ObjectMeta{
								Name: "already running job1",
							},
							Status: batchv1.JobStatus{
								Active: 1,
							},
						},
					}
				}).Return(nil)

				mock.On("Create", testify.Anything, testify.Anything, testify.Anything).Return(nil).Once()

				agentPod := &client.Pod{
					ObjectMeta: client.ObjectMeta{
						Labels: map[string]string{
							client.PodIndexLabel: "0",
						},
						Annotations: map[string]string{
							vald.LastTimeSaveIndexTimestampAnnotationsKey: saveTime.Format(vald.TimeFormat),
						},
					},
				}

				o, err := New(
					"namespace",
					"agentName",
					"rotatorName",
					"targetReadReplicaIDEnvName",
					nil,
					WithK8sClient(mock),
					WithReadReplicaEnabled(true),
				)
				require.NoError(t, err)

				op, ok := o.(*operator)
				require.True(t, ok)

				op.rotatorJob = &client.Job{
					ObjectMeta: client.ObjectMeta{
						Name: "foo job",
					},
				}

				res, err := op.podOnReconcile(context.Background(), agentPod)
				require.NoError(t, err)

				require.Equal(t, client.Result{
					Requeue: true,
				}, res)
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			test.testfunc(tt)
		})
	}
}
