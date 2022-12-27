#
# Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
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

.PHONY: binary/build
## build all binaries
binary/build: \
	cmd/agent/core/ngt/ngt \
	cmd/agent/sidecar/sidecar \
	cmd/discoverer/k8s/discoverer \
	cmd/gateway/lb/lb \
	cmd/gateway/filter/filter \
	cmd/manager/index/index \
	cmd/tools/benchmark/job/job \
	cmd/tools/benchmark/operator/operator

cmd/agent/core/ngt/ngt: \
	ngt/install \
	$(GO_SOURCES_INTERNAL) \
	$(PBGOS) \
	$(shell find ./cmd/agent/core/ngt -type f -name '*.go' -not -name '*_test.go' -not -name 'doc.go') \
	$(shell find ./pkg/agent/core/ngt ./pkg/agent/internal -type f -name '*.go' -not -name '*_test.go' -not -name 'doc.go')
	CFLAGS="$(CFLAGS)" \
	CXXFLAGS="$(CXXFLAGS)" \
	CGO_ENABLED=1 \
	CGO_CXXFLAGS="-g -Ofast -march=native" \
	CGO_FFLAGS="-g -Ofast -march=native" \
	CGO_LDFLAGS="-g -Ofast -march=native" \
	GO111MODULE=on \
	GOPRIVATE=$(GOPRIVATE) \
	go build \
		--ldflags "-w -linkmode 'external' \
		-extldflags '-static -fPIC -pthread -fopenmp -std=gnu++20 -lstdc++ -lm -z relro -z now $(EXTLDFLAGS)' \
		-X '$(GOPKG)/internal/info.Version=$(VERSION)' \
		-X '$(GOPKG)/internal/info.GitCommit=$(GIT_COMMIT)' \
		-X '$(GOPKG)/internal/info.BuildTime=$(DATETIME)' \
		-X '$(GOPKG)/internal/info.GoVersion=$(GO_VERSION)' \
		-X '$(GOPKG)/internal/info.GoOS=$(GOOS)' \
		-X '$(GOPKG)/internal/info.GoArch=$(GOARCH)' \
		-X '$(GOPKG)/internal/info.CGOEnabled=$${CGO_ENABLED}' \
		-X '$(GOPKG)/internal/info.NGTVersion=$(NGT_VERSION)' \
		-X '$(GOPKG)/internal/info.BuildCPUInfoFlags=$(CPU_INFO_FLAGS)' \
		-buildid=" \
		-mod=readonly \
		-modcacherw \
		-a \
		-tags "cgo osusergo netgo static_build" \
		-trimpath \
		-o $@ \
		$(dir $@)main.go
	$@ -version

cmd/agent/sidecar/sidecar: \
	$(GO_SOURCES_INTERNAL) \
	$(PBGOS) \
	$(shell find ./cmd/agent/sidecar -type f -name '*.go' -not -name '*_test.go' -not -name 'doc.go') \
	$(shell find ./pkg/agent/sidecar ./pkg/agent/internal -type f -name '*.go' -not -name '*_test.go' -not -name 'doc.go')
	CGO_ENABLED=0 \
	GO111MODULE=on \
	GOPRIVATE=$(GOPRIVATE) \
	go build \
		--ldflags "-w -extldflags=-static \
		-X '$(GOPKG)/internal/info.Version=$(VERSION)' \
		-X '$(GOPKG)/internal/info.GitCommit=$(GIT_COMMIT)' \
		-X '$(GOPKG)/internal/info.BuildTime=$(DATETIME)' \
		-X '$(GOPKG)/internal/info.GoVersion=$(GO_VERSION)' \
		-X '$(GOPKG)/internal/info.GoOS=$(GOOS)' \
		-X '$(GOPKG)/internal/info.GoArch=$(GOARCH)' \
		-X '$(GOPKG)/internal/info.CGOEnabled=$${CGO_ENABLED}' \
		-X '$(GOPKG)/internal/info.BuildCPUInfoFlags=$(CPU_INFO_FLAGS)' \
		-buildid=" \
		-mod=readonly \
		-modcacherw \
		-a \
		-tags "osusergo netgo static_build" \
		-trimpath \
		-o $@ \
		$(dir $@)main.go
	$@ -version

cmd/discoverer/k8s/discoverer: \
	$(GO_SOURCES_INTERNAL) \
	$(PBGOS) \
	$(shell find ./cmd/discoverer/k8s -type f -name '*.go' -not -name '*_test.go' -not -name 'doc.go') \
	$(shell find ./pkg/discoverer/k8s -type f -name '*.go' -not -name '*_test.go' -not -name 'doc.go')
	CGO_ENABLED=0 \
	GO111MODULE=on \
	GOPRIVATE=$(GOPRIVATE) \
	go build \
		--ldflags "-w -extldflags=-static \
		-X '$(GOPKG)/internal/info.Version=$(VERSION)' \
		-X '$(GOPKG)/internal/info.GitCommit=$(GIT_COMMIT)' \
		-X '$(GOPKG)/internal/info.BuildTime=$(DATETIME)' \
		-X '$(GOPKG)/internal/info.GoVersion=$(GO_VERSION)' \
		-X '$(GOPKG)/internal/info.GoOS=$(GOOS)' \
		-X '$(GOPKG)/internal/info.GoArch=$(GOARCH)' \
		-X '$(GOPKG)/internal/info.CGOEnabled=$${CGO_ENABLED}' \
		-X '$(GOPKG)/internal/info.BuildCPUInfoFlags=$(CPU_INFO_FLAGS)' \
		-buildid=" \
		-mod=readonly \
		-modcacherw \
		-a \
		-tags "osusergo netgo static_build" \
		-trimpath \
		-o $@ \
		$(dir $@)main.go
	$@ -version

cmd/gateway/lb/lb: \
	$(GO_SOURCES_INTERNAL) \
	$(PBGOS) \
	$(shell find ./cmd/gateway/lb -type f -name '*.go' -not -name '*_test.go' -not -name 'doc.go') \
	$(shell find ./pkg/gateway/lb -type f -name '*.go' -not -name '*_test.go' -not -name 'doc.go')
	CGO_ENABLED=0 \
	GO111MODULE=on \
	GOPRIVATE=$(GOPRIVATE) \
	go build \
		--ldflags "-w -extldflags=-static \
		-X '$(GOPKG)/internal/info.Version=$(VERSION)' \
		-X '$(GOPKG)/internal/info.GitCommit=$(GIT_COMMIT)' \
		-X '$(GOPKG)/internal/info.BuildTime=$(DATETIME)' \
		-X '$(GOPKG)/internal/info.GoVersion=$(GO_VERSION)' \
		-X '$(GOPKG)/internal/info.GoOS=$(GOOS)' \
		-X '$(GOPKG)/internal/info.GoArch=$(GOARCH)' \
		-X '$(GOPKG)/internal/info.CGOEnabled=$${CGO_ENABLED}' \
		-X '$(GOPKG)/internal/info.BuildCPUInfoFlags=$(CPU_INFO_FLAGS)' \
		-buildid=" \
		-mod=readonly \
		-modcacherw \
		-a \
		-tags "osusergo netgo static_build" \
		-trimpath \
		-o $@ \
		$(dir $@)main.go
	$@ -version

cmd/gateway/filter/filter: \
	$(GO_SOURCES_INTERNAL) \
	$(PBGOS) \
	$(shell find ./cmd/gateway/filter -type f -name '*.go' -not -name '*_test.go' -not -name 'doc.go') \
	$(shell find ./pkg/gateway/filter -type f -name '*.go' -not -name '*_test.go' -not -name 'doc.go')
	CGO_ENABLED=0 \
	GO111MODULE=on \
	GOPRIVATE=$(GOPRIVATE) \
	go build \
		--ldflags "-w -extldflags=-static \
		-X '$(GOPKG)/internal/info.Version=$(VERSION)' \
		-X '$(GOPKG)/internal/info.GitCommit=$(GIT_COMMIT)' \
		-X '$(GOPKG)/internal/info.BuildTime=$(DATETIME)' \
		-X '$(GOPKG)/internal/info.GoVersion=$(GO_VERSION)' \
		-X '$(GOPKG)/internal/info.GoOS=$(GOOS)' \
		-X '$(GOPKG)/internal/info.GoArch=$(GOARCH)' \
		-X '$(GOPKG)/internal/info.CGOEnabled=$${CGO_ENABLED}' \
		-X '$(GOPKG)/internal/info.BuildCPUInfoFlags=$(CPU_INFO_FLAGS)' \
		-buildid=" \
		-mod=readonly \
		-modcacherw \
		-a \
		-tags "osusergo netgo static_build" \
		-trimpath \
		-o $@ \
		$(dir $@)main.go
	$@ -version

cmd/manager/index/index: \
	$(GO_SOURCES_INTERNAL) \
	$(PBGOS) \
	$(shell find ./cmd/manager/index -type f -name '*.go' -not -name '*_test.go' -not -name 'doc.go') \
	$(shell find ./pkg/manager/index -type f -name '*.go' -not -name '*_test.go' -not -name 'doc.go')
	CGO_ENABLED=0 \
	GO111MODULE=on \
	GOPRIVATE=$(GOPRIVATE) \
	go build \
		--ldflags "-w -extldflags=-static \
		-X '$(GOPKG)/internal/info.Version=$(VERSION)' \
		-X '$(GOPKG)/internal/info.GitCommit=$(GIT_COMMIT)' \
		-X '$(GOPKG)/internal/info.BuildTime=$(DATETIME)' \
		-X '$(GOPKG)/internal/info.GoVersion=$(GO_VERSION)' \
		-X '$(GOPKG)/internal/info.GoOS=$(GOOS)' \
		-X '$(GOPKG)/internal/info.GoArch=$(GOARCH)' \
		-X '$(GOPKG)/internal/info.CGOEnabled=$${CGO_ENABLED}' \
		-X '$(GOPKG)/internal/info.BuildCPUInfoFlags=$(CPU_INFO_FLAGS)' \
		-buildid=" \
		-mod=readonly \
		-modcacherw \
		-a \
		-tags "osusergo netgo static_build" \
		-trimpath \
		-o $@ \
		$(dir $@)main.go
	$@ -version

cmd/tools/benchmark/job/job: \
	$(GO_SOURCES_INTERNAL) \
	$(PBGOS) \
	$(shell find ./cmd/tools/benchmark/job -type f -name '*.go' -not -name '*_test.go' -not -name 'doc.go') \
	$(shell find ./pkg/tools/benchmark/job -type f -name '*.go' -not -name '*_test.go' -not -name 'doc.go')
	CFLAGS="$(CFLAGS)" \
	CXXFLAGS="$(CXXFLAGS)" \
	CGO_ENABLED=1 \
	CGO_CXXFLAGS="-g -Ofast -march=native" \
	CGO_FFLAGS="-g -Ofast -march=native" \
	CGO_LDFLAGS="-g -Ofast -march=native" \
	GO111MODULE=on \
	GOPRIVATE=$(GOPRIVATE) \
	go build \
		--ldflags "-s -w \
		-X '$(GOPKG)/internal/info.Version=$(VERSION)' \
		-X '$(GOPKG)/internal/info.GitCommit=$(GIT_COMMIT)' \
		-X '$(GOPKG)/internal/info.BuildTime=$(DATETIME)' \
		-X '$(GOPKG)/internal/info.GoVersion=$(GO_VERSION)' \
		-X '$(GOPKG)/internal/info.GoOS=$(GOOS)' \
		-X '$(GOPKG)/internal/info.GoArch=$(GOARCH)' \
		-X '$(GOPKG)/internal/info.CGOEnabled=$${CGO_ENABLED}' \
		-X '$(GOPKG)/internal/info.NGTVersion=$(NGT_VERSION)' \
		-X '$(GOPKG)/internal/info.BuildCPUInfoFlags=$(CPU_INFO_FLAGS)' \
		-buildid=" \
		-mod=readonly \
		-modcacherw \
		-a \
		-tags "cgo osusergo netgo" \
		-trimpath \
		-o $@ \
		$(dir $@)main.go
	$@ -version

cmd/tools/benchmark/scenario/scenario: \
	$(GO_SOURCES_INTERNAL) \
	$(PBGOS) \
	$(shell find ./cmd/tools/benchmark/scenario -type f -name '*.go' -not -name '*_test.go' -not -name 'doc.go') \
	$(shell find ./pkg/tools/benchmark/scenario -type f -name '*.go' -not -name '*_test.go' -not -name 'doc.go')
	CFLAGS="$(CFLAGS)" \
	CXXFLAGS="$(CXXFLAGS)" \
	CGO_ENABLED=1 \
	CGO_CXXFLAGS="-g -Ofast -march=native" \
	CGO_FFLAGS="-g -Ofast -march=native" \
	CGO_LDFLAGS="-g -Ofast -march=native" \
	GO111MODULE=on \
	GOPRIVATE=$(GOPRIVATE) \
	go build \
		--ldflags "-w -extldflags=-static \
		-X '$(GOPKG)/internal/info.Version=$(VERSION)' \
		-X '$(GOPKG)/internal/info.GitCommit=$(GIT_COMMIT)' \
		-X '$(GOPKG)/internal/info.BuildTime=$(DATETIME)' \
		-X '$(GOPKG)/internal/info.GoVersion=$(GO_VERSION)' \
		-X '$(GOPKG)/internal/info.GoOS=$(GOOS)' \
		-X '$(GOPKG)/internal/info.GoArch=$(GOARCH)' \
		-X '$(GOPKG)/internal/info.CGOEnabled=$${CGO_ENABLED}' \
		-X '$(GOPKG)/internal/info.NGTVersion=$(NGT_VERSION)' \
		-X '$(GOPKG)/internal/info.BuildCPUInfoFlags=$(CPU_INFO_FLAGS)' \
		-buildid=" \
		-mod=readonly \
		-modcacherw \
		-a \
		-tags "cgo osusergo netgo" \
		-trimpath \
		-o $@ \
		$(dir $@)main.go
	$@ -version


.PHONY: binary/build/zip
## build all binaries and zip them
binary/build/zip: \
	artifacts/vald-agent-ngt-$(GOOS)-$(GOARCH).zip \
	artifacts/vald-agent-sidecar-$(GOOS)-$(GOARCH).zip \
	artifacts/vald-discoverer-k8s-$(GOOS)-$(GOARCH).zip \
	artifacts/vald-lb-gateway-$(GOOS)-$(GOARCH).zip \
	artifacts/vald-filter-gateway-$(GOOS)-$(GOARCH).zip \
	artifacts/vald-manager-index-$(GOOS)-$(GOARCH).zip

artifacts/vald-agent-ngt-$(GOOS)-$(GOARCH).zip: cmd/agent/core/ngt/ngt
	$(call mkdir, $(dir $@))
	zip --junk-paths $@ $<

artifacts/vald-agent-sidecar-$(GOOS)-$(GOARCH).zip: cmd/agent/sidecar/sidecar
	$(call mkdir, $(dir $@))
	zip --junk-paths $@ $<

artifacts/vald-discoverer-k8s-$(GOOS)-$(GOARCH).zip: cmd/discoverer/k8s/discoverer
	$(call mkdir, $(dir $@))
	zip --junk-paths $@ $<

artifacts/vald-lb-gateway-$(GOOS)-$(GOARCH).zip: cmd/gateway/lb/lb
	$(call mkdir, $(dir $@))
	zip --junk-paths $@ $<

artifacts/vald-filter-gateway-$(GOOS)-$(GOARCH).zip: cmd/gateway/filter/filter
	$(call mkdir, $(dir $@))
	zip --junk-paths $@ $<

artifacts/vald-manager-index-$(GOOS)-$(GOARCH).zip: cmd/manager/index/index
	$(call mkdir, $(dir $@))
	zip --junk-paths $@ $<

artifacts/vald-benchmark-job-$(GOOS)-$(GOARCH).zip: cmd/tools/benchmark/job/job
	$(call mkdir, $(dir $@))
	zip --junk-paths $@ $<

artifacts/vald-benchmark-operator-$(GOOS)-$(GOARCH).zip: cmd/tools/benchmark/operator/operator
	$(call mkdir, $(dir $@))
	zip --junk-paths $@ $<
