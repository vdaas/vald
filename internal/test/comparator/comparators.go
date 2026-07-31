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
package comparator

import (
	"reflect"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/atomic"
)

// deepEqualComparer builds a cmp.Option comparing two T values with
// reflect.DeepEqual. Test comparators for synchronization primitives
// (Mutex, Once, WaitGroup, atomic values) deliberately compare snapshots of
// the internal state by value; the copy is confined to the comparison and
// never handed back, so the usual copylocks hazard (mutating a copied lock)
// does not apply.
func deepEqualComparer[T any]() Option {
	return Comparer(func(x, y T) bool {
		return reflect.DeepEqual(x, y)
	})
}

// The comparers are package-level shared cmp.Options for tests; they are
// immutable after init, so the accidental mutation gochecknoglobals guards
// against cannot occur.
//
//nolint:gochecknoglobals
var (
	RWMutexComparer = deepEqualComparer[*sync.RWMutex]()

	MutexComparer = deepEqualComparer[sync.Mutex]()

	CondComparer = deepEqualComparer[*sync.Cond]()

	ErrorComparer = Comparer(func(x, y error) bool {
		return errors.Is(x, y)
	})

	OnceComparer = deepEqualComparer[sync.Once]()

	WaitGroupComparer = deepEqualComparer[sync.WaitGroup]()

	AtomicUint64Comparator = deepEqualComparer[atomic.Uint64]()
)
