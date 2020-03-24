package ngt

import (
	"io/ioutil"
	"os"

	"github.com/vdaas/vald/hack/benchmark/internal/core"
	"github.com/vdaas/vald/internal/core/ngt"
)

type ObjectType int

const (
	ObjectNone ObjectType = iota
	Uint8
	Float
)

type agent struct {
	idxPath    string
	tmpdir     string
	objectType ObjectType
	dimension  int
	ngt.NGT
}

func New(opts ...Option) (core.Core32, error) {
	a := new(agent)
	for _, opt := range append(defaultOptions, opts...) {
		opt(a)
	}

	tmpdir, err := ioutil.TempDir("", a.idxPath)
	if err != nil {
		return nil, err
	}
	a.tmpdir = tmpdir

	var typ = ngt.ObjectNone
	switch a.objectType {
	case Uint8:
		typ = ngt.Uint8
	case Float:
		typ = ngt.Float
	}

	n, err := ngt.New(
		ngt.WithIndexPath(tmpdir),
		ngt.WithDimension(a.dimension),
		ngt.WithObjectType(typ),
	)
	if err != nil {
		return nil, err
	}
	a.NGT = n

	return a, nil
}

func (a *agent) Search(vec []float32, size int, epsilon, radius float32) (interface{}, error) {
	return a.Search(vec, size, epsilon, radius)
}

func (a *agent) Close() {
	if len(a.tmpdir) != 0 {
		os.RemoveAll(a.tmpdir)
	}
	a.NGT.Close()
}
