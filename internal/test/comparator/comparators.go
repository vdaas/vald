package comparator

import (
	"reflect"
	"sync"
)

var (
	RWMutexComparer = Comparer(func(x, y *sync.RWMutex) bool {
		return reflect.DeepEqual(x, y)
	})
)
