package gongt

import (
	"errors"
	"io/ioutil"
	"os"

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

type NGT interface {
	Search(vec []float64, size int, epsilon float64) ([]SearchResult, error)
	Insert(vec []float64) (int, error)
	InsertCommit(vec []float64, poolSize int) (int, error)
	BulkInsert(vecs [][]float64) ([]int, []error)
	BulkInsertCommit(vecs [][]float64, poolSize int) ([]int, []error)
	CreateAndSaveIndex(poolSize int) error
	CreateIndex(poolSize int) error
	SaveIndex() error
	Remove(id int) error
	BulkRemove(ids ...uint) error
	GetVector(id int) ([]float64, error)
	Close()
}

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

func New(opts ...Option) (NGT, error) {
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

func (n *ngt) BulkRemove(ids ...uint) error {
	return ErrNotSupportedMethod
}

func (n *ngt) Close() {
	if len(n.indexPath) != 0 {
		os.RemoveAll(n.tmpdir)
	}
	n.NGT.Close()
}
