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
	"time"
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
