package gongt

import (
	"github.com/vdaas/vald/internal/core/ngt"
	"github.com/vdaas/vald/internal/errors"

	"github.com/yahoojapan/gongt"
)

// NGT wrapps gongt.NGT object.
type NGT struct {
	indexPath string
	*gongt.NGT
}

var (

	// ErrNotSupportedMethod is not supported error.
	ErrNotSupportedMethod = errors.New("not supported method")
)

// New returns ngt.NGT implementation.
func New(opts ...Option) ngt.NGT {
	n := new(NGT)
	for _, opt := range append(defaultOptions, opts...) {
		opt(n)
	}
	n.NGT = gongt.New(n.indexPath)
	return n
}

func (n *NGT) Search(vec []float32, size int, epsilon, radius float32) ([]ngt.SearchResult, error) {
	results, err := n.NGT.Search(tofloat64(vec), size, float64(epsilon))
	return toResults(results), err
}

func (n *NGT) Insert(vec []float32) (uint, error) {
	id, err := n.NGT.Insert(tofloat64(vec))
	return uint(id), err
}

func (n *NGT) InsertCommit(vec []float32, poolSize uint32) (uint, error) {
	id, err := n.NGT.InsertCommit(tofloat64(vec), int(poolSize))
	return uint(id), err
}

func (n *NGT) BulkInsert(vecs [][]float32) ([]uint, []error) {
	ids, err := n.NGT.BulkInsert(tofloats64(vecs))
	return touint(ids), err
}

func (n *NGT) BulkInsertCommit(vecs [][]float32, poolSize uint32) ([]uint, []error) {
	ids, errs := n.NGT.BulkInsertCommit(tofloats64(vecs), int(poolSize))
	return touint(ids), errs
}

func (n *NGT) CreateAndSaveIndex(poolSize uint32) error {
	return n.NGT.CreateAndSaveIndex(int(poolSize))
}

func (n *NGT) CreateIndex(poolSize uint32) error {
	return n.NGT.CreateIndex(int(poolSize))
}

func (n *NGT) Remove(id uint) error {
	return n.NGT.Remove(int(id))
}

func (n *NGT) BulkRemove(ids ...uint) error {
	return ErrNotSupportedMethod
}

func (n *NGT) GetVector(id uint) ([]float32, error) {
	vecs, err := n.NGT.GetVector(int(id))
	return tofloat32(vecs), err
}

func toResults(in []gongt.SearchResult) (out []ngt.SearchResult) {
	out = make([]ngt.SearchResult, len(in))
	for i := range in {
		out[i] = ngt.SearchResult{
			ID:       uint32(in[i].ID),
			Distance: float32(in[i].Distance),
		}
	}
	return
}

func tofloats64(in [][]float32) (out [][]float64) {
	out = make([][]float64, len(in))
	for i := range in {
		out[i] = tofloat64(in[i])
	}
	return
}

func tofloat64(in []float32) (out []float64) {
	out = make([]float64, len(in))
	for i := range in {
		out[i] = float64(in[i])
	}
	return
}

func tofloats32(in [][]float64) (out [][]float32) {
	out = make([][]float32, len(in))
	for i := range in {
		out[i] = tofloat32(in[i])
	}
	return
}

func tofloat32(in []float64) (out []float32) {
	out = make([]float32, len(in))
	for i := range in {
		out[i] = float32(in[i])
	}
	return
}

func touint(in []int) (out []uint) {
	out = make([]uint, len(in))
	for i := range in {
		out[i] = uint(in[i])
	}
	return
}
