//
// Copyright (C) 2019-2020 Vdaas.org Vald team ( kpango, rinx, kmrmt )
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

package errors

var (
	ErrAddrCouldNotDiscover = func(err error, record string) error {
		return Wrapf(err, "addr %s ip couldn't discover", record)
	}

	ErrNodeNotFound = func(node string) error {
		return Errorf("discover node %s not found", node)
	}

	ErrNamespaceNotFound = func(ns string) error {
		return Errorf("discover namespace %s not found", ns)
	}

	ErrPodNameNotFound = func(name string) error {
		return Errorf("discover pod %s not found", name)
	}

	ErrInvalidDiscoveryCache = New("cache type cast failed")
)
