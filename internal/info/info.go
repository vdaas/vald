//
// Copyright (C) 2019-2020 Vdaas.org Vald team ( kpango, rinx, kmrmt )
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
	"sync"
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

	keyvals = map[string]*string{
		"version":     &Version,
		"commit hash": &GitCommit,
		"build time":  &BuildTime,
		"go version":  &GoVersion,
		"os":          &GoOS,
		"arch":        &GoArch,
		"cgo enabled": &CGOEnabled,
		"ngt version": &NGTVersion,
	}

	once sync.Once
	name string
)

func Init(n string) {
	once.Do(func() {
		name = n
	})
}

func Info() {
	showVersionInfo(log.Info)
}

func Debug() {
	showVersionInfo(log.Debug)
}

func Warn() {
	showVersionInfo(log.Warn)
}

func Error() {
	showVersionInfo(log.Error)
}

func showVersionInfo(logfunc func(vals ...interface{})) {
	keys := make([]string, 0, len(keyvals))
	for k := range keyvals {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var l int
	maxlen := 0
	for _, k := range keys {
		l = len(k)
		if maxlen < l {
			maxlen = l
		}
	}

	infoFormat := fmt.Sprintf("%%-%ds -> %%s", maxlen)

	strs := make([]string, 0, len(keys))
	if name != "" {
		strs = append(strs, fmt.Sprintf("\nvald %s server", name))
	}
	for _, k := range keys {
		if k != "" && *keyvals[k] != "" {
			if k == "version" {
				strs = append(strs, fmt.Sprintf(infoFormat, k, log.Bold(*keyvals[k])))
			} else {
				strs = append(strs, fmt.Sprintf(infoFormat, k, *keyvals[k]))
			}
		}
	}

	logfunc(strings.Join(strs, "\n"))
}
