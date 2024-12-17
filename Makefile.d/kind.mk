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

SNAPSHOTTER_VERSION=v8.2.0

.PHONY: kind/install
## install KinD
kind/install: $(BINDIR)/kind

$(BINDIR)/kind:
	mkdir -p $(BINDIR)
	$(eval DARCH := $(subst aarch64,arm64,$(ARCH)))
	curl -fsSL https://github.com/kubernetes-sigs/kind/releases/download/v$(KIND_VERSION)/kind-$(OS)-$(subst x86_64,amd64,$(shell echo $(DARCH) | tr '[:upper:]' '[:lower:]')) -o $(BINDIR)/kind
	chmod a+x $(BINDIR)/kind

.PHONY: kind/start
## start kind (kubernetes in docker) cluster
kind/start: \
	$(BINDIR)/docker
	kind create cluster --name $(NAME)
	@make kind/login

.PHONY: kind/login
## login command for kind (kubernetes in docker) cluster
kind/login:
	kubectl cluster-info --context kind-$(NAME)

.PHONY: kind/stop
## stop kind (kubernetes in docker) cluster
kind/stop: \
	$(BINDIR)/docker
	kind delete cluster --name $(NAME)

.PHONY: kind/restart
## restart kind (kubernetes in docker) cluster
kind/restart: \
	kind/stop \
	kind/start

.PHONY: kind/cluster/start
## start kind (kubernetes in docker) multi node cluster
kind/cluster/start:
	sudo sysctl net/netfilter/nf_conntrack_max=524288
	kind create cluster --name $(NAME)-cluster --config $(ROOTDIR)/k8s/debug/kind/config.yaml
	kubectl apply -f https://projectcontour.io/quickstart/operator.yaml
	kubectl apply -f https://projectcontour.io/quickstart/contour-custom-resource.yaml

.PHONY: kind/cluster/stop
## stop kind (kubernetes in docker) multi node cluster
kind/cluster/stop:
	kind delete cluster --name $(NAME)-cluster

.PHONY: kind/cluster/login
## login command for kind (kubernetes in docker) multi node cluster
kind/cluster/login:
	kubectl cluster-info --context kind-$(NAME)-cluster

.PHONY: kind/cluster/restart
## restart kind (kubernetes in docker) multi node cluster
kind/cluster/restart: \
	kind/cluster/stop \
	kind/cluster/start

.PHONY: kind/vs/start
## start kind (kubernetes in docker) cluster with volume snapshot
kind/vs/start: \
	$(BINDIR)/docker
	sed -e 's/apiServerAddress: "127.0.0.1"/apiServerAddress: "$(shell grep host.docker.internal /etc/hosts | cut -f1)"/' $(ROOTDIR)/k8s/debug/kind/e2e.yaml | kind create cluster --name $(NAME)-vs --config - 
	@make kind/vs/login

	kubectl apply -f https://raw.githubusercontent.com/kubernetes-csi/external-snapshotter/$(SNAPSHOTTER_VERSION)/client/config/crd/snapshot.storage.k8s.io_volumesnapshotclasses.yaml
	kubectl apply -f https://raw.githubusercontent.com/kubernetes-csi/external-snapshotter/$(SNAPSHOTTER_VERSION)/client/config/crd/snapshot.storage.k8s.io_volumesnapshotcontents.yaml
	kubectl apply -f https://raw.githubusercontent.com/kubernetes-csi/external-snapshotter/$(SNAPSHOTTER_VERSION)/client/config/crd/snapshot.storage.k8s.io_volumesnapshots.yaml
	kubectl apply -f https://raw.githubusercontent.com/kubernetes-csi/external-snapshotter/$(SNAPSHOTTER_VERSION)/deploy/kubernetes/snapshot-controller/rbac-snapshot-controller.yaml
	kubectl apply -f https://raw.githubusercontent.com/kubernetes-csi/external-snapshotter/$(SNAPSHOTTER_VERSION)/deploy/kubernetes/snapshot-controller/setup-snapshot-controller.yaml

	mkdir -p $(TEMP_DIR)/csi-driver-hostpath \
		&& curl -fsSL https://github.com/kubernetes-csi/csi-driver-host-path/archive/refs/tags/v1.15.0.tar.gz | tar zxf - -C $(TEMP_DIR)/csi-driver-hostpath --strip-components 1 \
		&& cd $(TEMP_DIR)/csi-driver-hostpath \
		&& deploy/kubernetes-latest/deploy.sh \
		&& kubectl apply -f examples/csi-storageclass.yaml \
		&& kubectl apply -f examples/csi-pvc.yaml \
		&& rm -rf $(TEMP_DIR)/csi-driver-hostpath

	@make k8s/metrics/metrics-server/deploy
	helm upgrade --install --set args={--kubelet-insecure-tls} metrics-server metrics-server/metrics-server -n kube-system
	sleep $(K8S_SLEEP_DURATION_FOR_WAIT_COMMAND)

.PHONY: kind/vs/stop
## stop kind (kubernetes in docker) cluster with volume snapshot
kind/vs/stop: \
	$(BINDIR)/docker
	kind delete cluster --name $(NAME)-vs

.PHONY: kind/vs/login
## login command for kind (kubernetes in docker)  cluster with volume snapshot
kind/vs/login:
	kubectl cluster-info --context kind-$(NAME)-vs

.PHONY: kind/vs/restart
## restart kind (kubernetes in docker) cluster with volume snapshot
kind/vs/restart: \
	kind/vs/stop \
	kind/vs/start
