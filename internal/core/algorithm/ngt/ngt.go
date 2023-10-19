//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

// Package ngt provides implementation of Go API for https://github.com/yahoojapan/NGT
package ngt

/*
#cgo LDFLAGS: -lngt
#include <NGT/Capi.h>
#include <stdlib.h>
*/
import "C"

import (
	"context"
	"reflect"
	"runtime"
	"sync/atomic"
	"unsafe"

	"github.com/vdaas/vald/internal/core/algorithm"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/file"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/sync"
)

type (

	// NGT is core interface.
	NGT interface {
		// Search returns search result as []SearchResult
		Search(ctx context.Context, vec []float32, size int, epsilon, radius float32) ([]SearchResult, error)

		// Linear Search returns linear search result as []SearchResult
		LinearSearch(ctx context.Context, vec []float32, size int) ([]SearchResult, error)

		// Insert returns NGT object id.
		// This only stores not indexing, you must call CreateIndex and SaveIndex.
		Insert(vec []float32) (uint, error)

		// InsertCommit returns NGT object id.
		// This stores and indexes at the same time.
		InsertCommit(vec []float32, poolSize uint32) (uint, error)

		// BulkInsert returns NGT object ids.
		// This only stores not indexing, you must call CreateIndex and SaveIndex.
		BulkInsert(vecs [][]float32) ([]uint, []error)

		// BulkInsertCommit returns NGT object ids.
		// This stores and indexes at the same time.
		BulkInsertCommit(vecs [][]float32, poolSize uint32) ([]uint, []error)

		// CreateAndSaveIndex call  CreateIndex and SaveIndex in a row.
		CreateAndSaveIndex(poolSize uint32) error

		// CreateIndex creates NGT index.
		CreateIndex(poolSize uint32) error

		// SaveIndex stores NGT index to storage.
		SaveIndex() error

		// SaveIndexWithPath stores NGT index to specified storage.
		SaveIndexWithPath(path string) error

		// Remove removes from NGT index.
		Remove(id uint) error

		// BulkRemove removes multiple NGT index
		BulkRemove(ids ...uint) error

		// GetVector returns vector stored in NGT index.
		GetVector(id uint) ([]float32, error)

		// Close NGT index.
		Close()
	}

	ngt struct {
		idxPath             string
		inMemory            bool
		bulkInsertChunkSize int
		dimension           C.int32_t
		objectType          objectType
		radius              float32
		epsilon             float32
		poolSize            uint32
		cnt                 atomic.Uint64
		prop                C.NGTProperty
		epool               sync.Pool     // NGT error buffer pool
		eps                 atomic.Uint64 // NGT error buffer pool size
		epl                 uint64        // NGT error buffer pool size limit
		index               C.NGTIndex
		ospace              C.NGTObjectSpace
		mu                  *sync.RWMutex
		cmu                 *sync.RWMutex
	}

	ngtError struct {
		err       C.NGTError
		destroyed atomic.Bool
	}
)

func newNGTError() (n *ngtError) {
	n = &ngtError{
		err: C.ngt_create_error_object(),
	}
	n.destroyed.Store(false)
	runtime.SetFinalizer(n, func(ne *ngtError) {
		ne.close()
	})
	return n
}

func (n *ngtError) close() {
	if !n.destroyed.Load() {
		C.ngt_destroy_error_object(n.err)
		n.destroyed.Store(true)
	}
}

// ObjectType is alias of object type in NGT.
type objectType int

// DistanceType is alias of distance type in NGT.
type distanceType int

const (
	// -------------------------------------------------------------
	// Object Type Definition
	// -------------------------------------------------------------
	// ObjectNone is unknown object type.
	ObjectNone objectType = iota
	// Uint8 is 8bit unsigned integer.
	Uint8
	// Float is 32bit floating point number.
	Float
	// HalfFloat is 16bit floating point number.
	HalfFloat
	// -------------------------------------------------------------.

	// -------------------------------------------------------------
	// Distance Type Definition
	// -------------------------------------------------------------
	// DistanceNone is unknown distance type.
	DistanceNone distanceType = iota - 1
	// L1 is l1 norm.
	L1
	// L2 is l2 norm.
	L2
	// Angle is angle distance.
	Angle
	// Hamming is hamming distance.
	Hamming
	// Cosine is cosine distance.
	Cosine
	// Poincare is poincare distance.
	Poincare
	// Lorentz is lorenz distance.
	Lorentz
	// Jaccard is jaccard distance.
	Jaccard
	// SparseJaccard is sparse jaccard distance.
	SparseJaccard
	// NormalizedL2 is l2 distance with normalization.
	NormalizedL2
	// NormalizedAngle is angle distance with normalization.
	NormalizedAngle
	// NormalizedCosine is cosine distance with normalization.
	NormalizedCosine

	// -------------------------------------------------------------.

	// -------------------------------------------------------------
	// ErrorCode is false
	// -------------------------------------------------------------.
	ErrorCode = C._Bool(false)
	// -------------------------------------------------------------.
)

func (o objectType) String() string {
	switch o {
	case Uint8:
		return "Uint8"
	case HalfFloat:
		return "HalfFloat"
	case Float:
		return "Float"
	}
	return "Unknown"
}

func (d distanceType) String() string {
	switch d {
	case L1:
		return "L1"
	case L2:
		return "L2"
	case Angle:
		return "Angle"
	case Hamming:
		return "Hamming"
	case Cosine:
		return "Cosine"
	case Poincare:
		return "Poincare"
	case Lorentz:
		return "Lorentz"
	case Jaccard:
		return "Jaccard"
	case SparseJaccard:
		return "SparseJaccard"
	case NormalizedL2:
		return "NormalizedL2"
	case NormalizedAngle:
		return "NormalizedAngle"
	case NormalizedCosine:
		return "NormalizedCosine"
	}
	return "Unknown"
}

// New returns NGT instance with recreating empty index file.
func New(opts ...Option) (NGT, error) {
	return gen(false, opts...)
}

// Load returns NGT instance from existing index file.
func Load(opts ...Option) (NGT, error) {
	return gen(true, opts...)
}

func gen(isLoad bool, opts ...Option) (NGT, error) {
	var (
		n   = new(ngt)
		err error
	)
	n.mu = new(sync.RWMutex)
	n.cmu = new(sync.RWMutex)

	defer func() {
		if err != nil {
			n.Close()
		}
	}()

	err = n.setup()
	if err != nil {
		log.Warnf("failed to setup ngt core index\terr: %v", err)
		return nil, err
	}

	defer C.ngt_destroy_property(n.prop)

	err = n.loadOptions(opts...)
	if err != nil {
		log.Warnf("failed to load ngt core options\terr: %v", err)
		return nil, err
	}

	if isLoad {
		err = n.open()
		if err != nil {
			log.Warnf("failed to load ngt core index\terr: %v", err)
			return nil, err
		}
	} else {
		err = n.create()
		if err != nil {
			log.Warnf("failed to create new ngt core index\terr: %v", err)
			return nil, err
		}
	}

	err = n.loadObjectSpace()
	if err != nil {
		log.Warnf("failed to load ngt object space\terr: %v", err)
		return nil, err
	}

	return n, nil
}

func (n *ngt) setup() error {
	n.epool = sync.Pool{
		New: func() interface{} {
			return newNGTError()
		},
	}

	for i := uint64(0); i < n.epl; i++ {
		n.PutErrorBuffer(newNGTError())
	}

	ne := n.GetErrorBuffer()
	n.prop = C.ngt_create_property(ne.err)
	if n.prop == nil {
		return errors.ErrCreateProperty(n.newGoError(ne))
	}
	n.PutErrorBuffer(ne)
	return nil
}

func (n *ngt) loadOptions(opts ...Option) (err error) {
	for _, opt := range append(defaultOptions, opts...) {
		err = opt(n)
		if err != nil {
			werr := errors.ErrOptionFailed(err, reflect.ValueOf(opt))
			e := new(errors.ErrCriticalOption)
			if errors.As(err, &e) {
				log.Error(werr)
				return werr
			}
			ue := new(errors.ErrIgnoredOption)
			if errors.As(err, &ue) {
				log.Debug(werr)
			} else {
				log.Warn(werr)
			}
		}
	}
	return nil
}

func (n *ngt) create() (err error) {
	path := C.CString(n.idxPath)
	defer C.free(unsafe.Pointer(path))

	ne := n.GetErrorBuffer()
	if !n.inMemory {
		n.index = C.ngt_create_graph_and_tree(path, n.prop, ne.err)
		if n.index == nil {
			return n.newGoError(ne)
		}
		if C.ngt_save_index(n.index, path, ne.err) == ErrorCode {
			return n.newGoError(ne)
		}
	} else {
		n.index = C.ngt_create_graph_and_tree_in_memory(n.prop, ne.err)
		if n.index == nil {
			return n.newGoError(ne)
		}
	}
	n.PutErrorBuffer(ne)

	return nil
}

func (n *ngt) open() error {
	if !file.Exists(n.idxPath) {
		return errors.ErrIndexFileNotFound
	}

	path := C.CString(n.idxPath)
	defer C.free(unsafe.Pointer(path))

	ne := n.GetErrorBuffer()
	n.index = C.ngt_open_index(path, ne.err)
	if n.index == nil {
		return n.newGoError(ne)
	}

	if C.ngt_get_property(n.index, n.prop, ne.err) == ErrorCode {
		return n.newGoError(ne)
	}

	n.dimension = C.ngt_get_property_dimension(n.prop, ne.err)
	if int(n.dimension) == -1 {
		return n.newGoError(ne)
	}
	n.PutErrorBuffer(ne)
	return nil
}

func (n *ngt) loadObjectSpace() error {
	ne := n.GetErrorBuffer()
	n.ospace = C.ngt_get_object_space(n.index, ne.err)
	if n.ospace == nil {
		return n.newGoError(ne)
	}
	n.PutErrorBuffer(ne)
	return nil
}

// Search returns search result as []SearchResult.
func (n *ngt) Search(ctx context.Context, vec []float32, size int, epsilon, radius float32) (result []SearchResult, err error) {
	if len(vec) != int(n.dimension) {
		return nil, errors.ErrIncompatibleDimensionSize(len(vec), int(n.dimension))
	}

	ne := n.GetErrorBuffer()
	results := C.ngt_create_empty_results(ne.err)
	defer C.ngt_destroy_results(results)
	if results == nil {
		return nil, n.newGoError(ne)
	}

	if epsilon == 0 {
		epsilon = n.epsilon
	}

	if radius == 0 {
		radius = n.radius
	}

	n.rLock(true)
	ret := C.ngt_search_index_as_float(
		n.index,
		(*C.float)(&vec[0]),
		n.dimension,
		*(*C.size_t)(unsafe.Pointer(&size)),
		*(*C.float)(unsafe.Pointer(&epsilon)),
		*(*C.float)(unsafe.Pointer(&radius)),
		results,
		ne.err)
	vec = nil
	if ret == ErrorCode {
		n.rUnlock(true)
		return nil, n.newGoError(ne)
	}
	n.rUnlock(true)

	rsize := int(C.ngt_get_result_size(results, ne.err))
	if rsize <= 0 {
		if n.cnt.Load() == 0 {
			n.PutErrorBuffer(ne)
			return nil, errors.ErrSearchResultEmptyButNoDataStored
		}
		err = n.newGoError(ne)
		if err != nil {
			return nil, err
		}
		return nil, errors.ErrEmptySearchResult
	}
	result = make([]SearchResult, rsize)

	for i := range result {
		select {
		case <-ctx.Done():
			n.PutErrorBuffer(ne)
			return result[:i], nil
		default:
		}
		d := C.ngt_get_result(results, C.uint32_t(i), ne.err)
		if d.id == 0 && d.distance == 0 {
			result[i] = SearchResult{0, 0, n.newGoError(ne)}
			ne = n.GetErrorBuffer()
		} else {
			result[i] = SearchResult{uint32(d.id), float32(d.distance), nil}
		}
	}
	n.PutErrorBuffer(ne)

	return result, nil
}

// Linear Search returns linear search result as []SearchResult.
func (n *ngt) LinearSearch(ctx context.Context, vec []float32, size int) (result []SearchResult, err error) {
	if len(vec) != int(n.dimension) {
		return nil, errors.ErrIncompatibleDimensionSize(len(vec), int(n.dimension))
	}

	ne := n.GetErrorBuffer()
	results := C.ngt_create_empty_results(ne.err)
	defer C.ngt_destroy_results(results)
	if results == nil {
		return nil, n.newGoError(ne)
	}

	n.rLock(true)
	ret := C.ngt_linear_search_index_as_float(
		n.index,
		(*C.float)(&vec[0]),
		n.dimension,
		// C.size_t(size),
		*(*C.size_t)(unsafe.Pointer(&size)),
		results,
		ne.err)
	vec = nil

	if ret == ErrorCode {
		n.rUnlock(true)
		return nil, n.newGoError(ne)
	}
	n.rUnlock(true)

	rsize := int(C.ngt_get_result_size(results, ne.err))
	if rsize <= 0 {
		if n.cnt.Load() == 0 {
			n.PutErrorBuffer(ne)
			return nil, errors.ErrSearchResultEmptyButNoDataStored
		}
		err = n.newGoError(ne)
		if err != nil {
			return nil, err
		}
		return nil, errors.ErrEmptySearchResult
	}
	result = make([]SearchResult, rsize)
	for i := range result {
		select {
		case <-ctx.Done():
			n.PutErrorBuffer(ne)
			return result[:i], nil
		default:
		}
		d := C.ngt_get_result(results, C.uint32_t(i), ne.err)
		if d.id == 0 && d.distance == 0 {
			result[i] = SearchResult{0, 0, n.newGoError(ne)}
			ne = n.GetErrorBuffer()
		} else {
			result[i] = SearchResult{uint32(d.id), float32(d.distance), nil}
		}
	}
	n.PutErrorBuffer(ne)

	return result, nil
}

// Insert returns NGT object id.
// This only stores not indexing, you must call CreateIndex and SaveIndex.
func (n *ngt) Insert(vec []float32) (id uint, err error) {
	if len(vec) != int(n.dimension) {
		return 0, errors.ErrIncompatibleDimensionSize(len(vec), int(n.dimension))
	}
	dim := C.uint32_t(n.dimension)
	cvec := (*C.float)(&vec[0])
	ne := n.GetErrorBuffer()
	n.lock(true)
	oid := C.ngt_insert_index_as_float(n.index, cvec, dim, ne.err)
	n.unlock(true)
	id = uint(oid)
	cvec = nil
	vec = vec[:0:0]
	vec = nil
	if id == 0 {
		return 0, n.newGoError(ne)
	}
	n.PutErrorBuffer(ne)
	n.cnt.Add(1)

	return id, nil
}

// InsertCommit returns NGT object id.
// This stores and indexes at the same time.
func (n *ngt) InsertCommit(vec []float32, poolSize uint32) (uint, error) {
	id, err := n.Insert(vec)
	if err != nil {
		return id, err
	}

	err = n.CreateIndex(poolSize)
	if err != nil {
		return id, err
	}

	err = n.SaveIndex()
	if err != nil {
		return id, err
	}
	return id, nil
}

// BulkInsert returns NGT object ids.
// This only stores not indexing, you must call CreateIndex and SaveIndex.
func (n *ngt) BulkInsert(vecs [][]float32) ([]uint, []error) {
	ids := make([]uint, 0, len(vecs))
	errs := make([]error, 0, len(vecs))

	log.Infof("started to bulk insert %d of vectors", len(vecs))
	for i, vec := range vecs {
		id, err := n.Insert(vec)
		if err != nil {
			errs = append(errs, errors.Wrapf(err, "bulkinsert error detected index number: %d,\tid: %d", i, id))
		} else {
			ids = append(ids, id)
		}
	}

	return ids, errs
}

// BulkInsertCommit returns NGT object ids.
// This stores and indexes at the same time.
func (n *ngt) BulkInsertCommit(vecs [][]float32, poolSize uint32) ([]uint, []error) {
	ids := make([]uint, 0, len(vecs))
	errs := make([]error, 0, len(vecs))

	idx := 0
	var id uint
	var err error

	for i, vec := range vecs {
		if id, err = n.Insert(vec); err == nil {
			ids = append(ids, id)
			idx++
			if idx >= n.bulkInsertChunkSize {
				err = n.CreateAndSaveIndex(poolSize)
				if err != nil {
					errs = append(errs, err)
				}
				idx = 0
			}
		} else {
			errs = append(errs, errors.Wrapf(err, "bulkinsert error detected index number: %d,\tid: %d", i, id))
		}
	}

	if idx > 0 {
		err = n.CreateAndSaveIndex(poolSize)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return ids, errs
}

// CreateAndSaveIndex call  CreateIndex and SaveIndex in a row.
func (n *ngt) CreateAndSaveIndex(poolSize uint32) error {
	err := n.CreateIndex(poolSize)
	if err != nil {
		return err
	}
	return n.SaveIndex()
}

// CreateIndex creates NGT index.
func (n *ngt) CreateIndex(poolSize uint32) error {
	if poolSize == 0 {
		poolSize = n.poolSize
	}
	ne := n.GetErrorBuffer()
	n.lock(true)
	ret := C.ngt_create_index(n.index, C.uint32_t(poolSize), ne.err)
	n.unlock(true)
	if ret == ErrorCode {
		return n.newGoError(ne)
	}
	n.PutErrorBuffer(ne)

	return nil
}

// SaveIndex stores NGT index to storage.
func (n *ngt) SaveIndex() error {
	if !n.inMemory {
		path := C.CString(n.idxPath)
		defer C.free(unsafe.Pointer(path))
		ne := n.GetErrorBuffer()
		n.rLock(true)
		ret := C.ngt_save_index(n.index, path, ne.err)
		n.rUnlock(true)
		if ret == ErrorCode {
			return n.newGoError(ne)
		}
		n.PutErrorBuffer(ne)
	}

	return nil
}

// SaveIndexWithPath stores NGT index to specified storage.
func (n *ngt) SaveIndexWithPath(idxPath string) error {
	if !n.inMemory && len(idxPath) != 0 {
		path := C.CString(idxPath)
		defer C.free(unsafe.Pointer(path))
		ne := n.GetErrorBuffer()
		n.rLock(true)
		ret := C.ngt_save_index(n.index, path, ne.err)
		n.rUnlock(true)
		if ret == ErrorCode {
			return n.newGoError(ne)
		}
		n.PutErrorBuffer(ne)
	}

	return nil
}

// Remove removes from NGT index.
func (n *ngt) Remove(id uint) error {
	ne := n.GetErrorBuffer()
	n.lock(true)
	ret := C.ngt_remove_index(n.index, C.ObjectID(id), ne.err)
	n.unlock(true)
	if ret == ErrorCode {
		return n.newGoError(ne)
	}
	n.PutErrorBuffer(ne)

	n.cnt.Add(^uint64(0))

	return nil
}

// BulkRemove removes multiple index from NGT index.
func (n *ngt) BulkRemove(ids ...uint) (errs error) {
	for i, id := range ids {
		err := n.Remove(id)
		if err != nil {
			errs = errors.Wrapf(errs, "bulkremove error detected index number: %d,\tid: %d\terr: %v", i, id, err)
		}
	}
	return errs
}

// GetVector returns vector stored in NGT index.
func (n *ngt) GetVector(id uint) (ret []float32, err error) {
	dimension := int(n.dimension)
	ne := n.GetErrorBuffer()
	switch n.objectType {
	case Float:
		n.rLock(false)
		results := C.ngt_get_object_as_float(n.ospace, C.ObjectID(id), ne.err)
		n.rUnlock(false)
		if results == nil {
			return nil, n.newGoError(ne)
		}
		ret = (*[algorithm.MaximumVectorDimensionSize]float32)(unsafe.Pointer(results))[:dimension:dimension]
	case HalfFloat:
		n.rLock(false)
		results := C.ngt_get_allocated_object_as_float(n.ospace, C.ObjectID(id), ne.err)
		n.rUnlock(false)
		defer C.free(unsafe.Pointer(results))
		if results == nil {
			return nil, n.newGoError(ne)
		}
		ret = make([]float32, dimension)
		for i, elem := range (*[algorithm.MaximumVectorDimensionSize]float32)(unsafe.Pointer(results))[:dimension:dimension] {
			ret[i] = elem
		}
	case Uint8:
		n.rLock(false)
		results := C.ngt_get_object_as_integer(n.ospace, C.ObjectID(id), ne.err)
		n.rUnlock(false)
		if results == nil {
			return nil, n.newGoError(ne)
		}
		ret = make([]float32, 0, dimension)
		for _, elem := range (*[algorithm.MaximumVectorDimensionSize]C.uint8_t)(unsafe.Pointer(results))[:dimension:dimension] {
			ret = append(ret, float32(elem))
		}
	default:
		n.PutErrorBuffer(ne)
		return nil, errors.ErrUnsupportedObjectType
	}
	n.PutErrorBuffer(ne)
	return ret, nil
}

func (n *ngt) newGoError(ne *ngtError) (err error) {
	msg := C.GoString(C.ngt_get_error_string(ne.err))
	if len(msg) == 0 {
		n.PutErrorBuffer(ne)
		return nil
	}
	if n.epl == 0 || n.eps.Load() < n.epl {
		n.PutErrorBuffer(newNGTError())
	}
	ne.close()
	return errors.NewNGTError(msg)
}

func (n *ngt) GetErrorBuffer() (ne *ngtError) {
	var ok bool
	ne, ok = n.epool.Get().(*ngtError)
	if !ok {
		ne = newNGTError()
	}
	n.eps.Add(^uint64(0))
	return ne
}

func (n *ngt) PutErrorBuffer(ne *ngtError) {
	if n.epl != 0 && n.eps.Load() > n.epl {
		ne.close()
		return
	}
	n.epool.Put(ne)
	n.eps.Add(1)
}

func (n *ngt) lock(cLock bool) {
	if cLock {
		n.cmu.Lock()
	}
	n.mu.Lock()
}

func (n *ngt) unlock(cLock bool) {
	n.mu.Unlock()
	if cLock {
		n.cmu.Unlock()
	}
}

func (n *ngt) rLock(cLock bool) {
	if cLock {
		n.cmu.RLock()
	}
	n.mu.RLock()
}

func (n *ngt) rUnlock(cLock bool) {
	n.mu.RUnlock()
	if cLock {
		n.cmu.RUnlock()
	}
}

// Close NGT index.
func (n *ngt) Close() {
	if n.index != nil {
		C.ngt_close_index(n.index)
		n.index = nil
		n.prop = nil
		n.ospace = nil
	}
}
