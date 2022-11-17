//
// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
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
package location

import (
	"sync"
	"sync/atomic"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
)

var (
	poolSize uint64 = 10
	locPool         = &sync.Pool{
		New: func() interface{} {
			return make(map[string]*payload.Object_Location, atomic.LoadUint64(&poolSize))
		},
	}
)

func ReStructure(uuids []string, locs *payload.Object_Locations) *payload.Object_Locations {
	if locs == nil || locs.Locations == nil {
		return nil
	}
	ll := uint64(len(locs.GetLocations()))
	if ll > atomic.LoadUint64(&poolSize) {
		atomic.StoreUint64(&poolSize, uint64(len(locs.GetLocations())))
	}
	lp := locPool.Get()
	lmap, ok := lp.(map[string]*payload.Object_Location)
	if !ok || lmap == nil {
		lmap = make(map[string]*payload.Object_Location, ll)
	}
	for _, loc := range locs.Locations {
		uuid := loc.GetUuid()
		lm, ok := lmap[uuid]
		if !ok || lm == nil {
			lmap[uuid] = new(payload.Object_Location)
		}
		lmap[uuid].Ips = append(lmap[uuid].GetIps(), loc.GetIps()...)
	}
	locs.Locations = locs.Locations[:0]
	for _, id := range uuids {
		loc, ok := lmap[id]
		if !ok {
			loc = new(payload.Object_Location)
		} else {
			delete(lmap, id)
		}
		locs.Locations = append(locs.GetLocations(), loc)
	}
	for id, lm := range lmap {
		locs.Locations = append(locs.GetLocations(), lm)
		delete(lmap, id)
	}
	locPool.Put(lmap)
	return locs
}
