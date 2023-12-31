#
# Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
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

.PHONY: binary/build
## build all binaries
binary/build: \
	cmd/agent/core/ngt/ngt \
	cmd/agent/sidecar/sidecar \
	cmd/discoverer/k8s/discoverer \
	cmd/gateway/lb/lb \
	cmd/gateway/filter/filter \
	cmd/manager/index/index

cmd/agent/core/ngt/ngt: \
	ngt/install \
	$(GO_SOURCES_INTERNAL) \
	$(PBGOS) \
	$(shell find $(ROOTDIR)/cmd/agent/core/ngt -type f -name '*.go' -not -name '*_test.go' -not -name 'doc.go') \
	$(shell find $(ROOTDIR)/pkg/agent/core/ngt $(ROOTDIR)/pkg/agent/internal -type f -name '*.go' -not -name '*_test.go' -not -name 'doc.go')
	$(eval CGO_ENABLED = 1)
	CFLAGS="$(CFLAGS)" \
	CXXFLAGS="$(CXXFLAGS)" \
	CGO_ENABLED=$(CGO_ENABLED) \
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
		-X '$(GOPKG)/internal/info.CGOEnabled=$(CGO_ENABLED)' \
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
	$(shell find $(ROOTDIR)/cmd/agent/sidecar -type f -name '*.go' -not -name '*_test.go' -not -name 'doc.go') \
	$(shell find $(ROOTDIR)/pkg/agent/sidecar $(ROOTDIR)/pkg/agent/internal -type f -name '*.go' -not -name '*_test.go' -not -name 'doc.go')
	$(eval CGO_ENABLED = 0)
	CGO_ENABLED=$(CGO_ENABLED) \
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
		-X '$(GOPKG)/internal/info.CGOEnabled=$(CGO_ENABLED)' \
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
	$(shell find $(ROOTDIR)/cmd/discoverer/k8s -type f -name '*.go' -not -name '*_test.go' -not -name 'doc.go') \
	$(shell find $(ROOTDIR)/pkg/discoverer/k8s -type f -name '*.go' -not -name '*_test.go' -not -name 'doc.go')
	$(eval CGO_ENABLED = 0)
	CGO_ENABLED=$(CGO_ENABLED) \
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
		-X '$(GOPKG)/internal/info.CGOEnabled=$(CGO_ENABLED)' \
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
	$(shell find $(ROOTDIR)/cmd/gateway/lb -type f -name '*.go' -not -name '*_test.go' -not -name 'doc.go') \
	$(shell find $(ROOTDIR)/pkg/gateway/lb -type f -name '*.go' -not -name '*_test.go' -not -name 'doc.go')
	$(eval CGO_ENABLED = 0)
	CGO_ENABLED=$(CGO_ENABLED) \
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
		-X '$(GOPKG)/internal/info.CGOEnabled=$(CGO_ENABLED)' \
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
	$(shell find $(ROOTDIR)/cmd/gateway/filter -type f -name '*.go' -not -name '*_test.go' -not -name 'doc.go') \
	$(shell find $(ROOTDIR)/pkg/gateway/filter -type f -name '*.go' -not -name '*_test.go' -not -name 'doc.go')
	$(eval CGO_ENABLED = 0)
	CGO_ENABLED=$(CGO_ENABLED) \
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
		-X '$(GOPKG)/internal/info.CGOEnabled=$(CGO_ENABLED)' \
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
	$(shell find $(ROOTDIR)/cmd/manager/index -type f -name '*.go' -not -name '*_test.go' -not -name 'doc.go') \
	$(shell find $(ROOTDIR)/pkg/manager/index -type f -name '*.go' -not -name '*_test.go' -not -name 'doc.go')
	$(eval CGO_ENABLED = 0)
	CGO_ENABLED=$(CGO_ENABLED) \
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
		-X '$(GOPKG)/internal/info.CGOEnabled=$(CGO_ENABLED)' \
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

cmd/index/job/correction/index-correction: \
	$(GO_SOURCES_INTERNAL) \
	$(PBGOS) \
	$(shell find $(ROOTDIR)/cmd/index/job/correction -type f -name '*.go' -not -name '*_test.go' -not -name 'doc.go') \
	$(shell find $(ROOTDIR)/pkg/index/job/correction -type f -name '*.go' -not -name '*_test.go' -not -name 'doc.go')
	$(eval CGO_ENABLED = 0)
	CGO_ENABLED=$(CGO_ENABLED) \
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
		-X '$(GOPKG)/internal/info.CGOEnabled=$(CGO_ENABLED)' \
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

cmd/index/job/creation/index-creation: \
	$(GO_SOURCES_INTERNAL) \
	$(PBGOS) \
	$(shell find $(ROOTDIR)/cmd/index/job/creation -type f -name '*.go' -not -name '*_test.go' -not -name 'doc.go') \
	$(shell find $(ROOTDIR)/pkg/index/job/creation -type f -name '*.go' -not -name '*_test.go' -not -name 'doc.go')
	$(eval CGO_ENABLED = 0)
	CGO_ENABLED=$(CGO_ENABLED) \
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
		-X '$(GOPKG)/internal/info.CGOEnabled=$(CGO_ENABLED)' \
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

cmd/index/job/save/index-save: \
	$(GO_SOURCES_INTERNAL) \
	$(PBGOS) \
	$(shell find $(ROOTDIR)/cmd/index/job/save -type f -name '*.go' -not -name '*_test.go' -not -name 'doc.go') \
	$(shell find $(ROOTDIR)/pkg/index/job/save -type f -name '*.go' -not -name '*_test.go' -not -name 'doc.go')
	$(eval CGO_ENABLED = 0)
	CGO_ENABLED=$(CGO_ENABLED) \
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
		-X '$(GOPKG)/internal/info.CGOEnabled=$(CGO_ENABLED)' \
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

cmd/index/job/readreplica/rotate/readreplica-rotate: \
	$(GO_SOURCES_INTERNAL) \
	$(PBGOS) \
	$(shell find $(ROOTDIR)/cmd/index/job/readreplica/rotate -type f -name '*.go' -not -name '*_test.go' -not -name 'doc.go') \
	$(shell find $(ROOTDIR)/pkg/index/job/readreplica/rotate -type f -name '*.go' -not -name '*_test.go' -not -name 'doc.go')
	$(eval CGO_ENABLED = 0)
	CGO_ENABLED=$(CGO_ENABLED) \
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
		-X '$(GOPKG)/internal/info.CGOEnabled=$(CGO_ENABLED)' \
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

