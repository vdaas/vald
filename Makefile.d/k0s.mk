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

K0S_OPTIONS ?=
KUBECONFIG = ~/.kube/config

$(BINDIR)/k0s: update/k0s

.PHONY: k0s/start
## start k0s cluster
k0s/start:
	docker rm -f k0s-controller || true
	docker rm -f k0s-worker || true
	docker run -d --name k0s-controller --hostname k0s-controller --net=host \
		-v /var/lib/k0s -v /var/log/pods `# this is where k0s stores its data` \
		--tmpfs /run `# this is where k0s stores runtime data` \
		--privileged `# this is the easiest way to enable container-in-container workloads` \
		-p 6443:6443 `# publish the Kubernetes API server port` \
		docker.io/k0sproject/k0s:v1.33.2-k0s.0
	sleep 10
	mkdir -p ~/.kube
	docker exec k0s-controller k0s kubeconfig admin > $(KUBECONFIG)
	until docker exec k0s-controller k0s status | grep 'Kube-api probing successful: true'; do \
		echo "Waiting for k0s to be ready..."; \
		sleep 5; \
	done
	docker run -d --name k0s-worker --hostname k0s-worker --net=host \
		-v /var/lib/k0s -v /var/log/pods `# this is where k0s stores its data` \
		--tmpfs /run `# this is where k0s stores runtime data` \
		--privileged `# this is the easiest way to enable container-in-container workloads` \
		docker.io/k0sproject/k0s:v1.33.2-k0s.0 \
		k0s worker $$(docker exec k0s-controller k0s token create --role=worker) \
		--kubelet-root-dir=/var/lib/kubelet
	until docker exec k0s-worker k0s status | grep 'Kube-api probing successful: true'; do \
		echo "Waiting for k0s to be ready..."; \
		sleep 5; \
	done

k0s/vs/start: k0s/start
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

# 	@make k8s/metrics/metrics-server/deploy
# 	helm upgrade --install --set args={--kubelet-insecure-tls} metrics-server metrics-server/metrics-server -n kube-system
# 	sleep $(K8S_SLEEP_DURATION_FOR_WAIT_COMMAND)


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
