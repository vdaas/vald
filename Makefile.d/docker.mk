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
	docker/build/discoverer-k8s \
	docker/build/gateway-vald \
	docker/build/meta-redis \
	docker/build/meta-cassandra \
	docker/build/backup-manager-mysql \
	docker/build/backup-manager-cassandra \
	docker/build/manager-compressor \
	docker/build/manager-index

.PHONY: docker/name/base
docker/name/base:
	@echo "$(REPO)/$(BASE_IMAGE)"

.PHONY: docker/build/base
## build base image
docker/build/base:
	docker build -f dockers/base/Dockerfile -t $(REPO)/$(BASE_IMAGE) .

.PHONY: docker/name/agent-ngt
docker/name/agent-ngt:
	@echo "$(REPO)/$(AGENT_IMAGE)"

.PHONY: docker/build/agent-ngt
## build agent-ngt image
docker/build/agent-ngt: docker/build/base
	docker build -f dockers/agent/ngt/Dockerfile -t $(REPO)/$(AGENT_IMAGE) .

.PHONY: docker/name/discoverer-k8s
docker/name/discoverer-k8s:
	@echo "$(REPO)/$(DISCOVERER_IMAGE)"

.PHONY: docker/build/discoverer-k8s
## build discoverer-k8s image
docker/build/discoverer-k8s: docker/build/base
	docker build -f dockers/discoverer/k8s/Dockerfile -t $(REPO)/$(DISCOVERER_IMAGE) .

.PHONY: docker/name/gateway-vald
docker/name/gateway-vald:
	@echo "$(REPO)/$(GATEWAY_IMAGE)"

.PHONY: docker/build/gateway-vald
## build gateway-vald image
docker/build/gateway-vald: docker/build/base
	docker build -f dockers/gateway/vald/Dockerfile -t $(REPO)/$(GATEWAY_IMAGE) .

.PHONY: docker/name/meta-redis
docker/name/meta-redis:
	@echo "$(REPO)/$(KVS_IMAGE)"

.PHONY: docker/build/meta-redis
## build meta-redis image
docker/build/meta-redis: docker/build/base
	docker build -f dockers/meta/redis/Dockerfile -t $(REPO)/$(KVS_IMAGE) .

.PHONY: docker/name/meta-cassandra
docker/name/meta-cassandra:
	@echo "$(REPO)/$(NOSQL_IMAGE)"

.PHONY: docker/build/meta-cassandra
## build meta-cassandra image
docker/build/meta-cassandra: docker/build/base
	docker build -f dockers/meta/cassandra/Dockerfile -t $(REPO)/$(NOSQL_IMAGE) .

.PHONY: docker/name/backup-manager-mysql
docker/name/backup-manager-mysql:
	@echo "$(REPO)/$(BACKUP_MANAGER_MYSQL_IMAGE)"

.PHONY: docker/build/backup-manager-mysql
## build backup-manager-mysql image
docker/build/backup-manager-mysql: docker/build/base
	docker build -f dockers/manager/backup/mysql/Dockerfile -t $(REPO)/$(BACKUP_MANAGER_MYSQL_IMAGE) .

.PHONY: docker/name/backup-manager-cassandra
docker/name/backup-manager-cassandra:
	@echo "$(REPO)/$(BACKUP_MANAGER_CASSANDRA_IMAGE)"

.PHONY: docker/build/backup-manager-cassandra
## build backup-manager-cassandra image
docker/build/backup-manager-cassandra: docker/build/base
	docker build -f dockers/manager/backup/cassandra/Dockerfile -t $(REPO)/$(BACKUP_MANAGER_CASSANDRA_IMAGE) .

.PHONY: docker/name/manager-compressor
docker/name/manager-compressor:
	@echo "$(REPO)/$(MANAGER_COMPRESSOR_IMAGE)"

.PHONY: docker/build/manager-compressor
## build manager-compressor image
docker/build/manager-compressor: docker/build/base
	docker build -f dockers/manager/compressor/Dockerfile -t $(REPO)/$(MANAGER_COMPRESSOR_IMAGE) .

.PHONY: docker/name/manager-index
docker/name/manager-index:
	@echo "$(REPO)/$(MANAGER_INDEX_IMAGE)"

.PHONY: docker/build/manager-index
## build manager-index image
docker/build/manager-index: docker/build/base
	docker build -f dockers/manager/index/Dockerfile -t $(REPO)/$(MANAGER_INDEX_IMAGE) .

.PHONY: docker/name/ci-container
docker/name/ci-container:
	@echo "$(REPO)/$(CI_CONTAINER_IMAGE)"

.PHONY: docker/build/ci-container
## build ci-container image
docker/build/ci-container: docker/build/base
	docker build -f dockers/ci/base/Dockerfile -t $(REPO)/$(CI_CONTAINER_IMAGE) .
