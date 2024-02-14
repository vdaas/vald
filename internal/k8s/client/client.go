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

// Package client is Kubernetes client for getting resource from Kubernetes cluster.
package client

import (
	"context"
	"fmt"

	snapshotv1 "github.com/kubernetes-csi/external-snapshotter/client/v6/apis/volumesnapshot/v1"
	"github.com/vdaas/vald/internal/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apimachinery/pkg/watch"
	applycorev1 "k8s.io/client-go/applyconfigurations/core/v1"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	cli "sigs.k8s.io/controller-runtime/pkg/client"
)

type (
	Object             = cli.Object
	ObjectKey          = cli.ObjectKey
	DeleteAllOfOptions = cli.DeleteAllOfOptions
	DeleteOptions      = cli.DeleteOptions
	ListOptions        = cli.ListOptions
	ListOption         = cli.ListOption
	CreateOption       = cli.CreateOption
	CreateOptions      = cli.CreateOptions
	UpdateOptions      = cli.UpdateOptions
	MatchingLabels     = cli.MatchingLabels
	MatchingFields     = cli.MatchingFields
	InNamespace        = cli.InNamespace
	VolumeSnapshot     = snapshotv1.VolumeSnapshot
	PodList            = corev1.PodList
	Pod                = corev1.Pod
	PatchOptions       = cli.PatchOptions
)

const (
	DeletePropagationBackground = metav1.DeletePropagationBackground
	WatchDeletedEvent           = watch.Deleted
	SelectionOpEquals           = selection.Equals
)

var (
	ServerSideApply = cli.Apply
	MergePatch      = cli.Merge
)

type Client interface {
	// Get retrieves an obj for the given object key from the Kubernetes Cluster.
	// obj must be a struct pointer so that obj can be updated with the response
	// returned by the Server.
	Get(ctx context.Context, name string, namespace string, obj Object, opts ...cli.GetOption) error
	// List retrieves list of objects for a given namespace and list options. On a
	// successful call, Items field in the list will be populated with the
	// result returned from the server.
	List(ctx context.Context, list cli.ObjectList, opts ...ListOption) error

	// Create saves the object obj in the Kubernetes cluster. obj must be a
	// struct pointer so that obj can be updated with the content returned by the Server.
	Create(ctx context.Context, obj Object, opts ...CreateOption) error

	// Delete deletes the given obj from Kubernetes cluster.
	Delete(ctx context.Context, obj Object, opts ...cli.DeleteOption) error

	// Update updates the given obj in the Kubernetes cluster. obj must be a
	// struct pointer so that obj can be updated with the content returned by the Server.
	Update(ctx context.Context, obj Object, opts ...cli.UpdateOption) error

	// Patch patches the given obj in the Kubernetes cluster. obj must be a
	// struct pointer so that obj can be updated with the content returned by the Server.
	Patch(ctx context.Context, obj Object, patch cli.Patch, opts ...cli.PatchOption) error

	// Watch watches the given obj for changes and takes the appropriate callbacks.
	Watch(ctx context.Context, obj cli.ObjectList, opts ...ListOption) (watch.Interface, error)

	// LabelSelector creates labels.Selector for Options like ListOptions.
	LabelSelector(key string, op selection.Operator, vals []string) (labels.Selector, error)
}

type client struct {
	scheme    *runtime.Scheme
	withWatch cli.WithWatch
}

func New(opts ...Option) (_ Client, err error) {
	c := new(client)
	if c.scheme == nil {
		c.scheme = runtime.NewScheme()
	}
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	// Add the core schemes
	if err := clientgoscheme.AddToScheme(c.scheme); err != nil {
		return nil, err
	}
	if err := snapshotv1.AddToScheme(c.scheme); err != nil {
		return nil, err
	}

	c.withWatch, err = cli.NewWithWatch(ctrl.GetConfigOrDie(), cli.Options{
		Scheme: c.scheme,
	})
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *client) Get(ctx context.Context, name, namespace string, obj cli.Object, opts ...cli.GetOption) error {
	return c.withWatch.Get(
		ctx,
		cli.ObjectKey{
			Name:      name,
			Namespace: namespace,
		},
		obj,
		opts...,
	)
}

func (c *client) List(ctx context.Context, list cli.ObjectList, opts ...cli.ListOption) error {
	return c.withWatch.List(ctx, list, opts...)
}

func (c *client) Create(ctx context.Context, obj Object, opts ...CreateOption) error {
	return c.withWatch.Create(ctx, obj, opts...)
}

func (c *client) Delete(ctx context.Context, obj Object, opts ...cli.DeleteOption) error {
	return c.withWatch.Delete(ctx, obj, opts...)
}

func (c *client) Update(ctx context.Context, obj Object, opts ...cli.UpdateOption) error {
	return c.withWatch.Update(ctx, obj, opts...)
}

func (c *client) Patch(ctx context.Context, obj Object, patch cli.Patch, opts ...cli.PatchOption) error {
	return c.withWatch.Patch(ctx, obj, patch, opts...)
}

func (c *client) Watch(ctx context.Context, obj cli.ObjectList, opts ...ListOption) (watch.Interface, error) {
	return c.withWatch.Watch(ctx, obj, opts...)
}

func (*client) LabelSelector(key string, op selection.Operator, vals []string) (labels.Selector, error) {
	requirements, err := labels.NewRequirement(key, op, vals)
	if err != nil {
		return nil, fmt.Errorf("failed to create requirement on creating label selector: %w", err)
	}
	return labels.NewSelector().Add(*requirements), nil
}

type Patcher struct {
	client       Client
	fieldManager string
}

func NewPatcher(fieldManager string) (Patcher, error) {
	client, err := New()
	if err != nil {
		return Patcher{}, err
	}

	return Patcher{
		client:       client,
		fieldManager: fieldManager,
	}, nil
}

func (s *Patcher) ApplyPodAnnotations(ctx context.Context, name, namespace string, entries map[string]string) error {
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

	//nolint: gomnd
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
	for k, v := range entries {
		annotations[k] = v
	}
	expectPod := applycorev1.Pod(name, namespace).
		WithAnnotations(annotations)

	if equality.Semantic.DeepEqual(expectPod, curApplyConfig) {
		// no change found in the pod annotations
		return nil
	}

	// now we found the diffs, apply the changes
	obj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(expectPod)
	if err != nil {
		return err
	}

	patch := &unstructured.Unstructured{Object: obj}
	if err := s.client.Patch(ctx, patch, cli.Apply, &cli.PatchOptions{
		FieldManager: s.fieldManager,
		Force:        ptr.To(true),
	}); err != nil {
		return err
	}

	return nil
}
