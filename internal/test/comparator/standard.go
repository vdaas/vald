package comparator

import (
	"reflect"
	"sync/atomic"

	"github.com/vdaas/vald/internal/errgroup"
)

type (
	atomicValue = atomic.Value
	errorGroup  = errgroup.Group
)

var (
	AtomicValue = func(x, y atomicValue) bool {
		return reflect.DeepEqual(x.Load(), y.Load())
	}

	ErrorGroup = func(x, y errorGroup) bool {
		return reflect.DeepEqual(x, y)
	}

	// channel comparator

	ErrorChannel = func(x, y <-chan error) bool {
		if x == nil && y == nil {
			return true
		}
		if x == nil || y == nil {
			return false
		}

		chanToSlice := func(c <-chan error) []error {
			s := make([]error, 0)
			for v := range c {
				s = append(s, v)
			}
			return s
		}

		s1 := chanToSlice(x)
		s2 := chanToSlice(y)

		return reflect.DeepEqual(s1, s2)
	}
)
