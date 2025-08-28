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

package main

import (
	"bytes"
	"context"
	"fmt"
	"io/fs"
	"os"
	"os/signal"
	"slices"
	"syscall"
	"text/template"
	"time"

	"github.com/vdaas/vald/internal/conv"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/file"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"golang.org/x/tools/go/packages"
	"gopkg.in/yaml.v2"
)

const (
	agent               = "agent"
	agentFaiss          = agent + "-faiss"
	agentNGT            = agent + "-ngt"
	agentSidecar        = agent + "-sidecar"
	benchJob            = "benchmark-job"
	benchOperator       = "benchmark-operator"
	binfmt              = "binfmt"
	buildbase           = "buildbase"
	buildkit            = "buildkit"
	buildkitSyftScanner = buildkit + "-syft-scanner"
	ciContainer         = "ci-container"
	devContainer        = "dev-container"
	exampleContainer    = "example-client"
	discovererK8s       = "discoverer-k8s"
	gateway             = "gateway"
	gatewayFilter       = gateway + "-filter"
	gatewayLb           = gateway + "-lb"
	gatewayMirror       = gateway + "-mirror"
	helmOperator        = "helm-operator"
	indexCorrection     = "index-correction"
	indexCreation       = "index-creation"
	indexDeletion       = "index-deletion"
	indexOperator       = "index-operator"
	indexSave           = "index-save"
	managerIndex        = "manager-index"
	readreplicaRotate   = "readreplica-rotate"
	e2e                 = "e2e"

	organization          = "vdaas"
	repository            = "vald"
	defaultBinaryDir      = "/usr/bin"
	usrLocal              = "/usr/local"
	usrLocalBinaryDir     = usrLocal + "/bin"
	usrLocalLibDir        = usrLocal + "/lib"
	defaultBuilderImage   = "ghcr.io/" + organization + "/" + repository + "/" + repository + "-" + buildbase
	defaultBuilderTag     = "nightly"
	defaultLanguage       = "en_US.UTF-8"
	defaultMaintainer     = organization + ".org " + repository + " team <" + repository + "@" + organization + ".org>"
	defaultRuntimeImage   = "gcr.io/distroless/static"
	nonrootUser           = "nonroot"
	rootUser              = "root"
	defaultRuntimeTag     = nonrootUser
	defaultRuntimeUser    = nonrootUser + ":" + nonrootUser
	defaultBuildUser      = rootUser + ":" + rootUser
	defaultBuildStageName = "builder"
	maintainerKey         = "MAINTAINER"
	minimumArgumentLength = 2
	ubuntuVersion         = "24.04"

	goWorkdir   = "${GOPATH}/src/github.com"
	rustWorkdir = "${HOME}/rust/src/github.com"

	agentInernalPackage = "pkg/agent/internal"

	ngtPreprocess     = "make ngt/install"
	faissPreprocess   = "make faiss/install"
	usearchPreprocess = "make usearch/install"

	helmOperatorRootdir   = "/opt/helm"
	helmOperatorWatchFile = helmOperatorRootdir + "/watches.yaml"
	helmOperatorChartsDir = helmOperatorRootdir + "/charts"

	apisProtoPath = "apis/proto/**"

	hackPath = "hack/**"

	chartsValdPath            = "charts/vald"
	helmOperatorPath          = chartsValdPath + "-helm-operator"
	chartPath                 = chartsValdPath + "/Chart.yaml"
	valuesPath                = chartsValdPath + "/values.yaml"
	templatesPath             = chartsValdPath + "/templates/**"
	helmOperatorChartPath     = helmOperatorPath + "/Chart.yaml"
	helmOperatorValuesPath    = helmOperatorPath + "/values.yaml"
	helmOperatorTemplatesPath = helmOperatorPath + "/templates/**"

	goModPath = "go.mod"
	goSumPath = "go.sum"

	cargoLockPath       = "rust/Cargo.lock"
	cargoTomlPath       = "rust/Cargo.toml"
	rustBinAgentDirPath = "rust/bin/agent"
	rustNgtRsPath       = "rust/libs/ngt-rs/**"
	rustNgtPath         = "rust/libs/ngt/**"
	rustProtoPath       = "rust/libs/proto/**"

	excludeTestFilesPath = "!**/*_test.go"
	excludeMockFilesPath = "!**/*_mock.go"

	versionsPath           = "versions"
	operatorSDKVersionPath = versionsPath + "/OPERATOR_SDK_VERSION"
	goVersionPath          = versionsPath + "/GO_VERSION"
	rustVersionPath        = versionsPath + "/RUST_VERSION"
	faissVersionPath       = versionsPath + "/FAISS_VERSION"
	ngtVersionPath         = versionsPath + "/NGT_VERSION"
	usearchVersionPath     = versionsPath + "/USEARCH_VERSION"

	makefilePath    = "Makefile"
	makefileDirPath = makefilePath + ".d/**"

	amd64Platform  = "linux/amd64"
	arm64Platform  = "linux/arm64"
	multiPlatforms = amd64Platform + "," + arm64Platform

	header = `#
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
#`
)

var license = template.Must(template.New("license").Parse(header + `

# DO_NOT_EDIT this workflow file is generated by https://github.com/vdaas/vald/blob/main/hack/docker/gen/main.go

`))

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
}).Parse(fmt.Sprintf(`# syntax = docker/dockerfile:latest
# check=error=true
%s

# DO_NOT_EDIT this Dockerfile is generated by https://github.com/vdaas/vald/blob/main/hack/docker/gen/main.go

{{- if .AliasImage }}
FROM {{.BuilderImage}}:{{.BuilderTag}} AS {{.BuildStageName}}
{{- else}}
ARG UPX_OPTIONS=-9

{{- range $key, $value := .Arguments }}
ARG {{$key}}={{$value}}
{{- end}}
{{- range $image := .ExtraImages }}
# skipcq: DOK-DL3026,DOK-DL3007
FROM {{$image}}
{{- end}}
# skipcq: DOK-DL3026,DOK-DL3007
FROM {{.BuilderImage}}:{{.BuilderTag}}{{if and (not (eq (ContainerName .ContainerType) "%s")) (not (eq (ContainerName .ContainerType) "%s"))}} AS {{.BuildStageName}} {{- end}}
LABEL maintainer="{{.Maintainer}}"
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
RUN {{RunMounts .RunMounts}} \
    set -ex \
    && rm -f /etc/apt/apt.conf.d/docker-clean \
    && echo 'APT::Keep-Downloaded-Packages "true";' > /etc/apt/apt.conf.d/keep-cache \
    && echo 'APT::Install-Recommends "false";' > /etc/apt/apt.conf.d/no-install-recommends \
    && apt-get update -y \
    && apt-get install -y --no-install-recommends --fix-missing \
    build-essential \
    ca-certificates \
    curl \
{{- if eq (ContainerName .ContainerType) "%s"}}
    gnupg \
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
    && {{RunCommands .RunCommands}}
{{- if and (not (eq (ContainerName .ContainerType) "%s")) (not (eq (ContainerName .ContainerType) "%s"))}}
# skipcq: DOK-DL3026,DOK-DL3007
FROM {{.RuntimeImage}}:{{.RuntimeTag}}
LABEL maintainer="{{.Maintainer}}"
COPY --from=builder {{.BinDir}}/{{.AppName}} {{.BinDir}}/{{.AppName}}
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
{{- end}}
{{- end}}`, header, DevContainer.String(), CIContainer.String(),
	DevContainer.String(),
	DevContainer.String(), CIContainer.String(),
	DevContainer.String(), CIContainer.String())))

type (
	Workflow struct {
		Name string `yaml:"name"`
		On   On     `yaml:"on"`
		Jobs Jobs   `yaml:"jobs"`
	}

	On struct {
		Schedule          Schedule    `yaml:"schedule,omitempty"`
		Push              Push        `yaml:"push"`
		PullRequest       PullRequest `yaml:"pull_request"`
		PullRequestTarget PullRequest `yaml:"pull_request_target"`
	}

	Schedule []struct {
		Cron string `yaml:"cron,omitempty"`
	}

	Push struct {
		Branches []string `yaml:"branches"`
		Tags     []string `yaml:"tags"`
	}

	PullRequest struct {
		Types Types `yaml:"types,omitempty"`
		Paths Paths `yaml:"paths"`
	}

	Jobs struct {
		Build Build `yaml:"build"`
	}

	Build struct {
		Uses    string `yaml:"uses"`
		With    With   `yaml:"with"`
		Secrets string `yaml:"secrets"`
	}

	With struct {
		Target    string `yaml:"target"`
		Platforms string `yaml:"platforms,omitempty"`
	}

	Types []string
	Paths []string

	Data struct {
		AliasImage        bool
		ConfigExists      bool
		Year              int
		ContainerType     ContainerType
		AppName           string
		BinDir            string
		BuildPlatforms    string
		BuildStageName    string
		BuildUser         string
		BuilderImage      string
		BuilderTag        string
		Maintainer        string
		Name              string
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
		PullRequestPaths  []string
		RunCommands       []string
		RunMounts         []string
		StageFiles        []string
	}
	ContainerType int
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
		"HOME":            "/" + rootUser,
		"USER":            rootUser,
		"INITRD":          "No",
		"LANG":            defaultLanguage,
		"LANGUAGE":        defaultLanguage,
		"LC_ALL":          defaultLanguage,
		"ORG":             organization,
		"TZ":              "Etc/UTC",
		"PATH":            "${PATH}:" + usrLocalBinaryDir,
		"REPO":            repository,
		"SCCACHE_DIR":     "/_cache/sccache",
	}
	goDefaultEnvironments = map[string]string{
		"GOROOT":      "/opt/go",
		"GOPATH":      "/go",
		"GO111MODULE": "on",
		"PATH":        "${PATH}:${GOROOT}/bin:${GOPATH}/bin:" + usrLocalBinaryDir,
	}
	rustDefaultEnvironments = map[string]string{
		"RUST_HOME":     usrLocalLibDir + "/rust",
		"RUSTUP_HOME":   "${RUST_HOME}/rustup",
		"CARGO_HOME":    "${RUST_HOME}/cargo",
		"PATH":          "${PATH}:${RUSTUP_HOME}/bin:${CARGO_HOME}/bin:" + usrLocalBinaryDir,
		"RUSTC_WRAPPER": "/usr/bin/sccache",
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
		"make GOARCH=\"${TARGETARCH}\" GOOS=\"${TARGETOS}\" REPO=\"${ORG}/${REPO}\" NAME=\"${REPO}\" cmd/${PKG}/${APP_NAME}",
		"mv \"cmd/${PKG}/${APP_NAME}\" \"{{$.BinDir}}/${APP_NAME}\"",
	}
	goExampleBuildCommands = []string{
		"make GOARCH=\"${TARGETARCH}\" GOOS=\"${TARGETOS}\" REPO=\"${ORG}/${REPO}\" NAME=\"${REPO}\" ${PKG}/${APP_NAME}",
		"mv \"${PKG}/${APP_NAME}\" \"{{$.BinDir}}/${APP_NAME}\"",
	}
	rustBuildCommands = []string{
		"make rust/target/release/${APP_NAME}",
		"mv \"rust/target/release/${APP_NAME}\" \"{{$.BinDir}}/${APP_NAME}\"",
		"rm -rf rust/target",
	}
	e2eBuildCommands = []string{
		"make GOARCH=\"${TARGETARCH}\" GOOS=\"${TARGETOS}\" REPO=\"${ORG}/${REPO}\" NAME=\"${REPO}\" ${PKG}/${APP_NAME}",
		"mv \"${PKG}/${APP_NAME}\" \"{{$.BinDir}}/${APP_NAME}\"",
	}

	defaultMounts = []string{
		"--mount=type=bind,target=.,rw",
		"--mount=type=tmpfs,target=/tmp",
		"--mount=type=cache,target=/var/lib/apt,sharing=locked,id=${APP_NAME}-${TARGETARCH}",
		"--mount=type=cache,target=/var/cache/apt,sharing=locked,id=${APP_NAME}-${TARGETARCH}",
		"--mount=type=cache,target=/_cache/sccache,sharing=locked,id=sccache-${TARGETARCH}",
	}
	goDefaultMounts = []string{
		"--mount=type=cache,target=\"${GOPATH}/pkg\",id=\"go-build-${TARGETARCH}\"",
		"--mount=type=cache,target=\"${HOME}/.cache/go-build\",id=\"go-build-${TARGETARCH}\"",
		"--mount=type=tmpfs,target=\"${GOPATH}/src\"",
	}
	rustDefaultMounts = []string{
		"--mount=type=cache,target=\"${CARGO_HOME}/registry\",sharing=locked,id=\"cargo-registry-${TARGETARCH}\"",
		"--mount=type=cache,target=\"${CARGO_HOME}/git\",sharing=locked,id=\"cargo-git-${TARGETARCH}\"",
	}

	clangBuildDeps = []string{
		"cmake",
		"g++",
		"gcc",
		"libssl-dev",
		"unzip",
		"sccache",
		"ninja-build",
	}
	ngtBuildDeps = []string{
		"liblapack-dev",
		"libomp-dev",
		"libopenblas-dev",
		"gfortran",
	}
	rustBuildDeps = []string{
		"pkg-config",
	}
	devContainerDeps = []string{
		"file",
		"gawk",
		"git-lfs",
		"gnupg2",
		"graphviz",
		"jq",
		"libaec-dev",
		"sed",
		"zip",
	}

	ciContainerPreprocess = []string{
		"make GOARCH=${TARGETARCH} GOOS=${TARGETOS} deps GO_CLEAN_DEPS=false",
		"make GOARCH=${TARGETARCH} GOOS=${TARGETOS} golangci-lint/install",
		"make GOARCH=${TARGETARCH} GOOS=${TARGETOS} gotestfmt/install",
		"make buf/install",
		"make hdf5/install",
		"make helm-docs/install",
		"make helm/install",
		"make k3d/install",
		"make kind/install",
		"make kubectl/install",
		"make reviewdog/install",
		"make tparse/install",
		"make yq/install",
		"make docker-cli/install",
	}

	devContainerPreprocess = []string{
		"curl -fsSL https://deb.nodesource.com/setup_current.x | bash -",
		"apt-get update -y",
		"apt-get install -y --no-install-recommends --fix-missing nodejs",
		"npm install -g npm@latest",
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
			if ok && !strings.Contains(v, ev) {
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

// extractVariables efficiently extracts variables from strings
func extractVariables(value string) []string {
	var vars []string
	start := -1
	for i := 0; i < len(value); i++ {
		if value[i] == '$' && i+1 < len(value) && value[i+1] == '{' {
			start = i + 2
		} else if start != -1 && value[i] == '}' {
			vars = append(vars, value[start:i])
			start = -1
		} else if value[i] == '$' && start == -1 {
			start = i + 1
			for start < len(value) && (('a' <= value[start] && value[start] <= 'z') || ('A' <= value[start] && value[start] <= 'Z') || ('0' <= value[start] && value[start] <= '9') || value[start] == '_') {
				start++
			}
			vars = append(vars, value[i+1:start])
			i = start - 1
			start = -1
		}
	}
	return vars
}

// topologicalSort sorts the elements topologically and ensures that equal-level nodes are sorted by name
func topologicalSort(envMap map[string]string) []string {
	inDegree := make(map[string]int)         // Tracks the in-degree of each node
	graph := make(map[string][]string)       // Tracks the edges between nodes
	result := make([]string, 0, len(envMap)) // Result slice pre-allocated for efficiency

	gl := 0
	// Initialize the graph structure and in-degrees
	for key, value := range envMap {
		vars := extractVariables(value)
		for _, refKey := range vars {
			if refKey != key { // Prevent self-dependency
				graph[refKey] = append(graph[refKey], key)
				if len(graph[refKey]) > gl {
					gl = len(graph[refKey])
				}
				inDegree[key]++
			}
		}
	}

	// Initialize the queue with nodes having in-degree 0 (no dependencies)
	queue := make([]string, 0, len(envMap)-len(graph))
	for key := range envMap {
		if inDegree[key] == 0 {
			queue = append(queue, key)
		}
	}

	// Sort the initial queue to maintain lexicographical order for nodes with no dependencies
	slices.Sort(queue)

	// Preallocate a reusable slice for collecting new nodes
	newNodes := make([]string, 0, gl)
	// Topological sort process
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]

		// Append the result as `node=value`
		if value, exists := envMap[node]; exists {
			result = append(result, node+"="+value)
		}

		// Process all neighbors and decrement their in-degrees
		for _, neighbor := range graph[node] {
			inDegree[neighbor]--
			if inDegree[neighbor] == 0 {
				newNodes = append(newNodes, neighbor)
			}
		}

		// If new nodes were found, sort them and append to the queue
		if len(newNodes) > 0 {
			slices.Sort(newNodes) // Sort new nodes only once
			queue = append(queue, newNodes...)
			newNodes = newNodes[:0] // Reuse the slice by resetting it
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

	maintainer := os.Getenv(maintainerKey)
	if maintainer == "" {
		maintainer = defaultMaintainer
	}
	year := time.Now().Year()
	eg, egctx := errgroup.New(ctx)
	for n, d := range map[string]Data{
		"vald-" + agentNGT: {
			AppName:       "ngt",
			PackageDir:    agent + "/core/ngt",
			ExtraPackages: append(clangBuildDeps, ngtBuildDeps...),
			Preprocess:    []string{ngtPreprocess},
		},
		"vald-" + agentFaiss: {
			AppName:       "faiss",
			PackageDir:    agent + "/core/faiss",
			ExtraPackages: append(clangBuildDeps, ngtBuildDeps...),
			Preprocess:    []string{faissPreprocess},
		},
		"vald-" + agent: {
			AppName:       agent,
			PackageDir:    agent + "/core/" + agent,
			ContainerType: Rust,
			RuntimeImage:  "gcr.io/distroless/cc-debian12",
			ExtraPackages: append(clangBuildDeps,
				append(ngtBuildDeps, rustBuildDeps...)...),
			Preprocess: []string{
				ngtPreprocess,
				faissPreprocess,
			},
		},
		"vald-" + agentSidecar: {
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
		"vald-index-deletion": {
			AppName:    "index-deletion",
			PackageDir: "index/job/deletion",
		},
		"vald-index-exportation": {
			AppName:    "index-exportation",
			PackageDir: "index/job/exportation",
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
			ExtraPackages: append(clangBuildDeps, "libaec-dev"),
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
				"--from=operator " + usrLocalBinaryDir + "/${APP_NAME} {{$.BinDir}}/${APP_NAME}",
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
		"vald-ci-container": {
			AppName:       "ci-container",
			ContainerType: CIContainer,
			PackageDir:    "ci/base",
			RuntimeUser:   defaultBuildUser,
			ExtraPackages: append([]string{"npm", "sudo"}, append(clangBuildDeps,
				append(ngtBuildDeps,
					append(rustBuildDeps,
						devContainerDeps...)...)...)...),
			Preprocess:  append(ciContainerPreprocess, ngtPreprocess, faissPreprocess, usearchPreprocess),
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
					append(rustBuildDeps,
						devContainerDeps...)...)...),
			Preprocess: append(devContainerPreprocess,
				append(ciContainerPreprocess,
					ngtPreprocess,
					faissPreprocess)...),
		},
		"vald-example-client": {
			AppName:       "client",
			PackageDir:    "example/client",
			ExtraPackages: append(clangBuildDeps, "libaec-dev"),
			Preprocess: []string{
				"make hdf5/install",
			},
		},
		"vald-e2e": {
			AppName:       "e2e",
			PackageDir:    "tests/v2/e2e",
			ExtraPackages: append(clangBuildDeps, "libaec-dev"),
			Preprocess: []string{
				"make hdf5/install",
			},
		},
		"vald-buildbase": {
			AppName:      "buildbase",
			AliasImage:   true,
			PackageDir:   "buildbase",
			BuilderImage: "ubuntu",
			BuilderTag:   "devel",
		},
		"vald-buildkit": {
			AppName:      "buildkit",
			AliasImage:   true,
			PackageDir:   "buildkit",
			BuilderImage: "moby/buildkit",
			BuilderTag:   "master",
		},
		"vald-binfmt": {
			AppName:      "binfmt",
			AliasImage:   true,
			PackageDir:   "binfmt",
			BuilderImage: "tonistiigi/binfmt",
			BuilderTag:   "master",
		},
		"vald-buildkit-syft-scanner": {
			AppName:        "scanner",
			AliasImage:     true,
			PackageDir:     "buildkit/syft/scanner",
			BuilderImage:   "docker/buildkit-syft-scanner",
			BuilderTag:     "edge",
			BuildStageName: "scanner",
		},
	} {
		name := n
		data := d

		eg.Go(safety.RecoverFunc(func() error {
			data.Name = strings.TrimPrefix(name, "vald-")
			switch data.ContainerType {
			case HelmOperator:
				data.PullRequestPaths = append(data.PullRequestPaths,
					chartPath,
					valuesPath,
					templatesPath,
					helmOperatorChartPath,
					helmOperatorValuesPath,
					helmOperatorTemplatesPath,
					operatorSDKVersionPath,
				)
			case DevContainer, CIContainer:
				data.PullRequestPaths = append(data.PullRequestPaths,
					apisProtoPath,
					hackPath,
				)
			case Go:
				data.PullRequestPaths = append(data.PullRequestPaths,
					apisProtoPath,
					goModPath,
					goSumPath,
					goVersionPath,
					excludeTestFilesPath,
					excludeMockFilesPath,
				)
				mainFile := file.Join(os.Args[1], "cmd", data.PackageDir, "main.go")
				if file.Exists(mainFile) {
					ns, err := buildDependencyTree(os.Args[1], mainFile)
					if err != nil {
						log.Error(err)
					}
					pkgs := make([]string, 0, len(ns)+1)
					pkgs = append(pkgs, file.Join("cmd", data.PackageDir))
					for _, pnode := range ns {
						pkgs = append(pkgs, pnode.ToSlice()...)
					}
					slices.Sort(pkgs)
					pkgs = slices.Compact(pkgs)
					root, err := os.Getwd()
					if err != nil {
						root = os.Getenv("HOME")
					}
					if root != "" && !strings.HasSuffix(root, string(os.PathSeparator)) {
						root += string(os.PathSeparator)
					}
					for i, pkg := range pkgs {
						const splitWord = "/vdaas/vald/"
						pkg = file.Join(pkg, "*.go")
						index := strings.LastIndex(pkg, splitWord)
						if index != -1 {
							pkg = pkg[index+len(splitWord):]
						}
						if root != "" {
							pkg = strings.TrimPrefix(pkg, root)
						}
						pkgs[i] = pkg
					}
					data.PullRequestPaths = append(data.PullRequestPaths, pkgs...)
				}
			case Rust:
				data.PullRequestPaths = append(data.PullRequestPaths,
					apisProtoPath,
					cargoLockPath,
					cargoTomlPath,
					rustBinAgentDirPath,
					rustNgtRsPath,
					rustNgtPath,
					rustProtoPath,
					rustVersionPath,
				)
			}
			if strings.EqualFold(data.Name, agentFaiss) || data.ContainerType == Rust {
				data.PullRequestPaths = append(data.PullRequestPaths, faissVersionPath)
			}
			if strings.EqualFold(data.Name, agentNGT) || data.ContainerType == Rust {
				data.PullRequestPaths = append(data.PullRequestPaths, ngtVersionPath)
			}

			if !data.AliasImage {
				data.PullRequestPaths = append(data.PullRequestPaths, makefilePath, makefileDirPath)
			}

			if data.AliasImage {
				data.BuildPlatforms = multiPlatforms
			}
			if data.ContainerType == CIContainer {
				data.BuildPlatforms = amd64Platform
			}

			data.Year = time.Now().Year()
			if maintainer := os.Getenv(maintainerKey); maintainer != "" {
				data.Maintainer = maintainer
			} else {
				data.Maintainer = defaultMaintainer
			}

			log.Infof("Generating %s's workflow", data.Name)
			workflow := new(Workflow)
			err := yaml.Unmarshal(conv.Atob(`name: "Build docker image: `+data.Name+`"
on:
  schedule:
    - cron: "0 * * * *"
  push:
    branches:
      - "main"
      - "release/v*.*"
      - "!release/v*.*.*"
    tags:
      - "*.*.*"
      - "*.*.*-*"
      - "v*.*.*"
      - "v*.*.*-*"
  pull_request:
    paths:
      - ".github/actions/docker-build/action.yaml"
      - ".github/workflows/_docker-image.yaml"
      - ".github/workflows/dockers-`+data.Name+`-image.yaml"
      - "dockers/`+data.PackageDir+`/Dockerfile"
      - "hack/docker/gen/main.go"
  pull_request_target:
    types: [opened, reopened, synchronize, labeled]
    paths: []

jobs:
  build:
    uses: "./.github/workflows/_docker-image.yaml"
    with:
      target: "`+data.Name+`"
      platforms: ""
    secrets: "inherit"
`), &workflow)
			if err != nil {
				return fmt.Errorf("Error decoding YAML: %v", err)
			}

			if !data.AliasImage {
				workflow.On.Schedule = nil
			}
			workflow.On.PullRequest.Paths = append(workflow.On.PullRequest.Paths, data.PullRequestPaths...)
			if strings.EqualFold(data.Name, exampleContainer) {
				workflow.On.PullRequest.Paths = slices.DeleteFunc(workflow.On.PullRequest.Paths, func(path string) bool {
					return strings.HasPrefix(path, "cmd") || strings.HasPrefix(path, "pkg")
				})
				workflow.On.PullRequest.Paths = append(workflow.On.PullRequest.Paths, data.PackageDir+"/**")
			}
			slices.Sort(workflow.On.PullRequest.Paths)
			workflow.On.PullRequest.Paths = slices.Compact(workflow.On.PullRequest.Paths)

			workflow.On.PullRequestTarget.Paths = workflow.On.PullRequest.Paths
			workflow.Jobs.Build.With.Platforms = data.BuildPlatforms

			workflowYamlTmp, err := yaml.Marshal(workflow)
			if err != nil {
				return fmt.Errorf("error marshaling workflowStruct to YAML: %w", err)
			}

			// remove the double quotation marks from the generated key "on": (note that the word "on" is a reserved word in sigs.k8s.io/yaml)
			workflowYaml := strings.Replace(string(workflowYamlTmp), "\"on\":", "on:", 1)

			if len(header) > (int(^uint(0)>>1) - len(workflowYaml)) {
				return fmt.Errorf("size computation for allocation may overflow")
			}
			totalLen := len(header) + len(workflowYaml)

			buf := bytes.NewBuffer(make([]byte, 0, totalLen))
			err = license.Execute(buf, data)
			if err != nil {
				return fmt.Errorf("error executing template: %w", err)
			}
			buf.WriteString("\r\n")
			buf.WriteString(workflowYaml)
			fileName := file.Join(os.Args[1], ".github/workflows", "dockers-"+data.Name+"-image.yaml")
			_, err = file.OverWriteFile(egctx, fileName, buf, fs.ModePerm)
			if err != nil {
				return fmt.Errorf("error writing workflow file for %s error: %w", fileName, err)
			}
			return nil
		}))

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
			if data.BuildStageName == "" {
				data.BuildStageName = defaultBuildStageName
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
				} else if strings.HasPrefix(data.PackageDir, "example") && file.Exists(file.Join(os.Args[1], data.PackageDir)) {
					commands = append(commands, goExampleBuildCommands...)
				} else if strings.HasPrefix(data.PackageDir, "tests/v2/e2e") && file.Exists(file.Join(os.Args[1], data.PackageDir)) {
					commands = append(commands, e2eBuildCommands...)
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
				mounts := make([]string, 0, len(defaultMounts)+len(rustDefaultMounts))
				mounts = append(mounts, defaultMounts...)
				mounts = append(mounts, rustDefaultMounts...)
				data.RunMounts = mounts
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
				mounts := make([]string, 0, len(defaultMounts)+len(goDefaultMounts)+len(rustDefaultMounts))
				mounts = append(mounts, defaultMounts...)
				mounts = append(mounts, goDefaultMounts...)
				mounts = append(mounts, rustDefaultMounts...)
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
			if strings.Contains(data.BuildUser, rootUser) {
				data.Environments["HOME"] = "/" + rootUser
				data.Environments["USER"] = rootUser
			} else {
				user, _, _ := strings.Cut(data.BuildUser, ":")
				data.Environments["HOME"] = "/home/" + user
				data.Environments["USER"] = user
			}

			data.Environments["APP_NAME"] = data.AppName
			data.Environments["PKG"] = data.PackageDir
			data.EnvironmentsSlice = topologicalSort(data.Environments)
			data.ConfigExists = file.Exists(file.Join(os.Args[1], "cmd", data.PackageDir, "sample.yaml"))

			buf := bytes.NewBuffer(make([]byte, 0, 1024))
			log.Infof("Generating %s's Dockerfile", name)
			docker.Execute(buf, data)
			tpl := buf.String()
			buf.Reset()
			template.Must(template.New("Dockerfile").Parse(tpl)).Execute(buf, data)
			fileName := file.Join(os.Args[1], "dockers", data.PackageDir, "Dockerfile")
			_, err := file.OverWriteFile(egctx, fileName, buf, fs.ModePerm)
			if err != nil {
				return fmt.Errorf("error writing Dockerfile for %s error: %w", fileName, err)
			}
			return nil
		}))
	}
	eg.Wait()
}

// PackageNode represents a node in the dependency tree.
type PackageNode struct {
	Name    string
	Imports []*PackageNode
}

// ToSlice traverses the dependency tree and returns all dependencies as a slice.
func (n PackageNode) ToSlice() (pkgs []string) {
	pkgs = make([]string, 0, len(n.Imports)+1)
	if n.Name != "command-line-arguments" {
		pkgs = append(pkgs, n.Name)
	}
	for _, node := range n.Imports {
		pkgs = append(pkgs, node.ToSlice()...)
	}
	return pkgs
}

// String returns string of the dependency tree in a readable format.
func (n PackageNode) String() string {
	return n.string(0)
}

func (n PackageNode) string(depth int) (tree string) {
	tree = fmt.Sprintf("%s- %s\n", strings.Repeat("  ", depth), n.Name)
	for _, node := range n.Imports {
		tree += node.string(depth + 1)
	}
	return tree
}

// processDependencies processes package dependencies while avoiding duplicate processing.
func processDependencies(
	pkg *packages.Package,
	nodes map[string]*PackageNode,
	mu *sync.Mutex,
	checkList map[string]*PackageNode,
	wg *sync.WaitGroup,
) *PackageNode {
	if !strings.Contains(pkg.PkgPath, "vdaas/vald") && pkg.Name != "main" {
		return nil
	}
	if node, exists := checkList[pkg.PkgPath]; exists {
		return node
	}

	node := &PackageNode{Name: pkg.PkgPath}
	nodes[pkg.PkgPath] = node
	checkList[pkg.PkgPath] = node
	for _, imp := range pkg.Imports {
		if !strings.Contains(imp.PkgPath, "vdaas/vald") {
			continue
		}
		if child, exists := checkList[imp.PkgPath]; exists {
			node.Imports = append(node.Imports, child)
			continue
		}
		child := processDependencies(imp, nodes, mu, checkList, wg)
		if child != nil {
			node.Imports = append(node.Imports, child)
		}
	}

	return node
}

// buildDependencyTree constructs a dependency tree for multiple entry packages.
func buildDependencyTree(rootDir, entryFile string) ([]*PackageNode, error) {
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedImports | packages.NeedDeps,
		Dir:  rootDir,
	}

	// Use entry file (e.g., main.go) as the root for analysis.
	pkgs, err := packages.Load(cfg, entryFile)
	if err != nil {
		return nil, err
	}

	nodes := make(map[string]*PackageNode)
	checkList := make(map[string]*PackageNode, len(pkgs)) // Tracks processed packages
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Process all entry packages
	var roots []*PackageNode
	for _, pkg := range pkgs {
		root := processDependencies(pkg, nodes, &mu, checkList, &wg)
		if root != nil {
			roots = append(roots, root)
		}
	}
	wg.Wait()

	return roots, nil
}
