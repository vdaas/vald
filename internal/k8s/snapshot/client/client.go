// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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
package client

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	snapshotv1 "github.com/kubernetes-csi/external-snapshotter/client/v6/apis/volumesnapshot/v1"
	snapshotclient "github.com/kubernetes-csi/external-snapshotter/client/v6/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	watch "k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/clientcmd"
)

type (
	ListOptions    = metav1.ListOptions
	CreateOptions  = metav1.CreateOptions
	DeleteOptions  = metav1.DeleteOptions
	Lists          = snapshotv1.VolumeSnapshotList
	VolumeSnapshot = snapshotv1.VolumeSnapshot
	Wacher         = watch.Interface
)

type Client interface {
	List(context.Context, ListOptions) (*Lists, error)
	Create(context.Context, *VolumeSnapshot, CreateOptions) (*VolumeSnapshot, error)
	Watch(context.Context, ListOptions) (watch.Interface, error)
	Delete(context.Context, string, DeleteOptions) error
}

type client struct {
	sClientset *snapshotclient.Clientset
	namespace  string
}

// check if client implements Client interface
var _ Client = (*client)(nil)

func New(opts ...Option) (Client, error) {
	c := new(client)
	for _, opt := range append(defaultOpts, opts...) {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	// Get client
	// FIXME: use incluster config
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home dir: %w", err)
	}

	kubeconfig := filepath.Join(homeDir, ".kube", "config")

	cfg, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("failed to build config from flags: %w", err)
	}

	sclient, err := snapshotclient.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create snapshot client: %w", err)
	}
	c.sClientset = sclient

	return c, nil
}

func (c *client) List(ctx context.Context, options ListOptions) (*Lists, error) {
	return c.sClientset.SnapshotV1().VolumeSnapshots(c.namespace).List(ctx, options)
}

func (c *client) Create(ctx context.Context, snapshot *VolumeSnapshot, options CreateOptions) (*VolumeSnapshot, error) {
	return c.sClientset.SnapshotV1().VolumeSnapshots(c.namespace).Create(ctx, snapshot, options)
}

func (c *client) Watch(ctx context.Context, options ListOptions) (Wacher, error) {
	return c.sClientset.SnapshotV1().VolumeSnapshots(c.namespace).Watch(ctx, options)
}

func (c *client) Delete(ctx context.Context, name string, options DeleteOptions) error {
	return c.sClientset.SnapshotV1().VolumeSnapshots(c.namespace).Delete(ctx, name, options)
}
