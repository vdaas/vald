// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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
package comparator

import (
	"reflect"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/sync"
)

var (
	RWMutexComparer = Comparer(func(x, y *sync.RWMutex) bool {
		return reflect.DeepEqual(x, y)
	})

	// skipcq: VET-V0008
	MutexComparer = Comparer(func(x, y sync.Mutex) bool {
		// skipcq: VET-V0008
		return reflect.DeepEqual(x, y)
	})

	CondComparer = Comparer(func(x, y *sync.Cond) bool {
		return reflect.DeepEqual(x, y)
	})

	ErrorComparer = Comparer(func(x, y error) bool {
		return errors.Is(x, y)
	})

	// skipcq: VET-V0008
	OnceComparer = Comparer(func(x, y sync.Once) bool {
		// skipcq: VET-V0008
		return reflect.DeepEqual(x, y)
	})

	// skipcq: VET-V0008
	WaitGroupComparer = Comparer(func(x, y sync.WaitGroup) bool {
		// skipcq: VET-V0008
		return reflect.DeepEqual(x, y)
	})
)
