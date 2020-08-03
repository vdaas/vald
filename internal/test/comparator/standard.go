package comparator

import (
	"sync/atomic"

	"github.com/google/go-cmp/cmp"
	"github.com/vdaas/vald/internal/errgroup"
)

type (
	atomicValue = atomic.Value
	errorGroup  = errgroup.Group
	Option      = cmp.Option
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
