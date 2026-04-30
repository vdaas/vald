#
# Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
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
	binary/build/go \
	binary/build/rust

.PHONY: binary/build/go
## build go binaries
binary/build/go: \
	cmd/agent/sidecar/sidecar \
	cmd/discoverer/k8s/discoverer \
	cmd/gateway/filter/filter \
	cmd/gateway/lb/lb \
	cmd/gateway/mirror/mirror \
	cmd/index/job/correction/index-correction \
	cmd/index/job/creation/index-creation \
	cmd/index/job/deletion/index-deletion \
	cmd/index/job/exportation/index-exportation \
	cmd/index/job/readreplica/rotate/readreplica-rotate \
	cmd/index/job/save/index-save \
	cmd/index/operator/index-operator \
	cmd/manager/index/index \
	cmd/tools/benchmark/job/job \
	cmd/tools/benchmark/operator/operator \
	example/client/client \
	cmd/agent/core/ngt/ngt \
	cmd/agent/core/faiss/faiss

.PHONY: binary/build/rust
## build rust binaries
binary/build/rust: \
	rust/target/debug/agent \
	rust/target/release/agent

.PHONY: e2e/build
## build all e2e binaries
e2e/build: \
	tests/v2/e2e/e2e

cmd/agent/core/ngt/ngt: \
	ngt/install
	$(eval CGO_ENABLED = 1)
	$(call go-build,agent/core/ngt,-linkmode 'external',$(LDFLAGS) $(NGT_LDFLAGS) $(EXTLDFLAGS), cgo,NGT-$(NGT_VERSION),$@)

cmd/agent/core/faiss/faiss: \
	faiss/install
	$(eval CGO_ENABLED = 1)
	$(call go-build,agent/core/faiss,-linkmode 'external',$(LDFLAGS) $(FAISS_LDFLAGS), cgo,FAISS-$(FAISS_VERSION),$@)

cmd/agent/sidecar/sidecar:
	$(eval CGO_ENABLED = 1)
	$(call go-build,agent/sidecar,-linkmode 'external',$(LDFLAGS) $(EXTLDFLAGS), cgo,,$@)

cmd/discoverer/k8s/discoverer:
	$(eval CGO_ENABLED = 1)
	$(call go-build,discoverer/k8s,-linkmode 'external',$(LDFLAGS) $(EXTLDFLAGS), cgo,,$@)

cmd/gateway/lb/lb:
	$(eval CGO_ENABLED = 1)
	$(call go-build,gateway/lb,-linkmode 'external',$(LDFLAGS) $(EXTLDFLAGS), cgo,,$@)

cmd/gateway/filter/filter:
	$(eval CGO_ENABLED = 1)
	$(call go-build,gateway/filter,-linkmode 'external',$(LDFLAGS) $(EXTLDFLAGS), cgo,,$@)

cmd/gateway/mirror/mirror:
	$(eval CGO_ENABLED = 1)
	$(call go-build,gateway/mirror,-linkmode 'external',$(LDFLAGS) $(EXTLDFLAGS), cgo,,$@)

cmd/manager/index/index:
	$(eval CGO_ENABLED = 1)
	$(call go-build,manager/index,-linkmode 'external',$(LDFLAGS) $(EXTLDFLAGS), cgo,,$@)

cmd/index/job/correction/index-correction:
	$(eval CGO_ENABLED = 1)
	$(call go-build,index/job/correction,-linkmode 'external',$(LDFLAGS) $(EXTLDFLAGS), cgo,,$@)

cmd/index/job/creation/index-creation:
	$(eval CGO_ENABLED = 1)
	$(call go-build,index/job/creation,-linkmode 'external',$(LDFLAGS) $(EXTLDFLAGS), cgo,,$@)

cmd/index/job/deletion/index-deletion:
	$(eval CGO_ENABLED = 1)
	$(call go-build,index/job/deletion,-linkmode 'external',$(LDFLAGS) $(EXTLDFLAGS), cgo,,$@)

cmd/index/job/exportation/index-exportation:
	$(eval CGO_ENABLED = 1)
	$(call go-build,index/job/exportation,-linkmode 'external',$(LDFLAGS) $(EXTLDFLAGS), cgo,,$@)

cmd/index/job/save/index-save:
	$(eval CGO_ENABLED = 1)
	$(call go-build,index/job/save,-linkmode 'external',$(LDFLAGS) $(EXTLDFLAGS), cgo,,$@)

cmd/index/job/readreplica/rotate/readreplica-rotate:
	$(eval CGO_ENABLED = 1)
	$(call go-build,index/job/readreplica/rotate,-linkmode 'external',$(LDFLAGS) $(EXTLDFLAGS), cgo,,$@)

cmd/index/operator/index-operator:
	$(eval CGO_ENABLED = 1)
	$(call go-build,index/operator,-linkmode 'external',$(LDFLAGS) $(EXTLDFLAGS), cgo,,$@)

cmd/tools/benchmark/job/job:
	$(eval CGO_ENABLED = 1)
	$(call go-build,tools/benchmark/job,-linkmode 'external',$(HDF5_LDFLAGS), cgo,$(HDF5_VERSION),$@)

cmd/tools/benchmark/operator/operator:
	$(eval CGO_ENABLED = 1)
	$(call go-build,tools/benchmark/operator,-linkmode 'external',$(LDFLAGS) $(EXTLDFLAGS), cgo,,$@)

example/client/client:
	$(eval CGO_ENABLED = 1)
	$(call go-example-build,example/client,-linkmode 'external',$(HDF5_LDFLAGS), cgo,$(HDF5_VERSION),$@)

rust/target/release/agent:
	pushd rust && \
	$(CC_ENV_VARS) \
	cargo build -p agent --release && \
	popd

rust/target/debug/agent:
	pushd rust && \
	$(CC_ENV_VARS) \
	cargo build -p agent && \
	popd

tests/v2/e2e/e2e:
	$(eval CGO_ENABLED = 1)
	$(call go-e2e-build,tests/v2/e2e/crud,$(HDF5_LDFLAGS),$@)

.PHONY: binary/build/zip
## build all binaries and zip them
binary/build/zip: \
	artifacts/vald-agent-faiss-$(GOOS)-$(GOARCH).zip \
	artifacts/vald-agent-ngt-$(GOOS)-$(GOARCH).zip \
	artifacts/vald-agent-sidecar-$(GOOS)-$(GOARCH).zip \
	artifacts/vald-benchmark-job-$(GOOS)-$(GOARCH).zip \
	artifacts/vald-benchmark-operator-$(GOOS)-$(GOARCH).zip \
	artifacts/vald-discoverer-k8s-$(GOOS)-$(GOARCH).zip \
	artifacts/vald-example-client-$(GOOS)-$(GOARCH).zip \
	artifacts/vald-filter-gateway-$(GOOS)-$(GOARCH).zip \
	artifacts/vald-index-correction-$(GOOS)-$(GOARCH).zip \
	artifacts/vald-index-creation-$(GOOS)-$(GOARCH).zip \
	artifacts/vald-index-deletion-$(GOOS)-$(GOARCH).zip \
	artifacts/vald-index-exportation-$(GOOS)-$(GOARCH).zip \
	artifacts/vald-index-operator-$(GOOS)-$(GOARCH).zip \
	artifacts/vald-index-save-$(GOOS)-$(GOARCH).zip \
	artifacts/vald-lb-gateway-$(GOOS)-$(GOARCH).zip \
	artifacts/vald-manager-index-$(GOOS)-$(GOARCH).zip \
	artifacts/vald-mirror-gateway-$(GOOS)-$(GOARCH).zip \
	artifacts/vald-readreplica-rotate-$(GOOS)-$(GOARCH).zip

artifacts/%.zip:
	$(call mkdir, $(dir $@))
	@case "$*" in \
		"vald-agent-ngt-$(GOOS)-$(GOARCH)" ) binary="cmd/agent/core/ngt/ngt" ;; \
		"vald-agent-faiss-$(GOOS)-$(GOARCH)" ) binary="cmd/agent/core/faiss/faiss" ;; \
		"vald-agent-sidecar-$(GOOS)-$(GOARCH)" ) binary="cmd/agent/sidecar/sidecar" ;; \
		"vald-discoverer-k8s-$(GOOS)-$(GOARCH)" ) binary="cmd/discoverer/k8s/discoverer" ;; \
		"vald-lb-gateway-$(GOOS)-$(GOARCH)" ) binary="cmd/gateway/lb/lb" ;; \
		"vald-filter-gateway-$(GOOS)-$(GOARCH)" ) binary="cmd/gateway/filter/filter" ;; \
		"vald-manager-index-$(GOOS)-$(GOARCH)" ) binary="cmd/manager/index/index" ;; \
		"vald-benchmark-job-$(GOOS)-$(GOARCH)" ) binary="cmd/tools/benchmark/job/job" ;; \
		"vald-benchmark-operator-$(GOOS)-$(GOARCH)" ) binary="cmd/tools/benchmark/operator/operator" ;; \
		"vald-mirror-gateway-$(GOOS)-$(GOARCH)" ) binary="cmd/gateway/mirror/mirror" ;; \
		"vald-index-correction-$(GOOS)-$(GOARCH)" ) binary="cmd/index/job/correction/index-correction" ;; \
		"vald-index-creation-$(GOOS)-$(GOARCH)" ) binary="cmd/index/job/creation/index-creation" ;; \
		"vald-index-deletion-$(GOOS)-$(GOARCH)" ) binary="cmd/index/job/deletion/index-deletion" ;; \
		"vald-index-exportation-$(GOOS)-$(GOARCH)" ) binary="cmd/index/job/exportation/index-exportation" ;; \
		"vald-index-save-$(GOOS)-$(GOARCH)" ) binary="cmd/index/job/save/index-save" ;; \
		"vald-readreplica-rotate-$(GOOS)-$(GOARCH)" ) binary="cmd/index/job/readreplica/rotate/readreplica-rotate" ;; \
		"vald-index-operator-$(GOOS)-$(GOARCH)" ) binary="cmd/index/operator/index-operator" ;; \
		"vald-example-client-$(GOOS)-$(GOARCH)" ) binary="example/client/client" ;; \
	esac; \
	zip --junk-paths $@ $$binary