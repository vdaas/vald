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

// Package info provides build-time info
package info

import (
	"fmt"
	"reflect"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
)

// InformationProvider represents the runtime information provider.
type InformationProvider interface {
	String() string
	Get() Detail
}

type info struct {
	detail   Detail
	prepOnce sync.Once

	// runtime functions
	rtCaller    func(skip int) (pc uintptr, file string, line int, ok bool)
	rtFuncForPC func(pc uintptr) *runtime.Func
}

// Detail represents environment information of system and stacktrace information.
type Detail struct {
	Version           string       `json:"vald_version,omitempty" yaml:"vald_version,omitempty"`
	ServerName        string       `json:"server_name,omitempty" yaml:"server_name,omitempty"`
	GitCommit         string       `json:"git_commit,omitempty" yaml:"git_commit,omitempty"`
	BuildTime         string       `json:"build_time,omitempty" yaml:"build_time,omitempty"`
	GoVersion         string       `json:"go_version,omitempty" yaml:"go_version,omitempty"`
	GoOS              string       `json:"go_os,omitempty" yaml:"go_os,omitempty"`
	GoArch            string       `json:"go_arch,omitempty" yaml:"go_arch,omitempty"`
	CGOEnabled        string       `json:"cgo_enabled,omitempty" yaml:"cgo_enabled,omitempty"`
	NGTVersion        string       `json:"ngt_version,omitempty" yaml:"ngt_version,omitempty"`
	BuildCPUInfoFlags []string     `json:"build_cpu_info_flags,omitempty" yaml:"build_cpu_info_flags,omitempty"`
	StackTrace        []StackTrace `json:"stack_trace,omitempty" yaml:"stack_trace,omitempty"`
}

// StackTrace represents stacktrace information about url, function name, file, line ..etc.
type StackTrace struct {
	URL      string `json:"url,omitempty" yaml:"url,omitempty"`
	FuncName string `json:"function_name,omitempty" yaml:"func_name,omitempty"`
	File     string `json:"file,omitempty" yaml:"file,omitempty"`
	Line     int    `json:"line,omitempty" yaml:"line,omitempty"`
}

var (
	// injected from build script
	Version           = "v0.0.1"
	GitCommit         = "master"
	BuildTime         = ""
	GoVersion         string
	GoOS              string
	GoArch            string
	CGOEnabled        string
	NGTVersion        string
	BuildCPUInfoFlags string

	Organization = "vdaas"
	Repository   = "vald"

	reps = strings.NewReplacer("_", " ", ",omitempty", "")

	once         sync.Once
	infoProvider InformationProvider
)

// Init initializes Detail object only once.
func Init(name string) {
	once.Do(func() {
		infoProvider, _ = New(WithServerName(name))
	})
}

// New initialize and return the information provider or any error occurred.
func New(opts ...Option) (InformationProvider, error) {
	i := &info{
		detail: Detail{
			Version:           Version,
			GitCommit:         GitCommit,
			BuildTime:         BuildTime,
			GoVersion:         GoVersion,
			GoOS:              GoOS,
			GoArch:            GoArch,
			CGOEnabled:        CGOEnabled,
			NGTVersion:        NGTVersion,
			BuildCPUInfoFlags: strings.Split(strings.TrimSpace(BuildCPUInfoFlags), " "),
		},
	}

	for _, opt := range append(defaultOpts, opts...) {
		if err := opt(i); err != nil {
			werr := errors.ErrOptionFailed(err, reflect.ValueOf(opt))

			e := new(errors.ErrCriticalOption)
			if errors.As(err, &e) {
				log.Error(werr)
				return nil, werr
			}
			log.Warn(werr)
		}
	}

	i.prepare()

	return i, nil
}

// String calls String method of global detail object.
func String() string {
	return infoProvider.String()
}

// Get calls Get method of global detail object.
func Get() Detail {
	return infoProvider.Get()
}

// String returns summary of Detail object.
func (i info) String() string {
	if len(i.detail.StackTrace) == 0 {
		i.detail = i.Get()
	}

	d := i.detail
	d.Version = log.Bold(d.Version)
	maxlen, l := 0, 0
	rt, rv := reflect.TypeOf(d), reflect.ValueOf(d)
	info := make(map[string]string, rt.NumField())
	for i := 0; i < rt.NumField(); i++ {
		v := rv.Field(i).Interface()
		value, ok := v.(string)
		if !ok {
			sts, ok := v.([]StackTrace)
			if ok {
				tag := reps.Replace(rt.Field(i).Tag.Get("json"))
				l = len(tag) + 2
				if maxlen < l {
					maxlen = l
				}
				urlMaxLen := 0
				for _, st := range sts {
					ul := len(st.URL)
					if urlMaxLen < ul {
						urlMaxLen = ul
					}
				}
				urlFormat := fmt.Sprintf("%%-%ds\t%%s", urlMaxLen)
				for i, st := range sts {
					info[fmt.Sprintf("%s-%d", tag, i)] = fmt.Sprintf(urlFormat, st.URL, st.FuncName)
				}
			} else {
				strs, ok := v.([]string)
				if ok {
					tag := reps.Replace(rt.Field(i).Tag.Get("json"))
					l = len(tag)
					if maxlen < l {
						maxlen = l
					}
					info[tag] = fmt.Sprintf("%v", strs)
				}
			}
			continue
		}
		tag := reps.Replace(rt.Field(i).Tag.Get("json"))
		l = len(tag)
		if maxlen < l {
			maxlen = l
		}
		info[tag] = value
	}

	infoFormat := fmt.Sprintf("%%-%ds ->\t%%s", maxlen)
	strs := make([]string, 0, rt.NumField())
	for tag, value := range info {
		if len(tag) != 0 && len(value) != 0 {
			strs = append(strs, fmt.Sprintf(infoFormat, tag, value))
		}
	}
	sort.Strings(strs)
	return "\n" + strings.Join(strs, "\n")
}

// Get returns parased Detail object.
func (i info) Get() Detail {
	i.prepare()
	valdRepo := fmt.Sprintf("github.com/%s/%s", Organization, Repository)
	defaultURL := fmt.Sprintf("https://%s/tree/%s", valdRepo, i.detail.GitCommit)

	i.detail.StackTrace = make([]StackTrace, 0, 10)
	for j := 3; ; j++ {
		pc, file, line, ok := i.rtCaller(j)
		if !ok {
			break
		}
		funcName := i.rtFuncForPC(pc).Name()
		if funcName == "runtime.main" {
			break
		}
		url := defaultURL
		switch {
		case strings.HasPrefix(file, runtime.GOROOT()+"/src"):
			url = fmt.Sprintf("https://github.com/golang/go/blob/%s%s#L%d", i.detail.GoVersion, strings.TrimPrefix(file, runtime.GOROOT()), line)
		case strings.Contains(file, "go/pkg/mod/"):
			url = "https:/"
			for _, path := range strings.Split(strings.SplitN(file, "go/pkg/mod/", 2)[1], "/") {
				if strings.Contains(path, "@") {
					sv := strings.SplitN(path, "@", 2)
					if strings.Count(sv[1], "-") > 2 {
						path = sv[0] + "/blob/master"
					} else {
						path = sv[0] + "/blob/" + sv[1]
					}
				}
				url += "/" + path
			}
			url += "#L" + strconv.Itoa(line)
		case strings.Contains(file, "go/src/") && strings.Contains(file, valdRepo):
			url = strings.Replace(strings.SplitN(file, "go/src/", 2)[1]+"#L"+strconv.Itoa(line), valdRepo, "https://"+valdRepo+"/blob/"+i.detail.GitCommit, -1)
		}
		i.detail.StackTrace = append(i.detail.StackTrace, StackTrace{
			FuncName: funcName,
			File:     file,
			Line:     line,
			URL:      url,
		})
	}
	return i.detail
}

func (i *info) prepare() {
	i.prepOnce.Do(func() {
		if len(i.detail.GitCommit) == 0 {
			i.detail.GitCommit = "master"
		}
		if len(Version) == 0 && len(i.detail.Version) == 0 {
			i.detail.Version = GitCommit
		}
		if len(i.detail.BuildTime) == 0 {
			i.detail.BuildTime = BuildTime
		}
		if len(i.detail.GoVersion) == 0 {
			i.detail.GoVersion = runtime.Version()
		}
		if len(i.detail.GoOS) == 0 {
			i.detail.GoOS = runtime.GOOS
		}
		if len(i.detail.GoArch) == 0 {
			i.detail.GoArch = runtime.GOARCH
		}
		if len(i.detail.CGOEnabled) == 0 && len(CGOEnabled) != 0 {
			i.detail.CGOEnabled = CGOEnabled
		}
		if len(i.detail.NGTVersion) == 0 && len(NGTVersion) != 0 {
			i.detail.NGTVersion = NGTVersion
		}
		if len(i.detail.BuildCPUInfoFlags) == 0 && len(BuildCPUInfoFlags) != 0 {
			i.detail.BuildCPUInfoFlags = strings.Split(strings.TrimSpace(BuildCPUInfoFlags), " ")
		}
	})
}

// String calls String method of global detail object.
func String() string {
	return detail.String()
}

// Get calls Get method of global detail object.
func Get() Detail {
	return detail.Get()
}

func (s StackTrace) String() string {
	return fmt.Sprintf("URL: %s\tFile: %s\tLine: #%d\tFuncName: %s", s.URL, s.File, s.Line, s.FuncName)
}
