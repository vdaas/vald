//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
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

// Package errors provides error types and function
package errors

var (
	ErrInvalidReconcilerConfig = New("invalid reconciler config")

	ErrPodIsNotRunning = func(namespace, name string) error {
		return Errorf("pod %s/%s is not running", namespace, name)
	}

	ErrPortForwardAddressNotFound = New("port forward address not found")

	ErrPortForwardPortPairNotFound = New("port forward port pair not found")

	ErrKubernetesClientNotFound = New("kubernetes client not found")

	ErrStatusPatternNeverMatched = New("status pattern never matched")

	ErrUnsupportedKubernetesResourceType = func(obj any) error {
		return Errorf("unsupported kubernetes resource type %T", obj)
	}

	ErrPodTemplateNotFound = New("pod template not found")

	ErrNoAvailablePods = New("no available pods")

	ErrUndefinedNamespace = New("Undefined namespace")

	ErrUndefinedService = New("Undefined service")
)
