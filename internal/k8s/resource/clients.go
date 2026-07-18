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
	"github.com/vdaas/vald/internal/k8s/client"
	"github.com/vdaas/vald/internal/sync"
)

// Clients lazily constructs and caches one typed per-Kind client per
// namespace from a single shared client.Client, so repeated lookups for the
// same (Kind, namespace) pair — e.g. one per Kubernetes action in an e2e
// scenario — reuse the same client instead of reconstructing it every time.
// Namespace is still resolved per call, so callers that legitimately target
// different namespaces within one run are unaffected; only repeated (Kind,
// namespace) pairs are deduplicated.
type Clients struct {
	c client.Client

	pods                           map[string]PodClient
	deployments                    map[string]DeploymentClient
	daemonSets                     map[string]DaemonSetClient
	statefulSets                   map[string]StatefulSetClient
	jobs                           map[string]JobClient
	cronJobs                       map[string]CronJobClient
	services                       map[string]ServiceClient
	secrets                        map[string]SecretClient
	configMaps                     map[string]ConfigMapClient
	endpoints                      map[string]EndpointClient
	mutatingWebhookConfiguration   MutatingWebhookConfigurationClient
	validatingWebhookConfiguration ValidatingWebhookConfigurationClient
	// mu holds no pointers; kept last so it trails the scanned pointer region.
	mu sync.RWMutex
}

// NewClients returns a Clients bundle backed by c. Construction is cheap and
// lazy: no per-Kind client is built until its accessor is first called.
func NewClients(c client.Client) *Clients {
	return &Clients{c: c}
}

// cached is the double-checked-locking skeleton shared by cachedNamespaced
// and cachedSingleton: a read-locked fast path, then a write-locked re-check
// before building. lookup reports the cached value and whether it is present;
// build constructs, stores, and returns the value and only ever runs while
// the write lock is held.
func cached[T any](mu *sync.RWMutex, lookup func() (T, bool), build func() T) T {
	mu.RLock()
	if v, ok := lookup(); ok {
		mu.RUnlock()
		return v
	}
	mu.RUnlock()

	mu.Lock()
	defer mu.Unlock()
	if v, ok := lookup(); ok {
		return v
	}
	return build()
}

// cachedNamespaced returns the cached client for namespace, constructing and
// storing it via construct(cs.c, namespace) on first use — the factory is
// passed as a plain function value (Pod, Deployment, ...) so accessors need no
// binding closure. Safe for concurrent use. T is unconstrained beyond any: the
// cache never compares T values (only presence in the map matters), so no
// comparability requirement applies here.
func cachedNamespaced[T any](
	cs *Clients, cache *map[string]T, namespace string, construct func(client.Client, string) T,
) T {
	return cached(&cs.mu,
		func() (T, bool) {
			v, ok := (*cache)[namespace]
			return v, ok
		},
		func() T {
			v := construct(cs.c, namespace)
			if *cache == nil {
				*cache = make(map[string]T, 1)
			}
			(*cache)[namespace] = v
			return v
		})
}

// cachedSingleton returns *cache if it has already been set (compared
// against T's zero value, e.g. nil for an interface), constructing and
// storing it via construct(cs.c) on first use. Safe for concurrent use.
// Unlike cachedNamespaced, this does compare T values, so T must be
// comparable. The singleton stays a plain struct field rather than a
// single-entry map so its memory layout is unchanged by sharing cached.
func cachedSingleton[T comparable](cs *Clients, cache *T, construct func(client.Client) T) T {
	var zero T
	return cached(&cs.mu,
		func() (T, bool) {
			v := *cache
			return v, v != zero
		},
		func() T {
			*cache = construct(cs.c)
			return *cache
		})
}

func (cs *Clients) Pod(namespace string) PodClient {
	return cachedNamespaced(cs, &cs.pods, namespace, Pod)
}

func (cs *Clients) Deployment(namespace string) DeploymentClient {
	return cachedNamespaced(cs, &cs.deployments, namespace, Deployment)
}

func (cs *Clients) DaemonSet(namespace string) DaemonSetClient {
	return cachedNamespaced(cs, &cs.daemonSets, namespace, DaemonSet)
}

func (cs *Clients) StatefulSet(namespace string) StatefulSetClient {
	return cachedNamespaced(cs, &cs.statefulSets, namespace, StatefulSet)
}

func (cs *Clients) Job(namespace string) JobClient {
	return cachedNamespaced(cs, &cs.jobs, namespace, Job)
}

func (cs *Clients) CronJob(namespace string) CronJobClient {
	return cachedNamespaced(cs, &cs.cronJobs, namespace, CronJob)
}

func (cs *Clients) Service(namespace string) ServiceClient {
	return cachedNamespaced(cs, &cs.services, namespace, Service)
}

func (cs *Clients) Secret(namespace string) SecretClient {
	return cachedNamespaced(cs, &cs.secrets, namespace, Secret)
}

func (cs *Clients) ConfigMap(namespace string) ConfigMapClient {
	return cachedNamespaced(cs, &cs.configMaps, namespace, ConfigMap)
}

func (cs *Clients) Endpoints(namespace string) EndpointClient {
	return cachedNamespaced(cs, &cs.endpoints, namespace, Endpoints)
}

// MutatingWebhookConfiguration and ValidatingWebhookConfiguration are
// cluster-scoped, so they cache a single instance rather than one per
// namespace. Inside each method the bare factory identifier resolves to the
// package-level function of the same name (methods are not in identifier
// scope).

func (cs *Clients) MutatingWebhookConfiguration() MutatingWebhookConfigurationClient {
	return cachedSingleton(cs, &cs.mutatingWebhookConfiguration, MutatingWebhookConfiguration)
}

func (cs *Clients) ValidatingWebhookConfiguration() ValidatingWebhookConfigurationClient {
	return cachedSingleton(cs, &cs.validatingWebhookConfiguration, ValidatingWebhookConfiguration)
}
