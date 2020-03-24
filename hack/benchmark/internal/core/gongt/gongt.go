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

type agent struct {
	indexPath  string
	tmpdir     string
	objectType ObjectType
	dimension  int
	*gongt.NGT
}

func New(opts ...Option) (core.Core, error) {
	a := new(agent)
	for _, opt := range append(defaultOptions, opts...) {
		opt(a)
	}

	tmpdir, err := ioutil.TempDir("", a.indexPath)
	if err != nil {
		return nil, err
	}
	a.tmpdir = tmpdir

	a.NGT = gongt.New(tmpdir).
		SetObjectType(a.objectType).
		SetDimension(a.dimension).
		Open()

	return a, nil
}

func (a *agent) Search(vec, size, epsilon, radius interface{}) (interface{}, error) {
	return a.NGT.Search(vec.([]float64), size.(int), epsilon.(float64))
}

func (a *agent) Insert(vec interface{}) (interface{}, error) {
	return a.NGT.Insert(vec.([]float64))
}

func (a *agent) InsertCommit(vec, poolSize interface{}) (interface{}, error) {
	return a.NGT.InsertCommit(vec.([]float64), poolSize.(int))
}

func (a *agent) BulkInsert(vecs interface{}) (interface{}, []error) {
	return a.NGT.BulkInsert(vecs.([][]float64))
}

func (a *agent) BulkInsertCommit(vecs, poolSize interface{}) (interface{}, []error) {
	return a.NGT.BulkInsertCommit(vecs.([][]float64), poolSize.(int))
}

func (a *agent) CreateAndSaveIndex(poolSize interface{}) error {
	return a.NGT.CreateAndSaveIndex(poolSize.(int))
}

func (a *agent) CreateIndex(poolSize interface{}) error {
	return a.NGT.CreateIndex(poolSize.(int))
}

func (a *agent) SaveIndex() error {
	return a.NGT.SaveIndex()
}

func (a *agent) Remove(id interface{}) error {
	return a.NGT.Remove(id.(int))
}

func (a *agent) BulkRemove(ids interface{}) error {
	return ErrNotSupportedMethod
}

func (a *agent) GetVector(id interface{}) (interface{}, error) {
	return a.NGT.GetVector(id.(int))
}

func (a *agent) Close() {
	if len(a.indexPath) != 0 {
		os.RemoveAll(a.tmpdir)
	}
	a.NGT.Close()
}
