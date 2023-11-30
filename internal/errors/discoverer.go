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

// Package errors provides error types and function
package errors

var (
	// ErrAddrCouldNotDiscover represents a function to generate an error that address couldn't discover.
	ErrAddrCouldNotDiscover = func(err error, record string) error {
		return Wrapf(err, "addr %s ip couldn't discover", record)
	}

	// ErrNodeNotFound represents a function to generate an error of discover node not found.
	ErrNodeNotFound = func(node string) error {
		return Errorf("discover node %s not found", node)
	}

	// ErrNamespaceNotFound represents a function to generate an error of discover namespace not found.
	ErrNamespaceNotFound = func(ns string) error {
		return Errorf("discover namespace %s not found", ns)
	}

	// ErrPodNameNotFound represents a function to generate an error of discover pod not found.
	ErrPodNameNotFound = func(name string) error {
		return Errorf("discover pod %s not found", name)
	}

	// ErrInvalidDiscoveryCache represents an error that type conversion of discovery cache failed.
	ErrInvalidDiscoveryCache = New("cache type cast failed")
)
