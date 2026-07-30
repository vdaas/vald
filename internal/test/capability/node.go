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

package capability

import "testing"

// Node is the type-erased face of Runner[X]. Runner-generic code must
// thread its X type parameter through every function that eventually
// spawns a subtest, which turns whole orchestration trees generic and
// forbids methods (Go methods cannot have type parameters). Node removes
// that constraint: NewNode captures the concrete X exactly once inside a
// closure, and everything downstream works with this small non-generic
// value instead — it embeds the running testing.TB (so a Node can be
// passed to anything expecting a testing.TB and all TB methods promote)
// and Run spawns subtests/sub-benchmarks of the same underlying concrete
// type, handing each child a fresh Node.
//
// The zero value is invalid (its embedded testing.TB is nil); always
// construct a Node via NewNode.
type Node struct {
	testing.TB
	run func(name string, fn func(Node)) bool
}

// Node must be usable wherever a plain testing.TB is expected.
var _ testing.TB = Node{}

// NewNode wraps t into a Node. This is the single point where generics are
// required: the returned Node's Run re-wraps every child in the same way,
// so the concrete type information survives arbitrarily deep subtest trees
// without any further type parameters.
func NewNode[X Runner[X]](t X) Node {
	return Node{
		TB: t,
		run: func(name string, fn func(Node)) bool {
			return t.Run(name, func(tt X) {
				fn(NewNode(tt))
			})
		},
	}
}

// Run runs fn as a subtest (or sub-benchmark) named name, handing it a
// Node wrapping the child's own testing entry.
func (n Node) Run(name string, fn func(Node)) bool {
	n.Helper()
	return n.run(name, fn)
}

// Unwrap returns the underlying testing entry (e.g. *testing.T or
// *testing.B), following the errors.Unwrap naming convention so As and the
// capability helpers built on it can reach the concrete type through any
// number of wrapping layers.
func (n Node) Unwrap() testing.TB { return n.TB }
