// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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
	"os"
	"time"

	"github.com/kpango/fastime"
	"github.com/vdaas/vald/internal/strings"
)

func Set(loc string) {
	switch strings.ToLower(loc) {
	case strings.ToLower(locationUTC):
		time.Local = UTC()
		os.Setenv("TZ", "UTC")
	case strings.ToLower(locationGMT):
		time.Local = GMT()
		os.Setenv("TZ", "GMT")
	case strings.ToLower(locationJST), strings.ToLower(locationTokyo):
		time.Local = JST()
		os.Setenv("TZ", "JST")
	default:
		time.Local = location(loc, 0)
	}
	fastime.SetLocation(time.Local)
}
