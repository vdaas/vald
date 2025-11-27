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

package metrics

import (
	"context"
	"unsafe"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/sync"
)

// scale is a ring buffer of slots for windowed metrics.
// It implements the Scale interface and manages a collection of slots based on time or request ID.
type scale struct {
	mu        sync.RWMutex
	slots     []Slot // ring buffer of slots
	width     uint64 // width of each slot (e.g., seconds for TimeScale)
	capacity  uint64 // number of slots in the ring buffer
	name      string // name of the scale
	scaleType ScaleType
}

// newScale creates a new scale with the given configuration.
// It initializes the ring buffer of slots, cloning the provided histogram/exemplar prototypes.
func newScale(
	name string,
	width, capacity uint64,
	numCounters int,
	st ScaleType,
	lat, qw Histogram,
	ex Exemplar,
) (Scale, error) {
	if width == 0 {
		return nil, errors.New("scale width must be > 0")
	}
	if capacity == 0 {
		return nil, errors.New("scale capacity must be > 0")
	}
	slots := make([]Slot, capacity)
	for i := range slots {
		// Use Clone to create new instances for each slot.
		// Clone handles initialization via pool.
		var l, q Histogram
		var e Exemplar
		if lat != nil {
			l = lat.Clone()
		}
		if qw != nil {
			q = qw.Clone()
		}
		if ex != nil {
			e = ex.Clone()
		}

		slots[i] = newSlot(
			numCounters,
			l,
			q,
			e,
		)
	}
	return &scale{
		name:      name,
		width:     width,
		capacity:  capacity,
		slots:     slots,
		scaleType: st,
	}, nil
}

// Reset resets the scale and all its slots.
// It locks the scale and iterates through all slots to reset them.
func (s *scale) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, slot := range s.slots {
		slot.Reset()
	}
}

// getSlot returns the slot for the given index.
// It calculates the ring buffer index based on the scale's width and capacity.
func (s *scale) getSlot(idx uint64) Slot {
	return s.slots[int((idx/s.width)%s.capacity)]
}

// Record adds a request result to the appropriate slot in the scale.
// It determines the index based on the scale type (Time or Range) and delegates recording to the target slot.
func (s *scale) Record(ctx context.Context, rr *RequestResult) {
	var idx uint64
	var ok bool

	switch s.scaleType {
	case RangeScale:
		idx, ok = requestIDFromCtx(ctx)
		if !ok {
			return
		}
	case TimeScale:
		idx = uint64(rr.EndedAt.UnixNano())
	}

	s.getSlot(idx).Record(rr, idx/s.width)
}

// Merge merges another scale into this one.
// It validates compatibility (width, capacity) and merges slot by slot.
func (s *scale) Merge(other Scale) error {
	if s == other {
		return nil
	}

	os, ok := other.(*scale)
	if !ok {
		return errors.New("incompatible scale implementation")
	}

	if s.width != os.width || s.capacity != os.capacity {
		return errors.New("incompatible scales")
	}

	// To prevent deadlocks, always lock in a consistent order.
	if uintptr(unsafe.Pointer(s)) < uintptr(unsafe.Pointer(os)) {
		s.mu.Lock()
		os.mu.Lock()
	} else {
		os.mu.Lock()
		s.mu.Lock()
	}
	defer s.mu.Unlock()
	defer os.mu.Unlock()

	for i := range s.slots {
		ss := s.slots[i]
		oSlot := os.slots[i]

		if err := ss.Merge(oSlot); err != nil {
			return err
		}
	}
	return nil
}

// Snapshot returns a snapshot of the scale.
// It collects snapshots from all slots and packages them into a ScaleSnapshot.
func (s *scale) Snapshot() *ScaleSnapshot {
	s.mu.RLock()
	defer s.mu.RUnlock()

	slots := make([]*SlotSnapshot, len(s.slots))
	for i, slot := range s.slots {
		slots[i] = slot.Snapshot()
	}

	return &ScaleSnapshot{
		Name:     s.name,
		Width:    s.width,
		Capacity: s.capacity,
		Slots:    slots,
	}
}

// Clone returns a deep copy of the scale.
// It creates a new independent scale with cloned slots.
func (s *scale) Clone() Scale {
	s.mu.RLock()
	defer s.mu.RUnlock()

	slots := make([]Slot, len(s.slots))
	for i, slot := range s.slots {
		slots[i] = slot.Clone()
	}

	return &scale{
		name:      s.name,
		width:     s.width,
		capacity:  s.capacity,
		scaleType: s.scaleType,
		slots:     slots,
	}
}

// Type returns the type of the scale.
func (s *scale) Type() ScaleType {
	return s.scaleType
}

// Name returns the name of the scale.
func (s *scale) Name() string {
	return s.name
}
