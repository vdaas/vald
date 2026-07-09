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
	kclient "sigs.k8s.io/controller-runtime/pkg/client"
)

// Objectable constrains PT to a pointer to T implementing client.Object, so
// generic code can materialize fresh instances via PT(new(T)).
type Objectable[T any] interface {
	*T
	Object
}

// ObjectClient binds the scheme-aware ObjectAPI to a single object type such
// as a custom resource, so callers neither repeat type arguments nor
// materialize empty objects by hand:
//
//	vor := resource.NewObjectClient[v1.ValdOperatorRelease](mgr.GetClient())
//	obj, err := vor.Get(ctx, name, namespace)
type ObjectClient[T any, PT Objectable[T]] struct {
	api ObjectAPI
}

func NewObjectClient[T any, PT Objectable[T]](api ObjectAPI) *ObjectClient[T, PT] {
	return &ObjectClient[T, PT]{api: api}
}

// Get fetches the object identified by name and namespace. Pass an empty
// namespace for cluster-scoped objects.
func (c *ObjectClient[T, PT]) Get(ctx context.Context, name, namespace string) (PT, error) {
	return GetObject(ctx, c.api, name, namespace, PT(new(T)))
}

func (c *ObjectClient[T, PT]) UpdateStatus(ctx context.Context, obj PT) error {
	return c.api.Status().Update(ctx, obj)
}

const (
	defaultWaitInterval = 5 * time.Second
	defaultWaitTimeout  = 5 * time.Minute
)

// Wait polls the named object until eval reports done or the context or the
// default timeout expires. NotFound errors keep the poll running so callers
// can wait for objects that do not exist yet.
func (c *ObjectClient[T, PT]) Wait(
	ctx context.Context, name, namespace string, eval func(PT) (done bool, err error),
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
			return false, errors.ErrWaitTimeoutFor(namespace, name)
		case <-ticker.C:
			obj, err := c.Get(ctx, name, namespace)
			if err != nil {
				if kclient.IgnoreNotFound(err) == nil {
					continue
				}
				return false, err
			}
			done, err := eval(obj)
			if err != nil {
				return false, err
			}
			if done {
				return true, nil
			}
		}
	}
}
