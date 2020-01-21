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
	"encoding/json"
	"fmt"
	"reflect"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"
	"unsafe"

	"github.com/vdaas/vald/internal/log"
)

type Detail struct {
	Version    string `json:"version,omitempty" yaml:"version,omitempty"`
	ServerName string `json:"server_name,omitempty" yaml:"server_name,omitempty"`
	GitCommit  string `json:"git_commit,omitempty" yaml:"git_commit,omitempty"`
	SourceURL  string `json:"source_url,omitempty" yaml:"source_url,omitempty"`
	BuildTime  string `json:"build_time,omitempty" yaml:"build_time,omitempty"`
	GoVersion  string `json:"go_version,omitempty" yaml:"go_version,omitempty"`
	GoOS       string `json:"go_os,omitempty" yaml:"go_os,omitempty"`
	GoArch     string `json:"go_arch,omitempty" yaml:"go_arch,omitempty"`
	CGOEnabled string `json:"cgo_enabled,omitempty" yaml:"cgo_enabled,omitempty"`
	NGTVersion string `json:"ngt_version,omitempty" yaml:"ngt_version,omitempty"`
}

var (
	Version   = "v0.0.1"
	GitCommit = "no commit info available."
	BuildTime = time.Now().Format(time.RFC1123)

	GoVersion  string
	GoOS       string
	GoArch     string
	CGOEnabled string

	NGTVersion string

	reps = strings.NewReplacer("_", " ", ",omitempty", "")

	once sync.Once

	detail Detail
)

const (
	organization = "vdaas"
	repository   = "vald"
)

func String(callStack int) string {
	return detail.String(callStack)
}

func Object(callStack int) Detail {
	return detail.Object(callStack)
}

func JSONString(format bool, callStack int) string {
	return detail.JSONString(format, callStack)
}

func (d Detail) String(callStack int) string {
	d.SourceURL = sourceURL(organization, repository, d.GitCommit, callStack)
	maxlen, l := 0, 0
	rt, rv := reflect.TypeOf(d), reflect.ValueOf(d)
	info := make(map[string]string, rt.NumField())
	for i := 0; i < rt.NumField(); i++ {
		tag := reps.Replace(rt.Field(i).Tag.Get("json"))
		value, ok := rv.Field(i).Interface().(string)
		if !ok {
			continue
		}
		l = len(value)
		if maxlen < l {
			maxlen = l
		}
		info[tag] = value
	}

	infoFormat := fmt.Sprintf("%%-%ds -> %%s", maxlen)
	strs := make([]string, 0, rt.NumField())
	for tag, value := range info {
		strs = append(strs, fmt.Sprintf(infoFormat, tag, value))
	}
	sort.Strings(strs)
	return strings.Join(strs, "\n")
}

func (d Detail) Object(callStack int) Detail {
	d.SourceURL = sourceURL(organization, repository, d.GitCommit, callStack)
	return d
}

func (d Detail) JSONString(format bool, callStack int) (str string) {
	var err error
	var b []byte
	d.SourceURL = sourceURL(organization, repository, d.GitCommit, callStack)
	if !format {
		b, err = json.Marshal(d)
	} else {
		b, err = json.MarshalIndent(d, "", "\t")
	}
	if err != nil {
		return d.String(callStack)
	}
	return *(*string)(unsafe.Pointer(&b))
}

func Init(name string) {
	once.Do(func() {
		detail = Detail{
			Version:    log.Bold(Version),
			ServerName: name,
			GitCommit:  GitCommit,
			BuildTime:  BuildTime,
			GoVersion:  GoVersion,
			GoOS:       GoOS,
			GoArch:     GoArch,
			CGOEnabled: CGOEnabled,
			NGTVersion: NGTVersion,
		}
	})
}

func sourceURL(org, repo, commit string, caller int) string {
	if caller < 2 {
		return fmt.Sprintf("https://github.com/%s/%s/blob/%s/",
			org, repo, commit)
	}

	_, file, line, ok := runtime.Caller(caller)
	if !ok {
		return fmt.Sprintf("https://github.com/%s/%s/blob/%s/",
			org, repo, commit)
	}
	return fmt.Sprintf("https://github.com/%s/%s/blob/%s/%s#L%d",
		org, repo, commit, strings.SplitN(file,
			fmt.Sprintf("%s/%s/", org, repo), 2)[1], line)
}
