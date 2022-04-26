// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package location

import (
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/vdaas/vald/internal/strings"
)

const (
	locationTokyo = "Asia/Tokyo"
	locationJST   = "JST"
	locationUTC   = "UTC"
	locationGMT   = "GMT"
)

var (
	gmt = location(locationGMT, 0)
	utc = location(locationUTC, 0)
	jst = location(locationJST, 9*60*60)
)

func Set(loc string) {
	var local *time.Location

	switch strings.ToLower(loc) {
	case strings.ToLower(locationUTC):
		local = UTC()
	case strings.ToLower(locationGMT):
		local = GMT()
	case strings.ToLower(locationJST), strings.ToLower(locationTokyo):
		local = JST()
	default:
		local = location(loc, 0)
	}

	new := atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&local)))
	old := atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&time.Local)))

	for !atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&time.Local)), old, new) {
		old = atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&time.Local)))
	}
}

func GMT() *time.Location {
	return gmt
}

func UTC() *time.Location {
	return utc
}

func JST() *time.Location {
	return jst
}

func location(zone string, offset int) *time.Location {
	loc, err := time.LoadLocation(zone)
	if err != nil {
		return time.FixedZone(zone, offset)
	}
	return loc
}
