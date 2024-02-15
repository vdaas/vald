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
	cmd/agent/core/faiss/faiss \
	cmd/agent/core/ngt/ngt \
	cmd/agent/core/qbg/qbg \
	cmd/agent/core/faiss/faiss \
	cmd/agent/sidecar/sidecar \
	cmd/discoverer/k8s/discoverer \
	cmd/gateway/filter/filter \
	cmd/gateway/lb/lb \
	cmd/gateway/mirror/mirror \
	cmd/index/job/correction/index-correction \
	cmd/index/job/creation/index-creation \
	cmd/index/job/readreplica/rotate/readreplica-rotate \
	cmd/index/job/save/index-save \
	cmd/manager/index/index \
	cmd/tools/benchmark/job/job \
	cmd/tools/benchmark/operator/operator


cmd/agent/core/ngt/ngt: \
	ngt/install
	$(eval CGO_ENABLED = 1)
<<<<<<< HEAD
	$(call go-build,agent/core/ngt,-linkmode 'external',-static -fPIC -pthread -fopenmp -std=gnu++20 -lstdc++ -lm -z relro -z now $(EXTLDFLAGS), cgo,NGT-$(NGT_VERSION),$@)
=======
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
		-X '$(GOPKG)/internal/info.AlgorithmInfo=NGT-$(NGT_VERSION)' \
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

cmd/agent/core/qbg/qbg: \
	ngt/install \
	$(GO_SOURCES_INTERNAL) \
	$(PBGOS) \
	$(shell find $(ROOTDIR)/cmd/agent/core/qbg -type f -name '*.go' -not -name '*_test.go' -not -name 'doc.go') \
	$(shell find $(ROOTDIR)/pkg/agent/core/qbg $(ROOTDIR)/pkg/agent/internal -type f -name '*.go' -not -name '*_test.go' -not -name 'doc.go')
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
		-X '$(GOPKG)/internal/info.AlgorithmInfo=QBG-$(NGT_VERSION)' \
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

cmd/agent/core/faiss/faiss: \
	faiss/install \
	$(GO_SOURCES_INTERNAL) \
	$(PBGOS) \
	$(shell find ./cmd/agent/core/faiss -type f -name '*.go' -not -name '*_test.go' -not -name 'doc.go') \
	$(shell find ./pkg/agent/core/faiss ./pkg/agent/core/ngt/service/kvs ./pkg/agent/core/ngt/service/vqueue ./pkg/agent/internal -type f -name '*.go' -not -name '*_test.go' -not -name 'doc.go')
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
		-extldflags '-fPIC -pthread -fopenmp -std=gnu++20 -lstdc++ -lm -z relro -z now $(EXTLDFLAGS)' \
		-X '$(GOPKG)/internal/info.Version=$(VERSION)' \
		-X '$(GOPKG)/internal/info.GitCommit=$(GIT_COMMIT)' \
		-X '$(GOPKG)/internal/info.BuildTime=$(DATETIME)' \
		-X '$(GOPKG)/internal/info.GoVersion=$(GO_VERSION)' \
		-X '$(GOPKG)/internal/info.GoOS=$(GOOS)' \
		-X '$(GOPKG)/internal/info.GoArch=$(GOARCH)' \
		-X '$(GOPKG)/internal/info.CGOEnabled=$${CGO_ENABLED}' \
		-X '$(GOPKG)/internal/info.FaissVersion=$(FAISS_VERSION)' \
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
>>>>>>> feature/agent/qbg

cmd/agent/core/faiss/faiss: \
	faiss/install
	$(eval CGO_ENABLED = 1)
	$(call go-build,agent/core/faiss,-linkmode 'external',-fPIC -pthread -fopenmp -std=gnu++20 -lstdc++ -lm -z relro -z now, cgo,FAISS-$(FAISS_VERSION),$@)

cmd/agent/sidecar/sidecar:
	$(eval CGO_ENABLED = 0)
	$(call go-build,agent/sidecar,,-static,,,$@)

cmd/discoverer/k8s/discoverer:
	$(eval CGO_ENABLED = 0)
	$(call go-build,discoverer/k8s,,-static,,,$@)

cmd/gateway/lb/lb:
	$(eval CGO_ENABLED = 0)
	$(call go-build,gateway/lb,,-static,,,$@)

cmd/gateway/filter/filter:
	$(eval CGO_ENABLED = 0)
	$(call go-build,gateway/filter,,-static,,,$@)

cmd/gateway/mirror/mirror:
	$(eval CGO_ENABLED = 0)
	$(call go-build,gateway/mirror,,-static,,,$@)

cmd/manager/index/index:
	$(eval CGO_ENABLED = 0)
	$(call go-build,manager/index,,-static,,,$@)

cmd/index/job/correction/index-correction:
	$(eval CGO_ENABLED = 0)
	$(call go-build,index/job/correction,,-static,,,$@)

cmd/index/job/creation/index-creation:
	$(eval CGO_ENABLED = 0)
	$(call go-build,index/job/creation,,-static,,,$@)

cmd/index/job/save/index-save:
	$(eval CGO_ENABLED = 0)
	$(call go-build,index/job/save,,-static,,,$@)

cmd/index/job/readreplica/rotate/readreplica-rotate:
	$(eval CGO_ENABLED = 0)
	$(call go-build,index/job/readreplica/rotate,,-static,,,$@)

cmd/tools/benchmark/job/job:
	$(call go-build,tools/benchmark/job,-linkmode 'external',-static -fPIC -pthread -fopenmp -std=gnu++20 -lhdf5 -lhdf5_hl -lm -ldl, cgo,$(HDF5_VERSION),$@)

cmd/tools/benchmark/operator/operator:
	$(eval CGO_ENABLED = 0)
	$(call go-build,tools/benchmark/operator,,-static,,,$@)

.PHONY: binary/build/zip
## build all binaries and zip them
binary/build/zip: \
	artifacts/vald-agent-faiss-$(GOOS)-$(GOARCH).zip \
	artifacts/vald-agent-ngt-$(GOOS)-$(GOARCH).zip \
	artifacts/vald-agent-qbg-$(GOOS)-$(GOARCH).zip \
	artifacts/vald-agent-faiss-$(GOOS)-$(GOARCH).zip \
	artifacts/vald-agent-sidecar-$(GOOS)-$(GOARCH).zip \
	artifacts/vald-discoverer-k8s-$(GOOS)-$(GOARCH).zip \
	artifacts/vald-filter-gateway-$(GOOS)-$(GOARCH).zip \
	artifacts/vald-index-correction-$(GOOS)-$(GOARCH).zip \
	artifacts/vald-index-creation-$(GOOS)-$(GOARCH).zip \
	artifacts/vald-index-save-$(GOOS)-$(GOARCH).zip \
	artifacts/vald-lb-gateway-$(GOOS)-$(GOARCH).zip \
	artifacts/vald-manager-index-$(GOOS)-$(GOARCH).zip \
	artifacts/vald-mirror-gateway-$(GOOS)-$(GOARCH).zip \
	artifacts/vald-readreplica-rotate-$(GOOS)-$(GOARCH).zip

artifacts/vald-agent-ngt-$(GOOS)-$(GOARCH).zip: cmd/agent/core/ngt/ngt
	$(call mkdir, $(dir $@))
	zip --junk-paths $@ $<

<<<<<<< HEAD
=======
artifacts/vald-agent-qbg-$(GOOS)-$(GOARCH).zip: cmd/agent/core/qbg/qbg
	$(call mkdir, $(dir $@))
	zip --junk-paths $@ $<

>>>>>>> feature/agent/qbg
artifacts/vald-agent-faiss-$(GOOS)-$(GOARCH).zip: cmd/agent/core/faiss/faiss
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

artifacts/vald-mirror-gateway-$(GOOS)-$(GOARCH).zip: cmd/gateway/mirror/mirror
	$(call mkdir, $(dir $@))
	zip --junk-paths $@ $<

artifacts/vald-index-correction-$(GOOS)-$(GOARCH).zip: cmd/index/job/correction/index-correction
	$(call mkdir, $(dir $@))
	zip --junk-paths $@ $<

artifacts/vald-index-creation-$(GOOS)-$(GOARCH).zip: cmd/index/job/creation/index-creation
	$(call mkdir, $(dir $@))
	zip --junk-paths $@ $<

artifacts/vald-index-save-$(GOOS)-$(GOARCH).zip: cmd/index/job/save/index-save
	$(call mkdir, $(dir $@))
	zip --junk-paths $@ $<

artifacts/vald-readreplica-rotate-$(GOOS)-$(GOARCH).zip: cmd/index/job/readreplica/rotate/readreplica-rotate
	$(call mkdir, $(dir $@))
	zip --junk-paths $@ $<
