//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
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

// Package algorithm provides core interface
package algorithm

import "context"

type Mode uint32

const (
	Float32 Mode = iota
	Float64
)

type Closer interface {
	Close()
}

// Bit is the interface for bit operations.
type Bit[T float32 | float64] interface {
	// Search searches for a vector.
	Search(ctx context.Context, vec []T, size int, epsilon, radius float32) (any, error)
	// Insert inserts a vector.
	Insert(vec []T) (uint, error)
	// InsertCommit inserts a vector and commits.
	InsertCommit(vec []T, poolSize uint32) (uint, error)
	// BulkInsert inserts multiple vectors.
	BulkInsert(vecs [][]T) ([]uint, []error)
	// BulkInsertCommit inserts multiple vectors and commits.
	BulkInsertCommit(vecs [][]T, poolSize uint32) ([]uint, []error)
	// CreateAndSaveIndex creates and saves the index.
	CreateAndSaveIndex(poolSize uint32) error
	// CreateIndex creates the index.
	CreateIndex(poolSize uint32) error
	// SaveIndex saves the index.
	SaveIndex() error
	// Remove removes a vector.
	Remove(id uint) error
	// BulkRemove removes multiple vectors.
	BulkRemove(ids ...uint) error
	// GetVector gets a vector.
	GetVector(id uint) ([]T, error)
	Closer
}

type Bit32 Bit[float32]

type Bit64 Bit[float64]
