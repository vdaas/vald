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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

type ResourceType[T metav1.Object] interface {
	Get(ctx context.Context, client kubernetes.Interface) (T, error)
	List(ctx context.Context, client kubernetes.Interface) ([]T, error)
	Update(ctx context.Context, client kubernetes.Interface, resource T) (T, error)
	Create(ctx context.Context, client kubernetes.Interface, resource T) (T, error)
	Patch(ctx context.Context, client kubernetes.Interface, data []byte) (T, error)
	Delete(ctx context.Context, client kubernetes.Interface) error
	Watch(ctx context.Context, client kubernetes.Interface) (watch.Interface, error)
}

var (
	_ ResourceType[*corev1.Pod]         = (*Pod)(nil)
	_ ResourceType[*appsv1.Deployment]  = (*Deployment)(nil)
	_ ResourceType[*appsv1.StatefulSet] = (*StatefulSet)(nil)
	_ ResourceType[*appsv1.DaemonSet]   = (*DaemonSet)(nil)
	_ ResourceType[*batchv1.Job]        = (*Job)(nil)
	_ ResourceType[*batchv1.CronJob]    = (*CronJob)(nil)
)

type baseResource struct {
	Name          string
	Namespace     string
	LabelSelector string
	FieldSelector string
}

func listConvert[T any](slice []T) []*T {
	result := make([]*T, len(slice))
	for i := range slice {
		result[i] = &slice[i]
	}
	return result
}

type Pod baseResource

func (p Pod) Get(ctx context.Context, client kubernetes.Interface) (*corev1.Pod, error) {
	return client.CoreV1().Pods(p.Namespace).Get(ctx, p.Name, metav1.GetOptions{})
}

func (p Pod) List(ctx context.Context, client kubernetes.Interface) ([]*corev1.Pod, error) {
	list, err := client.CoreV1().Pods(p.Namespace).List(ctx, metav1.ListOptions{
		LabelSelector: p.LabelSelector,
		FieldSelector: p.FieldSelector,
	})
	if err != nil {
		return nil, err
	}
	return listConvert(list.Items), nil
}

func (p Pod) Update(ctx context.Context, client kubernetes.Interface, resource *corev1.Pod) (*corev1.Pod, error) {
	return client.CoreV1().Pods(p.Namespace).Update(ctx, resource, metav1.UpdateOptions{})
}

func (p Pod) Create(ctx context.Context, client kubernetes.Interface, resource *corev1.Pod) (*corev1.Pod, error) {
	return client.CoreV1().Pods(p.Namespace).Create(ctx, resource, metav1.CreateOptions{})
}

func (p Pod) Patch(ctx context.Context, client kubernetes.Interface, data []byte) (*corev1.Pod, error) {
	return client.CoreV1().Pods(p.Namespace).Patch(ctx, p.Name, types.JSONPatchType, data, metav1.PatchOptions{})
}

func (p Pod) Delete(ctx context.Context, client kubernetes.Interface) error {
	return client.CoreV1().Pods(p.Namespace).Delete(ctx, p.Name, metav1.DeleteOptions{})
}

func (p Pod) Watch(ctx context.Context, client kubernetes.Interface) (watch.Interface, error) {
	return client.CoreV1().Pods(p.Namespace).Watch(ctx, metav1.ListOptions{
		LabelSelector: p.LabelSelector,
		FieldSelector: p.FieldSelector,
	})
}

type Deployment baseResource

func (d Deployment) Get(ctx context.Context, client kubernetes.Interface) (*appsv1.Deployment, error) {
	return client.AppsV1().Deployments(d.Namespace).Get(ctx, d.Name, metav1.GetOptions{})
}

func (d Deployment) List(ctx context.Context, client kubernetes.Interface) ([]*appsv1.Deployment, error) {
	list, err := client.AppsV1().Deployments(d.Namespace).List(ctx, metav1.ListOptions{
		LabelSelector: d.LabelSelector,
		FieldSelector: d.FieldSelector,
	})
	if err != nil {
		return nil, err
	}
	return listConvert(list.Items), nil
}

func (d Deployment) Update(ctx context.Context, client kubernetes.Interface, resource *appsv1.Deployment) (*appsv1.Deployment, error) {
	return client.AppsV1().Deployments(d.Namespace).Update(ctx, resource, metav1.UpdateOptions{})
}

func (d Deployment) Create(ctx context.Context, client kubernetes.Interface, resource *appsv1.Deployment) (*appsv1.Deployment, error) {
	return client.AppsV1().Deployments(d.Namespace).Create(ctx, resource, metav1.CreateOptions{})
}

func (d Deployment) Patch(ctx context.Context, client kubernetes.Interface, data []byte) (*appsv1.Deployment, error) {
	return client.AppsV1().Deployments(d.Namespace).Patch(ctx, d.Name, types.JSONPatchType, data, metav1.PatchOptions{})
}

func (d Deployment) Delete(ctx context.Context, client kubernetes.Interface) error {
	return client.AppsV1().Deployments(d.Namespace).Delete(ctx, d.Name, metav1.DeleteOptions{})
}

func (d Deployment) Watch(ctx context.Context, client kubernetes.Interface) (watch.Interface, error) {
	return client.AppsV1().Deployments(d.Namespace).Watch(ctx, metav1.ListOptions{
		LabelSelector: d.LabelSelector,
		FieldSelector: d.FieldSelector,
	})
}

type DaemonSet baseResource

func (d DaemonSet) Get(ctx context.Context, client kubernetes.Interface) (*appsv1.DaemonSet, error) {
	return client.AppsV1().DaemonSets(d.Namespace).Get(ctx, d.Name, metav1.GetOptions{})
}

func (d DaemonSet) List(ctx context.Context, client kubernetes.Interface) ([]*appsv1.DaemonSet, error) {
	list, err := client.AppsV1().DaemonSets(d.Namespace).List(ctx, metav1.ListOptions{
		LabelSelector: d.LabelSelector,
		FieldSelector: d.FieldSelector,
	})
	if err != nil {
		return nil, err
	}
	return listConvert(list.Items), nil
}

func (d DaemonSet) Update(ctx context.Context, client kubernetes.Interface, resource *appsv1.DaemonSet) (*appsv1.DaemonSet, error) {
	return client.AppsV1().DaemonSets(d.Namespace).Update(ctx, resource, metav1.UpdateOptions{})
}

func (d DaemonSet) Create(ctx context.Context, client kubernetes.Interface, resource *appsv1.DaemonSet) (*appsv1.DaemonSet, error) {
	return client.AppsV1().DaemonSets(d.Namespace).Create(ctx, resource, metav1.CreateOptions{})
}

func (d DaemonSet) Patch(ctx context.Context, client kubernetes.Interface, data []byte) (*appsv1.DaemonSet, error) {
	return client.AppsV1().DaemonSets(d.Namespace).Patch(ctx, d.Name, types.JSONPatchType, data, metav1.PatchOptions{})
}

func (d DaemonSet) Delete(ctx context.Context, client kubernetes.Interface) error {
	return client.AppsV1().DaemonSets(d.Namespace).Delete(ctx, d.Name, metav1.DeleteOptions{})
}

func (d DaemonSet) Watch(ctx context.Context, client kubernetes.Interface) (watch.Interface, error) {
	return client.AppsV1().DaemonSets(d.Namespace).Watch(ctx, metav1.ListOptions{
		LabelSelector: d.LabelSelector,
		FieldSelector: d.FieldSelector,
	})
}

type StatefulSet baseResource

func (s StatefulSet) Get(ctx context.Context, client kubernetes.Interface) (*appsv1.StatefulSet, error) {
	return client.AppsV1().StatefulSets(s.Namespace).Get(ctx, s.Name, metav1.GetOptions{})
}

func (s StatefulSet) List(ctx context.Context, client kubernetes.Interface) ([]*appsv1.StatefulSet, error) {
	list, err := client.AppsV1().StatefulSets(s.Namespace).List(ctx, metav1.ListOptions{
		LabelSelector: s.LabelSelector,
		FieldSelector: s.FieldSelector,
	})
	if err != nil {
		return nil, err
	}
	return listConvert(list.Items), nil
}

func (s StatefulSet) Update(ctx context.Context, client kubernetes.Interface, resource *appsv1.StatefulSet) (*appsv1.StatefulSet, error) {
	return client.AppsV1().StatefulSets(s.Namespace).Update(ctx, resource, metav1.UpdateOptions{})
}

func (s StatefulSet) Create(ctx context.Context, client kubernetes.Interface, resource *appsv1.StatefulSet) (*appsv1.StatefulSet, error) {
	return client.AppsV1().StatefulSets(s.Namespace).Create(ctx, resource, metav1.CreateOptions{})
}

func (s StatefulSet) Patch(ctx context.Context, client kubernetes.Interface, data []byte) (*appsv1.StatefulSet, error) {
	return client.AppsV1().StatefulSets(s.Namespace).Patch(ctx, s.Name, types.JSONPatchType, data, metav1.PatchOptions{})
}

func (s StatefulSet) Delete(ctx context.Context, client kubernetes.Interface) error {
	return client.AppsV1().StatefulSets(s.Namespace).Delete(ctx, s.Name, metav1.DeleteOptions{})
}

func (s StatefulSet) Watch(ctx context.Context, client kubernetes.Interface) (watch.Interface, error) {
	return client.AppsV1().StatefulSets(s.Namespace).Watch(ctx, metav1.ListOptions{
		LabelSelector: s.LabelSelector,
		FieldSelector: s.FieldSelector,
	})
}

type Job baseResource

func (j Job) Get(ctx context.Context, client kubernetes.Interface) (*batchv1.Job, error) {
	return client.BatchV1().Jobs(j.Namespace).Get(ctx, j.Name, metav1.GetOptions{})
}

func (j Job) List(ctx context.Context, client kubernetes.Interface) ([]*batchv1.Job, error) {
	list, err := client.BatchV1().Jobs(j.Namespace).List(ctx, metav1.ListOptions{
		LabelSelector: j.LabelSelector,
		FieldSelector: j.FieldSelector,
	})
	if err != nil {
		return nil, err
	}
	return listConvert(list.Items), nil
}

func (j Job) Update(ctx context.Context, client kubernetes.Interface, resource *batchv1.Job) (*batchv1.Job, error) {
	return client.BatchV1().Jobs(j.Namespace).Update(ctx, resource, metav1.UpdateOptions{})
}

func (j Job) Create(ctx context.Context, client kubernetes.Interface, resource *batchv1.Job) (*batchv1.Job, error) {
	return client.BatchV1().Jobs(j.Namespace).Create(ctx, resource, metav1.CreateOptions{})
}

func (j Job) Patch(ctx context.Context, client kubernetes.Interface, data []byte) (*batchv1.Job, error) {
	return client.BatchV1().Jobs(j.Namespace).Patch(ctx, j.Name, types.JSONPatchType, data, metav1.PatchOptions{})
}

func (j Job) Delete(ctx context.Context, client kubernetes.Interface) error {
	return client.BatchV1().Jobs(j.Namespace).Delete(ctx, j.Name, metav1.DeleteOptions{})
}

func (j Job) Watch(ctx context.Context, client kubernetes.Interface) (watch.Interface, error) {
	return client.BatchV1().Jobs(j.Namespace).Watch(ctx, metav1.ListOptions{
		LabelSelector: j.LabelSelector,
		FieldSelector: j.FieldSelector,
	})
}

type CronJob baseResource

func (cj CronJob) Get(ctx context.Context, client kubernetes.Interface) (*batchv1.CronJob, error) {
	return client.BatchV1().CronJobs(cj.Namespace).Get(ctx, cj.Name, metav1.GetOptions{})
}

func (cj CronJob) List(ctx context.Context, client kubernetes.Interface) ([]*batchv1.CronJob, error) {
	list, err := client.BatchV1().CronJobs(cj.Namespace).List(ctx, metav1.ListOptions{
		LabelSelector: cj.LabelSelector,
		FieldSelector: cj.FieldSelector,
	})
	if err != nil {
		return nil, err
	}
	return listConvert(list.Items), nil
}

func (cj CronJob) Update(ctx context.Context, client kubernetes.Interface, resource *batchv1.CronJob) (*batchv1.CronJob, error) {
	return client.BatchV1().CronJobs(cj.Namespace).Update(ctx, resource, metav1.UpdateOptions{})
}

func (cj CronJob) Create(ctx context.Context, client kubernetes.Interface, resource *batchv1.CronJob) (*batchv1.CronJob, error) {
	return client.BatchV1().CronJobs(cj.Namespace).Create(ctx, resource, metav1.CreateOptions{})
}

func (cj CronJob) Patch(ctx context.Context, client kubernetes.Interface, data []byte) (*batchv1.CronJob, error) {
	return client.BatchV1().CronJobs(cj.Namespace).Patch(ctx, cj.Name, types.JSONPatchType, data, metav1.PatchOptions{})
}

func (cj CronJob) Delete(ctx context.Context, client kubernetes.Interface) error {
	return client.BatchV1().CronJobs(cj.Namespace).Delete(ctx, cj.Name, metav1.DeleteOptions{})
}

func (cj CronJob) Watch(ctx context.Context, client kubernetes.Interface) (watch.Interface, error) {
	return client.BatchV1().CronJobs(cj.Namespace).Watch(ctx, metav1.ListOptions{
		LabelSelector: cj.LabelSelector,
		FieldSelector: cj.FieldSelector,
	})
}
