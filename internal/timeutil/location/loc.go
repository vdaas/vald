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

import (
	"strings"
	"time"
)

const (
	locationTokyo = "Asia/Tokyo"
	locationJST   = "JST"
	locationUTC   = "UTC"
	locationGMT   = "GMT"
)

func Set(loc string) {
	switch strings.ToLower(loc) {
	case strings.ToLower(locationUTC):
		time.Local = UTC()
	case strings.ToLower(locationGMT):
		time.Local = GMT()
	case strings.ToLower(locationJST), strings.ToLower(locationTokyo):
		time.Local = JST()
	default:
		time.Local = location(loc, 0)
	}
}

func GMT() *time.Location {
	return location(locationGMT, 0)
}

func UTC() *time.Location {
	return location(locationUTC, 0)
}

func JST() *time.Location {
	return location(locationJST, 9*60*60)
}

func location(zone string, offset int) *time.Location {
	loc, err := time.LoadLocation(zone)
	if err != nil {
		return time.FixedZone(zone, offset)
	}
	return loc
}
