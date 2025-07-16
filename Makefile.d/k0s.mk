#
# Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
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

K0S_COMMAND = k0s
K0S_OPTIONS ?=

.PHONY: k0s/install
## install K0S
k0s/install: $(BINDIR)/k0s

$(BINDIR)/k0s: update/k0s

.PHONY: k0s/start
## start k0s cluster
k0s/start:
	sudo $(K0S_COMMAND) install controller
	sudo k0s start \
		$(K0S_OPTIONS)
	sudo k0s kubectl get nodes




# .PHONY: k3d/storage
# ## storage k3d (kubernetes in docker) cluster
# k3d/storage:
# 	kubectl apply -f https://raw.githubusercontent.com/rancher/local-path-provisioner/master/deploy/local-path-storage.yaml
# 	kubectl get storageclass
# 	kubectl patch storageclass local-path -p '{"metadata": {"annotations":{"storageclass.kubernetes.io/is-default-class":"true"}}}'

# .PHONY: k3d/config
# ## config k3d (kubernetes in docker) cluster
# k3d/config:
# 	export KUBECONFIG="$(shell $(K3D_COMMAND) kubeconfig merge -o $(TEMP_DIR)/k3d_$(K3D_CLUSTER_NAME)_kubeconfig.yaml $(K3D_CLUSTER_NAME) --kubeconfig-switch-context)"

# .PHONY: k3d/restart
# ## restart k3d (kubernetes in docker) cluster
# k3d/restart: \
# 	k3d/delete \
# 	k3d/start

# .PHONY: k3d/delete
# ## stop k3d (kubernetes in docker) cluster
# k3d/delete:
# 	-$(K3D_COMMAND) cluster delete $(K3D_CLUSTER_NAME)
