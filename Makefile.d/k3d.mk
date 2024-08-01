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

K3D_CLUSTER_NAME  = "vald-cluster"
K3D_COMMAND       = k3d
K3D_NODES         = 5
K3D_PORT          = 6550
K3D_HOST          = localhost
K3D_INGRESS_PORT  = 8081
K3D_HOST_PID_MODE = true
K3D_OPTIONS       = --port $(K3D_INGRESS_PORT):80@loadbalancer

.PHONY: k3d/install
## install K3D
k3d/install: $(BINDIR)/k3d

$(BINDIR)/k3d:
	mkdir -p $(BINDIR)
	curl -fsSL https://raw.githubusercontent.com/k3d-io/k3d/main/install.sh | bash
	chmod a+x $(BINDIR)/$(K3D_COMMAND)

.PHONY: k3d/start
## start k3d (kubernetes in docker) cluster
k3d/start:
	$(K3D_COMMAND) cluster create $(K3D_CLUSTER_NAME) \
	  --agents $(K3D_NODES) \
	  --image docker.io/rancher/k3s:$(K3S_VERSION) \
	  --host-pid-mode=$(K3D_HOST_PID_MODE) \
	  --api-port $(K3D_HOST):$(K3D_PORT) \
	  -v "/lib/modules:/lib/modules" \
	  $(K3D_OPTIONS)
	@make k3d/config

.PHONY: k3d/storage
## storage k3d (kubernetes in docker) cluster
k3d/storage:
	kubectl apply -f https://raw.githubusercontent.com/rancher/local-path-provisioner/master/deploy/local-path-storage.yaml
	kubectl get storageclass
	kubectl patch storageclass local-path -p '{"metadata": {"annotations":{"storageclass.kubernetes.io/is-default-class":"true"}}}'

.PHONY: k3d/config
## config k3d (kubernetes in docker) cluster
k3d/config:
	export KUBECONFIG="$(shell $(K3D_COMMAND) kubeconfig merge -o $(TEMP_DIR)/k3d_$(K3D_CLUSTER_NAME)_kubeconfig.yaml $(K3D_CLUSTER_NAME) --kubeconfig-switch-context)"

.PHONY: k3d/restart
## restart k3d (kubernetes in docker) cluster
k3d/restart: \
	k3d/delete \
	k3d/start

.PHONY: k3d/delete
## stop k3d (kubernetes in docker) cluster
k3d/delete:
	-$(K3D_COMMAND) cluster delete $(K3D_CLUSTER_NAME)
