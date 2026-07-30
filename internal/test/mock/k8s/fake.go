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

package k8s

import (
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/client/interceptor"
)

// InterceptorFuncs is an alias for controller-runtime's interceptor.Funcs,
// letting tests hook individual fake-client operations without importing
// sigs.k8s.io directly.
type InterceptorFuncs = interceptor.Funcs

// Aliases re-exporting controller-runtime's fake-client constructors so that
// tests outside internal/ can build fake clients without importing
// sigs.k8s.io directly.
//
//nolint:gochecknoglobals // immutable function aliases confining sigs.k8s.io imports to this package
var (
	NewFakeClientBuilder = fake.NewClientBuilder
	NewInterceptorClient = interceptor.NewClient
)
