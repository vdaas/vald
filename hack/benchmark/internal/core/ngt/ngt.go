package ngt

import (
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

func New(opts ...Option) (core.Core, error) {
	a := new(agent)
	for _, opt := range append(defaultOptions, opts...) {
		opt(a)
	}
	return a, nil
}

func (a *agent) InsertCommit(vec, poolSize interface{}) (interface{}, error) {
	return a.NGT.InsertCommit(vec.([]float32), poolSize.(uint32))
}

func (a *agent) BulkInsert(vecs interface{}) (interface{}, []error) {
	return a.NGT.BulkInsert(vecs.([][]float32))
}

func (a *agent) BulkInsertCommit(vecs, poolSize interface{}) (interface{}, []error) {
	return a.NGT.BulkInsertCommit(vecs.([][]float32), poolSize.(uint32))
}

func (a *agent) CreateAndSaveIndex(poolSize interface{}) error {
	return a.NGT.CreateAndSaveIndex(poolSize.(uint32))
}

func (a *agent) CreateIndex(poolSize interface{}) error {
	return a.NGT.CreateIndex(poolSize.(uint32))
}

func (a *agent) SaveIndex() error {
	return a.NGT.SaveIndex()
}

func (a *agent) Remove(id interface{}) error {
	return a.NGT.Remove(id.(uint))
}

func (a *agent) BulkRemove(ids interface{}) error {
	return a.NGT.BulkRemove(ids.([]uint)...)
}

func (a *agent) GetVector(id interface{}) (interface{}, error) {
	return a.NGT.GetVector(id.(uint))
}

func (a *agent) Close() {
	if len(a.tmpdir) != 0 {
		os.RemoveAll(a.tmpdir)
	}
	a.NGT.Close()
}
