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

package vald

import (
	"context"
	"reflect"
	"slices"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/k8s"
	"github.com/vdaas/vald/internal/k8s/reconciler"
	"github.com/vdaas/vald/internal/k8s/resource"
	"github.com/vdaas/vald/internal/log"
)

// WatcherConfig collects the settings shared by every CRD list watcher built
// through NewListWatcher; V is the value type of the name-keyed map handed to
// OnReconcile. The fields are exported so the thin per-CRD packages (and
// their white-box tests) can alias this type and keep validating options of
// their own.
type WatcherConfig[V any] struct {
	OnError     func(err error)
	OnReconcile func(ctx context.Context, resources map[string]V)
	Labels      map[string]string
	Name        string
	Namespaces  []string
}

// WatcherOption mutates WatcherConfig; every watcher package exposes its
// public Option type as an instantiation alias of this type.
type WatcherOption[V any] func(*WatcherConfig[V]) error

// WithWatcherControllerName returns the option to set the controller name; an
// empty name is rejected with errors.ErrInvalidOption.
func WithWatcherControllerName[V any](name string) WatcherOption[V] {
	return func(c *WatcherConfig[V]) error {
		if name == "" {
			return errors.NewErrInvalidOption("controllerName", name)
		}
		c.Name = name
		return nil
	}
}

// WithWatcherNamespaces returns the option to restrict List to the given
// namespaces, replacing any previously configured set. An empty set or an
// empty namespace element is rejected with errors.ErrInvalidOption.
func WithWatcherNamespaces[V any](nss ...string) WatcherOption[V] {
	return func(c *WatcherConfig[V]) error {
		if len(nss) == 0 {
			return errors.NewErrInvalidOption("namespaces", nss)
		}
		if slices.Contains(nss, "") {
			return errors.NewErrInvalidOption("namespaces", nss)
		}
		c.Namespaces = nss
		return nil
	}
}

// WithWatcherLabels returns the option to restrict List to objects matching
// every given label, replacing any previously configured set. An empty (or
// nil) label set is rejected with errors.ErrInvalidOption.
func WithWatcherLabels[V any](labels map[string]string) WatcherOption[V] {
	return func(c *WatcherConfig[V]) error {
		if len(labels) == 0 {
			return errors.NewErrInvalidOption("labels", labels)
		}
		c.Labels = labels
		return nil
	}
}

// WithWatcherOnErrorFunc returns the option to set the reconcile error
// callback; a nil callback is rejected with errors.ErrInvalidOption.
func WithWatcherOnErrorFunc[V any](f func(err error)) WatcherOption[V] {
	return func(c *WatcherConfig[V]) error {
		if f == nil {
			return errors.NewErrInvalidOption("onErrorFunc", f)
		}
		c.OnError = f
		return nil
	}
}

// WithWatcherOnReconcileFunc returns the option to set the callback invoked
// with the name-keyed resource map on every reconcile; a nil callback is
// rejected with errors.ErrInvalidOption.
func WithWatcherOnReconcileFunc[V any](
	f func(ctx context.Context, resources map[string]V),
) WatcherOption[V] {
	return func(c *WatcherConfig[V]) error {
		if f == nil {
			return errors.NewErrInvalidOption("onReconcileFunc", f)
		}
		c.OnReconcile = f
		return nil
	}
}

// WatcherOptionErrorPolicy decides how NewListWatcher reacts to an option
// application error: abort reports whether construction must stop and reason
// carries the error to return in that case.
type WatcherOptionErrorPolicy func(err error, opt any) (abort bool, reason error)

// SkipNonCriticalOptionError wraps the error with the failing option's
// identity and aborts only on an errors.ErrCriticalOption; any other option
// error is logged as a warning and construction continues (the unified policy
// of every watcher package in this tree).
func SkipNonCriticalOptionError(err error, opt any) (abort bool, reason error) {
	oerr := errors.ErrOptionFailed(err, reflect.ValueOf(opt))
	e := &errors.ErrCriticalOption{}
	if errors.As(err, &e) {
		log.Error(oerr)
		return true, oerr
	}
	log.Warn(oerr)
	return false, nil
}

// WatchedObject constrains PT to a CRD pointer type usable both as a
// k8s.Object (watch installation and name keying) and as the item pointer of
// resource.List (deepcopy), which every watched CRD list type in this tree
// aliases.
type WatchedObject[T any] interface {
	resource.Objectable[T]
	resource.DeepCopyIntoer[T]
}

// staticListOpts is the count of ListOption entries NewListWatcher may append
// (WithAddToScheme, WithOnError, WithLabels) on top of the per-namespace
// ones, used only to size the lopts slice up front.
const staticListOpts = 3

// NewListWatcher builds a batch-type ResourceController for the CRD T: on
// every reconcile it lists all matching objects, converts each item via
// convert and hands the result to the OnReconcile callback as a map keyed by
// object name. It consolidates the boilerplate previously copied across the
// benchmark job/scenario and mirror target watcher packages: option
// application under the given error policy, scheme registration, namespace
// scoping and label selection.
func NewListWatcher[T any, PT WatchedObject[T], V any](
	addToScheme func(s *k8s.Scheme) error,
	convert func(obj PT) V,
	policy WatcherOptionErrorPolicy,
	opts ...WatcherOption[V],
) (k8s.ResourceController, error) {
	c := new(WatcherConfig[V])
	for _, opt := range opts {
		if err := opt(c); err != nil {
			if abort, reason := policy(err, opt); abort {
				return nil, reason
			}
		}
	}

	lopts := make([]reconciler.ListOption, 0, len(c.Namespaces)+staticListOpts)
	lopts = append(lopts, reconciler.WithAddToScheme(addToScheme))
	if c.OnError != nil {
		lopts = append(lopts, reconciler.WithOnError(c.OnError))
	}
	// Reflect every configured namespace as its own ListOption so it scopes
	// the List call. controller-runtime keeps a single ListOptions.Namespace,
	// so with multiple namespaces the last one wins (known pre-existing
	// limitation, kept as-is).
	for _, ns := range c.Namespaces {
		lopts = append(lopts, reconciler.WithNamespace(ns))
	}
	if len(c.Labels) > 0 {
		lopts = append(lopts, reconciler.WithLabels(c.Labels))
	}

	onReconcile := c.OnReconcile
	return reconciler.NewListReconciler(
		c.Name,
		PT(new(T)),
		func(ctx context.Context, list *resource.List[T, PT]) {
			if onReconcile == nil {
				return
			}
			m := make(map[string]V, len(list.Items))
			for i := range list.Items {
				item := PT(&list.Items[i])
				m[item.GetName()] = convert(item)
			}
			onReconcile(ctx, m)
		},
		lopts...,
	), nil
}
