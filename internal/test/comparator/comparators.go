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

	ErrorComparer = Comparer(func(x, y error) bool {
		return errors.Is(x, y)
	})
)
