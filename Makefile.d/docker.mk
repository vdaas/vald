#
# Copyright (C) 2019-2020 Vdaas.org Vald team ( kpango, rinx, kmrmt )
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
	docker/build/base \
	docker/build/agent-ngt \
	docker/build/agent-sidecar \
	docker/build/discoverer-k8s \
	docker/build/gateway-vald \
	docker/build/meta-redis \
	docker/build/meta-cassandra \
	docker/build/backup-manager-mysql \
	docker/build/backup-manager-cassandra \
	docker/build/manager-compressor \
	docker/build/manager-index \
	docker/build/helm-operator

.PHONY: docker/platforms
docker/platforms:
	@echo "linux/amd64,linux/arm64"

.PHONY: docker/name/base
docker/name/base:
	@echo "$(REPO)/$(BASE_IMAGE)"

.PHONY: docker/build/base
## build base image
docker/build/base:
	$(DOCKER) build \
	    $(DOCKER_OPTS_BASE) \
	    -f dockers/base/Dockerfile \
	    -t $(REPO)/$(BASE_IMAGE):$(TAG) . \
	    --build-arg GO_VERSION=$(GO_VERSION)

.PHONY: docker/name/agent-ngt
docker/name/agent-ngt:
	@echo "$(REPO)/$(AGENT_IMAGE)"

.PHONY: docker/build/agent-ngt
## build agent-ngt image
docker/build/agent-ngt: docker/build/base
	$(DOCKER) build \
	    $(DOCKER_OPTS) \
	    -f dockers/agent/core/ngt/Dockerfile \
	    -t $(REPO)/$(AGENT_IMAGE):$(TAG) . \
	    --build-arg BASE_TAG=$(TAG) \
	    --build-arg DISTROLESS_IMAGE=$(DISTROLESS_IMAGE) \
	    --build-arg DISTROLESS_IMAGE_TAG=$(DISTROLESS_IMAGE_TAG) \
	    --build-arg UPX_OPTIONS=$(UPX_OPTIONS)

.PHONY: docker/name/agent-sidecar
docker/name/agent-sidecar:
	@echo "$(REPO)/$(AGENT_SIDECAR_IMAGE)"

.PHONY: docker/build/agent-sidecar
## build agent-sidecar image
docker/build/agent-sidecar: docker/build/base
	$(DOCKER) build \
	    $(DOCKER_OPTS) \
	    -f dockers/agent/sidecar/Dockerfile \
	    -t $(REPO)/$(AGENT_SIDECAR_IMAGE):$(TAG) . \
	    --build-arg BASE_TAG=$(TAG) \
	    --build-arg DISTROLESS_IMAGE=$(DISTROLESS_IMAGE) \
	    --build-arg DISTROLESS_IMAGE_TAG=$(DISTROLESS_IMAGE_TAG) \
	    --build-arg UPX_OPTIONS=$(UPX_OPTIONS)

.PHONY: docker/name/discoverer-k8s
docker/name/discoverer-k8s:
	@echo "$(REPO)/$(DISCOVERER_IMAGE)"

.PHONY: docker/build/discoverer-k8s
## build discoverer-k8s image
docker/build/discoverer-k8s: docker/build/base
	$(DOCKER) build \
	    $(DOCKER_OPTS) \
	    -f dockers/discoverer/k8s/Dockerfile \
	    -t $(REPO)/$(DISCOVERER_IMAGE):$(TAG) . \
	    --build-arg BASE_TAG=$(TAG) \
	    --build-arg DISTROLESS_IMAGE=$(DISTROLESS_IMAGE) \
	    --build-arg DISTROLESS_IMAGE_TAG=$(DISTROLESS_IMAGE_TAG) \
	    --build-arg UPX_OPTIONS=$(UPX_OPTIONS)

.PHONY: docker/name/gateway-vald
docker/name/gateway-vald:
	@echo "$(REPO)/$(GATEWAY_IMAGE)"

.PHONY: docker/build/gateway-vald
## build gateway-vald image
docker/build/gateway-vald: docker/build/base
	$(DOCKER) build \
	    $(DOCKER_OPTS) \
	    -f dockers/gateway/vald/Dockerfile \
	    -t $(REPO)/$(GATEWAY_IMAGE):$(TAG) . \
	    --build-arg BASE_TAG=$(TAG) \
	    --build-arg DISTROLESS_IMAGE=$(DISTROLESS_IMAGE) \
	    --build-arg DISTROLESS_IMAGE_TAG=$(DISTROLESS_IMAGE_TAG) \
	    --build-arg UPX_OPTIONS=$(UPX_OPTIONS)

.PHONY: docker/name/gateway-lb
docker/name/gateway-lb:
	@echo "$(REPO)/$(LB_GATEWAY_IMAGE)"

.PHONY: docker/build/gateway-lb
## build gateway-lb image
docker/build/gateway-lb: docker/build/base
	docker build \
	    -f dockers/gateway/lb/Dockerfile \
	    -t $(REPO)/$(LB_GATEWAY_IMAGE) . \
	    --build-arg DISTROLESS_IMAGE=$(DISTROLESS_IMAGE) \
	    --build-arg DISTROLESS_IMAGE_TAG=$(DISTROLESS_IMAGE_TAG) \
	    --build-arg UPX_OPTIONS=$(UPX_OPTIONS)

.PHONY: docker/name/gateway-meta
docker/name/gateway-meta:
	@echo "$(REPO)/$(META_GATEWAY_IMAGE)"

.PHONY: docker/build/gateway-meta
## build gateway-meta image
docker/build/gateway-meta: docker/build/base
	docker build \
	    -f dockers/gateway/meta/Dockerfile \
	    -t $(REPO)/$(META_GATEWAY_IMAGE) . \
	    --build-arg DISTROLESS_IMAGE=$(DISTROLESS_IMAGE) \
	    --build-arg DISTROLESS_IMAGE_TAG=$(DISTROLESS_IMAGE_TAG) \
	    --build-arg UPX_OPTIONS=$(UPX_OPTIONS)

.PHONY: docker/name/gateway-backup
docker/name/gateway-backup:
	@echo "$(REPO)/$(BACKUP_GATEWAY_IMAGE)"

.PHONY: docker/build/gateway-backup
## build gateway-backup image
docker/build/gateway-backup: docker/build/base
	docker build \
	    -f dockers/gateway/backup/Dockerfile \
	    -t $(REPO)/$(BACKUP_GATEWAY_IMAGE) . \
	    --build-arg DISTROLESS_IMAGE=$(DISTROLESS_IMAGE) \
	    --build-arg DISTROLESS_IMAGE_TAG=$(DISTROLESS_IMAGE_TAG) \
	    --build-arg UPX_OPTIONS=$(UPX_OPTIONS)

.PHONY: docker/name/meta-redis
docker/name/meta-redis:
	@echo "$(REPO)/$(META_REDIS_IMAGE)"

.PHONY: docker/build/meta-redis
## build meta-redis image
docker/build/meta-redis: docker/build/base
	$(DOCKER) build \
	    $(DOCKER_OPTS) \
	    -f dockers/meta/redis/Dockerfile \
	    -t $(REPO)/$(META_REDIS_IMAGE):$(TAG) . \
	    --build-arg BASE_TAG=$(TAG) \
	    --build-arg DISTROLESS_IMAGE=$(DISTROLESS_IMAGE) \
	    --build-arg DISTROLESS_IMAGE_TAG=$(DISTROLESS_IMAGE_TAG) \
	    --build-arg UPX_OPTIONS=$(UPX_OPTIONS)

.PHONY: docker/name/meta-cassandra
docker/name/meta-cassandra:
	@echo "$(REPO)/$(META_CASSANDRA_IMAGE)"

.PHONY: docker/build/meta-cassandra
## build meta-cassandra image
docker/build/meta-cassandra: docker/build/base
	$(DOCKER) build \
	    $(DOCKER_OPTS) \
	    -f dockers/meta/cassandra/Dockerfile \
	    -t $(REPO)/$(META_CASSANDRA_IMAGE):$(TAG) . \
	    --build-arg BASE_TAG=$(TAG) \
	    --build-arg DISTROLESS_IMAGE=$(DISTROLESS_IMAGE) \
	    --build-arg DISTROLESS_IMAGE_TAG=$(DISTROLESS_IMAGE_TAG) \
	    --build-arg UPX_OPTIONS=$(UPX_OPTIONS)

.PHONY: docker/name/backup-manager-mysql
docker/name/backup-manager-mysql:
	@echo "$(REPO)/$(MANAGER_BACKUP_MYSQL_IMAGE)"

.PHONY: docker/build/backup-manager-mysql
## build backup-manager-mysql image
docker/build/backup-manager-mysql: docker/build/base
	$(DOCKER) build \
	    $(DOCKER_OPTS) \
	    -f dockers/manager/backup/mysql/Dockerfile \
	    -t $(REPO)/$(MANAGER_BACKUP_MYSQL_IMAGE):$(TAG) . \
	    --build-arg BASE_TAG=$(TAG) \
	    --build-arg DISTROLESS_IMAGE=$(DISTROLESS_IMAGE) \
	    --build-arg DISTROLESS_IMAGE_TAG=$(DISTROLESS_IMAGE_TAG) \
	    --build-arg UPX_OPTIONS=$(UPX_OPTIONS)

.PHONY: docker/name/backup-manager-cassandra
docker/name/backup-manager-cassandra:
	@echo "$(REPO)/$(MANAGER_BACKUP_CASSANDRA_IMAGE)"

.PHONY: docker/build/backup-manager-cassandra
## build backup-manager-cassandra image
docker/build/backup-manager-cassandra: docker/build/base
	$(DOCKER) build \
	    $(DOCKER_OPTS) \
	    -f dockers/manager/backup/cassandra/Dockerfile \
	    -t $(REPO)/$(MANAGER_BACKUP_CASSANDRA_IMAGE):$(TAG) . \
	    --build-arg BASE_TAG=$(TAG) \
	    --build-arg DISTROLESS_IMAGE=$(DISTROLESS_IMAGE) \
	    --build-arg DISTROLESS_IMAGE_TAG=$(DISTROLESS_IMAGE_TAG) \
	    --build-arg UPX_OPTIONS=$(UPX_OPTIONS)

.PHONY: docker/name/manager-compressor
docker/name/manager-compressor:
	@echo "$(REPO)/$(MANAGER_COMPRESSOR_IMAGE)"

.PHONY: docker/build/manager-compressor
## build manager-compressor image
docker/build/manager-compressor: docker/build/base
	$(DOCKER) build \
	    $(DOCKER_OPTS) \
	    -f dockers/manager/compressor/Dockerfile \
	    -t $(REPO)/$(MANAGER_COMPRESSOR_IMAGE):$(TAG) . \
	    --build-arg BASE_TAG=$(TAG) \
	    --build-arg DISTROLESS_IMAGE=$(DISTROLESS_IMAGE) \
	    --build-arg DISTROLESS_IMAGE_TAG=$(DISTROLESS_IMAGE_TAG) \
	    --build-arg UPX_OPTIONS=$(UPX_OPTIONS)

.PHONY: docker/name/manager-index
docker/name/manager-index:
	@echo "$(REPO)/$(MANAGER_INDEX_IMAGE)"

.PHONY: docker/build/manager-index
## build manager-index image
docker/build/manager-index: docker/build/base
	$(DOCKER) build \
	    $(DOCKER_OPTS) \
	    -f dockers/manager/index/Dockerfile \
	    -t $(REPO)/$(MANAGER_INDEX_IMAGE):$(TAG) . \
	    --build-arg BASE_TAG=$(TAG) \
	    --build-arg DISTROLESS_IMAGE=$(DISTROLESS_IMAGE) \
	    --build-arg DISTROLESS_IMAGE_TAG=$(DISTROLESS_IMAGE_TAG) \
	    --build-arg UPX_OPTIONS=$(UPX_OPTIONS)

.PHONY: docker/name/ci-container
docker/name/ci-container:
	@echo "$(REPO)/$(CI_CONTAINER_IMAGE)"

.PHONY: docker/build/ci-container
## build ci-container image
docker/build/ci-container: docker/build/base
	$(DOCKER) build \
	    $(DOCKER_OPTS) \
	    -f dockers/ci/base/Dockerfile \
	    -t $(REPO)/$(CI_CONTAINER_IMAGE):$(TAG) . \
	    --build-arg BASE_TAG=$(TAG)

.PHONY: docker/name/dev-container
docker/name/dev-container:
	@echo "$(REPO)/$(DEV_CONTAINER_IMAGE)"

.PHONY: docker/build/dev-container
## build dev-container image
docker/build/dev-container: docker/build/ci-container
	$(DOCKER) build \
	    $(DOCKER_OPTS) \
	    -f dockers/dev/Dockerfile \
	    -t $(REPO)/$(DEV_CONTAINER_IMAGE):$(TAG) . \
	    --build-arg BASE_TAG=$(TAG)

.PHONY: docker/name/operator/helm
docker/name/operator/helm:
	@echo "$(REPO)/$(HELM_OPERATOR_IMAGE)"

.PHONY: docker/build/operator/helm
## build helm-operator image
docker/build/operator/helm:
	$(DOCKER) build \
	    $(DOCKER_OPTS) \
	    -f dockers/operator/helm/Dockerfile \
	    -t $(REPO)/$(HELM_OPERATOR_IMAGE):$(TAG) . \
	    --build-arg OPERATOR_SDK_VERSION=$(OPERATOR_SDK_VERSION)

.PHONY: docker/name/loadtest
docker/name/loadtest:
	@echo "$(REPO)/$(LOADTEST_IMAGE)"

.PHONY: docker/build/loadtest
## build loadtest image
docker/build/loadtest: docker/build/base
	$(DOCKER) build \
	    $(DOCKER_OPTS) \
	    -f dockers/tools/cli/loadtest/Dockerfile \
	    -t $(REPO)/$(LOADTEST_IMAGE):$(TAG) .
