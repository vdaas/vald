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

package client

import (
	"github.com/vdaas/vald/internal/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/restmapper"
)

type (
	// Dynamic is the dynamic client interface for operating on arbitrary
	// resources by GroupVersionResource.
	Dynamic = dynamic.Interface
	// DynamicResource operates on a single dynamic resource.
	DynamicResource = dynamic.ResourceInterface
	// RESTMapper resolves GroupVersionKinds to GroupVersionResources.
	RESTMapper = meta.RESTMapper
)

// RESTScopeNameNamespace identifies namespace-scoped resources in a RESTMapping.
const RESTScopeNameNamespace = meta.RESTScopeNameNamespace

// NewDynamicClient builds a dynamic client and a RESTMapper from the
// ClientSet's rest config so arbitrary GVKs can be resolved to GVRs.
func NewDynamicClient(c ClientSet) (Dynamic, RESTMapper, error) {
	if c == nil {
		return nil, nil, errors.ErrKubernetesClientNotFound
	}
	cfg := c.GetRESTConfig()
	if cfg == nil {
		return nil, nil, errors.ErrKubernetesClientNotFound
	}
	dyn, err := dynamic.NewForConfig(cfg)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to create dynamic client")
	}
	dc, err := discovery.NewDiscoveryClientForConfig(cfg)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to create discovery client")
	}
	groups, err := restmapper.GetAPIGroupResources(dc)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to fetch api group resources")
	}
	return dyn, restmapper.NewDiscoveryRESTMapper(groups), nil
}
