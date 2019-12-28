//
// Copyright (C) 2019 Vdaas.org Vald team ( kpango, kmrmt, rinx )
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

// Package info provides build-time info
package info

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/vdaas/vald/internal/log"
)

var (
	Version   = "v0.0.1"
	GitCommit = "no commit info available."
	BuildTime = time.Now().Format(time.RFC1123)

	GoVersion  string
	GoOS       string
	GoArch     string
	CGOEnabled string
	NGTVersion string
)

func ShowVersionInfo(extras map[string]string) func(name string) {
	return func(name string) {
		defaultKeys := []string{"version", "commit hash", "build time"}
		keys := make([]string, 0, len(extras))
		for k := range extras {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		var l int
		maxlen := 0
		for _, k := range append(defaultKeys, keys...) {
			l = len(k)
			if maxlen < l {
				maxlen = l
			}
		}

		infoFormat := fmt.Sprintf("%%-%ds -> %%s", maxlen)

		strs := make([]string, 0, len(keys)+4)
		strs = append(strs, fmt.Sprintf("vald %s server", name))
		strs = append(strs, fmt.Sprintf(infoFormat, defaultKeys[0], log.Bold(Version)))
		strs = append(strs, fmt.Sprintf(infoFormat, defaultKeys[1], GitCommit))
		strs = append(strs, fmt.Sprintf(infoFormat, defaultKeys[2], BuildTime))
		for _, k := range keys {
			if k != "" && extras[k] != "" {
				strs = append(strs, fmt.Sprintf(infoFormat, k, extras[k]))
			}
		}

		log.Info(strings.Join(strs, "\n"))
	}
}
