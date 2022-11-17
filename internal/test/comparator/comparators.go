//
// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
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
package comparator

import (
	"reflect"
	"sync"

	"github.com/vdaas/vald/internal/errors"
)

var (
	RWMutexComparer = Comparer(func(x, y *sync.RWMutex) bool {
		return reflect.DeepEqual(x, y)
	})

	MutexComparer = Comparer(func(x, y sync.Mutex) bool {
		return reflect.DeepEqual(x, y)
	})

	CondComparer = Comparer(func(x, y *sync.Cond) bool {
		return reflect.DeepEqual(x, y)
	})

	ErrorComparer = Comparer(func(x, y error) bool {
		return errors.Is(x, y)
	})

	OnceComparer = Comparer(func(x, y sync.Once) bool {
		return reflect.DeepEqual(x, y)
	})

	WaitGroupComparer = Comparer(func(x, y sync.WaitGroup) bool {
		return reflect.DeepEqual(x, y)
	})
)
