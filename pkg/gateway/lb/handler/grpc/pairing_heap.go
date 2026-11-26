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
package grpc

import "github.com/vdaas/vald/internal/sync"

const pairingNull int32 = -1

const defaultPairingArenaCap = 32

// pairingNode is a Left-Child Right-Sibling node stored by index in a pooled arena slice.
type pairingNode struct {
	dp      DistPayload
	child   int32
	sibling int32
}

var pairingArenaPool = sync.Pool[[]pairingNode]{ //nolint:gochecknoglobals // pairing-heap node arenas are reused per search request
	New: func() []pairingNode {
		return make([]pairingNode, 0, defaultPairingArenaCap)
	},
}

func mergeMinNodes(nodes []pairingNode, a, b int32) int32 {
	if a == pairingNull {
		return b
	}
	if b == pairingNull {
		return a
	}
	if nodes[b].dp.distance < nodes[a].dp.distance {
		a, b = b, a
	}
	nodes[b].sibling = nodes[a].child
	nodes[a].child = b
	return a
}

func mergeMaxNodes(nodes []pairingNode, a, b int32) int32 {
	if a == pairingNull {
		return b
	}
	if b == pairingNull {
		return a
	}
	if nodes[b].dp.distance > nodes[a].dp.distance {
		a, b = b, a
	}
	nodes[b].sibling = nodes[a].child
	nodes[a].child = b
	return a
}

func mergePairsMinNodes(nodes []pairingNode, first int32) int32 {
	if first == pairingNull {
		return pairingNull
	}
	mergedList := pairingNull
	cur := first
	for cur != pairingNull {
		a := cur
		b := nodes[a].sibling
		if b == pairingNull {
			nodes[a].sibling = mergedList
			mergedList = a
			break
		}
		next := nodes[b].sibling
		nodes[a].sibling = pairingNull
		nodes[b].sibling = pairingNull
		m := mergeMinNodes(nodes, a, b)
		nodes[m].sibling = mergedList
		mergedList = m
		cur = next
	}
	acc := mergedList
	cur = nodes[acc].sibling
	nodes[acc].sibling = pairingNull
	for cur != pairingNull {
		next := nodes[cur].sibling
		nodes[cur].sibling = pairingNull
		acc = mergeMinNodes(nodes, cur, acc)
		cur = next
	}
	return acc
}

func mergePairsMaxNodes(nodes []pairingNode, first int32) int32 {
	if first == pairingNull {
		return pairingNull
	}
	mergedList := pairingNull
	cur := first
	for cur != pairingNull {
		a := cur
		b := nodes[a].sibling
		if b == pairingNull {
			nodes[a].sibling = mergedList
			mergedList = a
			break
		}
		next := nodes[b].sibling
		nodes[a].sibling = pairingNull
		nodes[b].sibling = pairingNull
		m := mergeMaxNodes(nodes, a, b)
		nodes[m].sibling = mergedList
		mergedList = m
		cur = next
	}
	acc := mergedList
	cur = nodes[acc].sibling
	nodes[acc].sibling = pairingNull
	for cur != pairingNull {
		next := nodes[cur].sibling
		nodes[cur].sibling = pairingNull
		acc = mergeMaxNodes(nodes, cur, acc)
		cur = next
	}
	return acc
}
