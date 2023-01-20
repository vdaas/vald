#
# Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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
.PHONY: kind/install
## install KinD
kind/install: $(BINDIR)/kind

ifeq ($(UNAME),Darwin)
$(BINDIR)/kind:
	mkdir -p $(BINDIR)
	curl -L https://github.com/kubernetes-sigs/kind/releases/download/$(KIND_VERSION)/kind-darwin-amd64 -o $(BINDIR)/kind
	chmod a+x $(BINDIR)/kind
else
$(BINDIR)/kind:
	mkdir -p $(BINDIR)
	curl -L https://github.com/kubernetes-sigs/kind/releases/download/$(KIND_VERSION)/kind-linux-amd64 -o $(BINDIR)/kind
	chmod a+x $(BINDIR)/kind
endif

.PHONY: kind/start
## start kind (kubernetes in docker) cluster
kind/start:
	kind create cluster --name $(NAME)
	@make kind/login

.PHONY: kind/login
## login command for kind (kubernetes in docker) cluster
kind/login:
	kubectl cluster-info --context kind-$(NAME)

.PHONY: kind/stop
## stop kind (kubernetes in docker) cluster
kind/stop:
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
