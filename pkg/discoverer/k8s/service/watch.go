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

// Package service manages the main logic of server.
package service

import (
	"context"

	"github.com/vdaas/vald/internal/k8s"
	"github.com/vdaas/vald/internal/k8s/reconciler"
	"github.com/vdaas/vald/internal/k8s/resource"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/sync"
)

// listSyncController assembles the list-reconciler wiring shared by every
// discoverer watch: it returns a k8s.Option installing a list reconciler
// named name that watches obj and, on every reconcile, mirrors the
// conv-converted snapshot into dst via syncOnReconcile. The common error
// logging ("failed to reconcile <what>:") and namespace/fields/labels
// selector options are assembled here; per-resource options (scheme
// registration, field indexes) are passed through extra.
func listSyncController[L any, PL resource.ListPtr[L], V comparable](
	name, what string,
	obj k8s.Object,
	dst *sync.Map[string, V],
	conv func(list PL) map[string]V,
	namespace string,
	fields, labels map[string]string,
	extra ...reconciler.ListOption,
) k8s.Option {
	return k8s.WithResourceController(reconciler.NewListReconciler(
		name,
		obj,
		syncOnReconcile(dst, conv),
		append([]reconciler.ListOption{
			reconciler.WithOnError(func(err error) {
				log.Error("failed to reconcile "+what+":", err)
			}),
			reconciler.WithNamespace(namespace),
			reconciler.WithFields(fields),
			reconciler.WithLabels(labels),
		}, extra...)...,
	))
}

// syncOnReconcile returns the reconcile callback shared by every discoverer
// watch: it converts the listed snapshot via conv into a map keyed by name
// and mirrors it into dst — every entry is stored, then keys absent from the
// snapshot are deleted (stale pruning).
func syncOnReconcile[L any, V comparable](
	dst *sync.Map[string, V], conv func(list L) map[string]V,
) func(ctx context.Context, list L) {
	return func(_ context.Context, list L) {
		entries := conv(list)
		for name, val := range entries {
			dst.Store(name, val)
		}
		dst.Range(func(name string, _ V) bool {
			if _, ok := entries[name]; !ok {
				dst.Delete(name)
			}
			return true
		})
	}
}

// debugLogged wraps conv so that every conversion result is debug-logged as
// "<what> reconciled" before being returned, preserving the discoverer's
// per-resource reconcile logs.
func debugLogged[I, O any](what string, conv func(I) O) func(I) O {
	return func(in I) O {
		out := conv(in)
		log.Debugf("%s reconciled\t%#v", what, out)
		return out
	}
}

// podEntries converts the listed pods into the per-application-name pointer
// entries stored in d.pods, raising d.maxPods to the largest group seen.
func (d *discoverer) podEntries(list *k8s.PodList) map[string]*[]Pod {
	podList := podsByAppName(list, d.namespace)
	log.Debugf("pod resource reconciled\t%#v", podList)
	entries := make(map[string]*[]Pod, len(podList))
	for name, pods := range podList {
		if n := int64(len(pods)); n > d.maxPods.Load() {
			d.maxPods.Store(n)
		}
		entries[name] = &pods
	}
	return entries
}

// nodeEntries converts the listed nodes into the per-name pointer entries
// stored in d.nodes.
func nodeEntries(list *k8s.NodeList) map[string]*Node {
	nodes := toNodes(list)
	log.Debugf("node resource reconciled\t%#v", nodes)
	entries := make(map[string]*Node, len(nodes))
	for _, node := range nodes {
		entries[node.Name] = &node
	}
	return entries
}
