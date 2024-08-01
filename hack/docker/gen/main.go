//
// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
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

package main

import (
	"bytes"
	"context"
	"fmt"
	"io/fs"
	"os"
	"os/signal"
	"regexp"
	"slices"
	"syscall"
	"text/template"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/file"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/internal/sync/errgroup"
)

var tmpl = fmt.Sprintf(`# syntax = docker/dockerfile:latest
#
# Copyright (C) 2019-{{.Year}} {{.Maintainer}}
#
# Licensed under the Apache License, Version 2.0 (the "License");
# You may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

# DO_NOT_EDIT this Dockerfile is generated by hack/docker/gen/main.go

ARG UPX_OPTIONS=-9

{{- range $key, $value := .Arguments }}
ARG {{$key}}={{$value}}
{{- end}}
{{- range $image := .ExtraImages }}
# skipcq: DOK-DL3026,DOK-DL3007
FROM {{$image}}
{{- end}}
# skipcq: DOK-DL3026,DOK-DL3007
FROM {{.BuilderImage}}:{{.BuilderTag}}{{if and (not (eq (ContainerName .ContainerType) "%s")) (not (eq (ContainerName .ContainerType) "%s"))}} AS builder {{- end}}
ARG MAINTAINER="{{.Maintainer}}"
LABEL maintainer="${MAINTAINER}"

# skipcq: DOK-DL3002
USER {{.BuildUser}}

ARG TARGETARCH
ARG TARGETOS
ARG GO_VERSION
ARG RUST_VERSION

{{- range $keyValue := .EnvironmentsSlice }}
ENV {{$keyValue}}
{{- end}}

WORKDIR {{.RootDir}}/${ORG}/${REPO}
{{- range $files := .ExtraCopies }}
COPY {{$files}}
{{- end}}
SHELL ["/bin/bash", "-o", "pipefail", "-c"]
#skipcq: DOK-W1001, DOK-SC2046, DOK-SC2086, DOK-DL3008
RUN {{RunMounts .RunMounts}}\
    echo 'Binary::apt::APT::Keep-Downloaded-Packages "true";' > /etc/apt/apt.conf.d/keep-cache \
    && echo 'APT::Install-Recommends "false";' > /etc/apt/apt.conf.d/no-install-recommends \
    && apt-get clean \
    && apt-get update -y \
    && apt-get upgrade -y \
{{- if eq (ContainerName .ContainerType) "%s"}}
    && apt-get install -y --no-install-recommends --fix-missing \
    curl \
    gnupg \
    software-properties-common \
    && add-apt-repository ppa:ubuntu-toolchain-r/test -y \
    && apt-get update -y \
    && apt-get upgrade -y \
{{- end}}
    && apt-get install -y --no-install-recommends --fix-missing \
    build-essential \
    ca-certificates \
{{- if not (eq (ContainerName .ContainerType) "%s")}}
    curl \
{{- end}}
    tzdata \
    locales \
    git \
{{- range $epkg := .ExtraPackages }}
    {{$epkg}} \
{{- end}}
    && ldconfig \
    && echo "${LANG} UTF-8" > /etc/locale.gen \
    && ln -fs /usr/share/zoneinfo/${TZ} /etc/localtime \
    && locale-gen ${LANGUAGE} \
    && update-locale LANG=${LANGUAGE} \
    && dpkg-reconfigure -f noninteractive tzdata \
    && apt-get clean \
    && apt-get autoclean -y \
    && apt-get autoremove -y \
    && {{RunCommands .RunCommands}}

{{- if and (not (eq (ContainerName .ContainerType) "%s")) (not (eq (ContainerName .ContainerType) "%s"))}}
# skipcq: DOK-DL3026,DOK-DL3007
FROM {{.RuntimeImage}}:{{.RuntimeTag}}
ARG MAINTAINER="{{.Maintainer}}"
LABEL maintainer="${MAINTAINER}"

ENV APP_NAME={{.AppName}}

COPY --from=builder {{.BinDir}}/${APP_NAME} {{.BinDir}}/${APP_NAME}
{{- if .ConfigExists }}
COPY cmd/{{.PackageDir}}/sample.yaml /etc/server/config.yaml
{{- end}}
{{- range $from, $file := .StageFiles }}
COPY --from=builder {{$file}} {{$file}}
{{- end}}
{{- end}}
# skipcq: DOK-DL3002
USER {{.RuntimeUser}}

{{- if .Entrypoints}}
ENTRYPOINT [{{Entrypoint .Entrypoints}}]
{{- else if and (not (eq (ContainerName .ContainerType) "%s")) (not (eq (ContainerName .ContainerType) "%s"))}}
ENTRYPOINT ["{{.BinDir}}/{{.AppName}}"]
{{- end}}`, DevContainer.String(), CIContainer.String(),
	DevContainer.String(),
	DevContainer.String(),
	DevContainer.String(), CIContainer.String(),
	DevContainer.String(), CIContainer.String())

var docker = template.Must(template.New("Dockerfile").Funcs(template.FuncMap{
	"RunCommands": func(commands []string) string {
		if len(commands) == 0 {
			return ""
		}
		var b strings.Builder
		for i, cmd := range commands {
			if i > 0 {
				b.WriteString(" \\\n    && ")
			}
			b.WriteString(cmd)
		}
		return b.String()
	},
	"RunMounts": func(commands []string) string {
		if len(commands) == 0 {
			return ""
		}
		var b strings.Builder
		for i, cmd := range commands {
			if i > 0 {
				b.WriteString(" \\\n    ")
			}
			b.WriteString(cmd)
		}
		return b.String()
	},

	"Entrypoint": func(entries []string) string {
		if len(entries) == 0 {
			return "\"{{.BinDir}}/{{.AppName}}\""
		}
		return "\"" + strings.Join(entries, "\", \"") + "\""
	},
	"ContainerName": func(c ContainerType) string {
		return c.String()
	},
}).Parse(tmpl))

type Data struct {
	ConfigExists      bool
	Year              int
	ContainerType     ContainerType
	AppName           string
	BinDir            string
	BuildUser         string
	BuilderImage      string
	BuilderTag        string
	Maintainer        string
	PackageDir        string
	RootDir           string
	RuntimeImage      string
	RuntimeTag        string
	RuntimeUser       string
	Arguments         map[string]string
	Environments      map[string]string
	Entrypoints       []string
	EnvironmentsSlice []string
	ExtraCopies       []string
	ExtraImages       []string
	ExtraPackages     []string
	Preprocess        []string
	RunCommands       []string
	RunMounts         []string
	StageFiles        []string
}

type ContainerType int

const (
	organization          = "vdaas"
	repository            = "vald"
	defaultBinaryDir      = "/usr/bin"
	defaultBuilderImage   = "ghcr.io/vdaas/vald/vald-buildbase"
	defaultBuilderTag     = "nightly"
	defaultLanguage       = "en_US.UTF-8"
	defaultMaintainer     = organization + ".org " + repository + " team <" + repository + "@" + organization + ".org>"
	defaultRuntimeImage   = "gcr.io/distroless/static"
	defaultRuntimeTag     = "nonroot"
	defaultRuntimeUser    = "nonroot:nonroot"
	defaultBuildUser      = "root:root"
	maintainerKey         = "MAINTAINER"
	minimumArgumentLength = 2
	ubuntuVersion         = "22.04"

	goWorkdir   = "${GOPATH}/src/github.com"
	rustWorkdir = "${HOME}/rust/src/github.com"

	agentInernalPackage = "pkg/agent/internal"

	ngtPreprocess   = "make ngt/install"
	faissPreprocess = "make faiss/install"

	helmOperatorRootdir   = "/opt/helm"
	helmOperatorWatchFile = helmOperatorRootdir + "/watches.yaml"
	helmOperatorChartsDir = helmOperatorRootdir + "/charts"
)

const (
	Go ContainerType = iota
	Rust
	DevContainer
	HelmOperator
	CIContainer
	Other
)

func (c ContainerType) String() string {
	return containerTypeName[c]
}

var (
	containerTypeName = map[ContainerType]string{
		Go:           "Go",
		Rust:         "Rust",
		DevContainer: "DevContainer",
		HelmOperator: "HelmOperator",
		CIContainer:  "CIContainer",
		Other:        "Other",
	}

	defaultEnvironments = map[string]string{
		"DEBIAN_FRONTEND": "noninteractive",
		"HOME":            "/root",
		"USER":            "root",
		"INITRD":          "No",
		"LANG":            defaultLanguage,
		"LANGUAGE":        defaultLanguage,
		"LC_ALL":          defaultLanguage,
		"ORG":             organization,
		"TZ":              "Etc/UTC",
		"PATH":            "${PATH}:/usr/local/bin",
		"REPO":            repository,
	}
	goDefaultEnvironments = map[string]string{
		"GOROOT":      "/opt/go",
		"GOPATH":      "/go",
		"GO111MODULE": "on",
		"PATH":        "${PATH}:${GOROOT}/bin:${GOPATH}/bin:/usr/local/bin",
	}
	rustDefaultEnvironments = map[string]string{
		"RUST_HOME":   "/usr/loacl/lib/rust",
		"RUSTUP_HOME": "${RUST_HOME}/rustup",
		"CARGO_HOME":  "${RUST_HOME}/cargo",
		"PATH":        "${PATH}:${RUSTUP_HOME}/bin:${CARGO_HOME}/bin:/usr/local/bin",
	}
	clangDefaultEnvironments = map[string]string{
		"CC":  "gcc",
		"CXX": "g++",
	}
	goInstallCommands = []string{
		"make GOPATH=\"${GOPATH}\" GOROOT=\"${GOROOT}\" GO_VERSION=\"${GO_VERSION}\" go/install",
		"make GOPATH=\"${GOPATH}\" GOROOT=\"${GOROOT}\" GO_VERSION=\"${GO_VERSION}\" go/download",
	}
	rustInstallCommands = []string{
		"make RUST_VERSION=\"${RUST_VERSION}\" rust/install",
	}
	goBuildCommands = []string{
		"make GOARCH=\"${TARGETARCH}\" GOOS=\"${TARGETOS}\" REPO=\"${ORG}\" NAME=\"${REPO}\" cmd/${PKG}/${APP_NAME}",
		"mv \"cmd/${PKG}/${APP_NAME}\" \"{{$.BinDir}}/${APP_NAME}\"",
	}
	rustBuildCommands = []string{
		"make rust/target/release/${APP_NAME}",
		"mv \"rust/target/release/${APP_NAME}\" \"{{$.BinDir}}/${APP_NAME}\"",
		"rm -rf rust/target",
	}

	defaultMounts = []string{
		"--mount=type=bind,target=.,rw",
		"--mount=type=tmpfs,target=/tmp",
		"--mount=type=cache,target=/var/lib/apt,sharing=locked",
		"--mount=type=cache,target=/var/cache/apt,sharing=locked",
	}

	goDefaultMounts = []string{
		"--mount=type=cache,target=\"${GOPATH}/pkg\",id=\"go-build-${TARGETARCH}\"",
		"--mount=type=cache,target=\"${HOME}/.cache/go-build\",id=\"go-build-${TARGETARCH}\"",
	}

	clangBuildDeps = []string{
		"cmake",
		"gcc",
		"g++",
		"unzip",
		"libssl-dev",
	}
	ngtBuildDeps = []string{
		"liblapack-dev",
		"libomp-dev",
		"libopenblas-dev",
	}
	faissBuildDeps = []string{
		"gfortran",
		"libquadmath0",
	}
	devContainerDeps = []string{
		"gawk",
		"gnupg2",
		"graphviz",
		"jq",
		"libhdf5-dev",
		"libaec-dev",
		"sed",
		"zip",
	}

	ciContainerPreprocess = []string{
		"make GOARCH=${TARGETARCH} GOOS=${TARGETOS} deps GO_CLEAN_DEPS=false",
		"make GOARCH=${TARGETARCH} GOOS=${TARGETOS} golangci-lint/install",
		"make GOARCH=${TARGETARCH} GOOS=${TARGETOS} gotestfmt/install",
		"make cmake/install",
		"make buf/install",
		"make hdf5/install",
		"make helm-docs/install",
		"make helm/install",
		"make k3d/install",
		"make k9s/install",
		"make kind/install",
		"make kubectl/install",
		"make kubelinter/install",
		"make reviewdog/install",
		"make tparse/install",
		"make valdcli/install",
		"make yq/install",
		"make minikube/install",
		"make stern/install",
		"make telepresence/install",
	}

	devContainerPreprocess = []string{
		"curl -fsSL https://deb.nodesource.com/setup_current.x | bash -",
		"apt-get clean",
		"apt-get update -y",
		"apt-get upgrade -y",
		"apt-get install -y --no-install-recommends --fix-missing nodejs",
		"npm install -g npm@latest",
		"apt-get clean",
		"apt-get autoclean -y",
		"apt-get autoremove -y",
		"make delve/install",
		"make gomodifytags/install",
		"make gopls/install",
		"make gotests/install",
		"make impl/install",
		"make staticcheck/install",
	}
)

func appendM[K comparable](maps ...map[K]string) map[K]string {
	if len(maps) == 0 {
		return nil
	}
	result := maps[0]
	for _, m := range maps[1:] {
		for k, v := range m {
			ev, ok := result[k]
			if ok {
				v += ":" + ev
			}
			result[k] = v
		}
	}

	for k, v := range result {
		vs := strings.Split(v, ":")
		slices.Sort(vs)
		v = strings.Join(slices.Compact(vs), ":")
		if strings.Contains(v, "${PATH}:") {
			v = strings.TrimPrefix(strings.ReplaceAll(strings.ReplaceAll(v, "${PATH}", ""), "::", ":")+":${PATH}", ":")
		}
		if strings.Contains(v, ":unix") {
			v = "unix:" + strings.TrimSuffix(v, ":unix")
		}
		result[k] = v
	}
	return result
}

var re = regexp.MustCompile(`\$\{?(\w+)\}?`)

func extractVariables(value string) []string {
	matches := re.FindAllStringSubmatch(value, -1)
	vars := make([]string, 0, len(matches))
	for _, match := range matches {
		vars = append(vars, match[1])
	}
	return vars
}

func topologicalSort(envMap map[string]string) []string {
	// Graph structures
	inDegree := make(map[string]int)
	graph := make(map[string][]string)

	// Initialize the graph
	for key, value := range envMap {
		vars := extractVariables(value)
		for _, refKey := range vars {
			if refKey != key {
				graph[refKey] = append(graph[refKey], key)
				inDegree[key]++
			}
		}
	}

	queue := make([]string, 0, len(envMap)-len(graph))
	for key := range envMap {
		if inDegree[key] == 0 {
			queue = append(queue, key)
		}
	}

	slices.Sort(queue)

	// Topological sort
	result := make([]string, 0, len(envMap))
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		if value, exists := envMap[node]; exists {
			result = append(result, node+"="+value)
		}
		for _, neighbor := range graph[node] {
			inDegree[neighbor]--
			if inDegree[neighbor] == 0 {
				queue = append(queue, neighbor)
			}
		}
	}

	return result
}

func main() {
	log.Init()
	if len(os.Args) < minimumArgumentLength {
		// skipcq: RVV-A0003
		log.Fatal(errors.New("invalid argument"))
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGHUP,
		syscall.SIGALRM,
		syscall.SIGKILL,
		syscall.SIGTERM)
	defer cancel()
	log.Debug(tmpl)

	maintainer := os.Getenv(maintainerKey)
	if maintainer == "" {
		maintainer = defaultMaintainer
	}
	year := time.Now().Year()
	eg, egctx := errgroup.New(ctx)
	for n, d := range map[string]Data{
		"vald-agent-ngt": {
			AppName:       "ngt",
			PackageDir:    "agent/core/ngt",
			ExtraPackages: append(clangBuildDeps, ngtBuildDeps...),
			Preprocess:    []string{ngtPreprocess},
		},
		"vald-agent-faiss": {
			AppName:    "faiss",
			PackageDir: "agent/core/faiss",
			ExtraPackages: append(clangBuildDeps,
				append(ngtBuildDeps,
					faissBuildDeps...)...),
			Preprocess: []string{faissPreprocess},
		},
		"vald-agent": {
			AppName:       "agent",
			PackageDir:    "agent/core/agent",
			ContainerType: Rust,
			RuntimeImage:  "gcr.io/distroless/cc-debian12",
			ExtraPackages: append(clangBuildDeps,
				append(ngtBuildDeps,
					faissBuildDeps...)...),
			Preprocess: []string{
				ngtPreprocess,
				faissPreprocess,
			},
		},
		"vald-agent-sidecar": {
			AppName:    "sidecar",
			PackageDir: "agent/sidecar",
		},
		"vald-discoverer-k8s": {
			AppName:    "discoverer",
			PackageDir: "discoverer/k8s",
		},
		"vald-gateway-lb": {
			AppName:    "lb",
			PackageDir: "gateway/lb",
		},
		"vald-gateway-filter": {
			AppName:    "filter",
			PackageDir: "gateway/filter",
		},
		"vald-gateway-mirror": {
			AppName:    "mirror",
			PackageDir: "gateway/mirror",
		},
		"vald-manager-index": {
			AppName:    "index",
			PackageDir: "manager/index",
		},
		"vald-index-correction": {
			AppName:    "index-correction",
			PackageDir: "index/job/correction",
		},
		"vald-index-creation": {
			AppName:    "index-creation",
			PackageDir: "index/job/creation",
		},
		"vald-index-save": {
			AppName:    "index-save",
			PackageDir: "index/job/save",
		},
		"vald-readreplica-rotate": {
			AppName:    "readreplica-rotate",
			PackageDir: "index/job/readreplica/rotate",
		},
		"vald-index-operator": {
			AppName:    "index-operator",
			PackageDir: "index/operator",
		},
		"vald-benchmark-job": {
			AppName:       "job",
			PackageDir:    "tools/benchmark/job",
			ExtraPackages: append(clangBuildDeps, "libhdf5-dev", "libaec-dev"),
			Preprocess: []string{
				"make hdf5/install",
			},
		},
		"vald-benchmark-operator": {
			AppName:    "operator",
			PackageDir: "tools/benchmark/operator",
		},
		"vald-helm-operator": {
			AppName:       "helm-operator",
			PackageDir:    "operator/helm",
			ContainerType: HelmOperator,
			Arguments: map[string]string{
				"OPERATOR_SDK_VERSION": "latest",
			},
			ExtraCopies: []string{
				"--from=operator /usr/local/bin/${APP_NAME} {{$.BinDir}}/${APP_NAME}",
			},
			ExtraImages: []string{
				"quay.io/operator-framework/helm-operator:${OPERATOR_SDK_VERSION} AS operator",
			},
			ExtraPackages: []string{"upx"},
			Preprocess: []string{
				"mkdir -p " + helmOperatorChartsDir,
				`{ \
        echo "---"; \
        echo "- version: v1"; \
        echo "  group: vald.vdaas.org"; \
        echo "  kind: ValdRelease"; \
        echo "  chart: ` + helmOperatorChartsDir + `/vald"; \
        echo "- version: v1"; \
        echo "  group: vald.vdaas.org"; \
        echo "  kind: ValdHelmOperatorRelease"; \
        echo "  chart: ` + helmOperatorChartsDir + `/vald-helm-operator"; \
    } > ` + helmOperatorWatchFile,
				"make GOARCH=${TARGETARCH} GOOS=${TARGETOS} helm/schema/vald",
				"make GOARCH=${TARGETARCH} GOOS=${TARGETOS} helm/schema/vald-helm-operator",
				"cp -r charts/* " + helmOperatorChartsDir + "/",
				"upx \"{{$.BinDir}}/${APP_NAME}\"",
			},
			StageFiles: []string{
				helmOperatorWatchFile,
				helmOperatorChartsDir + "/vald",
				helmOperatorChartsDir + "/vald-helm-operator",
			},
			Entrypoints: []string{"{{$.BinDir}}/{{.AppName}}", "run", "--watches-file=" + helmOperatorWatchFile},
		},
		"vald-cli-loadtest": {
			AppName:       "loadtest",
			PackageDir:    "tools/cli/loadtest",
			ExtraPackages: append(clangBuildDeps, "libhdf5-dev", "libaec-dev"),
			Preprocess: []string{
				"make hdf5/install",
			},
		},
		"vald-ci-container": {
			AppName:       "ci-container",
			ContainerType: CIContainer,
			PackageDir:    "ci/base",
			RuntimeUser:   defaultBuildUser,
			ExtraPackages: append([]string{"npm"}, append(clangBuildDeps,
				append(ngtBuildDeps,
					append(faissBuildDeps,
						devContainerDeps...)...)...)...),
			Preprocess:  append(ciContainerPreprocess, ngtPreprocess, faissPreprocess),
			Entrypoints: []string{"/bin/bash"},
		},
		"vald-dev-container": {
			AppName:       "dev-container",
			BuilderImage:  "mcr.microsoft.com/devcontainers/base",
			BuilderTag:    "ubuntu" + ubuntuVersion,
			BuildUser:     defaultBuildUser,
			RuntimeUser:   defaultBuildUser,
			ContainerType: DevContainer,
			PackageDir:    "dev",
			ExtraPackages: append(clangBuildDeps,
				append(ngtBuildDeps,
					append(faissBuildDeps,
						devContainerDeps...)...)...),
			Preprocess: append(devContainerPreprocess,
				append(ciContainerPreprocess,
					ngtPreprocess,
					faissPreprocess)...),
		},
	} {
		name := n
		data := d

		eg.Go(safety.RecoverFunc(func() error {
			data.Maintainer = maintainer
			data.Year = year
			if data.BinDir == "" {
				data.BinDir = defaultBinaryDir
			}
			if data.RuntimeImage == "" {
				data.RuntimeImage = defaultRuntimeImage
			}
			if data.RuntimeTag == "" {
				data.RuntimeTag = defaultRuntimeTag
			}
			if data.BuilderImage == "" {
				data.BuilderImage = defaultBuilderImage
			}
			if data.BuilderTag == "" {
				data.BuilderTag = defaultBuilderTag
			}
			if data.RuntimeUser == "" {
				data.RuntimeUser = defaultRuntimeUser
			}

			if data.BuildUser == "" {
				data.BuildUser = defaultBuildUser
			}

			if data.Environments != nil {
				data.Environments = appendM(data.Environments, defaultEnvironments)
			} else {
				data.Environments = make(map[string]string, len(defaultEnvironments))
				data.Environments = appendM(data.Environments, defaultEnvironments)
			}
			switch data.ContainerType {
			case Go:
				data.Environments = appendM(data.Environments, goDefaultEnvironments)
				data.RootDir = goWorkdir
				commands := make([]string, 0, len(goInstallCommands)+len(data.Preprocess)+len(goBuildCommands))
				commands = append(commands, goInstallCommands...)
				if data.Preprocess != nil {
					commands = append(commands, data.Preprocess...)
				}
				if file.Exists(file.Join(os.Args[1], "cmd", data.PackageDir)) {
					commands = append(commands, goBuildCommands...)
				}
				data.RunCommands = commands
				mounts := make([]string, 0, len(defaultMounts)+len(goDefaultMounts))
				mounts = append(mounts, defaultMounts...)
				mounts = append(mounts, goDefaultMounts...)
				data.RunMounts = mounts
			case Rust:
				data.Environments = appendM(data.Environments, rustDefaultEnvironments)
				data.RootDir = rustWorkdir
				commands := make([]string, 0, len(rustInstallCommands)+len(data.Preprocess)+len(rustBuildCommands))
				commands = append(commands, rustInstallCommands...)
				if data.Preprocess != nil {
					commands = append(commands, data.Preprocess...)
				}
				commands = append(commands, rustBuildCommands...)
				data.RunCommands = commands
				data.RunMounts = defaultMounts
			case DevContainer, CIContainer:
				data.Environments = appendM(data.Environments, goDefaultEnvironments, rustDefaultEnvironments, clangDefaultEnvironments)
				data.RootDir = goWorkdir
				commands := make([]string, 0, len(goInstallCommands)+len(rustInstallCommands)+len(data.Preprocess)+1)
				commands = append(commands, append(goInstallCommands, rustInstallCommands...)...)
				if data.Preprocess != nil {
					commands = append(commands, data.Preprocess...)
				}
				commands = append(commands, "rm -rf {{.RootDir}}/${ORG}/${REPO}/*")
				data.RunCommands = commands
				mounts := make([]string, 0, len(defaultMounts)+len(goDefaultMounts))
				mounts = append(mounts, defaultMounts...)
				mounts = append(mounts, goDefaultMounts...)
				data.RunMounts = mounts
			case HelmOperator:
				data.Environments = appendM(data.Environments, goDefaultEnvironments)
				data.RootDir = goWorkdir
				commands := make([]string, 0, len(goInstallCommands)+len(data.Preprocess))
				commands = append(commands, goInstallCommands...)
				if data.Preprocess != nil {
					commands = append(commands, data.Preprocess...)
				}
				data.RunCommands = commands
				mounts := make([]string, 0, len(defaultMounts)+len(goDefaultMounts))
				mounts = append(mounts, defaultMounts...)
				mounts = append(mounts, goDefaultMounts...)
				data.RunMounts = mounts
			default:
				data.RootDir = "${HOME}"
				data.Environments["ROOTDIR"] = os.Args[1]
			}
			if strings.Contains(data.BuildUser, "root") {
				data.Environments["HOME"] = "/root"
				data.Environments["USER"] = "root"
			} else {
				user := data.BuildUser
				if strings.Contains(user, ":") {
					user = strings.SplitN(user, ":", 2)[0]
				}
				data.Environments["HOME"] = "/home/" + user
				data.Environments["USER"] = user
			}

			data.Environments["APP_NAME"] = data.AppName
			data.Environments["PKG"] = data.PackageDir
			data.EnvironmentsSlice = topologicalSort(data.Environments)
			data.ConfigExists = file.Exists(file.Join(os.Args[1], "cmd", data.PackageDir, "sample.yaml"))

			buf := bytes.NewBuffer(make([]byte, 0, len(tmpl)))
			log.Infof("Generating %s's Dockerfile", name)
			docker.Execute(buf, data)
			tpl := buf.String()
			buf.Reset()
			template.Must(template.New("Dockerfile").Parse(tpl)).Execute(buf, data)
			file.OverWriteFile(egctx, file.Join(os.Args[1], "dockers", data.PackageDir, "Dockerfile"), buf, fs.ModePerm)
			return nil
		}))
	}
	eg.Wait()
}
