package ngt

import (
	"io/ioutil"
	"os"

	icore "github.com/vdaas/vald/hack/benchmark/internal/core"
	"github.com/vdaas/vald/internal/core/ngt"
)

type ObjectType int

const (
	ObjectNone ObjectType = iota
	Uint8
	Float
)

type core struct {
	idxPath    string
	tmpdir     string
	objectType ObjectType
	dimension  int
	ngt.NGT
}

func New(opts ...Option) (icore.Core32, error) {
	c := new(core)
	for _, opt := range append(defaultOptions, opts...) {
		opt(c)
	}

	tmpdir, err := ioutil.TempDir("", c.idxPath)
	if err != nil {
		return nil, err
	}
	c.tmpdir = tmpdir

	var typ = ngt.ObjectNone
	switch c.objectType {
	case Uint8:
		typ = ngt.Uint8
	case Float:
		typ = ngt.Float
	}

	n, err := ngt.New(
		ngt.WithIndexPath(tmpdir),
		ngt.WithDimension(c.dimension),
		ngt.WithObjectType(typ),
	)
	if err != nil {
		return nil, err
	}
	c.NGT = n

	return c, nil
}

func (c *core) Search(vec []float32, size int, epsilon, radius float32) (interface{}, error) {
	return c.NGT.Search(vec, size, epsilon, radius)
}

func (c *core) Close() {
	if len(c.tmpdir) != 0 {
		os.RemoveAll(c.tmpdir)
	}
	// c.NGT.Close()
}
