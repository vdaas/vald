// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
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

package reconciler

import (
	"context"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/k8s"
	"github.com/vdaas/vald/internal/k8s/client"
	"github.com/vdaas/vald/internal/log"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

// baseReconciler holds the fields and methods common to listReconciler and
// objectReconciler, so each concrete type only embeds it instead of repeating
// them.
type baseReconciler struct {
	mgr           manager.Manager
	client        client.Client
	fieldIndexes  map[string]func(o k8s.Object) []string
	name          string
	maxConcurrent int
	// initErr records a NewReconciler initialization failure (missing
	// manager, scheme or index registration) so Reconcile can surface it
	// instead of silently continuing with a broken setup.
	initErr error
}

func (b *baseReconciler) GetName() string { return b.name }

// MaxConcurrentReconciles implements the optional k8s.ConcurrentReconciler
// interface: values greater than zero request that many reconcile workers.
func (b *baseReconciler) MaxConcurrentReconciles() int { return b.maxConcurrent }

// checkReady reports the recorded initialization failure, or a missing
// manager, so Reconcile implementations can fail fast instead of running
// against a broken setup.
func (b *baseReconciler) checkReady() error {
	if b.initErr != nil {
		return b.initErr
	}
	if b.mgr == nil {
		return errors.Errorf("manager is not registered for %s reconciler", b.name)
	}
	return nil
}

// For, Owns and Watches return nil by default; concrete reconcilers override
// only the hook they actually install (listReconciler overrides Watches,
// objectReconciler overrides For).

func (*baseReconciler) For() (k8s.Object, []builder.ForOption) { return nil, nil }

func (*baseReconciler) Owns() (k8s.Object, []builder.OwnsOption) { return nil, nil }

func (*baseReconciler) Watches() (k8s.Object, handler.EventHandler, []builder.WatchesOption) {
	return nil, nil, nil
}

// addFieldIndex registers a cache index. No-op when field or indexer is nil/empty.
func (b *baseReconciler) addFieldIndex(field string, indexer func(o k8s.Object) []string) {
	if field == "" || indexer == nil {
		return
	}
	if b.fieldIndexes == nil {
		b.fieldIndexes = make(map[string]func(o k8s.Object) []string, 1)
	}
	b.fieldIndexes[field] = indexer
}

// setup initializes the manager, registers the scheme, and registers field
// indexes. Failures are recorded in initErr (and surfaced by Reconcile via
// checkReady) because NewReconciler cannot return an error.
func (b *baseReconciler) setup(
	ctx context.Context,
	mgr manager.Manager,
	addToScheme func(*runtime.Scheme) error,
	indexOn k8s.Object,
) {
	if b.mgr == nil && mgr != nil {
		b.mgr = mgr
	}
	if b.mgr == nil {
		b.initErr = b.checkReady() // records the missing-manager error
		log.Error(b.initErr)
		return
	}
	if b.client == nil {
		c, err := client.NewFromManager(b.mgr)
		if err != nil {
			b.initErr = errors.Wrapf(err, "failed to build client for %s reconciler", b.name)
			log.Error(b.initErr)
			return
		}
		b.client = c
	}
	if addToScheme == nil {
		addToScheme = clientgoscheme.AddToScheme
	}
	if err := addToScheme(b.mgr.GetScheme()); err != nil {
		b.initErr = errors.Wrapf(err, "failed to register scheme for %s reconciler", b.name)
		log.Error(b.initErr)
		return
	}
	for field, indexer := range b.fieldIndexes {
		if err := b.mgr.GetFieldIndexer().IndexField(ctx, indexOn, field, indexer); err != nil {
			b.initErr = errors.Wrapf(err, "failed to register field index %s for %s reconciler", field, b.name)
			log.Error(b.initErr)
			return
		}
	}
}
