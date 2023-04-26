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

import "sigs.k8s.io/controller-runtime/pkg/scheme"

type Option func(*client) error

func WithSchemeBuilder(sb scheme.Builder) Option {
	return func(c *client) error {
		sb.AddToScheme(c.scheme)
		return nil
	}
}
