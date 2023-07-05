package sync

import gache "github.com/kpango/gache/v2"

type Map[K comparable, V any] struct {
	gache.Map[K, V]
}
