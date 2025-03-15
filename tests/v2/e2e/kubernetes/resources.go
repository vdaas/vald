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

type ResourceInterface[T Object, L ObjectList, C NamedObject] interface {
	Get(ctx context.Context, name string, opts metav1.GetOptions) (T, error)
	List(ctx context.Context, opts metav1.ListOptions) (L, error)

	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)

	Create(ctx context.Context, resource T, opts metav1.CreateOptions) (T, error)

	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error

	Update(ctx context.Context, resource T, opts metav1.UpdateOptions) (T, error)

	Apply(ctx context.Context, resource C, opts metav1.ApplyOptions) (result T, err error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result T, err error)
}

type ExtResourceInterface[T Object] interface {
	DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error
	UpdateStatus(ctx context.Context, resource T, opts metav1.UpdateOptions) (T, error)
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

type PodTemplateInterface[T Object] interface {
	GetPodTemplate(obj T) (*corev1.PodTemplateSpec, error)
	SetPodTemplate(obj T, pt *corev1.PodTemplateSpec) (T, error)
}

type PodAnnotationInterface[T Object] interface {
	GetPodAnnotations(ctx context.Context, name string, opts metav1.GetOptions) (map[string]string, error)
	SetPodAnnotations(ctx context.Context, name string, annotations map[string]string, gopts metav1.GetOptions, uopts metav1.UpdateOptions) (T, error)
}

type ClientControlInterface[T Object, L ObjectList, C NamedObject] interface {
	GetInterface() ResourceInterface[T, L, C]
	SetInterface(c ResourceInterface[T, L, C])

	GetClient() Client
	SetClient(c Client)

	GetNamespace() string
	SetNamespace(namespace string)
}

type ResourceClient[T Object, L ObjectList, C NamedObject] interface {
	ResourceInterface[T, L, C]
	ClientControlInterface[T, L, C]
}

type WorkloadResourceClient[T Object, L ObjectList, C NamedObject] interface {
	ResourceClient[T, L, C]
	ExtResourceInterface[T]
}

type WorkloadControllerResourceClient[T Object, L ObjectList, C NamedObject] interface {
	WorkloadResourceClient[T, L, C]
	PodTemplateInterface[T]
	PodAnnotationInterface[T]
}

type (
	PodClient interface {
		WorkloadResourceClient[*corev1.Pod, *corev1.PodList, *applyconfigurationscorev1.PodApplyConfiguration]
		PodExtendInterface
	}
	DeploymentClient interface {
		WorkloadControllerResourceClient[*appsv1.Deployment, *appsv1.DeploymentList, *applyconfigurationsappsv1.DeploymentApplyConfiguration]
		ScaleInterface
	}
	StatefulSetClient interface {
		WorkloadControllerResourceClient[*appsv1.StatefulSet, *appsv1.StatefulSetList, *applyconfigurationsappsv1.StatefulSetApplyConfiguration]
		ScaleInterface
	}
	DaemonSetClient             = WorkloadControllerResourceClient[*appsv1.DaemonSet, *appsv1.DaemonSetList, *applyconfigurationsappsv1.DaemonSetApplyConfiguration]
	JobClient                   = WorkloadControllerResourceClient[*batchv1.Job, *batchv1.JobList, *applyconfigurationsbatchv1.JobApplyConfiguration]
	CronJobClient               = WorkloadControllerResourceClient[*batchv1.CronJob, *batchv1.CronJobList, *applyconfigurationsbatchv1.CronJobApplyConfiguration]
	ServiceClient               = ResourceClient[*corev1.Service, *corev1.ServiceList, *applyconfigurationscorev1.ServiceApplyConfiguration]
	SecretClient                = ResourceClient[*corev1.Secret, *corev1.SecretList, *applyconfigurationscorev1.SecretApplyConfiguration]
	ConfigMapClient             = ResourceClient[*corev1.ConfigMap, *corev1.ConfigMapList, *applyconfigurationscorev1.ConfigMapApplyConfiguration]
	PersistentVolumeClaimClient = ResourceClient[*corev1.PersistentVolumeClaim, *corev1.PersistentVolumeClaimList, *applyconfigurationscorev1.PersistentVolumeClaimApplyConfiguration]
	PersistentVolumeClient      = ResourceClient[*corev1.PersistentVolume, *corev1.PersistentVolumeList, *applyconfigurationscorev1.PersistentVolumeApplyConfiguration]
	EndpointClient              = ResourceClient[*corev1.Endpoints, *corev1.EndpointsList, *applyconfigurationscorev1.EndpointsApplyConfiguration]
)

type (
	pod         = baseClient[*corev1.Pod, *corev1.PodList, *applyconfigurationscorev1.PodApplyConfiguration]
	deployment  = baseClient[*appsv1.Deployment, *appsv1.DeploymentList, *applyconfigurationsappsv1.DeploymentApplyConfiguration]
	daemonSet   = baseClient[*appsv1.DaemonSet, *appsv1.DaemonSetList, *applyconfigurationsappsv1.DaemonSetApplyConfiguration]
	statefulSet = baseClient[*appsv1.StatefulSet, *appsv1.StatefulSetList, *applyconfigurationsappsv1.StatefulSetApplyConfiguration]
	job         = baseClient[*batchv1.Job, *batchv1.JobList, *applyconfigurationsbatchv1.JobApplyConfiguration]
	cronJob     = baseClient[*batchv1.CronJob, *batchv1.CronJobList, *applyconfigurationsbatchv1.CronJobApplyConfiguration]
	service     = baseClient[*corev1.Service, *corev1.ServiceList, *applyconfigurationscorev1.ServiceApplyConfiguration]
	secret      = baseClient[*corev1.Secret, *corev1.SecretList, *applyconfigurationscorev1.SecretApplyConfiguration]
	configMap   = baseClient[*corev1.ConfigMap, *corev1.ConfigMapList, *applyconfigurationscorev1.ConfigMapApplyConfiguration]
	pvc         = baseClient[*corev1.PersistentVolumeClaim, *corev1.PersistentVolumeClaimList, *applyconfigurationscorev1.PersistentVolumeClaimApplyConfiguration]
	pv          = baseClient[*corev1.PersistentVolume, *corev1.PersistentVolumeList, *applyconfigurationscorev1.PersistentVolumeApplyConfiguration]
	endponts    = baseClient[*corev1.Endpoints, *corev1.EndpointsList, *applyconfigurationscorev1.EndpointsApplyConfiguration]
)

var (
	_ PodClient                   = (*pod)(nil)
	_ DeploymentClient            = (*deployment)(nil)
	_ DaemonSetClient             = (*daemonSet)(nil)
	_ StatefulSetClient           = (*statefulSet)(nil)
	_ JobClient                   = (*job)(nil)
	_ CronJobClient               = (*cronJob)(nil)
	_ ServiceClient               = (*service)(nil)
	_ SecretClient                = (*secret)(nil)
	_ ConfigMapClient             = (*configMap)(nil)
	_ PersistentVolumeClaimClient = (*pvc)(nil)
	_ PersistentVolumeClient      = (*pv)(nil)
	_ EndpointClient              = (*endponts)(nil)
)

func Pod(c Client, namespace string) PodClient {
	if c == nil {
		return nil
	}
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
		getPodTemplate: func(t *appsv1.Deployment) *corev1.PodTemplateSpec {
			return &t.Spec.Template
		},
		setPodTemplate: func(t *appsv1.Deployment, pt *corev1.PodTemplateSpec) *appsv1.Deployment {
			t.Spec.Template = *pt
			return t
		},
		Namespace: namespace,
	}
}

func DaemonSet(c Client, namespace string) DaemonSetClient {
	return &daemonSet{
		Interface: c.GetClientSet().AppsV1().DaemonSets(namespace),
		Client:    c,
		getPodTemplate: func(t *appsv1.DaemonSet) *corev1.PodTemplateSpec {
			return &t.Spec.Template
		},
		setPodTemplate: func(t *appsv1.DaemonSet, pt *corev1.PodTemplateSpec) *appsv1.DaemonSet {
			t.Spec.Template = *pt
			return t
		},
		Namespace: namespace,
	}
}

func StatefulSet(c Client, namespace string) StatefulSetClient {
	return &statefulSet{
		Interface: c.GetClientSet().AppsV1().StatefulSets(namespace),
		Client:    c,
		getPodTemplate: func(t *appsv1.StatefulSet) *corev1.PodTemplateSpec {
			return &t.Spec.Template
		},
		setPodTemplate: func(t *appsv1.StatefulSet, pt *corev1.PodTemplateSpec) *appsv1.StatefulSet {
			t.Spec.Template = *pt
			return t
		},
		Namespace: namespace,
	}
}

func Job(c Client, namespace string) JobClient {
	return &job{
		Interface: c.GetClientSet().BatchV1().Jobs(namespace),
		Client:    c,
		getPodTemplate: func(t *batchv1.Job) *corev1.PodTemplateSpec {
			return &t.Spec.Template
		},
		setPodTemplate: func(t *batchv1.Job, pt *corev1.PodTemplateSpec) *batchv1.Job {
			t.Spec.Template = *pt
			return t
		},
		Namespace: namespace,
	}
}

func CronJob(c Client, namespace string) CronJobClient {
	return &cronJob{
		Interface: c.GetClientSet().BatchV1().CronJobs(namespace),
		Client:    c,
		getPodTemplate: func(t *batchv1.CronJob) *corev1.PodTemplateSpec {
			return &t.Spec.JobTemplate.Spec.Template
		},
		setPodTemplate: func(t *batchv1.CronJob, pt *corev1.PodTemplateSpec) *batchv1.CronJob {
			t.Spec.JobTemplate.Spec.Template = *pt
			return t
		},
		Namespace: namespace,
	}
}

func Service(c Client, namespace string) ServiceClient {
	return &service{
		Interface: c.GetClientSet().CoreV1().Services(namespace),
		Client:    c,
		Namespace: namespace,
	}
}

func Secret(c Client, namespace string) SecretClient {
	return &secret{
		Interface: c.GetClientSet().CoreV1().Secrets(namespace),
		Client:    c,
		Namespace: namespace,
	}
}

func ConfigMap(c Client, namespace string) ConfigMapClient {
	return &configMap{
		Interface: c.GetClientSet().CoreV1().ConfigMaps(namespace),
		Client:    c,
		Namespace: namespace,
	}
}

func PersistentVolumeClaim(c Client, namespace string) PersistentVolumeClaimClient {
	return &pvc{
		Interface: c.GetClientSet().CoreV1().PersistentVolumeClaims(namespace),
		Client:    c,
		Namespace: namespace,
	}
}

func PersistentVolume(c Client) PersistentVolumeClient {
	return &pv{
		Interface: c.GetClientSet().CoreV1().PersistentVolumes(),
		Client:    c,
	}
}

func Endpoints(c Client, namespace string) EndpointClient {
	return &endponts{
		Interface: c.GetClientSet().CoreV1().Endpoints(namespace),
		Client:    c,
		Namespace: namespace,
	}
}

type baseClient[T Object, L ObjectList, C NamedObject] struct {
	Interface      ResourceInterface[T, L, C]
	Client         Client
	getPodTemplate func(t T) *corev1.PodTemplateSpec
	setPodTemplate func(t T, pt *corev1.PodTemplateSpec) T
	Namespace      string
	mu             sync.RWMutex
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
	if eri, ok := b.Interface.(ExtResourceInterface[T]); ok {
		return eri.UpdateStatus(ctx, resource, opts)
	}
	return t, errors.ErrUnimplemented("UpdateStatus")
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
	if eri, ok := b.Interface.(ExtResourceInterface[T]); ok {
		return eri.DeleteCollection(ctx, opts, listOpts)
	}
	return errors.ErrUnimplemented("DeleteCollection")
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

func (b *baseClient[T, L, C]) GetInterface() ResourceInterface[T, L, C] {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.Interface
}

func (b *baseClient[T, L, C]) SetInterface(c ResourceInterface[T, L, C]) {
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

func (b *baseClient[T, L, C]) GetPodTemplate(obj T) (*corev1.PodTemplateSpec, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	if b.getPodTemplate == nil {
		return nil, errors.ErrPodTemplateNotFound
	}
	return b.getPodTemplate(obj), nil
}

func (b *baseClient[T, L, C]) SetPodTemplate(obj T, pt *corev1.PodTemplateSpec) (T, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	if b.setPodTemplate == nil {
		return obj, errors.ErrPodTemplateNotFound
	}
	return b.setPodTemplate(obj, pt), nil
}

func (b *baseClient[T, L, C]) GetPodAnnotations(
	ctx context.Context, name string, opts metav1.GetOptions,
) (map[string]string, error) {
	obj, err := b.Get(ctx, name, opts)
	if err != nil {
		return nil, err
	}
	tmpl, err := b.GetPodTemplate(obj)
	if err != nil {
		return nil, err
	}
	if tmpl == nil || tmpl.Annotations == nil {
		return nil, errors.ErrPodTemplateNotFound
	}
	return tmpl.Annotations, nil
}

func (b *baseClient[T, L, C]) SetPodAnnotations(
	ctx context.Context,
	name string,
	annotations map[string]string,
	gopts metav1.GetOptions,
	uopts metav1.UpdateOptions,
) (T, error) {
	obj, err := b.Get(ctx, name, gopts)
	if err != nil {
		return obj, err
	}
	tmpl, err := b.GetPodTemplate(obj)
	if err != nil {
		return obj, err
	}
	if tmpl == nil {
		return obj, errors.ErrPodTemplateNotFound
	}
	if tmpl.Annotations == nil {
		tmpl.Annotations = make(map[string]string, len(annotations))
	}
	for key, val := range annotations {
		tmpl.Annotations[key] = val
	}
	obj, err = b.SetPodTemplate(obj, tmpl)
	if err != nil {
		return obj, err
	}
	return b.Update(ctx, obj, uopts)
}
