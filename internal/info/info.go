//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
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

package info

import (
	"fmt"
	"reflect"
	"runtime"
	"slices"
	"strconv"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/internal/sync"
)

// Info represents an interface to get the runtime information.
type Info interface {
	String() string
	Get() Detail
}

type info struct {
	baseURL      string // e.g https://github.com/vdaas/vald/tree/main
	detail       Detail
	prepOnce     sync.Once
	valdReplacer *strings.Replacer

	// runtime functions
	rtCaller    func(skip int) (pc uintptr, file string, line int, ok bool)
	rtFuncForPC func(pc uintptr) *runtime.Func
}

// Detail represents environment information of system and stacktrace information.
type Detail struct {
	AlgorithmInfo     string       `json:"algorithm_info,omitempty"       yaml:"algorithm_info,omitempty"`
	BuildTime         string       `json:"build_time,omitempty"           yaml:"build_time,omitempty"`
	CGOCall           string       `json:"cgo_call,omitempty"             yaml:"cgo_call"`
	CGOEnabled        string       `json:"cgo_enabled,omitempty"          yaml:"cgo_enabled,omitempty"`
	GitCommit         string       `json:"git_commit,omitempty"           yaml:"git_commit,omitempty"`
	GoArch            string       `json:"go_arch,omitempty"              yaml:"go_arch,omitempty"`
	GoMaxProcs        string       `json:"go_max_procs,omitempty"         yaml:"go_max_procs,omitempty"`
	GoOS              string       `json:"go_os,omitempty"                yaml:"go_os,omitempty"`
	GoRoot            string       `json:"go_root,omitempty"              yaml:"go_root,omitempty"`
	GoVersion         string       `json:"go_version,omitempty"           yaml:"go_version,omitempty"`
	GoroutineCount    string       `json:"goroutine_count,omitempty"      yaml:"goroutine_count"`
	RuntimeCPUCores   string       `json:"runtime_cpu_cores,omitempty"    yaml:"runtime_cpu_cores,omitempty"`
	ServerName        string       `json:"server_name,omitempty"          yaml:"server_name,omitempty"`
	Version           string       `json:"vald_version,omitempty"         yaml:"vald_version,omitempty"`
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
	// injected from build script.

	// Version represent Vald version.
	Version = "v0.0.1"
	// GitCommit represent the Vald GitCommit.
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
	// AlgorithmInfo represent the NGT version in Vald.
	AlgorithmInfo string
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
	goSrc        = "go/src/"
	goSrcLen     = len(goSrc)
	goMod        = "go/pkg/mod/"
	goModLen     = len(goMod)
	cgoTrue      = "true"
	cgoFalse     = "false"
	cgoUnknown   = "unknown"
	googleGolang = "google.golang.org"
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
			AlgorithmInfo:     AlgorithmInfo,
			BuildCPUInfoFlags: strings.Split(strings.TrimSpace(BuildCPUInfoFlags), " "),
			BuildTime:         BuildTime,
			CGOEnabled:        CGOEnabled,
			GitCommit:         GitCommit,
			GoArch:            GoArch,
			GoOS:              GoOS,
			GoRoot:            GoRoot,
			GoVersion:         GoVersion,
			RuntimeCPUCores:   strconv.Itoa(runtime.NumCPU()),
			ServerName:        "",
			StackTrace:        nil,
			Version:           Version,
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
		return nil, errors.ErrRuntimeFuncNil
	}

	i.prepare()

	return i, nil
}

// String calls String method of global detail object.
func String() string {
	if infoProvider == nil {
		Init(log.Bold("WARNING: uninitialized info provider"))
	}
	return infoProvider.String()
}

// Get calls Get method of global detail object.
func Get() Detail {
	if infoProvider == nil {
		Init(log.Bold("WARNING: uninitialized info provider"))
	}
	return infoProvider.Get()
}

// String returns summary of Detail object.
// The stacktrace will be initialized when the stacktrace is not initialized yet.
func (i *info) String() string {
	if len(i.detail.StackTrace) == 0 {
		i.prepare()
		i.detail = i.getDetail()
	}

	return i.detail.String()
}

// String returns summary of Detail object.
// skipcq: RVV-B0006
func (d Detail) String() string {
	// skipcq: RVV-B0006
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
		switch tag {
		case "cgo_call":
			value = strconv.FormatInt(runtime.NumCgoCall(), 10)
		case "goroutine_count":
			value = strconv.Itoa(runtime.NumGoroutine())
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
	slices.Sort(strs)
	return "\n" + strings.Join(strs, "\n")
}

// Get returns parsed Detail object.
func (i *info) Get() Detail {
	i.prepare()
	return i.getDetail()
}

// skipcq: VET-V0008
func (i info) getDetail() Detail {
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
		index := strings.LastIndex(funcName, "/")
		if index != -1 {
			funcName = funcName[index+1:]
		}
		url := i.baseURL
		var idx int
		switch {
		case strings.HasPrefix(file, i.detail.GoRoot+"/src"):
			url = "https://github.com/golang/go/blob/" + i.detail.GoVersion + strings.TrimPrefix(file, i.detail.GoRoot)
		case strings.HasPrefix(file, "runtime"):
			url = "https://github.com/golang/go/blob/" + i.detail.GoVersion + "/src/" + file
		case strings.HasPrefix(file, googleGolang+"/grpc"):
			// google.golang.org/grpc@v1.65.0/server.go to https://github.com/grpc/grpc-go/blob/v1.65.0/server.go
			url = "https://github.com/grpc/grpc-go/blob/"
			_, versionSource, ok := strings.Cut(file, "@")
			if ok && versionSource != "" {
				url += versionSource
			} else {
				url = strings.ReplaceAll(file, googleGolang+"/grpc@", url)
			}
		case strings.HasPrefix(file, googleGolang+"/protobuf"):
			// google.golang.org/protobuf@v1.34.0/proto/decode.go to https://github.com/protocolbuffers/protobuf-go/blob/v1.34.0/proto/decode.go
			url = "https://github.com/protocolbuffers/protobuf-go/blob/"
			_, versionSource, ok := strings.Cut(file, "@")
			if ok && versionSource != "" {
				url += versionSource
			} else {
				url = strings.ReplaceAll(file, googleGolang+"/protobuf@", url)
			}
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
		case func() bool {
			idx = strings.Index(file, goSrc)
			return idx >= 0 && strings.Index(file, valdRepo) >= 0
		}():
			url = i.valdReplacer.Replace(file[idx+goSrcLen:])
		case strings.HasPrefix(file, valdRepo):
			url = i.valdReplacer.Replace(file)
		}
		url += "#L" + strconv.Itoa(line)

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
		if i.detail.GitCommit == "" {
			i.detail.GitCommit = "main"
		}
		if Version == "" && i.detail.Version == "" {
			i.detail.Version = GitCommit
		}
		if i.detail.BuildTime == "" {
			i.detail.BuildTime = BuildTime
		}
		if i.detail.GoVersion == "" {
			i.detail.GoVersion = runtime.Version()
		}
		if i.detail.GoOS == "" {
			i.detail.GoOS = runtime.GOOS
		}
		if i.detail.GoArch == "" {
			i.detail.GoArch = runtime.GOARCH
		}
		if i.detail.GoRoot == "" {
			i.detail.GoRoot = runtime.GOROOT()
		}
		if i.detail.CGOEnabled == "" && CGOEnabled != "" {
			i.detail.CGOEnabled = CGOEnabled
		}
		switch CGOEnabled {
		case "0", cgoFalse:
			i.detail.CGOEnabled = cgoFalse
		case "1", cgoTrue:
			i.detail.CGOEnabled = cgoTrue
		default:
			i.detail.CGOEnabled = cgoUnknown
		}
		if i.detail.AlgorithmInfo == "" && AlgorithmInfo != "" {
			i.detail.AlgorithmInfo = AlgorithmInfo
		}
		if len(i.detail.BuildCPUInfoFlags) == 0 && BuildCPUInfoFlags != "" {
			i.detail.BuildCPUInfoFlags = strings.Split(strings.TrimSpace(BuildCPUInfoFlags), " ")
		}
		if i.baseURL == "" {
			i.baseURL = "https://" + valdRepo + "/tree/" + i.detail.GitCommit
		}
		if len(i.detail.GoMaxProcs) == 0 {
			i.detail.GoMaxProcs = strconv.Itoa(runtime.GOMAXPROCS(-1))
		}
		if len(i.detail.CGOCall) == 0 {
			i.detail.CGOCall = strconv.FormatInt(runtime.NumCgoCall(), 10)
		}
		if len(i.detail.GoroutineCount) == 0 {
			i.detail.GoroutineCount = strconv.Itoa(runtime.NumGoroutine())
		}
		if i.valdReplacer == nil {
			i.valdReplacer = strings.NewReplacer(valdRepo, "https://"+valdRepo+"/blob/"+i.detail.GitCommit)
		}
	})
}

func (s StackTrace) String() string {
	return "URL: " + s.URL + "\tFile: " + s.File + "\tLine: #" + strconv.Itoa(s.Line) + "\tFuncName: " + s.FuncName
}

func (s StackTrace) ShortString() string {
	return s.URL + " " + s.FuncName
}
