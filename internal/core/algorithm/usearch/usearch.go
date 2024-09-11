//
// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
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

// Package usearch provides Go API implementation for USearch library. https://github.com/unum-cloud/usearch
package usearch

import (
	"sync"

	core "github.com/unum-cloud/usearch/golang"
	"github.com/vdaas/vald/internal/core/algorithm"
	"github.com/vdaas/vald/internal/errors"
)

type (
	// Uuearch is the core interface for interacting with usearch index.
	Usearch interface {
		// SaveIndex saves the USearch index to storage.
		SaveIndex() error

		// SaveIndexWithPath saves the USearch index to the specified path.
		SaveIndexWithPath(path string) error

		// GetIndicesSize returns the number of vectors in index.
		GetIndicesSize() (indicesSize int, err error)

		// Reserve reserves memory for vectors of given number of arg.
		Reserve(vectorCount int) error

		// Add adds vectors to the USearch index and returns the total count.
		Add(key uint64, vec []float32) error

		// Search performs a nearest neighbor search and returns the results.
		Search(q []float32, k int) ([]algorithm.SearchResult, error)

		// GetObject retruns search result by id as []algorithm.SearchResult.
		GetObject(key core.Key, count int) ([]float32, error)

		// Remove removes vectors from the index by key.
		Remove(key uint64) error

		// Close frees the resources used by the USearch index.
		Close() error
	}

	usearch struct {
		// index struct
		index *core.Index

		// config
		quantizationType core.Quantization
		metricType       core.Metric
		dimension        uint
		connectivity     uint
		expansionAdd     uint
		expansionSearch  uint
		multi            bool

		idxPath string
		mu      *sync.RWMutex
	}
)

// New initializes a new USearch instance with the provided options.
func New(opts ...Option) (Usearch, error) {
	return gen(false, opts...)
}

func Load(opts ...Option) (Usearch, error) {
	return gen(true, opts...)
}

func gen(isLoad bool, opts ...Option) (Usearch, error) {
	var (
		u   = new(usearch)
		err error
	)
	u.mu = new(sync.RWMutex)

	for _, opt := range append(defaultOptions, opts...) {
		if err = opt(u); err != nil {
			return nil, errors.NewUsearchError("usarch option error :" + err.Error())
		}
	}

	if isLoad {
		conf := core.DefaultConfig(uint(u.dimension))
		u.index, err = core.NewIndex(conf)
		if err != nil {
			return nil, errors.NewUsearchError("usearch new index error for load index")
		}

		err = u.index.Load(u.idxPath)
		if err != nil {
			return nil, errors.NewUsearchError("usearch load index error")
		}
	} else {
		options := core.DefaultConfig(u.dimension)
		options.Quantization = u.quantizationType
		options.Metric = u.metricType
		options.Dimensions = u.dimension
		options.Connectivity = u.connectivity
		options.ExpansionAdd = u.expansionAdd
		options.ExpansionSearch = u.expansionSearch
		options.Multi = u.multi

		u.index, err = core.NewIndex(options)
		if err != nil {
			return nil, errors.NewUsearchError("usearch create index error")
		}
	}

	return u, nil
}

// SaveIndex stores usearch index to storage.
func (u *usearch) SaveIndex() error {
	u.mu.Lock()
	defer u.mu.Unlock()

	err := u.index.Save(u.idxPath)
	if err != nil {
		return errors.NewUsearchError("usarch save index error")
	}
	return nil
}

// SaveIndexWithPath stores usearch index to specified storage.
func (u *usearch) SaveIndexWithPath(idxPath string) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	err := u.index.Save(idxPath)
	if err != nil {
		return errors.NewUsearchError("usarch save index with path error")
	}
	return nil
}

// GetIndicesSize returns the number of vectors in index.
func (u *usearch) GetIndicesSize() (indicesSize int, err error) {
	u.mu.Lock()
	defer u.mu.Unlock()
	size, err := u.index.Len()
	if err != nil {
		return -1, errors.NewUsearchError("failed to usearch_size")
	}
	return int(size), err
}

// Add adds vectors to the index
func (u *usearch) Add(key core.Key, vec []float32) error {
	if len(vec) != int(u.dimension) {
		return errors.New("inconsistent dimensions")
	}

	u.mu.Lock()
	err := u.index.Add(key, vec)
	defer u.mu.Unlock()
	if err != nil {
		return errors.NewUsearchError("failed to usearch_add")
	}
	return nil
}

// Reserve reserves memory for vectors of given number of arg.
func (u *usearch) Reserve(vectorCount int) error {
	u.mu.Lock()
	err := u.index.Reserve(uint(vectorCount))
	defer u.mu.Unlock()
	if err != nil {
		return errors.NewUsearchError("failed to usearch_reserve")
	}
	return nil
}

// Search returns search result as []algorithm.SearchResult.
func (u *usearch) Search(q []float32, k int) ([]algorithm.SearchResult, error) {
	if len(q) != int(u.dimension) {
		return nil, errors.ErrIncompatibleDimensionSize(len(q), int(u.dimension))
	}
	u.mu.Lock()
	I, D, err := u.index.Search(q, uint(k))
	u.mu.Unlock()
	if err != nil {
		return nil, errors.NewUsearchError("failed to usearch_search")
	}

	if len(I) == 0 || len(D) == 0 {
		return nil, errors.ErrEmptySearchResult
	}

	result := make([]algorithm.SearchResult, min(len(I), k))
	for i := range result {
		result[i] = algorithm.SearchResult{ID: uint32(I[i]), Distance: D[i], Error: nil}
	}
	return result, nil
}

// GetObject retruns search result by id as []algorithm.SearchResult.
func (u *usearch) GetObject(key core.Key, count int) ([]float32, error) {
	u.mu.RLock()
	vectors, err := u.index.Get(key, uint(count))
	u.mu.RUnlock()
	if err != nil {
		return nil, errors.NewUsearchError("failed to usearch_get")
	}
	// ASK: 何か適切なerrorがある？
	if vectors == nil {
		return nil, nil
	}

	return vectors, nil
}

// Remove removes from usearch index.
func (u *usearch) Remove(key core.Key) error {
	u.mu.Lock()
	err := u.index.Remove(key)
	defer u.mu.Unlock()
	if err != nil {
		return errors.NewUsearchError("failed to usearch_remove")
	}

	return nil
}

// Close frees the resources associated with the USearch index.
func (u *usearch) Close() error {
	err := u.index.Destroy()
	if err != nil {
		return errors.NewUsearchError("failed to usearch_free")
	}
	u.index = nil
	return nil
}
