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

package client

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
)

type Option func(*client) error

// WithSchemeBuilder registers the given scheme builder's types on the
// client's scheme.
func WithSchemeBuilder(sb runtime.SchemeBuilder) Option {
	return func(c *client) error {
		return sb.AddToScheme(c.scheme)
	}
}

// WithRESTConfig sets the *rest.Config used to construct the underlying
// controller-runtime client and Clientset, so New can be used from
// environments that cannot rely on the KUBECONFIG/in-cluster
// auto-detection (e.g. tests or external tooling holding their own
// rest.Config).
func WithRESTConfig(cfg *rest.Config) Option {
	return func(c *client) error {
		c.restConfig = cfg
		return nil
	}
}

// WithKubeConfigPath sets an explicit kubeconfig file path, overriding the
// KUBECONFIG environment variable / recommended home file auto-detection
// that New falls back to when no RESTConfig is supplied.
func WithKubeConfigPath(path string) Option {
	return func(c *client) error {
		c.kubeConfigPath = path
		return nil
	}
}

// WithKubeContext selects a non-default context from the resolved
// kubeconfig.
func WithKubeContext(name string) Option {
	return func(c *client) error {
		c.kubeContext = name
		return nil
	}
}
