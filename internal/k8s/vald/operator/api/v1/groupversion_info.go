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

// Package v1 contains API Schema definitions for the controller v1 API group.
package v1

import (
	"github.com/vdaas/vald/internal/k8s/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	GroupVersion = schema.GroupVersion{Group: "vald.vdaas.org", Version: "v1"}

	// SchemeBuilder registers the item type via AddKnownTypes and the generic
	// list alias via AddListToScheme: generic instantiations carry mangled
	// reflect type names, so the list kind must be registered explicitly.
	SchemeBuilder = runtime.NewSchemeBuilder(func(s *runtime.Scheme) error {
		s.AddKnownTypes(GroupVersion, &ValdOperatorRelease{})
		resource.AddListToScheme[ValdOperatorRelease](s, GroupVersion, "ValdOperatorReleaseList")
		metav1.AddToGroupVersion(s, GroupVersion)
		return nil
	})

	AddToScheme = SchemeBuilder.AddToScheme
)
