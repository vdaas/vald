// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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

type PairingHeap struct {
	DistPayload *DistPayload
	Children    []*PairingHeap
}

func (ph *PairingHeap) IsEmpty() bool {
	return ph == nil
}

func (ph *PairingHeap) Insert(dp *DistPayload) *PairingHeap {
	return ph.Merge(&PairingHeap{
		DistPayload: dp,
	})
}

func (ph *PairingHeap) Merge(h2 *PairingHeap) *PairingHeap {
	if ph == nil || ph.DistPayload == nil || ph.DistPayload.distance == nil {
		return h2
	}
	if h2 == nil || h2.DistPayload == nil || h2.DistPayload.distance == nil {
		return ph
	}

	var smaller, larger *PairingHeap
	if ph.DistPayload.distance.Cmp(h2.DistPayload.distance) < 0 {
		smaller = ph
		larger = h2
	} else {
		smaller = h2
		larger = ph
	}

	smaller.Children = append(smaller.Children, larger)
	return smaller
}

func (ph *PairingHeap) ExtractMin() (*DistPayload, *PairingHeap) {
	if ph == nil {
		return nil, nil
	}
	minDistPayload := ph.DistPayload
	newHeap := ph.mergePairs(ph.Children)
	return minDistPayload, newHeap
}

func (ph *PairingHeap) mergePairs(pairs []*PairingHeap) *PairingHeap {
	switch len(pairs) {
	case 0:
		return nil
	case 1:
		return pairs[0]
	}
	first := pairs[0].Merge(pairs[1])
	remaining := pairs[2:]
	return first.Merge(ph.mergePairs(remaining))
}
