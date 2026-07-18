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
package target

import (
	"context"

	"github.com/vdaas/vald/internal/k8s/vald"
)

// Option represents the functional option for the settings collected by New;
// it is the Endpoint instantiation of the shared watcher factory option, which
// validates every value it is given.
type Option = vald.WatcherOption[Endpoint]

// WithControllerName returns the option to set the name of controller; an
// empty name is rejected with errors.ErrInvalidOption.
func WithControllerName(name string) Option {
	return vald.WithWatcherControllerName[Endpoint](name)
}

// WithOnErrorFunc returns the option to set the function to notify an error;
// a nil callback is rejected with errors.ErrInvalidOption.
func WithOnErrorFunc(f func(error)) Option {
	return vald.WithWatcherOnErrorFunc[Endpoint](f)
}

// WithOnReconcileFunc returns the option to set the function to get the
// reconciled result; a nil callback is rejected with errors.ErrInvalidOption.
func WithOnReconcileFunc(f func(context.Context, map[string]Endpoint)) Option {
	return vald.WithWatcherOnReconcileFunc(f)
}

// WithNamespaces returns the option to restrict List to the given namespaces,
// replacing any previously configured set. controller-runtime keeps a single
// ListOptions.Namespace, so when multiple namespaces are given the last one
// wins (known pre-existing limitation of the shared watcher factory, kept
// as-is).
func WithNamespaces(nss ...string) Option {
	return vald.WithWatcherNamespaces[Endpoint](nss...)
}

// WithLabels returns the option to set the label selector to get resources
// matching the given label; an empty (or nil) label set is rejected with
// errors.ErrInvalidOption.
func WithLabels(labels map[string]string) Option {
	return vald.WithWatcherLabels[Endpoint](labels)
}
