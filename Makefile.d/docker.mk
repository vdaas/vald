.PHONY: docker/build
## build all docker images
docker/build: \
	docker/build/base \
	docker/build/agent-ngt \
	docker/build/discoverer-k8s \
	docker/build/gateway-vald \
	docker/build/meta-redis

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

.PHONY: docker/name/backup-manager
docker/name/backup-manager:
	@echo "$(REPO)/$(BACKUP_MANAGER_IMAGE)"

.PHONY: docker/build/backup-manager
## build backup-manager image
docker/build/backup-manager: docker/build/base
	docker build -f dockers/manager/backup/Dockerfile -t $(REPO)/$(BACKUP_MANAGER_IMAGE) .
