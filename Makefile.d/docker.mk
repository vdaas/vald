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
.PHONY: docker/build
## build all docker images
docker/build: \
	docker/build/agent \
	docker/build/agent-faiss \
	docker/build/agent-ngt \
	docker/build/agent-sidecar \
	docker/build/benchmark-job \
	docker/build/benchmark-operator \
	docker/build/binfmt \
	docker/build/buildbase \
	docker/build/buildkit \
	docker/build/buildkit-syft-scanner \
	docker/build/ci-container \
	docker/build/dev-container \
	docker/build/discoverer-k8s \
	docker/build/gateway-filter \
	docker/build/gateway-lb \
	docker/build/gateway-mirror \
	docker/build/index-correction \
	docker/build/index-creation \
	docker/build/index-deletion \
	docker/build/index-operator \
	docker/build/index-save \
	docker/build/loadtest \
	docker/build/manager-index \
	docker/build/helm-operator \
	docker/build/readreplica-rotate

docker/xpanes/build:
	@xpanes -s -c "make -f $(ROOTDIR)/Makefile {}" \
		docker/build/agent \
		docker/build/agent-faiss \
		docker/build/agent-ngt \
		docker/build/agent-sidecar \
		docker/build/benchmark-job \
		docker/build/benchmark-operator \
		docker/build/binfmt \
		docker/build/buildbase \
		docker/build/buildkit \
		docker/build/buildkit-syft-scanner \
		docker/build/ci-container \
		docker/build/dev-container \
		docker/build/discoverer-k8s \
		docker/build/gateway-filter \
		docker/build/gateway-lb \
		docker/build/gateway-mirror \
		docker/build/index-correction \
		docker/build/index-creation \
		docker/build/index-deletion \
		docker/build/index-operator \
		docker/build/index-save \
		docker/build/loadtest \
		docker/build/manager-index \
		docker/build/operator/helm \
		docker/build/readreplica-rotate

.PHONY: docker/name/org
docker/name/org:
	@echo "$(ORG)"

.PHONY: docker/name/org/alter
docker/name/org/alter:
	@echo "$(GHCRORG)"

.PHONY: docker/platforms
docker/platforms:
	@echo "linux/amd64,linux/arm64"

.PHONY: docker/build/image
## Generalized docker build function
docker/build/image:
ifeq ($(REMOTE),true)
	@echo "starting remote build for $(IMAGE):$(TAG)"
	DOCKER_BUILDKIT=1 $(DOCKER) buildx build \
		$(DOCKER_OPTS) \
		--cache-to type=gha,scope=$(TAG)-buildcache,mode=max \
		--cache-to type=registry,ref=$(GHCRORG)/$(IMAGE):$(TAG)-buildcache,mode=max \
		--cache-from type=gha,scope=$(TAG)-buildcache \
		--cache-from type=registry,ref=$(GHCRORG)/$(IMAGE):$(TAG)-buildcache \
		--build-arg BUILDKIT_INLINE_CACHE=$(BUILDKIT_INLINE_CACHE) \
		--build-arg GO_VERSION=$(GO_VERSION) \
		--build-arg RUST_VERSION=$(RUST_VERSION) \
		--build-arg MAINTAINER=$(MAINTAINER) \
		--attest type=sbom,generator=$(DEFAULT_BUILDKIT_SYFT_SCANNER_IMAGE) \
		--provenance=mode=max \
		-t $(CRORG)/$(IMAGE):$(TAG) \
		-t $(GHCRORG)/$(IMAGE):$(TAG) \
		$(EXTRA_ARGS) \
		--output type=registry,oci-mediatypes=true,compression=zstd,compression-level=5,force-compression=true,push=true \
		-f $(DOCKERFILE) $(ROOTDIR)
else
	@echo "starting local build for $(IMAGE):$(TAG)"
	DOCKER_BUILDKIT=1 $(DOCKER) build \
		$(DOCKER_OPTS) \
		--build-arg BUILDKIT_INLINE_CACHE=$(BUILDKIT_INLINE_CACHE) \
		--build-arg GO_VERSION=$(GO_VERSION) \
		--build-arg RUST_VERSION=$(RUST_VERSION) \
		--build-arg MAINTAINER=$(MAINTAINER) \
		$(EXTRA_ARGS) \
		-t $(CRORG)/$(IMAGE):$(TAG) \
		-t $(GHCRORG)/$(IMAGE):$(TAG) \
		-f $(DOCKERFILE) $(ROOTDIR)
endif

.PHONY: docker/name/agent-ngt
docker/name/agent-ngt:
	@echo "$(ORG)/$(AGENT_NGT_IMAGE)"

.PHONY: docker/build/agent-ngt
## build agent-ngt image
docker/build/agent-ngt:
	@make DOCKERFILE="$(ROOTDIR)/dockers/agent/core/ngt/Dockerfile" \
		IMAGE=$(AGENT_NGT_IMAGE) \
		docker/build/image

.PHONY: docker/name/agent-faiss
docker/name/agent-faiss:
	@echo "$(ORG)/$(AGENT_FAISS_IMAGE)"

.PHONY: docker/build/agent-faiss
## build agent-faiss image
docker/build/agent-faiss:
	@make DOCKERFILE="$(ROOTDIR)/dockers/agent/core/faiss/Dockerfile" \
		IMAGE=$(AGENT_FAISS_IMAGE) \
		docker/build/image

.PHONY: docker/name/agent-sidecar
docker/name/agent-sidecar:
	@echo "$(ORG)/$(AGENT_SIDECAR_IMAGE)"

.PHONY: docker/build/agent-sidecar
## build agent-sidecar image
docker/build/agent-sidecar:
	@make DOCKERFILE="$(ROOTDIR)/dockers/agent/sidecar/Dockerfile" \
		IMAGE=$(AGENT_SIDECAR_IMAGE) \
		docker/build/image

.PHONY: docker/name/agent
docker/name/agent:
	@echo "$(ORG)/$(AGENT_IMAGE)"

.PHONY: docker/build/agent
docker/build/agent:
	@make DOCKERFILE="$(ROOTDIR)/dockers/agent/core/agent/Dockerfile" \
		IMAGE=$(AGENT_IMAGE) \
		docker/build/image

.PHONY: docker/name/discoverer-k8s
docker/name/discoverer-k8s:
	@echo "$(ORG)/$(DISCOVERER_IMAGE)"

.PHONY: docker/build/discoverer-k8s
## build discoverer-k8s image
docker/build/discoverer-k8s:
	@make DOCKERFILE="$(ROOTDIR)/dockers/discoverer/k8s/Dockerfile" \
		IMAGE=$(DISCOVERER_IMAGE) \
		docker/build/image

.PHONY: docker/name/gateway-lb
docker/name/gateway-lb:
	@echo "$(ORG)/$(LB_GATEWAY_IMAGE)"

.PHONY: docker/build/gateway-lb
## build gateway-lb image
docker/build/gateway-lb:
	@make DOCKERFILE="$(ROOTDIR)/dockers/gateway/lb/Dockerfile" \
		IMAGE=$(LB_GATEWAY_IMAGE) \
		docker/build/image

.PHONY: docker/name/gateway-filter
docker/name/gateway-filter:
	@echo "$(ORG)/$(FILTER_GATEWAY_IMAGE)"

.PHONY: docker/build/gateway-filter
## build gateway-filter image
docker/build/gateway-filter:
	@make DOCKERFILE="$(ROOTDIR)/dockers/gateway/filter/Dockerfile" \
		IMAGE=$(FILTER_GATEWAY_IMAGE) \
		docker/build/image

.PHONY: docker/name/gateway-mirror
docker/name/gateway-mirror:
	@echo "$(ORG)/$(MIRROR_GATEWAY_IMAGE)"

.PHONY: docker/build/gateway-mirror
## build gateway-mirror image
docker/build/gateway-mirror:
	@make DOCKERFILE="$(ROOTDIR)/dockers/gateway/mirror/Dockerfile" \
		IMAGE=$(MIRROR_GATEWAY_IMAGE) \
		docker/build/image

.PHONY: docker/name/manager-index
docker/name/manager-index:
	@echo "$(ORG)/$(MANAGER_INDEX_IMAGE)"

.PHONY: docker/build/manager-index
## build manager-index image
docker/build/manager-index:
	@make DOCKERFILE="$(ROOTDIR)/dockers/manager/index/Dockerfile" \
		IMAGE=$(MANAGER_INDEX_IMAGE) \
		docker/build/image

.PHONY: docker/name/buildbase
docker/name/buildbase:
	@echo "$(ORG)/$(BUILDBASE_IMAGE)"

.PHONY: docker/build/buildbase
## build buildbase image
docker/build/buildbase:
	@make DOCKERFILE="$(ROOTDIR)/dockers/buildbase/Dockerfile" \
		IMAGE=$(BUILDBASE_IMAGE) \
		docker/build/image

.PHONY: docker/name/buildkit
docker/name/buildkit:
	@echo "$(ORG)/$(BUILDKIT_IMAGE)"

.PHONY: docker/build/buildkit
## build buildkit image
docker/build/buildkit:
	@make DOCKERFILE="$(ROOTDIR)/dockers/buildkit/Dockerfile" \
		IMAGE=$(BUILDKIT_IMAGE) \
		docker/build/image

.PHONY: docker/name/binfmt
docker/name/binfmt:
	@echo "$(ORG)/$(BINFMT_IMAGE)"

.PHONY: docker/build/binfmt
## build binfmt image
docker/build/binfmt:
	@make DOCKERFILE="$(ROOTDIR)/dockers/binfmt/Dockerfile" \
		IMAGE=$(BINFMT_IMAGE) \
		docker/build/image

PHONY: docker/name/buildkit-syft-scanner
docker/name/buildkit-syft-scanner:
	@echo "$(ORG)/$(BUILDKIT_SYFT_SCANNER_IMAGE)"

.PHONY: docker/build/buildkit-syft-scanner
## build buildkit-syft-scanner image
docker/build/buildkit-syft-scanner:
	@make DOCKERFILE="$(ROOTDIR)/dockers/buildkit/syft/scanner/Dockerfile" \
		IMAGE=$(BUILDKIT_SYFT_SCANNER_IMAGE) \
		DEFAULT_BUILDKIT_SYFT_SCANNER_IMAGE="docker/buildkit-syft-scanner:edge" \
		docker/build/image

.PHONY: docker/name/ci-container
docker/name/ci-container:
	@echo "$(ORG)/$(CI_CONTAINER_IMAGE)"

.PHONY: docker/build/ci-container
## build ci-container image
docker/build/ci-container:
	@make DOCKERFILE="$(ROOTDIR)/dockers/ci/base/Dockerfile" \
		IMAGE=$(CI_CONTAINER_IMAGE) \
		EXTRA_ARGS="--add-host=registry.npmjs.org:104.16.20.35 $(EXTRA_ARGS)" \
		docker/build/image

.PHONY: docker/name/dev-container
docker/name/dev-container:
	@echo "$(ORG)/$(DEV_CONTAINER_IMAGE)"

.PHONY: docker/build/dev-container
## build dev-container image
docker/build/dev-container:
	@make DOCKERFILE="$(ROOTDIR)/dockers/dev/Dockerfile" \
		IMAGE=$(DEV_CONTAINER_IMAGE) \
		docker/build/image

.PHONY: docker/name/helm-operator
docker/name/helm-operator:
	@echo "$(ORG)/$(HELM_OPERATOR_IMAGE)"

.PHONY: docker/build/helm-operator
## build helm-operator image
docker/build/helm-operator:
	@make DOCKERFILE="$(ROOTDIR)/dockers/operator/helm/Dockerfile" \
		IMAGE=$(HELM_OPERATOR_IMAGE) \
		EXTRA_ARGS="--build-arg OPERATOR_SDK_VERSION=$(OPERATOR_SDK_VERSION) --build-arg UPX_OPTIONS=$(UPX_OPTIONS) $(EXTRA_ARGS)" \
		docker/build/image

.PHONY: docker/name/loadtest
docker/name/loadtest:
	@echo "$(ORG)/$(LOADTEST_IMAGE)"

.PHONY: docker/build/loadtest
## build loadtest image
docker/build/loadtest:
	@make DOCKERFILE="$(ROOTDIR)/dockers/tools/cli/loadtest/Dockerfile" \
		DOCKER_OPTS="--build-arg ZLIB_VERSION=$(ZLIB_VERSION) --build-arg HDF5_VERSION=$(HDF5_VERSION)" \
		IMAGE=$(LOADTEST_IMAGE) \
		docker/build/image

.PHONY: docker/name/index-correction
docker/name/index-correction:
	@echo "$(ORG)/$(INDEX_CORRECTION_IMAGE)"

.PHONY: docker/build/index-correction
## build index-correction image
docker/build/index-correction:
	@make DOCKERFILE="$(ROOTDIR)/dockers/index/job/correction/Dockerfile" \
		IMAGE=$(INDEX_CORRECTION_IMAGE) \
		docker/build/image

.PHONY: docker/name/index-creation
docker/name/index-creation:
	@echo "$(ORG)/$(INDEX_CREATION_IMAGE)"

.PHONY: docker/build/index-creation
## build index-creation image
docker/build/index-creation:
	@make DOCKERFILE="$(ROOTDIR)/dockers/index/job/creation/Dockerfile" \
		IMAGE=$(INDEX_CREATION_IMAGE) \
		docker/build/image

.PHONY: docker/name/index-save
docker/name/index-save:
	@echo "$(ORG)/$(INDEX_SAVE_IMAGE)"

.PHONY: docker/build/index-save
## build index-save image
docker/build/index-save:
	@make DOCKERFILE="$(ROOTDIR)/dockers/index/job/save/Dockerfile" \
		IMAGE=$(INDEX_SAVE_IMAGE) \
		docker/build/image

.PHONY: docker/name/index-deletion
docker/name/index-deletion:
	@echo "$(ORG)/$(INDEX_DELETION_IMAGE)"

.PHONY: docker/build/index-deletion
## build index-deletion image
docker/build/index-deletion:
	@make DOCKERFILE="$(ROOTDIR)/dockers/index/job/deletion/Dockerfile" \
		IMAGE=$(INDEX_DELETION_IMAGE) \
		docker/build/image

.PHONY: docker/name/index-operator
docker/name/index-operator:
	@echo "$(ORG)/$(INDEX_OPERATOR_IMAGE)"

.PHONY: docker/build/index-operator
## build index-operator image
docker/build/index-operator:
	@make DOCKERFILE="$(ROOTDIR)/dockers/index/operator/Dockerfile" \
		IMAGE=$(INDEX_OPERATOR_IMAGE) \
		docker/build/image

.PHONY: docker/name/readreplica-rotate
docker/name/readreplica-rotate:
	@echo "$(ORG)/$(READREPLICA_ROTATE_IMAGE)"

.PHONY: docker/build/readreplica-rotate
## build readreplica-rotate image
docker/build/readreplica-rotate:
	@make DOCKERFILE="$(ROOTDIR)/dockers/index/job/readreplica/rotate/Dockerfile" \
		IMAGE=$(READREPLICA_ROTATE_IMAGE) \
		docker/build/image

.PHONY: docker/name/benchmark-job
docker/name/benchmark-job:
	@echo "$(ORG)/$(BENCHMARK_JOB_IMAGE)"

.PHONY: docker/build/benchmark-job
## build benchmark job
docker/build/benchmark-job:
	@make DOCKERFILE="$(ROOTDIR)/dockers/tools/benchmark/job/Dockerfile" \
		IMAGE=$(BENCHMARK_JOB_IMAGE) \
		DOCKER_OPTS="--build-arg ZLIB_VERSION=$(ZLIB_VERSION) --build-arg HDF5_VERSION=$(HDF5_VERSION)" \
		docker/build/image

.PHONY: docker/name/benchmark-operator
docker/name/benchmark-operator:
	@echo "$(ORG)/$(BENCHMARK_OPERATOR_IMAGE)"

.PHONY: docker/build/benchmark-operator
## build benchmark operator
docker/build/benchmark-operator:
	@make DOCKERFILE="$(ROOTDIR)/dockers/tools/benchmark/operator/Dockerfile" \
		IMAGE=$(BENCHMARK_OPERATOR_IMAGE) \
		docker/build/image
