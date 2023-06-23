//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
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

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	cli "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

type (
	Object             = cli.Object
	ObjectKey          = cli.ObjectKey
	DeleteAllOfOptions = cli.DeleteAllOfOptions
	DeleteOptions      = cli.DeleteOptions
	ListOptions        = cli.ListOptions
	MatchingLabels     = cli.MatchingLabels
	InNamespace        = cli.InNamespace
)

const (
	DeletePropagationBackground = metav1.DeletePropagationBackground
)

type Client interface {
	// Get retrieves an obj for the given object key from the Kubernetes Cluster.
	// obj must be a struct pointer so that obj can be updated with the response
	// returned by the Server.
	Get(ctx context.Context, name string, namespace string, obj cli.Object, opts ...cli.GetOption) error
	// List retrieves list of objects for a given namespace and list options. On a
	// successful call, Items field in the list will be populated with the
	// result returned from the server.
	List(ctx context.Context, list cli.ObjectList, opts ...cli.ListOption) error
}

type client struct {
	scheme *runtime.Scheme
	reader cli.Reader
}

func New(opts ...Option) (Client, error) {
	c := new(client)
	if c.scheme == nil {
		c.scheme = runtime.NewScheme()
	}
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), manager.Options{
		Scheme: c.scheme,
	})
	if err != nil {
		return nil, err
	}
	c.reader = mgr.GetAPIReader()
	return c, nil
}

func (c *client) Get(ctx context.Context, name string, namespace string, obj cli.Object, opts ...cli.GetOption) error {
	return c.reader.Get(
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
	return c.reader.List(ctx, list, opts...)
}
