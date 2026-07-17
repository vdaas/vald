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

package resource

import (
	"context"
	"fmt"
	"maps"
	"math/rand/v2"
	"strconv"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/k8s/client"
	"github.com/vdaas/vald/internal/sync"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	appsv1 "k8s.io/api/apps/v1"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	applyconfigurationsadmissionregistrationv1 "k8s.io/client-go/applyconfigurations/admissionregistration/v1"
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

type extResourceInterface[T Object] interface {
	DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error
	UpdateStatus(ctx context.Context, resource T, opts metav1.UpdateOptions) (T, error)
}

type scaleInterface interface {
	GetScale(ctx context.Context, resourceName string, options metav1.GetOptions) (*autoscalingv1.Scale, error)
	UpdateScale(ctx context.Context, resourceName string, scale *autoscalingv1.Scale, opts metav1.UpdateOptions) (*autoscalingv1.Scale, error)
	ApplyScale(ctx context.Context, resourceName string, scale *applyconfigurationsautoscalingv1.ScaleApplyConfiguration, opts metav1.ApplyOptions) (*autoscalingv1.Scale, error)
}

type podExtendInterface interface {
	UpdateEphemeralContainers(ctx context.Context, podName string, pod *corev1.Pod, opts metav1.UpdateOptions) (*corev1.Pod, error)
	UpdateResize(ctx context.Context, podName string, pod *corev1.Pod, opts metav1.UpdateOptions) (*corev1.Pod, error)
}

type podTemplateInterface[T Object] interface {
	GetPodTemplate(obj T) (*corev1.PodTemplateSpec, error)
	SetPodTemplate(obj T, pt *corev1.PodTemplateSpec) (T, error)
}

type podAnnotationInterface[T Object] interface {
	GetPodAnnotations(ctx context.Context, name string, opts metav1.GetOptions) (map[string]string, error)
	SetPodAnnotations(ctx context.Context, name string, annotations map[string]string, gopts metav1.GetOptions, uopts metav1.UpdateOptions) (T, error)
}

type clientControlInterface[T Object, L ObjectList, C NamedObject] interface {
	GetInterface() ResourceInterface[T, L, C]
	SetInterface(c ResourceInterface[T, L, C])

	GetClient() client.Client
	SetClient(c client.Client)

	GetNamespace() string
	SetNamespace(namespace string)
}

type ResourceClient[T Object, L ObjectList, C NamedObject] interface {
	ResourceInterface[T, L, C]
	clientControlInterface[T, L, C]
}

type WorkloadResourceClient[T Object, L ObjectList, C NamedObject] interface {
	ResourceClient[T, L, C]
	extResourceInterface[T]
}

type WorkloadControllerResourceClient[T Object, L ObjectList, C NamedObject] interface {
	WorkloadResourceClient[T, L, C]
	podTemplateInterface[T]
	podAnnotationInterface[T]
	CreateJob(
		ctx context.Context,
		from string,
		gopts metav1.GetOptions,
		copts metav1.CreateOptions,
	) (*batchv1.Job, error)
}

type (
	PodClient interface {
		WorkloadResourceClient[*corev1.Pod, *corev1.PodList, *applyconfigurationscorev1.PodApplyConfiguration]
		podExtendInterface
	}
	DeploymentClient interface {
		WorkloadControllerResourceClient[*appsv1.Deployment, *appsv1.DeploymentList, *applyconfigurationsappsv1.DeploymentApplyConfiguration]
		scaleInterface
	}
	StatefulSetClient interface {
		WorkloadControllerResourceClient[*appsv1.StatefulSet, *appsv1.StatefulSetList, *applyconfigurationsappsv1.StatefulSetApplyConfiguration]
		scaleInterface
	}
	DaemonSetClient = WorkloadControllerResourceClient[*appsv1.DaemonSet, *appsv1.DaemonSetList, *applyconfigurationsappsv1.DaemonSetApplyConfiguration]
	JobClient       = WorkloadControllerResourceClient[*batchv1.Job, *batchv1.JobList, *applyconfigurationsbatchv1.JobApplyConfiguration]
	CronJobClient   = WorkloadControllerResourceClient[*batchv1.CronJob, *batchv1.CronJobList, *applyconfigurationsbatchv1.CronJobApplyConfiguration]
	ServiceClient   = ResourceClient[*corev1.Service, *corev1.ServiceList, *applyconfigurationscorev1.ServiceApplyConfiguration]
	SecretClient    = ResourceClient[*corev1.Secret, *corev1.SecretList, *applyconfigurationscorev1.SecretApplyConfiguration]
	ConfigMapClient = ResourceClient[*corev1.ConfigMap, *corev1.ConfigMapList, *applyconfigurationscorev1.ConfigMapApplyConfiguration]
	EndpointClient  = ResourceClient[*corev1.Endpoints, *corev1.EndpointsList, *applyconfigurationscorev1.EndpointsApplyConfiguration] //nolint:staticcheck // skipcq: GO-W1009 Endpoints is still served; EndpointSlice migration is tracked separately.

	MutatingWebhookConfigurationClient   = ResourceClient[*admissionregistrationv1.MutatingWebhookConfiguration, *admissionregistrationv1.MutatingWebhookConfigurationList, *applyconfigurationsadmissionregistrationv1.MutatingWebhookConfigurationApplyConfiguration]
	ValidatingWebhookConfigurationClient = ResourceClient[*admissionregistrationv1.ValidatingWebhookConfiguration, *admissionregistrationv1.ValidatingWebhookConfigurationList, *applyconfigurationsadmissionregistrationv1.ValidatingWebhookConfigurationApplyConfiguration]
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
	endpoints   = baseClient[*corev1.Endpoints, *corev1.EndpointsList, *applyconfigurationscorev1.EndpointsApplyConfiguration] //nolint:staticcheck // skipcq: GO-W1009

	mutatingWebhookConfiguration   = baseClient[*admissionregistrationv1.MutatingWebhookConfiguration, *admissionregistrationv1.MutatingWebhookConfigurationList, *applyconfigurationsadmissionregistrationv1.MutatingWebhookConfigurationApplyConfiguration]
	validatingWebhookConfiguration = baseClient[*admissionregistrationv1.ValidatingWebhookConfiguration, *admissionregistrationv1.ValidatingWebhookConfigurationList, *applyconfigurationsadmissionregistrationv1.ValidatingWebhookConfigurationApplyConfiguration]
)

var (
	_ PodClient         = (*pod)(nil)
	_ DeploymentClient  = (*deployment)(nil)
	_ DaemonSetClient   = (*daemonSet)(nil)
	_ StatefulSetClient = (*statefulSet)(nil)
	_ JobClient         = (*job)(nil)
	_ CronJobClient     = (*cronJob)(nil)
	_ ServiceClient     = (*service)(nil)
	_ SecretClient      = (*secret)(nil)
	_ ConfigMapClient   = (*configMap)(nil)
	_ EndpointClient    = (*endpoints)(nil)

	_ MutatingWebhookConfigurationClient   = (*mutatingWebhookConfiguration)(nil)
	_ ValidatingWebhookConfigurationClient = (*validatingWebhookConfiguration)(nil)
)

// Empty{Get,Create,Delete}Options are zero-value option sentinels referenced
// by value at call sites frozen as part of the tests/v2/e2e public API
// (e.g. tests/v2/e2e/crud/kubernetes_test.go), so they cannot be converted to
// functions without breaking that frozen surface.
//
//nolint:gochecknoglobals
var (
	EmptyGetOptions    = metav1.GetOptions{}
	EmptyCreateOptions = metav1.CreateOptions{}
	EmptyDeleteOptions = metav1.DeleteOptions{}
)

// withIface executes f under the read lock.
// Returns ErrKubernetesClientNotFound if b or b.Interface is nil.
func withIface[R any, T Object, L ObjectList, C NamedObject](
	b *baseClient[T, L, C], f func(ResourceInterface[T, L, C]) (R, error),
) (R, error) {
	var zero R
	if b == nil {
		return zero, errors.ErrKubernetesClientNotFound
	}
	b.mu.RLock()
	defer b.mu.RUnlock()
	if b.Interface == nil {
		return zero, errors.ErrKubernetesClientNotFound
	}
	return f(b.Interface)
}

// withExt executes f under the read lock after asserting the typed clientset
// interface to the optional extension interface I (extResourceInterface,
// scaleInterface, podExtendInterface, ...), returning ErrUnimplemented(name)
// when the resource does not support it. It centralizes the ok-check +
// zero-value + name-string boilerplate the extension dispatch methods would
// otherwise each repeat.
func withExt[I any, R any, T Object, L ObjectList, C NamedObject](
	b *baseClient[T, L, C], name string, f func(I) (R, error),
) (R, error) {
	return withIface(b, func(i ResourceInterface[T, L, C]) (R, error) {
		ei, ok := i.(I)
		if !ok {
			var zero R
			return zero, errors.ErrUnimplemented(name)
		}
		return f(ei)
	})
}

// newResourceClient builds the generic client implementation shared by every
// factory. The returned concrete type satisfies all per-resource client
// interfaces, so each factory only wires the typed clientset interface.
func newResourceClient[T Object, L ObjectList, C NamedObject](
	c client.Client, iface ResourceInterface[T, L, C], namespace string,
) *baseClient[T, L, C] {
	return &baseClient[T, L, C]{
		Interface: iface,
		Client:    c,
		Namespace: namespace,
	}
}

// newWorkloadClient builds the generic client implementation for workload
// controllers, which additionally need access to their pod template.
// podTemplate returns a pointer INTO the object, so SetPodTemplate can derive
// the setter by assigning through it; Go generics cannot express the field
// access itself (golang/go#48522), so one accessor per Kind is the floor.
func newWorkloadClient[T Object, L ObjectList, C NamedObject](
	c client.Client,
	iface ResourceInterface[T, L, C],
	namespace string,
	podTemplate func(t T) *corev1.PodTemplateSpec,
) *baseClient[T, L, C] {
	b := newResourceClient(c, iface, namespace)
	b.podTemplate = podTemplate
	return b
}

// factory guards every per-Kind constructor below against a nil ClientSet,
// returning R's zero value (nil, since every R here is a client interface)
// instead of invoking build. c is passed in only to be checked: build closes
// over the same value to reach the clientset, so callers never repeat the
// "if c == nil { return nil }" guard.
func factory[R any](c client.Client, build func() R) R {
	var zero R
	if c == nil {
		return zero
	}
	return build()
}

func Pod(c client.Client, namespace string) PodClient {
	return factory(c, func() PodClient {
		return newResourceClient(c, c.GetClientSet().CoreV1().Pods(namespace), namespace)
	})
}

func Deployment(c client.Client, namespace string) DeploymentClient {
	return factory(c, func() DeploymentClient {
		return newWorkloadClient(c, c.GetClientSet().AppsV1().Deployments(namespace), namespace,
			func(t *appsv1.Deployment) *corev1.PodTemplateSpec {
				return &t.Spec.Template
			})
	})
}

func DaemonSet(c client.Client, namespace string) DaemonSetClient {
	return factory(c, func() DaemonSetClient {
		return newWorkloadClient(c, c.GetClientSet().AppsV1().DaemonSets(namespace), namespace,
			func(t *appsv1.DaemonSet) *corev1.PodTemplateSpec {
				return &t.Spec.Template
			})
	})
}

func StatefulSet(c client.Client, namespace string) StatefulSetClient {
	return factory(c, func() StatefulSetClient {
		return newWorkloadClient(c, c.GetClientSet().AppsV1().StatefulSets(namespace), namespace,
			func(t *appsv1.StatefulSet) *corev1.PodTemplateSpec {
				return &t.Spec.Template
			})
	})
}

func Job(c client.Client, namespace string) JobClient {
	return factory(c, func() JobClient {
		return newWorkloadClient(c, c.GetClientSet().BatchV1().Jobs(namespace), namespace,
			func(t *batchv1.Job) *corev1.PodTemplateSpec {
				return &t.Spec.Template
			})
	})
}

func CronJob(c client.Client, namespace string) CronJobClient {
	return factory(c, func() CronJobClient {
		return newWorkloadClient(c, c.GetClientSet().BatchV1().CronJobs(namespace), namespace,
			func(t *batchv1.CronJob) *corev1.PodTemplateSpec {
				return &t.Spec.JobTemplate.Spec.Template
			})
	})
}

func Service(c client.Client, namespace string) ServiceClient {
	return factory(c, func() ServiceClient {
		return newResourceClient(c, c.GetClientSet().CoreV1().Services(namespace), namespace)
	})
}

func Secret(c client.Client, namespace string) SecretClient {
	return factory(c, func() SecretClient {
		return newResourceClient(c, c.GetClientSet().CoreV1().Secrets(namespace), namespace)
	})
}

func ConfigMap(c client.Client, namespace string) ConfigMapClient {
	return factory(c, func() ConfigMapClient {
		return newResourceClient(c, c.GetClientSet().CoreV1().ConfigMaps(namespace), namespace)
	})
}

func MutatingWebhookConfiguration(c client.Client) MutatingWebhookConfigurationClient {
	return factory(c, func() MutatingWebhookConfigurationClient {
		return newResourceClient(c, c.GetClientSet().AdmissionregistrationV1().MutatingWebhookConfigurations(), "")
	})
}

func ValidatingWebhookConfiguration(c client.Client) ValidatingWebhookConfigurationClient {
	return factory(c, func() ValidatingWebhookConfigurationClient {
		return newResourceClient(c, c.GetClientSet().AdmissionregistrationV1().ValidatingWebhookConfigurations(), "")
	})
}

func Endpoints(c client.Client, namespace string) EndpointClient {
	return factory(c, func() EndpointClient {
		return newResourceClient(c, c.GetClientSet().CoreV1().Endpoints(namespace), namespace)
	})
}

type baseClient[T Object, L ObjectList, C NamedObject] struct {
	Interface ResourceInterface[T, L, C]
	Client    client.Client
	// podTemplate returns a pointer to the object's pod template field; both
	// GetPodTemplate and SetPodTemplate derive from it, so getter and setter
	// cannot disagree about which field they access.
	podTemplate func(t T) *corev1.PodTemplateSpec
	Namespace   string
	mu          sync.RWMutex
}

func (b *baseClient[T, L, C]) Create(
	ctx context.Context, resource T, opts metav1.CreateOptions,
) (T, error) {
	return withIface(b, func(i ResourceInterface[T, L, C]) (T, error) {
		return i.Create(ctx, resource, opts)
	})
}

func (b *baseClient[T, L, C]) Update(
	ctx context.Context, resource T, opts metav1.UpdateOptions,
) (T, error) {
	return withIface(b, func(i ResourceInterface[T, L, C]) (T, error) {
		return i.Update(ctx, resource, opts)
	})
}

func (b *baseClient[T, L, C]) UpdateStatus(
	ctx context.Context, resource T, opts metav1.UpdateOptions,
) (T, error) {
	return withExt(b, "UpdateStatus", func(eri extResourceInterface[T]) (T, error) {
		return eri.UpdateStatus(ctx, resource, opts)
	})
}

func (b *baseClient[T, L, C]) Delete(
	ctx context.Context, name string, opts metav1.DeleteOptions,
) error {
	_, err := withIface(b, func(i ResourceInterface[T, L, C]) (struct{}, error) {
		return struct{}{}, i.Delete(ctx, name, opts)
	})
	return err
}

func (b *baseClient[T, L, C]) DeleteCollection(
	ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions,
) error {
	_, err := withExt(b, "DeleteCollection", func(eri extResourceInterface[T]) (struct{}, error) {
		return struct{}{}, eri.DeleteCollection(ctx, opts, listOpts)
	})
	return err
}

func (b *baseClient[T, L, C]) Get(
	ctx context.Context, name string, opts metav1.GetOptions,
) (T, error) {
	return withIface(b, func(i ResourceInterface[T, L, C]) (T, error) {
		return i.Get(ctx, name, opts)
	})
}

func (b *baseClient[T, L, C]) List(ctx context.Context, opts metav1.ListOptions) (L, error) {
	return withIface(b, func(i ResourceInterface[T, L, C]) (L, error) {
		return i.List(ctx, opts)
	})
}

func (b *baseClient[T, L, C]) Watch(
	ctx context.Context, opts metav1.ListOptions,
) (watch.Interface, error) {
	return withIface(b, func(i ResourceInterface[T, L, C]) (watch.Interface, error) {
		return i.Watch(ctx, opts)
	})
}

func (b *baseClient[T, L, C]) Patch(
	ctx context.Context,
	name string,
	pt types.PatchType,
	data []byte,
	opts metav1.PatchOptions,
	subresources ...string,
) (T, error) {
	return withIface(b, func(i ResourceInterface[T, L, C]) (T, error) {
		return i.Patch(ctx, name, pt, data, opts, subresources...)
	})
}

func (b *baseClient[T, L, C]) Apply(
	ctx context.Context, resource C, opts metav1.ApplyOptions,
) (T, error) {
	return withIface(b, func(i ResourceInterface[T, L, C]) (T, error) {
		return i.Apply(ctx, resource, opts)
	})
}

func (b *baseClient[T, L, C]) UpdateEphemeralContainers(
	ctx context.Context, podName string, pod *corev1.Pod, opts metav1.UpdateOptions,
) (*corev1.Pod, error) {
	return withExt(b, "UpdateEphemeralContainers", func(pc podExtendInterface) (*corev1.Pod, error) {
		return pc.UpdateEphemeralContainers(ctx, podName, pod, opts)
	})
}

func (b *baseClient[T, L, C]) UpdateResize(
	ctx context.Context, podName string, pod *corev1.Pod, opts metav1.UpdateOptions,
) (*corev1.Pod, error) {
	return withExt(b, "UpdateResize", func(pc podExtendInterface) (*corev1.Pod, error) {
		return pc.UpdateResize(ctx, podName, pod, opts)
	})
}

func (b *baseClient[T, L, C]) GetScale(
	ctx context.Context, resourceName string, options metav1.GetOptions,
) (*autoscalingv1.Scale, error) {
	return withExt(b, "GetScale", func(sc scaleInterface) (*autoscalingv1.Scale, error) {
		return sc.GetScale(ctx, resourceName, options)
	})
}

func (b *baseClient[T, L, C]) UpdateScale(
	ctx context.Context, resourceName string, scale *autoscalingv1.Scale, opts metav1.UpdateOptions,
) (*autoscalingv1.Scale, error) {
	return withExt(b, "UpdateScale", func(sc scaleInterface) (*autoscalingv1.Scale, error) {
		return sc.UpdateScale(ctx, resourceName, scale, opts)
	})
}

func (b *baseClient[T, L, C]) ApplyScale(
	ctx context.Context,
	resourceName string,
	scale *applyconfigurationsautoscalingv1.ScaleApplyConfiguration,
	opts metav1.ApplyOptions,
) (*autoscalingv1.Scale, error) {
	return withExt(b, "ApplyScale", func(sc scaleInterface) (*autoscalingv1.Scale, error) {
		return sc.ApplyScale(ctx, resourceName, scale, opts)
	})
}

// getLocked returns *field while holding mu for reading. Shared by every
// baseClient getter that just copies a field out under the read lock.
func getLocked[V any](mu *sync.RWMutex, field *V) V {
	mu.RLock()
	defer mu.RUnlock()
	return *field
}

// setLocked stores v into *field while holding mu for writing. Shared by
// every baseClient setter that just copies a value into a field under the
// write lock.
func setLocked[V any](mu *sync.RWMutex, field *V, v V) {
	mu.Lock()
	*field = v
	mu.Unlock()
}

// withRLocked runs f while holding mu for reading and returns its result.
// Unlike getLocked/setLocked, f is a closure rather than a field pointer, so
// this also fits getters whose body is more than a bare field copy (e.g.
// GetPodTemplate's nil-guard).
func withRLocked[R any](mu *sync.RWMutex, f func() (R, error)) (R, error) {
	mu.RLock()
	defer mu.RUnlock()
	return f()
}

// withLocked is withRLocked's write-lock counterpart, used by SetPodTemplate.
func withLocked[R any](mu *sync.RWMutex, f func() (R, error)) (R, error) {
	mu.Lock()
	defer mu.Unlock()
	return f()
}

func (b *baseClient[T, L, C]) GetInterface() ResourceInterface[T, L, C] {
	return getLocked(&b.mu, &b.Interface)
}

func (b *baseClient[T, L, C]) SetInterface(c ResourceInterface[T, L, C]) {
	setLocked(&b.mu, &b.Interface, c)
}

func (b *baseClient[T, L, C]) GetClient() client.Client {
	return getLocked(&b.mu, &b.Client)
}

func (b *baseClient[T, L, C]) SetClient(c client.Client) {
	setLocked(&b.mu, &b.Client, c)
}

func (b *baseClient[T, L, C]) GetNamespace() string {
	return getLocked(&b.mu, &b.Namespace)
}

func (b *baseClient[T, L, C]) SetNamespace(namespace string) {
	setLocked(&b.mu, &b.Namespace, namespace)
}

// podTemplateAccessor returns b.podTemplate, or ErrPodTemplateNotFound if no
// accessor was wired for this Kind. Callers must already hold b.mu (read or
// write); it centralizes the nil-guard GetPodTemplate and SetPodTemplate
// would otherwise each repeat.
func (b *baseClient[T, L, C]) podTemplateAccessor() (func(T) *corev1.PodTemplateSpec, error) {
	if b.podTemplate == nil {
		return nil, errors.ErrPodTemplateNotFound
	}
	return b.podTemplate, nil
}

func (b *baseClient[T, L, C]) GetPodTemplate(obj T) (*corev1.PodTemplateSpec, error) {
	return withRLocked(&b.mu, func() (*corev1.PodTemplateSpec, error) {
		pt, err := b.podTemplateAccessor()
		if err != nil {
			return nil, err
		}
		return pt(obj), nil
	})
}

func (b *baseClient[T, L, C]) SetPodTemplate(obj T, pt *corev1.PodTemplateSpec) (T, error) {
	return withLocked(&b.mu, func() (T, error) {
		fn, err := b.podTemplateAccessor()
		if err != nil {
			return obj, err
		}
		*fn(obj) = *pt
		return obj, nil
	})
}

// getWithPodTemplate fetches obj by name and resolves its pod template,
// centralizing the fetch-then-template-nil-guard prefix shared by
// GetPodAnnotations, SetPodAnnotations, and CreateJob. Callers needing a
// different response to a missing (but present-and-empty) template, such as
// SetPodAnnotations initializing rather than erroring on nil Annotations,
// still handle that themselves on the returned template.
func (b *baseClient[T, L, C]) getWithPodTemplate(
	ctx context.Context, name string, opts metav1.GetOptions,
) (T, *corev1.PodTemplateSpec, error) {
	obj, err := b.Get(ctx, name, opts)
	if err != nil {
		return obj, nil, err
	}
	tmpl, err := b.GetPodTemplate(obj)
	if err != nil {
		return obj, nil, err
	}
	if tmpl == nil {
		return obj, nil, errors.ErrPodTemplateNotFound
	}
	return obj, tmpl, nil
}

func (b *baseClient[T, L, C]) GetPodAnnotations(
	ctx context.Context, name string, opts metav1.GetOptions,
) (map[string]string, error) {
	_, tmpl, err := b.getWithPodTemplate(ctx, name, opts)
	if err != nil {
		return nil, err
	}
	if tmpl.Annotations == nil {
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
	obj, tmpl, err := b.getWithPodTemplate(ctx, name, gopts)
	if err != nil {
		return obj, err
	}
	if tmpl.Annotations == nil {
		tmpl.Annotations = make(map[string]string, len(annotations))
	}
	maps.Copy(tmpl.Annotations, annotations)
	obj, err = b.SetPodTemplate(obj, tmpl)
	if err != nil {
		return obj, err
	}
	return b.Update(ctx, obj, uopts)
}

func (b *baseClient[T, L, C]) CreateJob(
	ctx context.Context, from string, gopts metav1.GetOptions, copts metav1.CreateOptions,
) (*batchv1.Job, error) {
	_, tmpl, err := b.getWithPodTemplate(ctx, from, gopts)
	if err != nil {
		return nil, err
	}
	// Nanosecond timestamp plus a random component: the former second-level
	// precision made two jobs created within the same second collide. Name
	// uniqueness has no security relevance, so a non-cryptographic PRNG is fine.
	suffix := strconv.FormatInt(time.Now().UnixNano(), 36) + "-" + strconv.FormatUint(uint64(rand.Uint32()), 36) //nolint:gosec
	return Job(b.GetClient(), b.GetNamespace()).Create(ctx, &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-%s", from, suffix),
			Namespace: b.GetNamespace(),
		},
		Spec: batchv1.JobSpec{
			Template: *tmpl,
		},
	}, copts)
}
