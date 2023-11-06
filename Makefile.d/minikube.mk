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
.PHONY: minikube/install
minikube/install: $(BINDIR)/minikube

$(BINDIR)/minikube:
	mkdir -p $(BINDIR)
	curl -L https://storage.googleapis.com/minikube/releases/latest/minikube-$(shell echo $(UNAME) | tr '[:upper:]' '[:lower:]')-$(subst x86_64,amd64,$(shell echo $(ARCH) | tr '[:upper:]' '[:lower:]')) -o $(BINDIR)/minikube
	chmod a+x $(BINDIR)/minikube

# Start minikube with CSI Driver and Volume Snapshots support
# Only use this for development related to Volume Snapshots. Usually k3d is faster.
.PHONY: minikube/start
minikube/start:
	minikube start --force
	minikube addons enable volumesnapshots
	minikube addons enable csi-hostpath-driver
	minikube addons disable storage-provisioner
	minikube addons disable default-storageclass
	kubectl patch storageclass csi-hostpath-sc -p '{"metadata": {"annotations":{"storageclass.kubernetes.io/is-default-class":"true"}}}'

.PHONY: minikube/delete
minikube/delete:
	minikube delete

.PHONY: minikube/restart
minikube/restart:
	@make minikube/delete
	@make minikube/start
