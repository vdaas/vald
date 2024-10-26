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
	"syscall"
	"text/template"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/file"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"gopkg.in/yaml.v2"
)

var license string = `#
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

# DO_NOT_EDIT this workflow file is generated by https://github.com/vdaas/vald/blob/main/hack/actions/gen/main.go
`

var licenseTmpl *template.Template = template.Must(template.New("license").Parse(license))

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

	Paths []string
)

type Data struct {
	AliasImage        bool
	ConfigExists      bool
	Year              int
	ContainerType     ContainerType
	AppName           string
	BinDir            string
	BuildUser         string
	BuilderImage      string
	BuilderTag        string
	BuildStageName    string
	Maintainer        string
	PackageDir        string
	RootDir           string
	RuntimeImage      string
	RuntimeTag        string
	RuntimeUser       string
	Name              string
	BuildPlatforms    string
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
	PullRequestPaths  []string
}

type ContainerType int

const (
	Go ContainerType = iota
	Rust
	DevContainer
	HelmOperator
	CIContainer
	Other
)

const (
	organization          = "vdaas"
	repository            = "vald"
	defaultMaintainer     = organization + ".org " + repository + " team <" + repository + "@" + organization + ".org>"
	defaultBuildUser      = "root:root"
	maintainerKey         = "MAINTAINER"
	minimumArgumentLength = 2
	ubuntuVersion         = "22.04"

	ngtPreprocess   = "make ngt/install"
	faissPreprocess = "make faiss/install"

	helmOperatorRootdir   = "/opt/helm"
	helmOperatorWatchFile = helmOperatorRootdir + "/watches.yaml"
	helmOperatorChartsDir = helmOperatorRootdir + "/charts"
)

var (
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

const baseWorkflowTmpl string = `name: "Build docker image: %s"
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
      - "v*.*.*"
      - "*.*.*-*"
      - "v*.*.*-*"
  pull_request:
    paths:
      - "hack/docker/gen/main.go"
      - "dockers/%s/Dockerfile"
      - "hack/actions/gen/main.go"
      - ".github/workflows/dockers-%s-image.yaml"
      - ".github/actions/docker-build/action.yaml"
      - ".github/workflows/_docker-image.yaml"
      - "cmd/%s/**"
      - "pkg/%s/**"
  pull_request_target:
    paths: []

jobs:
  build:
    uses: "./.github/workflows/_docker-image.yaml"
    with:
      target: "%s"
      platforms: ""
    secrets: "inherit"
`

const (
	cmdBenchOperatorsPath = "cmd/tools/benchmark/operators/**"
	pkgBenchOperatorsPath = "pkg/tools/benchmark/operators/**"
	cmdBenchJobsPath      = "cmd/tools/benchmark/jobs/**"
	pkgBenchJobsPath      = "pkg/tools/benchmark/jobs/**"

	agentInternalPath   = "pkg/agent/internal/**"
	gatewayInternalPath = "pkg/gateway/internal/**"

	apisGrpcPath  = "apis/grpc/**"
	apisProtoPath = "apis/proto/**"

	hackPath = "hack/**"

	chartsValdPath            = "charts/vald"
	helmOperatorPath          = "charts/vald-helm-operator"
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

	internalPath         = "internal/**"
	internalStoragePath  = "internal/db/storage/blob/**"
	excludeTestFilesPath = "!internal/**/*_test.go"
	excludeMockFilesPath = "!internal/**/*_mock.go"
	excludeDbPath        = "!internal/db/**"
	excludeK8sPath       = "!internal/k8s/**"

	versionsPath           = "versions"
	operatorSDKVersionPath = versionsPath + "/OPERATOR_SDK_VERSION"
	goVersionPath          = versionsPath + "/GO_VERSION"
	rustVersionPath        = versionsPath + "/RUST_VERSION"
	faissVersionPath       = versionsPath + "/FAISS_VERSION"
	ngtVersionPath         = versionsPath + "/NGT_VERSION"
	usearchVersionPath     = versionsPath + "/USEARCH_VERSION"

	makefilePath    = "Makefile"
	makefileDirPath = "Makefile.d/**"
)

const (
	agentNgt     = "agent-ngt"
	agentFaiss   = "agent-faiss"
	agentSidecar = "agent-sidecar"
	agent        = "agent"

	discovererK8s = "discoverer-k8s"

	gateway       = "gateway"
	gatewayLb     = "gateway-lb"
	gatewayFilter = "gateway-filter"
	gatewayMirror = "gateway-mirror"

	managerIndex = "manager-index"

	indexCorrection = "index-correction"
	indexCreation   = "index-creation"
	indexDeletion   = "index-deletion"
	indexSave       = "index-save"
	indexOperator   = "index-operator"

	readreplicaRotate = "readreplica-rotate"

	benchJob      = "benchmark-job"
	benchOperator = "benchmark-operator"

	helmOperator = "helm-operator"

	loadtest = "loadtest"

	ciContainer  = "ci-container"
	devContainer = "dev-container"

	buildbase           = "buildbase"
	buildkit            = "buildkit"
	binfmt              = "binfmt"
	buildkitSyftScanner = "buildkit-syft-scanner"
)

const (
	multiPlatforms = amd64Platform + "," + arm64Platform
	amd64Platform  = "linux/amd64"
	arm64Platform  = "linux/arm64"
)

func (data *Data) initPullRequestPaths() {
	switch data.Name {
	// the benchmark components trigger each other, not just themselves
	case benchJob:
		data.PullRequestPaths = append(data.PullRequestPaths,
			cmdBenchOperatorsPath,
			pkgBenchOperatorsPath,
		)
	case benchOperator:
		data.PullRequestPaths = append(data.PullRequestPaths,
			cmdBenchJobsPath,
			pkgBenchJobsPath,
		)
	case agentFaiss, agentNgt:
		data.PullRequestPaths = append(data.PullRequestPaths, agentInternalPath)
	default:
		if strings.Contains(strings.ToLower(data.Name), gateway) {
			data.PullRequestPaths = append(data.PullRequestPaths, gatewayInternalPath)
		}
	}

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
			apisGrpcPath,
			apisProtoPath,
			hackPath,
		)
	case Go:
		data.PullRequestPaths = append(data.PullRequestPaths,
			apisGrpcPath,
			apisProtoPath,
			goModPath,
			goSumPath,
			goVersionPath,
			internalPath,
			excludeTestFilesPath,
			excludeMockFilesPath,
			excludeDbPath,
		)
		switch data.Name {
		case discovererK8s, indexOperator, gatewayMirror, readreplicaRotate, agentNgt, benchJob, benchOperator:
		default:
			data.PullRequestPaths = append(data.PullRequestPaths, excludeK8sPath)
		}
	case Rust:
		data.PullRequestPaths = append(data.PullRequestPaths,
			apisGrpcPath,
			apisProtoPath,
			cargoLockPath,
			cargoTomlPath,
			rustBinAgentDirPath,
			rustNgtRsPath,
			rustNgtPath,
			rustProtoPath,
			rustVersionPath,
			faissVersionPath,
			ngtVersionPath,
		)
	}
	if strings.EqualFold(data.Name, agentFaiss) || data.ContainerType == Rust {
		data.PullRequestPaths = append(data.PullRequestPaths, faissVersionPath)
	}
	if strings.EqualFold(data.Name, agentNgt) || data.ContainerType == Rust {
		data.PullRequestPaths = append(data.PullRequestPaths, ngtVersionPath)
	}

	if data.Name == agentSidecar {
		data.PullRequestPaths = append(data.PullRequestPaths, internalStoragePath)
	}
	if !data.AliasImage {
		data.PullRequestPaths = append(data.PullRequestPaths, makefilePath, makefileDirPath)
	}
}

func (data *Data) initData() {
	data.initPullRequestPaths()

	if data.AliasImage {
		data.BuildPlatforms = multiPlatforms
	}
	if data.ContainerType == CIContainer || data.Name == loadtest {
		data.BuildPlatforms = amd64Platform
	}

	data.Year = time.Now().Year()
	if maintainer := os.Getenv(maintainerKey); maintainer != "" {
		data.Maintainer = maintainer
	} else {
		data.Maintainer = defaultMaintainer
	}
}

func (data *Data) generateWorkflowStruct() (*Workflow, error) {
	workflow := &Workflow{}
	baseWorkflow := fmt.Sprintf(baseWorkflowTmpl,
		data.Name,
		data.PackageDir,
		data.Name,
		data.PackageDir,
		data.PackageDir,
		data.Name,
	)
	err := yaml.NewDecoder(strings.NewReader(baseWorkflow)).Decode(workflow)
	if err != nil {
		return nil, fmt.Errorf("Error decoding YAML: %v", err)
	}

	if !data.AliasImage {
		workflow.On.Schedule = nil
	}
	workflow.On.PullRequest.Paths = append(workflow.On.PullRequest.Paths, data.PullRequestPaths...)
	workflow.On.PullRequestTarget.Paths = workflow.On.PullRequest.Paths
	workflow.Jobs.Build.With.Platforms = data.BuildPlatforms

	return workflow, nil
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
		"vald-index-deletion": {
			AppName:    "index-deletion",
			PackageDir: "index/job/deletion",
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
		"vald-loadtest": { // note: this name is a little different from that of docker/gen/main.go
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
			data.initData()

			log.Infof("Generating %s's workflow", data.Name)
			workflow, err := data.generateWorkflowStruct()
			if err != nil {
				return fmt.Errorf("Error generating workflowStruct: %w", err)
			}
			workflowYamlTmp, err := yaml.Marshal(workflow)
			if err != nil {
				return fmt.Errorf("error marshaling workflowStruct to YAML: %w", err)
			}

			// remove the double quotation marks from the generated key "on": (note that the word "on" is a reserved word in gopkg.in/yaml.v2)
			workflowYaml := strings.Replace(string(workflowYamlTmp), "\"on\":", "on:", 1)

			buf := bytes.NewBuffer(make([]byte, 0, len(license)+len(workflowYaml)))
			err = licenseTmpl.Execute(buf, data)
			if err != nil {
				return fmt.Errorf("error executing template: %w", err)
			}
			buf.WriteString("\r\n")
			buf.WriteString(workflowYaml)
			_, err = file.OverWriteFile(egctx, file.Join(os.Args[1], ".github/workflows", "dockers-"+data.Name+"-image.yaml"), buf, fs.ModePerm)
			if err != nil {
				return fmt.Errorf("error writing file: %w", err)
			}
			return nil
		}))
	}
	eg.Wait()
}
