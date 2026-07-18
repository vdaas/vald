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
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/k8s"
	"github.com/vdaas/vald/internal/k8s/client"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	kclient "sigs.k8s.io/controller-runtime/pkg/client"
)

// cloneSeed returns a zero-value clone of seed via DeepCopyObject, asserting
// the result back to the concrete type S. The assertion always holds because
// DeepCopyObject on a value of concrete type S is generated to return that
// same concrete type. Shared by Client.fresh (S = T) and Client.freshList
// (S = L): both Object and ObjectListType embed runtime.Object, so either
// seed satisfies this constraint.
func cloneSeed[S runtime.Object](seed S) S {
	obj, _ := seed.DeepCopyObject().(S) //nolint:forcetypeassert
	return obj
}

// Objectable constrains PT to a pointer to T implementing client.Object, so
// generic code can materialize fresh instances via PT(new(T)). Kept for
// internal/k8s/reconciler, which has no seed value available at construction
// time (unlike Client below); it does not need to materialize a companion
// list from a distinct construction-time argument.
type Objectable[T any] interface {
	*T
	Object
}

// ListPtr constrains PL to a pointer to L implementing ObjectListType, so
// generic code can materialize fresh list instances via PL(new(L)) the same
// way Objectable does for single objects. Kept for internal/k8s/reconciler;
// see the note on Objectable.
type ListPtr[L any] interface {
	*L
	ObjectListType
}

// Client binds the scheme-aware k8s.Client to a single object Kind T (e.g.
// *appsv1.Deployment, a pointer type) and its companion list Kind L, so
// callers neither repeat type arguments nor materialize empty objects or
// lists by hand. Fresh instances for Get/Wait and List/Watch are produced via
// seed.DeepCopyObject(), so no pointer-constraint type parameters are needed
// — contrast with Objectable[T]'s two-parameter (element, pointer) split,
// which internal/k8s/reconciler still uses because it has no seed value to
// clone at construction time:
//
//	deployments := resource.NewClient(mgr.GetClient(), new(appsv1.Deployment), new(appsv1.DeploymentList))
//	list, err := deployments.List(ctx, kclient.MatchingLabels{"app": "foo"})
type Client[T Object, L ObjectListType] struct {
	api k8s.Client
	// watchAPI holds api re-asserted to kclient.WithWatch when the underlying
	// client supports watching (k8s.Client, controller-runtime's plain Client,
	// has no Watch method; only the wider kclient.WithWatch does). It stays
	// nil for watch-incapable clients such as what mgr.GetClient() returns, in
	// which case Watch and DeleteAndWait report
	// errors.ErrKubernetesClientWatchNotSupported instead of panicking.
	watchAPI kclient.WithWatch
	seed     T
	listSeed L
}

// NewClient binds api to a single object type and its list type. The watch
// capability of api is detected here once: Watch/DeleteAndWait are only
// usable when api also implements kclient.WithWatch.
func NewClient[T Object, L ObjectListType](api k8s.Client, tSeed T, lSeed L) *Client[T, L] {
	c := &Client[T, L]{api: api, seed: tSeed, listSeed: lSeed}
	if w, ok := api.(kclient.WithWatch); ok {
		c.watchAPI = w
	}
	return c
}

// NewClientOf builds a Client from the unified client.Client, for callers
// that construct their client via client.New() instead of receiving one from
// a manager. client.Client embeds k8s.Client directly, so unlike the former
// StandaloneClient split, no bridging through a separate Raw() accessor is
// needed: c already satisfies NewClient's k8s.Client parameter.
func NewClientOf[T Object, L ObjectListType](c client.Client, tSeed T, lSeed L) *Client[T, L] {
	return NewClient(c, tSeed, lSeed)
}

// fresh returns a zero-value clone of seed for Get/Wait to populate.
func (c *Client[T, L]) fresh() T {
	return cloneSeed(c.seed)
}

// freshList returns a zero-value clone of listSeed for List/Watch to
// populate.
func (c *Client[T, L]) freshList() L {
	return cloneSeed(c.listSeed)
}

// Get fetches the object identified by name and namespace. Pass an empty
// namespace for cluster-scoped objects.
func (c *Client[T, L]) Get(ctx context.Context, name, namespace string) (T, error) {
	return GetObject(ctx, c.api, name, namespace, c.fresh())
}

func (c *Client[T, L]) UpdateStatus(ctx context.Context, obj T) error {
	return c.api.Status().Update(ctx, obj)
}

// Create saves obj in the Kubernetes cluster.
func (c *Client[T, L]) Create(ctx context.Context, obj T, opts ...kclient.CreateOption) error {
	return c.api.Create(ctx, obj, opts...)
}

// Update updates obj in the Kubernetes cluster.
func (c *Client[T, L]) Update(ctx context.Context, obj T, opts ...kclient.UpdateOption) error {
	return c.api.Update(ctx, obj, opts...)
}

// Delete deletes obj from the Kubernetes cluster.
func (c *Client[T, L]) Delete(ctx context.Context, obj T, opts ...kclient.DeleteOption) error {
	return c.api.Delete(ctx, obj, opts...)
}

const (
	defaultWaitInterval = 5 * time.Second
	defaultWaitTimeout  = 5 * time.Minute
)

// waitLoop runs step every defaultWaitInterval until step reports done, the
// context is canceled, or defaultWaitTimeout elapses (returning onTimeout).
// It is the single polling skeleton shared by Client.Wait and WaitForStatus.
func waitLoop(
	ctx context.Context, onTimeout error, step func(context.Context) (done bool, err error),
) (bool, error) {
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

// Wait polls the named object until eval reports done or the context or the
// default timeout expires. NotFound errors keep the poll running so callers
// can wait for objects that do not exist yet.
func (c *Client[T, L]) Wait(
	ctx context.Context, name, namespace string, eval func(T) (done bool, err error),
) (bool, error) {
	return waitLoop(ctx, errors.ErrWaitTimeoutFor(namespace, name),
		func(ctx context.Context) (bool, error) {
			obj, err := c.Get(ctx, name, namespace)
			if err != nil {
				if kclient.IgnoreNotFound(err) == nil {
					return false, nil
				}
				return false, err
			}
			return eval(obj)
		})
}

// List retrieves every object matching opts into a freshly constructed list.
func (c *Client[T, L]) List(ctx context.Context, opts ...ListOption) (L, error) {
	return ListObjects(ctx, c.api, c.freshList(), opts...)
}

// Watch watches the Kind fixed by the type parameter L for changes scoped by
// opts. The freshly constructed list passed to the underlying client is only
// a GVK carrier — controller-runtime resolves the watched resource from the
// list's type and ignores its contents. When the client passed to NewClient
// does not support watching, ErrKubernetesClientWatchNotSupported is
// returned.
func (c *Client[T, L]) Watch(ctx context.Context, opts ...ListOption) (watch.Interface, error) {
	if c.watchAPI == nil {
		return nil, errors.ErrKubernetesClientWatchNotSupported
	}
	return c.watchAPI.Watch(ctx, c.freshList(), opts...)
}

// DeleteAndWait deletes obj and blocks until a Deleted event scoped by opts is
// observed or ctx is canceled. The watch is started before the delete so the
// deletion event cannot be missed; any Deleted event matching opts releases
// the wait, mirroring a label-selector-scoped watch.
func (c *Client[T, L]) DeleteAndWait(ctx context.Context, obj T, opts ...ListOption) error {
	watcher, err := c.Watch(ctx, opts...)
	if err != nil {
		return errors.Wrapf(err, "failed to watch %T(%s)", obj, obj.GetName())
	}
	defer watcher.Stop()

	eg, egctx := errgroup.New(ctx)
	eg.Go(func() error {
		log.Infof("deleting %T(%s)...", obj, obj.GetName())
		log.Debugf("%T detail: %#v", obj, obj)
		for {
			select {
			case <-egctx.Done():
				return egctx.Err()
			case event := <-watcher.ResultChan():
				if event.Type == watch.Deleted {
					log.Infof("%T(%s) deleted", obj, obj.GetName())
					return nil
				}
				log.Debugf("watching %T(%s) events. event: %v", obj, obj.GetName(), event.Type)
			}
		}
	})

	if err := c.Delete(ctx, obj); err != nil {
		return errors.Wrapf(err, "failed to delete %T(%s)", obj, obj.GetName())
	}
	return eg.Wait()
}
