// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
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
package testify

import (
	"github.com/stretchr/testify/mock"
)

type (
	Arguments = mock.Arguments
)

<<<<<<< HEAD:internal/test/testify/testify.go
const (
	Anything = mock.Anything
)

var AnythingOfType = mock.AnythingOfType
=======
func SortFunc[E any](x []E, less func(left, right E) bool) {
	slices.SortFunc(x, less)
}

func CompactFunc[S ~[]E, E any](s S, eq func(E, E) bool) S {
	return slices.CompactFunc(s, eq)
}
>>>>>>> feature/gateway-lb/add-search-ratio-for-limited-forwarding-to-agent-and-add-new-sort-algos:internal/slices/slices.go
