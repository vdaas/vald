//go:build e2e

//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
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

// package kubernetes provides kubernetes e2e tests
package kubernetes

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
)

type ResourceType[T metav1.Object] interface {
	Get(ctx context.Context, client kubernetes.Interface) (T, error)
	Update(ctx context.Context, client kubernetes.Interface, resource T) (T, error)
}

var (
	_ ResourceType[*appsv1.Deployment]  = (*Deployment)(nil)
	_ ResourceType[*appsv1.StatefulSet] = (*StatefulSet)(nil)
	_ ResourceType[*appsv1.DaemonSet]   = (*DaemonSet)(nil)
	_ ResourceType[*batchv1.Job]        = (*Job)(nil)
	_ ResourceType[*batchv1.CronJob]    = (*CronJob)(nil)
)

type Deployment struct {
	Name      string
	Namespace string
}

func (d Deployment) Get(ctx context.Context, client kubernetes.Interface) (*appsv1.Deployment, error) {
	return client.AppsV1().Deployments(d.Namespace).Get(ctx, d.Name, metav1.GetOptions{})
}

func (d Deployment) Update(ctx context.Context, client kubernetes.Interface, resource *appsv1.Deployment) (*appsv1.Deployment, error) {
	return client.AppsV1().Deployments(d.Namespace).Update(ctx, resource, metav1.UpdateOptions{})
}

type DaemonSet struct {
	Name      string
	NameSpace string
}

func (d DaemonSet) Get(ctx context.Context, client kubernetes.Interface) (*appsv1.DaemonSet, error) {
	return client.AppsV1().DaemonSets(d.NameSpace).Get(ctx, d.Name, metav1.GetOptions{})
}

func (d DaemonSet) Update(ctx context.Context, client kubernetes.Interface, resource *appsv1.DaemonSet) (*appsv1.DaemonSet, error) {
	return client.AppsV1().DaemonSets(d.NameSpace).Update(ctx, resource, metav1.UpdateOptions{})
}

type StatefulSet struct {
	Name      string
	Namespace string
}

func (s StatefulSet) Get(ctx context.Context, client kubernetes.Interface) (*appsv1.StatefulSet, error) {
	return client.AppsV1().StatefulSets(s.Namespace).Get(ctx, s.Name, metav1.GetOptions{})
}

func (s StatefulSet) Update(ctx context.Context, client kubernetes.Interface, resource *appsv1.StatefulSet) (*appsv1.StatefulSet, error) {
	return client.AppsV1().StatefulSets(s.Namespace).Update(ctx, resource, metav1.UpdateOptions{})
}

type Job struct {
	Name      string
	Namespace string
}

func (j Job) Get(ctx context.Context, client kubernetes.Interface) (*batchv1.Job, error) {
	return client.BatchV1().Jobs(j.Namespace).Get(ctx, j.Name, metav1.GetOptions{})
}

func (j Job) Update(ctx context.Context, client kubernetes.Interface, resource *batchv1.Job) (*batchv1.Job, error) {
	return client.BatchV1().Jobs(j.Namespace).Update(ctx, resource, metav1.UpdateOptions{})
}

type CronJob struct {
	Name      string
	Namespace string
}

func (cj CronJob) Get(ctx context.Context, client kubernetes.Interface) (*batchv1.CronJob, error) {
	return client.BatchV1().CronJobs(cj.Namespace).Get(ctx, cj.Name, metav1.GetOptions{})
}

func (cj CronJob) Update(ctx context.Context, client kubernetes.Interface, resource *batchv1.CronJob) (*batchv1.CronJob, error) {
	return client.BatchV1().CronJobs(cj.Namespace).Update(ctx, resource, metav1.UpdateOptions{})
}
