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

//go:build race

package sync

import (
	"runtime"
	"unsafe"
)

const poolRaceEnabled = true

// Only the pinned owner of a shard accesses its private slot or the head of
// its queue. The runtime does not report handoffs between goroutines on the
// same P to the race detector. Use the shard's privateSet address as the event
// key, including when GC has moved the shard to the victim cache.
//
// TSAN also propagates these events to callers. This mode is intended for
// testing the per-P algorithm. Queue atomics and ordinary memory accesses
// remain instrumented.
func poolRaceAcquire(shard unsafe.Pointer) { runtime.RaceAcquire(shard) }

func poolRaceRelease(shard unsafe.Pointer) { runtime.RaceRelease(shard) }
