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

K3D_CLUSTER_NAME = "vald-cluster"
K3D_COMMAND      = k3d
K3D_NODES        = 5

.PHONY: k3d/install
## install K3D
k3d/install: $(BINDIR)/k3d

$(BINDIR)/k3d:
	mkdir -p $(BINDIR)
	curl -sS https://raw.githubusercontent.com/k3d-io/k3d/main/install.sh | bash
	chmod a+x $(BINDIR)/$(K3D_COMMAND)

.PHONY: k3d/start
## start k3d (kubernetes in docker) cluster
k3d/start:
	$(K3D_COMMAND) cluster create $(K3D_CLUSTER_NAME) --agents $(K3D_NODES) --image docker.io/rancher/k3s:latest -v "/lib/modules:/lib/modules"
	# $(K3D_COMMAND) cluster create $(K3D_CLUSTER_NAME) --agents $(K3D_NODES) -v "/lib/modules:/lib/modules" --host-pid-mode=true
	# $(K3D_COMMAND) cluster create $(K3D_CLUSTER_NAME) -p "8081:80@loadbalancer" --agents $(K3D_NODES) --k3s-arg '--disable=traefik@all'
	export KUBECONFIG="$(shell sudo $(K3D_COMMAND) kubeconfig merge -o $(TEMP_DIR)/k3d_$(K3D_CLUSTER_NAME)_kubeconfig.yaml $(K3D_CLUSTER_NAME))"
	kubectl apply -f https://projectcontour.io/quickstart/operator.yaml
	kubectl apply -f https://projectcontour.io/quickstart/contour-custom-resource.yaml

.PHONY: k3d/restart
## restart k3d (kubernetes in docker) cluster
k3d/restart: \
	k3d/delete \
	k3d/start


.PHONY: k3d/delete
## stop k3d (kubernetes in docker) cluster
k3d/delete:
	$(K3D_COMMAND) cluster delete $(K3D_CLUSTER_NAME)
