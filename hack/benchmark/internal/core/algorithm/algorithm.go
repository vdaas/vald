//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
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

type Bit32 interface {
	Search(ctx context.Context, vec []float32, size int, epsilon, radius float32) (interface{}, error)
	Insert(vec []float32) (uint, error)
	InsertCommit(vec []float32, poolSize uint32) (uint, error)
	BulkInsert(vecs [][]float32) ([]uint, []error)
	BulkInsertCommit(vecs [][]float32, poolSize uint32) ([]uint, []error)
	CreateAndSaveIndex(poolSize uint32) error
	CreateIndex(poolSize uint32) error
	SaveIndex() error
	Remove(id uint) error
	BulkRemove(ids ...uint) error
	GetVector(id uint) ([]float32, error)
	Closer
}

type Bit64 interface {
	Search(vec []float64, size int, epsilon, radius float32) (interface{}, error)
	Insert(vec []float64) (uint, error)
	InsertCommit(vec []float64, poolSize uint32) (uint, error)
	BulkInsert(vecs [][]float64) ([]uint, []error)
	BulkInsertCommit(vecs [][]float64, poolSize uint32) ([]uint, []error)
	CreateAndSaveIndex(poolSize uint32) error
	CreateIndex(poolSize uint32) error
	SaveIndex() error
	Remove(id uint) error
	BulkRemove(ids ...uint) error
	GetVector(id uint) ([]float64, error)
	Closer
}
