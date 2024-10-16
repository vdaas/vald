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
	cmd/agent/sidecar/sidecar \
	cmd/discoverer/k8s/discoverer \
	cmd/gateway/filter/filter \
	cmd/gateway/lb/lb \
	cmd/gateway/mirror/mirror \
	cmd/index/job/correction/index-correction \
	cmd/index/job/creation/index-creation \
	cmd/index/job/deletion/index-deletion \
	cmd/index/job/readreplica/rotate/readreplica-rotate \
	cmd/index/job/save/index-save \
	cmd/index/operator/index-operator \
	cmd/manager/index/index \
	cmd/tools/benchmark/job/job \
	cmd/tools/benchmark/operator/operator \
	cmd/tools/cli/loadtest/loadtest \
	cmd/agent/core/ngt/ngt \
	cmd/agent/core/faiss/faiss \
	rust/target/debug/agent \
	rust/target/release/agent \


cmd/agent/core/ngt/ngt: \
	ngt/install
	$(eval CGO_ENABLED = 1)
	$(call go-build,agent/core/ngt,-linkmode 'external',$(LDFLAGS) $(NGT_LDFLAGS) $(EXTLDFLAGS), cgo,NGT-$(NGT_VERSION),$@)

cmd/agent/core/faiss/faiss: \
	faiss/install
	$(eval CGO_ENABLED = 1)
	$(call go-build,agent/core/faiss,-linkmode 'external',$(LDFLAGS) $(FAISS_LDFLAGS), cgo,FAISS-$(FAISS_VERSION),$@)

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

cmd/index/job/deletion/index-deletion:
	$(eval CGO_ENABLED = 0)
	$(call go-build,index/job/deletion,,-static,,,$@)

cmd/index/job/save/index-save:
	$(eval CGO_ENABLED = 0)
	$(call go-build,index/job/save,,-static,,,$@)

cmd/index/job/readreplica/rotate/readreplica-rotate:
	$(eval CGO_ENABLED = 0)
	$(call go-build,index/job/readreplica/rotate,,-static,,,$@)

cmd/index/operator/index-operator:
	$(eval CGO_ENABLED = 0)
	$(call go-build,index/operator,,-static,,,$@)

cmd/tools/benchmark/job/job:
	$(eval CGO_ENABLED = 1)
	$(call go-build,tools/benchmark/job,-linkmode 'external',$(LDFLAGS) $(HDF5_LDFLAGS), cgo,$(HDF5_VERSION),$@)

cmd/tools/benchmark/operator/operator:
	$(eval CGO_ENABLED = 0)
	$(call go-build,tools/benchmark/operator,,-static,,,$@)

cmd/tools/cli/loadtest/loadtest:
	$(eval CGO_ENABLED = 1)
	$(call go-build,tools/cli/loadtest,-linkmode 'external',$(LDFLAGS) $(HDF5_LDFLAGS), cgo,$(HDF5_VERSION),$@)

rust/target/release/agent:
	pushd rust && cargo build -p agent --release && popd

rust/target/debug/agent:
	pushd rust && cargo build -p agent && popd

.PHONY: binary/build/zip
## build all binaries and zip them
binary/build/zip: \
	artifacts/vald-agent-faiss-$(GOOS)-$(GOARCH).zip \
	artifacts/vald-agent-ngt-$(GOOS)-$(GOARCH).zip \
	artifacts/vald-agent-sidecar-$(GOOS)-$(GOARCH).zip \
	artifacts/vald-benchmark-job-$(GOOS)-$(GOARCH).zip \
	artifacts/vald-benchmark-operator-$(GOOS)-$(GOARCH).zip \
	artifacts/vald-cli-loadtest-$(GOOS)-$(GOARCH).zip \
	artifacts/vald-discoverer-k8s-$(GOOS)-$(GOARCH).zip \
	artifacts/vald-filter-gateway-$(GOOS)-$(GOARCH).zip \
	artifacts/vald-index-correction-$(GOOS)-$(GOARCH).zip \
	artifacts/vald-index-creation-$(GOOS)-$(GOARCH).zip \
	artifacts/vald-index-deletion-$(GOOS)-$(GOARCH).zip \
	artifacts/vald-index-operator-$(GOOS)-$(GOARCH).zip \
	artifacts/vald-index-save-$(GOOS)-$(GOARCH).zip \
	artifacts/vald-lb-gateway-$(GOOS)-$(GOARCH).zip \
	artifacts/vald-manager-index-$(GOOS)-$(GOARCH).zip \
	artifacts/vald-mirror-gateway-$(GOOS)-$(GOARCH).zip \
	artifacts/vald-readreplica-rotate-$(GOOS)-$(GOARCH).zip

artifacts/vald-agent-ngt-$(GOOS)-$(GOARCH).zip: cmd/agent/core/ngt/ngt
	$(call mkdir, $(dir $@))
	zip --junk-paths $@ $<

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

artifacts/vald-cli-loadtest-$(GOOS)-$(GOARCH).zip: cmd/tools/cli/loadtest/loadtest
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

artifacts/vald-index-deletion-$(GOOS)-$(GOARCH).zip: cmd/index/job/deletion/index-deletion
	$(call mkdir, $(dir $@))
	zip --junk-paths $@ $<

artifacts/vald-index-save-$(GOOS)-$(GOARCH).zip: cmd/index/job/save/index-save
	$(call mkdir, $(dir $@))
	zip --junk-paths $@ $<

artifacts/vald-readreplica-rotate-$(GOOS)-$(GOARCH).zip: cmd/index/job/readreplica/rotate/readreplica-rotate
	$(call mkdir, $(dir $@))
	zip --junk-paths $@ $<

artifacts/vald-index-operator-$(GOOS)-$(GOARCH).zip: cmd/index/operator/index-operator
	$(call mkdir, $(dir $@))
	zip --junk-paths $@ $<
