#
# Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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
	docker/build/agent-ngt \
	docker/build/agent-sidecar \
	docker/build/discoverer-k8s \
	docker/build/gateway-lb \
	docker/build/gateway-filter \
	docker/build/manager-index \
	docker/build/benchmark-job \
	docker/build/benchmark-operator \
	docker/build/helm-operator

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
		--cache-to type=registry,ref=$(GHCRORG)/$(IMAGE):$(TAG)-buildcache,mode=max \
		--cache-from type=registry,ref=$(GHCRORG)/$(IMAGE):$(TAG)-buildcache \
		--build-arg BUILDKIT_INLINE_CACHE=1 \
		--build-arg GO_VERSION=$(GO_VERSION) \
		--build-arg DISTROLESS_IMAGE=$(DISTROLESS_IMAGE) \
		--build-arg DISTROLESS_IMAGE_TAG=$(DISTROLESS_IMAGE_TAG) \
		--build-arg MAINTAINER=$(MAINTAINER) \
		$(EXTRA_ARGS) \
		--sbom=true \
		--provenance=mode=max \
		-t $(CRORG)/$(IMAGE):$(TAG) \
		-t $(GHCRORG)/$(IMAGE):$(TAG) \
		--output type=registry,oci-mediatypes=true,compression=zstd,compression-level=5,force-compression=true,push=true \
		-f $(DOCKERFILE) .
else
	@echo "starting local build for $(IMAGE):$(TAG)"
	DOCKER_BUILDKIT=1 $(DOCKER) build \
		$(DOCKER_OPTS) \
		--build-arg BUILDKIT_INLINE_CACHE=1 \
		--build-arg GO_VERSION=$(GO_VERSION) \
		--build-arg DISTROLESS_IMAGE=$(DISTROLESS_IMAGE) \
		--build-arg DISTROLESS_IMAGE_TAG=$(DISTROLESS_IMAGE_TAG) \
		--build-arg MAINTAINER=$(MAINTAINER) \
		$(EXTRA_ARGS) \
		-t $(IMAGE):$(TAG) \
		-f $(DOCKERFILE) .
endif

.PHONY: docker/name/agent-ngt
docker/name/agent-ngt:
	@echo "$(ORG)/$(AGENT_IMAGE)"

.PHONY: docker/build/agent-ngt
## build agent-ngt image
docker/build/agent-ngt:
	@make DOCKERFILE="$(ROOTDIR)/dockers/agent/core/ngt/Dockerfile" \
		IMAGE=$(AGENT_IMAGE) \
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

.PHONY: docker/name/manager-index
docker/name/manager-index:
	@echo "$(ORG)/$(MANAGER_INDEX_IMAGE)"

.PHONY: docker/build/manager-index
## build manager-index image
docker/build/manager-index:
	@make DOCKERFILE="$(ROOTDIR)/dockers/manager/index/Dockerfile" \
		IMAGE=$(MANAGER_INDEX_IMAGE) \
		docker/build/image

.PHONY: docker/name/ci-container
docker/name/ci-container:
	@echo "$(ORG)/$(CI_CONTAINER_IMAGE)"

.PHONY: docker/build/ci-container
## build ci-container image
docker/build/ci-container:
	@make DOCKERFILE="$(ROOTDIR)/dockers/ci/base/Dockerfile" \
		IMAGE=$(CI_CONTAINER_IMAGE) \
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

.PHONY: docker/name/operator/helm
docker/name/operator/helm:
	@echo "$(ORG)/$(HELM_OPERATOR_IMAGE)"

.PHONY: docker/build/operator/helm
## build helm-operator image
docker/build/operator/helm:
	@make DOCKERFILE="$(ROOTDIR)/dockers/operator/helm/Dockerfile" \
		IMAGE=$(HELM_OPERATOR_IMAGE) \
		EXTRA_ARGS="--build-arg OPERATOR_SDK_VERSION=$(OPERATOR_SDK_VERSION) --build-arg UPX_OPTIONS=$(UPX_OPTIONS)" \
		docker/build/image

.PHONY: docker/name/loadtest
docker/name/loadtest:
	@echo "$(ORG)/$(LOADTEST_IMAGE)"

.PHONY: docker/build/loadtest
## build loadtest image
docker/build/loadtest:
	@make DOCKERFILE="$(ROOTDIR)/dockers/tools/cli/loadtest/Dockerfile" \
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
	$(DOCKER) build \
	    $(DOCKER_OPTS) \
	    -f dockers/tools/benchmark/job/Dockerfile \
	    -t $(ORG)/$(BENCHMARK_JOB_IMAGE):$(TAG) . \
	    --build-arg GO_VERSION=$(GO_VERSION) \
	    --build-arg DISTROLESS_IMAGE=$(DISTROLESS_IMAGE) \
	    --build-arg DISTROLESS_IMAGE_TAG=$(DISTROLESS_IMAGE_TAG)

.PHONY: docker/name/benchmark-operator
docker/name/benchmark-operator:
	@echo "$(ORG)/$(BENCHMARK_OPERATOR_IMAGE)"

.PHONY: docker/build/benchmark-operator
## build benchmark operator
docker/build/benchmark-operator:
	$(DOCKER) build \
	    $(DOCKER_OPTS) \
	    -f dockers/tools/benchmark/operator/Dockerfile \
	    -t $(ORG)/$(BENCHMARK_OPERATOR_IMAGE):$(TAG) . \
	    --build-arg GO_VERSION=$(GO_VERSION) \
	    --build-arg DISTROLESS_IMAGE=$(DISTROLESS_IMAGE) \
	    --build-arg DISTROLESS_IMAGE_TAG=$(DISTROLESS_IMAGE_TAG)
