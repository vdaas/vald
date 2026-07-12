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

package client

import (
	"context"
	"maps"
	"time"

	snapshotv1 "github.com/kubernetes-csi/external-snapshotter/client/v6/apis/volumesnapshot/v1"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/k8s"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apimachinery/pkg/watch"
	applycorev1 "k8s.io/client-go/applyconfigurations/core/v1"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	cli "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

var NewSelector = labels.NewSelector

// Client is the unified generic Kubernetes client across the vald repository.
// It acts as the single source of truth for safe, scheme-aware CRUD operations,
// replacing previous distinct StandaloneClient, ObjectClient, ListClient,
// and legacy typed e2e clients.
type Client[T k8s.Object, L k8s.ObjectList] interface {
	Get(ctx context.Context, name, namespace string) (T, error)
	List(ctx context.Context, opts ...k8s.ListOption) (L, error)
	Create(ctx context.Context, obj T, opts ...k8s.CreateOption) error
	Update(ctx context.Context, obj T, opts ...k8s.UpdateOption) error
	UpdateStatus(ctx context.Context, obj T, opts ...k8s.SubResourceUpdateOption) error
	Delete(ctx context.Context, obj T, opts ...k8s.DeleteOption) error
	Watch(ctx context.Context, opts ...k8s.ListOption) (watch.Interface, error)
	Wait(ctx context.Context, name, namespace string, eval func(T) (done bool, err error)) (bool, error)
	Raw() k8s.ClientWithWatch
}

type client[T k8s.Object, L k8s.ObjectList] struct {
	api k8s.ClientWithWatch
	t   T
	l   L
}

func New[T k8s.Object, L k8s.ObjectList](t T, l L, opts ...Option) (Client[T, L], error) {
	o := new(options)
	o.scheme = runtime.NewScheme()
	for _, opt := range opts {
		if err := opt(o); err != nil {
			return nil, err
		}
	}

	if err := clientgoscheme.AddToScheme(o.scheme); err != nil {
		return nil, err
	}
	if err := snapshotv1.AddToScheme(o.scheme); err != nil {
		return nil, err
	}

	api, err := cli.NewWithWatch(ctrl.GetConfigOrDie(), cli.Options{
		Scheme: o.scheme,
	})
	if err != nil {
		return nil, err
	}

	return &client[T, L]{
		api: api,
		t:   t,
		l:   l,
	}, nil
}

// NewWithClient creates a new generic Client wrapping an existing controller-runtime client.
func NewWithClient[T k8s.Object, L k8s.ObjectList](api k8s.ClientWithWatch, t T, l L) Client[T, L] {
	return &client[T, L]{
		api: api,
		t:   t,
		l:   l,
	}
}

// NewWithConfig creates a new generic Client from a REST config.
func NewWithConfig[T k8s.Object, L k8s.ObjectList](restCfg *rest.Config, t T, l L, opts ...Option) (Client[T, L], error) {
	// this is a simplified version, actually let's just make it properly
	o := new(options)
	o.scheme = runtime.NewScheme()
	if err := clientgoscheme.AddToScheme(o.scheme); err != nil {
		return nil, err
	}
	if err := snapshotv1.AddToScheme(o.scheme); err != nil {
		return nil, err
	}
	
	api, err := cli.NewWithWatch(restCfg, cli.Options{
		Scheme: o.scheme,
	})
	if err != nil {
		return nil, err
	}
	return &client[T, L]{
		api: api,
		t:   t,
		l:   l,
	}, nil
}

func (c *client[T, L]) Get(ctx context.Context, name, namespace string) (T, error) {
	obj := c.t.DeepCopyObject().(T)
	err := c.api.Get(ctx, cli.ObjectKey{Name: name, Namespace: namespace}, obj)
	return obj, err
}

func (c *client[T, L]) List(ctx context.Context, opts ...k8s.ListOption) (L, error) {
	list := c.l.DeepCopyObject().(L)
	err := c.api.List(ctx, list, opts...)
	return list, err
}

func (c *client[T, L]) Create(ctx context.Context, obj T, opts ...k8s.CreateOption) error {
	return c.api.Create(ctx, obj, opts...)
}

func (c *client[T, L]) Update(ctx context.Context, obj T, opts ...cli.UpdateOption) error {
	return c.api.Update(ctx, obj, opts...)
}

func (c *client[T, L]) UpdateStatus(ctx context.Context, obj T, opts ...cli.SubResourceUpdateOption) error {
	return c.api.Status().Update(ctx, obj, opts...)
}

func (c *client[T, L]) Delete(ctx context.Context, obj T, opts ...k8s.DeleteOption) error {
	return c.api.Delete(ctx, obj, opts...)
}

func (c *client[T, L]) Watch(ctx context.Context, opts ...k8s.ListOption) (watch.Interface, error) {
	list := c.l.DeepCopyObject().(L)
	return c.api.Watch(ctx, list, opts...)
}

func (c *client[T, L]) Raw() k8s.ClientWithWatch {
	return c.api
}

func MatchingLabelsString(selector string) k8s.ListOption {
	sel, err := labels.Parse(selector)
	if err == nil {
		return cli.MatchingLabelsSelector{Selector: sel}
	}
	return cli.MatchingLabels{}
}

const (
	defaultWaitInterval = 5 * time.Second
	defaultWaitTimeout  = 5 * time.Minute
)

func WaitLoop(ctx context.Context, onTimeout error, step func(context.Context) (done bool, err error)) (bool, error) {
	ticker := time.NewTicker(defaultWaitInterval)
	defer ticker.Stop()
	timeout := time.NewTimer(defaultWaitTimeout)
	defer timeout.Stop()
	for {
		select {
		case <-ctx.Done():
			return false, ctx.Err()
		case <-timeout.C:
			return false, onTimeout
		case <-ticker.C:
			done, err := step(ctx)
			if err != nil {
				return false, err
			}
			if done {
				return true, nil
			}
		}
	}
}

func (c *client[T, L]) Wait(ctx context.Context, name, namespace string, eval func(T) (done bool, err error)) (bool, error) {
	return WaitLoop(ctx, errors.ErrWaitTimeoutFor(namespace, name), func(ctx context.Context) (bool, error) {
		obj, err := c.Get(ctx, name, namespace)
		if err != nil {
			if cli.IgnoreNotFound(err) == nil {
				return false, nil
			}
			return false, err
		}
		return eval(obj)
	})
}

// NewLabelSelector is a helper to construct a labels.Selector easily
func NewLabelSelector(key string, op selection.Operator, vals []string) (labels.Selector, error) {
	requirements, err := labels.NewRequirement(key, op, vals)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create requirement on creating label selector")
	}
	return labels.NewSelector().Add(*requirements), nil
}

// DeploymentListClient is the generic client instantiation for Deployment, shared
// by callers (index operator, read-replica rotator) that would otherwise
// repeat this 2-argument spelling verbatim.
type DeploymentListClient = Client[*k8s.Deployment, *k8s.DeploymentList]

// PodPredicates returns a builder.Predicates with the given filter function.
func PodPredicates(filter func(pod *corev1.Pod) bool) builder.Predicates {
	return builder.WithPredicates(predicate.Funcs{
		CreateFunc: func(e event.CreateEvent) bool {
			pod, ok := e.Object.(*corev1.Pod)
			if !ok {
				return false
			}
			return filter(pod)
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			pod, ok := e.Object.(*corev1.Pod)
			if !ok {
				return false
			}
			return filter(pod)
		},
		UpdateFunc: func(e event.UpdateEvent) bool {
			pod, ok := e.ObjectNew.(*corev1.Pod)
			if !ok {
				return false
			}
			return filter(pod)
		},
		GenericFunc: func(e event.GenericEvent) bool {
			pod, ok := e.Object.(*corev1.Pod)
			if !ok {
				return false
			}
			return filter(pod)
		},
	})
}

// Patcher is an interface for patching resources with controller-runtime client.
type Patcher interface {
	// ApplyPodAnnotations applies the given annotations to the agent pod with server-side apply.
	ApplyPodAnnotations(ctx context.Context, name, namespace string, entries map[string]string) error
}

type patcher struct {
	client       k8s.ClientWithWatch
	fieldManager string
}

func NewPatcher(fieldManager string) (Patcher, error) {
	scheme := runtime.NewScheme()
	if err := clientgoscheme.AddToScheme(scheme); err != nil {
		return nil, err
	}
	c, err := cli.NewWithWatch(ctrl.GetConfigOrDie(), cli.Options{
		Scheme: scheme,
	})
	if err != nil {
		return nil, err
	}

	return &patcher{
		client:       c,
		fieldManager: fieldManager,
	}, nil
}

func (s *patcher) ApplyPodAnnotations(
	ctx context.Context, name, namespace string, entries map[string]string,
) error {
	var podList corev1.PodList
	if err := s.client.List(ctx, &podList, &cli.ListOptions{
		Namespace:     namespace,
		FieldSelector: fields.OneTermEqualSelector("metadata.name", name),
	}); err != nil {
		return err
	}

	if len(podList.Items) == 0 {
		return errors.New("agent pod not found on exporting metrics")
	}

	//nolint:gomnd
	if len(podList.Items) >= 2 {
		return errors.New("multiple agent pods found on exporting metrics. pods with same name exist in the same namespace?")
	}
	pod := podList.Items[0]

	curApplyConfig, err := applycorev1.ExtractPod(&pod, s.fieldManager)
	if err != nil {
		return err
	}

	// check if there is any diffs in the annotations
	annotations := pod.GetObjectMeta().GetAnnotations()
	if annotations == nil {
		annotations = make(map[string]string)
	}
	maps.Copy(annotations, entries)
	expectPod := applycorev1.Pod(name, namespace).
		WithAnnotations(annotations)

	if equality.Semantic.DeepEqual(expectPod, curApplyConfig) {
		// no change found in the pod annotations
		return nil
	}

	// now we found the diffs, apply the changes
	return s.client.Patch(ctx, &pod, cli.Apply, cli.FieldOwner(s.fieldManager), cli.ForceOwnership)
}
