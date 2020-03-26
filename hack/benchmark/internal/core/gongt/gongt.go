package gongt

import (
	"io/ioutil"
	"os"

	icore "github.com/vdaas/vald/hack/benchmark/internal/core"
	"github.com/vdaas/vald/hack/benchmark/internal/errors"
	"github.com/yahoojapan/gongt"
)

type (
	ObjectType = gongt.ObjectType
)

const (
	ObjectNone ObjectType = iota
	Uint8
	Float
)

type core struct {
	indexPath  string
	tmpdir     string
	objectType ObjectType
	dimension  int
	*gongt.NGT
}

func New(opts ...Option) (icore.Core64, error) {
	c := new(core)
	for _, opt := range append(defaultOptions, opts...) {
		opt(c)
	}

	tmpdir, err := ioutil.TempDir("", c.indexPath)
	if err != nil {
		return nil, err
	}
	c.tmpdir = tmpdir

	c.NGT = gongt.New(tmpdir).
		SetObjectType(c.objectType).
		SetDimension(c.dimension).
		Open()

	return c, nil
}

func (c *core) Search(vec []float64, size int, epsilon, radius float32) (interface{}, error) {
	return c.NGT.Search(vec, size, float64(epsilon))
}

func (c *core) Insert(vec []float64) (uint, error) {
	id, err := c.NGT.Insert(vec)
	return uint(id), err
}

func (c *core) InsertCommit(vec []float64, poolSize uint32) (uint, error) {
	id, err := c.NGT.Insert(vec)
	return uint(id), err
}

func (c *core) BulkInsert(vecs [][]float64) ([]uint, []error) {
	ids, errs := c.NGT.BulkInsert(vecs)
	return toUint(ids), errs
}

func (c *core) BulkInsertCommit(vecs [][]float64, poolSize uint32) ([]uint, []error) {
	ids, errs := c.NGT.BulkInsertCommit(vecs, int(poolSize))
	return toUint(ids), errs
}

func (c *core) CreateAndSaveIndex(poolSize uint32) error {
	return c.NGT.CreateAndSaveIndex(int(poolSize))
}

func (c *core) CreateIndex(poolSize uint32) error {
	return c.NGT.CreateIndex(int(poolSize))
}

func (c *core) Remove(id uint) error {
	return c.NGT.StrictRemove(id)
}

func (c *core) BulkRemove(ids ...uint) error {
	return errors.ErrGoNGTNotSupportedMethod
}

func (c *core) GetVector(id uint) ([]float64, error) {
	return c.NGT.GetVector(int(id))
}

func (c *core) Close() {
	if len(c.indexPath) != 0 {
		os.RemoveAll(c.tmpdir)
	}
	c.NGT.Close()
}

func toUint(in []int) (out []uint) {
	out = make([]uint, len(in))
	for i := range in {
		out[i] = uint(in[i])
	}
	return
}
