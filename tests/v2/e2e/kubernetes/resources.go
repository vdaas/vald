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

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/sync"
	appsv1 "k8s.io/api/apps/v1"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	applyconfigurationsappsv1 "k8s.io/client-go/applyconfigurations/apps/v1"
	applyconfigurationsautoscalingv1 "k8s.io/client-go/applyconfigurations/autoscaling/v1"
	applyconfigurationsbatchv1 "k8s.io/client-go/applyconfigurations/batch/v1"
	applyconfigurationscorev1 "k8s.io/client-go/applyconfigurations/core/v1"
	kclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type (
	Object      = kclient.Object
	ObjectList  = runtime.Object
	NamedObject interface {
		comparable
		GetName() *string
	}
)

type ObjectInterface[T Object, L ObjectList, C NamedObject] interface {
	Get(ctx context.Context, name string, opts metav1.GetOptions) (T, error)
	List(ctx context.Context, opts metav1.ListOptions) (L, error)

	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)

	Create(ctx context.Context, resource T, opts metav1.CreateOptions) (T, error)

	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error

	Update(ctx context.Context, resource T, opts metav1.UpdateOptions) (T, error)
	UpdateStatus(ctx context.Context, resource T, opts metav1.UpdateOptions) (T, error)

	Apply(ctx context.Context, resource C, opts metav1.ApplyOptions) (result T, err error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result T, err error)
}

type ScaleInterface interface {
	GetScale(ctx context.Context, resourceName string, options metav1.GetOptions) (*autoscalingv1.Scale, error)
	UpdateScale(ctx context.Context, resourceName string, scale *autoscalingv1.Scale, opts metav1.UpdateOptions) (*autoscalingv1.Scale, error)
	ApplyScale(ctx context.Context, resourceName string, scale *applyconfigurationsautoscalingv1.ScaleApplyConfiguration, opts metav1.ApplyOptions) (*autoscalingv1.Scale, error)
}

type PodExtendInterface interface {
	UpdateEphemeralContainers(ctx context.Context, podName string, pod *corev1.Pod, opts metav1.UpdateOptions) (*corev1.Pod, error)
	UpdateResize(ctx context.Context, podName string, pod *corev1.Pod, opts metav1.UpdateOptions) (*corev1.Pod, error)
}

type ClientControlInterface[T Object, L ObjectList, C NamedObject] interface {
	GetInterface() ObjectInterface[T, L, C]
	SetInterface(c ObjectInterface[T, L, C])

	GetClient() Client
	SetClient(c Client)

	GetNamespace() string
	SetNamespace(namespace string)
}

type ResourceClient[T Object, L ObjectList, C NamedObject] interface {
	ObjectInterface[T, L, C]
	ClientControlInterface[T, L, C]
}

type (
	PodClient interface {
		ResourceClient[*corev1.Pod, *corev1.PodList, *applyconfigurationscorev1.PodApplyConfiguration]
		PodExtendInterface
	}
	DeploymentClient interface {
		ResourceClient[*appsv1.Deployment, *appsv1.DeploymentList, *applyconfigurationsappsv1.DeploymentApplyConfiguration]
		ScaleInterface
	}
	DaemonSetClient   = ResourceClient[*appsv1.DaemonSet, *appsv1.DaemonSetList, *applyconfigurationsappsv1.DaemonSetApplyConfiguration]
	StatefulSetClient interface {
		ResourceClient[*appsv1.StatefulSet, *appsv1.StatefulSetList, *applyconfigurationsappsv1.StatefulSetApplyConfiguration]
		ScaleInterface
	}
	JobClient     = ResourceClient[*batchv1.Job, *batchv1.JobList, *applyconfigurationsbatchv1.JobApplyConfiguration]
	CronJobClient = ResourceClient[*batchv1.CronJob, *batchv1.CronJobList, *applyconfigurationsbatchv1.CronJobApplyConfiguration]
)

type (
	pod         = baseClient[*corev1.Pod, *corev1.PodList, *applyconfigurationscorev1.PodApplyConfiguration]
	deployment  = baseClient[*appsv1.Deployment, *appsv1.DeploymentList, *applyconfigurationsappsv1.DeploymentApplyConfiguration]
	daemonSet   = baseClient[*appsv1.DaemonSet, *appsv1.DaemonSetList, *applyconfigurationsappsv1.DaemonSetApplyConfiguration]
	statefulSet = baseClient[*appsv1.StatefulSet, *appsv1.StatefulSetList, *applyconfigurationsappsv1.StatefulSetApplyConfiguration]
	job         = baseClient[*batchv1.Job, *batchv1.JobList, *applyconfigurationsbatchv1.JobApplyConfiguration]
	cronJob     = baseClient[*batchv1.CronJob, *batchv1.CronJobList, *applyconfigurationsbatchv1.CronJobApplyConfiguration]
)

var (
	_ PodClient         = (*pod)(nil)
	_ DeploymentClient  = (*deployment)(nil)
	_ DaemonSetClient   = (*daemonSet)(nil)
	_ StatefulSetClient = (*statefulSet)(nil)
	_ JobClient         = (*job)(nil)
	_ CronJobClient     = (*cronJob)(nil)
)

func Pod(c Client, namespace string) PodClient {
	return &pod{
		Interface: c.GetClientSet().CoreV1().Pods(namespace),
		Client:    c,
		Namespace: namespace,
	}
}

func Deployment(c Client, namespace string) DeploymentClient {
	return &deployment{
		Interface: c.GetClientSet().AppsV1().Deployments(namespace),
		Client:    c,
		Namespace: namespace,
	}
}

func DaemonSet(c Client, namespace string) DaemonSetClient {
	return &daemonSet{
		Interface: c.GetClientSet().AppsV1().DaemonSets(namespace),
		Client:    c,
		Namespace: namespace,
	}
}

func StatefulSet(c Client, namespace string) StatefulSetClient {
	return &statefulSet{
		Interface: c.GetClientSet().AppsV1().StatefulSets(namespace),
		Client:    c,
		Namespace: namespace,
	}
}

func Job(c Client, namespace string) JobClient {
	return &job{
		Interface: c.GetClientSet().BatchV1().Jobs(namespace),
		Client:    c,
		Namespace: namespace,
	}
}

func CronJob(c Client, namespace string) CronJobClient {
	return &cronJob{
		Interface: c.GetClientSet().BatchV1().CronJobs(namespace),
		Client:    c,
		Namespace: namespace,
	}
}

type baseClient[T Object, L ObjectList, C NamedObject] struct {
	Interface ObjectInterface[T, L, C]
	Client    Client
	Namespace string
	mu        sync.RWMutex
}

func (b *baseClient[T, L, C]) Create(
	ctx context.Context, resource T, opts metav1.CreateOptions,
) (t T, err error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	if b == nil || b.Interface == nil {
		return t, errors.ErrKubernetesClientNotFound
	}
	return b.Interface.Create(ctx, resource, opts)
}

func (b *baseClient[T, L, C]) Update(
	ctx context.Context, resource T, opts metav1.UpdateOptions,
) (t T, err error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	if b == nil || b.Interface == nil {
		return t, errors.ErrKubernetesClientNotFound
	}
	return b.Interface.Update(ctx, resource, opts)
}

func (b *baseClient[T, L, C]) UpdateStatus(
	ctx context.Context, resource T, opts metav1.UpdateOptions,
) (t T, err error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	if b == nil || b.Interface == nil {
		return t, errors.ErrKubernetesClientNotFound
	}
	return b.Interface.UpdateStatus(ctx, resource, opts)
}

func (b *baseClient[T, L, C]) Delete(
	ctx context.Context, name string, opts metav1.DeleteOptions,
) error {
	b.mu.RLock()
	defer b.mu.RUnlock()
	if b == nil || b.Interface == nil {
		return errors.ErrKubernetesClientNotFound
	}
	return b.Interface.Delete(ctx, name, opts)
}

func (b *baseClient[T, L, C]) DeleteCollection(
	ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions,
) error {
	b.mu.RLock()
	defer b.mu.RUnlock()
	if b == nil || b.Interface == nil {
		return errors.ErrKubernetesClientNotFound
	}
	return b.Interface.DeleteCollection(ctx, opts, listOpts)
}

func (b *baseClient[T, L, C]) Get(
	ctx context.Context, name string, opts metav1.GetOptions,
) (t T, err error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	if b == nil || b.Interface == nil {
		return t, errors.ErrKubernetesClientNotFound
	}
	return b.Interface.Get(ctx, name, opts)
}

func (b *baseClient[T, L, C]) List(ctx context.Context, opts metav1.ListOptions) (l L, err error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	if b == nil || b.Interface == nil {
		return l, errors.ErrKubernetesClientNotFound
	}
	return b.Interface.List(ctx, opts)
}

func (b *baseClient[T, L, C]) Watch(
	ctx context.Context, opts metav1.ListOptions,
) (w watch.Interface, err error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	if b == nil || b.Interface == nil {
		return w, errors.ErrKubernetesClientNotFound
	}
	return b.Interface.Watch(ctx, opts)
}

func (b *baseClient[T, L, C]) Patch(
	ctx context.Context,
	name string,
	pt types.PatchType,
	data []byte,
	opts metav1.PatchOptions,
	subresources ...string,
) (t T, err error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	if b == nil || b.Interface == nil {
		return t, errors.ErrKubernetesClientNotFound
	}
	return b.Interface.Patch(ctx, name, pt, data, opts, subresources...)
}

func (b *baseClient[T, L, C]) Apply(
	ctx context.Context, resource C, opts metav1.ApplyOptions,
) (t T, err error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	if b == nil || b.Interface == nil {
		return t, errors.ErrKubernetesClientNotFound
	}
	return b.Interface.Apply(ctx, resource, opts)
}

func (b *baseClient[T, L, C]) UpdateEphemeralContainers(
	ctx context.Context, podName string, pod *corev1.Pod, opts metav1.UpdateOptions,
) (*corev1.Pod, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	if b == nil || b.Interface == nil {
		return nil, errors.ErrKubernetesClientNotFound
	}
	if pc, ok := b.Interface.(PodExtendInterface); ok {
		return pc.UpdateEphemeralContainers(ctx, podName, pod, opts)
	}
	return nil, errors.ErrUnimplemented("UpdateEphemeralContainers")
}

func (b *baseClient[T, L, C]) UpdateResize(
	ctx context.Context, podName string, pod *corev1.Pod, opts metav1.UpdateOptions,
) (*corev1.Pod, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	if b == nil || b.Interface == nil {
		return nil, errors.ErrKubernetesClientNotFound
	}
	if pc, ok := b.Interface.(PodExtendInterface); ok {
		return pc.UpdateResize(ctx, podName, pod, opts)
	}
	return nil, errors.ErrUnimplemented("UpdateResize")
}

func (b *baseClient[T, L, C]) GetScale(
	ctx context.Context, resourceName string, options metav1.GetOptions,
) (*autoscalingv1.Scale, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	if b == nil || b.Interface == nil {
		return nil, errors.ErrKubernetesClientNotFound
	}
	if sc, ok := b.Interface.(ScaleInterface); ok {
		return sc.GetScale(ctx, resourceName, options)
	}
	return nil, errors.ErrUnimplemented("GetScale")
}

func (b *baseClient[T, L, C]) UpdateScale(
	ctx context.Context, resourceName string, scale *autoscalingv1.Scale, opts metav1.UpdateOptions,
) (*autoscalingv1.Scale, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	if b == nil || b.Interface == nil {
		return nil, errors.ErrKubernetesClientNotFound
	}
	if sc, ok := b.Interface.(ScaleInterface); ok {
		return sc.UpdateScale(ctx, resourceName, scale, opts)
	}
	return nil, errors.ErrUnimplemented("UpdateScale")
}

func (b *baseClient[T, L, C]) ApplyScale(
	ctx context.Context,
	resourceName string,
	scale *applyconfigurationsautoscalingv1.ScaleApplyConfiguration,
	opts metav1.ApplyOptions,
) (*autoscalingv1.Scale, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	if b == nil || b.Interface == nil {
		return nil, errors.ErrKubernetesClientNotFound
	}
	if sc, ok := b.Interface.(ScaleInterface); ok {
		return sc.ApplyScale(ctx, resourceName, scale, opts)
	}
	return nil, errors.ErrUnimplemented("ApplyScale")
}

func (b *baseClient[T, L, C]) GetInterface() ObjectInterface[T, L, C] {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.Interface
}

func (b *baseClient[T, L, C]) SetInterface(c ObjectInterface[T, L, C]) {
	b.mu.Lock()
	b.Interface = c
	b.mu.Unlock()
}

func (b *baseClient[T, L, C]) GetClient() Client {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.Client
}

func (b *baseClient[T, L, C]) SetClient(c Client) {
	b.mu.Lock()
	b.Client = c
	b.mu.Unlock()
}

func (b *baseClient[T, L, C]) GetNamespace() string {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.Namespace
}

func (b *baseClient[T, L, C]) SetNamespace(namespace string) {
	b.mu.Lock()
	b.Namespace = namespace
	b.mu.Unlock()
}
