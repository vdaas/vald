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

// Package info provides build-time info
package info

import (
	"fmt"
	"reflect"
	"runtime"
	"sort"
	"strconv"
	"sync"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/strings"
)

// Info represents an interface to get the runtime information.
type Info interface {
	String() string
	Get() Detail
}

type info struct {
	baseURL  string // e.g https://github.com/vdaas/vald/tree/main
	detail   Detail
	prepOnce sync.Once

	// runtime functions
	rtCaller    func(skip int) (pc uintptr, file string, line int, ok bool)
	rtFuncForPC func(pc uintptr) *runtime.Func
}

// Detail represents environment information of system and stacktrace information.
type Detail struct {
	Version           string       `json:"vald_version,omitempty"         yaml:"vald_version,omitempty"`
	ServerName        string       `json:"server_name,omitempty"          yaml:"server_name,omitempty"`
	GitCommit         string       `json:"git_commit,omitempty"           yaml:"git_commit,omitempty"`
	BuildTime         string       `json:"build_time,omitempty"           yaml:"build_time,omitempty"`
	GoVersion         string       `json:"go_version,omitempty"           yaml:"go_version,omitempty"`
	GoOS              string       `json:"go_os,omitempty"                yaml:"go_os,omitempty"`
	GoArch            string       `json:"go_arch,omitempty"              yaml:"go_arch,omitempty"`
	GoRoot            string       `json:"go_root,omitempty"              yaml:"go_root,omitempty"`
	CGOEnabled        string       `json:"cgo_enabled,omitempty"          yaml:"cgo_enabled,omitempty"`
	NGTVersion        string       `json:"ngt_version,omitempty"          yaml:"ngt_version,omitempty"`
	BuildCPUInfoFlags []string     `json:"build_cpu_info_flags,omitempty" yaml:"build_cpu_info_flags,omitempty"`
	StackTrace        []StackTrace `json:"stack_trace,omitempty"          yaml:"stack_trace,omitempty"`
}

// StackTrace represents stacktrace information about url, function name, file, line ..etc.
type StackTrace struct {
	URL      string `json:"url,omitempty"           yaml:"url,omitempty"`
	FuncName string `json:"function_name,omitempty" yaml:"func_name,omitempty"`
	File     string `json:"file,omitempty"          yaml:"file,omitempty"`
	Line     int    `json:"line,omitempty"          yaml:"line,omitempty"`
}

var (
	// injected from build script

	// Version represent Vald version.
	Version = "v0.0.1"
	// GitCommit represent the Vald GitCommit
	GitCommit = "main"
	// BuildTime represent the Vald Build time.
	BuildTime = ""
	// GoVersion represent the golang version to build Vald.
	GoVersion string
	// GoOS represent the OS version of golang to build Vald.
	GoOS string
	// GoArch represent the architecture target to build Vald.
	GoArch string
	// GoRoot represent the root of the Go tree.
	GoRoot string
	// CGOEnabled represent the cgo is enable or not to build Vald.
	CGOEnabled string
	// NGTVersion represent the NGT version in Vald.
	NGTVersion string
	// BuildCPUInfoFlags represent the CPU info flags to build Vald.
	BuildCPUInfoFlags string

	// Organization represent the organization of Vald.
	Organization = "vdaas"
	// Repository represent the repository of Vald.
	Repository = "vald"

	reps = strings.NewReplacer("_", " ", ",omitempty", "")

	once         sync.Once
	infoProvider Info

	rt = reflect.TypeOf(Detail{})

	rtNumField = rt.NumField()
	valdRepo   = fmt.Sprintf("github.com/%s/%s", Organization, Repository)
)

const (
	goSrc    = "go/src/"
	goSrcLen = len(goSrc)
	goMod    = "go/pkg/mod/"
	goModLen = len(goMod)
)

// Init initializes Detail object only once.
func Init(name string) {
	once.Do(func() {
		i, err := New(WithServerName(name))
		if err != nil {
			log.Init()
			// skipcq: RVV-A0003
			log.Fatal(errors.ErrFailedToInitInfo(err))
		}
		infoProvider = i
	})
}

// New initializes and returns the info object or any error occurred.
func New(opts ...Option) (Info, error) {
	i := &info{
		detail: Detail{
			ServerName:        "",
			Version:           Version,
			GitCommit:         GitCommit,
			BuildTime:         BuildTime,
			GoVersion:         GoVersion,
			GoOS:              GoOS,
			GoArch:            GoArch,
			GoRoot:            GoRoot,
			CGOEnabled:        CGOEnabled,
			NGTVersion:        NGTVersion,
			BuildCPUInfoFlags: strings.Split(strings.TrimSpace(BuildCPUInfoFlags), " "),
			StackTrace:        nil,
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

	if i.rtCaller == nil || i.rtFuncForPC == nil {
		return nil, errors.ErrRuntimeFuncNil()
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
// The stacktrace will be initialized when the stacktrace is not initialized yet.
func (i *info) String() string {
	if len(i.detail.StackTrace) == 0 {
		i.detail = i.Get()
	}

	return i.detail.String()
}

// String returns summary of Detail object.
func (d Detail) String() string {
	d.Version = log.Bold(d.Version)
	maxlen, l := 0, 0
	rv := reflect.ValueOf(d)
	info := make(map[string]string, rtNumField)
	for i := 0; i < rtNumField; i++ {
		rtField := rt.Field(i)
		v := rv.Field(i).Interface()
		value, ok := v.(string)
		if !ok {
			sts, ok := v.([]StackTrace)
			if ok {
				tag := reps.Replace(rtField.Tag.Get("json"))
				l = len(tag) + 2
				if maxlen < l {
					maxlen = l
				}
				urlMaxLen := 0
				fileMaxLen := 0
				for _, st := range sts {
					ul := len(st.URL)
					fl := len(st.File + "#L" + strconv.Itoa(st.Line))
					if urlMaxLen < ul {
						urlMaxLen = ul
					}
					if fileMaxLen < fl {
						fileMaxLen = fl
					}
				}
				urlFormat := fmt.Sprintf("%%-%ds\t%%-%ds\t", urlMaxLen, fileMaxLen)
				for i, st := range sts {
					info[fmt.Sprintf("%s-%03d", tag, i)] = fmt.Sprintf(urlFormat, st.URL, st.File+"#L"+strconv.Itoa(st.Line)) + st.FuncName
				}
			} else {
				strs, ok := v.([]string)
				if ok {
					tag := reps.Replace(rtField.Tag.Get("json"))
					l = len(tag)
					if maxlen < l {
						maxlen = l
					}
					info[tag] = fmt.Sprintf("%v", strs)
				}
			}
			continue
		}
		tag := reps.Replace(rtField.Tag.Get("json"))
		l = len(tag)
		if maxlen < l {
			maxlen = l
		}
		info[tag] = value
	}

	infoFormat := "%-" + strconv.Itoa(maxlen) + "s ->\t"
	strs := make([]string, 0, rtNumField)
	for tag, value := range info {
		if len(tag) != 0 && len(value) != 0 {
			strs = append(strs, fmt.Sprintf(infoFormat, tag)+value)
		}
	}
	sort.Strings(strs)
	return "\n" + strings.Join(strs, "\n")
}

// Get returns parased Detail object.
func (i *info) Get() Detail {
	i.prepare()
	return i.get()
}

func (i info) get() Detail {
	i.detail.StackTrace = make([]StackTrace, 0, 10)
	for j := 2; ; j++ {
		pc, file, line, ok := i.rtCaller(j)
		if !ok {
			break
		}
		funcName := i.rtFuncForPC(pc).Name()
		if funcName == "runtime.main" {
			break
		}
		url := i.baseURL
		var idx int
		switch {
		case strings.HasPrefix(file, i.detail.GoRoot+"/src"):
			url = "https://github.com/golang/go/blob/" + i.detail.GoVersion + strings.TrimPrefix(file, i.detail.GoRoot) + "#L" + strconv.Itoa(line)
		case func() bool {
			idx = strings.Index(file, goMod)
			return idx >= 0
		}():
			url = "https:/"
			for _, path := range strings.Split(file[idx+goModLen:], "/") {
				left, right, ok := strings.Cut(path, "@")
				if ok {
					if strings.Count(right, "-") > 2 {
						path = left + "/blob/main"
					} else {
						path = left + "/blob/" + right
					}
				}
				url += "/" + path
			}
			url += "#L" + strconv.Itoa(line)
		case func() bool {
			idx = strings.Index(file, goSrc)
			return idx >= 0 && strings.Contains(file, valdRepo)
		}():
			url = strings.Replace(file[idx+goSrcLen:]+"#L"+strconv.Itoa(line), valdRepo, "https://"+valdRepo+"/blob/"+i.detail.GitCommit, -1)
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
			i.detail.GitCommit = "main"
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
		if len(i.detail.GoRoot) == 0 {
			i.detail.GoRoot = runtime.GOROOT()
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
		if len(i.baseURL) == 0 {
			i.baseURL = "https://" + valdRepo + "/tree/" + i.detail.GitCommit
		}
	})
}

func (s StackTrace) String() string {
	return "URL: " + s.URL + "\tFile: " + s.File + "\tLine: #" + strconv.Itoa(s.Line) + "\tFuncName: " + s.FuncName
}
