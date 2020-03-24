package gongt

import (
	"errors"
	"io/ioutil"
	"os"

	"github.com/vdaas/vald/hack/benchmark/internal/core"
	"github.com/yahoojapan/gongt"
)

type (
	SearchResult = gongt.SearchResult
	ObjectType   = gongt.ObjectType
)

const (
	ObjectNone ObjectType = iota
	Uint8
	Float
)

var (
	ErrNotSupportedMethod = errors.New("not supported method")
)

type ngt struct {
	indexPath  string
	tmpdir     string
	objectType ObjectType
	dimension  int
	*gongt.NGT
}

func New(opts ...Option) (core.Core64, error) {
	n := new(ngt)
	for _, opt := range append(defaultOptions, opts...) {
		opt(n)
	}

	tmpdir, err := ioutil.TempDir("", n.indexPath)
	if err != nil {
		return nil, err
	}
	n.tmpdir = tmpdir

	n.NGT = gongt.New(tmpdir).
		SetObjectType(n.objectType).
		SetDimension(n.dimension).
		Open()

	return n, nil
}

func (n *ngt) Search(vec []float64, size int, epsilon, radius float64) (interface{}, error) {
	return n.NGT.Search(vec, size, epsilon)
}

func (n *ngt) Insert(vec []float64) (uint, error) {
	id, err := n.NGT.Insert(vec)
	return uint(id), err
}

func (n *ngt) InsertCommit(vec []float64, poolSize uint32) (uint, error) {
	id, err := n.NGT.Insert(vec)
	return uint(id), err
}

func (n *ngt) BulkInsert(vecs [][]float64) ([]uint, []error) {
	ids, errs := n.NGT.BulkInsert(vecs)
	return toUint(ids), errs
}

func (n *ngt) BulkInsertCommit(vecs [][]float64, poolSize uint32) ([]uint, []error) {
	ids, errs := n.NGT.BulkInsertCommit(vecs, int(poolSize))
	return toUint(ids), errs
}

func (n *ngt) CreateAndSaveIndex(poolSize uint32) error {
	return n.NGT.CreateAndSaveIndex(int(poolSize))
}

func (n *ngt) CreateIndex(poolSize uint32) error {
	return n.NGT.CreateIndex(int(poolSize))
}

func (n *ngt) Remove(id uint) error {
	return n.NGT.StrictRemove(id)
}

func (n *ngt) BulkRemove(ids ...uint) error {
	return ErrNotSupportedMethod
}

func (n *ngt) Close() {
	if len(n.indexPath) != 0 {
		os.RemoveAll(n.tmpdir)
	}
	n.NGT.Close()
}

func toUint(in []int) (out []uint) {
	out = make([]uint, len(in))
	for i := range in {
		out[i] = uint(in[i])
	}
	return
}
