// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package faiss

/*
#cgo LDFLAGS: -lfaiss
#include <Capi.h>
*/
import "C"

import (
	"unsafe"

	"github.com/vdaas/vald/internal/core/algorithm"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/sync"
)

// errIndexClosed is returned when an operation is requested on a closed (freed) index.
var errIndexClosed = errors.NewFaissError("faiss index is closed")

type (
	// Faiss is core interface.
	Faiss interface {
		// SaveIndex stores faiss index to strage.
		SaveIndex() error

		// SaveIndexWithPath stores faiss index to specified storage.
		SaveIndexWithPath(idxPath string) error

		// Train trains faiss index.
		Train(nb int, xb []float32) error

		// Add returns faiss ntotal.
		Add(nb int, xb []float32, xids []int64) (int, error)

		// Search returns search result as []algorithm.SearchResult.
		Search(k, nprobe, nq int, xq []float32) ([]algorithm.SearchResult, error)

		// Remove removes from faiss index.
		Remove(size int, ids []int64) (int, error)

		// Close faiss index.
		Close()
	}

	faiss struct {
		st          *C.FaissStruct
		dimension   C.int
		nlist       C.int
		m           C.int
		nbitsPerIdx C.int
		methodType  methodType
		metricType  metricType
		idxPath     string
		mu          *sync.RWMutex
	}
)

// methodType is alias of method type in Faiss(e.g. IVFPQ, BinaryIndex, ...).
type methodType int

// metricType is alias of metric type in Faiss.
type metricType int

const (
	// -------------------------------------------------------------
	// Method Type Definition
	// -------------------------------------------------------------
	IVFPQ = iota
	BinaryIndex
)

const (
	// -------------------------------------------------------------
	// Metric Type Definition
	// (https://github.com/facebookresearch/faiss/wiki/MetricType-and-distances)
	// -------------------------------------------------------------
	// DistanceNone is unknown distance type.
	DistanceNone metricType = iota - 1
	// InnerProduct is inner product.
	InnerProduct
	// L2 is l2 norm.
	L2
	// -------------------------------------------------------------

	// -------------------------------------------------------------
	// ErrorCode is false
	// -------------------------------------------------------------
	ErrorCode = C._Bool(false)
	// -------------------------------------------------------------
)

// New returns Faiss instance with recreating empty index file.
func New(opts ...Option) (Faiss, error) {
	return gen(false, opts...)
}

func Load(opts ...Option) (Faiss, error) {
	return gen(true, opts...)
}

func gen(isLoad bool, opts ...Option) (Faiss, error) {
	var (
		f   = new(faiss)
		err error
	)
	f.mu = new(sync.RWMutex)

	defer func() {
		if err != nil {
			f.Close()
		}
	}()

	for _, opt := range append(defaultOptions, opts...) {
		if err = opt(f); err != nil {
			return nil, errors.Wrap(err, "faiss option error")
		}
	}

	if isLoad {
		path := C.CString(f.idxPath)
		defer C.free(unsafe.Pointer(path))
		f.st = C.faiss_read_index(path, C.int(f.methodType))
		if f.st == nil {
			return nil, errors.NewFaissError("faiss load index error")
		}
	} else {
		switch f.metricType {
		case InnerProduct:
			f.st = C.faiss_create_index(f.dimension, f.nlist, f.m, f.nbitsPerIdx, C.int(f.methodType), C.int(InnerProduct))
		case L2:
			f.st = C.faiss_create_index(f.dimension, f.nlist, f.m, f.nbitsPerIdx, C.int(f.methodType), C.int(L2))
		default:
			return nil, errors.NewFaissError("faiss create index error: no metric type")
		}
		if f.st == nil {
			return nil, errors.NewFaissError("faiss create index error: nil pointer")
		}
	}

	return f, nil
}

// SaveIndex stores faiss index to storage.
func (f *faiss) SaveIndex() error {
	path := C.CString(f.idxPath)
	defer C.free(unsafe.Pointer(path))

	f.mu.Lock()
	if f.st == nil {
		f.mu.Unlock()
		return errIndexClosed
	}
	ret := C.faiss_write_index(f.st, path, C.int(f.methodType))
	f.mu.Unlock()
	if ret == ErrorCode {
		return errors.NewFaissError("failed to faiss_write_index")
	}

	return nil
}

// SaveIndexWithPath stores faiss index to specified storage.
func (f *faiss) SaveIndexWithPath(idxPath string) error {
	path := C.CString(idxPath)
	defer C.free(unsafe.Pointer(path))

	f.mu.Lock()
	if f.st == nil {
		f.mu.Unlock()
		return errIndexClosed
	}
	ret := C.faiss_write_index(f.st, path, C.int(f.methodType))
	f.mu.Unlock()
	if ret == ErrorCode {
		return errors.NewFaissError("failed to faiss_write_index")
	}

	return nil
}

// Train trains faiss index.
func (f *faiss) Train(nb int, xb []float32) error {
	if nb <= 0 || len(xb) == 0 {
		return errors.NewFaissError("failed to faiss_train: no training vectors")
	}
	if dim := int(f.dimension); len(xb) < dim*nb {
		return errors.ErrIncompatibleDimensionSize(len(xb)/nb, dim)
	}

	f.mu.Lock()
	if f.st == nil {
		f.mu.Unlock()
		return errIndexClosed
	}
	ret := C.faiss_train(f.st, C.int(nb), (*C.float)(&xb[0]), C.int(f.methodType))
	f.mu.Unlock()
	if ret == ErrorCode {
		return errors.NewFaissError("failed to faiss_train")
	}

	return nil
}

// Add returns faiss ntotal.
func (f *faiss) Add(nb int, xb []float32, xids []int64) (int, error) {
	if nb <= 0 || len(xb) == 0 || len(xids) == 0 {
		return -1, errors.NewFaissError("failed to faiss_add: no vectors")
	}
	dim := int(f.dimension)
	if len(xb) != dim*nb || len(xb) != dim*len(xids) {
		return -1, errors.ErrIncompatibleDimensionSize(len(xb)/nb, dim)
	}

	f.mu.Lock()
	if f.st == nil {
		f.mu.Unlock()
		return -1, errIndexClosed
	}
	ntotal := int(C.faiss_add(f.st, C.int(nb), (*C.float)(&xb[0]), (*C.long)(&xids[0]), C.int(f.methodType)))
	f.mu.Unlock()
	if ntotal < 0 {
		return ntotal, errors.NewFaissError("failed to faiss_add")
	}

	return ntotal, nil
}

// Search returns search result as []algorithm.SearchResult.
func (f *faiss) Search(k, nprobe, nq int, xq []float32) ([]algorithm.SearchResult, error) {
	if k <= 0 || nq <= 0 || len(xq) == 0 {
		return nil, errors.NewFaissError("failed to faiss_search: invalid k or no query vectors")
	}
	if len(xq) != nq*int(f.dimension) {
		return nil, errors.ErrIncompatibleDimensionSize(len(xq), int(f.dimension))
	}

	I := make([]int64, k*nq)
	D := make([]float32, k*nq)
	f.mu.RLock()
	if f.st == nil {
		f.mu.RUnlock()
		return nil, errIndexClosed
	}
	ret := C.faiss_search(f.st, C.int(k), C.int(nprobe), C.int(nq), (*C.float)(&xq[0]), (*C.long)(&I[0]), (*C.float)(&D[0]), C.int(f.methodType))
	f.mu.RUnlock()
	if ret == ErrorCode {
		return nil, errors.NewFaissError("failed to faiss_search")
	}

	return toSearchResult(k, I, D), nil
}

// toSearchResult converts the first k labels/distances of a faiss search into
// results. FAISS pads missing neighbours with label -1 at the tail, so the
// result stops at the first negative label instead of reporting uint32(-1).
func toSearchResult(k int, I []int64, D []float32) []algorithm.SearchResult {
	n := min(k, len(I), len(D))
	for i := range n {
		if I[i] < 0 {
			n = i
			break
		}
	}
	result := make([]algorithm.SearchResult, n)
	for i := range result {
		result[i] = algorithm.SearchResult{
			ID:       uint32(I[i]),
			Distance: D[i],
			Error:    nil,
		}
	}
	return result
}

// Remove removes from faiss index.
func (f *faiss) Remove(size int, ids []int64) (int, error) {
	if size <= 0 || len(ids) < size {
		return -1, errors.NewFaissError("failed to faiss_remove: invalid size or ids")
	}

	f.mu.Lock()
	if f.st == nil {
		f.mu.Unlock()
		return -1, errIndexClosed
	}
	ntotal := int(C.faiss_remove(f.st, C.int(size), (*C.long)(&ids[0]), C.int(f.methodType)))
	f.mu.Unlock()
	if ntotal < 0 {
		return ntotal, errors.NewFaissError("failed to faiss_remove")
	}

	return ntotal, nil
}

// Close faiss index.
func (f *faiss) Close() {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.st != nil {
		C.faiss_free(f.st, C.int(f.methodType))
		f.st = nil
	}
}
