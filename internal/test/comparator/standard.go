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
package comparator

import (
	"sync/atomic"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/vdaas/vald/internal/errgroup"
)

type (
	atomicValue = atomic.Value
	errorGroup  = errgroup.Group

	Option = cmp.Option
)

var (
	AllowUnexported  = cmp.AllowUnexported
	IgnoreUnexported = cmpopts.IgnoreUnexported
	Comparer         = cmp.Comparer
	Diff             = cmp.Diff
	Equal            = cmp.Equal
	IgnoreTypes      = cmpopts.IgnoreTypes
	IgnoreFields     = cmpopts.IgnoreFields
)

/*
var (
	AtomicValue = func(x, y atomicValue) bool {
		return reflect.DeepEqual(x.Load(), y.Load())
	}

	ErrorGroup = func(x, y errorGroup) bool {
		return reflect.DeepEqual(x, y)
	}

	// channel comparator

		ErrChannel := func(x, y <-chan error) bool {
			if x == nil && y == nil {
				return true
			}
			if x == nil || y == nil || len(x) != len(y) {
				return false
			}

			for e := range x {
				if e1 := <-y; !errors.Is(e, e1) {
					return false
				}
			}
			return true
		}
)
*/
