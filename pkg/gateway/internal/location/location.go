//
// Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
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

import "github.com/vdaas/vald/apis/grpc/v1/payload"

func ReStructure(uuids []string, locs *payload.Object_Locations) *payload.Object_Locations {
	if locs == nil {
		return nil
	}
	lmap := make(map[string]*payload.Object_Location, len(locs.Locations))
	for _, loc := range locs.Locations {
		uuid := loc.GetUuid()
		_, ok := lmap[uuid]
		if !ok {
			lmap[uuid] = new(payload.Object_Location)
		}
		lmap[uuid].Ips = append(lmap[uuid].GetIps(), loc.GetIps()...)
	}
	locs = &payload.Object_Locations{
		Locations: make([]*payload.Object_Location, 0, len(lmap)),
	}
	for _, id := range uuids {
		loc, ok := lmap[id]
		if !ok {
			loc = new(payload.Object_Location)
		}
		locs.Locations = append(locs.Locations, loc)
	}
	return locs
}
