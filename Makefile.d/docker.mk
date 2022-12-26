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
.PHONY: docker/build
## build all docker images
docker/build: \
	docker/build/agent-ngt \
	docker/build/agent-sidecar \
	docker/build/discoverer-k8s \
	docker/build/gateway-lb \
	docker/build/gateway-filter \
	docker/build/manager-index \
	docker/build/helm-operator

.PHONY: docker/name/org
docker/name/org:
	@echo "$(ORG)"

.PHONY: docker/name/org/alter
docker/name/org/alter:
	@echo "ghcr.io/vdaas/vald"

.PHONY: docker/platforms
docker/platforms:
	@echo "linux/amd64,linux/arm64"

.PHONY: docker/name/agent-ngt
docker/name/agent-ngt:
	@echo "$(ORG)/$(AGENT_IMAGE)"

.PHONY: docker/build/agent-ngt
## build agent-ngt image
docker/build/agent-ngt:
	$(DOCKER) build \
	    $(DOCKER_OPTS) \
	    -f dockers/agent/core/ngt/Dockerfile \
	    -t $(ORG)/$(AGENT_IMAGE):$(TAG) . \
	    --build-arg GO_VERSION=$(GO_VERSION) \
	    --build-arg DISTROLESS_IMAGE=$(DISTROLESS_IMAGE) \
	    --build-arg DISTROLESS_IMAGE_TAG=$(DISTROLESS_IMAGE_TAG) \
	    --build-arg MAINTAINER=$(MAINTAINER)

.PHONY: docker/name/agent-sidecar
docker/name/agent-sidecar:
	@echo "$(ORG)/$(AGENT_SIDECAR_IMAGE)"

.PHONY: docker/build/agent-sidecar
## build agent-sidecar image
docker/build/agent-sidecar:
	$(DOCKER) build \
	    $(DOCKER_OPTS) \
	    -f dockers/agent/sidecar/Dockerfile \
	    -t $(ORG)/$(AGENT_SIDECAR_IMAGE):$(TAG) . \
	    --build-arg GO_VERSION=$(GO_VERSION) \
	    --build-arg DISTROLESS_IMAGE=$(DISTROLESS_IMAGE) \
	    --build-arg DISTROLESS_IMAGE_TAG=$(DISTROLESS_IMAGE_TAG) \
	    --build-arg MAINTAINER=$(MAINTAINER)

.PHONY: docker/name/discoverer-k8s
docker/name/discoverer-k8s:
	@echo "$(ORG)/$(DISCOVERER_IMAGE)"

.PHONY: docker/build/discoverer-k8s
## build discoverer-k8s image
docker/build/discoverer-k8s:
	$(DOCKER) build \
	    $(DOCKER_OPTS) \
	    -f dockers/discoverer/k8s/Dockerfile \
	    -t $(ORG)/$(DISCOVERER_IMAGE):$(TAG) . \
	    --build-arg GO_VERSION=$(GO_VERSION) \
	    --build-arg DISTROLESS_IMAGE=$(DISTROLESS_IMAGE) \
	    --build-arg DISTROLESS_IMAGE_TAG=$(DISTROLESS_IMAGE_TAG) \
	    --build-arg MAINTAINER=$(MAINTAINER)

.PHONY: docker/name/gateway-lb
docker/name/gateway-lb:
	@echo "$(ORG)/$(LB_GATEWAY_IMAGE)"

.PHONY: docker/build/gateway-lb
## build gateway-lb image
docker/build/gateway-lb:
	$(DOCKER) build \
	    $(DOCKER_OPTS) \
	    -f dockers/gateway/lb/Dockerfile \
	    -t $(ORG)/$(LB_GATEWAY_IMAGE):$(TAG) . \
	    --build-arg GO_VERSION=$(GO_VERSION) \
	    --build-arg DISTROLESS_IMAGE=$(DISTROLESS_IMAGE) \
	    --build-arg DISTROLESS_IMAGE_TAG=$(DISTROLESS_IMAGE_TAG)

.PHONY: docker/name/gateway-filter
docker/name/gateway-filter:
	@echo "$(ORG)/$(FILTER_GATEWAY_IMAGE)"

.PHONY: docker/build/gateway-filter
## build gateway-filter image
docker/build/gateway-filter:
	$(DOCKER) build \
	    $(DOCKER_OPTS) \
	    -f dockers/gateway/filter/Dockerfile \
	    -t $(ORG)/$(FILTER_GATEWAY_IMAGE):$(TAG) . \
	    --build-arg GO_VERSION=$(GO_VERSION) \
	    --build-arg DISTROLESS_IMAGE=$(DISTROLESS_IMAGE) \
	    --build-arg DISTROLESS_IMAGE_TAG=$(DISTROLESS_IMAGE_TAG)

.PHONY: docker/name/gateway-mirror
docker/name/gateway-mirror:
	@echo "$(ORG)/$(MIRROR_GATEWAY_IMAGE)"

.PHONY: docker/build/gateway-mirror
## build gateway-mirror image
docker/build/gateway-mirror:
	$(DOCKER) build \
	    $(DOCKER_OPTS) \
	    -f dockers/gateway/mirror/Dockerfile \
	    -t $(ORG)/$(MIRROR_GATEWAY_IMAGE):$(TAG) . \
	    --build-arg GO_VERSION=$(GO_VERSION) \
	    --build-arg DISTROLESS_IMAGE=$(DISTROLESS_IMAGE) \
	    --build-arg DISTROLESS_IMAGE_TAG=$(DISTROLESS_IMAGE_TAG)

.PHONY: docker/name/manager-index
docker/name/manager-index:
	@echo "$(ORG)/$(MANAGER_INDEX_IMAGE)"

.PHONY: docker/build/manager-index
## build manager-index image
docker/build/manager-index:
	$(DOCKER) build \
	    $(DOCKER_OPTS) \
	    -f dockers/manager/index/Dockerfile \
	    -t $(ORG)/$(MANAGER_INDEX_IMAGE):$(TAG) . \
	    --build-arg GO_VERSION=$(GO_VERSION) \
	    --build-arg DISTROLESS_IMAGE=$(DISTROLESS_IMAGE) \
	    --build-arg DISTROLESS_IMAGE_TAG=$(DISTROLESS_IMAGE_TAG) \
	    --build-arg MAINTAINER=$(MAINTAINER)

.PHONY: docker/name/ci-container
docker/name/ci-container:
	@echo "$(ORG)/$(CI_CONTAINER_IMAGE)"

.PHONY: docker/build/ci-container
## build ci-container image
docker/build/ci-container:
	$(DOCKER) build \
	    $(DOCKER_OPTS) \
	    -f dockers/ci/base/Dockerfile \
	    -t $(ORG)/$(CI_CONTAINER_IMAGE):$(TAG) . \
	    --build-arg MAINTAINER=$(MAINTAINER) \
	    --build-arg GO_VERSION=$(GO_VERSION)

.PHONY: docker/name/dev-container
docker/name/dev-container:
	@echo "$(ORG)/$(DEV_CONTAINER_IMAGE)"

.PHONY: docker/build/dev-container
## build dev-container image
docker/build/dev-container:
	$(DOCKER) build \
	    $(DOCKER_OPTS) \
	    -f dockers/dev/Dockerfile \
	    -t $(ORG)/$(DEV_CONTAINER_IMAGE):$(TAG) . \
	    --build-arg MAINTAINER=$(MAINTAINER)

.PHONY: docker/name/operator/helm
docker/name/operator/helm:
	@echo "$(ORG)/$(HELM_OPERATOR_IMAGE)"

.PHONY: docker/build/operator/helm
## build helm-operator image
docker/build/operator/helm:
	$(DOCKER) build \
	    $(DOCKER_OPTS) \
	    -f dockers/operator/helm/Dockerfile \
	    -t $(ORG)/$(HELM_OPERATOR_IMAGE):$(TAG) . \
	    --build-arg GO_VERSION=$(GO_VERSION) \
	    --build-arg DISTROLESS_IMAGE=$(DISTROLESS_IMAGE) \
	    --build-arg DISTROLESS_IMAGE_TAG=$(DISTROLESS_IMAGE_TAG) \
	    --build-arg MAINTAINER=$(MAINTAINER) \
	    --build-arg OPERATOR_SDK_VERSION=$(OPERATOR_SDK_VERSION) \
	    --build-arg UPX_OPTIONS=$(UPX_OPTIONS)

.PHONY: docker/name/loadtest
docker/name/loadtest:
	@echo "$(ORG)/$(LOADTEST_IMAGE)"

.PHONY: docker/build/loadtest
## build loadtest image
docker/build/loadtest:
	$(DOCKER) build \
	    $(DOCKER_OPTS) \
	    -f dockers/tools/cli/loadtest/Dockerfile \
	    -t $(ORG)/$(LOADTEST_IMAGE):$(TAG) . \
	    --build-arg MAINTAINER=$(MAINTAINER) \
	    --build-arg GO_VERSION=$(GO_VERSION)
