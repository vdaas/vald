//
// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
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

// package iter
package iter

import (
	"context"
	"iter"
)

// -------------------------
// Iterator Type Definition |
// -------------------------
// Cycle provides an iterator abstraction over a slice.
type Cycle[S ~[]E, E any] interface {
	At(i uint64) E
	ForEach(ctx context.Context, fn func(uint64, E) bool)
	Len() uint64
	Raw() S
	Seq(context.Context) iter.Seq[E]
	Seq2(context.Context) iter.Seq2[uint64, E]
	Indexes(context.Context) iter.Seq[uint64]
	Values(context.Context) iter.Seq[E]
}

// cycle provides an iterator abstraction over a slice.
// It applies an optional modFunc to transform each element on‑the‑fly without precomputing the entire dataset.
type cycle[S ~[]E, E any] struct {
	array   S
	modFunc func(uint64, E) E
	start   uint64
	num     uint64
	size    uint64
	offset  uint64
}

// New creates a new cycle iterator instance. It validates the input array and computes the starting index (offset modulo array size).
func NewCycle[S ~[]E, E any](array S, num, offset uint64, mod func(uint64, E) E) Cycle[S, E] {
	if array == nil {
		return nil
	}
	size := uint64(len(array))
	if size == 0 {
		return nil
	}
	return &cycle[S, E]{
		start:   offset % size,
		num:     num,
		size:    size,
		offset:  offset,
		array:   array,
		modFunc: mod,
	}
}

// At returns the element at logical index i.
// If modFunc is provided, it applies the function on‑the‑fly.
func (c *cycle[_, E]) At(i uint64) E {
	idx := (c.start + i) % c.size
	if c.modFunc != nil {
		return c.modFunc(i, c.array[idx])
	}
	return c.array[idx]
}

// Seq2 returns an iterator sequence (iter.Seq2) that yields each element along with its index.
func (c *cycle[_, E]) Seq2(ctx context.Context) iter.Seq2[uint64, E] {
	return func(yield func(uint64, E) bool) {
		for i := uint64(0); i < c.num; i++ {
			select {
			case <-ctx.Done():
				return
			default:
			}
			if !yield(i, c.At(i)) {
				return
			}
		}
	}
}

func (c *cycle[_, E]) Seq(ctx context.Context) iter.Seq[E] {
	return c.Values(ctx)
}

// Values returns an iterator sequence (iter.Seq) that yields the values (without indexes).
func (c *cycle[_, E]) Values(ctx context.Context) iter.Seq[E] {
	return func(yield func(E) bool) {
		for i := uint64(0); i < c.num; i++ {
			select {
			case <-ctx.Done():
				return
			default:
			}
			if !yield(c.At(i)) {
				return
			}
		}
	}
}

// Indexes returns an iterator sequence (iter.Seq) that yields the indexes.
func (c cycle[_, _]) Indexes(ctx context.Context) iter.Seq[uint64] {
	return func(yield func(uint64) bool) {
		for i := uint64(0); i < c.num; i++ {
			select {
			case <-ctx.Done():
				return
			default:
			}
			if !yield(i) {
				return
			}
		}
	}
}

func (c cycle[_, E]) ForEach(ctx context.Context, fn func(uint64, E) bool) {
	for i := uint64(0); i < c.num; i++ {
		select {
		case <-ctx.Done():
			return
		default:
		}
		if !fn(i, c.At(i)) {
			return
		}
	}
}

func (c cycle[S, _]) Raw() S {
	return c.array
}

// Len returns the total number of elements in the iterator.
func (c cycle[_, _]) Len() uint64 {
	return c.num
}
